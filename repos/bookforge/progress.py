"""
Real-time Progress Tracking System for BookForge

Provides progress tracking for:
- Outline generation
- Chapter generation (per chapter and overall)
- TTS conversion
- Export jobs

Features:
- Server-Sent Events (SSE) for real-time updates
- In-memory storage with dict-based approach
- ETA calculation based on historical data
- Database logging for progress history
"""

import asyncio
import json
import time
import uuid
from datetime import datetime, timedelta
from enum import Enum
from typing import Optional, Dict, List, Any, AsyncGenerator
from dataclasses import dataclass, field, asdict
from collections import defaultdict
import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse
from pydantic import BaseModel

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


class JobType(str, Enum):
    """Types of jobs that can be tracked."""
    OUTLINE_GENERATION = "outline_generation"
    CHAPTER_GENERATION = "chapter_generation"
    TTS_CONVERSION = "tts_conversion"
    EXPORT = "export"


class JobStatus(str, Enum):
    """Status of a tracked job."""
    PENDING = "pending"
    IN_PROGRESS = "in_progress"
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"


@dataclass
class ProgressData:
    """Data structure for tracking job progress."""
    job_id: str
    job_type: JobType
    status: JobStatus
    total_steps: int
    current_step: int
    message: str
    started_at: float
    updated_at: float
    completed_at: Optional[float] = None
    error: Optional[str] = None
    eta_seconds: Optional[float] = None
    parent_job_id: Optional[str] = None  # For sub-jobs like individual chapters
    metadata: Dict[str, Any] = field(default_factory=dict)

    @property
    def progress_percent(self) -> float:
        """Calculate progress percentage."""
        if self.total_steps == 0:
            return 0.0
        return min(100.0, (self.current_step / self.total_steps) * 100)

    @property
    def elapsed_seconds(self) -> float:
        """Calculate elapsed time in seconds."""
        end_time = self.completed_at or time.time()
        return end_time - self.started_at

    def to_dict(self) -> Dict[str, Any]:
        """Convert to dictionary for JSON serialization."""
        return {
            "job_id": self.job_id,
            "job_type": self.job_type.value,
            "status": self.status.value,
            "total_steps": self.total_steps,
            "current_step": self.current_step,
            "progress_percent": round(self.progress_percent, 2),
            "message": self.message,
            "started_at": self.started_at,
            "updated_at": self.updated_at,
            "completed_at": self.completed_at,
            "elapsed_seconds": round(self.elapsed_seconds, 2),
            "eta_seconds": round(self.eta_seconds, 2) if self.eta_seconds else None,
            "error": self.error,
            "parent_job_id": self.parent_job_id,
            "metadata": self.metadata,
        }


class HistoricalStats:
    """Tracks historical job statistics for ETA calculation."""

    def __init__(self, max_samples: int = 100):
        self.max_samples = max_samples
        # Store average time per step for each job type
        self._step_times: Dict[JobType, List[float]] = defaultdict(list)
        # Store total job durations
        self._job_durations: Dict[JobType, List[float]] = defaultdict(list)

    def record_step_time(self, job_type: JobType, step_duration: float) -> None:
        """Record the duration of a single step."""
        times = self._step_times[job_type]
        times.append(step_duration)
        if len(times) > self.max_samples:
            times.pop(0)

    def record_job_completion(self, job_type: JobType, total_duration: float, total_steps: int) -> None:
        """Record a completed job's statistics."""
        if total_steps > 0:
            avg_step_time = total_duration / total_steps
            self.record_step_time(job_type, avg_step_time)

        durations = self._job_durations[job_type]
        durations.append(total_duration)
        if len(durations) > self.max_samples:
            durations.pop(0)

    def get_average_step_time(self, job_type: JobType) -> Optional[float]:
        """Get average time per step for a job type."""
        times = self._step_times[job_type]
        if not times:
            return None
        return sum(times) / len(times)

    def estimate_remaining_time(self, job_type: JobType, remaining_steps: int,
                                 current_avg_step_time: Optional[float] = None) -> Optional[float]:
        """Estimate remaining time based on historical data and current pace."""
        historical_avg = self.get_average_step_time(job_type)

        if current_avg_step_time and historical_avg:
            # Weighted average: 70% current pace, 30% historical
            avg_step_time = (current_avg_step_time * 0.7) + (historical_avg * 0.3)
        elif current_avg_step_time:
            avg_step_time = current_avg_step_time
        elif historical_avg:
            avg_step_time = historical_avg
        else:
            return None

        return avg_step_time * remaining_steps


class ProgressStore:
    """In-memory storage for progress data (Redis-like interface using dict)."""

    def __init__(self):
        self._jobs: Dict[str, ProgressData] = {}
        self._subscribers: Dict[str, List[asyncio.Queue]] = defaultdict(list)
        self._global_subscribers: List[asyncio.Queue] = []
        self._historical_stats = HistoricalStats()
        self._lock = asyncio.Lock()
        self._last_step_times: Dict[str, float] = {}  # Track time of last step update

    async def set(self, job_id: str, data: ProgressData) -> None:
        """Store progress data."""
        async with self._lock:
            self._jobs[job_id] = data

    async def get(self, job_id: str) -> Optional[ProgressData]:
        """Retrieve progress data."""
        return self._jobs.get(job_id)

    async def delete(self, job_id: str) -> bool:
        """Delete progress data."""
        async with self._lock:
            if job_id in self._jobs:
                del self._jobs[job_id]
                if job_id in self._last_step_times:
                    del self._last_step_times[job_id]
                return True
            return False

    async def get_all_active(self) -> List[ProgressData]:
        """Get all active (non-completed, non-failed) jobs."""
        return [
            job for job in self._jobs.values()
            if job.status in (JobStatus.PENDING, JobStatus.IN_PROGRESS)
        ]

    async def get_by_type(self, job_type: JobType) -> List[ProgressData]:
        """Get all jobs of a specific type."""
        return [
            job for job in self._jobs.values()
            if job.job_type == job_type
        ]

    async def get_children(self, parent_job_id: str) -> List[ProgressData]:
        """Get all child jobs of a parent job."""
        return [
            job for job in self._jobs.values()
            if job.parent_job_id == parent_job_id
        ]

    async def subscribe(self, job_id: Optional[str] = None) -> asyncio.Queue:
        """Subscribe to progress updates for a specific job or all jobs."""
        queue: asyncio.Queue = asyncio.Queue()
        async with self._lock:
            if job_id:
                self._subscribers[job_id].append(queue)
            else:
                self._global_subscribers.append(queue)
        return queue

    async def unsubscribe(self, queue: asyncio.Queue, job_id: Optional[str] = None) -> None:
        """Unsubscribe from progress updates."""
        async with self._lock:
            if job_id:
                if queue in self._subscribers[job_id]:
                    self._subscribers[job_id].remove(queue)
            else:
                if queue in self._global_subscribers:
                    self._global_subscribers.remove(queue)

    async def notify(self, job_id: str, data: ProgressData) -> None:
        """Notify all subscribers of a progress update."""
        event_data = data.to_dict()

        # Notify job-specific subscribers
        for queue in self._subscribers.get(job_id, []):
            try:
                await queue.put(event_data)
            except Exception as e:
                logger.error(f"Failed to notify subscriber: {e}")

        # Notify global subscribers
        for queue in self._global_subscribers:
            try:
                await queue.put(event_data)
            except Exception as e:
                logger.error(f"Failed to notify global subscriber: {e}")

    def calculate_eta(self, data: ProgressData) -> Optional[float]:
        """Calculate ETA based on current progress and historical data."""
        if data.current_step == 0:
            # No progress yet, use historical data only
            return self._historical_stats.estimate_remaining_time(
                data.job_type,
                data.total_steps
            )

        remaining_steps = data.total_steps - data.current_step
        if remaining_steps <= 0:
            return 0.0

        # Calculate current average step time
        elapsed = data.elapsed_seconds
        current_avg_step_time = elapsed / data.current_step if data.current_step > 0 else None

        return self._historical_stats.estimate_remaining_time(
            data.job_type,
            remaining_steps,
            current_avg_step_time
        )

    def record_step_completion(self, job_id: str, job_type: JobType) -> None:
        """Record step completion time for historical tracking."""
        current_time = time.time()
        if job_id in self._last_step_times:
            step_duration = current_time - self._last_step_times[job_id]
            self._historical_stats.record_step_time(job_type, step_duration)
        self._last_step_times[job_id] = current_time

    def record_job_completion(self, data: ProgressData) -> None:
        """Record job completion for historical tracking."""
        self._historical_stats.record_job_completion(
            data.job_type,
            data.elapsed_seconds,
            data.total_steps
        )


class DatabaseLogger:
    """Logs progress events to database for history tracking."""

    def __init__(self, db_session_factory=None):
        self._db_session_factory = db_session_factory
        self._log_buffer: List[Dict[str, Any]] = []
        self._buffer_size = 100

    async def log_event(self, event_type: str, data: ProgressData) -> None:
        """Log a progress event to the database."""
        log_entry = {
            "event_type": event_type,
            "job_id": data.job_id,
            "job_type": data.job_type.value,
            "status": data.status.value,
            "current_step": data.current_step,
            "total_steps": data.total_steps,
            "progress_percent": data.progress_percent,
            "message": data.message,
            "error": data.error,
            "timestamp": datetime.utcnow().isoformat(),
            "elapsed_seconds": data.elapsed_seconds,
        }

        self._log_buffer.append(log_entry)

        # Flush buffer if full
        if len(self._log_buffer) >= self._buffer_size:
            await self.flush()

        logger.info(f"Progress event: {event_type} - Job {data.job_id}: {data.message}")

    async def flush(self) -> None:
        """Flush buffered logs to database."""
        if not self._log_buffer:
            return

        if self._db_session_factory:
            try:
                # In a real implementation, this would write to the database
                # async with self._db_session_factory() as session:
                #     for entry in self._log_buffer:
                #         session.add(ProgressLog(**entry))
                #     await session.commit()
                pass
            except Exception as e:
                logger.error(f"Failed to flush logs to database: {e}")

        # For now, just clear the buffer (logs are already written to logger)
        self._log_buffer.clear()

    async def get_history(self, job_id: str, limit: int = 100) -> List[Dict[str, Any]]:
        """Get progress history for a job from the database."""
        # In a real implementation, this would query the database
        return [
            entry for entry in self._log_buffer
            if entry.get("job_id") == job_id
        ][-limit:]


class ProgressTracker:
    """Main progress tracking system."""

    def __init__(self, db_session_factory=None):
        self.store = ProgressStore()
        self.db_logger = DatabaseLogger(db_session_factory)

    async def start_tracking(
        self,
        job_id: str,
        job_type: JobType,
        total_steps: int,
        message: str = "Starting job...",
        parent_job_id: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None
    ) -> ProgressData:
        """
        Start tracking a new job.

        Args:
            job_id: Unique identifier for the job
            job_type: Type of job (outline, chapter, tts, export)
            total_steps: Total number of steps in the job
            message: Initial status message
            parent_job_id: Parent job ID for sub-jobs (e.g., chapters under outline)
            metadata: Additional metadata to store with the job

        Returns:
            ProgressData object for the new job
        """
        now = time.time()

        data = ProgressData(
            job_id=job_id,
            job_type=job_type,
            status=JobStatus.PENDING,
            total_steps=total_steps,
            current_step=0,
            message=message,
            started_at=now,
            updated_at=now,
            parent_job_id=parent_job_id,
            metadata=metadata or {}
        )

        await self.store.set(job_id, data)
        await self.db_logger.log_event("job_started", data)
        await self.store.notify(job_id, data)

        return data

    async def update_progress(
        self,
        job_id: str,
        current_step: int,
        message: str,
        metadata_update: Optional[Dict[str, Any]] = None
    ) -> Optional[ProgressData]:
        """
        Update progress for a job.

        Args:
            job_id: Job identifier
            current_step: Current step number
            message: Status message
            metadata_update: Additional metadata to merge

        Returns:
            Updated ProgressData or None if job not found
        """
        data = await self.store.get(job_id)
        if not data:
            logger.warning(f"Attempted to update non-existent job: {job_id}")
            return None

        # Record step completion time for ETA calculation
        if current_step > data.current_step:
            self.store.record_step_completion(job_id, data.job_type)

        # Update fields
        data.current_step = current_step
        data.message = message
        data.updated_at = time.time()
        data.status = JobStatus.IN_PROGRESS

        if metadata_update:
            data.metadata.update(metadata_update)

        # Calculate ETA
        data.eta_seconds = self.store.calculate_eta(data)

        await self.store.set(job_id, data)
        await self.db_logger.log_event("progress_updated", data)
        await self.store.notify(job_id, data)

        return data

    async def complete_job(
        self,
        job_id: str,
        message: str = "Job completed successfully"
    ) -> Optional[ProgressData]:
        """
        Mark a job as completed.

        Args:
            job_id: Job identifier
            message: Completion message

        Returns:
            Updated ProgressData or None if job not found
        """
        data = await self.store.get(job_id)
        if not data:
            logger.warning(f"Attempted to complete non-existent job: {job_id}")
            return None

        now = time.time()
        data.status = JobStatus.COMPLETED
        data.current_step = data.total_steps
        data.message = message
        data.updated_at = now
        data.completed_at = now
        data.eta_seconds = 0

        # Record completion for historical stats
        self.store.record_job_completion(data)

        await self.store.set(job_id, data)
        await self.db_logger.log_event("job_completed", data)
        await self.store.notify(job_id, data)

        return data

    async def fail_job(
        self,
        job_id: str,
        error: str
    ) -> Optional[ProgressData]:
        """
        Mark a job as failed.

        Args:
            job_id: Job identifier
            error: Error message describing the failure

        Returns:
            Updated ProgressData or None if job not found
        """
        data = await self.store.get(job_id)
        if not data:
            logger.warning(f"Attempted to fail non-existent job: {job_id}")
            return None

        now = time.time()
        data.status = JobStatus.FAILED
        data.error = error
        data.message = f"Failed: {error}"
        data.updated_at = now
        data.completed_at = now
        data.eta_seconds = None

        await self.store.set(job_id, data)
        await self.db_logger.log_event("job_failed", data)
        await self.store.notify(job_id, data)

        return data

    async def cancel_job(
        self,
        job_id: str,
        reason: str = "Cancelled by user"
    ) -> Optional[ProgressData]:
        """
        Cancel a job.

        Args:
            job_id: Job identifier
            reason: Cancellation reason

        Returns:
            Updated ProgressData or None if job not found
        """
        data = await self.store.get(job_id)
        if not data:
            return None

        now = time.time()
        data.status = JobStatus.CANCELLED
        data.message = reason
        data.updated_at = now
        data.completed_at = now
        data.eta_seconds = None

        await self.store.set(job_id, data)
        await self.db_logger.log_event("job_cancelled", data)
        await self.store.notify(job_id, data)

        return data

    async def get_progress(self, job_id: str) -> Optional[Dict[str, Any]]:
        """
        Get current progress for a job.

        Args:
            job_id: Job identifier

        Returns:
            Progress data as dictionary or None if not found
        """
        data = await self.store.get(job_id)
        if data:
            return data.to_dict()
        return None

    async def get_all_active_jobs(self) -> List[Dict[str, Any]]:
        """
        Get all currently active jobs.

        Returns:
            List of progress data dictionaries for active jobs
        """
        jobs = await self.store.get_all_active()
        return [job.to_dict() for job in jobs]

    async def get_job_with_children(self, job_id: str) -> Optional[Dict[str, Any]]:
        """
        Get a job with all its child jobs.

        Args:
            job_id: Parent job identifier

        Returns:
            Progress data with children or None if not found
        """
        data = await self.store.get(job_id)
        if not data:
            return None

        children = await self.store.get_children(job_id)
        result = data.to_dict()
        result["children"] = [child.to_dict() for child in children]

        # Calculate aggregate progress from children if applicable
        if children:
            total_progress = sum(child.progress_percent for child in children)
            result["aggregate_progress"] = round(total_progress / len(children), 2)

        return result

    async def subscribe_to_progress(
        self,
        job_id: Optional[str] = None
    ) -> AsyncGenerator[str, None]:
        """
        Subscribe to progress updates via Server-Sent Events.

        Args:
            job_id: Optional job ID to subscribe to. If None, subscribes to all jobs.

        Yields:
            SSE formatted strings with progress updates
        """
        queue = await self.store.subscribe(job_id)

        try:
            # Send initial state
            if job_id:
                data = await self.store.get(job_id)
                if data:
                    yield f"event: progress\ndata: {json.dumps(data.to_dict())}\n\n"
            else:
                # Send all active jobs
                jobs = await self.store.get_all_active()
                for job in jobs:
                    yield f"event: progress\ndata: {json.dumps(job.to_dict())}\n\n"

            # Stream updates
            while True:
                try:
                    event_data = await asyncio.wait_for(queue.get(), timeout=30.0)
                    yield f"event: progress\ndata: {json.dumps(event_data)}\n\n"
                except asyncio.TimeoutError:
                    # Send keepalive
                    yield f"event: keepalive\ndata: {json.dumps({'timestamp': time.time()})}\n\n"
        except asyncio.CancelledError:
            pass
        finally:
            await self.store.unsubscribe(queue, job_id)

    async def cleanup_old_jobs(self, max_age_hours: int = 24) -> int:
        """
        Clean up completed/failed jobs older than max_age_hours.

        Args:
            max_age_hours: Maximum age in hours for completed jobs

        Returns:
            Number of jobs cleaned up
        """
        cutoff = time.time() - (max_age_hours * 3600)
        jobs_to_delete = []

        for job_id, data in list(self.store._jobs.items()):
            if data.status in (JobStatus.COMPLETED, JobStatus.FAILED, JobStatus.CANCELLED):
                if data.completed_at and data.completed_at < cutoff:
                    jobs_to_delete.append(job_id)

        for job_id in jobs_to_delete:
            await self.store.delete(job_id)

        return len(jobs_to_delete)


# Global tracker instance
tracker = ProgressTracker()


# FastAPI application for SSE endpoints
@asynccontextmanager
async def lifespan(app: FastAPI):
    """Application lifespan handler."""
    # Startup
    logger.info("Progress tracking system started")
    yield
    # Shutdown
    await tracker.db_logger.flush()
    logger.info("Progress tracking system stopped")


app = FastAPI(
    title="BookForge Progress Tracking",
    description="Real-time progress tracking with SSE",
    lifespan=lifespan
)


class StartJobRequest(BaseModel):
    """Request model for starting a job."""
    job_id: Optional[str] = None
    job_type: JobType
    total_steps: int
    message: Optional[str] = "Starting job..."
    parent_job_id: Optional[str] = None
    metadata: Optional[Dict[str, Any]] = None


class UpdateProgressRequest(BaseModel):
    """Request model for updating progress."""
    current_step: int
    message: str
    metadata_update: Optional[Dict[str, Any]] = None


class FailJobRequest(BaseModel):
    """Request model for failing a job."""
    error: str


@app.post("/api/progress/start")
async def api_start_tracking(request: StartJobRequest):
    """Start tracking a new job."""
    job_id = request.job_id or str(uuid.uuid4())
    data = await tracker.start_tracking(
        job_id=job_id,
        job_type=request.job_type,
        total_steps=request.total_steps,
        message=request.message,
        parent_job_id=request.parent_job_id,
        metadata=request.metadata
    )
    return data.to_dict()


@app.put("/api/progress/{job_id}")
async def api_update_progress(job_id: str, request: UpdateProgressRequest):
    """Update progress for a job."""
    data = await tracker.update_progress(
        job_id=job_id,
        current_step=request.current_step,
        message=request.message,
        metadata_update=request.metadata_update
    )
    if not data:
        raise HTTPException(status_code=404, detail="Job not found")
    return data.to_dict()


@app.post("/api/progress/{job_id}/complete")
async def api_complete_job(job_id: str, message: str = "Job completed successfully"):
    """Mark a job as completed."""
    data = await tracker.complete_job(job_id, message)
    if not data:
        raise HTTPException(status_code=404, detail="Job not found")
    return data.to_dict()


@app.post("/api/progress/{job_id}/fail")
async def api_fail_job(job_id: str, request: FailJobRequest):
    """Mark a job as failed."""
    data = await tracker.fail_job(job_id, request.error)
    if not data:
        raise HTTPException(status_code=404, detail="Job not found")
    return data.to_dict()


@app.post("/api/progress/{job_id}/cancel")
async def api_cancel_job(job_id: str, reason: str = "Cancelled by user"):
    """Cancel a job."""
    data = await tracker.cancel_job(job_id, reason)
    if not data:
        raise HTTPException(status_code=404, detail="Job not found")
    return data.to_dict()


@app.get("/api/progress/{job_id}")
async def api_get_progress(job_id: str, include_children: bool = False):
    """Get progress for a specific job."""
    if include_children:
        data = await tracker.get_job_with_children(job_id)
    else:
        data = await tracker.get_progress(job_id)

    if not data:
        raise HTTPException(status_code=404, detail="Job not found")
    return data


@app.get("/api/progress")
async def api_get_all_active():
    """Get all active jobs."""
    return await tracker.get_all_active_jobs()


@app.get("/api/progress/stream/{job_id}")
async def api_stream_progress(job_id: str):
    """SSE endpoint for streaming progress updates for a specific job."""
    return StreamingResponse(
        tracker.subscribe_to_progress(job_id),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "X-Accel-Buffering": "no",
        }
    )


@app.get("/api/progress/stream")
async def api_stream_all_progress():
    """SSE endpoint for streaming all progress updates."""
    return StreamingResponse(
        tracker.subscribe_to_progress(),
        media_type="text/event-stream",
        headers={
            "Cache-Control": "no-cache",
            "Connection": "keep-alive",
            "X-Accel-Buffering": "no",
        }
    )


# Convenience functions for direct usage
async def start_tracking(
    job_id: str,
    job_type: JobType,
    total_steps: int,
    message: str = "Starting job...",
    parent_job_id: Optional[str] = None,
    metadata: Optional[Dict[str, Any]] = None
) -> ProgressData:
    """Start tracking a new job."""
    return await tracker.start_tracking(
        job_id=job_id,
        job_type=job_type,
        total_steps=total_steps,
        message=message,
        parent_job_id=parent_job_id,
        metadata=metadata
    )


async def update_progress(
    job_id: str,
    current_step: int,
    message: str,
    metadata_update: Optional[Dict[str, Any]] = None
) -> Optional[ProgressData]:
    """Update progress for a job."""
    return await tracker.update_progress(
        job_id=job_id,
        current_step=current_step,
        message=message,
        metadata_update=metadata_update
    )


async def complete_job(
    job_id: str,
    message: str = "Job completed successfully"
) -> Optional[ProgressData]:
    """Mark a job as completed."""
    return await tracker.complete_job(job_id, message)


async def fail_job(
    job_id: str,
    error: str
) -> Optional[ProgressData]:
    """Mark a job as failed."""
    return await tracker.fail_job(job_id, error)


async def get_progress(job_id: str) -> Optional[Dict[str, Any]]:
    """Get current progress for a job."""
    return await tracker.get_progress(job_id)


async def get_all_active_jobs() -> List[Dict[str, Any]]:
    """Get all currently active jobs."""
    return await tracker.get_all_active_jobs()


# Example usage and frontend JavaScript
FRONTEND_JS_EXAMPLE = """
// Frontend JavaScript for consuming SSE progress updates

class ProgressTracker {
    constructor(baseUrl = '') {
        this.baseUrl = baseUrl;
        this.eventSource = null;
        this.callbacks = {
            onProgress: [],
            onComplete: [],
            onError: [],
            onKeepalive: []
        };
    }

    // Subscribe to all progress updates
    subscribeAll() {
        return this._subscribe('/api/progress/stream');
    }

    // Subscribe to a specific job's progress
    subscribeToJob(jobId) {
        return this._subscribe(`/api/progress/stream/${jobId}`);
    }

    _subscribe(endpoint) {
        if (this.eventSource) {
            this.eventSource.close();
        }

        this.eventSource = new EventSource(this.baseUrl + endpoint);

        this.eventSource.addEventListener('progress', (event) => {
            const data = JSON.parse(event.data);
            this._emit('onProgress', data);

            if (data.status === 'completed') {
                this._emit('onComplete', data);
            } else if (data.status === 'failed') {
                this._emit('onError', data);
            }
        });

        this.eventSource.addEventListener('keepalive', (event) => {
            const data = JSON.parse(event.data);
            this._emit('onKeepalive', data);
        });

        this.eventSource.onerror = (error) => {
            console.error('SSE Error:', error);
            // Attempt to reconnect after 5 seconds
            setTimeout(() => this._subscribe(endpoint), 5000);
        };

        return this;
    }

    // Register event callbacks
    onProgress(callback) {
        this.callbacks.onProgress.push(callback);
        return this;
    }

    onComplete(callback) {
        this.callbacks.onComplete.push(callback);
        return this;
    }

    onError(callback) {
        this.callbacks.onError.push(callback);
        return this;
    }

    onKeepalive(callback) {
        this.callbacks.onKeepalive.push(callback);
        return this;
    }

    _emit(event, data) {
        this.callbacks[event].forEach(cb => cb(data));
    }

    // Disconnect from SSE
    disconnect() {
        if (this.eventSource) {
            this.eventSource.close();
            this.eventSource = null;
        }
    }

    // API methods
    async startJob(jobType, totalSteps, options = {}) {
        const response = await fetch(this.baseUrl + '/api/progress/start', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                job_type: jobType,
                total_steps: totalSteps,
                ...options
            })
        });
        return response.json();
    }

    async updateProgress(jobId, currentStep, message, metadata = null) {
        const response = await fetch(this.baseUrl + `/api/progress/${jobId}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                current_step: currentStep,
                message: message,
                metadata_update: metadata
            })
        });
        return response.json();
    }

    async completeJob(jobId, message = 'Job completed successfully') {
        const response = await fetch(
            this.baseUrl + `/api/progress/${jobId}/complete?message=${encodeURIComponent(message)}`,
            { method: 'POST' }
        );
        return response.json();
    }

    async failJob(jobId, error) {
        const response = await fetch(this.baseUrl + `/api/progress/${jobId}/fail`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ error: error })
        });
        return response.json();
    }

    async getProgress(jobId) {
        const response = await fetch(this.baseUrl + `/api/progress/${jobId}`);
        return response.json();
    }

    async getAllActive() {
        const response = await fetch(this.baseUrl + '/api/progress');
        return response.json();
    }
}

// Usage example:
// const tracker = new ProgressTracker('http://localhost:8000');
// tracker.subscribeAll()
//     .onProgress(data => {
//         console.log(`Job ${data.job_id}: ${data.progress_percent}% - ${data.message}`);
//         updateProgressBar(data.job_id, data.progress_percent);
//         if (data.eta_seconds) {
//             updateETA(data.job_id, data.eta_seconds);
//         }
//     })
//     .onComplete(data => {
//         console.log(`Job ${data.job_id} completed!`);
//         showSuccessNotification(data.message);
//     })
//     .onError(data => {
//         console.error(`Job ${data.job_id} failed: ${data.error}`);
//         showErrorNotification(data.error);
//     });
"""


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)

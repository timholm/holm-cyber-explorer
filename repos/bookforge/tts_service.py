"""
TTS Service using Piper for audio generation.

Provides text-to-speech conversion with support for:
- Single text conversion
- Batch chapter conversion
- Markdown cleaning for natural speech
- Progress callbacks
- Async processing
- Queue management
"""

import asyncio
import os
import re
import subprocess
import tempfile
import wave
from concurrent.futures import ThreadPoolExecutor
from dataclasses import dataclass, field
from enum import Enum
from pathlib import Path
from queue import PriorityQueue
from threading import Lock
from typing import Callable, Dict, List, Optional, Tuple, Union
import logging
import time


# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


# Configuration
PIPER_BINARY = "/home/tim/audiobook/venv/bin/piper"
DEFAULT_VOICE_MODEL = "/home/tim/audiobook/voices/en_US-lessac-medium.onnx"
MAX_CHUNK_LENGTH = 5000  # Maximum characters per chunk for processing
DEFAULT_SAMPLE_RATE = 22050  # Piper's default sample rate


class ConversionStatus(Enum):
    """Status of a conversion job."""
    PENDING = "pending"
    PROCESSING = "processing"
    COMPLETED = "completed"
    FAILED = "failed"
    CANCELLED = "cancelled"


@dataclass
class ConversionJob:
    """Represents a single TTS conversion job."""
    job_id: str
    text: str
    output_path: str
    priority: int = 0
    status: ConversionStatus = ConversionStatus.PENDING
    progress: float = 0.0
    error: Optional[str] = None
    created_at: float = field(default_factory=time.time)

    def __lt__(self, other):
        """Compare jobs by priority for queue ordering."""
        return self.priority < other.priority


@dataclass
class Chapter:
    """Represents a chapter for batch conversion."""
    title: str
    content: str
    chapter_number: int
    output_filename: Optional[str] = None


class TTSService:
    """
    Text-to-Speech service using Piper.

    Provides synchronous and asynchronous methods for converting text to audio,
    with support for batch processing and queue management.
    """

    def __init__(
        self,
        piper_binary: str = PIPER_BINARY,
        voice_model: str = DEFAULT_VOICE_MODEL,
        max_workers: int = 2,
        max_chunk_length: int = MAX_CHUNK_LENGTH
    ):
        """
        Initialize the TTS service.

        Args:
            piper_binary: Path to the piper executable
            voice_model: Path to the ONNX voice model
            max_workers: Maximum concurrent conversion workers
            max_chunk_length: Maximum characters per chunk
        """
        self.piper_binary = piper_binary
        self.voice_model = voice_model
        self.max_workers = max_workers
        self.max_chunk_length = max_chunk_length

        # Job management
        self._job_queue: PriorityQueue = PriorityQueue()
        self._jobs: Dict[str, ConversionJob] = {}
        self._job_lock = Lock()
        self._job_counter = 0

        # Thread pool for sync operations
        self._executor = ThreadPoolExecutor(max_workers=max_workers)

        # Progress callbacks
        self._progress_callbacks: Dict[str, Callable[[str, float, ConversionStatus], None]] = {}

        # Processing state
        self._is_processing = False
        self._should_stop = False

    def clean_markdown(self, text: str) -> str:
        """
        Strip markdown formatting for natural speech.

        Removes or converts markdown elements that would sound unnatural
        when read aloud by a TTS system.

        Args:
            text: Text with potential markdown formatting

        Returns:
            Clean text suitable for TTS
        """
        if not text:
            return ""

        cleaned = text

        # Remove code blocks (both fenced and inline)
        cleaned = re.sub(r'```[\s\S]*?```', '', cleaned)
        cleaned = re.sub(r'`[^`]+`', '', cleaned)

        # Convert headers to plain text with pause
        cleaned = re.sub(r'^#{1,6}\s*(.+)$', r'\1.', cleaned, flags=re.MULTILINE)

        # Remove bold/italic markers but keep text
        cleaned = re.sub(r'\*\*\*(.+?)\*\*\*', r'\1', cleaned)
        cleaned = re.sub(r'\*\*(.+?)\*\*', r'\1', cleaned)
        cleaned = re.sub(r'\*(.+?)\*', r'\1', cleaned)
        cleaned = re.sub(r'___(.+?)___', r'\1', cleaned)
        cleaned = re.sub(r'__(.+?)__', r'\1', cleaned)
        cleaned = re.sub(r'_(.+?)_', r'\1', cleaned)

        # Convert links to just the link text
        cleaned = re.sub(r'\[([^\]]+)\]\([^\)]+\)', r'\1', cleaned)

        # Remove images
        cleaned = re.sub(r'!\[([^\]]*)\]\([^\)]+\)', '', cleaned)

        # Convert bullet points to sentences
        cleaned = re.sub(r'^\s*[-*+]\s+', '', cleaned, flags=re.MULTILINE)

        # Convert numbered lists
        cleaned = re.sub(r'^\s*\d+\.\s+', '', cleaned, flags=re.MULTILINE)

        # Remove horizontal rules
        cleaned = re.sub(r'^[-*_]{3,}\s*$', '', cleaned, flags=re.MULTILINE)

        # Remove blockquotes markers but keep text
        cleaned = re.sub(r'^>\s*', '', cleaned, flags=re.MULTILINE)

        # Remove HTML tags
        cleaned = re.sub(r'<[^>]+>', '', cleaned)

        # Clean up extra whitespace
        cleaned = re.sub(r'\n{3,}', '\n\n', cleaned)
        cleaned = re.sub(r' {2,}', ' ', cleaned)

        # Expand common abbreviations for better pronunciation
        abbreviations = {
            r'\betc\.': 'et cetera',
            r'\be\.g\.': 'for example',
            r'\bi\.e\.': 'that is',
            r'\bvs\.': 'versus',
            r'\bDr\.': 'Doctor',
            r'\bMr\.': 'Mister',
            r'\bMrs\.': 'Missus',
            r'\bMs\.': 'Miss',
            r'\bProf\.': 'Professor',
        }

        for pattern, replacement in abbreviations.items():
            cleaned = re.sub(pattern, replacement, cleaned, flags=re.IGNORECASE)

        return cleaned.strip()

    def _chunk_text(self, text: str) -> List[str]:
        """
        Split long text into manageable chunks for processing.

        Attempts to split at natural boundaries (sentences, paragraphs)
        rather than arbitrary positions.

        Args:
            text: Text to split into chunks

        Returns:
            List of text chunks
        """
        if len(text) <= self.max_chunk_length:
            return [text]

        chunks = []
        current_chunk = ""

        # Split by paragraphs first
        paragraphs = text.split('\n\n')

        for paragraph in paragraphs:
            # If a single paragraph is too long, split by sentences
            if len(paragraph) > self.max_chunk_length:
                sentences = re.split(r'(?<=[.!?])\s+', paragraph)
                for sentence in sentences:
                    if len(current_chunk) + len(sentence) + 1 <= self.max_chunk_length:
                        current_chunk += (" " if current_chunk else "") + sentence
                    else:
                        if current_chunk:
                            chunks.append(current_chunk.strip())
                        # If a single sentence is too long, force split
                        if len(sentence) > self.max_chunk_length:
                            words = sentence.split()
                            current_chunk = ""
                            for word in words:
                                if len(current_chunk) + len(word) + 1 <= self.max_chunk_length:
                                    current_chunk += (" " if current_chunk else "") + word
                                else:
                                    chunks.append(current_chunk.strip())
                                    current_chunk = word
                        else:
                            current_chunk = sentence
            else:
                if len(current_chunk) + len(paragraph) + 2 <= self.max_chunk_length:
                    current_chunk += ("\n\n" if current_chunk else "") + paragraph
                else:
                    if current_chunk:
                        chunks.append(current_chunk.strip())
                    current_chunk = paragraph

        if current_chunk:
            chunks.append(current_chunk.strip())

        return chunks

    def _run_piper(self, text: str, output_path: str) -> Tuple[bool, Optional[str]]:
        """
        Run the piper binary to convert text to audio.

        Args:
            text: Text to convert
            output_path: Path for the output WAV file

        Returns:
            Tuple of (success, error_message)
        """
        try:
            # Ensure output directory exists
            os.makedirs(os.path.dirname(output_path) or '.', exist_ok=True)

            # Run piper with subprocess
            process = subprocess.run(
                [
                    self.piper_binary,
                    "--model", self.voice_model,
                    "--output_file", output_path
                ],
                input=text,
                text=True,
                capture_output=True,
                timeout=300  # 5 minute timeout
            )

            if process.returncode != 0:
                error_msg = process.stderr or f"Piper exited with code {process.returncode}"
                return False, error_msg

            if not os.path.exists(output_path):
                return False, "Output file was not created"

            return True, None

        except subprocess.TimeoutExpired:
            return False, "Conversion timed out"
        except Exception as e:
            return False, str(e)

    def convert_to_audio(
        self,
        text: str,
        output_path: str,
        clean_markdown: bool = True,
        progress_callback: Optional[Callable[[float], None]] = None
    ) -> bool:
        """
        Convert text to a WAV audio file.

        Handles long texts by chunking and concatenating the results.

        Args:
            text: Text to convert to speech
            output_path: Path for the output WAV file
            clean_markdown: Whether to clean markdown formatting
            progress_callback: Optional callback for progress updates (0.0-1.0)

        Returns:
            True if conversion succeeded, False otherwise
        """
        if not text or not text.strip():
            logger.warning("Empty text provided for conversion")
            return False

        # Clean the text if requested
        if clean_markdown:
            text = self.clean_markdown(text)

        # Chunk the text for processing
        chunks = self._chunk_text(text)
        total_chunks = len(chunks)

        logger.info(f"Converting text to audio: {len(text)} chars, {total_chunks} chunks")

        if total_chunks == 1:
            # Simple case: single chunk
            success, error = self._run_piper(text, output_path)
            if not success:
                logger.error(f"Conversion failed: {error}")
                return False
            if progress_callback:
                progress_callback(1.0)
            return True

        # Multiple chunks: convert each and concatenate
        temp_files = []
        try:
            for i, chunk in enumerate(chunks):
                temp_file = tempfile.NamedTemporaryFile(
                    suffix='.wav',
                    delete=False
                )
                temp_files.append(temp_file.name)
                temp_file.close()

                success, error = self._run_piper(chunk, temp_file.name)
                if not success:
                    logger.error(f"Chunk {i+1}/{total_chunks} failed: {error}")
                    return False

                if progress_callback:
                    progress_callback((i + 1) / total_chunks * 0.9)  # 90% for conversion

            # Concatenate all chunks
            self._concatenate_wav_files(temp_files, output_path)

            if progress_callback:
                progress_callback(1.0)

            return True

        finally:
            # Clean up temp files
            for temp_file in temp_files:
                try:
                    os.unlink(temp_file)
                except OSError:
                    pass

    def _concatenate_wav_files(self, input_files: List[str], output_path: str) -> None:
        """
        Concatenate multiple WAV files into a single output file.

        Args:
            input_files: List of input WAV file paths
            output_files: Path for the concatenated output
        """
        if not input_files:
            raise ValueError("No input files provided")

        # Read parameters from first file
        with wave.open(input_files[0], 'rb') as first_wav:
            params = first_wav.getparams()

        # Create output file and write all audio data
        with wave.open(output_path, 'wb') as output_wav:
            output_wav.setparams(params)

            for input_file in input_files:
                with wave.open(input_file, 'rb') as input_wav:
                    output_wav.writeframes(input_wav.readframes(input_wav.getnframes()))

    def get_audio_duration(self, wav_path: str) -> float:
        """
        Get the duration of a WAV audio file in seconds.

        Args:
            wav_path: Path to the WAV file

        Returns:
            Duration in seconds
        """
        try:
            with wave.open(wav_path, 'rb') as wav_file:
                frames = wav_file.getnframes()
                rate = wav_file.getframerate()
                return frames / float(rate)
        except Exception as e:
            logger.error(f"Failed to get duration for {wav_path}: {e}")
            return 0.0

    def batch_convert(
        self,
        chapters: List[Chapter],
        output_dir: str,
        progress_callback: Optional[Callable[[int, int, str, float], None]] = None,
        file_prefix: str = "chapter"
    ) -> Dict[int, str]:
        """
        Convert multiple chapters to audio files.

        Args:
            chapters: List of Chapter objects to convert
            output_dir: Directory for output files
            progress_callback: Callback(chapter_num, total, status, progress)
            file_prefix: Prefix for output filenames

        Returns:
            Dictionary mapping chapter numbers to output file paths
        """
        os.makedirs(output_dir, exist_ok=True)
        results = {}
        total_chapters = len(chapters)

        for i, chapter in enumerate(chapters):
            chapter_num = chapter.chapter_number

            # Determine output filename
            if chapter.output_filename:
                output_file = os.path.join(output_dir, chapter.output_filename)
            else:
                safe_title = re.sub(r'[^\w\s-]', '', chapter.title)[:50]
                safe_title = safe_title.strip().replace(' ', '_')
                output_file = os.path.join(
                    output_dir,
                    f"{file_prefix}_{chapter_num:03d}_{safe_title}.wav"
                )

            logger.info(f"Converting chapter {chapter_num}: {chapter.title}")

            # Create chapter-specific progress callback
            def chapter_progress(progress: float, ch_num=chapter_num):
                if progress_callback:
                    progress_callback(ch_num, total_chapters, "converting", progress)

            if progress_callback:
                progress_callback(chapter_num, total_chapters, "starting", 0.0)

            # Add chapter title as intro if present
            full_text = f"{chapter.title}.\n\n{chapter.content}" if chapter.title else chapter.content

            success = self.convert_to_audio(
                full_text,
                output_file,
                clean_markdown=True,
                progress_callback=chapter_progress
            )

            if success:
                results[chapter_num] = output_file
                if progress_callback:
                    progress_callback(chapter_num, total_chapters, "completed", 1.0)
            else:
                logger.error(f"Failed to convert chapter {chapter_num}")
                if progress_callback:
                    progress_callback(chapter_num, total_chapters, "failed", 0.0)

        return results

    # Async methods

    async def convert_to_audio_async(
        self,
        text: str,
        output_path: str,
        clean_markdown: bool = True,
        progress_callback: Optional[Callable[[float], None]] = None
    ) -> bool:
        """
        Asynchronously convert text to audio.

        Args:
            text: Text to convert
            output_path: Output file path
            clean_markdown: Whether to clean markdown
            progress_callback: Progress callback

        Returns:
            True if successful
        """
        loop = asyncio.get_event_loop()
        return await loop.run_in_executor(
            self._executor,
            lambda: self.convert_to_audio(text, output_path, clean_markdown, progress_callback)
        )

    async def batch_convert_async(
        self,
        chapters: List[Chapter],
        output_dir: str,
        progress_callback: Optional[Callable[[int, int, str, float], None]] = None,
        file_prefix: str = "chapter",
        concurrent: bool = False
    ) -> Dict[int, str]:
        """
        Asynchronously convert multiple chapters.

        Args:
            chapters: List of chapters to convert
            output_dir: Output directory
            progress_callback: Progress callback
            file_prefix: Prefix for filenames
            concurrent: If True, convert chapters concurrently (up to max_workers)

        Returns:
            Dictionary mapping chapter numbers to file paths
        """
        if not concurrent:
            # Sequential processing
            loop = asyncio.get_event_loop()
            return await loop.run_in_executor(
                self._executor,
                lambda: self.batch_convert(chapters, output_dir, progress_callback, file_prefix)
            )

        # Concurrent processing
        os.makedirs(output_dir, exist_ok=True)
        results = {}

        async def convert_chapter(chapter: Chapter) -> Tuple[int, Optional[str]]:
            chapter_num = chapter.chapter_number

            if chapter.output_filename:
                output_file = os.path.join(output_dir, chapter.output_filename)
            else:
                safe_title = re.sub(r'[^\w\s-]', '', chapter.title)[:50]
                safe_title = safe_title.strip().replace(' ', '_')
                output_file = os.path.join(
                    output_dir,
                    f"{file_prefix}_{chapter_num:03d}_{safe_title}.wav"
                )

            full_text = f"{chapter.title}.\n\n{chapter.content}" if chapter.title else chapter.content

            success = await self.convert_to_audio_async(
                full_text,
                output_file,
                clean_markdown=True
            )

            return chapter_num, output_file if success else None

        # Use semaphore to limit concurrent conversions
        semaphore = asyncio.Semaphore(self.max_workers)

        async def limited_convert(chapter: Chapter) -> Tuple[int, Optional[str]]:
            async with semaphore:
                return await convert_chapter(chapter)

        tasks = [limited_convert(ch) for ch in chapters]
        completed = await asyncio.gather(*tasks)

        for chapter_num, output_path in completed:
            if output_path:
                results[chapter_num] = output_path

        return results

    # Queue management methods

    def _generate_job_id(self) -> str:
        """Generate a unique job ID."""
        with self._job_lock:
            self._job_counter += 1
            return f"job_{self._job_counter}_{int(time.time() * 1000)}"

    def queue_conversion(
        self,
        text: str,
        output_path: str,
        priority: int = 0,
        progress_callback: Optional[Callable[[str, float, ConversionStatus], None]] = None
    ) -> str:
        """
        Add a conversion job to the queue.

        Args:
            text: Text to convert
            output_path: Output file path
            priority: Job priority (lower = higher priority)
            progress_callback: Callback for job progress updates

        Returns:
            Job ID for tracking
        """
        job_id = self._generate_job_id()

        job = ConversionJob(
            job_id=job_id,
            text=text,
            output_path=output_path,
            priority=priority,
            status=ConversionStatus.PENDING
        )

        with self._job_lock:
            self._jobs[job_id] = job
            self._job_queue.put((priority, job_id))
            if progress_callback:
                self._progress_callbacks[job_id] = progress_callback

        logger.info(f"Queued job {job_id} with priority {priority}")

        return job_id

    def get_job_status(self, job_id: str) -> Optional[ConversionJob]:
        """
        Get the status of a queued job.

        Args:
            job_id: The job ID

        Returns:
            ConversionJob or None if not found
        """
        with self._job_lock:
            return self._jobs.get(job_id)

    def cancel_job(self, job_id: str) -> bool:
        """
        Cancel a pending job.

        Args:
            job_id: The job ID to cancel

        Returns:
            True if cancelled, False if not found or already processing
        """
        with self._job_lock:
            job = self._jobs.get(job_id)
            if job and job.status == ConversionStatus.PENDING:
                job.status = ConversionStatus.CANCELLED
                return True
        return False

    def _update_job_progress(
        self,
        job_id: str,
        progress: float,
        status: ConversionStatus,
        error: Optional[str] = None
    ) -> None:
        """Update job progress and notify callback."""
        with self._job_lock:
            job = self._jobs.get(job_id)
            if job:
                job.progress = progress
                job.status = status
                job.error = error

                callback = self._progress_callbacks.get(job_id)
                if callback:
                    try:
                        callback(job_id, progress, status)
                    except Exception as e:
                        logger.warning(f"Progress callback error: {e}")

    def _process_job(self, job: ConversionJob) -> None:
        """Process a single job from the queue."""
        job_id = job.job_id

        self._update_job_progress(job_id, 0.0, ConversionStatus.PROCESSING)

        def progress_callback(progress: float):
            self._update_job_progress(job_id, progress, ConversionStatus.PROCESSING)

        try:
            success = self.convert_to_audio(
                job.text,
                job.output_path,
                clean_markdown=True,
                progress_callback=progress_callback
            )

            if success:
                self._update_job_progress(job_id, 1.0, ConversionStatus.COMPLETED)
            else:
                self._update_job_progress(
                    job_id, job.progress, ConversionStatus.FAILED,
                    "Conversion failed"
                )
        except Exception as e:
            self._update_job_progress(
                job_id, job.progress, ConversionStatus.FAILED,
                str(e)
            )

    def start_queue_processing(self) -> None:
        """Start processing the job queue in background threads."""
        if self._is_processing:
            logger.warning("Queue processing already running")
            return

        self._is_processing = True
        self._should_stop = False

        def process_loop():
            while not self._should_stop:
                try:
                    # Get next job (with timeout to check stop flag)
                    try:
                        priority, job_id = self._job_queue.get(timeout=1.0)
                    except:
                        continue

                    with self._job_lock:
                        job = self._jobs.get(job_id)

                    if job and job.status == ConversionStatus.PENDING:
                        self._process_job(job)

                except Exception as e:
                    logger.error(f"Queue processing error: {e}")

            self._is_processing = False

        # Start worker threads
        for _ in range(self.max_workers):
            self._executor.submit(process_loop)

        logger.info(f"Started queue processing with {self.max_workers} workers")

    def stop_queue_processing(self) -> None:
        """Stop processing the job queue."""
        self._should_stop = True
        logger.info("Stopping queue processing...")

    def get_queue_stats(self) -> Dict:
        """
        Get statistics about the job queue.

        Returns:
            Dictionary with queue statistics
        """
        with self._job_lock:
            stats = {
                "total_jobs": len(self._jobs),
                "pending": 0,
                "processing": 0,
                "completed": 0,
                "failed": 0,
                "cancelled": 0
            }

            for job in self._jobs.values():
                stats[job.status.value] += 1

            return stats

    def clear_completed_jobs(self) -> int:
        """
        Remove completed and cancelled jobs from tracking.

        Returns:
            Number of jobs cleared
        """
        with self._job_lock:
            to_remove = [
                job_id for job_id, job in self._jobs.items()
                if job.status in (ConversionStatus.COMPLETED, ConversionStatus.CANCELLED, ConversionStatus.FAILED)
            ]

            for job_id in to_remove:
                del self._jobs[job_id]
                self._progress_callbacks.pop(job_id, None)

            return len(to_remove)

    def shutdown(self) -> None:
        """Shutdown the service and clean up resources."""
        self.stop_queue_processing()
        self._executor.shutdown(wait=True)
        logger.info("TTS Service shutdown complete")


# Convenience functions for simple usage

def convert_text_to_audio(
    text: str,
    output_path: str,
    voice_model: str = DEFAULT_VOICE_MODEL,
    clean_markdown: bool = True
) -> bool:
    """
    Simple function to convert text to audio.

    Args:
        text: Text to convert
        output_path: Output WAV file path
        voice_model: Path to voice model
        clean_markdown: Whether to clean markdown

    Returns:
        True if successful
    """
    service = TTSService(voice_model=voice_model)
    return service.convert_to_audio(text, output_path, clean_markdown)


def get_duration(wav_path: str) -> float:
    """
    Get duration of a WAV file.

    Args:
        wav_path: Path to WAV file

    Returns:
        Duration in seconds
    """
    service = TTSService()
    return service.get_audio_duration(wav_path)


# Example usage and testing
if __name__ == "__main__":
    import sys

    # Create service instance
    tts = TTSService()

    # Test markdown cleaning
    test_markdown = """
    # Chapter One: The Beginning

    This is **bold** and *italic* text with a [link](http://example.com).

    - First item
    - Second item

    Here's some `inline code` and a code block:

    ```python
    print("Hello, World!")
    ```

    > This is a blockquote.

    Dr. Smith said, "The results are e.g. quite impressive, etc."
    """

    print("Original text:")
    print(test_markdown)
    print("\nCleaned text:")
    print(tts.clean_markdown(test_markdown))

    # Test simple conversion
    if len(sys.argv) > 1 and sys.argv[1] == "--test":
        test_text = "Hello! This is a test of the Piper text to speech service."
        output = "/tmp/tts_test.wav"

        print(f"\nConverting test text to {output}...")

        def progress(p):
            print(f"Progress: {p*100:.1f}%")

        if tts.convert_to_audio(test_text, output, progress_callback=progress):
            duration = tts.get_audio_duration(output)
            print(f"Success! Duration: {duration:.2f} seconds")
        else:
            print("Conversion failed!")

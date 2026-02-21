#!/usr/bin/env python3
"""
BookForge Worker Launcher

A convenient script to start and manage the job queue workers.
Provides various startup modes and monitoring capabilities.

Usage:
    python start_workers.py                    # Start with default settings
    python start_workers.py -n 8               # Start with 8 workers
    python start_workers.py --status           # Show queue status
    python start_workers.py --reset-stale      # Reset stale jobs
    python start_workers.py --test-jobs 5      # Create test jobs
"""

import os
import sys
import time
import signal
import argparse
import logging
from datetime import datetime
from pathlib import Path

# Add the current directory to path for imports
sys.path.insert(0, os.path.dirname(os.path.abspath(__file__)))

from queue_worker import (
    WorkerPool,
    WorkerConfig,
    DatabaseManager,
    JobQueueManager,
    setup_logging,
    setup_signal_handlers,
    JobStatus,
    JobType,
    create_job
)


def print_banner():
    """Print startup banner."""
    banner = r"""
    ____              __   ______                    
   / __ )____  ____  / /__/ ____/___  _________ ____ 
  / __  / __ \/ __ \/ //_/ /_  / __ \/ ___/ __ \/ _ \
 / /_/ / /_/ / /_/ / ,< / __/ / /_/ / /  / /_/ /  __/
/_____/\____/\____/_/|_/_/    \____/_/   \__, /\___/ 
                                        /____/       
                Job Queue Worker System
    """
    print(banner)


def print_config():
    """Print current configuration."""
    print("\n Configuration:")
    print(f"   Database Path:     {WorkerConfig.DATABASE_PATH}")
    print(f"   Number of Workers: {WorkerConfig.NUM_WORKERS}")
    print(f"   Poll Interval:     {WorkerConfig.POLL_INTERVAL}s")
    print(f"   Max Retries:       {WorkerConfig.MAX_RETRIES}")
    print(f"   Job Timeout:       {WorkerConfig.JOB_TIMEOUT}s")
    print(f"   Log Level:         {WorkerConfig.LOG_LEVEL}")
    print(f"   Log File:          {WorkerConfig.LOG_FILE}")
    print()


def show_status():
    """Show current queue status."""
    print("\n Queue Status:")
    print("-" * 50)
    
    try:
        DatabaseManager.initialize()
        manager = JobQueueManager('status')
        stats = manager.get_queue_stats()
        
        total = sum(stats.values())
        
        status_icons = {
            'pending': '[PEND]',
            'running': '[RUN ]',
            'completed': '[DONE]',
            'failed': '[FAIL]',
            'cancelled': '[CANC]'
        }
        
        for status_name, count in stats.items():
            icon = status_icons.get(status_name, '[----]')
            bar_length = int((count / max(total, 1)) * 30)
            bar = '#' * bar_length + '.' * (30 - bar_length)
            print(f"   {icon} {status_name.capitalize():12} [{bar}] {count:5}")
        
        print("-" * 50)
        print(f"   Total: {total} jobs")
        print()
        
    except Exception as e:
        print(f"   Error getting status: {e}")
        print()


def show_recent_jobs(limit: int = 10):
    """Show recent jobs."""
    print(f"\n Recent Jobs (last {limit}):")
    print("-" * 80)
    
    try:
        DatabaseManager.initialize()
        conn = DatabaseManager.get_connection()
        cursor = conn.cursor()
        
        cursor.execute('''
            SELECT id, job_type, status, project_id, created_at, completed_at
            FROM job_queue
            ORDER BY created_at DESC
            LIMIT ?
        ''', (limit,))
        
        rows = cursor.fetchall()
        
        if not rows:
            print("   No jobs found")
        else:
            print(f"   {'ID':36} {'Type':20} {'Status':12} {'Created'}")
            print("   " + "-" * 76)
            for row in rows:
                job_id = row['id'][:8] + '...'
                job_type = row['job_type'][:18]
                status = row['status']
                created = row['created_at'][:19] if row['created_at'] else 'N/A'
                print(f"   {job_id:36} {job_type:20} {status:12} {created}")
        
        print()
        
    except Exception as e:
        print(f"   Error getting jobs: {e}")
        print()


def reset_stale_jobs():
    """Reset any stale jobs."""
    print("\n Resetting stale jobs...")
    
    try:
        DatabaseManager.initialize()
        manager = JobQueueManager('reset')
        manager.reset_stale_jobs()
        print("   Done!")
    except Exception as e:
        print(f"   Error: {e}")
    print()


def start_workers(num_workers: int = None, foreground: bool = True):
    """Start the worker pool."""
    if num_workers:
        WorkerConfig.NUM_WORKERS = num_workers
    
    logger = setup_logging()
    
    print_config()
    show_status()
    
    print(" Starting workers...")
    print("   Press Ctrl+C to stop\n")
    
    pool = WorkerPool(num_workers=WorkerConfig.NUM_WORKERS)
    setup_signal_handlers(pool)
    
    try:
        pool.start()
        
        if foreground:
            # Monitor loop
            while pool.is_running:
                time.sleep(10)
                status = pool.get_status()
                if status['active_workers'] < status['num_workers']:
                    logger.warning(
                        f"Only {status['active_workers']}/{status['num_workers']} "
                        f"workers active"
                    )
        else:
            pool.wait()
            
    except KeyboardInterrupt:
        print("\n\n Shutdown requested...")
    finally:
        pool.stop()
        print("\n Workers stopped. Goodbye!")


def create_test_jobs(count: int = 5):
    """Create test jobs for development/testing."""
    import sqlite3
    import uuid
    
    print(f"\n Creating {count} test jobs...")
    
    DatabaseManager.initialize()
    conn = DatabaseManager.get_connection()
    cursor = conn.cursor()
    
    # Check for existing projects
    cursor.execute('SELECT id, title FROM projects LIMIT 1')
    project = cursor.fetchone()
    
    if not project:
        # Create a test project
        project_id = str(uuid.uuid4())
        cursor.execute('''
            INSERT INTO projects (id, title, description, genre, status)
            VALUES (?, ?, ?, ?, ?)
        ''', (project_id, 'Test Project', 'A test project for job queue', 'Test', 'draft'))
        print(f"   Created test project (ID: {project_id[:8]}...)")
    else:
        project_id = project['id']
        print(f"   Using existing project: {project['title']}")
    
    # Create test jobs
    job_types = [
        JobType.OUTLINE_GENERATION.value,
        JobType.CHAPTER_GENERATION.value,
        JobType.TTS_CONVERSION.value,
        JobType.STYLE_GUIDE.value,
        JobType.REVISION.value
    ]
    
    for i in range(count):
        job_type = job_types[i % len(job_types)]
        
        # Create chapter for chapter-specific jobs
        chapter_id = None
        if job_type in (JobType.CHAPTER_GENERATION.value, JobType.TTS_CONVERSION.value, JobType.REVISION.value):
            chapter_id = str(uuid.uuid4())
            cursor.execute('''
                INSERT INTO chapters (id, project_id, chapter_number, title, status)
                VALUES (?, ?, ?, ?, ?)
            ''', (chapter_id, project_id, i + 1, f'Test Chapter {i + 1}', 'pending'))
        
        job_id = create_job(
            job_type=job_type,
            project_id=project_id,
            chapter_id=chapter_id,
            payload={'test': True, 'index': i},
            priority=i
        )
        print(f"   Created {job_type} job (ID: {job_id[:8]}...)")
    
    print("   Done!\n")


def cleanup_completed_jobs(days: int = 7):
    """Clean up completed jobs older than specified days."""
    from datetime import datetime, timedelta
    
    print(f"\n Cleaning up completed jobs older than {days} days...")
    
    try:
        DatabaseManager.initialize()
        conn = DatabaseManager.get_connection()
        cursor = conn.cursor()
        
        threshold = (datetime.utcnow() - timedelta(days=days)).isoformat()
        
        cursor.execute('''
            DELETE FROM job_queue 
            WHERE status IN ('completed', 'cancelled') 
            AND completed_at < ?
        ''', (threshold,))
        
        deleted = cursor.rowcount
        print(f"   Deleted {deleted} old jobs")
        
    except Exception as e:
        print(f"   Error: {e}")
    print()


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description='BookForge Job Queue Worker Launcher',
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog="""
Examples:
  %(prog)s                      Start with default settings
  %(prog)s -n 8                 Start with 8 workers
  %(prog)s --status             Show queue status
  %(prog)s --recent             Show recent jobs
  %(prog)s --reset-stale        Reset stale jobs
  %(prog)s --test-jobs 10       Create 10 test jobs
  %(prog)s --cleanup 30         Clean up jobs older than 30 days

Environment Variables:
  BOOKFORGE_DATABASE_PATH       Database file path
  BOOKFORGE_NUM_WORKERS         Number of workers (default: 4)
  BOOKFORGE_POLL_INTERVAL       Poll interval in seconds (default: 2.0)
  BOOKFORGE_MAX_RETRIES         Max retry attempts (default: 3)
  BOOKFORGE_LOG_LEVEL           Logging level (default: INFO)
  BOOKFORGE_LOG_FILE            Log file path
  OLLAMA_URL                    Ollama API URL (default: http://localhost:11434)
  OLLAMA_MODEL                  Ollama model name (default: mistral-nemo)
"""
    )
    
    parser.add_argument(
        '-n', '--num-workers',
        type=int,
        default=WorkerConfig.NUM_WORKERS,
        help=f'Number of worker processes (default: {WorkerConfig.NUM_WORKERS})'
    )
    
    parser.add_argument(
        '-d', '--database',
        type=str,
        default=WorkerConfig.DATABASE_PATH,
        help='Database path'
    )
    
    parser.add_argument(
        '-l', '--log-level',
        type=str,
        default=WorkerConfig.LOG_LEVEL,
        choices=['DEBUG', 'INFO', 'WARNING', 'ERROR'],
        help=f'Logging level (default: {WorkerConfig.LOG_LEVEL})'
    )
    
    parser.add_argument(
        '--status',
        action='store_true',
        help='Show queue status and exit'
    )
    
    parser.add_argument(
        '--recent',
        nargs='?',
        const=10,
        type=int,
        metavar='N',
        help='Show recent N jobs (default: 10)'
    )
    
    parser.add_argument(
        '--reset-stale',
        action='store_true',
        help='Reset stale jobs and exit'
    )
    
    parser.add_argument(
        '--test-jobs',
        type=int,
        metavar='N',
        help='Create N test jobs and exit'
    )
    
    parser.add_argument(
        '--cleanup',
        type=int,
        metavar='DAYS',
        help='Clean up completed jobs older than DAYS'
    )
    
    parser.add_argument(
        '--no-banner',
        action='store_true',
        help='Suppress the startup banner'
    )
    
    args = parser.parse_args()
    
    # Update config from args
    WorkerConfig.DATABASE_PATH = args.database
    WorkerConfig.LOG_LEVEL = args.log_level
    WorkerConfig.NUM_WORKERS = args.num_workers
    
    # Print banner
    if not args.no_banner:
        print_banner()
    
    # Handle commands
    if args.status:
        show_status()
        return
    
    if args.recent is not None:
        show_recent_jobs(args.recent)
        return
    
    if args.reset_stale:
        reset_stale_jobs()
        return
    
    if args.test_jobs:
        create_test_jobs(args.test_jobs)
        return
    
    if args.cleanup:
        cleanup_completed_jobs(args.cleanup)
        return
    
    # Start workers
    start_workers(num_workers=args.num_workers)


if __name__ == '__main__':
    main()

#!/usr/bin/env python3
"""
Bulk import TikToks from various sources.

Usage:
    python bulk_import.py --file urls.txt
    python bulk_import.py --liked @username  # Requires TikTok cookies
    python bulk_import.py --user @username   # Download all from a user
"""

import argparse
import requests
import time
import sys

API_URL = "http://localhost:8000"


def import_from_file(filepath):
    """Import URLs from a text file (one per line)."""
    with open(filepath, 'r') as f:
        urls = [line.strip() for line in f if line.strip() and not line.startswith('#')]

    print(f"Found {len(urls)} URLs to import")

    # Send to API
    response = requests.post(f"{API_URL}/api/download/bulk", json={"urls": urls})

    if response.status_code == 200:
        data = response.json()
        queued = sum(1 for q in data['queued'] if q['status'] == 'queued')
        exists = sum(1 for q in data['queued'] if q['status'] == 'exists')
        print(f"Queued: {queued}, Already archived: {exists}")
    else:
        print(f"Error: {response.text}")


def monitor_queue():
    """Monitor download progress."""
    print("\nMonitoring download queue (Ctrl+C to stop)...")

    while True:
        try:
            response = requests.get(f"{API_URL}/api/queue")
            data = response.json()

            pending = len(data['pending'])
            if pending == 0:
                print("\nAll downloads complete!")
                break

            downloading = [p for p in data['pending'] if p['status'] == 'downloading']
            print(f"\rDownloading: {len(downloading)}, Pending: {pending - len(downloading)}", end='')

            time.sleep(2)
        except KeyboardInterrupt:
            break


def main():
    parser = argparse.ArgumentParser(description='Bulk import TikToks')
    parser.add_argument('--file', '-f', help='File with URLs (one per line)')
    parser.add_argument('--monitor', '-m', action='store_true', help='Monitor queue after import')

    args = parser.parse_args()

    if args.file:
        import_from_file(args.file)
        if args.monitor:
            monitor_queue()
    else:
        parser.print_help()


if __name__ == '__main__':
    main()

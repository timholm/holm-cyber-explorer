#!/usr/bin/env python3
"""
Export archive data for backup or migration.

Usage:
    python export.py --format json --output backup.json
    python export.py --format csv --output backup.csv
"""

import argparse
import requests
import json
import csv
import sys

API_URL = "http://localhost:8000"


def fetch_all_videos():
    """Fetch all videos from the API."""
    videos = []
    page = 1

    while True:
        response = requests.get(f"{API_URL}/api/videos", params={"page": page, "limit": 100})
        data = response.json()

        videos.extend(data['videos'])

        if page >= data['pages']:
            break
        page += 1
        print(f"\rFetching page {page}/{data['pages']}...", end='')

    print(f"\nFetched {len(videos)} videos")
    return videos


def export_json(videos, output):
    """Export to JSON format."""
    with open(output, 'w') as f:
        json.dump(videos, f, indent=2, default=str)
    print(f"Exported to {output}")


def export_csv(videos, output):
    """Export to CSV format."""
    if not videos:
        print("No videos to export")
        return

    fields = ['id', 'tiktok_id', 'url', 'title', 'uploader', 'upload_date',
              'duration', 'view_count', 'like_count', 'archived_at']

    with open(output, 'w', newline='') as f:
        writer = csv.DictWriter(f, fieldnames=fields, extrasaction='ignore')
        writer.writeheader()
        writer.writerows(videos)

    print(f"Exported to {output}")


def main():
    parser = argparse.ArgumentParser(description='Export TikTok archive')
    parser.add_argument('--format', '-f', choices=['json', 'csv'], default='json')
    parser.add_argument('--output', '-o', required=True, help='Output file path')

    args = parser.parse_args()

    videos = fetch_all_videos()

    if args.format == 'json':
        export_json(videos, args.output)
    elif args.format == 'csv':
        export_csv(videos, args.output)


if __name__ == '__main__':
    main()

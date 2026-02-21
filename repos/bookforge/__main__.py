#!/usr/bin/env python3
"""
BookForge CLI Entry Point

This module allows running the BookForge CLI using:
    python -m bookforge <command>

Or directly via the installed entry point:
    bookforge <command>
"""

import sys
import os

# Ensure the package directory is in the path
package_dir = os.path.dirname(os.path.abspath(__file__))
parent_dir = os.path.dirname(package_dir)
if parent_dir not in sys.path:
    sys.path.insert(0, parent_dir)

from bookforge.cli import main

if __name__ == "__main__":
    main()

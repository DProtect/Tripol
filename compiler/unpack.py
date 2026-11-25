"""Transport only: decode GENESIS hex output, with no compilation or code generation."""
import argparse
import os
from pathlib import Path
import sys

def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("output", type=Path)
    args = parser.parse_args()
    data = bytes.fromhex(sys.stdin.read())
    if len(data) < 153 or data[:4] != b"\x7fELF":
        raise SystemExit("Not a GENESIS ELF image")
    with args.output.open("xb") as out:
        out.write(data)
    if os.name == "posix":
        args.output.chmod(0o755)

if __name__ == "__main__":
    main()

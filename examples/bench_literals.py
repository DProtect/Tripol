# Maintainer-only benchmark input generator, not required by TRIPOL users.
import sys
from pathlib import Path
Path(sys.argv[1]).write_text(''.join(f'show "literal number {i}"\n' for i in range(3000)), encoding='utf-8')

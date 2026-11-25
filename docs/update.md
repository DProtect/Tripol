# TRIPOL UPDATE R1

Interactive commands can print a notice on stderr when a newer numbered release is available. Checks are cached for 24 hours (including failed offline attempts), with a 700 ms HTTP deadline on a cache miss. Redirected output, CI, benchmarks, machine-oriented commands and --version are unaffected. Set TRIPOL_NO_UPDATE=1 to disable automatic checks. Explicit commands still work.

- `tripol update --check`: refresh the official GitHub release list and report status.
- `tripol update`: download the platform-specific official distribution into a private cache directory, validate its exact published size and GitHub SHA-256 digest, then print installation instructions.
- No installation, execution or source-code changes occur automatically. Windows setup remains unsigned. Numbered public pre-releases are eligible because TRIPOL is a public preview; drafts and unrecognized version tags are excluded.

No tokens are required. Release checks contact api.github.com and explicit downloads contact GitHub release assets. No source files are sent. The cache holds public release metadata. SHA-256 detects changes to the downloaded asset relative to GitHub's release metadata; it does not replace code signing.

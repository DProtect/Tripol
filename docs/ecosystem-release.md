# Public ecosystem release / Publicación

Repository: https://github.com/Tripollas/Tripol

## Release

After the ecosystem PR is reviewed and merged, open Actions → Public ecosystem → Run workflow on `main`. The workflow scans reachable Git history, tests builds on Linux and Windows, builds the unsigned Windows setup, Linux compiler archive and VSIX, and creates a **draft pre-release** with SHA-256 checksums. Review the actual logs, artifacts and unsigned installer before publishing the draft. Never replace existing release assets silently. If the workflow is not available, follow [build.md](build.md).

## VS Code Marketplace

The repository includes an extension package. Install its VSIX locally for testing. To publish publicly, create/verify a Visual Studio Marketplace publisher and set `vscode/package.json` publisher to that owned ID; use `vsce publish` from a trusted machine with a scoped token stored outside source. Do not expose tokens in CI or logs.

## Community

An administrator must enable Discussions in repository Settings → General → Features → Discussions. Suggested categories: Announcements, Q&A, Ideas, Show and tell. Use the existing code of conduct and security policy. The API connector cannot change this setting.

## Security

The workflow scans Git history of all fetched refs with Gitleaks in redacted mode; it does not scan every external artifact, release, Actions log, unreachable commit or binary. Review those separately. Enable GitHub secret scanning and push protection if available in Settings. Rotate any exposed credentials immediately, even after deleting them from a file. Never publish a real secret in a ticket or discussion. Keep private vulnerability reporting enabled where available.

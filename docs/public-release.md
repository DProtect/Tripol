# Public preview release gate

- [ ] Synchronize the complete ORIGIN snapshot and verify the remote tree matches it.
- [ ] Review all branches, tags, history, releases, artifacts and logs for secrets, personal information and material not authorized for publication; revoke any leaked credentials. A current-tree scan is not a history audit.
- [ ] Confirm all public files are authorized for Apache-2.0 distribution; keep private research outside this repository.
- [ ] Verify Windows and Linux CI, then run live Discord and graphics smoke tests on an isolated test server; no credentials in CI logs.
- [ ] Publish verified installer/portable/source archives with SHA-256 checksums as a pre-release, not v1.0. Installer is unsigned.
- [ ] Enable branch protection, required CI, private vulnerability reporting and community discussions where available.
- [ ] Change visibility in Settings only after approval of these checks. Public Git history may be copied permanently.

The current repository is a private staging area until the maintainer completes these steps. No claim of a public release or full audit is made by this checklist.

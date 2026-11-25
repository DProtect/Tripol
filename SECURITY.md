# Security policy

TRIPOL is pre-1.0 experimental software. Only the latest preview is considered for security fixes; there is no production support guarantee. Do not execute untrusted TRIPOL programs: they can access local files and network APIs.

Report vulnerabilities privately using GitHub private vulnerability reporting if enabled. If unavailable, request a private reporting channel from the maintainer without posting exploit details or secrets publicly. Include version, OS, minimal reproduction with synthetic data and impact. Never include actual Discord tokens or other credentials. Revoke exposed credentials promptly. We will triage reports on a best-effort basis.

Do not publish binaries containing embedded credentials. Use `discord.login_env()` for shareable code. Installer builds are unsigned previews; verify release SHA-256 checksums and use only trusted sources.

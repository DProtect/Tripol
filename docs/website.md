# Website / GitHub Pages

The website is plain HTML, CSS and JavaScript. It uses a black background, white type, no external fonts, no generated images, and a text-only wordmark until the maintainer supplies `assets/raven.svg`. All links and installers are relative, except the public Pages command URLs and GitHub links. The benchmark table is labeled ORIGIN v0.23 (Windows user results), not a v0.24 Rust comparison.

## Local preview
Extract `TRIPOL-v0.24-Website.zip` and open `index.html` or run `python -m http.server 8000` inside that directory. The distribution ZIP contains `downloads/`; the source repository excludes generated binaries. The install scripts verify a pinned SHA-256 before installing; Windows is unsigned. Linux installs to `~/.local/bin` without changing shell profiles.

## Publish manually (no compiler CI required)
Keep the private code repo private until the history audit and permissions review are complete. Once public, create a `gh-pages` branch containing ONLY the extracted website distribution, including `.nojekyll` and `downloads/`. In Settings > Pages select Deploy from a branch > gh-pages > /(root). The project URL is expected to be https://dprotect.github.io/Tripol/; deployment is not verified. GitHub Pages may still require a Pages deployment workflow. If quota prevents it, wait for the reset or resolve the billing issue. Never put secrets in the website branch. Use a clean local directory for the website to avoid accidentally publishing private source/history. Alternatively publish the website in a separate public repository and adjust the URLs in index.html and both scripts.

Before each new version, regenerate the downloads/checksums and update the pinned hashes in both installers. Run `tripol --version` after installation. No automatic updater exists yet.

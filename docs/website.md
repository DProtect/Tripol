# Website / GitHub Pages

The public website is maintained exclusively on the [gh-pages branch](https://github.com/Tripollas/Tripol/tree/gh-pages), not in main. Its public URL is https://tripollas.github.io/Tripol/.

## Local preview

Clone the website branch into a separate directory:

```sh
git clone --branch gh-pages --single-branch https://github.com/Tripollas/Tripol.git tripol-website
cd tripol-website
python -m http.server 8000
```

Open http://localhost:8000. The branch includes static files, documentation and versioned downloads. The installation scripts verify pinned SHA-256 hashes; Windows setup is unsigned.

## Deployment

In Settings > Pages, use Deploy from a branch > gh-pages > /(root). Keep changes to the website on gh-pages and compiler changes on main. Before publishing a new version, regenerate verified downloads and checksums, update installer hashes and links, and test both platforms. Never commit secrets. GitHub Releases and website downloads must be verified independently; no automatic synchronization is implied.

#!/bin/sh
set -eu
[ "$(uname -s)" = Linux ] && [ "$(uname -m)" = x86_64 ] || { echo 'TRIPOL requires Linux x86-64' >&2; exit 1; }
name='TRIPOL-v0.24-Linux-x64.tar.gz'
expected='b39d4833dc246cac77d21052a62e313d7f105a292f235be599aead9bbbd2ee21'
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
base=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
if [ -f "$base/downloads/$name" ]; then cp "$base/downloads/$name" "$tmp/$name"; else curl --fail --location --proto '=https' -o "$tmp/$name" "https://tripollas.github.io/Tripol/downloads/$name"; fi
printf '%s  %s\n' "$expected" "$tmp/$name" | sha256sum -c -
tar -xzf "$tmp/$name" -C "$tmp" tripol
mkdir -p "$HOME/.local/bin"
install -m 755 "$tmp/tripol" "$HOME/.local/bin/tripol"
"$HOME/.local/bin/tripol" --version
printf '%s\n' 'If needed add $HOME/.local/bin to your PATH (no shell configuration modified).'

#!/bin/sh
set -eu
[ "$(uname -s)" = Linux ] && [ "$(uname -m)" = x86_64 ] || { echo 'TRIPOL requires Linux x86-64' >&2; exit 1; }
name='TRIPOL-v0.27-Linux-x64.tar.gz'
expected='4c9afc6ca1d31ce8ca962f18f97090799031f6b9e1a92eb8280f7f1ceb5403d0'
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
base=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
if [ -f "$base/downloads/$name" ]; then cp "$base/downloads/$name" "$tmp/$name"; else curl --fail --location --proto '=https' -o "$tmp/$name" "https://github.com/Tripollas/Tripol/releases/download/v0.27.0/$name"; fi
printf '%s  %s\n' "$expected" "$tmp/$name" | sha256sum -c -
tar -xzf "$tmp/$name" -C "$tmp" tripol
mkdir -p "$HOME/.local/bin"
install -m 755 "$tmp/tripol" "$HOME/.local/bin/tripol"
"$HOME/.local/bin/tripol" --version
printf '%s\n' 'If needed add $HOME/.local/bin to your PATH (no shell configuration modified).'

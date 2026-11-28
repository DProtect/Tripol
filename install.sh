#!/bin/sh
set -eu
[ "$(uname -s)" = Linux ] && [ "$(uname -m)" = x86_64 ] || { echo 'TRIPOL requires Linux x86-64' >&2; exit 1; }
name='TRIPOL-v0.28-Linux-x64.tar.gz'
expected='f0761eb2c1678fb135d470cea6776829aac6a38882109911af86a758b29f55d3'
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
curl --fail --location --proto '=https' -o "$tmp/$name" "https://github.com/Tripollas/Tripol/releases/download/v0.28.0/$name"
printf '%s  %s\n' "$expected" "$tmp/$name" | sha256sum -c -
tar -xzf "$tmp/$name" -C "$tmp" tripol
mkdir -p "$HOME/.local/bin"
install -m 755 "$tmp/tripol" "$HOME/.local/bin/tripol"
"$HOME/.local/bin/tripol" --version
printf '%s\n' 'If needed add $HOME/.local/bin to your PATH (no shell configuration modified).'

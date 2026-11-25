#!/bin/sh
set -eu
[ "$(uname -s)" = Linux ] && [ "$(uname -m)" = x86_64 ] || { echo 'TRIPOL requires Linux x86-64' >&2; exit 1; }
name='TRIPOL-v0.25-Linux-x64.tar.gz'
expected='05673ebe7e9fe9f3eba152859f7d025828cdada295180c48172890d171181855'
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
base=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
if [ -f "$base/downloads/$name" ]; then cp "$base/downloads/$name" "$tmp/$name"; else curl --fail --location --proto '=https' -o "$tmp/$name" "https://github.com/Tripollas/Tripol/releases/download/v0.25.0/$name"; fi
printf '%s  %s\n' "$expected" "$tmp/$name" | sha256sum -c -
tar -xzf "$tmp/$name" -C "$tmp" tripol
mkdir -p "$HOME/.local/bin"
install -m 755 "$tmp/tripol" "$HOME/.local/bin/tripol"
"$HOME/.local/bin/tripol" --version
printf '%s\n' 'If needed add $HOME/.local/bin to your PATH (no shell configuration modified).'

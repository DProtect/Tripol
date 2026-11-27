#!/bin/sh
set -eu
[ "$(uname -s)" = Linux ] && [ "$(uname -m)" = x86_64 ] || { echo 'TRIPOL requires Linux x86-64' >&2; exit 1; }
name='TRIPOL-v0.27-Linux-x64.tar.gz'
expected='708b1cbde74e5d816e04321fda8863b2e795bad4927a4e93d66e811383c5e024'
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT HUP INT TERM
curl --fail --location --proto '=https' -o "$tmp/$name" "https://github.com/Tripollas/Tripol/releases/download/v0.27.1/$name"
printf '%s  %s\n' "$expected" "$tmp/$name" | sha256sum -c -
tar -xzf "$tmp/$name" -C "$tmp" tripol
mkdir -p "$HOME/.local/bin"
install -m 755 "$tmp/tripol" "$HOME/.local/bin/tripol"
"$HOME/.local/bin/tripol" --version
printf '%s\n' 'If needed add $HOME/.local/bin to your PATH (no shell configuration modified).'

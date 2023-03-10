#!/bin/bash
out=$1
if [ "$out" == "" ]; then
    echo "Usage: $0 <output file>"
    exit 1
fi
tmpFile=$(openssl rand -base64 20 | tr -dc 'a-zA-Z0-9' | head -c 5)

termsvg rec --skip-first-line /tmp/${tmpFile}.cast
termsvg export --background-color "#1e1e1e" --text-color "#dfa500" /tmp/${tmpFile}.cast -o ${out} && rm /tmp/${tmpFile}.cast
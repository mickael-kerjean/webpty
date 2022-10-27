#/bin/bash

find `dirname $0`/src -type f | grep -v '\.ts$' | grep -v '\.map$' | xargs gzip -f -k --best

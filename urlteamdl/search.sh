#!/bin/bash
find . -type f -name "*.txt.xz" | xargs -n1 -I{} sh -c "xzcat {} | grep -v '^#' | cut -d'|' -f2 | grep $1"

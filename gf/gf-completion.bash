#!/usr/bin/env bash
complete -W "$(ls ~/.gf | sed -r 's/\.json$//')" gf

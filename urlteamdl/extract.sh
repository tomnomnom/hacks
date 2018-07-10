#!/bin/bash
ls *.zip | xargs -n1 -I{} sh -c 'unzip {} -d $(basename {} .zip)'

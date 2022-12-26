#!/bin/bash

files=(
    "/.terraform"
    "/.terraform.lock.hcl"
    "/*.tfstate"
    "/*.txt"
    "/*.backup"
)

for d in $(find ./examples -maxdepth 4 -type d)
    do
        for i in "${files[@]}"
            do
                rm -rfv $d$i
            done
        echo "Cleaned {$d}"
    done

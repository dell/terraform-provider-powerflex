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

for d in $(find ./examples/ -maxdepth 4 -type d)
    do
        echo "Removing sensitive data - $d/main.tf"
        grep -vs "username" $d"/main.tf" > tmpfile && mv tmpfile $d"/main.tf"
        grep -vs "password" $d"/main.tf" > tmpfile && mv tmpfile $d"/main.tf"
        grep -vs "host" $d"/main.tf" > tmpfile && mv tmpfile $d"/main.tf"
    done
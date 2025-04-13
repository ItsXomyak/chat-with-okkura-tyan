#!/bin/bash

commit_message="$*"

git add .

git commit -m "$commit_message"

git push

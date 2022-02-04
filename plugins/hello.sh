#!/usr/bin/env bash

INPUT="${@}"

if [[ -z "${INPUT}" ]]; then
  echo "Hi there! :)"
else
  echo "Hi ${INPUT}"
fi

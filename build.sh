#!/bin/bash
for micro in `ls -d ./micros/*/`; do
  cd $micro
  echo "Building ${micro}:"
  go build
  status=$?
  if [ $status -eq 0 ]; then
    echo " ${micro} successfully built"
  else
    echo " Failed to build ${micro}"
  fi
  cd -
done

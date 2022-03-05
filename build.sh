#!/bin/bash

ACTION=$1

function build() {
  cd ./micros
  for micro in `ls`; do
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
}

function cleanup() {
  cd ./micros
  for micro in `ls`; do
    cd $micro
    echo "Cleaning ${micro}:"
    rm $micro
    cd -
  done
}

case "$ACTION" in
  "build")
    build
    ;;
  "cleanup")
    cleanup
    ;;
esac

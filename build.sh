#!/bin/bash

ACTION=${1:-"build"}
BUILDDIR="./micros"

function build() {
  cd $BUILDDIR
  local bd=$(pwd)
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
    cd $bd
  done
}

function cleanup() {
  cd $BUILDDIR
  local bd=$(pwd)
  for micro in `ls`; do
    cd $micro
    echo "Cleaning ${micro}:"
    rm $micro
    cd $bd
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

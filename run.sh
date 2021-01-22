#!/usr/bin/env bash

set -e

# Working directory to mount on docker run
WD="$(
  cd "$(dirname "$0")" >/dev/null 2>&1
  pwd -P
)"

# Docker image tag
IMAGE=recipe-collector

VOLUMES=(
  "${WD}"
)

# Check if docker is installed
if ! type docker &>/dev/null; then
  echo "docker is required to run this script!"
  exit 1
fi

function print_usage() {
  echo "Usage: $0 <build|run>"
  echo ""
  echo "Commands:"
  echo "  build    Build docker image with binary. Use it to rebuild image"
  echo "  run      Run recipe-collector (image will be built implicitly)"
}

function build_image() {
  set +e
  docker images | grep "${IMAGE}" >/dev/null
  if [[ "$1" == "1" || "$?" == "1" ]]; then
    docker build --target runtime --tag "${IMAGE}" "${WD}"
  fi
  set -e
}

function run() {
  local cmd=(
    "run"
    "--rm"
    "-w"
    "${WD}"
  )

  for mount in "${VOLUMES[@]}"; do
    cmd+=("-v")
    cmd+=("${mount}:${mount}")
  done

  cmd+=("${IMAGE}" "$@")

  docker "${cmd[@]}"
}

MODE=$1
if [[ -z ${MODE} ]]; then
  MODE="run"
else
  shift
fi

# Run commands
if [[ ${MODE} == "build" ]]; then
  build_image "1"
elif [[ ${MODE} == "run" ]]; then
  build_image "0"
  run "$@"
else
  print_usage
fi

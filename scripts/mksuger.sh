#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


set -o errexit

declare -A map=(
  # 11 useable components (@2023.07.21)
  ["demo"]="zero-demo" # a demo web server
  ["uc"]="zero-usercenter"
  ["api"]="zero-apiserver"
  ["gw"]="zero-gateway"
  ["nw"]="zero-nightwatch"
  ["pmp"]="zero-pump"
  ["tblc"]="zero-toyblc"
  ["cm"]="zero-controller-manager"
  ["msc"]="zero-minerset-controller"
  ["mc"]="zero-miner-controller"
  ["zctl"]="zeroctl"
)

ZROOT=$(dirname "${BASH_SOURCE[0]}")/..
ZVERSION=${ZVERSION:-$(date +'%Y%m%d%H%M%S')}

SHORTNAME=${SHORTNAME:=uc}
COMPONENTS=()
DEPLOY=false
IMAGE=false
LOAD=false
ZERO_KIND_CLUSTER_NAME=${ZERO_KIND_CLUSTER_NAME:-zero}
NODES=${KIND_LOAD_NODES:-$(kubectl get nodes|awk '/Ready/ && !/SchedulingDisabled/{nodes=nodes$1","} END{gsub(/,$/,"",nodes);print nodes}')}

# Get file names from COMMAND LINE arguments
function getcomponents() {
  for f in "$@"; do
    COMPONENTS[${#COMPONENTS[*]}]=${map[$1]:-$1}
  done
}

# Load docker images to kind cluster
function load_docker_images() {
  for comp in "${COMPONENTS[@]}"
  do
    kind load docker-image --name ${ZERO_KIND_CLUSTER_NAME} --nodes ${NODES} ccr.ccs.tencentyun.com/superproj/${comp}-amd64:${ZVERSION}
  done
}

# Build docker images
function build_image() {
  for comp in "${COMPONENTS[@]}"
  do
    make -C ${ZROOT} image IMAGES=${comp} VERSION=${ZVERSION} MULTISTAGE=0
    [[ "$LOAD" == true ]] && load_docker_images
  done
}

# Only build component
function build() {
  for comp in "${COMPONENTS[@]}"
  do
    make -C ${ZROOT} build BINS=${comp} VERSION=${ZVERSION}
  done
}

# Build docker images and deploy them
function deploy() {
  for comp in "${COMPONENTS[@]}"
  do
    make -C ${ZROOT} deploy DEPLOYS=${comp} VERSION=${ZVERSION}
    load_docker_images
    kubectl rollout restart deployment ${comp}
  done
}

# Print usage infomation
function Usage()
{
  cat << EOF

Usage: $0 [ OPTIONS ] SHORTNAME [-d]
build suger script.

  SHORTNAME              short name for zero component.

OPTIONS:
  -h, --help             usage information.
  -d, --deploy           whether to deploy component to kind cluster (build image and deploy).
  -i, --image            build image only.
      --load             load docker image to kind cluster. Only work when `-i` options is specified.
  -v, --version          build or deploy version.

Reprot bugs to <colin404@foxmail.com>.
EOF
}

# Print message to standerr
function die()
{
  echo "$@" >&2
  exit 1
}

# Check the argument associate with a option
function requiredarg()
{
  [ -z "$2" -o "$(echo $2 | awk '$0~/^-/{print 1}')" == "1" ] && die "$0: option $1 requires an argument"
  ((args++))
}


### read cli options
# separate groups of short options. replace --foo=bar with --foo bar
while [[ -n $1 ]]; do
  case "$1" in
    -- )
      for arg in "$@"; do
        ARGS[${#ARGS[*]}]="$arg"
      done
      break
      ;;
    --*=?* )
      ARGS[${#ARGS[*]}]="${1%%=*}"
      ARGS[${#ARGS[*]}]="${1#*=}"
      ;;
    --* )
      #die "$0: option $1 requires a value"
      ARGS[${#ARGS[*]}]="$1"
      ;;
    -* )
      for shortarg in $(sed -e 's|.| -&|g' <<< "${1#-}"); do
        ARGS[${#ARGS[*]}]="$shortarg"
      done
      ;;
    * )
      ARGS[${#ARGS[*]}]="$1"
  esac
  shift
done

# set the separated options as input options.
set -- "${ARGS[@]}"

while [[ -n $1 ]]; do
  ((args=1))
  case "$1" in
    -h | --help )
      Usage
      exit 0
      ;;
    -d | --deploy )
      DEPLOY="true"
      ;;
    -i | --image )
      IMAGE="true"
      ;;
    --load )
      LOAD="true"
      ;;
    -v | --version )
      requiredarg "$@"
      ZVERSION="$2"
      ;;
    -* )
      die "$0: unrecognized option '$1'"
      ;;
    *)
      getcomponents "$1"
      ;;
  esac
  shift $args
done

if [  "${#COMPONENTS[@]}" -eq 0 ];then
  COMPONENTS=("${map[@]}")
fi

if [ "${DEPLOY}" == true ];then
  deploy
elif [ "${IMAGE}" == true ];then
  build_image
else
  build
fi

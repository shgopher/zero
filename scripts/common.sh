#!/usr/bin/env bash

# Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/superproj/zero.
#


# shellcheck disable=SC2034 # Variables sourced in other scripts.

# Common utilities, variables and checks for all build scripts.
set -o errexit
set -o nounset
set -o pipefail

# Unset CDPATH, having it set messes up with script import paths
unset CDPATH

USER_ID=$(id -u)
GROUP_ID=$(id -g)

DOCKER_OPTS=${DOCKER_OPTS:-""}
IFS=" " read -r -a DOCKER <<< "docker ${DOCKER_OPTS}"
DOCKER_HOST=${DOCKER_HOST:-""}
DOCKER_MACHINE_NAME=${DOCKER_MACHINE_NAME:-"miner-dev"}
readonly DOCKER_MACHINE_DRIVER=${DOCKER_MACHINE_DRIVER:-"virtualbox --virtualbox-cpu-count -1"}

# This will canonicalize the path
ZROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd -P)

source "${ZROOT}/scripts/lib/init.sh"

# Constants
readonly PROJ_BUILD_IMAGE_REPO=miner-build
#readonly PROJ_BUILD_IMAGE_CROSS_TAG="$(cat "${ZROOT}/build/build-image/cross/VERSION")"

readonly PROJ_DOCKER_REGISTRY="${PROJ_DOCKER_REGISTRY:-k8s.gcr.io}"
readonly PROJ_BASE_IMAGE_REGISTRY="${PROJ_BASE_IMAGE_REGISTRY:-us.gcr.io/k8s-artifacts-prod/build-image}"

# This version number is used to cause everyone to rebuild their data containers
# and build image.  This is especially useful for automated build systems like
# Jenkins.
#
# Increment/change this number if you change the build image (anything under
# build/build-image) or change the set of volumes in the data container.
#readonly PROJ_BUILD_IMAGE_VERSION_BASE="$(cat "${ZROOT}/build/build-image/VERSION")"
#readonly PROJ_BUILD_IMAGE_VERSION="${PROJ_BUILD_IMAGE_VERSION_BASE}-${PROJ_BUILD_IMAGE_CROSS_TAG}"

# Here we map the output directories across both the local and remote _output
# directories:
#
# *_OUTPUT_ROOT    - the base of all output in that environment.
# *_OUTPUT_SUBPATH - location where golang stuff is built/cached.  Also
#                    persisted across docker runs with a volume mount.
# *_OUTPUT_BINPATH - location where final binaries are placed.  If the remote
#                    is really remote, this is the stuff that has to be copied
#                    back.
# OUT_DIR can come in from the Makefile, so honor it.
readonly LOCAL_OUTPUT_ROOT="${ZROOT}/${OUT_DIR:-_output}"
readonly LOCAL_OUTPUT_SUBPATH="${LOCAL_OUTPUT_ROOT}/platforms"
readonly LOCAL_OUTPUT_BINPATH="${LOCAL_OUTPUT_SUBPATH}"
readonly LOCAL_OUTPUT_GOPATH="${LOCAL_OUTPUT_SUBPATH}/go"
readonly LOCAL_OUTPUT_IMAGE_STAGING="${LOCAL_OUTPUT_ROOT}/images"

# This is the port on the workstation host to expose RSYNC on.  Set this if you
# are doing something fancy with ssh tunneling.
readonly PROJ_RSYNC_PORT="${PROJ_RSYNC_PORT:-}"

# This is the port that rsync is running on *inside* the container. This may be
# mapped to PROJ_RSYNC_PORT via docker networking.
readonly PROJ_CONTAINER_RSYNC_PORT=8730

# Get the set of master binaries that run in Docker (on Linux)
# Entry format is "<name-of-binary>,<base-image>".
# Binaries are placed in /usr/local/bin inside the image.
#
# $1 - server architecture
zero::build::get_docker_wrapped_binaries() {
  local arch=$1
  local debian_base_version=v2.1.0
  local debian_iptables_version=v12.1.0
  ### If you change any of these lists, please also update DOCKERIZED_BINARIES
  ### in build/BUILD. And zero::golang::server_image_targets
  local targets=(
    "zero-apiserver,${PROJ_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "zero-controller-manager,${PROJ_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "zero-scheduler,${PROJ_BASE_IMAGE_REGISTRY}/debian-base-${arch}:${debian_base_version}"
    "zero-proxy,${PROJ_BASE_IMAGE_REGISTRY}/debian-iptables-${arch}:${debian_iptables_version}"
  )

  echo "${targets[@]}"
}

# ---------------------------------------------------------------------------
# Basic setup functions

# Verify that the right utilities and such are installed for building miner. Set
# up some dynamic constants.
# Args:
#   $1 - boolean of whether to require functioning docker (default true)
#
# Vars set:
#   ZROOT_HASH
#   PROJ_BUILD_IMAGE_TAG_BASE
#   PROJ_BUILD_IMAGE_TAG
#   PROJ_BUILD_IMAGE
#   PROJ_BUILD_CONTAINER_NAME_BASE
#   PROJ_BUILD_CONTAINER_NAME
#   PROJ_DATA_CONTAINER_NAME_BASE
#   PROJ_DATA_CONTAINER_NAME
#   PROJ_RSYNC_CONTAINER_NAME_BASE
#   PROJ_RSYNC_CONTAINER_NAME
#   DOCKER_MOUNT_ARGS
#   LOCAL_OUTPUT_BUILD_CONTEXT
function zero::build::verify_prereqs() {
  local -r require_docker=${1:-true}
  zero::log::status "Verifying Prerequisites...."
  zero::build::ensure_tar || return 1
  zero::build::ensure_rsync || return 1
  if ${require_docker}; then
    zero::build::ensure_docker_in_path || return 1
    zero::util::ensure_docker_daemon_connectivity || return 1

    if (( PROJ_VERBOSE > 6 )); then
      zero::log::status "Docker Version:"
      "${DOCKER[@]}" version | zero::log::info_from_stdin
    fi
  fi

  PROJ_GIT_BRANCH=$(git symbolic-ref --short -q HEAD 2>/dev/null || true)
  ZROOT_HASH=$(zero::build::short_hash "${HOSTNAME:-}:${ZROOT}:${PROJ_GIT_BRANCH}")
  PROJ_BUILD_IMAGE_TAG_BASE="build-${ZROOT_HASH}"
  #PROJ_BUILD_IMAGE_TAG="${PROJ_BUILD_IMAGE_TAG_BASE}-${PROJ_BUILD_IMAGE_VERSION}"
  #PROJ_BUILD_IMAGE="${PROJ_BUILD_IMAGE_REPO}:${PROJ_BUILD_IMAGE_TAG}"
  PROJ_BUILD_CONTAINER_NAME_BASE="miner-build-${ZROOT_HASH}"
  #PROJ_BUILD_CONTAINER_NAME="${PROJ_BUILD_CONTAINER_NAME_BASE}-${PROJ_BUILD_IMAGE_VERSION}"
  PROJ_RSYNC_CONTAINER_NAME_BASE="miner-rsync-${ZROOT_HASH}"
  #PROJ_RSYNC_CONTAINER_NAME="${PROJ_RSYNC_CONTAINER_NAME_BASE}-${PROJ_BUILD_IMAGE_VERSION}"
  PROJ_DATA_CONTAINER_NAME_BASE="miner-build-data-${ZROOT_HASH}"
  #PROJ_DATA_CONTAINER_NAME="${PROJ_DATA_CONTAINER_NAME_BASE}-${PROJ_BUILD_IMAGE_VERSION}"
  #DOCKER_MOUNT_ARGS=(--volumes-from "${PROJ_DATA_CONTAINER_NAME}")
  #LOCAL_OUTPUT_BUILD_CONTEXT="${LOCAL_OUTPUT_IMAGE_STAGING}/${PROJ_BUILD_IMAGE}"

  zero::version::get_version_vars
  #zero::version::save_version_vars "${ZROOT}/.dockerized-miner-version-defs"
}

# ---------------------------------------------------------------------------
# Utility functions

function zero::build::docker_available_on_osx() {
  if [[ -z "${DOCKER_HOST}" ]]; then
    if [[ -S "/var/run/docker.sock" ]]; then
      zero::log::status "Using Docker for MacOS"
      return 0
    fi

    zero::log::status "No docker host is set. Checking options for setting one..."
    if [[ -z "$(which docker-machine)" ]]; then
      zero::log::status "It looks like you're running Mac OS X, yet neither Docker for Mac nor docker-machine can be found."
      zero::log::status "See: https://docs.docker.com/engine/installation/mac/ for installation instructions."
      return 1
    elif [[ -n "$(which docker-machine)" ]]; then
      zero::build::prepare_docker_machine
    fi
  fi
}

function zero::build::prepare_docker_machine() {
  zero::log::status "docker-machine was found."

  local available_memory_bytes
  available_memory_bytes=$(sysctl -n hw.memsize 2>/dev/null)

  local bytes_in_mb=1048576

  # Give virtualbox 1/2 the system memory. Its necessary to divide by 2, instead
  # of multiple by .5, because bash can only multiply by ints.
  local memory_divisor=2

  local virtualbox_memory_mb=$(( available_memory_bytes / (bytes_in_mb * memory_divisor) ))

  docker-machine inspect "${DOCKER_MACHINE_NAME}" &> /dev/null || {
    zero::log::status "Creating a machine to build PROJ"
    docker-machine create --driver "${DOCKER_MACHINE_DRIVER}" \
      --virtualbox-memory "${virtualbox_memory_mb}" \
      --engine-env HTTP_PROXY="${PROJRNETES_HTTP_PROXY:-}" \
      --engine-env HTTPS_PROXY="${PROJRNETES_HTTPS_PROXY:-}" \
      --engine-env NO_PROXY="${PROJRNETES_NO_PROXY:-127.0.0.1}" \
      "${DOCKER_MACHINE_NAME}" > /dev/null || {
      zero::log::error "Something went wrong creating a machine."
      zero::log::error "Try the following: "
      zero::log::error "docker-machine create -d ${DOCKER_MACHINE_DRIVER} --virtualbox-memory ${virtualbox_memory_mb} ${DOCKER_MACHINE_NAME}"
      return 1
    }
  }
  docker-machine start "${DOCKER_MACHINE_NAME}" &> /dev/null
  # it takes `docker-machine env` a few seconds to work if the machine was just started
  local docker_machine_out
  while ! docker_machine_out=$(docker-machine env "${DOCKER_MACHINE_NAME}" 2>&1); do
    if [[ ${docker_machine_out} =~ "Error checking TLS connection" ]]; then
      echo "${docker_machine_out}"
      docker-machine regenerate-certs "${DOCKER_MACHINE_NAME}"
    else
      sleep 1
    fi
  done
  eval "$(docker-machine env "${DOCKER_MACHINE_NAME}")"
  zero::log::status "A Docker host using docker-machine named '${DOCKER_MACHINE_NAME}' is ready to go!"
  return 0
}

function zero::build::is_gnu_sed() {
  [[ $(sed --version 2>&1) == *GNU* ]]
}

function zero::build::ensure_rsync() {
  if [[ -z "$(which rsync)" ]]; then
    zero::log::error "Can't find 'rsync' in PATH, please fix and retry."
    return 1
  fi
}

function zero::build::update_dockerfile() {
  if zero::build::is_gnu_sed; then
    sed_opts=(-i)
  else
    sed_opts=(-i '')
  fi
  sed "${sed_opts[@]}" "s/PROJ_BUILD_IMAGE_CROSS_TAG/${PROJ_BUILD_IMAGE_CROSS_TAG}/" "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
}

function  zero::build::set_proxy() {
  if [[ -n "${PROJRNETES_HTTPS_PROXY:-}" ]]; then
    echo "ENV https_proxy $PROJRNETES_HTTPS_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
  if [[ -n "${PROJRNETES_HTTP_PROXY:-}" ]]; then
    echo "ENV http_proxy $PROJRNETES_HTTP_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
  if [[ -n "${PROJRNETES_NO_PROXY:-}" ]]; then
    echo "ENV no_proxy $PROJRNETES_NO_PROXY" >> "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  fi
}

function zero::build::ensure_docker_in_path() {
  if [[ -z "$(which docker)" ]]; then
    zero::log::error "Can't find 'docker' in PATH, please fix and retry."
    zero::log::error "See https://docs.docker.com/installation/#installation for installation instructions."
    return 1
  fi
}

function zero::build::ensure_tar() {
  if [[ -n "${TAR:-}" ]]; then
    return
  fi

  # Find gnu tar if it is available, bomb out if not.
  TAR=tar
  if which gtar &>/dev/null; then
      TAR=gtar
  else
      if which gnutar &>/dev/null; then
	  TAR=gnutar
      fi
  fi
  if ! "${TAR}" --version | grep -q GNU; then
    echo "  !!! Cannot find GNU tar. Build on Linux or install GNU tar"
    echo "      on Mac OS X (brew install gnu-tar)."
    return 1
  fi
}

function zero::build::has_docker() {
  which docker &> /dev/null
}

function zero::build::has_ip() {
  which ip &> /dev/null && ip -Version | grep 'iproute2' &> /dev/null
}

# Detect if a specific image exists
#
# $1 - image repo name
# $2 - image tag
function zero::build::docker_image_exists() {
  [[ -n $1 && -n $2 ]] || {
    zero::log::error "Internal error. Image not specified in docker_image_exists."
    exit 2
  }

  [[ $("${DOCKER[@]}" images -q "${1}:${2}") ]]
}

# Delete all images that match a tag prefix except for the "current" version
#
# $1: The image repo/name
# $2: The tag base. We consider any image that matches $2*
# $3: The current image not to delete if provided
function zero::build::docker_delete_old_images() {
  # In Docker 1.12, we can replace this with
  #    docker images "$1" --format "{{.Tag}}"
  for tag in $("${DOCKER[@]}" images "${1}" | tail -n +2 | awk '{print $2}') ; do
    if [[ "${tag}" != "${2}"* ]] ; then
      V=3 zero::log::status "Keeping image ${1}:${tag}"
      continue
    fi

    if [[ -z "${3:-}" || "${tag}" != "${3}" ]] ; then
      V=2 zero::log::status "Deleting image ${1}:${tag}"
      "${DOCKER[@]}" rmi "${1}:${tag}" >/dev/null
    else
      V=3 zero::log::status "Keeping image ${1}:${tag}"
    fi
  done
}

# Stop and delete all containers that match a pattern
#
# $1: The base container prefix
# $2: The current container to keep, if provided
function zero::build::docker_delete_old_containers() {
  # In Docker 1.12 we can replace this line with
  #   docker ps -a --format="{{.Names}}"
  for container in $("${DOCKER[@]}" ps -a | tail -n +2 | awk '{print $NF}') ; do
    if [[ "${container}" != "${1}"* ]] ; then
      V=3 zero::log::status "Keeping container ${container}"
      continue
    fi
    if [[ -z "${2:-}" || "${container}" != "${2}" ]] ; then
      V=2 zero::log::status "Deleting container ${container}"
      zero::build::destroy_container "${container}"
    else
      V=3 zero::log::status "Keeping container ${container}"
    fi
  done
}

# Takes $1 and computes a short has for it. Useful for unique tag generation
function zero::build::short_hash() {
  [[ $# -eq 1 ]] || {
    zero::log::error "Internal error.  No data based to short_hash."
    exit 2
  }

  local short_hash
  if which md5 >/dev/null 2>&1; then
    short_hash=$(md5 -q -s "$1")
  else
    short_hash=$(echo -n "$1" | md5sum)
  fi
  echo "${short_hash:0:10}"
}

# Pedantically kill, wait-on and remove a container. The -f -v options
# to rm don't actually seem to get the job done, so force kill the
# container, wait to ensure it's stopped, then try the remove. This is
# a workaround for bug https://github.com/docker/docker/issues/3968.
function zero::build::destroy_container() {
  "${DOCKER[@]}" kill "$1" >/dev/null 2>&1 || true
  if [[ $("${DOCKER[@]}" version --format '{{.Server.Version}}') = 17.06.0* ]]; then
    # Workaround https://github.com/moby/moby/issues/33948.
    # TODO: remove when 17.06.0 is not relevant anymore
    DOCKER_API_VERSION=v1.29 "${DOCKER[@]}" wait "$1" >/dev/null 2>&1 || true
  else
    "${DOCKER[@]}" wait "$1" >/dev/null 2>&1 || true
  fi
  "${DOCKER[@]}" rm -f -v "$1" >/dev/null 2>&1 || true
}

# ---------------------------------------------------------------------------
# Building


function zero::build::clean() {
  if zero::build::has_docker ; then
    zero::build::docker_delete_old_containers "${PROJ_BUILD_CONTAINER_NAME_BASE}"
    zero::build::docker_delete_old_containers "${PROJ_RSYNC_CONTAINER_NAME_BASE}"
    zero::build::docker_delete_old_containers "${PROJ_DATA_CONTAINER_NAME_BASE}"
    zero::build::docker_delete_old_images "${PROJ_BUILD_IMAGE_REPO}" "${PROJ_BUILD_IMAGE_TAG_BASE}"

    V=2 zero::log::status "Cleaning all untagged docker images"
    "${DOCKER[@]}" rmi "$("${DOCKER[@]}" images -q --filter 'dangling=true')" 2> /dev/null || true
  fi

  if [[ -d "${LOCAL_OUTPUT_ROOT}" ]]; then
    zero::log::status "Removing _output directory"
    rm -rf "${LOCAL_OUTPUT_ROOT}"
  fi
}

# Set up the context directory for the miner-build image and build it.
function zero::build::build_image() {
  mkdir -p "${LOCAL_OUTPUT_BUILD_CONTEXT}"
  # Make sure the context directory owned by the right user for syncing sources to container.
  chown -R "${USER_ID}":"${GROUP_ID}" "${LOCAL_OUTPUT_BUILD_CONTEXT}"

  cp /etc/localtime "${LOCAL_OUTPUT_BUILD_CONTEXT}/"

  cp "${ZROOT}/build/build-image/Dockerfile" "${LOCAL_OUTPUT_BUILD_CONTEXT}/Dockerfile"
  cp "${ZROOT}/build/build-image/rsyncd.sh" "${LOCAL_OUTPUT_BUILD_CONTEXT}/"
  dd if=/dev/urandom bs=512 count=1 2>/dev/null | LC_ALL=C tr -dc 'A-Za-z0-9' | dd bs=32 count=1 2>/dev/null > "${LOCAL_OUTPUT_BUILD_CONTEXT}/rsyncd.password"
  chmod go= "${LOCAL_OUTPUT_BUILD_CONTEXT}/rsyncd.password"

  zero::build::update_dockerfile
  zero::build::set_proxy
  zero::build::docker_build "${PROJ_BUILD_IMAGE}" "${LOCAL_OUTPUT_BUILD_CONTEXT}" 'false'

  # Clean up old versions of everything
  zero::build::docker_delete_old_containers "${PROJ_BUILD_CONTAINER_NAME_BASE}" "${PROJ_BUILD_CONTAINER_NAME}"
  zero::build::docker_delete_old_containers "${PROJ_RSYNC_CONTAINER_NAME_BASE}" "${PROJ_RSYNC_CONTAINER_NAME}"
  zero::build::docker_delete_old_containers "${PROJ_DATA_CONTAINER_NAME_BASE}" "${PROJ_DATA_CONTAINER_NAME}"
  zero::build::docker_delete_old_images "${PROJ_BUILD_IMAGE_REPO}" "${PROJ_BUILD_IMAGE_TAG_BASE}" "${PROJ_BUILD_IMAGE_TAG}"

  zero::build::ensure_data_container
  zero::build::sync_to_container
}

# Build a docker image from a Dockerfile.
# $1 is the name of the image to build
# $2 is the location of the "context" directory, with the Dockerfile at the root.
# $3 is the value to set the --pull flag for docker build; true by default
function zero::build::docker_build() {
  local -r image=$1
  local -r context_dir=$2
  local -r pull="${3:-true}"
  local -ra build_cmd=("${DOCKER[@]}" build -t "${image}" "--pull=${pull}" "${context_dir}")

  zero::log::status "Building Docker image ${image}"
  local docker_output
  docker_output=$("${build_cmd[@]}" 2>&1) || {
    cat <<EOF >&2
+++ Docker build command failed for ${image}

${docker_output}

To retry manually, run:

${build_cmd[*]}

EOF
    return 1
  }
}

function zero::build::ensure_data_container() {
  # If the data container exists AND exited successfully, we can use it.
  # Otherwise nuke it and start over.
  local ret=0
  local code=0

  code=$(docker inspect \
      -f '{{.State.ExitCode}}' \
      "${PROJ_DATA_CONTAINER_NAME}" 2>/dev/null) || ret=$?
  if [[ "${ret}" == 0 && "${code}" != 0 ]]; then
    zero::build::destroy_container "${PROJ_DATA_CONTAINER_NAME}"
    ret=1
  fi
  if [[ "${ret}" != 0 ]]; then
    zero::log::status "Creating data container ${PROJ_DATA_CONTAINER_NAME}"
    # We have to ensure the directory exists, or else the docker run will
    # create it as root.
    mkdir -p "${LOCAL_OUTPUT_GOPATH}"
    # We want this to run as root to be able to chown, so non-root users can
    # later use the result as a data container.  This run both creates the data
    # container and chowns the GOPATH.
    #
    # The data container creates volumes for all of the directories that store
    # intermediates for the Go build. This enables incremental builds across
    # Docker sessions. The *_cgo paths are re-compiled versions of the go std
    # libraries for true static building.
    local -ra docker_cmd=(
      "${DOCKER[@]}" run
      --volume "${REMOTE_ROOT}"   # white-out the whole output dir
      --volume /usr/local/go/pkg/linux_386_cgo
      --volume /usr/local/go/pkg/linux_amd64_cgo
      --volume /usr/local/go/pkg/linux_arm_cgo
      --volume /usr/local/go/pkg/linux_arm64_cgo
      --volume /usr/local/go/pkg/linux_ppc64le_cgo
      --volume /usr/local/go/pkg/darwin_amd64_cgo
      --volume /usr/local/go/pkg/darwin_386_cgo
      --volume /usr/local/go/pkg/windows_amd64_cgo
      --volume /usr/local/go/pkg/windows_386_cgo
      --name "${PROJ_DATA_CONTAINER_NAME}"
      --hostname "${HOSTNAME}"
      "${PROJ_BUILD_IMAGE}"
      chown -R "${USER_ID}":"${GROUP_ID}"
        "${REMOTE_ROOT}"
        /usr/local/go/pkg/
    )
    "${docker_cmd[@]}"
  fi
}

# Build all miner commands.
function zero::build::build_command() {
  zero::log::status "Running build command..."
  make -C "${ZROOT}" build.multiarch BINS="minerctl miner-apiserver miner-authz-server miner-pump miner-watcher"
}

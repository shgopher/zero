#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# This file creates release artifacts (tar files, container images) that are
# ready to distribute to install or distribute to end users.

###############################################################################
# Most of the ::release:: namespace functions have been moved to
# github.com/zero/release.  Have a look in that repo and specifically in
# lib/releaselib.sh for ::release::-related functionality.
###############################################################################

# Tencent cos configuration
readonly BUCKET="superproj-1254073058"
readonly REGION="ap-beijing"
readonly COS_RELEASE_DIR="zero-release"
readonly COSTOOL="coscmd"

# This is where the final release artifacts are created locally
readonly RELEASE_STAGE="${LOCAL_OUTPUT_ROOT}/release-stage"
readonly RELEASE_TARS="${LOCAL_OUTPUT_ROOT}/release-tars"
readonly RELEASE_IMAGES="${LOCAL_OUTPUT_ROOT}/release-images"

# zero github account info
readonly PROJ_GITHUB_ORG=superproj
readonly PROJ_GITHUB_REPO=zero

readonly ARTIFACT=zero.tar.gz
readonly CHECKSUM=${ARTIFACT}.sha1sum

PROJ_BUILD_CONFORMANCE=${PROJ_BUILD_CONFORMANCE:-y}
PROJ_BUILD_PULL_LATEST_IMAGES=${PROJ_BUILD_PULL_LATEST_IMAGES:-y}

# Validate a ci version
#
# Globals:
#   None
# Arguments:
#   version
# Returns:
#   If version is a valid ci version
# Sets:                    (e.g. for '1.2.3-alpha.4.56+abcdef12345678')
#   VERSION_MAJOR          (e.g. '1')
#   VERSION_MINOR          (e.g. '2')
#   VERSION_PATCH          (e.g. '3')
#   VERSION_PRERELEASE     (e.g. 'alpha')
#   VERSION_PRERELEASE_REV (e.g. '4')
#   VERSION_BUILD_INFO     (e.g. '.56+abcdef12345678')
#   VERSION_COMMITS        (e.g. '56')
function zero::release::parse_and_validate_ci_version() {
  # Accept things like "v1.2.3-alpha.4.56+abcdef12345678" or "v1.2.3-beta.4"
  local -r version_regex="^v(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)-([a-zA-Z0-9]+)\\.(0|[1-9][0-9]*)(\\.(0|[1-9][0-9]*)\\+[0-9a-f]{7,40})?$"
  local -r version="${1-}"
  [[ "${version}" =~ ${version_regex} ]] || {
    zero::log::error "Invalid ci version: '${version}', must match regex ${version_regex}"
    return 1
  }

  # The VERSION variables are used when this file is sourced, hence
  # the shellcheck SC2034 'appears unused' warning is to be ignored.

  # shellcheck disable=SC2034
  VERSION_MAJOR="${BASH_REMATCH[1]}"
  # shellcheck disable=SC2034
  VERSION_MINOR="${BASH_REMATCH[2]}"
  # shellcheck disable=SC2034
  VERSION_PATCH="${BASH_REMATCH[3]}"
  # shellcheck disable=SC2034
  VERSION_PRERELEASE="${BASH_REMATCH[4]}"
  # shellcheck disable=SC2034
  VERSION_PRERELEASE_REV="${BASH_REMATCH[5]}"
  # shellcheck disable=SC2034
  VERSION_BUILD_INFO="${BASH_REMATCH[6]}"
  # shellcheck disable=SC2034
  VERSION_COMMITS="${BASH_REMATCH[7]}"
}

# ---------------------------------------------------------------------------
# Build final release artifacts
function zero::release::clean_cruft() {
  # Clean out cruft
  find "${RELEASE_STAGE}" -name '*~' -exec rm {} \;
  find "${RELEASE_STAGE}" -name '#*#' -exec rm {} \;
  find "${RELEASE_STAGE}" -name '.DS*' -exec rm {} \;
}

function zero::release::package_tarballs() {
  # Clean out any old releases
  rm -rf "${RELEASE_STAGE}" "${RELEASE_TARS}" "${RELEASE_IMAGES}"
  mkdir -p "${RELEASE_TARS}"
  zero::release::package_src_tarball &
  zero::release::package_client_tarballs &
  zero::release::package_zero_manifests_tarball &
  zero::release::package_server_tarballs &
  zero::util::wait-for-jobs || { zero::log::error "previous tarball phase failed"; return 1; }

  zero::release::package_final_tarball & # _final depends on some of the previous phases
  zero::util::wait-for-jobs || { zero::log::error "previous tarball phase failed"; return 1; }
}

function zero::release::updload_tarballs() {
  zero::log::info "upload ${RELEASE_TARS}/* to cos bucket ${BUCKET}."
  for file in $(ls ${RELEASE_TARS}/*)
  do
    if [ "${COSTOOL}" == "coscli" ];then
      coscli cp "${file}" "cos://${BUCKET}/${COS_RELEASE_DIR}/${PROJ_GIT_VERSION}/"
      coscli cp "${file}" "cos://${BUCKET}/${COS_RELEASE_DIR}/latest/"
    else
      coscmd upload  "${file}" "${COS_RELEASE_DIR}/${PROJ_GIT_VERSION}/"
      coscmd upload  "${file}" "${COS_RELEASE_DIR}/latest/"
    fi
  done
}

# Package the source code we built, for compliance/licensing/audit/yadda.
function zero::release::package_src_tarball() {
  local -r src_tarball="${RELEASE_TARS}/zero-src.tar.gz"
  zero::log::status "Building tarball: src"
  if [[ "${PROJ_GIT_TREE_STATE-}" = 'clean' ]]; then
    git archive -o "${src_tarball}" HEAD
  else
    find "${ZROOT}" -mindepth 1 -maxdepth 1 \
      ! \( \
      \( -path "${ZROOT}"/_\* -o \
      -path "${ZROOT}"/.git\* -o \
      -path "${ZROOT}"/.gitignore\* -o \
      -path "${ZROOT}"/.gsemver.yaml\* -o \
      -path "${ZROOT}"/.config\* -o \
      -path "${ZROOT}"/.chglog\* -o \
      -path "${ZROOT}"/.gitlint -o \
      -path "${ZROOT}"/.golangci.yaml -o \
      -path "${ZROOT}"/.goreleaser.yml -o \
      -path "${ZROOT}"/.note.md -o \
      -path "${ZROOT}"/.todo.md \
      \) -prune \
      \) -print0 \
      | "${TAR}" czf "${src_tarball}" --transform "s|${ZROOT#/*}|zero|" --null -T -
  fi
}

# Package up all of the server binaries
function zero::release::package_server_tarballs() {
  # Find all of the built client binaries
  local long_platforms=("${LOCAL_OUTPUT_BINPATH}"/*/*)
  if [[ -n ${PROJ_BUILD_PLATFORMS-} ]]; then
    read -ra long_platforms <<< "${PROJ_BUILD_PLATFORMS}"
  fi

  for platform_long in "${long_platforms[@]}"; do
    local platform
    local platform_tag
    platform=${platform_long##${LOCAL_OUTPUT_BINPATH}/} # Strip LOCAL_OUTPUT_BINPATH
    platform_tag=${platform/\//-} # Replace a "/" for a "-"
    zero::log::status "Starting tarball: server $platform_tag"

    (
    local release_stage="${RELEASE_STAGE}/server/${platform_tag}/zero"
    rm -rf "${release_stage}"
    mkdir -p "${release_stage}/server/bin"

    local server_bins=("${PROJ_SERVER_BINARIES[@]}")

      # This fancy expression will expand to prepend a path
      # (${LOCAL_OUTPUT_BINPATH}/${platform}/) to every item in the
      # server_bins array.
      cp "${server_bins[@]/#/${LOCAL_OUTPUT_BINPATH}/${platform}/}" \
        "${release_stage}/server/bin/"

      zero::release::clean_cruft

      local package_name="${RELEASE_TARS}/zero-server-${platform_tag}.tar.gz"
      zero::release::create_tarball "${package_name}" "${release_stage}/.."
      ) &
    done

    zero::log::status "Waiting on tarballs"
    zero::util::wait-for-jobs || { zero::log::error "server tarball creation failed"; exit 1; }
  }

# Package up all of the cross compiled clients. Over time this should grow into
# a full SDK
function zero::release::package_client_tarballs() {
  # Find all of the built client binaries
  local long_platforms=("${LOCAL_OUTPUT_BINPATH}"/*/*)
  if [[ -n ${PROJ_BUILD_PLATFORMS-} ]]; then
    read -ra long_platforms <<< "${PROJ_BUILD_PLATFORMS}"
  fi

  for platform_long in "${long_platforms[@]}"; do
    local platform
    local platform_tag
    platform=${platform_long##${LOCAL_OUTPUT_BINPATH}/} # Strip LOCAL_OUTPUT_BINPATH
    platform_tag=${platform/\//-} # Replace a "/" for a "-"
    zero::log::status "Starting tarball: client $platform_tag"

    (
    local release_stage="${RELEASE_STAGE}/client/${platform_tag}/zero"
    rm -rf "${release_stage}"
    mkdir -p "${release_stage}/client/bin"

    local client_bins=("${PROJ_CLIENT_BINARIES[@]}")

      # This fancy expression will expand to prepend a path
      # (${LOCAL_OUTPUT_BINPATH}/${platform}/) to every item in the
      # client_bins array.
      cp "${client_bins[@]/#/${LOCAL_OUTPUT_BINPATH}/${platform}/}" \
        "${release_stage}/client/bin/"

      zero::release::clean_cruft

      local package_name="${RELEASE_TARS}/zero-client-${platform_tag}.tar.gz"
      zero::release::create_tarball "${package_name}" "${release_stage}/.."
    ) &
  done

  zero::log::status "Waiting on tarballs"
  zero::util::wait-for-jobs || { zero::log::error "client tarball creation failed"; exit 1; }
}

# Package up all of the server binaries in docker images
function zero::release::build_server_images() {
  # Clean out any old images
  rm -rf "${RELEASE_IMAGES}"
  local platform
  for platform in "${PROJ_SERVER_PLATFORMS[@]}"; do
    local platform_tag
    local arch
    platform_tag=${platform/\//-} # Replace a "/" for a "-"
    arch=$(basename "${platform}")
    zero::log::status "Building images: $platform_tag"

    local release_stage
    release_stage="${RELEASE_STAGE}/server/${platform_tag}/zero"
    rm -rf "${release_stage}"
    mkdir -p "${release_stage}/server/bin"

    # This fancy expression will expand to prepend a path
    # (${LOCAL_OUTPUT_BINPATH}/${platform}/) to every item in the
    # PROJ_SERVER_IMAGE_BINARIES array.
    cp "${PROJ_SERVER_IMAGE_BINARIES[@]/#/${LOCAL_OUTPUT_BINPATH}/${platform}/}" \
      "${release_stage}/server/bin/"

    zero::release::create_docker_images_for_server "${release_stage}/server/bin" "${arch}"
  done
}

function zero::release::md5() {
  if which md5 >/dev/null 2>&1; then
    md5 -q "$1"
  else
    md5sum "$1" | awk '{ print $1 }'
  fi
}

function zero::release::sha1() {
  if which sha1sum >/dev/null 2>&1; then
    sha1sum "$1" | awk '{ print $1 }'
  else
    shasum -a1 "$1" | awk '{ print $1 }'
  fi
}

function zero::release::build_conformance_image() {
  local -r arch="$1"
  local -r registry="$2"
  local -r version="$3"
  local -r save_dir="${4-}"
  zero::log::status "Building conformance image for arch: ${arch}"
  ARCH="${arch}" REGISTRY="${registry}" VERSION="${version}" \
    make -C cluster/images/conformance/ build >/dev/null

  local conformance_tag
  conformance_tag="${registry}/conformance-${arch}:${version}"
  if [[ -n "${save_dir}" ]]; then
    "${DOCKER[@]}" save "${conformance_tag}" > "${save_dir}/conformance-${arch}.tar"
  fi
  zero::log::status "Deleting conformance image ${conformance_tag}"
  "${DOCKER[@]}" rmi "${conformance_tag}" &>/dev/null || true
}

# This builds all the release docker images (One docker image per binary)
# Args:
#  $1 - binary_dir, the directory to save the tared images to.
#  $2 - arch, architecture for which we are building docker images.
function zero::release::create_docker_images_for_server() {
  # Create a sub-shell so that we don't pollute the outer environment
  (
    local binary_dir
    local arch
    local binaries
    local images_dir
    binary_dir="$1"
    arch="$2"
    binaries=$(zero::build::get_docker_wrapped_binaries "${arch}")
    images_dir="${RELEASE_IMAGES}/${arch}"
    mkdir -p "${images_dir}"

    # k8s.gcr.io is the constant tag in the docker archives, this is also the default for config scripts in GKE.
    # We can use PROJ_DOCKER_REGISTRY to include and extra registry in the docker archive.
    # If we use PROJ_DOCKER_REGISTRY="k8s.gcr.io", then the extra tag (same) is ignored, see release_docker_image_tag below.
    local -r docker_registry="k8s.gcr.io"
    # Docker tags cannot contain '+'
    local docker_tag="${PROJ_GIT_VERSION/+/_}"
    if [[ -z "${docker_tag}" ]]; then
      zero::log::error "git version information missing; cannot create Docker tag"
      return 1
    fi

    # provide `--pull` argument to `docker build` if `PROJ_BUILD_PULL_LATEST_IMAGES`
    # is set to y or Y; otherwise try to build the image without forcefully
    # pulling the latest base image.
    local docker_build_opts
    docker_build_opts=
    if [[ "${PROJ_BUILD_PULL_LATEST_IMAGES}" =~ [yY] ]]; then
        docker_build_opts='--pull'
    fi

    for wrappable in $binaries; do

      local binary_name=${wrappable%%,*}
      local base_image=${wrappable##*,}
      local binary_file_path="${binary_dir}/${binary_name}"
      local docker_build_path="${binary_file_path}.dockerbuild"
      local docker_file_path="${docker_build_path}/Dockerfile"
      local docker_image_tag="${docker_registry}/${binary_name}-${arch}:${docker_tag}"

      zero::log::status "Starting docker build for image: ${binary_name}-${arch}"
      (
        rm -rf "${docker_build_path}"
        mkdir -p "${docker_build_path}"
        ln "${binary_file_path}" "${docker_build_path}/${binary_name}"
        ln "${ZROOT}/build/nsswitch.conf" "${docker_build_path}/nsswitch.conf"
        chmod 0644 "${docker_build_path}/nsswitch.conf"
        cat <<EOF > "${docker_file_path}"
FROM ${base_image}
COPY ${binary_name} /usr/local/bin/${binary_name}
EOF
        # ensure /etc/nsswitch.conf exists so go's resolver respects /etc/hosts
        if [[ "${base_image}" =~ busybox ]]; then
          echo "COPY nsswitch.conf /etc/" >> "${docker_file_path}"
        fi

        "${DOCKER[@]}" build ${docker_build_opts:+"${docker_build_opts}"} -q -t "${docker_image_tag}" "${docker_build_path}" >/dev/null
        # If we are building an official/alpha/beta release we want to keep
        # docker images and tag them appropriately.
        local -r release_docker_image_tag="${PROJ_DOCKER_REGISTRY-$docker_registry}/${binary_name}-${arch}:${PROJ_DOCKER_IMAGE_TAG-$docker_tag}"
        if [[ "${release_docker_image_tag}" != "${docker_image_tag}" ]]; then
          zero::log::status "Tagging docker image ${docker_image_tag} as ${release_docker_image_tag}"
          "${DOCKER[@]}" rmi "${release_docker_image_tag}" 2>/dev/null || true
          "${DOCKER[@]}" tag "${docker_image_tag}" "${release_docker_image_tag}" 2>/dev/null
        fi
        "${DOCKER[@]}" save -o "${binary_file_path}.tar" "${docker_image_tag}" "${release_docker_image_tag}"
        echo "${docker_tag}" > "${binary_file_path}.docker_tag"
        rm -rf "${docker_build_path}"
        ln "${binary_file_path}.tar" "${images_dir}/"

        zero::log::status "Deleting docker image ${docker_image_tag}"
        "${DOCKER[@]}" rmi "${docker_image_tag}" &>/dev/null || true
      ) &
    done

    if [[ "${PROJ_BUILD_CONFORMANCE}" =~ [yY] ]]; then
      zero::release::build_conformance_image "${arch}" "${docker_registry}" \
        "${docker_tag}" "${images_dir}" &
    fi

    zero::util::wait-for-jobs || { zero::log::error "previous Docker build failed"; return 1; }
    zero::log::status "Docker builds done"
  )

}

# This will pack zero-system manifests files for distros such as COS.
function zero::release::package_zero_manifests_tarball() {
  zero::log::status "Building tarball: manifests"

  local src_dir="${ZROOT}/deployments"

  local release_stage="${RELEASE_STAGE}/manifests/zero"
  rm -rf "${release_stage}"

  local dst_dir="${release_stage}"
  mkdir -p "${dst_dir}"
  cp -r ${src_dir}/* "${dst_dir}"
  #cp "${src_dir}/zero-apiserver.yaml" "${dst_dir}"
  #cp "${src_dir}/zero-authz-server.yaml" "${dst_dir}"
  #cp "${src_dir}/zero-pump.yaml" "${dst_dir}"
  #cp "${src_dir}/zero-watcher.yaml" "${dst_dir}"
  #cp "${ZROOT}/cluster/gce/gci/health-monitor.sh" "${dst_dir}/health-monitor.sh"

  zero::release::clean_cruft

  local package_name="${RELEASE_TARS}/zero-manifests.tar.gz"
  zero::release::create_tarball "${package_name}" "${release_stage}/.."
}

# This is all the platform-independent stuff you need to run/install zero.
# Arch-specific binaries will need to be downloaded separately (possibly by
# using the bundled cluster/get-zero-binaries.sh script).
# Included in this tarball:
#   - Cluster spin up/down scripts and configs for various cloud providers
#   - Tarballs for manifest configs that are ready to be uploaded
#   - Examples (which may or may not still work)
#   - The remnants of the docs/ directory
function zero::release::package_final_tarball() {
  zero::log::status "Building tarball: final"

  # This isn't a "full" tarball anymore, but the release lib still expects
  # artifacts under "full/zero/"
  local release_stage="${RELEASE_STAGE}/full/zero"
  rm -rf "${release_stage}"
  mkdir -p "${release_stage}"

  mkdir -p "${release_stage}/client"
  cat <<EOF > "${release_stage}/client/README"
Client binaries are no longer included in the zero final tarball.

Run release/get-zero-binaries.sh to download client and server binaries.
EOF

  # We want everything in /scripts.
  mkdir -p "${release_stage}/release"
  cp -R "${ZROOT}/scripts/release" "${release_stage}/"
  cat <<EOF > "${release_stage}/release/get-zero-binaries.sh"
#!/usr/bin/env bash

# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# This file download zero client and server binaries from tencent cos bucket.

os=linux arch=amd64 version=${PROJ_GIT_VERSION} && wget https://${BUCKET}.cos.${REGION}.myqcloud.com/${COS_RELEASE_DIR}/\$version/{zero-client-\$os-\$arch.tar.gz,zero-server-\$os-\$arch.tar.gz}
EOF
  chmod +x ${release_stage}/release/get-zero-binaries.sh

  mkdir -p "${release_stage}/server"
  cp "${RELEASE_TARS}/zero-manifests.tar.gz" "${release_stage}/server/"
  cat <<EOF > "${release_stage}/server/README"
Server binary tarballs are no longer included in the zero final tarball.

Run release/get-zero-binaries.sh to download client and server binaries.
EOF

  # Include hack/lib as a dependency for the cluster/ scripts
  #mkdir -p "${release_stage}/hack"
  #cp -R "${ZROOT}/hack/lib" "${release_stage}/hack/"

  cp -R ${ZROOT}/{docs,configs,scripts,deployments,init,README.md,LICENSE} "${release_stage}/"

  echo "${PROJ_GIT_VERSION}" > "${release_stage}/version"

  zero::release::clean_cruft

  local package_name="${RELEASE_TARS}/${ARTIFACT}"
  zero::release::create_tarball "${package_name}" "${release_stage}/.."
}

# Build a release tarball.  $1 is the output tar name.  $2 is the base directory
# of the files to be packaged.  This assumes that ${2}/zerois what is
# being packaged.
function zero::release::create_tarball() {
  zero::build::ensure_tar

  local tarfile=$1
  local stagingdir=$2

  "${TAR}" czf "${tarfile}" -C "${stagingdir}" zero --owner=0 --group=0
}

function zero::release::install_github_release(){
  GO111MODULE=on go install github.com/github-release/github-release@latest
}

# Require the following tools:
# - github-release
# - gsemver
# - git-chglog
# - coscmd or coscli
function zero::release::verify_prereqs(){
  if [ -z "$(which github-release 2>/dev/null)" ]; then
    zero::log::info "'github-release' tool not installed, try to install it."

    if ! zero::release::install_github_release; then
      zero::log::error "failed to install 'github-release'"
      return 1
    fi
  fi

  if [ -z "$(which git-chglog 2>/dev/null)" ]; then
    zero::log::info "'git-chglog' tool not installed, try to install it."

    if ! go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest &>/dev/null; then
      zero::log::error "failed to install 'git-chglog'"
      return 1
    fi
  fi

  if [ -z "$(which gsemver 2>/dev/null)" ]; then
    zero::log::info "'gsemver' tool not installed, try to install it."

    if ! go install github.com/arnaud-deprez/gsemver@latest &>/dev/null; then
      zero::log::error "failed to install 'gsemver'"
      return 1
    fi
  fi


  if [ -z "$(which ${COSTOOL} 2>/dev/null)" ]; then
    zero::log::info "${COSTOOL} tool not installed, try to install it."

    if ! make -C "${ZROOT}" tools.install.${COSTOOL}; then
      zero::log::error "failed to install ${COSTOOL}"
      return 1
    fi
  fi

  if [ -z "${TENCENT_SECRET_ID}" -o -z "${TENCENT_SECRET_KEY}" ];then
      zero::log::error "can not find env: TENCENT_SECRET_ID and TENCENT_SECRET_KEY"
      return 1
  fi

  if [ "${COSTOOL}" == "coscli" ];then
    if [ ! -f "${HOME}/.cos.yaml" ];then
      cat << EOF > "${HOME}/.cos.yaml"
cos:
  base:
    secretid: ${TENCENT_SECRET_ID}
    secretkey: ${TENCENT_SECRET_KEY}
    sessiontoken: ""
  buckets:
  - name: ${BUCKET}
    alias: ${BUCKET}
    region: ${REGION}
EOF
    fi
  else
    if [ ! -f "${HOME}/.cos.conf" ];then
      cat << EOF > "${HOME}/.cos.conf"
[common]
secret_id = ${TENCENT_SECRET_ID}
secret_key = ${TENCENT_SECRET_KEY}
bucket = ${BUCKET}
region =${REGION}
max_thread = 5
part_size = 1
schema = https
EOF
    fi
  fi
}

# Create a github release with specified tarballs.
# NOTICE: Must export 'GITHUB_TOKEN' env in the shell, details:
# https://github.com/github-release/github-release
function zero::release::github_release() {
  # create a github release
  zero::log::info "create a new github release with tag ${PROJ_GIT_VERSION}"
  github-release release \
    --user ${PROJ_GITHUB_ORG} \
    --repo ${PROJ_GITHUB_REPO} \
    --tag ${PROJ_GIT_VERSION} \
    --description "" \
    --pre-release

  # update zero tarballs
  zero::log::info "upload ${ARTIFACT} to release ${PROJ_GIT_VERSION}"
  github-release upload \
    --user ${PROJ_GITHUB_ORG} \
    --repo ${PROJ_GITHUB_REPO} \
    --tag ${PROJ_GIT_VERSION} \
    --name ${ARTIFACT} \
    --file ${RELEASE_TARS}/${ARTIFACT}

  zero::log::info "upload zero-src.tar.gz to release ${PROJ_GIT_VERSION}"
  github-release upload \
    --user ${PROJ_GITHUB_ORG} \
    --repo ${PROJ_GITHUB_REPO} \
    --tag ${PROJ_GIT_VERSION} \
    --name "zero-src.tar.gz" \
    --file ${RELEASE_TARS}/zero-src.tar.gz
}

function zero::release::generate_changelog() {
  zero::log::info "generate CHANGELOG-${PROJ_GIT_VERSION#v}.md and commit it"

  git-chglog ${PROJ_GIT_VERSION} > ${ZROOT}/CHANGELOG/CHANGELOG-${PROJ_GIT_VERSION#v}.md

  set +o errexit
  git add ${ZROOT}/CHANGELOG/CHANGELOG-${PROJ_GIT_VERSION#v}.md
  git commit -a -m "docs(changelog): add CHANGELOG-${PROJ_GIT_VERSION#v}.md"
  git push -f origin master # 最后将 CHANGELOG 也 push 上去
}


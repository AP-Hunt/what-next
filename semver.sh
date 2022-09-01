#!/usr/bin/env sh

incr_major(){
  VER="${1}"
  echo $VER | awk -F. '{print ($1+1 ".0.0")}'
}

incr_minor() {
    VER="${1}"
    echo $VER | awk -F. '{print ($1 "." $2+1 "." $3)}'
}

incr_patch() {
      VER="${1}"
      echo $VER | awk -v pre_release="$(get_pre_release "$VER") "-F. '{print ($1 "." $2 "." $3+1 pre_release)}'
}

get_pre_release() {
  VER="${1}"
  echo $VER | awk -F- 'NF==1 { print "" } NF>1 { print "-" $2}'
}

set_pre_release() {
  VER="${1}"
  NEW_VALUE="${2}"
  echo $VER | awk -v new_value="${NEW_VALUE}" -F- 'NF==1 { print $1 "-" new_value } NF>1 { print $1 "-" new_value}'
}

usage() {
  echo "Usage: $0 -v version {-M|-m|-p} [-r  name]"
  echo "  -v        The version to operate on"
  echo "  -M        Increment the major version"
  echo "  -m        Increment the minor version"
  echo "  -p        Increment the patch version"
  echo "  -r name   Set the pre-release version"
}

VERSION_SET=false
VERSION=""
MAJOR=false
MINOR=false
PATCH=false
PRE_RELEASE=false
PRE_RELEASE_VALUE=""

while getopts "Mmpr:v:" options; do
  case "${options}" in
    v)
      VERSION_SET=true
      VERSION="${OPTARG}"
      ;;
    M)
      if ${MINOR}; then
        echo "Cannot bump major and minor at the same time" >&2
        exit 1
      fi

      if ${PATCH}; then
        echo "Cannot bump major and patch at the same time" >&2
        exit 1
      fi

      MAJOR=true
      ;;
    m)
      if ${MAJOR}; then
        echo "Cannot bump minor and major at the same time" >&2
        exit 1
      fi

      if ${PATCH}; then
        echo "Cannot bump minor and patch at the same time" >&2
        exit 1
      fi

      MINOR=true
      ;;

    p)
      if ${MAJOR}; then
        echo "Cannot bump patch and major at the same time" >&2
        exit 1
      fi

      if ${MINOR}; then
        echo "Cannot bump patch and minor at the same time" >&2
        exit 1
      fi

      PATCH=true
      ;;

    r)
      PRE_RELEASE=true
      PRE_RELEASE_VALUE="${OPTARG}"
      ;;

    :)
      echo "-${OPTARG} requires an argument" >&2
      exit 1
      ;;

    *)
      usage
      exit 1
  esac
done

if [ "$VERSION_SET" = false ]; then
  echo "Version must be set with -v flag" >&2
  exit 1
fi

if [ "$MAJOR" = false ] && [ "$MINOR" = false ] && [ "$PATCH" = false ] && [ "$PRE_RELEASE" = false ]; then
  echo "Must bump one of major, minor, and patch, or set the pre-release version" >&2
  exit 1
fi

if [ "$MAJOR" = true ]; then
  NEW_VERSION=$(incr_major "$VERSION")

  if [ "$PRE_RELEASE" = true ]; then
    NEW_VERSION=$(set_pre_release "$NEW_VERSION" "$PRE_RELEASE_VALUE")
  fi
  printf "%s" "$NEW_VERSION"
elif [ "$MINOR" = true ]; then
    NEW_VERSION=$(incr_minor "$VERSION")

    if [ "$PRE_RELEASE" = true ]; then
      NEW_VERSION=$(set_pre_release "$NEW_VERSION" "$PRE_RELEASE_VALUE")
    fi
    printf "%s" "$NEW_VERSION"
elif [ "$PATCH" = true ]; then
    NEW_VERSION=$(incr_patch "$VERSION")

    if [ "$PRE_RELEASE" = true ]; then
      NEW_VERSION=$(set_pre_release "$NEW_VERSION" "$PRE_RELEASE_VALUE")
    fi
    printf "%s" "$NEW_VERSION"
elif [ "$PRE_RELEASE" = true ]; then
    NEW_VERSION=$(set_pre_release "$VERSION" "$PRE_RELEASE_VALUE")
    printf "%s" "$NEW_VERSION"
fi
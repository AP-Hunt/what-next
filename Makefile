## Make settings
.DEFAULT_GOAL := $OUT_PATH

## Vars
BIN_NAME := $(shell basename $$(pwd))
OUT_PATH := bin/$(BIN_NAME)
GO_SRC := $(shell find . -type f -name "*.go" -not -path "./vendor/*")
VENDOR_DIRS := $(shell find vendor/ -mindepth 1 -maxdepth 3 -type d 2>/dev/null | sort | uniq)
VERSION_FILE := ./version
VERSION := $(shell cat "${VERSION_FILE}")

## Build targets
$OUT_PATH: $(GO_SRC) ./vendor/ $(VENDOR_DIRS)
	@echo "Compiling to ${OUT_PATH}"
	go build \
		-o "${OUT_PATH}" \
		-ldflags="-X 'main.Version=${VERSION}'" \
		.

go.mod:
	go mod init

go.sum:
	go mod tidy

./vendor/: go.mod go.sum
	go mod vendor

## Test targets
.PHONY: test
test: ./vendor/
	go test ./...

## Versioning targets
.PHONY: version
version:
	@echo "Use the bump_major, bump_minor, bump_patch, and set_pre_release targets to manage the project version"
	@echo "To set the pre-release version in each, set the P variable e.g."
	@echo "    make bump_minor P=beta-1"
.PHONY: bump_major
bump_major:
	@EXTRA_ARGS=""
	@if [ ! -z "${P:-}" ]; then EXTRA_ARGS="-r ${P}"; fi
	@./semver.sh -v "${VERSION}" -M ${EXTRA_ARGS} > "${VERSION_FILE}"

.PHONY: bump_minor
bump_minor:
	@EXTRA_ARGS=""
	@if [ ! -z "${P:-}" ]; then EXTRA_ARGS="-r ${P}"; fi
	@./semver.sh -v "${VERSION}" -m ${EXTRA_ARGS} > "${VERSION_FILE}"

.PHONY: bump_patch
bump_patch:
	@EXTRA_ARGS=""
	@if [ ! -z "${P:-}" ]; then EXTRA_ARGS="-r ${P}"; fi
	@./semver.sh -v "${VERSION}" -p ${EXTRA_ARGS} > "${VERSION_FILE}"

.PHONY: set_pre_release
set_pre_release:
	@if [ -z "$P" ]; then \
      echo "Set the value with P=value"; \
    else \
	  ./semver.sh -v "${VERSION}" -r "${P}" > "${VERSION_FILE}"; \
  	fi

push_version:

# Release targets
.PHONY: release
release:
	@echo "Pushing version ${VERSION}"
	git add version
	git commit -m "Bump to version ${VERSION}"
	git tag "${VERSION}"
	@echo "\n\nVersion bumped in commit $$(git rev-parse HEAD)"
	@echo "Run the following to push the new version"
	@echo "\t git push origin main"
	@echo "\t git push origin ${VERSION}"
	@echo "\nThe GitHub Actions workflow will then produce a new release"
	@echo "on GitHub and you can edit the release notes from there."

.PHONY: dist
dist: ./vendor/
	if [ -d release ]; then \
  		rm -rf release; \
  	fi; \
	mkdir release; \
	for os in linux windows darwin; do \
  		for arch in amd64 arm64; do \
  		  	GOOS="$${OS}" GOARCH="$${arch}" \
			go build \
				-o "release/${BIN_NAME}-${VERSION}-$${os}-$${arch}" \
				-ldflags="-X 'main.Version=${VERSION}'" \
				.; \
  		done \
  	done; \
  	tar czf "${BIN_NAME}-${VERSION}.tar.gz" -C release/ .; \
  	mv "./${BIN_NAME}-${VERSION}.tar.gz" "./release/${BIN_NAME}-${VERSION}.tar.gz"
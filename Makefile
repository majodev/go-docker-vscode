### -----------------------
# --- Building
### -----------------------

# first is default target when running "make" without args
build: ##- Default 'make' target: sql, swagger, go-generate, go-format, go-build and lint.
	@$(MAKE) go-format
	@$(MAKE) go-build
	@$(MAKE) lint

# useful to ensure that everything gets resetuped from scratch
all: init ##- Runs all of our common make targets: clean, init, build and test.
	@$(MAKE) build
	@$(MAKE) test

info: ##- Prints info about go.mod updates and current go version.
	@$(MAKE) get-go-outdated-modules
	@go version

lint: ##- (opt) Runs golangci-lint.
	golangci-lint run --timeout 5m

go-format: ##- (opt) Runs go format.
	go fmt ./...

go-build: ##- (opt) Runs go build.
	go build -o bin/app

# https://github.com/gotestyourself/gotestsum#format 
# w/o cache https://github.com/golang/go/issues/24573 - see "go help testflag"
# note that these tests should not run verbose by default (e.g. use your IDE for this)
# TODO: add test shuffling/seeding when landed in go v1.15 (https://github.com/golang/go/issues/28592)
# tests by pkgname
test: ##- Run tests, output by package, print coverage.
	@$(MAKE) go-test-by-pkg
	@$(MAKE) go-test-print-coverage

# tests by testname
test-by-name: ##- Run tests, output by testname, print coverage.
	@$(MAKE) go-test-by-name
	@$(MAKE) go-test-print-coverage

# note that we explicitly don't want to use a -coverpkg=./... option, per pkg coverage take precedence
go-test-by-pkg: ##- (opt) Run tests, output by package.
	gotestsum --format pkgname-and-test-fails --jsonfile /tmp/test.log -- -race -cover -count=1 -coverprofile=/tmp/coverage.out ./...

go-test-by-name: ##- (opt) Run tests, output by testname.
	gotestsum --format testname --jsonfile /tmp/test.log -- -race -cover -count=1 -coverprofile=/tmp/coverage.out ./...

go-test-print-coverage: ##- (opt) Print overall test coverage (must be done after running tests).
	@printf "coverage "
	@go tool cover -func=/tmp/coverage.out | tail -n 1 | awk '{$$1=$$1;print}'

go-test-print-slowest: ##- Print slowest running tests (must be done after running tests).
	gotestsum tool slowest --jsonfile /tmp/test.log --threshold 2s

# TODO: switch to "-m direct" after go 1.17 hits: https://github.com/golang/go/issues/40364
get-go-outdated-modules: ##- (opt) Prints outdated (direct) go modules (from go.mod). 
	@((go list -u -m -f '{{if and .Update (not .Indirect)}}{{.}}{{end}}' all) 2>/dev/null | grep " ") || echo "go modules are up-to-date."

watch-tests: ##- Watches *.go files and runs package tests on modifications.
	gotestsum --format testname --watch -- -race -count=1

### -----------------------
# --- Initializing
### -----------------------

init: ##- Runs make modules, tools and tidy.
	@$(MAKE) modules
	@$(MAKE) tools
	@$(MAKE) tidy

# cache go modules (locally into .pkg)
modules: ##- (opt) Cache packages as specified in go.mod.
	go mod download

# https://marcofranssen.nl/manage-go-tools-via-go-modules/
tools: ##- (opt) Install packages as specified in tools.go.
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -P $$(nproc) -tI % go install %

tidy: ##- (opt) Tidy our go.sum file.
	go mod tidy

### -----------------------
# --- Binary checks
### -----------------------

# Got license issues with some dependencies? Provide a custom lichen --config
# see https://github.com/uw-labs/lichen#config 
get-licenses: ##- Prints licenses of embedded modules in the compiled bin/app.
	lichen bin/app

get-embedded-modules: ##- Prints embedded modules in the compiled bin/app.
	go version -m -v bin/app

get-embedded-modules-count: ##- (opt) Prints count of embedded modules in the compiled bin/app.
	go version -m -v bin/app | grep $$'\tdep' | wc -l

### -----------------------
# --- Helpers
### -----------------------

# https://gist.github.com/prwhite/8168133 - based on comment from @m000
help: ##- Show common make targets.
	@echo "usage: make <target>"
	@echo "note: use 'make help-all' to see all make targets."
	@echo ""
	@sed -e '/#\{2\}-/!d; s/\\$$//; s/:[^#\t]*/@/; s/#\{2\}- *//' $(MAKEFILE_LIST) | grep --invert "(opt)" | sort | column -t -s '@'

help-all: ##- Show all make targets.
	@echo "usage: make <target>"
	@echo "note: make targets flagged with '(opt)' are part of a main target."
	@echo ""
	@sed -e '/#\{2\}-/!d; s/\\$$//; s/:[^#\t]*/@/; s/#\{2\}- *//' $(MAKEFILE_LIST) | sort | column -t -s '@'


### -----------------------
# --- Special targets
### -----------------------

# https://www.gnu.org/software/make/manual/html_node/Special-Targets.html
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
# ignore matching file/make rule combinations in working-dir
.PHONY: test help

# https://unix.stackexchange.com/questions/153763/dont-stop-makeing-if-a-command-fails-but-check-exit-status
# https://www.gnu.org/software/make/manual/html_node/One-Shell.html
# required to ensure make fails if one recipe fails (even on parallel jobs) and on pipefails
.ONESHELL:

# # normal POSIX bash shell mode
# SHELL = /bin/bash
# .SHELLFLAGS = -cEeuo pipefail

# wrapped make time tracing shell, use it via MAKE_TRACE_TIME=true make <target>
SHELL = /app/rksh
.SHELLFLAGS = $@
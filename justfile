# gnmatcher project justfile

# Variables
app := "gnmatcher"
org := "github.com/gnames/"
build_dir := "out/"
release_dir := build_dir + "releases/"
test_opts := "-parallel=1 -shuffle=on -count=1 -race -coverprofile=coverage.txt -covermode=atomic"

no_c := "CGO_ENABLED=0"
x86 := "GOARCH=amd64"
arm := "GOARCH=arm64"
linux := "GOOS=linux"
mac := "GOOS=darwin"
win := "GOOS=windows"

# Colors
green := `tput -Txterm setaf 2`
yellow := `tput -Txterm setaf 3`
white := `tput -Txterm setaf 7`
cyan := `tput -Txterm setaf 6`
reset := `tput -Txterm sgr0`

# Get version from git
version_full := `git describe --tags`
version := `git describe --tags --abbrev=0`
date := `date -u '+%Y-%m-%d_%H:%M:%S%Z'`

# LD flags with version and build date
flags_ld := "-ldflags '-X " + org + app + \
    "/pkg.Build=" + date + " -X " + org + app + \
    "/pkg.Version=" + version_full + "'"
flags_rel := "-trimpath -ldflags '-s -w -X " + org + app + \
    "/pkg.Build=" + date + "'"

# Default recipe (just install)
default: install

# Show this help
help:
    @echo ''
    @echo 'Usage:'
    @echo '  {{yellow}}just{{reset}} {{green}}<target>{{reset}}'
    @echo ''
    @echo 'Targets:'
    @just --list --unsorted

# Display current version
version:
    @echo {{version_full}}

# Clean up and sync dependencies
tidy:
    @go mod tidy
    @go mod verify

# Install tools
tools: tidy
    @go install tool
    @echo "✅ tools of the project are installed"

# Build binary
build:
    {{no_c}} go build -o {{build_dir}}{{app}} {{flags_ld}}
    @echo "✅ {{app}} built to {{build_dir}}{{app}}"

# Build binary without debug info and with hardcoded version
buildrel:
    {{no_c}} go build -o {{build_dir}}{{app}} {{flags_rel}}
    @echo "✅ {{app}} release binary built to {{build_dir}}{{app}}"

# Build and install binary
install:
    {{no_c}} go install {{flags_ld}}
    @echo "✅ {{app}} installed to ~/go/bin/{{app}}"

# Build docker image
docker: buildrel
    @mkdir -p {{build_dir}}
    {{no_c}} {{linux}} {{x86}} go build {{flags_rel}} -o {{build_dir}}{{app}}
    docker buildx build -t gnames/{{app}}:latest -t gnames/{{app}}:{{version_full}} .
    @echo "✅ docker image gnames/{{app}}:{{version_full}} built"

# Push docker image to dockerhub
dockerhub: docker
    docker push gnames/{{app}}
    docker push gnames/{{app}}:{{version_full}}
    @echo "✅ docker image pushed to dockerhub"

# Build and package binaries for a release
release: dockerhub
    @echo "Building releases for Linux, Mac, Windows (Intel and Arm)"
    @mkdir -p {{release_dir}}

    {{no_c}} {{linux}} {{x86}} go build {{flags_rel}} -o {{release_dir}}{{app}}
    tar zcvf {{release_dir}}{{app}}-{{version}}-linux-amd64.tar.gz {{release_dir}}{{app}}
    rm {{release_dir}}{{app}}

    {{no_c}} {{linux}} {{arm}} go build {{flags_rel}} -o {{release_dir}}{{app}}
    tar zcvf {{release_dir}}{{app}}-{{version}}-linux-arm64.tar.gz {{release_dir}}{{app}}
    rm {{release_dir}}{{app}}

    {{no_c}} {{mac}} {{x86}} go build {{flags_rel}} -o {{release_dir}}{{app}}
    tar zcvf {{release_dir}}{{app}}-{{version}}-mac-amd64.tar.gz {{release_dir}}{{app}}
    rm {{release_dir}}{{app}}

    {{no_c}} {{mac}} {{arm}} go build {{flags_rel}} -o {{release_dir}}{{app}}
    tar zcvf {{release_dir}}{{app}}-{{version}}-mac-arm64.tar.gz {{release_dir}}{{app}}
    rm {{release_dir}}{{app}}

    {{no_c}} {{win}} {{x86}} go build {{flags_rel}} -o {{release_dir}}{{app}}.exe
    cd {{release_dir}} && zip -9 {{app}}-{{version}}-win-amd64.zip {{app}}.exe
    rm {{release_dir}}{{app}}.exe

    {{no_c}} {{win}} {{arm}} go build {{flags_rel}} -o {{release_dir}}{{app}}.exe
    cd {{release_dir}} && zip -9 {{app}}-{{version}}-win-arm64.zip {{app}}.exe
    rm {{release_dir}}{{app}}.exe

    @echo "✅ Release binaries created in {{release_dir}}"

# Clean all the files and binaries generated
clean:
    @rm -rf ./{{build_dir}}

# Lint the code
lint:
    golangci-lint run

# Run the tests of the project
test: install
    go test {{test_opts}} ./...

# Run the tests of the project and export the coverage
coverage:
    @go test -p 1 -cover -covermode=count -coverprofile=profile.cov ./...
    @go tool cover -func profile.cov

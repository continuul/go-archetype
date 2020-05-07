NAME=go-archetype
VERSION:=$(shell semtag final -s minor -o)

# Git version information
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
GIT_DESCRIBE ?= $(shell git describe --tags --always --match "v*")
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_IMPORT ?= github.com/continuul/go-template/lib/version
GOLDFLAGS = -X $(GIT_IMPORT).GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) -X $(GIT_IMPORT).GitDescribe=$(GIT_DESCRIBE)
export GOLDFLAGS

# Go tools
GOTOOLS=\
	github.com/hashicorp/go-bindata/... \
	golang.org/x/lint/golint

CGO_ENABLED ?= 0

# cross-compilation variables
ARCHITECTURES ?= amd64
PLATFORMS ?= darwin linux windows

#:help: help        | Displays the GNU makefile help
.PHONY: help
help: ; @sed -n 's/^#:help://p' Makefile

#:help: build       | Build the executables, including all cross-compilation targets
.PHONY: build
build: generate
	@echo "==> Building go-archetype for all GOARCH/GOOS..."
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -ldflags '$(GOLDFLAGS)' -v -o release/$(NAME)-$(GOOS)-$(GOARCH) )))

#:help: check       | Runs the pre-commit hooks to check the files
.PHONY: check
check:
	pre-commit run --all-files

#:help: changelog   | Build the changelog
.PHONY: changelog
changelog:
	git-chglog -o CHANGELOG.md --next-tag $(VERSION)
	git add CHANGELOG.md && git commit -m "Updated CHANGELOG"

#:help: changever   | Change the product version to the next consecutive version number.
.PHONY: changever
changever:
	find bin -type f -name ws -exec sed -i "" "s/VERSION=.*/VERSION=\"$(VERSION)\"/g" {} \;
	git add bin/ws && git commit -m "Updated VERSION"

#:help: generate    | Generate archetype static assets
.PHONY: generate
generate:
	@go-bindata --prefix archetypes -o lib/archetype/archetype.go -pkg archetype ./archetypes/...

#:help: release     | Release the product, setting the tag and pushing.
.PHONY: release
release:
	semtag final -s minor
	git push --follow-tags

#:help: tools       | Installs go tools
.PHONY: tools
tools:
	@mkdir -p .gotools
	@cd .gotools && if [[ ! -f go.mod ]]; then \
		go mod init terraform-archetype-tools ; \
	fi
	cd .gotools && go get -v $(GOTOOLS)

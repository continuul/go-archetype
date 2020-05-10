NAME=go-archetype
VERSION:=$(shell semtag getfinal)

# Git version information
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
GIT_DESCRIBE ?= $(shell git describe --tags --always --match "v*")
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
GIT_IMPORT="github.com/continuul/go-archetype"
GOLDFLAGS="-X $(GIT_IMPORT)/lib/version.GitCommit=$(GIT_COMMIT)$(GIT_DIRTY) \
		   -X $(GIT_IMPORT)/lib/version.GitDescribe=$(GIT_DESCRIBE)"
export GOLDFLAGS

# Go tools
GOTAGS ?=
GOTOOLS=\
	github.com/go-bindata/go-bindata/... \
	golang.org/x/lint/golint

CGO_ENABLED ?= 0

# cross-compilation variables
ARCHITECTURES ?= amd64
PLATFORMS ?= darwin linux windows

#:help: help        | Displays the GNU makefile help
.PHONY: help
help: ; @sed -n 's/^#:help://p' Makefile

#:help: assets      | Generate archetype static assets
.PHONY: assets
assets:
	@go-bindata --prefix archetypes -o generated/archetype/archetype.go -pkg archetype ./archetypes/...

#:help: build       | Build the executables, including all cross-compilation targets
.PHONY: build
build: assets
	@echo "==> Building go-archetype for all GOARCH/GOOS..."
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -ldflags $(GOLDFLAGS) -v -o release/$(NAME)-$(GOOS)-$(GOARCH) )))

CHECK_FILES:=$(shell find . -type f | grep -v generated | grep -v .idea | grep -v .git)

#:help: check       | Runs the pre-commit hooks to check the files
.PHONY: check
check:
	@pre-commit run --files $(CHECK_FILES)

#:help: changelog   | Build the changelog
.PHONY: changelog
changelog:
	git-chglog -o CHANGELOG.md --next-tag $(VERSION)
	git add CHANGELOG.md && git commit -m "Updated CHANGELOG"

#:help: changever   | Change the product version to the next consecutive version number
.PHONY: changever
changever:
	find lib/version -type f -name version.go -exec sed -i "" "s/Version = .*/Version = \"$(VERSION)\"/g" {} \;
	git add lib/version/version.go && git commit -m "Updated VERSION"

#:help: clean       | Clean the build artifacts
.PHONY: clean
clean:
	@go clean
	@rm -fr release

#:help: release     | Release the product, setting the tag and pushing.
.PHONY: release
release:
	semtag final -s minor
	git push --follow-tags

#:help: test        | Installs go tools
.PHONY: test
test:
	@echo "===================== lib =====================" | tee -a test.log
	@cd lib && { go test -v $(GOTEST_FLAGS) -tags '$(GOTAGS)' ./... 2>&1 ; echo $$? >> ../exit.log ; } | tee -a ../test.log | egrep '^(ok|FAIL|panic:|--- FAIL|--- PASS)'

#:help: tools       | Installs go tools
.PHONY: tools
tools:
	@mkdir -p .gotools
	@cd .gotools && if [[ ! -f go.mod ]]; then \
		go mod init terraform-archetype-tools ; \
	fi
	cd .gotools && go get -v $(GOTOOLS)

u := $(if $(update),-u)

export GO111MODULE=on

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: devel-deps
devel-deps: deps
	sh -c '\
      tmpdir=$$(mktemp -d); \
	  cd $$tmpdir; \
	  go get ${u} \
	    golang.org/x/lint/golint               \
	    github.com/Songmu/godzil/cmd/godzil    \
	    github.com/tcnksm/ghr                  \
	    github.com/Songmu/statikp/cmd/statikp; \
	  rm -rf $$tmpdir \
    '

.PHONY: test
test: deps
	go test

.PHONY: lint
lint: devel-deps
	golint -set_exit_status

.PHONY: cover
cover: devel-deps
	goveralls

.PHONY: release
release: devel-deps
	godzil release

u := $(if $(update),-u)

export GO111MODULE=on

.PHONY: deps
deps:
	go get ${u} -d
	go mod tidy

.PHONY: devel-deps
devel-deps: deps
	GO111MODULE=off go get ${u} \
	  golang.org/x/lint/golint            \
	  github.com/mattn/goveralls          \
	  github.com/Songmu/godzil/cmd/godzil

.PHONY: test
test: deps
	go test

.PHONY: cover
cover: devel-deps
	goveralls

.PHONY: release
release: devel-deps
	godzil release

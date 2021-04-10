# Stream
VELA_VERSION ?= master
# Repo info
GIT_COMMIT          ?= git-$(shell git rev-parse --short HEAD)
GIT_COMMIT_LONG     ?= $(shell git rev-parse HEAD)
VELA_VERSION_VAR    := github.com/oam-dev/kubevela/version.VelaVersion
VELA_GITVERSION_VAR := github.com/oam-dev/kubevela/version.GitRevision
LDFLAGS             ?= "-X $(VELA_VERSION_VAR)=$(VELA_VERSION) -X $(VELA_GITVERSION_VAR)=$(GIT_COMMIT)"

GOX         = go run github.com/mitchellh/gox
TARGETS     := darwin/amd64 linux/amd64 windows/amd64

TIME_LONG	= `date +%Y-%m-%d' '%H:%M:%S`
TIME_SHORT	= `date +%H:%M:%S`
TIME		= $(TIME_SHORT)

BLUE         := $(shell printf "\033[34m")
YELLOW       := $(shell printf "\033[33m")
RED          := $(shell printf "\033[31m")
GREEN        := $(shell printf "\033[32m")
CNone        := $(shell printf "\033[0m")

INFO	= echo ${TIME} ${BLUE}[ .. ]${CNone}
WARN	= echo ${TIME} ${YELLOW}[WARN]${CNone}
ERR		= echo ${TIME} ${RED}[FAIL]${CNone}
OK		= echo ${TIME} ${GREEN}[ OK ]${CNone}
FAIL	= (echo ${TIME} ${RED}[FAIL]${CNone} && false)

test:
	go test -count=1 -coverpkg=models -covermode=count -coverprofile=coverprofile.cov -run="^Test" -coverpkg=./... ./...
	go tool cover -func=coverprofile.cov -o coverage.txt

upload-test:
	bash <(curl -s https://codecov.io/bash)



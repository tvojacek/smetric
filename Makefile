PKG_NAME=azure-picker

PLATFORM := $(shell [ $$(uname) = "Darwin" ] && echo darwin || echo linux)

VERSION := $(if ${VERSION},${VERSION},0.0.0)
COMMIT_SHA := $(if ${CI_COMMIT_SHA},${CI_COMMIT_SHA},$(shell git rev-parse HEAD))

BUILD_TIME := $(shell date "+%FT%T")

#The name of the package should be in PACKAGE_INFO
PKG_NAME := $(if ${PKG_NAME},${PKG_NAME},"UNKOWN")
IS_TAG := $(shell git describe --tags --exact-match 2>&1 >/dev/null && echo 'ok')
ifeq ($(IS_TAG),ok)
BRANCH := $(shell git describe --tags --exact-match)
BRANCH_ESCAPED := $(BRANCH)
VERSION:=$(BRANCH)
else
BRANCH := $(if ${CI_COMMIT_REF_NAME},${CI_COMMIT_REF_NAME},$(if ${GIT_BRANCH},$(GIT_BRANCH:origin/%=%),$(shell git symbolic-ref --short HEAD)))
BRANCH_ESCAPED := $(shell echo $(BRANCH) |sed "s@[[:punct:]]@-@g" )
endif

ifeq ($(BRANCH),$(filter $(BRANCH),master develop))
BUILD_DATE:=$(shell date "+%s")
BUILD_SEQ:= $(if ${CI_PIPELINE_IID},${CI_PIPELINE_IID},$(if ${BUILD_NUMBER},${BUILD_NUMBER},$(BUILD_DATE)))
BUILD_NO := $(if ${CI_PIPELINE_IID},${CI_PIPELINE_IID},$(if ${BUILD_NUMBER},${BUILD_NUMBER},0~$(BUILD_DATE)))
VERSION:=$(BRANCH)
else
BUILD_SEQ := $(if ${CI_PIPELINE_IID},${CI_PIPELINE_IID},$(if ${BUILD_NUMBER},${BUILD_NUMBER},0))
BUILD_NO := $(BUILD_SEQ)~$(BRANCH_ESCAPED)
endif

BUILD_DIR:=$(if ${BUILD_DIR},${BUILD_DIR},build)

.DEFAULT_GOAL:= test
clean:
	rm -rf $(BUILD_DIR)

include Makefile.golang



.PHONY: all clean $(GOLANG_PHONY) $(DOCKER_PHONY)
all: $(GOLANG_ALL) $(DOCKER_ALL)


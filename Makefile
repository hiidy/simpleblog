# ==============================================================================
# 전역 Makefile 변수 정의로 이후 참조에 편리

COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 프로젝트 루트 디렉토리
PROJ_ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/ && pwd -P))
# 빌드 산출물, 임시 파일 저장 디렉토리
OUTPUT_DIR := $(PROJ_ROOT_DIR)/_output

APIROOT=$(PROJ_ROOT_DIR)/pkg/api

# ==============================================================================
# 기본 타겟을 all로 정의
.DEFAULT_GOAL := all

# Makefile all 가상 타겟 정의, `make` 실행 시 기본적으로 all 가상 타겟 실행
.PHONY: all
all: tidy format build

# ==============================================================================
# 기타 필요한 가상 타겟 정의

.PHONY: build
build: tidy  # 소스 코드 컴파일, tidy 타겟에 의존하여 자동으로 의존성 패키지 추가/제거
	@go build -v -o $(OUTPUT_DIR)/sb-apiserver $(PROJ_ROOT_DIR)/cmd/sb-apiserver/main.go

.PHONY: format
format:  # Go 소스 코드 포맷팅
	@gofmt -s -w ./

.PHONY: tidy
tidy:  # 자동으로 의존성 패키지 추가/제거
	@go mod tidy

.PHONY: clean
clean:  # 빌드 산출물, 임시 파일 등 정리
	@-rm -vrf $(OUTPUT_DIR)

.PHONY: protoc
protoc:
	@echo "===========> Generate protobuf files"
	@protoc                                              \
		--proto_path=$(APIROOT)                          \
		--proto_path=$(PROJ_ROOT_DIR)/third_party/protobuf    \
		--go_out=paths=source_relative:$(APIROOT)        \
		--go-grpc_out=paths=source_relative:$(APIROOT)   \
		$(shell find $(APIROOT) -name *.proto)

# simpleblog

간단한 Go 기반 블로그 API 서버입니다. 모듈화된 구조와 재사용 가능한 패키지 설계로 다른 프로젝트의 boilerplate로 활용할 수 있습니다.


## Why simpleblog?

- **명확한 계층 분리**: cmd(진입점), internal(내부 로직), pkg(공개 API)로 구조화하여 유지보수성 확보
- **project-layout 준수**: Go 커뮤니티 표준 디렉토리 구조 적용
- **다양한 서버 모드**: gRPC, gRPC-Gateway (HTTP/JSON), Gin HTTP를 동일한 인터페이스로 전환 가능
- **재사용 가능한 로깅 모듈**: Zap 기반 structured logging 패키지 제공
- **Makefile 기반 빌드**: 빌드, 코드 생성, 린트 등 프로젝트 관리 자동화
- **API 문서화**: Protobuf 기반 API 정의 및 OpenAPI/Swagger spec 제공

## Features

- Health Check endpoint (`GET /healthz`) - gRPC, HTTP 모두 제공
- Cobra/Viper 기반 설정 및 CLI flags 지원
- pprof endpoint (Gin 모드)

## Quick Start

### Prerequisites

- Go 1.25+
- Make
- [buf](https://buf.build/docs/installation) (protobuf codegen)

### Build
```bash
make build
```

바이너리는 `_output/sb-apiserver`에 생성됩니다.

Protobuf 코드 재생성
```bash
make buf
```

### Run

기본값은 gRPC-Gateway 모드입니다. gRPC는 `:6666`, HTTP는 `:5555`에서 실행됩니다.
```bash
./_output/sb-apiserver
```

서버 모드 변경
```bash
./_output/sb-apiserver --server-mode=grpc
./_output/sb-apiserver --server-mode=grpc-gateway
./_output/sb-apiserver --server-mode=gin
```

### Health Check

HTTP (Gin 또는 gRPC-Gateway 모드)
```bash
curl http://127.0.0.1:5555/healthz
```

gRPC
```bash
go run examples/client/health/main.go -addr 127.0.0.1:6666
```

## Configuration

| Flag | Environment Variable | Default | Description |
|------|---------------------|---------|-------------|
| `--server-mode` | `SIMPLEBLOG_SERVER_MODE` | `grpc-gateway` | 서버 모드 |
| `--http-addr` | `SIMPLEBLOG_HTTP_ADDR` | `:5555` | HTTP 주소 |
| `--grpc-addr` | `SIMPLEBLOG_GRPC_ADDR` | `:6666` | gRPC 주소 |

## Project Structure
```
cmd/sb-apiserver        # Entrypoint
internal/apiserver      # HTTP/gRPC server wiring
internal/pkg            # Logging, server utilities
pkg/api                 # Protobuf definitions & generated code
api/openapi             # OpenAPI spec artifacts
```

## Tech Stack

- Go 1.25
- gRPC / Protocol Buffers
- gRPC-Gateway
- Gin
- Cobra / Viper / Pflag
- Zap
- 
## Roadmap

- [ ] 비즈니스 계층 구현 (Store / Biz / Handler)
- [ ] 미들웨어 및 인터셉터 (Request ID, Recovery, CORS)
- [ ] 요청 파라미터 검증
- [ ] JWT 인증 / Casbin 권한 관리
- [ ] 단위 테스트 및 정적 분석
- [ ] Wire 의존성 주입


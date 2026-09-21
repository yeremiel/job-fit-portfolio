# Turn 014 — Go HTTP API Vertical Slice

- 날짜: 2026-09-21
- 기반: Turn 013 HTTP API Contract, Turn 011 accepted v0 evaluator.
- 시작 HEAD: `635a0c2` (`docs: define job-fit HTTP API contract`).
- 시작 시 README/Evaluation Contract 수정과 untracked API Contract가 이미 있었다. 이 작업을 보존했으며 Evaluation Contract에는 이번 Turn의 편집을 하지 않았다.

## Completed

`POST /v1/job-fit/evaluate`의 최소 실행 가능한 Go HTTP server를 구현했다. 요청 검증 → 공유 application service → 기존 evaluator/Jev → HTTP DTO → 공통 envelope 흐름을 연결했다. CLI도 동일 service를 사용하며 stdout domain JSON은 유지했다. 외부 API 없는 테스트, 정적 검사, CLI/server 빌드 및 실제 Prop Tech plus HTTP/Jev 요청 한 번을 검증했다.

사용자 공통 개발 문서 **API Response & Error Standard (v1)**（repository 외부 참조 문서, 원본 미포함）와 [Turn 013 API Contract](../api-contract.md)를 기준으로 response envelope와 오류 분류를 구현했다. `status`, HTTP status와 같은 숫자 `code`, 실패 시 문자열 `error`, human message, timestamp/requestId를 사용한다. 공통 문서 원본은 변경하지 않았다.

## Changed Files

### 신규

- `cmd/job-fit-server/main.go`, `main_test.go`: startup 설정/검증, loopback listener, server lifecycle 및 startup failure tests.
- `internal/application/service.go`, `service_test.go`: profile snapshot, 공유 평가 호출, metadata, concurrent snapshot/parity/cancellation tests.
- `internal/failure/failure.go`: transport-neutral error category.
- `internal/httpapi/handler.go`, `request.go`, `response.go`, `errors.go`, `handler_test.go`: request validation, DTO/envelope, UUID/context/logging, 오류 mapper, HTTP contract tests.
- `internal/jev/http_boundary_test.go`: typed provider failure mapping, 실제 로컬 HTTP 연결을 통한 Jev context 취소 검증. Provider transport는 fake다.
- 본 문서와 [Prop Tech plus verification artifact](turn-014-results/proptech-plus.json).

### 변경

- `cmd/job-fit/main.go`, `main_test.go`: 공유 service 사용, 기존 JSON 성공 출력 regression.
- `internal/evaluator/input.go`: 기존 parsing/normalization 재사용 및 단일-read profile fingerprint.
- `internal/evaluator/evaluator.go`: stage 사이 context 취소 확인, no-requirements 오류 category. 판단 로직은 그대로다.
- `internal/jev/client.go`: 안전한 typed 오류, provider timeout 구분, configured engine identity.
- `internal/jev/mapper.go`: 기존 pair capacity 오류에 category만 추가. Instructions/rationales/evidence mapping 변경 없음.
- `README.md`: server 실행/호출/설정/제한 및 shared architecture 안내.
- `docs/api-contract.md`: 구현 상태와 확정 timeout/error 전달 상태만 갱신. Wire contract 재설계 없음.

## Server Architecture

표준 라이브러리 `net/http`, `encoding/json`, `log/slog`, `crypto/rand`, `crypto/sha256`만 사용한다. Framework/dependency, health endpoint, 새 business endpoint는 추가하지 않았다.

```text
cmd/job-fit-server
  → load canonical profile + construct Jev client
  → application.Service
  → httpapi.Handler
      request validation / deadline
        → service.Evaluate
          → evaluator.Evaluate (existing stages)
            → Jev
        → response DTO / common envelope
```

설정은 필수 `--profile`, 선택 `--listen`(기본 `127.0.0.1:8080`), `--evaluation-timeout`(기본 `5m`)이다. Numeric loopback 주소만 허용한다. Profile/key/listen 오류는 startup nonzero로 종료하며 상세 path/key를 로그에 출력하지 않는다. Key는 `TYPESAFE_API_KEY` 환경변수에서만 읽는다. Server는 `.env`를 직접 로딩하지 않는다.

## Shared Application Boundary

기존 CLI의 file loading → evaluator 직접 호출을 file loading → snapshot/service → evaluator 호출로 바꿨다. 동일한 `PrepareJob`과 `evaluator.Evaluate`가 사용되며 profile/JD/metadata에 따른 별도 HTTP 평가 정책은 없다. Metadata company/title/sourceUrl은 handler에 남기고 service에는 description만 전달한다.

Service는 validated domain result와 evaluation metadata를 반환한다. CLI는 domain result만 기존 형태로 encode한다. HTTP presenter는 category를 sourceType으로 표시하고 Trace를 ID 목록으로 표현한다. Evidence ID/text/relation과 explanation provenance는 보존한다. 잘못된 Trace를 수정하거나 label을 보정하지 않는다. HTTP handler는 Jev client를 조립하지 않는다.

## Profile Loading / Versioning

Profile 파일은 startup에서 한 번 읽는다. 같은 bytes를 기존 parser로 검증하고 SHA-256을 계산한다. Service가 private snapshot을 보유하며 engine에 전달하는 evidence slice는 요청마다 복사한다. 동시 요청이나 engine의 slice mutation이 저장된 profile에 영향을 주지 않도록 했다. Hot reload는 없으며 변경은 재시작 후 반영된다.

`data.metadata.profileVersion`은 `sha256:<64 lowercase hex>`다. 파일을 변경해도 이미 생성한 service의 version/evidence가 유지되는 테스트와 동시 요청 race 검사를 통과했다. Version은 semantic version이나 persistence ID가 아니다.

## HTTP Endpoint / Envelope

`POST /v1/job-fit/evaluate`, JSON synchronous request/response, 성공 200. 필수 `job.description`; company/title/sourceUrl optional. Unknown field, duplicate key, null/type 오류, 잘못된 UTF-8 및 unpaired Unicode surrogate를 거부한다. Source URL은 absolute HTTP/HTTPS이고 credentials를 포함할 수 없으며 fetch하지 않는다.

Success는 `status/code/message/data/metadata`, error는 `status/code/error/message/metadata`이고 data를 생략한다. 422만 field/reason/message의 errors 배열을 포함한다. 모든 array는 빈 경우 `[]`로 출력한다. Response Content-Type은 application/json, Cache-Control은 no-store다.

각 handler request에 UUID v4를 생성한다. Caller의 ID header는 사용하지 않는다. Context의 requestId는 service/Jev request까지 이어지고 success/error response와 로그에서 연결된다. Persistence identifier가 아니다. Unknown route는 404, 정의된 route의 non-POST는 405 및 Allow: POST다. Server가 handler 이전에 거부한 malformed HTTP나 끊어진 연결은 JSON envelope 전달을 보장하지 않는다.

## Error Mapping

| HTTP | Error code | 분류 |
| --- | --- | --- |
| 400 | REQUEST_MALFORMED_JSON / REQUEST_INVALID | 문법·encoding 또는 구조/타입 오류 |
| 413 | REQUEST_INPUT_TOO_LARGE | Body/description/metadata/line limit |
| 415 | REQUEST_UNSUPPORTED_MEDIA_TYPE | JSON/UTF-8 이외 media type 또는 content encoding |
| 422 | REQUEST_VALIDATION_FAILED | 필수/blank/URL/NUL, no requirements, pair/payload capacity |
| 500 | INTERNAL_CONFIGURATION_ERROR / INTERNAL_ERROR | 서버 설정, upstream 401/403, 내부 실패 |
| 502 | INTERNAL_EVALUATION_FAILED | Network/upstream failure, invalid response/result/trace |
| 503 | INTERNAL_DEPENDENCY_UNAVAILABLE | Upstream 429/529 |
| 504 | INTERNAL_EVALUATION_TIMEOUT | 전체 deadline 또는 provider timeout |

작은 `failure.Kind`/typed error를 사용하고 context errors는 wrapping을 유지한다. Provider status → failure category → HTTP mapper 순서이며 raw string matching이 없다. Upstream 400/422는 caller 잘못으로 단정하지 않고 502로 반환한다. Domain consistency 오류도 service에서 evaluation failure로 분류한다. Unexpected internal error/panic은 안전한 500 응답이다.

오류 원문/상위 body/Authorization/key는 public message에 사용하지 않는다. 일부 오류에는 이미 외부 호출 비용이 발생했을 수 있으며 partial result나 자동 재시도는 없다.

## Timeout / Cancellation

| 설정 | 값 | 이유 |
| --- | --- | --- |
| Total evaluation deadline | 기본 5분, positive~1시간 설정 허용 | 기존 최대 5개 순차 provider 호출과 호출당 60초 제한을 고려한 bounded synchronous MVP |
| Jev client timeout | 기존 60초/call | 기존 동작 유지 |
| ReadHeaderTimeout | 5초 | Header 대기 제한 |
| ReadTimeout | 15초 | 작은 JSON body 수신 제한 |
| WriteTimeout | evaluation deadline + 30초 | Body read/응답 여유를 두고 timeout envelope 전달 가능하도록 구성 |
| IdleTimeout | 60초 | 유휴 연결 정리 |
| Shutdown grace | 10초 | SIGINT/SIGTERM 시 context 취소 후 연결 정리 |

Evaluation deadline은 요청 검증 후 평가 작업에 적용한다. HTTP request context → deadline child → service → core stages → Jev HTTP request로 전달한다. Client disconnect와 timeout을 로컬 HTTP test에서 각각 검증했다. 종료 신호는 server BaseContext를 취소하고 Shutdown/필요 시 Close로 정리한다. 원격 provider 처리/과금의 중단까지 보장하지 않는다. Context를 무시하는 임의 Engine을 강제로 중단시키는 별도 worker/goroutine 구조는 도입하지 않았다.

## Request Limits

Turn 013의 body 128 KiB, decoded description 16 KiB(trim 전), 최대 64 non-empty lines를 그대로 적용했다. Metadata는 company 256/title 512/sourceUrl 2,048 UTF-8 bytes다. Body 한도는 JSON escape expansion까지 고려한 계약 값이다. Profile 1–32 evidence, pair 512, generated Jev request 192 KiB 제한은 유지한다. 평가 threshold가 아니며 자동 truncation/decomposition/parser를 추가하지 않았다.

## Logging

`log/slog` JSON으로 requestId, method, route, code, error category, outcome, elapsed를 기록한다. 일치하지 않는 route는 `unmatched`로 기록해 임의 path/query에 실린 정보를 노출하지 않는다. Body, headers, profile, JD, raw upstream 오류는 로깅하지 않는다. 취소된 요청은 응답을 만들지 않고 outcome=caller_canceled/code=0으로 기록한다. Code=0은 로그의 미전송 상태이지 HTTP 응답 코드가 아니다. Error code와 outcome을 안전한 원인 맥락으로 사용하며 stack/raw error를 덤프하지 않는다.

## Tests

다음 명령이 통과했다.

```sh
go test ./...
go test -race ./...
go vet ./...
go build -o ./bin/job-fit ./cmd/job-fit
go build -o ./bin/job-fit-server ./cmd/job-fit-server
```

- 기존 evaluator/Trace/qualification regression 유지.
- HTTP 200/400/413/415/422/500/502/503/504 및 404/405 검증.
- Status/code 일치, error에 data 없음, UUID/timestamp, array/DTO/sourceType/evidence mapping, requestId context/log correlation.
- Malformed/trailing/duplicate/unknown/type/null/UTF-8/surrogate, profile override, body/description/line/metadata limits, URL validation.
- Strong+Limiting, unknown Supporting, duplicate/nonexistent/mismatched Trace 거부.
- Safe response/log와 upstream secret redaction.
- Service snapshot/version 동시 재사용, core/service parity, stage cancellation 및 CLI 성공 JSON 회귀.
- 실제 로컬 HTTP connection → service → Jev fake transport까지 disconnect/deadline 전달. 일반 테스트는 외부 API/key 없이 실행.

초기 sandbox 실행은 `log/slog`를 Go 표준 라이브러리에서 찾지 못하는 환경 오류로 실패했다. 승인된 로컬 Go 실행 환경에서 동일 테스트를 실행해 통과했다. 라이브러리 dependency를 추가하여 우회하지 않았다.

## Manual HTTP / Jev Verification

기존 `.env`를 실행 shell의 환경변수로 내보내고 loopback `127.0.0.1:18080`에서 server를 실행했다. Key를 대화/로그/응답에 출력하지 않았다. 기존 Prop Tech plus sample 전체를 description으로, company/title만 metadata로 넣은 임시 request 파일을 사용했다.

```sh
./bin/job-fit-server --profile ./data/profile.json --listen 127.0.0.1:18080
curl --silent --show-error --max-time 330 \
  --request POST http://127.0.0.1:18080/v1/job-fit/evaluate \
  --header 'Content-Type: application/json' \
  --data-binary @/tmp/jobfit-turn014-request.json
```

실제 evaluation request는 **한 번만** 실행했다. 결과를 reference에 맞춰 재실행/tuning하지 않았다.

| 항목 | 결과 |
| --- | --- |
| HTTP / Envelope | 200 / success / code 200 |
| Curl elapsed | 2.513035초 |
| Overall | Partial |
| Requirements | 15 |
| Strong / Partial / Weak / Unknown | 3 / 7 / 1 / 4 |
| Supporting / Limiting / Non-decisive | 3 / 3 / 9 |
| Engine / model | jev / jev-1.13.0 |
| evaluatedAt | 2026-09-21T11:32:25.266377Z |
| requestId | 0e9e452b-c711-49c9-994e-3e0d2facbb46 |
| profileVersion | sha256:cae86ad4906124b0eee8283c05bc74d439543f75d9b2e7496f8d6895e32f09cd |

응답 fingerprint가 실제 profile 파일과 일치하고, 15개 requirement 모두 Trace에 한 번씩 연결되며, sourceType/evidence arrays/metadata가 유지되는 것을 확인했다. Server log의 requestId도 동일했다. SIGINT 후 정상 종료했다.

[저장된 응답](turn-014-results/proptech-plus.json)의 SHA-256: `5c7b0e8d738060d616b2ecb925a001e34aa0c96baa37e77dd7d5764d9c4dd0f3`. 이는 개발 검증 artifact이며 service의 runtime persistence 기능이 아니다.

## CLI Regression

기존 입력 오류/missing-key 테스트와 offline 성공 테스트를 통과했다. 성공 출력은 기존 overallMatch/overallReasoning/requirements/summary/explanationSource/decisionTrace이며 HTTP envelope나 새 metadata를 CLI에 추가하지 않았다. 공유 service와 기존 core의 domain 결과 parity도 확인했다. 별도의 실제 CLI Jev 호출은 하지 않았다.

## Decisions Made

- net/http + slog, 새 dependency 없음.
- `--profile`, numeric loopback `--listen`, 기본 5분 `--evaluation-timeout` 및 명시적 server lifecycle timeouts.
- 기존 engine/core를 그대로 호출하는 얇은 service와 typed failure category.
- Opaque profile snapshot + 요청별 evidence slice 복사.
- RequestId는 context 전파, 로그/응답 correlation용이며 persistence 없음.
- 공통 code 상수와 단일 HTTP error mapper, 422 상세 errors, safe panic recovery.
- Health endpoint 없이 유일한 business endpoint만 구현.

## Problems / Findings

계약을 변경해야 하는 충돌은 발견하지 않았다. 일반 error 문자열만 있던 Jev/client와 일부 capacity 오류에 category를 추가해야 했으며 timeout 읽기/전송 오류와 context cancellation을 분리했다. 이는 transport 관측과 오류 전달 변경이며 evaluation semantics 변경이 아니다.

한 번의 실제 성공과 offline failure tests는 장기 운영/부하/다수 caller 안정성 검증을 뜻하지 않는다. 총 request concurrency 제한, 비용 제한, auth, idempotency, persist/recover는 없다. Personal loopback 전제를 유지한다. HTTP contract의 엄격한 Unicode/duplicate-key validation은 구현했으나 JD 의미 해석 오류를 방지하는 parser는 아니다.

기존 Composite Public / Energy Context ambiguity와 post-hoc Trace의 비인과성은 그대로다. Human Reference 또는 이전 실행과 label/distribution 차이는 이번 HTTP 구현의 tuning 근거로 사용하지 않았다.

## Deferred / Suggested Next Step

OpenClaw workflow에서 호출할 때의 입력 준비, timeout/실패 처리, 결과 사용 경계를 planning conversation에서 검토한다. MCP, OpenClaw/Notion integration, auth/public deployment, DB/history, batch, crawler/URL fetch, profile CRUD, multi-user 및 evaluation tuning은 시작하지 않았다.

기존 Turn 000–013, profile, JD, evaluation contract/semantics, runtime policy, qualification/trace instructions는 작업 시작 시 hash와 비교해 보존했다. Commit은 생성하지 않았다.

# Turn 013 — HTTP API Contract

- 날짜: 2026-09-21
- 범위: 문서 설계만. HTTP server, dependency, CLI refactoring, evaluator/prompt/profile 변경 및 Jev 실행 없음.
- 시작 상태: `c5f1b31` (Turn 012 documentation), clean working tree.
- 기반: Turn 011 accepted v0 evaluation baseline, Turn 012 제품 책임. 이번 계약의 v1은 별도의 public HTTP version이다.

## Completed

현재 CLI, input loader, domain result, evaluator workflow, Jev client/mapper와 Trace validation을 확인하고 [HTTP API Contract](../api-contract.md)를 작성했다. 사용자 prompt의 여덟 결정 질문을 모두 문서화했다. 코드가 이미 이 계약을 구현한다고 주장하지 않는다.

## Endpoint / Request Contract

`POST /v1/job-fit/evaluate`, synchronous JSON request/response, 성공 200. 필수 `job.description`, 선택 `company`, `title`, `sourceUrl`로 결정했다. Company/title을 필수로 만들면 본문만 확보한 CLI/script workflow에 불필요한 제약이 생기므로 optional이다. SourceUrl은 provenance metadata이고 fetch하지 않는다.

Company/title을 JD에 자동 결합하면 현재 줄 기반 identity와 평가 state가 달라진다. 따라서 세 field 모두 식별 metadata로만 사용하고 role/company context가 평가에 필요하면 description에 포함하도록 명시했다. Metadata와 본문의 충돌을 service가 조사하거나 수정하지 않는다.

Unknown/duplicate fields, 타입 오류, null 및 요청별 profile override를 거부한다. 기존 16 KiB JD, 64 lines, 512 pairs 제한을 보존하고 HTTP envelope 128 KiB 및 metadata field 한도를 별도로 정했다. JSON escaping을 고려한 envelope 한도이며 Match threshold가 아니다. 상세 규칙과 실패 시점은 API Contract에 있다.

## Response Contract / Alternatives

검토한 방향:

| 선택지 | 장점 | 비용 / 판단 |
| --- | --- | --- |
| Domain JSON 직접 노출 | 구현량 최소 | Go domain 변경과 외부 contract가 결합됨 |
| 작은 API-specific DTO | 의미를 보존하며 external naming/metadata를 독립 관리 | 소규모 명시적 mapping 필요. **채택** |
| Label와 text만 대폭 축약 | 작은 payload | Evidence relation/identity, provenance와 Trace 해석 정보 손실. 제외 |

Requirement `category`를 public `sourceType`으로 표현한다. Evidence는 ID·원문·relation을 그대로 유지한다. Trace는 ID 목록으로 반환하여 requirements에서 Match/source/evidence를 resolve한다. 중복 summary와 Trace의 반복 text는 제외한다. Overall/requirement reasoning, explanationSource, post-hoc Trace source는 보존한다. API가 설명을 새로 생성하거나 Match/Trace를 고치지 않는다.

모든 array는 `[]`, 필수 field는 항상 반환한다. 원본 JD 전체는 echo하지 않고 request metadata만 `job` object로 반환한다. Requirement ID는 result-local이고 수정 전후 영구 ID가 아니다. Evidence를 응답하므로 caller에게 개인 경력 정보가 전달된다는 점도 명시했다.

## Candidate Profile Handling / Metadata

Server는 시작 시 canonical profile 하나를 읽고 검증하여 immutable snapshot으로 사용한다. 변경 반영은 restart이며 hot reload/profile selection API는 없다. `profileVersion`은 실제 읽은 동일 원본 bytes의 `sha256:<64 lowercase hex>`다. 별도 version 관리 파일보다 단순하지만 formatting-only 변경도 version 변경으로 보인다는 trade-off를 기록했다.

Evaluation metadata는 profileVersion, 성공 완료 UTC evaluatedAt, engine, 실제 사용한 configured engineModel이다. Request metadata는 job에 분리했다. 공통 응답 표준 반영 후 response timestamp/requestId를 채택했다. EvaluationId/persistenceId는 제외했다. Model/profile metadata만으로 evaluator build나 전체 실행의 재현성을 증명하지 않는다.

## Application Boundary

현재 `evaluator.Evaluate(ctx, engine, profile, job)`가 Extract → Assess → Overall → Trace → validation을 이미 공유한다. Engine은 provider interface다. 향후 얇은 application service는 snapshot, completion metadata, categorized error만 보탠다. 입력은 동일 normalization/limits를 적용한 JD, 출력은 validated domain result와 evaluation metadata다.

HTTP adapter가 DTO와 HTTP status를 담당하고 CLI는 파일/flags/stdout/stderr를 담당한다. CLI의 explicit profile 선택과 기존 JSON presentation을 유지할 수 있다. Profile 선택 위치가 다를 뿐 평가 workflow는 하나다. Generic mapping framework, repository 계층 또는 두 번째 evaluation pipeline은 만들지 않는다. 향후 MCP는 REST를 통해 동일 boundary를 사용한다.

## Error Contract

사용자 공통 개발 문서 **API Response & Error Standard (v1)**（repository 외부 참조 문서, 원본 미포함）를 참조했다. 성공은 `status/code/message/data/metadata`, 오류는 `status/code/error/message/metadata` envelope다. Code는 HTTP status와 같은 숫자이고 error는 문자열 분기 코드다. 422에는 field/reason/message를 가진 errors 배열을 사용한다.

| Status | 분류 |
| --- | --- |
| 400 | JSON/형식/구조 오류 |
| 413 / 415 | 입력 크기 / media type 오류 |
| 422 | 필수값/URL 등의 의미 검증 실패, 추출 결과 없음, pair/payload 제한 |
| 500 | Server configuration 또는 unexpected internal error; 별도 code |
| 502 | External evaluation failure, invalid upstream/result 또는 consistency failure |
| 503 | Explicit upstream busy/rate limit |
| 504 | Connected caller에 반환하는 evaluation timeout |

404/405도 동일 envelope로 정의했다. Upstream의 HTTP status나 raw error를 그대로 전달하지 않는다. 인증 거절은 caller 인증 오류가 아니며 upstream 400/422도 caller invalid request로 단정하지 않는다. Safe error/message 및 HTTP status mirror code만 전달하며 secret, local path, request state, raw upstream body는 노출하지 않는다.

## Timeout / Cancellation

HTTP context → finite total evaluation deadline → evaluator → Jev로 전달한다. CLI도 signal-aware context를 공유한다. Caller disconnect 시 후속 단계를 중단하되 응답 전달이나 원격 연산 중단을 보장하지 않는다. 연결된 caller의 timeout은 504다. 총 deadline과 server timeout의 구체 duration은 구현 Turn에서 정한다. 현재 provider call마다 60초 제한이 있다는 사실을 총 evaluation timeout으로 오해하지 않도록 기록했다.

## Decisions Made

1. Company/title optional, sourceUrl 포함 세 field는 metadata only.
2. Canonical profile startup snapshot + raw-byte SHA-256 version; request override 금지.
3. Domain을 직접 expose하지 않는 작은 HTTP DTO; trace ID references 및 evidence relation/provenance 유지.
4. Request metadata, evaluation metadata, 공통 response metadata 분리. 공통 로그 지침에 따라 server-generated requestId만 채택하며 persistence ID는 미도입.
5. Explicit error categories/status, partial success 금지, automatic retry 없음.
6. Synchronous `/v1`, shared evaluator operation + thin service boundary, API 평가 logic 추가 없음.
7. Strict request와 additive response compatibility, loopback-default 구현 방향 및 server-side secret 원칙.

## Problems / Findings

- 현재 Jev client는 독립 provider timeout과 network error를 하나의 문자열 오류로 축약한다. Context cancellation은 wrapping되지만 HTTP status별 원인은 typed category로 노출되지 않는다. 향후 typed/sentinel error 전달이 필요하며 문자열 parsing으로 해결하지 않는다.
- 현재 job normalization은 file loader 안에 있다. HTTP implementation에서 같은 규칙을 공유할 작은 입력 준비 경계가 필요하다. 이번에 refactor하지 않았다.
- No-capability/pair/payload 오류는 evaluation 중 발생할 수 있다. 입력 validation 성공이 evaluation 성공 또는 무비용 실패를 보장하지 않는다.
- DTO가 보존하는 Trace는 post-hoc attribution이다. 새로운 인과 설명이나 더 강한 의미적 보장을 만들지 않는다.
- Composite Public / Energy Context는 여전히 Ambiguous이며 본 contract로 해결하지 않는다.

## Security / Non-goals

Personal/internal tool이며 auth/public exposure는 구현하지 않는다. 향후 service 기본값은 loopback으로 하고 노출 범위 변경 시 별도 access-control 검토가 필요하다. Jev key는 server-side에만 존재한다. Evaluation 시 profile/JD가 Jev로 전달되는 기존 흐름은 유지한다.

검색, crawler/URL fetching, 회사/SES/派遣 분류, Notion/OpenClaw integration, profile CRUD/multi-user, DB/persistence/batch/async, MCP, UI, score/probability/recommendation은 제외했다. Application code, dependency, evaluator, prompt, profile, samples 및 과거 Turn 기록을 수정하지 않았다.

## Changed Files / Validation

- [API Contract](../api-contract.md): 신규 public HTTP contract.
- 본 Turn 기록: 결정·근거·구현 전제 기록.
- [README](../../README.md): contract link와 현재 결정 상태 최소 갱신.
- [Evaluation Contract](../evaluation-contract.md): Turn 013 연결 note 추가; 평가 규칙 변경 없음.

검증 완료: JSON 예시 4개의 parsing과 success example의 trace 참조/Strong consistency/profile fingerprint 형식, 변경 문서의 로컬 링크 및 code fences 검사 모두 통과했다. `git diff --check`가 통과했고 변경 파일이 위 문서 4개뿐임을 확인했다. 문서-only 변경이므로 Go test/vet/build 또는 실제 Jev API를 실행하지 않았다. Commit은 생성하지 않았다.

## Open Questions

외부 contract의 필수 질문은 결정했다. 다음 구현 Turn에서 total deadline 값, HTTP server read/write/shutdown 설정, 로컬 profile path 설정 방식, 패키지 배치 및 typed error의 최소 구현을 선택한다. Framework는 선택하지 않았다. 이후 network exposure 변경에 필요한 auth, evaluator build metadata 필요성은 별도 범위다.

## Suggested Next Step

Planning conversation에서 계약을 검토한 뒤 **Go HTTP server 최소 vertical slice**를 구현한다. DTO validation/mapping, canonical snapshot/version, shared input/service boundary, categorized failures 및 cancellation을 포함한다. Mocked Engine/HTTP 기반으로 성공·오류·DTO parity·timeout/cancellation을 검증하고 기존 CLI regression을 유지한다. Evaluation tuning 또는 integration 확장은 하지 않는다. 이번 Turn에서는 구현을 시작하지 않았다.

## 공통 표준 반영 후속 요청

사용자가 대부분의 프로젝트에 참조하는 공통 개발 문서의 응답 표준을 이번 API 계약에도 적용하도록 요청했다. 지정한 `api-responde.md` 대신 실제 존재하는 `api-response.md`를 확인해 읽었다. 원본 공통 문서는 변경하지 않았으며, job-fit용 성공/실패/422 템플릿은 API Contract에 포함했다.

초안의 flat success와 nested error를 공통 envelope로 대체했다. 공통 timestamp, 로그 correlation용 requestId를 추가하고 validation은 400 형식 오류와 422 의미 오류로 구분했다. Profile/engine/evaluatedAt은 data.metadata에 유지했다. 외부 공개/구현 전 contract 수정이므로 `/v1`을 유지하며 운영 client migration은 없다. 구현 Turn에서 envelope/status 일치, error mapping, requestId correlation 테스트를 추가해야 한다. 이번에도 application code나 평가 의미는 변경하지 않았다.

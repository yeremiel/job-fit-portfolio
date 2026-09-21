# Turn 017 — Architecture & Project Documentation

- 날짜: 2026-09-21
- 범위: 현 구현의 architecture 설명만 작성. 새 설계/기능/code/API behavior 변경 없음.
- 시작 상태: clean working tree, HEAD `b65b4e7` (`docs: add job-fit API documentation`).

## Completed

[Architecture Overview](../architecture/overview.md) 한 문서에 system overview, component/runtime/dependency boundaries, REST/CLI flow, evaluation/Trace/profile, validation/error/context/security, project structure와 기존 주요 결정을 정리했다. README에는 문서 링크만 추가했다.

## Changed Files

- `docs/architecture/overview.md`: 신규 architecture 문서.
- 본 `docs/development/turn-017.md`: 확인 근거와 작업 기록.
- `README.md`: Architecture Overview와 Turn 017 링크.

## Architecture Overview / Runtime Components

Current code를 우선 확인했다. Entry points가 Engine/service를 조립하고 HTTP와 CLI가 같은 `application.Service`와 `evaluator.Evaluate`를 호출한다. 그림에서 REST가 Jev에 직접 연결되는 것으로 보이지 않도록 실제 호출 방향을 표시했다.

Application service는 snapshot, input 준비, core 호출, metadata 및 오류 분류를 담당한다. Evaluation stage sequencing은 core, semantic question/Choice와 response/rationale mapping은 Jev adapter에 있다. Core의 Engine interface와 adapter의 구현 관계를 설명했다. `httpapi`가 domain types/input/Trace validation helpers도 참조하므로 모든 dependency가 service package만 거친다는 이상적인 계층을 주장하지 않았다.

## Request Flow

Startup의 key/client 구성 → profile 한 번 로딩/hash → listener 순서와 요청별 requestId → route/validation → deadline → service/core → DTO/envelope 흐름을 구분했다. Request 시 profile 파일을 재로딩하지 않는다. Metadata는 handler에 남고 description만 evaluation에 들어간다.

CLI는 profile/JD 파일과 signal context로 같은 service/core를 호출하며 REST를 경유하지 않는다. Domain JSON과 HTTP DTO/envelope의 presentation 차이, CLI에 HTTP의 전체 5분 deadline이 적용되지 않는 operational 차이도 명시했다.

## Evaluation / Jev Boundary

`Extract → Assess → Overall → Trace → ResolveTrace`를 실제 메서드 순서대로 설명했다. Assess 내부의 evidence relation과 Match/rationale는 별도 provider 호출이며 정상 전체 흐름은 5회다. 모든 항목에서 linked evidence가 없으면 Match 질문을 생략해 4회가 될 수 있다.

Overall은 Jev Choice로 결정하며 평균/다수결이 아니다. Profile/JD/policy/실제 assessments와 category가 전달된다. Trace는 이미 선택된 Overall에 대한 별도 structured post-hoc attribution이며 Overall로 feedback되지 않는다. 원문 evidence 보존과 제한된 rationale 렌더링도 구분했다.

## Profile Boundary

Submission resume와 canonical evaluation evidence의 목적 차이, 단일-read raw bytes의 SHA-256, opaque snapshot, 요청별 slice copy, restart 반영 및 CLI `--profile` 선택을 기록했다. Generic caller profile input은 현재 API에 없다.

## Error / Cancellation

HTTP validation, Jev wire/choice 검증, core consistency, HTTP presenter의 경계를 구분했다. Typed failure/context wrapping → service 분류 → HTTP status/error envelope를 설명했다. 모든 기존 core 오류가 처음부터 typed였다고 표현하지 않았다.

400/413/415/422/500/502/503/504 및 routing 404/405를 요약하고 상세 응답은 API docs로 연결했다. HTTP context → service → core → Jev request 취소 전파, provider 60초와 HTTP 전체 기본 5분의 차이, disconnect/shutdown 및 remote cancellation 보장의 한계를 설명했다.

## Security Boundary / Project Structure

Loopback-only, 기본 127.0.0.1:8080, auth/public deployment 없음, server-side env key/profile, secret-safe HTTP logs와 requestId correlation을 기록했다. Profile/JD가 외부 Jev로 전송된다는 점과 candidate evidence가 result에 포함된다는 점을 유지했다.

실제 `cmd`, `internal/application`, `evaluator`, `failure`, `httpapi`, `jev`, `data`, `samples`, `docs` 구조를 확인했다. Tests는 package-local이고 history JSON은 runtime persistence가 아니다. go.mod는 Go 1.24.0이며 third-party module dependency는 없다.

## Architecture Decisions

Go/CLI-first, Engine/adapter 분리, structured output/Trace, Human Reference 비정답 원칙, confirmed violation 수정 후 baseline 수용, reusable REST/shared service, snapshot/DTO, portfolio scope freeze를 해당 Turn/test에 연결했다. 새 architecture 결정을 추가하지 않았다.

## Out of Scope / Extension Points

Search/discovery/company classification/crawling/OpenClaw/Notion/MCP/DB/UI 등을 현재 system 밖으로 표시했다. Generic input, REST를 호출하는 optional thin MCP adapter, external automation은 확장 가능성이지 구현 약속 또는 MVP 완료 조건이 아니다.

## Evidence Reviewed

- Current CLI/server entry points, application service, evaluator orchestration/input/types/trace.
- Jev mapper/trace/client/qualification 및 failure categories, HTTP request/handler/response/errors.
- CLI 성공/오류 tests, application snapshot/parity/cancellation tests, Jev pipeline/HTTP cancellation tests, HTTP validation/DTO/error tests.
- API Contract, Evaluation Contract, Match Semantics 및 Turn 014–016 문서. Go 선택 등 과거 결정은 Turn 004/006/007/010/011/012–015 근거를 참조했다.

## Problems / Findings

현재 구현을 바꿔야 할 architecture 충돌은 발견하지 않았다. 다만 문서가 과장되기 쉬운 경계를 명시했다: shared service는 두 프로세스가 같은 instance를 공유한다는 뜻이 아니며, HTTP가 core helper를 직접 사용하는 dependency가 있고, service 자체가 모든 semantic 단계의 orchestration을 구현하지 않는다.

기존 Composite Public / Energy Context ambiguity, limited rationale, 줄 기반 extraction, language/education 제외, provider variability 및 Trace의 비인과성은 보존했다. Structural validation을 semantic 정확성 보장으로 설명하지 않았다.

## Validation

문서 링크/anchor/code fences, 실제 project directory와 source links, `git diff --check` 검증이 통과했다. 작업 시작 시 hash와 비교해 README 외 기존 파일이 모두 보존되고 신규 문서가 두 개뿐임을 확인했다. Code 변경이 없으므로 Go tests/vet/build와 Jev 호출은 수행하지 않았다. 기존 기록/계약/입력/artifact를 수정하거나 commit하지 않았다.

## Deferred

Portfolio README polish, AI-assisted development story 및 final regression은 다음 Turn으로 남긴다. Dependency, OpenAPI/Swagger, MCP/OpenClaw, DB/UI/crawler/auth 및 evaluator 변경은 수행하지 않았다.

## Suggested Next Step

다음 Turn에서 **최종 README와 portfolio narrative를 정리하고 final regression을 수행**한다. 이번에는 이를 미리 시작하지 않았다.

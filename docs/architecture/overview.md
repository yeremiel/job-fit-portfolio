# Architecture Overview

> **Server-side canonical candidate profile과 제공된 Job Description을 비교하여 Jev 기반 structured Job Fit evaluation을 반환하는 Go application.**

여러 JD를 같은 사실 기반 경력과 비교할 때 평가 기준·근거·미확인 정보를 구조적으로 확인하는 개인용 도구다. Tailored resume가 평가 근거를 바꾸는 bias를 피하고, 판단 결과를 재사용 가능한 CLI/REST workflow로 제공한다. Match는 채용 확률이나 지원 추천이 아니며 최종 판단은 사람이 한다.

이 문서는 Turn 017 시점의 **실제 구현**을 설명한다. 새 architecture나 기능 계획이 아니다. 제품 범위는 [Turn 015](../development/turn-015.md)에서 동결했다. API 사용법은 [Specification](../api/specification.md)과 [Endpoint Documentation](../api/job-fit.md), 설계 이유는 [API Contract](../api-contract.md)를 참조한다.

## System Overview

```text
Local profile file                 Local profile file
        |                                  |
CLI process                         HTTP server process
cmd/job-fit                         cmd/job-fit-server
  file input                          startup configuration
  domain JSON output                  httpapi: validation / DTO / envelope
        |                                  |
        +---------- same implementation ---+
                           |
                  application.Service
                  snapshot + normalized JD
                           |
                   evaluator.Evaluate
                   Engine interface calls
                           |
                      jev.Client
                 implements evaluator.Engine
                           |
       ---------------- network / system boundary ----------------
                           |
                       Jev API
                  external dependency
```

CLI와 HTTP server는 별도 실행 파일/프로세스다. 같은 service/core **코드를 공유**하며 하나의 실행 중 service instance를 원격으로 공유하는 것은 아니다. CLI는 REST를 경유하지 않는다. HTTP handler도 Jev를 직접 호출하거나 별도 평가 pipeline을 만들지 않는다. 각 entry point가 Jev client와 service를 조립한다.

### Runtime calls vs Go dependencies

```text
cmd/job-fit         -> application + evaluator(input loading) + jev
cmd/job-fit-server  -> application + httpapi + jev
httpapi            -> application output / evaluator types + validation / failure
application        -> evaluator + failure
jev                -> evaluator (Engine implementation and domain types) + failure
evaluator          -> failure
```

Evaluator는 Jev package를 import하지 않고 `Engine` interface를 호출한다. Jev adapter가 core의 types/interface에 의존하는 방향이다. HTTP도 shared input 준비와 DTO consistency 검증에 evaluator types/helpers를 사용한다. 따라서 모든 Go import가 application package만 거치는 엄격한 계층이라고 표현하지 않는다. Transport 규칙이나 provider wire types가 core에 들어가지 않는 것이 현재 경계다.

## Runtime Components

| Component | 실제 책임 | 담당하지 않는 일 |
| --- | --- | --- |
| [CLI](../../cmd/job-fit/main.go) | `--profile`/`--job` 파일 로딩, 환경변수, signal context, service 호출, domain JSON stdout 및 error stderr/exit code | HTTP 호출, CLI 전용 평가 정책 |
| [Server entry point](../../cmd/job-fit-server/main.go) | Flags/env, Jev client/service 구성, profile startup loading, loopback listener, HTTP timeouts/shutdown | Requirement/Match 판단 |
| [HTTP API](../../internal/httpapi/handler.go) | Route/method, requestId/context, 요청 검증, evaluation deadline, response DTO/common envelope, error mapping, sanitized logging | Engine 조립, Overall 계산, JD metadata의 임의 결합 |
| [Application Service](../../internal/application/service.go) | Opaque canonical snapshot/version, shared JD 준비, 요청별 evidence slice 복사, `evaluator.Evaluate` 호출, error 분류 보존/보완, 완료 metadata | 단계별 Choice 생성, HTTP status 결정 |
| [Evaluator](../../internal/evaluator/evaluator.go) | Domain/Engine 정의, Extract→Assess→Overall→Trace 순서, context 확인, assessment identity/Match 검증, summary, Trace resolve/validation | Provider HTTP/인증, REST envelope |
| [Jev Adapter](../../internal/jev/mapper.go) | JD line category, evidence relation, Match/rationale, Overall, Trace에 대한 structured questions/choices; provider 응답을 domain으로 mapping; 제한된 설명 렌더링 | 자유 텍스트 경력 생성, company별 target label, numeric score |
| [Jev Client](../../internal/jev/client.go) | HTTP 요청/인증, provider response shape/allowed-choice 검증, payload/response limits, transport 오류 category | Caller HTTP response/status |
| [Failure categories](../../internal/failure/failure.go) | 작은 transport-neutral typed error 및 wrapping을 통한 분류 | 거대한 exception hierarchy나 raw upstream error 노출 |
| [Canonical profile](../../data/profile.json) | ID/text evidence의 factual source. HTTP startup snapshot 또는 CLI 실행별 파일 로딩 | 제출용 resume, DB/profile 관리 서비스 |

Core가 실행 순서를 정하지만 실제 semantic classification과 rationale mapping은 Engine 구현인 Jev adapter에 있다. 이 역할을 HTTP 또는 application service에 복제하지 않는다. 여기서 orchestration은 **평가 단계 실행**이며 외부 회사 검색/취업 활동 orchestration과 다르다.

## REST Request Flow

### Startup — 요청 처리 전 한 번

1. Server가 `--profile`, `--listen`, `--evaluation-timeout`과 `TYPESAFE_API_KEY`를 확인한다.
2. Jev client를 만들고 configured engine/model identity를 얻는다.
3. Application service가 canonical profile 파일을 한 번 읽어 검증하고 같은 bytes를 hash한다. Snapshot이 준비된 뒤 listener를 시작한다.
4. Profile/key/config 실패는 startup 실패다. 이때는 정상 HTTP listener가 없으므로 client error envelope가 나오지 않는다.

### Request — 동기 실행

```text
HTTP request
  → server-generated requestId → request context
  → exact path / POST 확인
  → media type / query / bounded body / JSON 구조 / field 의미 검증
  → evaluation deadline child context
  → Application Service.Evaluate(description)
      → shared PrepareJob
      → stored snapshot에서 요청용 evidence slice 복사
      → evaluator.Evaluate(ctx, engine, profile, job)
          → Extract: JD line classification
          → Assess: evidence relation → Match/rationale
          → requirement identity / count / Match 검증 및 summary
          → Overall: Jev Choice
          → Trace: selected Overall에 대한 별도 Jev Choice
          → ResolveTrace / consistency validation
          → EvaluationResult
      → profileVersion / evaluatedAt / engine / engineModel
  → caller cancellation / deadline / failure 확인
  → HTTP DTO mapping 및 Trace/source consistency 확인
  → success envelope + response timestamp/requestId
  → HTTP 200 JSON + sanitized request log
```

오류가 발생하면 이후 평가 단계를 진행하거나 불완전한 결과를 성공으로 반환하지 않는다. HTTP validation은 evaluation 전에 `PrepareJob`으로 같은 제한을 확인하고, service도 shared input을 다시 준비한다. 이는 다른 semantic workflow가 아니라 input boundary 검증이다.

Company/title/sourceUrl은 handler에 보존되어 결과의 `data.job`으로 반환한다. Service에 전달되는 JD는 description뿐이며 metadata를 앞에 붙이지 않는다. Profile을 매 요청마다 파일에서 읽는 단계는 없다.

HTTP presenter는 domain `category`를 `sourceType`으로 표현하고 Trace entry를 requirement ID 목록으로 연결한다. 반복되는 Trace의 requirement text/Match/source와 domain `summary`는 HTTP에서 생략한다. Core evidence의 ID/text/relation, missing/unknown, reasoning 및 explanation provenance는 보존한다. 자세한 wire shape는 [API 문서](../api/job-fit.md#response-field-description)에 있다.

## CLI Flow

```text
--profile file → LoadSnapshot
--job text file → LoadJob / PrepareJob
TYPESAFE_API_KEY → Jev client
        ↓
application.WithSnapshot → Service.Evaluate(job.Text)
        ↓
same evaluator.Evaluate / same Jev adapter
        ↓
output.Result → domain JSON stdout
failure       → stderr + nonzero exit
```

CLI는 매 실행 시 지정한 profile을 읽고 HTTP는 startup 때 선택한 snapshot을 재사용한다. 같은 profile/JD의 evaluation 단계는 동일하다. CLI 출력에는 기존 summary와 상세 domain Trace가 있고 HTTP common envelope/evaluation metadata는 없다. 이는 presentation 차이다.

CLI context는 interrupt signal을 받는다. HTTP 전용 5분 전체 evaluation deadline은 CLI에 추가하지 않았으며 provider 호출별 60초 timeout은 공통이다. 따라서 두 interface의 operational 설정까지 동일하다고 주장하지 않는다. [CLI regression](../../cmd/job-fit/main_test.go)과 [service/core parity test](../../internal/application/service_test.go)가 이 경계를 다룬다.

## Evaluation / Jev Boundary

Jev는 외부 structured decision component다. 현재 adapter는 configured model `jev-1.13.0`을 사용하여 `https://api.typesafe.ai/v1/systemone`으로 요청한다. 이는 [현재 client 구현](../../internal/jev/client.go)의 설정이며 외부 서비스의 최신 지원 모델 목록에 대한 주장이 아니다.

| 순서 | Local preparation / validation | Jev 판단 |
| --- | --- | --- |
| 1. Extract | Non-empty line 순서로 `L001` 등 ID, 원문 line 보존. Ignore 제외 | required/preferred/responsibility/stack/context/ignore |
| 2. Evidence relation (Assess 내부) | Requirement × profile evidence 조합 구성, ID/text는 profile에서 복사 | direct/transferable/limited/unrelated |
| 3. Match/rationale (Assess 내부) | 실제 연결 결과를 질문에 반영, direct 없으면 Strong 제외, 조건 없는 version/duration rationale 제외 | 허용된 Match와 rationale 조합 중 선택 |
| 4. Overall | 전체 실제 assessments, profile, JD, policy를 state로 전달. Source category 유지 | Strong/Partial/Weak/Unknown 중 summary 판단 |
| 5. Trace | 이미 선택된 Overall과 전체 assessment를 전달. 새로운 Match를 계산하지 않음 | Supporting/Limiting/Non-decisive attribution |

정상적인 전체 경로는 provider HTTP 호출 5회다. 연결 evidence가 없는 항목은 application이 Unknown으로 유지하고 해당 Match 질문을 생략한다. 모든 항목이 그런 경우 Match 호출 전체가 생략되어 4회가 된다. 추출 대상 없음/오류에서는 더 일찍 종료할 수 있다. 자동 retry/redirect는 하지 않는다.

Provider answers는 질문 ID, answer type, 선택지 유효성을 검사한다. Adapter는 probability/confidence를 Match Score로 내보내지 않는다. Overall에 평균·다수결·가중치 계산은 없고 Jev가 정해진 기준에서 직접 선택한다. 설명은 선택된 rationale/criterion에 기반한 application-rendered 문장이다. 자유 텍스트 LLM 설명 또는 전체 내부 추론을 얻는 구조가 아니다. State/target은 untrusted data로 다루도록 instructions에 명시하지만 이것이 모든 의미적 오류를 막는 보장은 아니다.

### Evaluation concepts

| 개념 | 의미 |
| --- | --- |
| Strong | 핵심 requirement를 직접 지지하는 evidence와 작은 role/context 차이 |
| Partial | 관련·transferable evidence는 있으나 일부 scope/depth/context gap 존재 |
| Weak | 관련 evidence는 있으나 요구 수준 또는 핵심 role과 gap이 큼 |
| Unknown | 판단할 evidence가 부족함. No Evidence ≠ No Experience |
| Supporting Evidence | Profile의 ID/원문과 판정된 relation. 새로운 경력 사실이 아님 |
| Missing Evidence / Unknowns | 현재 문서로 입증되지 않거나 확인할 수 없는 정보. 실제 경험 부재와 구분 |
| Overall Match | Summary signal. 객관적 ground truth 또는 hiring probability가 아님 |
| Decision Trace | 선택된 Overall에 대한 requirement별 사후 attribution |

전체 정의는 [Evaluation Contract](../evaluation-contract.md)와 [Match Semantics](../match-semantics.md)를 참조한다. 현재 줄 기반 추출, 유한 rationale의 제한, language/education 제외 및 Composite Public / Energy Context ambiguity는 그대로다.

## Decision Trace Position / Limitation

```text
Requirement assessments → Overall Match (selected)
                                ↓
                         separate Trace Choice
                                ↓
                    ID resolution / validation
```

Supporting은 Overall을 높은 방향으로 지지하는 주요 근거, Limiting은 더 높은 Overall을 제한하는 gap/uncertainty, Non-decisive는 핵심 결정 요인이 아닌 항목이다. Source type이나 Match 개수만으로 기계적으로 분류하지 않는다.

[ResolveTrace](../../internal/evaluator/trace.go)는 missing/duplicate/unknown references, invalid influence, Unknown/evidence-free Supporting, Strong Overall의 Limiting을 거부한다. Strong일 때 adapter Choice에서도 Limiting을 제외한다. Validation failure는 Overall이나 requirement Match를 보정해 해결하지 않는다.

Trace는 **structured post-hoc attribution**이다. 최초 Overall의 causal trace가 아니고 Overall로 feedback되지 않는다. 별도 호출의 비용/실패/변동 가능성이 있으며 implicit weighting을 입증하지 않는다. ID는 동일 JD line 순서에서 연결 가능하지만 JD 편집 전후 영구 ID나 persisted evaluation ID는 아니다.

## Profile Boundary

**Submission Resume ≠ Evaluation Profile.** 특정 공고에 맞춘 resume를 직접 평가 입력으로 쓰면 문서 tailoring이 결과에 영향을 줄 수 있으므로 별도의 factual canonical evidence source를 사용한다.

[LoadProfileSnapshot](../../internal/evaluator/input.go)은 한 번 읽은 raw bytes를 기존 profile parser로 검증하고 SHA-256을 계산한다. Service의 opaque snapshot은 private하게 유지하고 evidence slice를 요청별로 복사한다. Hot reload는 없으며 HTTP profile 갱신은 restart 후 반영한다. 서버가 살아 있는 동안 파일이 바뀌어도 기존 snapshot/version을 유지한다.

`profileVersion`은 `sha256:<64 lowercase hex>`이며 formatting-only 변경도 달라질 수 있다. Profile 조회 ID나 semantic version이 아니다. CLI는 `--profile`로 로컬 파일을 명시하지만 REST caller는 profile/ID/version을 요청에 넣을 수 없다. Future generic candidate input은 현재 scope 밖이다.

## Validation / Error / Cancellation Boundaries

| Boundary | 확인하는 내용 |
| --- | --- |
| Startup / input | Canonical profile shape/unique evidence, JSON/type/duplicate/Unicode, URL, byte/line limits |
| Jev client / adapter | Provider wire shape/allowed choices, request/response sizes, evidence/rationale constraints |
| Evaluator | Assessment count/identity/Match, Trace 참조와 기본 consistency |
| HTTP presenter | Public field/source/relation 및 Trace entry/order consistency, array/DTO mapping |

Body 128 KiB, description 16 KiB/64줄, profile evidence 1–32, pairs 512 및 generated provider request 192 KiB 제한은 resource safeguards다. 평가 threshold가 아니다. Pair/payload 제한은 일부 external 작업 이후 발견될 수 있다. Structural validation은 semantic 정확성 전체를 증명하지 않는다.

```text
Jev transport / provider failure
        ↓
failure.Error category (context errors retain identity)
        ↓
Evaluator wrapping → Application service
  (uncategorized evaluation errors → Evaluation failure)
        ↓
HTTP mapError → status + public error code/message
        ↓
Common error envelope (no partial data)
```

HTTP request validation은 provider를 거치지 않고 직접 error envelope를 만든다. Application 내부의 모든 오류가 원래부터 typed였던 것은 아니며, service가 기존 evaluator/provider 검증 오류를 evaluation failure로 보완한다. HTTP mapper는 raw string matching을 하지 않는다.

| HTTP | 경계 / 대표 원인 |
| --- | --- |
| 400 | JSON/구조/타입/unknown·duplicate fields/query |
| 413 | Body/field bytes 또는 line 제한 |
| 415 | Unsupported/missing media type 또는 encoding |
| 422 | 의미 검증, no capability requirements, pair/payload capacity |
| 500 | Configuration(포함: upstream 401/403) 또는 unexpected internal failure |
| 502 | Provider network/response/evaluation/consistency failure |
| 503 | Upstream busy/rate limit (429/529) |
| 504 | Evaluation deadline/provider timeout |

404/405 routing 처리도 HTTP 경계에 있다. 상세 code/응답은 [Error Summary](../api/specification.md#error-summary)를 참조한다. 실패 응답에는 data가 없고 raw upstream body/secret을 표시하지 않는다. CLI는 HTTP status 대신 stderr/nonzero exit를 사용한다.

```text
HTTP request context (+ requestId)
    → WithTimeout (default 5 minutes)
    → Application Service
    → evaluator stage checks
    → Jev http.NewRequestWithContext
```

Provider request별 60초 제한은 별도로 유지한다. `--evaluation-timeout`은 positive~1h로 설정 가능하다. Caller disconnect 시 이후 작업을 중단하지만 response/499 또는 remote 처리·과금 중단을 보장하지 않는다. 서버 SIGINT/SIGTERM은 BaseContext 취소와 shutdown으로 이어진다. Server read-header/read/write/idle/shutdown 값은 [API Timeout](../api/job-fit.md#timeout)에 있다.

## Runtime / Security Boundary

- 기본 `127.0.0.1:8080`, numeric loopback bind만 허용. Auth/authorization/public deployment는 구현하지 않았다.
- `TYPESAFE_API_KEY`는 실행 프로세스 환경변수다. Server-side canonical profile과 key를 HTTP 요청으로 받지 않는다. `.env` 자동 로딩은 없다.
- Profile/JD/key/Authorization/raw upstream 오류를 HTTP 로그에 직접 출력하지 않는다. RequestId, method, matched route 또는 unmatched, status/error category, outcome, elapsed를 기록한다.
- UUID v4 requestId는 server가 생성해 context와 응답/로그에 전달한다. Distributed tracing SDK나 history/조회 ID는 아니다.
- Jev API는 **system 밖의 network dependency**다. Profile/JD는 evaluation 단계에서 Jev로 전송되고 결과에는 연결된 candidate evidence가 포함된다. Local server라는 사실이 local-only inference를 뜻하지 않는다.
- 자동 retry, persistence/recovery, queue, request concurrency/cost controls는 없다. 통신 오류와 내부 오류는 분류하지만 외부 서비스 availability나 장기 운영 안정성을 보장하지 않는다.

[go.mod](../../go.mod)는 Go 1.24.0을 지정하며 third-party module dependency가 없다. HTTP/JSON/context/hash/logging은 Go standard library를 사용한다.

## Project Structure

주요 추적 파일/디렉터리만 표시했다. 생성된 binary, local secret 파일 및 test 파일 개별 목록은 생략했다. Tests는 각 package의 `*_test.go`에 있다.

```text
job-fit/
├── go.mod
├── README.md
├── cmd/
│   ├── job-fit/          # CLI entry point
│   └── job-fit-server/   # HTTP process startup/lifecycle
├── internal/
│   ├── application/      # Shared use case, snapshot, metadata
│   ├── evaluator/        # Domain, Engine, input/trace validation, stage order
│   ├── failure/          # Transport-neutral typed error categories
│   ├── httpapi/          # Request/DTO/envelope/error mapping
│   └── jev/              # Provider client, questions, mapping, qualification checks
├── data/
│   └── profile.json      # Canonical evidence
├── samples/              # 6 prepared JD text samples
└── docs/
    ├── api-contract.md   # Public contract and design rationale
    ├── evaluation-contract.md
    ├── match-semantics.md
    ├── api/              # API overview and endpoint usage
    ├── architecture/
    │   └── overview.md   # This document
    └── development/      # Turn history and verification JSON artifacts
```

`internal/evaluator/instructions.md`는 binary에 embed되는 evaluation policy다. Development records/expected labels는 runtime에 읽거나 Jev에 전송하지 않는다. `docs/development/*-results`는 검증 artifact이며 runtime result store가 아니다.

## Architectural Decisions

| 결정 | 현재 구조에 미친 영향 / 근거 |
| --- | --- |
| Go + CLI first | 작은 HTTP/JSON tool을 standard library로 검증하고 단일 binary로 사용. [Turn 004](../development/turn-004.md) |
| Core / Engine / Jev adapter 분리 | Provider wire 형식과 HTTP 인증을 core 밖에 두고 fake Engine/transport로 offline 검증. [Turn 004](../development/turn-004.md), [pipeline test](../../internal/jev/pipeline_test.go) |
| Structured result와 관측 우선 | 자유 생성 설명 대신 evidence IDs/Choice/rationale 및 post-hoc Trace. [Turn 006](../development/turn-006.md) |
| Human Reference는 ground truth 아님 | Label tuning보다 evidence/requirement fidelity와 Unknown/consistency를 검토. [Turn 007](../development/turn-007.md) |
| Confirmed violation만 수정 후 baseline 수용 | Strong+Limiting 및 synthetic qualification의 제한된 수정, 6개 JD regression 후 tuning 종료. [Turn 010](../development/turn-010.md), [Turn 011](../development/turn-011.md) |
| REST reusable boundary + shared service | 특정 agent protocol에 core를 묶지 않으며 CLI workflow를 재사용. [Turn 012](../development/turn-012.md)–[014](../development/turn-014.md) |
| Canonical server-side snapshot + DTO | Personal profile consistency, file-byte fingerprint, HTTP와 domain presentation 분리. [Turn 013](../development/turn-013.md), [snapshot/parity tests](../../internal/application/service_test.go) |
| Portfolio feature freeze | Integration/generalization을 완료 조건에서 제외하고 문서화·최종 검증으로 전환. [Turn 015](../development/turn-015.md) |

이는 이전 결정의 요약이며 새로운 architecture 선택이 아니다. [HTTP tests](../../internal/httpapi/handler_test.go)와 [cancellation boundary tests](../../internal/jev/http_boundary_test.go)는 DTO/error/cancellation 경계를 검증한다. 이번 문서 Turn에서는 테스트나 실제 API를 재실행하지 않았다.

## Out of Scope / Extension Points

```text
External discovery / company analysis / filtering / JD acquisition
                              |
                       provided JD
                              v
                   [ current job-fit boundary ]
                              |
                       structured result
                              v
                 External storage / human review
```

현재 system 밖: Company Search, Position Discovery, SES/派遣/客先常駐 classification, JD crawling/URL fetching, OpenClaw orchestration, Notion, MCP, DB/history, Web UI, batch, auth/multi-user/public deployment, score/probability/recommendation. 이들을 job-fit이 호출하거나 운영한다고 암시하지 않는다.

가능한 확장 지점은 다음 정도이며 구현 계획/약속 또는 이번 MVP 완료 조건이 아니다.

- **Generic candidate input:** 별도로 검증된 profile + JD를 같은 evaluator에 전달할 가능성. 현재 endpoint는 이를 거부하며 schema/privacy/범위 검토가 먼저 필요하다.
- **Thin MCP adapter:** 필요하면 기존 REST를 호출하는 외부 adapter. 현재 code나 tool은 없다. Adapter에 evaluation policy를 복제하지 않는 방향이다.
- **External automation:** OpenClaw 등이 확보한 JD로 REST를 호출하고 결과를 활용할 가능성. Search/filtering/Notion 책임은 외부에 남는다.

향후 변경은 별도 범위 결정이 필요하다. 남은 현재 portfolio 작업은 README polish, AI-assisted development story 및 final regression이다.

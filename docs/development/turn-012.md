# Turn 012 — Product Boundary & Interface Direction

- 날짜: 2026-09-21
- 범위: 제품 책임 및 외부 interface 방향 문서화. Application code 변경·API 호출 없음.
- 기반: Turn 011의 **Evaluation Baseline Accepted**. 현재 evaluator를 v0 baseline으로 유지한다. 이는 새 API 버전이나 release tag 결정은 아니다.

## Completed

사용자가 제시한 실제 job-search workflow에 job-fit이 들어갈 위치를 정리했다. 기존 CLI/core 경계를 확인하고, evaluator에 검색·회사 분석·저장 책임을 추가하지 않는 방향을 기록했다. REST를 primary integration interface 후보로 채택하고 CLI 유지·MCP optional adapter 역할을 구분했다. 구체 HTTP 계약과 구현은 다음 Turn으로 남겼다.

## Existing Workflow

아래는 사용자가 설명한 현재 운영 workflow다. OpenClaw/Notion 연결이나 cron 설정을 이번에 조회·검증·변경한 것은 아니다.

```text
OpenClaw cron / search workflow
  → 회사 또는 포지션 후보 탐색
  → 홈페이지·회사 정보 조사
  → LLM 기반 회사 성격 분석
      自社開発 / SIer / SES / 派遣 / 客先常駐 등
  → 지원 대상이 아닌 회사 hard filtering
  → 채용 페이지 및 모집 Position / JD 확인·확보
  → job-fit evaluation
  → Notion enrichment / 조사 결과 저장
  → Human review
```

회사/포지션 탐색과 회사 filtering의 구체 순서는 외부 workflow가 결정한다. job-fit의 시작점은 **평가할 JD가 이미 확보된 시점**이다. 기존에 LLM과 사람이 수행하던 적합성 검토에 구조화된 평가 component를 넣는 것이며, 최종 사람 검토를 대체하지 않는다.

## Product Boundary

핵심 책임:

> 이미 확보된 Job Description을 canonical candidate profile과 비교하여 structured Job Fit evaluation을 반환한다.

```text
Provided JD + Canonical candidate evidence
                   ↓
             Evaluator core
                   ↓
       Structured evaluation JSON
```

| 책임 주체 | 맡는 일 |
| --- | --- |
| Orchestration layer (OpenClaw 등) | 탐색 대상, 회사 조사·제외, JD acquisition, 평가 호출 시점, 저장 위치/Notion enrichment 결정 |
| job-fit | 제공된 JD requirement와 기록된 candidate evidence의 capability/experience match, missing/unknown 및 structured result 반환 |
| Human | 결과 해석, 추가 확인 및 최종 지원 여부 판단 |

회사 분류/고용형태 filtering을 job-fit의 evaluation policy로 가져오지 않는다. Source URL이 제공된다고 job-fit이 해당 사이트를 가져오거나 크롤링할 책임을 갖는 것도 아니다. Context를 평가에 사용하는 기존 semantics와 회사 성격 조사·hard filtering의 책임은 구분한다.

## Current Implementation / Interface Direction

현재 [CLI](../../cmd/job-fit/main.go)는 `--profile`/`--job` 파일을 읽고 Jev client를 구성한 뒤 [evaluator.Evaluate](../../internal/evaluator/evaluator.go)를 호출하여 domain JSON을 출력한다. Core는 Engine interface를 통해 Jev adapter를 사용하며 CLI의 flag/stdout 또는 HTTP transport를 직접 처리하지 않는다. 여기서 core의 orchestration은 **평가 단계 실행 순서**이며 회사 검색 등 외부 job-search orchestration과는 다른 책임이다.

현재 HTTP server나 MCP server는 없다. 앞으로의 방향은 다음과 같다.

```text
OpenClaw / Codex / scripts / curl / optional MCP adapter
                          ↓
                 REST API (후보, 미구현)
                          ↓
Local CLI ────────→ Shared evaluator core
                          ↓
                         Jev
```

CLI를 반드시 REST 호출 client로 바꾼다는 결정은 아니다. 현재 직접 core 호출을 유지할 수 있고, REST와 CLI가 같은 core를 재사용하는 것이 핵심이다. MCP가 필요해지면 REST를 호출하는 얇은 adapter로 고려한다. MCP protocol이나 client별 동작에 core를 종속시키지 않고 adapter에 evaluation logic을 넣지 않는다.

## Why REST API

OpenClaw, Codex, CLI, conversational MCP client, 일반 script/curl 및 향후 automation처럼 caller가 다양하다. 공통 structured request/response를 사용하는 HTTP 경계는 caller와 evaluator를 분리하고 특정 AI client protocol에 대한 종속을 줄일 수 있다. 현재 Go evaluator를 공통 backend 기능으로 재사용할 방향과도 맞는다.

이는 요구된 usage workflow에 따른 architecture 판단이다. REST가 MCP보다 본질적으로 우월하다거나 현재 OpenClaw에 연결 검증을 마쳤다는 결론은 아니다. 서비스 형태에는 lifecycle·설정·오류·timeout 등의 계약이 필요하며, 이를 이번에 임의 확정하지 않는다.

## Role of CLI / REST / MCP

| Interface | 역할 | 현재 상태 |
| --- | --- | --- |
| CLI | 개발·테스트·로컬 수동 사용 | 구현됨. 명시 profile/JD 파일 입력 유지 |
| REST API | automation 및 다양한 외부 caller의 integration | Primary integration interface 후보로 채택. 계약·구현 미정 |
| MCP | conversational AI client용 얇은 adapter | 향후 필요 시 검토. REST 호출 방향, 구현 없음 |

제품 경계는 프로토콜 이름이 아니라 **JD 단위 요청 → structured evaluation result**다. API request가 내부 `Job` 또는 현재 CLI JSON과 완전히 같은 schema여야 하는지까지 결정한 것은 아니다.

## Candidate Profile Handling

초기 사용자는 한 명이고 canonical profile이 있으므로, API에서는 매 요청마다 profile 전체를 보내기보다 **server-side canonical profile 사용**을 우선 검토한다. 입력 개념은 확보된 JD와 회사명·직무명·source URL 같은 metadata, 평가에 필요한 candidate evidence는 server 쪽 canonical source라는 방향이다.

이는 기존 profile 파일·내용 또는 CLI의 `--profile` 동작을 변경하지 않는다. Server가 어느 파일을 언제 읽고 갱신하는지, 결과에 profile 식별/버전을 표시할지 등은 다음 contract 설계의 질문이다. DB 도입이나 새 profile 관리 API를 뜻하지 않는다. Multi-user, profile ID, 요청별 explicit profile은 미래 확장 가능성일 뿐 현재 MVP 요구로 확정하지 않는다.

## Decisions Made

이번에 확정한 방향:

- job-fit은 crawler나 job-search orchestration system이 아닌 JD 단위 reusable evaluation component다.
- Caller는 OpenClaw/MCP client에 한정되지 않는다.
- REST API를 primary integration interface **후보**로 채택한다.
- CLI를 유지하고 CLI/REST에서 동일 evaluator core를 재사용한다.
- MCP는 향후 optional adapter이며 evaluation logic을 넣지 않는다.
- Evaluation은 transport/interface 및 외부 orchestration과 분리한다.

Server-side canonical profile은 초기 API의 검토 방향이다. 구체 endpoint, request/response schema, framework, server implementation은 확정하지 않았다. 사용자 prompt의 `POST /job-fit/evaluate` 및 JSON은 개념 예시이며 구현할 확정 계약으로 승격하지 않았다.

## Open Questions

다음 HTTP API Contract Turn에서 결정할 사항:

- Endpoint와 HTTP method, API versioning을 어떻게 표현할 것인가?
- JD 본문과 company/title/source URL metadata의 필수·선택 범위, validation 및 evaluator 입력 변환은 무엇인가? Metadata 결합이 기존 줄 기반 ID/평가 의미에 미치는 영향도 검토한다.
- 현재 domain JSON과 API response의 관계, metadata/provenance 및 식별 정보는 무엇인가?
- Canonical profile의 선택·로딩·갱신 시점과 오류 처리를 어떻게 정할 것인가?
- 입력 오류, Jev 실패, cancellation/timeout을 caller에게 어떤 구조로 전달할 것인가? 현재 여러 순차 호출을 하나의 요청 lifecycle로 어떻게 다룰 것인가?
- CLI/API가 공유하는 application boundary와 transport별 입력/오류 처리의 경계는 어디인가?
- Service의 사용·노출 환경에 맞는 접근 범위와 credential 배치는 어떻게 정할 것인가? Public deployment나 인증 체계를 이번에 전제하지 않는다.

이 질문들은 결정이나 구현이 아니며, 기존 evaluation semantics의 재설계나 tuning 안건도 아니다.

## Out of Scope

회사/Web 검색, 회사 성격 조사, SES/派遣 판별, 채용 페이지 crawling, JD discovery, Notion 저장 및 최종 지원 판단은 job-fit의 책임 밖이다. HTTP server/framework 선택, MCP, OpenClaw/Notion integration, DB, crawler, evaluator/prompt/profile 변경은 이번 Turn에서 하지 않았다.

## Changed Files / Validation

- [README](../../README.md): accepted v0 상태와 현재 제품/interface 방향 반영. 완료된 전체 sample 평가를 미래 작업처럼 설명하던 문구와 “다른 5개 미실행” 문구를 갱신했다.
- [Evaluation Contract](../evaluation-contract.md#product-boundary--interface-direction--turn-012): 기존 평가 원칙을 유지한 채 제품 책임·interface 방향의 현재 결정 추가.
- 본 `docs/development/turn-012.md` 생성.

과거 Turn 000–011 및 accepted baseline의 code/profile/sample/policy/semantics/result artifact는 작업 전후 hash로 보존 확인했다. 문서 링크·공백·코드블록을 검증했다. 문서만 변경했으므로 Go regression이나 Jev 실행은 반복하지 않았다. Commit은 생성하지 않았다.

## Suggested Next Step

이 경계를 기준으로 **HTTP API Contract**를 설계한다. Endpoint, request/response, canonical profile, metadata, errors, timeout/external Jev failure, API versioning과 CLI/API 공유 application boundary를 검토한 뒤 구현한다. 이번 Turn에서는 server나 integration 구현을 시작하지 않았다.

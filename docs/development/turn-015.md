# Turn 015 — Final MVP Scope Freeze

- 날짜: 2026-09-21
- 범위: 최종 portfolio MVP 제품 범위와 남은 문서 작업 확정. Application code/API behavior 변경 없음.
- 시작 HEAD: `c585cec` (`feat: expose job-fit evaluator over HTTP`).
- 시작 상태: 이전 Turn의 README, Evaluation Contract, CLI/core/Jev 변경과 server/application/httpapi 등 untracked 파일이 존재했다. 기존 작업은 그대로 보존했다.

## Completed

Turn 000–014의 개발 흐름과 현재 repository 상태를 확인하고, 이번 사용자 prompt에 따라 final portfolio MVP scope를 확정했다. 새 문서 계층 대신 본 기록과 README의 active scope 요약을 사용한다. Current MVP, 명시적 제외 항목, 별도 Future Work 및 남은 portfolio 작업을 구분했다. 기능 구현이나 API 재실행은 하지 않았다.

## Project Origin

Jev라는 structured decision tool을 발견한 뒤, 여러 회사의 Job Description을 검토할 때 자신의 실제 경력과 요구 경험의 부합 정도를 구조적으로 확인하는 개인 작업에 적용하며 시작했다. 개인 실험과 도구가 구현 과정에서 AI-assisted development portfolio로 발전했다. 채용 성공 확률을 예측하거나 지원 여부를 대신 결정하는 시스템은 목표가 아니다. 이전 disposable PoC 코드는 재사용하지 않았다.

## Final MVP Definition

> **Server-side canonical candidate profile과 제공된 Job Description을 비교하여 Jev 기반의 structured Job Fit evaluation을 반환하는 Go application.**

```text
Provided Job Description + Canonical Candidate Profile
                         ↓
                 Application Service
                         ↓
                     Evaluator
                         ↓
                        Jev
                         ↓
            Validated Structured Evaluation
                         ↓
                  CLI / REST API
```

REST는 primary external integration interface이고 CLI는 로컬 개발·테스트·수동 사용 interface다. 둘은 동일 service와 evaluator workflow를 공유한다. HTTP의 profile은 server-side snapshot이며 CLI의 기존 explicit `--profile` 선택은 유지한다. 정의를 위해 CLI 동작까지 server 호출 방식으로 바꾸지 않는다.

## Repository / Development Review

| Turn | 확인한 흐름 |
| --- | --- |
| [000](turn-000.md)–[001](turn-001.md) | Bootstrap, 문제 정의와 personal Software Engineering capability/experience positioning |
| [002](turn-002.md)–[003](turn-003.md) | 6개 JD의 사전 해석 및 Required/Preferred, core role, Weak/Unknown, transferable semantics |
| [004](turn-004.md)–[006](turn-006.md) | Go/Jev CLI 첫 slice, Prop Tech 진단, Overall observability를 위한 post-hoc Trace |
| [007](turn-007.md)–[008](turn-008.md) | Human Reference는 ground truth가 아니라는 철학, frozen evaluation 관찰 |
| [009](turn-009.md)–[010](turn-010.md) | 세 후보 triage와 두 confirmed violation의 제한된 수정 |
| [011](turn-011.md) | 6개 sample 최종 regression, Evaluation Baseline Accepted, iterative tuning 종료 |
| [012](turn-012.md)–[013](turn-013.md) | Orchestration/evaluation 분리, REST 방향 및 공통 표준 기반 HTTP contract |
| [014](turn-014.md) | 공유 service, Go HTTP server, snapshot/version, errors/cancellation, 실제 HTTP/Jev 성공 |

기존 기록과 저장된 검증 결과를 참조한 것이며, 이번 Turn에서 Go regression이나 실제 API 검증을 다시 수행한 것은 아니다.

## Included Scope

| 영역 | 포함 내용 / 경계 |
| --- | --- |
| Canonical candidate profile | 회사별 tailored resume와 분리한 factual evidence source 재사용. HTTP startup snapshot/profileVersion, CLI file input. Profile에 없는 경험을 만들어내지 않는 원칙 |
| Job Fit evaluator | 현재 줄 기반 JD requirement 분석, candidate evidence 연결, Strong/Partial/Weak/Unknown, supporting/missing evidence와 Unknown preservation |
| Overall / Decision Trace | Overall은 summary signal. Requirement identity와 Supporting/Limiting/Non-decisive, evidence 추적 및 contract validation. Trace는 post-hoc attribution이며 causal reasoning 기록이 아님 |
| Jev integration | Structured decision component. 일반 LLM 대비 객관적 label accuracy 우위를 주장하지 않음 |
| CLI | 로컬 개발·테스트·수동 평가. REST와 동일 application service/evaluator 사용, 별도 평가 정책 없음 |
| REST API | `POST /v1/job-fit/evaluate`, 확보된 JD 입력, server-side canonical profile 사용. HTTP validation/presentation은 담당하지만 evaluation logic은 core에 유지 |
| API contract / errors | Common envelope, requestId, profileVersion, typed error/status mapping, body/description limits, timeout/cancellation 전파, secret-safe logging |
| Tests / evaluation | Unit/handler/regression/contract tests, go test/race/vet/build, 6개 실제 JD regression 및 실제 Jev HTTP 연동 검증의 기록 |

Canonical evidence는 JD별 resume tailoring에 따른 평가 bias를 피하기 위한 기준이다. 구조화된 결과를 제공한다고 해서 semantic error나 evidence hallucination 가능성을 완전히 제거했다고 주장하지 않는다. `Match ≠ Hiring Probability`이며 최종 판단은 사람이 한다.

Jev 통합의 탐구 질문은 다음과 같다.

> Structured AI decision을 reusable software workflow 내부의 constrained and observable component로 사용할 수 있는가?

품질 기준은 Evidence Faithfulness, Requirement Fidelity, Unknown Preservation, Evidence Traceability, Structural Consistency, Repeatability / Variability observation, User Usefulness다. Human Reference와의 label 일치율은 품질 목표가 아니다. Variability 관찰은 통계적 repeatability 실험이 이미 완료되었다는 뜻이 아니다.

## Explicitly Out of Scope

| 영역 | 이번 portfolio MVP에서 구현하지 않을 항목 |
| --- | --- |
| Orchestration / discovery | OpenClaw integration, company search, position discovery, company type analysis, SES/派遣/客先常駐 classification, JD crawler, URL fetching, Web scraping |
| External workflow | Notion integration, cron workflow, automated job discovery pipeline |
| AI client adapters | MCP server, Claude-specific integration, Codex-specific integration |
| Product features | Web UI, database, evaluation history, batch evaluation, authentication, authorization, multi-user accounts, public deployment |
| Generalization | Caller-provided candidate profile endpoint 및 그에 필요한 candidate schema generalization |
| Scoring / decisions | Match percentage, hiring probability, offer probability, application recommendation, numeric weighting/threshold calibration |

일반 script/curl 또는 AI agent가 기존 REST를 호출할 수 있다는 것과 특정 client adapter/integration을 구현했다는 것은 구분한다. 서버의 runtime history 기능은 없으며 개발 기록의 result JSON은 검증 artifact다. 기존 API sourceUrl은 provenance metadata이며 fetch 책임을 추가하지 않는다.

## Personal vs Generalized Direction

현재 primary user는 프로젝트 작성자다. HTTP caller는 이미 확보한 JD를 보내고 server-side canonical profile을 사용한다. 동일 evidence를 여러 JD에 재사용하며 profile override endpoint는 없다.

미래 범용 버전에서 candidate profile + JD를 caller가 모두 제공하는 형태는 고려할 수 있다. 그러나 summary/skills/experiences/roles 같은 예시 field는 승인된 schema가 아니며 현재 `/v1` request에 추가하지 않는다. Candidate schema generalization, 직군별 평가 차이, 개인정보 처리, validation 및 multi-user 제품 범위는 이번 MVP 핵심을 벗어나므로 별도 작업으로 남긴다.

## Future Work

다음 항목은 가능성으로만 남기며 **이번 portfolio 완료 조건이나 확정된 구현 일정이 아니다**.

| 별도 후속 작업 | 가능한 방향 |
| --- | --- |
| Generic Evaluation API | Caller-provided canonical evidence와 JD를 받는 범용 endpoint. 새 schema/범위를 먼저 검토 |
| OpenClaw integration | Position discovery → company filtering → JD acquisition 이후 REST 호출, 결과를 Notion/human review용 enrichment로 사용 |
| MCP adapter | 필요할 때 기존 REST를 호출하는 optional thin adapter. MCP 계층에 evaluation logic 없음 |
| Automation | 다수 JD batch 평가 또는 job-search pipeline 연결 |

Turn 014의 다음 단계였던 OpenClaw integration boundary 검토는 **현재 MVP 작업 순서에서 분리**한다. 기존 Turn 014 기록은 당시 제안으로 보존하며, 이번 결정 이후의 우선순위는 문서화와 portfolio polish다.

## Remaining Portfolio Work

제품 기능 범위는 확정했다. 남은 문서화와 final regression까지 완료했다고 선언하는 것은 아니다. Evaluation subsystem은 accepted v0 baseline을 유지하고 iterative tuning 대상으로 다시 열지 않는다. Final regression도 Human Reference label에 맞추는 작업이 아니다.

| 예상 Turn | 작업 | 범위 |
| --- | --- | --- |
| 016 | API Documentation | MyLib 스타일의 API 사용자용 reference와 endpoint 문서. 실제 스타일 자료를 확인한 뒤 작성 |
| 017 | Architecture / Project Documentation | CLI/REST/application/evaluator/Jev 역할·dependency boundary 정리 |
| 018 | Portfolio README / AI-assisted Development Story / Final Regression | 문제·선택·검증·한계와 AI-assisted 과정 정리, 현 기능의 최종 회귀 확인 |

Turn 번호는 planning에서 조정할 수 있다. 이번 Turn에서 위 문서를 미리 만들거나 구현·regression을 시작하지 않았다.

## Documentation Direction

| 문서 | 역할 |
| --- | --- |
| README | 프로젝트 소개, 빠른 시작, architecture summary, documentation links |
| API Contract | 왜 이 계약을 선택했는지와 public design/contract 설명 |
| API Specification | 실제 API 사용자의 빠른 reference |
| Endpoint Documentation | Request/response, field/error 및 curl examples |
| Architecture | CLI/REST/application/evaluator/Jev 관계와 dependency boundary |
| Development History | Turn 000부터의 AI-assisted 과정, 판단·수정·검증 기록 |

API Specification/Endpoint Documentation은 현재 구현 완료 산출물로 표시하지 않는다. 구체 파일 배치와 MyLib 형식은 다음 문서화 Turn에서 다룬다. 기존 공통 응답 표준을 재설계하지 않는다.

## Decisions Made

- 이번 portfolio의 제품 기능 범위를 현재 canonical-profile evaluator + CLI/REST로 동결한다.
- REST를 후보 표현에서 현재 primary external integration interface로 명확히 한다.
- 이전 README의 optional Match Score 고려 문구를 이번 MVP 제외 범위로 갱신한다. 과거 Turn 기록은 수정하지 않는다.
- OpenClaw/MCP/generic candidate API/batch는 별도 Future Work이며 MVP 완료 조건이 아니다.
- Evaluator tuning 종료를 유지하고 이후 작업을 documentation/portfolio polish/final regression에 집중한다.
- 별도 scope 문서 계층을 추가하지 않고 본 Turn 기록과 README 요약으로 범위를 관리한다.

## Problems / Findings

README의 “Match Score를 선택적으로 고려”와 “REST API 후보”는 final scope와 맞지 않는 과거 표현이었다. Active README만 갱신했다. 과거 Evaluation Contract의 검토 이력이나 이전 Turn의 다음 단계 제안을 구현 약속으로 승격하지 않는다.

Composite Public / Energy Context ambiguity, 모델 variability, 줄 단위 extraction의 한계, Trace의 post-hoc 성격은 그대로다. 기능 동결은 이런 제한이 해결됐다는 뜻이 아니다. MyLib 스타일의 구체 template은 이번 prompt에 포함되지 않았으며 이번에는 찾거나 임의로 정의하지 않았다. 다음 API documentation Turn에서 실제 reference를 확인한다.

## Changed Files / Validation

- [README](../../README.md): final scope와 origin 요약, included/out-of-scope/future work, 남은 문서 단계 및 현재 REST/scoring 표현 갱신.
- 본 `docs/development/turn-015.md`: 최종 범위 확정 기록.

작업 시작/종료 SHA-256 비교로 README 외 모든 기존 파일이 동일하며 신규 파일은 본 Turn 기록뿐임을 확인했다. 문서 링크/code fences와 `git diff --check` 검증이 통과했다. Code/API/semantics 변경이 없으므로 Go tests/vet/build와 Jev 실행은 반복하지 않았다. 기존 작업을 덮어쓰거나 commit하지 않았다.

## Suggested Next Step

다음 Turn에서 **MyLib 스타일의 API documentation**을 작성한다. 먼저 실제 MyLib reference를 확인하고 API Specification과 endpoint 설명의 역할을 구분한다. 현재 `/v1/job-fit/evaluate`와 공통 응답 표준을 문서화하며 새 endpoint, adapter 또는 integration 구현은 시작하지 않는다.

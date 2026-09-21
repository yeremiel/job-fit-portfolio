# AI-assisted Development: From a Personal Experiment to a Bounded MVP

job-fit은 Jev라는 structured decision tool을 개인적인 JD 검토에 적용해보며 시작했다. 여러 회사 공고와 같은 경력을 반복 비교할 때 근거와 미확인 정보를 일관된 형태로 보고 싶었다. 결과는 canonical evidence와 제공된 JD를 비교하는 작은 Go CLI/REST application이며, 채용 성공 예측이나 자동 지원 시스템은 아니다.

이 문서는 Turn 000–017의 전환점을 요약한다. 당시 판단과 검증의 원문은 [Development History](development/), 현재 구조는 [Architecture Overview](architecture/overview.md)에 있다. 기록되지 않은 인간의 검토나 모델 내부 추론을 추정하여 서술하지 않는다.

## Development Workflow

| 작업 공간 | 역할 |
| --- | --- |
| ChatGPT planning conversation | Brainstorming, 문제 정의, architecture discussion, Turn scope/acceptance criteria, 결과 검토와 다음 결정 |
| Codex workspace | Repository 분석, implementation planning, 코드·테스트·문서 작성, 실행 검증, code review/refactoring, regression, Turn Result |
| Human review | 문제와 제약 제공, 제안/결과 판단, scope와 정책 선택, 다음 Turn 진행 여부 결정 |

```text
Brainstorm → Scope → Task Prompt → AI Implementation
                                      ↓
                              Test / Review
                                      ↓
                               Human Review
                                      ↓
                                Next Decision
```

한 번의 큰 구현 요청 대신 Turn마다 목표·금지 범위·완료 조건을 지정했다. Usage limit으로 작업이 중단되었을 때에도 기존 workspace와 미완료 항목부터 확인하도록 했고, 결과를 planning conversation으로 가져가 다음 행동을 정했다. Turn 기록은 성공만 나열하는 일지가 아니라 가정·오류·보류 사항을 남기는 장치였다.

## 1. Disposable PoC → Greenfield Implementation

처음의 shell PoC는 Jev로 structured evaluation이 가능한지 확인하는 실험이었다. Portfolio implementation에는 기존 PoC 코드를 재사용하지 않고 새 repository 구조와 Go 구현으로 시작했다. Feasibility 확인용 코드를 곧바로 제품 architecture로 간주하지 않았다.

[Turn 000](development/turn-000.md)에서 application code 없이 시작 상태를 기록하고, [Turn 004](development/turn-004.md)에서 standard library 중심의 CLI 첫 end-to-end slice를 구현했다. 작은 Go CLI로 profile/JD → Jev → JSON 경로를 먼저 검증한 뒤 interface를 확장했다.

## 2. Evaluation Semantics Before Implementation

코드 전에 Required/Preferred, core role, transferable experience와 Strong/Partial/Weak/Unknown을 정의했다. 단순 기술 keyword가 겹친다고 frontend Tech Lead 같은 전문 역할 경험을 만들어내지 않으며, 기록이 없다는 사실을 무경험으로 바꾸지 않는 것이 중요했다.

Canonical profile도 회사별 제출 resume와 분리했다. JD에 맞춘 resume를 입력으로 쓰면 표현을 바꾼 효과가 적합성 평가에 섞일 수 있기 때문이다. 이것이 의미적 정확성을 보장하지는 않지만 무엇을 근거로 판단했는지 검토할 기준을 만들었다. [Turn 002](development/turn-002.md)–[003](development/turn-003.md)에 초기 해석과 규칙이 있다.

## 3. The Human Baseline Trap

Prop Tech plus의 사전 해석은 Strong이었지만 첫 Jev 결과는 Partial이었다. 이 차이를 줄이는 것을 개선 목표로 삼기 쉬웠다. 그러나 “누가 정답인가?”부터 불확실했다. 같은 evidence와 규칙에서도 사람·LLM·Jev가 scope와 context를 다르게 해석할 수 있었다.

개발 기록은 **실제로 reference label에 맞춰 tuning했다**고 말하지 않는다. [Turn 005](development/turn-005.md)는 차이를 진단했고, [Turn 007](development/turn-007.md)은 차이를 없애려는 방향으로 계속 진행하면 특정 해석에 overfit할 위험을 명시했다. 이 점을 calibration을 했다가 되돌린 이야기로 과장하지 않는다.

전환은 다음과 같다.

> Reference ≠ Ground Truth. Label disagreement alone ≠ Contract Violation.

“Human Baseline”도 독립적인 인간 정답 데이터가 아니었다. Codex 초안에서 출발해 사용자 지시로 확정한 **Pre-evaluation Reference Interpretation**이다. 비교와 관찰에 쓰되 target label로 사용하지 않는다.

품질 기준은 label 일치보다 Evidence Faithfulness, Requirement Fidelity, Unknown Preservation, Evidence Traceability, Structural Consistency, Repeatability / Variability observation, User Usefulness에 두었다. 이는 새로운 accuracy 수치나 모든 기준을 완전히 충족했다는 인증이 아니다.

## 4. Observability Before Calibration

“왜 Partial인가?”에 답하기 위해 먼저 관측 정보를 추가했다. [Turn 006](development/turn-006.md)의 Decision Trace는 requirement가 Supporting/Limiting/Non-decisive 중 어떤 역할로 분류되는지를 보여준다. Stable requirement ID로 Match와 evidence를 연결했다.

여기서 설명 가능성의 범위를 제한하는 것도 engineering decision이었다. Trace는 **이미 선택된 Overall에 대해 별도 Jev 호출이 만든 structured post-hoc attribution**이다. 최초 판단의 causal trace나 내부 추론을 추출한 것이 아니다. 추가 호출 비용과 변동/실패 가능성이 있으며, trace를 보고 암묵적 weighting이 입증됐다고 할 수 없다.

## 5. Contract Violation Triage

[Turn 008](development/turn-008.md)에서 여러 실제 JD의 결과를 관찰한 뒤, 모든 불일치를 고치지 않고 [Turn 009](development/turn-009.md)에서 세 후보를 기존 규칙과 비교했다. 선택 가능한 판정은 Confirmed Violation / Not a Violation / Ambiguous였으며 실제 세 후보는 두 Confirmed와 한 Ambiguous로 나뉘었다.

| 후보 | 판정 / 후속 처리 |
| --- | --- |
| Strong Overall + Limiting Trace | 최고 level인데 더 높은 Overall을 제한한다는 정의와 모순. Confirmed |
| Public/Energy 복합 Context 전체 Strong | 일부 domain evidence만으로 Context 정합성을 어떻게 해석할지 기존 규칙만으로 확정하기 어려움. Ambiguous로 유지 |
| JD에 없는 version/duration gap | 명시하지 않은 qualification을 부족한 근거처럼 생성. Confirmed |

[Turn 010](development/turn-010.md)은 두 Confirmed만 수정했다. Strong Overall의 Limiting 선택지를 제외하고 domain validation에서도 거부했다. Version/duration rationale은 해당 requirement에 명시된 조건이 있을 때만 선택·적용하도록 제한했다. Overall을 낮추거나 모든 항목을 Supporting으로 바꾸는 보정은 하지 않았다.

이 단계의 학습은 **결정하지 않을 항목을 남기는 것**이었다. Ambiguous는 정상 인증도 수정 약속도 아니다. Evidence가 허용하는 결론과 허용하지 않는 결론을 구분했다.

## 6. Freeze the Evaluator

[Turn 011](development/turn-011.md)에서 6개 JD를 수정된 evaluator로 각각 한 번 실행했다. 기록된 검토에서 새로운 Confirmed Contract Violation이 발견되지 않았고 두 regression이 재발하지 않아 baseline을 수용했다. 반복 실행으로 원하는 label을 얻지 않고 iterative tuning을 종료했다.

이는 제한된 sample에서의 contract integrity 확인이다. 일반화된 정확도, statistical repeatability 또는 모든 ambiguity 해결을 증명한 것은 아니다. 종료 기준을 명시함으로써 애매한 평가 품질을 끝없이 조정하는 작업에서 재사용 가능한 도구 완성으로 이동했다.

## 7. CLI → REST Product Boundary

초기에는 MCP tool을 생각했지만 실제 workflow를 검토하니 caller가 OpenClaw, Codex, 일반 script 등으로 다양했다. 제품 책임은 protocol이 아니라 **이미 확보된 JD + canonical evidence → structured result**였다.

[Turn 012](development/turn-012.md)–[014](development/turn-014.md)에서 REST를 reusable external boundary로 선택하고, CLI와 HTTP가 같은 application service/evaluator를 사용하도록 했다. HTTP DTO/common envelope, typed errors, timeout/cancellation, snapshot/version을 추가하면서 evaluation logic을 HTTP에 넣지 않았다.

Company discovery/filtering, crawling, Notion 저장과 최종 지원 결정은 system 밖에 남겼다. MCP는 REST를 호출할 수 있는 optional future adapter다. [Turn 015](development/turn-015.md)의 final scope freeze 이후 API/architecture 문서를 정리했고 추가 integration을 MVP의 미완료 항목으로 취급하지 않았다.

## Human / AI Contribution

AI가 상당 부분의 application code, tests, 문서를 작성하고 도구로 검증했다. 사람이 모든 코드를 직접 작성했다거나 자동 생성만으로 결과가 승인됐다고 설명하지 않는다. 정확한 기여 percentage는 측정하지 않았다.

| 영역 | Human role | AI role |
| --- | --- | --- |
| Problem / scope | 실제 문제·경력·workflow 제공, 범위와 non-goals 결정 | 질문·대안·문제 구조화 지원 |
| Evaluation philosophy | Ground truth 가정, Unknown, 종료 기준 등 최종 방향 선택 | 비교·비판·한계 및 대안 정리 |
| Architecture | 제품 경계와 interface 방향, scope 수용 | 현재 코드 분석, 설계 제안, trade-off 설명 |
| Implementation | 요구사항·acceptance criteria 제공, 결과 검토/수용 | Codex의 코드 작성·수정·refactoring |
| Tests / regression | 검증 범위와 성공 기준 결정 | 테스트 작성/실행, 실패 분석, 결과 기록 |
| Contract triage | 무엇을 수정 대상으로 삼을지 최종 판단과 다음 Turn 지시 | Source/result/code 근거 수집, 판정 분석과 문서 초안 |
| Documentation | Portfolio 방향·reference/style·공개 범위 결정 | 초안 작성·유지, 링크/예시 검증 |

Planning conversation의 모든 세부 대화를 이 repository가 담는 것은 아니다. 위 역할은 전달된 Turn prompts와 개발 기록에 기반한다. 모든 코드에 대한 인간의 독립적인 line-by-line review나 외부 전문가 검증을 주장하지 않는다.

## What Was Not Delegated as a Final Decision

프로젝트/MVP scope, evaluation philosophy, Reference를 ground truth로 삼지 않는 방향, 수정할 violation의 범위, tuning 종료, REST product boundary와 final portfolio 범위는 사람이 다음 Turn의 지시·수용 기준으로 결정했다. AI는 분석과 구현을 수행했지만 그 결정을 대신 승인하는 주체는 아니었다. 실제 입사 지원 여부 역시 도구 밖의 사람 판단이다.

## Limitations and Future Work

Objective Job Fit ground truth는 가정하지 않으며 provider 판단은 변동할 수 있다. Trace는 post-hoc이고, compound/context coverage ambiguity와 줄 기반 추출·유한 rationale의 한계가 남아 있다. 현재 profile은 한 개인의 Software Engineering 경험 중심이고 API는 loopback personal service다. Linked evidence는 개인 경력 정보이며 공개 범위 검토가 필요하다.

Generic candidate input, auth, persistence, OpenClaw/MCP/batch integration은 없다. OpenAPI/Swagger도 현재 구현으로 표시하지 않는다. 이런 확장은 별도 future work이며 이번 portfolio의 완료 조건이 아니다. [최종 검증/공개 전 점검](development/turn-018.md)에 현재 상태와 필요한 조치를 기록한다.

이 프로젝트의 결과는 AI가 Job Fit의 정답을 찾았다는 증명이 아니다. 불확실한 AI 판단을 입력·선택지·검증·관측·명시적 한계가 있는 software component로 다루고, 사람이 다음 결정을 내릴 수 있게 만든 과정이다.

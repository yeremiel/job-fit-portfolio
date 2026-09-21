# Evaluation Contract

Turn 001에서 합의한 평가 목적, 범위와 원칙을 기록한다. Turn 002에서 구체화하고 Turn 003에서 확정한 match level, 평가 규칙과 baseline은 [Match Semantics](match-semantics.md)에 기록한다. 이 문서는 개념 수준의 계약이며, 구현된 기능이나 확정된 API / JSON schema를 나타내지 않는다.

Turn 004에서는 이 계약의 첫 구현으로 Go CLI와 Jev adapter를 도입했다. Profile은 최소 `evidence: [{id, text}]` JSON, JD는 raw text file로 읽는다. JSON Schema와 범용 문서 parser는 도입하지 않았다. [Runtime instructions](../internal/evaluator/instructions.md)는 이 문서와 Match Semantics의 규칙을 코드 실행에 필요한 형태로 요약한 단일 instruction 파일이며, baseline 값과 sample별 판단은 포함하지 않는다. 구현 범위와 검증 결과는 [Turn 004](development/turn-004.md)를 참조한다.

## Project Positioning

일본 및 한국 회사의 채용공고를 직접 검토하면서 반복하던 질문에서 출발한다.

> 이 JD가 요구하는 경험과 내가 실제로 가진 경험은 얼마나 잘 맞는가?

특정 채용 플랫폼 내부의 추천 기능에 종속되지 않고, 사용자가 직접 가져온 공개 Job Description과 자신의 Career Profile을 비교하는 개인용 Job Match 도구를 목표로 한다. 대규모 채용 플랫폼이나 범용 AI recruiting SaaS를 만들려는 프로젝트가 아니다.

Jev라는 structured AI evaluation 도구를 발견하고 이를 개인의 반복적인 JD 검토에 적용해보려는 것이 출발점이다. Turn 004에서 첫 evaluation engine으로 Jev를 채택했다. 개발 과정에서 AI evaluation, bias, evidence, explainability, scope control에 관한 판단과 실험을 기록하며, AI-assisted development portfolio로 활용한다.

```text
Personal Problem
    ↓
Personal Tool
    ↓
Structured AI Evaluation
    ↓
Engineering Decisions / Experiments
    ↓
Portfolio
```

## Primary User / Use Case

초기 primary user는 프로젝트 작성자 본인이다. Backend Engineer / Tech Lead 경력을 reference profile로 삼고, 초기 구현과 검증은 Software Engineering 포지션을 중심으로 수행한다.

사용자는 Candidate Career Profile과 Job Description을 제공한다. 시스템이 답하려는 핵심 질문은 다음과 같다.

> How well does the candidate's documented experience match the requirements described in this job description?

평가 대상은 JD에 기술된 요구 경험·역량과 profile에 기록된 실제 경험 사이의 정합성이다.

## Candidate Profile

Candidate Profile은 사실 기반 career evidence를 담는 canonical source다. 회사별 제출용 resume와 역할을 구분한다. Resume는 target JD에 맞춰 강조되는 정보가 달라질 수 있어, 이를 그대로 평가 입력으로 사용하면 evaluation bias가 발생할 수 있다.

Experience, Skills, Roles, Domains, Leadership, Achievements, Architecture experience, Operational experience 등을 기록할 수 있다는 개념을 고려한다. 이 목록은 필수 필드나 확정된 schema가 아니다. 구체적인 JSON schema와 수집·표현 방식은 이후 검토한다.

## Evaluation Scope

초기 MVP는 **Capability / Experience Match**만 평가한다. 다음은 평가 대상에서 제외한다.

- Hiring Probability / Expected Offer Probability
- Application Recommendation
- Competitor Comparison
- Salary Fit
- Visa Fit
- Location Fit
- Remote Work Preference
- Company Preference
- Employment Type Preference

따라서 “이 회사에 지원해야 하나?”, “합격 확률이 몇 퍼센트인가?”, “이 회사가 나에게 좋은 회사인가?”에는 답하지 않는다. 기록된 경력과 JD 요구사항이 어느 정도 부합하는지를 설명한다.

## Evaluation Principles

1. **Evidence First**: Candidate Profile에 실제로 존재하는 evidence를 근거로 평가한다. 기록되지 않은 경험을 추측하지 않는다.
2. **Unknown Remains Unknown**: JD나 profile에서 확인할 수 없는 사실을 임의로 보완하지 않는다. 정보가 없으면 Unknown으로 유지한다. Evidence가 없다는 이유만으로 해당 경험이 실제로 없다고 단정하지 않는다.
3. **Explainable Result**: 숫자만 제공하지 않는다. 어떤 경험이 강하게 매칭되고 어떤 근거가 부족한지 사람이 이해할 수 있어야 한다.
4. **No Hiring Probability**: Match Score는 실제 채용 성공 확률을 의미하지 않는다.
5. **Human Decision**: 시스템은 structured matching information을 제공하고, 지원 여부 등 최종 의사결정은 사용자가 한다.

## Evaluation Output

초기 계약은 다음 정보를 제공하는 것을 목표로 한다. 명칭은 개념적인 출력 항목이며, 직렬화 형식이나 필수 API 필드를 확정한 것은 아니다.

| 정보 | 전달하려는 내용 |
| --- | --- |
| Overall Match | 기록된 경험과 JD 요구사항의 전반적인 정합성 요약 |
| Match Level | Strong / Partial / Weak / Unknown. Turn 003의 정성적 정의와 Overall 해석은 Match Semantics 참조. 수치 집계 규칙 없음 |
| Strong Matches | 요구사항과 강하게 부합하는 경험 및 근거 |
| Partial Matches | 일부 부합하는 경험 및 근거 |
| Weak / Missing Evidence | 요구사항을 뒷받침하기에 약하거나 기록에서 찾을 수 없는 근거 |
| Unknowns | JD 또는 profile에서 확인할 수 없는 사실 |

위 표의 Weak / Missing Evidence는 출력 설명을 묶은 명칭이다. Weak는 match level이고 Missing Evidence는 모든 level과 별도로 기록할 수 있는 설명 정보다. 비교 가능한 관련 근거는 있으나 핵심 요구·수준과 gap이 크면 Weak, 평가할 근거가 부족하면 Unknown이다. Missing Evidence는 현재 profile만으로 충분히 입증하지 못하는 정보이며, 실제 경험 부재를 뜻하지 않는다.

Turn 003에서 다음 정성적 규칙을 확정했다.

- Strong은 핵심 요구에 직접 근거가 있고 role/context 차이가 크지 않은 경우다. Partial은 충분한 관련·전이 가능 경험이 있으나 scope, depth, technology, domain 또는 운영 맥락이 부족한 경우다.
- Required mismatch는 Preferred 부족보다 중요하다. Preferred 부족만으로 Overall을 크게 낮추지 않는다.
- Core role mismatch는 개별 skill overlap보다 우선한다. 일반 리딩과 특정 기술 사용을 합쳐 해당 분야의 Tech Lead 경험으로 만들지 않는다.
- Transferable experience는 문제의 유사성, 책임 수준, 설계 깊이, 운영 맥락, leadership/ownership 및 기술 전환 가능성으로 설명한다. Specialized experience의 존재를 추정하지 않는다.
- Overall은 core role, Required, 주요 capability, transferable experience, critical gap, Preferred를 함께 해석한다. 단순 평균이 아니며 core Required mismatch는 상한에 영향을 줄 수 있다. 모든 공고에 적용할 수치 weighting이나 고정 상한표는 만들지 않는다.

Requirement-level 결과는 다음 형태를 고려한다. 아래 예시는 설명용이며 실제 candidate에 대한 평가, 확정된 rubric 또는 schema가 아니다.

```text
Requirement:
Production cloud architecture experience

Matching Evidence:
- Basic AWS experience

Weak / Missing Evidence:
- Large-scale AWS architecture
- Production cloud operations

Unknown:
- Exact production traffic scale

Match:
Weak
```

실제 평가에서는 JD에 존재하는 요구사항과 profile의 근거를 기준으로 판단해야 한다. 이 예시만으로 모든 cloud architecture 요구에 대규모 운영 경험이 필요하다고 일반화하지 않는다.

## Optional Match Score

직관적인 Match Score를 제공하는 방향을 고려한다. 예를 들어 `AI Match 82%`, `Match Level HIGH`와 같은 표현을 생각할 수 있으나, 이는 표시 예시일 뿐 실제 평가 결과나 확정된 출력 형식이 아니다.

**Match Score ≠ Hiring Probability**

점수가 의미하는 것은 다음뿐이다.

> The degree of alignment between the documented candidate experience and the requirements described in the job description.

점수 제공 여부, 계산 공식, 가중치, threshold와 점수 기반 등급 판정 방식은 미정이다. 위의 `HIGH`는 초기 표시 예시이며, Turn 002 baseline은 Strong / Partial / Weak / Unknown을 사용하고 percentage를 부여하지 않는다. 실제 여러 JD를 평가하고 결과를 관찰한 뒤 calibration하는 것을 고려한다. 점수만으로 설명 가능한 evidence를 대신할 수 없다.

## Initial Conceptual Flow

```text
Candidate Profile + Job Description
                  ↓
Job Requirements Extraction
                  ↓
Candidate Evidence Matching
                  ↓
Requirement-level Assessment
                  ↓
Strong / Partial / Weak / Unknown
                  ↓
Overall Match Summary
                  ↓
Optional Match Score
```

이는 conceptual model이다. 구현 단계, 모듈 구성, API 호출 순서 또는 implementation architecture로 확정하지 않는다. Deterministic rule과 semantic judgment의 역할 분담도 미정이다.

## Generalization / Non-goals

Java/Spring 전용 evaluator로 만들거나 특정 개인의 필드명을 schema에 직접 넣지 않는다. Candidate Profile 개념은 다른 사람에게도 적용할 수 있도록 고려한다. 다만 모든 직군에 대한 평가 품질이나 범용 채용 서비스 완성은 약속하지 않는다.

| 구분 | 범위 |
| --- | --- |
| Concept | General job matching approach |
| MVP | Personal software-engineering job matching tool |
| Reference User | Backend Engineer / Tech Lead |
| Portfolio Scope | AI evaluation design and development process |

다음은 현재 프로젝트 범위 밖이다.

- 모든 직군에 대한 평가 품질 보장 및 직군별 evaluation strategy
- 사용자 onboarding 및 resume parsing
- 계정/인증 및 개인정보 저장 서비스
- Web SaaS product 및 recruiting platform 구축
- 자동 입사 지원 및 AI의 최종 지원 결정
- 대규모 production infrastructure 구축

## Development Philosophy / Open Decisions

별도의 ChatGPT 대화에서 기획, 범위 결정과 결과 검토를 수행하고, Codex workspace에서는 해당 Turn의 구현 및 검증을 수행한다.

```text
Brainstorm → Scope → Task → Implementation → Test / Review → Result → Next Decision
```

각 Turn 결과를 검토한 뒤 다음 결정을 내린다. 필요성이 확인되지 않은 기술과 infrastructure를 미리 추가하지 않는다. Turn 004의 Go·CLI·Jev 선택 외에 MCP, Database, Web UI 및 AI coding agent integration은 이후 필요성과 trade-off를 검토할 대상이다. 초기 TypeScript 아이디어 대신 첫 구현에는 Go를 선택했다. 이전 disposable PoC 코드는 재사용하지 않는다.

이 문서는 평가 범위와 정성적 규칙을 정하지만, profile schema, 복합 requirement의 표준 분해 단위, 점수 산출 방법 및 구현 architecture는 정하지 않는다. 언어·학력 조건의 일반적인 평가 범위도 후속 검토 사항이다.

## Turn 006 Implementation Note

Overall은 기존과 같이 전체 requirement assessments와 profile/JD/policy를 받은 Jev Choice로 결정된다. 이후 별도 Choice 호출로 `decisionTrace`의 Supporting / Limiting / Non-decisive 역할을 분류한다. ID·Match·source type은 기존 assessment에서 복사하며 evidence는 ID로 연결한다. 이는 선택된 결과에 대한 사후 attribution이며 원래 판단의 인과 기록은 아니다. Trace는 Overall 계산에 사용되지 않고 기존 match semantics를 바꾸지 않는다. 상세 설계와 실제 결과는 [Turn 006](development/turn-006.md)을 참조한다.

## Evaluation Philosophy — Turn 007

### Reference와 Overall의 역할

Job Fit의 Overall Strong / Partial / Weak / Unknown에 객관적 ground truth가 있다고 가정하지 않는다. **Human Reference ≠ Ground Truth, Jev Result ≠ Ground Truth.** 기존 Human Baseline / Expected Match는 **Pre-evaluation Reference Interpretation**, 즉 evaluator 결과를 보기 전에 기록한 하나의 사전 해석이다. Codex 초안에서 출발해 Turn 003에서 비교 기록을 확정했다는 provenance를 유지한다. 목적은 사전에 중요하게 본 요소를 보존하고 판단 차이와 evaluation behavior를 분석하는 것이며, target label 제공이 아니다.

Human Reference와 단순 label 일치율을 primary metric으로 사용하지 않는다. Prop Tech plus의 Reference Strong / Jev Partial만으로 수정하지 않으며, 두 해석이 현재 evidence와 semantics 안에서 합리적이면 disagreement로 기록하고 끝낼 수 있다. Reference에 맞을 때까지 prompt나 evaluator를 조정하는 것은 특정 해석에 대한 overfitting 위험이 있다. 이는 과거에 실제 tuning을 했다는 주장이 아니라 앞으로의 변경 기준이다.

**Overall Match is a summary signal, not an objective truth.** 빠른 JD 비교를 돕는 structured judgment summary로 유지한다. 객관적 적합률이나 채용 확률을 뜻하지 않으며 requirement-level 결과·evidence·Unknown이 더 중요한 정보다. 이전 Optional Match Score 절은 과거 검토 아이디어로 남기며 이번 결정으로 점수, 가중치 또는 threshold를 도입하지 않는다.

### Evaluator quality criteria

| 기준 | 관찰할 질문 / failure 또는 한계의 예 |
| --- | --- |
| Evidence Faithfulness | Profile에 없는 경력을 생성하는가? Kubernetes 근거 없이 production 운영 경력이 있다고 주장하면 failure다. 원문 복사만으로 의미적 충실성까지 보장되지는 않는다. |
| Requirement Fidelity | JD의 Required / Preferred / Responsibility / Stack / Context 의미를 보존하는가? Azure 우대를 필수로 바꾸면 failure다. |
| Unknown Preservation | 확인할 수 없는 사실을 추측하거나 No Evidence를 No Experience로 바꾸지 않는가? Weak와 Unknown의 기존 구분을 유지한다. |
| Evidence Traceability | 어떤 evidence가 Match를 지지하고 무엇이 missing/unknown인지 연결해 확인할 수 있는가? Trace가 있더라도 원래 판단의 인과성을 보장하지 않는다. |
| Structural Consistency | 여러 JD에서 level 의미, source type 구분과 output contract를 일관되게 유지하는가? 같은 label을 강제하는 것이 아니다. |
| Repeatability / Variability | 동일 입력에서 어떤 항목과 결과가 변하는지 관찰·기록하는가? 회계 Weak → Partial, Azure 환경 Partial → Unknown이 이미 관찰되었다. 완전한 결정론을 필수로 요구하지 않으며 변동률은 아직 측정하지 않았다. |
| User Usefulness | 사용자가 강점, 부분적 부합, 약한 evidence, profile만으로 알 수 없는 내용을 이해하는 데 도움이 되는가? Overall label 하나보다 항목별 정보가 중요할 수 있다. |

객관적 Overall 정답을 가정하지 않는다고 해서 사실·계약 위반도 취향 차이로 취급하지는 않는다. 위 기준은 앞으로의 검토 기준이며 모든 기준을 이미 충족했다는 인증이나 새 numeric metric이 아니다.

### Jev와 software의 역할

탐구 질문은 **“Can structured AI decision primitives be used as a constrained and observable component inside a reusable job-matching workflow?”**다. Jev가 일반 LLM보다 Job Fit 정답을 더 정확히 맞힌다는 것을 증명하는 프로젝트가 아니다.

```text
Canonical Profile + JD
→ Application-controlled workflow
→ Jev structured judgment
→ Validation / Traceability
→ Structured Match Result
→ Human interpretation
```

Engineering value는 모델의 지능 우위가 아니라 profile을 지속적으로 재사용하고, explicit semantics, consistent structured output, requirement 단위 평가, evidence tracing, missing evidence와 Unknown visibility, 반복 가능한 workflow 및 machine-readable 결과를 제공하는 데 있다. 현재 canonical profile의 지속성은 로컬 JSON 파일 재사용이며 DB/result persistence를 뜻하지 않는다. Decomposition은 현재 줄 단위라는 한계가 있다. 향후 MCP/automation integration은 가능성이며 아직 구현하지 않았다. 최종 판단은 사람이 한다.

### Frozen evaluation baseline / 변경 조건

Turn 006의 현재 구현을 **evaluation baseline implementation**으로 정하고 다음 전체 sample evaluation이 완료될 때까지 유지한다. 이는 정답 label 집합과 구분되는 구현 기준선이다. Candidate Profile, match semantics, Jev requirement instructions, Overall instructions, requirement decomposition logic, Decision Trace semantics, Overall Match logic을 변경하지 않는다. Evidence mapping을 포함한 현재 평가 동작도 한 sample에 맞춰 조정하지 않는다.

Freeze 기간 중 명백한 execution bug 또는 contract violation만 예외로 검토한다. 예외 수정 시 문제 근거·변경 범위·비교 가능성 영향을 기록하며 조용히 기준선을 교체하지 않는다. 여러 sample의 systematic 문제 후보는 먼저 관찰하고 evaluation 이후 planning에서 변경을 검토한다.

Evaluation logic 변경 근거는 다음 두 범주다.

- **Contract Violation:** 없는 경력 생성, Required/Preferred 역전, Unknown을 실제 무경험으로 단정, core role mismatch를 무시한 keyword Strong, stack/context의 명시 필수 승격, evidence와 Match 모순 등. 구체적인 입력·출력과 계약 근거로 확인한다.
- **Systematic Problem:** 여러 sample에서 동일한 문제 패턴이 반복됨. 예를 들어 Preferred Unknown이 반복적으로 과도한 제한으로 작용하거나 business context가 여러 JD에서 Limiting으로 나타나는 패턴을 검토한다. 반복된 Limiting 자체만으로 위반이 확정되지는 않으며 JD 의미·근거와 함께 판단한다.

한 sample의 label 불일치나 단일 post-hoc attribution만으로 수정하지 않는다. Decision Trace는 별도 Jev 호출의 **structured post-hoc attribution**으로, causal trace나 implicit weighting의 증거가 아니다.

### 다음 evaluation phase

기존 6개 JD: BlueMeme / Micro Court, Prop Tech plus, PIA TECH LAB, CORE, teamLab, Nitori Digital Base를 사용한다. Prop Tech plus의 기존 관찰을 보존하고 나머지 5개를 현재 evaluator로 실행한 뒤 전체 behavior를 비교한다. 목적은 reference 일치 개수가 아니라 evidence hallucination, requirement meaning, Required/Preferred 구분, Unknown, core role mismatch, evidence traceability, trace usefulness, 예상 밖 동작과 variability를 검토하는 것이다.

Nitori는 negative/control case다. Angular/TypeScript와 일반 Tech Lead 경험을 합쳐 frontend Tech Lead 경력을 생성하거나 core role mismatch를 무시하면 contract violation 후보로 본다. 단지 Overall이 reference Weak와 다르다는 사실만으로 실패로 간주하지 않는다.

현재 실행용 sample 파일은 Prop Tech plus만 있다. 나머지 5개는 기존 문서의 출처·검토일·JD 의미를 보존하여 입력을 준비할 후속 작업이며 이번에 준비하거나 실행하지 않는다. teamLab 언어와 runtime 언어/학력 제외 범위 등 비교상 차이도 결과 해석에 기록해야 하며 이 Turn에서 scope를 바꾸지 않는다.

## Turn 010 Implementation Note

Turn 009에서 확정한 두 위반만 수정했다. 기존 Limiting 정의에 따라 Strong Overall에서는 Limiting을 Choice에서 제외하고 domain 검증에서도 거부한다. Overall이나 requirement Match를 보정하지 않는다. Version/duration rationale은 해당 requirement에 명시된 표현이 인식될 때만 선택·적용할 수 있다. Profile의 기간 미기록을 JD 조건으로 바꾸지 않으며, 다른 gap도 명시 조건 또는 해당 평가에 직접 필요한 정보에 근거해야 한다. 사전 경험·자격·규모·ownership의 새 qualification을 추정하지 않도록 지시한다. 표현 인식 및 의미적 보장의 한계는 [Turn 010](development/turn-010.md)에 기록한다. Public/Energy Context의 coverage 의미는 여전히 Ambiguous이며 이번에 결정하거나 변경하지 않았다.

## Product Boundary / Interface Direction — Turn 012

Turn 011에서 accepted된 evaluator를 v0 evaluation baseline으로 유지하고 iterative tuning을 종료했다. 이제 제품 책임은 **이미 확보된 JD + canonical candidate profile → structured evaluation result**로 한정한다. Caller는 OpenClaw 또는 MCP client로 제한하지 않는다.

Orchestration과 evaluation은 분리한다. 회사/포지션 탐색, 회사 성격 조사와 SES/派遣 등의 hard filtering, 채용 페이지 crawling/JD acquisition, 실행 시점 결정, Notion 등 결과 저장은 외부 orchestration의 책임이다. job-fit은 전달된 JD와 기록된 evidence의 capability/experience match만 평가하며 최종 지원 판단은 사람이 한다. Company metadata가 전달되더라도 이를 회사 조사나 지원 대상 filtering 책임으로 확장하지 않는다.

REST API를 primary integration interface **후보**로 채택한다. CLI는 로컬 개발·테스트·수동 사용에 유지하고, MCP는 필요할 때 REST API를 호출하는 optional adapter로 고려한다. CLI와 REST는 같은 evaluator core를 사용하며 transport 또는 MCP 계층에 evaluation logic을 복제하지 않는다. 이는 현재의 Go core를 여러 caller에서 재사용하기 위한 방향이며 HTTP server나 MCP가 이미 구현되었다는 뜻은 아니다.

초기 단일 사용자 API에서는 server-side canonical profile을 쓰는 방향을 검토한다. 현재 CLI는 profile/job 파일을 명시적으로 받는다. Server의 profile 선택·로딩·갱신 방식은 미정이며 multi-user/profile ID/요청별 explicit profile은 현재 MVP의 필요 범위가 아니다.

구체 endpoint, request/response schema, framework, HTTP server 구현은 이번에 확정하지 않는다. User prompt의 endpoint/JSON 예시는 개념 설명이며 승인된 wire contract가 아니다. 다음 Turn에서 metadata, canonical profile handling, errors, timeout/external Jev failure, API versioning, CLI/API 공유 application boundary를 포함한 **HTTP API Contract**를 설계하고 검토한 뒤 구현한다. [Turn 012](development/turn-012.md)에 현재 흐름과 결정 상태를 기록했다.

## HTTP API Contract — Turn 013

[API Contract](api-contract.md)에 `POST /v1/job-fit/evaluate`의 동기 request/response, server-side canonical profile snapshot/version, 오류/status 및 shared application boundary를 정의했다. 이는 구현 대상 계약이며 HTTP service는 아직 없다. Company/title/sourceUrl은 선택 metadata이고 `description`만 JD 평가 입력이다. API DTO는 기존 Match/evidence/Unknown/사후 Trace 의미를 보존하며 새 평가 정책을 넣지 않는다. API v1과 accepted evaluation baseline v0은 별개다. [Turn 013](development/turn-013.md)에 선택 근거와 향후 구현 시 필요한 오류 전달 개선을 기록한다. 기존 평가 semantics 및 ambiguity는 변경하지 않았다.

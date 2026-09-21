# job-fit — AI 활용 포트폴리오

반복하던 채용공고 검토를 위해, 동일한 경력 근거와 JD를 비교해 구조화된 결과를 반환하는 Go CLI/REST 도구를 만들었습니다. ChatGPT를 planning partner로, Codex를 구현·검증 도구로 사용했고, application 내부에서는 Jev를 제한된 선택지에 응답하는 structured decision component로 사용했습니다.

이 문서는 **AI의 한계를 겪고 해결한 경험**과 **LLM 기반 기능 개발 경험**을 중심으로, 문제 해결 과정과 개발 과정에서의 AI 활용을 함께 설명합니다.

**검토 경로:** 이 문서 → [두 위반의 수정 기록](development/turn-010.md) → [regression test](../internal/jev/contract_regression_test.go) → [실제 평가 결과](development/turn-011.md). 실행은 [README Quick Start](../README.md#quick-start)를 참조하시면 됩니다. 이 문서가 포함된 repository에서 전체 코드를 확인할 수 있습니다.

## 1. 해결하려던 문제와 결과

여러 회사의 JD를 검토할 때마다 같은 경력을 반복해서 비교했습니다. 일반 LLM의 답변만 사용하면 기준이 달라지거나, 기록하지 않은 경험을 추정하거나, 공고 간 결과를 비교하기 어려웠습니다. Jev를 발견한 뒤 이를 개인 workflow에 적용하는 실험으로 시작했습니다.

| 이전 방식의 문제 | 구현한 대응 |
| --- | --- |
| 질문마다 경력 표현과 판단 기준이 달라질 수 있음 | 회사별 제출 resume와 분리한 canonical profile 및 명시적 evaluation semantics |
| 적합성 결론만으로 근거를 확인하기 어려움 | Requirement별 Match와 원문 evidence ID 연결 |
| 정보 누락을 경험 부재로 해석할 위험 | Unknown을 별도 결과로 유지 |
| 반복 작업에서 답변 형식이 달라짐 | CLI/REST의 structured JSON과 validation |

결과는 Overall 및 requirement별 Strong / Partial / Weak / Unknown, supporting/missing evidence, unknowns, Decision Trace입니다. **채용 확률이나 지원 추천을 제공하지 않습니다.** 시간 절감률이나 평가 정확도 향상률은 측정하지 않았습니다. 개선은 입력·출력·검증 절차를 코드로 반복 가능하게 만든 데 있습니다.

## 2. 본인과 AI의 역할

| 영역 | 본인 | AI |
| --- | --- | --- |
| 문제·범위 | 실제 JD 검토 workflow 제공, 목표·비목표·MVP 범위 결정 | ChatGPT의 문제 구조화·대안 제시 |
| 평가 철학 | Reference를 정답으로 사용하지 않는 방향, Unknown 및 종료 기준 결정 | 해석 차이·한계 분석과 비판 지원 |
| 설계 | 문제 구조, 평가 계약, CLI/REST 경계와 단계별 구조를 설계하고 결정 | 설계 대안, 기존 코드 분석 및 장단점 검토 지원 |
| 구현·테스트 | Turn별 요구사항·제약·acceptance criteria 제공, 결과 검토와 수정·종료 결정 | Codex가 실제 코드와 테스트를 모두 작성하고 실행·수정·문서화 |
| 오류 처리 | 수정할 두 위반과 보류할 ambiguity의 범위 확정 | Source/result/code 근거 수집, 분석 및 수정 구현 |

이 프로젝트에서는 설계와 구현의 역할을 의도적으로 분리했습니다. 저는 문제 정의, 평가 규칙, 아키텍처, 단계별 요구사항과 수용 기준을 설계하고 결과의 채택·수정·종료를 결정했습니다. **실제 코드와 테스트 작성은 모두 Codex에 맡겼습니다.** 아래 수정 사례 역시 제가 직접 코드를 고쳤다는 의미가 아니라, 수정 범위와 검증 조건을 정하고 Codex가 구현한 결과를 실행·검토한 사례입니다.

Workflow는 `ChatGPT planning → 범위/Task Prompt → Codex 구현·검증 → 본인 결과 검토 → 다음 Turn 결정`이었습니다. [상세 개발 narrative](ai-assisted-development.md)에 주요 전환점과 역할을 기록했습니다.

## 3. AI의 한계를 겪고 해결한 경험

### 출발: 결과가 다르다고 모델을 정답에 맞추지 않기

Prop Tech plus의 사전 Reference는 Strong, Jev 결과는 Partial이었습니다. 처음에는 이 차이가 개선 대상처럼 보였지만, Job Fit에는 객관적 정답을 가정하기 어렵다는 문제가 있었습니다. Reference 역시 Codex 초안을 사람이 검토·확정한 해석이지 독립적인 ground truth가 아니었습니다.

먼저 Decision Trace를 추가해 관측 가능성을 확보했습니다. 이후 label 일치 대신 evidence faithfulness, requirement fidelity, unknown preservation, traceability와 consistency를 검토했습니다. 실제로 reference에 맞춰 tuning했다가 되돌렸다는 이야기는 아닙니다. [진단](development/turn-005.md), [관측 추가](development/turn-006.md), [철학 재정리](development/turn-007.md)에 그 과정을 남겼습니다.

### 사례 A: Strong Overall인데 Limiting이 존재함

- **관찰:** CORE 결과가 Overall Strong인데 Limiting requirement를 포함했습니다.
- **계약 위반:** Limiting은 더 높은 Overall을 제한한다는 의미인데 Strong은 최고 단계였습니다.
- **원인:** Jev에 모든 Overall에서 같은 선택지를 제공했고, domain validation은 Overall을 전달받지 못했습니다.
- **수정:** Strong에서는 Limiting 선택지만 제외하고 domain에서도 해당 조합을 거부했습니다. Overall을 낮추거나 trace를 자동으로 Supporting으로 바꾸지 않았습니다. Non-decisive는 유지했습니다.
- **검증:** 저장된 오류 fixture와 mock 응답으로 차단을 확인했습니다. Turn 010 live CORE는 Partial이어서 Strong 분기의 live 검증이라고 주장하지 않습니다. 후속 Turn 011에서는 CORE가 Strong이며 Limiting 없이 반환됐습니다.

근거: [수정 상세](development/turn-010.md#fix-1--strong-overall--limiting-trace), [domain regression](../internal/evaluator/trace_contract_test.go), [Jev Choice regression](../internal/jev/contract_regression_test.go), [후속 CORE JSON](development/turn-011-results/core.json).

### 사례 B: JD에 없는 version/duration gap 생성

- **관찰:** Nitori의 Angular/TypeScript requirement에 명시되지 않은 version/duration 부족이 표시됐습니다.
- **원인:** 모든 requirement에 duration rationale 선택지를 제공했고 application이 선택 결과를 정형 문구로 렌더링했습니다. 모델이 자유 문장을 지어낸 것만의 문제가 아니라 application이 허용한 선택지와 template 적용에도 원인이 있었습니다.
- **수정:** 해당 requirement 자체에 명시된 version/기간 조건이 인식되는 경우에만 선택지를 제공하고, template 적용 단계에서도 재검증했습니다.
- **검증:** Angular/TypeScript만 있으면 duration 선택·적용을 거부합니다. 반대로 `Java 8+ development for at least one year`는 허용합니다. Nitori 재실행에서도 합성 version/duration gap이 사라졌습니다.

근거: [수정 상세](development/turn-010.md#fix-2--synthetic-version--duration), [eligibility 구현](../internal/jev/qualification.go), [regression tests](../internal/jev/contract_regression_test.go), [Nitori 결과](development/turn-011-results/nitori.json).

이 검사는 제한된 lexical rule이며 모든 언어·복합 문장의 qualification을 완벽히 이해하는 parser가 아닙니다.

### 해결하지 않은 항목과 종료 판단

Public/Energy 복합 context의 일부 evidence만으로 전체를 Strong으로 볼 수 있는지는 **Ambiguous**로 남겼습니다. 세 후보 중 두 개만 Confirmed Violation으로 수정했고, 사람의 기대와 다른 결과를 모두 버그로 취급하지 않았습니다.

이후 [6개 실제 JD를 각각 한 번 평가](development/turn-011.md)해 두 위반의 재발 및 새로운 명시적 위반이 발견되지 않아 baseline을 수용했습니다. 이는 제한된 sample의 contract 검토이며 통계적 정확도나 무오류를 보장하지 않습니다. 원하는 label이 나올 때까지 반복하지 않았습니다.

## 4. LLM 기반 기능 개발 경험

### Jev를 사용하는 위치

```text
CLI ───┐
       ├─ Application Service ─ Evaluator ─ Jev Adapter ─ Jev API
REST ──┘                         │
                                └─ 검증된 structured result
```

Jev(TypeSafe AI)의 구조화된 의사결정 API를 외부 AI 판단 컴포넌트로 연동했습니다. 선택한 이유는 질문과 선택지를 이용해 구조화된 판단을 실험할 수 있었기 때문입니다. 다른 모델보다 정확하다는 benchmark에 근거한 선택은 아닙니다. Application이 orchestration, 허용 선택지, mapping 및 validation을 책임지고 Jev는 semantic judgment를 담당합니다.

평가는 `JD 줄 분류 → evidence relation → requirement Match → Overall → Decision Trace → consistency validation`으로 진행됩니다. Overall은 평균이나 점수 계산이 아닙니다. Decision Trace는 이미 선택된 Overall에 대해 별도 요청으로 생성한 **structured post-hoc attribution**이며 최초 판단의 causal trace가 아닙니다.

### 구현한 경계와 근거

| 관심사 | 구현 / 확인 위치 |
| --- | --- |
| AI API 연동 | [Jev client](../internal/jev/client.go): structured question/choice, 응답 mapping, provider 오류 처리 |
| 재사용 가능한 use case | [Application Service](../internal/application/service.go): CLI와 REST가 동일 evaluator 사용 |
| 입력과 resource 제한 | HTTP body/description, evidence pair 및 provider payload 크기 제한. Token 예산 관리 기능과는 구분 |
| Timeout/cancellation | 요청 context를 provider까지 전달. [경계 테스트](../internal/jev/http_boundary_test.go) |
| 안전한 오류 반환 | Typed error를 HTTP envelope로 mapping. Raw upstream 오류·key는 노출하지 않음 |
| 근거 추적 | Canonical snapshot, SHA-256 profileVersion, requirement/evidence ID, trace consistency 검증 |
| 실제 integration | [Turn 014 HTTP/Jev 응답](development/turn-014-results/proptech-plus.json), [API 문서](api/job-fit.md) |

RAG, fine-tuning, cache 기반 비용 절감, 월 예산 관리 기능은 구현하지 않았습니다. 추가 live 호출은 필요한 sample 검증으로 제한했고 통상 unit test는 외부 API 없이 실행합니다. API 비용 절감률을 산출하지는 않았습니다.

## 5. 실행 및 검증 자료

[README Quick Start](../README.md#quick-start)에 Go 요구사항, key 환경변수, CLI/server build와 curl 예제가 있습니다. Live evaluation에는 사용자의 Jev key와 provider 비용이 필요합니다. Key 없이도 다음 테스트를 실행할 수 있습니다.

```sh
go test ./...
go test -race ./...
go vet ./...
```

기존 [최종 regression](development/turn-018.md)과 [공개 준비 검증](development/turn-019.md)에 test/race/vet 및 두 실행 파일 build 결과를 기록했습니다. 이번 제출 문서 작성에서는 application 변경이나 Jev 재호출을 하지 않았습니다.

## 6. 데이터 출처·공개 범위·한계

- JD sample은 공개 채용공고에서 평가에 필요한 내용을 요약한 자료이며, [sample 설명과 source URL](match-semantics.md)에 출처가 있습니다.
- Candidate profile은 **본인의 공개용 경력 요약**입니다. 가상 데이터라고 표시하지 않습니다. 이전 회사의 소스 코드, 특정 고객·비공개 프로젝트명, credential은 포함하지 않는 범위로 검토했습니다. [공개 정보 검토](development/turn-019.md)를 참조하십시오.
- 공개용 repository는 과거 개발 Git history를 가져오지 않고 정리된 현재 파일로 별도 생성했습니다. 이는 작성자가 보고한 공개본 최종 검증에 근거하며, 과거 Turn 기록의 당시 미완료 상태와 구분합니다. 개발 판단의 history는 문서로 보존했습니다.
- Jev judgment의 변동성과 compound/context ambiguity는 남아 있습니다. Profile도 한 개인의 Software Engineering 경력에 한정됩니다.
- 현재 API는 인증 없는 loopback personal service입니다. 공개 인터넷 서비스나 전체 지원 자격 판정 시스템이 아닙니다.
- OpenClaw, generic candidate endpoint, MCP adapter, OpenAPI/Swagger는 Future Work이며 제출물에 구현된 기능으로 주장하지 않습니다.

## 더 살펴보기

- [Architecture Overview](architecture/overview.md): component, runtime 및 error boundary
- [API Specification](api/specification.md) / [Job Fit API](api/job-fit.md): 실행 가능한 외부 계약
- [Evaluation Contract](evaluation-contract.md) / [Match Semantics](match-semantics.md): 판단 기준과 한계
- [AI-assisted Development](ai-assisted-development.md): 사람과 AI의 협업 과정
- [Development History](development/): 당시 판단·실패·검증 결과

이 프로젝트에서 보여드리고 싶은 것은 AI가 적합성의 정답을 찾았다는 주장이 아닙니다. AI가 만든 결과를 검토 가능한 계약 안에 두고, 실제 오류와 해석 차이를 구분하며, 수정과 종료의 기준을 책임지는 개발 과정입니다.

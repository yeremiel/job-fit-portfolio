# Turn 006 — Overall Decision Observability

- 날짜: 2026-09-21
- 범위: 기존 Overall 선택을 유지하고 사후 structured Decision Trace 추가.
- 시작 상태: `main`, HEAD `d5145a7d4095d0df7081f4f7f98750143974ea3f`. 기존 untracked `turn-005.md` 보존.

## Completed

Turn 005에서는 Overall Partial을 재현했지만 결정적 항목을 식별할 수 없었다. 이번에는 기존 평가 후 각 requirement의 Supporting / Limiting / Non-decisive 역할을 별도 Jev Choice로 관찰한다. 실제 Prop Tech plus 한 번의 실행은 성공했고 Overall Partial을 유지했다. Baseline Strong을 목표로 조정하지 않았다.

## Before Implementation / Plan

기존 코드·입력·테스트·문서를 먼저 확인하고 사용자에게 판단 흐름과 두 대안을 제시했다. 같은 요청에 Overall과 trace 질문을 함께 넣는 방법은 호출 수를 유지하지만, 독립 질문이 실제 Overall 응답을 참조하지 못하고 기존 요청도 바뀐다. 기존 Overall 선택 후 별도 요청하는 방법은 한 번의 호출을 추가하지만 기존 판단 요청을 보존한다. 사용자가 후자를 승인했다.

구현 계획은 기존 ID 재사용 → domain trace/reference 및 검증 → 별도 Jev 분류 → JSON 연결 → offline 회귀/실패 테스트 → 동일 sample 실제 1회 → 결과 기록 순서였다. 별도 계획 문서 계층이나 commit은 만들지 않았다.

## Current Overall Flow

```text
JD line classification
→ evidence relations
→ requirement match + rationale
→ 기존 Jev Overall Choice (여기에서 Overall 결정)
→ 새 Jev post-hoc influence Choice (각 requirement)
→ application 참조/정합성 검증 및 JSON 출력
```

[기존 mapper](../../internal/jev/mapper.go)의 `Overall`은 policy, 전체 profile, JD, 모든 assessments를 전달한다. 각 assessment에는 ID/원문/category/Match/evidence/missing/unknown/reasoning이 포함된다. Source type은 이 시점에도 보존된다. Summary 개수는 전달하지 않으며 application이 Overall을 계산하지 않는다. 기존 Overall 질문, 선택지, state와 기존 `overallReasoning` 생성 코드는 변경하지 않았다.

기존 설명은 선택한 Overall 정의에 행별 label/category 목록을 붙이는 방식이다. 그 응답만으로 결정적 항목을 복구할 수 없으므로 [새 trace adapter](../../internal/jev/trace.go)를 추가했다. 동일 입력·평가 전체에 실제 선택된 `overallMatch`를 더해 별도 요청한다. Trace는 앞 단계로 되돌아가거나 결과를 수정하지 않는다. 기본 호출은 4회에서 5회로 늘었다. 모든 evidence가 미연결이면 항목 판단 단계를 건너뛰므로 4회다. 자동 retry는 없다.

## Decision Trace Design

`decisionTrace`는 `source`, `supporting`, `limiting`, `nonDecisive`를 가진다. 각 목록의 항목은 `requirementId`, `requirement`, `match`, `sourceType`이다. 역할은 소속 목록으로 표현한다. 원문/Match/sourceType은 provider에게 새로 생성시키지 않고 기존 assessment에서 복사한다.

```json
{
  "requirementId": "L006",
  "requirement": "Required: At least one year of development experience with Java 8 or later.",
  "match": "Partial",
  "sourceType": "required"
}
```

위는 실제 Limiting 항목이다. `requirements`에서 `id == "L006"`인 항목을 찾아 E1과 버전·기간 missing/unknown을 확인할 수 있다. Evidence를 중복 저장하지 않는다. 기존 L번호는 동일 입력의 비어 있지 않은 줄 순서에서 안정적이며, 중간 ignore 행이 있어도 번호를 다시 매기지 않는다. JD 수정 시 같은 의미의 영구 ID를 보장하지 않는다.

역할은 중요한 긍정 근거(Supporting), 더 높은 match를 제약하는 중요한 gap/uncertainty(Limiting), 결정적이지 않은 항목(Non-decisive)이다. 회사·기술·ID별 역할이나 category별 고정 변환은 없다. 작은 reason code taxonomy는 이번 목적에 필요하지 않아 추가하지 않았다.

**Provenance:** `source`는 `jev-post-hoc-attribution; not a causal record of the original decision`이다. 이것은 모델이 이미 선택된 Overall에 부여한 사후 attribution이지 원래 요청 내부의 인과 기록이 아니다. 따라서 “왜 Partial인가”를 검토할 구조화된 단서는 제공하지만 그 설명이 최초 판단의 실제 원인이었음을 증명하지 않는다. 자유 텍스트 LLM explanation은 도입하지 않았다.

## Consistency / Failure Handling

[ResolveTrace](../../internal/evaluator/trace.go)는 모든 assessment가 정확히 한 번 참조되는지 검사한다. 빈/중복 source ID, 미존재 ID, 중복 reference, 누락, invalid influence를 거부한다. Unknown 또는 연결 evidence가 없는 항목의 Supporting도 거부한다. 전체 Match/sourceType/원문은 canonical assessment에서 복사하여 불일치를 방지한다. 빈 목록은 `[]`, 목록 순서는 기존 assessment 순서다.

Preferred Unknown을 Limiting으로 선택하는 것은 허용한다. 이를 application이 Non-decisive로 보정하면 관찰 목적을 훼손하기 때문이다. 나머지 모든 semantic 조합을 고정 표로 제한하지 않는다. 예를 들어 Strong인데 Limiting이라고 반환하는 등 의미적 의문은 현재 기본 guard가 모두 해결하지 않는다. Trace 의미와 category별 중요도는 사람 검토가 필요하다.

Trace HTTP 오류 또는 검증 실패 시 기존 오류 흐름으로 non-zero 종료하며 성공 JSON을 출력하지 않는다. 이미 Overall을 얻었어도 trace 없는 성공처럼 보이지 않도록 한 선택이다. 이는 평가 semantics 변경은 아니지만 추가 호출의 비용·지연·실패 경로라는 실행상 변화다. 결과를 보정하거나 API를 재호출하지 않는다.

## Tests

- `go test ./...`: 통과. 기존 테스트 유지 및 새 테스트 추가. 외부 Jev API 호출 없음.
- `go vet ./...`: 통과.
- `go build -o /tmp/job-fit-turn006 ./cmd/job-fit`: 통과.

Domain tests: 세 역할, source identity/Match/category 복사, evidence 연결 가능한 JSON, 빈 배열, 누락·중복·미존재 참조, invalid influence, Unknown/무근거 Supporting 거부, Preferred Unknown Limiting 허용.

Adapter tests: 세 Choice의 mapping, invalid provider choice 거부, 실제 Overall과 source category/unknown 정보를 포함하는 state. Pipeline test는 5단계 메모리 HTTP로 변경하고 기존 평가/evidence assertion을 유지했다. 무연결 evidence 경로는 4단계에서 Unknown을 보존한다. Orchestration tests는 trace가 실제 Overall을 받고 이를 변경하지 않으며 provider 실패와 invalid trace가 성공 result로 노출되지 않음을 검사한다.

변경 전후 hash로 Turn 000–005, candidate profile, sample JD, match semantics, runtime policy, 기존 Jev mapper/HTTP client가 동일함을 확인했다. Formatting/whitespace, 문서 링크, 실제 trace 15개 참조 및 evidence 원문 대응을 확인했다. 키를 출력하거나 문서에 기록하지 않았고 `.env`는 Git ignored 상태다.

## Prop Tech plus Result

동일 profile/sample과 기존 `jev-1.13.0`을 사용해 실제 CLI를 **한 번** 실행했다. Exit 0. 키는 기존 사용자 `.env`에서 자식 프로세스 환경변수로 전달했다. 임시 진단 JSON만 repository 밖 `/tmp/job-fit-turn006-result.json`에 보관했으며 application persistence 기능은 없다. 원시 provider 응답이나 내부 reasoning은 수집하지 않았다.

**Overall: Partial.** 15개 항목, Strong 3 / Partial 7 / Weak 0 / Unknown 5. Turn 005 대비 L012 Azure 복합 환경이 Partial에서 Unknown으로 바뀌었다. 이번에는 evidence 연결이 없어 application Unknown 경로를 거쳤다. 기존 추출·mapping·판단 코드는 유지되었으나 API 결과의 결정론적 동일성을 보장하지 않는다. 한 번의 관측으로 변동 원인이나 확률은 추정하지 않는다.

### Supporting

| ID | Requirement | Match | Source type |
| --- | --- | --- | --- |
| L002 | Server-side engineer 역할 | Strong | context |
| L007 | Backend API 설계·개발 | Strong | responsibility |
| L011 | Java / Spring Boot 환경 | Strong | stack |

### Limiting

| ID | Requirement | Match / Source type | 기존 evidence | Missing / Unknown |
| --- | --- | --- | --- | --- |
| L005 | 부동산펀드 운용사 자사 패키지라는 사업 맥락 | Partial / context | E2 transferable, E6 limited | 해당 specific domain 경험 미입증; 기록 밖 보유 여부 미확인 |
| L006 | Java 8+ 개발 1년 이상 | Partial / required | E1 limited | 버전 또는 기간 미기록; 실제 충족 여부 미확인 |
| L008 | Backend API 운영 | Partial / responsibility | E1 limited, E3 limited | operational context 미입증; 기록 밖 운영 경험 미확인 |

### Non-decisive

| ID | Requirement | Match | Source type |
| --- | --- | --- | --- |
| L009 | 패키지 개선·결함 수정 | Partial | responsibility |
| L010 | 신규 고객 도입·customization | Partial | responsibility |
| L012 | Azure App Service / Linux / Docker / Azure SQL | Unknown | stack |
| L013 | Git / Azure DevOps CI/CD | Unknown | stack |
| L014 | Azure 운영·관리 우대 | Unknown | preferred |
| L015 | 부기 3급 상당 회계 지식 | Partial | preferred |
| L016 | 고객 사용성을 고려한 구현 | Unknown | preferred |
| L017 | 고객 업무 이해·주도성 | Unknown | preferred |
| L018 | Angular / Angular Material / TypeScript | Partial | preferred |

## Findings

1. **Azure Preferred Unknown:** L014는 사후 분류에서 Non-decisive다. 따라서 이번 trace는 Azure 우대 부족을 주요 제한 이유로 설명하지 않는다. 이것이 원래 Overall에서 Azure의 영향이 전혀 없었다는 인과 증거는 아니다.
2. **Java version/duration:** L006은 Limiting이다. 일반 Java 능력은 E1 및 L011 Strong에 남아 있지만 명시 필수의 세부 미확인은 제한 요인으로 설명된다.
3. **Product operation context:** L008 API 운영은 Limiting이다. 그러나 L009 개선과 L010 도입은 Non-decisive여서 제품 관련 Partial이 전부 제한 요인으로 분류되지는 않았다. 별도로 L005 산업 context가 Limiting이다.
4. **Repeated limiting signal:** 같은 Java/버전/기간이 별도 여러 행으로 제한되는 사례는 없다. 제품 관련 운영·개선·도입이 모두 Limiting인 것도 아니다. L005 산업 맥락과 L008 운영 맥락은 연관될 수 있지만 동일 gap이라고 단정할 수 없다. De-duplication/grouping은 하지 않았으며 암묵적 count weighting의 존재 여부도 여전히 미입증이다.
5. **새 검토 후보:** 필수라고 적히지 않은 L005 사업 context가 Limiting이다. 정당한 중요한 맥락 차이인지 도메인 특화 요구의 과대해석인지 planning 검토가 필요하다. 이번에 역할을 보정하지 않았다.
6. 동일 Partial 중 L006/L008은 Limiting, L009/L010/L015/L018은 Non-decisive다. 단순 label 변환으로 trace를 만든 것이 아니다. Requirement category를 보존하여 위 차이를 확인할 수 있다.

## Decisions Made

사용자 승인에 따라 Overall과 trace를 분리하고 기존 네 평가 단계의 요청은 유지했다. 전역 ID를 만들지 않고 기존 L번호를 재사용했다. 모든 항목을 정확히 한 역할에 포함하도록 했고 원문/label/category는 application이 원본에서 가져온다. 사후 attribution 출처를 출력에 명시하고 trace 실패는 전체 CLI 오류로 처리한다. 별도 free-text rationale, reason taxonomy, presentation layer는 추가하지 않았다.

## Changed Files

- `internal/evaluator/types.go`, `evaluator.go`: 출력과 Engine/orchestration 확장.
- `internal/evaluator/trace.go`, `trace_test.go`, `trace_flow_test.go`: domain trace, 검증 및 테스트.
- `internal/jev/trace.go`, `trace_test.go`: 사후 Choice adapter와 테스트.
- `internal/jev/pipeline_test.go`: 기존 pipeline에 trace 단계 검증 추가.
- `README.md`, `docs/evaluation-contract.md`: 현재 출력·호출 수·한계 최소 갱신.
- 본 `docs/development/turn-006.md`. 기존 untracked Turn 005는 이번 변경이 아니다.

## Problems / Limitations

Trace도 추가 모델 판단이므로 변동하거나 최초 선택에 맞춘 사후 설명일 수 있다. 내부 인과관계, 실제 중요도 크기, counterfactual 영향은 알 수 없다. 각 항목은 하나의 역할만 가지므로 동시에 지지하면서 일부 제한하는 복합 의미를 모두 표현하지 못한다. Missing/Unknown 정형 문구의 해상도는 기존 그대로다. Guard는 구조적 오류와 일부 명백한 모순을 막을 뿐 모든 semantic 모순을 검증하지 않는다. 추가 호출 실패 시 기존 Overall만 따로 반환하지 않는다.

## Deferred

Match semantics/기존 evaluation logic 수정, Strong 강제, percentage·weighting·threshold·calibration, grouping/de-duplication, profile/JD 변경, 다른 5개 JD, 반복 안정성 실험, MCP/DB/UI/crawler/parser/result persistence는 수행하지 않았다. Commit은 만들지 않았다.

## Suggested Next Step

Planning conversation에서 L005 사업 context, L006 필수의 세부 미확인, L008 운영 맥락을 핵심 backend 강점과 함께 어떻게 해석할지 검토한다. 특히 context의 Limiting이 기존 semantics에 적절한지 먼저 확인한다. Trace의 사후 attribution 한계를 유지한 채 후속 검증 범위를 정하고, 이번 결과를 근거로 즉시 Strong을 강제하거나 로직을 바꾸지 않는다.

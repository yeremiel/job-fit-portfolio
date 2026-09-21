# Turn 010 — Fix Confirmed Contract Violations

- 날짜: 2026-09-21
- 범위: Turn 009의 Confirmed Violation 1과 3만 수정. Candidate 2는 Ambiguous로 보존.
- 시작 상태: clean `main`, HEAD `00fad1e684a5ce206d1181b68731d5cc922ac2f5`.
- 실제 재검증: CORE와 Nitori 각각 1회, 둘 다 exit 0. 다른 JD 또는 reference에 맞춘 반복 실행 없음.

## Completed

기존 코드 경로를 분석하고 offline 회귀 테스트에서 두 위반을 먼저 재현했다. Strong Overall의 Limiting Choice를 제거하고 domain에서도 거부한다. JD에 없는 version/duration을 rationale로 선택·출력하지 않도록 requirement별 eligibility 검사와 application validation을 추가했다. 기타 qualification gap도 명시된 JD 조건 또는 평가에 직접 필요한 정보에 근거하도록 항목 판단 지시를 제한했다.

전체 테스트·vet·build를 통과한 뒤 동일 CORE/Nitori 입력으로 한 번씩 실행했다. 결과를 보정하거나 reference에 맞추지 않았다. 기존 Turn 000–009와 결과 artifact는 수정하지 않았다.

## Fix 1 — Strong Overall / Limiting Trace

### Root Cause

[기존 Jev trace](../../internal/jev/trace.go)는 실제 `overallMatch`를 state로 받지만 모든 level에 같은 세 Choice를 제공했다. Strong에서도 Limiting을 선택할 수 있었다. [ResolveTrace](../../internal/evaluator/trace.go)는 rows와 refs만 받아 Overall을 알 수 없었으며, ID/enum 및 Supporting 일부 모순만 검사했다.

Turn 009에서 확인한 현재 의미는 “더 높은 Overall Match를 제약하는 중요한 gap”이다. 최고 level Strong보다 높은 match는 없다. 개별 Partial이나 missing evidence가 존재하는 것 자체는 모순이 아니며, 그 항목의 Limiting attribution이 문제였다.

### Change

- Jev `Trace`: Overall Strong이면 해당 요청의 `Limiting` Choice만 제거한다. Supporting과 Non-decisive는 그대로 둔다.
- Domain `ResolveTrace`: 실제 Overall을 인자로 받아 유효한 level인지 확인하고 Strong + Limiting을 거부한다. Orchestration이 기존 선택값을 전달한다.
- Provider가 제거된 Choice를 반환하면 기존 allowed-choice validation에서 실패한다. Domain 검증도 별도로 보호한다.
- 잘못된 trace를 자동으로 Supporting/Non-decisive로 바꾸거나 Overall/requirement Match를 낮추지 않는다. 기존 fail-closed 흐름으로 종료한다.

Limiting의 정의, Overall 선택 질문·state·policy, requirement-level semantics는 바꾸지 않았다.

### Regression Test

- `TestStrongTraceExcludesLimiting`: Strong 요청에 Limiting이 없고 Non-decisive가 존재함을 확인. 정상 Non-decisive 응답은 허용하고 불법 Limiting 응답은 거부.
- `TestStoredCOREStrongLimitingIsRejected`: Turn 008 실제 CORE JSON의 Strong + L011 Limiting을 domain에 전달하면 거부. 테스트 안에서 Non-decisive 대안을 구성하면 허용하고 source assessments/Overall은 변경되지 않음을 확인. 실제 결과 파일을 수정하거나 fallback 구현을 넣은 것이 아니다.
- 기존 trace mapping/참조/Unknown Supporting 검증과 Partial에서 Preferred Unknown Limiting 허용 테스트는 그대로 유지했다. 함수에 Overall을 전달하도록 호출부만 갱신했다.

### CORE Verification

[실제 JSON](turn-010-results/core.json), [동일 입력](../../samples/core.txt).

| 관찰 | Turn 008 | Turn 010 |
| --- | --- | --- |
| Overall | Strong | Partial |
| 항목 수 | 8 | 8 |
| Strong / Partial / Weak / Unknown | 7 / 1 / 0 / 0 | 6 / 2 / 0 / 0 |
| Limiting | L011 | L010, L011 |

- Supporting: L002, L005, L006, L008.
- Limiting: L010 business context, L011 lifecycle scope.
- Non-decisive: L007, L009.

이번 실제 실행은 Partial이므로 **Strong 분기를 직접 검증한 live 사례가 아니다**. Strong + Limiting 차단은 위 모의 응답과 저장된 Strong fixture의 regression tests로 확인했다. 실제 결과에는 모순이 없지만 이를 위해 Overall을 강제로 Partial로 바꾼 것은 아니다. CORE를 Strong이 나올 때까지 재실행하지 않았다.

L010 Public/Energy가 이번에 Partial/scope로 바뀌었으나 Candidate 2를 고친 것으로 간주하지 않는다. 해당 coverage 처리나 direct evidence 기준은 수정하지 않았다. 실행 변동과 qualification 지시 보강의 간접 영향을 이 한 번의 관측으로 분리할 수 없다.

## Fix 2 — Synthetic Version / Duration

### Root Cause

`rationales["duration"]`은 Partial과 `specified version or duration` gap을 함께 정의한다. 기존 `assessmentOptions`는 direct evidence 유무로 Strong만 제한하고 duration은 모든 requirement에 제공했다. Candidate E5의 기간 미기록 정보만 있어도 Jev가 duration을 선택할 수 있었다.

`applyRationale`은 그 선택으로 Match와 정형 reasoning/missingEvidence/unknowns를 만들며 JD에 해당 조건이 있는지 확인하지 않았다. 따라서 provider가 임의 문장을 생성한 문제가 아니라 **무조건 제공된 결합 Choice + application의 무조건 template 적용** 경로에서 synthetic qualification이 발생했다.

### Change

[qualification.go](../../internal/jev/qualification.go)에 현재 requirement 문구만 사용하는 작은 lexical eligibility check를 추가했다. Profile, 다른 requirement 또는 reference에서 조건을 가져오지 않는다.

- 명시 duration 표현(예: `one year`, `3 years`, `18 months`, `years of experience`, 숫자+년/年 등) 또는 명시 version 표현(예: `version 5`, `Java 8+`, `Java 8 or later`, `Python 3.12`)이 인식될 때만 duration Choice를 제공한다.
- 지원하는 부정 표현(`not required`, `not mandatory`, `not documented`, `no specific version`, 일부 기간 불문 표현 등)이 있는 clause는 조건 근거로 삼지 않는다. 세미콜론/개행으로 분리된 다른 긍정 clause는 유지한다.
- `applyRationale`에서도 같은 검사를 하여 Choice 경로를 우회한 잘못된 duration 적용을 거부한다. 오류 전에 Match/설명/배열을 변경하지 않는다.
- 항목 판단 질문에 gap의 출처 제한을 추가했다. Version/duration/years/certification/scale/ownership를 profile의 미기록에서 새 qualification으로 만들지 않고, 명시 JD 조건 또는 해당 평가에 직접 필요한 정보만 사용하도록 한다.

기존 기타 rationale, level mapping, source category, evidence relation 질문, extraction, 전체 compound 평가 지시 및 Overall 정책은 유지했다. 명시 certification/scale/ownership 요구를 삭제하는 규칙도 없다. 이들은 기존 scope/depth/operations/insufficient 등의 선택으로 계속 평가할 수 있다.

**보장 범위:** version/duration template에는 deterministic eligibility 및 거부 검사가 있다. 그 외 qualification의 의미적 grounding은 보강한 지시와 기존 유한 출력 형식에 의존하며 모든 semantic 오판을 자동 검출하는 validator를 새로 만든 것은 아니다. Template에 없는 certification/scale/ownership 세부값을 application이 자유 생성하지는 않는다. 일부 scope/depth 해석이 여전히 지나칠 가능성은 남는다.

### Regression Test

- `TestAssessmentQualificationGrounding`: Angular/TypeScript만 있는 행에는 duration Choice가 없고, provider가 이를 반환하면 오류. Profile에 기간 미기록이 있고 JD 다른 행에 Java version 17/3년이 있어도 현재 행에는 전파하지 않는다.
- 동일 테스트에서 `Java 8+ development for at least one year`에는 duration 선택과 Partial/missing 출력을 허용한다. 실제 Jev API 대신 메모리 transport를 사용한다.
- `TestRationaleRejectsUngroundedDurationWithoutMutatingAssessment`: application-side 직접 적용도 거부하고 명시 조건은 허용.
- `TestQualificationEligibility`: 긍정적인 version/기간, 버전·기간 미기록/불문, 단순 기술·certification·scale/ownership 문구, 일부 영어/일본어/한국어 표현을 확인.
- `TestVerificationInputsHaveNoDurationOption`: 현재 CORE/Nitori의 모든 입력 줄에서 duration이 제공되지 않음을 확인.

수정 전 테스트에서 Strong의 Limiting 제공/수락, Angular/TypeScript의 duration 제공/수락이 실제로 실패했다. 수정 후 기존 테스트와 함께 통과한다.

### Nitori Verification

[실제 JSON](turn-010-results/nitori.json), [동일 입력](../../samples/nitori.txt).

| 관찰 | Turn 008 | Turn 010 |
| --- | --- | --- |
| Overall | Weak | Weak |
| 항목 수 | 11 | 11 |
| Strong / Partial / Weak / Unknown | 3 / 1 / 4 / 3 | 3 / 2 / 3 / 3 |
| L010 Match | Partial | Partial |
| L010 gap | specified version or duration | requested depth |

L010은 E3/E5 limited 근거를 사용하고 depth 미확인을 표시한다. JD에 없는 version/duration gap은 전체 새 Nitori 결과의 reasoning/missing/unknown에 없다. Depth rationale의 보편적인 정답 여부까지 이번 검증으로 확정하지 않는다.

- Supporting: L005, L007.
- Limiting: L002, L006, L008.
- Non-decisive: L009, L010, L011, L012, L013, L014.

Frontend Tech Lead L006은 Weak로 남았고 새 특화 경력을 생성하지 않았다. L009는 Weak에서 Partial로 바뀌었으나 재조정하지 않았다.

## Tests

- 수정 전 targeted regression tests: 의도대로 실패하여 두 경로 재현.
- 최종 `go test ./...`: 통과, 외부 API 없음.
- 최종 `go vet ./...`: 통과.
- 최종 `go build ./cmd/job-fit`: 통과.
- 실제 CORE/Nitori는 각각 1회, retry 없음. 기존 사용자 `.env` 키를 자식 프로세스 환경변수로 전달하고 출력하지 않았다. 기존 5단계 호출 수를 유지한다.
- Live 실행 후 `version 5.0 ... not required` 같은 부정 표현 경계 처리를 보완하고 offline tests를 추가했다. CORE/Nitori에는 애초 version/duration 인식 표현이 없어 전후 Choice가 동일함을 확인했으며 API를 추가 실행하지 않았다. 이는 해당 두 live 결과에 영향을 주는 요청 변경이 아니다.
- 기존 Turn 000–009, profile, 6개 sample, match semantics, runtime policy, 과거 JSON artifact를 hash로 보존 확인. 기존 Overall 함수와 evidence relation/extraction 경로는 원문 비교로 불변 확인.
- 새 JSON의 evidence 원문, trace 참조 및 분포, 문서 링크·format·키 미포함 확인.

## Side Effects / Limitations

- `ResolveTrace` 내부 Go 호출 계약에 Overall 인자가 추가되었다. 외부 CLI JSON schema는 바뀌지 않았다.
- 부적합 Choice를 조용히 수정하지 않으므로 provider가 허용되지 않은 응답을 주면 기존과 같이 non-zero 종료한다.
- 보수적 lexical check는 일반 JD parser가 아니다. 인식하지 못하는 언어·표기나 복잡한 부정/혼합 문장에서 false negative/positive 가능성이 있다. 예를 들어 단순 `Java 8`처럼 version marker/비교 표현 없이 적힌 경우 전용 duration Choice가 제공되지 않을 수 있다. 다른 technology/scope/insufficient 판단 경로는 유지된다. 명시 조건의 지원 범위를 무제한이라고 주장하지 않는다.
- 기타 qualification의 의미적 검증은 완전하지 않다. Numeric scoring/threshold나 새로운 평가 의미를 추가한 것은 아니다.
- 실제 CORE Overall/L010, Nitori L009가 변동했다. Prompt 보강과 모델 variability 중 원인을 확정하지 못하며 label을 다시 맞추지 않았다.
- Turn 008의 동결 실행과 수정 후 실행은 동일 구현의 repeatability 실험이 아니다. 결과 비교는 두 위반의 검증 범위로 제한한다.

## Unchanged / Out of Scope

Candidate 2 Public/Energy Context는 여전히 **Ambiguous**다. Coverage 의미, decomposition, direct evidence 판정 기준 또는 해당 sample 전용 처리를 추가하지 않았다. Candidate profile, Human Reference, match semantics, Required/Preferred 정책, Overall evaluation policy, weighting/threshold/percentage도 변경하지 않았다. 전체 sample 재평가나 productization은 하지 않았다.

## Changed Files

- `internal/evaluator/evaluator.go`, `trace.go`: Overall 전달 및 trace 검증.
- `internal/evaluator/trace_test.go`, `trace_contract_test.go`: 기존 테스트 호출 갱신 및 실제 저장 결과 회귀.
- `internal/jev/trace.go`, `trace_test.go`: Strong의 Choice 제한 및 테스트 호출 갱신.
- `internal/jev/mapper.go`, `qualification.go`: qualification eligibility와 gap grounding 지시/검증.
- `internal/jev/client_test.go`, `contract_regression_test.go`, `qualification_test.go`: 기존 및 신규 회귀 테스트.
- `README.md`, `docs/evaluation-contract.md`: 현재 구현의 최소 설명 추가.
- 본 `docs/development/turn-010.md`, `turn-010-results/core.json`, `turn-010-results/nitori.json`.

## Problems / Findings

Confirmed 위반의 원인은 요청 Choice와 검증 조건에서 재현 가능했다. Fix 1의 실제 Strong branch 증거는 regression fixture/mock이고 live CORE는 Partial이었다는 한계를 분리했다. Candidate 2의 새로운 결과를 해결 증거로 쓰지 않았다. Free-text LLM이나 새 dependency/API 단계를 도입하지 않았다. Commit은 생성하지 않았다.

## Deferred

Candidate 2 semantics 결정, 일반 JD parser/완전한 qualification extractor, 모든 semantic 모순 검증, 전체 sample API regression, 추가 calibration/tuning, MCP/DB/UI/productization은 수행하지 않았다.

## Suggested Next Step

Planning conversation에서 수정된 evaluator를 전체 sample에 한 번 regression 실행할 필요가 있는지 결정한다. 목적은 qualification 제약의 영향과 기존 계약 보존 확인이며 reference 일치가 아니다. 다음 Turn의 추가 tuning이나 productization은 시작하지 않았다.

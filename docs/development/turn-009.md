# Turn 009 — Contract Violation Triage

- 날짜: 2026-09-21
- 범위: Turn 008의 후보 3개만 기존 계약·semantics와 저장된 결과로 판정.
- Application code, prompt, semantics, profile, sample, 기존 artifact 변경 없음. API 호출 없음.

## Completed / Triage Results

| Candidate | Verdict | Reason |
| --- | --- | --- |
| 1 — CORE Overall Strong + Limiting | **Confirmed Violation** | 현재 Limiting은 더 높은 Overall Match를 제약하는 항목이다. 최고 level Strong에 대해 이 의미를 적용한 trace는 모순이다. |
| 2 — CORE Public / Energy Context Strong | **Ambiguous** | Energy evidence는 없지만, 이 Context의 Strong이 두 도메인의 직접 경력 충족을 뜻하는지 폭넓은 업무 맥락의 정합성을 뜻하는지 현재 계약만으로 확정할 수 없다. |
| 3 — Nitori Version / Duration Gap | **Confirmed Violation** | 입력에 없는 “specified version or duration”을 Match와 결합된 rationale로 선택하여 요구된 qualification처럼 설명했다. |

판정은 Human Reference와의 일치 여부를 사용하지 않는다. Confirmed는 관찰된 출력의 명시적 의미 위반이며, 해당 sample의 Overall label 자체가 틀렸거나 다른 level로 바뀌어야 한다는 판정이 아니다. Ambiguous는 정상 동작 인증이 아니다.

## Evidence Reviewed

- [Evaluation Contract](../evaluation-contract.md): Evaluation Principles, Turn 006 Implementation Note, Turn 007 philosophy/quality criteria/change conditions.
- [Match Semantics](../match-semantics.md): Match Level, Missing Evidence, source type 구분, specialized experience 비추정 및 복합 requirement 분해 미정 사항. Sample별 Expected Match 값은 판정 근거에서 제외했다.
- [Turn 006](turn-006.md): Decision Trace Design, Consistency / Failure Handling 및 post-hoc 한계.
- [Turn 008](turn-008.md): 후보의 원래 관찰 기록. 당시 후보 표기를 확정 판정으로 간주하지 않고 재검토했다.
- [CORE 입력](../../samples/core.txt), [CORE 결과](turn-008-results/core.json), [Nitori 입력](../../samples/nitori.txt), [Nitori 결과](turn-008-results/nitori.json), [Candidate Profile](../../data/profile.json).
- [Jev trace](../../internal/jev/trace.go)의 `Trace`, [domain trace](../../internal/evaluator/trace.go)의 `ResolveTrace`, [mapper](../../internal/jev/mapper.go)의 `rationales`, `Assess`, `assessmentOptions`, `applyRationale`, `Overall`, [runtime policy](../../internal/evaluator/instructions.md), [MatchLevel](../../internal/evaluator/types.go).

JD 관련 사실 확인 범위는 당시 보존된 sample 및 JD 검토 기록이다. 공개 사이트를 새로 조회하거나 원문이 현재도 같다고 주장하지 않는다. 결과의 설명은 provider 자유 생성문이 아니라 application의 closed-set template임을 고려했다.

## Candidate 1 — Overall Strong + Limiting

### Observed Behavior

CORE JSON의 `overallMatch`는 `Strong`이다. `decisionTrace.limiting`에는 L011이 있다.

- Requirement: `Responsibility: Work across requirements, system testing, introduction and maintenance.`
- Match: `Partial`; source type: `responsibility`.
- E1/E2/E3/E4 모두 limited. Gap은 `full scope`다.

문제는 개별 Partial이나 gap이 존재한다는 사실이 아니라 그 항목을 **Overall Strong의 Limiting**으로 분류한 것이다.

### Relevant Contract / Semantic Rule

Turn 006의 정의는 “더 높은 match를 제약하는 중요한 gap/uncertainty(Limiting)”이다. 실제 `Trace` 선택지도 다음과 같다.

> An important gap or uncertainty that constrains a higher overall match.

질문은 이미 선택된 Overall에 대한 역할을 분류하도록 한다. 현재 MatchLevel은 Strong/Partial/Weak/Unknown이며 Strong보다 높은 level이나 Strong 내부의 점수/세부 등급은 정의되어 있지 않다. Overall은 단순 평균이 아니므로 개별 Partial이 있다는 이유만으로 Overall Strong을 금지하는 규칙도 없다.

### Analysis

“더 높은 Overall로 올라가지 못하게 하는 중요한 gap”을 이미 최고 level인 Strong에 부여하면 현재 출력 언어에서 제약 대상이 없다. Limiting을 “약점이 있지만 강점에 의해 상쇄된 부정적 요소”라고 읽으면 공존 가능하지만, 이는 현재 문서/질문의 정의보다 넓은 의미다. Confidence 저하나 Strong 내부 수준 제한으로 해석할 수도 없으며 그러한 축은 현재 계약에 없다.

따라서 일반적인 약점과 실제 더 높은 level의 제한은 구별해야 한다. L011 Partial과 Overall Strong의 조합 자체는 허용될 수 있고, 이 판정은 L011의 gap을 지우거나 Overall을 낮추라는 뜻이 아니다.

구현상 `ResolveTrace`는 rows와 references만 받아 참조/일부 Supporting 정합성을 검사한다. Overall level을 인자로 받지 않으며 이 조합을 검증하지 않는다. Turn 006에서 semantic guard의 한계를 기록한 것은 이 출력의 의미를 새로 허용한 것이 아니다. Post-hoc attribution이라는 provenance 역시 정의와의 모순을 면제하지 않는다.

### Verdict

**Confirmed Violation**

### Reason

문서와 실제 Choice가 사용하는 **higher overall match 제한**이라는 동일 정의에 비추어 CORE Strong + Limiting은 trace semantics의 모순이다. 원래 Overall의 causal reasoning이나 정답 여부를 판정한 것이 아니다. 다음 수정 후보는 이 trace 정합성 문제로 한정하며, Limiting을 단순 부정적 영향으로 재정의하지 않는다.

## Candidate 2 — Composite Public / Energy Context

### Observed Behavior

CORE L010:

```text
Business context: Public-sector and energy business systems.
category: context
match: Strong
evidence: E2 transferable, E6 direct
missingEvidence: []
unknowns: []
```

E2는 RDB/SQL과 복잡한 업무 시스템 경험, E6는 ERP·회계·공공 등 도메인 경험이다. Energy-domain 경력은 어느 evidence에도 명시되어 있지 않다. 출력 reasoning은 `Direct documented evidence supports the requirement with no material role/context gap.`이며, trace에서는 Supporting이다.

### Relevant Contract / Semantic Rule

- Evidence First / Unknown Remains Unknown: 없는 경험을 추정하지 않고 미확인 사실은 미확인으로 유지한다.
- Runtime policy: `Do not invent specialized experience` 및 `An unknown fact must remain unknown.`
- Strong: 핵심 requirement에 직접 근거가 있고 material role/context 차이가 크지 않다.
- Match Semantics는 Stack/Context를 전부 필수 경험으로 만들지 않으며, Strong에서도 더 좁은 세부 사실이 Unknown으로 남을 수 있다고 설명한다.
- `Assess`는 `For compound requirements assess the entire line.`을 지시한다. 다만 Context의 접속된 도메인이 각각 필수인 AND 조건인지, 가능한 사업 영역들의 나열인지에 대한 별도 결합 규칙은 없다. 표준 분해 단위도 미정이다.

### Analysis

Public-sector와 energy는 구별되는 두 영역이다. **이를 둘 다 직접 경험해야 하는 독립 조건으로 읽는다면**, E6의 공공 직접 근거만으로 전체 Strong을 주고 energy 미확인을 생략한 것은 문제가 된다. Context라고 해서 specialized experience 비추정 원칙의 예외가 되는 것은 아니다. 전체 line 평가 지시는 일부 coverage만 보고 판단했을 가능성을 검토할 근거다.

그러나 실제 line은 `Required: experience in both ...`가 아니라 사업 context 나열이다. 현재 Strong 정의는 모든 명사의 직접 경험을 요구하지 않고 핵심 역량과 **material** 차이를 기준으로 한다. 공공 경험 및 복잡한 업무 시스템 경험을 바탕으로 넓은 업무 맥락에 직접 정합성이 있고 나머지 차이는 material하지 않다고 판단했을 가능성도 현재 계약만으로 배제할 수 없다. 전체 line을 고려하라는 지시가 Context를 모두 mandatory로 바꾸거나 모든 sub-domain을 직접 경험해야 한다는 규칙은 아니다.

출력에 energy 경력이 있다고 명시한 문장이나 생성된 evidence는 없다. E6 direct가 line 전체의 어느 부분을 가리키는지 분리되지 않으며 Strong template도 개별 도메인 충족을 명시하지 않는다. Empty unknowns는 정보 손실 우려가 크지만, 현재 schema가 모든 하위 사실의 exhaustive coverage를 보장하거나 빈 배열을 모든 조건 검증 완료로 정의하지는 않는다. 따라서 **energy 미확인**은 확인 가능해도, 그것을 실제 충족으로 간주했는지 또는 정합성 요약에서 material하지 않다고 본 것인지는 출력으로 확정할 수 없다.

No Evidence ≠ No Experience는 energy 경험을 있다고 또는 없다고 단정하지 못하게 한다. 이 원칙만으로 전체 Context에 특정 Match label을 강제할 수는 없다. 공개 JD 전체 snapshot이 아닌 요약 표현이라는 입력 한계도 엄격한 논리적 AND를 단정하기 어렵게 한다.

### Verdict

**Ambiguous**

### Reason

과대 coverage 및 미확인 정보 누락의 우려는 있으나, 현재 정의에서 Context Strong이 모든 domain의 직접 경력을 뜻한다고 확정할 근거가 부족하다. 반대로 정상이라고 확정할 만큼 세부 설명도 충분하지 않다. 이 판정은 energy evidence를 인정하거나 Unknown 보존 원칙을 완화하지 않는다. 다음 코드 수정 후보에서는 제외하고, 필요하다면 planning에서 Context의 coverage 의미를 명확히 할 질문으로만 남긴다. Decomposition·semantics는 변경하지 않았다.

## Candidate 3 — Synthetic Version / Duration Gap

### Observed Behavior

Nitori L010 입력은 `Technology environment: Angular and TypeScript for frontend development.`다. 저장된 sample 전체와 기존 JD 검토 기록에는 이 기술의 version 또는 duration 조건이 없다.

출력은 `stack / Partial`, E3 limited + E5 limited다. E5에는 기술 사용 경험과 개별 기간·깊이 미기록이 함께 있다. 출력은 다음과 같다.

```text
reasoning:
Selected rationale: Related experience exists but the specified version or duration is not documented.

missingEvidence:
Not established by the profile: specified version or duration for: Technology environment: Angular and TypeScript for frontend development.

unknowns:
Whether the candidate has specified version or duration beyond the documented evidence.
```

Trace 역할은 Non-decisive, Overall은 Weak다.

### Relevant Contract / Semantic Rule

Requirement Fidelity는 JD 의미를 보존하도록 하고 `Extract`는 `Do not invent requirements`를 명시한다. Unknown Preservation은 미확인 사실을 채우지 말라는 규칙이지, JD에 없는 qualification을 만들어 요구해도 된다는 허용이 아니다.

실제 `rationales["duration"]`은 **Partial level과 “specified version or duration” 설명을 함께** 정의한다. `applyRationale`은 이 선택으로 Match를 설정하고 Missing Evidence/Unknowns 문구도 생성한다. 단순한 부가 설명을 독립적으로 붙이는 구조가 아니다.

### Analysis

“E5에 사용 기간은 기록되지 않았다”는 서술 자체는 사실이다. 하지만 출력은 단순히 그 미확인을 알려주는 데 그치지 않고 **specified** version/duration이 요구되는데 이를 입증하지 못했다는 gap으로 연결한다. 저장된 JD에는 그러한 명시 조건이 없다. Candidate profile의 기간 미기록이 JD의 기간 요구를 만들어낼 근거가 될 수 없다.

따라서 새로운 requirement 행이나 수치 기준을 추가한 것은 아니지만, 기존 기술 항목의 평가 rationale에 없는 qualification을 삽입했다. 이 synthetic qualification은 `duration → Partial`이라는 결합 선택을 통해 **항목 Match를 정당화하는 이유**로 사용되었다. 독립적인 uncertainty note로 볼 수 없다.

다만 이 관찰로 “해당 선택지가 없었다면 Strong이었을 것”이라고 추정하지 않는다. Limited evidence 때문에 Strong 선택지가 이미 제외될 수 있고 다른 rationale로 Partial이 유지될 수도 있다. Non-decisive는 사후 attribution이므로 Overall Weak에 실제 영향이 없었다는 증명도 아니다. 확정 가능한 위반은 **요구되지 않은 버전/기간을 명시된 요건처럼 설명한 requirement-level rationale**에 한정된다.

### Verdict

**Confirmed Violation**

### Reason

입력에 없는 qualification을 “specified” gap으로 표현하고 Match와 결합한 것은 Requirement Fidelity 위반이다. Profile의 기간 미확인 정보를 그대로 보존하는 것과 구별된다. Partial 또는 Overall Weak를 다른 label로 바꾸는 것이 수정 목표는 아니다.

## Confirmed Violations

- Candidate 1: 최고 Overall과 higher-match Limiting attribution의 모순.
- Candidate 3: 미명시 version/duration을 요구된 gap처럼 사용하는 rationale.

두 항목만 다음 제한된 수정 Turn의 후보로 제안한다. 구체적인 구현 방식, prompt 변경, fallback 정책이나 추가 검증 규칙은 이번에 확정하지 않았다.

## Ambiguous Items

Candidate 2. Energy evidence 부재는 확인되지만 Context Strong의 coverage 의미가 충분히 명시되지 않아 위반 또는 정상으로 확정하지 않는다. Human Reference의 label이나 reasoning을 새 규칙 대신 사용하지 않았다. Not a Violation으로 판정한 후보는 없다.

## Problems / Findings

- Turn 008의 candidate 표시는 이번 검토의 출발점이며 결론이 아니었다. 이번 판정은 그 기록을 수정하지 않고 progression으로 남긴다.
- 구조 검증을 통과했다고 semantics까지 보장되는 것은 아니다. 반대로 validation이 없다는 사실만으로 위반을 확정한 것도 아니다. 위반 판정은 현재 정의와 실제 출력의 모순에 근거한다.
- Post-hoc trace로 최초 판단의 인과관계를 확정할 수 없고, closed-set 문구로 Context의 상세 coverage를 복원할 수도 없다.

## Validation / Changed Files

`docs/development/turn-009.md`만 생성했다. 기존 미커밋 작업은 보존했다. 기존 문서·Go 코드·정책·profile·sample·결과 JSON의 작업 전후 SHA-256이 같음을 확인했다. 두 artifact의 관련 ID/Match/trace/문구와 문서 판정 근거를 대조하고 상대 링크·공백·코드블록을 검증했다. 실행 동작 변경이 없어 Go test/build를 반복하지 않았고 Jev API 또는 다른 evaluator도 호출하지 않았다. Commit은 생성하지 않았다.

## Deferred

두 Confirmed 항목의 실제 수정·테스트 추가, Candidate 2 coverage semantics 결정, decomposition 변경, tuning/calibration, API 재실행, 새 JD, weighting/threshold/percentage 및 MCP는 수행하지 않았다.

## Suggested Next Step

Planning에서 Candidate 1과 Candidate 3만 다루는 제한된 수정 Turn을 정한다. 성공 기준은 기존 trace 정의와 requirement fidelity 준수이며 Human Reference 또는 특정 Overall에 맞추는 것이 아니다. Candidate 2는 코드 수정 대상에서 제외하고 semantics clarification 질문으로 보존한다. 다음 Turn의 수정 작업은 시작하지 않았다.

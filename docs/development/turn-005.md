# Turn 005 — First Evaluation Diagnosis

- 날짜: 2026-09-21
- 범위: Prop Tech plus 한 건의 진단. 코드·입력·정책·baseline 변경 없음.
- 기준 commit: `d5145a7d4095d0df7081f4f7f98750143974ea3f` (시작 시 clean `main`).

## Completed

동일 profile/JD로 실제 evaluation을 한 번 재실행하고, 출력된 15개 항목 전체를 기존 Human Baseline 및 구현과 대조했다. Overall **Partial**은 재현되었으나 회계 지식은 **Weak → Partial**로 바뀌었다. 따라서 회계 Weak 하나를 Overall 차이의 원인으로 단정할 수 없다.

가장 유력한 설명은 **핵심 역량의 직접 근거와 세부 미확인을 Overall에서 해석하는 방식의 차이**다. 다만 현재 Overall 설명에는 결정적 requirement가 없으므로, 중요도 위반이나 숨은 다수결을 입증한 것은 아니다. 기존 baseline을 정답으로 가정하거나 Jev 결과에 맞춰 수정하지 않았다.

## Sources / Reproduction

- [Turn 004 기록](turn-004.md): 당시 성공 결과와 알려진 제약.
- [Match Semantics, Sample 02](../match-semantics.md#sample-02--prop-tech-plus): Turn 003 확정 baseline 7행과 Overall Strong. Codex 초안에서 출발한 비교 기준이며 독립 검증된 ground truth는 아니다.
- [Profile](../../data/profile.json), [고정 sample](../../samples/proptech-plus.txt): 이번 실행의 실제 입력. JD 재수집 없이 같은 요약 text 사용.
- [Runtime policy](../../internal/evaluator/instructions.md), [orchestration](../../internal/evaluator/evaluator.go), [Jev mapping](../../internal/jev/mapper.go), [HTTP client](../../internal/jev/client.go): 진단 근거. 회사별 baseline은 runtime에 전달되지 않는다.

기존 코드를 `/tmp/job-fit-turn005`로 빌드하여 다음 인자로 실행했다.

```sh
/tmp/job-fit-turn005 --profile data/profile.json --job samples/proptech-plus.txt
```

사용자 제공 `.env`에서 키를 출력 없이 읽어 자식 프로세스 환경변수로 전달했다. CLI의 `.env` 지원을 추가하지 않았다. 요청 모델은 기존 `jev-1.13.0`, endpoint도 기존 TypeSafe System One API다. 한 번의 CLI 실행은 추출 → evidence 관계 → 항목 판단 → Overall의 4단계 호출 경로를 거치며 자동 retry가 없다. 별도 반복·다른 JD 호출은 하지 않았다.

진단용 stdout JSON은 repository 밖 `/tmp/job-fit-turn005-result.json`에 제한된 파일 권한으로 임시 보관했다. 이는 CLI 결과 저장 기능 추가가 아니며 영구 재현 자료는 아래 표다. Provider 원시 응답이나 내부 추론은 수집하지 않았다.

| 보존 대상 | SHA-256 |
| --- | --- |
| `data/profile.json` | `cae86ad4906124b0eee8283c05bc74d439543f75d9b2e7496f8d6895e32f09cd` |
| `samples/proptech-plus.txt` | `97068536167d4b21fc9d4f53bc92deedaee922cc98a08c7275e774049829fd19` |
| `internal/evaluator/instructions.md` | `f919adc667e09c3bfd1d7f35e94e3b04b0e67ed9ffdb84ca46785e997cc0ee2c` |

## Reproduced Result

| 관찰 | Turn 004 성공 실행 | Turn 005 재실행 |
| --- | --- | --- |
| Overall | Partial | Partial |
| Requirement count | 15 | 15 |
| Strong | 3 | 3 |
| Partial | 7 | 8 |
| Weak | 1 | 0 |
| Unknown | 4 | 4 |

이번 실행은 exit 0, 유효한 domain JSON을 반환했다.

- Strong: L002, L007, L011
- Partial: L005, L006, L008, L009, L010, L012, L015, L018
- Weak: 없음
- Unknown: L013, L014, L016, L017

회계 지식 L015는 Turn 004 기록의 Weak와 달리 Partial로 기존 baseline에 일치한다. 당시 문서에 명시된 Java 기간 Partial, API 설계 Strong, Azure 운영 Unknown도 이번에 동일하다. Turn 004 문서에는 모든 evidence 관계의 원본 snapshot이 없으므로, 과거 관계 선택까지 동일했다고 주장하지 않는다. 같은 입력의 관측 결과에 변동이 있었으며 결과를 결정론적으로 보장할 수 없다. 두 성공 관측만으로 변동률, 원인, 모델 내부 sampling 방식 또는 장기 안정성을 추정하지 않는다.

## Requirement-level Diagnostic Table

**H**는 기존 확정 baseline, **R**은 이번 진단자의 semantics 해석이다. R은 새 human label 확정이나 baseline 수정이 아니다. “H 없음”인 행을 기존 불일치로 집계하지 않는다. Source Type의 Other는 실제 출력 `context`다. `d/t/l`은 각각 실제 연결의 `direct/transferable/limited`이며, 나열하지 않은 evidence는 이번 출력에서 연결되지 않았다.

Evidence 요약: E1 전문 Java/Spring/Boot·REST API·약 16년(버전별 기간 미기록), E2 RDB/SQL·복잡한 업무 로직, E3 요구 분석·API/DB 설계·legacy 개선, E4 리딩·고객 요구/명세 조율, E5 Angular/TypeScript 등 사용(깊이·기간 미기록), E6 ERP/회계 등 도메인, E7 기본 AWS. 요약은 읽기 편의를 위한 것이며 실제 출력은 profile 원문을 복사한다.

Missing / Unknown 열은 **출력의 gap 문구를 축약한 것**이다. Strong 이외에는 “profile이 해당 gap을 입증하지 않음 / 기록 밖에서 보유하는지는 미확인”이라는 정형 문장이다. 독립적으로 추출된 새로운 사실이 아니다. Review의 추가 미확인은 진단자의 발견으로 구분한다.

| ID / Requirement (실제 JD line) | Source Type | Human Baseline | Jev Result | Candidate Evidence | Missing / Unknown (출력) | Review |
| --- | --- | --- | --- | --- | --- | --- |
| L002 Position: Server-side engineer | Other / context | H 독립행 없음; R Strong (기존 핵심 backend 해석) | Strong | E1d, E2t, E3d, E4l, E5l | 없음 | 일반 backend 역할과 직접 경력이 일치. 역할 identity이지 별도 필수 경력 조건은 아니다. |
| L005 Business context: In-house package development for real-estate fund management companies. | Other / context | H 독립행 없음; R Partial | Partial | E2t, E6l | specific domain experience | 업무 시스템의 전이는 인정, 부동산펀드 직접 경험은 미확인. 이 domain을 입사 전 필수로 바꾸면 안 된다. |
| L006 Required: At least one year of development experience with Java 8 or later. | Required | H Partial | Partial | E1l | specified version or duration | 일치. Java 전문 경력은 명확하지만 버전·기간을 전체 경력으로 추론할 수 없다. E1을 통째로 limited로 표현하여 직접 Java 역량과 미확인 세부가 분리되지 않는다. |
| L007 Responsibility: Design and develop backend APIs. | Responsibility | H Strong | Strong | E1d, E2l, E3d | 없음 | 일치. E1/E3 직접 API 근거를 연결했다. E2의 RDB는 API 전부를 단독 입증하지 않으므로 limited만으로 오류라 할 수 없다. |
| L008 Responsibility: Operate backend APIs. | Responsibility | H 운영·개선·도입 묶음 Partial; R 이 행 Partial | Partial | E1l, E3l | operational context | 묶음 baseline과 양립. 개발 경력을 production 운영 책임으로 바꾸지 않았다. 개별 운영 행의 과거 human label은 없다. |
| L009 Responsibility: Enhance the package and fix defects. | Responsibility | H 같은 묶음 Partial; R Partial | Partial | E1l, E2l, E3l, E4l | full scope | Legacy 개선 E3가 연결되어 누락은 아니다. 일반 개선의 직접성과 해당 제품 맥락을 구분하지 못한다. limited 대신 transferable도 가능한 해석이나 이 행 Strong이 확정되는 것은 아니다. |
| L010 Responsibility: Set up and customize the package for new customers. | Responsibility | H 같은 묶음 Partial; R Partial | Partial | E1l, E2l, E3l, E4l, E6l | full scope | 고객 조율 E4 및 업무 시스템 E2/E6를 사용. 구체적인 도입 책임은 미기록. H의 한 gap이 세 행으로 표현되었다. |
| L011 Technology environment: Java and Spring Boot. | Stack | H Spring Boot·RDB Strong의 일부; R Strong | Strong | E1d, E3l | 없음 | Boot 직접 근거는 유지. 이 행에는 RDB가 없으므로 H 전체 축과 1:1 비교는 불가. |
| L012 Technology environment: Azure App Service on Linux/Docker and Azure SQL Database. | Stack | H 일반 Spring Boot·RDB Strong, Azure SQL 미확인; R 이 복합행 Partial | Partial | E2t | exact technology experience | 일반 RDB가 특정 Azure 실행환경과 묶였다. H Strong과의 표면 차이는 비교 단위 차이(D). RDB 전이는 합리적이나 App Service/Linux/Docker 각각은 판단 근거가 부족하다는 세부가 출력에 없다. |
| L013 Technology environment: Git and Azure DevOps for CI/CD. | Stack | H 없음; R Unknown | Unknown | 없음 | sufficient relevant evidence | Git/DevOps/CI/CD 명시 근거가 없다. 경력 연수나 코드 리뷰로 도구 사용을 추론하지 않았다. Application이 연결 없음으로 Unknown 처리. |
| L014 Preferred: Hands-on Azure operation/administration experience. | Preferred | H Unknown | Unknown | 없음 | sufficient relevant evidence | 일치. H는 E7 기본 AWS를 한계 설명에 인용하나 실제 Azure 근거로 인정하지 않는다. E7 미연결은 유효한 직접 근거 누락이라고 볼 수 없다. Application Unknown. |
| L015 Preferred: Accounting knowledge comparable to bookkeeping level 3. | Preferred | H Partial | Partial | E6l | requested depth | 이번에는 일치. 도메인 경험과 부기 3급 상당 지식의 차이를 남김. T004 Weak는 아래 A/B/E 가설로 검토하며 확정 오류로 판정하지 않는다. |
| L016 Preferred: Program with the customer's usability in mind. | Preferred | H 없음; R Unknown 또는 Partial 검토 대상 | Unknown | E3l, E4l | sufficient relevant evidence | 요구 분석/고객 조율은 관련되나 실제 usability 고려 사례는 없다. 전이 근거를 충분히 인정할지(E)가 남음. Evidence 연결과 평가 가능성은 같지 않다. |
| L017 Preferred: Willingness to understand the customer's business and work proactively. | Preferred | H 없음; R 전체 행 Unknown, 업무 이해 부분 Partial 검토 대상 | Unknown | E3l, E4l, E6l | sufficient relevant evidence | 업무 이해 근거와 proactive 성향을 한 행에 결합했다. 리딩 경험만으로 의지·성향을 확정하지 않는 것은 타당. D/E 후보이며 기존 baseline mismatch는 아니다. |
| L018 Optional adjacent web development: Angular, Angular Material and TypeScript. | Preferred (원문 Optional adjacent) | H Angular/TypeScript Partial; R 이 행 Partial | Partial | E3l, E5l | requested depth | Label은 일치. Angular Material은 profile에 없지만 출력 gap은 depth만 표시한다. 진단자가 발견한 tool 미확인이 rationale 하나에 가려진 사례. |

L001 회사명, L003 출처, L004 검토일은 제외되었다. Category 수는 Required 1, Responsibility 4, Stack 3, Preferred 5, Other/context 2다. Preferred에 Optional adjacent L018을 포함하는 분류는 필수 승격이 아니다.

## Human vs Jev Mismatches

분류: A human baseline이 낙관적일 가능성, B Jev가 엄격할 가능성, C evidence mapping, D decomposition, E semantics ambiguity. 복수 가설은 인과 확정이 아니다.

| 차이 | 분류 / 판단 | 근거와 제한 |
| --- | --- | --- |
| Overall H Strong / 두 실행 Partial | E 우선; A/B Possible, D 기여 Possible | 사람은 핵심 backend/API/Boot/RDB를 중심으로 세부 미확인을 남겨도 Strong. Runtime은 모든 행을 보고 중요한 gap을 판단하지만 결정적 행을 공개하지 않는다. 어느 쪽이 과한지 확정 불가. |
| T004 회계 H Partial / Weak | A/B/E Possible; C 원인 미확정 | ERP/회계 경험으로 지식 수준 전이를 인정하는 H가 낙관적인지, 미확인을 큰 gap으로 판단한 Jev가 엄격한지 분리 필요. 이번에는 같은 E6 근거로 Partial이며 지속적 불일치가 아니다. |
| H Boot/RDB Strong / L011 Strong + L012 Partial | D Observed; label 오류 미확정 | 일반 역량과 특정 플랫폼의 비교 단위가 달라졌다. L012 전체를 H Strong으로 올려야 한다는 결론은 부적절하다. |

L008–L010은 기존 Partial 묶음의 분할이며 세 개의 독립 human label과 일치했다고 주장하지 않는다. L005/L013/L016/L017은 기존 개별 label이 없어 신규 차이로 단정할 수 없다. L016/L017의 C/D/E 후보는 진단자의 검토 의견이다. 비교 가능한 기존 항목들은 이번에 대체로 양립하므로 **Case A(Overall 해석 우선 검토)**에 가깝다. 다만 decomposition 차이 때문에 모든 requirement가 동일하게 정렬되었다고 할 수 없고, Overall 수정부터 바로 시작할 근거도 충분하지 않다.

## Diagnosis

### Requirement Decomposition

**Observed:** `Extract`는 의미 단위로 새로 분해하지 않고 비어 있지 않은 JD 줄을 분류·보존한다. 따라서 15행은 sample 작성 단계의 줄 구분에도 의존한다. Java/Java 8/기간을 세 행으로 자동 분해한 사례는 없다. 대신 L002 역할, L006 Java 조건, L007 API, L011 Java/Boot가 겹치는 backend 강점을 표현한다.

H의 운영·개선·도입 한 행은 L008–L010 세 행으로 나뉜다. 서로 다른 업무여서 분리 자체가 오류는 아니지만, 제품 맥락 미확인 하나가 반복해서 나타난다. L005 domain, L010 package 도입과도 관련된다. 이는 **over-decomposition 및 duplicate influence의 후보**다. 반복은 Strong 쪽에도 있으므로 Partial을 늘리는 방향만 가정할 수 없다.

**Under-decomposition observed:** L012 플랫폼/OS/container/DB, L013 Git/CI-CD 플랫폼, L017 업무 이해/주도성, L018 Angular/Material/TypeScript가 한 줄에서 판단된다. L006도 직접 Java 능력과 버전·기간 미확인이 단일 관계·label로 압축된다. 모든 복합행을 쪼개면 또 개수 문제가 생기므로 분해 규칙의 변경 방향은 planning 검토로 남긴다.

**Requirement inflation:** 모든 비메타 행이 `requirements` 배열에서 평가되지만, 실제 category에 Required로 승격된 Responsibility/Stack/Preferred/Context는 **없다**. L005 사업 context와 L013 환경 Unknown 등이 Overall 입력을 확장한 것은 확인되지만, 실제 필수처럼 불이익을 줬는지는 현재 설명으로 확인할 수 없다.

### Evidence Mapping

모든 선택된 행과 E1–E7의 조합을 평가하므로 일부 evidence가 검색 후보에서 빠지는 구조는 아니다. 이후 Jev의 `unrelated` 선택으로 제외될 수 있다. 이번 핵심 API/Java/Boot의 직접 근거는 연결되었고, 명백히 존재하는 Azure 운영·Git/DevOps 경력을 놓쳤다는 증거는 없다(해당 경력이 profile에 없다).

L009의 E3 legacy 개선, L010의 E4 고객 조율, L017의 E3/E4/E6는 실제로 연결되어 있다. 다만 전부 limited로 표시하여 일반 역량의 직접성/전이와 해당 제품·성향의 미확인을 구별하지 못한다. **Overly weak relation Possible**, 그것이 잘못된 최종 label을 야기했다고 입증하지는 못했다. L005에는 E3 설계가 미연결이지만 특정 산업을 뒷받침하는 정도는 간접적이고 E2/E6가 남아 있어 핵심 누락이라고 단정하지 않는다.

**Unsupported evidence text Not observed:** 모든 출력 ID/text를 profile 원문과 대조했다. 새로운 경력·버전·기간을 만들어 쓴 항목은 없다. 다만 원문 복사는 relevance 보장이 아니다. E4 리딩/리뷰를 L009의 limited 근거로 삼는 선택 등은 약한 연관일 수 있으며, 복사된 긴 문단 전체가 해당 requirement를 지지하는 것으로 읽히지 않도록 검토해야 한다.

E1의 전문 Java 직접성 + 버전/기간 미확인, E2의 일반 RDB 직접성 + Azure SQL 미확인을 두 층으로 설명할 필요성이 드러났다. 현재는 coarse evidence 하나당 relation 하나이므로 이 차이를 보존하기 어렵다. 새로운 schema나 mapping 규칙은 만들지 않았다.

### Jev Judgment

L006 Partial, L007/L011 Strong, L014 Unknown, L015/L018 Partial은 기존 판단과 양립한다. 운영/제품 scope의 Partial도 근거가 있다. 이번에 광범위한 requirement-level 오판은 확인되지 않았다.

T004 회계 Weak는 미확인을 큰 gap으로 해석했는지 의문이지만, 이번 Partial로 바뀌어 **Jev가 회계 전이를 항상 과소평가한다**는 주장은 지지되지 않는다. L016/L017은 관련 경험이 있어도 구체적인 usability/태도 판단에 충분한지 애매하다. Unknown을 무경험으로 바꾸거나 Partial을 강제하지 않는다.

선택지는 level과 rationale이 결합되어 있다. Partial에는 scope/depth/technology/domain/operations/duration의 6개 선택지, Strong 1개, Weak 2개, Unknown 1개다. 선택지 표현이 판단에 영향을 줄 가능성은 있으나 선택지 수만으로 확률적 편향을 입증할 수 없다. 이 Turn에서는 변경·반복 실험을 하지 않았다.

### Overall Aggregation / Required and Preferred Influence

`Evaluate`의 summary는 출력용 ID 목록이며 `Overall`에는 summary나 count가 전달되지 않는다. 실제 rows, profile, JD, policy를 전달하고 Jev가 Overall Choice를 반환한다. **Application의 평균·다수결·숫자 weighting은 없음**을 코드로 확인했다. 세부 label을 숫자로 바꾸는 계산도 없다.

Category와 “Required를 우선, Preferred gap만으로 크게 낮추지 않음”이라는 지시는 보존된다. 그러나 핵심 capability 묶음, critical gap, 최종 판단에 사용한 행은 별도 결과로 식별되지 않는다. **표현·지시 수준의 중요도 보존은 확인, 실제 판단에서의 보존 여부는 미확인**이다.

관측은 단순한 `Core Strong + Preferred Unknown`만의 사례가 아니다. 유일한 명시 Required인 L006도 Partial이고, 운영·제품 업무와 domain/environment에도 Partial이 있다. 사람 baseline은 이미 이런 미확인을 알고도 Strong을 유지했다. Runtime의 Strong 기준은 “core role and major Required capabilities have direct evidence and no material role/context gap”이고, Partial 기준은 “some important ... gaps remain”이다. 버전·기간 미확인 또는 운영 맥락을 **중요한 gap으로 볼지**에서 다른 해석이 가능하다. 이 문구만으로 현재 Jev가 어떤 행을 중요하게 봤는지는 알 수 없다.

따라서 Preferred Unknown 때문에 부당하게 내려갔다거나, 8개의 Partial이 많아서 Partial이 되었다고 결론내릴 수 없다. 회계 Weak가 사라져도 Overall이 유지된 사실은 **그 Weak가 유일한 원인이라는 설명을 약화**하지만, 통제된 ablation이 아니므로 회계가 전혀 영향을 주지 않는다는 증명도 아니다.

**Is requirement count implicitly acting as weighting?** 코드상 명시적 count weighting은 **Not supported**. 반복된 행의 수·표현이 모델의 중요도 해석에 영향을 주는 것은 **Possible, unproven**. 증명하려면 동일한 의미를 보존한 묶음/표현 비교를 별도 설계해야 한다. 이번에 수행하지 않았다.

### Explainability

| 사람이 알고 싶은 것 | 현재 출력으로 알 수 있는 것 | 부족한 정보 |
| --- | --- | --- |
| 왜 Strong인가? | 직접 연결 evidence 원문과 일반적 Strong 문구 | 어떤 문구가 어느 부분을 충족하는지, 복합행의 전체 충족 여부 |
| 왜 Partial인가? | 선택한 주된 gap 종류와 관련 evidence | 각 하위 조건의 충족/미확인 구분, 다른 label보다 적절한 근거 |
| 무엇이 missing인가? | “requested depth”, “full scope” 같은 정형 문장 | 확인할 구체적 사실/사례. L018 Material 미기록처럼 선택하지 않은 gap |
| 왜 Overall Partial인가? | Partial 정의와 모든 행의 label/category 나열 | 결정적 requirement, 중요도 판단, core 강점과 gap의 관계, Preferred 영향 |

`applyRationale`이 missingEvidence/unknowns를 template로 만들며, `Overall` 역시 선택지 설명에 행 목록을 붙인다. 이것은 Jev의 자유 설명 또는 내부 reasoning이 아니다. 특히 “full scope”를 반복한 문장은 원인 진단에 충분하지 않다. 이 한계 때문에 Overall 중요도 해석이 유력한 후보여도 확정 root cause로 선언할 수 없다.

### Consistency Guard

| 검증 대상 | 현재 보장 | 한계 / 과보수성 |
| --- | --- | --- |
| Invalid Strong | direct relation 없으면 Strong 선택지를 제거하고 `applyRationale`에서도 거부 | direct 자체가 잘못 분류되면 정상 Strong도 막힐 수 있다. 잘못된 direct 하나면 복합행 전체 Strong을 허용할 수 있다. |
| Evidence 없음 | 연결이 없으면 application Unknown, 후속 항목 질문 생략 | mapper false negative를 확정한다. Profile에 실제로 근거가 없는지 재검증하지 않는다. |
| Unsupported evidence | ID/text를 profile에서 복사 | 잘못된 연관성/과장된 direct는 막지 못한다. Text provenance와 semantic validity는 다르다. |
| Contradictory result | 허용하지 않은 Choice, source requirement 불일치, non-Unknown 무근거 등을 거부 | Core mismatch와 Overall Strong의 모든 의미적 모순을 막는 guard는 없다. Direct+Unknown도 자동 오류는 아니다. |

이번 direct 연결은 L002/L007/L011에 있고 모두 Strong이다. 다른 12행은 Strong 선택지가 차단되지만, 그 행들이 정상 Strong이어야 한다는 근거는 확정되지 않았다. 특히 L006은 버전·기간이 미확인이어서 전체 조건 Strong 차단이 baseline과도 일치한다. L009의 일반 개선 직접성이 limited로 손실될 가능성은 있으나 제품 scope까지 Strong으로 만들 근거는 부족하다. **실제 false-negative Strong은 미입증, 구조적 가능성은 있음**이다.

기존 `TestAssessmentOptionsRespectEvidence`, `TestRationaleMapping`, `TestNoLinkedEvidenceRemainsUnknown`은 위 분기와 거부 동작을 검증한다. `TestVerticalSliceOffline`의 Strong Overall + Weak Preferred는 그러한 결과를 application이 허용함을 보여줄 뿐 실제 Jev의 중요도 판단이 옳다는 증명이 아니다. 테스트를 수정하거나 새로운 규칙을 추가하지 않았다.

## Diagnostic Conclusion

### Primary Cause — Likely, not proven

**Overall 중요도 해석의 차이:** Human Baseline은 직접적인 backend/API/Java/Boot/RDB 강점을 중심으로 세부 미확인을 남기고 Strong을 허용한다. 현재 실행은 그 미확인을 포함한 여러 행에서 Partial을 반환한다. 항목 수준의 비교 가능한 판단은 대체로 양립하므로, Overall에서 어떤 gap을 중요하다고 보는지가 가장 유력한 원인 후보다. 다만 결정적 행을 설명하지 않아 evaluator의 잘못인지, baseline이 낙관적인지 확정할 수 없다.

### Secondary Causes

- **Observed representation difference / Possible causal contribution:** 7행 baseline과 15행 평가의 단위 차이. 운영 scope의 반복, 일반 RDB와 Azure platform의 결합으로 핵심 역량과 세부 gap의 표현이 달라진다.
- **Observed granularity limitation / Possible contribution:** 하나의 evidence relation과 하나의 rationale로 직접 역량/미확인 세부를 함께 압축한다. Missing 정보도 일반 문구여서 영향 구분이 어렵다.
- **Observed variability:** 회계 Weak가 Partial로 바뀌었다. 세부 판단의 안정성 한계이며 Overall 불일치의 단독 원인으로는 지지되지 않는다.
- **Observed diagnosis limitation:** Overall 설명이 정의+목록이라 원인 확정을 막는다. 설명 부족 자체가 label을 낮췄다는 뜻은 아니다.

### Not Observed / Not Supported

명시적 count-based aggregation, category의 Required 승격, 새로운 candidate evidence text 생성, 핵심 API/Java/Boot 직접 evidence의 중대한 누락은 관찰되지 않았다. Preferred 부족만으로 Overall이 낮아졌다는 주장, guard가 이번 정상 Strong을 막았다는 주장, Human Baseline 또는 Jev가 확실히 틀렸다는 주장은 현재 증거로 지지되지 않는다.

## Decisions Made

관찰 사실과 원인 가설을 분리하고 H/R provenance를 표에 표시했다. 재현은 한 번으로 제한하고 기존 일곱 baseline 행을 그대로 유지했다. 진단 자료는 본 문서 하나에 모았다. Semantics ambiguity는 planning 검토 대상으로 남겼으며 동작 변경은 결정하지 않았다.

## Validation / Problems and Findings

- 변경 전후 기존 tracked 파일 전체의 SHA-256을 대조하여 동일함을 확인했다. 신규 파일은 본 문서 하나다.
- 기존 `go test ./...`, `go vet ./...`, CLI build 통과. 테스트는 외부 API 없이 수행한다.
- 실제 stdout의 15개 고유 ID, JD 원문 일치, category/level 분포, summary와 행의 대응, evidence ID/text 원문 일치를 확인했다.
- 문서의 15행이 실제 결과 전체와 대응하는지, relative link와 whitespace를 확인했다. 키 값이 문서에 포함되지 않았고 `.env`는 ignored 상태다.
- T004는 요약 기록이라 모든 provider 선택을 재구성할 수 없다. 이번에도 provider 내부 추론이나 중요도 근거는 출력되지 않는다. 이는 확정 인과 진단의 한계다.
- 분석을 막는 application bug는 발견하지 않았다. 기존 baseline의 Strong을 자동 정답으로 간주하지 않았다.

## Changed Files

- `docs/development/turn-005.md` 생성. 기존 파일 변경 및 commit 생성 없음.

## Deferred

평가 로직·policy·prompt·mapping·Overall·CLI·schema·profile·JD 수정, calibration, Match percentage/weighting formula/threshold, 추가 sample 평가, 통계적 반복·ablation, MCP/DB/UI/crawler 및 결과 저장 기능은 수행하지 않았다.

## Suggested Next Step

Planning conversation에서 먼저 **명시 Required의 버전·기간 미확인과 제품 운영 scope 미확인을 Overall에 어떻게 해석할지**를 검토한다. 이어 일반 역량/세부 조건과 반복 capability를 구분하는 비교 단위를 정하고, 결정적 requirement 및 그 중요도를 추적할 설명에 무엇이 필요한지 다음 Turn 범위로 선택한다. L016/L017의 행동·성향 우대에 대한 근거 충분성도 열린 질문이다. 숫자 규칙이나 특정 sample의 Strong을 강제하는 calibration을 전제로 삼지 않는다. 다음 변경이나 실험은 시작하지 않았다.

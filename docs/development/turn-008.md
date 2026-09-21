# Turn 008 — Evaluate Frozen Baseline Across Samples

- 날짜: 2026-09-21
- 범위: 동결 evaluator로 새 sample 5개를 각각 1회 평가. Prop Tech plus는 Turn 006 결과 재사용.
- 결과: 새 실행 5개 모두 exit 0. Code/policy/profile/semantics/validation 변경, tuning 또는 재시도 없음.

## Frozen Baseline Verification

Turn 007의 fingerprint `845acdfb12ebb507598e88b9fd60ce1759525750c765609a3e1dd3ff2968f425`와 실행 전 현재 작업 트리가 일치했다. 산출 범위/방식은 [Turn 007](turn-007.md)의 기록과 같다. HEAD `d5145a7d4095d0df7081f4f7f98750143974ea3f`만으로는 Turn 006 구현을 식별할 수 없으며 기존 미커밋 코드까지 포함한 기준이다. 기존 변경사항을 그대로 보존했다.

## Input Preparation / Provenance

[기존 JD 검토 기록](../match-semantics.md)의 source URL, reviewed date와 requirement 요약을 사용해 다섯 text file을 준비했다. 새 사이트 수집이나 현재 원문 변경분을 섞지 않았다. 각 파일에 Turn 002 검토일 2026-09-21, URL, 전체 JD가 아닌 당시 기록의 영문 paraphrase임을 명시했다. Candidate evidence, reference label, 당시 evaluator 예상 reasoning은 입력에 넣지 않았다. 결과를 본 뒤 입력을 수정하지 않았다.

입력의 줄 구분은 기존 sample 요약의 항목을 기준으로 했으며 runtime decomposition은 바꾸지 않았다. BlueMeme는 OutSystems 사전 경험 불문/교육을 명시했고, PIA는 AWS 구축·운영 3년 이상을 Required로, Nitori는 frontend Tech Lead를 Required로 보존했다. CORE는 leader 채용 분기다. teamLab은 기술 목록이 전부 필수라는 의미로 바꾸지 않았고 일본어 Required도 입력에 포함했다.

**제한:** 전체 원문 snapshot이 아닌 보존된 요약의 재표현이다. 따라서 이번 Requirement Fidelity 검토는 입력/기존 검토 기록에 대한 것이며 공개 JD 전체를 완전하게 재현했다는 보장이 아니다. CORE 학력의 세부 조건 등 기존 기록에 상세 값이 없는 비기술 조건은 새로 추측하지 않았다. 요약/언어/줄 분할의 영향과 evaluator 자체의 영향을 완전히 분리할 수 없다.

## Runs / Artifacts

기존 `./job-fit --profile data/profile.json --job samples/<sample>.txt` CLI를 사용했다. 사용자 `.env`의 키를 자식 프로세스 환경변수로만 전달했다. 키를 문서·artifact에 넣지 않았다. 새 실행은 BlueMeme → PIA → CORE → teamLab → Nitori 순서로 각각 1회이며, 각 성공 실행은 기존 5단계 HTTP 경로를 거친다. 다른 evaluator나 통계적 반복 호출은 없다.

[결과 manifest](turn-008-results/manifest.json)에 frozen fingerprint, profile/input/result SHA-256 및 각 결과의 출처를 남겼다. JSON들은 **CLI의 raw structured domain result**이며 provider 원시 응답은 아니다. 단순 평가 artifact 저장이며 application persistence/DB 기능을 추가하지 않았다. Prop Tech plus JSON은 남아 있던 `/tmp/job-fit-turn006-result.json`을 당시 Turn 006 요약과 대조한 뒤 byte-for-byte 복사했다. 재실행하지 않았다.

## Sample Results / Reference Comparison

S/P/W/U는 Strong/Partial/Weak/Unknown의 관찰 개수다. Relationship은 label 관찰이며 accuracy 또는 failure 판정이 아니다.

| Sample | Pre-evaluation Reference | Jev Overall | Relationship | Requirements | S / P / W / U |
| --- | --- | --- | --- | --- | --- |
| BlueMeme / Micro Court | Strong | Partial | Disagreement | 7 | 2 / 3 / 1 / 1 |
| Prop Tech plus | Strong | Partial | Disagreement | 15 | 3 / 7 / 0 / 5 |
| PIA TECH LAB | Partial | Weak | Disagreement | 10 | 3 / 3 / 1 / 3 |
| CORE | Strong | Strong | Agreement | 8 | 7 / 1 / 0 / 0 |
| teamLab | Partial | Weak | Disagreement | 8 | 2 / 4 / 1 / 1 |
| Nitori Digital Base | Weak | Weak | Agreement | 11 | 3 / 1 / 4 / 3 |

## Complete Requirement / Trace Observations

각 표는 모든 반환 requirement를 포함한다. S/L/N은 Supporting/Limiting/Non-decisive, evidence의 d/t/l은 direct/transferable/limited다. 원문 evidence, 정형 reasoning, Missing Evidence, Unknowns는 연결된 JSON에 모두 보존되어 있다. 아래 gap은 실제 reasoning의 주된 선택 코드이며 자유 생성 설명이 아니다. `insufficient`의 무연결 항목은 application Unknown이다. ID는 sample 내부 식별자이며 서로 다른 JD의 같은 L번호가 동일 역량을 뜻하지 않는다.

### BlueMeme / Micro Court

[입력](../../samples/bluememe.txt) · [전체 JSON](turn-008-results/bluememe.json)

- supporting: L005, L008
- limiting: L002
- nonDecisive: L006, L007, L009, L010

| ID / Requirement | Type | Match | Trace | Evidence | Missing / rationale gap |
| --- | --- | --- | --- | --- | --- |
| L002 Position: Low-code engineer | context | Weak | L | E1t, E2t, E3t, E4t, E5t, E6t | core role experience |
| L005 Required: Practical experience developing systems. | required | Strong | S | E1d, E2d, E3d, E4l, E5d | 없음 |
| L006 Responsibility: Interview customers about their needs and design systems. | responsibility | Partial | N | E3l, E4l | full scope |
| L007 Preferred: Web development experience and upstream customer coordination experience. | preferred | Partial | N | E3l, E4l, E5l, E6l | full scope |
| L008 Responsibility: Support project execution. | responsibility | Strong | S | E1l, E2l, E3l, E4d | 없음 |
| L009 Responsibility: Develop using OutSystems; prior OutSystems experience is not required and training is provided. | responsibility | Partial | N | E1t, E2t, E3t, E5t | exact technology experience |
| L010 Responsibility: Handle testing and releases. | responsibility | Unknown | N | 없음 | sufficient relevant evidence |

### Prop Tech plus

[입력](../../samples/proptech-plus.txt) · [전체 JSON](turn-008-results/proptech-plus.json)

- supporting: L002, L007, L011
- limiting: L005, L006, L008
- nonDecisive: L009, L010, L012, L013, L014, L015, L016, L017, L018

| ID / Requirement | Type | Match | Trace | Evidence | Missing / rationale gap |
| --- | --- | --- | --- | --- | --- |
| L002 Position: Server-side engineer | context | Strong | S | E1d, E2t, E3d, E4l, E5l | 없음 |
| L005 Business context: In-house package development for real-estate fund management companies. | context | Partial | L | E2t, E6l | specific domain experience |
| L006 Required: At least one year of development experience with Java 8 or later. | required | Partial | L | E1l | specified version or duration |
| L007 Responsibility: Design and develop backend APIs. | responsibility | Strong | S | E1d, E2l, E3d, E5l | 없음 |
| L008 Responsibility: Operate backend APIs. | responsibility | Partial | L | E1l, E3l | operational context |
| L009 Responsibility: Enhance the package and fix defects. | responsibility | Partial | N | E1l, E2l, E3l, E4l | full scope |
| L010 Responsibility: Set up and customize the package for new customers. | responsibility | Partial | N | E3l, E4l, E6l | full scope |
| L011 Technology environment: Java and Spring Boot. | stack | Strong | S | E1d, E3l | 없음 |
| L012 Technology environment: Azure App Service on Linux/Docker and Azure SQL Database. | stack | Unknown | N | 없음 | sufficient relevant evidence |
| L013 Technology environment: Git and Azure DevOps for CI/CD. | stack | Unknown | N | 없음 | sufficient relevant evidence |
| L014 Preferred: Hands-on Azure operation/administration experience. | preferred | Unknown | N | 없음 | sufficient relevant evidence |
| L015 Preferred: Accounting knowledge comparable to bookkeeping level 3. | preferred | Partial | N | E6l | requested depth |
| L016 Preferred: Program with the customer's usability in mind. | preferred | Unknown | N | E3l, E4l | sufficient relevant evidence |
| L017 Preferred: Willingness to understand the customer's business and work proactively. | preferred | Unknown | N | E3l, E4l, E6l | sufficient relevant evidence |
| L018 Optional adjacent web development: Angular, Angular Material and TypeScript. | preferred | Partial | N | E3l, E5l | requested depth |

### PIA TECH LAB

[입력](../../samples/pia-tech-lab.txt) · [전체 JSON](turn-008-results/pia-tech-lab.json)

- supporting: L002, L006
- limiting: L005, L007, L010
- nonDecisive: L008, L009, L011, L012, L013

| ID / Requirement | Type | Match | Trace | Evidence | Missing / rationale gap |
| --- | --- | --- | --- | --- | --- |
| L002 Position: Server-side engineer | context | Strong | S | E1d, E2l, E3l, E5l, E7l | 없음 |
| L005 Required: At least three years of AWS infrastructure construction and operations experience. | required | Weak | L | E7l | evidence at the required level |
| L006 Required: Web service design and development experience. | required | Strong | S | E1d, E3d, E5l | 없음 |
| L007 Required: Knowledge of infrastructure, middleware, communications and browsers. | required | Partial | L | E1l, E3l, E5l, E7l | full scope |
| L008 Responsibility: Establish monitoring and backup operations. | responsibility | Unknown | N | 없음 | sufficient relevant evidence |
| L009 Responsibility: Select and apply security patches. | responsibility | Unknown | N | 없음 | sufficient relevant evidence |
| L010 Responsibility: Handle performance and incidents from the network through middleware layers. | responsibility | Partial | L | E1l, E2l, E3l, E5l | operational context |
| L011 Preferred: Application development experience with Java or similar languages. | preferred | Strong | N | E1d, E3l, E5l | 없음 |
| L012 Preferred: AWS certification. | preferred | Unknown | N | E7l | sufficient relevant evidence |
| L013 Preferred: Project management experience for infrastructure construction. | preferred | Partial | N | E4l | full scope |

### CORE

[입력](../../samples/core.txt) · [전체 JSON](turn-008-results/core.json)

- supporting: L002, L005, L006, L008, L009, L010
- limiting: L011
- nonDecisive: L007

| ID / Requirement | Type | Match | Trace | Evidence | Missing / rationale gap |
| --- | --- | --- | --- | --- | --- |
| L002 Position: Systems engineer for web and business systems in Kyushu (Fukuoka), leader recruitment track | context | Strong | S | E1t, E2t, E3d, E4d, E5l, E6d | 없음 |
| L005 Required: Development experience from basic design onward. | required | Strong | S | E1l, E2l, E3d, E4l | 없음 |
| L006 Required: Communicate with customers to identify their problems. | required | Strong | S | E4d | 없음 |
| L007 Preferred: Upstream process experience and coordination of customer issues. | preferred | Strong | N | E3l, E4d | 없음 |
| L008 Responsibility: Act as a sub-leader, with project leader work possible if desired. | responsibility | Strong | S | E4d | 없음 |
| L009 Technology environment: Java and Oracle are examples of the development environment. | stack | Strong | S | E1d, E2l | 없음 |
| L010 Business context: Public-sector and energy business systems. | context | Strong | S | E2t, E6d | 없음 |
| L011 Responsibility: Work across requirements, system testing, introduction and maintenance. | responsibility | Partial | L | E1l, E2l, E3l, E4l | full scope |

### teamLab

[입력](../../samples/teamlab.txt) · [전체 JSON](turn-008-results/teamlab.json)

- supporting: L010
- limiting: L002, L006
- nonDecisive: L005, L007, L008, L009, L011

| ID / Requirement | Type | Match | Trace | Evidence | Missing / rationale gap |
| --- | --- | --- | --- | --- | --- |
| L002 Position: Web application and smartphone application engineer | context | Weak | L | E1t, E2t, E3l, E4t, E5l, E6t, E7l | core role experience |
| L005 Responsibility: Define requirements and design based on customer needs. | responsibility | Strong | N | E3d, E4d | 없음 |
| L006 Responsibility: Establish overall development direction and architecture. | responsibility | Partial | L | E1l, E2l, E3l, E4l | full scope |
| L007 Responsibility: Improve processing efficiency and generalize reusable functionality. | responsibility | Partial | N | E1l, E2l, E3l, E4l | requested depth |
| L008 Responsibility: Perform performance tuning. | responsibility | Partial | N | E1l, E2l | requested depth |
| L009 Responsibility: Build CI/CD pipelines. | responsibility | Unknown | N | 없음 | sufficient relevant evidence |
| L010 Technology environment: Java, Spring and SQL for backend development. | stack | Strong | S | E1d, E2l, E3l | 없음 |
| L011 Technology environment: Frontend and mobile technologies include React, Vue and Swift; the technology listing does not establish that all are mandatory together. | stack | Partial | N | E3l | exact technology experience |

### Nitori Digital Base

[입력](../../samples/nitori.txt) · [전체 JSON](turn-008-results/nitori.json)

- supporting: L005, L007
- limiting: L002, L006, L008, L009
- nonDecisive: L010, L011, L012, L013, L014

| ID / Requirement | Type | Match | Trace | Evidence | Missing / rationale gap |
| --- | --- | --- | --- | --- | --- |
| L002 Position: EC frontend engineer (NDB0055) | context | Weak | L | E3l, E4t, E5l | core role experience |
| L005 Required: Project team leadership experience. | required | Strong | S | E4d | 없음 |
| L006 Required: Web frontend Tech Lead experience. | required | Weak | L | E3l, E4t, E5l | core role experience |
| L007 Required: Experience developing and mentoring team members. | required | Strong | S | E4d | 없음 |
| L008 Responsibility: Design, develop and operate frontend products. | responsibility | Weak | L | E3l, E5l | core role experience |
| L009 Responsibility: Maintain frontend quality and conduct design and code reviews. | responsibility | Weak | L | E3l, E4t, E5l | core role experience |
| L010 Technology environment: Angular and TypeScript for frontend development. | stack | Partial | N | E3l, E5l | specified version or duration |
| L011 Technology environment: Java and Spring Boot in the organization's backend stack. | stack | Strong | N | E1d | 없음 |
| L012 Preferred: Large-scale e-commerce construction and operations experience. | preferred | Unknown | N | 없음 | sufficient relevant evidence |
| L013 Preferred: Knowledge of SEO, advertising, search, location and payment services. | preferred | Unknown | N | 없음 | sufficient relevant evidence |
| L014 Preferred: Experience using design tools. | preferred | Unknown | N | 없음 | sufficient relevant evidence |

## Quality Evaluation Matrix

OK는 이번 입력/출력에서 문제를 관찰하지 못했다는 뜻이며 완전한 정확성 인증이 아니다. Concern과 Violation candidate는 후속 검토 대상이며 이번에 수정하지 않았다. Agreement인 sample에도 문제 후보가 있을 수 있다.

| Sample | Faithfulness | Requirement Fidelity | Unknown Preservation | Traceability | Usefulness | Notable Issue |
| --- | --- | --- | --- | --- | --- | --- |
| BlueMeme | OK: 경력 원문 생성 없음 | Concern: 사전 도구 경험 불문인데 role gap 강조 | OK: 테스트/릴리스 Unknown | OK: ID/원문 연결; scope 설명은 포괄적 | Concern: 직무명 Weak가 업무 강점을 가릴 수 있음 | L002 role Weak/Limiting, L009 도구 Partial/Non-decisive |
| Prop Tech plus | OK: 기존 E1–E7만 사용 | OK: category 보존; domain 영향은 검토 필요 | OK: Azure 등 미확인 보존 | OK: 전체 artifact로 추적 | OK: 버전/기간·운영 gap을 구분 | L005 domain Limiting은 한 sample 관찰 |
| PIA | OK: basic AWS를 전문 운영 경력으로 생성하지 않음 | OK: AWS 3년 필수와 Java 우대 구분 | OK: 자격/백업/패치 Unknown | OK: 각 gap/근거 연결 | OK: AWS 중심 제한을 볼 수 있음 | Reference Partial과 달리 Weak지만 계약 위반 근거는 없음 |
| CORE | Violation candidate: 복합 context 전체 Strong | Concern: 공공+에너지 범위의 압축 | Concern: L010 에너지 미확인이 출력되지 않음 | OK: E6 복사; 전체 충족 근거는 불충분 | Concern: 유용한 강점 목록이나 trace 모순 | Overall Strong인데 L011 Limiting |
| teamLab | OK: mobile 경력을 생성하지 않음 | Concern: 복합 직무명의 role gap 해석 | OK: CI/CD Unknown; 일본어는 제외 | OK: ID/근거 연결; depth 정형 문구 한계 | Concern: 넓은 직무명 Weak가 핵심 제한 | Web/mobile 배정 경로 ambiguity, 일본어 미평가 |
| Nitori | OK: frontend Tech Lead 경력 합성 없음 | Violation candidate: L010에 없는 version/duration을 gap으로 설명 | OK: EC/서비스/디자인 도구 Unknown | OK: 원문 연결; L010 gap 부정확 | OK: 핵심 역할 gap과 일반 리딩 강점 구별 | 같은 역할 gap이 L002/L006/L008/L009에 반복 |

### Evidence Faithfulness

반환된 모든 Supporting Evidence ID/text를 candidate profile과 비교했으며 원문이 일치했다. 새 경력 text 생성은 관찰되지 않았다. 그러나 원문 복사가 곧 의미적 충실성은 아니다. CORE L010은 E6의 public-sector 경험에 direct를 부여하고 public-sector **and energy** 전체를 Strong으로 설명한다. Energy는 profile에 없다. 직접 에너지 경력 문장을 만들어낸 것은 아니지만 전체 조건 충족을 암시하는 과도한 요약 후보다.

Nitori는 일반 리딩 E4를 frontend Tech Lead의 transferable로 연결하되 L006을 role_gap Weak로 남겼다. E3/E5를 합쳐 frontend 리딩 경험이 있다고 생성하지 않았다. Backend Java/Boot L011 Strong도 Non-decisive라서 그것이 frontend 역할을 입증한다고 설명하지 않는다. PIA 역시 E7 basic AWS를 limited로 연결하고 L005 Weak를 유지했다.

### Requirement Fidelity

입력에 명시된 Required/Preferred/Responsibility/Stack은 반환된 행에서 동일 category로 보존되었다. 직무명은 context다. Metadata는 ignore되었고 teamLab L012의 일본어만 frozen language exclusion에 따라 빠졌다. 이는 실제 언어 조건을 충족했다는 뜻이 아니다. CORE 학력은 input에 구체 조건을 재구성하지 않았으므로 평가 대상이 아니다. 결과를 전체 지원 자격 평가로 읽지 않는다.

Category 보존과 영향의 타당성은 별개다. BlueMeme에서는 OutSystems 사전 경험 불문이 입력에 있지만 L002 Low-code engineer가 Weak/Limiting이다. 해당 업무의 준비도를 좁게 본 해석인지, 요구하지 않은 사전 역할 경력을 사실상 전제했는지 명확하지 않다. teamLab도 Web/mobile 전체 역할을 동시에 입증해야 하는지 불명확한 상태에서 직무명 Weak/Limiting이다. 두 경우 모두 source category를 Required로 실제 바꾼 사례는 아니다.

Nitori L010의 기술 환경에는 버전/기간 요건이 없는데 `specified version or duration` gap이 선택되었다. 요구하지 않은 세부조건이 명시된 것처럼 설명되는 requirement fidelity/설명 계약 위반 후보다. Partial label이 reference와 같다는 사실은 이 문제를 없애지 않는다.

### Unknown Preservation

근거 없는 자격/특화 서비스/도구, CI/CD, 운영 정보는 Unknown 또는 미확인 gap으로 남았다. No Evidence를 실제 무경험으로 단정한 출력은 관찰되지 않았다. Unknown이 아닌 Partial/Weak도 미확인 항목을 포함하므로 모든 누락을 무조건 Unknown label로 바꿔야 한다고 보지는 않는다.

CORE L010 Strong의 missing/unknown 배열은 비어 있어 에너지 도메인 미확인이 가려지는 후보가 있다. BlueMeme OutSystems는 Partial이나 기술 미확인이 명시되어 직접 사용 경험 주장으로 읽어서는 안 된다. PIA AWS 자격은 E7 limited 연결이 있어도 Unknown을 유지한다. 관련 근거 연결 자체와 판단 충분성을 구분한 사례다.

### Evidence Traceability

모든 결과에서 requirement ID/text, relation을 포함한 profile 원문, missing/unknown, trace 역할을 연결할 수 있다. 모든 ID 참조와 summary 분포를 검증했다. 다만 `full scope`, `requested depth`, `core role experience` 등 closed-set 문구는 어떤 부분이 충족되었고 어떤 부분이 부족한지까지 설명하지 못한다. BlueMeme L006의 고객 분석/설계, teamLab L007의 처리 개선은 각각 기존 E3/E4, E2/E3 근거가 연결되어도 limited+포괄적 gap으로 압축된다. CORE 에너지처럼 Strong이면 아예 세부 미확인이 빠질 수 있다.

### Structural Consistency

모든 6개 artifact는 동일 schema와 네 Match level, 세 trace 역할을 사용한다. 배열과 원문/ID 대응은 유효하고 Unknown Supporting은 없다. 다만 구조 검증을 통과해도 semantic 정합성은 별도다. CORE의 최고 level Strong에 Limiting이 존재하는 문제와 Nitori의 잘못 선택된 duration rationale는 현재 guard가 잡지 못했다.

Compound requirement에서 일부 직접 근거 하나만으로 전체 Strong이 허용되는 구조가 CORE L009(Java/Oracle), L010(공공/에너지), teamLab L010(Java/Spring/SQL)에 보인다. Java/Oracle와 Java/Spring/SQL은 profile에 각각 기술이 존재하므로 Strong 자체를 오류로 보지 않는다. 공공/에너지는 에너지 근거가 없어 다르다. 동일 출력 구조가 같은 정도의 evidence coverage를 보장하지 않는다.

### Decision Trace Usefulness

PIA는 AWS 및 인프라 scope/운영이 Limiting이고 Java 우대 Strong은 Non-decisive다. Nitori는 일반 리딩/육성을 Supporting으로 두면서 frontend 역할 gap을 Limiting으로 구분한다. Prop Tech plus는 Azure 우대와 Java 버전/기간을 구분한다. 이 정보는 Overall label만 있을 때보다 유용하다.

반면 CORE의 Overall Strong + L011 Limiting은 “더 높은 Match를 제약한다”는 현재 역할 정의와 충돌하는 후보다. 최고 level에서 어떤 방향을 제한한다는 뜻인지 현재 출력으로 해석할 수 없다. 이를 confidence 저하 같은 새 의미로 임의 재정의하지 않는다. teamLab의 요구 분석 Strong은 Non-decisive이고 backend stack만 Supporting인 것도 정보 손실 가능성을 보여주지만, 모든 Strong이 Supporting이어야 하는 것은 아니므로 그 자체를 오류로 판정하지 않는다.

모든 trace는 **structured post-hoc attribution**이다. 별도 요청의 분류이지 원래 Overall의 causal explanation은 아니다. 따라서 “Preferred가 실제로 영향을 전혀 주지 않았다”, “직무명이 실제로 결과를 낮췄다”라고 인과관계를 확정할 수 없다.

### Repeatability / Variability

새 sample은 각각 1회이므로 sample별 repeatability를 측정하지 않았다. 서로 다른 JD 간 결과 차이는 동일 입력 variability의 증거가 아니다. Turn 004–006의 회계 Weak → Partial, Azure 환경 Partial → Unknown 관찰을 유지한다. 이번에는 Prop Tech plus를 재실행하지 않아 추가적인 동일 입력 variability 관측도 없다.

### User Usefulness

PIA의 AWS gap, Nitori의 frontend 특화 리딩 gap, CORE의 기본설계/고객 조율 강점은 사람이 후속 확인할 위치를 찾는 데 도움이 된다. Weak도 지원 여부 결정을 뜻하지 않는다. BlueMeme/teamLab의 넓은 직무명 gap, CORE의 energy coverage, Nitori의 잘못된 duration 설명은 사람에게 오해를 줄 수 있어 label보다 evidence와 source를 함께 읽어야 한다. 여기서 usefulness는 문서 검토자의 정성적 판단이며 실제 사용자 연구 결과는 아니다.

## Contract Violation Candidates

| 후보 | 관찰 근거 | 판정 범위 |
| --- | --- | --- |
| CV-1 CORE trace의 상한 모순 | Overall Strong이면서 L011 Partial을 Limiting으로 분류 | 현재 “더 높은 Match 제약” 정의와 충돌 후보. Gap 존재 자체가 오류라는 뜻은 아님. |
| CV-2 CORE 복합 context coverage | L010 public-sector and energy Strong, E6 direct, missing/unknown 없음 | 일부 직접 근거를 전체 충족으로 확대했을 가능성. Energy 경력 원문 생성은 없으므로 literal hallucination과 구분. |
| CV-3 Nitori 미명시 duration/version | L010 Angular/TypeScript stack에 기간/버전 조건이 없지만 `specified version or duration` rationale | Explanation/requirement fidelity 위반 후보. Partial 자체가 틀렸다는 주장은 아님. |

BlueMeme 사전 도구 경험 불문 대비 role_gap은 **별도 Concern**이다. Mandatory 승격이 입증되지 않아 확정 violation으로 올리지 않는다. teamLab 직무 범위와 경로 해석도 같은 방식으로 ambiguity를 남긴다. 위 세 후보는 즉시 수정하거나 결과를 재생성하지 않았다.

## Interpretation Differences

- **Prop Tech plus Strong / Partial:** core backend 강점과 Java 세부조건·제품 운영/domain 미확인의 중요도 차이. 기존 결과에서 label 차이 자체를 위반으로 보지 않는다.
- **PIA Partial / Weak:** AWS 3년 구축·운영은 명시 필수이고 candidate는 basic AWS다. Web 강점과 필수 운영 gap을 요약하는 방식이 다를 수 있다. 이번 결과에는 근거 있는 해석 차이로 기록하며 계약 위반 후보는 식별하지 못했다.
- **BlueMeme Strong / Partial:** 도구 경험 불문과 low-code role gap 해석이 충돌할 가능성이 있어 단순한 합리적 disagreement로 확정하지 않는다. 우선 Concern으로 남기되 label 차이 때문에 실패 처리하지 않는다.
- **teamLab Partial / Weak:** architecture 및 Web/mobile 역할 범위를 중요하게 본 해석일 수 있다. 대체 배정 경로와 복합 직무 해석이 미해결이므로 우선 Concern이다. Reference의 일본어 Unknown은 runtime에서 제외되어 비교 범위도 완전히 같지 않다.

CORE와 Nitori는 Overall label이 reference와 같지만 위반 후보가 있다. 이는 label agreement를 quality 인증으로 삼지 않는 이유다. Accuracy/Precision/Recall/label 점수는 계산하지 않았다.

## Systematic Pattern Candidates

여러 sample에서 실제 관찰된 behavior만 다음에 적는다. 반복된 behavior가 곧 수정할 systematic problem이라는 뜻은 아니다.

1. **직무명 context의 role_gap/Limiting:** BlueMeme L002, teamLab L002, Nitori L002. Nitori에는 별도 frontend Tech Lead Required가 있지만 BlueMeme는 도구 경험 불문이고 teamLab은 경로 ambiguity가 있다. Broad title을 specialized prior experience처럼 좁혀 읽는지 후속 검토할 반복 후보다. Nitori의 타당한 역할 제약까지 일괄 제거할 근거는 아니다.
2. **직접/전이 가능한 세부 능력을 limited + generic gap으로 압축:** BlueMeme 고객 분석·설계 L006 및 우대 L007, teamLab 효율화 L007/성능 L008, Prop Tech 제품 개선 L009 등에서 관찰된다. E3 legacy/설계와 E2 개선 정보가 있어도 정확히 어떤 세부가 부족한지 식별하기 어렵다. 반복적인 설명 해상도 한계이며 모든 Partial이 오판이라는 뜻은 아니다. BlueMeme L007에서 E1 전문 REST/Java가 미연결인 것은 개별 누락 Concern으로, 같은 누락의 반복이 입증되지는 않았다.
3. **Preferred Unknown은 일관되게 Non-decisive:** Prop Tech Azure/태도 항목, PIA AWS 자격, Nitori EC/도구 등. 이번 관측은 “모든 Preferred Unknown을 과도하게 제한한다”는 가설을 지지하지 않는다. Post-hoc 분류이므로 최초 판단의 무영향 증명은 아니다.

다음은 반복 문제로 단정하지 않는다: business/domain context 자체의 Limiting은 Prop Tech에서만 확인되고 CORE domain은 Supporting이다(직무명 context와 구분). Nitori의 role gap이 네 행에서 반복되지만 동일 gap의 다중 제한이 여러 JD에서 반복된다는 증거는 부족하다. Version/duration 제한의 과도한 반복도 미입증이다. Prop Tech는 실제 기간 조건이고 Nitori는 잘못된 rationale이나 Non-decisive여서 같지 않다. 줄 단위 추출은 모든 JD에 공통이지만 이것만으로 over-decomposition 또는 implicit weighting을 증명할 수 없다.

## Problems / Findings

새 실행 모두 성공하여 evaluation을 막는 execution bug는 없었다. 발견한 계약 후보와 Concern은 기록만 하고 동결 기준을 유지했다. 두 가지 입력/실험 한계가 중요하다: 역사적 source summary 기반이라는 점, 각 신규 sample 1회와 기존 Prop Tech 1회만 비교한다는 점이다. 전체 공개 원문에 대한 일반화, 모델 간 우위 또는 통계적 repeatability를 주장하지 않는다.

## Tests / Validation

- `go test ./...`: 통과. 외부 API 없는 기존 regression tests.
- `go vet ./...`: 통과.
- `go build ./cmd/job-fit`: 통과. Git ignored local binary로 실행.
- 실행 전후 frozen fingerprint 및 기존 파일 SHA-256 일치 확인. Input hash도 평가 전 고정값과 일치.
- 6개 JSON의 requirement 고유 ID·JD line 원문·source category·summary 분포·evidence 원문·trace 전체 일대일 참조 및 복사된 Match/source type 검증.
- 새 sample마다 예상 metadata 제외, teamLab 일본어 제외 확인. Missing/Unknown을 포함한 원래 출력 그대로 보존.
- Prop Tech 결과는 Turn 006 JSON과 byte-for-byte 동일. Manifest hash, Markdown 링크/표/공백 및 API 키 미포함 확인.

## Changed Files

- `samples/bluememe.txt`, `pia-tech-lab.txt`, `core.txt`, `teamlab.txt`, `nitori.txt` 신규.
- `docs/development/turn-008-results/`의 6개 JSON 및 `manifest.json` 신규.
- 본 `docs/development/turn-008.md` 신규.

README/Contract/Match Semantics/이전 Turn/code/profile/정책은 이번 Turn에서 수정하지 않았다. 이전 Turn의 미커밋 변경도 보존했고 commit은 생성하지 않았다.

## Deferred

계약 후보 수정, calibration/tuning, profile/prompt/logic/validation 변경, grouping/de-duplication, 수치 weighting/threshold/percentage, 추가 반복·Prop Tech 재실행, MCP/DB/UI/crawler/parser/persistence architecture는 수행하지 않았다.

## Suggested Next Step

Planning conversation에서 CV-1~3의 실제 계약 위반 여부와 입력 요약의 한계를 먼저 검토한 뒤, evaluator 유지 / 확인된 계약 위반만 수정 / 반복 pattern에 대해 한 번의 calibration pass 중 다음 범위를 결정한다. 특히 직무명 role_gap 패턴을 검토하되 reference label에 맞추는 목적과 분리한다. 다음 Turn의 수정이나 재실행을 시작하지 않았다.

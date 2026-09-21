# Match Semantics & Human Baseline

- 작성 및 Source Reviewed Date: 2026-09-21 (Asia/Seoul)
- 범위: Capability / Experience Match의 의미와 공개 JD 6개에 대한 예상 판단 기록
- 상태: **Turn 003 지시에 따라 평가 규칙과 6개 Overall baseline 확정** (2026-09-21). Turn 002의 Codex 작성 초안에서 발전한 기록이다.

## Turn 007 — 현재 문서의 해석

이 문서의 기존 `Human Baseline`, `Expected Match`, `Overall Expected Match`는 모두 **Pre-evaluation Reference Interpretation**을 뜻한다. Evaluator 결과를 보기 전에 중요하게 본 요소와 하나의 해석을 기록한 것이며, evaluator의 target label 또는 객관적 ground truth가 아니다. Turn 003의 “확정”은 비교 기록의 확정이지 정답 인증이 아니다. Codex 초안과 사용자 지시에 따른 provenance도 유지한다.

아래 sample별 값·근거·당시 표현은 보존한다. “Strong 유지”, “Weak 확정” 등의 sample 해석을 runtime에 강제할 정답으로 사용하지 않는다. Human Reference와 Jev의 label 차이만으로 오류를 선언하지 않으며, 둘 다 evidence와 semantics 안에서 합리적이면 disagreement로 남길 수 있다. Overall은 summary signal이고 requirement-level evidence와 Unknown이 더 중요한 정보다.

Strong / Partial / Weak / Unknown 및 Required/Preferred, core role, transferable experience의 기존 semantics는 변경하지 않는다. Turn 006 implementation과 함께 다음 전체 sample evaluation 동안 동결한다. 품질 기준과 변경 조건은 [Evaluation Contract](evaluation-contract.md#evaluation-philosophy--turn-007), 전환 경위는 [Turn 007](development/turn-007.md)에 기록한다.

## 1. 목적과 기록의 성격

평가 engine의 결과를 본 뒤 기준을 맞추는 것을 피하기 위해, engine 실행 전에 의미와 예상 판단을 기록한다. Jev 호출, 별도 AI evaluator 실행, 점수 계산은 수행하지 않았다. 이 문서는 정답 label이나 채용 가능성의 예측이 아니다.

사용자가 제공한 정의와 가정된 경력 evidence를 바탕으로 Codex가 공개 JD를 읽어 초안을 작성했다. Turn 003에서 사용자가 제시한 평가 규칙과 baseline 확정 지시에 따라 기존 6개 Overall을 유지하여 비교 기준으로 확정했다. Requirement-level 근거와 Unknown은 그대로 보존했다. 이는 모든 세부 사실을 사용자가 독립적으로 검증했다는 의미나 AI의 도움 없이 사람이 작성한 ground truth라는 의미가 아니다. 이후 근거가 달라지면 변경 이유를 새 Turn에 기록한다.

## 2. Match Level

| Level | 의미 | 경계 확인 |
| --- | --- | --- |
| Strong | JD의 핵심 requirement를 직접 뒷받침하는 명확한 evidence가 있고 role/context 차이가 크지 않다. | 장기간 backend 경력이 모든 backend 세부 기술을 입증하지는 않는다. |
| Partial | 관련 evidence와 transferable experience가 충분하지만 일부 scope, depth, technology, domain 또는 operational context가 부족하다. | Transferable experience를 설명하되 직접 경험으로 바꾸지 않는다. |
| Weak | 관련 evidence가 일부 있으나 핵심 요구와 상당한 gap이 있다. | 단순 keyword overlap으로 평가를 높이지 않는다. Gap은 우선 기록된 evidence의 한계다. |
| Unknown | 충분한 관련 정보가 없어 판단할 수 없다. | No Evidence ≠ No Experience. 누락을 경험 부재로 단정하지 않는다. |

사용자 제공 예시: Java/Spring 전문 backend 경력은 동일 기술의 개발 요구에 Strong 근거가 된다. Backend/API/RDB 설계와 Tech Lead 경력은 직접적인 대규모 distributed SaaS 설계 evidence가 없다면 Partial을 검토한다. 기본 AWS 경험만으로 production AWS architecture를 뒷받침하면 Weak이며, Kubernetes 운영에 관한 정보가 전혀 없으면 Unknown이다. 이 예시들은 아래 JD에 해당 요구가 있다는 뜻이 아니다.

**Missing Evidence**는 match level과 별개이며, 현재 Candidate Profile만으로 requirement를 충분히 입증하지 못하는 정보를 뜻한다. Partial과 Unknown 모두 Missing Evidence를 가질 수 있다. Strong이라도 더 좁은 세부 사실은 Unknown으로 남을 수 있다.

아래 표의 `Missing / Weak Evidence`는 필요한 근거 또는 현재 근거의 한계를, `Unknowns`는 확인되지 않은 사실을 적는다. 동일한 정보가 두 관점에 걸쳐 나타날 수 있으며, 배타적인 필드나 JSON schema로 정의한 것은 아니다. `추가 없음`은 해당 requirement의 표현 수준에서 추가 근거 부족을 식별하지 못했다는 뜻이며, 완전한 검증을 뜻하지 않는다.

## 3. 적용 시 구분할 사항

- **Required**: JD가 명시한 필수 경험·역량으로, 해당 역할 수행에 핵심적으로 기대된다. Mismatch는 우대 사항 부족보다 중요하게 취급하되 만족 여부를 추측하지 않는다.
- **Preferred / Nice-to-have**: 있으면 유리하지만 부재 자체가 핵심 mismatch는 아니다. 부족만으로 Overall을 크게 낮추지 않으며 필수 조건으로 바꾸지 않는다.
- **Responsibility**: 맡을 업무. 입사 전 필수 경험과 동일하지 않다.
- **Stack / Context**: 개발환경과 사업 맥락. 기술이 나열됐다는 이유만으로 전부 필수 경험으로 만들지 않는다.

이 구분은 출처의 표현을 보존하기 위한 것이며 weighting formula가 아니다. Overall Strong도 모든 지원 자격 충족을 뜻하지 않는다.

### Core role alignment

Frontend Tech Lead, Backend Architect, SRE, Data Engineer, Engineering Manager처럼 역할 수준의 요구는 기술 keyword보다 중요할 수 있다. **Core role mismatch는 개별 skill overlap보다 우선한다.** 관련 기술이 일부 겹쳐도 핵심 역할 경험의 큰 gap을 다른 skill의 수로 상쇄하여 높은 match를 부여하지 않는다.

Angular/TypeScript 사용과 일반 Tech Lead 경험을 결합하여 frontend Tech Lead 경력을 만들어내지 않는다. Nitori에서는 그 역할을 입증하는 evidence가 제한적이므로 Overall Weak를 유지한다. 이는 기록된 evidence와 역할의 gap에 대한 판단이다. 실제 해당 경험이 없다고 확인된 경우와 profile에 직접 기록이 없는 경우는 다르며, 이번 reference evidence에 경험 부재라는 새 사실을 추가하지 않는다.

### Weak와 Unknown의 경계

Weak는 비교할 관련 근거가 존재하지만 핵심 요구 또는 요구 수준과의 gap이 큰 경우다. 기본 AWS 사용은 production AWS architecture의 Weak 근거가 될 수 있다. Unknown은 해당 requirement를 평가할 관련 근거가 충분하지 않은 경우다. Kubernetes 운영 정보가 전혀 없다면 Unknown이며 경험 부재를 뜻하지 않는다.

단순 keyword만으로 비교 근거가 충분하다고 판단하지 않는다. Nitori의 Weak는 기술명뿐 아니라 frontend 개발과 일반 리딩이라는 인접 경험의 범위가 확인된 상황에 적용한다. 이를 모든 직접 evidence 누락에 Weak를 부여하는 규칙으로 일반화하지 않는다. 확인되지 않은 기간·책임·숙련도는 level과 별도로 Unknowns에 남긴다.

### Transferable experience

Exact technology/domain이 다르다는 이유만으로 Weak 처리하지 않는다. 전이 가능성을 판단할 때 다음을 함께 설명한다.

- Underlying engineering problem의 유사성
- Responsibility level
- Architecture/design depth
- Operational context
- Leadership / ownership
- 직접적인 기술 전환 가능성

충분한 관련 경험이 있으나 특정 context나 깊이가 부족하면 Partial이 적절할 수 있다. 일반 backend scalability 경험을 직접적인 hyperscale distributed SaaS 경험으로 바꾸지는 않는다. 동일 도구 경험이 필요 없는 JD의 일반 설계·개발 요구에는 직접 evidence를 기준으로 Strong을 부여할 수 있다.

### Overall Match 해석

Overall은 requirement-level 결과의 단순 평균이 아니다. Core role alignment, Required requirements, Major capability areas, Transferable experience, Critical gaps, Preferred requirements를 함께 보고 서술형 이유로 요약한다.

Core Required mismatch는 Overall의 상한에 영향을 줄 수 있다. 핵심 frontend Tech Lead 역할의 근거가 부족한 상황에서 기술 keyword overlap만으로 Overall Strong을 주는 것은 부적절하다. Preferred 부족은 Overall을 크게 낮추는 근거로 사용하지 않는다.

이 원칙은 모든 Required Weak에 동일한 상한을 적용하는 자동 규칙이 아니다. PIA의 AWS 운영 gap은 Strong을 지지하지 못하게 하지만, Web 설계·개발 및 관련 engineering 경험을 함께 고려하여 이 sample은 Partial로 확정한다. Nitori는 포지션을 규정하는 frontend 리딩의 gap 때문에 Weak로 확정한다. 각각의 판단 근거를 보존하며 점수, 평균, weighting formula 또는 threshold는 정의하지 않는다.

JD 원문을 확인할 수 없으면 requirement를 추측하지 않고 Overall Unknown과 source 확인 실패를 기록해야 한다. 이는 candidate evidence가 부족한 Unknown과 원인을 구분한다. 이번에는 접근 방법을 바꾸어 6개 모두 본문을 확인했다.

## 4. Reference Candidate Evidence

이번 Turn에서 사용자가 존재한다고 가정하도록 제공한 정보다. 독립적으로 경력의 사실 여부를 검증한 것은 아니다. 아래 번호는 문서 안의 참조용이며 Candidate Profile 필드나 schema가 아니다.

| 참조 | 제공된 evidence |
| --- | --- |
| E1 | 약 16년 software/backend 개발. Java / Spring / Spring Boot 중심 professional backend experience. REST API 설계·구현. |
| E2 | RDB 설계, SQL/DB 처리 개선. DB2 / Oracle / MySQL / MariaDB. 복잡한 business logic을 포함하는 업무 시스템. |
| E3 | 요구사항 분석, 기본·상세설계, API·DB 설계, legacy 분석·개선, 공통 모듈 설계, 코드리뷰. Backend / frontend / database를 아우르는 개발. |
| E4 | Tech Lead / Development Lead, PM / PL, 팀 구축, OJT, 개발 표준 정비, 코드리뷰·기술지원. 최대 약 12명 규모 조직/팀 경험. 고객 요구사항·specification 조율. |
| E5 | TypeScript, Node.js, NestJS, Angular, Redis, Netty, C++, TCP/IP 경험. 개별 기간·깊이는 제공되지 않음. |
| E6 | ERP, 회계, 개인정보 DB 암호화, 선박 모니터링, 물류/관제, 공공 Web / 업무 시스템 도메인. |
| E7 | AWS 기본 경험. 대규모 cloud-native production architecture/operation 경험은 명확하게 입증되어 있지 않음. |

총 경력 기간을 Java 특정 버전, frontend, cloud 또는 개별 도메인의 경력 기간으로 변환하지 않는다. 기술 목록과 leadership 목록을 결합하여 “frontend Tech Lead 경력”을 새로 만들지 않는다. 이 정보 밖의 경험·자격·언어 능력은 추측하지 않는다.

## 5. Sample Evaluations

각 Source URL의 공개 본문에서 관련 요구만 요약했다. 전체 JD 원문은 repository에 복제하지 않는다. 채용공고는 변경·삭제될 수 있으며, 이 기록은 검토일에 확보한 내용에 한정된다. 사용자 prompt의 관찰 대상은 검색 단서로만 사용하고, 실제 requirement는 원문에서 확인했다. 급여·비자·지역·근무 선호와 지원 추천은 평가하지 않는다.

### Sample 01 — BlueMeme / Micro Court

- Company / Position: BlueMeme 그룹 Micro Court / low-code engineer
- Source URL: [공개 공고](https://hrmos.co/pages/bluememe/jobs/0000005)
- Reviewed Date: 2026-09-21
- Source 확인: 웹 읽기 도구 실패 후 동일 URL의 공개 HTML 본문 확인.
- Relevant Requirements Summary: 실무 시스템 개발을 필수로 요구하고, Web 개발과 고객 조율은 우대한다. OutSystems 기반 업무를 맡지만 해당 기술의 사전 경험은 불문이며 교육을 제공한다.
- 출처 위치: 필수/우대 조건, 구체 업무 및 OutSystems 교육 설명.

| Requirement / 구분 | Supporting Evidence | Missing / Weak Evidence | Unknowns | Expected Match | Reasoning |
| --- | --- | --- | --- | --- | --- |
| 시스템 개발 실무 / Required | E1, E2, E3 | 추가 없음 | 프로젝트별 기간 | Strong | 업무 시스템 개발과 설계 경력이 직접 대응한다. |
| 고객 hearing·시스템 설계 / Responsibility | E3, E4 | 추가 없음 | 개별 산출물 | Strong | 요구사항 분석과 고객 specification 조율 근거가 있다. |
| Web 개발·상류 고객 조율 / Preferred | E1, E3, E4, E6 | 추가 없음 | 제안영업 수행 여부 | Strong | Java/Web 개발과 고객 조율을 입증하며 영업 경험은 추론하지 않는다. |
| 프로젝트 추진 지원 / Responsibility | E4 | 추가 없음 | 해당 조직 업무 방식 | Strong | PM/PL 및 개발 리딩이 직접 관련된다. |
| OutSystems 개발 / Responsibility, 사전 경험 불문 | 직접 근거 없음 | OutSystems 사용 이력 | 사용 여부·숙련도 | Unknown | 일반 개발 경험을 특정 도구 경험으로 치환하지 않는다. |
| 테스트·릴리스 대응 / Responsibility | E3의 개발 전반 경험 | 테스트·릴리스 담당 이력 | 담당 범위 | Partial | 인접 개발 근거는 있으나 직접 수행 범위가 불명확하다. |

**Overall Expected Match: Strong (Turn 003 확정)**

Why: 필수 시스템 개발 및 설계·고객 조율 업무에 직접 근거가 있다. Primary Gaps: OutSystems 숙련과 테스트·릴리스 역할은 확인되지 않았다. 사전 경험을 요구하지 않는 도구의 Unknown을 필수 결격처럼 취급하지 않아 Strong을 유지한다. 이 요약은 즉시 생산성이나 학습 속도를 보장하지 않는다.

### Sample 02 — Prop Tech plus

- Company / Position: Prop Tech plus / server-side engineer
- Source URL: [공개 공고](https://proptech.plus/recruit/occupations/20220203-568/)
- Reviewed Date: 2026-09-21
- Source 확인: 공개 본문 확인.
- Relevant Requirements Summary: Java 8 이상 개발 1년 이상이 필수다. 자사 부동산펀드 패키지의 backend API를 맡는다. Azure 운영은 우대이고 Angular/TypeScript Web 개발은 추가 가능 역량으로 소개한다.
- 출처 위치: 업무, 개발환경, 필수/우대 조건.

| Requirement / 구분 | Supporting Evidence | Missing / Weak Evidence | Unknowns | Expected Match | Reasoning |
| --- | --- | --- | --- | --- | --- |
| Java 8 이상 개발 1년 이상 / Required | E1 전문 Java 경력 | 버전별 기간 | Java 8+ 사용 기간 | Partial | 언어 경험은 직접 관련되지만 버전·기간 조건을 확인할 수 없다. |
| Backend API 설계·개발 / Responsibility | E1, E2, E3 | 추가 없음 | 개별 서비스 규모 | Strong | API/RDB 설계와 구현이 직접 대응한다. |
| API 운영·패키지 개선·도입 customization / Responsibility | E2, E3, E6 | 제품 운영·도입 사례 | 해당 패키지 운영 책임 | Partial | 업무 시스템과 legacy 개선은 관련되나 제품 맥락은 미입증이다. |
| Spring Boot·RDB / Stack | E1, E2 | Azure SQL 직접 경험 | Azure SQL 숙련도 | Strong | 일반 Spring Boot/RDB 축의 판단이며 Azure SQL 경험을 주장하지 않는다. |
| Azure 조작·운영 / Preferred | E7 기본 AWS만 존재 | Azure 직접 근거 | Azure 사용·운영 여부 | Unknown | 타 cloud의 기초 경험만으로 Azure 운영을 판단하지 않는다. |
| 회계 지식 / Preferred | E6 회계·ERP | 부기 3급 상당 지식의 구체 근거 | 회계 지식의 범위 | Partial | 관련 도메인 경험은 있으나 지식 수준의 동등성은 미확인이다. |
| Angular·TypeScript Web 개발 / Optional adjacent | E3, E5 | 기간·프로젝트 상세 | 숙련 깊이 | Partial | 사용 경험은 있으나 역할 수행 깊이까지 입증하지 않는다. |

**Overall Expected Match: Strong (Turn 003 확정)**

Why: 중심 업무인 backend/API/업무 데이터 처리와 Spring Boot에 직접 근거가 집중되어 Strong을 유지한다. Primary Gaps: Java 버전별 기간, Azure 운영, 제품 운영 맥락. Java 버전·기간의 미확인을 확인된 mismatch로 바꾸지 않으며 해당 행의 Partial은 보존한다. Azure 우대의 Unknown만으로 전체를 크게 낮추지 않는다. Overall Strong이 Java의 명시된 기간 조건 충족을 확정하지는 않는다.

### Sample 03 — PIA TECH LAB

- Company / Position: PIA TECH LAB / server-side engineer
- Source URL: [공개 공고](https://tech-lab.pia.jp/recruit/job05.html)
- Reviewed Date: 2026-09-21
- Source 확인: 웹 읽기 도구에서 502 발생. 동일 URL의 공개 HTML 본문에서 업무·필수·우대 조건 확인.
- Relevant Requirements Summary: AWS 인프라 구축·운영 3년 이상, Web 설계·개발, 인프라/미들웨어/통신/브라우저 지식이 필수다. 운영 체계, 보안 패치, 성능·장애 대응 업무가 명시된다.
- 출처 위치: 업무 내용, 필수 조건, 우대 조건. 직무명만으로 일반 Java backend 역할이라고 해석하지 않았다.

| Requirement / 구분 | Supporting Evidence | Missing / Weak Evidence | Unknowns | Expected Match | Reasoning |
| --- | --- | --- | --- | --- | --- |
| AWS 구축·운영 3년 이상 / Required | E7 기본 AWS | 지속적인 production 구축·운영 근거 | 기간·책임·규모 | Weak | 기초 경험과 명시된 운영 책임 사이 evidence gap이 크다. |
| Web 서비스 설계·개발 / Required | E1, E2, E3, E6 | 추가 없음 | BtoB/BtoC별 수행 내역 | Strong | API와 업무 Web 설계·개발이 직접 관련된다. |
| 인프라·미들웨어·통신·브라우저 지식 / Required | E3, E5의 TCP/IP·Netty·Redis | 인프라·브라우저 지식 깊이 | 운영 수준 지식 | Partial | 일부 기술 영역 근거는 충분하나 전체 범위는 미입증이다. |
| 모니터링·백업 운영 체계 구축 / Responsibility | E6 선박 모니터링은 다른 맥락 | 서비스 운영·백업 구축 이력 | 운영 체계 담당 여부 | Unknown | 도메인 이름의 모니터링을 인프라 관측·백업 경험으로 바꾸지 않는다. |
| 보안 패치 선정·적용 / Responsibility | E6 DB 암호화는 인접 보안 도메인 | 패치 운영 근거 | 패치 관리 수행 여부 | Unknown | 암호화 경력만으로 패치 업무를 판단할 수 없다. |
| 네트워크~미들웨어 성능·장애 대응 / Responsibility | E2 처리 개선, E5 통신 기술 | 해당 계층의 tuning·장애 대응 사례 | 대상 계층·운영 책임 | Partial | 개선 및 통신 근거가 관련되지만 DB 개선과 동일한 업무는 아니다. |
| Java 등 application 개발 / Preferred | E1 | 추가 없음 | 다른 언어 경험 상세 | Strong | 전문 Java 경험이 직접 대응한다. |
| AWS 자격 / Preferred | 식별된 근거 없음 | 자격 이력 | 보유 여부 | Unknown | AWS 기본 경험으로 자격을 추측하지 않는다. |
| 인프라 구축 PM / Preferred | E4 PM/PL | 인프라 프로젝트 관리 이력 | PM 대상 프로젝트 | Partial | 일반 PM 경험은 전이 가능하나 인프라 맥락이 미확인이다. |

**Overall Expected Match: Partial (Turn 003 확정)**

Why: Web 설계·개발과 일부 engineering 지식에는 직접 또는 전이 가능한 근거가 있어 Partial로 확정한다. Primary Gaps: AWS 3년 구축·운영, 백업·패치 책임. 중요한 AWS 필수의 Weak를 Java 우대나 Oracle 경험으로 상쇄하여 Strong으로 올리지 않는다. 이는 이 sample의 전체 capability 요약이며, 필수 Weak가 있는 모든 공고를 Partial로 만드는 규칙은 아니다.

### Sample 04 — CORE

- Company / Position: CORE / 九州(福岡) 시스템 엔지니어, Web·업무 시스템
- Source URL: [공개 공고](https://www.core.co.jp/recruitment/career/kc1)
- Reviewed Date: 2026-09-21
- Source 확인: 공개 본문 확인. 리더 채용 분기의 필수·우대 및 입사 직후 역할을 기준으로 읽었다.
- Relevant Requirements Summary: 기본설계 이후 경험과 고객 문제 파악이 필수다. 상류 경험·과제 조율은 우대이며 sub-leader 역할을 기대한다. Java/Oracle은 예시 개발환경이다.
- 출처 위치: 리더 지원 조건, 초기 역할, 배정 예시와 개발환경.

| Requirement / 구분 | Supporting Evidence | Missing / Weak Evidence | Unknowns | Expected Match | Reasoning |
| --- | --- | --- | --- | --- | --- |
| 기본설계 이후 개발 / Required | E1, E3 | 추가 없음 | 단계별 산출물 | Strong | 기본·상세설계와 개발 근거가 직접 대응한다. |
| 고객과 대화하여 과제 파악 / Required | E3, E4 | 추가 없음 | 개별 고객 사례 | Strong | 요구사항 분석·specification 조율이 직접 관련된다. |
| 상류 과정·고객 과제 조율 / Preferred | E3, E4 | 추가 없음 | 조율 결과 사례 | Strong | 기술과 고객 협의 양쪽 근거가 있다. |
| Sub-leader, 희망 시 project leader / Responsibility | E4 | 해당 프로젝트 규모 리딩 근거 | 실제 배정 범위 | Strong | PL/Tech Lead가 역할에 직접 관련되나 같은 규모를 주장하지 않는다. |
| Java·Oracle / Stack | E1, E2 | 추가 없음 | 버전·적용 내역 | Strong | 두 기술 모두 직접 경험이 있다. |
| 공공·에너지 업무 시스템 / Context | E2, E6 공공·회계·ERP | 에너지·Polisys 직접 경험 | 해당 도메인 숙련 | Partial | 공공은 직접, 에너지는 인접 업무 근거만 있다. |
| 요구사항부터 시스템 테스트·도입·보수 / Responsibility | E3, E4 | 테스트·현장 도입·보수 책임 | 전 과정 담당 범위 | Partial | 상류·개발 근거는 있으나 전 lifecycle 책임은 미입증이다. |

**Overall Expected Match: Strong (Turn 003 확정)**

Why: 명시된 핵심 필수인 기본설계와 고객 과제 파악, 기대 리딩 역할에 직접 근거가 있다. Primary Gaps: 에너지 특화 제품, 테스트·도입·보수 및 동일 규모 경험. 팀 경험 약 12명을 공고 프로젝트 규모와 같다고 보지 않았다. 학력 조건도 본문에 있으나 이 요약은 capability match이며 전체 지원 자격 충족 여부를 판단하지 않는다.

### Sample 05 — teamLab

- Company / Position: teamLab / 통년 채용 Web app·smartphone app engineer
- Source URL: [공개 공고](https://open.talentio.com/r/1/c/teamlab/pages/109375)
- Reviewed Date: 2026-09-21
- Source 확인: 웹 읽기 및 정적 HTML에는 제목 수준만 나타났다. 동일 URL을 브라우저로 렌더링하여 본문 확인.
- Relevant Requirements Summary: 고객 요구 기반 기획·설계, architecture, 개발 효율화, 성능 개선·CI/CD 업무를 다룬다. 기술 목록은 backend/frontend/mobile별로 제시되며, 필요 skill에는 business-level 일본어가 명시된다.
- 출처 위치: 개요, 주요 업무, 주요 언어, 필요 skill. 모든 나열 기술이 필수라는 문구는 확인하지 못했다.

| Requirement / 구분 | Supporting Evidence | Missing / Weak Evidence | Unknowns | Expected Match | Reasoning |
| --- | --- | --- | --- | --- | --- |
| 고객 요구를 반영한 요구정의·설계 / Responsibility | E3, E4 | 견적·제안 산출물 | 견적 담당 범위 | Strong | 요구정의·설계 축의 직접 근거이며 견적까지 입증한 것은 아니다. |
| 전체 개발방침·architecture 수립 / Responsibility | E2, E3, E4 | 시스템 전체 architecture 결정 사례 | 결정 범위·책임 | Partial | API/DB 설계와 리딩은 관련되나 전체 범위는 미확인이다. |
| 처리 효율화·범용화 / Responsibility | E2, E3 공통 모듈·legacy 개선 | 자동화 구축 사례 | 자동화 범위 | Strong | 효율화·범용화에 직접 근거가 있다. 자동화는 별도 미확인이다. |
| 성능 tuning / Responsibility | E2 SQL/DB 개선 | 서비스 전체 성능 개선 | 개선 대상·지표 | Partial | 한 계층의 근거는 있으나 전체 성능 책임을 보장하지 않는다. |
| CI/CD 구축 / Responsibility | 식별된 근거 없음 | Pipeline 구축 이력 | 구축·운영 여부 | Unknown | 개발 경력만으로 CI/CD를 추정하지 않는다. |
| Java·Spring·SQL backend / Stack | E1, E2 | 추가 없음 | 배정 기술·역할 | Strong | 확인된 backend 기술과 직접 맞는다. |
| Frontend/mobile 목록 / Stack, 대체 경로 여부 미정 | E3, E5 frontend·TypeScript·Node.js | 나열된 UI framework·native mobile 경험 | React/Vue/Swift 등 경험 | Partial | 인접 frontend 근거는 있으나 mobile 경험은 Unknown이다. 모든 기술을 필수로 합산하지 않는다. |
| Business-level 일본어 / 명시 필요 skill | 식별된 근거 없음 | 업무 일본어 evidence | 언어 수준 | Unknown | 고객 조율 경력만으로 일본어 능력을 추측하지 않는다. |

**Overall Expected Match: Partial (Turn 003 확정)**

Why: backend와 고객 요구·설계 경험은 강하지만 전체 architecture·성능 범위는 Partial, CI/CD는 Unknown이다. Primary Gaps: 전체 engineering 역할 범위, CI/CD, 일본어 evidence. Backend와 mobile 중 배정 경로가 정해지지 않았다. 언어 조건을 capability overall에 포함하는 방법도 미정이며, Unknown을 충족으로 처리하지 않았다.

### Sample 06 — Nitori Digital Base

- Company / Position: Nitori 그룹 / EC frontend engineer (NDB0055)
- Source URL: [공개 공고](https://hrmos.co/pages/nitorihd/jobs/NDB0055)
- Reviewed Date: 2026-09-21
- Source 확인: 공개 본문 확인.
- Relevant Requirements Summary: 팀 리딩, Web frontend Tech Lead, 멤버 육성이 필수다. Angular/TypeScript는 frontend stack, Java/Spring Boot는 조직 backend stack이다. 대규모 EC 운영 등은 우대다.
- 출처 위치: 업무 내용, 기술 stack, 필수/우대 조건. Negative control 후보지만 결과를 미리 고정한 정답은 아니다.

| Requirement / 구분 | Supporting Evidence | Missing / Weak Evidence | Unknowns | Expected Match | Reasoning |
| --- | --- | --- | --- | --- | --- |
| 프로젝트 팀 리딩 / Required | E4 | 추가 없음 | 해당 조직 규모 리딩 | Strong | 팀 리딩 경험은 직접 관련된다. |
| Web frontend Tech Lead / Required | E3 frontend 개발, E4 일반 리딩, E5 Angular/TypeScript | Frontend 전담 기술 리딩 근거 | 리딩 분야·기간·책임 | Weak | 인접 근거는 있지만 이를 결합해 frontend Tech Lead 경험을 만들 수 없다. |
| 개발 멤버 육성 / Required | E4 OJT·기술지원·표준 정비 | 추가 없음 | 육성 결과 상세 | Strong | 명시된 육성 활동이 직접 대응한다. |
| Frontend 제품 설계·개발·운영 / Responsibility | E3, E5 | EC frontend 설계·운영 깊이 | 운영 책임·제품 규모 | Partial | Frontend 경험은 있으나 역할 전체 범위는 미입증이다. |
| Frontend 품질·설계/코드 리뷰 / Responsibility | E3, E4 일반 review | Frontend 특화 review 사례 | 대상 기술·품질 판단 범위 | Partial | 일반 리뷰는 전이 가능하지만 분야별 깊이는 확인되지 않는다. |
| Angular·TypeScript / Stack | E5 | 숙련·기간 상세 | 전문성 깊이 | Partial | 기술명 일치는 있지만 전문 리딩을 뒷받침하지는 못한다. |
| 대규모 EC 구축·운영 / Preferred | 식별된 직접 근거 없음 | EC 구축·운영 이력 | 수행 여부·규모 | Unknown | 업무 시스템 경력을 대규모 EC 경험으로 치환하지 않는다. |
| SEO·광고·검색·위치·결제 지식 / Preferred | 식별된 구체 근거 없음 | 해당 서비스 지식 사례 | 실제 보유 지식 | Unknown | 물류·회계 도메인명으로 위치·결제 구현 지식을 추측하지 않는다. |
| 디자인 도구 / Preferred | 식별된 근거 없음 | 도구 사용 이력 | 사용 여부 | Unknown | Frontend 경험만으로 특정 도구 사용을 추론하지 않는다. |

**Overall Expected Match: Weak (Turn 003 확정)**

Why: 일반 리딩·육성은 Strong이지만 이 포지션의 핵심인 frontend Tech Lead를 지지하는 evidence는 제한적이다. Primary Gaps: frontend 전문 기술 의사결정·리딩, EC 운영. 조직의 Java/Spring Boot stack은 candidate의 강점과 겹쳐도 frontend 역할을 입증하지 않으므로 별도 가점처럼 쓰지 않았다. 인접 frontend 개발·일반 리딩 근거와 핵심 역할의 gap을 비교하여 Weak로 확정한다. “Frontend 리딩 경험이 없다”는 사실을 새로 단정한 것은 아니다.

## 6. Overall Baseline Summary

Turn 003의 사용자 제공 규칙과 확정 지시에 따라 기존 6개 값을 유지한 baseline이다. Percentage를 부여하지 않았고 requirement level의 평균도 아니다. 근거가 불충분한 개별 사실은 여전히 Unknown이다.

| Sample | Overall Expected Match | 핵심 근거 / 한계 |
| --- | --- | --- |
| 01 BlueMeme / Micro Court | Strong | 시스템 개발·고객 조율 직접 근거. OutSystems 사전 경험은 불문. |
| 02 Prop Tech plus | Strong | Backend/API 직접 근거. Java 버전별 기간은 미확인. |
| 03 PIA TECH LAB | Partial | Web 설계와 전이 가능한 경험을 인정하되 AWS 필수 운영 gap으로 Strong은 지지하지 않음. |
| 04 CORE | Strong | 기본설계·고객 조율·리딩 직접 근거. 특화 도메인/전 lifecycle은 제한적. |
| 05 teamLab | Partial | Backend 강점, 전체 역할 범위·CI/CD·일본어 미확인. |
| 06 Nitori Digital Base | Weak | 일반 리딩·기술 overlap이 frontend Tech Lead를 입증하지 않음. |

## 7. Turn 002 당시 Ambiguities / Planning Review

아래 A–H는 당시 초안의 질문과 상태를 보존한 기록이다. 현재 결정은 3절 및 8절을 따른다. 아래의 “잠정”, “검토 필요”는 현재 baseline 상태를 뜻하지 않는다.

### A. Transferable experience와 Strong / Partial

BlueMeme는 일반 시스템 개발을 필수로 두고 OutSystems 사전 경험을 요구하지 않는다. 따라서 일반 개발 requirement는 Strong, 도구 자체는 Unknown으로 분리했다. Overall Strong 초안을 작성했으나, 당장 수행할 업무의 도구 숙련까지 overall이 표현해야 한다면 Partial 해석도 가능하다. Overall이 입사 전 요구 충족과 현재 업무 수행 범위 중 무엇을 요약하는지 결정이 필요하다.

### B. 핵심 specialized requirement 하나가 부족할 때

Prop Tech의 Azure 우대 Unknown과 PIA의 AWS 필수 Weak는 성격이 다르다. PIA는 Web 개발 근거가 충분해 Partial 초안으로 기록했지만 AWS 중심 책임을 보면 Weak도 가능하다. 특정 필수의 gap이 overall을 제한하는 방식은 미정이다. Prop Tech의 Java 버전·기간 미확인 상태에서 Strong 요약이 과도한지도 검토한다.

### C. No evidence와 Weak / Unknown

Kubernetes처럼 관련 정보가 전혀 없으면 Unknown이다. Nitori처럼 일반 리딩과 frontend 사용은 각각 존재하지만 둘의 결합 역할이 입증되지 않는 경우, 인접 근거가 있으므로 Weak인지 해당 역할 직접 정보가 없으므로 Unknown인지 경계가 남는다. 이번 Weak는 잠정 해석이다. 직무 경험이 실제로 없다는 negative evidence를 받았을 때의 처리도 아직 정의하지 않았다.

### D. Keyword overlap과 핵심 role

Nitori의 일반 리딩·육성은 Strong으로 보존하되, backend stack이나 Angular 이름으로 frontend Tech Lead evidence를 보충하지 않았다. Negative control은 낮은 점수를 강제하는 label이 아니라 과대평가가 발생하는지 관찰할 사례다. 추가적인 직접 frontend 리딩 evidence가 생기면 판단을 다시 검토해야 한다.

### E. Required / Preferred / Responsibility / Stack

필수와 우대는 같지 않다. 동시에 실제 업무와 단순 환경도 구분해야 한다. CORE의 Java/Oracle 환경을 필수 조건으로, teamLab의 모든 frontend/mobile 기술을 동시 필수로 바꾸지 않았다. 이 구분을 전체 요약에 반영하는 상세 규칙은 미정이며 weighting은 만들지 않았다.

### F. 정보의 해상도와 requirement 분해

“Java 경험”은 있어도 “Java 8+ 1년”은 확인되지 않는다. “개발 효율화”를 만족해도 “자동화/CI/CD”까지 입증되는 것은 아니다. 복합 요구의 분해 단위와 버전·기간·scope 미확인 시 level을 표현하는 방식은 추가 검토가 필요하다.

### G. 언어와 지원 자격의 범위

teamLab의 일본어는 원문 필요 skill이므로 Unknown으로 기록했지만, 현재 기술 중심 Capability / Experience Match에서 overall에 반영하는 방식은 정의되지 않았다. CORE 학력처럼 capability와 구분되는 지원 조건까지 이 도구가 검증하는 것은 아니다. 언어 조건을 다룰 범위를 planning conversation에서 명시해야 한다.

### H. Baseline의 독립성 및 재현성

이번 판단은 AI assistance를 받은 문서 초안이므로 독립적인 human label은 아니다. 사용자가 engine 결과를 보기 전에 Expected Match와 이유를 검토·수정한 기록이 필요하다. 향후 JD 변경 시 다른 조건으로 비교하지 않도록 URL·날짜·요약을 함께 확인해야 한다. 전체 원문 snapshot이나 별도 dataset 체계는 이번에 도입하지 않았다.

## 8. Turn 003 Review Outcome

| 기존 질문 | 이번 결정 또는 남은 범위 |
| --- | --- |
| A: Transferable experience | 문제·책임·설계 깊이·운영·ownership·기술 전환을 고려한다. BlueMeme는 사전 도구 경험을 필수화하지 않아 Strong 유지. |
| B: Specialized requirement gap | Required gap은 중요하고 Preferred 부족은 전체를 크게 낮추지 않는다. PIA Partial, Prop Tech Strong 유지. 사실의 미확인은 계속 기록한다. |
| C: Weak / Unknown | 비교 가능한 관련 근거와 큰 gap은 Weak, 평가 근거 부족은 Unknown. Nitori의 인접 근거와 역할 gap은 Weak로 판단. |
| D: Keyword와 role | Core role mismatch가 skill overlap보다 우선. Nitori Overall Weak 유지. |
| E: 요구 구분 | Required/Preferred의 정성적 중요도 차이 확정. Responsibility/Stack을 필수로 승격하지 않는 구분 유지. |
| F: 세부 정보와 분해 | 기존 판단·Unknown 보존. 복합 requirement의 표준 분해 단위는 미정. |
| G: 언어·지원 자격 | 일본어 Unknown 보존. 언어와 학력의 일반적인 평가 범위는 이번 지시에서 결정하지 않아 후속 검토로 남김. |
| H: Baseline 상태 | Turn 003 지시에 따라 비교 기준 확정. Codex 초안에서 출발했다는 provenance 유지. 독립 검증·ground truth를 주장하지 않음. |

이번 확정은 기록된 JD와 reference evidence에 대한 정성적 기준이다. 채용 확률, 지원 여부, 모든 필수 지원 자격 충족 또는 실제 경험의 외부 검증을 뜻하지 않는다. JD를 다시 수집하거나 Jev/별도 evaluator를 실행하지 않았다.

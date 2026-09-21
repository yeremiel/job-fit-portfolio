# Job Fit Evaluation API

확보한 JD를 서버의 canonical candidate profile과 비교하여 structured evaluation을 반환한다. 공통 응답 field와 전체 error 목록은 [API Specification](specification.md)을 참조한다.

---

## Job Fit 평가

```http
POST /v1/job-fit/evaluate
```

기본 URL: `http://127.0.0.1:8080/v1/job-fit/evaluate`.

## Request

### Headers

```http
Content-Type: application/json
```

`application/json; charset=utf-8`도 허용한다. Content-Type은 필수이며 다른 parameter/charset은 거부한다. Content-Encoding은 보내지 않는다(압축뿐 아니라 명시적 `identity` 값도 현재 구현에서 거부). Authentication header나 API client ID는 필요하지 않다. Jev key를 client header/body에 보내지 않는다. Query parameter는 지원하지 않는다.

### Body

```json
{
  "job": {
    "company": "Example Corp",
    "title": "Backend Engineer",
    "sourceUrl": "https://example.com/jobs/123",
    "description": "Position: Backend Engineer\nRequired: Backend API design and development"
  }
}
```

### Fields

| Field | Type | Required | Description / validation |
| --- | --- | --- | --- |
| `job` | object | Yes | 최상위 유일한 field. Null/array 불가 |
| `job.description` | string | Yes | 평가할 JD 본문. Trim 후 비어 있지 않아야 함. Trim 전 UTF-8 16 KiB 이하, 최대 64 non-empty lines |
| `job.company` | string | No | 회사명 metadata. Trim 전 UTF-8 256 bytes 이하 |
| `job.title` | string | No | 포지션명 metadata. Trim 전 UTF-8 512 bytes 이하 |
| `job.sourceUrl` | string | No | Provenance metadata. Trim 전 UTF-8 2,048 bytes 이하. Absolute HTTP/HTTPS URL, hostname 필수, userinfo 금지 |

Optional field는 생략할 수 있지만 제공 시 null 또는 blank string이면 거부한다. Metadata는 trim하여 반환한다. 올바른 UTF-8 JSON object 하나만 보내며 duplicate/unknown keys, trailing JSON, NUL, 잘못된 Unicode surrogate를 허용하지 않는다. Field 이름은 대소문자를 구분한다. `candidateProfile`, `profileId`, `profileVersion` 같은 추가 input field는 받지 않는다.

JD는 한 줄에 한 요구·업무·환경 항목을 준비한다. 전체 앞뒤 공백 제거 후 각 줄을 trim하고 빈 줄을 제외한다. 자동 HTML parsing, 요약, truncation 또는 일반 문서 requirement decomposition을 제공하지 않는다.

## Behavior

- Profile은 request로 받지 않는다. Server startup에서 읽은 canonical snapshot을 모든 요청에 재사용하며 profile 변경은 restart 후 반영된다.
- `data.metadata.profileVersion`은 실제 읽은 profile file bytes의 SHA-256 fingerprint다. Formatting-only 변경도 fingerprint를 바꿀 수 있으며 semantic version이나 profile 조회 ID가 아니다.
- Company/title/sourceUrl은 식별용 metadata다. Evaluation description에 삽입하거나 URL을 fetch하지 않는다. 평가에 필요한 role context는 description에 포함한다.
- Evaluation은 synchronous다. 성공은 전체 결과 하나이며 실패 시 partial result를 반환하지 않는다.
- 결과를 저장하지 않으며 같은 요청도 새 평가다. 결과 label의 동일성, idempotency, 자동 retry를 보장하지 않는다.

---

## Success Response

**200 OK**. 아래는 [Turn 014 Prop Tech plus 실제 HTTP 응답](../development/turn-014-results/proptech-plus.json)을 생략 없이 복사한 예시다. 위 Example Corp request의 결과가 아니며, 해당 실행은 [Prop Tech plus sample](../../samples/proptech-plus.txt) 전체를 description으로 사용했다. Timestamp/requestId/profileVersion은 당시 값이다. Caller는 같은 label·문장·ID 목록이 반복된다고 가정하지 않는다.

```json
{
  "status": "success",
  "code": 200,
  "message": "ok",
  "data": {
    "job": {
      "company": "Prop Tech plus",
      "title": "Server-side engineer"
    },
    "overallMatch": "Partial",
    "overallReasoning": "Selected overall criterion: Substantial related/transferable capabilities exist but some important scope/depth/technology/domain/operational gaps remain. Assessments: L002=Strong (context); L005=Partial (context); L006=Partial (required); L007=Strong (responsibility); L008=Partial (responsibility); L009=Partial (responsibility); L010=Partial (responsibility); L011=Strong (stack); L012=Partial (stack); L013=Unknown (stack); L014=Unknown (preferred); L015=Weak (preferred); L016=Unknown (preferred); L017=Unknown (preferred); L018=Partial (preferred)",
    "requirements": [
      {
        "id": "L002",
        "requirement": "Position: Server-side engineer",
        "sourceType": "context",
        "match": "Strong",
        "supportingEvidence": [
          {
            "id": "E1",
            "text": "Approximately 16 years of software/backend development, centered on professional Java, Spring and Spring Boot. REST API design and implementation. Version-specific durations are not documented.",
            "relation": "direct"
          },
          {
            "id": "E2",
            "text": "RDB design and SQL/database processing improvements. DB2, Oracle, MySQL and MariaDB experience. Business systems with complex business logic.",
            "relation": "transferable"
          },
          {
            "id": "E3",
            "text": "Requirements analysis, basic and detailed design, API/database design, legacy analysis and improvement, common module design and code review. Development across backend, frontend and database.",
            "relation": "direct"
          },
          {
            "id": "E4",
            "text": "Tech Lead, Development Lead, PM and PL experience. Team building, OJT, development standards, code review and technical support. Experience in organizations/teams up to approximately 12 people. Customer requirements and specification coordination. The profile does not specify frontend-specialist leadership.",
            "relation": "limited"
          },
          {
            "id": "E5",
            "text": "Experience with TypeScript, Node.js, NestJS, Angular, Redis, Netty, C++ and TCP/IP. Individual durations and depth are not documented.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [],
        "unknowns": [],
        "reasoning": "Selected rationale: Direct documented evidence supports the requirement with no material role/context gap."
      },
      {
        "id": "L005",
        "requirement": "Business context: In-house package development for real-estate fund management companies.",
        "sourceType": "context",
        "match": "Partial",
        "supportingEvidence": [
          {
            "id": "E2",
            "text": "RDB design and SQL/database processing improvements. DB2, Oracle, MySQL and MariaDB experience. Business systems with complex business logic.",
            "relation": "transferable"
          },
          {
            "id": "E6",
            "text": "Domain experience in ERP, accounting, personal-information database encryption, ship monitoring, logistics/control systems and public-sector web/business systems.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [
          "Not established by the profile: specific domain experience for: Business context: In-house package development for real-estate fund management companies."
        ],
        "unknowns": [
          "Whether the candidate has specific domain experience beyond the documented evidence."
        ],
        "reasoning": "Selected rationale: Transferable engineering experience exists but the specific domain experience is not established."
      },
      {
        "id": "L006",
        "requirement": "Required: At least one year of development experience with Java 8 or later.",
        "sourceType": "required",
        "match": "Partial",
        "supportingEvidence": [
          {
            "id": "E1",
            "text": "Approximately 16 years of software/backend development, centered on professional Java, Spring and Spring Boot. REST API design and implementation. Version-specific durations are not documented.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [
          "Not established by the profile: specified version or duration for: Required: At least one year of development experience with Java 8 or later."
        ],
        "unknowns": [
          "Whether the candidate has specified version or duration beyond the documented evidence."
        ],
        "reasoning": "Selected rationale: Related experience exists but the specified version or duration is not documented."
      },
      {
        "id": "L007",
        "requirement": "Responsibility: Design and develop backend APIs.",
        "sourceType": "responsibility",
        "match": "Strong",
        "supportingEvidence": [
          {
            "id": "E1",
            "text": "Approximately 16 years of software/backend development, centered on professional Java, Spring and Spring Boot. REST API design and implementation. Version-specific durations are not documented.",
            "relation": "direct"
          },
          {
            "id": "E2",
            "text": "RDB design and SQL/database processing improvements. DB2, Oracle, MySQL and MariaDB experience. Business systems with complex business logic.",
            "relation": "limited"
          },
          {
            "id": "E3",
            "text": "Requirements analysis, basic and detailed design, API/database design, legacy analysis and improvement, common module design and code review. Development across backend, frontend and database.",
            "relation": "direct"
          },
          {
            "id": "E5",
            "text": "Experience with TypeScript, Node.js, NestJS, Angular, Redis, Netty, C++ and TCP/IP. Individual durations and depth are not documented.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [],
        "unknowns": [],
        "reasoning": "Selected rationale: Direct documented evidence supports the requirement with no material role/context gap."
      },
      {
        "id": "L008",
        "requirement": "Responsibility: Operate backend APIs.",
        "sourceType": "responsibility",
        "match": "Partial",
        "supportingEvidence": [
          {
            "id": "E1",
            "text": "Approximately 16 years of software/backend development, centered on professional Java, Spring and Spring Boot. REST API design and implementation. Version-specific durations are not documented.",
            "relation": "limited"
          },
          {
            "id": "E3",
            "text": "Requirements analysis, basic and detailed design, API/database design, legacy analysis and improvement, common module design and code review. Development across backend, frontend and database.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [
          "Not established by the profile: operational context for: Responsibility: Operate backend APIs."
        ],
        "unknowns": [
          "Whether the candidate has operational context beyond the documented evidence."
        ],
        "reasoning": "Selected rationale: Related evidence exists but the required operational context is not established."
      },
      {
        "id": "L009",
        "requirement": "Responsibility: Enhance the package and fix defects.",
        "sourceType": "responsibility",
        "match": "Partial",
        "supportingEvidence": [
          {
            "id": "E1",
            "text": "Approximately 16 years of software/backend development, centered on professional Java, Spring and Spring Boot. REST API design and implementation. Version-specific durations are not documented.",
            "relation": "limited"
          },
          {
            "id": "E2",
            "text": "RDB design and SQL/database processing improvements. DB2, Oracle, MySQL and MariaDB experience. Business systems with complex business logic.",
            "relation": "limited"
          },
          {
            "id": "E3",
            "text": "Requirements analysis, basic and detailed design, API/database design, legacy analysis and improvement, common module design and code review. Development across backend, frontend and database.",
            "relation": "limited"
          },
          {
            "id": "E4",
            "text": "Tech Lead, Development Lead, PM and PL experience. Team building, OJT, development standards, code review and technical support. Experience in organizations/teams up to approximately 12 people. Customer requirements and specification coordination. The profile does not specify frontend-specialist leadership.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [
          "Not established by the profile: full scope for: Responsibility: Enhance the package and fix defects."
        ],
        "unknowns": [
          "Whether the candidate has full scope beyond the documented evidence."
        ],
        "reasoning": "Selected rationale: Related evidence is sufficient but the full scope is not documented."
      },
      {
        "id": "L010",
        "requirement": "Responsibility: Set up and customize the package for new customers.",
        "sourceType": "responsibility",
        "match": "Partial",
        "supportingEvidence": [
          {
            "id": "E1",
            "text": "Approximately 16 years of software/backend development, centered on professional Java, Spring and Spring Boot. REST API design and implementation. Version-specific durations are not documented.",
            "relation": "limited"
          },
          {
            "id": "E2",
            "text": "RDB design and SQL/database processing improvements. DB2, Oracle, MySQL and MariaDB experience. Business systems with complex business logic.",
            "relation": "limited"
          },
          {
            "id": "E3",
            "text": "Requirements analysis, basic and detailed design, API/database design, legacy analysis and improvement, common module design and code review. Development across backend, frontend and database.",
            "relation": "limited"
          },
          {
            "id": "E4",
            "text": "Tech Lead, Development Lead, PM and PL experience. Team building, OJT, development standards, code review and technical support. Experience in organizations/teams up to approximately 12 people. Customer requirements and specification coordination. The profile does not specify frontend-specialist leadership.",
            "relation": "limited"
          },
          {
            "id": "E6",
            "text": "Domain experience in ERP, accounting, personal-information database encryption, ship monitoring, logistics/control systems and public-sector web/business systems.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [
          "Not established by the profile: full scope for: Responsibility: Set up and customize the package for new customers."
        ],
        "unknowns": [
          "Whether the candidate has full scope beyond the documented evidence."
        ],
        "reasoning": "Selected rationale: Related evidence is sufficient but the full scope is not documented."
      },
      {
        "id": "L011",
        "requirement": "Technology environment: Java and Spring Boot.",
        "sourceType": "stack",
        "match": "Strong",
        "supportingEvidence": [
          {
            "id": "E1",
            "text": "Approximately 16 years of software/backend development, centered on professional Java, Spring and Spring Boot. REST API design and implementation. Version-specific durations are not documented.",
            "relation": "direct"
          },
          {
            "id": "E3",
            "text": "Requirements analysis, basic and detailed design, API/database design, legacy analysis and improvement, common module design and code review. Development across backend, frontend and database.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [],
        "unknowns": [],
        "reasoning": "Selected rationale: Direct documented evidence supports the requirement with no material role/context gap."
      },
      {
        "id": "L012",
        "requirement": "Technology environment: Azure App Service on Linux/Docker and Azure SQL Database.",
        "sourceType": "stack",
        "match": "Partial",
        "supportingEvidence": [
          {
            "id": "E2",
            "text": "RDB design and SQL/database processing improvements. DB2, Oracle, MySQL and MariaDB experience. Business systems with complex business logic.",
            "relation": "transferable"
          }
        ],
        "missingEvidence": [
          "Not established by the profile: exact technology experience for: Technology environment: Azure App Service on Linux/Docker and Azure SQL Database."
        ],
        "unknowns": [
          "Whether the candidate has exact technology experience beyond the documented evidence."
        ],
        "reasoning": "Selected rationale: Transferable engineering experience exists but the exact technology experience is not established."
      },
      {
        "id": "L013",
        "requirement": "Technology environment: Git and Azure DevOps for CI/CD.",
        "sourceType": "stack",
        "match": "Unknown",
        "supportingEvidence": [],
        "missingEvidence": [
          "Not established by the profile: sufficient relevant evidence for: Technology environment: Git and Azure DevOps for CI/CD."
        ],
        "unknowns": [
          "Whether the candidate has sufficient relevant evidence beyond the documented evidence."
        ],
        "reasoning": "No supporting evidence was linked by the evaluator. The application leaves the match Unknown; this does not mean no experience."
      },
      {
        "id": "L014",
        "requirement": "Preferred: Hands-on Azure operation/administration experience.",
        "sourceType": "preferred",
        "match": "Unknown",
        "supportingEvidence": [],
        "missingEvidence": [
          "Not established by the profile: sufficient relevant evidence for: Preferred: Hands-on Azure operation/administration experience."
        ],
        "unknowns": [
          "Whether the candidate has sufficient relevant evidence beyond the documented evidence."
        ],
        "reasoning": "No supporting evidence was linked by the evaluator. The application leaves the match Unknown; this does not mean no experience."
      },
      {
        "id": "L015",
        "requirement": "Preferred: Accounting knowledge comparable to bookkeeping level 3.",
        "sourceType": "preferred",
        "match": "Weak",
        "supportingEvidence": [
          {
            "id": "E6",
            "text": "Domain experience in ERP, accounting, personal-information database encryption, ship monitoring, logistics/control systems and public-sector web/business systems.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [
          "Not established by the profile: evidence at the required level for: Preferred: Accounting knowledge comparable to bookkeeping level 3."
        ],
        "unknowns": [
          "Whether the candidate has evidence at the required level beyond the documented evidence."
        ],
        "reasoning": "Selected rationale: Some related evidence exists but it falls substantially short of the requirement's level."
      },
      {
        "id": "L016",
        "requirement": "Preferred: Program with the customer's usability in mind.",
        "sourceType": "preferred",
        "match": "Unknown",
        "supportingEvidence": [
          {
            "id": "E3",
            "text": "Requirements analysis, basic and detailed design, API/database design, legacy analysis and improvement, common module design and code review. Development across backend, frontend and database.",
            "relation": "limited"
          },
          {
            "id": "E4",
            "text": "Tech Lead, Development Lead, PM and PL experience. Team building, OJT, development standards, code review and technical support. Experience in organizations/teams up to approximately 12 people. Customer requirements and specification coordination. The profile does not specify frontend-specialist leadership.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [
          "Not established by the profile: sufficient relevant evidence for: Preferred: Program with the customer's usability in mind."
        ],
        "unknowns": [
          "Whether the candidate has sufficient relevant evidence beyond the documented evidence."
        ],
        "reasoning": "Selected rationale: The profile contains insufficient relevant evidence to judge this requirement; this does not mean no experience."
      },
      {
        "id": "L017",
        "requirement": "Preferred: Willingness to understand the customer's business and work proactively.",
        "sourceType": "preferred",
        "match": "Unknown",
        "supportingEvidence": [
          {
            "id": "E3",
            "text": "Requirements analysis, basic and detailed design, API/database design, legacy analysis and improvement, common module design and code review. Development across backend, frontend and database.",
            "relation": "limited"
          },
          {
            "id": "E4",
            "text": "Tech Lead, Development Lead, PM and PL experience. Team building, OJT, development standards, code review and technical support. Experience in organizations/teams up to approximately 12 people. Customer requirements and specification coordination. The profile does not specify frontend-specialist leadership.",
            "relation": "limited"
          },
          {
            "id": "E6",
            "text": "Domain experience in ERP, accounting, personal-information database encryption, ship monitoring, logistics/control systems and public-sector web/business systems.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [
          "Not established by the profile: sufficient relevant evidence for: Preferred: Willingness to understand the customer's business and work proactively."
        ],
        "unknowns": [
          "Whether the candidate has sufficient relevant evidence beyond the documented evidence."
        ],
        "reasoning": "Selected rationale: The profile contains insufficient relevant evidence to judge this requirement; this does not mean no experience."
      },
      {
        "id": "L018",
        "requirement": "Optional adjacent web development: Angular, Angular Material and TypeScript.",
        "sourceType": "preferred",
        "match": "Partial",
        "supportingEvidence": [
          {
            "id": "E3",
            "text": "Requirements analysis, basic and detailed design, API/database design, legacy analysis and improvement, common module design and code review. Development across backend, frontend and database.",
            "relation": "limited"
          },
          {
            "id": "E5",
            "text": "Experience with TypeScript, Node.js, NestJS, Angular, Redis, Netty, C++ and TCP/IP. Individual durations and depth are not documented.",
            "relation": "limited"
          }
        ],
        "missingEvidence": [
          "Not established by the profile: requested depth for: Optional adjacent web development: Angular, Angular Material and TypeScript."
        ],
        "unknowns": [
          "Whether the candidate has requested depth beyond the documented evidence."
        ],
        "reasoning": "Selected rationale: Related evidence is sufficient but the requested depth is not documented."
      }
    ],
    "decisionTrace": {
      "source": "jev-post-hoc-attribution; not a causal record of the original decision",
      "supporting": [
        "L002",
        "L007",
        "L011"
      ],
      "limiting": [
        "L005",
        "L006",
        "L012"
      ],
      "nonDecisive": [
        "L008",
        "L009",
        "L010",
        "L013",
        "L014",
        "L015",
        "L016",
        "L017",
        "L018"
      ]
    },
    "explanationSource": "application-rendered rationale from provider choices or absence of linked evidence; evidence copied from profile",
    "metadata": {
      "profileVersion": "sha256:cae86ad4906124b0eee8283c05bc74d439543f75d9b2e7496f8d6895e32f09cd",
      "evaluatedAt": "2026-09-21T11:32:25.266377Z",
      "engine": "jev",
      "engineModel": "jev-1.13.0"
    }
  },
  "metadata": {
    "timestamp": "2026-09-21T11:32:25.266443Z",
    "requestId": "0e9e452b-c711-49c9-994e-3e0d2facbb46"
  }
}
```

## Response Field Description

공통 envelope는 [Common Response Envelope](specification.md#common-response-envelope)를 참조한다. 아래는 `data` 내부 field다.

| Field | Meaning |
| --- | --- |
| `job` | 제공한 company/title/sourceUrl을 normalize하여 반환. 생략한 field는 반환하지 않으며 모두 없으면 `{}`. Description 전체는 echo하지 않음 |
| `overallMatch` | Strong/Partial/Weak/Unknown summary signal. 객관적 정답·hiring probability·지원 추천이 아님 |
| `overallReasoning` | Provider의 선택을 application이 렌더링한 설명. 자유 생성 rationale이나 chain of thought가 아님 |
| `requirements` | Source order로 정렬된 requirement별 assessment 배열 |
| `requirements[].id` | 평가 결과 안에서 unique한 ID. 현재 non-empty JD line 순서 기반 `L001` 형태. Ignore된 줄로 번호가 건너뛸 수 있으며 JD 변경 전후 영구 ID가 아님 |
| `requirements[].requirement` | 평가한 JD line text |
| `requirements[].sourceType` | `required`, `preferred`, `responsibility`, `stack`, `context`. Internal `category` 대신 사용하는 public field |
| `requirements[].match` | 네 Match level 중 하나 |
| `requirements[].supportingEvidence` | Profile evidence의 `id`, 원문 `text`, 관계 `relation` 배열. Relation은 `direct`, `transferable`, `limited`이며 모델의 연결 판단을 뜻함 |
| `requirements[].missingEvidence` | 현재 profile로 충분히 입증하지 못하는 정보. 실제 경험 부재를 뜻하지 않음 |
| `requirements[].unknowns` | 현재 입력/profile로 확인할 수 없는 정보 |
| `requirements[].reasoning` | 해당 항목의 application-rendered 설명 |
| `decisionTrace.source` | Structured post-hoc attribution이라는 설명 출처/한계 |
| `decisionTrace.supporting` | Overall을 높은 방향으로 지지하는 requirement ID 목록 |
| `decisionTrace.limiting` | 더 높은 Overall을 제한하는 requirement ID 목록 |
| `decisionTrace.nonDecisive` | Overall의 핵심 결정 요인이 아닌 requirement ID 목록 |
| `explanationSource` | 설명 생성 방식과 evidence 원문의 출처 |
| `metadata.profileVersion` | 평가에 사용한 canonical snapshot의 `sha256:<64 lowercase hex>` |
| `metadata.evaluatedAt` | 평가 및 검증 완료 시각(UTC RFC 3339). 공통 metadata.timestamp는 응답 생성 시각 |
| `metadata.engine` | 현재 `jev` |
| `metadata.engineModel` | 실제 요청에 설정한 모델 identifier. Turn 014는 `jev-1.13.0` |

Requirement 각각은 세 Trace 목록 중 정확히 하나에 존재한다. ID로 requirements를 찾아 Match, sourceType, supporting/missing evidence, Unknowns를 함께 읽는다. Strong Overall은 limiting이 비어 있고, Unknown 또는 근거가 없는 항목은 Supporting이 될 수 없다. 오류 결과를 고쳐 성공으로 반환하지 않는다.

**Decision Trace는 causal explanation이 아니라 post-hoc attribution이다.** 원래 Overall 판단의 내부 추론이나 implicit weighting 증거로 읽지 않는다. String 설명은 display용이며 client 분기 조건으로 exact-match하지 않는다. Empty arrays는 `[]`이며 `null`이 아니다. Profile 전체, provider 원문, probability/confidence, CLI의 `summary`는 HTTP response에 포함하지 않는다. 연결된 경력 evidence 원문은 포함하므로 결과를 다룰 때 개인 정보임을 고려한다.

## Match Semantics

| Level | Meaning |
| --- | --- |
| Strong | 핵심 requirement를 직접 뒷받침하는 명확한 evidence가 있고 role/context 차이가 크지 않음 |
| Partial | 관련·transferable evidence가 있으나 일부 scope/depth/technology/domain/operational gap이 있음 |
| Weak | 관련 evidence는 있으나 핵심 requirement 또는 요구 수준과 gap이 큼 |
| Unknown | 평가할 evidence가 충분하지 않음. 경험 부재를 의미하지 않음 |

상세 정의는 [Match Semantics](../match-semantics.md)를 참조한다.

---

## Error Responses

다음 JSON은 현재 mapper의 실제 code/message/field 구조를 사용한 예시다. Timestamp/requestId는 예시값이며 이번 문서 작업에서 실패 요청을 실행해 얻은 artifact는 아니다. 실패에 data는 없고 secret/upstream raw body를 노출하지 않는다.

### 400 Bad Request

Malformed JSON 또는 빈 body, trailing JSON, 잘못된 encoding:

```json
{
  "status": "error",
  "code": 400,
  "error": "REQUEST_MALFORMED_JSON",
  "message": "Request must contain one valid UTF-8 JSON object.",
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

문법은 올바르지만 잘못된 구조·타입·null·unknown/duplicate field, query parameter이면 code는 400, error는 `REQUEST_INVALID`, message는 `Request structure is invalid.`다. 예를 들어 `job` 자체가 없으면 400이고, `job: {}`에서 description만 누락된 경우는 422다.

### 413 Payload Too Large

Body뿐 아니라 description/metadata bytes 또는 description line limit 초과에도 적용한다.

```json
{
  "status": "error",
  "code": 413,
  "error": "REQUEST_INPUT_TOO_LARGE",
  "message": "Input exceeds the allowed size.",
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

### 415 Unsupported Media Type

Content-Type 누락/JSON이 아닌 media type/미지원 parameter 또는 Content-Encoding:

```json
{
  "status": "error",
  "code": 415,
  "error": "REQUEST_UNSUPPORTED_MEDIA_TYPE",
  "message": "Use application/json with UTF-8 and no content encoding.",
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

### 422 Unprocessable Entity

`{"job":{}}`처럼 description이 누락된 요청:

```json
{
  "status": "error",
  "code": 422,
  "error": "REQUEST_VALIDATION_FAILED",
  "message": "Request validation failed.",
  "errors": [
    {
      "field": "job.description",
      "reason": "required",
      "message": "Description is required."
    }
  ],
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

다른 validation도 동일 envelope/errors 배열을 사용한다. 현재 검출된 오류 하나를 반환하며 모든 오류를 모아 반환하지 않는다.

| 상황 | field | reason | errors[].message |
| --- | --- | --- | --- |
| Description 누락 | `job.description` | `required` | `Description is required.` |
| Description 또는 제공된 optional field가 blank | 해당 `job.*` field | `blank` | `Field must not be blank.` |
| 상대 URL, HTTP/HTTPS 이외, hostname 없음, userinfo 포함 | `job.sourceUrl` | `invalid_format` | `Source URL must be an absolute HTTP or HTTPS URL without credentials.` |
| NUL 포함 | 해당 `job.*` field | `invalid_character` | `NUL is not allowed.` |
| Shared input validation 실패 | `job.description` | `invalid_format` | `Description is invalid.` |
| 평가 대상 capability requirement 없음 | `job.description` | `no_capability_requirements` | `No in-scope capability requirements were identified.` |
| Requirement/evidence pairs 또는 생성된 provider payload 제한 초과 | `job.description` | `evaluation_capacity_exceeded` | `Description exceeds evaluation capacity.` |

마지막 두 경우는 이미 일부 Jev 호출이 진행됐을 수 있다. 구조상 null/non-string은 422가 아닌 400이다.

### 500 Internal Server Error

Profile/credential 설정 문제 또는 upstream 401/403:

```json
{
  "status": "error",
  "code": 500,
  "error": "INTERNAL_CONFIGURATION_ERROR",
  "message": "Evaluation service is not configured correctly.",
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

Startup에서 key/profile 문제가 발견되면 서버는 시작하지 않는다. 요청 처리 중 예상하지 못한 내부 실패는 다음과 같다.

```json
{
  "status": "error",
  "code": 500,
  "error": "INTERNAL_ERROR",
  "message": "An internal error occurred.",
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

### 502 Bad Gateway

Jev network/응답/평가 실패 또는 invalid/consistency result:

```json
{
  "status": "error",
  "code": 502,
  "error": "INTERNAL_EVALUATION_FAILED",
  "message": "Evaluation could not be completed.",
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

### 503 Service Unavailable

현재 Jev upstream 429/529 busy/rate limit을 아래로 mapping한다. Upstream HTTP 코드를 그대로 노출하지 않는다.

```json
{
  "status": "error",
  "code": 503,
  "error": "INTERNAL_DEPENDENCY_UNAVAILABLE",
  "message": "Evaluation dependency is unavailable.",
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

### 504 Gateway Timeout

전체 평가 deadline 또는 개별 Jev 호출 timeout:

```json
{
  "status": "error",
  "code": 504,
  "error": "INTERNAL_EVALUATION_TIMEOUT",
  "message": "Evaluation timed out.",
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

Caller가 연결을 끊으면 downstream 요청을 취소하지만 JSON/499 응답 전달을 보장하지 않는다. 오류 발생 전 remote 작업/비용이 이미 발생했을 수 있다. 결과 persistence와 자동 retry가 없으므로 무조건 재전송하는 workflow로 해석하지 않는다.

등록되지 않은 path는 `404 RESOURCE_NOT_FOUND` (`Resource not found.`), 정의된 path의 다른 method는 `405 REQUEST_METHOD_NOT_ALLOWED` (`Method not allowed.` 및 `Allow: POST`)다.

---

## curl Example

아래는 정상 실행 중인 local server에 보낼 수 있는 형식 예제다. 실제 요청은 Jev 호출을 수행한다. 결과 label은 예측하거나 고정하지 않는다.

```sh
curl --request POST http://127.0.0.1:8080/v1/job-fit/evaluate \
  --header 'Content-Type: application/json' \
  --data-binary '{
    "job": {
      "company": "Example Corp",
      "title": "Backend Engineer",
      "sourceUrl": "https://example.com/jobs/123",
      "description": "Position: Backend Engineer\nRequired: Backend API design and development"
    }
  }'
```

실제 JD의 JSON encoding은 request.json 파일로 준비할 수도 있다. 위 Body 형식의 파일을 저장한 뒤:

```sh
curl --request POST http://127.0.0.1:8080/v1/job-fit/evaluate \
  --header 'Content-Type: application/json' \
  --data-binary @request.json
```

## Server Requirements

Server 실행 shell에 `TYPESAFE_API_KEY`가 설정되어 있어야 하며 읽을 수 있는 canonical profile 파일이 필요하다. 빌드된 server 기준, repository root에서:

```sh
./bin/job-fit-server --profile ./data/profile.json
```

빌드/환경변수 준비는 [README](../../README.md#quick-start)를 참조한다. Server는 `.env`를 직접 읽지 않는다. Profile은 기존 evidence `{id,text}` 형식과 validation을 만족해야 한다. 요청별 profile override, profile CRUD 또는 reload endpoint는 없다.

## Request Limits

| 대상 | 실제 한도 / 동작 |
| --- | --- |
| HTTP request body | 128 KiB (131,072 bytes), JSON escape를 포함한 wire body |
| Decoded description | Trim 전 UTF-8 16 KiB (16,384 bytes) |
| Description lines | Trim한 non-empty line 최대 64개. 빈 줄은 제외 |
| Company / title / sourceUrl | 각각 trim 전 UTF-8 256 / 512 / 2,048 bytes |
| Canonical profile file / evidence | 기존 16 KiB file limit, evidence 1–32개 |
| Requirement × evidence pairs | 최대 512개. 추출 후 확인하므로 422가 evaluation 도중 발생 가능 |
| 생성된 Jev request payload | 기존 192 KiB. 초과 시 422 capacity failure |

입력 제한을 만족해도 외부 평가 성공을 보장하지 않는다. 이 값은 resource safeguards이며 Match threshold가 아니다.

## Timeout

| 대상 | 설정 |
| --- | --- |
| Overall evaluation deadline | 기본 5분. Request validation 이후 평가 작업에 적용 |
| Individual Jev request | 60초/call. 여러 순차 호출이므로 전체 60초 응답 보장은 아님 |
| Server read-header / read | 5초 / 15초 |
| Server write | 설정한 evaluation deadline + 30초 |
| Idle / shutdown grace | 60초 / 10초 |

전체 evaluation deadline은 다음처럼 변경할 수 있다. 값은 0보다 크고 1시간 이하여야 한다.

```sh
./bin/job-fit-server --profile ./data/profile.json --evaluation-timeout 2m
```

Client timeout은 caller가 정한다. Client disconnect/취소는 Jev request까지 전달되지만 remote 처리/과금 중단까지 보장하지 않는다. SIGINT/SIGTERM으로 서버를 종료하면 진행 중 요청도 취소된다.

## Security Notes

기본 bind는 `127.0.0.1:8080`이고 `--listen 127.0.0.1:18080`처럼 numeric loopback 주소만 허용한다. Hostname `localhost`는 현재 server flag에서 허용하지 않는다. Authentication/authorization은 없으며 public deployment를 전제로 하지 않는다.

Jev key는 server-side environment variable이고 canonical profile은 server-side resource다. Client request/response에 key를 노출하지 않으며 로그에는 requestId/route/status/error category/elapsed 등만 기록한다. Profile/JD는 평가 과정에서 Jev로 전송되고 response에는 연결된 candidate evidence가 포함된다. Server-side profile이 local-only inference를 의미하지 않는다.

[API Specification으로 돌아가기](specification.md) · [설계 계약](../api-contract.md)

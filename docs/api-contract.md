# HTTP API Contract — v1

Turn 013에서 정한 구현 대상 contract다. **Turn 014에서 이 계약의 최소 HTTP service를 구현했다.** [Evaluation Contract](evaluation-contract.md), [Match Semantics](match-semantics.md)의 accepted v0 evaluator를 재사용한다. API v1은 evaluation baseline v0과 별개의 버전이다.

## 공통 개발 문서 참조 / 응답 템플릿

사용자가 대부분의 개발에 공통으로 적용하는 **API Response & Error Standard (v1)**（repository 외부 참조 문서, 원본 미포함）를 읽고 이 계약에 적용했다. 요청에 적힌 `api-responde.md`는 존재하지 않았으며 실제 파일명은 `api-response.md`다. 공통 원본은 수정하지 않았다. 이 절과 아래 JSON 예시는 job-fit용 응답 템플릿이며 별도 복제 표준을 만들지 않는다. 외부 원본 없이도 이 문서의 요약과 템플릿으로 job-fit 계약을 확인할 수 있다.

성공 envelope는 `status`, 숫자 `code`, `message`, `data`, `metadata`를 갖는다. 실패 envelope는 `status`, 숫자 `code`, 문자열 `error`, `message`, `metadata`를 갖고 `data`는 생략한다. `code`는 항상 실제 HTTP status와 같으며 2xx error는 금지한다. Non-i18n MVP에서는 plain message를 사용하고 caller는 message 대신 error로 분기한다.

공통 `metadata.timestamp`는 응답 생성 시점 UTC RFC 3339이고, `metadata.requestId`는 서버가 요청마다 생성한 UUID v4다. Incoming ID는 채택하지 않으며 client가 지정할 field도 없다. 공통 문서의 오류 로그 correlation 지침을 위해 requestId를 채택한다. 이는 persistence ID가 아니며 조회 endpoint를 만들지 않는다. 서비스가 처리하는 성공/오류 응답과 sanitized 로그에 동일 ID를 사용한다. TraceId는 도입하지 않는다. Evaluation 완료 시각 및 profile/engine 정보는 `data.metadata`에 별도로 둔다.

## Endpoint / execution

`POST /v1/job-fit/evaluate` — 하나의 이미 확보된 JD를 server-side canonical profile과 동기적으로 평가한다. 성공은 `200 OK`, 완전한 JSON result 하나다. 실패 시 error JSON만 반환하며 partial result나 임시 Overall을 반환하지 않는다.

Request/response media type은 `application/json`이다. Request의 `charset=utf-8` parameter는 허용한다. Request body는 UTF-8 JSON object 하나이며 압축 body는 지원하지 않는다. 성공·오류 응답은 `Cache-Control: no-store`를 사용한다. 저장, resource ID, 조회 endpoint, batch, async job, queue, polling, webhook은 없다. 자동 재시도는 하지 않으며 재요청은 새 evaluation이다. 동일 입력의 동일 label이나 idempotency를 보장하지 않는다.

## Request

```json
{
  "job": {
    "company": "Example Corp",
    "title": "Backend Engineer",
    "sourceUrl": "https://example.com/jobs/123",
    "description": "Backend Engineer\nRequired: Backend API design and development"
  }
}
```

| Field | 필수 | Contract |
| --- | --- | --- |
| `job` | 예 | Object. 아래 네 field만 허용 |
| `job.description` | 예 | 평가에 쓰는 유일한 JD text. Trim 후 비어 있지 않은 string |
| `job.company` | 아니오 | 식별용 metadata, UTF-8 256 bytes 이하 |
| `job.title` | 아니오 | 식별용 metadata, UTF-8 512 bytes 이하 |
| `job.sourceUrl` | 아니오 | Provenance metadata. UTF-8 2,048 bytes 이하의 absolute HTTP/HTTPS URL, host 필수, userinfo 금지. Fetch/검증 요청을 보내지 않음 |

Company/title/sourceUrl은 evaluation state에 보내거나 description 앞에 붙이지 않는다. **Role/title/company context가 평가에 필요하면 caller가 description 안에 포함해야 한다.** Metadata 충돌을 근거로 JD를 수정하거나 어느 것이 사실인지 조사하지 않는다. Metadata만 변경하면 evaluator에 전달하는 JD는 같다. 결과 label의 동일성까지 보장하는 것은 아니다.

Optional field는 생략 가능하지만 제공 시 non-null string이며 trim 후 비어 있으면 거부한다. Metadata는 앞뒤 공백을 제거하여 반환한다. 모든 string의 NUL과 잘못된 UTF-8/Unicode encoding을 거부한다. Unknown field, duplicate object key, null, 잘못된 타입, trailing JSON, 빈 body를 거부한다. 따라서 `candidateProfile`, `profileId`, `profileVersion` 같은 요청별 profile override도 허용하지 않는다. Endpoint는 query parameter를 받지 않는다.

### Input preparation / limits

- JSON body 최대 **128 KiB**. JSON escape expansion과 아래 field limits를 수용하기 위한 transport 제한이다.
- Decoded description은 trim 전 UTF-8 **16 KiB** 이하. 기존 CLI와 같이 전체 앞뒤 공백을 제거한 값을 `Job.Text`로 사용한다.
- `\n`으로 나눈 뒤 각 줄을 trim하고 빈 줄을 제외하여 `Job.Lines`를 만든다. 1–64개 non-empty line, 한 줄에 한 항목을 caller가 준비한다. HTML/Markdown parser, 자동 요약, 분해 또는 truncation은 하지 않는다.
- 기존 profile 1–32 evidence 및 추출 requirement × evidence **512개** 제한을 유지한다. Pair 제한은 추출 이후에 확인되므로 일부 Jev 호출 후 실패할 수 있다.
- 기존 Jev request payload 192 KiB 제한도 유지한다. 외부 단계별 state 크기에 따라 위 입력 제한을 만족해도 평가가 불가능할 수 있다.

이는 transport/resource safeguards이며 Match threshold가 아니다. Metadata를 줄 목록에 추가하지 않으므로 기존 `L001` 등의 identity 생성 방식을 보존한다.

## Canonical profile

Server 시작 시 운영자가 지정한 로컬 canonical profile 파일 하나를 읽고 검증한 **동일 snapshot**을 service에 주입한다. HTTP caller는 profile을 선택하거나 바꾸지 못한다. Hot reload는 없고 profile 갱신은 service 재시작으로 반영한다. Profile 내용/schema는 이번 Turn에서 변경하지 않는다.

`metadata.profileVersion`은 **`sha256:` + 로딩한 profile 파일 원본 bytes의 SHA-256 소문자 64자리 hex**다. Hash와 parsing은 동일하게 한 번 읽은 bytes에서 수행한다. 별도 version file, DB, 수동 version 관리가 필요 없다. 공백·순서만 바뀌어도 version이 달라질 수 있는 content fingerprint이며 semantic version, profile ID, 비밀값, 재현성 보장 또는 조회 가능한 resource가 아니다. Version 값으로 profile을 요청할 수 없다.

Profile read/validation 또는 credential 설정 실패는 startup 실패로 처리한다. 요청 처리 시 configuration failure가 발생하는 경우의 wire mapping도 아래에 정의한다. 성공 응답에는 실제 evaluation에 사용한 snapshot의 version만 기록한다.

## Success response

HTTP-specific DTO를 사용한다. Domain result의 의미를 보존하는 작은 명시적 mapping만 두고 새 evaluation logic은 만들지 않는다. 아래는 형식 설명용 가상 예시이며 실제 sample 결과가 아니다.

```json
{
  "status": "success",
  "code": 200,
  "message": "ok",
  "data": {
    "job": {
      "company": "Example Corp",
      "title": "Backend Engineer",
      "sourceUrl": "https://example.com/jobs/123"
    },
    "overallMatch": "Strong",
    "overallReasoning": "Core capabilities have direct documented support.",
    "requirements": [
      {
        "id": "L002",
        "requirement": "Required: Backend API design and development",
        "sourceType": "required",
        "match": "Strong",
        "supportingEvidence": [
          {
            "id": "E001",
            "text": "REST API design and implementation",
            "relation": "direct"
          }
        ],
        "missingEvidence": [],
        "unknowns": [],
        "reasoning": "Direct documented evidence supports the requirement."
      }
    ],
    "decisionTrace": {
      "source": "jev-post-hoc-attribution; not a causal record of the original decision",
      "supporting": [
        "L002"
      ],
      "limiting": [],
      "nonDecisive": []
    },
    "explanationSource": "application-rendered rationale from provider choices or absence of linked evidence; evidence copied from profile",
    "metadata": {
      "profileVersion": "sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
      "evaluatedAt": "2026-09-21T09:00:00Z",
      "engine": "jev",
      "engineModel": "jev-1.13.0"
    }
  },
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

### Public fields / mapping

아래 표의 경로는 모두 success envelope의 `data` 기준이다. 공통 envelope metadata와 evaluation metadata를 혼동하지 않는다.

| Public field | Type / source / meaning |
| --- | --- |
| `job` | Object containing supplied normalized company/title/sourceUrl only. Always present, `{}` if none. Description is not echoed; caller retains original input |
| `overallMatch` | `Strong`, `Partial`, `Weak`, `Unknown`, case-sensitive; unchanged domain value |
| `overallReasoning` | Domain string unchanged; application-rendered explanation, not model chain of thought |
| `requirements` | Non-empty array in source order; every retained requirement included |
| `requirements[].id` | Domain requirement ID unchanged. Current format `L` + at least three decimal digits, e.g. `L007`. Unique within result, gaps allowed after ignored lines. Not a global/persistent ID; editing/reordering JD can change it |
| `requirements[].requirement` | Domain text, normalized original JD line |
| `requirements[].sourceType` | Domain `category` renamed: `required`, `preferred`, `responsibility`, `stack`, `context`. `ignore` does not appear |
| `requirements[].match` | Same four Match values as Overall |
| `requirements[].supportingEvidence` | Array of `{id, text, relation}`; ID/text copied from profile, relation is `direct`, `transferable` or `limited`. An evidence link is a judgment, not a new candidate fact |
| `requirements[].missingEvidence`, `unknowns` | Arrays of domain strings unchanged. Missing evidence does not establish no experience |
| `requirements[].reasoning` | Domain string unchanged; no new explanation generation |
| `decisionTrace.source` | Preserve domain post-hoc attribution disclosure |
| `decisionTrace.supporting`, `limiting`, `nonDecisive` | Arrays of requirement ID strings, mapped from domain trace entries' `requirementId`; full Match/source/evidence resolved through `requirements` |
| `explanationSource` | Preserve domain explanation provenance disclosure |
| `metadata.profileVersion` | Profile snapshot fingerprint defined above |
| `metadata.evaluatedAt` | Successful evaluation/validation completion time, UTC RFC 3339 string; not request arrival time |
| `metadata.engine` | `jev` for v1 MVP |
| `metadata.engineModel` | Configured model identifier actually used for this evaluation, currently `jev-1.13.0`. Not inferred from result or hard-coded by HTTP mapper; not a promise of deterministic behavior |

All fields above except optional fields inside `job` are required on success. Arrays are `[]`, never `null`; JSON member order is insignificant. Human explanation strings (`reasoning`, `overallReasoning`, missing/unknown text, disclosure text) are for display, not client branching or exact-text matching. Provider model strings are descriptive metadata, not fixed enums.

Each requirement occurs **exactly once** across the three trace lists, preserving source order within each list. Missing, duplicate or nonexistent references are errors. Strong Overall must have empty `limiting`. Unknown or evidence-free requirements cannot be Supporting. Mapping must not silently drop an entry, reclassify influence, change a Match, or repair an invalid result. Remaining semantic validation stays with the accepted evaluator; API adds no stricter Match/influence policy.

Trace is **structured post-hoc attribution**, not a causal record of the original Overall decision. `Overall Match` is a summary signal, not a hiring probability, score or recommendation. Unknown does not mean no experience.

Domain `summary` is omitted because clients can group requirements by Match. Repeated text/Match/source fields in domain trace entries are omitted because IDs resolve them losslessly. Provider questions/state/choice IDs, raw responses, confidence/probability and full profile are not exposed. Linked candidate evidence **is** included and must be treated as personal data by callers. No new duplicate evidence model is introduced.

Request metadata is echoed under `data.job`; evaluation metadata is under `data.metadata`. Common response metadata contains timestamp/requestId. No evaluationId, persistenceId, duration metric or duplicated sourceUrl is added. This minimal response does not identify a complete evaluator build/policy revision; retain deployment information externally if required for comparisons.

## Errors

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

HTTP status와 숫자 `code`는 같아야 한다. `error`는 아래 문자열 코드이며 같은 조건에는 같은 코드를 반환한다. Message는 안전한 human summary이며 submitted content를 echo하지 않는다. 실패에는 `data`나 evaluation metadata를 넣지 않지만 공통 metadata는 항상 포함한다. 공통 코드 목록을 재사용하고 필요한 413/415/500/502/504 코드는 naming convention에 맞춰 확장했다.

| HTTP status | `error` | Conditions |
| --- | --- | --- |
| 400 | `REQUEST_MALFORMED_JSON` | Empty body, JSON syntax/encoding error, trailing JSON |
| 400 | `REQUEST_INVALID` | Object shape/type/null/unknown or duplicate key, unsupported query parameters |
| 422 | `REQUEST_VALIDATION_FAILED` | Missing/blank description, blank optional metadata, invalid sourceUrl, NUL; no in-scope requirements or pair/generated payload limit |
| 413 | `REQUEST_INPUT_TOO_LARGE` | Body/description/metadata byte limit or 64-line input limit |
| 415 | `REQUEST_UNSUPPORTED_MEDIA_TYPE` | Missing/unsupported Content-Type or Content-Encoding |
| 500 | `INTERNAL_CONFIGURATION_ERROR` | Profile unavailable/invalid, missing key, known upstream 401/403 access rejection |
| 500 | `INTERNAL_ERROR` | Unexpected application or DTO mapping/programming failure |
| 502 | `INTERNAL_EVALUATION_FAILED` | Jev network failure, other upstream non-success, invalid/oversized response, evaluator/trace consistency failure |
| 503 | `INTERNAL_DEPENDENCY_UNAVAILABLE` | Explicit upstream busy/rate limited (currently 429/529) |
| 504 | `INTERNAL_EVALUATION_TIMEOUT` | Total evaluation deadline/provider timeout with writable caller connection |

Validation error template:

```json
{
  "status": "error",
  "code": 422,
  "error": "REQUEST_VALIDATION_FAILED",
  "message": "Request validation failed.",
  "errors": [
    { "field": "job.description", "reason": "required", "message": "Description is required." }
  ],
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

422에만 `errors`를 포함한다. Field는 request 경로, reason은 `required`, `blank`, `invalid_format`, `invalid_character`, `no_capability_requirements`, `evaluation_capacity_exceeded` 중 하나다. 마지막 두 오류는 `job.description`에 연결하며 이미 일부 외부 호출이 수행됐을 수 있다. 검출된 오류 하나만 반환해도 되며 전체 validation 오류 수집을 보장하지 않는다. 문법 → 구조/타입 → 의미 검증 순서이며 byte limit은 읽는 중 먼저 적용할 수 있다. Error mapping은 enum/const와 하나의 mapper에서 관리한다.

Unknown routes return `404 RESOURCE_NOT_FOUND`; unsupported methods on the defined route return `405 REQUEST_METHOD_NOT_ALLOWED` with `Allow: POST`, using the same JSON error envelope when handled by the service. Malformed HTTP rejected before the handler or broken connections cannot guarantee this envelope. No additional endpoint is introduced.

Upstream status is not blindly forwarded: upstream 400/422 may indicate adapter incompatibility, so it is `502 INTERNAL_EVALUATION_FAILED`, not caller `400`; known local input limits use 413/422 above. Configuration errors are never described as caller authentication failures. Upstream 5xx is 502 unless explicitly recognized as busy above. No automatic retries or Retry-After promise; a caller may have incurred work/cost even on an error. No success result is persisted for recovery.

Never serialize raw Go/upstream error strings, upstream body, Authorization headers, API key, local profile path or internal request state. Error logs include numeric code, error, requestId and sanitized cause context; do not log raw JD/profile/secrets by default. Internal classification must use typed/sentinel errors or equivalent explicit categories, **not string matching**. Turn 014 implements categorized transport/application errors shared by CLI and HTTP.

## Timeout / cancellation

Application operation receives `context.Context` as its first argument. HTTP request context → service deadline child context → evaluator stages → Jev client request. CLI supplies its signal-aware context to the same operation. Cancel/deadline checks stop further stages; in-flight upstream calls receive cancellation. Local cancellation cannot guarantee that remote Jev processing/cost has stopped.

Use a finite configurable total evaluation deadline in addition to provider-call timeout. Current client has a 60-second timeout per call and evaluation can make several sequential calls; this is **not** a 60-second total API latency promise. Turn 014 defaults: total evaluation deadline 5 minutes (configurable, positive and at most 1 hour), read-header 5 seconds, read 15 seconds, write evaluation deadline + 30 seconds, idle 60 seconds, shutdown grace 10 seconds. Avoid a response write timeout that expires before the evaluation deadline/error can be sent.

Caller disconnect/cancellation stops work and ordinarily produces no deliverable HTTP response; do not promise a 499 response. If service/provider timeout occurs while connected, return 504. Preserve context cancellation/deadline identity and provider timeout categories through error wrapping. Turn 014 distinguishes provider timeout from network failure using typed errors without changing evaluation policy.

## Shared application boundary

```text
CLI files ──→ local profile snapshot + normalized JD ──┐
                                                     ↓
HTTP DTO ──→ server profile snapshot + normalized JD → Application evaluation operation
                                                     ↓
                                               evaluator.Evaluate
                                                     ↓
                                                Engine / Jev
                                                     ↓
                                       validated Domain Result + metadata
                                          ↙                      ↘
                                CLI JSON presenter           HTTP response DTO
```

Conceptual operation: `Evaluate(ctx, input) → output or categorized error`.

- Service dependencies: Engine plus immutable validated profile snapshot/version and engine identity; clock for completion metadata. No HTTP request/response types in service/core.
- Input: normalized JD (`Text`, `Lines`). Shared preparation/validation must give CLI text and HTTP description the same normalization and limits. Company/title/sourceUrl remain transport metadata; not evaluation inputs.
- Output: existing validated `EvaluationResult` plus profileVersion/evaluatedAt/engine/engineModel. HTTP mapper adds request metadata; CLI may keep existing domain JSON presentation without metadata.
- CLI can construct its service with the explicitly selected `--profile` file snapshot. REST uses its server-configured snapshot. Identical profile/JD run the same evaluation stages; neither interface adds its own evaluation policy.
- Errors: explicit invalid-input, unsupported-input, configuration, unavailable, timeout, canceled, evaluation-failure and internal categories. HTTP maps categories to the table; CLI retains stderr/nonzero behavior. File loading errors are handled by CLI or server bootstrap respectively.

Current `evaluator.Evaluate(ctx, engine, profile, job)` already sequences Extract → Assess → Overall → Trace → validation. Reuse this operation rather than create a second pipeline. Current Engine interface is the provider boundary, not an HTTP DTO. A thin service facade may own snapshot/metadata/error propagation; no generic repository framework or multiple mapping tiers are needed. Future MCP-over-REST reuses this same boundary indirectly; no MCP-specific evaluation logic.

## Versioning / security / non-goals

`/v1` versions the public request/response/error contract, independently from profile fingerprint, engine model and evaluation baseline. Breaking field/type/enum/meaning changes require a new major path. Additive optional response fields are allowed; clients ignore unknown response fields. Requests remain strict. Textual explanation changes are not a versioned machine interface. No complex deprecation policy is introduced.

Personal/internal use only. Authentication, authorization and public internet exposure are not implemented or implied by this contract. Initial service implementation should default to loopback; changing network exposure requires a separate access-control review. Jev credentials exist only server-side (CLI process for local use), never in API input/output. Profile/JD are still sent to Jev as part of evaluation; server-side profile storage does not mean local-only inference.

No company/Web search, crawling, URL fetching, company-type/SES/派遣 classification, Notion/OpenClaw integration, profile CRUD, multiple profiles/users, auth system, DB/persistence, batch, async, UI, percentage/weighting/threshold, hiring probability or application recommendation. Accepted evaluation semantics and unresolved Composite Public / Energy Context ambiguity remain unchanged.

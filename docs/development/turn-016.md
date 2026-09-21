# Turn 016 — API Documentation

- 날짜: 2026-09-21
- 범위: 현재 구현된 REST API의 사용자용 문서 작성. Code/API behavior 변경·Jev 호출 없음.
- 시작 HEAD: `c585cec` (`feat: expose job-fit evaluator over HTTP`). 기존 README/code/contract 수정 및 untracked server/application/httpapi/Turn 015 문서 등이 존재했다. 이 상태를 보존했다.

## Completed

Frozen portfolio MVP의 유일한 endpoint를 실제 caller가 사용할 수 있도록 overview와 endpoint 상세 문서 두 개를 작성했다. Existing API Contract는 설계 이유를 설명하는 문서로 유지한다. README에는 문서 탐색 링크만 추가했다.

## Changed Files

- [API Specification](../api/specification.md): 신규 overview/빠른 참조.
- [Job Fit Evaluation API](../api/job-fit.md): 신규 endpoint 상세 사용 문서.
- 본 Turn 기록.
- [README](../../README.md): API Specification/endpoint/contract/semantics 링크와 Turn 016 링크 추가.

기존 contract 문서는 수정하지 않았다.

## Reference / Source of Truth

실제 MyLib 문서를 읽어 overview와 domain endpoint 문서의 분리, endpoint table, Method/Path, Headers/Body/field 설명, success/error 예시 구조를 참고했다.

- **MyLib API Specification**（repository 외부 참조 문서, 원본 미포함）
- **MyLib Wishlist API**（repository 외부 참조 문서, 원본 미포함）
- **공통 API Response & Error Standard (v1)**（repository 외부 참조 문서, 원본 미포함）

참조 문서는 read-only로 사용했다. MyLib의 JWT/client ID, health, Swagger 또는 다수 domain 구조를 job-fit에 가져오지 않았다. 외부 reference의 원본은 repository에 포함하지 않는다. 사용자용 문서 간 링크는 repository-relative로 제공한다.

내용의 우선순위는 current code/tests → Turn 014 actual response → API Contract → Evaluation Contract → Match Semantics다. 구체적으로 HTTP request/response/handler/error mapper와 handler tests, server flags/lifecycle, application metadata, evaluator input loader, Jev timeout/payload/pair limit을 확인했다.

## API Specification

Base URL `http://127.0.0.1:8080`, path version v1, 단일 `POST /v1/job-fit/evaluate`, synchronous/personal/loopback 전제를 정리했다. Common success/error envelope와 field 의미, 실제 error code 전체 목록, Match level 요약 및 상세 문서 링크를 제공했다.

Envelope 예시의 `data: {}`는 공통 구조를 설명하는 placeholder임을 명시했다. 실제 success의 필수 field가 생략 가능하다고 읽히지 않도록 상세 응답으로 연결한다.

## Endpoint Documentation

Headers, request body, field별 required/type/limit/validation, canonical profile 동작, 실제 success response, response field 의미, Match 요약, 주요 오류 예시, inline/file curl, 최소 server requirements, limits, timeout, security를 문서화했다.

Company/title/sourceUrl은 metadata이며 evaluation에 자동 결합하지 않는다. Profile은 startup snapshot이고 fingerprint는 실제 파일 bytes의 SHA-256이다. Decision Trace는 post-hoc attribution이며 causal explanation이 아니라는 한계를 명시했다.

## Request / Response Verification

Request example은 parser가 허용하는 정확한 job/description/company/title/sourceUrl 구조를 사용한다. Success example은 [Turn 014 Prop Tech plus artifact](turn-014-results/proptech-plus.json)의 15개 requirement 전체를 수정/축약 없이 복사했다. Overall Partial, sourceType, evidence id/text/relation, missingEvidence, unknowns, trace IDs, evaluation/common metadata를 그대로 유지한다.

상단 Example Corp request와 저장된 Prop Tech response가 같은 실행인 것처럼 보이지 않도록 다른 예시임을 명시하고, 실제 Prop Tech description으로 사용한 sample을 연결했다. 신규 실제 호출은 수행하지 않았다. Curl 예제는 JSON/인자 형식을 정적으로 확인하며 재실행 결과나 고정 label을 주장하지 않는다.

## Error Documentation

400/413/415/422/500/502/503/504 각각 실제 code/message를 사용하는 전체 envelope 예시를 제공했다. 404/405 route/method 오류도 실제 code와 Allow header로 설명했다. 422는 누락 description의 errors 배열과 blank/URL/NUL/no-requirements/capacity field/reason/message 표를 제공한다.

Errors 예시의 timestamp/requestId는 설명용 값이며 이번 Turn에서 실패 요청을 실제로 발생시켰다는 뜻이 아니다. Code는 HTTP status와 같고 실패에 data는 없다. Upstream 400/422→502, 401/403→configuration 500, 429/529→503 구분도 확인했다. Secret이나 raw upstream error를 예제에 넣지 않았다.

## Limits / Timeout

Body 128 KiB, decoded description 16 KiB/64 non-empty lines, metadata 256/512/2,048 bytes, profile file 16 KiB/evidence 1–32, pairs 512, generated Jev payload 192 KiB를 코드와 대조했다. Byte 제한은 문자 수와 다르며 trim 전 측정이라는 점을 명시했다.

Overall evaluation 기본 5분, Jev request별 60초, 실제 존재하는 `--evaluation-timeout`(positive~1h)과 server read/write/idle/shutdown 값을 기록했다. Deadline은 검증 이후 evaluation에 적용하며 caller disconnect는 전달 가능한 오류 JSON을 보장하지 않는다.

## Security Notes

Loopback-only, 기본 127.0.0.1:8080, auth 없음, public exposure 전제 없음, server-side Jev env key/profile을 명시했다. `--listen`은 numeric loopback만 허용하므로 localhost hostname을 설정값 예제로 만들지 않았다. Response는 secret을 노출하지 않지만 linked candidate evidence를 포함하고 evaluation 시 profile/JD가 Jev로 전송된다는 기존 동작도 설명했다.

## Documentation Links

README → Specification → Endpoint 상세 문서로 이동할 수 있고 각 문서에서 API Contract, Evaluation Contract, Match Semantics 및 README로 돌아갈 수 있다. 기존 contract의 링크나 설계 결정을 재작성하지 않았다.

## Problems / Findings

검토한 endpoint/DTO/error/limit/timeout 항목에서 명시적 동작 계약을 바꿔야 할 mismatch는 발견하지 않았다. 다음은 caller가 혼동하기 쉬운 구현 세부사항이므로 문서에 구체적으로 표시했다.

- 413은 body만이 아니라 field byte/line limit에도 사용한다.
- Job object 부재/null은 400, job 안의 description 누락/blank는 422다.
- Content-Encoding은 어떤 non-empty 값도 거부하며 `identity`도 예외가 아니다. Content-Type은 JSON과 optional UTF-8 charset만 허용한다.
- Provider credential 오류는 client auth 401/403이 아니며 startup 실패는 HTTP response로 받을 수 없다.
- No-capability/pair/payload capacity 오류는 이미 일부 외부 호출 후 422가 될 수 있다.
- CLI result와 HTTP DTO는 동일 JSON shape가 아니다. HTTP trace는 ID 목록이고 summary는 없다.

기존 evaluation ambiguity와 모델 variability를 문서 작업으로 해결했다고 주장하지 않는다.

## Validation

JSON 예시 14개의 parsing, response envelope/status/errors 구조, success example과 실제 artifact의 완전한 동등성, Trace 참조, 코드의 error 문자열/메시지, curl 인자/JSON 정적 확인, 링크/anchor/code fence 및 `git diff --check`가 통과했다. Constants/flags와 문서의 제한·설정값도 소스로 대조했다. 기존 파일 hash 비교로 README 외 기존 파일의 변경이 없고 신규 파일이 위 문서 3개뿐임을 확인했다. Code/behavior 변경이 없으므로 Go tests/build나 API 호출은 반복하지 않았다. Commit하지 않았다.

## Deferred

OpenAPI/Swagger, architecture/project documentation 및 portfolio polish는 이번에 수행하지 않았다. Turn 015에 따라 OpenAPI/Swagger 도입은 별도 scope 결정이 필요하며, OpenClaw/MCP/Notion/DB/auth/UI/crawler 등 product 기능도 추가하지 않았다.

## Suggested Next Step

다음 Turn에서 현재 구현을 기준으로 **architecture / project documentation**을 작성한다. CLI/REST/application/evaluator/Jev 관계와 dependency boundary, 책임 및 한계를 정리하며 새 기능을 구현하지 않는다.

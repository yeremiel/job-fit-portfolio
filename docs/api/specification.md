# API Specification

현재 REST API의 빠른 참조 문서다. 호출 예제와 field 설명은 [Job Fit Evaluation API](job-fit.md), 설계 이유는 [API Contract](../api-contract.md)를 참조한다.

---

## Overview

Server-side canonical candidate profile과 제공된 Job Description을 비교하여 structured Job Fit evaluation을 반환한다. Caller가 JD를 이미 확보한 상태를 전제로 하며 회사 검색, URL fetching, 지원 추천은 수행하지 않는다.

**Job Fit ≠ Hiring Probability. Job Fit ≠ Application Recommendation.** 최종 판단은 사람이 한다.

| 항목 | 값 |
| --- | --- |
| Base URL | `http://127.0.0.1:8080` |
| API version | `v1` (path에 포함) |
| Content-Type / encoding | `application/json` / UTF-8 |
| 실행 방식 | Synchronous, JD 한 건씩 |
| 사용 범위 | Loopback-only personal/internal MVP. 인증 없음, public deployment 전제 없음 |

## Endpoint List

| API | Method | Endpoint | Description |
| --- | --- | --- | --- |
| [Job Fit 평가](job-fit.md) | POST | `/v1/job-fit/evaluate` | JD와 canonical candidate profile의 capability/experience match 평가 |

현재 business endpoint는 하나다. Profile CRUD, history, batch, health, Swagger/OpenAPI endpoint는 없다.

---

## Common Response Envelope

사용자의 공통 개발 문서 **API Response & Error Standard (v1)**（repository 외부 참조 문서, 원본 미포함）를 참조하고 실제 HTTP implementation에 적용된 형태를 사용한다. 외부 원본은 repository에 포함하지 않으며 [API Contract](../api-contract.md)의 공통 표준 요약도 참조할 수 있다.

### Success Response

아래 `data: {}`는 공통 envelope를 보여주는 자리이며, 실제 평가 성공에서는 [상세 응답](job-fit.md#success-response)의 필수 field가 채워진다.

```json
{
  "status": "success",
  "code": 200,
  "message": "ok",
  "data": {},
  "metadata": {
    "timestamp": "2026-09-21T09:00:01Z",
    "requestId": "7f65f5a0-9ed1-4c3b-8afe-46d03b69a5c1"
  }
}
```

### Error Response

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

| Field | Meaning |
| --- | --- |
| `status` | `success` 또는 `error` |
| `code` | 실제 HTTP status와 같은 숫자. HTTP 200에 error를 반환하지 않음 |
| `message` | Human-readable plain message. 분기용 고정 문자열이 아님 |
| `data` | 성공에서만 존재하는 evaluation object |
| `error` | 실패에서만 존재하는 machine-readable code |
| `errors` | 422의 field별 `{field, reason, message}` 배열 |
| `metadata.timestamp` | 응답 생성 시각, UTC RFC 3339 |
| `metadata.requestId` | 서버 생성 UUID v4. 로그 correlation용이며 조회 가능한 evaluation ID가 아님 |

성공/오류 응답에 `Content-Type: application/json`, `Cache-Control: no-store`를 사용한다. RequestId는 caller가 지정하지 않는다. Evaluation metadata는 공통 metadata와 별도로 `data.metadata`에 있다. 실패 응답에는 data와 partial result가 없다. 연결 종료나 handler 이전에 거부된 malformed HTTP는 JSON 응답 전달을 보장하지 않는다.

---

## Match Levels

| Level | Meaning |
| --- | --- |
| Strong | 핵심 requirement를 직접 뒷받침하는 명확한 evidence가 있고 role/context 차이가 크지 않음 |
| Partial | 관련·transferable evidence가 있으나 일부 scope/depth/technology/domain/operational gap이 있음 |
| Weak | 관련 evidence는 있으나 핵심 requirement 또는 요구 수준과 gap이 큼 |
| Unknown | 평가할 evidence가 충분하지 않음. 경험 부재를 의미하지 않음 |

상세 정의는 [Match Semantics](../match-semantics.md)를 참조한다.

Overall은 summary signal이며 requirement Match의 단순 평균이 아니다. Decision Trace는 **post-hoc attribution**으로 원래 판단의 인과 기록이 아니다.

## Error Summary

| HTTP | Error code | Meaning |
| --- | --- | --- |
| 400 | `REQUEST_MALFORMED_JSON` | 빈 body, JSON 문법/encoding 오류, trailing JSON |
| 400 | `REQUEST_INVALID` | 잘못된 구조·타입·null·unknown/duplicate field 또는 query parameter |
| 404 | `RESOURCE_NOT_FOUND` | 등록되지 않은 path |
| 405 | `REQUEST_METHOD_NOT_ALLOWED` | 정의된 path의 non-POST 요청. `Allow: POST` |
| 413 | `REQUEST_INPUT_TOO_LARGE` | Body/description/metadata bytes 또는 JD 줄 수 제한 초과 |
| 415 | `REQUEST_UNSUPPORTED_MEDIA_TYPE` | 누락/미지원 Content-Type, charset/parameter 또는 Content-Encoding |
| 422 | `REQUEST_VALIDATION_FAILED` | 필수/blank/URL/NUL 검증, 평가 대상 없음 또는 평가 capacity 초과 |
| 500 | `INTERNAL_CONFIGURATION_ERROR` | 서버 profile/credential 설정 문제, upstream 401/403 |
| 500 | `INTERNAL_ERROR` | 예상하지 못한 내부 실패 |
| 502 | `INTERNAL_EVALUATION_FAILED` | Jev network/응답/평가/consistency 실패. 기타 upstream non-success 포함 |
| 503 | `INTERNAL_DEPENDENCY_UNAVAILABLE` | Jev busy/rate limit (upstream 429/529) |
| 504 | `INTERNAL_EVALUATION_TIMEOUT` | 전체 evaluation deadline 또는 Jev 호출 timeout |

Caller는 HTTP status와 `error`로 분기하고 `message`는 표시용으로 사용한다. Upstream status를 그대로 전달하지 않는다. 예를 들어 upstream 400/422는 502, upstream 401/403은 server configuration 500이다. Key/profile 오류로 startup 자체가 실패하면 HTTP listener가 없으므로 error JSON을 받을 수 없다.

## Limits / Timeout Quick Reference

Body 128 KiB, decoded description 16 KiB/64 non-empty lines, evidence pairs 512. 전체 평가 deadline 기본 5분, Jev 호출별 60초. 설정과 추가 제한은 [상세 문서](job-fit.md#request-limits)를 참조한다. 자동 retry/history는 없고 재요청은 새로운 evaluation이다.

## Detailed Documentation

- [Job Fit Evaluation API](job-fit.md): headers/body/fields, 실제 성공 응답, 오류 예시, curl, 서버 설정.
- [API Contract](../api-contract.md): 설계 선택과 public contract.
- [Evaluation Contract](../evaluation-contract.md): 평가 책임과 품질 기준.
- [Match Semantics](../match-semantics.md): Match 의미와 제한.
- [README](../../README.md): 프로젝트 소개와 실행 준비.

# Turn 019 — Public Portfolio Release Preparation

- 날짜: 2026-09-21
- 범위: 공개 정보 검토, 개인 절대경로 정리, repository hygiene 및 regression. Application behavior 변경 없음.
- 시작 상태: HEAD `027869c`; Turn 018의 README/API 링크 변경과 narrative/Turn 기록이 미커밋 상태였다. 기존 작업을 보존했다.

## Completed

현재 profile, 이를 인용한 문서와 저장된 평가 결과를 검토하고 **Option A — Keep Current Profile Public**을 선택했다. 공개용 가상 profile을 만들거나 원래 evidence를 변경하지 않았다. 문서의 machine-local 경로를 정리하고 민감정보·링크·파일 상태를 점검했다. Jev 호출 없이 요청한 regression을 모두 통과했다.

현재 파일과 Git history의 공개 범위는 구분했다. 현재 파일 정리는 완료했으나, 과거 commit의 개인 경로가 남아 있어 전체 history 공개 준비는 아래와 같이 별도 판단했다.

## Changed Files

이번 Turn의 변경:

- `README.md`: 공개 준비 점검 기록 링크 추가.
- `docs/api-contract.md`, `docs/api/specification.md`: 외부 공통 개발 문서의 개인 경로 제거. 표준을 참조했다는 사실과 repository 내 계약 요약은 유지.
- `docs/development/turn-000.md`: workspace를 repository root로 표현.
- `docs/development/turn-013.md`, `turn-014.md`, `turn-016.md`: 외부 reference의 제목·참조 사실을 유지하고 개인 절대경로 링크 제거.
- `docs/development/turn-018.md`: 경로 예시를 일반 표현으로 치환. 당시 privacy 판단과 미완료 상태는 역사적 기록으로 보존.
- 본 Turn 기록 생성.

기존 Turn 018 작업인 `docs/api/job-fit.md` 링크 수정 및 `docs/ai-assisted-development.md` 생성은 이번 Turn의 신규 작업이 아니다. Code/tests/profile/JD/평가 artifacts는 변경하지 않았다.

## Candidate Profile Review

| 분류 | 확인한 내용 | 판단 |
| --- | --- | --- |
| Publicly Acceptable | E1: 약 16년 backend, Java/Spring/API; E2: RDB/SQL; E3: 분석/설계/개선; E4: Lead/PM/PL 및 약 12명 팀; E5: 기술 stack; E6: ERP/회계/암호화/선박/물류/공공 도메인; E7: 기본 AWS | 일반적인 공개 이력서 수준의 전문 경력 요약. 유지 |
| Review Required | 실경력이라는 점, 경력·팀 규모·도메인의 조합 및 문서/artifact에 반복되는 원문 | 이번 요청의 공개 분류 기준으로 검토 완료. 특정 과거 고용주·고객·비공개 프로젝트·계약·내부 시스템명은 발견하지 못함 |
| Must Not Publish | 실제 연락처, 주소, credential, 계정 식별자, private repository 및 구체적 confidential 정보 | 현재 profile/파생 evidence에서 발견 없음. 과거 Git blob의 로컬 계정 경로는 별도 잔여 항목 |

저장된 evaluation JSON의 evidence `text` 353건은 모두 현재 profile의 7개 원문 중 하나였다. Match Semantics의 E1–E7과 API 예제도 같은 경력 범위를 사용한다. 여섯 회사명은 공개 채용공고의 평가 대상이며 candidate의 과거 고객/고용주를 나타내는 정보가 아니다. `personal-information database encryption`은 일반 업무 도메인 설명이며 실제 개인정보 데이터가 아니다.

이 검토는 repository 내용에 근거한 분류다. 경력의 사실성 또는 저장소에 없는 계약/NDA의 유무를 독립 검증했다는 뜻은 아니다. 소유자와 연결되는 실경력이라는 점도 README에서 계속 명시한다.

## Public Profile Strategy

**Option A — Keep Current Profile Public.** 사용자에게서 제공된 공개 허용 범주의 경력 요약에 해당하며 직접 식별정보나 특정 confidential detail을 찾지 못했다. Profile과 파생 결과를 함께 그대로 유지해 평가 이력의 일관성을 보존한다.

`data/profile.example.json`을 만들거나 private profile로 전환하지 않는다. Evidence가 그대로이므로 sample 재평가도 필요하지 않다. 이후 소유자가 특정 경력을 비공개로 지정하면 profile뿐 아니라 파생 문서·artifacts·Git history까지 함께 검토해야 한다.

## Local Path Cleanup

Repository 내부 workspace는 `.`로 표현했다. Common **API Response & Error Standard (v1)** 및 MyLib 문서는 repository 밖의 참조 자료이므로 제목과 원본 미포함 사실만 남겼다. 존재하지 않는 상대 링크를 만들거나 외부 원본을 복사하지 않았다. 사용자 공통 개발 문서를 참조해 API 표준을 적용했다는 기록은 유지된다.

현재 공개 대상 파일에서는 사용자 홈 절대경로가 더 이상 발견되지 않는다. 역사 문서의 당시 결정·검증 결과는 변경하지 않았다.

**Git history는 별개다.** Reachable historical blobs에는 API Contract/Specification, Turn 000/013/014/016/018의 이전 로컬 경로가 남아 있다. 현재 문서를 commit해도 과거 blob은 사라지지 않는다. History 재작성, force push, repository 공개는 수행하지 않았다.

## Sensitive Data Scan

현재 82개 tracked 파일과 Turn 018/019의 신규 문서를 대상으로 내용·파일명을 검사했다. 민감 값은 출력하거나 결과 artifact에 저장하지 않았다.

| 점검 | 결과 / false positive 구분 |
| --- | --- |
| API key/token/private key | 일반 credential 패턴 및 실제 로컬 Jev key와의 비출력 비교에서 실제 secret 발견 없음 |
| Bearer/Authorization/password | HTTP header 처리, redaction 설명, 명시적인 fake test credential. 실제 값 아님 |
| Email / phone / address | 실제 연락처 발견 없음. `example.com` URL user-info 검증 fixture와 이를 설명한 Turn 018만 email 패턴에 해당 |
| Private repository URL | GitHub/GitLab/Bitbucket 및 SSH repository URL 패턴 발견 없음. 공고 URL은 공개 source reference |
| 고객/회사 confidential 명칭 | Profile/파생 evidence에서 특정 명칭 발견 없음. 채용 대상 회사 및 외부 문서 제목과 구분 |
| Personal absolute path | 현재 파일 정리 완료. 과거 blob에는 잔존 |
| `.env` | 로컬 ignored 파일. Reachable history에서 env 경로 추적 기록 없음 |
| Generated binary / temp | Tracked executable binary 또는 임시/backup 파일 발견 없음 |

추가로 reachable Git blobs에 대해 private-key/token 패턴과 현재 로컬 Jev key 일치를 확인했고 발견하지 못했다. 이는 알려진 패턴과 현재 key에 대한 검사이며 모든 형태의 과거 credential 부재를 보장하는 보안 인증은 아니다.

## Repository Hygiene

- `.gitignore`는 `.env`, `.env.*`(example 예외), `/bin/`, `/job-fit`, `.DS_Store`를 제외한다. 현재 파일 상태에 필요한 변경은 없다.
- 기존 ignored `.env`, root CLI binary, `bin/` binaries 및 `.DS_Store`는 공개 대상이 아니며 임의 삭제하지 않았다. 로컬 디렉터리 전체를 압축해 공개하는 것과 Git 파일 공개는 다르다.
- README의 profile/sample 경로는 실제 파일과 일치한다. Go, Jev key 및 네트워크/비용, `.env` 자동 로딩 없음이라는 prerequisite가 명시되어 있다. 새 clone에는 `.env`가 없으며 해당 source 예시는 로컬 파일을 준비한 경우에만 사용한다.
- 외부 공통 원본 없이도 repository의 API Contract/Specification으로 실행과 response 구조를 확인할 수 있다.
- Repository Markdown의 내부 파일 링크·anchors·fences, README JSON/쉘 예시를 정적으로 검증했다. 외부 채용 URL의 현재 availability는 재조회하지 않았다.
- Turn 시작 시점의 hash와 비교하여 code/tests/profile/JD/JSON artifacts가 그대로임을 확인했다.

## Final Regression

| Command | Result |
| --- | --- |
| `go test ./...` | PASS (Go test cache 사용) |
| `go test -race ./...` | PASS (Go test cache 사용) |
| `go vet ./...` | PASS |
| `go build -o /tmp/job-fit ./cmd/job-fit` | PASS |
| `go build -o /tmp/job-fit-server ./cmd/job-fit-server` | PASS |
| `git diff --check` | PASS |

외부 Jev API는 호출하지 않았다. Build output은 요청대로 `/tmp`에만 생성했다.

## Remaining Privacy Concerns

Profile의 public/private 분리는 현재 검토 결과 필요하지 않다. 남은 항목은 **전체 Git history를 공개할 경우 과거 개인 로컬 계정 경로가 노출된다는 점**이다.

공개 전 아래 중 한 방식으로 범위를 정리해야 한다.

1. 기존 비공개 history는 보존하고, 정리된 tracked/new documentation 파일만으로 새 공개 repository snapshot을 만든다. Ignored `.env`나 binary는 포함하지 않는다.
2. 전체 개발 Git history를 공개하려면 별도 승인된 작업에서 과거 blob의 개인 경로를 정리하고 history 검사를 반복한다. Commit ID 변경과 기존 remote/clone에 미치는 영향을 함께 다룬다.

이 조치는 evaluation 재실행이나 새 기능을 요구하지 않는다. 이번에는 어느 방식으로도 Git history 또는 remote를 변경하지 않았다.

## Final Status

**Additional Privacy Work Required** — 현재 공개 대상 파일과 profile 검토, 경로 정리 및 regression은 완료했다. 전체 Git history의 개인 경로는 남아 있으므로, 공개 전에 위 snapshot/history 정리 중 하나만 처리하면 된다. Profile 분리·evaluator 변경·Jev 재평가는 필요하지 않다.

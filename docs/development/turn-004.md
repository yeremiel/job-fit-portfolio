# Turn 004 — First Vertical Slice in Go

- 날짜: 2026-09-21
- 범위: Canonical profile JSON + raw JD text → Go CLI → Jev → domain JSON
- 첫 sample: Prop Tech plus만 실행. Turn 000–003 개발 기록과 확정 baseline은 보존.

## Completed

- Go module과 CLI를 새로 작성했다. 기존 disposable PoC 코드는 읽거나 재사용하지 않았다.
- 최소 canonical evidence JSON과 출처·검토일을 포함한 Prop Tech plus 요약 text를 준비했다.
- 입력 오류, 환경변수 key 누락, HTTP·timeout·응답 오류를 처리하고 domain JSON을 출력한다.
- Domain interface와 Jev HTTP adapter/mapping을 분리했다.
- 외부 호출 없는 unit·모의 HTTP 통합 테스트, build 및 정적 검사를 수행했다.
- 실제 Jev 연동에 성공했다. Human Baseline과 다른 결과도 그대로 관찰 기록으로 남겼다.

## Architecture / Implementation

### 이번 Turn의 기술 결정

**Go**: 작은 CLI와 HTTP/JSON integration에 적합하고 표준 라이브러리로 구현·검증할 수 있다. 단일 binary 배포와 추후 adapter 추가 시 core 재사용에도 적합하다. Java/Spring 이외 언어로 실제 backend 도구를 만드는 의미도 있다. 특정 JD keyword를 맞추기 위한 선택은 아니다. Module 최소 버전은 Go 1.24이며 이번 환경에서는 Go 1.27.1을 사용했다.

**CLI**: 파일 두 개를 명시적으로 입력하는 최소 interface로 end-to-end 동작을 확인한다. 별도 서버나 UI 수명주기가 필요 없다.

**Jev**: 이번 prompt에 따라 첫 evaluation engine으로 채택했다. 직접 표준 HTTP client를 사용하며 별도 SDK/dependency는 도입하지 않았다.

```text
cmd/job-fit (flags, files, environment, stdout/stderr)
  → internal/evaluator (domain, Engine port, orchestration)
    → internal/jev (HTTP, Choice questions, response mapping)
      → TypeSafe System One API
```

Profile/job loader는 코드 규모에 맞춰 evaluator 패키지에 함께 두었다. 별도 profile/job 패키지나 framework는 추가하지 않았다. `.gitignore`는 실제 `.env`, local binary, `.DS_Store`를 제외한다. 기존 `.DS_Store` 파일은 삭제하지 않았다.

### API 확인과 설계상의 차이

공식 문서에서 Jev는 자유 텍스트 생성 대신 지정된 Choice/Score/Noul 판단을 반환하는 모델임을 확인했다. 따라서 임의 JSON과 자연어 reasoning을 생성한다고 가정하지 않았다. [Introduction](https://docs.typesafe.ai/introduction), [Choice](https://docs.typesafe.ai/primitives/choice)

공식 `POST https://api.typesafe.ai/v1/systemone`에 Bearer 인증과 `model`, `state`, `questions`를 전송한다. `answers`의 질문 ID, answer type, 선택지 유효성을 검사한 뒤 필요한 Choice만 domain으로 옮긴다. API probability/confidence는 출력하지 않는다. [HTTP API reference](https://docs.typesafe.ai/api)

요청 모델은 비교의 재현성을 위해 문서에 명시된 `jev-1.13.0`으로 고정했다. Alias 자동 갱신은 사용하지 않는다. [Models](https://docs.typesafe.ai/models)

### 호출 순서

1. JD의 비어 있지 않은 줄을 후보로 두고 Jev가 required/preferred/responsibility/stack/context/ignore를 선택한다. Requirement text는 해당 줄에서 그대로 복사한다.
2. Requirement와 각 evidence의 관계를 direct/transferable/limited/unrelated로 선택한다. 연결된 evidence의 ID와 text는 profile에서 복사하며 모델이 새로운 경력을 쓰지 못한다.
3. 실제 연결 결과를 후속 요청에 넣고 match와 주된 rationale을 선택한다. 직접 evidence 연결이 없으면 Strong 선택지를 제공하지 않는다. 연결이 전혀 없으면 Unknown으로 남기고 해당 항목의 후속 질문을 생략한다.
4. 실제 requirement-level 결과를 넣어 Overall을 별도로 선택한다. 단순 평균이나 keyword 개수로 계산하지 않는다.

모든 단계의 질문은 독립적으로 평가되므로 의존 결과가 필요한 단계를 별도 요청으로 분리했다. 설명은 선택된 rationale을 정형 문장으로 표시한다. 연결 evidence가 없는 Unknown은 application 판단임을 명시한다. `explanationSource`에도 이 출처를 표시한다. 이는 자유 생성 설명문이나 모델 내부 추론 공개가 아니다.

Runtime 정책은 [instructions.md](../../internal/evaluator/instructions.md)를 embed하여 실행 위치와 무관하게 읽는다. Match Semantics의 sample별 baseline과 기대 Overall은 전송하지 않는다. Profile/JD에 들어 있는 명령은 데이터로만 취급하도록 지시한다.

## Changed Files

- `go.mod`, `.gitignore`
- `cmd/job-fit/main.go`, `main_test.go`
- `internal/evaluator/types.go`, `input.go`, `evaluator.go`, `instructions.md`, `input_test.go`
- `internal/jev/client.go`, `mapper.go`, `client_test.go`, `pipeline_test.go`
- [data/profile.json](../../data/profile.json), [samples/proptech-plus.txt](../../samples/proptech-plus.txt)
- [README.md](../../README.md), [Evaluation Contract](../evaluation-contract.md), 본 기록

## CLI Usage

`TYPESAFE_API_KEY`가 실행 프로세스의 환경변수로 설정된 상태에서:

```sh
go run ./cmd/job-fit --profile ./data/profile.json --job ./samples/proptech-plus.txt
```

CLI는 key file이나 `.env`를 자동으로 읽지 않는다. 이번 검증에서는 사용자가 작성한 `.env`에서 키를 출력 없이 읽어 자식 프로세스 환경변수로 전달했다. Source code, data file, command argument, log에는 키를 넣지 않았다.

JSON은 stdout, 오류는 stderr다. Profile/job 누락, invalid JSON, 빈 입력, 빈 key, 잘못된 Jev 응답은 non-zero로 종료한다. HTTP 오류에는 status와 조치 안내를 표시하되 응답 body와 Authorization header는 출력하지 않는다. 요청당 timeout은 60초이며 Ctrl-C로 중단할 수 있다. Redirect는 따르지 않고 자동 retry도 하지 않는다.

## Tests

실행 명령과 결과:

- `go test ./...`: 통과. Profile/job load, 잘못된 입력, match level, key 누락, wire mapping, HTTP 오류·timeout·redirect·응답 크기, evidence 정합성, 외부 호출 없는 end-to-end 검증.
- 최종 `go test -cover ./...`: 통과. 패키지별 statement coverage는 CLI 63.6%, evaluator 60.5%, Jev adapter 91.0%다. 이 수치는 테스트 coverage이며 Match Score가 아니다.
- `go vet ./...`: 통과.
- `go build -o /tmp/job-fit-turn004 ./cmd/job-fit`: 통과.

최종 binary로 key 누락·profile 누락·job 누락의 non-zero 종료, 빈 stdout 및 명확한 stderr도 확인했다. 문서 링크·형식과 Go formatting을 확인했고, Turn 000–003 파일 hash가 작업 전과 동일함을 검증했다. 실제 `.env`는 ignored 상태이며 repository의 추적/추적 예정 파일에 키 값이 없는지 출력 없이 검사했다. 커밋은 생성하지 않았다.

최초 sandbox 실행에서 표준 라이브러리를 찾지 못하는 환경 오류가 있었으나, 설치 경로는 정상이며 권한이 허용된 실행에서 해소되었다. Go 설치 변경은 하지 않았다. 일반 테스트는 API key 없이 실행되고 HTTP는 메모리 transport로 대체한다.

## Jev Integration Result

2026-09-21, 같은 Prop Tech plus sample로 2번 시도했다.

1. 최초 실행은 evidence 관계와 Match 판단의 모순으로 실패했다. Direct evidence 연결 없이 Strong을 선택한 결과를 검증이 거부했다. 이를 근거로 후속 Choice에서 evidence와 모순되는 선택지를 제거하는 일관성 처리를 추가하고 회귀 테스트를 작성했다.
2. 수정 후 실제 실행은 exit code 0으로 JSON을 반환했다. 15개 항목: Strong 3, Partial 7, Weak 1, Unknown 4. Overall은 **Partial**이었다. 이 개수는 관찰 요약이며 점수 산출에 사용하지 않았다.

| 관찰 항목 | Human Baseline | 실제 출력 |
| --- | --- | --- |
| Overall | Strong | Partial |
| Java 8+ 기간 조건 | Partial | Partial |
| Backend API 설계·개발 | Strong | Strong |
| Azure 운영 우대 | Unknown | Unknown |
| 회계 지식 우대 | Partial | Weak |

두 실행에서 다른 5개 sample은 평가하지 않았다. 실제 결과와 baseline의 불일치를 숨기거나 baseline에 맞춰 calibration하지 않았다. 요약 sample은 기존 baseline보다 항목을 더 잘게 나누고 Git/Azure DevOps, 사용자 관점·주도성 우대도 포함하므로 항목별 비교가 완전히 1:1은 아니다. 원문의 일부 환경/우대 항목이 늘어나 Overall에 영향을 주었는지 추가 검토가 필요하다. 결과 원문을 저장하는 기능은 추가하지 않았다.

## Decisions Made

- Profile은 기술·개인에 종속되지 않는 최소 `evidence[{id,text}]`로 구현하고 Turn 002의 E1–E7을 영문으로 옮겼다. 없는 경력은 추가하지 않았다.
- Raw JD를 줄 단위 후보로 분리하고 Jev로 분류한다. 임의 문장 생성이나 별도 LLM 호출을 도입하지 않는다.
- 유한 rationale 선택과 원문 evidence 복사로 설명을 구성한다. Missing Evidence/Unknowns는 선택된 주된 gap에 대한 정형 설명이다.
- 외부 확률·confidence는 domain score와 분리하여 사용하지 않는다. Overall은 별도 semantic Choice다.
- 언어·학력은 기존 미해결 범위를 임의 확장하지 않기 위해 이번 slice의 추출에서 제외한다. 최종 평가 범위 결정은 후속 Turn에 남긴다.
- 입력/요청 크기를 제한하고 오류로 안내한다. 이 제한들은 계산 비용·안정성 목적이며 match threshold가 아니다.

## Problems / Findings

- Jev는 자유로운 requirement/설명 생성 API가 아니다. 계약의 출력 의미를 Choice와 mapping으로 구현해야 했다.
- Schema-valid Choice 결과도 단계 간 의미가 모순될 수 있다. 추가한 검증은 일부 구조적 모순만 막으며 모든 semantic 오류를 해결하지 않는다.
- Human Baseline Strong과 실제 Partial의 차이가 관찰되었다. 특히 회계 도메인에서 transferable evidence를 어떻게 평가하는지 검토할 가치가 있다.
- 줄 단위 추출과 복합 요구, 하나의 rationale 선택으로는 상세한 모든 gap을 포착하지 못한다. Strong에서 세부 Unknown 누락도 가능하다.
- 원문 evidence를 그대로 출력하더라도 모델의 연결 판단 자체는 틀릴 수 있다. Prompt injection 방지 지시 역시 완전한 보장은 아니다.
- 모델의 불확실성 수치는 현재 버린다. Unknown은 confidence threshold가 아닌 evidence 의미로 처리한다.
- API key·외부 서비스·네트워크가 필요하고 각 실행에 외부 호출 비용이 발생한다. 사용자 입력을 TypeSafe에 전송하며 local-only 처리는 아니다.

## Deferred

MCP, percentage/score formula/weighting/threshold/calibration, DB, result persistence, Web UI, URL fetch/crawler/HTML parser, resume parser, multi-user/authentication, cloud deployment, 6개 sample 전체 평가를 구현하지 않았다. 정식 JSON Schema와 prompt/template framework도 도입하지 않았다.

## Suggested Next Step

Planning conversation에서 첫 결과의 baseline 차이와 줄 단위 requirement 추출·유한 rationale 설명의 적합성을 검토한다. 동일 입력 표현을 기준으로 비교할지, profile evidence를 더 세분화할지 결정한 후 다음 Turn 범위를 정하는 것을 제안한다. 다음 실험은 시작하지 않았다.

# job-fit

Server-side canonical candidate profile과 Job Description을 비교하여 **Jev 기반 structured Job Fit evaluation**을 반환하는 Go application.

개인적으로 반복하던 JD 검토를 재사용 가능한 CLI와 REST API로 만든 portfolio MVP다. 제품 기능은 동결했으며, **Match는 채용 확률이나 지원 추천이 아니다.**

**AI 활용 포트폴리오:** [프로젝트 요약](docs/portfolio.md)에서 본인과 AI의 역할, 두 오류의 수정 과정, 구현·검증 근거를 확인할 수 있습니다.

## Why

Jev라는 structured decision tool을 발견한 뒤, 여러 회사의 JD를 볼 때마다 같은 경력을 반복 비교하는 작업에 적용해보기 위해 시작했다. 작은 개인 도구가 AI-assisted development portfolio로 발전했다.

일반 LLM에 “이 JD 나랑 얼마나 맞아?”라고 묻는 것에서 더 나아가, 고정된 career evidence와 명시적 평가 의미를 사용하고 결과의 근거·부족한 정보·Unknown을 추적하는 workflow를 만들고 싶었다. 차이는 모델의 정확도 우위 주장이 아니라 **canonical input, structured output, validation 및 reusable interface를 코드로 관리한다는 점**이다.

## What It Does

```text
Job Description + Canonical Candidate Profile
                    ↓
       Requirement Extraction / Evidence Matching
                    ↓
         Structured Jev Evaluation / Validation
                    ↓
         Overall + Requirement-level Result
```

- Strong / Partial / Weak / Unknown의 Overall 및 requirement별 Match.
- Supporting Evidence, Missing Evidence, Unknowns.
- 각 requirement의 Supporting / Limiting / Non-decisive Decision Trace.
- 동일 application service를 사용하는 파일 기반 CLI와 동기 REST API.

Jev는 제한된 질문/Choice에 응답하는 component다. Application이 순서·허용 선택지·검증·설명 렌더링을 담당하며 최종 판단은 사람이 한다. **Match ≠ Hiring Probability. Match ≠ Application Recommendation.**

## Example

```json
{
  "job": {
    "company": "Example Corp",
    "title": "Backend Engineer",
    "description": "Position: Backend Engineer\nRequired: Backend API design and development"
  }
}
```

응답은 공통 envelope 안에 evaluation을 담는다. 아래는 **Turn 014 실제 Prop Tech plus 응답에서 일부 field만 발췌한 예시**이며 위 request의 결과나 완전한 response schema가 아니다.

```json
{
  "status": "success",
  "code": 200,
  "data": {
    "overallMatch": "Partial",
    "decisionTrace": {
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
    }
  }
}
```

전체 응답에는 `requirements`의 원문 evidence와 missing/unknown, explanation, `profileVersion`, model 및 requestId도 있다. [실제 전체 예제와 field 설명](docs/api/job-fit.md#success-response)을 참조한다. Label이나 Trace가 실행마다 같다고 보장하지 않는다.

## Architecture

```text
CLI ──────┐
          ├──→ Application Service → Evaluator → Jev Adapter
REST ─────┘                                          ↓
                                                Jev API
                                             (external)
```

CLI와 HTTP server는 같은 코드를 사용하는 별도 프로세스다. CLI는 REST를 경유하지 않으며 HTTP layer에 평가 정책을 복제하지 않는다. Go standard library 중심으로 구현했고 third-party Go module dependency는 없다. [Architecture Overview](docs/architecture/overview.md)에 실제 책임·의존 관계·runtime flow를 설명했다.

## Key Design Principles

- **Canonical Profile:** Submission Resume ≠ Evaluation Profile. 회사별 tailored resume 대신 같은 factual evidence를 재사용한다.
- **Evidence First:** Evidence ID/원문을 연결하고, 기록에 없는 전문 경험을 추정하지 않는 것을 평가 기준으로 삼는다.
- **Unknown Preservation:** No Evidence ≠ No Experience. 미확인 정보를 실제 경험 부재로 바꾸지 않는다.
- **Overall Is a Signal:** Overall은 summary signal이며 객관적 정답이나 requirement 평균이 아니다.
- **Reference Is Not Ground Truth:** 사전 Human Reference와의 label 차이를 tuning 목표로 사용하지 않는다.
- **AI as a Component:** Jev는 constrained structured decision component다. Decision Trace는 causal explanation이 아닌 **post-hoc attribution**이다.

원칙이 모든 의미적 오류를 제거한다는 주장은 아니다. [Evaluation Contract](docs/evaluation-contract.md)와 [Match Semantics](docs/match-semantics.md)에 기준과 한계를 기록했다.

## Quick Start

**Requirements:** Go 1.24 이상, Jev/TypeSafe API 이용이 가능한 `TYPESAFE_API_KEY`, 준비된 canonical profile과 JD. 외부 호출에는 네트워크와 provider 사용 비용이 필요하다.

Repository root에서, key를 환경변수로 설정한 shell을 사용한다. 이미 로컬 `.env`를 준비했다면 다음처럼 읽을 수 있다. 실행 파일 자체는 `.env`를 자동 로딩하지 않는다. Key를 Git이나 API request에 넣지 않는다.

```sh
set -a
. ./.env
set +a

go build -o /tmp/job-fit ./cmd/job-fit
go build -o /tmp/job-fit-server ./cmd/job-fit-server
```

로컬 파일 기반 평가:

```sh
/tmp/job-fit --profile ./data/profile.json --job ./samples/proptech-plus.txt
```

HTTP server 실행:

```sh
/tmp/job-fit-server --profile ./data/profile.json
```

다른 terminal에서 요청:

```sh
curl --request POST http://127.0.0.1:8080/v1/job-fit/evaluate \
  --header 'Content-Type: application/json' \
  --data-binary '{"job":{"description":"Position: Backend Engineer\nRequired: Backend API design and development"}}'
```

Profile/JD는 Jev로 전송된다. 저장된 profile은 작성자의 실제 경력을 요약한 personal evidence이며 synthetic demo profile이 아니다. 자신의 경력을 평가하려면 별도 로컬 profile 파일을 `--profile`에 지정한다. [입력 형식과 상세 사용법](docs/api/job-fit.md)을 참조한다.

## API

`POST /v1/job-fit/evaluate` — JD 한 건을 server-side canonical snapshot과 동기 평가한다.

`job.description`은 필수다. Company/title/sourceUrl은 optional metadata이며 평가 본문에 자동 삽입하거나 URL을 fetch하지 않는다. Profile은 startup에서 한 번 로딩하며 변경은 restart 후 반영한다. 응답의 SHA-256 `profileVersion`은 사용한 파일 bytes를 식별한다.

기본 bind는 `127.0.0.1:8080`이며 numeric loopback만 허용한다. Auth는 없다. Body 128 KiB, description 16 KiB/64줄, evidence pairs 512, 전체 evaluation deadline 기본 5분/Jev 호출별 60초다. HTTP error는 안전한 common envelope로 반환하고 cancellation을 provider 요청까지 전달한다. [API Specification](docs/api/specification.md)에서 오류와 설정을 확인할 수 있다.

## Testing

일반 테스트는 외부 Jev API/key 없이 실행한다.

```sh
go test ./...
go test -race ./...
go vet ./...
go build -o /tmp/job-fit ./cmd/job-fit
go build -o /tmp/job-fit-server ./cmd/job-fit-server
```

[6개 실제 JD regression](docs/development/turn-011.md) 후 evaluator baseline을 수용하고 iterative tuning을 종료했다. [실제 HTTP/Jev 검증](docs/development/turn-014.md)도 기록했다. 이는 label accuracy나 일반화 성능의 인증이 아니다. [최종 검증 기록](docs/development/turn-018.md)과 [공개 준비 점검](docs/development/turn-019.md)에 regression 및 repository/privacy 검토 결과를 남긴다.

## AI-assisted Development

ChatGPT 대화에서 문제·범위·architecture·다음 Turn을 검토하고, Codex workspace에서 코드·테스트·문서 작성과 실행 검증을 수행했다. AI가 상당 부분을 구현했다는 점을 공개하며, 사람은 scope·판단 기준·수용 여부를 결정했다. 기여 비율이나 모든 코드에 대한 사람의 line-by-line 검토를 주장하지 않는다.

핵심 학습은 “Reference Strong / Jev Partial을 어떻게 같게 만들까?”보다 **어떤 차이가 합리적 해석이고 어떤 차이가 명시적 계약 위반인가**를 구분하는 것이었다. [AI-assisted Development Story](docs/ai-assisted-development.md)에 관측, triage, 두 위반 수정, baseline freeze와 REST 경계 선택 과정을 정리했다.

## Limitations

- 객관적 Job Fit ground truth를 가정하지 않는다. Jev 판단은 실행 간 변동할 수 있다.
- Decision Trace는 사후 attribution이며 최초 판단의 인과성을 입증하지 않는다.
- 줄 단위 extraction과 유한 rationale은 모든 compound/context gap을 표현하지 못한다. Public/Energy Context ambiguity가 남아 있다.
- 현재 profile은 한 개인의 Software Engineering 경력 중심이고, runtime은 language/education 조건을 평가에서 제외한다. 전체 지원 자격 판정이 아니다.
- Personal loopback service이며 generic candidate endpoint, auth, persistence/history, batch, UI는 없다. 장기 운영·부하·통계적 정확도 검증을 주장하지 않는다.

## Future Work

Generic candidate+JD API, OpenClaw enrichment, optional thin MCP adapter, batch/automation, OpenAPI/Swagger는 별도 검토 가능성이다. **현재 구현이나 portfolio MVP 완료 조건이 아니다.** 검색·회사 분류·crawler·Notion 저장·최종 지원 결정은 job-fit core 책임 밖이다. [Final Scope Freeze](docs/development/turn-015.md)를 따른다.

## Documentation

| 문서 | 내용 |
| --- | --- |
| [API Specification](docs/api/specification.md) | Endpoint, 공통 response/error 빠른 참조 |
| [Job Fit API](docs/api/job-fit.md) | Request/response, 실제 전체 응답, curl, limits/timeout |
| [API Contract](docs/api-contract.md) | 외부 계약과 설계 이유 |
| [Evaluation Contract](docs/evaluation-contract.md) | 평가 책임·원칙·품질 기준 |
| [Match Semantics](docs/match-semantics.md) | Match 의미와 사전 reference의 역할 |
| [Architecture Overview](docs/architecture/overview.md) | Component/dependency/runtime 경계 |
| [AI-assisted Development](docs/ai-assisted-development.md) | Turning points와 Human/AI 역할 |
| [Development History](docs/development/) | Turn 기록과 verification artifacts |

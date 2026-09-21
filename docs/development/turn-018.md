# Turn 018 — Portfolio Polish & Final Verification

- 날짜: 2026-09-21
- 범위: README/narrative 정리, final regression, repository/publication hygiene 검토. 제품 기능·API·평가 behavior 변경 없음.
- 시작 상태: clean working tree, HEAD `027869c` (`docs: document job-fit architecture`).

## Completed

README를 portfolio landing page로 재구성하고 [AI-assisted development narrative](../ai-assisted-development.md)를 작성했다. Current code, 실제 artifact 및 Turn 기록에 근거하여 목적·가치·실행 방법·한계를 정리했다. Scope는 [Turn 015](turn-015.md)의 frozen MVP를 유지한다.

최종 unit/race/vet/build와 help smoke check를 완료했다. 실제 Jev 호출은 반복하지 않았으며 기존 [Turn 014 artifact](turn-014-results/proptech-plus.json)를 사용했다. 개인정보는 삭제/변형하지 않고 공개 전 검토 결과를 아래에 기록한다.

## Changed Files

- [README](../../README.md): Why/What/Example/Architecture/Principles/Quick Start/API/Testing/AI-assisted/Limitations/Future Work/Documentation 구조.
- [AI-assisted Development](../ai-assisted-development.md): 신규 narrative, turning points, Human/AI 책임과 한계.
- 본 Turn 기록.
- `docs/api/job-fit.md`: README 재구성에 따라 실행 안내 anchor를 `#quick-start`로 변경한 링크 수정만 수행.

기존 code/tests/profile/JD/prompt/계약/architecture/Turn history/artifacts는 수정하지 않았다. API 문서는 README anchor 링크 하나만 수정했고 내용/behavior는 유지했다. Commit 또는 publication을 수행하지 않았다.

## README

상단에서 Go/Jev personal evaluator라는 정의와 hiring probability/recommendation이 아니라는 경계를 제시했다. General LLM과의 차이는 정확도 우위가 아닌 canonical evidence, explicit semantics, structured output/validation/traceability 및 재사용 가능한 interface다.

기존 핵심 철학·한계·실행 명령은 유지하고 상세 API/runtime/history 설명은 전문 문서 링크로 이동했다. 짧은 JSON response는 실제 Turn 014에서 일부 field를 발췌한 것으로 명시했다. 완전한 response schema나 위 Example Corp request의 실행 결과라고 주장하지 않는다. CLI/REST가 Jev를 직접 각각 호출하는 그림 대신 공유 service/core 경계를 표시했다.

## AI-assisted Development Narrative / Key Portfolio Story

기록의 중요한 전환점만 선택했다: disposable PoC와 greenfield 구현 분리, semantics 먼저 정의, Human Baseline trap, observability before calibration, contract triage, baseline 수용/tuning 종료, CLI에서 reusable REST boundary로 발전한 과정.

Turn 007은 과거에 reference에 맞춰 tuning했다는 뜻이 아니라고 명시한다. Narrative도 Strong/Partial 불일치를 없애려는 개선 방향의 위험을 발견했다고 설명하고, 실제 calibration을 했다가 되돌렸다는 이야기를 만들지 않았다. 초기 Human Baseline이 Codex 초안에서 출발했다는 provenance와 실제 triage 결과(Confirmed 2, Ambiguous 1)도 보존했다.

AI가 상당 부분의 코드·테스트·문서를 작성하고 실행 검증했음을 공개한다. Human의 역할은 실제 문제/경력/workflow 제공, scope·철학·architecture·수정 범위·acceptance/종료 결정이다. 기여 percentage, 모든 코드의 human line-by-line review 또는 기록 밖의 planning 세부 검토를 주장하지 않는다.

## Documentation Index

README에서 API Specification, Job Fit API, API Contract, Evaluation Contract, Match Semantics, Architecture Overview, AI-assisted Development 및 Development History로 연결한다. 전체 Turn 목록을 landing page에서 반복하지 않는다.

## Final Regression

모두 통과했다.

```sh
go test ./...
go test -race ./...
go vet ./...
go build -o /tmp/job-fit ./cmd/job-fit
go build -o /tmp/job-fit-server ./cmd/job-fit-server
```

추가 확인:

- `gofmt -l cmd internal`: 출력 없음, source formatting 변경 불필요.
- `/tmp/job-fit --help`, `/tmp/job-fit-server --help`: exit 0, 실제 flags/defaults와 README 일치.
- 별도 Makefile/lint configuration/CI formatter command는 발견하지 않았다.
- 일반 테스트는 fake Engine/transport와 로컬 httptest를 사용하며 실제 Jev 호출 없이 실행했다.
- 로컬 Go 실행 환경에서 검증했으며 일부 기존 테스트는 Go cache 결과를 사용했다. Cache를 숨기거나 새로운 live evaluation으로 표현하지 않는다.

새 behavior 변경이 없어 외부 API 비용을 발생시키는 재호출은 하지 않았다.

## Repository Hygiene

| 점검 | 결과 / 범위 |
| --- | --- |
| Tracked secret | 현재 82개 tracked 파일에서 private-key header/고정 credential assignment 패턴 발견 없음. 실제 로컬 Jev key와 tracked 내용의 비출력 비교도 일치 없음 |
| Env tracking | `.env`는 ignored. Reachable history에서 `.env`/`.env.*` 경로를 추적한 기록 없음 |
| Tracked binary | Executable magic 검사에서 없음. Final build는 요청대로 `/tmp`에 생성 |
| Tracked temporary files | tmp/bak/log/backup 패턴 없음. 이번 검사 스크립트/중간 자료는 `/tmp` 사용 |
| Existing ignored files | `.env`, `bin/job-fit`, `bin/job-fit-server`, 루트 `job-fit`, `.DS_Store` 파일들이 로컬에 존재하고 ignore됨. 임의 삭제하지 않음 |
| .gitignore | `.env`, `.env.*`(example 예외), `/bin/`, `/job-fit`, `.DS_Store`가 현재 사용 패턴을 보호함. 변경하지 않음 |
| Email scan | Test fixture의 `secret@example.com`만 발견. 실제 연락처가 아님 |
| Documentation | 새 링크/anchors/JSON/shell 예시 및 formatting 정적 검증. Local absolute reference는 portability 제한으로 별도 기록 |

이는 현재 tracked content와 특정 env-path history의 제한된 점검이다. 모든 과거 blob의 모든 가능한 secret을 탐지한 보안 인증은 아니다. Secret 값은 출력·새 artifact 저장하지 않았다.

## Public Portfolio / Privacy Review

`data/profile.json`에는 이름·연락처·생년월일·주소·credential이 없다. 그러나 **실제 개인 경력 요약**이며 synthetic sample은 아니다. 약 16년 경력, Tech Lead/PM/PL, 팀 규모, 기술과 업무 도메인의 조합은 작성자와 연결되면 개인 정보가 된다. NDA/고객 비밀에 해당하는지는 저장소만으로 확인할 수 없다.

같은 evidence는 과거 evaluation 결과 JSON, API 상세 success example, Match Semantics/Turn 기록에도 반복된다. 따라서 공개 범위를 줄이려면 profile 파일 하나만 바꾸는 것으로 충분하지 않다. 이번 Turn에서는 input/result/history를 수정하지 않았다.

공개 전 필요한 결정:

1. 소유자가 현재 경력 요약 및 반복된 evidence의 공개 가능 여부를 확인한다. 이번 Turn에서 선택 질문을 전달했으며 응답 전에는 동의로 간주하지 않는다.
2. 공개 불가라면 별도 작업에서 public synthetic/sample profile과 private local profile 분리 및 파생 examples/artifacts/history의 공개 범위를 정한다. 현재 평가 baseline을 조용히 변경하지 않는다.
3. `docs/api-contract.md`, `docs/api/specification.md` 및 일부 과거 Turn에는 사용자 홈 디렉터리의 절대경로 로컬 사용자 경로가 있다. 이는 공통 문서 reference로 요청되어 기록된 것이며 다른 checkout에서 동작하지 않는다. 공개 시 계정명 노출/로컬 링크를 그대로 유지할지 검토한다. 이번에 임의 삭제하지 않았다.

소스와 테스트의 완료 여부와 공개할 개인정보 범위의 승인은 구분한다.

## Known Limitations

Objective ground truth 부재, Jev variability, post-hoc Trace, compound/context ambiguity, 개인 Software Engineering profile, language/education 미평가, 줄 기반 extraction/유한 rationale, loopback-only/no auth/no persistence/no generic candidate endpoint를 README와 narrative에서 숨기지 않았다. 통계적 정확도·장기 운영·부하 테스트를 완료했다고 주장하지 않는다.

## Future Work

Generic candidate+JD API, OpenClaw enrichment, optional thin MCP adapter, batch/automation, OpenAPI/Swagger는 별도 scope 검토 가능성이다. 현재 구현이나 완료 조건이 아니며 이번에 시작하지 않았다.

## Portfolio Readiness

| Reviewer 질문 | 확인 |
| --- | --- |
| 2–3분 안에 프로젝트를 이해할 수 있는가? | 상단 정의/Why/What/Example/Architecture로 목적·결과·경계를 요약했다. 실제 독자 시간 측정은 하지 않음 |
| History 없이 실행할 수 있는가? | Requirements/env/build/CLI/server/curl/tests를 README에 제공하고 도움말과 대조 |
| Jev를 쓰는 이유가 보이는가? | Workflow 안의 constrained structured decision component라는 목적 명시 |
| Evaluator가 주장하지 않는 것을 알 수 있는가? | 확률/지원 추천/정답/정확도 우위 부정, reference의 역할 명시 |
| API/architecture를 더 볼 수 있는가? | 전문 문서 index와 source/기록 링크 제공 |
| AI 사용을 이해할 수 있는가? | 별도 narrative와 Human/AI 역할 표, 실제 turning points/provenance |
| 한계가 눈에 보이는가? | README Limitations 및 narrative의 ambiguity/variability/security 범위 |

## Problems / Findings

Application/regression 실패는 없다. 공개 준비 측면에서는 실경력 요약과 반복된 evidence의 공개 허용 여부 및 로컬 사용자 경로 노출 여부가 아직 확인되지 않았다. 이를 제품 기능 추가나 tuning 요구로 바꾸지 않는다. 파일을 삭제/변형하지 않고 필요한 결정만 남겼다.

## Validation

README와 API 문서의 실행 안내 링크 수정을 제외한 기존 파일의 SHA-256 보존, 신규 문서 두 개 추가, repository 문서 링크/anchor/code fence, README JSON과 실제 응답 발췌 일치, curl 인자/JSON 및 shell syntax, `git diff --check`를 검증했다. README 재구성으로 끊어진 기존 API 문서 anchor 한 곳을 Quick Start로 갱신했다. URL 또는 로컬 server를 호출하지 않고 정적으로 확인했다.

## Deferred

OpenAPI/Swagger 및 모든 future integration은 MVP 이후 작업이다. 개인정보 분리/공개 여부는 소유자 판단 없이 실행하지 않는다. Repository 공개·배포·commit/push도 수행하지 않는다.

## Final Status

**Additional Work Required** — Application MVP와 final regression은 완료했다. 공개 전에 필요한 작업은 현재 실경력 evidence 및 로컬 사용자 경로의 공개 범위 확인뿐이다. 공개 불가로 판단한 정보가 있다면 별도 범위에서 public sample/파생 문서 분리 후 공개한다. 새 기능 구현이나 evaluator tuning은 필요 작업으로 제안하지 않는다.

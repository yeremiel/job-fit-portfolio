# Turn 000 — Project Bootstrap

- 날짜: 2026-09-21
- 범위: workspace 및 Git 상태 확인, 최소 문서 구조 준비, 프로젝트 시작 상태 기록
- 결과: README와 본 개발 기록 생성. Application code 구현 없음.

## Background

여러 채용공고를 검토할 때 동일한 경력 정보를 기준으로 각 공고와의 적합성을 반복적으로 판단해야 한다. 일반적인 LLM에게 단순히 적합성을 질문하면 다음 문제가 발생할 수 있다.

- 평가 기준이 매번 달라질 수 있다.
- Candidate에게 실제로 없는 경험을 추론할 수 있다.
- 결과가 구조화되어 있지 않아 여러 공고를 비교하기 어렵다.
- 확률처럼 보이는 값을 실제 채용 성공 확률로 오해할 수 있다.
- 공고에 맞춰 수정된 resume를 평가 입력으로 사용하면 evaluation bias가 발생할 수 있다.

이 프로젝트는 동일한 candidate evidence와 명시적인 평가 기준을 사용하여 Job Fit을 구조적으로 평가하는 도구를 만드는 것을 목표로 한다. AI-assisted development 과정 자체도 포트폴리오의 일부로 기록한다.

사용자가 제공한 배경에 따르면, 이전에는 shell script와 Jev API를 사용한 disposable PoC로 structured job-fit evaluation의 feasibility를 확인했다. 이 PoC는 가능성 확인 용도였으며, 기존 코드는 재사용하지 않는다. 이번 프로젝트는 새로운 구현으로 시작한다. 이번 Turn에서 PoC 자체를 검토하거나 그 결과를 재검증하지는 않았다.

## Initial Ideas — 미확정 검토 후보

다음은 초기 아이디어이며 채택된 architecture decision이 아니다. 이후 Turn에서 필요성과 trade-off를 검토한다.

- Canonical candidate profile
- Explicit evaluation criteria
- Jev를 structured evaluation engine으로 사용하는 방법
- Deterministic rule과 semantic judgment의 분리
- MCP interface
- TypeScript implementation
- AI coding agent와의 integration

구체적인 schema, 평가 방식, 기술 스택 및 인터페이스는 아직 결정하지 않았다.

## Non-goals

- 실제 채용 합격 확률 예측
- 다른 지원자와의 경쟁력 예측
- 자동 입사 지원
- AI가 지원 여부를 최종 결정하는 시스템
- 대규모 production infrastructure 구축

최종 판단은 항상 사람이 한다.

## Development Process

기획 및 검토와 구현의 역할을 분리한다.

| 환경 | 역할 |
| --- | --- |
| 별도의 ChatGPT planning conversation | Brainstorming, problem definition, architecture discussion, Turn scope 결정, 이전 Turn 결과 검토, 다음 Turn 계획 |
| Codex workspace | Repository 분석, implementation planning, coding, testing, code review, refactoring, commit 준비, Turn Result 작성 |

각 Turn 결과는 planning conversation에서 검토한 뒤 다음 Turn을 결정한다. Codex가 다음 Turn 범위를 임의로 구현하지 않는다.

개발 기록은 우선 `docs/development/turn-000.md`와 같은 Turn 단위 문서로 관리한다. 중요한 architecture decision이나 experiment가 생기면 필요에 따라 별도 문서 구조를 추가한다. 처음부터 과도한 문서 계층을 만들지 않는다.

## Workspace / Git 확인 결과

작업 시작 시점에 확인한 상태다.

- Workspace: repository root (`.`)
- Workspace에는 `.git` 디렉터리만 존재했다.
- Git repository는 이미 초기화되어 있었고 현재 브랜치는 `main`이었다.
- 기존 커밋, 추적 파일 및 작업 트리 변경 사항은 없었다.
- 등록된 원격 저장소는 없었다.
- Workspace 및 상위 경로에서 적용할 `AGENTS.md`는 발견되지 않았다.

커밋이 없어 Git 로그 조회가 실패했으며, 이는 초기 저장소 상태에 따른 결과다.

## Completed / Changed Files

최소 기본 구조로 다음 두 문서만 생성했다.

```text
README.md
docs/
└── development/
    └── turn-000.md
```

- `README.md`: 프로젝트 목적, 현재 단계, 주요 비목표, 작업 방식 및 개발 기록 링크
- `docs/development/turn-000.md`: 배경, 미확정 아이디어, 비목표, 개발 방식, 초기 상태 및 이번 Turn 결과

## Decisions Made

- 기본 구조는 요청한 두 문서와 이를 담는 디렉터리로 한정했다.
- 별도 architecture/experiment 문서 계층이나 빈 application 디렉터리는 추가하지 않았다.
- 추가적인 제품 또는 architecture decision은 내리지 않았다.
- 기존 Git repository를 유지했으며, 커밋 생성과 원격 저장소 설정은 수행하지 않았다.

## Validation / Problems

문서 내용이 요청 범위와 일치하는지, 미확정 아이디어가 결정 사항으로 표현되지 않았는지 확인했다. README의 상대 링크 대상과 최종 파일 목록, Git 변경 상태 및 공백 오류를 확인했다. Application code나 실행 환경을 추가하지 않았으므로 실행 테스트는 해당하지 않는다.

진행을 막는 문제는 발견하지 않았다. 커밋 이력과 원격 저장소는 없는 상태이며, 필요 시 이후 Turn에서 다룬다.

## Deferred

이번 Turn에서는 다음 작업을 의도적으로 수행하지 않았다.

- TypeScript project 초기화 및 `package.json` 생성
- Dependency 설치
- Jev API client 구현
- MCP server 구현
- Candidate profile schema 설계
- Evaluation criteria 설계
- Database 도입
- Application code 작성

초기 아이디어의 채택 여부와 구체적인 architecture는 이후 논의에 맡긴다.

## Suggested Next Step

Planning conversation에서 이번 Turn 결과를 검토하고, 첫 사용 시나리오와 평가 결과가 사람의 판단에 제공해야 할 정보를 구체화하는 것을 제안한다. 이를 바탕으로 다음 Turn의 범위와 완료 기준을 정한다. 이 제안은 확정된 계획이나 architecture decision이 아니다.

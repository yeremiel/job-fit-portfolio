# Turn 001 — Project Positioning & Evaluation Contract

- 날짜: 2026-09-21
- 범위: 프로젝트 포지셔닝과 개념 수준 Evaluation Contract 문서화
- 결과: README 갱신, Evaluation Contract 및 본 기록 생성. Application code 작성 없음.

## Completed

- 특정 채용 플랫폼에 종속되지 않는 개인용 Software Engineering job matching 도구로 포지셔닝을 정리했다.
- 작성자 본인을 primary user로, Backend Engineer / Tech Lead 경력을 reference profile로 기록했다.
- Canonical Candidate Profile과 회사별 제출용 resume의 역할 및 evaluation bias 문제를 구분했다.
- MVP 평가 범위를 Capability / Experience Match로 한정하고, 평가 제외 항목을 명시했다.
- Evidence First, Unknown Remains Unknown, Explainable Result, No Hiring Probability, Human Decision 원칙을 기록했다.
- 설명 가능한 출력 정보와 requirement-level 결과 예시, optional Match Score의 의미를 정리했다.
- Conceptual flow와 구현 architecture를 구분하고, 미정 사항을 명시했다.
- Turn 000과 현재 문서의 일관성, 문서 링크, 파일 범위 및 공백 오류를 확인했다. 문서만 변경하여 application 실행 테스트는 해당하지 않는다.

## Changed Files

- [README.md](../../README.md): 현재 목적, 사용자, 평가 범위와 단계 및 문서 링크 갱신
- [Evaluation Contract](../evaluation-contract.md): 평가 목적, 입력의 역할, 범위, 원칙, 출력 개념, 비목표 및 미정 사항 기록
- `docs/development/turn-001.md`: 이번 Turn의 합의 내용과 작업 결과 기록

## Decisions Made

Prompt에서 주어진 결정 외에 추가적인 제품 또는 architecture decision은 없다. 권장 경로인 `docs/evaluation-contract.md`를 사용했다.

이번 prompt에서 확정한 내용은 개인용 도구라는 포지셔닝, 초기 사용자와 검증 대상, canonical evidence source의 역할, Capability / Experience Match라는 평가 범위 및 다섯 가지 평가 원칙이다. 출력 항목은 제공 목표이며, 상세 판정 규칙이나 구현 schema를 확정한 것은 아니다.

## Problems / Findings

[Turn 000](turn-000.md)과의 모순은 발견하지 않았다. Turn 000은 당시의 미확정 상태와 작업 범위를 기록하므로 수정하지 않았다. 이번 Turn에서 canonical profile의 역할과 평가 범위·원칙을 구체화한 것은 이후 합의의 진전이다. Schema와 상세 평가 기준, 기술 선택은 여전히 미정이다.

Turn 000의 최소 문서 구조 원칙에 따라, 이번에 요청된 Evaluation Contract만 별도 문서로 추가했다. 이전 PoC 코드 미재사용, 사람의 최종 판단 및 planning/implementation 역할 분리는 유지했다.

Weak / Missing Evidence와 Unknown의 구체적인 판정 규칙은 불명확한 상태로 남아 있다. 현재 문서는 evidence 부족을 실제 경험 부재로 단정하지 않는다는 원칙을 명시하며, 별도 rubric을 임의로 만들지 않았다. 점수와 `HIGH` 등의 표시는 예시이고, 실제 결과 또는 확정된 threshold가 아니다.

작업 시작 시 Git은 `main`에 커밋이 없는 상태였으며, Turn 000의 두 문서는 untracked 상태였다. 이번 Turn에서도 커밋은 생성하지 않았다.

## Deferred

- Candidate Profile JSON schema 및 데이터 모델
- 상세 evaluation criteria, requirement-level 판정 규칙 및 Unknown 처리의 구체화
- Match Score 도입 여부, 계산 공식, 가중치, threshold 및 calibration 실행
- Jev 채택 여부와 API client 구현
- Deterministic rule / semantic judgment 역할 분담
- MCP, TypeScript, Database, Web UI 및 AI coding agent integration 도입
- Package 초기화, dependency 설치 및 모든 application code 작성

## Suggested Next Step

Planning conversation에서 Evaluation Contract를 검토하고, 소수의 JD와 사실 기반 career evidence 사례를 통해 Strong / Partial / Weak / Unknown을 어떻게 구분해야 하는지 논의하는 것을 제안한다. 특히 evidence 누락과 확인된 불일치를 구분하는 데 필요한 정보를 검토하면 다음 Turn의 범위와 완료 기준을 정하는 데 도움이 된다.

이는 다음 Turn의 검토 제안이며, 이번 Turn에서 샘플 평가나 schema 설계 및 구현을 시작하지 않았다.

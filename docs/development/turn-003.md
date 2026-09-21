# Turn 003 — Baseline Review & Evaluation Rules

- 날짜: 2026-09-21
- 목적: 사용자 제공 평가 규칙을 명시하고 Turn 002 baseline을 비교 기준으로 확정한다.
- 범위: 문서 검토·갱신. Application code, Jev 호출 및 별도 AI evaluator 실행 없음.

## Completed

- Required와 Preferred의 중요도 차이 및 core role mismatch 우선 원칙을 기록했다.
- Strong / Partial / Weak / Unknown 정의와 Missing Evidence의 별도 설명 역할을 정리했다.
- Transferable experience를 검토할 여섯 관점과 specialized experience를 추정하지 않는 원칙을 명시했다.
- Overall을 단순 평균이 아닌 역할·필수 요구·주요 capability·전이 경험·critical gap·우대 사항의 정성적 요약으로 정의했다.
- 기존 46개 requirement-level 행과 reference evidence를 보존하고, 6개 Overall을 유지하여 확정 근거와 상태를 갱신했다.
- README와 Evaluation Contract를 현재 규칙에 맞게 연결했다. Turn 000–002 개발 기록은 수정하지 않았다.

## Changed Files

- [README.md](../../README.md): 현재 규칙·baseline 확정 상태와 Turn 링크
- [docs/evaluation-contract.md](../evaluation-contract.md): 정성적 평가 규칙 반영
- [docs/match-semantics.md](../match-semantics.md): 규칙, sample별 확정 근거, 과거 ambiguity와 현재 결정 구분
- `docs/development/turn-003.md`: 이번 Turn 결과

## Human Baseline Summary

Turn 003 prompt의 확정 지시에 따라 기존 값을 유지했다. 새로운 candidate 사실이나 JD 요구사항을 추가한 결과가 아니다.

| Sample | Turn 002 → Turn 003 | 확정 근거 |
| --- | --- | --- |
| BlueMeme / Micro Court | Strong → Strong | 필수 개발·설계·고객 조율 근거가 직접적이며 OutSystems 사전 경험은 필수가 아님. |
| Prop Tech plus | Strong → Strong | 핵심 backend/API 역량에 직접 근거. Azure 우대 부족으로 크게 낮추지 않음. Java 버전·기간 미확인은 계속 기록. |
| PIA TECH LAB | Partial → Partial | AWS 필수 운영 gap으로 Strong은 지지하지 않지만 Web 설계·개발 및 전이 가능한 경험은 인정. |
| CORE | Strong → Strong | 기본설계, 고객 과제 파악 및 리딩의 직접 근거. 특화 도메인·전 lifecycle gap은 별도 기록. |
| teamLab | Partial → Partial | Backend·고객 설계 경험은 관련되나 전체 architecture·운영 맥락·CI/CD 등의 evidence는 제한적. |
| Nitori Digital Base | Weak → Weak | 일반 리딩·기술 overlap보다 frontend Tech Lead라는 핵심 역할의 evidence gap을 우선. |

Baseline 확정은 모든 세부 사실 검증이나 필수 지원 자격 충족을 뜻하지 않는다. Codex가 작성한 초안에 사용자 제공 규칙을 반영한 기록이며, 사람이 AI assistance 없이 독립 작성한 ground truth로 표현하지 않는다.

## Decisions Made

사용자 지시의 정성적 규칙을 채택했다. Core Required mismatch가 Overall의 상한에 영향을 줄 수 있다는 원칙은 적용하지만, 모든 Required Weak를 동일한 Overall로 변환하는 자동 규칙은 만들지 않았다. PIA Partial과 Nitori Weak의 차이는 각 sample의 역할·근거 관계로 설명한다.

Turn 002의 ambiguity는 삭제하지 않고 당시 기록으로 표시한 뒤 현재 review outcome을 추가했다. Nitori 예시를 실제 경력에 대한 새 부정적 사실로 해석하지 않았다. 직접 frontend 리딩 evidence 부족과 실제 경험 부재를 구분했다.

## Problems / Findings

- Git은 여전히 `main`에 커밋이 없는 상태다. Untracked 문서를 검토했고 커밋은 생성하지 않았다.
- Turn 002의 사람 검토 전 상태와 현재 확정 상태가 섞이지 않도록 README 및 현재 의미 문서를 갱신했다. 과거 Turn의 초안 상태 표시는 역사적 기록으로 보존했다.
- Prop Tech의 Java 버전·기간, PIA의 AWS 운영 기간·책임, Nitori의 frontend 리딩 범위 등 미확인 사실은 확정 이후에도 미확인이다.
- teamLab의 언어 조건 및 복합 requirement 분해 방식은 이번 prompt에서 결정하지 않았다.
- JD 재수집은 하지 않았다. Turn 002에서 확인한 URL·검토일·요약을 그대로 사용했다.

## Validation

- 6개 sample과 46개 requirement-level 행, reference evidence 및 출처 정보가 보존되었는지 확인했다.
- Sample별 Overall과 요약의 일치, 확정 상태, 과거 질문과 현재 결정의 구분을 검토했다.
- 문서 상대 링크, 공백, 코드블록, 파일 범위 및 Turn 000–002 파일 보존을 확인했다.
- Application code나 schema, dependency는 추가하지 않았다. 실행 대상 application이 없어 application 테스트는 해당하지 않는다.

## Deferred

- 언어·학력 조건을 다룰 일반적인 평가 범위와 복합 requirement 분해 방식
- Match Score, percentage, weighting formula, threshold 및 calibration
- Candidate Profile schema, 기술 스택과 application architecture
- Jev/별도 evaluator 실행 및 실제 baseline 비교 실험

## Suggested Next Step

Planning conversation에서 이 규칙과 baseline을 바탕으로 첫 비교 실험의 입력 범위·관찰 항목·완료 기준을 정하는 것을 제안한다. 언어 조건과 requirement 분해 단위가 그 실험에 필요한지도 함께 검토한다. 실험이나 구현은 이번 Turn에서 시작하지 않았다.

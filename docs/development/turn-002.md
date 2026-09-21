# Turn 002 — Match Semantics & Human Baseline

- 날짜: 2026-09-21
- 목적: 평가 engine 실행 전에 match semantics와 공개 JD 기반 예상 판단을 기록한다.
- 상태: Codex 작성 baseline 초안 완료. 사람의 독립 판단 또는 승인 완료를 뜻하지 않는다.

## Completed

- [Match Semantics](../match-semantics.md)에 Strong / Partial / Weak / Unknown과 별도 설명 정보인 Missing Evidence를 정의했다.
- 사용자 제공 reference evidence만 사용하고, 문서 내부 참조 번호로 판단 근거를 연결했다. Profile schema는 만들지 않았다.
- 지정된 공개 JD 6개의 현재 본문을 확인하고 company/position, URL, reviewed date, 관련 요구 요약을 기록했다.
- Requirement별 Supporting Evidence, Missing / Weak Evidence, Unknowns, Expected Match, Reasoning을 작성했다.
- Required / Preferred / Responsibility / Stack을 구분하고, Overall을 수치 없이 정성적으로 요약했다.
- 계약의 미정 표현을 이번 정의와 연결하고 README에 현재 단계와 문서 링크를 반영했다.

## Changed Files

- [docs/match-semantics.md](../match-semantics.md): 의미 정의, reference evidence, sample별 baseline, 경계 문제
- [docs/evaluation-contract.md](../evaluation-contract.md): level·Missing Evidence의 관계 및 Turn 002 링크 반영
- [README.md](../../README.md): 현재 단계와 baseline 초안 상태, Turn 링크 갱신
- `docs/development/turn-002.md`: 이번 작업 결과 기록

Turn 000과 Turn 001 문서는 당시 기록으로 보존했다.

## Human Baseline Summary

아래는 사람 검토 전 잠정 예상이다. 각 원문과 상세 근거는 Match Semantics의 해당 sample에 연결되어 있다.

| Sample | Overall Expected Match |
| --- | --- |
| BlueMeme / Micro Court | Strong |
| Prop Tech plus | Strong |
| PIA TECH LAB | Partial |
| CORE | Strong |
| teamLab | Partial |
| Nitori Digital Base | Weak |

정답 label, 지원 추천 또는 채용 확률이 아니다. Jev 및 별도 AI evaluator를 실행하지 않았다. 단, Codex가 판단 초안을 작성했으므로 AI assistance 없이 사람이 작성한 baseline이라고 표현하지 않았다.

## Ambiguities Found

- BlueMeme의 도구 경험 불문 조건과 transferable experience가 전체 Strong에 충분한가?
- Prop Tech의 필수 Java 버전·기간이 미확인일 때 Overall Strong을 유지할 수 있는가?
- PIA의 AWS 필수 Weak와 Web 설계 Strong 조합을 Overall Partial 또는 Weak 중 어떻게 표현할 것인가?
- Nitori의 인접 frontend·리딩 근거를 Weak로 볼지, 직접 역할 정보가 없어 Unknown으로 볼지?
- 필수·우대·업무·환경의 차이와 복합 requirement의 분해 단위를 어떻게 유지할 것인가?
- teamLab의 일본어와 같이 명시된 비기술 역량을 현재 evaluation scope에서 어떻게 다룰 것인가?
- Codex 초안을 사람이 검토한 baseline과 어떻게 구분하고 기록할 것인가?

중요한 경계는 확정 규칙으로 만들지 않고 planning conversation의 검토 대상으로 남겼다.

## Decisions Made

사용자가 제시한 네 level과 Missing Evidence 분리를 문서에 반영했다. 추가적인 architecture나 수치 판정 규칙은 결정하지 않았다.

기록 방식으로 다음을 적용했다.

- 원문의 필수·우대와 업무·환경을 구분하여, stack을 자동으로 필수 조건으로 승격하지 않는다.
- Overall은 개별 근거를 설명하는 잠정 요약으로 작성하고, 중요한 대안 해석을 함께 남긴다.
- 작성 주체와 검토 상태를 명시하여 “Human Baseline”을 사람 검토가 끝난 결과로 오인하지 않게 한다.

이 방식은 이번 문서의 추적 가능성을 위한 것이며, 자동 evaluator의 최종 semantic policy를 확정한 것은 아니다.

## Problems / Findings

- BlueMeme: 웹 읽기 도구에서 접근 실패. 동일 공개 URL의 HTML을 읽어 복구했다. OutSystems 사전 경험 불문 조건을 확인했다.
- PIA: 웹 읽기 도구는 502를 반환했지만 직접 공개 HTML에서는 본문 확인이 가능했다. 직무명과 달리 AWS 인프라 구축·운영 3년 이상이 중요한 필수다.
- teamLab: 정적 읽기는 제목 수준만 반환했다. 브라우저 렌더링으로 본문을 확보했다. 기술 목록을 모두 필수로 해석할 근거는 없었다.
- Nitori: Frontend Tech Lead가 실제 필수임을 확인했다. Java/Spring Boot는 조직 backend stack이며 해당 역할의 직접 근거를 대신하지 않는다.
- 모든 sample의 requirement를 확인했으므로 접근 실패만으로 Unknown 처리한 sample은 없다. Candidate evidence의 Unknown은 별도로 유지했다.
- 경력은 사용자 제공 가정이며, 기간·범위·언어·자격 등을 추가 추론하지 않았다.
- 기존 원칙과 모순은 발견하지 않았다. Turn 001의 미정 의미를 이번 Turn에서 구체화했으며, 과거 기록은 수정하지 않았다.

## Validation

6개 sample의 출처·날짜·판단 항목 및 overall의 일관성을 검토했다. 문서 링크, 공백·코드블록 형식, 최종 파일 범위 및 이전 Turn 기록 보존을 확인했다. Application code, dependency, schema, 점수 공식 또는 JD 전체 원문 파일은 repository에 추가하지 않았다. 실행할 application이 없으므로 application 테스트는 해당하지 않는다. Git 커밋은 생성하지 않았다.

### 중단 후 재개 확인

Usage limit 중단 후 기존 workspace부터 확인했다. Git은 `main`에 커밋이 없는 상태이며 README와 docs는 모두 untracked다. 따라서 Git diff만으로 Turn별 변경을 구분할 수 없어, 현재 파일 내용과 중단 전 작업 기록을 대조했다.

재개 시점에 의미 정의, 6개 sample 평가, 전체 요약, ambiguity 기록, README 및 계약 연결은 이미 작성되어 있었다. 미완료 작업은 최종 validation과 Turn Result 보고였다. 기존 sample과 level은 보존하고 이 재개·검증 기록만 추가했다. 이전에 확인한 JD를 다시 수집하거나 평가를 새로 실행하지 않았다.

- 6개 sample 모두 company/position, source URL, reviewed date, requirement 요약, Why 및 Primary Gaps를 포함한다.
- Requirement-level 판단은 총 46개이며 sample별로 6 / 7 / 9 / 7 / 8 / 9개다. 모든 행에 요청된 6개 항목이 채워져 있고 Expected Match는 네 level 중 하나다.
- Sample별 overall과 두 문서의 요약 표가 일치한다. 사람 검토 전 초안 표시와 중요한 semantic ambiguity도 유지되어 있다.
- Repository에는 예상한 Markdown 문서 6개만 있다. 상대 링크, 줄 끝 공백, 마지막 개행과 코드블록 짝을 확인했다.
- Turn 000·001 파일의 SHA-256이 중단 전 확인값과 같아 과거 기록이 보존되었음을 확인했다.
- Baseline에는 percentage가 없다. Evaluation Contract의 기존 숫자 표시는 Turn 001의 설명용 예시이며 이번 평가 결과나 계산 규칙이 아니다. 이번 Turn에서 score, weighting formula, threshold를 정의하지 않았다.

## Deferred

- 사람의 baseline 검토·확정 및 중요 semantic ambiguity 해소
- Jev API 호출, AI evaluator 실행 및 baseline 비교 실험
- Match Score, percentage, weighting formula, threshold 및 calibration
- Candidate Profile JSON schema, TypeScript project, MCP server, database 및 application code
- 전체 JD 원문 보관과 별도 dataset infrastructure

## Suggested Next Step

Planning conversation에서 먼저 requirement-level 초안과 overall을 검토하고, PIA의 Partial / Weak 및 Nitori의 Weak / Unknown 경계, Prop Tech 필수 정보 누락과 Strong 요약의 관계를 논의한다. 그 뒤 사람의 판단과 수정 이유를 engine 결과를 보기 전에 기록하여 baseline 상태를 확정하는 것을 제안한다. 다음 Turn의 실험이나 구현은 시작하지 않았다.

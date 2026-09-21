# Turn 007 — Reframe Evaluation Philosophy

- 날짜: 2026-09-21
- 범위: documentation only. Application code, prompt, profile, JD 및 evaluation logic 변경 없음. API 호출 없음.

## Completed

Turn 000–006 기록과 현재 repository를 검토하고 evaluation philosophy를 active documentation에 명시했다. 기존부터 evidence/Unknown/사람의 최종 판단과 ground truth가 아니라는 설명은 있었지만, Expected Match와 baseline 확정이라는 표현은 정답으로 오해될 여지가 있었다. 이번에는 reference 일치가 개선 목표가 아님을 명시하고, 품질 기준과 구현 freeze 및 변경 조건을 확정했다.

## Changed Files

- [README](../../README.md): 현재 철학, workflow 가치, freeze 및 다음 phase 요약.
- [Evaluation Contract](../evaluation-contract.md#evaluation-philosophy--turn-007): 품질 기준 7개, reference/Overall/Jev 역할, 변경 조건과 전체 sample evaluation 방향.
- [Match Semantics](../match-semantics.md): 현재 해석 안내만 추가. 기존 정의·46개 requirement 행·6개 Overall·근거·당시 표현 보존.
- 본 `docs/development/turn-007.md` 생성. 과거 Turn 문서는 수정하지 않았다.

## Philosophy Change

Turn 000–003은 문제·scope·semantics와 사전 reference를 정리했다. Turn 004의 Prop Tech plus Reference Strong / Jev Partial을 계기로 Turn 005는 차이를 진단했고, Turn 006은 관찰을 위해 사후 Decision Trace를 추가했다. 이전 Turn에서 reference에 맞춰 tuning했다는 뜻은 아니다. 다만 이 차이를 없애려는 방향으로 계속 개선하면 특정 해석에 overfit할 위험이 있다.

앞으로 Overall의 객관적인 ground truth를 가정하지 않는다. Human Reference와 Jev 모두 정답이 아니며, evidence와 semantics 안에서 둘 다 합리적이면 disagreement로 기록하고 끝낼 수 있다. Overall은 빠른 비교를 돕는 summary signal이다. Requirement-level 근거와 미확인 정보가 더 중요할 수 있다. Label 일치율은 primary metric이 아니다.

Jev는 일반 LLM보다 정답을 잘 맞히는지 증명할 대상이 아니라 reusable workflow 안의 constrained, observable structured decision component다. Profile 재사용, explicit semantics, consistent machine-readable output, evidence tracing, missing/unknown visibility와 사람의 해석을 연결하는 것이 engineering value다. 현재 profile은 로컬 JSON으로 유지하며, DB나 result persistence를 새로 도입한 것이 아니다. MCP/automation 연계는 미래 가능성이다.

## Human Reference Role

현재 의미는 **Pre-evaluation Reference Interpretation**이다. Evaluator 실행 전에 중요하게 본 요소와 하나의 사전 해석을 보존하고 판단 차이·evaluation behavior를 관찰하는 비교 자료다. Target label 또는 ground truth가 아니다. Codex 작성 초안에서 출발해 Turn 003 사용자 지시로 비교 기록을 확정한 provenance를 유지한다.

| Sample | 보존한 Pre-evaluation Reference |
| --- | --- |
| BlueMeme / Micro Court | Strong |
| Prop Tech plus | Strong |
| PIA TECH LAB | Partial |
| CORE | Strong |
| teamLab | Partial |
| Nitori Digital Base | Weak |

과거 문서의 Human Baseline / Expected Match 및 “확정”은 역사적 표현으로 보존한다. Active semantics 문서 앞부분의 안내가 이 표현들의 현재 용도를 설명한다. Sample별 label을 runtime에 강제할 규칙으로 읽지 않는다.

## Evaluator Quality Criteria

| 기준 | 검토 내용 |
| --- | --- |
| Evidence Faithfulness | Profile에 없는 경험을 만들어내지 않는가? |
| Requirement Fidelity | Required/Preferred/Responsibility/Stack/Context 등 JD 의미를 보존하는가? |
| Unknown Preservation | No Evidence를 No Experience로 바꾸거나 미확인 사실을 추정하지 않는가? |
| Evidence Traceability | Match와 evidence, missing/unknown을 연결해 확인할 수 있는가? |
| Structural Consistency | 서로 다른 JD에서도 동일한 semantics와 output contract를 유지하는가? |
| Repeatability / Variability | 같은 입력의 결과 변동을 관찰하고 기록하는가? 완전한 deterministic 결과는 필수 조건이 아니다. |
| User Usefulness | 강점·부분 부합·약한 근거·알 수 없는 정보를 이해하는 데 도움이 되는가? |

이는 품질 검토 기준이며 지금 모두 만족했다는 인증이나 새로운 수치 metric이 아니다. 이미 관찰한 회계 Weak → Partial, Azure 환경 Partial → Unknown은 특성으로 기록하며, 통계적 반복 실험을 수행한 것으로 표현하지 않는다.

## Frozen Evaluation Baseline

**Turn 006 implementation을 다음 전체 sample evaluation 완료까지 freeze한다.** Candidate Profile, match semantics, requirement instructions, Overall instructions, decomposition logic, Decision Trace semantics, Overall logic과 현재 evidence mapping을 유지한다. 이번 active documentation 갱신은 철학·운영 기준의 명확화이며 frozen match semantics 또는 실행 prompt 수정이 아니다.

시작 시 `main`의 HEAD는 `d5145a7d4095d0df7081f4f7f98750143974ea3f`이나 이는 Turn 004 커밋이다. Turn 006 코드와 Turn 005/006 문서는 아직 미커밋 변경사항으로 존재한다. **Freeze 대상은 HEAD만이 아니라 해당 변경을 포함한 현재 작업 트리**다. Commit/tag 또는 파일 잠금 기능은 추가하지 않았다.

구현 식별을 위한 content fingerprint:

`845acdfb12ebb507598e88b9fd60ce1759525750c765609a3e1dd3ff2968f425`

산출 방식: `cmd/**/*.go`, `internal/**/*.go`, `internal/evaluator/instructions.md`, `go.mod`, `data/profile.json`, `samples/proptech-plus.txt`를 상대 경로 문자열 순으로 정렬한다. 각 파일의 SHA-256과 상대 경로를 `hash + 두 공백 + path + LF` 형식으로 결합한 UTF-8 manifest의 SHA-256이다. 테스트도 포함하며 documentation은 제외한다. 이는 정답 label hash가 아니라 미커밋 구현을 식별하는 기록이다.

## Change Conditions

Evaluation logic 변경은 **Contract Violation** 또는 **Systematic Problem**에 근거해 검토한다.

- Contract Violation: 없는 경력 생성, 필수/우대 역전, Unknown을 실제 경험 부재로 단정, core role mismatch 무시, stack/context의 명시 Required 승격, evidence/Match 모순 등 구체적인 계약 위반.
- Systematic Problem: 여러 sample에서 동일한 문제 패턴이 반복되는 경우. Preferred Unknown의 과도한 영향이나 context의 반복 Limiting은 검토 후보이며 반복 자체만으로 오류를 확정하지 않는다.

Freeze 기간의 예외는 명백한 execution bug 또는 contract violation이다. 예외 변경은 근거·범위·비교 가능성 영향을 기록해야 한다. Systematic 후보는 우선 전체 sample에서 관찰하고 다음 planning에서 변경 여부를 검토한다. 단일 reference 불일치는 변경 근거가 아니다.

## Problems / Findings

- 기존 문서도 ground truth를 부정했지만, “확정”, “Expected”와 sample별 “유지” 표현이 target으로 읽힐 수 있었다. 기존 기록을 재작성하지 않고 현재 해석 안내로 구분했다.
- Turn 006 L005 사업 context의 Limiting은 관찰이다. Required로 실제 승격됐다는 증거가 아니며 이것만으로 수정하지 않는다.
- Decision Trace는 별도 호출의 **structured post-hoc attribution**이다. 실제 causal trace가 아니고 implicit count weighting도 입증하지 않는다.
- README/Contract의 optional score 관련 내용은 과거 아이디어로 남아 있다. 이번 철학은 점수나 calibration을 채택하지 않는다.
- 기존 6개 JD는 문서에 기록되어 있지만 실행용 sample은 Prop Tech plus만 존재한다. 나머지 5개 입력 준비는 다음 evaluation Turn의 작업이다. JD 출처·검토일·의미와 현재 runtime의 언어/학력 제외 범위를 함께 확인하여 비교 차이를 기록해야 한다.
- 현재 Git diff에는 앞선 Turn의 코드 변경도 보인다. 이번 변경은 시작 시점 snapshot과 비교하여 네 문서로 한정했으며 이전 작업을 되돌리거나 commit하지 않았다.

## Validation

작업 시작 시 파일 hash를 보관하고, 완료 시 지정한 네 문서 외 기존 파일이 동일한지 확인했다. Turn 000–006과 코드·prompt·profile·sample을 포함한다. Match Semantics는 추가한 해석 안내를 제외한 전체 본문이 작업 전과 동일하여 46개 requirement 행, 6개 Overall 및 근거가 보존된다. 상대 링크·공백·Markdown 표와 코드블록을 확인했다. Application 변경이 없으므로 Go 테스트를 재실행하지 않았으며 Jev 및 별도 evaluator 호출도 없다.

## Deferred

남은 5개 JD 준비·실행, Prop Tech 재실행, 통계적 repeatability test, code/prompt/logic 변경, tuning/calibration, decomposition/mapping 변경, score/weighting/threshold, MCP/DB/UI/crawler는 수행하지 않았다.

## Suggested Next Step

현재 evaluator를 변경하지 않은 채 BlueMeme / Micro Court, PIA TECH LAB, CORE, teamLab, Nitori Digital Base를 평가하고 기존 Prop Tech plus와 함께 6개 전체의 behavior를 비교하는 evaluation Turn을 제안한다. Nitori에서는 Angular/TypeScript와 일반 리딩을 frontend Tech Lead 경력으로 합성하거나 core role mismatch를 무시하는지 검토한다. Reference Weak와 다른 label 자체를 실패로 삼지 않는다. 품질 기준별 관찰과 disagreement, 위반 후보를 구분하고 다음 Turn의 API 실행이나 코드 변경은 시작하지 않았다.

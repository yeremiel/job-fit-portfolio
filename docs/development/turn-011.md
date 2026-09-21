# Turn 011 — Final Evaluation Regression

- 날짜: 2026-09-21
- 목적: Turn 010 수정 후 전체 6개 sample에서 contract integrity 확인. Reference label 비교·calibration 목적 아님.
- 시작 상태: clean `main`, HEAD `8f9ec89aca2e1026eaa61a4bc5c9097ff0088b82`.
- 범위: 동일 evaluator/profile/policy/semantics/JD로 각각 1회. 모두 exit 0. 실패·retry·추가 실행 없음.

## Completed

기존 regression tests/vet/build 후 여섯 JD를 새로 실행했다. Prop Tech plus도 이번에는 이전 결과 재사용이 아닌 Turn 011 신규 1회다. 모든 결과를 저장하고 requirement 원문·source type·evidence·missing/unknown·trace를 검토했다. Application code, profile, semantics, prompt, JD input 및 과거 기록은 변경하지 않았다.

## Regression Results

S/P/W/U는 Strong/Partial/Weak/Unknown 개수이며 Match 점수가 아니다. Human Reference 일치율은 계산하지 않는다.

| Sample | Overall | Requirements | S / P / W / U |
| --- | --- | --- | --- |
| BlueMeme / Micro Court | Partial | 7 | 2 / 3 / 1 / 1 |
| Prop Tech plus | Partial | 15 | 3 / 8 / 0 / 4 |
| PIA TECH LAB | Weak | 10 | 3 / 3 / 1 / 3 |
| CORE | Strong | 8 | 7 / 1 / 0 / 0 |
| teamLab | Partial | 8 | 2 / 5 / 0 / 1 |
| Nitori Digital Base | Weak | 11 | 3 / 2 / 3 / 3 |

## Stored Results / Decision Trace

[Manifest](turn-011-results/manifest.json)은 구현·입력·결과 hash와 실행 횟수/성공 여부를 기록한다. 구현 fingerprint는 경로순 정렬한 파일별 SHA-256 manifest의 SHA-256이다. 구성 파일 목록을 함께 보존하므로 이전 Turn의 fingerprint 범위와 혼동하지 않는다.

결과는 CLI stdout의 원래 structured domain JSON이며 provider 원시 응답은 아니다. 기존 artifact 구조를 재사용했고 persistence architecture는 추가하지 않았다. ID는 각 sample 내부에서만 유효하다. 아래 S/L/N은 Supporting/Limiting/Non-decisive다. 각 JSON에서 해당 requirement의 evidence와 missing/unknown을 찾을 수 있다.

| Sample / JSON | Supporting IDs | Limiting IDs | Non-decisive IDs |
| --- | --- | --- | --- |
| [BlueMeme / Micro Court](turn-011-results/bluememe.json) | L005, L008 | L002 | L006, L007, L009, L010 |
| [Prop Tech plus](turn-011-results/proptech-plus.json) | L002, L007, L011 | L005, L006 | L008, L009, L010, L012, L013, L014, L015, L016, L017, L018 |
| [PIA TECH LAB](turn-011-results/pia-tech-lab.json) | L002, L006 | L005, L007, L010 | L008, L009, L011, L012, L013 |
| [CORE](turn-011-results/core.json) | L002, L005, L006, L008, L009, L010 | 없음 | L007, L011 |
| [teamLab](turn-011-results/teamlab.json) | L005, L010 | L002, L006 | L007, L008, L009, L011 |
| [Nitori Digital Base](turn-011-results/nitori.json) | L005, L007 | L002, L006, L008 | L009, L010, L011, L012, L013, L014 |

## Contract Integrity

검토 근거: [Evaluation Contract](../evaluation-contract.md), [Match Semantics](../match-semantics.md), [Turn 009 triage](turn-009.md), [Turn 010 수정](turn-010.md), 동일 profile/sample 및 여섯 새 JSON. 새 source 수집이나 외부 evaluator는 사용하지 않았다. OK는 이번 관찰에서 명시적 위반을 식별하지 못했다는 뜻이지 모든 가능한 입력의 정확성 보장이 아니다.

| Sample | Evidence Faithfulness | Requirement Fidelity | Unknown Preservation | Decision Trace Consistency |
| --- | --- | --- | --- | --- |
| BlueMeme | 원문 evidence 보존, OutSystems 경력 합성 없음 | category 보존; 도구 경험 불문과 role gap의 기존 Concern 유지 | 테스트/릴리스 Unknown, 도구 미확인 기록 | Supporting은 근거 있는 Strong, 직무명 Weak가 Limiting; 새 명시 모순 없음 |
| Prop Tech plus | API/Boot 직접 근거와 전이 근거 구분 | Java 8+/1년 명시 조건만 duration gap; Azure 우대 보존 | Azure 운영·CI/CD·태도 정보 Unknown | domain/명시 Java 조건 Limiting, 우대 Unknown Non-decisive |
| PIA | Basic AWS를 전문 운영 경험으로 생성하지 않음 | AWS 3년 Required, Java Preferred 보존 | 자격·백업·패치 Unknown | AWS/인프라 scope/운영 gap Limiting, Java 우대 Strong Non-decisive |
| CORE | 원문 evidence 보존; Public/Energy coverage는 기존 Ambiguous | category 보존, 새 qualification 없음 | Energy 미확인의 표현 범위는 여전히 unresolved | Overall Strong, Limiting 없음; lifecycle Partial은 Non-decisive |
| teamLab | mobile/특화 architecture 경험 생성 없음 | source type 보존; 언어 제외는 기존 scope | CI/CD Unknown, 기술/범위 gap 보존 | Partial 역할·architecture가 Limiting, 요구 분석/Backend Strong Supporting |
| Nitori | Frontend Tech Lead 경력 합성 없음 | 미명시 version/duration 생성 없음; Required/Stack 구분 유지 | EC/서비스/디자인 도구 Unknown | core frontend role Weak Limiting, 일반 리딩/육성 Supporting |

### Evidence Faithfulness

모든 Supporting Evidence ID/text가 canonical profile 원문과 일치한다. 새로운 경력 사실·수치·자격을 생성한 문장은 발견하지 못했다. Nitori L006에서는 일반 리딩 E4를 transferable로, frontend 관련 E3/E5를 limited로 사용하고 role-gap Weak를 유지한다. PIA L005는 E7의 basic AWS를 limited로 연결하고 Weak로 평가한다.

원문 일치가 모든 relation의 타당성을 보장하지는 않는다. CORE L010의 E6 direct와 전체 context Strong은 기존 Candidate 2의 ambiguity이며, 에너지 직접 경험이 검증되었다고 주장하지 않는다. 이 한계를 이유 없이 “해결됨” 또는 새로운 Confirmed Violation으로 바꾸지 않았다.

### Requirement Fidelity

모든 반환 행의 원문은 입력의 동일 L번호와 일치하고 Required/Preferred/Responsibility/Stack/Context category를 보존했다. 새 certification/scale/ownership/기간 수치 조건의 추가도 발견하지 못했다. `requested depth`와 `full scope`는 기존 정형 gap이며 별도의 특정 연수·규모·자격 충족을 주장하지 않는다. 그 정성적 해석이 유일한 정답이라는 의미는 아니다.

명시된 version/duration gap은 Prop Tech plus L006뿐이며 실제 입력에 Java 8+ 및 1년 조건이 있다. PIA의 AWS 3년은 실제 Required로 유지되고 주된 rationale은 large-gap이다. Duration 조건이 있어도 반드시 duration rationale을 선택해야 한다는 규칙은 없다. Nitori 기술 행에는 그런 조건이 없고 이번 출력에도 없다.

Metadata는 모두 제외되었다. teamLab L012의 business Japanese는 frozen runtime의 언어 제외 정책에 따라 평가에서 빠졌다. 이는 언어 조건 충족 판정이나 새 regression이 아니며 결과를 전체 지원 자격 평가로 읽을 수 없다는 한계다.

### Unknown Preservation

Unknown은 No Experience로 변환되지 않았다. PIA AWS 자격에는 E7 limited 연결이 있으나 Unknown을 유지한다. Nitori의 EC/디자인 도구, teamLab CI/CD 등에는 경험 존재를 추정하지 않는다. Partial/Weak의 missing/unknown도 실제 경험 부재가 아닌 profile 미입증으로 설명한다.

CORE L010은 여전히 empty missing/unknown과 Strong이지만, Context Strong이 모든 개별 domain의 직접 경력 충족을 의미하는지는 Turn 009에서 Ambiguous로 남겼다. 이번 출력도 같은 해석 한계다. 다른 sample에서 에너지 같은 미입증 sub-domain을 명시적으로 보유했다고 생성한 반복 사례는 발견하지 못했다.

### Decision Trace Consistency

모든 requirement는 정확히 한 trace 역할에 포함되며 ID·Match·source type·원문이 일치한다. Unknown Supporting 또는 무근거 Supporting은 없다. 모든 Supporting은 이번에 Strong이고 직접 evidence가 연결된다. 이는 관찰 결과이며 모든 Supporting을 Strong으로 강제하는 새 규칙이 아니다.

CORE의 Overall Strong은 실제 live 결과이고 `limiting: []`다. L011 lifecycle의 Partial은 Non-decisive로 남아 있어 Strong일 때 모든 행을 Supporting으로 만들지 않는다. 다른 Overall의 Limiting은 gap을 가진 Partial/Weak이고, trace 정의와 직접 충돌하는 조합은 확인하지 못했다. Preferred Strong이 Non-decisive일 수 있으며 개별 Strong과 trace Supporting은 같은 개념이 아니다.

Trace는 **structured post-hoc attribution**, 즉 별도 요청의 사후 분류다. 선택된 Overall의 실제 causal reasoning을 관측한 것으로 표현하지 않는다.

## Turn 010 Regression Verification

### 1. Overall Strong + Limiting

**재발 없음.** 이번에는 CORE가 실제 Strong을 반환했으며 Limiting은 비어 있다. Turn 010 live CORE가 Partial이라 Strong 분기는 fixture/mock으로 검증했던 한계를 이번 실제 관측에서 보완했다. 기존 unit test의 Strong Choice exclusion과 domain rejection도 모두 통과했다. CORE의 L011 Partial과 missing scope를 지우거나 Overall을 낮춰 문제를 회피하지 않았다.

### 2. Synthetic version / duration / years

**재발 없음.** 6개 결과의 reasoning/missing/unknown을 검토했다. `specified version or duration`은 실제 조건이 있는 Prop Tech plus L006에만 나타난다. Nitori L010은 depth rationale이며 미명시 연수나 버전이 없다. E5 evidence 원문에 있는 “durations ... not documented”는 source text 복사이므로 새 JD qualification 생성과 구분한다.

이 결과는 미명시 조건 차단과 실제 조건 평가가 함께 유지됨을 보여준다. 모든 자연어 표현과 향후 입력에서의 완전한 grounding 보장은 아니며 lexical check의 알려진 한계는 그대로다.

## New Contract Violation Candidates

**새로운 후보 및 새로운 Confirmed Contract Violation 없음.** Reference와의 불일치, 이전 결과와의 변동은 failure로 세지 않았다. 기존 unresolved Concern을 새로운 확정 위반으로 승격하지 않았다.

## Known Ambiguities / Limitations

- **Composite Public / Energy Context:** CORE L010은 다시 Strong이고 E6만 direct 연결되었다. 여전히 Ambiguous다. 같은 CORE에서의 재관찰은 여러 sample에서 반복된 systematic issue의 증거가 아니다. 일반 복합행은 다른 JD에도 있으나 이번에 동일한 명시적 coverage 위반이 반복되었다고 확정할 수 없다.
- **직무 범위 해석:** BlueMeme는 도구 경험 불문임에도 Low-code 직무명이 Weak/Limiting이다. 기존 Concern이지만 OutSystems 사용 경력을 생성하거나 category를 필수로 바꾼 것은 아니며 현재 의미만으로 명백한 위반으로 확정하지 않는다. teamLab의 Web/mobile 역할 scope도 기존 해석 한계를 유지한다.
- **Variability:** CORE는 Turn 010 Partial에서 Strong으로, teamLab은 Turn 008 Weak에서 Partial로 변했다. Prop Tech Azure 환경은 이번 Partial이고 API 운영 trace는 Non-decisive다. 이 차이 때문에 재실행하거나 tuning하지 않았다. Turn 008과는 구현도 달라 변동 원인을 모델만으로 설명할 수 없고, 각 sample 1회로 통계적 안정성을 추정하지 않는다.
- **설명 해상도:** Scope/depth template와 post-hoc attribution은 상세한 인과 설명이 아니다. Profile 원문 복사만으로 연결의 의미적 정확성까지 검증되지 않는다.
- **검증 범위:** 기존 영문 paraphrased JD 6개에 대한 결과다. 전체 공개 원문, 모든 언어, 새로운 JD 및 모든 qualification 표현에 대한 보장은 아니다. 언어/학력 제외와 lexical eligibility의 false positive/negative 가능성은 유지된다.

## Tests / Validation

- `go test ./...`: 통과. 기존 Turn 010 regression 포함, unit test에서 외부 API 호출 없음.
- `go vet ./...`: 통과.
- `go build ./cmd/job-fit`: 통과.
- 실제 실행은 BlueMeme → Prop Tech plus → PIA → CORE → teamLab → Nitori 순서, 각 1회 성공. 재사용·재시도 없음. 키는 사용자 `.env`에서 자식 프로세스 환경변수로 전달하고 출력/저장하지 않았다.
- 각 CLI 성공은 현재 5단계 흐름이며 별도의 evaluator 호출이나 새 검증 API는 없다.
- 기존 코드·prompt·profile·6개 입력·semantics·Turn 기록·과거 artifact의 SHA-256이 시작 시점과 같음을 확인했다.
- 6개 JSON의 source line/category, evidence ID/text, level별 summary, trace 일대일 참조, Strong/Limiting 및 qualification gap을 확인했다.
- Manifest/result hash, 문서 상대 링크/공백 및 API 키 미포함 검증 통과.

## Changed Files

- `docs/development/turn-011.md` 생성.
- `docs/development/turn-011-results/`에 6개 원래 CLI JSON과 `manifest.json` 생성.

기존 파일은 수정하지 않았다. Commit, runtime persistence 기능 또는 새 architecture도 추가하지 않았다.

## Evaluation Phase Decision

**Evaluation Baseline Accepted**

이번 전체 sample regression에서 새로운 Confirmed Contract Violation을 발견하지 못했고 Turn 010의 두 위반이 재발하지 않았다. 명시된 completion 기준에 따라 evaluation subsystem의 iterative tuning을 종료한다. 이는 모든 semantic ambiguity가 해결되었거나 객관적인 Job Fit 정답을 검증했다는 뜻이 아니다. Known limitations는 위 기록대로 보존한다.

## Suggested Next Step

Evaluator tuning을 종료하고 다음 planning conversation에서 productization 단계의 범위를 결정하는 것을 제안한다. 이번 Turn에서는 다음 단계의 구현, MCP/UI/DB/crawler, calibration, 추가 반복 평가를 시작하지 않았다.

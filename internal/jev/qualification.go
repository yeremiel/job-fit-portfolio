package jev

import (
	"regexp"
	"strings"
)

// A conservative lexical eligibility check, not a JD parser or evidence inference.
// Only the current requirement is inspected: profile omissions and other rows
// must not create a version/duration qualification for this row.
var durationCondition = regexp.MustCompile(`(?i)\b(?:[0-9]+(?:\.[0-9]+)?|one|two|three|four|five|six|seven|eight|nine|ten|several)\s*(?:\+\s*)?(?:years?|months?|yrs?)\b|\b(?:years?|months?)\s+of\s+(?:relevant\s+)?experience\b|[0-9]+\s*(?:年間?|か月|ヶ月|개월|년)`)
var versionCondition = regexp.MustCompile(`(?i)\b(?:version\s*|v)[0-9]+(?:\.[0-9]+)*\b|\b[A-Za-z][A-Za-z0-9_.-]*\s+[0-9]+(?:\.[0-9]+)*(?:\+|\s+(?:or|and)\s+(?:later|newer|above|higher))|\b[A-Za-z][A-Za-z0-9_.-]*\s+[0-9]+\.[0-9]+\b|(?:バージョン|버전)\s*[0-9]+`)
var negatedQualification = regexp.MustCompile(`(?i)\b(?:no|without)\s+(?:minimum\s+|specific\s+|specified\s+)?(?:version|duration|years?\s+of\s+experience)\b|\b(?:version|duration|experience)\b[^;\n]*\b(?:not required|not mandatory|not specified|not documented|irrelevant)\b|(?:年数|期間|バージョン).{0,12}(?:不問|不要)|(?:기간|버전).{0,12}(?:무관|불필요)`)

func hasExplicitVersionOrDuration(text string) bool {
	// Keep a negated clause from admitting a qualification, without discarding
	// an independent affirmative clause in the same requirement.
	for _, clause := range strings.FieldsFunc(text, func(r rune) bool { return r == ';' || r == '\n' }) {
		if negatedQualification.MatchString(clause) {
			continue
		}
		if durationCondition.MatchString(clause) || versionCondition.MatchString(clause) {
			return true
		}
	}
	return false
}

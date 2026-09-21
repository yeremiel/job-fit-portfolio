package jev

import (
	"context"
	"job-fit/internal/evaluator"
)

// Trace runs after Overall. Its answers never feed back into the selected match.
func (c *Client) Trace(ctx context.Context, profile evaluator.Profile, job evaluator.Job, rows []evaluator.Assessment, overall evaluator.MatchLevel) ([]evaluator.InfluenceReference, error) {
	questions := map[string]question{}
	criteria := map[string]string{
		"Supporting":   "An important documented capability that supports the selected overall match in a higher direction. Unknown or absent usable evidence cannot support capability alignment.",
		"Limiting":     "An important gap or uncertainty that constrains a higher overall match. Missing information is not proof of no experience.",
		"Non-decisive": "Included in the evaluation but not a decisive supporting or limiting factor for the selected overall match.",
	}
	if overall == evaluator.Strong {
		delete(criteria, "Limiting")
	}
	for _, row := range rows {
		questions[row.ID] = choice("Classify this assessment's role in the already selected overall match, considering all assessments, source categories and the evaluation policy. This is post-hoc attribution, not access to the original decision's internal reasoning. Do not re-evaluate matches or infer new evidence. Do not classify mechanically by match, source category or number of rows. Do not assume every Preferred is non-decisive or every Required is limiting. Return one provided influence.", row, criteria)
	}
	answers, err := c.choices(ctx, map[string]any{"evaluationPolicy": evaluator.Policy, "profile": profile, "job": job.Text, "assessments": rows, "overallMatch": overall}, questions)
	if err != nil {
		return nil, err
	}
	refs := make([]evaluator.InfluenceReference, 0, len(rows))
	for _, row := range rows {
		refs = append(refs, evaluator.InfluenceReference{RequirementID: row.ID, Influence: evaluator.Influence(answers[row.ID])})
	}
	return refs, nil
}

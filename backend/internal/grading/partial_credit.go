package grading

import (
	"math"

	"harama/internal/domain"
)

type PartialCreditEngine struct{}

func NewPartialCreditEngine() *PartialCreditEngine {
	return &PartialCreditEngine{}
}

// CalculateScore computes the score based on the rubric and the criteria/rules identified as met.
// It effectively "enforces" the rubric's point values, correcting any arithmetic errors from the AI.
func (e *PartialCreditEngine) CalculateScore(rubric domain.Rubric, criteriaMet []string) (float64, []string) {
	totalScore := 0.0
	appliedRules := []string{}
	metSet := make(map[string]bool)

	for _, id := range criteriaMet {
		metSet[id] = true
	}

	// 1. Check Full Credit Criteria
	for _, criterion := range rubric.FullCreditCriteria {
		// Match by ID (preferred) or Description (fallback)
		if metSet[criterion.ID] || (criterion.Description != "" && metSet[criterion.Description]) {
			totalScore += criterion.Points
			appliedRules = append(appliedRules, criterion.ID)
		}
	}

	// 2. Check Partial Credit Rules
	for _, rule := range rubric.PartialCreditRules {
		// Check if the rule itself is marked as met by the AI
		if metSet[rule.ID] || (rule.Description != "" && metSet[rule.Description]) {
			// Check dependencies if any
			dependenciesMet := true
			for _, depID := range rule.Dependencies {
				if !metSet[depID] {
					dependenciesMet = false
					break
				}
			}

			if dependenciesMet {
				totalScore += rule.Points
				appliedRules = append(appliedRules, rule.ID)
			}
		}
	}

	// 3. Apply Common Mistakes penalties
	for _, mistake := range rubric.CommonMistakes {
		if metSet[mistake.ID] {
			totalScore -= mistake.Penalty
			appliedRules = append(appliedRules, mistake.ID)
		}
	}

	// 4. Clamp score
	// Assuming we can't go below 0
	totalScore = math.Max(0, totalScore)
	
	// We don't have MaxScore readily available in the rubric struct directly (it's usually on the Question),
	// but we should ensure we don't exceed it if we knew it.
	// For now, we trust the additive nature or rely on the caller to clamp to MaxScore.

	return totalScore, appliedRules
}

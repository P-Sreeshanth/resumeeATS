package models

// AnalysisResult represents the complete analysis result
type AnalysisResult struct {
	Score            float64            `json:"score"`
	SkillMatch       SkillMatchResult   `json:"skill_match"`
	ExperienceMatch  ExperienceResult   `json:"experience_match"`
	EducationMatch   EducationResult    `json:"education_match"`
	FormatScore      FormatResult       `json:"format_score"`
	MissingKeywords  []string           `json:"missing_keywords"`
	Suggestions      []string           `json:"suggestions"`
	MatchedKeywords  []string           `json:"matched_keywords"`
	ScoreBreakdown   ScoreBreakdown     `json:"score_breakdown"`
}

// SkillMatchResult contains skill matching details
type SkillMatchResult struct {
	Percentage      float64  `json:"percentage"`
	MatchedSkills   []string `json:"matched_skills"`
	MissingSkills   []string `json:"missing_skills"`
	TotalRequired   int      `json:"total_required"`
	TotalMatched    int      `json:"total_matched"`
}

// ExperienceResult contains experience matching details
type ExperienceResult struct {
	Score           float64 `json:"score"`
	YearsRequired   int     `json:"years_required"`
	YearsCandidate  float64 `json:"years_candidate"`
	MeetsRequirement bool   `json:"meets_requirement"`
}

// EducationResult contains education matching details
type EducationResult struct {
	Score       float64  `json:"score"`
	MatchedDegrees []string `json:"matched_degrees"`
	HasRequiredEducation bool `json:"has_required_education"`
}

// FormatResult contains ATS formatting analysis
type FormatResult struct {
	Score  float64  `json:"score"`
	Issues []string `json:"issues"`
	IsATSFriendly bool `json:"is_ats_friendly"`
}

// ScoreBreakdown shows how the final score was calculated
type ScoreBreakdown struct {
	SkillWeight      float64 `json:"skill_weight"`
	ExperienceWeight float64 `json:"experience_weight"`
	EducationWeight  float64 `json:"education_weight"`
	FormatWeight     float64 `json:"format_weight"`
	SkillScore       float64 `json:"skill_score"`
	ExperienceScore  float64 `json:"experience_score"`
	EducationScore   float64 `json:"education_score"`
	FormatScore      float64 `json:"format_score"`
}

// AnalysisRequest represents the request payload for analysis
type AnalysisRequest struct {
	JobDescription string `json:"job_description" binding:"required"`
}

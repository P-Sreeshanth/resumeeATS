package models

// JobDescription represents a parsed job description
type JobDescription struct {
	Title             string   `json:"title"`
	Company           string   `json:"company"`
	RequiredSkills    []string `json:"required_skills"`
	PreferredSkills   []string `json:"preferred_skills"`
	MinExperience     int      `json:"min_experience"`
	Education         []string `json:"education"`
	Location          string   `json:"location"`
	Description       string   `json:"description"`
	Keywords          []string `json:"keywords"`
	RawText           string   `json:"raw_text"`
}

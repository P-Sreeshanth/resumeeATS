package models

import "time"

// Resume represents the parsed resume data
type Resume struct {
	PersonalInfo PersonalInfo `json:"personal_info"`
	Education    []Education  `json:"education"`
	Experience   []Experience `json:"experience"`
	Skills       []string     `json:"skills"`
	Projects     []Project    `json:"projects"`
	Certifications []string   `json:"certifications"`
	RawText      string       `json:"raw_text"`
	FormatIssues []string     `json:"format_issues"`
}

// PersonalInfo contains basic personal information
type PersonalInfo struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// Education represents educational background
type Education struct {
	Degree      string `json:"degree"`
	Institution string `json:"institution"`
	Year        int    `json:"year"`
	GPA         string `json:"gpa,omitempty"`
}

// Experience represents work experience
type Experience struct {
	Company     string    `json:"company"`
	Position    string    `json:"position"`
	StartDate   time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	Description string    `json:"description"`
	IsCurrent   bool      `json:"is_current"`
}

// Project represents a project
type Project struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Technologies []string `json:"technologies"`
}

// CalculateExperienceYears calculates total years of experience
func (r *Resume) CalculateExperienceYears() float64 {
	var totalYears float64
	now := time.Now()

	for _, exp := range r.Experience {
		endDate := now
		if exp.EndDate != nil {
			endDate = *exp.EndDate
		}
		
		duration := endDate.Sub(exp.StartDate)
		years := duration.Hours() / (24 * 365.25)
		totalYears += years
	}

	return totalYears
}

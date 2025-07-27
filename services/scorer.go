package services

import (
        "ats-analyzer/models"
        "ats-analyzer/utils"
        "strings"
)

// Scorer handles resume scoring and analysis
type Scorer struct {
        nlp *NLPService
}

// NewScorer creates a new scorer instance
func NewScorer() *Scorer {
        return &Scorer{
                nlp: NewNLPService(),
        }
}

// ScoringWeights defines the weights for different scoring components
type ScoringWeights struct {
        SkillWeight      float64
        ExperienceWeight float64
        EducationWeight  float64
        FormatWeight     float64
}

// DefaultWeights returns the default scoring weights
func DefaultWeights() ScoringWeights {
        return ScoringWeights{
                SkillWeight:      0.4,
                ExperienceWeight: 0.3,
                EducationWeight:  0.2,
                FormatWeight:     0.1,
        }
}

// AnalyzeResume performs comprehensive resume analysis
func (s *Scorer) AnalyzeResume(resume *models.Resume, jobDesc *models.JobDescription) *models.AnalysisResult {
        weights := DefaultWeights()

        // Calculate individual scores
        skillMatch := s.calculateSkillMatch(resume, jobDesc)
        experienceMatch := s.calculateExperienceMatch(resume, jobDesc)
        educationMatch := s.calculateEducationMatch(resume, jobDesc)
        formatScore := s.calculateFormatScore(resume)

        // Calculate overall score
        overallScore := (skillMatch.Percentage/100)*weights.SkillWeight +
                experienceMatch.Score*weights.ExperienceWeight +
                educationMatch.Score*weights.EducationWeight +
                formatScore.Score*weights.FormatWeight

        // Convert to 0-100 scale
        overallScore *= 100

        // Generate suggestions
        suggestions := s.generateSuggestions(resume, jobDesc, skillMatch, experienceMatch, educationMatch, formatScore)

        return &models.AnalysisResult{
                Score:           overallScore,
                SkillMatch:      skillMatch,
                ExperienceMatch: experienceMatch,
                EducationMatch:  educationMatch,
                FormatScore:     formatScore,
                MissingKeywords: skillMatch.MissingSkills,
                MatchedKeywords: skillMatch.MatchedSkills,
                Suggestions:     suggestions,
                ScoreBreakdown: models.ScoreBreakdown{
                        SkillWeight:      weights.SkillWeight,
                        ExperienceWeight: weights.ExperienceWeight,
                        EducationWeight:  weights.EducationWeight,
                        FormatWeight:     weights.FormatWeight,
                        SkillScore:       skillMatch.Percentage,
                        ExperienceScore:  experienceMatch.Score * 100,
                        EducationScore:   educationMatch.Score * 100,
                        FormatScore:      formatScore.Score * 100,
                },
        }
}

// calculateSkillMatch calculates skill matching score
func (s *Scorer) calculateSkillMatch(resume *models.Resume, jobDesc *models.JobDescription) models.SkillMatchResult {
        // Combine required and preferred skills
        allJobSkills := append(jobDesc.RequiredSkills, jobDesc.PreferredSkills...)
        allJobSkills = utils.RemoveDuplicates(allJobSkills)

        percentage, matched, missing := s.nlp.CalculateSkillMatch(resume.Skills, allJobSkills)

        return models.SkillMatchResult{
                Percentage:    percentage,
                MatchedSkills: matched,
                MissingSkills: missing,
                TotalRequired: len(allJobSkills),
                TotalMatched:  len(matched),
        }
}

// calculateExperienceMatch calculates experience matching score
func (s *Scorer) calculateExperienceMatch(resume *models.Resume, jobDesc *models.JobDescription) models.ExperienceResult {
        candidateYears := resume.CalculateExperienceYears()
        requiredYears := float64(jobDesc.MinExperience)

        var score float64
        meetsRequirement := candidateYears >= requiredYears

        if requiredYears == 0 {
                score = 1.0 // No experience requirement
        } else if candidateYears >= requiredYears {
                score = 1.0 // Meets or exceeds requirement
        } else {
                // Partial score based on how close they are
                score = candidateYears / requiredYears
                if score > 1.0 {
                        score = 1.0
                }
        }

        return models.ExperienceResult{
                Score:            score,
                YearsRequired:    jobDesc.MinExperience,
                YearsCandidate:   candidateYears,
                MeetsRequirement: meetsRequirement,
        }
}

// calculateEducationMatch calculates education matching score
func (s *Scorer) calculateEducationMatch(resume *models.Resume, jobDesc *models.JobDescription) models.EducationResult {
        if len(jobDesc.Education) == 0 {
                return models.EducationResult{
                        Score:                1.0, // No education requirement
                        HasRequiredEducation: true,
                }
        }

        var matchedDegrees []string
        hasMatch := false

        // Check if any resume education matches job requirements
        for _, resumeEd := range resume.Education {
                for _, requiredEd := range jobDesc.Education {
                        if s.educationMatches(resumeEd.Degree, requiredEd) {
                                matchedDegrees = append(matchedDegrees, resumeEd.Degree)
                                hasMatch = true
                        }
                }
        }

        score := 0.0
        if hasMatch {
                score = 1.0
        } else if len(resume.Education) > 0 {
                score = 0.5 // Has some education but not exact match
        }

        return models.EducationResult{
                Score:                score,
                MatchedDegrees:       matchedDegrees,
                HasRequiredEducation: hasMatch,
        }
}

// calculateFormatScore analyzes resume formatting for ATS compatibility
func (s *Scorer) calculateFormatScore(resume *models.Resume) models.FormatResult {
        issues := resume.FormatIssues
        
        // Additional format checks
        additionalIssues := s.analyzeAdditionalFormatIssues(resume)
        issues = append(issues, additionalIssues...)

        // Calculate score based on number of issues
        score := 1.0
        if len(issues) > 0 {
                // Reduce score by 0.2 for each issue, minimum 0.3
                score = 1.0 - float64(len(issues))*0.2
                if score < 0.3 {
                        score = 0.3
                }
        }

        isATSFriendly := len(issues) <= 1 // Allow for one minor issue

        return models.FormatResult{
                Score:         score,
                Issues:        issues,
                IsATSFriendly: isATSFriendly,
        }
}

// analyzeAdditionalFormatIssues performs additional format analysis
func (s *Scorer) analyzeAdditionalFormatIssues(resume *models.Resume) []string {
        var issues []string
        text := resume.RawText

        // Check for contact information
        if resume.PersonalInfo.Email == "" {
                issues = append(issues, "Missing email address")
        }
        if resume.PersonalInfo.Phone == "" {
                issues = append(issues, "Missing phone number")
        }

        // Check for section organization
        hasExperience := len(resume.Experience) > 0
        hasEducation := len(resume.Education) > 0
        hasSkills := len(resume.Skills) > 0

        if !hasExperience && !hasEducation {
                issues = append(issues, "Missing key sections (experience or education)")
        }
        if !hasSkills {
                issues = append(issues, "No skills section identified")
        }

        // Check for excessive length (heuristic)
        if len(strings.Split(text, " ")) > 1000 {
                issues = append(issues, "Resume may be too long (consider condensing)")
        }

        return issues
}

// generateSuggestions creates actionable suggestions for resume improvement
func (s *Scorer) generateSuggestions(resume *models.Resume, jobDesc *models.JobDescription, 
        skillMatch models.SkillMatchResult, experienceMatch models.ExperienceResult,
        educationMatch models.EducationResult, formatScore models.FormatResult) []string {
        
        var suggestions []string

        // Skill-related suggestions
        if skillMatch.Percentage < 50 {
                suggestions = append(suggestions, "Your skill match is low. Consider adding more relevant skills from the job description.")
                
                if len(skillMatch.MissingSkills) > 0 {
                        topMissing := skillMatch.MissingSkills
                        if len(topMissing) > 5 {
                                topMissing = topMissing[:5]
                        }
                        suggestions = append(suggestions, "Key missing skills: "+strings.Join(topMissing, ", "))
                }
        } else if skillMatch.Percentage < 75 {
                maxSkills := 3
                if len(skillMatch.MissingSkills) < maxSkills {
                        maxSkills = len(skillMatch.MissingSkills)
                }
                suggestions = append(suggestions, "Good skill match! Consider adding: "+strings.Join(skillMatch.MissingSkills[:maxSkills], ", "))
        }

        // Experience-related suggestions
        if !experienceMatch.MeetsRequirement {
                if experienceMatch.YearsCandidate < float64(experienceMatch.YearsRequired) {
                        suggestions = append(suggestions, "You may not meet the minimum experience requirement. Highlight relevant internships, projects, or transferable skills.")
                }
        }

        // Education-related suggestions
        if !educationMatch.HasRequiredEducation && len(jobDesc.Education) > 0 {
                suggestions = append(suggestions, "Consider highlighting relevant coursework, certifications, or continuing education if you don't have the preferred degree.")
        }

        // Format-related suggestions
        for _, issue := range formatScore.Issues {
                switch {
                case strings.Contains(issue, "table"):
                        suggestions = append(suggestions, "Avoid using tables - use bullet points and clear headings instead.")
                case strings.Contains(issue, "column"):
                        suggestions = append(suggestions, "Use a single-column layout for better ATS readability.")
                case strings.Contains(issue, "email"):
                        suggestions = append(suggestions, "Add your email address to the contact section.")
                case strings.Contains(issue, "phone"):
                        suggestions = append(suggestions, "Include your phone number in the contact information.")
                case strings.Contains(issue, "skills"):
                        suggestions = append(suggestions, "Add a clear skills section with relevant technical and soft skills.")
                case strings.Contains(issue, "too long"):
                        suggestions = append(suggestions, "Consider condensing your resume to 1-2 pages for better readability.")
                }
        }

        // General suggestions based on overall score
        overallScore := (skillMatch.Percentage/100)*0.4 + experienceMatch.Score*0.3 + educationMatch.Score*0.2 + formatScore.Score*0.1
        overallScore *= 100

        if overallScore < 60 {
                suggestions = append(suggestions, "Consider tailoring your resume more closely to this specific job description.")
        }

        // Add quantification suggestion
        hasQuantifiedResults := strings.Contains(strings.ToLower(resume.RawText), "%") || 
                strings.Contains(strings.ToLower(resume.RawText), "increased") ||
                strings.Contains(strings.ToLower(resume.RawText), "reduced") ||
                strings.Contains(strings.ToLower(resume.RawText), "improved")
        
        if !hasQuantifiedResults {
                suggestions = append(suggestions, "Add quantified achievements (e.g., 'Increased sales by 20%', 'Managed team of 5 people').")
        }

        return suggestions
}

// educationMatches checks if education levels match
func (s *Scorer) educationMatches(candidateEd, requiredEd string) bool {
        candidate := strings.ToLower(candidateEd)
        required := strings.ToLower(requiredEd)

        // Direct match
        if strings.Contains(candidate, required) || strings.Contains(required, candidate) {
                return true
        }

        // Common degree equivalents
        equivalents := map[string][]string{
                "bachelor": {"bs", "ba", "btech", "bsc", "bachelor's"},
                "master":   {"ms", "ma", "mtech", "msc", "master's", "mba"},
                "phd":      {"doctorate", "doctoral", "ph.d"},
        }

        for degree, aliases := range equivalents {
                candidateHasDegree := strings.Contains(candidate, degree)
                requiredHasDegree := strings.Contains(required, degree)

                for _, alias := range aliases {
                        candidateHasDegree = candidateHasDegree || strings.Contains(candidate, alias)
                        requiredHasDegree = requiredHasDegree || strings.Contains(required, alias)
                }

                if candidateHasDegree && requiredHasDegree {
                        return true
                }
        }

        return false
}



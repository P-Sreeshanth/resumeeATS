package services

import (
        "ats-analyzer/models"
        "ats-analyzer/utils"
        "fmt"
        "path/filepath"
        "regexp"
        "strconv"
        "strings"
        "time"

        "github.com/ledongthuc/pdf"
        "github.com/unidoc/unioffice/document"
)

// Parser handles document parsing
type Parser struct {
        nlp *NLPService
}

// NewParser creates a new parser instance
func NewParser() *Parser {
        return &Parser{
                nlp: NewNLPService(),
        }
}

// ParseResume parses a resume file and extracts structured data
func (p *Parser) ParseResume(filename string) (*models.Resume, error) {
        ext := strings.ToLower(filepath.Ext(filename))
        var text string
        var err error

        switch ext {
        case ".pdf":
                text, err = p.parsePDF(filename)
        case ".docx":
                text, err = p.parseDOCX(filename)
        default:
                return nil, fmt.Errorf("unsupported file format: %s", ext)
        }

        if err != nil {
                return nil, fmt.Errorf("failed to extract text: %v", err)
        }

        resume := &models.Resume{
                RawText: text,
        }

        // Extract structured data from text
        p.extractPersonalInfo(resume, text)
        p.extractEducation(resume, text)
        p.extractExperience(resume, text)
        p.extractSkills(resume, text)
        p.extractProjects(resume, text)
        p.extractCertifications(resume, text)
        p.analyzeFormat(resume, text)

        return resume, nil
}

// ParseJobDescription parses job description text
func (p *Parser) ParseJobDescription(text string) (*models.JobDescription, error) {
        jd := &models.JobDescription{
                RawText: text,
        }

        p.extractJDTitle(jd, text)
        p.extractJDCompany(jd, text)
        p.extractJDSkills(jd, text)
        p.extractJDExperience(jd, text)
        p.extractJDEducation(jd, text)
        p.extractJDLocation(jd, text)
        p.extractJDKeywords(jd, text)

        return jd, nil
}

// parsePDF extracts text from PDF file
func (p *Parser) parsePDF(filename string) (string, error) {
        file, reader, err := pdf.Open(filename)
        if err != nil {
                return "", err
        }
        defer file.Close()

        var text strings.Builder
        totalPages := reader.NumPage()

        for i := 1; i <= totalPages; i++ {
                page := reader.Page(i)
                if page.V.IsNull() {
                        continue
                }

                pageText, err := page.GetPlainText(nil)
                if err != nil {
                        continue
                }
                text.WriteString(pageText)
                text.WriteString("\n")
        }

        return text.String(), nil
}

// parseDOCX extracts text from DOCX file
func (p *Parser) parseDOCX(filename string) (string, error) {
        doc, err := document.Open(filename)
        if err != nil {
                return "", err
        }
        defer doc.Close()

        var text strings.Builder
        for _, para := range doc.Paragraphs() {
                for _, run := range para.Runs() {
                        text.WriteString(run.Text())
                }
                text.WriteString("\n")
        }

        return text.String(), nil
}

// extractPersonalInfo extracts personal information from resume text
func (p *Parser) extractPersonalInfo(resume *models.Resume, text string) {
        lines := strings.Split(text, "\n")
        
        // Email regex
        emailRegex := regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`)
        if email := emailRegex.FindString(text); email != "" {
                resume.PersonalInfo.Email = email
        }

        // Phone regex
        phoneRegex := regexp.MustCompile(`(\+?1?[-.\s]?)?\(?([0-9]{3})\)?[-.\s]?([0-9]{3})[-.\s]?([0-9]{4})`)
        if phone := phoneRegex.FindString(text); phone != "" {
                resume.PersonalInfo.Phone = phone
        }

        // Name is typically in the first few lines
        for i, line := range lines {
                if i > 5 { // Don't look beyond first few lines
                        break
                }
                cleanLine := strings.TrimSpace(line)
                if len(cleanLine) > 2 && len(cleanLine) < 50 && 
                   !strings.Contains(cleanLine, "@") && 
                   !phoneRegex.MatchString(cleanLine) {
                        // Simple name detection - could be improved
                        if regexp.MustCompile(`^[A-Za-z\s.]{2,}$`).MatchString(cleanLine) {
                                resume.PersonalInfo.Name = cleanLine
                                break
                        }
                }
        }
}

// extractEducation extracts education information
func (p *Parser) extractEducation(resume *models.Resume, text string) {
        degreeRegex := regexp.MustCompile(`(?i)(bachelor|master|phd|b\.?s\.?|m\.?s\.?|b\.?a\.?|m\.?a\.?|b\.?tech|m\.?tech|mba|diploma)`)
        yearRegex := regexp.MustCompile(`(19|20)\d{2}`)
        
        lines := strings.Split(text, "\n")
        
        for i, line := range lines {
                if degreeRegex.MatchString(line) {
                        education := models.Education{}
                        
                        // Extract degree
                        if match := degreeRegex.FindString(line); match != "" {
                                education.Degree = strings.TrimSpace(match)
                        }
                        
                        // Look for institution in current and next few lines
                        for j := i; j < len(lines) && j < i+3; j++ {
                                currentLine := strings.TrimSpace(lines[j])
                                if len(currentLine) > 5 && !degreeRegex.MatchString(currentLine) {
                                        education.Institution = currentLine
                                        break
                                }
                        }
                        
                        // Extract year
                        if match := yearRegex.FindString(line); match != "" {
                                if year, err := strconv.Atoi(match); err == nil {
                                        education.Year = year
                                }
                        }
                        
                        if education.Degree != "" {
                                resume.Education = append(resume.Education, education)
                        }
                }
        }
}

// extractExperience extracts work experience
func (p *Parser) extractExperience(resume *models.Resume, text string) {
        // Simple experience extraction - look for date patterns and company names
        dateRegex := regexp.MustCompile(`(?i)(jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\s+(19|20)\d{2}`)
        lines := strings.Split(text, "\n")
        
        for i, line := range lines {
                if dateRegex.MatchString(line) {
                        experience := models.Experience{}
                        
                        // Try to extract dates
                        dates := dateRegex.FindAllString(line, -1)
                        if len(dates) > 0 {
                                startDate, err := p.parseDate(dates[0])
                                if err == nil {
                                        experience.StartDate = startDate
                                }
                                
                                if len(dates) > 1 {
                                        endDate, err := p.parseDate(dates[1])
                                        if err == nil {
                                                experience.EndDate = &endDate
                                        }
                                }
                        }
                        
                        // Look for company and position in surrounding lines
                        startIdx := i - 2
                        if startIdx < 0 {
                                startIdx = 0
                        }
                        endIdx := i + 3
                        if endIdx > len(lines) {
                                endIdx = len(lines)
                        }
                        for j := startIdx; j < endIdx; j++ {
                                currentLine := strings.TrimSpace(lines[j])
                                if len(currentLine) > 2 && !dateRegex.MatchString(currentLine) {
                                        if experience.Company == "" {
                                                experience.Company = currentLine
                                        } else if experience.Position == "" {
                                                experience.Position = currentLine
                                        }
                                }
                        }
                        
                        if experience.Company != "" {
                                resume.Experience = append(resume.Experience, experience)
                        }
                }
        }
}

// extractSkills extracts skills from resume text
func (p *Parser) extractSkills(resume *models.Resume, text string) {
        // Common technical skills - this could be expanded with a larger dictionary
        skillKeywords := []string{
                "python", "java", "javascript", "typescript", "go", "golang", "rust", "c++", "c#",
                "react", "angular", "vue", "nodejs", "express", "django", "flask", "spring",
                "sql", "mysql", "postgresql", "mongodb", "redis", "elasticsearch",
                "aws", "azure", "gcp", "docker", "kubernetes", "terraform", "ansible",
                "git", "github", "gitlab", "jenkins", "ci/cd", "devops",
                "machine learning", "deep learning", "tensorflow", "pytorch", "scikit-learn",
                "html", "css", "bootstrap", "tailwind", "sass", "less",
        }
        
        textLower := strings.ToLower(text)
        var foundSkills []string
        
        for _, skill := range skillKeywords {
                if strings.Contains(textLower, strings.ToLower(skill)) {
                        foundSkills = append(foundSkills, skill)
                }
        }
        
        resume.Skills = utils.RemoveDuplicates(foundSkills)
}

// extractProjects extracts project information
func (p *Parser) extractProjects(resume *models.Resume, text string) {
        projectRegex := regexp.MustCompile(`(?i)(project|projects?)[\s:]*`)
        lines := strings.Split(text, "\n")
        
        for i, line := range lines {
                if projectRegex.MatchString(line) {
                        // Extract projects from next few lines
                        for j := i + 1; j < len(lines) && j < i + 10; j++ {
                                projectLine := strings.TrimSpace(lines[j])
                                if len(projectLine) > 10 {
                                        project := models.Project{
                                                Name:        projectLine,
                                                Description: projectLine,
                                        }
                                        resume.Projects = append(resume.Projects, project)
                                }
                        }
                        break
                }
        }
}

// extractCertifications extracts certifications
func (p *Parser) extractCertifications(resume *models.Resume, text string) {
        certRegex := regexp.MustCompile(`(?i)(certification|certified|certificate)`)
        lines := strings.Split(text, "\n")
        
        for _, line := range lines {
                if certRegex.MatchString(line) {
                        cleanLine := strings.TrimSpace(line)
                        if len(cleanLine) > 5 {
                                resume.Certifications = append(resume.Certifications, cleanLine)
                        }
                }
        }
}

// analyzeFormat analyzes resume formatting for ATS compatibility
func (p *Parser) analyzeFormat(resume *models.Resume, text string) {
        var issues []string
        
        // Check for tables (simple heuristic)
        if strings.Contains(text, "\t") || regexp.MustCompile(`\s{5,}`).MatchString(text) {
                issues = append(issues, "Document may contain tables or complex formatting")
        }
        
        // Check for special characters that might indicate formatting
        if regexp.MustCompile(`[│┌┐└┘├┤┬┴┼]`).MatchString(text) {
                issues = append(issues, "Document contains table borders or special formatting")
        }
        
        // Check for multiple columns (heuristic)
        lines := strings.Split(text, "\n")
        for _, line := range lines {
                if len(strings.Fields(line)) > 10 {
                        issues = append(issues, "Possible multi-column layout detected")
                        break
                }
        }
        
        resume.FormatIssues = issues
}

// Helper functions for job description parsing
func (p *Parser) extractJDTitle(jd *models.JobDescription, text string) {
        lines := strings.Split(text, "\n")
        for _, line := range lines {
                cleanLine := strings.TrimSpace(line)
                if len(cleanLine) > 5 && len(cleanLine) < 100 {
                        jd.Title = cleanLine
                        break
                }
        }
}

func (p *Parser) extractJDCompany(jd *models.JobDescription, text string) {
        // Simple company extraction - this could be improved
        companyRegex := regexp.MustCompile(`(?i)(company|organization|corp|inc|ltd)`)
        lines := strings.Split(text, "\n")
        
        for _, line := range lines {
                if companyRegex.MatchString(line) {
                        jd.Company = strings.TrimSpace(line)
                        break
                }
        }
}

func (p *Parser) extractJDSkills(jd *models.JobDescription, text string) {
        // Extract skills using similar logic as resume
        skillKeywords := []string{
                "python", "java", "javascript", "typescript", "go", "golang", "rust", "c++", "c#",
                "react", "angular", "vue", "nodejs", "express", "django", "flask", "spring",
                "sql", "mysql", "postgresql", "mongodb", "redis", "elasticsearch",
                "aws", "azure", "gcp", "docker", "kubernetes", "terraform", "ansible",
                "git", "github", "gitlab", "jenkins", "ci/cd", "devops",
                "machine learning", "deep learning", "tensorflow", "pytorch", "scikit-learn",
                "html", "css", "bootstrap", "tailwind", "sass", "less",
        }
        
        textLower := strings.ToLower(text)
        var requiredSkills []string
        
        for _, skill := range skillKeywords {
                if strings.Contains(textLower, strings.ToLower(skill)) {
                        requiredSkills = append(requiredSkills, skill)
                }
        }
        
        jd.RequiredSkills = utils.RemoveDuplicates(requiredSkills)
}

func (p *Parser) extractJDExperience(jd *models.JobDescription, text string) {
        expRegex := regexp.MustCompile(`(\d+)\s*\+?\s*year`)
        if match := expRegex.FindStringSubmatch(strings.ToLower(text)); len(match) > 1 {
                if years, err := strconv.Atoi(match[1]); err == nil {
                        jd.MinExperience = years
                }
        }
}

func (p *Parser) extractJDEducation(jd *models.JobDescription, text string) {
        degreeRegex := regexp.MustCompile(`(?i)(bachelor|master|phd|b\.?s\.?|m\.?s\.?|b\.?a\.?|m\.?a\.?|b\.?tech|m\.?tech|mba|diploma)`)
        matches := degreeRegex.FindAllString(text, -1)
        jd.Education = utils.RemoveDuplicates(matches)
}

func (p *Parser) extractJDLocation(jd *models.JobDescription, text string) {
        locationRegex := regexp.MustCompile(`(?i)(location|based in|located in)[\s:]*([a-zA-Z\s,]+)`)
        if match := locationRegex.FindStringSubmatch(text); len(match) > 2 {
                jd.Location = strings.TrimSpace(match[2])
        }
}

func (p *Parser) extractJDKeywords(jd *models.JobDescription, text string) {
        // Extract important keywords using TF-IDF
        jd.Keywords = p.nlp.ExtractKeywords(text, 20)
}

// parseDate parses date string to time.Time
func (p *Parser) parseDate(dateStr string) (time.Time, error) {
        formats := []string{
                "Jan 2006", "January 2006", "2006",
                "Jan 02, 2006", "January 02, 2006",
        }
        
        for _, format := range formats {
                if t, err := time.Parse(format, dateStr); err == nil {
                        return t, nil
                }
        }
        
        return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// Helper functions


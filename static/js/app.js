class ATSAnalyzer {
    constructor() {
        this.form = document.getElementById('analysisForm');
        this.loadingState = document.getElementById('loadingState');
        this.errorAlert = document.getElementById('errorAlert');
        this.resultsSection = document.getElementById('resultsSection');
        this.analyzeBtn = document.getElementById('analyzeBtn');
        
        this.scoreChart = null;
        this.breakdownChart = null;
        
        this.initializeEventListeners();
        this.setupFileUpload();
    }

    initializeEventListeners() {
        this.form.addEventListener('submit', (e) => this.handleFormSubmit(e));
        
        // File input change event
        document.getElementById('resumeFile').addEventListener('change', (e) => {
            this.validateFile(e.target.files[0]);
        });
    }

    setupFileUpload() {
        const fileInput = document.getElementById('resumeFile');
        const uploadArea = fileInput.parentElement;

        // Add drag and drop functionality
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            uploadArea.addEventListener(eventName, this.preventDefaults, false);
        });

        ['dragenter', 'dragover'].forEach(eventName => {
            uploadArea.addEventListener(eventName, () => uploadArea.classList.add('dragover'), false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            uploadArea.addEventListener(eventName, () => uploadArea.classList.remove('dragover'), false);
        });

        uploadArea.addEventListener('drop', (e) => {
            const files = e.dataTransfer.files;
            if (files.length > 0) {
                fileInput.files = files;
                this.validateFile(files[0]);
            }
        });
    }

    preventDefaults(e) {
        e.preventDefault();
        e.stopPropagation();
    }

    validateFile(file) {
        if (!file) return;

        const validTypes = ['application/pdf', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'];
        const maxSize = 10 * 1024 * 1024; // 10MB

        if (!validTypes.includes(file.type)) {
            this.showError('Please select a PDF or DOCX file.');
            return false;
        }

        if (file.size > maxSize) {
            this.showError('File size must be less than 10MB.');
            return false;
        }

        return true;
    }

    async handleFormSubmit(e) {
        e.preventDefault();
        
        const formData = new FormData(this.form);
        const file = formData.get('resume');
        const jobDescription = formData.get('job_description');

        // Validate inputs
        if (!file || file.size === 0) {
            this.showError('Please select a resume file.');
            return;
        }

        if (!jobDescription.trim()) {
            this.showError('Please enter a job description.');
            return;
        }

        if (jobDescription.trim().length < 50) {
            this.showError('Job description is too short. Please provide more details.');
            return;
        }

        if (!this.validateFile(file)) {
            return;
        }

        await this.analyzeResume(formData);
    }

    async analyzeResume(formData) {
        try {
            this.showLoading();
            this.hideError();
            this.hideResults();

            const response = await fetch('/api/v1/analyze', {
                method: 'POST',
                body: formData
            });

            const result = await response.json();

            if (!response.ok) {
                throw new Error(result.error || 'Analysis failed');
            }

            if (result.success && result.data) {
                this.displayResults(result.data);
            } else {
                throw new Error('Invalid response format');
            }

        } catch (error) {
            console.error('Analysis error:', error);
            this.showError(error.message || 'An error occurred during analysis. Please try again.');
        } finally {
            this.hideLoading();
        }
    }

    displayResults(analysis) {
        this.hideLoading();
        this.hideError();
        
        // Show results section with animation
        this.resultsSection.classList.remove('d-none');
        this.resultsSection.classList.add('animate-fade-in');
        
        // Display overall score
        this.displayOverallScore(analysis.score);
        
        // Display detailed analysis
        this.displaySkillMatch(analysis.skill_match);
        this.displayExperienceMatch(analysis.experience_match);
        this.displayEducationMatch(analysis.education_match);
        this.displayFormatScore(analysis.format_score);
        
        // Display suggestions
        this.displaySuggestions(analysis.suggestions);
        
        // Display score breakdown
        this.displayScoreBreakdown(analysis.score_breakdown);
        
        // Scroll to results
        this.resultsSection.scrollIntoView({ behavior: 'smooth' });
    }

    displayOverallScore(score) {
        const scoreElement = document.getElementById('overallScore');
        const descriptionElement = document.getElementById('scoreDescription');
        
        const roundedScore = Math.round(score);
        scoreElement.textContent = roundedScore;
        
        // Update score description and styling
        let description, className;
        if (roundedScore >= 80) {
            description = 'Excellent Match!';
            className = 'score-excellent';
        } else if (roundedScore >= 70) {
            description = 'Good Match';
            className = 'score-good';
        } else if (roundedScore >= 50) {
            description = 'Average Match';
            className = 'score-average';
        } else {
            description = 'Needs Improvement';
            className = 'score-poor';
        }
        
        descriptionElement.textContent = description;
        scoreElement.className = `display-3 fw-bold ${className}`;
        
        // Create score chart
        this.createScoreChart(roundedScore);
    }

    createScoreChart(score) {
        const ctx = document.getElementById('scoreChart').getContext('2d');
        
        if (this.scoreChart) {
            this.scoreChart.destroy();
        }
        
        this.scoreChart = new Chart(ctx, {
            type: 'doughnut',
            data: {
                datasets: [{
                    data: [score, 100 - score],
                    backgroundColor: [
                        this.getScoreColor(score),
                        '#e9ecef'
                    ],
                    borderWidth: 0
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                cutout: '70%',
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        enabled: false
                    }
                }
            }
        });
    }

    getScoreColor(score) {
        if (score >= 80) return '#198754';
        if (score >= 70) return '#20c997';
        if (score >= 50) return '#ffc107';
        return '#dc3545';
    }

    displaySkillMatch(skillMatch) {
        document.getElementById('skillPercentage').textContent = `${Math.round(skillMatch.percentage)}%`;
        document.getElementById('skillProgress').style.width = `${skillMatch.percentage}%`;
        document.getElementById('matchedSkillsCount').textContent = skillMatch.total_matched;
        document.getElementById('totalSkillsCount').textContent = skillMatch.total_required;
        
        // Display matched skills
        const matchedContainer = document.getElementById('matchedSkills');
        matchedContainer.innerHTML = '';
        if (skillMatch.matched_skills && skillMatch.matched_skills.length > 0) {
            skillMatch.matched_skills.forEach(skill => {
                const tag = document.createElement('span');
                tag.className = 'skill-tag';
                tag.textContent = skill;
                matchedContainer.appendChild(tag);
            });
        } else {
            matchedContainer.innerHTML = '<span class="text-muted">None</span>';
        }
        
        // Display missing skills
        const missingContainer = document.getElementById('missingSkills');
        missingContainer.innerHTML = '';
        if (skillMatch.missing_skills && skillMatch.missing_skills.length > 0) {
            skillMatch.missing_skills.slice(0, 10).forEach(skill => { // Limit to 10 skills
                const tag = document.createElement('span');
                tag.className = 'skill-tag missing';
                tag.textContent = skill;
                missingContainer.appendChild(tag);
            });
        } else {
            missingContainer.innerHTML = '<span class="text-muted">None</span>';
        }
    }

    displayExperienceMatch(experienceMatch) {
        document.getElementById('candidateYears').textContent = 
            experienceMatch.years_candidate ? `${experienceMatch.years_candidate.toFixed(1)} years` : 'N/A';
        document.getElementById('requiredYears').textContent = 
            experienceMatch.years_required ? `${experienceMatch.years_required} years` : 'Not specified';
        
        const statusElement = document.getElementById('experienceStatus');
        if (experienceMatch.meets_requirement) {
            statusElement.textContent = 'Meets Requirements';
            statusElement.className = 'badge bg-success fs-6';
        } else {
            statusElement.textContent = 'Below Requirements';
            statusElement.className = 'badge bg-warning fs-6';
        }
    }

    displayEducationMatch(educationMatch) {
        const statusElement = document.getElementById('educationStatus');
        const degreesContainer = document.getElementById('matchedDegrees');
        
        if (educationMatch.has_required_education) {
            statusElement.textContent = 'Education Matched';
            statusElement.className = 'badge bg-success fs-6';
        } else {
            statusElement.textContent = 'No Exact Match';
            statusElement.className = 'badge bg-secondary fs-6';
        }
        
        degreesContainer.innerHTML = '';
        if (educationMatch.matched_degrees && educationMatch.matched_degrees.length > 0) {
            educationMatch.matched_degrees.forEach(degree => {
                const tag = document.createElement('span');
                tag.className = 'skill-tag';
                tag.textContent = degree;
                degreesContainer.appendChild(tag);
            });
        }
    }

    displayFormatScore(formatScore) {
        const statusElement = document.getElementById('formatStatus');
        const issuesContainer = document.getElementById('formatIssues');
        
        if (formatScore.is_ats_friendly) {
            statusElement.textContent = 'ATS Friendly';
            statusElement.className = 'badge bg-success fs-6';
        } else {
            statusElement.textContent = 'Format Issues Found';
            statusElement.className = 'badge bg-warning fs-6';
        }
        
        issuesContainer.innerHTML = '';
        if (formatScore.issues && formatScore.issues.length > 0) {
            const issuesList = document.createElement('ul');
            issuesList.className = 'list-unstyled mb-0';
            
            formatScore.issues.forEach(issue => {
                const listItem = document.createElement('li');
                listItem.className = 'text-danger small mb-1';
                listItem.innerHTML = `<i class="fas fa-exclamation-triangle me-1"></i>${issue}`;
                issuesList.appendChild(listItem);
            });
            
            issuesContainer.appendChild(issuesList);
        } else {
            issuesContainer.innerHTML = '<span class="text-success small"><i class="fas fa-check me-1"></i>No formatting issues detected</span>';
        }
    }

    displaySuggestions(suggestions) {
        const container = document.getElementById('suggestions');
        container.innerHTML = '';
        
        if (suggestions && suggestions.length > 0) {
            suggestions.forEach((suggestion, index) => {
                const suggestionElement = document.createElement('div');
                suggestionElement.className = 'suggestion-item animate-slide-in';
                suggestionElement.style.animationDelay = `${index * 0.1}s`;
                suggestionElement.innerHTML = `
                    <i class="fas fa-lightbulb"></i>
                    ${suggestion}
                `;
                container.appendChild(suggestionElement);
            });
        } else {
            container.innerHTML = '<div class="text-center text-muted">No specific suggestions at this time.</div>';
        }
    }

    displayScoreBreakdown(breakdown) {
        if (!breakdown) return;
        
        document.getElementById('skillScoreBreakdown').textContent = `${breakdown.skill_score.toFixed(1)}%`;
        document.getElementById('experienceScoreBreakdown').textContent = `${breakdown.experience_score.toFixed(1)}%`;
        document.getElementById('educationScoreBreakdown').textContent = `${breakdown.education_score.toFixed(1)}%`;
        document.getElementById('formatScoreBreakdown').textContent = `${breakdown.format_score.toFixed(1)}%`;
        
        this.createBreakdownChart(breakdown);
    }

    createBreakdownChart(breakdown) {
        const ctx = document.getElementById('breakdownChart').getContext('2d');
        
        if (this.breakdownChart) {
            this.breakdownChart.destroy();
        }
        
        this.breakdownChart = new Chart(ctx, {
            type: 'bar',
            data: {
                labels: ['Skills', 'Experience', 'Education', 'Format'],
                datasets: [{
                    label: 'Score',
                    data: [
                        breakdown.skill_score,
                        breakdown.experience_score,
                        breakdown.education_score,
                        breakdown.format_score
                    ],
                    backgroundColor: [
                        '#0dcaf0',
                        '#ffc107',
                        '#6c757d',
                        '#212529'
                    ],
                    borderRadius: 4,
                    borderSkipped: false
                }]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        beginAtZero: true,
                        max: 100,
                        ticks: {
                            callback: function(value) {
                                return value + '%';
                            }
                        }
                    }
                },
                plugins: {
                    legend: {
                        display: false
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return `${context.label}: ${context.parsed.y.toFixed(1)}%`;
                            }
                        }
                    }
                }
            }
        });
    }

    showLoading() {
        this.loadingState.classList.remove('d-none');
        this.analyzeBtn.disabled = true;
        this.analyzeBtn.innerHTML = '<i class="fas fa-spinner fa-spin me-2"></i>Analyzing...';
    }

    hideLoading() {
        this.loadingState.classList.add('d-none');
        this.analyzeBtn.disabled = false;
        this.analyzeBtn.innerHTML = '<i class="fas fa-search me-2"></i>Analyze Resume';
    }

    showError(message) {
        document.getElementById('errorMessage').textContent = message;
        this.errorAlert.classList.remove('d-none');
        this.errorAlert.scrollIntoView({ behavior: 'smooth' });
    }

    hideError() {
        this.errorAlert.classList.add('d-none');
    }

    hideResults() {
        this.resultsSection.classList.add('d-none');
    }
}

// Initialize the application when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    new ATSAnalyzer();
});

// Add some utility functions for better UX
document.addEventListener('DOMContentLoaded', () => {
    // Add smooth scrolling for anchor links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            }
        });
    });

    // Add auto-resize to textarea
    const textarea = document.getElementById('jobDescription');
    if (textarea) {
        textarea.addEventListener('input', function() {
            this.style.height = 'auto';
            this.style.height = Math.min(this.scrollHeight, 300) + 'px';
        });
    }

    // Add keyboard shortcuts
    document.addEventListener('keydown', (e) => {
        // Ctrl/Cmd + Enter to submit form
        if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
            const form = document.getElementById('analysisForm');
            if (form) {
                form.dispatchEvent(new Event('submit'));
            }
        }
    });
});

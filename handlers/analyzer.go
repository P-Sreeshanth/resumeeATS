package handlers

import (
        "ats-analyzer/services"
        "ats-analyzer/utils"
        "fmt"
        "net/http"
        "path/filepath"

        "github.com/gin-gonic/gin"
        "github.com/sirupsen/logrus"
)

// AnalyzeResume handles the resume analysis request
func AnalyzeResume(c *gin.Context) {
        // Parse multipart form
        form, err := c.MultipartForm()
        if err != nil {
                logrus.Errorf("Failed to parse multipart form: %v", err)
                c.JSON(http.StatusBadRequest, gin.H{
                        "error": "Failed to parse form data",
                })
                return
        }

        // Get resume file
        files := form.File["resume"]
        if len(files) == 0 {
                c.JSON(http.StatusBadRequest, gin.H{
                        "error": "Resume file is required",
                })
                return
        }

        file := files[0]
        
        // Validate file
        if !utils.IsValidResumeFile(file.Filename) {
                c.JSON(http.StatusBadRequest, gin.H{
                        "error": "Invalid file format. Only PDF and DOCX files are supported",
                })
                return
        }

        // Get job description
        jobDescText := c.PostForm("job_description")
        if jobDescText == "" {
                c.JSON(http.StatusBadRequest, gin.H{
                        "error": "Job description is required",
                })
                return
        }

        // Save uploaded file temporarily
        filename := fmt.Sprintf("uploads/%d_%s", 
                utils.GenerateTimestamp(), 
                filepath.Base(file.Filename))
        
        if err := c.SaveUploadedFile(file, filename); err != nil {
                logrus.Errorf("Failed to save uploaded file: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{
                        "error": "Failed to save uploaded file",
                })
                return
        }

        // Parse resume
        parser := services.NewParser()
        resume, err := parser.ParseResume(filename)
        if err != nil {
                logrus.Errorf("Failed to parse resume: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{
                        "error": "Failed to parse resume: " + err.Error(),
                })
                return
        }

        // Parse job description
        jobDesc, err := parser.ParseJobDescription(jobDescText)
        if err != nil {
                logrus.Errorf("Failed to parse job description: %v", err)
                c.JSON(http.StatusInternalServerError, gin.H{
                        "error": "Failed to parse job description: " + err.Error(),
                })
                return
        }

        // Analyze and score
        scorer := services.NewScorer()
        analysis := scorer.AnalyzeResume(resume, jobDesc)

        // Clean up temporary file
        utils.CleanupFile(filename)

        logrus.Infof("Analysis completed with score: %.2f", analysis.Score)
        
        c.JSON(http.StatusOK, gin.H{
                "success": true,
                "data": analysis,
        })
}

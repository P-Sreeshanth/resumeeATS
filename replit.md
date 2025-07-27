# ATS Resume Analyzer

## Overview

This is a rule-based ATS (Applicant Tracking System) Resume Analyzer that evaluates how well resumes match job descriptions without using Large Language Models (LLMs). The system uses traditional NLP techniques, keyword matching, and scoring algorithms to provide instant feedback on resume-job fit, similar to how enterprise ATS systems like Taleo, Workday, and Greenhouse operate.

## User Preferences

Preferred communication style: Simple, everyday language.

## System Architecture

The application follows a traditional web application architecture with:

- **Frontend**: Static HTML/CSS/JavaScript with Bootstrap for UI components
- **Backend**: Python-based (Flask/FastAPI) for resume processing and analysis
- **Processing Engine**: Rule-based NLP using libraries like Spacy, NLTK, and scikit-learn
- **No Database Required**: Stateless processing with file-based input/output

## Key Components

### Frontend Layer
- **Static Web Interface**: Bootstrap-based responsive design with drag-and-drop file upload
- **Real-time Feedback**: JavaScript-driven UI updates and progress indicators
- **Chart Visualization**: Score displays and breakdown charts for analysis results
- **File Validation**: Client-side file type and size validation for PDF/DOCX files

### Processing Engine
- **Resume Parser**: Extracts structured data (name, email, experience, skills, education) from PDF/DOCX files
- **Job Description Parser**: Analyzes job requirements, required skills, and qualifications
- **Keyword Matching**: Uses TF-IDF and cosine similarity for skill matching
- **Scoring Algorithm**: Rule-based scoring system that calculates resume-job fit percentage
- **Suggestion Engine**: Generates actionable recommendations for resume improvement

### Core Libraries (No LLMs)
- **Document Processing**: pyresparser, PyPDF2, docx for file parsing
- **NLP Analysis**: Spacy, NLTK for text processing and entity extraction
- **Similarity Matching**: scikit-learn, FuzzyWuzzy for keyword and skill matching
- **Mathematical Operations**: TF-IDF vectorization and cosine similarity calculations

## Data Flow

1. **File Upload**: User uploads resume (PDF/DOCX) and pastes job description
2. **Document Parsing**: Extract structured data from both resume and job description
3. **Skill Extraction**: Identify and categorize skills, experience, and qualifications
4. **Matching Algorithm**: Calculate similarity scores using NLP techniques
5. **Score Generation**: Apply rule-based scoring logic to generate overall match percentage
6. **Recommendations**: Generate specific suggestions for resume improvement
7. **Results Display**: Present scores, breakdowns, and suggestions in the web interface

## External Dependencies

### Frontend Dependencies
- **Bootstrap 5.3.0**: UI framework for responsive design
- **Font Awesome 6.4.0**: Icon library for user interface elements
- **Chart.js** (implied): For score visualization and breakdown charts

### Backend Dependencies (from analysis)
- **pyresparser**: Resume parsing and data extraction
- **PyPDF2/pdfminer.six**: PDF document processing
- **python-docx**: DOCX document handling
- **Spacy**: Advanced NLP processing and entity recognition
- **NLTK**: Natural language processing utilities
- **scikit-learn**: Machine learning utilities for TF-IDF and similarity
- **FuzzyWuzzy**: Fuzzy string matching for skill variations
- **Flask/FastAPI**: Web framework for API endpoints

## Deployment Strategy

The application is designed for simple deployment with minimal infrastructure requirements:

- **Stateless Architecture**: No database required, processes files on-demand
- **Static Asset Serving**: CSS, JavaScript, and HTML files served directly
- **API Endpoints**: Backend serves analysis results via REST API
- **File Processing**: Temporary file handling for uploaded resumes
- **Scalable Design**: Can be deployed on various platforms (Replit, Heroku, AWS, etc.)

### Key Architectural Decisions

1. **No LLM Dependency**: Uses traditional NLP to avoid API costs and latency issues
2. **Rule-Based Scoring**: Deterministic results that are explainable and consistent
3. **Client-Side Validation**: Reduces server load by validating files before upload
4. **Bootstrap Framework**: Rapid UI development with mobile-responsive design
5. **Modular JavaScript**: Object-oriented frontend code for maintainability
6. **File-Based Processing**: Eliminates database complexity for this use case

The system prioritizes accuracy, speed, and cost-effectiveness over advanced AI capabilities, making it suitable for high-volume resume screening scenarios.
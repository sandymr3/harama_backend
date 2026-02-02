# HARaMA — Skills & Competencies Guide

> **Purpose:** This document breaks down the technical skills, domain knowledge, and implementation expertise needed to build HARaMA successfully. Organized by phase, problem domain, and difficulty level.

---

## Table of Contents

1. [Phase 1: Core Grading Engine Skills](#phase-1-core-grading-engine-skills)
2. [Phase 2: Multimodal Processing Skills](#phase-2-multimodal-processing-skills)
3. [Phase 3: Trust & Reliability Skills](#phase-3-trust--reliability-skills)
4. [Phase 4: Learning Systems Skills](#phase-4-learning-systems-skills)
5. [Phase 5: Platform Engineering Skills](#phase-5-platform-engineering-skills)
6. [Cross-Cutting Skills](#cross-cutting-skills)
7. [Gemini 3 Specific Skills](#gemini-3-specific-skills)
8. [Team Composition Guide](#team-composition-guide)

---

## Phase 1: Core Grading Engine Skills

### 1.1 Backend Development (Go)

**Skill Level Required:** Intermediate to Advanced

**Core Competencies:**

#### Go Language Fundamentals
- **Structs & Interfaces:** Domain modeling, repository patterns
- **Concurrency:** Goroutines for parallel OCR/grading, channels for job queues
- **Error Handling:** Idiomatic error wrapping, custom error types
- **Package Organization:** Clean architecture, internal vs public packages
- **Context Management:** Request timeouts, cancellation propagation

**Example Skills in Action:**
```go
// Repository Pattern
type ExamRepository interface {
    Create(ctx context.Context, exam domain.Exam) error
    GetByID(ctx context.Context, id uuid.UUID) (*domain.Exam, error)
}

// Concurrent Processing
func ProcessSubmissions(submissions []domain.Submission) {
    var wg sync.WaitGroup
    results := make(chan domain.GradingResult, len(submissions))
    
    for _, sub := range submissions {
        wg.Add(1)
        go func(s domain.Submission) {
            defer wg.Done()
            result := gradeSubmission(s)
            results <- result
        }(sub)
    }
    
    wg.Wait()
    close(results)
}
```

**Learning Resources:**
- Go by Example: https://gobyexample.com/
- Effective Go: https://go.dev/doc/effective_go
- "Let's Go" by Alex Edwards (book)

---

#### REST API Design
- **HTTP Routing:** Chi, Gorilla Mux, or Gin frameworks
- **Middleware:** Auth, logging, rate limiting, tenant isolation
- **Request Validation:** Struct tags, custom validators
- **Response Formatting:** Consistent JSON structure, error responses
- **Versioning:** `/api/v1/` patterns, backward compatibility

**API Design Principles:**
```go
// Good: Clear resource hierarchy
POST   /api/v1/exams
GET    /api/v1/exams/{id}
POST   /api/v1/exams/{id}/submissions
GET    /api/v1/submissions/{id}/grades

// Bad: Unclear nesting
POST   /api/v1/create-exam
GET    /api/v1/get-exam-by-id
```

**Skills Needed:**
- RESTful resource naming conventions
- HTTP status code selection (200, 201, 400, 401, 404, 500)
- Pagination strategies (cursor vs offset)
- CORS handling for frontend integration

---

#### Database Design & PostgreSQL
- **Schema Design:** Normalization, foreign keys, indexes
- **Migrations:** golang-migrate, version control for schema
- **Querying:** Raw SQL vs ORMs (GORM, sqlx)
- **Transactions:** ACID compliance for grade overrides
- **Performance:** Query optimization, EXPLAIN ANALYZE

**Critical Tables:**
```sql
-- Exams table with tenant isolation
CREATE TABLE exams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    subject VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    CONSTRAINT fk_tenant FOREIGN KEY (tenant_id) 
        REFERENCES tenants(id) ON DELETE CASCADE
);

-- Grades with immutable AI decision tracking
CREATE TABLE grades (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID NOT NULL,
    question_id UUID NOT NULL,
    ai_score DECIMAL(5,2),
    final_score DECIMAL(5,2) NOT NULL,
    confidence DECIMAL(3,2),
    status VARCHAR(20) NOT NULL,
    
    CONSTRAINT unique_submission_question 
        UNIQUE (submission_id, question_id)
);
```

**Skills Needed:**
- Database normalization (3NF understanding)
- Index strategy (B-tree, composite indexes)
- JSON column usage for flexible data
- Full-text search for rubric matching

---

### 1.2 AI Integration (Gemini 3)

**Skill Level Required:** Intermediate

**Core Competencies:**

#### Gemini API Fundamentals
- **Authentication:** API key management, secure storage
- **Request Structure:** Content parts, generation config
- **Response Parsing:** Handling streaming vs batch responses
- **Error Handling:** Rate limits, quota management, retries

**Basic Gemini Call:**
```go
import (
    "github.com/google/generative-ai-go/genai"
    "google.golang.org/api/option"
)

func callGemini(ctx context.Context, prompt string) (string, error) {
    client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return "", err
    }
    defer client.Close()
    
    model := client.GenerativeModel("gemini-3-pro")
    model.SetTemperature(0.2)
    
    resp, err := model.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        return "", err
    }
    
    return extractText(resp), nil
}
```

**Skills Needed:**
- SDK integration (Go client library)
- Async request handling
- Cost optimization (caching, batching)
- Fallback strategies for API failures

---

#### Prompt Engineering for Grading
- **System Prompts:** Role definition, behavior constraints
- **Few-Shot Learning:** Example answers for calibration
- **Chain-of-Thought:** Forcing step-by-step reasoning
- **Structured Outputs:** JSON schema enforcement
- **Context Management:** Rubric injection, answer formatting

**Grading Prompt Template:**
```
SYSTEM:
You are an expert educator grading student responses.
You must evaluate systematically using the provided rubric.
Think step-by-step and explain your reasoning clearly.

RUBRIC:
{
  "full_credit_criteria": ["Correct formula applied", "Units specified"],
  "partial_credit_rules": [
    {"condition": "Correct method, calculation error", "points": 3},
    {"condition": "Concept understood, formula wrong", "points": 2}
  ],
  "common_mistakes": ["Forgetting to convert units", "Sign error in final answer"]
}

STUDENT ANSWER:
"F = ma = (10 kg)(5 m/s²) = 50 N"

EVALUATE:
1. Check each criterion systematically
2. Identify mistakes if any
3. Assign partial credit if applicable
4. Provide reasoning for score

OUTPUT FORMAT (JSON):
{
  "score": <number>,
  "confidence": <0.0-1.0>,
  "reasoning": "<explanation>",
  "criteria_met": ["<criterion1>", ...],
  "mistakes_found": ["<mistake1>", ...]
}
```

**Skills Needed:**
- Understanding of educational rubrics
- JSON schema design for structured outputs
- Prompt iteration and testing
- Bias detection in AI responses

---

#### Response Validation & Safety
- **Schema Validation:** Ensuring AI returns expected structure
- **Confidence Calibration:** Mapping AI certainty to actionable thresholds
- **Hallucination Detection:** Cross-checking AI reasoning with rubric
- **Fallback Logic:** Human escalation when AI fails

**Validation Example:**
```go
type GradingResponse struct {
    Score          float64  `json:"score" validate:"min=0,max=10"`
    Confidence     float64  `json:"confidence" validate:"min=0,max=1"`
    Reasoning      string   `json:"reasoning" validate:"required,min=10"`
    CriteriaMet    []string `json:"criteria_met"`
    MistakesFound  []string `json:"mistakes_found"`
}

func validateResponse(resp GradingResponse, maxScore int) error {
    if resp.Score > float64(maxScore) {
        return errors.New("score exceeds maximum")
    }
    if resp.Confidence < 0.5 && len(resp.Reasoning) < 50 {
        return errors.New("low confidence requires detailed reasoning")
    }
    return nil
}
```

---

### 1.3 Document Processing (OCR)

**Skill Level Required:** Intermediate

**Core Competencies:**

#### Google Cloud Vision API
- **API Setup:** Project creation, service account credentials
- **Document Text Detection:** Batch vs online processing
- **Confidence Scoring:** Character-level confidence extraction
- **Bounding Boxes:** Spatial coordinate mapping
- **Language Detection:** Multi-language support

**OCR Implementation:**
```go
import (
    vision "cloud.google.com/go/vision/apiv1"
)

func extractTextWithConfidence(ctx context.Context, imageBytes []byte) (*OCRResult, error) {
    client, err := vision.NewImageAnnotatorClient(ctx)
    if err != nil {
        return nil, err
    }
    defer client.Close()
    
    image := vision.NewImageFromReader(bytes.NewReader(imageBytes))
    annotations, err := client.DetectDocumentText(ctx, image, nil)
    if err != nil {
        return nil, err
    }
    
    result := &OCRResult{
        FullText:   annotations.Text,
        Confidence: calculateConfidence(annotations.Pages),
        Words:      extractWords(annotations.Pages),
    }
    
    return result, nil
}
```

**Skills Needed:**
- Google Cloud Platform basics
- Image preprocessing (rotation, contrast)
- Multi-page PDF handling
- Confidence score interpretation

---

#### Tesseract Integration (Fallback)
- **Installation:** Docker image or system package
- **Configuration:** Language packs, PSM modes
- **Performance:** Speed vs accuracy tradeoffs
- **Cost Optimization:** When to use free vs paid OCR

**Tesseract Usage:**
```go
import "github.com/otiai10/gosseract/v2"

func tesseractOCR(imagePath string) (string, error) {
    client := gosseract.NewClient()
    defer client.Close()
    
    client.SetImage(imagePath)
    client.SetLanguage("eng")
    client.SetPageSegMode(gosseract.PSM_AUTO)
    
    text, err := client.Text()
    return text, err
}
```

---

#### Answer Segmentation Algorithms
- **Heuristic Detection:** Question number patterns (1., 2., Q1, etc.)
- **Spatial Analysis:** Using bounding boxes to detect boundaries
- **ML-Based:** Training models for complex layouts
- **Cross-Page Handling:** Answers spanning multiple pages

**Segmentation Logic:**
```go
func segmentAnswers(ocrResult OCRResult, template Exam) []AnswerSegment {
    segments := []AnswerSegment{}
    
    // Method 1: Pattern matching
    questionPattern := regexp.MustCompile(`^\d+\.`)
    
    // Method 2: Spatial analysis
    for i, word := range ocrResult.Words {
        if questionPattern.MatchString(word.Text) {
            segment := extractAnswerText(ocrResult.Words[i:])
            segments = append(segments, segment)
        }
    }
    
    // Method 3: Gemini verification for ambiguous cases
    if len(segments) != len(template.Questions) {
        segments = geminiVerifySegmentation(ocrResult, template)
    }
    
    return segments
}
```

**Skills Needed:**
- Regular expressions for pattern matching
- Geometric algorithms (bounding box overlap)
- Understanding of document layouts
- Edge case handling (handwritten margins)

---

### 1.4 Frontend Development (React/Next.js)

**Skill Level Required:** Intermediate

**Core Competencies:**

#### React Fundamentals
- **Hooks:** useState, useEffect, useContext, custom hooks
- **Component Design:** Separation of concerns, reusability
- **State Management:** Local state vs global (Context, Zustand)
- **Performance:** Memoization, lazy loading, code splitting

**Example Component:**
```tsx
'use client'

import { useState, useEffect } from 'react'
import { useGrading } from '@/hooks/useGrading'

export function GradingView({ submissionId }: { submissionId: string }) {
    const [overrideScore, setOverrideScore] = useState<number | null>(null)
    const { grading, loading, applyOverride } = useGrading(submissionId)
    
    if (loading) return <LoadingSpinner />
    
    return (
        <div className="grid grid-cols-2 gap-4">
            <AnswerDisplay answer={grading.answer} />
            <AIReasoningPanel 
                reasoning={grading.aiReasoning}
                confidence={grading.confidence}
            />
            <OverrideForm 
                currentScore={grading.score}
                onSubmit={(score) => applyOverride(score)}
            />
        </div>
    )
}
```

**Skills Needed:**
- TypeScript for type safety
- Tailwind CSS for styling
- Form handling and validation
- Error boundaries for crash prevention

---

#### API Integration
- **Fetch/Axios:** HTTP request management
- **Error Handling:** Network failures, 4xx/5xx responses
- **Loading States:** Skeleton screens, spinners
- **Optimistic Updates:** UI responsiveness

**API Client:**
```typescript
// lib/api.ts
export async function gradeSubmission(id: string) {
    const response = await fetch(`/api/v1/submissions/${id}/grade`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${getToken()}`,
            'Content-Type': 'application/json',
        },
    })
    
    if (!response.ok) {
        throw new Error(`Grading failed: ${response.statusText}`)
    }
    
    return response.json()
}
```

---

#### File Upload Handling
- **Drag & Drop:** React Dropzone
- **Progress Tracking:** Upload percentage, multi-file
- **Validation:** File type, size limits
- **Preview:** PDF rendering, image thumbnails

**Upload Component:**
```tsx
import { useDropzone } from 'react-dropzone'

export function UploadZone({ onUpload }: { onUpload: (files: File[]) => void }) {
    const { getRootProps, getInputProps, isDragActive } = useDropzone({
        accept: {
            'application/pdf': ['.pdf'],
            'image/*': ['.png', '.jpg', '.jpeg']
        },
        maxSize: 10 * 1024 * 1024, // 10MB
        onDrop: onUpload,
    })
    
    return (
        <div {...getRootProps()} className="border-2 border-dashed p-8">
            <input {...getInputProps()} />
            {isDragActive ? 
                <p>Drop files here...</p> : 
                <p>Drag & drop PDFs, or click to select</p>
            }
        </div>
    )
}
```

---

## Phase 2: Multimodal Processing Skills

### 2.1 Computer Vision

**Skill Level Required:** Intermediate to Advanced

**Core Competencies:**

#### Image Processing
- **OpenCV Basics:** Reading, resizing, color conversion
- **Preprocessing:** Noise reduction, contrast enhancement
- **Region Detection:** Edge detection (Canny), contour finding
- **Cropping:** Extracting diagram regions from submissions

**Diagram Extraction:**
```go
import "gocv.io/x/gocv"

func extractDiagram(imagePath string) ([]byte, error) {
    img := gocv.IMRead(imagePath, gocv.IMReadColor)
    defer img.Close()
    
    // Convert to grayscale
    gray := gocv.NewMat()
    defer gray.Close()
    gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
    
    // Apply Canny edge detection
    edges := gocv.NewMat()
    defer edges.Close()
    gocv.Canny(gray, &edges, 50, 150)
    
    // Find contours (diagram boundaries)
    contours := gocv.FindContours(edges, gocv.RetrievalExternal, gocv.ChainApproxSimple)
    
    // Extract largest non-text region
    diagramRegion := findLargestContour(contours)
    cropped := img.Region(diagramRegion)
    
    return gocv.IMEncode(".png", cropped)
}
```

**Skills Needed:**
- Basic image processing theory
- GoCV or similar library usage
- Bounding box calculations
- Image quality assessment

---

#### Gemini Vision API
- **Multimodal Inputs:** Combining text + images
- **Image Understanding:** Diagram type classification
- **Visual Question Answering:** "What is labeled in this diagram?"
- **OCR on Diagrams:** Reading labels, annotations

**Vision Grading:**
```go
func gradeWithDiagram(ctx context.Context, text string, diagramBytes []byte) (*GradingResult, error) {
    client, _ := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    defer client.Close()
    
    model := client.GenerativeModel("gemini-3-pro-vision")
    
    prompt := genai.Text(`
QUESTION: Draw a labeled diagram of a plant cell.

EVALUATE:
1. Are all major organelles present? (nucleus, mitochondria, chloroplast, etc.)
2. Are labels accurate and legible?
3. Is the structure scientifically correct?

STUDENT TEXT: ` + text)
    
    imagePart := genai.ImageData("image/png", diagramBytes)
    
    resp, err := model.GenerateContent(ctx, prompt, imagePart)
    if err != nil {
        return nil, err
    }
    
    return parseGradingResponse(resp), nil
}
```

---

### 2.2 Partial Credit Logic

**Skill Level Required:** Intermediate

**Core Competencies:**

#### Rule-Based Scoring
- **Criterion Decomposition:** Breaking rubrics into atomic checks
- **Dependency Graphs:** "Correct method" → "Execution score"
- **Weight Assignment:** Different criteria worth different points
- **Combination Logic:** AND vs OR conditions

**Partial Credit Engine:**
```go
type PartialCreditRule struct {
    ID          string
    Condition   func(answer AnswerSegment) bool
    Points      float64
    Description string
}

func calculatePartialCredit(answer AnswerSegment, rules []PartialCreditRule) float64 {
    totalScore := 0.0
    
    for _, rule := range rules {
        if rule.Condition(answer) {
            totalScore += rule.Points
            log.Printf("Rule %s: +%.1f points", rule.ID, rule.Points)
        }
    }
    
    return totalScore
}

// Example rules
var mathRules = []PartialCreditRule{
    {
        ID: "correct_formula",
        Condition: func(a AnswerSegment) bool {
            return strings.Contains(a.Text, "F = ma")
        },
        Points: 3.0,
        Description: "Correct formula identified",
    },
    {
        ID: "units_specified",
        Condition: func(a AnswerSegment) bool {
            return regexp.MustCompile(`\d+\s*(N|Newtons)`).MatchString(a.Text)
        },
        Points: 1.0,
        Description: "Proper units used",
    },
}
```

**Skills Needed:**
- Domain knowledge (math, science rubrics)
- Boolean logic and predicates
- Rubric analysis and decomposition
- Edge case identification

---

## Phase 3: Trust & Reliability Skills

### 3.1 Statistical Analysis

**Skill Level Required:** Intermediate

**Core Competencies:**

#### Variance & Disagreement Detection
- **Standard Deviation:** Measuring evaluator spread
- **Outlier Detection:** Identifying suspicious scores
- **Agreement Metrics:** Krippendorff's alpha, Cohen's kappa
- **Threshold Selection:** Balancing sensitivity vs specificity

**Variance Calculation:**
```go
func calculateVariance(scores []float64) float64 {
    if len(scores) < 2 {
        return 0.0
    }
    
    mean := calculateMean(scores)
    
    sumSquares := 0.0
    for _, score := range scores {
        sumSquares += math.Pow(score-mean, 2)
    }
    
    variance := sumSquares / float64(len(scores)-1)
    return math.Sqrt(variance) // Standard deviation
}

func shouldEscalate(variance float64, maxScore float64) bool {
    // Variance threshold as % of max score
    threshold := 0.15 * maxScore
    return variance > threshold
}
```

**Skills Needed:**
- Basic statistics (mean, variance, std dev)
- Probability theory (confidence intervals)
- Threshold tuning through experimentation
- Visualization of distributions

---

#### Confidence Calibration
- **Bayesian Thinking:** Prior + evidence → posterior confidence
- **Multi-Factor Scoring:** Combining OCR, rubric, historical data
- **Calibration Curves:** Plotting predicted vs actual confidence
- **Threshold Optimization:** ROC curves for decision boundaries

**Confidence Formula:**
```go
type ConfidenceFactors struct {
    OCRQuality          float64 // 0.0-1.0
    RubricClarity       float64 // How unambiguous the rubric is
    AnswerCompleteness  float64 // All required parts present
    HistoricalAccuracy  float64 // Past AI performance on similar
}

func calculateConfidence(factors ConfidenceFactors) float64 {
    weights := map[string]float64{
        "ocr":        0.25,
        "rubric":     0.35,
        "complete":   0.20,
        "historical": 0.20,
    }
    
    return (factors.OCRQuality * weights["ocr"]) +
           (factors.RubricClarity * weights["rubric"]) +
           (factors.AnswerCompleteness * weights["complete"]) +
           (factors.HistoricalAccuracy * weights["historical"])
}
```

---

### 3.2 Audit & Compliance

**Skill Level Required:** Intermediate

**Core Competencies:**

#### Immutable Logging
- **Append-Only Logs:** Never delete AI decisions
- **Event Sourcing:** Storing state changes, not just final state
- **Tamper Detection:** Cryptographic hashing of audit entries
- **Retention Policies:** GDPR compliance, data lifecycle

**Audit Table Design:**
```sql
CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    event_id UUID NOT NULL DEFAULT gen_random_uuid(),
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    entity_type VARCHAR(50) NOT NULL, -- 'grade', 'submission', etc.
    entity_id UUID NOT NULL,
    
    event_type VARCHAR(50) NOT NULL, -- 'created', 'updated', 'deleted'
    actor_id UUID, -- user or 'system'
    actor_type VARCHAR(20), -- 'teacher', 'student', 'ai'
    
    changes JSONB NOT NULL, -- {"before": {...}, "after": {...}}
    metadata JSONB, -- Additional context
    
    hash VARCHAR(64) NOT NULL -- SHA-256 of previous row + this row
);

-- Index for fast querying
CREATE INDEX idx_audit_entity ON audit_log(entity_type, entity_id, timestamp DESC);
CREATE INDEX idx_audit_actor ON audit_log(actor_id, timestamp DESC);
```

**Audit Logging Code:**
```go
func logGradeEvent(ctx context.Context, event AuditEvent) error {
    // Calculate hash chain for tamper detection
    previousHash := getLastAuditHash(ctx)
    currentData := fmt.Sprintf("%s|%v|%s", event.EventID, event.Timestamp, event.Changes)
    newHash := sha256.Sum256([]byte(previousHash + currentData))
    
    _, err := db.ExecContext(ctx, `
        INSERT INTO audit_log 
        (event_id, entity_type, entity_id, event_type, actor_id, changes, hash)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, event.EventID, event.EntityType, event.EntityID, event.EventType, 
       event.ActorID, event.Changes, hex.EncodeToString(newHash[:]))
    
    return err
}
```

**Skills Needed:**
- Database transaction management
- JSONB querying in PostgreSQL
- Cryptographic hashing (SHA-256)
- Compliance knowledge (FERPA, GDPR)

---

## Phase 4: Learning Systems Skills

### 4.1 Feedback Analysis

**Skill Level Required:** Advanced

**Core Competencies:**

#### Pattern Recognition in Overrides
- **Clustering:** Grouping similar teacher corrections
- **Trend Analysis:** Detecting systematic AI biases
- **Anomaly Detection:** Finding outlier feedback
- **Temporal Patterns:** Changes in accuracy over time

**Pattern Detection:**
```go
func analyzeOverridePatterns(ctx context.Context, questionID uuid.UUID) (*PatternReport, error) {
    // Fetch all overrides for this question
    overrides, err := repo.GetOverrides(ctx, questionID)
    if err != nil {
        return nil, err
    }
    
    // Group by delta direction
    tooLenient := []Override{}
    tooStrict := []Override{}
    
    for _, override := range overrides {
        delta := override.TeacherScore - override.AIScore
        if delta > 0 {
            tooLenient = append(tooLenient, override)
        } else if delta < 0 {
            tooStrict = append(tooStrict, override)
        }
    }
    
    // Analyze common reasons
    reasonClusters := clusterReasons(overrides)
    
    return &PatternReport{
        TotalOverrides:   len(overrides),
        SystematicBias:   calculateBias(overrides),
        CommonReasons:    reasonClusters,
        Recommendation:   generateRecommendation(reasonClusters),
    }, nil
}

func clusterReasons(overrides []Override) map[string]int {
    reasons := make(map[string]int)
    
    // Use Gemini to extract themes from teacher reasons
    for _, override := range overrides {
        theme := extractTheme(override.TeacherReason)
        reasons[theme]++
    }
    
    return reasons
}
```

**Skills Needed:**
- SQL aggregation queries
- Basic ML (k-means clustering)
- NLP for text analysis (Gemini-assisted)
- Data visualization

---

#### Adaptive Rubric Weighting
- **Dynamic Adjustment:** Updating criterion weights based on feedback
- **A/B Testing:** Comparing rubric variations
- **Version Control:** Tracking rubric changes over time
- **Rollback Mechanisms:** Reverting problematic changes

**Adaptive Weighting:**
```go
type CriterionWeight struct {
    CriterionID    string
    BaseWeight     float64
    AdjustedWeight float64
    Reason         string
    UpdatedAt      time.Time
}

func adaptRubricWeights(ctx context.Context, questionID uuid.UUID) error {
    patterns := analyzeOverridePatterns(ctx, questionID)
    
    for theme, frequency := range patterns.CommonReasons {
        if frequency > 10 { // Threshold for adaptation
            criterionID := mapThemeToCriterion(theme)
            
            // Increase weight for consistently missed criteria
            err := updateWeight(ctx, criterionID, +0.1, theme)
            if err != nil {
                return err
            }
            
            log.Printf("Increased weight for %s due to: %s", criterionID, theme)
        }
    }
    
    return nil
}
```

---

### 4.2 Personalized Feedback Generation

**Skill Level Required:** Intermediate

**Core Competencies:**

#### Student History Analysis
- **Performance Tracking:** Tracking improvement over time
- **Weakness Identification:** Recurring mistake patterns
- **Strength Recognition:** Consistently correct areas
- **Learning Velocity:** Rate of improvement metrics

**History Analysis:**
```go
func analyzeStudentHistory(ctx context.Context, studentID uuid.UUID) (*StudentProfile, error) {
    // Fetch last 10 submissions
    submissions, err := repo.GetStudentSubmissions(ctx, studentID, 10)
    if err != nil {
        return nil, err
    }
    
    profile := &StudentProfile{
        StudentID: studentID,
        Strengths: []string{},
        Weaknesses: []string{},
    }
    
    // Analyze by topic/skill
    topicScores := make(map[string][]float64)
    for _, sub := range submissions {
        for _, grade := range sub.Grades {
            topic := grade.Question.Topic
            topicScores[topic] = append(topicScores[topic], grade.Score)
        }
    }
    
    // Identify strengths (avg > 80%) and weaknesses (avg < 60%)
    for topic, scores := range topicScores {
        avg := calculateMean(scores)
        if avg > 0.8 {
            profile.Strengths = append(profile.Strengths, topic)
        } else if avg < 0.6 {
            profile.Weaknesses = append(profile.Weaknesses, topic)
        }
    }
    
    return profile, nil
}
```

---

#### Contextual Feedback Prompting
- **Tone Calibration:** Encouraging vs critical
- **Specificity:** Actionable advice, not generic praise
- **Resource Linking:** Suggesting specific study materials
- **Progress Highlighting:** Comparing to past performance

**Feedback Prompt:**
```go
func generatePersonalizedFeedback(grade GradingResult, history StudentProfile) string {
    prompt := fmt.Sprintf(`
ROLE: Supportive educator providing personalized feedback

STUDENT PERFORMANCE:
Current Answer Score: %d/%d
Reasoning: %s
Mistakes: %s

STUDENT HISTORY:
Strengths: %s
Weaknesses: %s
Recent Trend: %s

GENERATE FEEDBACK:
1. Acknowledge what they did well
2. Explain specific mistakes clearly
3. Provide actionable improvement steps
4. Reference their history (e.g., "You're improving in X")
5. Suggest specific resources if needed

TONE: Encouraging but honest
LENGTH: 3-4 sentences
`, grade.Score, grade.MaxScore, grade.Reasoning, 
   strings.Join(grade.MistakesFound, ", "),
   strings.Join(history.Strengths, ", "),
   strings.Join(history.Weaknesses, ", "),
   history.Trend)
    
    return callGemini(context.Background(), prompt)
}
```

---

## Phase 5: Platform Engineering Skills

### 5.1 Multi-Tenancy

**Skill Level Required:** Advanced

**Core Competencies:**

#### Tenant Isolation
- **Row-Level Security:** PostgreSQL policies
- **Data Partitioning:** Separate schemas vs shared tables
- **Query Filtering:** Automatic tenant_id injection
- **Resource Quotas:** Limiting per-tenant usage

**RLS Implementation:**
```sql
-- Enable RLS on all tables
ALTER TABLE exams ENABLE ROW LEVEL SECURITY;
ALTER TABLE submissions ENABLE ROW LEVEL SECURITY;

-- Policy: Users only see their tenant's data
CREATE POLICY tenant_isolation_policy ON exams
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant')::UUID);

-- Set tenant context in application
SET app.current_tenant = 'uuid-here';

-- All queries now automatically filtered
SELECT * FROM exams; -- Only returns exams for current tenant
```

**Application-Level Enforcement:**
```go
type TenantContext struct {
    TenantID uuid.UUID
}

func WithTenant(ctx context.Context, tenantID uuid.UUID) context.Context {
    return context.WithValue(ctx, tenantKey, TenantContext{TenantID: tenantID})
}

func GetTenantID(ctx context.Context) (uuid.UUID, error) {
    tenant, ok := ctx.Value(tenantKey).(TenantContext)
    if !ok {
        return uuid.Nil, errors.New("no tenant in context")
    }
    return tenant.TenantID, nil
}

// Middleware to set tenant from JWT
func TenantMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        claims := extractJWTClaims(r)
        ctx := WithTenant(r.Context(), claims.TenantID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**Skills Needed:**
- PostgreSQL advanced features
- Context propagation in Go
- Security best practices
- Database connection pooling per tenant

---

### 5.2 API Design & Versioning

**Skill Level Required:** Intermediate

**Core Competencies:**

#### API Versioning Strategies
- **URL Versioning:** `/api/v1/`, `/api/v2/`
- **Header Versioning:** `Accept: application/vnd.harrama.v1+json`
- **Backward Compatibility:** Deprecation policies
- **Documentation:** OpenAPI/Swagger specs

**Versioned Router:**
```go
func NewRouter() *chi.Mux {
    r := chi.NewRouter()
    
    // V1 routes
    r.Route("/api/v1", func(r chi.Router) {
        r.Use(AuthMiddleware)
        r.Use(TenantMiddleware)
        
        r.Post("/exams", handlers.CreateExam)
        r.Get("/exams/{id}", handlers.GetExam)
        r.Post("/exams/{id}/submissions", handlers.UploadSubmission)
    })
    
    // V2 routes (future)
    r.Route("/api/v2", func(r chi.Router) {
        // Enhanced features
        r.Post("/exams/{id}/submissions/batch", handlers.BatchUpload)
    })
    
    return r
}
```

---

#### Rate Limiting
- **Token Bucket:** Per-tenant request limits
- **Redis-Based:** Distributed rate limiting
- **Graceful Degradation:** Returning 429 with retry headers
- **Quota Management:** Daily/monthly limits

**Rate Limiter:**
```go
import "github.com/go-redis/redis_rate/v10"

func RateLimitMiddleware(limiter *redis_rate.Limiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            tenantID := GetTenantID(r.Context())
            
            res, err := limiter.Allow(r.Context(), fmt.Sprintf("tenant:%s", tenantID), redis_rate.PerMinute(60))
            if err != nil {
                http.Error(w, "Rate limit error", 500)
                return
            }
            
            if res.Allowed == 0 {
                w.Header().Set("Retry-After", fmt.Sprintf("%d", res.RetryAfter.Seconds()))
                http.Error(w, "Rate limit exceeded", 429)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
```

---

### 5.3 Observability

**Skill Level Required:** Intermediate

**Core Competencies:**

#### Structured Logging
- **JSON Logs:** Machine-parseable formats
- **Context Propagation:** Request IDs, tenant IDs
- **Log Levels:** DEBUG, INFO, WARN, ERROR
- **Log Aggregation:** ELK stack, Loki, CloudWatch

**Structured Logger:**
```go
import "github.com/rs/zerolog/log"

func gradeSubmission(ctx context.Context, id uuid.UUID) error {
    logger := log.With().
        Str("request_id", getRequestID(ctx)).
        Str("tenant_id", getTenantID(ctx).String()).
        Str("submission_id", id.String()).
        Logger()
    
    logger.Info().Msg("Starting grading")
    
    result, err := gradingEngine.Grade(ctx, id)
    if err != nil {
        logger.Error().Err(err).Msg("Grading failed")
        return err
    }
    
    logger.Info().
        Float64("score", result.Score).
        Float64("confidence", result.Confidence).
        Msg("Grading completed")
    
    return nil
}
```

---

#### Metrics & Monitoring
- **Prometheus:** Time-series metrics
- **Grafana:** Dashboards and alerting
- **Custom Metrics:** Grading latency, AI costs, override rates
- **Health Checks:** Liveness and readiness probes

**Metrics Collection:**
```go
import "github.com/prometheus/client_golang/prometheus"

var (
    gradingDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "harrama_grading_duration_seconds",
            Help: "Time spent grading submissions",
        },
        []string{"subject", "evaluator"},
    )
    
    overrideRate = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "harrama_overrides_total",
            Help: "Number of teacher overrides",
        },
        []string{"reason"},
    )
)

func init() {
    prometheus.MustRegister(gradingDuration)
    prometheus.MustRegister(overrideRate)
}

func gradeWithMetrics(subject string, evalID string) {
    timer := prometheus.NewTimer(gradingDuration.WithLabelValues(subject, evalID))
    defer timer.ObserveDuration()
    
    // Grading logic...
}
```

---

## Cross-Cutting Skills

### 6.1 DevOps & Deployment

**Skill Level Required:** Intermediate

**Core Competencies:**

#### Docker & Containerization
- **Dockerfile:** Multi-stage builds for Go
- **Docker Compose:** Local development orchestration
- **Image Optimization:** Layer caching, minimal base images
- **Health Checks:** Container-level monitoring

**Optimized Dockerfile:**
```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /harrama ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /harrama .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./harrama"]
```

---

#### CI/CD Pipelines
- **GitHub Actions:** Automated testing and deployment
- **Build Automation:** Linting, testing, building
- **Deployment Strategies:** Blue-green, canary releases
- **Rollback Procedures:** Quick recovery from bad deploys

**GitHub Actions Workflow:**
```yaml
name: Deploy

on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go test ./...
  
  deploy:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: flyctl deploy --remote-only
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}
```

---

### 6.2 Security

**Skill Level Required:** Intermediate

**Core Competencies:**

#### Authentication & Authorization
- **JWT Tokens:** Stateless auth
- **Role-Based Access Control (RBAC):** Teacher vs student permissions
- **API Key Management:** Secure storage, rotation
- **Password Hashing:** bcrypt, Argon2

**JWT Middleware:**
```go
import "github.com/golang-jwt/jwt/v5"

type Claims struct {
    UserID   uuid.UUID `json:"user_id"`
    TenantID uuid.UUID `json:"tenant_id"`
    Role     string    `json:"role"`
    jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := extractToken(r)
        
        claims := &Claims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
        if err != nil || !token.Valid {
            http.Error(w, "Unauthorized", 401)
            return
        }
        
        ctx := context.WithValue(r.Context(), "claims", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

---

#### Data Security
- **Encryption at Rest:** Database encryption
- **Encryption in Transit:** TLS/HTTPS
- **Sensitive Data Handling:** PII masking in logs
- **Key Management:** Environment variables, secret managers

**Secrets Management:**
```go
import "os"

type Config struct {
    DatabaseURL     string
    GeminiAPIKey    string // Never log this!
    JWTSecret       string
}

func LoadConfig() *Config {
    return &Config{
        DatabaseURL:  getEnv("DATABASE_URL", ""),
        GeminiAPIKey: getEnv("GEMINI_API_KEY", ""),
        JWTSecret:    getEnv("JWT_SECRET", ""),
    }
}

func getEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
}
```

---

## Gemini 3 Specific Skills

### 7.1 Advanced Prompting Techniques

**Skill Level Required:** Advanced

**Techniques to Master:**

#### Chain-of-Thought Reasoning
```
SYSTEM: Think step-by-step and show your work.

USER: Evaluate this physics answer...

GEMINI RESPONSE:
Step 1: Identify the formula used
- Student wrote F = ma ✓

Step 2: Check substitution
- Mass: 10 kg ✓
- Acceleration: 5 m/s² ✓

Step 3: Verify calculation
- 10 × 5 = 50 ✓

Step 4: Check units
- Result in Newtons ✓

CONCLUSION: Full credit (10/10)
```

---

#### Few-Shot Learning
```
SYSTEM: Here are examples of how to grade:

EXAMPLE 1:
Answer: "Photosynthesis makes food"
Score: 2/5
Reason: Correct concept but lacks detail (no mention of light, CO2, glucose)

EXAMPLE 2:
Answer: "6CO2 + 6H2O + light → C6H12O6 + 6O2"
Score: 5/5
Reason: Complete chemical equation with reactants and products

NOW GRADE:
Answer: "Plants use sunlight to create glucose from CO2 and water"
```

---

#### Structured Output Enforcement
```go
type GradingSchema struct {
    Type       string                 `json:"type"`
    Properties map[string]interface{} `json:"properties"`
    Required   []string               `json:"required"`
}

var gradingSchema = GradingSchema{
    Type: "object",
    Properties: map[string]interface{}{
        "score": map[string]interface{}{
            "type":    "number",
            "minimum": 0,
            "maximum": 10,
        },
        "confidence": map[string]interface{}{
            "type":    "number",
            "minimum": 0.0,
            "maximum": 1.0,
        },
        "reasoning": map[string]string{
            "type":      "string",
            "minLength": "20",
        },
    },
    Required: []string{"score", "confidence", "reasoning"},
}

// In Gemini call:
model.ResponseSchema = gradingSchema
```

---

### 7.2 Cost Optimization

**Techniques:**

#### Prompt Caching
```go
var rubricCache = make(map[uuid.UUID]string)

func buildGradingPrompt(rubric Rubric) string {
    if cached, ok := rubricCache[rubric.ID]; ok {
        return cached
    }
    
    // Expensive: Format rubric into prompt
    prompt := formatRubric(rubric)
    rubricCache[rubric.ID] = prompt
    
    return prompt
}
```

---

#### Batching Requests
```go
func gradeBatch(answers []AnswerSegment) []GradingResult {
    // Instead of N API calls, make 1 with all answers
    prompt := "Grade the following student answers:\n\n"
    for i, ans := range answers {
        prompt += fmt.Sprintf("ANSWER %d: %s\n\n", i+1, ans.Text)
    }
    
    response := callGemini(prompt)
    return parseMultipleResults(response)
}
```

---

## Team Composition Guide

### For a 4-Person Hackathon Team

**Recommended Roles:**

#### Role 1: Backend Lead (Go + AI)
- **Skills:** Go, PostgreSQL, Gemini API, OCR
- **Responsibilities:**
  - API server setup
  - Grading engine implementation
  - Gemini integration
  - Database schema design

**Time Allocation:**
- Day 1: API scaffold + DB setup
- Day 2: Grading engine + Gemini
- Day 3: Multi-evaluator logic
- Day 4: Testing + bug fixes

---

#### Role 2: Frontend Lead (React/Next.js)
- **Skills:** React, TypeScript, Tailwind, API integration
- **Responsibilities:**
  - Dashboard UI
  - Grading interface
  - File upload component
  - Real-time updates

**Time Allocation:**
- Day 1: Component library setup
- Day 2: Upload + exam creation UI
- Day 3: Grading review interface
- Day 4: Polish + responsive design

---

#### Role 3: AI/ML Specialist
- **Skills:** Prompt engineering, OCR, computer vision
- **Responsibilities:**
  - Prompt template design
  - OCR pipeline (Google Vision)
  - Diagram extraction
  - Confidence calibration

**Time Allocation:**
- Day 1: OCR integration
- Day 2: Prompt engineering
- Day 3: Multimodal grading
- Day 4: Quality assurance

---

#### Role 4: DevOps + Full-Stack
- **Skills:** Docker, deployment, testing, generalist coding
- **Responsibilities:**
  - Docker Compose setup
  - Deployment to Fly.io
  - Integration testing
  - Demo preparation

**Time Allocation:**
- Day 1: Local dev environment
- Day 2: Background workers
- Day 3: Deployment pipeline
- Day 4: Demo script + video

---

## Learning Path for Beginners

### If You're New to This Stack

**Week 1: Foundations**
- [ ] Go basics (Tour of Go)
- [ ] PostgreSQL tutorial
- [ ] React fundamentals
- [ ] REST API concepts

**Week 2: Integration**
- [ ] Build a simple CRUD API in Go
- [ ] Connect React to API
- [ ] Add PostgreSQL persistence
- [ ] Deploy to Fly.io

**Week 3: AI Integration**
- [ ] Gemini API quickstart
- [ ] Prompt engineering basics
- [ ] Build a simple AI grading prototype
- [ ] Test confidence scoring

**Week 4: Advanced Topics**
- [ ] Multi-evaluator pattern
- [ ] OCR integration
- [ ] Audit logging
- [ ] Performance optimization

---

## Recommended Tools & Libraries

### Backend (Go)
```
Chi Router:          github.com/go-chi/chi/v5
PostgreSQL Driver:   github.com/lib/pq
Gemini SDK:          github.com/google/generative-ai-go
Google Vision:       cloud.google.com/go/vision
Validator:           github.com/go-playground/validator/v10
Zerolog (logging):   github.com/rs/zerolog
UUID:                github.com/google/uuid
GORM (ORM):          gorm.io/gorm
```

### Frontend (React)
```
Next.js:             next
Tailwind CSS:        tailwindcss
React Dropzone:      react-dropzone
Recharts:            recharts
SWR (data fetching): swr
Zod (validation):    zod
React Hook Form:     react-hook-form
```

### Infrastructure
```
Docker:              docker.com
PostgreSQL:          postgresql.org
MinIO:               min.io
Qdrant:              qdrant.tech
Fly.io:              fly.io
```

---

## Success Metrics

### How to Know You've Mastered These Skills

**Phase 1 Complete:**
- [ ] Can build a REST API in Go from scratch
- [ ] Understand database normalization
- [ ] Can integrate Gemini API with proper prompting
- [ ] Built a working OCR → grading pipeline

**Phase 2 Complete:**
- [ ] Can process images with OpenCV/GoCV
- [ ] Implemented multimodal Gemini calls
- [ ] Built partial credit logic
- [ ] Handle diagrams in grading

**Phase 3 Complete:**
- [ ] Implemented multi-evaluator pattern
- [ ] Calculate variance and confidence
- [ ] Built audit trail system
- [ ] Understand escalation logic

**Phase 4 Complete:**
- [ ] Analyze teacher feedback patterns
- [ ] Adapt rubric weights dynamically
- [ ] Generate personalized student feedback
- [ ] Track learning progress over time

**Phase 5 Complete:**
- [ ] Implemented multi-tenant architecture
- [ ] Built versioned APIs
- [ ] Deployed to production
- [ ] Set up monitoring and logging

---

## Final Notes

**Key Takeaways:**

1. **Start Simple:** Phase 1 is 80% of the value
2. **Master Prompting:** Good prompts > complex code
3. **Trust is Everything:** Transparency beats accuracy
4. **Iterate Quickly:** Ship, test, learn, repeat

**When Stuck:**
- Read the error message carefully
- Check the documentation
- Search GitHub issues
- Ask Gemini for help (meta!)
- Pair program with teammates

**Remember:**
You don't need to be an expert in everything.
Focus on your strengths, collaborate, and ship.

---

## Additional Resources

### Documentation
- Go Docs: https://go.dev/doc/
- Gemini API: https://ai.google.dev/docs
- PostgreSQL: https://www.postgresql.org/docs/
- React: https://react.dev/

### Tutorials
- Go Web Apps: https://lets-go.alexedwards.net/
- Next.js Learn: https://nextjs.org/learn
- Prompt Engineering Guide: https://www.promptingguide.ai/

### Communities
- Go Forum: https://forum.golangbridge.org/
- r/golang: https://reddit.com/r/golang
- Gemini Discord: [Check Google AI Discord]
- Hackathon Slack: [Your event's channel]

---

**Good luck building HARaMA! 🚀**
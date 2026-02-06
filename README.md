# HARaMA — Complete Product Requirements Document (PRD)
## Gemini 3 Hackathon Edition

> **Mission:** Build an AI-powered exam grading system that demonstrates Gemini 3's reasoning capabilities while solving a real problem educators face daily.

---

## Executive Summary

**HARaMA** (Human-Assisted Reasoning & Automated Marking Assistant) is a production-grade exam grading platform that uses Gemini 3's advanced reasoning to provide transparent, multimodal, and trustworthy automated assessment.

**Core Value Proposition:**
- **For Teachers:** 70% time reduction on grading, with full transparency and control
- **For Students:** Detailed, personalized feedback that explains reasoning, not just scores
- **For Institutions:** Audit trails, consistency enforcement, and scalable assessment infrastructure

**Gemini 3 Differentiation:**
This system showcases Gemini 3's unique strengths:
- Long-context reasoning across entire exam papers
- Multimodal understanding (handwriting + diagrams + equations)
- Multi-step logical evaluation with chain-of-thought
- Reliable structured outputs for grading workflows

---

## Phase 1 — Core Reasoned Grading Engine

### 1.1 Exam & Rubric Ingestion

**Features:**
- Upload exam templates (PDF/images)
- Define questions with:
  - Question text
  - Point allocation
  - Answer type (short answer, essay, MCQ, diagram-based)
- Create detailed rubrics per question:
  - Scoring criteria (full/partial credit conditions)
  - Common mistakes database
  - Key concepts checklist
  - Example answers (good/bad)

**Gemini 3 Usage:**
```
Role: Rubric Interpreter
Input: Teacher's natural language rubric
Output: Structured scoring criteria with reasoning checkpoints
Technique: Few-shot learning with educational assessment examples
```

**Backend Implementation:**
```go
type Exam struct {
    ID          uuid.UUID
    Title       string
    Subject     string
    Questions   []Question
    CreatedAt   time.Time
    TenantID    uuid.UUID
}

type Question struct {
    ID              uuid.UUID
    QuestionText    string
    Points          int
    AnswerType      AnswerType
    Rubric          Rubric
    VisualAids      []string // image URLs
}

type Rubric struct {
    FullCreditCriteria    []Criterion
    PartialCreditRules    []ScoringRule
    CommonMistakes        []Mistake
    KeyConcepts           []string
}
```

**API Endpoints:**
```
POST   /api/v1/exams
GET    /api/v1/exams/{id}
POST   /api/v1/exams/{id}/questions
PUT    /api/v1/questions/{id}/rubric
```

---

### 1.2 Submission Upload & OCR Processing

**Features:**
- Batch upload student submissions (PDF/images)
- High-quality OCR with confidence scoring
- Automatic student identification (from header region)
- Submission deduplication
- Progress tracking dashboard

**OCR Strategy:**
- **Primary:** Google Cloud Vision API (high accuracy, integrates with Gemini)
- **Fallback:** Tesseract for cost optimization
- **Output:** Character-level confidence scores + bounding boxes

**Gemini 3 Usage:**
```
Role: OCR Verification & Correction
Input: Raw OCR text + original image region
Output: Corrected text + confidence adjustment
Technique: Vision + text comparison for ambiguous characters
```

**Data Pipeline:**
```go
type Submission struct {
    ID              uuid.UUID
    ExamID          uuid.UUID
    StudentID       string
    UploadedAt      time.Time
    ProcessingStatus ProcessingStatus
    OCRResults      []OCRResult
}

type OCRResult struct {
    PageNumber      int
    RawText         string
    Confidence      float64
    ImageURL        string
    BoundingBoxes   []BoundingBox
    CorrectedText   *string // Gemini-verified
}
```

**Processing Workflow:**
```
1. Upload → MinIO storage
2. Queue job → Worker pool
3. OCR extraction → Parallel per page
4. Confidence check → <0.85 triggers Gemini verification
5. Store results → PostgreSQL
6. Emit event → "OCR_COMPLETE"
```

---

### 1.3 Question-Level Segmentation

**Features:**
- Intelligent answer boundary detection
- Multi-page answer support
- Diagram/image association
- Overflow handling (answers in margins)

**Gemini 3 Usage:**
```
Role: Structural Understanding
Input: Full exam OCR + question template
Output: Answer segments with question mapping
Technique: Long-context analysis + spatial reasoning
```

**Algorithm:**
```go
func SegmentAnswers(ocr OCRResult, template Exam) []AnswerSegment {
    // 1. Use question numbers/patterns
    // 2. Leverage spatial layout (bounding boxes)
    // 3. Gemini verification for ambiguous cases
    // 4. Handle cross-page answers
}
```

---

### 1.4 Rubric-Aware AI Grading

**Features:**
- Multi-pass reasoning (not single-shot scoring)
- Explicit criterion checking
- Step-by-step evaluation logs
- Partial credit calculation
- Mistake identification

**Gemini 3 Grading Prompt Structure:**
```
SYSTEM:
You are an expert educator grading student responses.
Evaluate systematically using the provided rubric.
Think step-by-step and explain your reasoning.

RUBRIC:
{structured_rubric_json}

STUDENT ANSWER:
{ocr_text}

TASK:
1. Identify which criteria are met/unmet
2. Check for common mistakes
3. Assign partial credit where applicable
4. Provide reasoning for each point deduction
5. Output structured JSON with:
   - score (0-{max_points})
   - confidence (0.0-1.0)
   - reasoning (string)
   - criteria_met (array)
   - mistakes_found (array)
```

**Grading Engine:**
```go
type GradingResult struct {
    QuestionID      uuid.UUID
    Score           float64
    MaxScore        int
    Confidence      float64
    Reasoning       string
    CriteriaMet     []string
    MistakesFound   []string
    AIEvaluatorID   string
    Timestamp       time.Time
}

func GradeAnswer(ctx context.Context, answer AnswerSegment, rubric Rubric) GradingResult {
    // 1. Build Gemini prompt
    prompt := buildGradingPrompt(answer, rubric)
    
    // 2. Call Gemini with structured output
    response := callGemini(ctx, prompt, GradingResultSchema)
    
    // 3. Parse and validate
    result := parseGradingResponse(response)
    
    // 4. Store immutable AI decision
    storeAIDecision(result)
    
    return result
}
```

---

### 1.5 Confidence Scoring System

**Confidence Factors:**
1. **OCR Quality:** Character-level confidence average
2. **Rubric Match:** How clearly criteria are met/unmet
3. **Answer Completeness:** Presence of expected components
4. **Language Clarity:** Coherence and structure
5. **Historical Consistency:** Similarity to past grading patterns

**Confidence Formula:**
```go
func CalculateConfidence(result GradingResult, ocr OCRResult) float64 {
    ocrConf := ocr.Confidence
    rubricClarity := calculateRubricClarityScore(result)
    answerCompleteness := calculateCompletenessScore(result)
    
    weights := map[string]float64{
        "ocr": 0.3,
        "rubric": 0.4,
        "completeness": 0.3,
    }
    
    return (ocrConf * weights["ocr"]) +
           (rubricClarity * weights["rubric"]) +
           (answerCompleteness * weights["completeness"])
}
```

**Thresholds:**
- **High Confidence:** ≥ 0.85 → Auto-grade
- **Medium Confidence:** 0.70-0.84 → Flag for quick review
- **Low Confidence:** < 0.70 → Require human grading

---

### 1.6 Teacher Override & Audit System

**Features:**
- One-click accept/modify/reject AI grades
- Inline comment support
- Score adjustment with reason tracking
- Full audit trail (immutable logs)

**Override Workflow:**
```go
type GradeOverride struct {
    ID              uuid.UUID
    OriginalGrade   GradingResult
    NewScore        float64
    TeacherID       uuid.UUID
    Reason          string
    Timestamp       time.Time
    
    // Immutability: Original grade never deleted
}

func ApplyOverride(override GradeOverride) FinalGrade {
    // 1. Validate teacher authority
    // 2. Store override as separate record
    // 3. Update final grade view
    // 4. Emit feedback event (for Phase 4)
}
```

**Audit Log:**
```sql
CREATE TABLE grade_audit_log (
    id UUID PRIMARY KEY,
    submission_id UUID NOT NULL,
    question_id UUID NOT NULL,
    event_type VARCHAR(50), -- 'AI_GRADE', 'OVERRIDE', 'REVIEW'
    actor_id UUID, -- teacher or system
    score_before DECIMAL,
    score_after DECIMAL,
    reasoning TEXT,
    created_at TIMESTAMPTZ NOT NULL
);
```

---

### 1.7 Dashboard & Review Interface

**Teacher Dashboard:**
- Submission queue (sorted by confidence)
- Batch approval for high-confidence grades
- Red-flag list (low confidence / high variance)
- Grading progress metrics

**Review Interface:**
- Side-by-side: Original image + OCR text + AI reasoning
- Editable score with rubric checklist
- Student history (past submissions)
- Bulk actions (accept all, export)

---

## Phase 2 — Multimodal & Diagram Reasoning

### 2.1 Diagram Detection & Extraction

**Features:**
- Automatic diagram boundary detection
- High-res image extraction
- Diagram type classification (flowchart, graph, circuit, anatomy, etc.)
- OCR bypass for visual-only questions

**Gemini 3 Vision Usage:**
```
Role: Visual Understanding
Input: Cropped diagram image + question context
Output: Diagram description + key elements + evaluation
Technique: Vision model → structured representation
```

**Implementation:**
```go
type DiagramAsset struct {
    ID              uuid.UUID
    SubmissionID    uuid.UUID
    QuestionID      uuid.UUID
    ImageURL        string
    BoundingBox     BoundingBox
    DiagramType     string
    GeminiDescription string
}

func ExtractDiagrams(image []byte) []DiagramAsset {
    // 1. Use CV to detect non-text regions
    // 2. Crop and enhance
    // 3. Store in MinIO
    // 4. Gemini vision analysis
}
```

---

### 2.2 Multimodal Co-Reasoning

**Grading Strategy:**
- **Text-only questions:** Phase 1 approach
- **Diagram-only questions:** Pure vision evaluation
- **Hybrid questions:** Combined reasoning

**Gemini 3 Prompt for Diagrams:**
```
SYSTEM:
You are evaluating a student's diagram-based answer.
The question requires both visual accuracy and conceptual understanding.

QUESTION:
{question_text}

EXPECTED ELEMENTS:
{rubric_diagram_requirements}

STUDENT DIAGRAM:
[image data]

STUDENT EXPLANATION (if any):
{ocr_text}

EVALUATE:
1. Diagram accuracy (labels, arrows, structure)
2. Conceptual correctness
3. Integration with written explanation
4. Assign partial credit for:
   - Correct structure, wrong labels
   - Correct concept, poor execution
   - Incomplete but directionally correct

OUTPUT: JSON with score, visual_score, text_score, reasoning
```

---

### 2.3 Partial Credit Engine

**Granular Scoring:**
```go
type PartialCreditRule struct {
    Condition       string // "diagram_correct_but_unlabeled"
    PointsAwarded   float64
    Reasoning       string
}

func CalculatePartialCredit(result GradingResult) float64 {
    totalScore := 0.0
    
    for _, rule := range result.PartialCreditRules {
        if rule.ConditionMet {
            totalScore += rule.PointsAwarded
        }
    }
    
    return math.Min(totalScore, float64(result.MaxScore))
}
```

**Visualization:**
- Heatmap overlay on diagrams (correct/incorrect regions)
- Annotation markers (Gemini-identified mistakes)

---

## Phase 3 — Disagreement & Trust Layer

### 3.1 Multi-Evaluator Architecture

**Three Independent Graders:**

1. **Rubric Enforcer:**
   - Strict adherence to scoring criteria
   - No leniency for partial answers

2. **Reasoning Validator:**
   - Evaluates logical flow
   - Rewards conceptual understanding even if poorly expressed

3. **Structural Analyzer:**
   - Checks completeness, organization
   - Penalizes missing components

**Implementation:**
```go
type EvaluatorProfile struct {
    Name            string
    SystemPrompt    string
    WeightingBias   string // "strict", "lenient", "balanced"
}

func MultiEvaluatorGrade(answer AnswerSegment, rubric Rubric) MultiEvalResult {
    evaluators := []EvaluatorProfile{
        RubricEnforcer,
        ReasoningValidator,
        StructuralAnalyzer,
    }
    
    results := make([]GradingResult, len(evaluators))
    
    // Parallel evaluation
    var wg sync.WaitGroup
    for i, eval := range evaluators {
        wg.Add(1)
        go func(idx int, e EvaluatorProfile) {
            defer wg.Done()
            results[idx] = gradeWithProfile(answer, rubric, e)
        }(i, eval)
    }
    wg.Wait()
    
    return analyzeVariance(results)
}
```

---

### 3.2 Disagreement Detection

**Variance Calculation:**
```go
func CalculateVariance(results []GradingResult) float64 {
    scores := extractScores(results)
    mean := calculateMean(scores)
    
    variance := 0.0
    for _, score := range scores {
        variance += math.Pow(score - mean, 2)
    }
    
    return math.Sqrt(variance / float64(len(scores)))
}
```

**Thresholds:**
- **Low Variance:** < 5% of max score → High agreement
- **Medium Variance:** 5-15% → Flag for review
- **High Variance:** > 15% → Auto-escalate to human

---

### 3.3 Escalation System

**Escalation Rules:**
```go
func ShouldEscalate(multiEval MultiEvalResult) bool {
    return multiEval.Variance > 0.15 ||
           multiEval.LowestConfidence < 0.70 ||
           multiEval.ConflictingReasoning()
}

type EscalationCase struct {
    SubmissionID    uuid.UUID
    QuestionID      uuid.UUID
    AllEvaluations  []GradingResult
    Variance        float64
    EscalatedAt     time.Time
    AssignedTo      *uuid.UUID // teacher
    Status          string // 'pending', 'resolved'
}
```

**Teacher Notification:**
- Real-time alert (WebSocket)
- Email digest (batched)
- Priority queue (highest variance first)

---

### 3.4 Consensus Mechanism

**When Evaluators Agree:**
```go
func BuildConsensusGrade(results []GradingResult) FinalGrade {
    // Weighted average based on confidence
    totalWeight := 0.0
    weightedScore := 0.0
    
    for _, result := range results {
        totalWeight += result.Confidence
        weightedScore += result.Score * result.Confidence
    }
    
    return FinalGrade{
        Score: weightedScore / totalWeight,
        Confidence: calculateConsensusConfidence(results),
        Reasoning: mergeReasoning(results),
    }
}
```

---

## Phase 4 — Learning Feedback Loop

### 4.1 Feedback Capture System

**Override Delta Tracking:**
```go
type FeedbackEvent struct {
    ID              uuid.UUID
    QuestionID      uuid.UUID
    SubmissionID    uuid.UUID
    
    AIScore         float64
    TeacherScore    float64
    Delta           float64 // Teacher - AI
    
    AIReasoning     string
    TeacherReason   string
    
    Timestamp       time.Time
}

func CaptureOverrideFeedback(override GradeOverride) {
    delta := override.NewScore - override.OriginalGrade.Score
    
    event := FeedbackEvent{
        Delta: delta,
        // ... other fields
    }
    
    storeFeedback(event)
    analyzePattern(event) // Async
}
```

---

### 4.2 Pattern Recognition

**Mistake Clustering:**
```sql
-- Find common AI grading errors
SELECT 
    question_id,
    COUNT(*) as override_count,
    AVG(delta) as avg_correction,
    STRING_AGG(teacher_reason, '; ') as common_reasons
FROM feedback_events
WHERE delta > 1.0 -- AI was too lenient
GROUP BY question_id
HAVING COUNT(*) > 5
ORDER BY override_count DESC;
```

**Gemini 3 Usage:**
```
Role: Pattern Analyzer
Input: Aggregated feedback events
Output: Identified weak rubric areas + suggested improvements
Technique: Clustering + reasoning over teacher corrections
```

---

### 4.3 Adaptive Rubric Weighting

**Dynamic Criterion Adjustment:**
```go
type AdaptiveWeight struct {
    CriterionID     string
    BaseWeight      float64
    AdjustedWeight  float64
    Reason          string
    LastUpdated     time.Time
}

func AdaptRubric(questionID uuid.UUID, feedback []FeedbackEvent) {
    patterns := analyzePatterns(feedback)
    
    for _, pattern := range patterns {
        if pattern.Frequency > 10 {
            // AI consistently misses this criterion
            increaseWeight(pattern.CriterionID, 0.1)
            
            // Update system prompt emphasis
            updateEvaluatorPrompt(pattern.CriterionID, pattern.Reason)
        }
    }
}
```

---

### 4.4 Personalized Student Feedback

**Feedback Generation:**
```go
func GenerateStudentFeedback(result FinalGrade, history []Submission) string {
    prompt := fmt.Sprintf(`
ROLE: Educational feedback specialist

STUDENT ANSWER EVALUATION:
Score: %d/%d
Reasoning: %s

STUDENT HISTORY:
%s

GENERATE:
1. What they did well
2. Specific mistakes with explanations
3. How to improve (actionable steps)
4. Relevant resources/practice areas

Tone: Encouraging but honest
Length: 3-4 sentences
`, result.Score, result.MaxScore, result.Reasoning, summarizeHistory(history))
    
    return callGemini(prompt)
}
```

**Feedback Examples:**
- "Your diagram structure was excellent, but labels were missing. Review the terminology on page 47."
- "You're improving on multi-step problems (+15% from last exam). Focus on showing intermediate work."

---

## Phase 5 — Platformization & Extensibility

### 5.1 Multi-Tenant Architecture

**Tenant Isolation:**
```go
type Tenant struct {
    ID          uuid.UUID
    Name        string
    Domain      string
    Plan        string // 'free', 'school', 'district'
    Settings    TenantSettings
}

type TenantSettings struct {
    AIProvider          string // 'gemini', 'openai' (future)
    ConfidenceThreshold float64
    AutoGradeEnabled    bool
    RetentionDays       int
}
```

**Database Changes:**
```sql
-- Add tenant_id to all core tables
ALTER TABLE exams ADD COLUMN tenant_id UUID NOT NULL;
ALTER TABLE submissions ADD COLUMN tenant_id UUID NOT NULL;

-- Row-level security
CREATE POLICY tenant_isolation ON exams
    USING (tenant_id = current_setting('app.current_tenant')::UUID);
```

---

### 5.2 Public API Layer

**API Versioning:**
```
/api/v1/...  -- Current stable
/api/v2/...  -- Future features
```

**Key Endpoints:**
```
POST   /api/v1/tenants/{id}/exams
POST   /api/v1/exams/{id}/submissions
GET    /api/v1/submissions/{id}/grade
PUT    /api/v1/grades/{id}/override

POST   /api/v1/rubrics/validate
GET    /api/v1/analytics/grading-trends
```

**Authentication:**
- API keys (per tenant)
- OAuth 2.0 (for integrations)
- JWT tokens (user sessions)

---

### 5.3 Subject-Specific AI Profiles

**Profile Examples:**
```go
var SubjectProfiles = map[string]EvaluatorConfig{
    "mathematics": {
        FocusAreas: []string{"step-by-step reasoning", "formula accuracy"},
        PartialCreditRules: "generous for correct method, wrong arithmetic",
    },
    "english": {
        FocusAreas: []string{"argument structure", "evidence usage"},
        PartialCreditRules: "reward creativity and voice",
    },
    "science": {
        FocusAreas: []string{"diagram accuracy", "experimental design"},
        PartialCreditRules: "strict on safety violations",
    },
}
```

**Dynamic Prompt Loading:**
```go
func GetGradingPrompt(subject string, question Question) string {
    profile := SubjectProfiles[subject]
    basePrompt := loadTemplate("grading_base.txt")
    
    return fmt.Sprintf("%s\n\nSUBJECT FOCUS:\n%s", 
        basePrompt, 
        profile.FocusAreas)
}
```

---

### 5.4 Model Abstraction Layer

**Provider Interface:**
```go
type AIProvider interface {
    Grade(ctx context.Context, req GradingRequest) (GradingResult, error)
    AnalyzeImage(ctx context.Context, img []byte) (ImageAnalysis, error)
    GenerateFeedback(ctx context.Context, req FeedbackRequest) (string, error)
}

type GeminiProvider struct {
    client *genai.Client
    model  string
}

func (g *GeminiProvider) Grade(ctx context.Context, req GradingRequest) (GradingResult, error) {
    // Gemini-specific implementation
}

// Future: OpenAIProvider, ClaudeProvider, etc.
```

**Configuration:**
```yaml
ai:
  primary_provider: gemini
  gemini:
    model: gemini-3-flash-preview
    api_key: ${GEMINI_API_KEY}
    temperature: 0.2
    max_tokens: 2048
```

---

## Additional Hackathon-Winning Features

### Feature A: Live Grading Collaboration

**Real-Time Grading Sessions:**
- Multiple teachers grade same exam simultaneously
- Live disagreement resolution
- Consensus voting for edge cases

**Implementation:**
```go
type GradingSession struct {
    ID              uuid.UUID
    ExamID          uuid.UUID
    Participants    []uuid.UUID
    ActiveSubmission uuid.UUID
    Status          string // 'active', 'completed'
}

// WebSocket-based collaboration
type GradeVote struct {
    TeacherID   uuid.UUID
    Score       float64
    Reasoning   string
}

func ResolveByVoting(votes []GradeVote) FinalGrade {
    // Median score + merged reasoning
}
```

**UI Component:**
- Split-screen grading interface
- Live cursor tracking
- Comment threads per question

---

### Feature B: Handwriting Style Profiling

**Student-Specific OCR Tuning:**
```go
type HandwritingProfile struct {
    StudentID       uuid.UUID
    CommonCharMistakes map[rune]rune // 'a' confused with 'o'
    WritingStyle    string // 'cursive', 'print', 'mixed'
    OCRAdjustments  map[string]float64
}

func ImproveOCR(text string, profile HandwritingProfile) string {
    // Apply student-specific corrections
    for wrong, correct := range profile.CommonCharMistakes {
        text = strings.ReplaceAll(text, string(wrong), string(correct))
    }
    return text
}
```

**Learning:**
- Gemini analyzes past submissions
- Identifies recurring OCR errors
- Builds correction dictionary

---

### Feature C: Plagiarism & Collaboration Detection

**Similarity Analysis:**
```go
func DetectSimilarity(submission1, submission2 Submission) SimilarityReport {
    // 1. Text-based (cosine similarity on embeddings)
    textSim := calculateTextSimilarity(submission1.Answers, submission2.Answers)
    
    // 2. Diagram-based (visual similarity)
    visualSim := compareImages(submission1.Diagrams, submission2.Diagrams)
    
    // 3. Gemini reasoning comparison
    reasoningSim := compareReasoningPatterns(submission1, submission2)
    
    return SimilarityReport{
        OverallScore: (textSim + visualSim + reasoningSim) / 3,
        FlaggedSections: identifyMatchingSections(),
    }
}
```

**Gemini 3 Usage:**
```
Role: Academic Integrity Analyst
Input: Two student answers + similarity metrics
Output: Likelihood of collaboration + suspicious patterns
Technique: Reasoning over writing style, mistake patterns, unusual similarities
```

---

### Feature D: Exam Difficulty Calibration

**Post-Grading Analytics:**
```sql
SELECT 
    question_id,
    AVG(score) as avg_score,
    STDDEV(score) as score_variance,
    COUNT(CASE WHEN score = 0 THEN 1 END) as zero_scores,
    COUNT(CASE WHEN score = max_score THEN 1 END) as perfect_scores
FROM grades
GROUP BY question_id;
```

**Gemini 3 Recommendations:**
```
Role: Assessment Designer
Input: Grading statistics + question rubrics
Output: Difficulty rating + suggestions for improvement
Example: "Question 3 has 78% zero scores. Consider:
- Adding a scaffolding sub-question
- Clarifying the prompt wording
- Providing a worked example"
```

---

### Feature E: Voice-Based Grading Review

**Teacher Voice Commands:**
- "Accept all grades above 85%"
- "Show me the most confused answers"
- "Override question 5 to 8 points because..."

**Implementation:**
```go
func ProcessVoiceCommand(audio []byte) GradingAction {
    // 1. Speech-to-text (Gemini or Whisper)
    transcript := transcribeAudio(audio)
    
    // 2. Intent classification
    action := parseGradingIntent(transcript)
    
    // 3. Execute with confirmation
    return action
}
```

---

### Feature F: Export & Integration

**Formats:**
- CSV (gradebook import)
- PDF (detailed reports with AI reasoning)
- Google Classroom integration
- Canvas LMS integration

**API Example:**
```go
POST /api/v1/exams/{id}/export
{
    "format": "pdf",
    "include_ai_reasoning": true,
    "include_images": false,
    "group_by": "student"
}
```

---

## Technical Architecture

### System Diagram

```
┌─────────────┐
│   Teacher   │
│  Dashboard  │
└──────┬──────┘
       │
       v
┌─────────────────────────────────────┐
│         API Gateway (Go)            │
│  - Auth & Rate Limiting             │
│  - Request Routing                  │
└──────┬──────────────────────────────┘
       │
       ├───────────────┬───────────────┬──────────────┐
       v               v               v              v
┌──────────┐    ┌──────────┐    ┌──────────┐   ┌──────────┐
│  Upload  │    │ Grading  │    │Analytics │   │ Admin    │
│  Service │    │  Engine  │    │ Service  │   │ Service  │
└────┬─────┘    └────┬─────┘    └────┬─────┘   └──────────┘
     │               │               │
     v               v               v
┌────────────────────────────────────────┐
│           PostgreSQL                   │
│  - Exams, Rubrics, Submissions         │
│  - Grades, Overrides, Audit Logs       │
└────────────────────────────────────────┘

     │               │               │
     v               v               v
┌─────────┐    ┌──────────┐    ┌─────────┐
│  MinIO  │    │  Qdrant  │    │ Gemini  │
│ (Files) │    │ (Vectors)│    │   API   │
└─────────┘    └──────────┘    └─────────┘
```

---

### Data Flow: Upload → Grade → Feedback

```
1. Teacher uploads exam
   └─> MinIO storage + OCR job queued

2. Worker picks up job
   ├─> Google Vision OCR
   ├─> Gemini verification (low confidence)
   └─> Store OCR results

3. Segmentation job
   ├─> Detect question boundaries
   ├─> Extract diagrams
   └─> Associate answers with questions

4. Grading job (per question)
   ├─> Load rubric
   ├─> Call Gemini 3x (multi-evaluator)
   ├─> Calculate variance
   ├─> If low variance → Auto-grade
   └─> If high variance → Escalate

5. Teacher review (if escalated)
   ├─> View AI reasoning
   ├─> Override score
   └─> Capture feedback event

6. Feedback processing
   ├─> Analyze patterns
   ├─> Update rubric weights
   └─> Generate student feedback
```

---

## Deployment Strategy

### Docker Compose (Local Development)

```yaml
version: '3.8'
services:
  api:
    build: ./backend
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://haramma:pass@postgres:5432/haramma
      - GEMINI_API_KEY=${GEMINI_API_KEY}
    depends_on:
      - postgres
      - minio
      - qdrant
  
  postgres:
    image: postgres:15
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment:
      - POSTGRES_PASSWORD=pass
  
  minio:
    image: minio/minio
    command: server /data --console-address ":9001"
    ports:
      - "9000:9000"
      - "9001:9001"
  
  qdrant:
    image: qdrant/qdrant
    ports:
      - "6333:6333"
  
  worker:
    build: ./backend
    command: ./worker
    depends_on:
      - postgres
      - minio

volumes:
  postgres_data:
```

---

### Production (Fly.io)

```toml
app = "haramma"

[build]
  dockerfile = "Dockerfile"

[[services]]
  internal_port = 8080
  protocol = "tcp"

  [[services.ports]]
    port = 80
    handlers = ["http"]

  [[services.ports]]
    port = 443
    handlers = ["tls", "http"]

[env]
  PORT = "8080"

[[vm]]
  cpu_kind = "shared"
  cpus = 2
  memory_mb = 2048
```

---

## Success Metrics (Hackathon Demo)

### Quantitative Metrics

1. **Grading Speed:**
   - Target: 50 submissions/minute
   - Metric: Time from upload → final grades

2. **Accuracy:**
   - Target: 90% agreement with human graders
   - Metric: Mean absolute error on score

3. **Confidence Calibration:**
   - Target: 85% of high-confidence grades accepted without override
   - Metric: Override rate by confidence bucket

4. **Cost Efficiency:**
   - Target: <$0.10 per exam graded
   - Metric: Gemini API costs / submissions

### Qualitative Metrics

1. **Trust:** Teachers can explain AI decisions to students
2. **Transparency:** Full audit trail accessible
3. **Fairness:** Consistent grading across similar answers
4. **Usefulness:** Teachers prefer this over manual grading

---

## Demo Script (5-Minute Pitch)

### Slide 1: Problem (30s)
"Teachers spend 40% of their time grading. This time should be spent teaching. But automated grading today is a black box — teachers don't trust it, students don't understand it."

### Slide 2: Solution (45s)
"HARaMA uses Gemini 3's reasoning to grade exams transparently. Every score has an explanation. Every decision is auditable. Teachers stay in control."

**Live Demo:** Upload exam → Show grading in progress

### Slide 3: Gemini 3 Differentiation (60s)
"We use Gemini 3's unique strengths:
- **Long-context:** Understands entire exam, not just isolated questions
- **Multimodal:** Grades diagrams, equations, handwriting together
- **Reasoning:** Explains *why* a score was assigned, not just *what*"

**Live Demo:** Show multi-evaluator reasoning side-by-side

### Slide 4: Trust Layer (60s)
"Most AI grading fails because it pretends to be certain. We embrace uncertainty. When our three evaluators disagree, we escalate to humans. Low confidence? Teacher review required."

**Live Demo:** Show variance detection → escalation

### Slide 5: Impact (45s)
"Phase 1 already works — we graded 50 real exams with 92% teacher acceptance. But this is just the start."

**Show roadmap:** Feedback loops → Personalized student learning

### Slide 6: Ask (30s)
"We're building the grading infrastructure education deserves. One that's transparent, trustworthy, and built for teachers, not replacing them."

---

## Implementation Checklist (4-Day Build)

### Day 1: Foundation
- [x] Go API scaffold with PostgreSQL
- [x] Exam + Rubric data models
- [x] Upload → MinIO pipeline
- [x] Basic OCR (Google Vision / Mock)
- [x] Gemini API integration (hello world)

### Day 2: Grading Engine
- [x] Question segmentation
- [x] Rubric → structured prompt
- [x] Single-evaluator grading (Gemini)
- [x] Confidence scoring
- [x] Teacher override UI (basic)

### Day 3: Multimodal + Trust
- [x] Diagram extraction
- [x] Multimodal grading prompts
- [x] Multi-evaluator architecture
- [x] Variance calculation
- [x] Escalation logic

### Day 4: Polish + Demo
- [ ] Dashboard UI
- [ ] Audit log viewer
- [ ] Export to CSV
- [ ] Demo data preparation
- [ ] Video recording
- [ ] Documentation

---

## Risk Mitigation

### Risk 1: OCR Accuracy
**Mitigation:**
- Use Google Vision (>95% accuracy)
- Gemini verification for low-confidence
- Always show original image to teachers

### Risk 2: Gemini API Costs
**Mitigation:**
- Cache rubric interpretations
- Batch requests where possible
- Use Gemini Flash for simple questions

### Risk 3: Grading Fairness
**Mitigation:**
- Multi-evaluator consensus
- Audit trails
- Teacher override always wins

### Risk 4: Complexity Creep
**Mitigation:**
- Ship Phase 1 cleanly first
- Phase 2-5 are stretch goals, not requirements

---

## Conclusion

HARaMA is designed to win the Gemini 3 Hackathon by:

1. **Solving a real problem:** Teacher time scarcity
2. **Showcasing Gemini 3:** Long-context, multimodal, reasoning
3. **Building trust:** Transparent, auditable, human-in-the-loop
4. **Production-ready:** Clean architecture, scalable, extensible

**This is not a demo. This is a product.**

Ship Phase 1. Demo Phase 2. Pitch Phase 3-5.

Win.



# HARaMA — Project Structure

## Overview Directory Tree

```
harrama/
├── backend/                      # Go backend services
│   ├── cmd/                      # Application entrypoints
│   │   ├── api/                  # Main API server
│   │   │   └── main.go
│   │   ├── worker/               # Background job processor
│   │   │   └── main.go
│   │   └── migrate/              # Database migrations runner
│   │       └── main.go
│   │
│   ├── internal/                 # Private application code
│   │   ├── api/                  # HTTP handlers & routing
│   │   │   ├── router.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── tenant.go
│   │   │   │   ├── ratelimit.go
│   │   │   │   └── logging.go
│   │   │   └── handlers/
│   │   │       ├── exam.go       # Exam CRUD endpoints
│   │   │       ├── submission.go # Upload & submission endpoints
│   │   │       ├── grading.go    # Grading endpoints
│   │   │       ├── override.go   # Teacher override endpoints
│   │   │       ├── analytics.go  # Dashboard & reports
│   │   │       └── health.go     # Health checks
│   │   │
│   │   ├── domain/               # Core business logic (entities)
│   │   │   ├── exam.go
│   │   │   ├── question.go
│   │   │   ├── rubric.go
│   │   │   ├── submission.go
│   │   │   ├── grade.go
│   │   │   ├── student.go
│   │   │   └── teacher.go
│   │   │
│   │   ├── service/              # Business logic layer
│   │   │   ├── exam_service.go
│   │   │   ├── grading_service.go
│   │   │   ├── ocr_service.go
│   │   │   ├── segmentation_service.go
│   │   │   ├── feedback_service.go
│   │   │   └── analytics_service.go
│   │   │
│   │   ├── grading/              # Core grading engine
│   │   │   ├── engine.go         # Main grading orchestrator
│   │   │   ├── evaluator.go      # Single evaluator logic
│   │   │   ├── multi_eval.go     # Multi-evaluator coordinator
│   │   │   ├── confidence.go     # Confidence calculation
│   │   │   ├── variance.go       # Variance & disagreement detection
│   │   │   ├── partial_credit.go # Partial credit engine
│   │   │   └── profiles/         # Subject-specific configs
│   │   │       ├── mathematics.go
│   │   │       ├── english.go
│   │   │       └── science.go
│   │   │
│   │   ├── ai/                   # AI provider abstraction
│   │   │   ├── provider.go       # Interface definition
│   │   │   ├── gemini/
│   │   │   │   ├── client.go     # Gemini API client
│   │   │   │   ├── grading.go    # Grading prompts
│   │   │   │   ├── vision.go     # Multimodal/diagram
│   │   │   │   ├── feedback.go   # Student feedback generation
│   │   │   │   └── prompts/      # Prompt templates
│   │   │   │       ├── base_grading.txt
│   │   │   │       ├── rubric_enforcer.txt
│   │   │   │       ├── reasoning_validator.txt
│   │   │   │       ├── structural_analyzer.txt
│   │   │   │       └── multimodal_grading.txt
│   │   │   └── cache.go          # Response caching
│   │   │
│   │   ├── ocr/                  # OCR processing
│   │   │   ├── processor.go      # OCR orchestrator
│   │   │   ├── google_vision.go  # Google Vision API
│   │   │   ├── tesseract.go      # Tesseract fallback
│   │   │   ├── confidence.go     # OCR confidence scoring
│   │   │   └── correction.go     # Gemini-based correction
│   │   │
│   │   ├── segmentation/         # Answer segmentation
│   │   │   ├── detector.go       # Question boundary detection
│   │   │   ├── diagram.go        # Diagram extraction
│   │   │   └── layout.go         # Spatial analysis
│   │   │
│   │   ├── storage/              # File & object storage
│   │   │   ├── minio.go          # MinIO client
│   │   │   ├── upload.go         # Upload handling
│   │   │   └── retrieval.go      # File retrieval
│   │   │
│   │   ├── repository/           # Data access layer
│   │   │   ├── postgres/
│   │   │   │   ├── exam_repo.go
│   │   │   │   ├── submission_repo.go
│   │   │   │   ├── grade_repo.go
│   │   │   │   ├── override_repo.go
│   │   │   │   ├── feedback_repo.go
│   │   │   │   └── audit_repo.go
│   │   │   └── qdrant/
│   │   │       └── vector_repo.go
│   │   │
│   │   ├── worker/               # Background job processing
│   │   │   ├── pool.go           # Worker pool
│   │   │   ├── jobs/
│   │   │   │   ├── ocr_job.go
│   │   │   │   ├── segmentation_job.go
│   │   │   │   ├── grading_job.go
│   │   │   │   └── feedback_job.go
│   │   │   └── queue.go          # Job queue (PostgreSQL-based)
│   │   │
│   │   ├── auth/                 # Authentication & authorization
│   │   │   ├── jwt.go
│   │   │   ├── permissions.go
│   │   │   └── tenant.go
│   │   │
│   │   ├── config/               # Configuration management
│   │   │   ├── config.go
│   │   │   └── validator.go
│   │   │
│   │   └── pkg/                  # Shared utilities
│   │       ├── logger/
│   │       │   └── logger.go
│   │       ├── errors/
│   │       │   └── errors.go
│   │       ├── validators/
│   │       │   └── validators.go
│   │       └── utils/
│   │           ├── image.go      # Image processing helpers
│   │           └── math.go       # Math utilities
│   │
│   ├── migrations/               # Database migrations
│   │   ├── 001_initial_schema.up.sql
│   │   ├── 001_initial_schema.down.sql
│   │   ├── 002_add_multimodal.up.sql
│   │   ├── 002_add_multimodal.down.sql
│   │   ├── 003_add_feedback.up.sql
│   │   └── 003_add_feedback.down.sql
│   │
│   ├── tests/                    # Tests
│   │   ├── unit/
│   │   │   ├── grading_test.go
│   │   │   ├── confidence_test.go
│   │   │   └── variance_test.go
│   │   ├── integration/
│   │   │   ├── api_test.go
│   │   │   └── grading_flow_test.go
│   │   └── fixtures/
│   │       ├── sample_exam.pdf
│   │       └── sample_submission.pdf
│   │
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── Makefile
│   └── README.md
│
├── frontend/                     # React/Next.js frontend
│   ├── public/
│   │   ├── icons/
│   │   └── images/
│   │
│   ├── src/
│   │   ├── app/                  # Next.js App Router
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx          # Landing page
│   │   │   ├── dashboard/
│   │   │   │   └── page.tsx      # Teacher dashboard
│   │   │   ├── exams/
│   │   │   │   ├── page.tsx      # Exam list
│   │   │   │   ├── [id]/
│   │   │   │   │   └── page.tsx  # Exam detail
│   │   │   │   └── create/
│   │   │   │       └── page.tsx  # Create exam
│   │   │   ├── grading/
│   │   │   │   └── [id]/
│   │   │   │       └── page.tsx  # Grading interface
│   │   │   └── analytics/
│   │   │       └── page.tsx      # Analytics dashboard
│   │   │
│   │   ├── components/           # Reusable components
│   │   │   ├── ui/               # Base UI components
│   │   │   │   ├── Button.tsx
│   │   │   │   ├── Input.tsx
│   │   │   │   ├── Card.tsx
│   │   │   │   ├── Modal.tsx
│   │   │   │   └── Badge.tsx
│   │   │   ├── exam/
│   │   │   │   ├── ExamCard.tsx
│   │   │   │   ├── QuestionEditor.tsx
│   │   │   │   └── RubricBuilder.tsx
│   │   │   ├── grading/
│   │   │   │   ├── GradingView.tsx
│   │   │   │   ├── AnswerDisplay.tsx
│   │   │   │   ├── AIReasoningPanel.tsx
│   │   │   │   ├── ConfidenceMeter.tsx
│   │   │   │   ├── VarianceIndicator.tsx
│   │   │   │   └── OverrideForm.tsx
│   │   │   ├── submission/
│   │   │   │   ├── UploadZone.tsx
│   │   │   │   └── SubmissionList.tsx
│   │   │   └── analytics/
│   │   │       ├── GradingTrends.tsx
│   │   │       └── QuestionDifficulty.tsx
│   │   │
│   │   ├── lib/                  # Utilities & helpers
│   │   │   ├── api.ts            # API client
│   │   │   ├── auth.ts           # Auth helpers
│   │   │   └── utils.ts          # Common utilities
│   │   │
│   │   ├── hooks/                # Custom React hooks
│   │   │   ├── useExams.ts
│   │   │   ├── useGrading.ts
│   │   │   └── useAuth.ts
│   │   │
│   │   ├── types/                # TypeScript types
│   │   │   ├── exam.ts
│   │   │   ├── grading.ts
│   │   │   └── api.ts
│   │   │
│   │   └── styles/
│   │       └── globals.css
│   │
│   ├── package.json
│   ├── tsconfig.json
│   ├── next.config.js
│   └── tailwind.config.js
│
├── docs/                         # Documentation
│   ├── API.md                    # API documentation
│   ├── ARCHITECTURE.md           # System architecture
│   ├── DEPLOYMENT.md             # Deployment guide
│   ├── GRADING_LOGIC.md          # Grading algorithm details
│   └── PROMPTS.md                # Prompt engineering guide
│
├── scripts/                      # Utility scripts
│   ├── seed_data.go              # Database seeding
│   ├── test_grading.sh           # Grading test script
│   └── deploy.sh                 # Deployment script
│
├── docker-compose.yml            # Local development stack
├── docker-compose.prod.yml       # Production stack
├── .env.example                  # Environment variables template
├── .gitignore
├── LICENSE
└── README.md                     # Project overview
```

---
# HARaMA — Problem Domain Skills Breakdown

> **Purpose:** This document organizes skills by specific problem domains and technical challenges you'll encounter while building HARaMA.

---

## Table of Contents

1. [Problem Domain 1: Document Intelligence](#problem-domain-1-document-intelligence)
2. [Problem Domain 2: AI Orchestration & Reliability](#problem-domain-2-ai-orchestration--reliability)
3. [Problem Domain 3: Educational Assessment Logic](#problem-domain-3-educational-assessment-logic)
4. [Problem Domain 4: Real-Time Systems & Performance](#problem-domain-4-real-time-systems--performance)
5. [Problem Domain 5: Data Integrity & Auditability](#problem-domain-5-data-integrity--auditability)
6. [Problem Domain 6: User Experience & Trust Building](#problem-domain-6-user-experience--trust-building)

---

## Problem Domain 1: Document Intelligence

### Challenge: "How do we accurately extract and understand handwritten exam submissions?"

#### Skill Set 1: OCR Pipeline Engineering

**The Problem:**
Raw exam scans contain:
- Handwritten text (varying quality)
- Printed questions
- Diagrams and drawings
- Margin notes
- Coffee stains, crumpled corners, shadows

**Skills Required:**

##### 1.1 Image Preprocessing
```go
import "gocv.io/x/gocv"

func preprocessImage(imagePath string) (gocv.Mat, error) {
    img := gocv.IMRead(imagePath, gocv.IMReadColor)
    
    // Step 1: Deskew (fix rotation)
    angle := detectSkewAngle(img)
    rotated := rotateImage(img, angle)
    
    // Step 2: Denoise
    denoised := gocv.NewMat()
    gocv.FastNlMeansDenoisingColored(rotated, &denoised, 10, 10, 7, 21)
    
    // Step 3: Enhance contrast (CLAHE)
    lab := gocv.NewMat()
    gocv.CvtColor(denoised, &lab, gocv.ColorBGRToLab)
    
    channels := gocv.Split(lab)
    clahe := gocv.NewCLAHE()
    clahe.SetClipLimit(2.0)
    clahe.Apply(channels[0], &channels[0])
    
    gocv.Merge(channels, &lab)
    gocv.CvtColor(lab, &denoised, gocv.ColorLabToBGR)
    
    return denoised, nil
}

func detectSkewAngle(img gocv.Mat) float64 {
    // Use Hough Line Transform to detect text lines
    gray := gocv.NewMat()
    gocv.CvtColor(img, &gray, gocv.ColorBGRToGray)
    
    edges := gocv.NewMat()
    gocv.Canny(gray, &edges, 50, 150)
    
    lines := gocv.NewMat()
    gocv.HoughLinesP(edges, &lines, 1, math.Pi/180, 100)
    
    // Calculate dominant angle
    angles := []float64{}
    for i := 0; i < lines.Rows(); i++ {
        line := lines.GetVecfAt(i, 0)
        x1, y1, x2, y2 := line[0], line[1], line[2], line[3]
        angle := math.Atan2(y2-y1, x2-x1) * 180 / math.Pi
        angles = append(angles, angle)
    }
    
    return calculateMedianAngle(angles)
}
```

**What You're Learning:**
- Computer vision fundamentals
- Image quality assessment
- Geometric transformations
- Noise reduction techniques

**Resources:**
- OpenCV tutorials: https://docs.opencv.org/
- GoCV examples: https://github.com/hybridgroup/gocv
- Image processing fundamentals: "Digital Image Processing" by Gonzalez

---

##### 1.2 OCR Quality Assessment

**The Problem:**
Not all OCR results are equal. You need to know when to trust OCR output.

```go
type OCRConfidenceMetrics struct {
    CharacterLevelConfidence float64 // Google Vision provides this
    WordLevelConfidence      float64 // Calculated from characters
    RegionClarity            float64 // Based on image quality
    OverallConfidence        float64 // Weighted combination
}

func assessOCRQuality(visionResponse *vision.TextAnnotation, imageRegion gocv.Mat) OCRConfidenceMetrics {
    metrics := OCRConfidenceMetrics{}
    
    // 1. Character-level confidence (from Google Vision)
    charConfidences := []float64{}
    for _, symbol := range visionResponse.Pages[0].Blocks[0].Paragraphs[0].Words[0].Symbols {
        charConfidences = append(charConfidences, float64(symbol.Confidence))
    }
    metrics.CharacterLevelConfidence = calculateMean(charConfidences)
    
    // 2. Image quality metrics
    metrics.RegionClarity = assessImageClarity(imageRegion)
    
    // 3. Consistency check
    metrics.WordLevelConfidence = checkWordConsistency(visionResponse.Text)
    
    // 4. Weighted overall
    weights := map[string]float64{
        "char":    0.5,
        "clarity": 0.3,
        "word":    0.2,
    }
    
    metrics.OverallConfidence = 
        (metrics.CharacterLevelConfidence * weights["char"]) +
        (metrics.RegionClarity * weights["clarity"]) +
        (metrics.WordLevelConfidence * weights["word"])
    
    return metrics
}

func assessImageClarity(region gocv.Mat) float64 {
    // Laplacian variance - higher = sharper
    gray := gocv.NewMat()
    gocv.CvtColor(region, &gray, gocv.ColorBGRToGray)
    
    laplacian := gocv.NewMat()
    gocv.Laplacian(gray, &laplacian, gocv.MatTypeCV64F, 1, 1, 0, gocv.BorderDefault)
    
    mean, stddev := laplacian.MeanStdDev()
    variance := stddev.Val1 * stddev.Val1
    
    // Normalize to 0-1 range (empirically determined)
    normalized := math.Min(variance/1000.0, 1.0)
    return normalized
}

func checkWordConsistency(text string) float64 {
    // Use dictionary to check if words are real
    words := strings.Fields(text)
    validWords := 0
    
    for _, word := range words {
        if isDictionaryWord(word) || isNumber(word) || isMathSymbol(word) {
            validWords++
        }
    }
    
    return float64(validWords) / float64(len(words))
}
```

**What You're Learning:**
- Quality metrics for OCR output
- Statistical confidence calculation
- Image sharpness detection
- Heuristic validation

---

##### 1.3 Gemini-Assisted OCR Correction

**The Problem:**
OCR makes mistakes. "6" becomes "G", "1" becomes "I", etc.

```go
func correctOCRWithGemini(ctx context.Context, ocrText string, confidence float64, imageBytes []byte) (string, error) {
    // Only use Gemini for low-confidence results
    if confidence > 0.85 {
        return ocrText, nil // Trust Google Vision
    }
    
    prompt := fmt.Sprintf(`
ROLE: OCR correction specialist

RAW OCR OUTPUT (confidence: %.2f):
"%s"

ORIGINAL IMAGE: [attached]

TASK:
1. Compare OCR text with the actual image
2. Identify likely OCR errors (common mistakes: 0→O, 1→I, 5→S)
3. Provide corrected text
4. Explain what was fixed

OUTPUT JSON:
{
  "corrected_text": "...",
  "corrections": [
    {"position": 5, "original": "G", "corrected": "6", "reason": "Clearly a digit not letter"}
  ],
  "new_confidence": 0.95
}
`, confidence, ocrText)
    
    response := callGeminiVision(ctx, prompt, imageBytes)
    
    var correction OCRCorrection
    json.Unmarshal([]byte(response), &correction)
    
    return correction.CorrectedText, nil
}

type OCRCorrection struct {
    CorrectedText  string                `json:"corrected_text"`
    Corrections    []CorrectionDetail    `json:"corrections"`
    NewConfidence  float64               `json:"new_confidence"`
}

type CorrectionDetail struct {
    Position  int    `json:"position"`
    Original  string `json:"original"`
    Corrected string `json:"corrected"`
    Reason    string `json:"reason"`
}
```

**What You're Learning:**
- Combining rule-based and AI approaches
- Cost-benefit analysis (only use Gemini when needed)
- Multimodal AI application
- Error recovery strategies

---

#### Skill Set 2: Document Segmentation

**The Problem:**
An exam has multiple questions. You need to:
1. Find where each question's answer starts/ends
2. Handle multi-page answers
3. Associate diagrams with correct questions
4. Deal with margin overflows

##### 2.1 Spatial Layout Analysis

```go
type DocumentLayout struct {
    Pages     []PageLayout
    Questions []QuestionRegion
}

type PageLayout struct {
    PageNumber int
    Width      int
    Height     int
    TextBlocks []TextBlock
    Images     []ImageRegion
}

type TextBlock struct {
    BoundingBox BoundingBox
    Text        string
    Confidence  float64
}

type QuestionRegion struct {
    QuestionNumber int
    StartPage      int
    StartY         int
    EndPage        int
    EndY           int
    TextBlocks     []TextBlock
    Diagrams       []ImageRegion
}

func segmentByLayout(layout DocumentLayout, template ExamTemplate) []QuestionRegion {
    regions := []QuestionRegion{}
    
    // Strategy 1: Look for question number markers
    markers := findQuestionMarkers(layout)
    
    // Strategy 2: Use template to predict boundaries
    if len(markers) != len(template.Questions) {
        markers = predictBoundaries(layout, template)
    }
    
    // Build regions
    for i, marker := range markers {
        region := QuestionRegion{
            QuestionNumber: i + 1,
            StartPage:      marker.Page,
            StartY:         marker.Y,
        }
        
        // Find end (next marker or page end)
        if i+1 < len(markers) {
            region.EndPage = markers[i+1].Page
            region.EndY = markers[i+1].Y
        } else {
            region.EndPage = layout.Pages[len(layout.Pages)-1].PageNumber
            region.EndY = layout.Pages[len(layout.Pages)-1].Height
        }
        
        // Extract text and diagrams in this region
        region.TextBlocks = extractTextInRegion(layout, region)
        region.Diagrams = extractDiagramsInRegion(layout, region)
        
        regions = append(regions, region)
    }
    
    return regions
}

type QuestionMarker struct {
    Page   int
    Y      int
    Number int
}

func findQuestionMarkers(layout DocumentLayout) []QuestionMarker {
    markers := []QuestionMarker{}
    
    // Patterns: "1.", "Q1", "Question 1", "1)", etc.
    patterns := []*regexp.Regexp{
        regexp.MustCompile(`^\d+\.`),
        regexp.MustCompile(`^Q\d+`),
        regexp.MustCompile(`^Question\s+\d+`),
        regexp.MustCompile(`^\d+\)`),
    }
    
    for _, page := range layout.Pages {
        for _, block := range page.TextBlocks {
            for _, pattern := range patterns {
                if pattern.MatchString(strings.TrimSpace(block.Text)) {
                    questionNum := extractQuestionNumber(block.Text)
                    markers = append(markers, QuestionMarker{
                        Page:   page.PageNumber,
                        Y:      block.BoundingBox.Y,
                        Number: questionNum,
                    })
                    break
                }
            }
        }
    }
    
    return markers
}
```

**What You're Learning:**
- Document structure analysis
- Heuristic algorithm design
- Spatial reasoning
- Pattern recognition

---

##### 2.2 Gemini-Assisted Segmentation (Ambiguous Cases)

```go
func geminiVerifySegmentation(ctx context.Context, layout DocumentLayout, template ExamTemplate) []QuestionRegion {
    prompt := fmt.Sprintf(`
ROLE: Document analysis expert

EXAM TEMPLATE:
Questions: %d
Expected structure:
%s

DETECTED LAYOUT:
Total pages: %d
Text blocks found: %d
Potential question markers: %s

TASK:
The automatic segmentation found %d regions but expected %d.
Analyze the document and provide the correct segmentation.

OUTPUT JSON:
{
  "regions": [
    {
      "question_number": 1,
      "start_page": 1,
      "start_y": 150,
      "end_page": 1,
      "end_y": 450,
      "reasoning": "Clear '1.' marker at Y=150"
    },
    ...
  ],
  "issues": ["Question 3 spans two pages", "No marker for Q4, inferred from spacing"]
}
`, len(template.Questions), 
   formatTemplateStructure(template),
   len(layout.Pages),
   countTextBlocks(layout),
   describeMarkers(layout))
    
    // Attach page images for visual verification
    images := renderLayoutImages(layout)
    
    response := callGeminiVision(ctx, prompt, images...)
    
    var result SegmentationResult
    json.Unmarshal([]byte(response), &result)
    
    return result.Regions
}
```

**What You're Learning:**
- When to escalate to AI
- Multimodal reasoning
- Structured data extraction
- Verification workflows

---

## Problem Domain 2: AI Orchestration & Reliability

### Challenge: "How do we make AI grading reliable, explainable, and cost-effective?"

#### Skill Set 3: Multi-Evaluator Architecture

**The Problem:**
A single AI call can be wrong. How do we build confidence through consensus?

##### 3.1 Evaluator Profile Design

```go
type EvaluatorProfile struct {
    ID           string
    Name         string
    SystemPrompt string
    Temperature  float64
    Perspective  string // "strict", "lenient", "balanced"
    FocusAreas   []string
}

var Evaluators = map[string]EvaluatorProfile{
    "rubric_enforcer": {
        ID:          "rubric_enforcer",
        Name:        "Rubric Enforcer",
        SystemPrompt: `You are a strict grader who follows the rubric exactly.
Award full credit only when ALL criteria are explicitly met.
Do not give partial credit unless the rubric specifically allows it.
Your job is to ensure consistency and fairness.`,
        Temperature: 0.1, // Very deterministic
        Perspective: "strict",
        FocusAreas:  []string{"rubric_compliance", "completeness"},
    },
    
    "reasoning_validator": {
        ID:          "reasoning_validator",
        Name:        "Reasoning Validator",
        SystemPrompt: `You are an educator who values logical thinking.
Reward students for correct reasoning even if execution has minor errors.
Look for conceptual understanding, not just correct final answers.
Partial credit should be generous for good reasoning with small mistakes.`,
        Temperature: 0.3, // More flexible
        Perspective: "lenient",
        FocusAreas:  []string{"logical_flow", "conceptual_understanding"},
    },
    
    "structural_analyzer": {
        ID:          "structural_analyzer",
        Name:        "Structural Analyzer",
        SystemPrompt: `You evaluate answer structure and organization.
Check for: clear introduction, step-by-step work, labeled diagrams.
Penalize disorganized answers even if content is correct.
Reward well-structured answers with clear explanations.`,
        Temperature: 0.2,
        Perspective: "balanced",
        FocusAreas:  []string{"organization", "clarity", "presentation"},
    },
}
```

**What You're Learning:**
- Designing complementary AI perspectives
- Prompt engineering for specific behaviors
- Temperature tuning for consistency vs creativity
- System architecture for parallel processing

---

##### 3.2 Parallel Evaluation Orchestration

```go
func multiEvaluatorGrade(ctx context.Context, answer AnswerSegment, rubric Rubric) (*MultiEvalResult, error) {
    evaluatorIDs := []string{"rubric_enforcer", "reasoning_validator", "structural_analyzer"}
    
    // Create channels for results
    type EvalTask struct {
        EvaluatorID string
        Result      GradingResult
        Error       error
    }
    
    results := make(chan EvalTask, len(evaluatorIDs))
    
    // Launch parallel evaluations
    var wg sync.WaitGroup
    for _, evalID := range evaluatorIDs {
        wg.Add(1)
        go func(id string) {
            defer wg.Done()
            
            profile := Evaluators[id]
            result, err := gradeWithProfile(ctx, answer, rubric, profile)
            
            results <- EvalTask{
                EvaluatorID: id,
                Result:      result,
                Error:       err,
            }
        }(evalID)
    }
    
    // Wait for all evaluations
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    evalResults := []GradingResult{}
    for task := range results {
        if task.Error != nil {
            log.Printf("Evaluator %s failed: %v", task.EvaluatorID, task.Error)
            continue
        }
        evalResults = append(evalResults, task.Result)
    }
    
    // Analyze variance
    return analyzeMultiEval(evalResults)
}

func gradeWithProfile(ctx context.Context, answer AnswerSegment, rubric Rubric, profile EvaluatorProfile) (GradingResult, error) {
    prompt := buildGradingPrompt(answer, rubric, profile)
    
    client, _ := genai.NewClient(ctx, option.WithAPIKey(apiKey))
    defer client.Close()
    
    model := client.GenerativeModel("gemini-3-pro")
    model.SetTemperature(profile.Temperature)
    model.SystemInstruction = &genai.Content{
        Parts: []genai.Part{genai.Text(profile.SystemPrompt)},
    }
    
    resp, err := model.GenerateContent(ctx, genai.Text(prompt))
    if err != nil {
        return GradingResult{}, err
    }
    
    result := parseGradingResponse(resp)
    result.AIEvaluatorID = profile.ID
    result.EvaluatorPerspective = profile.Perspective
    
    return result, nil
}
```

**What You're Learning:**
- Concurrent programming in Go (goroutines, channels)
- Error handling in parallel systems
- AI model configuration
- Result aggregation

---

##### 3.3 Consensus Building & Disagreement Detection

```go
type MultiEvalResult struct {
    Evaluations    []GradingResult
    MeanScore      float64
    MedianScore    float64
    Variance       float64
    StdDev         float64
    ConsensusScore float64
    Confidence     float64
    ShouldEscalate bool
    Reasoning      string
}

func analyzeMultiEval(evaluations []GradingResult) (*MultiEvalResult, error) {
    if len(evaluations) == 0 {
        return nil, errors.New("no evaluation results")
    }
    
    scores := extractScores(evaluations)
    
    result := &MultiEvalResult{
        Evaluations: evaluations,
        MeanScore:   calculateMean(scores),
        MedianScore: calculateMedian(scores),
    }
    
    // Calculate variance
    variance := 0.0
    for _, score := range scores {
        variance += math.Pow(score-result.MeanScore, 2)
    }
    result.Variance = variance / float64(len(scores))
    result.StdDev = math.Sqrt(result.Variance)
    
    // Determine consensus score (confidence-weighted average)
    result.ConsensusScore = calculateWeightedConsensus(evaluations)
    
    // Calculate overall confidence
    result.Confidence = calculateConsensusConfidence(evaluations, result.Variance)
    
    // Escalation decision
    maxScore := evaluations[0].MaxScore
    varianceThreshold := 0.15 * float64(maxScore)
    result.ShouldEscalate = result.Variance > varianceThreshold || result.Confidence < 0.70
    
    // Generate reasoning
    result.Reasoning = generateConsensusReasoning(evaluations, result)
    
    return result, nil
}

func calculateWeightedConsensus(evaluations []GradingResult) float64 {
    totalWeight := 0.0
    weightedSum := 0.0
    
    for _, eval := range evaluations {
        weight := eval.Confidence
        totalWeight += weight
        weightedSum += eval.Score * weight
    }
    
    if totalWeight == 0 {
        return calculateMean(extractScores(evaluations))
    }
    
    return weightedSum / totalWeight
}

func calculateConsensusConfidence(evaluations []GradingResult, variance float64) float64 {
    // Confidence decreases with:
    // 1. High variance between evaluators
    // 2. Low individual confidences
    
    avgIndividualConfidence := 0.0
    for _, eval := range evaluations {
        avgIndividualConfidence += eval.Confidence
    }
    avgIndividualConfidence /= float64(len(evaluations))
    
    // Variance penalty (normalized)
    variancePenalty := math.Max(0, 1.0-(variance/10.0))
    
    // Combined confidence
    return (avgIndividualConfidence * 0.6) + (variancePenalty * 0.4)
}

func generateConsensusReasoning(evaluations []GradingResult, result *MultiEvalResult) string {
    if result.Variance < 1.0 {
        return fmt.Sprintf(
            "All evaluators agree (variance: %.2f). Consensus score: %.1f",
            result.Variance, result.ConsensusScore,
        )
    }
    
    // Identify disagreements
    strictScore := findEvaluatorScore(evaluations, "rubric_enforcer")
    lenientScore := findEvaluatorScore(evaluations, "reasoning_validator")
    
    if math.Abs(strictScore-lenientScore) > 2.0 {
        return fmt.Sprintf(
            "Significant disagreement: Strict grader gave %.1f, lenient gave %.1f. "+
            "This suggests ambiguity in the answer. Human review recommended.",
            strictScore, lenientScore,
        )
    }
    
    return fmt.Sprintf(
        "Moderate variance (%.2f). Consensus reached through weighted average: %.1f",
        result.Variance, result.ConsensusScore,
    )
}
```

**What You're Learning:**
- Statistical analysis (variance, std dev)
- Weighted averaging algorithms
- Decision tree logic
- Explainability in AI systems

---

#### Skill Set 4: Prompt Engineering Mastery

**The Problem:**
The quality of AI output is 90% prompt quality.

##### 4.1 Structured Prompt Templates

```go
const baseGradingPromptTemplate = `
ROLE: {{.Role}}

SYSTEM INSTRUCTIONS:
{{.SystemInstructions}}

EXAM CONTEXT:
Subject: {{.Subject}}
Question Type: {{.QuestionType}}
Max Score: {{.MaxScore}} points

RUBRIC:
{{.RubricJSON}}

STUDENT ANSWER:
{{.AnswerText}}

{{if .HasDiagram}}
STUDENT DIAGRAM: [see attached image]
{{end}}

EVALUATION INSTRUCTIONS:
1. Read the rubric criteria carefully
2. Evaluate each criterion systematically
3. {{.SpecialInstructions}}
4. Assign score with clear reasoning
5. Identify specific mistakes if any

OUTPUT REQUIRED (JSON):
{
  "score": <number 0-{{.MaxScore}}>,
  "confidence": <float 0.0-1.0>,
  "reasoning": "<detailed explanation>",
  "criteria_met": ["<criterion1>", "<criterion2>", ...],
  "criteria_not_met": ["<criterion3>", ...],
  "mistakes_found": [
    {
      "mistake": "<description>",
      "impact": "<how it affects score>",
      "suggestion": "<how to improve>"
    }
  ],
  "partial_credit_awarded": [
    {
      "criterion": "<which criterion>",
      "points": <amount>,
      "reason": "<why partial credit>"
    }
  ]
}

IMPORTANT:
- Be objective and fair
- Explain your reasoning clearly
- Use the rubric as the source of truth
- If unsure, indicate lower confidence
`

type PromptData struct {
    Role                string
    SystemInstructions  string
    Subject             string
    QuestionType        string
    MaxScore            int
    RubricJSON          string
    AnswerText          string
    HasDiagram          bool
    SpecialInstructions string
}

func buildGradingPrompt(answer AnswerSegment, rubric Rubric, profile EvaluatorProfile) string {
    data := PromptData{
        Role:               profile.Name,
        SystemInstructions: profile.SystemPrompt,
        Subject:            answer.Question.Subject,
        QuestionType:       string(answer.Question.Type),
        MaxScore:           answer.Question.Points,
        RubricJSON:         rubric.ToJSON(),
        AnswerText:         answer.Text,
        HasDiagram:         len(answer.Diagrams) > 0,
        SpecialInstructions: getSpecialInstructions(profile),
    }
    
    tmpl := template.Must(template.New("grading").Parse(baseGradingPromptTemplate))
    
    var buf bytes.Buffer
    tmpl.Execute(&buf, data)
    
    return buf.String()
}

func getSpecialInstructions(profile EvaluatorProfile) string {
    switch profile.Perspective {
    case "strict":
        return "Only award points when criteria are EXPLICITLY and COMPLETELY met"
    case "lenient":
        return "Give partial credit generously for correct reasoning, even with errors"
    case "balanced":
        return "Balance accuracy with fair assessment of student effort"
    default:
        return "Evaluate fairly according to the rubric"
    }
}
```

**What You're Learning:**
- Template-based prompt generation
- Dynamic content injection
- Structured output enforcement
- Perspective-based prompting

---

##### 4.2 Few-Shot Learning for Calibration

```go
func addFewShotExamples(prompt string, subject string) string {
    examples := getFewShotExamples(subject)
    
    exampleSection := "
CALIBRATION EXAMPLES:

Here are some graded examples to calibrate your expectations:

"
    
    for i, ex := range examples {
        exampleSection += fmt.Sprintf(`
EXAMPLE %d:
Question: %s
Student Answer: %s
Correct Grade: %d/%d
Reasoning: %s

`, i+1, ex.Question, ex.StudentAnswer, ex.Score, ex.MaxScore, ex.Reasoning)
    }
    
    return exampleSection + "

NOW GRADE THE ACTUAL ANSWER:

" + prompt
}

type FewShotExample struct {
    Question      string
    StudentAnswer string
    Score         int
    MaxScore      int
    Reasoning     string
}

func getFewShotExamples(subject string) []FewShotExample {
    // Subject-specific examples
    examples := map[string][]FewShotExample{
        "physics": {
            {
                Question:      "Calculate the force on a 10kg object accelerating at 5m/s²",
                StudentAnswer: "F = ma = (10)(5) = 50 N",
                Score:         5,
                MaxScore:      5,
                Reasoning:     "Correct formula, correct calculation, proper units. Full credit.",
            },
            {
                Question:      "Calculate the force on a 10kg object accelerating at 5m/s²",
                StudentAnswer: "F = ma = (10)(5) = 50",
                Score:         4,
                MaxScore:      5,
                Reasoning:     "Correct formula and calculation but missing units (-1 point).",
            },
            {
                Question:      "Calculate the force on a 10kg object accelerating at 5m/s²",
                StudentAnswer: "Force equals mass times acceleration which is 50",
                Score:         3,
                MaxScore:      5,
                Reasoning:     "Understands concept but no formula shown (-1), no units (-1). Partial credit for reasoning.",
            },
        },
        
        "english": {
            {
                Question:      "Analyze the theme of isolation in Frankenstein",
                StudentAnswer: "The monster is isolated because everyone rejects him. Victor also isolates himself in his lab. This shows how isolation leads to tragedy.",
                Score:         3,
                MaxScore:      5,
                Reasoning:     "Identifies theme correctly (+2) with basic examples (+1) but lacks depth and textual evidence (-2).",
            },
        },
    }
    
    return examples[subject]
}
```

**What You're Learning:**
- Few-shot prompting techniques
- Domain-specific calibration
- Example selection strategies
- Prompt augmentation

---

## Problem Domain 3: Educational Assessment Logic

### Challenge: "How do we implement fair, nuanced grading that respects educational principles?"

#### Skill Set 5: Rubric Engineering

**The Problem:**
Converting teacher intent into machine-readable grading criteria.

##### 5.1 Rubric Data Model

```go
type Rubric struct {
    ID                    uuid.UUID
    QuestionID            uuid.UUID
    FullCreditCriteria    []Criterion
    PartialCreditRules    []PartialCreditRule
    CommonMistakes        []CommonMistake
    KeyConcepts           []string
    GradingNotes          string
    StrictMode            bool
}

type Criterion struct {
    ID          string
    Description string
    Points      float64
    Required    bool // Must be met for any credit
    Category    string // "content", "methodology", "presentation"
}

type PartialCreditRule struct {
    ID          string
    Condition   string // Natural language or code
    Points      float64
    Description string
    Dependencies []string // Which criteria this depends on
}

type CommonMistake struct {
    ID          string
    Description string
    Penalty     float64
    Category    string
    Frequency   int // How often students make this
}

// Example rubric for a physics problem
var physicsRubric = Rubric{
    FullCreditCriteria: []Criterion{
        {
            ID:          "correct_formula",
            Description: "Uses the correct formula (F = ma)",
            Points:      2.0,
            Required:    false,
            Category:    "methodology",
        },
        {
            ID:          "correct_substitution",
            Description: "Substitutes values correctly",
            Points:      1.0,
            Required:    false,
            Category:    "content",
        },
        {
            ID:          "correct_calculation",
            Description: "Performs calculation accurately",
            Points:      1.0,
            Required:    false,
            Category:    "content",
        },
        {
            ID:          "units_specified",
            Description: "Includes proper units in answer",
            Points:      1.0,
            Required:    false,
            Category:    "presentation",
        },
    },
    
    PartialCreditRules: []PartialCreditRule{
        {
            ID:          "method_correct_calc_error",
            Condition:   "correct_formula AND NOT correct_calculation",
            Points:      3.0, // 2 for formula + 1 for attempt
            Description: "Right method but arithmetic error",
        },
        {
            ID:          "conceptual_understanding",
            Condition:   "describes_correct_concept AND NOT correct_formula",
            Points:      1.5,
            Description: "Shows understanding without formula",
        },
    },
    
    CommonMistakes: []CommonMistake{
        {
            ID:          "wrong_formula",
            Description: "Uses F = mv instead of F = ma",
            Penalty:     2.0,
            Category:    "conceptual_error",
        },
        {
            ID:          "unit_conversion_error",
            Description: "Fails to convert units before calculating",
            Penalty:     0.5,
            Category:    "procedural_error",
        },
    },
    
    KeyConcepts: []string{
        "Newton's Second Law",
        "Force as mass × acceleration",
        "Importance of units",
    },
}
```

**What You're Learning:**
- Educational assessment design
- Criterion decomposition
- Dependency modeling
- Structured data design

---

##### 5.2 Natural Language Rubric Parser

```go
func parseNaturalLanguageRubric(ctx context.Context, rubricText string) (*Rubric, error) {
    prompt := fmt.Sprintf(`
ROLE: Educational assessment expert

TEACHER'S RUBRIC (natural language):
"""
%s
"""

TASK:
Convert this rubric into a structured format.

OUTPUT JSON:
{
  "full_credit_criteria": [
    {
      "id": "unique_id",
      "description": "Clear criterion",
      "points": <number>,
      "required": <boolean>,
      "category": "content|methodology|presentation"
    }
  ],
  "partial_credit_rules": [
    {
      "id": "unique_id",
      "condition": "When to apply this rule",
      "points": <number>,
      "description": "What this covers"
    }
  ],
  "common_mistakes": [
    {
      "id": "unique_id",
      "description": "What students often get wrong",
      "penalty": <number>
    }
  ],
  "key_concepts": ["concept1", "concept2"]
}
`, rubricText)
    
    response := callGemini(ctx, prompt)
    
    var parsed ParsedRubric
    if err := json.Unmarshal([]byte(response), &parsed); err != nil {
        return nil, fmt.Errorf("failed to parse rubric: %w", err)
    }
    
    // Convert to internal Rubric structure
    rubric := &Rubric{
        FullCreditCriteria: parsed.Criteria,
        PartialCreditRules: parsed.PartialRules,
        CommonMistakes:     parsed.Mistakes,
        KeyConcepts:        parsed.KeyConcepts,
    }
    
    return rubric, nil
}
```

**What You're Learning:**
- NLP for educational content
- Structured data extraction
- Prompt design for parsing
- Schema validation

---

## Problem Domain 4: Real-Time Systems & Performance

### Challenge: "How do we grade 500 submissions in minutes, not hours?"

#### Skill Set 6: Concurrent Processing Architecture

**The Problem:**
Sequential processing is too slow. Need parallel pipeline.

##### 6.1 Worker Pool Pattern

```go
type WorkerPool struct {
    workers    int
    jobs       chan Job
    results    chan Result
    wg         sync.WaitGroup
    ctx        context.Context
    cancel     context.CancelFunc
}

type Job struct {
    ID           uuid.UUID
    SubmissionID uuid.UUID
    QuestionID   uuid.UUID
    Answer       AnswerSegment
    Rubric       Rubric
}

type Result struct {
    JobID   uuid.UUID
    Grade   GradingResult
    Error   error
}

func NewWorkerPool(numWorkers int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    
    return &WorkerPool{
        workers: numWorkers,
        jobs:    make(chan Job, numWorkers*2), // Buffer
        results: make(chan Result, numWorkers*2),
        ctx:     ctx,
        cancel:  cancel,
    }
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        p.wg.Add(1)
        go p.worker(i)
    }
    
    log.Printf("Started %d workers", p.workers)
}

func (p *WorkerPool) worker(id int) {
    defer p.wg.Done()
    
    log.Printf("Worker %d started", id)
    
    for {
        select {
        case <-p.ctx.Done():
            log.Printf("Worker %d stopping", id)
            return
            
        case job, ok := <-p.jobs:
            if !ok {
                return // Channel closed
            }
            
            log.Printf("Worker %d processing job %s", id, job.ID)
            
            // Grade the answer
            grade, err := gradeAnswer(p.ctx, job.Answer, job.Rubric)
            
            p.results <- Result{
                JobID: job.ID,
                Grade: grade,
                Error: err,
            }
        }
    }
}

func (p *WorkerPool) Submit(job Job) {
    p.jobs <- job
}

func (p *WorkerPool) Stop() {
    close(p.jobs)
    p.wg.Wait()
    close(p.results)
}

// Usage
func gradeSubmissionBatch(submissions []Submission) {
    pool := NewWorkerPool(10) // 10 concurrent graders
    pool.Start()
    defer pool.Stop()
    
    // Submit all jobs
    go func() {
        for _, sub := range submissions {
            for _, answer := range sub.Answers {
                pool.Submit(Job{
                    ID:           uuid.New(),
                    SubmissionID: sub.ID,
                    QuestionID:   answer.QuestionID,
                    Answer:       answer,
                    Rubric:       getRubric(answer.QuestionID),
                })
            }
        }
    }()
    
    // Collect results
    processed := 0
    total := countTotalAnswers(submissions)
    
    for result := range pool.results {
        if result.Error != nil {
            log.Printf("Job %s failed: %v", result.JobID, result.Error)
        } else {
            saveGrade(result.Grade)
        }
        
        processed++
        log.Printf("Progress: %d/%d (%.1f%%)", processed, total, float64(processed)/float64(total)*100)
        
        if processed == total {
            break
        }
    }
}
```

**What You're Learning:**
- Worker pool pattern
- Channel-based communication
- Context-based cancellation
- Progress tracking

---

##### 6.2 Rate Limiting for API Calls

```go
import "golang.org/x/time/rate"

type RateLimitedAIProvider struct {
    provider ai.Provider
    limiter  *rate.Limiter
}

func NewRateLimitedProvider(provider ai.Provider, requestsPerSecond int) *RateLimitedAIProvider {
    return &RateLimitedAIProvider{
        provider: provider,
        limiter:  rate.NewLimiter(rate.Limit(requestsPerSecond), requestsPerSecond*2),
    }
}

func (r *RateLimitedAIProvider) Grade(ctx context.Context, req ai.GradingRequest) (domain.GradingResult, error) {
    // Wait for rate limiter
    if err := r.limiter.Wait(ctx); err != nil {
        return domain.GradingResult{}, fmt.Errorf("rate limit wait failed: %w", err)
    }
    
    // Make actual call
    return r.provider.Grade(ctx, req)
}

// Advanced: Per-tenant rate limiting
type TenantAwareRateLimiter struct {
    limiters sync.Map // map[uuid.UUID]*rate.Limiter
    perTenantRPS int
}

func (t *TenantAwareRateLimiter) GetLimiter(tenantID uuid.UUID) *rate.Limiter {
    if limiter, ok := t.limiters.Load(tenantID); ok {
        return limiter.(*rate.Limiter)
    }
    
    limiter := rate.NewLimiter(rate.Limit(t.perTenantRPS), t.perTenantRPS*2)
    t.limiters.Store(tenantID, limiter)
    
    return limiter
}
```

**What You're Learning:**
- Rate limiting algorithms
- Token bucket implementation
- Per-tenant resource management
- Sync.Map for concurrent access

---

## Problem Domain 5: Data Integrity & Auditability

### Challenge: "How do we ensure every grading decision is traceable and tamper-proof?"

#### Skill Set 7: Immutable Audit Logging

(Continue with detailed implementation...)

---

**[Document continues with remaining problem domains...]**

Would you like me to continue with the remaining problem domains (5 and 6) in detail?
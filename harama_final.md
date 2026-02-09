# HARaMA - Product Requirements Document

**Version:** 3.0  
**Last Updated:** February 07, 2026  
**Document Owner:** Product Team  
**Status:** Backend Implemented - Frontend Development Ready

---

## 🚀 Implementation Status

### Backend (Go) - IMPLEMENTED ✅

| Component | Status | Notes |
|-----------|--------|-------|
| **Core API** | ✅ Complete | Chi Router, RESTful endpoints |
| **Database** | ✅ Complete | PostgreSQL with bun ORM, migrations ready |
| **Exam Management** | ✅ Complete | Create, list, get exams with questions |
| **Question/Rubric System** | ✅ Complete | Full credit criteria, partial credit rules |
| **Submission Processing** | ✅ Complete | OCR with Gemini Vision |
| **Multi-Evaluator Grading** | ✅ Complete | 3 AI evaluators with consensus |
| **Grade Override** | ✅ Complete | Teacher feedback loop |
| **Audit Trail** | ✅ Complete | All actions logged with hash |
| **Analytics** | ✅ Complete | Grading trends, CSV export |
| **Worker Pool** | ✅ Complete | Async OCR and grading jobs |
| **Rate Limiting** | ✅ Complete | Per-IP rate limiter |
| **File Storage** | ✅ Complete | MinIO (S3-compatible) |

### Backend - NEEDS ENHANCEMENT 🔄

| Component | Status | Notes |
|-----------|--------|-------|
| Authentication | 🔄 Pending | Replace X-Tenant-ID with Supabase JWT |
| User/Role Tables | 🔄 Pending | Add users, roles, permissions |
| Optional Question Groups | 🔄 Pending | Add question_group to schema |
| Batch Upload Endpoint | 🔄 Pending | Multi-file upload support |
| Answer Key Storage | 🔄 Pending | Store expected answers per question |

### Frontend (Next.js) - TO BE CREATED 📋

| Component | Priority | Description |
|-----------|----------|-------------|
| Auth Pages | P0 | Supabase login/signup/forgot password |
| Dashboard | P0 | Teacher overview, recent exams |
| Exam Creation | P0 | Multi-step wizard, question builder |
| Answer Key Upload | P0 | Per-question expected answers |
| Batch Submission Upload | P0 | Drag-drop, student mapping |
| Grading Review | P0 | Side-by-side comparison, override |
| Optional Q Schema UI | P1 | 11a/11b selection interface |
| Analytics Dashboard | P1 | Charts, trends, export |
| Student Results View | P2 | Per-student breakdown |

---

## Executive Summary

HARaMA (Handwritten Analysis and Mark Allocation) is an AI-powered academic evaluation platform that automates the grading of handwritten exam scripts. By leveraging advanced OCR technology, natural language processing, and intelligent scoring algorithms, HARaMA reduces manual grading time by up to 80% while maintaining evaluation consistency and providing detailed analytics.

**Target Market:** Educational institutions (K-12, higher education), competitive exam boards, tutoring centers  
**Expected Impact:** Grade 500+ answer sheets in under 2 hours (vs. 40+ hours manually)  
**Revenue Model:** SaaS subscription (per-teacher/per-institution), API access for enterprise

---

## 1. Product Vision & Objectives

### 1.1 Vision Statement
To revolutionize academic assessment by making exam evaluation faster, fairer, and more insightful through intelligent automation.

### 1.2 Business Objectives
- **Primary:** Reduce teacher grading workload by 70-85%
- **Secondary:** Provide data-driven insights into student performance patterns
- **Tertiary:** Enable standardized evaluation across multiple evaluators

### 1.3 Success Metrics
| Metric | Target | Measurement Period |
|--------|--------|-------------------|
| Grading Time Reduction | 75% | Per exam cycle |
| Accuracy vs Manual Grading | 90%+ correlation | Continuous |
| Teacher Adoption Rate | 60% active users | 6 months post-launch |
| Student Satisfaction Score | 4.2+/5.0 | Quarterly survey |
| Platform Uptime | 99.5% | Monthly |

---

## 2. User Personas & Use Cases

### 2.1 Primary Personas

#### Persona 1: Professor Sarah (University Lecturer)
- **Age:** 35-50 | **Context:** Teaches 200+ students, conducts 3-4 exams per semester
- **Pain Points:** Spends 60+ hours grading per exam, inconsistent scoring across fatigue
- **Goals:** Fast turnaround, consistent evaluation, time for research
- **Tech Savviness:** Moderate (comfortable with cloud tools)

#### Persona 2: Mr. Ramesh (High School Teacher)
- **Age:** 28-45 | **Context:** Handles 5 sections, 40 students each
- **Pain Points:** Pressure for quick result publishing, managing optional questions
- **Goals:** Accuracy in partial marking, easy result generation
- **Tech Savviness:** Basic to Moderate

#### Persona 3: Exam Board Coordinator
- **Age:** 40-60 | **Context:** Oversees standardized testing for thousands of students
- **Pain Points:** Inter-rater reliability, scalability, audit trails
- **Goals:** Consistent standards, bulk processing, compliance
- **Tech Savviness:** High

### 2.2 Key Use Cases

#### Use Case 1: Standard Exam Evaluation
**Actor:** Teacher  
**Precondition:** Question paper and answer key prepared  
**Flow:**
1. Upload question paper template (parts, marks distribution)
2. Upload answer key with expected keywords/concepts
3. Set correction mode (Strict/Moderate/Easy)
4. Upload scanned student answer sheets (batch)
5. System processes: OCR → Text extraction → Similarity analysis → Mark allocation
6. Review AI-suggested marks and feedback
7. Make manual adjustments if needed
8. Publish results to students

**Success Criteria:** 95% of answers require no manual adjustment

#### Use Case 2: Optional Question Handling
**Actor:** Teacher  
**Precondition:** Exam has choice questions (e.g., "Answer 11a OR 11b")  
**Flow:**
1. Define question paper with optional sections
2. System detects which option student attempted
3. Evaluates only the attempted question
4. Marks others as "Not Attempted"

**Success Criteria:** 100% accurate detection of attempted vs skipped questions

#### Use Case 3: Multi-Evaluator Standardization
**Actor:** Exam Board Coordinator  
**Precondition:** 3+ teachers evaluating same exam  
**Flow:**
1. Coordinator sets master answer key and rubric
2. All evaluators use same correction mode
3. System generates consistency report across evaluators
4. Flags answers with >20% mark variation for review

**Success Criteria:** Inter-rater reliability coefficient > 0.85

---

## 3. Functional Requirements

### 3.1 User Management

#### FR-1.1: Authentication & Authorization
- **Priority:** P0 (Critical)
- Support role-based access: Admin, Teacher, Student, Coordinator
- Firebase Authentication with email/password and OAuth (Google, Microsoft)
- Multi-factor authentication (MFA) for admin accounts
- Custom claims for role-based access control

#### FR-1.2: User Profile Management
- **Priority:** P1 (High)
- Teachers: Name, subject expertise, institution, contact
- Students: Roll number, class/section, photo
- Bulk user import via CSV (500+ users at once)

### 3.2 Question Paper & Answer Key Setup

#### FR-2.1: Question Paper Template Creator
- **Priority:** P0
- Support multi-part structure:
  - Part A: Multiple choice/Short answer (1-2 marks)
  - Part B: Medium answer (3-5 marks)
  - Part C: Long answer (7-10 marks)
  - Optional sections with "Attempt X out of Y" logic
- Define marks distribution per question
- Save as reusable templates

#### FR-2.2: Answer Key Upload
- **Priority:** P0
- Format options:
  - Text entry (for keyword-based evaluation)
  - Document upload (PDF/DOCX) with reference answers
  - Point-based rubrics (e.g., "Concept A: 2 marks, Example: 1 mark")
- Support multiple acceptable answers per question
- Version control for answer key updates

#### FR-2.3: Custom Correction Modes
- **Priority:** P1
- **Strict Mode:** Similarity ≥ 85% for full marks
- **Moderate Mode:** Similarity ≥ 70% for full marks
- **Easy Mode:** Similarity ≥ 55% for full marks
- Custom mode: Teacher defines threshold ranges
- Per-question mode override option

### 3.3 Answer Sheet Processing

#### FR-3.1: Document Upload & Preprocessing
- **Priority:** P0
- Accept formats: PDF, JPEG, PNG (max 10MB per file)
- Batch upload: Up to 100 answer sheets at once
- Auto-rotation and deskewing of scanned images
- Page segmentation by question number
- Error handling: Unreadable pages flagged for re-upload

#### FR-3.2: OCR Integration
- **Priority:** P0
- Use Google Gemini Vision API for OCR
- Minimum confidence threshold: 75%
- Low-confidence sections highlighted for manual review
- Support for:
  - English (primary)
  - Numbers, mathematical symbols
  - Basic diagrams (image recognition)
- Language support roadmap: Hindi, Tamil, Spanish

#### FR-3.3: Text Extraction Output
- **Priority:** P0
- Structured JSON format per question:
```json
{
  "student_id": "2025-CS-101",
  "exam_id": "MID-TERM-CS101",
  "answers": [
    {
      "question_no": "1a",
      "extracted_text": "Newton's first law states...",
      "ocr_confidence": 0.92,
      "word_count": 47
    }
  ]
}
```

### 3.4 AI-Powered Evaluation

#### FR-4.1: Semantic Similarity Engine
- **Priority:** P0
- Use sentence-transformers model: `all-MiniLM-L6-v2`
- Calculate cosine similarity between:
  - Student answer embedding
  - Answer key embedding
- Similarity score: 0.0 (no match) to 1.0 (perfect match)

#### FR-4.2: Relevancy Validation (Optional LLM)
- **Priority:** P2 (Medium)
- Integrate Gemini 3 Flash Preview API for logical consistency check
- Validate: Does answer address the question context?
- Output: Relevancy score (0-1) + brief feedback
- Fallback: Rule-based keyword matching if API unavailable

#### FR-4.3: Mark Calculation Algorithm
- **Priority:** P0
```python
# Pseudocode
base_score = cosine_similarity(student_answer, answer_key)
relevancy_weight = 0.9 if LLM_enabled else 1.0
mode_factor = {
    'strict': 0.9,
    'moderate': 1.0,
    'easy': 1.1
}

final_marks = min(
    max_marks,
    max_marks * base_score * relevancy_weight * mode_factor[correction_mode]
)

# Apply threshold-based cutoffs
if base_score < threshold[correction_mode]:
    final_marks *= 0.5  # Half marks for below-threshold answers
```

#### FR-4.4: Partial Credit Logic
- **Priority:** P1
- Detect partial correctness via keyword matching:
  - All keywords present: Full marks
  - 70-99% keywords: Proportional marks
  - <70% keywords: Similarity-based scoring
- Subject-specific rules (e.g., physics: formula = 40%, derivation = 60%)

### 3.5 Review & Override Interface

#### FR-5.1: Teacher Dashboard
- **Priority:** P0
- Display per student:
  - Question-by-question comparison (Answer Key | Student Answer | AI Marks)
  - Confidence indicators (Green: >85%, Yellow: 70-85%, Red: <70%)
  - Edit marks manually with reason code
  - Add text/voice feedback per question
- Bulk actions: Accept all, Reject all, Accept >threshold

#### FR-5.2: Audit Trail
- **Priority:** P1
- Log all mark changes: Original AI score, Modified score, Timestamp, Teacher ID
- Export audit log as CSV
- Display on student result: "AI Suggested: 7/10, Final: 8/10 (Reviewed)"

### 3.6 Reporting & Analytics

#### FR-6.1: Individual Student Report
- **Priority:** P0
- PDF/Excel export containing:
  - Student details, Total marks, Grade
  - Question-wise breakdown with feedback
  - Comparison with class average
  - Areas of improvement (keyword gap analysis)

#### FR-6.2: Class Analytics Dashboard
- **Priority:** P1
- Visualizations:
  - Score distribution histogram
  - Question difficulty analysis (avg. score per question)
  - Time-to-grade metrics (AI vs Manual time)
  - Top/bottom performers
- Export charts as PNG/PDF

#### FR-6.3: Institutional Insights (Enterprise)
- **Priority:** P2
- Cross-course performance trends
- Teacher grading pattern consistency
- Predictive analytics: At-risk students identification
- Curriculum effectiveness heatmaps

---

## 4. Non-Functional Requirements

### 4.1 Performance

| Requirement | Specification | Priority |
|-------------|--------------|----------|
| OCR Processing Time | <45 seconds per 10-page answer sheet | P0 |
| Similarity Calculation | <5 seconds per answer | P0 |
| Concurrent Users | Support 200 simultaneous uploads | P1 |
| Page Load Time | <2 seconds for dashboard | P1 |
| Batch Processing | 100 answer sheets in <20 minutes | P0 |

### 4.2 Scalability
- Horizontal scaling via Railway auto-scaling
- Firebase Firestore auto-sharding for documents
- CDN for static assets (Vercel Edge Network)
- Queue-based processing (Celery + Redis Labs) for OCR jobs

### 4.3 Security

#### Data Protection
- **Encryption at Rest:** Firebase default encryption
- **Encryption in Transit:** TLS 1.3 for all API calls
- **Data Retention:** Answer sheets auto-deleted after 365 days (configurable)
- **Access Control:** Firebase Security Rules with role-based access

#### Compliance
- FERPA compliant (US education data privacy)
- GDPR ready (EU - right to deletion, data portability)
- Firebase compliance certifications (SOC 2, ISO 27001)

### 4.4 Availability & Reliability
- **Uptime SLA:** 99.0% (leveraging Railway and Firebase uptime)
- **Backup:** Firebase automatic daily backups
- **Disaster Recovery:** Firebase multi-region replication
- **Monitoring:** Firebase Analytics + Google Cloud Logging

### 4.5 Usability
- **Accessibility:** WCAG 2.1 AA compliance
- **Mobile Responsiveness:** Full functionality on tablets (teachers), read-only on phones
- **Browser Support:** Chrome 90+, Firefox 88+, Safari 14+, Edge 90+
- **Onboarding:** Interactive tutorial for first-time users (<10 mins)

### 4.6 Maintainability
- **Code Coverage:** Min 80% unit test coverage
- **Documentation:** API docs (FastAPI auto-generated), admin guides, user manuals
- **Logging:** Firebase Cloud Logging with 30-day retention
- **Monitoring:** Firebase Performance Monitoring

---

## 5. System Architecture

### 5.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     CLIENT LAYER                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ React Web App│  │ Mobile Web   │  │ Admin Panel  │      │
│  │ (Vercel)     │  │ (Responsive) │  │              │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ▼ HTTPS/REST
┌─────────────────────────────────────────────────────────────┐
│                  APPLICATION LAYER                           │
│  ┌───────────────────────────────────────────────────────┐  │
│  │  FastAPI Backend (Python 3.11+) - Railway            │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐           │  │
│  │  │ Firebase │  │ Upload   │  │ Grading  │           │  │
│  │  │ Auth API │  │ Handler  │  │ API      │           │  │
│  │  └──────────┘  └──────────┘  └──────────┘           │  │
│  │                                                       │  │
│  │  ┌────────────────────────────────────────────────┐  │  │
│  │  │   Background Job Queue (Celery + Redis Labs)   │  │  │
│  │  │  ┌─────────┐  ┌─────────┐  ┌─────────┐        │  │  │
│  │  │  │ OCR Job │  │ NLP Job │  │ Report  │        │  │  │
│  │  │  │         │  │         │  │ Gen Job │        │  │  │
│  │  │  └─────────┘  └─────────┘  └─────────┘        │  │  │
│  │  └────────────────────────────────────────────────┘  │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    AI/ML SERVICES                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Gemini Vision│  │ Transformers │  │ Gemini 3     │      │
│  │ (OCR)        │  │ (Sentence    │  │ Flash Preview│      │
│  │              │  │ Embeddings)  │  │ (Optional)   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    DATA LAYER                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Firebase     │  │ Redis Labs   │  │ Cloudinary   │      │
│  │ Firestore    │  │ (Cache)      │  │ (Files)      │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

### 5.2 Data Models

#### Why Cloud Firestore?

| Feature | Cloud Firestore | Realtime Database | Rationale for HARaMA |
|---------|----------------|-------------------|---------------------|
| Data Structure | Collections → Documents (JSON-like) | Single JSON tree | ✅ Structured exam/submission data |
| Query Power | Complex queries, indexing, filters | Limited querying | ✅ Filter by teacher, status, date |
| Scalability | Horizontally scalable | Limited scalability | ✅ Handle growing institutions |
| Pricing | Per read/write | Per bandwidth | ✅ Better for frequent small queries |
| Offline Support | Built-in | Built-in | ✅ Teachers can work offline |

**Decision:** Cloud Firestore + Firebase Storage for files

#### Storage Strategy

| Data Type | Storage Solution | Reason |
|-----------|-----------------|--------|
| User profiles, exam metadata | **Firestore** | Structured, queryable data |
| Question papers (PDFs) | **Firebase Storage** | Large files (5-10MB) |
| Answer sheets (images/PDFs) | **Firebase Storage** | Large files (1-3MB per page) |
| Extracted OCR text | **Firestore** | Text strings (<1KB), fast retrieval |
| AI analysis results | **Firestore** | Scores, feedback, marks |
| Cropped answer images | **Firebase Storage** | Individual question images |

#### Firebase Storage Structure
```
harama-storage/
├── question-papers/{examId}/{filename}.pdf
├── answer-keys/{examId}/{filename}.pdf
├── submissions/{examId}/{studentId}/{page}.jpg
├── extracted-texts/{examId}/{studentId}/ocr.json
└── reports/{examId}/{studentId}/report.pdf
```

#### Firebase Firestore Collections Structure

**users** (Collection)
```javascript
users/{userId}
{
  // Basic Info
  email: "teacher@school.edu",
  name: "John Doe",
  role: "teacher", // enum: "admin", "teacher", "student", "coordinator"
  
  // Authentication
  phone: "+91-9876543210",
  profilePhotoUrl: "firebase_storage_url", // optional
  
  // Role-specific fields
  // For Teachers:
  institutionId: "inst_abc123",
  subjectExpertise: ["Physics", "Chemistry"],
  teacherCode: "TCH-2025-001",
  
  // For Students:
  studentId: "STU-2025-101", // Roll number
  class: "12th",
  section: "A",
  
  // Subscription (for teachers)
  subscriptionTier: "free", // enum: "free", "pro", "institution", "enterprise"
  subscriptionExpiresAt: Timestamp, // null for free tier
  monthlyUsage: {
    answersheetsProcessed: 15,
    limit: 50 // based on tier
  },
  
  // Metadata
  createdAt: Timestamp,
  updatedAt: Timestamp,
  lastLoginAt: Timestamp,
  isActive: true
}
```

**institutions** (Collection)
```javascript
institutions/{institutionId}
{
  name: "ABC High School",
  code: "CIT123",
  address: "123 Main St, City, State, ZIP",
  contactEmail: "admin@school.edu",
  contactPhone: "+91-1234567890",
  
  // Subscription
  subscriptionTier: "institution",
  totalTeachers: 50,
  totalStudents: 2000,
  
  // Branding (for paid plans)
  logoUrl: "firebase_storage_url",
  primaryColor: "#1E40AF",
  
  // Metadata
  createdBy: "teacherUid",
  createdAt: Timestamp,
  isActive: true
}
```

**exams** (Collection)
```javascript
exams/{examId}
{
  // Basic Info
  title: "AI - Mid Term Exam 2025",
  description: "Chapter 1-5 coverage",
  subject: "Artificial Intelligence",
  class: "12th",
  section: "A",
  
  // Ownership
  createdBy: "userId_of_teacher",
  institutionId: "inst_abc123",
  
  // Exam Configuration
  totalMarks: 100,
  totalQuestions: 10,
  duration: 180, // minutes
  examDate: Timestamp,
  
  // Question Paper (stored in Firebase Storage)
  questionPaperUrl: "https://firebasestorage.googleapis.com/.../question_paper.pdf",
  questionPaperStoragePath: "question-papers/exam123/paper.pdf",
  
  // Grading Configuration
  correctionMode: "moderate", // enum: "strict", "moderate", "easy", "custom"
  customThreshold: null, // only if correctionMode is "custom"
  
  // Optional Questions Configuration
  hasOptionalQuestions: true,
  optionalSections: [
    {
      section: "Part C",
      instruction: "Answer any 2 out of 3",
      questions: ["8", "9", "10"],
      requiredCount: 2
    }
  ],
  
  // Status
  status: "published", // enum: "draft", "published", "grading", "completed"
  
  // Statistics (updated as submissions come in)
  stats: {
    totalSubmissions: 45,
    gradingCompleted: 40,
    averageMarks: 72.5,
    highestMarks: 95,
    lowestMarks: 42
  },
  
  // Metadata
  createdAt: Timestamp,
  updatedAt: Timestamp,
  publishedAt: Timestamp
}
```

**exams/{examId}/answerKeys** (Subcollection)
```javascript
exams/{examId}/answerKeys/{answerKeyId}
{
  questionNo: "1a",
  questionText: "Explain Newton's First Law of Motion",
  expectedAnswer: "Newton's first law states that an object at rest stays at rest...",
  
  // File reference (if uploaded as PDF/image)
  fileUrl: "https://firebasestorage.googleapis.com/.../answer-key.pdf",
  extractedText: "AI is the science of...", // if OCR extracted from uploaded file
  
  // Keywords for partial credit
  keywords: [
    { word: "inertia", weight: 0.2 },
    { word: "unbalanced force", weight: 0.3 },
    { word: "rest", weight: 0.15 },
    { word: "motion", weight: 0.15 }
  ],
  
  // Marks allocation
  maxMarks: 5,
  partialCreditEnabled: true,
  
  // Subject-specific rules
  evaluationRules: {
    formulaWeight: 0.4, // for math/physics
    explanationWeight: 0.6
  },
  
  // Alternative acceptable answers
  alternativeAnswers: [
    "An object continues in its state of rest or uniform motion unless compelled by an external force"
  ],
  
  // Metadata
  uploadedAt: Timestamp,
  updatedAt: Timestamp
}
```

**exams/{examId}/submissions** (Subcollection)
```javascript
exams/{examId}/submissions/{submissionId}
{
  // Student Info
  studentId: "userId_of_student",
  studentName: "Jane Smith", // denormalized for quick access
  studentRollNo: "STU-2025-101",
  
  // Submission Details
  submittedAt: Timestamp,
  answerSheetUrl: "https://firebasestorage.googleapis.com/.../answer_sheet.pdf",
  answerSheetStoragePath: "submissions/exam123/student456/sheet.pdf",
  totalPages: 12,
  
  // Processing Status
  status: "evaluated", // enum: "uploaded", "processing", "ocr_completed", "grading_completed", "reviewed", "published"
  processingStartedAt: Timestamp,
  processingCompletedAt: Timestamp,
  
  // OCR Job Info
  ocrJobId: "celery_task_id_123",
  ocrStatus: "completed", // enum: "pending", "processing", "completed", "failed"
  ocrError: null, // error message if failed
  
  // Grading Results
  totalMarks: 87.5,
  totalMaxMarks: 100,
  percentage: 87.5,
  grade: "A", // calculated based on percentage
  
  // AI Analysis Summary
  analyzedResult: {
    similarityScore: 0.87,
    marks: 87.5,
    comments: "Good overall performance"
  },
  
  // Review Status
  isReviewed: true,
  reviewedBy: "userId_of_teacher",
  reviewedAt: Timestamp,
  reviewNotes: "Good overall performance, check calculation in Q5",
  
  // Publishing
  isPublished: false,
  publishedAt: null,
  
  // Metadata
  createdAt: Timestamp,
  updatedAt: Timestamp
}
```

**exams/{examId}/submissions/{submissionId}/answers** (Subcollection)
```javascript
exams/{examId}/submissions/{submissionId}/answers/{answerId}
{
  questionNo: "1a",
  
  // Extracted Text from OCR
  extractedText: "Newton's first law tells us that objects at rest stay at rest unless a force acts on them.",
  ocrConfidence: 0.92, // 0-1 scale
  wordCount: 24,
  
  // Image of answer (cropped from full answer sheet)
  answerImageUrl: "https://firebasestorage.googleapis.com/.../q1a_crop.jpg",
  answerImageStoragePath: "submissions/exam123/student456/q1a.jpg",
  
  // AI Evaluation
  similarityScore: 0.78, // cosine similarity with answer key
  relevancyScore: 0.85, // from LLM (optional)
  keywordsMatched: ["rest", "motion", "force"],
  keywordsMissed: ["inertia", "unbalanced force"],
  keywordCoverage: 0.5, // 3/6 keywords matched
  
  // Marks
  maxMarks: 5,
  aiSuggestedMarks: 3.9,
  finalMarks: 4.0, // after teacher review
  partialCreditApplied: true,
  
  // Teacher Feedback
  feedback: "Good understanding, but mention 'inertia' and be more specific about 'unbalanced force'",
  aiGeneratedFeedback: "Your answer demonstrates understanding of the concept but lacks key terminology.",
  
  // Review History
  isReviewed: true,
  reviewedBy: "userId_of_teacher",
  reviewedAt: Timestamp,
  wasModified: true, // true if teacher changed AI marks
  
  // Confidence Indicators
  confidenceLevel: "medium", // enum: "high" (>85%), "medium" (70-85%), "low" (<70%)
  requiresManualReview: false, // flagged if confidence < 70%
  
  // Metadata
  createdAt: Timestamp,
  updatedAt: Timestamp
}
```

**auditLogs** (Collection)
```javascript
auditLogs/{logId}
{
  // What was changed
  entityType: "answer", // enum: "answer", "submission", "exam"
  entityId: "answerId_or_submissionId",
  action: "marks_updated", // enum: "marks_updated", "status_changed", "published", "deleted"
  
  // Change details
  changes: {
    field: "finalMarks",
    oldValue: 3.9,
    newValue: 4.0
  },
  
  // Who made the change
  changedBy: "userId_of_teacher",
  changedByName: "John Doe", // denormalized
  changedByRole: "teacher",
  
  // Why (optional)
  reason: "Student showed understanding of core concept",
  
  // Context
  examId: "exam_abc123",
  submissionId: "sub_xyz789",
  studentId: "student_user_id",
  
  // Metadata
  timestamp: Timestamp,
  ipAddress: "192.168.1.1", // for security
  userAgent: "Mozilla/5.0..."
}
```

**notifications** (Collection)
```javascript
notifications/{notificationId}
{
  // Recipient
  userId: "userId_of_recipient",
  
  // Notification details
  type: "result_published", // enum: "result_published", "grading_completed", "submission_received", "system_alert"
  title: "Your exam results are published",
  message: "Results for Midterm Physics Exam 2025 are now available",
  
  // Related entities
  examId: "exam_abc123",
  submissionId: "sub_xyz789",
  
  // Action
  actionUrl: "/submissions/sub_xyz789",
  actionText: "View Results",
  
  // Status
  isRead: false,
  readAt: null,
  
  // Metadata
  createdAt: Timestamp,
  expiresAt: Timestamp // auto-delete after 30 days
}
```

**templates** (Collection)
```javascript
templates/{templateId}
{
  name: "CBSE Physics 12th Standard Template",
  description: "Standard CBSE pattern with 3 sections",
  subject: "Physics",
  class: "12th",
  
  // Template structure
  sections: [
    {
      name: "Part A - Short Answers",
      marks: 20,
      questions: 10,
      marksPerQuestion: 2
    },
    {
      name: "Part B - Medium Answers",
      marks: 30,
      questions: 6,
      marksPerQuestion: 5
    }
  ],
  
  // Usage
  createdBy: "userId_of_teacher",
  isPublic: true, // available to all users
  usageCount: 145,
  
  // Metadata
  createdAt: Timestamp,
  updatedAt: Timestamp
}
```

**analytics** (Collection - Daily Aggregated)
```javascript
analytics/{analyticsId} // format: "YYYY-MM-DD_scopeId"
{
  date: "2025-11-09",
  
  // Scope
  scope: "teacher", // enum: "teacher", "institution", "system"
  scopeId: "userId_of_teacher",
  
  // Metrics
  totalSubmissions: 45,
  totalAnswersheetsProcessed: 45,
  totalMarksAwarded: 3825,
  averageMarks: 85,
  
  // Processing stats
  averageOcrTime: 32, // seconds
  averageGradingTime: 8, // seconds per answer
  totalProcessingTime: 1440, // seconds
  
  // Quality metrics
  averageOcrConfidence: 0.89,
  manualReviewRate: 0.15, // 15% required manual review
  markModificationRate: 0.08, // 8% marks were changed by teacher
  
  // Metadata
  createdAt: Timestamp
}
```

**systemConfig** (Collection)
```javascript
systemConfig/settings
{
  // Free tier limits
  freeTierLimits: {
    answerSheetsPerMonth: 50,
    examsPerMonth: 2,
    storageGB: 1
  },
  
  // Correction mode thresholds
  correctionModes: {
    strict: { threshold: 0.85, passingScore: 0.90 },
    moderate: { threshold: 0.70, passingScore: 0.80 },
    easy: { threshold: 0.55, passingScore: 0.70 }
  },
  
  // Feature flags
  features: {
    llmRelevancyCheck: false, // disabled in free tier
    diagramEvaluation: false, // coming soon
    plagiarismDetection: false // coming soon
  },
  
  // System status
  maintenanceMode: false,
  
  // Updated
  updatedAt: Timestamp
}
```

### 5.3 API Architecture

**Base URL:** `https://api.harama.ai/v1`

#### Authentication (Firebase Auth Integration)
```
POST /auth/register → Creates Firebase user + Firestore document
POST /auth/login → Firebase custom token generation
POST /auth/refresh-token → Firebase token refresh
GET  /auth/me → Get current user profile
```

#### Exam Management
```
POST   /exams → Create new exam
GET    /exams → List all exams (filtered by teacher)
GET    /exams/{exam_id}
PUT    /exams/{exam_id}/settings
DELETE /exams/{exam_id}
```

#### Answer Key Management
```
POST /exams/{exam_id}/answer-keys → Upload answer key
GET  /exams/{exam_id}/answer-keys
PUT  /exams/{exam_id}/answer-keys/{key_id}
```

#### Submission & Grading
```
POST /exams/{exam_id}/submissions → Upload answer sheets (multipart/form-data to Cloudinary)
GET  /submissions/{submission_id}/status → Check processing status
POST /submissions/{submission_id}/analyze → Trigger AI evaluation
GET  /submissions/{submission_id}/results → Get graded results
PUT  /answers/{answer_id} → Override marks
```

#### Reports
```
GET  /exams/{exam_id}/analytics → Class-level analytics
GET  /submissions/{submission_id}/report → Individual student report (PDF)
POST /exams/{exam_id}/bulk-report → Generate reports for all students
```

---

## 6. Technology Stack (UPDATED - Actual Implementation)

### 6.1 Frontend (To Be Implemented)

| Component | Technology | Justification |
|-----------|-----------|---------------|
| Framework | Next.js 14+ (App Router) | Server components, file-based routing |
| UI Library | Tailwind CSS + shadcn/ui | Rapid prototyping, consistency |
| State Management | TanStack Query (React Query) | Server state caching, mutations |
| File Upload | react-dropzone | Drag-drop, batch validation |
| Charts | Recharts | Rich visualizations, free |
| PDF Generation | jsPDF / react-pdf | Client-side report generation |
| Auth UI | Supabase Auth Helpers | Pre-built components |
| Hosting | Vercel | **Free tier:** Unlimited projects, 100GB bandwidth/month |

### 6.2 Backend (IMPLEMENTED ✅)

| Component | Technology | Status | Notes |
|-----------|-----------|--------|-------|
| Framework | **Go 1.21+** with Chi Router | ✅ Done | High performance, statically typed |
| ORM | **Bun (uptrace/bun)** | ✅ Done | Fast PostgreSQL ORM |
| Database | **PostgreSQL** (via Supabase) | ✅ Done | 5 migrations implemented |
| Authentication | **Supabase Auth** | 🔄 Pending | Currently using X-Tenant-ID header |
| File Storage | **MinIO** (S3-compatible) | ✅ Done | Can switch to Supabase Storage |
| Worker Pool | **Custom Go Worker Pool** | ✅ Done | 5 workers, 100 job buffer |
| Rate Limiting | **Token Bucket (in-memory)** | ✅ Done | 50 req/min, 100 burst |

### 6.3 AI/ML (IMPLEMENTED ✅)

| Component | Technology | Status | Notes |
|-----------|-----------|--------|-------|
| OCR | **Google Gemini Vision (gemini-3-flash-preview)** | ✅ Done | Handwriting extraction |
| Grading | **Google Gemini (gemini-3-flash-preview)** | ✅ Done | Multi-evaluator consensus |
| Evaluators | 3 AI Personas | ✅ Done | rubric_enforcer, reasoning_validator, structural_analyzer |
| Partial Credit | **Custom Go Engine** | ✅ Done | Rubric-based scoring with penalties |
| Confidence | **Variance Calculator** | ✅ Done | Auto-escalation on high variance |

### 6.4 Database Schema (IMPLEMENTED ✅)

```
Tables:
├── tenants              # Multi-tenant support
├── exams                # Exam metadata (title, subject, tenant)
├── questions            # Question text, points, answer_type
├── rubrics              # Full credit criteria, partial rules, common mistakes
├── submissions          # Student answers, OCR results, processing status
├── grades               # AI scores, confidence, reasoning, status
├── escalations          # High-variance cases for review
├── feedback_events      # Teacher override tracking
└── audit_log            # All actions with hash chain
```

### 6.5 API Endpoints (IMPLEMENTED ✅)

```
Exam Management:
  POST   /api/v1/exams                          # Create exam
  GET    /api/v1/exams                          # List exams (by tenant)
  GET    /api/v1/exams/{id}                     # Get exam with questions
  POST   /api/v1/exams/{id}/questions           # Add question
  PUT    /api/v1/questions/{id}/rubric          # Set rubric

Submissions:
  POST   /api/v1/exams/{id}/submissions         # Create submission
  GET    /api/v1/submissions/{id}               # Get status + OCR results
  POST   /api/v1/submissions/{id}/trigger-grading  # Start AI grading

Grading:
  GET    /api/v1/submissions/{id}/grades        # Get all grades
  POST   /api/v1/.../override                   # Teacher override

Feedback:
  GET    /api/v1/.../feedback                   # Get student feedback
  GET    /api/v1/questions/{id}/analysis        # Pattern analysis
  POST   /api/v1/questions/{id}/adapt-rubric    # AI rubric refinement

Analytics:
  GET    /api/v1/analytics/grading-trends       # Stats per question
  POST   /api/v1/exams/{id}/export              # CSV export

Audit:
  GET    /api/v1/audit/{id}                     # Audit logs for entity
```

### 6.6 DevOps & Infrastructure

| Component | Technology | Justification |
|-----------|-----------|---------------|
| Backend Hosting | Railway / Fly.io | Go binary deployment |
| Frontend Hosting | Vercel | **Free tier:** Edge CDN, zero config |
| Database | **Supabase PostgreSQL** | **Free tier:** 500MB, RLS, Auth built-in |
| File Storage | **Supabase Storage** | **Free tier:** 1GB, CDN included |
| CI/CD | GitHub Actions | **Free:** 2,000 minutes/month |
| Monitoring | Supabase Dashboard | Built-in query stats |

---

## 7. Free Tier Limitations & Scaling Strategy

### 7.1 Current Free Tier Limits

| Service | Free Tier Limit | Usage Estimate (100 teachers) |
|---------|----------------|-------------------------------|
| Railway | $5 credit/month (~500 hours) | ~350 hours (with sleep mode) |
| Firebase Firestore | 50K reads/day, 20K writes/day, 1GB storage | ~30K reads/day, ~8K writes/day, ~200MB storage |
| Firebase Storage | 5GB storage, 1GB/day download | ~3GB storage (answer sheets), ~500MB/day download |
| Cloudinary (Backup) | 25GB storage, 25GB bandwidth/month | Use only if Firebase Storage exceeded |
| Gemini Vision | 60 requests/minute | ~20 requests/minute peak |
| Vercel | 100GB bandwidth/month | ~40GB/month |
| Redis Labs | 30MB storage | ~15MB (cache only) |

### 7.2 Optimization Strategies

**Railway (Backend Hosting):**
- Implement auto-sleep after 15 minutes of inactivity
- Use serverless functions for non-critical endpoints
- Optimize Docker image size (<500MB)

**Firebase Firestore:**
- Cache frequently accessed data in Redis
- Batch read/write operations
- Use pagination for large queries
- Implement data archiving after 90 days

**Firebase Storage:**
- Compress images before upload (target: 1-2MB per answer sheet page)
- Use lazy loading for thumbnails
- Delete temporary files after processing (keep only processed results)
- Move old files to Cloudinary if approaching 5GB limit

**Cloudinary (Backup/Overflow):**
- Use only when Firebase Storage > 4GB
- Store large historical archives
- Serves as CDN for frequently accessed files

**Gemini Vision API:**
- Queue-based processing to avoid rate limits
- Retry with exponential backoff
- Batch similar requests

### 7.3 Scaling Roadmap

#### Phase 1: Free Tier (0-100 teachers)
- **Status:** Fully operational on free services
- **Monthly Cost:** $0
- **Capacity:** 500 answer sheets/month

#### Phase 2: Starter Paid (100-500 teachers)
- **Upgrades:**
  - Railway Pro: $20/month
  - Firebase Blaze (pay-as-you-go): ~$30/month
- **Monthly Cost:** ~$50/month
- **Capacity:** 5,000 answer sheets/month

#### Phase 3: Growth (500-2000 teachers)
- **Upgrades:**
  - Railway Pro: $50/month (increased resources)
  - Cloudinary Plus: $89/month
  - Firebase Blaze: ~$150/month
- **Monthly Cost:** ~$290/month
- **Capacity:** 25,000 answer sheets/month

#### Phase 4: Enterprise (2000+ teachers)
- **Migration Path:**
  - Consider AWS/GCP with reserved instances
  - Self-hosted OCR infrastructure
  - Enterprise Firebase plan
- **Monthly Cost:** ~$1,500+/month
- **Capacity:** Unlimited (with proper infrastructure)

---

## 8. Development Roadmap

### Phase 1: Foundation (Weeks 1-4)
**Goal:** MVP infrastructure and authentication

**Deliverables:**
- [ ] Project setup (frontend + backend repos)
- [ ] Firebase Firestore schema design
- [ ] Firebase Authentication integration (email + Google OAuth)
- [ ] Basic UI: Login, Dashboard skeleton
- [ ] Cloudinary integration for file uploads
- [ ] CI/CD pipeline setup (GitHub Actions)

**Team:** 2 Full-stack developers

### Phase 2: Question Paper & Answer Key Setup (Weeks 5-7)
**Goal:** Enable teachers to define exams

**Deliverables:**
- [ ] Question paper template builder UI
- [ ] Answer key upload (text + document formats)
- [ ] Correction mode selector
- [ ] API: CRUD for exams and answer keys
- [ ] Template library (10+ common formats)

**Team:** 1 Frontend + 1 Backend developer

### Phase 3: OCR Integration (Weeks 8-10)
**Goal:** Extract text from handwritten answers

**Deliverables:**
- [ ] Google Gemini Vision API integration
- [ ] File upload handler (batch processing via Cloudinary)
- [ ] Image preprocessing pipeline
- [ ] Celery task queue for async OCR (Redis Labs)
- [ ] Progress indicator UI (real-time updates)
- [ ] Error handling: Low-confidence detection

**Team:** 1 Backend + 1 ML engineer

### Phase 4: AI Evaluation Engine (Weeks 11-14)
**Goal:** Automated mark allocation

**Deliverables:**
- [ ] Sentence embedding model deployment (local)
- [ ] Cosine similarity calculation
- [ ] Mark calculation algorithm (all 3 modes)
- [ ] Gemini 3 Flash Preview relevancy check (optional, v1.1)
- [ ] Partial credit logic
- [ ] Performance optimization (<5s per answer)

**Team:** 1 ML engineer + 1 Backend developer

### Phase 5: Review Interface (Weeks 15-17)
**Goal:** Teacher oversight and manual corrections

**Deliverables:**
- [ ] Side-by-side comparison view (Key vs Student)
- [ ] Mark override interface
- [ ] Bulk actions (Accept/Reject)
- [ ] Feedback text editor
- [ ] Audit trail logging (Firestore collection)
- [ ] Confidence-based highlighting

**Team:** 2 Frontend developers

### Phase 6: Reporting & Analytics (Weeks 18-20)
**Goal:** Insights and export functionality

**Deliverables:**
- [ ] Individual student PDF report (jsPDF)
- [ ] Class analytics dashboard (Recharts)
- [ ] Excel/CSV export
- [ ] Email notifications (Firebase Cloud Functions)
- [ ] Performance visualizations
- [ ] Downloadable answer sheets with annotations

**Team:** 1 Frontend + 1 Backend developer

### Phase 7: Testing & Optimization (Weeks 21-23)
**Goal:** Production-ready stability

**Deliverables:**
- [ ] Load testing (200 concurrent users)
- [ ] Security audit (Firebase Security Rules)
- [ ] Accessibility testing (WCAG 2.1)
- [ ] Beta testing with 5 pilot schools
- [ ] Bug fixes and performance tuning
- [ ] Documentation: API, User guides

**Team:** Full team + QA engineer

### Phase 8: Launch (Week 24)
**Goal:** Public release

**Deliverables:**
- [ ] Production deployment (Railway + Vercel)
- [ ] Firebase Analytics dashboards live
- [ ] Marketing website
- [ ] Onboarding video tutorials
- [ ] Customer support setup (Intercom/Zendesk)
- [ ] Post-launch feedback collection

**Team:** Full team + Marketing

---

## 9. Pricing Strategy

### 9.1 Subscription Tiers

| Tier | Target | Price (USD/month) | Features |
|------|--------|-------------------|----------|
| **Free** | Individual teachers | $0 | 50 answer sheets/month, 2 exams, Basic reports |
| **Pro** | Teachers | $29 | 500 sheets/month, Unlimited exams, Advanced analytics, Priority OCR |
| **Institution** | Schools/Colleges | $299/100 teachers | Unlimited sheets, Multi-evaluator, Custom branding, API access |
| **Enterprise** | Exam boards | Custom | White-label, SLA, Dedicated support, On-premise option |

### 9.2 Add-ons
- Additional answer sheets: $0.10/sheet
- LLM-based relevancy check: $0.05/answer
- Long-term storage (>1 year): $50/TB/year

---

## 10. Risk Management

| Risk | Impact | Probability | Mitigation |
|------|--------|-------------|------------|
| OCR accuracy <85% for cursive handwriting | High | Medium | Fallback to manual review, Use Gemini's multi-modal capabilities |
| Gemini API rate limiting during peak | High | Medium | Queue-based processing, implement exponential backoff |
| Railway sleep time affecting UX | Medium | High | Implement warm-up pings, use serverless for critical endpoints |
| Firebase free tier exceeded | Medium | Medium | Implement caching, data archiving, monitor usage closely, overflow to Cloudinary |
| Firebase Storage limits (5GB) | Medium | Medium | Compress files, archive old data, use Cloudinary for overflow |
| Teachers resist AI grading | Medium | Medium | Emphasize "AI-assisted" not "replacement", Gradual rollout |
| Data breach | Critical | Low | Firebase Security Rules, encryption, compliance certifications |
| Network egress limits | Low | Medium | Image compression, lazy loading, CDN caching |

---

## 11. Success Criteria & KPIs

### 11.1 Product KPIs (First Year)

| KPI | Target | Tracking Method |
|-----|--------|----------------|
| Active Teachers | 1,000 | Firebase Analytics |
| Answer Sheets Processed | 100,000 | Firestore query |
| Average Grading Time Reduction | 75% | User surveys + time logs |
| Teacher Satisfaction (NPS) | 50+ | Quarterly survey |
| System Accuracy vs Manual | 90% correlation | A/B testing (sample of 500 sheets) |
| Churn Rate | <10% monthly | Subscription analytics |

### 11.2 Technical KPIs

| Metric | Target | Tool |
|--------|--------|------|
| API Uptime | 99.0% | Firebase Performance Monitoring |
| P95 Response Time | <500ms | Firebase Performance Monitoring |
| OCR Processing Speed | <45s per 10 pages | Custom logging |
| Error Rate | <0.5% | Sentry |
| Customer Support Tickets | <5% of active users/month | Zendesk |

---

## 12. Assumptions & Dependencies

### 12.1 Assumptions
- Students write legibly enough for 75%+ OCR confidence
- Teachers have reliable internet (min 5 Mbps upload for scans)
- Answer keys are comprehensive (not just "Refer textbook")
- Schools can digitize answer sheets (scanner/phone camera available)
- Teachers are willing to review AI suggestions before finalizing grades

### 12.2 Dependencies
- **Google Gemini API:** Pricing stability, API availability, rate limits
- **Hugging Face Models:** Model weights remain open-source
- **Cloud Providers:** Railway/Vercel/Firebase uptime, no major price hikes
- **Regulatory:** No sudden education data privacy law changes
- **Internet Connectivity:** Teachers and students have stable internet access

### 12.3 Technical Constraints
- Railway free tier sleep time (app may take 5-10s to wake up)
- Gemini API rate limits (60 requests/minute) may cause delays during peak hours
- Firebase Firestore read/write limits on free tier
- Cloudinary bandwidth limits for image uploads

---

## 13. Open Questions

1. **Multi-language support priority:** Should we launch with English-only or delay 3 months for Hindi/Tamil?
   - **Recommendation:** Launch with English, add Hindi/Tamil in v1.1 (3 months post-launch)

2. **Diagram evaluation:** How to handle students who draw diagrams (e.g., biology, physics)? Image similarity models?
   - **Recommendation:** Phase 2 feature - Use Gemini Vision for diagram comparison

3. **Plagiarism detection:** Should this be Phase 1 or post-launch enhancement?
   - **Recommendation:** Post-launch (v1.2) - Add similarity checking between student answers

4. **Offline mode:** Do teachers need desktop app for areas with poor connectivity?
   - **Recommendation:** Not in MVP - Monitor user feedback first

5. **Answer sheet format standardization:** Should we enforce QR codes on sheets for auto-identification?
   - **Recommendation:** Optional feature - Use OCR-based student ID detection first

6. **Payment gateway integration:** Which payment providers to support?
   - **Recommendation:** Stripe (international), Razorpay (India-specific)

---

## 14. Appendices

### Appendix A: Competitor Analysis

| Competitor | Strengths | Weaknesses | HARaMA Differentiator |
|-----------|-----------|-----------|----------------------|
| **Gradescope** | Strong in STEM, bubble sheet support | Weak handwriting OCR, expensive ($200+/year per teacher) | Better handwriting recognition, affordable pricing |
| **Turnitin** | Plagiarism detection, large database | No grading automation | Full grading + plagiarism in one platform |
| **Examsoft** | Enterprise features, secure testing | Enterprise-only, very expensive | Accessible to individual teachers, freemium model |
| **Zipgrade** | Simple bubble sheet scanning | Only MCQs, no subjective answers | Handles descriptive answers with AI |

### Appendix B: Sample User Flows

#### Teacher Onboarding Flow
1. Sign up with email or Google OAuth
2. Complete profile (name, subject, institution)
3. Interactive tutorial (5 mins)
   - Upload sample question paper
   - Create answer key
   - Test with sample answer sheet
4. Receive welcome email with resources

#### Student Result Viewing Flow
1. Receive email notification: "Your results are published"
2. Click link → Login with credentials
3. View dashboard:
   - Overall score and grade
   - Question-wise breakdown
   - Teacher feedback
   - Compare with class average
4. Download detailed PDF report

#### Grading Workflow
```
Teacher uploads answer sheets
         ↓
Queue job in Celery
         ↓
Cloudinary stores images
         ↓
Gemini Vision extracts text (OCR)
         ↓
Sentence-transformers calculates similarity
         ↓
Mark calculation algorithm
         ↓
Store results in Firestore
         ↓
Teacher reviews suggestions
         ↓
Manual adjustments (if needed)
         ↓
Publish results to students
         ↓
Email notifications sent
```

### Appendix C: Glossary

- **OCR Confidence:** Probability (0-1) that extracted text is correct
- **Cosine Similarity:** Measure of semantic similarity between two text embeddings (range: -1 to 1, typically 0 to 1 for text)
- **Correction Mode:** Grading strictness level (Strict/Moderate/Easy)
- **Inter-rater Reliability:** Agreement level between multiple evaluators (measured by correlation coefficient)
- **Embedding:** Vector representation of text that captures semantic meaning
- **FAISS:** Facebook AI Similarity Search - library for efficient similarity search
- **Celery:** Distributed task queue for asynchronous job processing
- **Firebase Firestore:** NoSQL document database with real-time sync
- **Railway:** Cloud platform for deploying and scaling applications

### Appendix D: API Rate Limits Summary

| Service | Free Tier Limit | Recommended Buffer |
|---------|----------------|-------------------|
| Gemini Vision | 60 req/min | Use queue to limit to 45 req/min |
| Gemini 3 Flash Preview | 15 req/min | Use queue to limit to 10 req/min |
| Firebase Firestore | 50K reads/day | Monitor at 40K/day, implement caching |
| Cloudinary | 25GB bandwidth/month | Monitor at 20GB/month, compress images |
| Railway | $5 credit/month | Monitor usage weekly, optimize resource usage |

### Appendix E: Security Best Practices

#### Firebase Security Rules Example
```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    // Users can only read their own profile
    match /users/{userId} {
      allow read: if request.auth.uid == userId;
      allow write: if request.auth.uid == userId;
    }
    
    // Teachers can read/write their own exams
    match /exams/{examId} {
      allow read: if request.auth.uid == resource.data.createdBy;
      allow write: if request.auth.uid == request.resource.data.createdBy;
      
      // Students can only read published submissions
      match /submissions/{submissionId} {
        allow read: if request.auth.uid == resource.data.studentId 
                    && resource.data.status == 'published';
      }
    }
    
    // Admin-only access to audit logs
    match /auditLogs/{logId} {
      allow read: if request.auth.token.role == 'admin';
    }
  }
}
```

#### Data Encryption Standards
- **At Rest:** Firebase default AES-256 encryption
- **In Transit:** TLS 1.3 for all API communications
- **Sensitive Data:** Additional encryption for PII (email, phone) using crypto-js
- **API Keys:** Stored in environment variables, never in code
- **Password Storage:** Firebase Auth handles password hashing (bcrypt)

### Appendix F: Performance Optimization Checklist

**Frontend Optimization:**
- [ ] Implement code splitting (React.lazy)
- [ ] Use CDN for static assets (Vercel automatic)
- [ ] Optimize images (WebP format, lazy loading)
- [ ] Minimize bundle size (<500KB initial load)
- [ ] Implement service workers for offline support

**Backend Optimization:**
- [ ] Use connection pooling for Firebase
- [ ] Implement Redis caching for frequent queries
- [ ] Optimize Celery worker count based on load
- [ ] Use batch operations for Firestore writes
- [ ] Implement request rate limiting

**Database Optimization:**
- [ ] Create composite indexes for common queries
- [ ] Use pagination for large result sets
- [ ] Implement data archiving for old records
- [ ] Denormalize frequently accessed data
- [ ] Monitor query performance with Firebase console

### Appendix G: Testing Strategy

#### Unit Testing (Target: 80% coverage)
- **Frontend:** Jest + React Testing Library
- **Backend:** Pytest + FastAPI TestClient
- **Models:** Unittest for ML model predictions

#### Integration Testing
- **API Endpoints:** Test all CRUD operations
- **Firebase Integration:** Test auth, Firestore CRUD, Security Rules
- **Cloudinary Integration:** Test upload, retrieval, deletion
- **Celery Tasks:** Test job queuing and processing

#### End-to-End Testing
- **Cypress:** Critical user flows (login, upload, grading, review)
- **Load Testing:** Artillery or Locust (200 concurrent users)
- **Security Testing:** OWASP ZAP for vulnerability scanning

#### User Acceptance Testing (UAT)
- Beta testing with 5 pilot schools
- Collect feedback via surveys and interviews
- Monitor Firebase Analytics for usage patterns
- Track error rates and performance metrics

### Appendix H: Deployment Checklist

**Pre-Launch:**
- [ ] Complete security audit
- [ ] Set up monitoring and alerts
- [ ] Configure Firebase Security Rules
- [ ] Set up backup and recovery procedures
- [ ] Create user documentation and video tutorials
- [ ] Set up customer support channels
- [ ] Configure error tracking (Sentry)
- [ ] Implement analytics (Firebase Analytics)
- [ ] Load testing completed successfully
- [ ] All critical bugs resolved

**Launch Day:**
- [ ] Deploy to production (Railway + Vercel)
- [ ] Verify all services are running
- [ ] Test critical user flows
- [ ] Monitor error rates and performance
- [ ] Prepare rollback plan
- [ ] Customer support team on standby
- [ ] Social media announcements
- [ ] Email marketing campaign

**Post-Launch (Week 1):**
- [ ] Daily monitoring of system health
- [ ] Collect user feedback
- [ ] Address critical issues immediately
- [ ] Monitor usage against free tier limits
- [ ] Analyze user behavior patterns
- [ ] Schedule retrospective meeting

### Appendix I: Support & Maintenance Plan

#### Tier 1 Support (Users)
- **Response Time:** 24 hours
- **Channels:** Email, in-app chat
- **Common Issues:** Login problems, upload errors, grading questions
- **Resources:** FAQ, video tutorials, knowledge base

#### Tier 2 Support (Technical)
- **Response Time:** 8 hours
- **Issues:** API errors, integration problems, performance issues
- **Team:** Backend developers + DevOps

#### Tier 3 Support (Critical)
- **Response Time:** 2 hours
- **Issues:** Security breaches, data loss, system outages
- **Team:** Full technical team + management

#### Maintenance Schedule
- **Daily:** Automated backups, log review, performance monitoring
- **Weekly:** Security updates, dependency updates, usage reports
- **Monthly:** Feature releases, user feedback review, capacity planning
- **Quarterly:** Major updates, security audits, user surveys

### Appendix J: Marketing & Growth Strategy

#### Launch Strategy
- **Target Audience:** Individual teachers in STEM subjects
- **Channels:** 
  - LinkedIn ads targeting educators
  - Educational technology forums and communities
  - YouTube tutorials and demonstrations
  - Twitter/X education community engagement
  - Product Hunt launch
  - Reddit r/Teachers, r/education

#### Content Marketing
- Blog posts: "How AI is transforming exam grading"
- Case studies from beta testers
- Webinars: "Getting started with HARaMA"
- Email newsletter with tips and best practices

#### Growth Tactics
- Referral program: Free month for each referral
- Educational institution partnerships
- Conference presentations at EdTech events
- Free tier with upgrade prompts
- Community building (Discord/Slack)

#### Metrics for Success
- Customer Acquisition Cost (CAC): Target <$50 per paid user
- Conversion Rate (Free to Paid): Target 15%
- User Retention: Target 85% month-over-month
- Net Promoter Score (NPS): Target 50+

---

## 15. Document Change Log

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0 | Oct 29, 2025 | Initial PRD with AWS/Azure stack | Product Team |
| 2.0 | Nov 09, 2025 | Updated to free tech stack (Firebase, Railway, Gemini, Cloudinary) | Product Team |

---

## 16. Document Approval

| Role | Name | Date | Signature |
|------|------|------|-----------|
| Product Manager | ________________ | __________ | ________________ |
| Engineering Lead | ________________ | __________ | ________________ |
| Design Lead | ________________ | __________ | ________________ |
| CEO/Stakeholder | ________________ | __________ | ________________ |

---

## 17. Next Steps

### Immediate Actions (This Week)
1. Set up GitHub repository (frontend + backend)
2. Create Firebase project and configure authentication
3. Set up Railway account and deploy "Hello World" API
4. Create Vercel project and deploy starter React app
5. Sign up for Cloudinary and Gemini API access
6. Create project timeline in project management tool

### Short-term Goals (Month 1)
1. Complete Phase 1 (Foundation)
2. Recruit beta testers (5 teachers)
3. Design UI mockups in Figma
4. Set up CI/CD pipeline
5. Create developer documentation

### Medium-term Goals (Months 2-6)
1. Complete Phases 2-7
2. Launch beta with pilot schools
3. Collect user feedback and iterate
4. Prepare for public launch
5. Create marketing materials

### Long-term Goals (Year 1)
1. Achieve 1,000 active teachers
2. Process 100,000 answer sheets
3. Expand to 2 additional languages
4. Build community of educators
5. Reach profitability

---

**End of Document**  

**For questions or feedback, contact:**  
- **Email:** product@harama.ai  
- **Website:** https://harama.ai  
- **Documentation:** https://docs.harama.ai  
- **Support:** support@harama.ai

---

**Document Status:** ✅ Ready for Development  
**Last Review Date:** November 09, 2025  
**Next Review Date:** December 09, 2025
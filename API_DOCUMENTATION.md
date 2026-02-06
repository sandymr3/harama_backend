# HARaMA API Documentation for Frontend Developers

## Base URL
```
http://localhost:8080
```

## Authentication
All API requests require a tenant ID header:
```
X-Tenant-ID: 00000000-0000-0000-0000-000000000001
```

---

## 📋 Exams

### Create Exam
```http
POST /api/v1/exams
Content-Type: application/json
X-Tenant-ID: {tenant-id}

{
  "title": "Physics Midterm",
  "subject": "Physics",
  "questions": [
    {
      "question_text": "Calculate F=ma with m=10kg, a=5m/s²",
      "points": 5,
      "answer_type": "short_answer",
      "rubric": {
        "full_credit_criteria": [
          {"description": "Correct formula", "points": 2},
          {"description": "Correct calculation", "points": 2},
          {"description": "Units specified", "points": 1}
        ],
        "partial_credit_rules": [],
        "common_mistakes": [],
        "key_concepts": ["Newton's Second Law"]
      }
    }
  ]
}
```

**Response:**
```json
{
  "id": "exam-uuid",
  "title": "Physics Midterm",
  "subject": "Physics",
  "created_at": "2024-02-06T10:00:00Z"
}
```

---

### List Exams
```http
GET /api/v1/exams
X-Tenant-ID: {tenant-id}
```

**Response:**
```json
[
  {
    "id": "exam-uuid",
    "title": "Physics Midterm",
    "subject": "Physics",
    "created_at": "2024-02-06T10:00:00Z"
  }
]
```

---

### Get Exam Details
```http
GET /api/v1/exams/{exam-id}
X-Tenant-ID: {tenant-id}
```

**Response:**
```json
{
  "id": "exam-uuid",
  "title": "Physics Midterm",
  "subject": "Physics",
  "questions": [
    {
      "id": "question-uuid",
      "question_text": "Calculate F=ma...",
      "points": 5,
      "answer_type": "short_answer",
      "rubric": {...}
    }
  ],
  "created_at": "2024-02-06T10:00:00Z"
}
```

---

### Add Question to Exam
```http
POST /api/v1/exams/{exam-id}/questions
Content-Type: application/json
X-Tenant-ID: {tenant-id}

{
  "question_text": "Explain Newton's Second Law",
  "points": 10,
  "answer_type": "essay"
}
```

---

### Set Rubric for Question
```http
PUT /api/v1/questions/{question-id}/rubric
Content-Type: application/json
X-Tenant-ID: {tenant-id}

{
  "full_credit_criteria": [
    {"description": "Defines force", "points": 3},
    {"description": "Explains relationship", "points": 4},
    {"description": "Provides example", "points": 3}
  ],
  "key_concepts": ["Force", "Mass", "Acceleration"]
}
```

---

## 📄 Submissions

### Upload Submission
```http
POST /api/v1/exams/{exam-id}/submissions
Content-Type: multipart/form-data
X-Tenant-ID: {tenant-id}

Form Data:
- file: [image/pdf file]
- student_id: "student123"
```

**Response:**
```json
{
  "id": "submission-uuid",
  "exam_id": "exam-uuid",
  "student_id": "student123",
  "processing_status": "pending",
  "uploaded_at": "2024-02-06T10:30:00Z"
}
```

---

### Get Submission Status
```http
GET /api/v1/submissions/{submission-id}
X-Tenant-ID: {tenant-id}
```

**Response:**
```json
{
  "id": "submission-uuid",
  "exam_id": "exam-uuid",
  "student_id": "student123",
  "processing_status": "completed",
  "ocr_results": [
    {
      "page_number": 1,
      "raw_text": "F = ma = (10)(5) = 50 N",
      "confidence": 0.90,
      "image_url": "submissions/abc123.png"
    }
  ],
  "answers": [
    {
      "question_id": "question-uuid",
      "text": "F = ma = (10)(5) = 50 N"
    }
  ],
  "uploaded_at": "2024-02-06T10:30:00Z"
}
```

**Status Values:**
- `pending` - Just uploaded
- `processing` - OCR/grading in progress
- `completed` - Ready for review
- `failed` - Error occurred

---

### Trigger Grading
```http
POST /api/v1/submissions/{submission-id}/trigger-grading
X-Tenant-ID: {tenant-id}
```

**Response:**
```json
{
  "message": "Grading started",
  "submission_id": "submission-uuid"
}
```

---

## 📊 Grading

### Get Grades for Submission
```http
GET /api/v1/submissions/{submission-id}/grades
X-Tenant-ID: {tenant-id}
```

**Response:**
```json
[
  {
    "id": "grade-uuid",
    "submission_id": "submission-uuid",
    "question_id": "question-uuid",
    "final_score": 4.5,
    "ai_score": 4.5,
    "confidence": 0.85,
    "reasoning": "Correct formula and calculation. Missing units (-0.5 points).",
    "status": "auto_graded",
    "updated_at": "2024-02-06T10:35:00Z"
  }
]
```

**Grade Status:**
- `pending` - Not graded yet
- `auto_graded` - AI graded with high confidence
- `needs_review` - Low confidence or high variance
- `overridden` - Teacher modified
- `final` - Approved by teacher

---

### Override Grade (Teacher)
```http
POST /api/v1/submissions/{submission-id}/questions/{question-id}/override
Content-Type: application/json
X-Tenant-ID: {tenant-id}

{
  "new_score": 5.0,
  "reason": "Student showed correct understanding, units not critical here"
}
```

**Response:**
```json
{
  "id": "grade-uuid",
  "final_score": 5.0,
  "override_score": 5.0,
  "status": "overridden",
  "updated_at": "2024-02-06T11:00:00Z"
}
```

---

## 💬 Feedback

### Get Student Feedback
```http
GET /api/v1/submissions/{submission-id}/questions/{question-id}/feedback
X-Tenant-ID: {tenant-id}
```

**Response:**
```json
{
  "question_id": "question-uuid",
  "submission_id": "submission-uuid",
  "feedback": "Your calculation was correct and you showed good understanding of F=ma. Remember to always include units in your final answer. Practice: Review unit conversions on page 47.",
  "score": 4.5,
  "max_score": 5,
  "strengths": ["Correct formula", "Accurate calculation"],
  "improvements": ["Include units"]
}
```

---

### Analyze Question Patterns
```http
GET /api/v1/questions/{question-id}/analysis
X-Tenant-ID: {tenant-id}
```

**Response:**
```json
{
  "question_id": "question-uuid",
  "total_submissions": 50,
  "average_score": 3.8,
  "common_mistakes": [
    {
      "mistake": "Missing units",
      "frequency": 15,
      "impact": -0.5
    }
  ],
  "difficulty_rating": "medium"
}
```

---

### Adapt Rubric (Based on Feedback)
```http
POST /api/v1/questions/{question-id}/adapt-rubric
X-Tenant-ID: {tenant-id}
```

**Response:**
```json
{
  "message": "Rubric adapted based on teacher feedback",
  "adjustments": [
    {
      "criterion": "units_specified",
      "old_weight": 1.0,
      "new_weight": 0.5,
      "reason": "Teachers consistently give partial credit"
    }
  ]
}
```

---

## 📈 Analytics

### Get Grading Trends
```http
GET /api/v1/analytics/grading-trends
X-Tenant-ID: {tenant-id}
```

**Response:**
```json
{
  "total_submissions": 150,
  "auto_graded": 120,
  "needs_review": 20,
  "overridden": 10,
  "average_confidence": 0.87,
  "average_score": 7.5,
  "trends": [
    {
      "date": "2024-02-06",
      "submissions": 50,
      "avg_score": 7.8
    }
  ]
}
```

---

### Export Grades (CSV)
```http
POST /api/v1/exams/{exam-id}/export
X-Tenant-ID: {tenant-id}
```

**Response:**
```
Content-Type: text/csv
Content-Disposition: attachment; filename="grades.csv"

student_id,question_1,question_2,question_3,total
student123,5,8,7,20
student456,4,9,8,21
```

---

## 🔍 Audit

### Get Audit Logs
```http
GET /api/v1/audit/{entity-id}
X-Tenant-ID: {tenant-id}
```

**Response:**
```json
[
  {
    "id": "audit-uuid",
    "entity_type": "grade",
    "entity_id": "grade-uuid",
    "event_type": "overridden",
    "actor_id": "teacher-uuid",
    "actor_type": "teacher",
    "changes": {
      "score_before": 4.5,
      "score_after": 5.0,
      "reason": "Student showed understanding"
    },
    "created_at": "2024-02-06T11:00:00Z",
    "hash": "abc123..."
  }
]
```

---

## 🏥 Health Check

### Check API Health
```http
GET /health
```

**Response:**
```
OK
```

---

## Error Responses

All errors follow this format:

```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": {}
}
```

**Common Status Codes:**
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

---

## Example Frontend Integration

### React/Next.js Example

```typescript
// lib/api.ts
const API_BASE = 'http://localhost:8080';
const TENANT_ID = '00000000-0000-0000-0000-000000000001';

export async function createExam(data: any) {
  const response = await fetch(`${API_BASE}/api/v1/exams`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-Tenant-ID': TENANT_ID,
    },
    body: JSON.stringify(data),
  });
  
  if (!response.ok) {
    throw new Error('Failed to create exam');
  }
  
  return response.json();
}

export async function uploadSubmission(examId: string, file: File, studentId: string) {
  const formData = new FormData();
  formData.append('file', file);
  formData.append('student_id', studentId);
  
  const response = await fetch(`${API_BASE}/api/v1/exams/${examId}/submissions`, {
    method: 'POST',
    headers: {
      'X-Tenant-ID': TENANT_ID,
    },
    body: formData,
  });
  
  return response.json();
}

export async function getGrades(submissionId: string) {
  const response = await fetch(`${API_BASE}/api/v1/submissions/${submissionId}/grades`, {
    headers: {
      'X-Tenant-ID': TENANT_ID,
    },
  });
  
  return response.json();
}
```

### Usage in Component

```tsx
'use client'

import { useState } from 'react';
import { uploadSubmission, getGrades } from '@/lib/api';

export function SubmissionUpload({ examId }: { examId: string }) {
  const [file, setFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);
  
  const handleSubmit = async () => {
    if (!file) return;
    
    setLoading(true);
    try {
      const submission = await uploadSubmission(examId, file, 'student123');
      console.log('Uploaded:', submission.id);
      
      // Poll for grades
      const grades = await getGrades(submission.id);
      console.log('Grades:', grades);
    } catch (error) {
      console.error('Error:', error);
    } finally {
      setLoading(false);
    }
  };
  
  return (
    <div>
      <input type="file" onChange={(e) => setFile(e.target.files?.[0] || null)} />
      <button onClick={handleSubmit} disabled={loading}>
        {loading ? 'Uploading...' : 'Submit'}
      </button>
    </div>
  );
}
```

---

## WebSocket (Future)

For real-time grading updates:

```typescript
const ws = new WebSocket('ws://localhost:8080/ws');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  
  if (data.type === 'grading_complete') {
    console.log('Grading done:', data.submission_id);
  }
};
```

---

## Rate Limits

- **60 requests per minute** per tenant
- **Burst:** 100 requests

**Headers:**
```
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 45
X-RateLimit-Reset: 1612345678
```

---

## Support

**Issues:** https://github.com/your-repo/issues  
**Docs:** https://docs.haramma.dev  
**API Status:** https://status.haramma.dev

# HARaMA API Documentation for Frontend Developers

> **Last Updated:** February 7, 2026
> **Backend:** Go (Chi Router + Bun ORM)
> **Database:** PostgreSQL (Supabase)
> **Authentication:** Supabase Auth (JWT)
> **File Storage:** Supabase Storage

---

## Implementation Status

### ✅ Implemented (Backend Ready)
| Feature | Endpoint | Status |
|---------|----------|--------|
| Create Exam | POST /api/v1/exams | ✅ Ready |
| List Exams | GET /api/v1/exams | ✅ Ready |
| Get Exam Details | GET /api/v1/exams/{id} | ✅ Ready |
| Add Question | POST /api/v1/exams/{id}/questions | ✅ Ready |
| Set Rubric | PUT /api/v1/questions/{id}/rubric | ✅ Ready |
| Create Submission | POST /api/v1/exams/{id}/submissions | ✅ Ready |
| Get Submission | GET /api/v1/submissions/{id} | ✅ Ready |
| Trigger Grading | POST /api/v1/submissions/{id}/trigger-grading | ✅ Ready |
| Get Grades | GET /api/v1/submissions/{id}/grades | ✅ Ready |
| Override Grade | POST /api/v1/.../override | ✅ Ready |
| Get Feedback | GET /api/v1/.../feedback | ✅ Ready |
| Pattern Analysis | GET /api/v1/questions/{id}/analysis | ✅ Ready |
| Adapt Rubric | POST /api/v1/questions/{id}/adapt-rubric | ✅ Ready |
| Grading Trends | GET /api/v1/analytics/grading-trends | ✅ Ready |
| Export Grades | POST /api/v1/exams/{id}/export | ✅ Ready |
| Audit Logs | GET /api/v1/audit/{id} | ✅ Ready |
| Health Check | GET /health | ✅ Ready |
| OCR Processing | Background Job (Gemini Vision) | ✅ Ready |
| Multi-Evaluator Grading | Background Job (3 AI Evaluators) | ✅ Ready |
| Escalation on High Variance | Automatic | ✅ Ready |

### 🔄 Needs Frontend Implementation
| Feature | Notes |
|---------|-------|
| User Authentication | Supabase Auth (sign up, login, logout) |
| Exam Question Paper Upload | File upload to Supabase Storage |
| Answer Key Upload | Store expected answers per question |
| Answer Sheet Batch Upload | Multiple files upload with student metadata |
| Optional Question Schema | UI for 11a1+11a2 OR 11b1+11b2 selection |
| Student Management | List students, assign to exams |
| Dashboard & Analytics Views | Charts, reports |
| Real-time Processing Status | Polling or WebSocket |

### ⏳ Backend Enhancement Needed
| Feature | Status |
|---------|--------|
| User/Role Tables | Add migration for users/roles |
| Supabase Auth Middleware | Replace X-Tenant-ID with JWT |
| Optional Question Groups | Add question_group field |
| Batch Submission Endpoint | Enhance for bulk processing |

---

## Base URL
```
http://localhost:8080
```

## Authentication

### Current (Development)
All API requests require a tenant ID header:
```
X-Tenant-ID: 00000000-0000-0000-0000-000000000001
```

### Target (Production with Supabase)
```
Authorization: Bearer <supabase_jwt_token>
```

The JWT token is obtained from Supabase Auth after login. The backend will:
1. Verify the JWT with Supabase
2. Extract user_id and tenant_id from claims
3. Apply Row Level Security (RLS) policies

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

## 📝 Optional Question Schema (Answer 11a OR 11b)

The backend supports optional question groups where students can answer one of multiple options. This is commonly used for questions like "Answer any 2 out of 3" or "Answer 11a OR 11b".

### Question Structure with Groups

When creating questions, use the `question_group` and `group_max_score` fields:

```http
POST /api/v1/exams/{exam-id}/questions
Content-Type: application/json

{
  "question_text": "11a(i) - Explain kinetic energy",
  "points": 5,
  "answer_type": "short_answer",
  "question_group": "Q11_OPTION_A",
  "question_number": "11a1"
}

{
  "question_text": "11a(ii) - Calculate kinetic energy for given mass",
  "points": 5,
  "answer_type": "short_answer", 
  "question_group": "Q11_OPTION_A",
  "question_number": "11a2"
}

{
  "question_text": "11b(i) - Explain potential energy",
  "points": 5,
  "answer_type": "short_answer",
  "question_group": "Q11_OPTION_B",
  "question_number": "11b1"
}

{
  "question_text": "11b(ii) - Calculate potential energy for given height",
  "points": 5,
  "answer_type": "short_answer",
  "question_group": "Q11_OPTION_B",
  "question_number": "11b2"
}
```

### Setting Question Group Rules

```http
POST /api/v1/exams/{exam-id}/question-groups
Content-Type: application/json

{
  "parent_question": "11",
  "groups": [
    {
      "group_id": "Q11_OPTION_A",
      "questions": ["11a1", "11a2"],
      "max_score": 10
    },
    {
      "group_id": "Q11_OPTION_B", 
      "questions": ["11b1", "11b2"],
      "max_score": 10
    }
  ],
  "selection_rule": "SELECT_ONE",
  "total_max_score": 10
}
```

### Frontend UI Requirements

The frontend should display:
1. A visual grouping showing "11a" and "11b" as alternatives
2. Radio selection or automatic detection based on which option the student answered
3. Scoring that only considers the selected option (max of 11a1+11a2 OR 11b1+11b2)

---

## 📦 Batch Submission Upload

For uploading multiple student answer sheets at once:

### Batch Upload Endpoint

```http
POST /api/v1/exams/{exam-id}/submissions/batch
Content-Type: multipart/form-data
X-Tenant-ID: {tenant-id}

Form Data:
- files[]: [file1.pdf, file2.pdf, file3.pdf, ...]
- student_mapping: JSON string mapping filename to student_id
```

**student_mapping example:**
```json
{
  "file1.pdf": "STU001",
  "file2.pdf": "STU002",
  "file3.pdf": "STU003"
}
```

**Response:**
```json
{
  "batch_id": "batch-uuid",
  "total_files": 30,
  "submissions_created": 30,
  "status": "processing",
  "submissions": [
    {"student_id": "STU001", "submission_id": "sub-uuid-1", "status": "pending"},
    {"student_id": "STU002", "submission_id": "sub-uuid-2", "status": "pending"}
  ]
}
```

### Check Batch Status

```http
GET /api/v1/batches/{batch-id}/status
X-Tenant-ID: {tenant-id}
```

**Response:**
```json
{
  "batch_id": "batch-uuid",
  "total": 30,
  "completed": 25,
  "processing": 3,
  "failed": 2,
  "submissions": [...]
}
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

## Example Frontend Integration with Supabase

### Setup Supabase Client

```typescript
// lib/supabase.ts
import { createClient } from '@supabase/supabase-js';

const supabaseUrl = process.env.NEXT_PUBLIC_SUPABASE_URL!;
const supabaseAnonKey = process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!;

export const supabase = createClient(supabaseUrl, supabaseAnonKey);
```

### API Client with Supabase Auth

```typescript
// lib/api.ts
import { supabase } from './supabase';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

async function getAuthHeaders() {
  const { data: { session } } = await supabase.auth.getSession();
  if (!session?.access_token) {
    throw new Error('Not authenticated');
  }
  return {
    'Authorization': `Bearer ${session.access_token}`,
    'Content-Type': 'application/json',
  };
}

export async function createExam(data: CreateExamRequest) {
  const headers = await getAuthHeaders();
  const response = await fetch(`${API_BASE}/api/v1/exams`, {
    method: 'POST',
    headers,
    body: JSON.stringify(data),
  });
  
  if (!response.ok) {
    throw new Error('Failed to create exam');
  }
  
  return response.json();
}

export async function listExams() {
  const headers = await getAuthHeaders();
  const response = await fetch(`${API_BASE}/api/v1/exams`, {
    headers,
  });
  return response.json();
}

export async function uploadSubmission(examId: string, file: File, studentId: string) {
  const { data: { session } } = await supabase.auth.getSession();
  
  const formData = new FormData();
  formData.append('file', file);
  formData.append('student_id', studentId);
  
  const response = await fetch(`${API_BASE}/api/v1/exams/${examId}/submissions`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${session?.access_token}`,
    },
    body: formData,
  });
  
  return response.json();
}

export async function uploadBatchSubmissions(
  examId: string, 
  files: File[], 
  studentMapping: Record<string, string>
) {
  const { data: { session } } = await supabase.auth.getSession();
  
  const formData = new FormData();
  files.forEach(file => formData.append('files[]', file));
  formData.append('student_mapping', JSON.stringify(studentMapping));
  
  const response = await fetch(`${API_BASE}/api/v1/exams/${examId}/submissions/batch`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${session?.access_token}`,
    },
    body: formData,
  });
  
  return response.json();
}

export async function getGrades(submissionId: string) {
  const headers = await getAuthHeaders();
  const response = await fetch(`${API_BASE}/api/v1/submissions/${submissionId}/grades`, {
    headers,
  });
  return response.json();
}

export async function overrideGrade(
  submissionId: string, 
  questionId: string, 
  newScore: number, 
  reason: string
) {
  const headers = await getAuthHeaders();
  const response = await fetch(
    `${API_BASE}/api/v1/submissions/${submissionId}/questions/${questionId}/override`,
    {
      method: 'POST',
      headers,
      body: JSON.stringify({ score: newScore, reason }),
    }
  );
  return response.json();
}
```

### Authentication Components

```tsx
// components/auth/Login.tsx
'use client'

import { useState } from 'react';
import { supabase } from '@/lib/supabase';
import { useRouter } from 'next/navigation';

export function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    
    const { error } = await supabase.auth.signInWithPassword({
      email,
      password,
    });
    
    if (error) {
      alert(error.message);
    } else {
      router.push('/dashboard');
    }
    setLoading(false);
  };

  return (
    <form onSubmit={handleLogin} className="space-y-4">
      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        className="w-full px-4 py-2 border rounded"
      />
      <input
        type="password"
        placeholder="Password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        className="w-full px-4 py-2 border rounded"
      />
      <button 
        type="submit" 
        disabled={loading}
        className="w-full px-4 py-2 bg-blue-600 text-white rounded"
      >
        {loading ? 'Signing in...' : 'Sign In'}
      </button>
    </form>
  );
}
```

### Batch Upload Component

```tsx
// components/grading/BatchUpload.tsx
'use client'

import { useState, useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import { uploadBatchSubmissions } from '@/lib/api';

interface BatchUploadProps {
  examId: string;
  students: Array<{ id: string; name: string; rollNumber: string }>;
}

export function BatchUpload({ examId, students }: BatchUploadProps) {
  const [files, setFiles] = useState<File[]>([]);
  const [mapping, setMapping] = useState<Record<string, string>>({});
  const [uploading, setUploading] = useState(false);
  const [progress, setProgress] = useState({ total: 0, completed: 0 });

  const onDrop = useCallback((acceptedFiles: File[]) => {
    setFiles(prev => [...prev, ...acceptedFiles]);
    // Auto-map files to students if filenames match roll numbers
    const newMapping: Record<string, string> = {};
    acceptedFiles.forEach(file => {
      const match = students.find(s => 
        file.name.includes(s.rollNumber) || file.name.includes(s.id)
      );
      if (match) {
        newMapping[file.name] = match.id;
      }
    });
    setMapping(prev => ({ ...prev, ...newMapping }));
  }, [students]);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: { 'image/*': [], 'application/pdf': [] }
  });

  const handleUpload = async () => {
    setUploading(true);
    setProgress({ total: files.length, completed: 0 });
    
    try {
      const result = await uploadBatchSubmissions(examId, files, mapping);
      console.log('Batch uploaded:', result);
      // Poll for status updates...
    } catch (error) {
      console.error('Upload failed:', error);
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="space-y-4">
      <div 
        {...getRootProps()} 
        className={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer
          ${isDragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300'}`}
      >
        <input {...getInputProps()} />
        <p>Drag & drop answer sheets here, or click to select files</p>
        <p className="text-sm text-gray-500">Supports PDF and images</p>
      </div>

      {files.length > 0 && (
        <div className="space-y-2">
          <h3 className="font-semibold">Files to upload: {files.length}</h3>
          <div className="max-h-64 overflow-y-auto">
            {files.map((file, idx) => (
              <div key={idx} className="flex justify-between items-center p-2 bg-gray-50 rounded">
                <span>{file.name}</span>
                <select
                  value={mapping[file.name] || ''}
                  onChange={(e) => setMapping(prev => ({ 
                    ...prev, 
                    [file.name]: e.target.value 
                  }))}
                  className="border rounded px-2 py-1"
                >
                  <option value="">Select student</option>
                  {students.map(s => (
                    <option key={s.id} value={s.id}>{s.name} ({s.rollNumber})</option>
                  ))}
                </select>
              </div>
            ))}
          </div>
          <button
            onClick={handleUpload}
            disabled={uploading || Object.keys(mapping).length !== files.length}
            className="w-full px-4 py-2 bg-green-600 text-white rounded disabled:bg-gray-400"
          >
            {uploading ? `Uploading... ${progress.completed}/${progress.total}` : 'Upload All'}
          </button>
        </div>
      )}
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

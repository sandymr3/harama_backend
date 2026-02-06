import { 
  Exam, 
  Question, 
  Rubric, 
  Submission, 
  FinalGrade
} from '@/types';

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
const TENANT_ID = 'd015c777-09b6-45a8-929e-63300557429f'; // Demo Tenant

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || response.statusText);
  }
  return response.json();
}

async function fetchWithTenant(endpoint: string, options: RequestInit = {}): Promise<Response> {
    const headers = {
        'Content-Type': 'application/json',
        'X-Tenant-ID': TENANT_ID,
        ...options.headers,
    } as HeadersInit;

    return fetch(`${BASE_URL}${endpoint}`, {
        ...options,
        headers,
    });
}

export const api = {
  // Exams
  async listExams(): Promise<Exam[]> {
    const response = await fetchWithTenant('/exams');
    return handleResponse<Exam[]>(response);
  },

  async createExam(data: Partial<Exam>): Promise<Exam> {
    const response = await fetchWithTenant('/exams', {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return handleResponse<Exam>(response);
  },

  async getExam(id: string): Promise<Exam> {
    const response = await fetchWithTenant(`/exams/${id}`);
    return handleResponse<Exam>(response);
  },

  async addQuestion(examId: string, data: Partial<Question>): Promise<Question> {
    const response = await fetchWithTenant(`/exams/${examId}/questions`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return handleResponse<Question>(response);
  },

  async setRubric(questionId: string, data: Partial<Rubric>): Promise<Rubric> {
    const response = await fetchWithTenant(`/questions/${questionId}/rubric`, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
    return handleResponse<Rubric>(response);
  },

  // Submissions
  async createSubmission(examId: string, studentId: string, file: File): Promise<Submission> {
    const formData = new FormData();
    formData.append('student_id', studentId);
    formData.append('file', file);

    // Note: Content-Type header is not set for FormData, fetch handles it
    const response = await fetch(`${BASE_URL}/exams/${examId}/submissions`, {
      method: 'POST',
      headers: {
          'X-Tenant-ID': TENANT_ID,
      },
      body: formData,
    });
    return handleResponse<Submission>(response);
  },

  async triggerGrading(submissionId: string): Promise<{ status: string }> {
    const response = await fetchWithTenant(`/submissions/${submissionId}/trigger-grading`, {
      method: 'POST',
    });
    return handleResponse(response);
  },

  // Grading
  async getGrades(submissionId: string): Promise<FinalGrade[]> {
    const response = await fetchWithTenant(`/submissions/${submissionId}/grades`);
    return handleResponse<FinalGrade[]>(response);
  },

  async captureOverride(
    submissionId: string, 
    questionId: string, 
    data: { score: number; reason: string }
  ): Promise<FinalGrade> {
    const response = await fetchWithTenant(`/submissions/${submissionId}/questions/${questionId}/override`, {
      method: 'POST',
      body: JSON.stringify(data),
    });
    return handleResponse<FinalGrade>(response);
  },

  async getStudentFeedback(submissionId: string, questionId: string): Promise<{ feedback: string }> {
    const response = await fetchWithTenant(`/submissions/${submissionId}/questions/${questionId}/feedback`);
    return handleResponse(response);
  },

  async analyzePatterns(questionId: string): Promise<any> {
    const response = await fetchWithTenant(`/questions/${questionId}/analysis`);
    return handleResponse(response);
  },
};
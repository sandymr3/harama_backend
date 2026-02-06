import { 
  Exam, 
  Question, 
  Rubric, 
  Submission, 
  FinalGrade, 
  MultiEvalResult 
} from '@/types';

const BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

async function handleResponse<T>(response: Response): Promise<T> {
  if (!response.ok) {
    const error = await response.text();
    throw new Error(error || response.statusText);
  }
  return response.json();
}

export const api = {
  // Exams
  async createExam(data: Partial<Exam>): Promise<Exam> {
    const response = await fetch(`${BASE_URL}/exams`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return handleResponse<Exam>(response);
  },

  async getExam(id: string): Promise<Exam> {
    const response = await fetch(`${BASE_URL}/exams/${id}`);
    return handleResponse<Exam>(response);
  },

  async addQuestion(examId: string, data: Partial<Question>): Promise<Question> {
    const response = await fetch(`${BASE_URL}/exams/${examId}/questions`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return handleResponse<Question>(response);
  },

  async setRubric(questionId: string, data: Partial<Rubric>): Promise<Rubric> {
    const response = await fetch(`${BASE_URL}/questions/${questionId}/rubric`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return handleResponse<Rubric>(response);
  },

  // Submissions
  async createSubmission(examId: string, studentId: string, file: File): Promise<Submission> {
    const formData = new FormData();
    formData.append('student_id', studentId);
    formData.append('file', file);

    const response = await fetch(`${BASE_URL}/exams/${examId}/submissions`, {
      method: 'POST',
      body: formData,
    });
    return handleResponse<Submission>(response);
  },

  async triggerGrading(submissionId: string): Promise<{ status: string }> {
    const response = await fetch(`${BASE_URL}/submissions/${submissionId}/trigger-grading`, {
      method: 'POST',
    });
    return handleResponse(response);
  },

  // Grading
  async getGrades(submissionId: string): Promise<FinalGrade[]> {
    const response = await fetch(`${BASE_URL}/submissions/${submissionId}/grades`);
    return handleResponse<FinalGrade[]>(response);
  },

  async captureOverride(
    submissionId: string, 
    questionId: string, 
    data: { score: number; reason: string }
  ): Promise<FinalGrade> {
    const response = await fetch(`${BASE_URL}/submissions/${submissionId}/questions/${questionId}/override`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data),
    });
    return handleResponse<FinalGrade>(response);
  },

  async getStudentFeedback(submissionId: string, questionId: string): Promise<{ feedback: string }> {
    const response = await fetch(`${BASE_URL}/submissions/${submissionId}/questions/${questionId}/feedback`);
    return handleResponse(response);
  },

  async analyzePatterns(questionId: string): Promise<any> {
    const response = await fetch(`${BASE_URL}/questions/${questionId}/analysis`);
    return handleResponse(response);
  },
};

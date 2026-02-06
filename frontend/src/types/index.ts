export type ProcessingStatus = 'pending' | 'processing' | 'completed' | 'failed';
export type AnswerType = 'short_answer' | 'essay' | 'mcq' | 'diagram';
export type GradeStatus = 'pending' | 'auto_graded' | 'needs_review' | 'overridden' | 'final';

export interface BoundingBox {
  x: number;
  y: number;
  width: number;
  height: number;
}

export interface OCRResult {
  page_number: number;
  raw_text: string;
  confidence: number;
  image_url: string;
  bounding_boxes: BoundingBox[];
  corrected_text?: string;
}

export interface AnswerSegment {
  id: string;
  submission_id: string;
  question_id: string;
  text: string;
  page_indices: number[];
  bounding_box: BoundingBox[];
  diagrams: string[];
}

export interface Criterion {
  id: string;
  description: string;
  points: number;
  required: boolean;
  category: string;
}

export interface PartialCreditRule {
  id: string;
  condition: string;
  points: number;
  description: string;
  dependencies?: string[];
}

export interface CommonMistake {
  id: string;
  description: string;
  penalty: number;
  category: string;
  frequency: number;
}

export interface Rubric {
  id: string;
  question_id: string;
  full_credit_criteria: Criterion[];
  partial_credit_rules: PartialCreditRule[];
  common_mistakes: CommonMistake[];
  key_concepts?: string[];
  grading_notes?: string;
  strict_mode: boolean;
}

export interface Question {
  id: string;
  exam_id: string;
  question_text: string;
  points: number;
  answer_type: AnswerType;
  rubric?: Rubric;
  visual_aids?: string[];
}

export interface Exam {
  id: string;
  title: string;
  subject: string;
  questions?: Question[];
  created_at: string;
  tenant_id: string;
}

export interface Submission {
  id: string;
  exam_id: string;
  student_id: string;
  uploaded_at: string;
  processing_status: ProcessingStatus;
  ocr_results?: OCRResult[];
  answers?: AnswerSegment[];
  tenant_id: string;
}

export interface GradingResult {
  id: string;
  submission_id: string;
  question_id: string;
  score: number;
  max_score: number;
  confidence: number;
  reasoning: string;
  criteria_met: string[];
  mistakes_found: string[];
  ai_evaluator_id: string;
  created_at: string;
}

export interface MultiEvalResult {
  evaluations: GradingResult[];
  variance: number;
  mean_score: number;
  consensus_score: number;
  confidence: number;
  should_escalate: boolean;
  reasoning: string;
}

export interface FinalGrade {
  id: string;
  submission_id: string;
  question_id: string;
  final_score: number;
  max_score: number;
  ai_score?: number;
  override_score?: number;
  confidence: number;
  status: GradeStatus;
  reasoning: string;
  updated_at: string;
}

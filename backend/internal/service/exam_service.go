package service

import (
	"context"
	"harama/internal/domain"
	"harama/internal/repository/postgres"

	"github.com/google/uuid"
)

type ExamService struct {
	repo      *postgres.ExamRepo
	auditRepo *postgres.AuditRepo
}

func NewExamService(repo *postgres.ExamRepo, auditRepo *postgres.AuditRepo) *ExamService {
	return &ExamService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (s *ExamService) CreateExam(ctx context.Context, exam *domain.Exam) error {
	err := s.repo.Create(ctx, exam)
	if err == nil {
		_ = s.auditRepo.Save(ctx, &domain.AuditLog{
			EntityType: "exam",
			EntityID:   exam.ID,
			EventType:  "created",
			Changes: map[string]interface{}{
				"title":   exam.Title,
				"subject": exam.Subject,
			},
		})
	}
	return err
}

func (s *ExamService) GetExam(ctx context.Context, id uuid.UUID) (*domain.Exam, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ExamService) AddQuestion(ctx context.Context, examID uuid.UUID, question *domain.Question) error {
	question.ExamID = examID
	err := s.repo.CreateQuestion(ctx, question)
	if err == nil {
		_ = s.auditRepo.Save(ctx, &domain.AuditLog{
			EntityType: "question",
			EntityID:   question.ID,
			EventType:  "created",
			Changes: map[string]interface{}{
				"exam_id": examID,
				"text":    question.QuestionText,
				"points":  question.Points,
			},
		})
	}
	return err
}

func (s *ExamService) SetRubric(ctx context.Context, questionID uuid.UUID, rubric *domain.Rubric) error {
	rubric.QuestionID = questionID
	err := s.repo.UpdateRubric(ctx, rubric)
	if err == nil {
		_ = s.auditRepo.Save(ctx, &domain.AuditLog{
			EntityType: "rubric",
			EntityID:   rubric.ID,
			EventType:  "updated",
			Changes: map[string]interface{}{
				"question_id": questionID,
			},
		})
	}
	return err
}

func (s *ExamService) ListExams(ctx context.Context, tenantID uuid.UUID) ([]domain.Exam, error) {
	return s.repo.ListByTenant(ctx, tenantID)
}

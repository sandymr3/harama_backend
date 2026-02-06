package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"harama/internal/repository/postgres"

	"github.com/google/uuid"
)

type AnalyticsService struct {
	gradeRepo *postgres.GradeRepo
	examRepo  *postgres.ExamRepo
	subRepo   *postgres.SubmissionRepo
}

func NewAnalyticsService(gradeRepo *postgres.GradeRepo, examRepo *postgres.ExamRepo, subRepo *postgres.SubmissionRepo) *AnalyticsService {
	return &AnalyticsService{
		gradeRepo: gradeRepo,
		examRepo:  examRepo,
		subRepo:   subRepo,
	}
}

func (s *AnalyticsService) GetGradingTrends(ctx context.Context, examID uuid.UUID) (interface{}, error) {
	// 1. Get raw stats from DB
	stats, err := s.gradeRepo.GetExamStats(ctx, examID)
	if err != nil {
		return nil, err
	}

	// 2. Enrich with Question text (optional, but good for UI)
	type EnrichedStat struct {
		postgres.QuestionStat
		QuestionText string `json:"question_text"`
	}

	enriched := make([]EnrichedStat, len(stats))
	for i, stat := range stats {
		q, err := s.examRepo.GetQuestionByID(ctx, stat.QuestionID)
		text := "Unknown Question"
		if err == nil {
			text = q.QuestionText
		}
		
		enriched[i] = EnrichedStat{
			QuestionStat: stat,
			QuestionText: text,
		}
	}

	return enriched, nil
}

func (s *AnalyticsService) ExportGrades(ctx context.Context, examID uuid.UUID, format string) ([]byte, string, error) {
	if format != "csv" {
		return nil, "", fmt.Errorf("unsupported format: %s", format)
	}

	// 1. Get Exam (for Questions)
	exam, err := s.examRepo.GetByID(ctx, examID)
	if err != nil {
		return nil, "", err
	}

	// 2. Get Submissions
	submissions, err := s.subRepo.ListByExam(ctx, examID)
	if err != nil {
		return nil, "", err
	}

	// 3. Prepare CSV
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Header
	header := []string{"Student ID", "Submission ID", "Total Score"}
	for _, q := range exam.Questions {
		header = append(header, fmt.Sprintf("Q: %s (Max: %d)", q.QuestionText, q.Points))
	}
	writer.Write(header)

	// Rows
	for _, sub := range submissions {
		grades, err := s.gradeRepo.GetBySubmission(ctx, sub.ID)
		if err != nil {
			continue // Skip or log error?
		}

		gradeMap := make(map[uuid.UUID]float64)
		totalScore := 0.0
		for _, g := range grades {
			gradeMap[g.QuestionID] = g.FinalScore
			totalScore += g.FinalScore
		}

		row := []string{sub.StudentID, sub.ID.String(), fmt.Sprintf("%.2f", totalScore)}
		for _, q := range exam.Questions {
			score, ok := gradeMap[q.ID]
			if ok {
				row = append(row, fmt.Sprintf("%.2f", score))
			} else {
				row = append(row, "N/A")
			}
		}
		writer.Write(row)
	}

	writer.Flush()
	return buf.Bytes(), "text/csv", nil
}

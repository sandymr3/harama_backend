package unit_test

import (
	"context"
	"testing"
	"time"

	"harama/internal/domain"
	"harama/internal/repository/postgres"
	"harama/internal/service"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func TestExamService_CreateExam(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	bunDB := bun.NewDB(db, pgdialect.New())
	examRepo := postgres.NewExamRepo(bunDB)
	auditRepo := postgres.NewAuditRepo(bunDB)
	examService := service.NewExamService(examRepo, auditRepo)

	ctx := context.Background()
	exam := &domain.Exam{
		ID:       uuid.New(),
		Title:    "Final Exam",
		Subject:  "Science",
		TenantID: uuid.New(),
	}

	// 1. Expectation for ExamRepo.Create (Insert into exams)
	mock.ExpectQuery(`INSERT INTO "exams" .*`).
		WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(time.Now()))

	// 2. Expectation for AuditRepo.Save (GetLastHash then Insert into audit_log)
	mock.ExpectQuery(`SELECT .* FROM "audit_log" .*`).
		WillReturnRows(sqlmock.NewRows([]string{"hash"}))
	
	mock.ExpectQuery(`INSERT INTO "audit_log" .*`).
		WillReturnRows(sqlmock.NewRows([]string{"actor_id"}).AddRow(nil))

	// Execute
	err = examService.CreateExam(ctx, exam)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestExamService_ListExams(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	bunDB := bun.NewDB(db, pgdialect.New())
	examRepo := postgres.NewExamRepo(bunDB)
	// AuditRepo not needed for List
	examService := service.NewExamService(examRepo, nil)

	ctx := context.Background()
	tenantID := uuid.New()

	// Expectation: Select exams
	mock.ExpectQuery(`SELECT .* FROM "exams" .* WHERE \(tenant_id = .*\)`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).
			AddRow(uuid.New(), "Exam 1").
			AddRow(uuid.New(), "Exam 2"))

	// Execute
	exams, err := examService.ListExams(ctx, tenantID)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, exams, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}

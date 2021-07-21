package repository

import (
	"database/sql"
	"github.com/ag-students/support-service/internal/microservices/communication/models"
	"github.com/ag-students/support-service/internal/microservices/communication/repository/postgres"
)

type CommunicationHistoryRepository interface {
	CreateCommunication(comm models.Communication) (int, error)
}

type Repository struct {
	CommunicationHistoryRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		CommunicationHistoryRepository: postgres.NewCommunicationPSQL(db),
	}
}

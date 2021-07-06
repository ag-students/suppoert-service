package repository

import (
	"database/sql"
	"github.com/ag-students/support-service/internal/microservices/communication/models"
	"github.com/ag-students/support-service/internal/microservices/communication/repository/postgres"
)

type CommunicationRepository interface {
	GetCommunication(id int) (models.Communication, error)
	CreateCommunication(comm models.Communication) (int, error)
}

type Repository struct {
	CommunicationRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		CommunicationRepository: postgres.NewCommunicationPSQL(db),
	}
}

package repository

import (
	"database/sql"
	"github.com/ag-students/support-service/internal/microservices/communication/models"
	"github.com/ag-students/support-service/internal/microservices/communication/repository/miniorepo"
	"github.com/ag-students/support-service/internal/microservices/communication/repository/postgres"
)

type CommunicationHistoryRepository interface {
	CreateCommunication(comm models.Communication) (int, error)
}

type DocumentsRepository interface {
	GetDocument(docHash string) (string, error)
}

type Repository struct {
	CommunicationHistoryRepository
	DocumentsRepository
}

func NewRepository(db *sql.DB, cnf *miniorepo.FileServerConfig) *Repository {
	return &Repository{
		CommunicationHistoryRepository: postgres.NewCommunicationPSQL(db),
		DocumentsRepository:            miniorepo.NewDocumentsMinio(cnf),
	}
}

package postgres

import (
	"database/sql"
	"github.com/ag-students/support-service/internal/microservices/communication/models"
)

type CommunicationPSQL struct {
	conn *sql.DB
}

func NewCommunicationPSQL(conn *sql.DB) *CommunicationPSQL {
	return &CommunicationPSQL{conn: conn}
}

func (r *CommunicationPSQL) CreateCommunication(comm models.Communication) (int, error) {
	return 0, nil
}

func (r *CommunicationPSQL) GetCommunication(id int) (models.Communication, error) {
	return models.Communication{
		Id:                0,
		UserId:            0,
		CommunicationType: "sd",
	}, nil
}
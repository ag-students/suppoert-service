package postgres

import (
	"database/sql"
	"fmt"
	"github.com/ag-students/support-service/internal/microservices/communication/models"
)

const communicationHistoryTableName string = "MONITORING.COMMUNICATION_HISTORY"

type CommunicationPSQL struct {
	conn *sql.DB
}

func NewCommunicationPSQL(conn *sql.DB) *CommunicationPSQL {
	return &CommunicationPSQL{conn: conn}
}

func (r *CommunicationPSQL) CreateCommunication(comm models.Communication) (int, error) {
	tx, err := r.conn.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (communication_type, delayed, phone_number, email) "+
		"values ($1, $2, $3, $4) RETURNING id", communicationHistoryTableName)

	row := tx.QueryRow(createItemQuery, comm.CommunicationType, comm.Delayed, comm.Phone, comm.Email)
	err = row.Scan(&itemId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

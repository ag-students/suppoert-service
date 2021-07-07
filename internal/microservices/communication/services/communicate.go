package services

import (
	"github.com/ag-students/support-service/internal/microservices/communication/repository"
	"log"
)

type PatientNotifier struct {
	repo *repository.Repository
}

func (r *PatientNotifier) Notify() error {
	log.Print(r.repo)

	return nil
}



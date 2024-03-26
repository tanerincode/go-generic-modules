package service

import (
	responsehandler "github.com/tanerincode/go-generic-modules/pkg/responseHandler"
	"github.com/tanerincode/go-generic-modules/pkg/storage"
	"net/http"
)

type HealthCheckService interface {
	Check(w http.ResponseWriter)
}

type healthCheckService struct {
	db storage.Storage
}

func (s *healthCheckService) Check(w http.ResponseWriter) {
	err := s.db.Disconnect()
	if err != nil {
		response := &responsehandler.Response{
			Status:  false,
			Code:    500,
			Message: "Database connection failed",
		}
		responsehandler.RespondAsJSON(w, response, err)
		return
	}

	response := &responsehandler.Response{
		Status:  true,
		Code:    200,
		Message: "Database connection successful",
	}
	responsehandler.RespondAsJSON(w, response, nil)
}

func NewHealthCheckService(db storage.Storage) HealthCheckService {
	return &healthCheckService{db: db}
}

package responsehandler

import (
	"encoding/json"
	"net/http"
)

type BadRequestError struct {
	msg string
}

func (e *BadRequestError) Error() string { return e.msg }

type InternalServerError struct {
	msg string
}

func (e *InternalServerError) Error() string { return e.msg }

type Response struct {
	Status  bool        `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespondAsJSON(w http.ResponseWriter, payload interface{}, err error) error {
	var res Response

	switch e := err.(type) {
	case *BadRequestError:
		res.Status = false
		res.Code = http.StatusBadRequest
		res.Message = e.Error()
	case *InternalServerError:
		res.Status = false
		res.Code = http.StatusInternalServerError
		res.Message = e.Error()
	case nil:
		res.Status = true
		res.Code = http.StatusOK
		res.Message = "Success"
	default:
		res.Status = false
		res.Code = http.StatusInternalServerError
		res.Message = "Unknown error"
	}

	res.Data = payload

	response, err := json.Marshal(res)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	_, err = w.Write(response)
	return err
}

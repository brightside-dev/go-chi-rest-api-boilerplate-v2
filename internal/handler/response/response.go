package response

import (
	"encoding/json"
	"net/http"
)

type APIResponseDTO struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type APIResponseManager interface {
	ErrorResponse(w http.ResponseWriter, r *http.Request, err error, statusCode int)
	SuccessResponse(w http.ResponseWriter, r *http.Request, response interface{}, statusCode ...int)
	ClientErrorResponse(w http.ResponseWriter, r *http.Request, err error)
}

type apiResponseManager struct {
}

func NewAPIResponseManager() APIResponseManager {
	return &apiResponseManager{}
}

func (rm *apiResponseManager) ErrorResponse(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonResp, err := json.Marshal(APIResponseDTO{
		Success: false,
		Data:    nil,
		Error:   err.Error(),
	})
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func (rm *apiResponseManager) SuccessResponse(w http.ResponseWriter, r *http.Request, response interface{}, statusCode ...int) {
	w.Header().Set("Content-Type", "application/json")

	if len(statusCode) == 0 {
		statusCode = append(statusCode, http.StatusOK)
	}

	// Set the response status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode[0])

	jsonResp, err := json.Marshal(APIResponseDTO{
		Success: true,
		Data:    response,
		Error:   "",
	})
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

func (rm *apiResponseManager) ClientErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	rm.ErrorResponse(w, r, err, http.StatusBadRequest)
}

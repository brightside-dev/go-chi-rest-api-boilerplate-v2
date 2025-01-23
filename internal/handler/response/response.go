package response

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	AddSource: true,
	// }))

	//logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	jsonResp, err := json.Marshal(APIResponse{
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

func SuccessResponse(w http.ResponseWriter, r *http.Request, response interface{}, statusCode ...int) {
	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	AddSource: true,
	// }))

	//logger.Info("success", "method", r.Method, "uri", r.URL.RequestURI())

	w.Header().Set("Content-Type", "application/json")

	if len(statusCode) == 0 {
		statusCode = append(statusCode, http.StatusOK)
	}

	// Set the response status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode[0])

	jsonResp, err := json.Marshal(APIResponse{
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

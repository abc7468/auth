package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func StringToInt(inputVal string) int {
	res, err := strconv.Atoi(inputVal)
	if err != nil {
		fmt.Print(err)
	}
	return res
}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
func RespondError(w http.ResponseWriter, code int, message string) {
	RespondJSON(w, code, map[string]string{"error": message})
}

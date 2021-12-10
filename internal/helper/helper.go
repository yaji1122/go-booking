package helper

import (
	"fmt"
	"github.com/yaji1122/bookings-go/internal/logger"
	"net/http"
	"runtime/debug"
)

var log *logger.Logger

// NewHelper sets up app config for helper
func NewHelper(logger *logger.Logger) {
	log = logger
}

func ClientError(w http.ResponseWriter, status int) {
	log.InfoLogger.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s \n %s", err.Error(), debug.Stack())
	log.ErrorLogger.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

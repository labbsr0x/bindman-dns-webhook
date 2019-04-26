package hook

import (
	"encoding/json"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"
	"github.com/sirupsen/logrus"
	"net/http"
)

// write200Response writes the response to be sent
func writeJSONResponse(payload interface{}, statusCode int, w http.ResponseWriter) {
	// Headers must be set before call WriteHeader or Write. see https://golang.org/pkg/net/http/#ResponseWriter
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if payload != nil {
		types.PanicIfError(json.NewEncoder(w).Encode(payload))
	}

	logrus.Infof("%d Response sent. Payload: %#v", statusCode, payload)
}

// handleError recovers from a panic
func handleError(w http.ResponseWriter) {
	r := recover()
	if r != nil {
		err := types.InternalServerError("An internal server error occurred, please contact the system administrator.", nil)
		if e, ok := r.(*types.Error); ok {
			err = e
		}
		logrus.Error(err)
		writeJSONResponse(err, err.Code, w)
	}
}

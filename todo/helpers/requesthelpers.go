package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/phalpin/liberr"
	"github.com/phalpin/liberr/errortype"
	"net/http"
)

func GetRouteVariable(r *http.Request, name string) (string, error) {
	vars := mux.Vars(r)
	if val, ok := vars[name]; ok {
		return val, nil
	}

	return "", liberr.NewKnown(fmt.Sprintf("route variable named '%v' not found", name), "Attempted to retrieve a route variable that did not exist.", liberr.WithErrorType(errortype.InvalidArgument))
}

type errorReturnVal struct {
	Message      string `json:"Message"`
	DebugMessage string `json:"DebugMessage"`
	StackTrace   string `json:"StackTrace"`
}

func getResponseValues(err error) (*errorReturnVal, int) {
	statusCode := http.StatusInternalServerError
	retVal := &errorReturnVal{
		Message:      "",
		DebugMessage: err.Error(),
	}

	knownCast, knownCastOk := err.(*liberr.KnownError)
	if knownCastOk {
		retVal.Message = knownCast.FriendlyMessage
	}

	baseCast, baseCastOk := err.(*liberr.BaseError)
	if baseCastOk {
		retVal.StackTrace = baseCast.StackTrace
		statusCode = baseCast.ErrorType.ToHttpStatusCode()
	}

	return retVal, statusCode
}

func WriteErrorResponse(w http.ResponseWriter, err error) {
	respVal, statusCode := getResponseValues(err)

	res, _ := json.Marshal(respVal)

	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	_, _ = w.Write(res)
}

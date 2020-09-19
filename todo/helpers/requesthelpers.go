package helpers

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouteVariable(r *http.Request, name string) (string, error) {
	vars := mux.Vars(r)
	if val, ok := vars[name]; ok {
		return val, nil
	}

	return "", errors.New(fmt.Sprintf("route variable named '%v' not found", name))
}

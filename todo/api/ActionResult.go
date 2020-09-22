package api

import (
	"encoding/json"
	"github.com/phalpin/GoGoGadgetRESTAPI/todo/helpers"
	"github.com/phalpin/liberr"
	"github.com/phalpin/liberr/errortype"
	"net/http"
)

//#region ActionResult stuff
type ActionResult struct {
	Result           interface{}
	ErrorEncountered error
}

func ObjResult(retVal interface{}) (*ActionResult, error) {
	return &ActionResult{
		Result: retVal,
	}, nil
}

func ErrResult(errorEnc error) (*ActionResult, error) {
	return nil, errorEnc
}

func NilResult() (*ActionResult, error) {
	return &ActionResult{
		Result: nil,
	}, nil
}

//#endregion

//#region HandlerPackage Stuff
type HandlerPackage struct {
	Path    string
	Name    string
	Handler func(*http.Request) (*ActionResult, error)
	Methods []string
}

func HandlerRecord(name string, path string, handlerFunc func(*http.Request) (*ActionResult, error), methods ...string) *HandlerPackage {
	pkg := &HandlerPackage{
		Path:    path,
		Name:    name,
		Handler: handlerFunc,
		Methods: methods,
	}

	return pkg
}

func (hp *HandlerPackage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	res, err := hp.Handler(r)
	if err != nil {
		helpers.WriteErrorResponse(w, err)
	}

	if res != nil {
		if res.Result != nil {
			encObj, marshalErr := json.Marshal(res.Result)
			if marshalErr != nil {
				helpers.WriteErrorResponse(w, liberr.NewBase("fatal error occurred, please try again later", liberr.WithErrorType(errortype.Unknown)))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write(encObj)
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

//#endregion

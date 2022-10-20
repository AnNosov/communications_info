package v1

import (
	"net/http"

	"github.com/AnNosov/communications_info/internal/usecase"

	"encoding/json"

	"github.com/gorilla/mux"
)

type CommunicationRouter struct {
	cAction usecase.CommunicationAction
}

func setErrorResponse(w http.ResponseWriter, status int, response string) {
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func NewCommunicationActionRoutes(handler *mux.Router, cAction usecase.CommunicationAction) {
	c := &CommunicationRouter{cAction}
	handler.HandleFunc("/", c.HandleConnection).Methods("GET")
}

func (c *CommunicationRouter) HandleConnection(w http.ResponseWriter, r *http.Request) {
	data := c.cAction.GetCommunicationResult()
	answer, err := json.Marshal(data)
	if err != nil {
		setErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(answer)
}

package handler

import (
	"L0test/pkg/repository"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type RegOrder struct {
	Order repository.Repository
}

type Id struct {
	idOrder string `json:"idOrder"`
}

func (o *RegOrder) Response(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")
	order := o.Order.FindById(id)
	jsonOrder, err := json.Marshal(order)
	if err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(jsonOrder)
}

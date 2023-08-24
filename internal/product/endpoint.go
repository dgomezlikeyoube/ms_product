package product

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
 	"github.com/dgomezlikeyoube/ms_meta/meta"
	"github.com/gorilla/mux"

)

type (
	Controller func(w http.ResponseWriter, r *http.Request)

	EndPoints struct {
		Create Controller
		Get    Controller
		GetAll Controller
		Update Controller
		Delete Controller
	}

	CreateRequest struct {
		Name        string  `json:"name"`
		Sku         string  `json:"sku"`
		Quantity    int32   `json:"quantity"`
		Price       float32 `json:"price"`
		CostPrice   float32 `json:"costprice"`
		Weight      int32   `json:"wight"`
		Enabled     bool    `json:"enabled"`
		Descripcion string  `json:"descripcion"`
		Category    string  `json:"category"`
	}

	UpdateRequest struct {
		Name        *string  `json:"name"`
		Sku         *string  `json:"sku"`
		Quantity    *int32   `json:"quantity"`
		Price       *float32 `json:"price"`
		CostPrice   *float32 `json:"costprice"`
		Weight      *int32   `json:"wight"`
		Enabled     *bool    `json:"enabled"`
		Descripcion *string  `json:"descripcion"`
		Category    *string  `json:"category"`
	}

	Response struct {
		Status int         `json:"status"`
		Data   interface{} `json:"data,omitempty"`
		Err    string      `json:"error,omitempty"`
		Meta   *meta.Meta  `json:"meta,omitempty"`
	}

	Config struct {
		LimPageDef string
	}
)

func MakeEndpoints(s Service, config Config) EndPoints {
	return EndPoints{
		Create: makeCreateEndpoint(s),
		Get:    makeGetEndpoint(s),
		GetAll: makeGetAllEndpoint(s, config),
		Update: makeUpdateEndpoint(s),
		Delete: makeDeleteEndpoint(s),
	}

}

func makeGetAllEndpoint(s Service, config Config) Controller {
	return func(w http.ResponseWriter, r *http.Request) {

		v := r.URL.Query()

		filters := Filters{
			Sku: v.Get("sku"),
		}

		limit, _ := strconv.Atoi(v.Get("limit"))
		page, _ := strconv.Atoi(v.Get("page"))

		count, err := s.Count(filters)

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}
		meta, err := meta.New(page, limit, count, config.LimPageDef)

		products, err := s.GetAll(filters, meta.Offset(), meta.Limit())
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		json.NewEncoder(w).Encode(&Response{Status: 200, Data: products, Meta: meta})
	}
}

func makeGetEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]
		product, err := s.Get(id)

		if err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: product})
	}
}

func makeDeleteEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		path := mux.Vars(r)
		id := path["id"]

		if err := s.Delete(id); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "success"})
	}
}

func makeUpdateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		path := mux.Vars(r)
		id := path["id"]

		if err := s.Update(id, req.Name, req.Sku, req.Quantity, req.Price, req.CostPrice, req.Weight, req.Enabled, req.Descripcion, req.Category); err != nil {
			w.WriteHeader(404)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "id producto no existe"})
			return
		}
		json.NewEncoder(w).Encode(&Response{Status: 200, Data: "success"})
	}
}

func makeCreateEndpoint(s Service) Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateRequest

		fmt.Println("Create products")

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
			return
		}

		if req.Sku == "0" || req.Sku == "" {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: "sku no encontrado"})
			return
		}
		product, err := s.Create(req.Name, req.Sku, req.Quantity, req.Price, req.CostPrice, req.Weight, req.Enabled, req.Descripcion, req.Category)
		if err != nil {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(&Response{Status: 400, Err: err.Error()})
		}
		json.NewEncoder(w).Encode(Response{Status: 200, Data: product})

	}
}

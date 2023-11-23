package httphandler2

import (
	"log"
	"middleware/internal/controllers"
	"middleware/internal/httphandler"
	"net/http"

	"github.com/go-chi/chi"
)

func NewCustomerHttpHandler2(prop httphandler.HTTPHandlerProperty) http.Handler {
	log.Println("TEST")
	r := chi.NewRouter()
	r.Post("/", controllers.CreateCustomer)

	return r
}

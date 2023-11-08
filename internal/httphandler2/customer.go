package httphandler2

import (
	"github.com/go-chi/chi"
	"gitlab.com/mohamadikbal/project-privy/internal/controllers"
	"gitlab.com/mohamadikbal/project-privy/internal/httphandler"
	"log"
	"net/http"
)

func NewCustomerHttpHandler2(prop httphandler.HTTPHandlerProperty) http.Handler {
	log.Println("TEST")
	r := chi.NewRouter()
	r.Post("/", controllers.CreateCustomer)

	return r
}

package cashbook

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"

	"github.com/gossie/router"
)

var templates *template.Template

func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	err := templates.ExecuteTemplate(w, fmt.Sprintf("%s.html", tmpl), data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) routes() {
	currencyFunc := func(f float64) string {
		return fmt.Sprintf("%.2f €", math.Round(f*100.0)/100.0)
	}

	pathFunc := func(path string) string {
		return s.basePath + path
	}

	templates = template.Must(template.New("").Funcs(template.FuncMap{
		"currency": currencyFunc,
		"path":     pathFunc,
	}).ParseFS(s.htmlTemplates, "templates/scaffold.html", "templates/index.html", "templates/costs.html", "templates/participants.html", "templates/checkout.html"))

	s.router.Handle("/assets/:dir/:file", http.FileServer(http.FS(s.assets)))
	s.router.Get("/", s.indexHandler())
	s.router.Post("/cashbooks", s.cashbooksHandler())
	s.router.Get("/cashbooks/:cashbookId/participants", s.getParticipantsHandler())
	s.router.Post("/cashbooks/:cashbookId/participants", s.postParticipantsHandler())
	s.router.Get("/cashbooks/:cashbookId/costs", s.singleCashbookHandler())
	s.router.Post("/cashbooks/:cashbookId/payments", s.paymentsHandler())
	s.router.Get("/cashbooks/:cashbookId/checkout", s.checkoutHandler())

	s.router.Get("/health", s.healthHandler())
}

func (s *Server) indexHandler() router.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request, _ router.Context) {
		renderTemplate(w, "index", nil)
	}
}

func (s *Server) cashbooksHandler() router.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request, _ router.Context) {
		tripName := r.FormValue("tripName")
		theCashbook, err := s.createNewCashbook(r.Context(), tripName)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		renderTemplate(w, "participants", theCashbook)
	}
}

func (s *Server) singleCashbookHandler() router.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		theCashbook, err := s.findById(r.Context(), ctx.PathParameter("cashbookId"))
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		renderTemplate(w, "costs", theCashbook)
	}
}

func (s *Server) getParticipantsHandler() router.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		theCashbook, err := s.findById(r.Context(), ctx.PathParameter("cashbookId"))
		if err != nil {
			http.Error(w, "", http.StatusNotFound)
			return
		}
		renderTemplate(w, "participants", theCashbook)
	}
}

func (s *Server) postParticipantsHandler() router.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		participantName := r.FormValue("participantName")
		theCashbook, err := s.createNewParticipant(r.Context(), ctx.PathParameter("cashbookId"), participantName)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		renderTemplate(w, "participants", theCashbook)
	}
}

func (s *Server) paymentsHandler() router.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
		p := Payment{
			Payer:       r.FormValue("payer"),
			Amount:      amount,
			Description: r.FormValue("description"),
		}
		theCashbook, err := s.createNewPayment(r.Context(), ctx.PathParameter("cashbookId"), &p)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		renderTemplate(w, "costs", theCashbook)
	}
}

func (s *Server) checkoutHandler() router.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		theCashbook, err := s.findById(r.Context(), ctx.PathParameter("cashbookId"))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		renderTemplate(w, "checkout", theCashbook.Checkout())
	}
}

func (s *Server) healthHandler() router.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		w.WriteHeader(200)
	}
}

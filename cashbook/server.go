package cashbook

import (
	"embed"
	"log"
	"net/http"
	"time"

	"github.com/gossie/router"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	db                    *mongo.Client
	router                *router.HttpRouter
	assets, htmlTemplates embed.FS
	basePath              string
}

func NewServer(mongoClient *mongo.Client, assets, htmlTemplates embed.FS, basePath string) *Server {
	s := Server{
		db:            mongoClient,
		router:        router.New(),
		assets:        assets,
		htmlTemplates: htmlTemplates,
		basePath:      basePath,
	}

	s.router.Use(profile)
	s.routes()

	return &s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func profile(handler router.HttpHandler) router.HttpHandler {
	return func(w http.ResponseWriter, r *http.Request, ctx router.Context) {
		start := time.Now()
		defer func() {
			log.Default().Println("request took", time.Since(start).Milliseconds(), "ms")
		}()

		handler(w, r, ctx)
	}
}

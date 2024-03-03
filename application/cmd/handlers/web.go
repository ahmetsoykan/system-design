package handlers

import (
	"log"
	"url-shortener/internal/cache"
	"url-shortener/internal/db"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// App is the entrypoint into our application and what controls the context of
// each request. Feel free to add any configuration data/logic on this type.
type Server struct {
	Router *chi.Mux
	Cache  *cache.RedisCache
	DB     *db.DynamoDB
}

// NewServer constructs a Server to handle a set of routes.
func NewServer(s Config) *Server {

	srv := &Server{
		Router: chi.NewRouter(),
		Cache:  cache.NewRedisCache(s.RedisHost),
		DB:     db.NewDynamoDBClient(s.Region),
	}
	srv.DB.DDBInit()

	err := srv.Cache.Set("test", "ok")
	if err != nil {
		log.Fatalf("main : no redis connection")
	} else {
		log.Print("main : connected to redis")
	}

	srv.routes()

	return srv
}

func (s *Server) routes() {

	s.Router.Use(render.SetContentType(render.ContentTypeJSON))
	s.Router.Get("/", s.health)
	s.Router.Get("/health", s.health)
	s.Router.Post("/shorten", s.shorten)
	s.Router.Get("/{shorturl}", s.redirect)
}

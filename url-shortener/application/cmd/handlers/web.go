package handlers

import (
	"log"
	"url-shortener/internal/cache"
	"url-shortener/internal/db"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

// App is the entrypoint into our application and what controls the context of
// each request. Feel free to add any configuration data/logic on this type.
type Server struct {
	Router *mux.Router
	Cache  *cache.RedisCache
	DB     *db.DynamoDB
}

// NewServer constructs a Server to handle a set of routes.
func NewServer(s Config) *Server {

	srv := &Server{
		Router: mux.NewRouter(),
		Cache:  cache.NewRedisCache(s.RedisHost),
		DB:     db.NewDynamoDBClient(s.Region),
	}
	// database init
	srv.DB.DDBInit()
	// cache access test
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

	s.Router.Use(otelmux.Middleware("url-shortener"))

	s.Router.HandleFunc("/", s.health).Methods("GET")
	s.Router.HandleFunc("/health", s.health).Methods("GET")
	s.Router.HandleFunc("/shorten", s.shorten).Methods("POST")
	s.Router.HandleFunc("/{short}", s.redirect).Methods("GET")

}

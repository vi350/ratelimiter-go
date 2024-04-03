package app

import (
	"floodcontrol/internal/floodcontrol"
	"floodcontrol/internal/storage"
	"floodcontrol/internal/storage/redis"
	redisLibrary "github.com/redis/go-redis/v9"
	"log"
	"net/http"
)

func Run(cfg *Config) {
	rd := redisLibrary.NewClient(&redisLibrary.Options{
		Addr:     cfg.Cache.Host + ":" + cfg.Cache.Port,
		Password: cfg.Cache.Pass,
	})
	var s storage.Storage
	s = redis.New(rd)
	var fc floodcontrol.FloodControl // интерфейс данный по заданию
	fc = floodcontrol.NewRateLimiter(s, cfg.Limit.Count, cfg.Limit.Period)

	router := NewRouter(fc)
	server := &http.Server{
		Addr:    ":" + cfg.HTTPServer.RunPort,
		Handler: router,
	}
	log.Printf("starting http app on port %s", cfg.HTTPServer.RunPort)
	log.Fatal(server.ListenAndServe())
}

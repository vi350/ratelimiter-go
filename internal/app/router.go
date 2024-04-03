package app

import (
	"context"
	"floodcontrol/internal/floodcontrol"
	"net/http"
	"strconv"
	"strings"
)

type FloodControlHandlerInterface interface {
	Check() http.HandlerFunc
}

type Router struct {
	fc floodcontrol.FloodControl
}

func NewRouter(fc floodcontrol.FloodControl) (router *Router) {
	router = &Router{fc: fc}
	return
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if path == "/ping" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
		return
	} else if strings.HasPrefix(path, "/check") {
		check(w, req, r.fc)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("use /check/{id}"))
}

func check(w http.ResponseWriter, r *http.Request, fc floodcontrol.FloodControl) {
	userStrID := strings.TrimPrefix(r.URL.Path, "/check/")
	if userStrID == "" {
		userStrID = "0"
	}
	userID, err := strconv.ParseInt(userStrID, 10, 64)
	if err != nil {
		userID = 0
	}

	allowed, err := fc.Check(context.Background(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	if allowed {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("allowed"))
		return
	} else {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte("denied"))
		return
	}
}

package router

import (
	"labyrinth/internal/state"
	"log/slog"
	"net"
	"net/http"
)

type Router struct {
	http.ServeMux
	State     *state.State
	Logger    *slog.Logger
	SrvConfig ServerConfig
}

func NewRouter() *Router {
	return &Router{ServeMux: *http.NewServeMux(), Logger: slog.Default()}
}

func (r *Router) WithState(state *state.State) *Router {

	r.State = state
	return r

}

func (r *Router) Run() error {

	return http.ListenAndServe(r.SrvConfig.Host.String(), r)

}

type ServerConfig struct {
	Host net.Addr
}

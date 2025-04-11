package router

import (
	"labyrinth/internal/state"
	"log/slog"
	"net"
	"net/http"

	"github.com/rs/cors"
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

func (r *Router) WithServerConfig(serverConfig ServerConfig) *Router {

	r.SrvConfig = serverConfig
	return r

}

func (r *Router) WithState(state *state.State) *Router {

	r.State = state
	return r

}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {

	cors.AllowAll().HandlerFunc(res, req)
	r.ServeHTTP(res, req)

}

func (r *Router) Run() error {
	r.Logger.Info("Labyrinth backend running at", "port", r.SrvConfig.Host.String())
	return http.ListenAndServe(r.SrvConfig.Host.String(), r)

}

type ServerConfig struct {
	Host net.Addr
}

package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/thegoodparticle/vehicle-data-layer/internal/config"
	"github.com/thegoodparticle/vehicle-data-layer/internal/handler"
	"github.com/thegoodparticle/vehicle-data-layer/internal/store"
)

type Router struct {
	config *Config
	router *chi.Mux
}

func NewRouter() *Router {
	return &Router{
		config: NewConfig().SetTimeout(config.GetConfig().Timeout),
		router: chi.NewRouter(),
	}
}

func (r *Router) SetRouters(repository store.Interface) *chi.Mux {
	r.setConfigsRouters()

	r.RouterVehicle(repository)

	return r.router
}

func (r *Router) setConfigsRouters() {
	r.EnableCORS()
	r.EnableLogger()
	r.EnableTimeout()
	r.EnableRecover()
	r.EnableRequestID()
	r.EnableRealIP()
}

func (r *Router) RouterVehicle(repository store.Interface) {
	handler := handler.NewHandler(repository)

	r.router.Route("/", func(route chi.Router) {
		route.Get("/health-check", handler.HealthCheck)
		route.Get("/", handler.HealthCheck)
	})

	r.router.Route("/vehicle-violation", func(route chi.Router) {
		route.Post("/", handler.Post)
		route.Get("/", handler.Get)
		route.Get("/{RegID}", handler.Get)
		route.Put("/{RegID}", handler.Put)
		route.Delete("/{RegID}", handler.Delete)
		route.Options("/", handler.Options)
	})
}

func (r *Router) EnableLogger() *Router {
	r.router.Use(middleware.Logger)
	return r
}

func (r *Router) EnableTimeout() *Router {
	r.router.Use(middleware.Timeout(r.config.GetTimeout()))
	return r
}

func (r *Router) EnableCORS() *Router {
	r.router.Use(r.config.Cors)
	return r
}

func (r *Router) EnableRecover() *Router {
	r.router.Use(middleware.Recoverer)
	return r
}

func (r *Router) EnableRequestID() *Router {
	r.router.Use(middleware.RequestID)
	return r
}

func (r *Router) EnableRealIP() *Router {
	r.router.Use(middleware.RealIP)
	return r
}

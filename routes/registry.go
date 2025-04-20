package routes

import (
	"user-service/controllers"

	"github.com/gin-gonic/gin"

	routes "user-service/routes/user"
)

type Registry struct {
	controller controllers.IControllerRegistry
	group      *gin.RouterGroup
}

type IRouteRegistry interface {
	Serve()
}

func NewRouteRegsitry(controller controllers.IControllerRegistry, group *gin.RouterGroup) IRouteRegistry {
	return &Registry{controller: controller, group: group}
}

func (r *Registry) Serve() {
	r.userRoute().Run()
}

func (r *Registry) userRoute() routes.IUserRoute {
	return routes.NewUserRoute(r.controller, r.group)
}

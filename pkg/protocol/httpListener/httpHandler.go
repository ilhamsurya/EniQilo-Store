package httpListener

import (
	"net/http"
	"projectsphere/eniqlo-store/config"
	productHandler "projectsphere/eniqlo-store/internal/product/handler"
	userHandler "projectsphere/eniqlo-store/internal/staff/handler"
	"projectsphere/eniqlo-store/pkg/middleware/auth"
	"projectsphere/eniqlo-store/pkg/middleware/logger"
	"projectsphere/eniqlo-store/pkg/protocol/msg"

	"github.com/gin-gonic/gin"
)

type HttpHandlerImpl struct {
	productHandler productHandler.ProductHandler
	userHandler    userHandler.UserHandler
	jwtAuth        auth.JWTAuth
}

func NewHttpHandler(
	productHandler productHandler.ProductHandler,
	userHandler userHandler.UserHandler,
	jwtAuth auth.JWTAuth,

) *HttpHandlerImpl {
	return &HttpHandlerImpl{
		productHandler: productHandler,
		userHandler:    userHandler,
		jwtAuth:        jwtAuth,
	}
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Disposition, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (h *HttpHandlerImpl) Router() *gin.Engine {
	server := gin.New()
	server.Use(gin.Recovery(), logger.Logger(), CORSMiddleware())
	server.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, msg.NotFound(msg.ErrPageNotFound))
	})

	server.Static("/v1/docs", "./dist")

	r := server.Group(config.GetString("APPLICATION_GROUP"))

	staff := r.Group("/staff")
	staff.POST("/register", h.userHandler.Register)
	staff.POST("/login", h.userHandler.Login)

	product := r.Group("/product")
	product.Use(h.jwtAuth.JwtAuthUserMiddleware())
	product.POST("/", h.productHandler.Create)
	product.PUT("/:id", h.productHandler.Update)
	product.DELETE("/:id", h.productHandler.Delete)

	return server
}

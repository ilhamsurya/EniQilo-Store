package httpListener

import (
	"net/http"

	"projectsphere/eniqlo-store/pkg/middleware/auth"
	"projectsphere/eniqlo-store/pkg/middleware/logger"
	"projectsphere/eniqlo-store/pkg/protocol/msg"

	"github.com/gin-gonic/gin"
)

type HttpHandlerImpl struct {
	jwtAuth auth.JWTAuth
}

func NewHttpHandler(

	jwtAuth auth.JWTAuth,
) *HttpHandlerImpl {
	return &HttpHandlerImpl{
		jwtAuth: jwtAuth,
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

	return server
}
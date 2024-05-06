package httpListener

import (
	"context"
	"fmt"
	"net/http"

	"projectsphere/eniqlo-store/config"
	"projectsphere/eniqlo-store/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type HttpImpl struct {
	HttpRouter *HttpRouterImpl
	httpServer *http.Server
}

func NewHttpProtocol(
	HttpRouter *HttpRouterImpl,
) *HttpImpl {
	return &HttpImpl{
		HttpRouter: HttpRouter,
	}
}

func (p *HttpImpl) setupRouter() *gin.Engine {
	return p.HttpRouter.Router()
}

func (p *HttpImpl) Listen() {
	app := p.setupRouter()

	serverPort := fmt.Sprintf(":%v", config.GetString("APP_PORT"))
	p.httpServer = &http.Server{
		Addr:    serverPort,
		Handler: app,
	}

	log.Info().Msgf("Server started on Port %s ", serverPort)
	err := p.httpServer.ListenAndServe()
	if err != nil {
		log.Printf(err.Error())
	}
}

func (p *HttpImpl) Shutdown(ctx context.Context) error {
	if err := p.httpServer.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func Start() *HttpImpl {

	// db connection
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal()
	}
	defer db.Close()

	return nil
}

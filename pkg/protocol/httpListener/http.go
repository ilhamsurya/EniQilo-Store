package httpListener

import (
	"context"
	"fmt"
	"net/http"

	"projectsphere/eniqlo-store/config"
	productHandler "projectsphere/eniqlo-store/internal/product/handler"
	productRepository "projectsphere/eniqlo-store/internal/product/repository"
	productService "projectsphere/eniqlo-store/internal/product/service"
	"projectsphere/eniqlo-store/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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

	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%v/%s?%s", config.GetString("DB_USERNAME"), config.GetString("DB_PASSWORD"), config.GetString("DB_HOST"), config.GetString("DB_PORT"), config.GetString("DB_NAME"), config.GetString("DB_PARAMS")))
	if err != nil {
		panic(err.Error())
	}
	postgresConnector := database.NewPostgresConnector(context.TODO(), db)

	productRepo := productRepository.NewProductRepo(postgresConnector)
	productSvc := productService.NewProductService(productRepo)
	productHandler := productHandler.NewProductHandler(productSvc)

	httpHandlerImpl := NewHttpHandler(
		productHandler,
	)
	httpRouterImpl := NewHttpRoute(httpHandlerImpl)
	httpImpl := NewHttpProtocol(httpRouterImpl)

	return httpImpl
}

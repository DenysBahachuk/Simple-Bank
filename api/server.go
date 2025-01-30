package api

import (
	"fmt"

	db "github.com/DenysBahachuk/Simple_Bank/db/sqlc"
	_ "github.com/DenysBahachuk/Simple_Bank/docs"
	"github.com/DenysBahachuk/Simple_Bank/token"
	"github.com/DenysBahachuk/Simple_Bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

//	@title			SimpleBank API
//	@description	API for SmpleBank project
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description

type Server struct {
	config     utils.Config
	store      db.Store
	router     *gin.Engine
	logger     *zap.SugaredLogger
	tokenMaker token.Maker
}

func NewServer(store db.Store, logger *zap.SugaredLogger, config utils.Config) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create a token maker: %w", err)
	}

	server := Server{
		config:     config,
		store:      store,
		logger:     logger,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validateCurrency)
	}

	return &server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(server.authMiddleware())

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)

	authRoutes.POST("/transfers", server.createTransfer)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	server.router = router
}

func (server *Server) Start(addr string) error {
	return server.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

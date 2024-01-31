package api

import (
	"fmt"

	db "github.com/Streamfair/streamfair-idp-svc/db/sqlc"
	"github.com/Streamfair/streamfair-idp-svc/token"
	"github.com/Streamfair/streamfair-idp-svc/util"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// Server serves HTTP requests for the streamfair backend service.
type Server struct {
	config          util.Config
	store           db.Store
	localTokenMaker token.LocalMaker
	router          *gin.Engine
}

// NewServer creates a new HTTP server and setup routing.
func NewServer(config util.Config, store db.Store) (*Server, error) {
	localTokenMaker, err := token.NewLocalPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		panic(fmt.Sprintf("Failed to create local token maker: %v", err))
	}

	server := &Server{
		config:          config,
		store:           store,
		localTokenMaker: localTokenMaker,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/readiness", server.readinessCheck)

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.GET("/users/id/:id", server.getUserByID)
	router.GET("/users/id", server.handleMissingID)
	router.GET("/users/username/:username", server.getUserByUsername)
	router.GET("/users/username", server.handleMissingUsername)
	router.GET("/users/list", server.listUsers)
	router.PUT("/users/update/:id", server.updateUser)
	router.PUT("/users/update", server.handleMissingID)
	router.PUT("/users/update/email/:id", server.updateUserEmail)
	router.PUT("/users/update/email", server.handleMissingID)
	router.PUT("/users/update/username/:id", server.updateUsername)
	router.PUT("/users/update/username", server.handleMissingUsername)
	router.PUT("/users/update/password/:id", server.updateUserPassword)
	router.PUT("/users/update/password", server.handleMissingID)
	router.DELETE("/users/delete/:id", server.deleteUser)
	router.DELETE("/users/delete", server.handleMissingID)

	server.router = router
}

// StartServer starts a new HTTP server on the specified address.
func (server *Server) StartServer(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	switch err := err.(type) {
	case *pgconn.PgError:
		// Handle pgconn.PgError
		switch err.Code {
		case "23505": // unique_violation
			return gin.H{"error": fmt.Sprintf("Unique violation error: %v: %v", err.Message, err.Hint)}
		case "23503": // foreign_key_violation
			return gin.H{"error": fmt.Sprintf("Foreign key violation error: %v: %v", err.Message, err.Hint)}
		default:
			return gin.H{"error": fmt.Sprintf("error: %v", err.Message)}
		}
	default:
		// Handle other types of errors
		return gin.H{"error": err.Error()}
	}
}

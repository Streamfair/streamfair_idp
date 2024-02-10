package api

import (
	"fmt"

	db "github.com/Streamfair/streamfair_idp/db/sqlc"
	"github.com/Streamfair/streamfair_idp/token"
	"github.com/Streamfair/streamfair_idp/util"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

// Server serves HTTP requests for the streamfair backend service.
type Server struct {
	config          util.Config
	store           db.Store
	localTokenMaker token.Maker
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

	// router.POST("/users/login", server.loginUser)

	// for later use
	// authRoutes := router.Group("/").Use(authMiddleware(server.localTokenMaker))

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

// func (server *Server) RunGinServer(config util.Config, store db.Store) {
// 	err := server.StartServer(config.HttpServerAddress)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "server: error while starting server: %v\n", err)
// 	}
// }

// func errorResponse(err error) gin.H {
// 	switch err := err.(type) {
// 	case *pgconn.PgError:
// 		// Handle pgconn.PgError
// 		switch err.Code {
// 		case "23505": // unique_violation
// 			return gin.H{"error": fmt.Sprintf("Unique violation error: %v: %v", err.Message, err.Hint)}
// 		case "23503": // foreign_key_violation
// 			return gin.H{"error": fmt.Sprintf("Foreign key violation error: %v: %v", err.Message, err.Hint)}
// 		default:
// 			return gin.H{"error": fmt.Sprintf("error: %v", err.Message)}
// 		}
// 	default:
// 		// Handle other types of errors
// 		return gin.H{"error": err.Error()}
// 	}
// }

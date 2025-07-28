package main

import (
	"database/sql"
	"ecpc-league/internal/league"
	"ecpc-league/internal/messages"
	"ecpc-league/middleware"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var ENV = os.Getenv("ENV")

func registerLeagueRoutes(r *gin.Engine) {

	leagueGroup := r.Group("/league")

	{
		leagueGroup.GET("/list", func(c *gin.Context) {
			resp, err := league.List(c.Request.Context())
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
			}
			c.JSON(200, resp)
		}) // list all leagues
		leagueGroup.POST("/create", func(c *gin.Context) {
			var input messages.CreateLeagueRequest

			if err := c.ShouldBindJSON(&input); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			if err := league.Create(c.Request.Context(), input.LeagueName, input.Description, input.Scoring); err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gin.H{"message": "success"})
		}) // create a league

		leagueGroup.GET("/get/:id_league", func(c *gin.Context) {
			idLeague := c.Param("id_league")
			id, err := strconv.ParseInt(idLeague, 10, 64)
			if err != nil {
				c.JSON(400, gin.H{"error": "invalid league id"})
				return
			}
			resp, err := league.Get(c, id)

			if err != nil {
				c.JSON(500, gin.H{"erro": err.Error()})
				return
			}

			c.JSON(200, resp)
		})
	}

}

func getDB() (*sql.DB, error) {
	fmt.Println("HERE COMES JOHNNY", os.Getenv("DATABASE_URL"))
	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// Set Gin to release mode in production for better performance and less logging
	// You can change this to gin.DebugMode or gin.TestMode for development

	if ENV == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize db
	db, err := getDB()

	if err != nil {
		panic(err)
	}

	// Initialize Gin router
	router := gin.Default()
	router.Use(middleware.TransactionMiddleware(db))
	registerLeagueRoutes(router)

	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

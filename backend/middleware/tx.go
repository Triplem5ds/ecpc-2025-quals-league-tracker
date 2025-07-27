package middleware

import (
	"database/sql"
	"ecpc-league/engines"
	"log"

	"github.com/gin-gonic/gin"
)

func TransactionMiddleware(db *sql.DB) gin.HandlerFunc {
	/// If I go any further I should probably make this a Context Middleware that has that tx as one of it's stuff
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			log.Printf("failed to start tx %v", err)
			c.AbortWithStatusJSON(500, gin.H{"error": "failed to start tx"})
			return
		}

		ctx = engines.WithTx(ctx, tx)
		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if len(c.Errors) > 0 {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}

	}
}

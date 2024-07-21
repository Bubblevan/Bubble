package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	LoginRoutes(r, db)
	RegisterRoutes(r, db)

	TicketsRoutes(r, db)
	// OrdersRoutes(r, db)
	UsersRoutes(r, db)
	ProfileRoutes(r, db)
}

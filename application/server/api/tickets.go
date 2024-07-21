package api

import (
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TicketsRoutes(r *gin.Engine, db *gorm.DB) {
	// 受保护的路由
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/createTicket", CreateTicketHandler)
	r.GET("/tickets/:id", QueryTicketHandler)
	protected.POST("/createOrder", CreateOrderHandler)
	r.GET("/queryOrder/:id", QueryOrderHandler)

	r.GET("/tickets", func(c *gin.Context) {
		var tickets []models.Ticket
		result := db.Find(&tickets)
		if result.Error != nil {
			utils.Fail(c, http.StatusInternalServerError, "获取门票列表失败")
			return
		}
		utils.Success(c, tickets, "获取门票列表成功")
	})
}

package main

import (
	"backend/middleware"
	"backend/models"
	"backend/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库
	utils.InitDB()
	db := utils.DB

	r := gin.Default()

	// 配置CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 允许访问的源
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/login", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			utils.Fail(c, http.StatusBadRequest, "无效请求")
			return
		}

		var foundUser models.User
		result := db.Where("username = ? AND role = ?", user.Username, user.Role).First(&foundUser)
		if result.Error != nil {
			utils.Fail(c, http.StatusUnauthorized, "用户名或密码错误")
			return
		}

		if foundUser.Password != user.Password {
			utils.Fail(c, http.StatusUnauthorized, "用户名或密码错误")
			return
		}

		// 生成JWT token
		token, err := utils.GenerateToken(foundUser.ID)
		if err != nil {
			utils.Fail(c, http.StatusInternalServerError, "生成token失败")
			return
		}

		utils.Success(c, gin.H{"token": token}, "登录成功")
	})

	r.POST("/register", func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			utils.Fail(c, http.StatusBadRequest, "请求失败: "+err.Error())
			return
		}

		result := db.Create(&user)
		if result.Error != nil {
			utils.Fail(c, http.StatusInternalServerError, "注册失败: "+result.Error.Error())
			return
		}

		if result.RowsAffected == 0 {
			utils.Fail(c, http.StatusConflict, "用户名或邮箱已存在")
			return
		}

		utils.Success(c, nil, "注册成功")
	})

	// 受保护的路由
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/buyTicket", func(c *gin.Context) {
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			utils.Fail(c, http.StatusBadRequest, "无效请求")
			return
		}

		// 假设这里调用了智能合约进行交易

		result := db.Create(&order)
		if result.Error != nil {
			utils.Fail(c, http.StatusInternalServerError, "购买失败")
			return
		}

		utils.Success(c, order, "购买成功")
	})

	r.GET("/tickets", func(c *gin.Context) {
		var tickets []models.Ticket
		result := db.Find(&tickets)
		if result.Error != nil {
			utils.Fail(c, http.StatusInternalServerError, "获取门票列表失败")
			return
		}

		utils.Success(c, tickets, "获取门票列表成功")
	})

	protected.POST("/createTicket", func(c *gin.Context) {
		var ticket models.Ticket
		if err := c.ShouldBindJSON(&ticket); err != nil {
			utils.Fail(c, http.StatusBadRequest, "无效请求")
			log.Printf("Binding error: %v", err)
			return
		}

		// 获取当前登录用户
		loggedInUser, exists := c.Get("currentUser")
		if !exists {
			utils.Fail(c, http.StatusUnauthorized, "用户未登录")
			log.Println("User not logged in")
			return
		}
		user := loggedInUser.(models.User)

		// 检查角色
		if user.Role == nil || !*user.Role {
			utils.Fail(c, http.StatusUnauthorized, "无权限")
			log.Println("User does not have permission")
			return
		}

		log.Printf("Creating ticket: %+v", ticket)

		result := utils.DB.Create(&ticket)
		if result.Error != nil {
			utils.Fail(c, http.StatusInternalServerError, "创建门票失败")
			log.Printf("Database error: %v", result.Error)
			return
		}

		utils.Success(c, ticket, "创建门票成功")
	})

	protected.GET("/orders/:userID", func(c *gin.Context) {
		userID, err := strconv.Atoi(c.Param("userID"))
		if err != nil {
			utils.Fail(c, http.StatusBadRequest, "无效用户ID")
			return
		}

		var orders []models.Order
		result := db.Where("user_id = ?", userID).Find(&orders)
		if result.Error != nil {
			utils.Fail(c, http.StatusInternalServerError, "获取订单失败")
			return
		}

		utils.Success(c, orders, "获取订单成功")
	})

	r.GET("/profile", func(c *gin.Context) {
		log.Println("Received request for /profile")
		loggedInUser, exists := c.Get("currentUser")
		if !exists {
			utils.Fail(c, http.StatusUnauthorized, "用户未登录")
			log.Println("User not logged in")
			return
		}
		user := loggedInUser.(models.User)
		log.Printf("User: %+v", user)

		utils.Success(c, user, "获取个人信息成功")
	})

	r.Run(":8080")
}
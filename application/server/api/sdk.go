package api

import (
	"backend/models"
	"backend/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	SDK           *fabsdk.FabricSDK
	channelClient *channel.Client
	channelName   = "mychannel"
	chaincodeName = "fabasg"
	orgName       = "Org1"
	orgAdmin      = "Admin"
	org1Peer0     = "peer0.org1.example.com"
	org2Peer0     = "peer0.org2.example.com"
)

type Ticket = models.Ticket
type Order = models.Order

func ChannelExecute(funcName string, args [][]byte) (channel.Response, error) {
	configPath := "./config/config.yaml"
	configProvider := config.FromFile(configPath)
	SDK, err := fabsdk.New(configProvider)
	if err != nil {
		log.Fatalf("Failed to create new SDK: %s", err)
	}
	ctx := SDK.ChannelContext(channelName, fabsdk.WithOrg(orgName), fabsdk.WithUser(orgAdmin))
	channelClient, err = channel.New(ctx)
	if err != nil {
		return channel.Response{}, err
	}

	response, err := channelClient.Execute(channel.Request{
		ChaincodeID: chaincodeName,
		Fcn:         funcName,
		Args:        args,
	})
	if err != nil {
		return response, err
	}
	SDK.Close()
	return response, nil
}

// CreateTicketHandler 创建门票的处理器
func CreateTicketHandler(c *gin.Context) {
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

	// 调用智能合约创建门票
	args := [][]byte{
		[]byte(fmt.Sprintf("%d", ticket.ID)),
		[]byte(ticket.EventName),
		[]byte(ticket.Venue),
		[]byte(ticket.EventDate.Format(time.RFC3339)),
		[]byte(fmt.Sprintf("%f", ticket.Price)),
		[]byte(fmt.Sprintf("%d", ticket.Num)),
	}
	response, err := ChannelExecute("CreateTicket", args)
	if err != nil {
		utils.Fail(c, http.StatusInternalServerError, "创建门票失败")
		log.Printf("Chaincode error: %v", err)
		return
	}
	fmt.Println(string(response.Payload))

	// 将门票信息存入数据库
	result := utils.DB.Create(&ticket)
	if result.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "创建门票失败")
		log.Printf("Database error: %v", result.Error)
		return
	}

	utils.Success(c, ticket, "创建门票成功")
}

// QueryTicketHandler 查询门票的处理器
func QueryTicketHandler(c *gin.Context) {
	ticketID := c.Param("id")
	response, err := ChannelExecute("QueryTicket", [][]byte{[]byte(ticketID)})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"result": string(response.Payload)})
}

// CreateOrderHandler 创建订单的处理器
func CreateOrderHandler(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 类型转换
	userIDUint, err := strconv.ParseUint(c.PostForm("userID"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid userID"})
		return
	}
	order.UserID = uint(userIDUint)

	ticketIDUint, err := strconv.ParseUint(c.PostForm("ticketID"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ticketID"})
		return
	}
	order.TicketID = uint(ticketIDUint)

	// 查询对应的票
	var ticket models.Ticket
	if err := utils.DB.First(&ticket, order.TicketID).Error; err != nil {
		log.Printf("Error finding ticket with ID %d: %v", order.TicketID, err)
		c.JSON(404, gin.H{"error": "Ticket not found"})
		return
	}
	log.Printf("Found ticket: %+v", ticket)

	// 判断剩余票数是否足够
	if ticket.Num < uint(order.Quantity) {
		c.JSON(400, gin.H{"error": "Not enough tickets available"})
		return
	}

	// 更新票的数量
	ticket.Num -= uint(order.Quantity)
	if err := utils.DB.Save(&ticket).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update ticket quantity"})
		return
	}

	// 调用智能合约创建订单
	args := [][]byte{
		[]byte(fmt.Sprintf("%d", order.UserID)),       // 转换为字符串
		[]byte(fmt.Sprintf("%d", order.TicketID)),     // 转换为字符串
		[]byte(fmt.Sprintf("%d", order.Quantity)),     // 转换为字符串
		[]byte(fmt.Sprintf("%.2f", order.TotalPrice)), // 保留两位小数，转换为字符串
		[]byte(order.OrderDate.Format(time.RFC3339)),  // 转换为字符串
	}
	response, err := ChannelExecute("CreateOrder", args)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// 将订单信息存入数据库
	result := utils.DB.Create(&order)
	if result.Error != nil {
		utils.Fail(c, http.StatusInternalServerError, "创建订单失败")
		log.Printf("Database error: %v", result.Error)
		return
	}

	c.JSON(200, gin.H{"message": "Order created", "result": string(response.Payload)})
}

// QueryOrderHandler 查询订单的处理器
func QueryOrderHandler(c *gin.Context) {
	orderID := c.Param("id")
	response, err := ChannelExecute("QueryOrder", [][]byte{[]byte(orderID)})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"result": string(response.Payload)})
}

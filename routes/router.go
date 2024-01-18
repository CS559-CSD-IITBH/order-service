package routes

import (
	"github.com/CS559-CSD-IITBH/order-service/controllers"
	"github.com/CS559-CSD-IITBH/order-service/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRouter(order *mongo.Collection, cart *mongo.Collection, store *sessions.FilesystemStore) *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	v1 := r.Group("/api/v1")
	{
		customers := v1.Group("/customer")
		{
			auth := middlewares.SessionAuth(store, "customer")
			customers.Use(auth)

			customers.POST("/savecart", func(c *gin.Context) {
				controllers.SaveCart(c, cart, store)
			})
			customers.GET("/getcart", func(c *gin.Context) {
				controllers.GetCart(c, cart, store)
			})
			customers.POST("/place", func(c *gin.Context) {
				controllers.PlaceOrder(c, order, store)
			})
			customers.POST("/cancel/:orderID", func(c *gin.Context) {
				controllers.CancelOrder(c, order, store)
			})
			customers.GET("/track/:orderID", func(c *gin.Context) {
				controllers.TrackOrder(c, order, store)
			})
		}

		merchants := v1.Group("/merchant")
		{
			auth := middlewares.SessionAuth(store, "merchant")
			merchants.Use(auth)

			merchants.GET("/get", func(c *gin.Context) {
				controllers.GetOrdersForMerchant(c, order, store)
			})
			merchants.POST("/confirm/:orderID", func(c *gin.Context) {
				controllers.ConfirmOrder(c, order, store)
			})
			merchants.POST("/ready/:orderID", func(c *gin.Context) {
				controllers.OrderReadyForPickup(c, order, store)
			})
			merchants.POST("/verify/:orderID", func(c *gin.Context) {
				controllers.VerifyPickup(c, order, store)
			})
		}

		deliveryAgents := v1.Group("/deliveryagent")
		{
			auth := middlewares.SessionAuth(store, "delivery_agent")
			deliveryAgents.Use(auth)

			deliveryAgents.GET("/get", func(c *gin.Context) {
				controllers.GetOrdersForDelivery(c, order, store)
			})
			deliveryAgents.POST("/accept/:orderID", func(c *gin.Context) {
				controllers.AcceptOrder(c, order, store)
			})
			deliveryAgents.POST("/verify/:orderID", func(c *gin.Context) {
				controllers.VerifyDelivery(c, order, store)
			})
		}
	}

	return r
}

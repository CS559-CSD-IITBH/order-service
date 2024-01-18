package controllers

import (
	"context"
	"net/http"

	"github.com/CS559-CSD-IITBH/order-service/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// GetOrdersForDelivery handles the endpoint for retrieving orders for a delivery agent.
func GetOrdersForDelivery(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	deliveryAgentID, _ := session.Values["user_id"].(uint)

	// Query orders for the specific delivery agent
	cursor, err := collection.Find(context.Background(), bson.M{"deliveryInfo.deliveryAgentID": deliveryAgentID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve orders"})
		return
	}
	defer cursor.Close(context.Background())

	var orders []models.Order
	if err := cursor.All(context.Background(), &orders); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode orders"})
		return
	}

	c.JSON(http.StatusOK, orders)
}

// AcceptOrder handles the endpoint for a delivery agent accepting an order.
func AcceptOrder(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	deliveryAgentID, _ := session.Values["user_id"].(uint)

	orderID := c.Param("orderID")

	// Check if the order exists and is assigned to the delivery agent
	existingOrder := models.Order{}
	err := collection.FindOne(context.Background(), bson.M{"_id": orderID, "deliveryInfo.deliveryAgentID": deliveryAgentID}).Decode(&existingOrder)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or does not belong to the delivery agent"})
		return
	}

	// Check if the order is in the correct status for acceptance
	if existingOrder.Status != "Ready" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot accept order. Order status is not ready"})
		return
	}

	// Update the order status to "Assigned"
	update := bson.M{"$set": bson.M{"status": "Assigned"}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": existingOrder.OrderID, "deliveryInfo.deliveryAgentID": deliveryAgentID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order accepted successfully"})
}

// VerifyDelivery handles the endpoint for verifying the delivery of an order by a delivery agent.
func VerifyDelivery(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	deliveryAgentID, _ := session.Values["user_id"].(uint)

	orderID := c.Param("orderID")
	// otp := c.PostForm("otp")

	// Check if the order exists and is assigned to the delivery agent
	existingOrder := models.Order{}
	err := collection.FindOne(context.Background(), bson.M{"_id": orderID, "deliveryInfo.deliveryAgentID": deliveryAgentID}).Decode(&existingOrder)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or does not belong to the delivery agent"})
		return
	}

	// Check if the order is in the correct status for verifying delivery
	if existingOrder.Status != "In-Transit" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot verify delivery. Order status is not in transit"})
		return
	}

	// Verify the OTP (you need to implement this part)
	// expectedOTP := getExpectedOTP(existingOrder) // Implement a function to retrieve the expected OTP
	// if otp != expectedOTP {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
	// 	return
	// }

	// Update the order status to "Delivered"
	update := bson.M{"$set": bson.M{"status": "Delivered"}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": existingOrder.OrderID, "deliveryInfo.deliveryAgentID": deliveryAgentID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify delivery"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Delivery verified successfully"})
}

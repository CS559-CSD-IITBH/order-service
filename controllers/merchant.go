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

// GetOrdersForMerchant handles the endpoint for retrieving orders for a specific merchant.
func GetOrdersForMerchant(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	merchantID, _ := session.Values["user_id"].(uint)

	// Query orders for the specific merchant
	cursor, err := collection.Find(context.Background(), bson.M{"storeID": merchantID})
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

// ConfirmOrder handles the endpoint for confirming an order.
func ConfirmOrder(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	merchantID, _ := session.Values["user_id"].(uint)

	orderID := c.Param("orderID")

	// Check if the order exists and belongs to the merchant.
	existingOrder := models.Order{}
	err := collection.FindOne(context.Background(), bson.M{"_id": orderID, "storeID": merchantID}).Decode(&existingOrder)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or does not belong to the merchant"})
		return
	}

	// Check if the order is in the correct status for confirmation.
	if existingOrder.Status != "Paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot confirm order. Order status is not Paid"})
		return
	}

	// Update the order status to confirmed.
	update := bson.M{"$set": bson.M{"status": "Confirmed"}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": existingOrder.OrderID, "storeID": merchantID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order confirmed successfully"})
}

// OrderReadyForPickup handles the endpoint for marking an order as ready for pickup.
func OrderReadyForPickup(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	merchantID, _ := session.Values["user_id"].(uint)

	orderID := c.Param("orderID")

	// Check if the order exists and belongs to the merchant.
	existingOrder := models.Order{}
	err := collection.FindOne(context.Background(), bson.M{"_id": orderID, "storeID": merchantID}).Decode(&existingOrder)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or does not belong to the merchant"})
		return
	}

	// Check if the order is in the correct status for marking as ready for pickup.
	if existingOrder.Status != "Confirmed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot mark order as ready for pickup. Order status is not confirmed"})
		return
	}

	// Update the order status to ready for pickup.
	update := bson.M{"$set": bson.M{"status": "Ready"}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": existingOrder.OrderID, "storeID": merchantID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark order as ready for pickup"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order marked as ready for pickup"})
}

// VerifyPickup handles the endpoint for verifying pickup by a delivery agent.
func VerifyPickup(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	merchantID, _ := session.Values["user_id"].(uint)

	orderID := c.Param("orderID")
	// otp := c.PostForm("otp")

	// Retrieve the order from the MongoDB collection
	var order models.Order
	err := collection.FindOne(context.Background(), bson.M{"_id": orderID, "storeID": merchantID}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or does not belong to the merchant"})
		return
	}

	// Check if the order is in the correct status for verifying pickup
	if order.Status != "Assigned" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot verify pickup. Order status is not assigned for pickup"})
		return
	}

	// Check if the provided OTP matches the expected OTP (you need to implement this part)
	// expectedOTP := generateExpectedOTP() // You need to implement this function
	// if otp != expectedOTP {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
	// 	return
	// }

	// Update the order status to In-Transit
	update := bson.M{"$set": bson.M{"status": "In-Transit"}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": order.OrderID, "storeID": merchantID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
		return
	}

	// Send OTP to customer (you need to implement this part)

	c.JSON(http.StatusOK, gin.H{"message": "Pickup verified successfully"})
}

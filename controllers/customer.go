package controllers

import (
	"context"
	"net/http"

	"github.com/CS559-CSD-IITBH/order-service/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SaveCart handles the endpoint for saving the user's cart.
func SaveCart(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	userID, _ := session.Values["user_id"].(uint)

	var cart models.Order
	if err := c.BindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Define the filter to find the existing cart
	filter := bson.M{"userID": userID}

	// Replace the existing cart or insert a new one
	_, err := collection.ReplaceOne(context.Background(), filter, cart, options.Replace().SetUpsert(true))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save cart to MongoDB"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart saved successfully"})
}

// GetCart handles the endpoint for retrieving the user's cart.
func GetCart(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	userID, _ := session.Values["user_id"].(uint)

	// Define the filter to find the existing cart
	filter := bson.M{"userID": userID}

	// Find the cart based on the user ID
	var existingCart models.Order
	err := collection.FindOne(context.Background(), filter).Decode(&existingCart)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	c.JSON(http.StatusOK, existingCart)
}

// PlaceOrder handles the endpoint for placing a new order.
func PlaceOrder(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	userID, _ := session.Values["user_id"].(uint)

	var newOrder models.Order
	if err := c.BindJSON(&newOrder); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	newOrder.UserID = userID
	newOrder.Status = "Paid"

	_, err := collection.InsertOne(context.Background(), newOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order placed successfully"})
}

// CancelOrder handles the endpoint for canceling an existing order.
func CancelOrder(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	userID, _ := session.Values["user_id"].(uint)

	orderID := c.Param("orderID")

	existingOrder := models.Order{}
	err := collection.FindOne(context.Background(), bson.M{"_id": orderID, "userID": userID}).Decode(&existingOrder)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or does not belong to the user"})
		return
	}

	if existingOrder.Status == "Cancelled" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order status is already cancelled"})
		return
	}

	update := bson.M{"$set": bson.M{"status": "Cancelled"}}
	_, err = collection.UpdateOne(context.Background(), bson.M{"_id": existingOrder.OrderID, "userID": userID}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order canceled successfully"})
}

// TrackOrder handles the endpoint for tracking the status of an order.
func TrackOrder(c *gin.Context, collection *mongo.Collection, storeSession *sessions.FilesystemStore) {
	session, _ := storeSession.Get(c.Request, "session-name")
	userID, _ := session.Values["user_id"].(uint)

	orderID := c.Param("orderID")

	// Retrieve the order from the MongoDB collection
	var order models.Order
	err := collection.FindOne(context.Background(), bson.M{"_id": orderID, "userID": userID}).Decode(&order)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found or does not belong to the user"})
		return
	}

	// Return the order status in the response
	c.JSON(http.StatusOK, gin.H{"status": order.Status})
}

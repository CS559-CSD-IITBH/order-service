package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	razorpay "github.com/razorpay/razorpay-go"
)

type PageVariables struct {
	OrderId string
	Email   string
	Name    string
	Amount  string
	Contact string
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("*.html")
	router.GET("/", App)
	router.GET("/payment-success", PaymentSuccess)
	router.Run(":8089")
}

func App(c *gin.Context) {

	page := &PageVariables{}
	page.Amount = "11000"
	page.Email = "vemuganti72@gmail.com"
	page.Name = "Vmz Itnagumev"
	page.Contact = "7013962027"
	//Create order_id from the server
	client := razorpay.NewClient("rzp_test_3uZ9y8He9JMg7M", "wPVCmxyR4hNaLWATOWZsbTtT")

	data := map[string]interface{}{
		"amount":   page.Amount,
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	body, err := client.Order.Create(data, nil)
	if err != nil {
		fmt.Println("Problem getting the repository information", err)
		os.Exit(1)
	}

	value := body["id"]

	str := value.(string)
	HomePageVars := PageVariables{ //store the order_id in a struct
		OrderId: str,
		Amount:  page.Amount,
		Email:   page.Email,
		Name:    page.Name,
		Contact: page.Contact,
	}

	c.HTML(http.StatusOK, "app.html", HomePageVars)
}

func PaymentSuccess(c *gin.Context) {

	paymentid := c.Query("paymentid")
	orderid := c.Query("orderid")
	signature := c.Query("signature")

	fmt.Println(paymentid, "paymentid")
	fmt.Println(orderid, "orderid")
	fmt.Println(signature, "signature")
}

func PaymentFaliure(c *gin.Context) {

}

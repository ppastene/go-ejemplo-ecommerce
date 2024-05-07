package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/ppastene/go-shopping-lite/payments"
)

type Response struct {
	Subscription string   `json:"subscription"`
	Optional     []string `json:"optional"`
	Payment      string   `json:"payment"`
}

type Product struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Payment struct {
	Code      string `json:"code"`
	Name      string `json:"name"`
	Available bool   `json:"available"`
}

type JSONData struct {
	Subscriptions []Product `json:"subscriptions"`
	Optionals     []Product `json:"optionals"`
	Payments      []Payment `json:"payments"`
}

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Static("/", "./public")

	// Read json
	jsonFile, err := os.Open("data.json")
	if err != nil {
		log.Fatal("Cannot open json file")
	}
	defer jsonFile.Close()

	jsonRead, _ := io.ReadAll(jsonFile)

	var jsonData JSONData
	json.Unmarshal([]byte(jsonRead), &jsonData)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Subscriptions": jsonData.Subscriptions,
			"Optionals":     jsonData.Optionals,
			"Payments":      jsonData.Payments,
		})
	})

	app.Post("/checkout", func(c *fiber.Ctx) error {
		response := new(Response)
		if err := c.BodyParser(response); err != nil {
			return err
		}
		var products []Product
		var total int
		var payment Payment
		var paymentData map[string]string = make(map[string]string)

		for _, v := range jsonData.Subscriptions {
			if v.Code == response.Subscription {
				products = append(products, v)
				break
			}
		}
		if len(response.Optional) > 0 {
			for _, v := range response.Optional {
				for _, v2 := range jsonData.Optionals {
					if v == v2.Code {
						products = append(products, v2)
						break
					}
				}
			}
		}
		for _, v := range products {
			total += v.Price
		}

		for _, v := range jsonData.Payments {
			if v.Code == response.Payment {
				payment = v
				break
			}
		}

		if payment.Code == "webpay" {
			var buyOrder string = time.Now().Format("20060102_150405")
			var sessionId string = fmt.Sprintf("%s_session_id", time.Now().Format("20060102_150405"))
			webpay := payments.NewWebpay("597055555532", "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C")
			res, _ := webpay.Create(buyOrder, sessionId, float64(total), "http://localhost:3000/return?payment=webpay")
			paymentData["url"] = res.Url
			paymentData["token"] = res.Token
		}

		return c.Render("checkout", fiber.Map{
			"Products":    products,
			"Total":       total,
			"Payment":     payment,
			"PaymentData": paymentData,
		})
	})

	app.Get("/return", func(c *fiber.Ctx) error {
		if c.Query("payment") == "webpay" {
			webpay := payments.NewWebpay("597055555532", "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C")
			status, _ := webpay.Commit(c.Query("token_ws"))
			if status.ResponseCode != 0 {
				return c.Redirect("/failure", 301)
			}
			return c.Redirect("/success", 301)
		}
		return c.Redirect("/failure", 301)
	})

	app.Get("/success", func(c *fiber.Ctx) error {
		return c.Render("finish", fiber.Map{
			"Message": "The payment has been successful",
		})
	})

	app.Get("/failure", func(c *fiber.Ctx) error {
		return c.Render("finish", fiber.Map{
			"Message": "The payment has failed. Please try again",
		})
	})

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

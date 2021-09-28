package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"go-short-url/help"
	"log"
	"time"
)

var rdb *redis.Client
var ctx = context.Background()

type CreateParma struct {
	Url string `json:"url" xml:"url" form:"url"`
}

func init() {
	rdb = newClient()
}

func main() {
	app := fiber.New()

	app.Get("/s", func(c *fiber.Ctx) error {
		c.Status(500)
		return c.SendString("miss hash code")
	})

	app.Get("/s/:hash", func(c *fiber.Ctx) error {
		key := c.Params("hash")

		val, err := rdb.Get(ctx, key).Result()
		if err == redis.Nil {
			fmt.Println("key2 does not exist")
			c.Status(404)
			return c.SendString("not found")
		} else if err != nil {
			panic(err)
		}

		return c.Redirect(val)
	})

	app.Post("/create", func(c *fiber.Ctx) error {
		cp := new(CreateParma)

		if err := c.BodyParser(cp); err != nil {
			fmt.Errorf(err.Error())
			c.Status(400)
			return c.SendString(err.Error())
		}

		uniqCode := generateUniqCode(cp.Url)

		s := fmt.Sprintf("http:localhost:3000/s/%s", uniqCode)
		return c.SendString(s)
	})

	app.Listen(":3000")
}

func newClient() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping(ctx).Result()
	log.Println(pong)
	if err != nil {
		log.Fatalln(err)
	}
	return client
}

func generateUniqCode(url string) string {
	uniqCode := help.Encode(uint64(time.Now().Nanosecond()))

	result := rdb.SetNX(ctx, uniqCode, url, 0)
	err := result.Err()
	if err != nil {
		panic(err)
	}

	isPass, _ := result.Result()
	if !isPass {
		uniqCode = generateUniqCode(url)
		fmt.Println(uniqCode)
	}

	return uniqCode
}

package main

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/jasonlvhit/gocron"
	"log"
	"os"
	"project_restapi/routes"
)

func main() {
	defer routes.DB.Close()

	log.Println("Web Server")
	app := fiber.New()
	routes.Route(app)
	port := fmt.Sprintf(":%v", os.Getenv("PORT"))

	goUpdate()
	
	app.Listen(port)
	
}

func goUpdate() {
	log.Println("Cron Job")
	gocron.Every(1).Hour().Do(func() {
		_, err := routes.DB.Exec(context.Background(), `CALL updateCoin()`)
		if err != nil {
			log.Println(err)
		}

		log.Println("Cron job done!")
	})

	go func() {
		<-gocron.Start()
	}()
}

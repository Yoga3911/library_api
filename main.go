package main

import (
	// "context"
	"github.com/gofiber/fiber/v2"
	// "github.com/jasonlvhit/gocron"
	"log"
	"os"
	// "os/signal"
	"project_restapi/routes"
	// "syscall"
)

func main() {
	defer routes.DB.Close()
	app := fiber.New()
	routes.Route(app)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
	// go func() {
	// 	log.Println("Web Server")
	// }()

	// go func() {
	// 	goUpdate()
	// }()

	// // Block main goroutine process
	// osCh := make(chan os.Signal)
	// stopCh := make(chan bool)
	// signal.Notify(osCh, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	<-osCh
	// 	log.Println("exiting process")
	// 	stopCh <- true
	// 	os.Exit(0)
	// }()
	// <-stopCh
}

// func goUpdate() {
// 	log.Println("Cron Job")
// 	gocron.Every(1).Hour().Do(func() {
// 		_, err := routes.DB.Exec(context.Background(), `CALL updateCoin()`)
// 		if err != nil {
// 			log.Println(err)
// 		}

// 		log.Println("Cron job done!")
// 	})

// 	<-gocron.Start()
// }

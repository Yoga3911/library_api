package routes

import (
	"project_restapi/cache"
	"project_restapi/configs"
	"project_restapi/controllers"
	"project_restapi/repository"
	"project_restapi/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	DB *pgxpool.Pool = configs.DatabaseConnection()

	redisS cache.RedisC = cache.NewRedisC("redis-15271.c56.east-us.azure.cloud.redislabs.com:15271", "dnxQS1TFH2i7EAQnWjGDoC8dNBPTdNbh", 0)

	cacheS cache.Cache = cache.NewCache()

	jwtS services.JWTS = services.NewJWTS()

	userR repository.UserR  = repository.NewUserR(DB)
	userS services.UserS    = services.NewUserS(DB, userR, jwtS)
	userC controllers.UserC = controllers.NewUserC(userS, cacheS)

	authR repository.AuthR  = repository.NewAuthR(DB)
	authS services.AuthS    = services.NewAuthS(authR, jwtS)
	authC controllers.AuthC = controllers.NewAuthC(authS, cacheS, redisS)

	bookR repository.BookR  = repository.NewBookR(DB)
	bookS services.BookS    = services.NewBookS(bookR, jwtS)
	bookC controllers.BookC = controllers.NewBookC(bookS, cacheS)
)

func Route(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/auth/login", authC.Login)
	api.Post("/auth/register", authC.Register)
	api.Get("/auth/logout", authC.Logout)

	api.Get("/auth/otp/:email", authC.SendVerif)
	api.Post("/auth/otp", authC.Verif)

	api.Get("/users", userC.GetAll)
	api.Get("/user", userC.GetByToken)
	api.Put("/user", userC.UpdateUser)
	api.Patch("/user/pass", userC.ChangePassword)
	api.Patch("/user", userC.DeleteUser)
	api.Get("/user/book/:id", userC.UserBookById)
	api.Get("/user/req", userC.GetRequest)
	api.Post("/user/req", userC.RequestAdmin)
	api.Post("/user/rev", userC.RequestAnswer)

	api.Get("/books", bookC.GetAllBook)
	api.Get("/books/:id", bookC.GetBookByGenre)
	api.Post("/book", bookC.AddBook)
	api.Put("/book/:id", bookC.UpdateBook)
	api.Delete("/book/:id", bookC.DeleteBook)

	api.Get("/book/reviews/:id", bookC.GetReview)
	api.Post("/book/review/:id", bookC.AddReview)
	api.Put("/book/review/:id", bookC.UpdateReview)

	api.Post("/get/book/:id", userC.TakeBook)
	api.Get("/get/book", userC.GetBookById)
	api.Delete("/get/book/:id", userC.DeleteBookById)

	api.Get("/monitor", monitor.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK!")
	})
}

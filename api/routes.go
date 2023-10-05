package api

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"hyperpage/controllers"
	"hyperpage/initializers"
	"hyperpage/middleware"
)

func Register(micro *fiber.App) {

	micro.Route("/auth", func(router fiber.Router) {

		router.Post("/register", controllers.SignUpUser)

		router.Post("/login", controllers.SignInUser)
		router.Post("/forgotpassword", controllers.ForgotPassword)

		router.Patch("/resetpassword/:resetToken", controllers.ResetPassword)

		router.Get("/verifyemail/:verificationCode", controllers.VerifyEmail)
		router.Get("/logout", controllers.LogoutUser)
		router.Get("/refresh/:refreshToken", controllers.RefreshAccessToken)

	})

	micro.Route("/followers", func(router fiber.Router) {
		router.Post("/scribe", middleware.DeserializeUser, controllers.Scribe)
		router.Post("/unscribe", middleware.DeserializeUser, controllers.Unscribe)
		router.Get("/get", middleware.DeserializeUser, controllers.GetFollowers)

	})

	micro.Route("/domains", func(router fiber.Router) {
		router.Get("/get", controllers.GetDomain)
	})

	micro.Route("/site", func(router fiber.Router) {
		router.Post("/update", middleware.DeserializeUser, controllers.UpdateSite)
		router.Get("/get", middleware.DeserializeUser, controllers.GetSite)

	})

	micro.Route("/users", func(router fiber.Router) {
		router.Get("/myTime", controllers.MyTime)
		router.Post("/deletme", middleware.DeserializeUser, controllers.DeleteUserWithRelations)
		router.Post("/setvip", middleware.DeserializeUser, controllers.SetVipUser)

		router.Post("/sendrequestcall", controllers.SendBotCallRequest)
		router.Get("/me", middleware.DeserializeUser, controllers.GetMe)
		router.Get("/getmefirst", middleware.DeserializeUser, controllers.GetMeFirst)
		router.Post("/addbalance", middleware.DeserializeUser, controllers.AddBalance)
		router.Post("/plan", middleware.DeserializeUser, controllers.Plan)
	})

	micro.Route("/billing", func(router fiber.Router) {
		router.Get("/transactions", middleware.DeserializeUser, controllers.GetTransactions)
	})

	micro.Route("/calls", func(router fiber.Router) {
		router.Get("/makecall", middleware.DeserializeUser, controllers.MakeCall)
	})

	micro.Route("/cities", func(router fiber.Router) {
		router.Get("/all", controllers.GetCities)
		router.Get("/query", controllers.GetName)
	})

	micro.Route("/guilds", func(router fiber.Router) {
		router.Get("/all", controllers.GetGuilds)
	})

	micro.Route("/profile", func(router fiber.Router) {
		router.Get("/get", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.GetProfile)
		router.Patch("/save", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.UpdateProfile)
		router.Patch("/saveAdditional", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.UpdateProfileAdditional)

		router.Patch("/photos", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.UpdateProfilePhotos)
		router.Post("/documents", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.NewProfileDocuments)
		router.Patch("/documents", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.UpdateProfileDocuments)
		router.Delete("/documents/:id", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.DeleteProfileDocuments)

		router.Get("/getdocuments", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.GetDocuments)
	})

	micro.Route("/profiles", func(router fiber.Router) {
		router.Get("/get", controllers.GetAllProfile)
		router.Get("/get/:name", controllers.GetProfileGuest)
	})

	micro.Route("/payment", func(router fiber.Router) {
		router.Post("/invoice", middleware.DeserializeUser, controllers.CreateInvoice)
		router.Post("/pending", controllers.Pending)

	})

	micro.Route("/profilehashtags", func(router fiber.Router) {
		router.Post("/addhashtag", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.AddHashTagProfile)
		router.Get("/findTag", controllers.SearchHashTagProfile)
	})

	micro.Route("/blog", func(router fiber.Router) {
		router.Get("/list", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.GetAllBlogs)
		router.Post("/makearchive/:id", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.SendToArchive)
		router.Post("/search", middleware.DeserializeUser, controllers.SearchBlogByTitle)
		router.Post("/addblogtime", middleware.DeserializeUser, controllers.AddBlogTime)
		router.Post("/addhashtag", middleware.DeserializeUser, controllers.AddHashTag)
		router.Get("/findTag", controllers.SearchHashTag)

		router.Get("/getAllByUser/:id", controllers.GetAllByUser)

		router.Get("/listAll", controllers.GetAll)

		router.Get("/random", controllers.GetRandom)

		router.Get("/:id", controllers.GetBlogById)
		router.Post("/create", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), middleware.CheckProfileFilled(), controllers.CreateBlog)
		router.Post("/create/photos", middleware.DeserializeUser, controllers.CreateBlogPhoto)
		router.Get("/edit/:id", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.EditBlogGetId)
		router.Patch("/patch/:id", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.UpdateBlog)
		router.Delete("/delete/:id", middleware.DeserializeUser, middleware.CheckRole([]string{"admin", "user", "vip"}), controllers.DeleteBlog)
	})

	micro.Route("/files", func(router fiber.Router) {
		router.Post("/upload/file", middleware.DeserializeUser, middleware.CheckProfileFilled(), controllers.UploadPdf)
		router.Post("/upload", middleware.DeserializeUser, middleware.CheckProfileFilled(), controllers.UploadImage)
		router.Post("/upload/images", middleware.DeserializeUser, middleware.CheckProfileFilled(), controllers.UploadImages)

	})

	micro.Route("/server", func(router fiber.Router) {

		ctx := context.TODO()
		value, err := initializers.RedisClient.Get(ctx, "statusHealth").Result()

		router.Get("/healthchecker", func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"status":  "success",
				"message": value,
			})
		})

		if err == redis.Nil {
			fmt.Println("key: statusHealth does not exist")
		} else if err != nil {
			panic(err)
		}

	})

	micro.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exists", path),
		})
	})
}
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"github.com/Eco-Sort/eco_sort_backend/config"
	"github.com/Eco-Sort/eco_sort_backend/delivery/http_api"
	"github.com/Eco-Sort/eco_sort_backend/domain"
	"github.com/Eco-Sort/eco_sort_backend/library/db"
	"github.com/Eco-Sort/eco_sort_backend/library/middleware"
	"github.com/Eco-Sort/eco_sort_backend/repository/mariadb"
	"github.com/Eco-Sort/eco_sort_backend/service/auth"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var wgInstance *sync.WaitGroup

// Service
var (
	authService domain.AuthService
)

func NewWaitGroup() *sync.WaitGroup {
	if wgInstance == nil {
		wgInstance = &sync.WaitGroup{}
	}
	return wgInstance
}

var wg = NewWaitGroup()

func bootstrap() {
	godotenv.Load()

}

func bootstrapServices() {
	serviceTimeOut := 180 * time.Second

	masterRepoUser := mariadb.NewMariadbUserRepository(db.Mariadb)

	authService = auth.NewAuthService(
		serviceTimeOut, masterRepoUser)
}

func bootstrapFiber() *fiber.App {

	app := fiber.New(
		fiber.Config{
			DisableStartupMessage: false,
			JSONEncoder:           sonic.Marshal,
			JSONDecoder:           sonic.Unmarshal,
			Prefork:               false,
			ServerHeader:          "ECO SORT",
			AppName:               config.GetAppName(),
			ReadTimeout:           time.Second * 60,
			CaseSensitive:         true,
			BodyLimit:             25 * 1024 * 1024,
			Network:               "tcp",
			ErrorHandler: func(c *fiber.Ctx, err error) error {
				// Status code defaults to 500
				code := fiber.StatusInternalServerError

				// Retrieve the custom status code if it's an fiber.*Error
				if e, ok := err.(*fiber.Error); ok {
					code = e.Code
				}

				if code == 404 {
					//Do something
				}

				// Return from handler
				return nil
			},
		},
	)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("OK")
	})

	app.Use(

		func(c *fiber.Ctx) error {
			defer func() {
				if r := recover(); r != nil {
					id := time.Now().Unix()
					fmt.Printf("Server Panic Occured [CrashID: %d]\n[PanicMessage]:%s\n[PanicStack]:%s\n", id, r, string(debug.Stack()))
					c.Status(fiber.StatusInternalServerError).SendString("Server Error")
				}
			}()
			return c.Next()
		},

		cors.New(cors.Config{
			AllowCredentials: true,
			AllowHeaders:     "*",
		}),

		func(c *fiber.Ctx) error {
			c.Locals("auth_service", &authService)
			c.Locals("wg", wg)
			return c.Next()
		},

		func(c *fiber.Ctx) error {
			err := c.Next()
			if c.Response().StatusCode() == 500 {
				fmt.Println(c.Body())
			}
			return err
		},
	)
	return app
}

func migrateMariadb(db *gorm.DB) {
	db.AutoMigrate(
		&domain.User{},
	)
}

func main() {
	bootstrap()
	db.InitMariadb()
	migrateMariadb(db.Mariadb)
	bootstrapServices()

	startHttp := flag.Bool("start-http", false, "Start HTTP server")

	flag.Parse()

	if *startHttp {
		initHttp()
	}
}

func initHttp() {
	app := bootstrapFiber()

	apiRoute := app.Group("/api")
	// Web API Route V1
	wV1ApiRoute := apiRoute.Group("/v1")
	http_api.NewAuthHttpApiDelivery(wV1ApiRoute, authService)

	// Admin Route
	adminApiRoute := wV1ApiRoute.Group("/admin", middleware.ValidateJWT)
	http_api.NewAdminUserHttpApiDelivery(adminApiRoute)

	//Client Route
	// clientApiRoute := wV1ApiRoute.Group("/app")

	//Public Route
	// publicApiRoute := wV1ApiRoute.Group("/public")

	go func() {
		if err := app.Listen(":" + config.GetAppPort()); err != nil {
			fmt.Println("Error starting HTTP server: ", err)
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	<-interrupt

	wg.Wait()

	if err := app.Shutdown(); err != nil {
		fmt.Println(err)
	}
}

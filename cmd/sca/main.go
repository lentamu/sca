package main

import (
	"flag"
	"log"
	"time"

	"sca/internal/config"
	"sca/internal/handler"
	"sca/internal/service"
	"sca/internal/storage"
	"sca/pkg/cache"
	"sca/pkg/database/mysql"
	"sca/pkg/errors"
	pkgvalidator "sca/pkg/validator"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
	configPath := flag.String("config", "configs/stub.toml", "path to config file")
	flag.Parse()

	conf, err := config.Load(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	db, err := mysql.Connect(conf.Mysql.Username, conf.Mysql.Password, conf.Mysql.Host, conf.Mysql.Database)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	redisCache := cache.NewRedisCache(cache.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	v := validator.New(validator.WithRequiredStructEnabled())
	pkgvalidator.RegisterValidators(v)
	pkgvalidator.InitBreedValidator(redisCache, conf.Breeds.Url, "breeds", time.Hour)

	store := storage.NewStorage(db)

	s := service.NewService(&service.Depends{
		Storage: store,
		Cache:   redisCache,
	})

	app := fiber.New(fiber.Config{
		ErrorHandler:    errors.ErrorHandler,
		JSONEncoder:     json.Marshal,
		JSONDecoder:     json.Unmarshal,
		StructValidator: pkgvalidator.NewStructValidator(v),
	})

	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${ip} ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	h := handler.NewHandler(s)
	h.RegisterRoutes(app)

	log.Fatal(app.Listen(conf.ListenAddr))
}

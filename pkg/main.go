package main

import (
	"errors"
	"fmt"
	"github.com/carlosrodriguesf/bank-api/pkg/api"
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	apimodel "github.com/carlosrodriguesf/bank-api/pkg/api/model"
	"github.com/carlosrodriguesf/bank-api/pkg/app"
	"github.com/carlosrodriguesf/bank-api/pkg/repository"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/cache"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/closer"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/db"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/go-redis/redis/v8"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	emiddleware "github.com/labstack/echo/v4/middleware"

	"os"
	"runtime"
	"strings"
)

func runMigrations(log logger.Logger) {
	log = log.WithPreffix("migration")

	log.Info("starting")
	if err := getMigration(log).Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}
	log.Info("ended")
}

func getMigration(log logger.Logger) *migrate.Migrate {
	dir, err := os.Getwd()
	if err != nil {
		logger.New("").Fatal(err)
	}
	m, err := migrate.New(fmt.Sprintf("file://%s/migrations", dir), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return m
}

func startDB(log logger.Logger) (db.ExtendedDB, error) {
	log = log.WithPreffix("postgres")

	log.Info("connection")
	sql, err := sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Error(err)
		return nil, err
	}
	log.Info("connected")
	return db.NewExtendedDB(sql), nil
}

func startCache(log logger.Logger) (cache.Cache, error) {
	log = log.WithPreffix("redis")

	log.Info("connecting")
	opts, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Info(err)
		return nil, err
	}
	log.Info("connected")

	return cache.NewRedisCache(redis.NewClient(opts)), nil
}

func startEcho(log logger.Logger) *echo.Echo {
	e := echo.New()
	e.Use(emiddleware.CORS())
	e.Use(emiddleware.RequestID())
	e.Use(emiddleware.Logger())
	e.Use(emiddleware.BodyLimit("1M"))
	e.Use(emiddleware.Recover())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		apiErr, ok := err.(*apierror.ApiError)
		if !ok {
			log.Error(err)
			apiErr = apierror.ErrInternal
		}
		err = c.JSON(apiErr.Code, apimodel.Response{Error: apiErr})
		if err != nil {
			log.Error(err)
		}
	}

	return e
}

func getProjectDir() string {
	_, file, _, _ := runtime.Caller(0)
	return strings.Replace(file, "main.go", "", 1)
}

func main() {
	log := logger.New(getProjectDir())

	err := godotenv.Load(".env")
	if err != nil {
		log.Info(err)
	}

	runMigrations(log)

	connDB, err := startDB(log)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.MustClose(log, connDB)

	connCache, err := startCache(log)
	if err != nil {
		log.Fatal(err)
	}
	defer closer.MustClose(log, connCache)

	e := startEcho(log)

	repositoryContainer := repository.NewContainer(repository.Options{
		Logger: log,
		DB:     connDB,
	})
	appContainer := app.NewContainer(app.Options{
		Logger:     log,
		Cache:      connCache,
		Repository: repositoryContainer,
	})
	api.Register(e, apimodel.Options{
		Logger: log,
		App:    appContainer,
	})

	if err = e.Start(":" + os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}

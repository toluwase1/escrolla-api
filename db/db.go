package db

import (
	"escrolla-api/config"
	"escrolla-api/models"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormDB struct {
	DB *gorm.DB
}

func GetDB(c *config.Config) *GormDB {
	gormDB := &GormDB{}
	gormDB.Init(c)
	return gormDB
}

func (g *GormDB) Init(c *config.Config) {
	g.DB = getPostgresDB(c)

	if err := migrate(g.DB); err != nil {
		log.Fatalf("unable to run migrations: %v", err)
	}
}

func getPostgresDB(c *config.Config) *gorm.DB {
	log.Printf("Connecting to postgres: %+v", c)
	//postgresDSN := "postgres://postgres:toluwase@localhost:5432/dbname?sslmode=disable"

	//postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d TimeZone=Africa/Lagos",
	//	c.PostgresHost, c.PostgresUser, c.PostgresPassword, c.PostgresDB, c.PostgresPort)

	log.Printf("Connecting to postgres")
	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		"localhost", "postgres", "toluwase", "escrolla", 5432) //, sslmode)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level Info, Silent, Warn, Error
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
	gormConfig := &gorm.Config{
		Logger: newLogger,
	}
	if c.Env == "prod" {
		gormConfig = &gorm.Config{}
	}
	postgresDB, err := gorm.Open(postgres.Open(postgresDSN), gormConfig)
	if err != nil {
		log.Fatal(err)
	}
	return postgresDB
}

func migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{}, &models.BlackList{}, &models.Transaction{}, &models.Order{})
	if err != nil {
		return fmt.Errorf("migrations error: %v", err)
	}

	return nil
}

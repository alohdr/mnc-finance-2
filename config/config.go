package config

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"mnc-finance/entity"
	"os"
	"time"
)

var (
	DB *gorm.DB
)

// Setup Database

func SetupDatabase() *gorm.DB {
	dbLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	gormConfig := &gorm.Config{
		// enhance performance config
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		Logger:                 dbLogger,
	}

	var DB_USERNAME = os.Getenv(`POSTGRES_READ_USER`)
	var DB_PASSWORD = os.Getenv(`POSTGRES_READ_PASSWORD`)
	var DB_NAME = os.Getenv(`POSTGRES_READ_DB`)
	var DB_HOST = os.Getenv(`POSTGRES_READ_HOST`)
	var DB_PORT = os.Getenv(`POSTGRES_READ_PORT`)

	var err error
	dsn := "host=" + DB_HOST + " user=" + DB_USERNAME + " password=" + DB_PASSWORD + " dbname=" + DB_NAME + " port=" + DB_PORT + " sslmode=disable TimeZone=Asia/Shanghai"

	openConn, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Println("Error connecting to database : error=", err)
		return nil
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: openConn, // data source name
	}), gormConfig)

	if err != nil {
		log.Println("Error connecting to database : error=", err)
		return nil
	}

	//db.Exec(`SET timezone TO '+07';`)
	//db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	db.AutoMigrate(&entity.User{}, &entity.Transaction{})

	return db
}

// Setup RabbitMQ

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func SetUpRabbitMQ() RabbitMQ {
	url := fmt.Sprintf("amqp://%v:%v@%v:%v/", os.Getenv("RabbitUsername"), os.Getenv("RabbitPassword"), os.Getenv("RabbitHost"), os.Getenv("RabbitPort"))

	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println("Url: ", url)
		fmt.Println("RabbitMq Connection Refused")
	}

	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}

	err = channel.ExchangeDeclare(
		os.Getenv("RabbitExchange"),
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Url: ", url)
	fmt.Println("RabbitMq Connection Established")

	return RabbitMQ{
		Conn:    conn,
		Channel: channel,
	}
}

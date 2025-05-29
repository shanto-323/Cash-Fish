package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	walletservice "cash-fish/wallet-service"

	"github.com/kelseyhightower/envconfig"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"github.com/tinrab/retry"
)

type Config struct {
	DatabaseUrl string `envconfig:"DATABASE_DSN"`
	RedisUrl    string `envconfig:"REDIS_URL"`
	RabbitmqUrl string `envconfig:"RABBITMQ_URL"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal("could not get env variabls", err)
	}

	var repository walletservice.Repository
	var redisClient *redis.Client
	var publisher *amqp.Channel
	retry.ForeverSleep(
		2*time.Second,
		func(_ int) error {
			repository, err = walletservice.NewTransactionRepository(cfg.DatabaseUrl)
			if err != nil {
				log.Print("main", err)
				return err
			}

			redisClient, err = setUpRedis(cfg.RedisUrl)
			if err != nil {
				log.Print(err)
				return err
			}

			publisher, err = setUpRabbitmq(cfg.RabbitmqUrl)
			if err != nil {
				log.Print(err)
				return err
			}

			log.Print("everything working")
			return nil
		},
	)

	s := walletservice.NewService(repository, redisClient, publisher)
	log.Fatal(walletservice.ListenGRPC(s, 8080))
}

func setUpRedis(redisUrl string) (*redis.Client, error) {
	parseUrl, err := url.Parse(redisUrl)
	if err != nil {
		return nil, err
	}
	password, _ := parseUrl.User.Password()
	config := &redis.Options{
		Addr:     parseUrl.Host,
		Password: password,
		DB:       0,
	}
	redisClient := redis.NewClient(config)
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("cant connect redis database: %s", err)
	}

	return redisClient, nil
}

func setUpRabbitmq(url string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, nil
}

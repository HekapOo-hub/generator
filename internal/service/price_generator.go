package service

import (
	"fmt"
	"github.com/HekapOo-hub/generator/internal/config"
	"github.com/HekapOo-hub/generator/internal/model"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type PriceGenerator struct {
	client    *redis.Client
	avgPrices map[string]float64
}

func NewPriceGenerator() (*PriceGenerator, error) {
	redisCfg, err := config.NewRedisConfig()
	if err != nil {
		log.Warnf("new price generator: %v", err)
		return nil, err
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password,
		DB:       redisCfg.DB,
	})

	pg := PriceGenerator{client: redisClient}
	pg.avgPrices = map[string]float64{
		"silver":    25.955,
		"gold":      1986,
		"bitcoin":   0.066,
		"eBay":      51.3,
		"crude oil": 106,
	}
	return &pg, nil
}

func (pg *PriceGenerator) send(msg map[string]interface{}) error {
	err := pg.client.XAdd(&redis.XAddArgs{
		Stream:       config.RedisStream,
		MaxLen:       0,
		MaxLenApprox: 0,
		ID:           "",
		Values: map[string]interface{}{
			"silver":    msg["silver"],
			"gold":      msg["gold"],
			"bitcoin":   msg["bitcoin"],
			"eBay":      msg["eBay"],
			"crude oil": msg["crude oil"],
		},
	}).Err()
	if err != nil {
		return fmt.Errorf("redis stream sending create human request error %w", err)
	}
	return nil
}
func (pg *PriceGenerator) GeneratePrices() {
	generator := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		generatedPrices := make(map[string]interface{}, len(pg.avgPrices))
		for key, val := range pg.avgPrices {
			gp := model.GeneratedPrice{Symbol: key}
			gp.Bid = 0.95*val + generator.Float64()*0.15*val
			gp.Ask = 0.95*val + generator.Float64()*0.15*val
			generatedPrices[key] = gp
		}
		err := pg.send(generatedPrices)
		if err != nil {
			log.Warnf("generate prices: %v", err)
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

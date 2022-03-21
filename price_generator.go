package main

import (
	"github.com/HekapOo-hub/generator/internal/service"
	"time"
)

func main() {
	priceGenerator, err := service.NewPriceGenerator()
	if err != nil {
		return
	}
	priceGenerator.GeneratePrices()
	time.Sleep(time.Minute)
}

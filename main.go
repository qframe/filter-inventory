package main

import (
	"time"
	"fmt"
	"golang.org/x/net/context"

	"github.com/docker/docker/client"
	dt "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/zpatrick/go-config"

	"github.com/qnib/qframe-types"
	"github.com/qframe/filter-inventory/lib"
	"sync"
)

func Run(qChan qtypes.QChan, cfg *config.Config, name string) {
	p,_ := filter_inventory.New(qChan, cfg, name)
	p.Run()
}

func main() {
	qChan := qtypes.NewQChan()
	qChan.Broadcast()
	cfgMap := map[string]string{
		"filter.inventory.inputs": "docker-events",
		"log.level": "debug",
		"filter.inventory.ticker-ms": "2500",
	}

	cfg := config.NewConfig(
		[]config.Provider{
			config.NewStatic(cfgMap),
		},
	)
	// Inventory Filter
	p, _ := filter_inventory.New(qChan, cfg, "inventory")
	go p.Run()
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}


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

func initConfig() (config *container.Config) {
	return &container.Config{Image: "alpine", Volumes: nil, Cmd: []string{"/bin/sleep", "5"}, AttachStdout: false}
}

func hConfig() (config *container.HostConfig) {
	return &container.HostConfig{AutoRemove: true}
}

func startCnt(cli *client.Client, name string, sec int) {
	time.Sleep(time.Duration(sec)*time.Second)
	// Start container
	create, err := cli.ContainerCreate(context.Background(), initConfig(), hConfig(), nil, name)
	if err != nil {
		fmt.Println(err)
	}
	err = cli.ContainerStart(context.Background(), create.ID, dt.ContainerStartOptions{})
	if err != nil {
		fmt.Println(err)
	}
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


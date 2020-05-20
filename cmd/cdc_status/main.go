package main

import (
	"fmt"
	"log"
	"net/http"
	"os/user"
	"time"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"

	"github.com/michaelmosher/monitoring/pkg/cdc"
	"github.com/michaelmosher/monitoring/pkg/octopus"
	octopus_http "github.com/michaelmosher/monitoring/pkg/octopus/http"
)

type octopusConfig struct {
	InstanceURL string   `hcl:"instanceURL"`
	APIKey      string   `hcl:"apiKey"`
	Space       string   `hcl:"space"`
	CDCProjects []string `hcl:"cdcProjects"`
	Extra       hcl.Body `hcl:",remain"`
}

type metriclyConfig struct {
	Username string `hcl:"Username"`
	Password string `hcl:"Password"`
}

type mainConfig struct {
	Metricly metriclyConfig `hcl:"Metricly,block"`
	Octopus  octopusConfig  `hcl:"Octopus,block"`
}

func main() {
	var config mainConfig
	readConfigFile(&config)

	service := cdc.Service{
		Octo: octopus.New(
			octopus_http.New(
				&http.Client{
					Timeout: 10 * time.Second,
				},
				config.Octopus.InstanceURL,
				config.Octopus.Space,
				config.Octopus.APIKey,
			),
		),
	}

	fmt.Println("Current CDC Install/Replication status:")

	offline := make(chan string)

	go func() {
		nucs, _ := service.CheckOfflineNUCs(config.Octopus.CDCProjects...)
		for _, n := range nucs {
			offline <- n
		}
		close(offline)
	}()

	printOfflineNUCs(offline)
}

func readConfigFile(cfg *mainConfig) {
	usr, _ := user.Current()
	configFile := fmt.Sprintf("%s/.monitoring/cdc_status.hcl", usr.HomeDir)

	err := hclsimple.DecodeFile(configFile, nil, cfg)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
}

func printOfflineNUCs(nucsChan <-chan string) {
	first := <-nucsChan

	if first == "" {
		// unset value means the channel closed without sending data
		fmt.Println("  - NUCs offline this morning: none")
	} else {
		fmt.Println("  - NUCs offline this morning:")
		fmt.Printf("    - %s\n", first)
	}

	for n := range nucsChan {
		fmt.Printf("    - %s\n", n)
	}
}

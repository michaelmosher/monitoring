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

	"github.com/michaelmosher/monitoring/pkg/metricly"
	metricly_http "github.com/michaelmosher/monitoring/pkg/metricly/http"
	"github.com/michaelmosher/monitoring/pkg/octopus"
	octopus_http "github.com/michaelmosher/monitoring/pkg/octopus/http"
)

type octopusCredentials struct {
	Label       string `hcl:",label"`
	InstanceURL string `hcl:"instanceURL"`
	APIKey      string `hcl:"apiKey"`
	Space       string `hcl:"space"`
}

type octopusConfig struct {
	Credentials []octopusCredentials `hcl:"credentials,block"`
	CDCProjects []string             `hcl:"cdcProjects"`
	Extra       hcl.Body             `hcl:",remain"`
}

type metriclyConfig struct {
	Username string `hcl:"Username"`
	Password string `hcl:"Password"`
}

type mainConfig struct {
	Metricly metriclyConfig `hcl:"Metricly,block"`
	Octopus  octopusConfig  `hcl:"Octopus,block"`
}

type unhealthyResult struct {
	name     string
	duration float64
}

func main() {
	var config mainConfig
	readConfigFile(&config)

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	service := cdc.Service{
		Metricly: metricly.New(
			metricly_http.Service{
				HTTPClient: httpClient,
				Username:   config.Metricly.Username,
				Password:   config.Metricly.Password,
			},
		),
	}

	var asiOcto, aosOcto octopus.Service

	for _, block := range config.Octopus.Credentials {
		if block.Label == "ASI" {
			asiOcto = octopus.New(
				octopus_http.New(httpClient, block.InstanceURL, block.Space, block.APIKey),
			)
		}

		if block.Label == "AOS" {
			aosOcto = octopus.New(
				octopus_http.New(httpClient, block.InstanceURL, block.Space, block.APIKey),
			)
		}
	}

	fmt.Println("Current CDC Install/Replication status:")

	offline := make(chan unhealthyResult)

	go func() {
		defer close(offline)

		results, _ := service.CheckOfflineNUCs(asiOcto, config.Octopus.CDCProjects...)
		for name, duration := range results {
			offline <- unhealthyResult{name, duration}
		}
	}()

	idle := make(chan string)

	go func() {
		defer close(idle)

		machines, _ := service.CheckIdleMachines(asiOcto, config.Octopus.CDCProjects...)
		for _, n := range machines {
			idle <- n
		}
	}()

	aosIdle := make(chan string)

	go func() {
		defer close(aosIdle)

		machines, _ := service.CheckIdleMachines(aosOcto, config.Octopus.CDCProjects...)
		for _, n := range machines {
			aosIdle <- n
		}
	}()

	printOfflineNUCs(offline)
	printIdleASIMachines(idle)
	printIdleAOSMachines(aosIdle)
}

func readConfigFile(cfg *mainConfig) {
	usr, _ := user.Current()
	configFile := fmt.Sprintf("%s/.monitoring/cdc_status.hcl", usr.HomeDir)

	err := hclsimple.DecodeFile(configFile, nil, cfg)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
}

func printOfflineNUCs(nucsChan <-chan unhealthyResult) {
	offlineStringChan := make(chan string)

	go func() {
		defer close(offlineStringChan)
		for ir := range nucsChan {
			offlineStringChan <- fmt.Sprintf("%s (offline for %.1f hours)", ir.name, ir.duration)
		}
	}()

	printFromChannel("NUCs offline this morning", offlineStringChan)
}

func printIdleASIMachines(idleChan <-chan string) {
	printFromChannel("NUCs or VMs Online but not replicating", idleChan)
}

func printIdleAOSMachines(idleChan <-chan string) {
	printFromChannel("AOS Systems not replicating", idleChan)
}

func printFromChannel(summary string, channel <-chan string) {
	first := <-channel

	fmt.Printf("  - %s:", summary)

	if first == "" {
		// unset value means the channel closed without sending data
		fmt.Println(" none")
	} else {
		fmt.Printf("\n    - %s\n", first)
	}

	for n := range channel {
		fmt.Printf("    - %s\n", n)
	}
}

package main

import (
	"fmt"
	"log"
	"os/user"

	"github.com/hashicorp/hcl/v2/hclsimple"

	"github.com/michaelmosher/monitoring/pkg/cdc"
)

type octopusConfig struct {
	InstanceURL string `hcl:"instanceURL"`
	APIKey      string `hcl:"apiKey"`
	Space       string `hcl:"space"`
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

	statuses, err := cdc.CheckStatus(
		config.Metricly.Username,
		config.Metricly.Password,
		config.Octopus.InstanceURL,
		config.Octopus.Space,
		config.Octopus.APIKey,
	)

	if err != nil {
		log.Fatalf("cdc.CheckStatus error: %s", err)
	}

	fmt.Println("Current CDC Install/Replication status:")

	printStatus(offlineStatus, statuses)
	printStatus(notReplicatingStatus, statuses)
	printStatus(aosStatus, statuses)
}

func readConfigFile(cfg *mainConfig) {
	usr, _ := user.Current()
	configFile := fmt.Sprintf("%s/.monitoring/cdc_status.hcl", usr.HomeDir)

	err := hclsimple.DecodeFile(configFile, nil, cfg)
	if err != nil {
		log.Fatalf("Failed to load configuration: %s", err)
	}
}

func printStatus(statusFn func(map[string]cdc.Status) bool, statuses map[string]cdc.Status) {
	if ok := statusFn(statuses); ok {
		// if everything is good, print that nothing is bad.
		fmt.Println("none")
		return
	}

	fmt.Println("")
}

func offlineStatus(statuses map[string]cdc.Status) bool {
	fmt.Print("  - NUCs offline this morning: ")

	allGood := true
	for _, status := range statuses {
		if !status.Online && !status.AOS {
			allGood = false
			fmt.Printf("\n    - %s", status.Name)
		}
	}

	return allGood
}

func notReplicatingStatus(statuses map[string]cdc.Status) bool {
	fmt.Print("  - NUCs or VMs Online but not replicating: ")

	allGood := true
	for _, status := range statuses {
		onlineButNotReplicating := status.Online && !status.Replicating
		if onlineButNotReplicating && !status.AOS {
			allGood = false
			fmt.Printf("\n    - %s", status.Name)
		}
	}

	return allGood
}

func aosStatus(statuses map[string]cdc.Status) bool {
	fmt.Print("  - AOS Systems not replicating (name lookup coming soon): ")

	allGood := true
	for uaid, status := range statuses {
		if !status.Replicating && status.AOS {
			allGood = false
			fmt.Printf("\n    - %s", uaid)
		}
	}

	return allGood
}

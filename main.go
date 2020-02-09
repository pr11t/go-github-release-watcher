package main

import (
	"flag"
	"fmt"
	"github.com/pr11t/go-github-release-watcher/internal/cmd"
	"time"
)

func main() {
	confPtr := flag.String("config", "configs/config.yml", "Configuration file path")
	flag.Parse()
	conf := cmd.LoadConfig(*confPtr)
	fmt.Println("Started program")
	for {
		cmd.GetLatestReleases(*conf)
		fmt.Printf("Sleeping for %d\n", conf.WaitBetweenChecks)
		time.Sleep(time.Duration(conf.WaitBetweenChecks) * time.Second)
	}
}

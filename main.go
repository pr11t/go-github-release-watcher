package main

import (
	"fmt"
	"github.com/pr11t/go-github-release-watcher/internal/cmd"
	"time"
)

func main() {
	conf := cmd.LoadConfig("configs/config.yml")
	fmt.Println("Started program")
	for {
		cmd.GetLatestReleases(*conf)
		fmt.Printf("Sleeping for %d\n", conf.WaitBetweenChecks)
		time.Sleep(time.Duration(conf.WaitBetweenChecks) * time.Second)
	}
}

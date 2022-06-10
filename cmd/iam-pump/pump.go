// iam-pump is a data collection service. It is responsible for transferring
// the authorization audit logs in the queue (redis list) to the storage (mongo).
package main

import (
	"math/rand"
	"time"

	"github.com/che-kwas/iam-kit/config"
	"github.com/spf13/pflag"

	"iam-pump/internal/pumpserver"
)

var (
	name = "iam-pump"
	cfg  = pflag.StringP("config", "c", "", "config file")
	help = pflag.BoolP("help", "h", false, "show help message")
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	// parse flag
	pflag.Parse()
	if *help {
		pflag.Usage()
		return
	}

	if err := config.LoadConfig(*cfg, name); err != nil {
		panic(err)
	}

	pumpserver.NewServer(name).Run()
}

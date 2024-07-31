package cmd

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/mmcloughlin/professor"
	"EncrypteDL/BlockIRC/pkg"

)

func main() {
	var (
		version    bool
		debug      bool
		configfile string
	)

	flag.BoolVar(&version, "v", false, "display version information")
	flag.BoolVar(&debug, "d", false, "enable debug logging")
	flag.StringVar(&configfile, "c", "ircd.yml", "config file")
	flag.Parse()

	if version {
		fmt.Printf(pkg.FullVersion())
		os.Exit(0)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	if debug {
		go professor.Launch(":6060")
	}

	config, err := pkg.LoadConfig(configfile)
	if err != nil {
		log.Fatal("Config file did not load successfully:", err.Error())
	}

	pkg.NewServer(config).Run()
}

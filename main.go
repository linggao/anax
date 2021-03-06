package main

import (
	"flag"
	"github.com/boltdb/bolt"
	"github.com/golang/glog"
	"github.com/open-horizon/anax/agreement"
	"github.com/open-horizon/anax/agreementbot"
	"github.com/open-horizon/anax/api"
	"github.com/open-horizon/anax/config"
	"github.com/open-horizon/anax/container"
	"github.com/open-horizon/anax/ethblockchain"
	"github.com/open-horizon/anax/exchange"
	"github.com/open-horizon/anax/governance"
	"github.com/open-horizon/anax/policy"
	"github.com/open-horizon/anax/torrent"
	"github.com/open-horizon/anax/worker"
	"os"
	"os/signal"
	"path"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

// The core of anax is an event handling system that distributes events to workers, where the workers
// process events that they are about. However, to get started, anax needs to do a bunch of initialization
// tasks. The config file has to be read in, the databases have to get created, and then the eventing system
// and the workers can be fired up.
//
func main() {
	configFile := flag.String("config", "/etc/colonus/anax.config", "Config file location")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		glog.V(2).Infof("Started CPU profiling. Writing to: %v", f.Name())
	}

	cfg, err := config.Read(*configFile)
	if err != nil {
		panic(err)
	}
	glog.V(2).Infof("Using config: %v", cfg)
	glog.V(2).Infof("GOMAXPROCS: %v", runtime.GOMAXPROCS(-1))

	// open edge DB if necessary
	var db *bolt.DB
	if len(cfg.Edge.DBPath) != 0 {
		if err := os.MkdirAll(cfg.Edge.DBPath, 0700); err != nil {
			panic(err)
		}

		edgeDB, err := bolt.Open(path.Join(cfg.Edge.DBPath, "anax.db"), 0600, &bolt.Options{Timeout: 10 * time.Second})
		if err != nil {
			panic(err)
		}
		db = edgeDB

	}

	// open Agreement Bot DB if necessary
	var agbotdb *bolt.DB
	if len(cfg.AgreementBot.DBPath) != 0 {
		if err := os.MkdirAll(cfg.AgreementBot.DBPath, 0700); err != nil {
			panic(err)
		}

		agdb, err := bolt.Open(path.Join(cfg.AgreementBot.DBPath, "agreementbot.db"), 0600, &bolt.Options{Timeout: 10 * time.Second})
		if err != nil {
			panic(err)
		}
		agbotdb = agdb
	}

	// start control signal handler
	control := make(chan os.Signal, 1)
	signal.Notify(control, os.Interrupt)
	signal.Notify(control, syscall.SIGTERM)

	// This routine does not need to be a subworker because it has no parent worker and it will terminate on its own
	// when the main anax process terminates.
	go func() {
		<-control
		glog.Infof("Closing up shop.")

		pprof.StopCPUProfile()
		if db != nil {
			db.Close()
		}
		if agbotdb != nil {
			agbotdb.Close()
		}

		os.Exit(0)
	}()

	// Get the device side policy manager started early so that all the workers can use it.
	// Make sure the policy directory is in place.
	var pm *policy.PolicyManager
	if cfg.Edge.PolicyPath == "" {
		// nothing to initialize
	} else if err := os.MkdirAll(cfg.Edge.PolicyPath, 0644); err != nil {
		glog.Errorf("Cannot create edge policy file path %v, terminating.", cfg.Edge.PolicyPath)
		panic(err)
	} else if policyManager, err := policy.Initialize(cfg.Edge.PolicyPath, cfg.ArchSynonyms, nil, true); err != nil {
		glog.Errorf("Unable to initialize policy manager, terminating.")
		panic(err)
	} else {
		pm = policyManager
	}

	// start workers
	workers := worker.NewMessageHandlerRegistry()

	workers.Add(agreementbot.NewAgreementBotWorker("AgBot", cfg, agbotdb))
	if cfg.AgreementBot.APIListen != "" {
		workers.Add(agreementbot.NewAPIListener("AgBot API", cfg, agbotdb))
	}
	workers.Add(ethblockchain.NewEthBlockchainWorker("Blockchain", cfg))

	if db != nil {
		workers.Add(api.NewAPIListener("API", cfg, db, pm))
		workers.Add(agreement.NewAgreementWorker("Agreement", cfg, db, pm))
		workers.Add(governance.NewGovernanceWorker("Governance", cfg, db, pm))
		workers.Add(exchange.NewExchangeMessageWorker("Exchange", cfg, db))
		workers.Add(container.NewContainerWorker("Container", cfg, db))
		workers.Add(torrent.NewTorrentWorker("Torrent", cfg, db))
	} else {
		workers.Add(container.NewContainerWorker("Container", cfg, agbotdb))
		workers.Add(torrent.NewTorrentWorker("Torrent", cfg, agbotdb))
	}

	// Get into the event processing loop until anax shuts itself down.
	workers.ProcessEventMessages()

	if db != nil {
		db.Close()
	}
	if agbotdb != nil {
		agbotdb.Close()
	}

	glog.Info("Main process terminating")

}

package main

import (
	"context"
	"flag"
	"github.com/jamiealquiza/envy"
	"github.com/simpicapp/simpic/internal"
	"github.com/simpicapp/simpic/internal/api"
	"github.com/simpicapp/simpic/internal/processing"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	dataDir = flag.String("path", "data", "the path to store data in")

	db    *internal.Database
	sm    *internal.SessionManager
	pr    *processing.Processor
	store internal.DiskStore
	srv   api.Server
	wg    = &sync.WaitGroup{}
)

func main() {
	envy.Parse("SIMPIC")
	flag.Parse()

	log.Printf("%s has started\n", internal.GetVersionString())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	makeDatabase()
	makeStore()
	makeSessionManager()
	makeProcessor()
	makeServer()

	log.Println("Checking for photos needing migration...")

	pr.MigrateAll()

	startServer()
	startPruningSessions()

	log.Println("Server is up.")

	<-c

	log.Println("Signal received, stopping")
	shutdown()
}

func makeDatabase() {
	var err error
	db, err = internal.OpenDatabase()
	if err != nil {
		log.Panicf("unable to connect to database: %v\n", err)
	}
}

func makeStore() {
	store = internal.DiskStore{Path: *dataDir}
}

func makeSessionManager() {
	sm = internal.NewSessionManager(db)
}

func makeProcessor() {
	pr = processing.NewProcessor(db, store, 220, 2160)
}

func makeServer() {
	userManager := internal.NewUserManager(db)
	userManager.CreateAdmin()

	srv = api.NewServer(
		db,
		userManager,
		store,
		pr)
}

func startServer() {
	wg.Add(1)
	go func() {
		if err := srv.Start(); err != nil {
			log.Panicf("Unable to start http server: %v\n", err)
		}
		wg.Done()
	}()
}

func startPruningSessions() {
	wg.Add(1)
	go func() {
		sm.PeriodicallyPruneInactiveSessions()
		wg.Done()
	}()
}

func shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	go func() {
		if err := srv.Stop(ctx); err != nil {
			log.Printf("Error shutting down http server: %v\n", err)
		}
	}()

	go sm.Stop()

	wg.Wait()
	log.Println("Simpic is stopping")
	os.Exit(0)
}

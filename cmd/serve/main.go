package main

import (
	"context"
	"flag"
	"github.com/jamiealquiza/envy"
	"github.com/simpicapp/simpic/internal"
	"github.com/simpicapp/simpic/internal/api"
	"github.com/simpicapp/simpic/internal/storage"
	"log"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
	"time"
)

var (
	dataDir = flag.String("path", "data", "the path to store data in")

	db  *internal.Database
	sm  *internal.SessionManager
	srv api.Server
	wg  = &sync.WaitGroup{}
)

func main() {
	envy.Parse("SIMPIC")
	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	makeDatabase()
	makeSessionManager()
	makeServer()

	startServer()
	startPruningSessions()

	log.Println("Simpic has started")

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

func makeSessionManager() {
	sm = internal.NewSessionManager(db)
}

func makeServer() {
	userManager := internal.NewUserManager(db)
	userManager.CreateAdmin()

	driver := storage.DiskDriver{Path: *dataDir}

	thumbnailer := internal.NewThumbnailer(driver, storage.DiskDriver{Path: path.Join(*dataDir, "thumbnails")}, 220)

	srv = api.NewServer(
		db,
		thumbnailer,
		userManager,
		driver,
		internal.NewStorer(db, driver))
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

package main

import (
	"context"
	"flag"
	"github.com/jamiealquiza/envy"
	"github.com/simpicapp/simpic"
	"github.com/simpicapp/simpic/http"
	"github.com/simpicapp/simpic/storage"
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
)

func main() {
	envy.Parse("SIMPIC")
	flag.Parse()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	server := makeServer()

	go func() {
		if err := server.Start(); err != nil {
			log.Panicf("Unable to start http server: %v\n", err)
		}
	}()

	log.Println("Simpic has started")

	<-c

	log.Println("Signal received, stopping")
	shutdown(server)
}

func makeServer() http.Server {
	db, err := simpic.OpenDatabase()
	if err != nil {
		log.Panicf("unable to connect to database: %v\n", err)
	}

	userManager := simpic.NewUserManager(db)
	userManager.CreateAdmin()

	driver := storage.DiskDriver{Path: *dataDir}

	thumbnailer := simpic.NewThumbnailer(driver, storage.DiskDriver{Path: path.Join(*dataDir, "thumbnails")}, 220)

	return http.NewServer(
		db,
		thumbnailer,
		userManager,
		driver,
		simpic.NewStorer(db, driver))
}

func shutdown(srv http.Server) {
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	wg.Add(1)
	go func() {
		if err := srv.Stop(ctx); err != nil {
			log.Printf("Error shutting down http server: %v\n", err)
		}
		wg.Done()
	}()

	wg.Wait()
	log.Println("Simpic is stopping")
	os.Exit(0)
}

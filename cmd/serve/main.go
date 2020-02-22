package main

import (
	"flag"
	"github.com/jamiealquiza/envy"
	"github.com/simpicapp/simpic"
	"github.com/simpicapp/simpic/http"
	"github.com/simpicapp/simpic/storage"
	"log"
	"path"
)

var (
	port        = flag.Int("port", 8080, "the port to listen on")
	dataDir     = flag.String("path", "data", "the path to store data in")
	frontendDir = flag.String("frontend", "dist", "the path to serve frontend files from")
)

func main() {
	envy.Parse("SIMPIC")
	flag.Parse()

	db, err := simpic.OpenDatabase()
	if err != nil {
		log.Fatalf("unable to connect to database: %v\n", err)
		return
	}

	userManager := simpic.NewUserManager(db)
	userManager.CreateAdmin()

	driver := storage.DiskDriver{Path: *dataDir}

	thumbnailer := simpic.NewThumbnailer(driver, storage.DiskDriver{Path: path.Join(*dataDir, "thumbnails")}, 220)

	http.Start(
		db,
		thumbnailer,
		userManager,
		driver,
		simpic.NewStorer(db, driver),
		*frontendDir,
		*port)
}

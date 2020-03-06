package api

import (
	"fmt"
	"github.com/simpicapp/simpic/internal"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

var (
	htmlContent    = ""
	htmlLoaderOnce sync.Once
)

func (s *server) handleFrontendPath(title string) http.HandlerFunc {
	htmlLoaderOnce.Do(func() {
		b, err := ioutil.ReadFile(filepath.Join(*frontendDir, "index.html"))
		if err != nil {
			log.Panicf("Unable to read frontend index.html file: %v", err)
		}
		htmlContent = string(b)
	})

	getPhoto := func(r *http.Request) *internal.Photo {
		id := uuidUrlParam(r, "photoId")
		if id != nil {
			photo, _ := s.db.GetPhoto(*id, visForAccess(r))
			return photo
		}
		return nil
	}

	getAlbum := func(r *http.Request) *internal.Album {
		id := uuidUrlParam(r, "albumId")
		if id != nil {
			album, _ := s.db.GetAlbum(*id, visForAccess(r))
			return album
		}
		return nil
	}

	reverse := func(s []string) []string {
		for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
			s[i], s[j] = s[j], s[i]
		}
		return s
	}

	return func(w http.ResponseWriter, r *http.Request) {
		titleParts := []string{"", title}
		if album := getAlbum(r); album != nil {
			titleParts = append(titleParts, album.Name)
		}

		if photo := getPhoto(r); photo != nil {
			titleParts = append(titleParts, photo.FileName)
		}

		content := fmt.Sprintf("<title>%s", html.EscapeString(strings.Join(reverse(titleParts), " - ")))
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(strings.Replace(htmlContent, "<title>", content, 1)))
	}
}

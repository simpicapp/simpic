package api

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
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

	buildOpenGraphImage := func(id uuid.UUID) string {
		return fmt.Sprintf(`
				<meta property="og:image" content="/data/%[1]s/2.webp">
				<meta property="og:image:type" content="image/webp">
				<meta property="og:image" content="/data/%[1]s/2.jpeg">
				<meta property="og:image:type" content="image/jpeg">
		`, id)
	}

	buildOpenGraphGeneral := func(title string) string {
		return fmt.Sprintf(`
				<meta property="og:site_name" content="Simpic">
				<meta property="og:type" content="article">
				<meta property="og:title" content="%[1]s">
		`, html.EscapeString(title))
	}

	buildOpenGraphTags := func(photo *internal.Photo, album *internal.Album) string {
		if photo != nil {
			return buildOpenGraphImage(photo.Id) + buildOpenGraphGeneral(photo.FileName)
		}
		if album != nil {
			content := buildOpenGraphGeneral(album.Name)
			if album.Cover != nil {
				content += buildOpenGraphImage(*album.Cover)
			}
			return content
		}
		return ""
	}

	return func(w http.ResponseWriter, r *http.Request) {
		titleParts := []string{"", title}

		album := getAlbum(r)
		if album != nil {
			titleParts = append(titleParts, album.Name)
		}

		photo := getPhoto(r)
		if photo != nil {
			titleParts = append(titleParts, photo.FileName)
		}

		content := fmt.Sprintf("%s<title>%s", buildOpenGraphTags(photo, album), html.EscapeString(strings.Join(reverse(titleParts), " - ")))
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(strings.Replace(htmlContent, "<title>", content, 1)))
	}
}

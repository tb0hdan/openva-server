package library

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/dhowden/tag"
	"github.com/tb0hdan/openva-server/api"
)

// https://stackoverflow.com/questions/20401873/remove-invalid-utf-8-characters-from-a-string-go-lang
func fixUTF(r rune) rune {
	if r == utf8.RuneError {
		return -1
	}
	return r
}

type LocalLibrary struct {
	MusicDir          string
	HTTPServerAddress string
}

func (l *LocalLibrary) Library(criteria, token, serverIP string) (libraryItems []*api.LibraryItem, err error) {
	dir, err := filepath.EvalSymlinks(l.MusicDir)
	if err != nil {
		log.Fatal(err)
	}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".mp3") {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}

		artist := ""
		album := ""
		track := ""

		m, err := tag.ReadFrom(file)
		if err != nil {
			log.Println(path, err)

		} else {
			artist = strings.Map(fixUTF, m.Artist())
			album = strings.Map(fixUTF, m.Album())
			track = strings.Map(fixUTF, m.Title())
		}

		if !libraryFilterPassed(criteria, artist, album, track, pathWords(path)) {
			return nil
		}

		escapedPath := ""
		for _, r := range strings.Split(strings.TrimPrefix(path, dir), "/") {
			escapedPath += "/" + url.PathEscape(r)
		}

		if strings.HasPrefix(escapedPath, "//") {
			escapedPath = strings.TrimPrefix(escapedPath, "/")
		}

		item := &api.LibraryItem{
			URL:    fmt.Sprintf("http://%s", serverIP) + l.HTTPServerAddress + "/music" + escapedPath + fmt.Sprintf("?token=%s", token),
			Artist: artist,
			Album:  album,
			Track:  track,
		}
		libraryItems = append(libraryItems, item)

		return nil
	})

	return
}

func libraryFilterPassed(criteria string, args ...string) bool { // nolint gocyclo
	var (
		artist string
		//album        string
		track        string
		searchArtist string
		searchTrack  string
	)
	if criteria == "" {
		return true
	}

	if len(args) == 0 {
		return true
	}
	criteria = strings.ToLower(criteria)

	if len(strings.Split(criteria, " - ")) >= 2 {
		searchArtist = strings.TrimSpace(strings.Split(criteria, " - ")[0])
		searchTrack = strings.TrimSpace(strings.Split(criteria, " - ")[1])
	}

	// artist, track
	if len(args) > 3 {
		artist = strings.TrimSpace(args[0])
		track = strings.TrimSpace(args[2])
	}

	for _, arg := range args {
		arg = strings.ToLower(arg)
		if len(arg) > 0 && strings.Contains(arg, criteria) {
			return true
		}
		// Special case: Artist Name - Track Name
		if searchArtist == "" || searchTrack == "" {
			continue
		}
		if artist == "" || track == "" {
			continue
		}
		if strings.EqualFold(searchArtist, artist) {
			continue
		}
		if strings.EqualFold(searchTrack, track) {
			return true

		}
	}
	return false
}

func pathWords(path string) (newString string) {
	re := regexp.MustCompile(`[/|_|-|-|(|)|\.]`)
	for _, str := range strings.Split(re.ReplaceAllString(path, " "), " ") {
		if strings.TrimSpace(str) == "" {
			continue
		}
		newString += " " + str
	}
	return
}

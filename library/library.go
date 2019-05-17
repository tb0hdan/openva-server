package library

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
	"unsafe"

	"github.com/dhowden/tag"
	"github.com/dustin/go-humanize"
	log "github.com/sirupsen/logrus"
	"github.com/tb0hdan/openva-server/api"
)

// https://stackoverflow.com/questions/20401873/remove-invalid-utf-8-characters-from-a-string-go-lang
func fixUTF(r rune) rune {
	if r == utf8.RuneError {
		return -1
	}
	return r
}

type File struct {
	Artist,
	Album,
	Track,
	EvaluatedDirectory,
	FilePath string
	Length int32
}

type Library struct {
	MusicDir          string
	HTTPServerAddress string
	Files             []File
}

func (l *Library) UpdateIndex() {
	files, err := l.IndexFiles()
	if err != nil {
		log.Error("library index failed")
		return
	}
	l.Files = files
	log.Println("Size of library in memory: ", humanize.Bytes(uint64(int(unsafe.Sizeof(l.Files))*len(l.Files))))
}

func (l *Library) IndexFiles() (libraryFiles []File, err error) {
	dir, err := filepath.EvalSymlinks(l.MusicDir)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
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
			log.Printf("%+v\n", err)
			return err
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

		libraryFiles = append(libraryFiles, File{
			Artist:             artist,
			Album:              album,
			Track:              track,
			Length:             0,
			FilePath:           path,
			EvaluatedDirectory: dir,
		})
		return nil
	})
	return libraryFiles, err
}

func (l *Library) Library(criteria, token, serverIP string) (libraryItems []*api.LibraryItem, err error) {

	for _, libraryFile := range l.Files {

		artist, album, track, path, dir := libraryFile.Artist,
			libraryFile.Album, libraryFile.Track, libraryFile.FilePath, libraryFile.EvaluatedDirectory

		if !libraryFilterPassed(criteria, artist, album, track, pathWords(path)) {
			continue
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

	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(libraryItems), func(i, j int) {
		libraryItems[i], libraryItems[j] = libraryItems[j], libraryItems[i]
	})

	return libraryItems, err
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

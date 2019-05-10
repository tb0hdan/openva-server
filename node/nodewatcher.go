package node

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/shkh/lastfm-go/lastfm"
)

const DeadInterval = 15

type NodeWatcher struct {
	LastFMAPIKey    string
	LastFMAPISecret string
	LastFMUsername  string
	LastFMPassword  string
	lastFMLoggedIn  bool
	lastFMAPI       *lastfm.Api
}

func NewNodeWatcher(lastFMAPIKey, lastFMAPISecret, lastFMUsername, lastFMPassword string) *NodeWatcher {
	return &NodeWatcher{
		LastFMAPIKey:    lastFMAPIKey,
		LastFMAPISecret: lastFMAPISecret,
		LastFMUsername:  lastFMUsername,
		LastFMPassword:  lastFMPassword,
	}
}

func NewNodeWatcherWithConnection(lastFMAPIKey, lastFMAPISecret, lastFMUsername, lastFMPassword string) *NodeWatcher {
	nodeWatcher := NewNodeWatcher(lastFMAPIKey, lastFMAPISecret, lastFMUsername, lastFMPassword)
	nodeWatcher.LoginLastFm()
	return nodeWatcher
}

func (nw *NodeWatcher) LoginLastFm() {
	lfmAPI := lastfm.New(nw.LastFMAPIKey, nw.LastFMAPISecret)

	err := lfmAPI.Login(nw.LastFMUsername, nw.LastFMPassword)
	if err != nil {
		fmt.Println(err)
		return
	}
	nw.lastFMLoggedIn = true
	nw.lastFMAPI = lfmAPI
}

func (nw *NodeWatcher) Scrobble(nowPlaying string) {
	var (
		artist,
		track string
	)
	if !nw.lastFMLoggedIn {
		log.Debug("Scrobble exit, not logged in")
		return
	}
	if len(strings.Split(nowPlaying, " - ")) < 2 {
		artist = nowPlaying
		track = nowPlaying
	} else {
		artist = strings.Split(nowPlaying, " - ")[0]
		track = strings.Split(nowPlaying, " - ")[1]
	}
	p := lastfm.P{"artist": artist, "track": track}
	start := time.Now().Unix()

	p["timestamp"] = start
	_, err := nw.lastFMAPI.Track.Scrobble(p)
	if err != nil {
		log.Debug(err)
		return
	}

}

func (nw *NodeWatcher) SetNowPlaying(nowPlaying string) {
	var (
		artist,
		track string
	)
	if !nw.lastFMLoggedIn {
		log.Debug("SetNowPlaying exit, not logged in")
		return
	}
	if len(strings.Split(nowPlaying, " - ")) < 2 {
		artist = nowPlaying
		track = nowPlaying
	} else {
		artist = strings.Split(nowPlaying, " - ")[0]
		track = strings.Split(nowPlaying, " - ")[1]
	}
	p := lastfm.P{"artist": artist, "track": track}
	_, err := nw.lastFMAPI.Track.UpdateNowPlaying(p)
	if err != nil {
		log.Debug(err)
		return
	}
}

func (nw *NodeWatcher) ClientWatcher(ctx context.Context, key string, nodeState *NodeStateType) {
	var (
		passed        int
		previousTrack string
		trackChanged  bool
	)

	if tmpState, ok := nodeState.Get(key); ok {
		previousTrack = tmpState.PlayerState.NowPlaying
	}

	if len(previousTrack) > 0 {
		nw.SetNowPlaying(previousTrack)
	}

	for {
		select {
		case <-ctx.Done():
			break
		default:
		}

		state, ok := nodeState.Get(key)

		if !ok {
			time.Sleep(time.Second)
			continue
		}

		// FIXME: make this duration instead of int64
		if time.Now().Unix()-state.LastUpdatedTS >= DeadInterval {
			log.Println(state.SystemInformation.SystemUUID, " is dead, removing...")
			nodeState.Delete(key)
			continue
		}

		if state.PlayerState.NowPlaying != previousTrack && len(state.PlayerState.NowPlaying) > 0 {
			previousTrack = state.PlayerState.NowPlaying
			passed = 0
			trackChanged = true
			// track change
			nw.SetNowPlaying(previousTrack)
		}
		log.Println(key, state.PlayerState.NowPlaying)
		if passed >= 35 && len(state.PlayerState.NowPlaying) > 0 && trackChanged {
			log.Println("Can scrobble: ", state.PlayerState.NowPlaying)
			nw.Scrobble(state.PlayerState.NowPlaying)
			trackChanged = false
		}

		passed++
		time.Sleep(time.Second)
	}
}

func (nw *NodeWatcher) ServerWatcher(nodeStates *NodeStateType) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	oldmap := make(map[string]*NodeState)
	for {
		time.Sleep(time.Second)
		if nodeStates.Len() == 0 {
			log.Println("Waiting for first client...")
			continue
		}
		for key := range nodeStates.All() {
			if value, ok := nodeStates.Get(key); ok && oldmap[key] == nil {
				oldmap[key] = value
				log.Println("Got new client ", key)
				go nw.ClientWatcher(ctx, key, nodeStates)
			}
		}
	}
}

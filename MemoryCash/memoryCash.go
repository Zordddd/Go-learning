package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func httpGetBody(url string) (interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	return io.ReadAll(resp.Body)
}

type Func[K comparable, T any] func(K) (T, error)
type result[V any] struct {
	value V
	err   error
}

type entry[V any] struct {
	ready chan struct{}
	res   result[V]
}
type Shard[K comparable, V any] struct {
	f     Func[K, V]
	mu    sync.Mutex
	cache map[K]*entry[V]
}
type Memo[K comparable, V any] struct {
	shards    []*Shard[K, V]
	numShards int
}

func NewMemo[K comparable, V any](f Func[K, V], numShards int) *Memo[K, V] {
	m := &Memo[K, V]{shards: make([]*Shard[K, V], 0, numShards), numShards: numShards}
	for i := 0; i < numShards; i++ {
		m.shards = append(m.shards, newShard[K, V](f))
	}
	return m
}

func (m *Memo[K, V]) Get(url K) (interface{}, error) {
	hashURL, err := hasher(url)
	if err != nil {
		return nil, err
	}
	shardID := hashURL % m.numShards
	res, err := m.shards[shardID].Get(url)
	return res, err
}

func newShard[K comparable, V any](f Func[K, V]) *Shard[K, V] {
	return &Shard[K, V]{f: f, cache: make(map[K]*entry[V])}
}

func (s *Shard[K, V]) Get(url K) (V, error) {
	s.mu.Lock()
	e := s.cache[url]
	if e == nil {
		e = &entry[V]{ready: make(chan struct{})}
		s.cache[url] = e
		s.mu.Unlock()

		e.res.value, e.res.err = s.f(url)
		close(e.ready)
	} else {
		s.mu.Unlock()
		<-e.ready
	}

	return e.res.value, e.res.err
}

func hasher[K comparable](url K) (int, error) {
	h := fnv.New32a()
	urlBytes, err := json.Marshal(url)
	if err != nil {
		return 0, err
	}
	_, _ = h.Write(urlBytes)
	return int(h.Sum32()), nil
}

func main() {
	sites := []string{
		"https://www.google.com",
		"https://www.youtube.com",
		"https://www.facebook.com",
		"https://www.twitter.com",
		"https://www.instagram.com",
		"https://www.linkedin.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
		"https://www.wikipedia.org",
		"https://www.amazon.com",
		"https://www.netflix.com",
		"https://www.instagram.com",
		"https://www.apple.com",
		"https://www.adobe.com",
		"https://www.cloudflare.com",
		"https://www.digitalocean.com",
		"https://www.heroku.com",
		"https://www.gitlab.com",
		"https://www.wikipedia.org",
		"https://www.docker.com",
		"https://www.kubernetes.io",
		"https://www.python.org",
		"https://www.golang.org",
		"https://www.nodejs.org",
		"https://www.adobe.com",
		"https://www.vuejs.org",
		"https://www.angular.io",
		"https://www.jquery.com",
		"https://www.bootstrap.com",
		"https://www.tailwindcss.com",
		"https://www.sass-lang.com",
		"https://www.lesscss.org",
		"https://www.webpack.js.org",
		"https://www.babeljs.io",
		"https://www.typescriptlang.org",
		"https://www.coffeescript.org",
		"https://www.d3js.org",
		"https://www.chartjs.org",
		"https://www.threejs.org",
		"https://www.unity.com",
		"https://www.unrealengine.com",
		"https://www.blender.org",
		"https://www.gimp.org",
		"https://www.inkscape.org",
		"https://www.figma.com",
		"https://www.jquery.com",
		"https://www.invisionapp.com",
		"https://www.zeplin.io",
		"https://www.abstract.com",
		"https://www.notion.so",
		"https://www.trello.com",
		"https://www.asana.com",
		"https://www.slack.com",
		"https://www.discord.com",
		"https://www.google.com",
		"https://www.youtube.com",
		"https://www.facebook.com",
		"https://www.twitter.com",
		"https://www.instagram.com",
		"https://www.linkedin.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
		"https://www.wikipedia.org",
		"https://www.amazon.com",
		"https://www.netflix.com",
		"https://www.instagram.com",
		"https://www.apple.com",
		"https://www.adobe.com",
		"https://www.cloudflare.com",
		"https://www.digitalocean.com",
		"https://www.heroku.com",
		"https://www.gitlab.com",
		"https://www.wikipedia.org",
		"https://www.docker.com",
		"https://www.kubernetes.io",
		"https://www.python.org",
		"https://www.buffer.com",
		"https://www.airtable.com",
		"https://www.adobe.com/creativecloud.html",
		"https://www.unsplash.com",
		"https://www.pexels.com",
		"https://www.pixabay.com",
		"https://www.shutterstock.com",
		"https://www.gettyimages.com",
		"https://www.istockphoto.com",
		"https://www.500px.com",
		"https://www.flickr.com",
		"https://www.deviantart.com",
		"https://www.behance.net",
		"https://www.dribbble.com",
		"https://www.artstation.com",
		"https://www.soundcloud.com",
		"https://www.spotify.com",
		"https://www.apple.com/music",
		"https://www.youtube.com/music",
		"https://www.deezer.com",
		"https://www.tidal.com",
		"https://www.pandora.com",
		"https://www.last.fm",
		"https://www.bandcamp.com",
		"https://www.vimeo.com",
		"https://www.dailymotion.com",
		"https://www.twitch.tv",
		"https://www.tiktok.com",
		"https://www.snapchat.com",
		"https://www.pinterest.com",
		"https://www.tumblr.com",
		"https://www.medium.com",
		"https://www.quora.com",
		"https://www.producthunt.com",
	}
	wg := sync.WaitGroup{}

	cacheGetBody := NewMemo[string, interface{}](httpGetBody, 10)
	for _, site := range sites {
		wg.Add(1)
		go func(site string) {
			defer wg.Done()
			start := time.Now()
			_, err := cacheGetBody.Get(site)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%s Get took %v\n", site, time.Now().Sub(start).Milliseconds())
		}(site)
	}
	wg.Wait()
}

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func httpGetBody(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	return readAllWithContext(ctx, resp.Body)
}

func readAllWithContext(ctx context.Context, r io.Reader) ([]byte, error) {
	var result []byte
	buf := make([]byte, 32*1024)

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			n, err := r.Read(buf)
			if n > 0 {
				result = append(result, buf[:n]...)
			}
			if err == io.EOF {
				return result, nil
			}
			if err != nil {
				return nil, err
			}
		}
	}
}

func zero[T any]() (res T) {
	return res
}

type Func[K comparable, T any] func(context.Context, K) (T, error)
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

func (m *Memo[K, V]) Get(ctx context.Context, url K) (res V, err error) {
	hashURL, err := hasher(url)
	if err != nil {
		return res, err
	}
	shardID := hashURL % m.numShards
	res, err = m.shards[shardID].Get(ctx, url)
	return res, err
}

func newShard[K comparable, V any](f Func[K, V]) *Shard[K, V] {
	return &Shard[K, V]{f: f, cache: make(map[K]*entry[V])}
}

func (s *Shard[K, V]) Get(ctx context.Context, url K) (V, error) {
	s.mu.Lock()
	e := s.cache[url]
	if e == nil {
		e = &entry[V]{ready: make(chan struct{})}
		s.cache[url] = e
		s.mu.Unlock()

		go func() {
			e.res.value, e.res.err = s.f(ctx, url)
			close(e.ready)

			if errors.Is(e.res.err, context.Canceled) || errors.Is(e.res.err, context.DeadlineExceeded) {
				s.mu.Lock()
				delete(s.cache, url)
				s.mu.Unlock()
			}
		}()
	} else {
		s.mu.Unlock()
	}
	select {
	case <-ctx.Done():
		return zero[V](), ctx.Err()
	case <-e.ready:
		return e.res.value, e.res.err
	}
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
	}
	wg := sync.WaitGroup{}

	cacheGetBody := NewMemo[string, []byte](httpGetBody, 10)
	for _, site := range sites {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		wg.Add(1)
		go func(ctx context.Context, site string) {
			defer func() {
				wg.Done()
				cancel()
			}()
			start := time.Now()
			_, err := cacheGetBody.Get(ctx, site)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%s Get took %v\n", site, time.Now().Sub(start).Milliseconds())
		}(ctx, site)
	}

	wg.Wait()

	// Медленная функция для теста таймаутов
	//slowFunc := func(ctx context.Context, key string) (string, error) {
	//	select {
	//	case <-time.After(2 * time.Second): // Всегда занимает 2 секунды
	//		return "result", nil
	//	case <-ctx.Done():
	//		return "", ctx.Err()
	//	}
	//}
	//
	//cache := NewMemo(slowFunc, 4)
	//
	//// Тест с коротким таймаутом (должен сработать)
	//ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	//defer cancel()
	//
	//start := time.Now()
	//result, err := cache.Get(ctx, "test")
	//elapsed := time.Since(start)
	//
	//if err != nil {
	//	fmt.Printf("✅ Таймаут сработал за %v: %v\n", elapsed, err)
	//} else {
	//	fmt.Printf("❌ Таймаут не сработал: %v\n", result)
	//}
	//
	//// Тест с длинным таймаутом (должен успешно завершиться)
	//ctx2, cancel2 := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel2()
	//
	//start = time.Now()
	//result, err = cache.Get(ctx2, "test")
	//elapsed = time.Since(start)
	//
	//if err != nil {
	//	fmt.Printf("❌ Неожиданная ошибка: %v\n", err)
	//} else {
	//	fmt.Printf("✅ Успешное завершение за %v: %v\n", elapsed, result)
	//}
}

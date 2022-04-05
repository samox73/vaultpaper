package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type redditLocation struct {
	Subreddit string
}

func NewRedditLocation(sub string) redditLocation {
	subBase := strings.TrimPrefix(sub, "r/")
	return redditLocation{Subreddit: subBase}
}

func (l *redditLocation) Scan() {}

func (l redditLocation) Print() {
	fmt.Printf("Subreddit: r/%s\n", l.Subreddit)
}

func (l redditLocation) PrintFiles(indent int) {}

func (l redditLocation) Save() {}

func (l *redditLocation) Load() {}

func (l redditLocation) GetRandomFilePath() (string, error) {
	return l.saveFileToTemporary()
}

func (l redditLocation) saveFileToTemporary() (string, error) {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	// The reddit API does not like HTTP/2
	// Per https://pkg.go.dev/net/http?utm_source=gopls#pkg-overview ,
	// I'm copying http.DefaultTransport and replacing the HTTP/2 stuff
	transport := &http.Transport{
		Proxy:       http.ProxyFromEnvironment,
		DialContext: dialer.DialContext,

		// change from default
		ForceAttemptHTTP2: false,

		MaxIdleConns:        100,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,

		// use an empty map instead of nil per the link above
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),

		ExpectContinueTimeout: 1 * time.Second,
	}

	httpClient := &http.Client{
		Timeout:   time.Second * 30,
		Transport: transport,
	}

	client, err := reddit.NewReadonlyClient(
		reddit.WithHTTPClient(httpClient),
	)
	posts, _, err := client.Subreddit.TopPosts(context.Background(), l.Subreddit, &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: 1,
		},
		Time: "today",
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received %d posts.\n", len(posts))
	return "", nil
}

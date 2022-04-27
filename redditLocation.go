package main

import (
	"context"
	"crypto/tls"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type redditLocation struct {
	Subreddit  string
	Files      []string
	client     *reddit.Client
	UrlsToFile map[string]string
}

func NewRedditLocation(sub string) redditLocation {
	subBase := strings.TrimPrefix(sub, "r/")
	l := redditLocation{Subreddit: subBase, UrlsToFile: make(map[string]string)}
	l.initClient()
	return l
}

func (l *redditLocation) Scan() {}

func (l redditLocation) Print() {
	fmt.Printf("Subreddit: r/%s\n", l.Subreddit)
}

func (l redditLocation) PrintFiles(indent int) {}

func (l redditLocation) getConfigPath() string {
	return filepath.Join(GetStoreDir(), l.Subreddit, ".vaultlocation")
}

func (l redditLocation) Save() {
	filepath := l.getConfigPath()
	dataFile, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	dataEncoder := gob.NewEncoder(dataFile)
	err = dataEncoder.Encode(l)
	if err != nil {
		log.Fatal(err)
	}
	dataFile.Close()
	fmt.Printf("Saved location to file '%s'\n", filepath)
}

func (l *redditLocation) Load() {}

func (l redditLocation) GetRandomFilePath() (string, error) {
	l.getTopPosts(5, "today")
	return "", nil
}

func (l *redditLocation) initClient() {
	dialer := &net.Dialer{
		Timeout:   3 * time.Second,
		KeepAlive: 3 * time.Second,
	}

	// reddit is not working with HTTP/2 yet https://pkg.go.dev/net/http?utm_source=gopls#pkg-overview ,
	transport := &http.Transport{
		Proxy:       http.ProxyFromEnvironment,
		DialContext: dialer.DialContext,

		// change from default
		ForceAttemptHTTP2: false,

		MaxIdleConns:        100,
		IdleConnTimeout:     3 * time.Second,
		TLSHandshakeTimeout: 1 * time.Second,

		// use an empty map instead of nil per the link above
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),

		ExpectContinueTimeout: 1 * time.Second,
	}

	httpClient := &http.Client{
		Timeout:   3 * time.Second,
		Transport: transport,
	}

	client, err := reddit.NewReadonlyClient(
		reddit.WithHTTPClient(httpClient),
	)
	if err != nil {
		log.Fatal(err)
	}
	l.client = client
}

func (l redditLocation) getTopPosts(limit int, time string) {
	fmt.Printf("Getting top posts from r/%s\n", l.Subreddit)
	if l.client == nil {
		log.Fatal("Client has not been initialized")
	}
	posts, _, err := l.client.Subreddit.TopPosts(context.Background(), l.Subreddit, &reddit.ListPostOptions{
		ListOptions: reddit.ListOptions{
			Limit: limit,
		},
		Time: time,
	})
	if err != nil {
		panic(err)
	}
	for _, post := range posts {
		urlFileName, err := validateImageURL(post.URL)
		if err != nil {
			log.Fatal(err)
		}
		filePath, err := genFilePath(l.Subreddit, post.Title, urlFileName)
		if err != nil {
			log.Fatal(err)
		}
		err = l.saveImage(post.URL, filePath)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("Received %d posts.\n", len(posts))
	fmt.Printf("Location r/%s contains:\n", l.Subreddit)
	for url, file := range l.UrlsToFile {
		fmt.Printf("  - %s => %s\n", url, file)
	}
}

func validateImageURL(fullURL string) (string, error) {
	fileURL, err := url.Parse(fullURL)
	if err != nil {
		log.Fatal(err)
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")

	urlFileName := segments[len(segments)-1]
	allowedImageExtensions := []string{".jpg", ".jpeg", ".png"}
	for _, suffix := range allowedImageExtensions {
		if strings.HasSuffix(urlFileName, suffix) {
			return urlFileName, nil
		}

	}
	return "", fmt.Errorf("urlFileName doesn't end in allowed extension: %#v , %#v\n ", urlFileName, allowedImageExtensions)
}

func genFilePath(subredditName string, title string, urlFileName string) (string, error) {
	directory := filepath.Join(GetStoreDir(), subredditName)
	os.Mkdir(directory, os.ModePerm)

	for _, s := range []string{" ", "/", "\\", "\n", "\r", "\x00"} {
		urlFileName = strings.ReplaceAll(urlFileName, s, "_")
		subredditName = strings.ReplaceAll(subredditName, s, "_")
		title = strings.ReplaceAll(title, s, "_")
	}

	fullFileName := subredditName + "_" + title + "_" + urlFileName
	filePath := filepath.Join(directory, fullFileName)

	// remove chars from title if it's too long for the OS to handle
	const maxPathLength = 250
	if len(filePath) > maxPathLength {
		toChop := len(filePath) - maxPathLength
		if toChop > len(title) {
			return "", fmt.Errorf("filePath to long and title too short: %#v\n", filePath)
		}

		title = title[:len(title)-toChop]
		fullFileName = subredditName + "_" + title + "_" + urlFileName
		filePath = filepath.Join(directory, fullFileName)
	}
	return filePath, nil
}

func (l *redditLocation) saveImage(url string, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	contentBytes := make([]byte, 512)

	_, err = response.Body.Read(contentBytes)
	if err != nil {
		return err
	}

	// https://golang.org/pkg/net/http/#DetectContentType
	contentType := http.DetectContentType(contentBytes)

	if !(contentType == "image/jpeg" || contentType == "image/png") {
		err = fmt.Errorf("contentType is not 'image/jpeg' or 'image/png': %+v\n", contentType)
		return err
	}

	_, err = file.Write(contentBytes)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	l.Files = append(l.Files, fileName)
	l.UrlsToFile[url] = fileName
	return nil
}

package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"github.com/SlyMarbo/rss"
	"io"
	"log"
	"os"
	"path"
	"sync"
)

var CacheDir string = "./data/cache/"
var cacheLock sync.Mutex

// returns the cache file path for the given feed url
func feedFile(url string) string {
	h := sha1.New()
	io.WriteString(h, url)
	encodedUrl := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return path.Join(CacheDir, encodedUrl)
}

func ReadFeed(task int, url string, resChan chan *rss.Feed) {
	log.Printf("(%03d) start fetch for %s", task, url)

	var feed *rss.Feed

	file := feedFile(url)

	// stat the file so we can tell if it exists
	_, err := os.Stat(file)
	if err != nil {
		log.Printf("(%03d) no cache file, fetching for first time", task)
		goto first_fetch
	} else {
		fd, err := os.Open(file)
		if err != nil {
			log.Printf("(%03d) could not read cache, refetching: %v", task, err)
			goto first_fetch
		} else {
			log.Printf("(%03d) reading feed cache from %v", task, file)
			dec := gob.NewDecoder(fd)
			err = dec.Decode(&feed)
			fd.Close()
			if err != nil {
				log.Printf("(%03d) invalid cache file, refetching: %v", task, err)
				goto first_fetch
			} else {
				err = feed.Update()
				if err != nil {
					log.Printf("(%03d) failed to update: %v", task, err)
                    resChan <- nil
                    return
				}
				goto loaded
			}
		}
	}

first_fetch:
	feed, err = rss.Fetch(url)
	if err != nil {
		log.Printf("(%03d) failed to fetch: %v", task, err)
        resChan <- nil
		return
	}

loaded:
	fd, err := os.Create(file)
	defer fd.Close()
	if err != nil {
		log.Println("(%03d) could not create cache file: %v", task, err)
	} else {
		enc := gob.NewEncoder(fd)
		err = enc.Encode(feed)
		if err != nil {
			log.Println("(%03d) could not write feed cache: %v", task, err)
		}
		log.Printf("(%03d) updated cache for %s", task, feed.Title)
	}

	resChan <- feed
}

// retreive all feeds in parallel
func ReadFeeds(feeds []string) []*rss.Feed {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	err := os.MkdirAll(CacheDir, 0700)
	if err != nil {
		log.Println("could not create cache dir, aborting:", err)
		return make([]*rss.Feed, 0)
	}
	log.Printf("start syncing %d feeds\n", len(feeds))

	res := make([]*rss.Feed, 0, len(feeds))
	resChan := make(chan *rss.Feed, len(feeds))
	for i, url := range feeds {
		go ReadFeed(i, url, resChan)
	}
	for i := 0; i < len(feeds); i++ {
		f := <-resChan
		if f != nil {
			res = append(res, f)
		}
	}

	log.Printf("successfully synced %d/%d feeds\n", len(res), len(feeds))
	return res
}

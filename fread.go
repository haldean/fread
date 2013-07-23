package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

var FeedFile string = "./data/feeds"
var OpmlFile string = "./data/feeds.opml"

var useOpml = flag.Bool(
	"use_opml", false,
	"if true, reads from data/feeds.opml to get feed list")

var useFeedFile = flag.Bool(
	"use_feed_file", true,
	"if true, reads from data/feeds to get feed list. "+
		"Feed list should be newline-separated URLs of XML feeds")

var showHeadlines = flag.Bool(
	"headlines", false, "if true, prints headlines from feeds after sync")

func ReadConfig() ([]string, error) {
	feedFile, err := ioutil.ReadFile(FeedFile)
	if err != nil {
		return nil, err
	}
	feeds := make([]string, 0)
	buf := bytes.NewBuffer(feedFile)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		line = strings.Trim(line, " \n")
		if len(line) == 0 {
			continue
		}
		feeds = append(feeds, line)
	}
	return feeds, nil
}

func main() {
	flag.Parse()

	urls := make([]string, 0)
	var err error

	if *useFeedFile {
		urls, err = ReadConfig()
		if err != nil {
			fmt.Println("couldn't read config file:", err)
			return
		}
	}

	if *useOpml {
		opml, err := ParseOpml(OpmlFile)
		if err != nil {
			fmt.Println("couldn't read opml file, ignoring:", err)
		} else {
			opmlFeeds := opml.ExtractFeeds()
			for _, f := range opmlFeeds {
				urls = append(urls, f)
			}
		}
	}

	if len(urls) == 0 {
		fmt.Println("no feeds")
		return
	}

	feeds := ReadFeeds(urls)
	timeline := BuildTimeline(feeds)
	if *showHeadlines {
		ShowHeadlines(timeline)
	}
}

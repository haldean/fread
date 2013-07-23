package main

import "io/ioutil"
import "encoding/xml"

type Opml struct {
	XMLName xml.Name `xml:"opml"`
	Head    OpmlHead `xml:"head"`
	Body    OpmlBody `xml:"body"`
}

type OpmlHead struct {
	Title string `xml:"title"`
}

type OpmlBody struct {
	Categories []OpmlCategory `xml:"outline"`
}

type OpmlCategory struct {
	Name  string     `xml:"text,attr"`
	Feeds []OpmlFeed `xml:"outline"`
}

type OpmlFeed struct {
	Text    string `xml:"text,attr"`
	HtmlUrl string `xml:"htmlUrl,attr"`
	FeedUrl string `xml:"xmlUrl,attr"`
	Type    string `xml:"type,attr"`
}

func ParseOpml(file string) (*Opml, error) {
	res := Opml{}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(data, &res)
	return &res, err
}

func (o *Opml) ExtractFeeds() []string {
	feeds := make([]string, 0)
	for _, cat := range o.Body.Categories {
		for _, feed := range cat.Feeds {
			feeds = append(feeds, feed.FeedUrl)
		}
	}
	return feeds
}

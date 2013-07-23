package main

import (
	"github.com/SlyMarbo/rss"
	"time"
)

type FeedItem struct {
	rss.Item
	Feed    string
	FeedUrl string
}

type Timeline []FeedItem

type mergeStack struct {
	Feed *rss.Feed
	Index int
}

func toFeedItem(f mergeStack) FeedItem {
	rssItem := f.Feed.Items[f.Index]
	return FeedItem {
		*rssItem, f.Feed.Title, f.Feed.UpdateURL,
	}
}

func BuildTimeline(feeds []*rss.Feed) Timeline {
	// merge sort! it's like coms1007 all over again
	size := 0
	stacks := make([]mergeStack, len(feeds))
	for i, f := range feeds {
		stacks[i] = mergeStack { f, 0 }
		size += len(f.Items)
	}

	minval, minindex := time.Time{}, -1
	tl := make(Timeline, size)
	idx := 0

	for {
		done := true
		minval, minindex = time.Time{}, -1
		for i, st := range stacks {
			if st.Index >= len(st.Feed.Items) {
				continue
			}
			done = false
			if minindex == -1 || st.Feed.Items[st.Index].Date.After(minval) {
				minval, minindex = st.Feed.Items[st.Index].Date, i
			}
		}
		if done { break }

		tl[idx] = toFeedItem(stacks[minindex])
		idx++
		stacks[minindex].Index++
	}

	return tl
}

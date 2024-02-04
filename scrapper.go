package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/AvivKermann/BlogAgg/internal/database"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func startScraping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	ticker := time.NewTicker(timeBetweenRequest)
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Printf("error fetching feeds %v", err)
			continue
		}
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			go scrapeFeed(db, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("error marking feed as fetched %v\n", feed.ID)
	}
	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Printf("error fetching feed %v\n", feed.ID)

	}

	for _, item := range rssFeed.Channel.Item {
		log.Println("found post", item.Title, "on feed", feed.Name)
		parsedTime := formatTimestamp(item.PubDate)
		_, err := db.CreatePosts(context.Background(), database.CreatePostsParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: item.Description,
			PublishedAt: parsedTime,
			FeedID:      feed.ID,
		})
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			continue

		} else if err != nil {
			log.Println(err.Error())
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}

func formatTimestamp(timestamp string) time.Time {
	const layoutInput = "Mon, 02 Jan 2006 15:04:05 -0700"
	const layoutOutput = "2006-01-02 15:04:05.99999"

	parsedInput, err := time.Parse(layoutInput, timestamp)
	if err != nil {
		return time.Time{}
	}
	parsedTimestamp, err := time.Parse(layoutOutput, parsedInput.Format(layoutOutput))
	if err != nil {
		return time.Time{}
	}
	return parsedTimestamp
}

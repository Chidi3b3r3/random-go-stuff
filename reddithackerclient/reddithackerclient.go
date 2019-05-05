package reddithackerclient

import (
	"fmt"
	"os"
	"sync"

	"github.com/caser/gophernews"
	"github.com/jzelinskie/geddit"
)

var redditSession *geddit.LoginSession
var hackerNewsClient *gophernews.Client

func init() {
	hackerNewsClient = gophernews.NewClient()
	var err error
	redditSession, err = geddit.NewLoginSession("username", "password", "gdAgent v1")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type Story struct {
	Title  string
	Url    string
	Author string
	Source string
}

func getHnStoryDetails(id int, ch chan<- Story, wg *sync.WaitGroup) {
	defer wg.Done()

	story, err := hackerNewsClient.GetStory(id)

	if err != nil {
		return
	}

	ch <- Story{
		Title:  story.Title,
		Url:    story.URL,
		Author: story.By,
		Source: "HackerNews",
	}
}

func NewHnStories(ch chan<- Story) {
	defer close(ch)

	changes, err := hackerNewsClient.GetChanges()
	if err != nil {
		fmt.Println(err)
	}

	var wg sync.WaitGroup

	for _, id := range changes.Items {
		wg.Add(1)
		go getHnStoryDetails(id, ch, &wg)
	}

	wg.Wait()
}

func NewRedditStories(ch chan<- Story) {
	defer close(ch)
	sort := geddit.PopularitySort(geddit.NewSubmissions)

	var listingOptions geddit.ListingOptions

	submissions, err := redditSession.SubredditSubmissions("programming", sort, listingOptions)

	if err != nil {
		fmt.Println(err)
	}

	for _, story := range submissions {
		ch <- Story{
			Title:  story.Title,
			Url:    story.URL,
			Author: story.Author,
			Source: "Reddit /r/programming",
		}
	}
}

func outputToConsole(ch <-chan Story) {
	for {
		s := <-ch
		fmt.Printf("%s - %s\n%s - %s\n\n", s.Title, s.Source, s.Author, s.Url)
	}
}

func Run() {
	hnChan := make(chan Story, 8)
	rdChan := make(chan Story, 8)
	consoleChan := make(chan Story, 8)

	go NewHnStories(hnChan)
	go NewRedditStories(rdChan)

	go outputToConsole(consoleChan)

	hnOpen := true
	rdOpen := true

	for hnOpen || rdOpen {
		select {
		case s, open := <-hnChan:
			if open {
				consoleChan <- s
			} else {
				hnOpen = false
			}
		case s, open := <-rdChan:
			if open {
				consoleChan <- s
			} else {
				rdOpen = false
			}
		}
	}
}

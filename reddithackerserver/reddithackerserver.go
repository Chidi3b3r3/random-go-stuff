package reddithackerserver

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/chidi3b3r3/gocode/reddithackerclient"
)

var stories []reddithackerclient.Story

func init() {
	stories = append(stories,
		reddithackerclient.Story{"Go Programming", "golang.com", "Josh Go", "Web"},
		reddithackerclient.Story{"Ruby Programming", "ruby.com", "Josh Ruby", "Book"},
		reddithackerclient.Story{"Java Programming", "java.com", "Josh Java", "Library"},
		reddithackerclient.Story{"Lisp Programming", "lisp.com", "Josh Lisp", "oral"},
		reddithackerclient.Story{"Python Programming", "python.com", "Josh Python", "oral"},
	)
}

func searchStories(query string) []reddithackerclient.Story {
	var foundStories []reddithackerclient.Story
	for _, story := range stories {
		if strings.Contains(strings.ToLower(story.Title), strings.ToLower(query)) {
			foundStories = append(foundStories, story)
		}
	}
	return foundStories
}

func search(w http.ResponseWriter, r *http.Request) {
	query := r.FormValue("q")

	if query == "" {
		http.Error(w, "pass in a search query param", http.StatusNotAcceptable)
	}

	w.Write([]byte("<html><body>"))
	s := searchStories(query)
	if len(s) == 0 {
		w.Write([]byte(fmt.Sprintf("No result for '%s'\n<br.", query)))
	} else {
		for _, story := range s {
			w.Write([]byte(fmt.Sprintf("<a href='%s'>%s</a> by %s on %s <br><br>", story.Url, story.Title, story.Author, story.Source)))
		}
	}
	w.Write([]byte("<a href='../'>Back</a>"))
	w.Write([]byte("</body></html>"))
}

func topTen(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<html><body>"))
	form := "<form action='search' method='get'> Search: <input type='text' name='q'> <input type='submit'></form>"
	w.Write([]byte(form))

	for i := len(stories) - 1; i >= 0 && len(stories) - 10 <= i; i-- {
		story := stories[i]
		w.Write([]byte(fmt.Sprintf("<a href='%s'>%s</a> by %s on %s <br><br>", story.Url, story.Title, story.Author, story.Source)))
	}
	w.Write([]byte("</body></html>"))
}

func Run() {
	http.HandleFunc("/", topTen)
	http.HandleFunc("/search", search)
	go runStoriesSources()

	if err := http.ListenAndServe(":9090", nil); err != nil {
		panic(err)
	}
}

func OutputToStories(ch <-chan reddithackerclient.Story) {
	for {
		s := <-ch
		stories = append(stories, s)
	}
}

func runStoriesSources() {
	go func() {
		for {
			fmt.Println("fetching new stories")
			hnChan := make(chan reddithackerclient.Story, 8)
			rdChan := make(chan reddithackerclient.Story, 8)
			toStoriesChan := make(chan reddithackerclient.Story, 8)

			go reddithackerclient.NewHnStories(hnChan)
			go reddithackerclient.NewRedditStories(rdChan)

			go OutputToStories(toStoriesChan)

			hnOpen := true
			rdOpen := true

			for hnOpen || rdOpen {
				select {
				case s, open := <-hnChan:
					if open {
						toStoriesChan <- s
					} else {
						hnOpen = false
					}
				case s, open := <-rdChan:
					if open {
						toStoriesChan <- s
					} else {
						rdOpen = false
					}
				}
			}

			fmt.Println("Done fetching new stories")
			time.Sleep(30 * time.Second)
		}
	}()
}

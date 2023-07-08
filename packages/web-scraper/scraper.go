//package main
//
//import (
//	"fmt"
//	"strings"
//
//	// importing Colly
//	"github.com/gocolly/colly"
//)
//
//func main() {
//	alphabetCollector := colly.NewCollector()
//	artistCollector := alphabetCollector.Clone()
//	lyricCollector := artistCollector.Clone()
//
//	alphabetCollector.OnRequest(func(r *colly.Request) {
//		fmt.Println("Visiting: ", r.URL)
//	})
//
//	alphabetCollector.OnHTML("a", func(e *colly.HTMLElement) {
//		href := e.Attr("href")
//		if strings.Contains(href, ".html") && strings.Contains(href, "a/") {
//			fmt.Println(href)
//			artistCollector.Visit(e.Request.AbsoluteURL(href))
//
//		}
//	})
//
//	artistCollector.OnRequest(func(r *colly.Request) {
//		fmt.Println("Visiting: ", r.URL)
//	})
//
//	artistCollector.OnError(func(_ *colly.Response, err error) {
//		fmt.Println("Something went wrong: ", err)
//	})
//
//	artistCollector.OnHTML("a", func(e *colly.HTMLElement) {
//		href := e.Attr("href")
//
//		// Filter hrefs for those that contain "/lyrics/a1"
//		if strings.Contains(href, "/lyrics/a1") {
//			// Printing only URLs associated with the a links in the page that contain "/lyrics/a1"
//			fmt.Println(href)
//
//			// Visit the URL
//			lyricCollector.Visit(e.Request.AbsoluteURL(href))
//		}
//	})
//
//	lyricCollector.OnRequest(func(r *colly.Request) {
//		fmt.Println("Visiting: ", r.URL)
//	})
//
//	lyricCollector.OnHTML("html", func(e *colly.HTMLElement) {
//		title := e.ChildText("div + b")
//		//lyrics := e.ChildText("br + br + div")
//		fmt.Println(title)
//	})
//
//	alphabetCollector.Visit("https://web.archive.org/web/20170330154346/http://www.azlyrics.com/a.html")
//	//artistCollector.Visit("https://web.archive.org/web/20170330154346/http://www.azlyrics.com/a/a1.html")
//	//lyricCollector.Visit("https://web.archive.org/web/20170228164408/http://www.azlyrics.com/lyrics/a1/ify
//	//ouweremygirl.html")
//}

package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/gocolly/colly"
	"golang.org/x/sync/semaphore"
)

func main() {
	sem := semaphore.NewWeighted(10) // Allow up to 10 concurrent calls
	ctx := context.TODO()

	alphabetCollector := colly.NewCollector()
	artistCollector := alphabetCollector.Clone()
	lyricCollector := artistCollector.Clone()

	alphabetCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	alphabetCollector.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.Contains(href, ".html") && strings.Contains(href, "a/") {
			sem.Acquire(ctx, 1)
			go func(link string) {
				defer sem.Release(1)
				artistCollector.Visit(e.Request.AbsoluteURL(link))
			}(href)
		}
	})

	artistCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	artistCollector.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	artistCollector.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")

		// Filter hrefs for those that contain "/lyrics/a1"
		if strings.Contains(href, "/lyrics/a1") {
			sem.Acquire(ctx, 1)
			go func(link string) {
				defer sem.Release(1)
				lyricCollector.Visit(e.Request.AbsoluteURL(link))
			}(href)
		}
	})

	lyricCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	lyricCollector.OnHTML("html", func(e *colly.HTMLElement) {
		title := e.ChildText("div + b")
		//lyrics := e.ChildText("br + br + div")
		fmt.Println(title)
	})

	for i := 'a'; i <= 'z'; i++ {
		alphabetCollector.Visit(fmt.Sprintf("https://web.archive.org/web/20170330154346/http://www.azlyrics.com/%c.html", i))
	}
	//alphabetCollector.Visit("https://web.archive.org/web/20170330154346/http://www.azlyrics.com/a.html")
}

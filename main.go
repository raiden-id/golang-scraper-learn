package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

type Anime struct {
	Title        string `json:"title"`
	Cover        string `json:"cover"`
	TotalEpisode string `json:"total_episode"`
	UpdatedOn    string `json:"updated_on"`
	UpdatedDay   string `json:"updated_day"`
	Endpoint     string `json:"endpoint"`
}

func timer(name string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %s\n", name, time.Since(start))
	}
}
func main(){
	defer timer("main")()
	var baseUrl string = "https://otakudesu.wiki"

	c := colly.NewCollector(
		colly.AllowedDomains("otakudesu.wiki"),
		colly.Async(true),

	)
	

	var animes []Anime

	c.OnHTML(".rapi ul > li", func(e *colly.HTMLElement) {
		title := e.DOM.Find("h2").Text()
		cover, _ := e.DOM.Find("img").Attr("src")
		totalEpisode := strings.Replace(e.DOM.Find(".epz").Text(), " ", "", -1)
		updatedOn := e.DOM.Find(".newnime").Text()
		updatedDay := e.DOM.Find(".epztipe").Text()
		endpointExists, _ := e.DOM.Find(".thumb > a").Attr("href")

		var endpoint string
		if endpointExists != "" {
				endpoint = strings.Replace(endpointExists, baseUrl+"/anime/", "", 1)
				endpoint = strings.TrimSuffix(endpoint, "/")
		}

		anime := Anime{
			Title:        title,
			Cover:        cover,
			TotalEpisode: totalEpisode,
			UpdatedOn:    updatedOn,
			UpdatedDay:   updatedDay,
			Endpoint:     endpoint,
		}

		animes = append(animes, anime)
	})

	c.OnHTML(".next.page-numbers", func(e *colly.HTMLElement) {
		next_page := e.Request.AbsoluteURL(e.Attr("href"))
		c.Visit(next_page)
	})

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	
	c.Visit("https://otakudesu.wiki/ongoing-anime/")
	c.Wait()

	
	content,err := json.Marshal(animes)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(string(content))

	

}
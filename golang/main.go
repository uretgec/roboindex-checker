package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly/v2"
	"log"
	"strings"
)

type Result struct {
	Url          	string `json:"url"`
	Status          int `json:"status"`
	MetaGooglebot   string `json:"meta_googlebot"`
	MetaRobot   	string `json:"meta_robot"`
	XRobotTag   	string `json:"x_robots_tag"`
}

type CheckData []struct {
	URL  string `json:"url"`
	Data string `json:"data"`
}

func main() {
	router := gin.Default()

	//router.Static("/assets", "./web/assets")

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(context *gin.Context) {
		context.HTML(200, "index.tmpl", gin.H{
			"title": "Agaaaaa",
		})
	})

	router.POST("/check", func(context *gin.Context) {
		var checkData CheckData

		if err := context.ShouldBindJSON(&checkData); err != nil {
			context.JSON(200, gin.H{"status":false,"error": "Wrong data"})
			return
		}
		fmt.Println("Request Body", checkData[0].URL)

		var result *Result
		results := make([]Result, 0, len(checkData))

		c := colly.NewCollector(
			colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
			/*colly.Async(true),*/
		)

		c.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)

			result = new(Result)
			result.Url = r.URL.String()
		})

		c.OnError(func(_ *colly.Response, err error) {
			log.Println("Something went wrong:", err)
		})

		c.OnResponseHeaders(func(r *colly.Response) {
			fmt.Println("VisitedHeaders", r.Headers)

			result.Status = r.StatusCode
			result.XRobotTag = r.Headers.Get("x-robots-tag")
		})

		c.OnResponse(func(r *colly.Response) {
			fmt.Println("VisitedResponse", r.Request.URL)
		})

		/*c.OnHTML("a[href]", func(e *colly.HTMLElement) {
			e.Request.Visit(e.Attr("href"))
		})*/

		c.OnHTML("meta[name]", func(e *colly.HTMLElement) {
			// fmt.Println("Meta Google Bot:", e.Attr("name"), e.Attr("content"))

			if strings.EqualFold(e.Attr("name"), "googlebot") {
				result.MetaGooglebot = e.Attr("content")
			}

			if strings.EqualFold(e.Attr("name"), "robots") {
				result.MetaRobot = e.Attr("content")
			}
		})

		c.OnScraped(func(r *colly.Response) {
			fmt.Println("Finished", r.Request.URL)

			results = append(results, *result)
		})

		for i := 0; i < len(checkData); i++ {
			// TODO: return de err kontrolü gerkeiyor
			c.Visit(checkData[i].URL)
		}

		/*c.Wait()*/
		context.JSON(200, gin.H{
			"status":true,
			"message":"",
			"result":results,
		})
	})

	router.Run()
}
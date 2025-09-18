package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"slices"
	"sync"

	"github.com/chaso-pa/real-estate-tracker/internal/models"
	"github.com/chaso-pa/real-estate-tracker/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"golang.org/x/sync/semaphore"
)

func SampleCrawl(c *gin.Context) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://suumo.jp/jj/bukken/ichiran/JJ012FC001/?ar=040&bs=030&ekTjCd=&ekTjNm=&kb=1&kj=9&km=1&kt=9999999&sc=20217&ta=20&tb=0&tj=0&tt=9999999&rssFlg=1")
	c.JSON(http.StatusOK, gin.H{
		"message": feed.Items,
	})
}

func CrawlSuumo(c *gin.Context) {
	urls := []string{
		"https://suumo.jp/jj/bukken/ichiran/JJ012FC001/?ar=040&bs=020&ekTjCd=&ekTjNm=&hb=0&ht=9999999&kb=1&km=1&kt=9999999&sc=20217&ta=20&tb=0&tj=0&tt=9999999&rssFlg=1",
		"https://suumo.jp/jj/bukken/ichiran/JJ012FC001/?ar=040&bs=021&cn=9999999&cnb=0&ekTjCd=&ekTjNm=&hb=0&ht=9999999&kb=1&kt=9999999&sc=20217&ta=20&tb=0&tj=0&tt=9999999&rssFlg=1",
		"https://suumo.jp/jj/bukken/ichiran/JJ012FC001/?ar=040&bs=030&ekTjCd=&ekTjNm=&kb=1&kj=9&km=1&kt=9999999&sc=20217&ta=20&tb=0&tj=0&tt=9999999&rssFlg=1",
	}
	for _, rawUrl := range urls {
		if err := crawlSuumoUrl(rawUrl); err != nil {
			log.Printf("Error crawling to %v: %v", url, err)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func crawlSuumoUrl(rawUrl string) error {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(rawUrl)
	openai := services.NewOpenAIService()
	schema := models.EstatesSchema()

	var wg sync.WaitGroup
	sem := semaphore.NewWeighted(int64(20))

	for c := range slices.Chunk(feed.Items, 10) {
		wg.Add(1)
		if err := sem.Acquire(context.TODO(), 1); err != nil {
			log.Println("Error when sem.Acquire")
			wg.Done()
			continue
		}
		go func(i []*gofeed.Item) {
			defer func() {
				sem.Release(1)
				wg.Done()
			}()

			jsonData, err := json.Marshal(i)
			if err != nil {
				log.Printf("Error marshaling to JSON: %v", err)
			}
			messages := []services.Message{
				{Role: "user", Content: "次の" + urlToEstateType(rawUrl) + "のJSONを解釈してJSONの配列を作成してください" + string(jsonData)},
			}
			response, err := openai.ChatCompletionWithStructuredOutput(messages, schema, "gpt-5-mini")
			if err != nil {
				log.Printf("Error getting to OpenAi Response: %v", err)
			}
			var estateResponse models.EstateResponse
			err = json.Unmarshal([]byte(response), &estateResponse)
			if err != nil {
				log.Printf("Error unmarshaling to OpenAi Response: %v", err)
			}
			models.EstatesSetValues(estateResponse.Estates)
			models.EstatesUpsert(estateResponse.Estates)
		}(c)
	}
	wg.Wait()
	return nil
}

func urlToEstateType(rawUrl string) string {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		log.Printf("Error parsing URL: %v", err)
		return ""
	}
	params := parsedUrl.Query()
	bs := params.Get("bs")

	switch bs {
	case "010":
		return "new_apartment"
	case "011":
		return "used_apartment"
	case "020":
		return "new_house"
	case "021":
		return "used_house"
	case "030":
		return "land"
	default:
		return "land"
	}
}

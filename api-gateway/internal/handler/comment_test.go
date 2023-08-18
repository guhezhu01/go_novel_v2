package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func TestGetComments(t *testing.T) {

	url := fmt.Sprintf("http://%s:%d/%s", "172.16.16.95", 3000, "api/v1/comments?article_id=1035&title=杂质少少年")
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			client := http.Client{Timeout: 5 * 60 * time.Second}

			response, err := client.Get(url)
			defer func() {
				_ = response.Body.Close()
			}()
			if err != nil {
				log.Println(err)
				return
			}
			result, _ := io.ReadAll(response.Body)

			respData := &gin.H{}
			err = json.Unmarshal(result, respData)
			if err != nil {
				log.Println(err)
			}
			log.Println(respData)
			defer wg.Done()
		}()
	}
	wg.Wait()
}

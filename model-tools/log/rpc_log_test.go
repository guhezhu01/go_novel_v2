package log

import (
	"sync"
	"testing"
)

func TestA(t *testing.T) {
	ok := InitRpcLog("log/log", "comment")
	if !ok {

	}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			Println("qweoqoeqoffsfsdf", "jdfjsdjf")
			wg.Done()
		}()

	}
	wg.Wait()
}

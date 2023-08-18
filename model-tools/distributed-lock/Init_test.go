package distributed_lock

import (
	"sync"
	"testing"
	"time"
)

func TestA(t *testing.T) {
	lock, _ := InitLock("123.249.88.132:6379", "3077267500zJ.")
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			lock.TryLock("aod", time.Millisecond*1000)
			lock.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
}

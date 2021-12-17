package broadcast_linkedlist

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestBroadcast(t *testing.T) {
	b := NewBroadcastLinkedList()

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			node := b.Node()
			count := 0
			wg.Done()
			for {
				select {
				// 即使chan已经close过了,这里依然有效
				case <-node.C():
					count++
					println(fmt.Sprintf("Data:%v count:%v", node.Data, count))
					node = node.Next()
				}
			}
		}()
	}
	wg.Wait()

	for i := 0; i < 100; i++ {
		b.Broadcast(i+1)
	}

	time.Sleep(time.Second)
}

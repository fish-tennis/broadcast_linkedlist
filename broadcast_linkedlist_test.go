package broadcast_linkedlist

import (
	"fmt"
	"runtime"
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

func TestBroadcastMem(t *testing.T) {
	println("begin")
	traceMemStats()
	b := NewBroadcastLinkedList()

	waitForBroadcast := make(chan interface{})
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			node := b.Node()
			wg.Done()
			<-waitForBroadcast
			for {
				select {
				case <-node.C():
					node = node.Next()
				}
			}
		}()
	}
	wg.Wait()

	println()
	println("before Broadcast")
	traceMemStats()

	for i := 0; i < 100; i++ {
		data := make([]byte,1024*1024,1024*1024)
		for j := 0; j < len(data); j++ {
			data[j] = 1
		}
		b.Broadcast(data)
	}
	println()
	println("after Broadcast")
	traceMemStats()
	close(waitForBroadcast)

	time.Sleep(time.Second)
	runtime.GC()
	time.Sleep(time.Second)

	println()
	println("after runtime.GC")
	traceMemStats()
}

func traceMemStats() {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	println("HeapObjects", ms.HeapObjects)
	println(fmt.Sprintf("HeapAlloc %.2f", toMB(ms.HeapAlloc)))
	println(fmt.Sprintf("TotalAlloc %.2f", toMB(ms.TotalAlloc)))
	println(fmt.Sprintf("HeapSys %.2f", toMB(ms.HeapSys)))
	println(fmt.Sprintf("HeapInuse %.2f", toMB(ms.HeapInuse)))
	println(fmt.Sprintf("StackInuse %.2f", toMB(ms.StackInuse)))
	println(fmt.Sprintf("HeapIdle %.2f", toMB(ms.HeapIdle)))
	println(fmt.Sprintf("HeapReleased %.2f", toMB(ms.HeapReleased)))
	println(fmt.Sprintf("HeapIdle-HeapReleased %.2f", toMB(ms.HeapIdle-ms.HeapReleased)))
}

func toMB(bytes uint64) float64 {
	return float64(bytes)/1024/1024
}

func BenchmarkBrodcast(b *testing.B) {
	bc := NewBroadcastLinkedList()

	go func() {
		node := bc.Node()
		for {
			select {
			case <-node.C():
				node = node.Next()
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		bc.Broadcast(nil)
	}
}
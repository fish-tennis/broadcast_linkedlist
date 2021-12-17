# broadcast_linkedlist
nonblocking broadcast chan

## example
```go
b := NewBroadcastLinkedList()

wg := sync.WaitGroup{}
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func() {
        node := b.Node()
        wg.Done()
        for {
            select {
            case <-node.C():
                println(node.Data)
                node = node.Next()
            }
        }
    }()
}
wg.Wait()

for i := 0; i < 100; i++ {
    b.Broadcast(i+1)
}
```
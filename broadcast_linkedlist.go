package broadcast_linkedlist

import (
	"sync"
)

// 广播链表
type BroadcastLinkedList struct {
	node      *BroadcastNode
	nodeMutex sync.Mutex
}

// 节点
type BroadcastNode struct {
	c    chan interface{}
	// Data将被所有监听者共享,监听者修改Data需要谨慎
	Data interface{}
	next *BroadcastNode
}

// 返回一个只读chan,供监听者监听
func (b *BroadcastNode) C() <-chan interface{} {
	return b.c
}

// 下一个监听节点
func (b *BroadcastNode) Next() *BroadcastNode {
	return b.next
}

// 构造一个新的BroadcastLinkedList
// example:
//	for i := 0; i < 10; i++ {
//		go func() {
//			node := b.Node()
//			for {
//				select {
//				case <-node.C():
//					println(node.Data)
//					node = node.Next()
//				}
//			}
//		}()
//	}
//
//	for i := 0; i < 100; i++ {
//		b.Broadcast(i+1)
//	}
func NewBroadcastLinkedList() *BroadcastLinkedList {
	return &BroadcastLinkedList{
		node: &BroadcastNode{
			c: make(chan interface{}),
		},
	}
}

// 广播一条消息
func (b *BroadcastLinkedList) Broadcast(data interface{}) {
	b.nodeMutex.Lock()
	curNode := b.node
	curNode.Data = data
	b.node = &BroadcastNode{
		c: make(chan interface{}),
	}
	curNode.next = b.node
	b.nodeMutex.Unlock()
	// 所有<-c的地方都会收到通知
	close(curNode.c)
}

// 当前节点
func (b *BroadcastLinkedList) Node() *BroadcastNode {
	// 这里不需要加锁
	return b.node
}
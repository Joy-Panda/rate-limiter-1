package types

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
)

// NewQueue 构造一个非并发安全的队列实例
func NewQueue() *PB_Queue {
	q := &PB_Queue{}
	return q.Init()
}

// Init 清空重置链表
func (q *PB_Queue) Init() *PB_Queue {
	q.Reset()
	return q
}

// PushFront 向队列头部添加一个值
func (q *PB_Queue) PushFront(value proto.Message) (*PB_Element, error) {
	v, err := MarshalAny(value)
	if err != nil {
		return nil, err
	}

	if q.Size == 0 {
		q.Head = &PB_Element{
			Prev:  nil,
			Next:  nil,
			Value: v,
		}
		q.Tail = q.Head
	} else {
		e := &PB_Element{
			Prev:  nil,
			Next:  q.Head.Next,
			Value: v,
		}
		q.Head = e
	}
	q.Size++

	return q.Head, nil
}

// PushBack 向队列中添加一个值
func (q *PB_Queue) PushBack(value proto.Message) (*PB_Element, error) {
	v, err := MarshalAny(value)
	if err != nil {
		return nil, err
	}

	return q.pushBackAny(v), nil
}

// pushBackAny 添加一个类型是*any.Any的元素到队列
func (q *PB_Queue) pushBackAny(v *any.Any) *PB_Element {
	element := &PB_Element{Value: v}

	if q.Size == 0 {
		element.Prev, element.Next = nil, nil
		q.Head, q.Tail = element, element
	} else {
		element.Prev, element.Next = q.Tail, nil
		q.Tail.Next = element
		q.Tail = element
	}

	q.Size++

	return element
}

// PushBackList 向队列中添加另一个队列的元素拷贝
func (q *PB_Queue) PushBackList(other *PB_Queue) {
	if other == nil {
		return
	}
	for itr := other.Front(); itr != nil; itr = itr.Next {
		q.pushBackAny(itr.Value)
	}
}

// PopFront 弹出队首元素，如果队列为空，则返回nil
func (q *PB_Queue) PopFront() *PB_Element {
	e := q.Front()
	if e != nil {
		q.Remove(e)
	}
	return e
}

// Front 获取链表第一个元素，如果为空链表则返回nil
func (q *PB_Queue) Front() *PB_Element {
	return q.Head
}

// Back 获取链表的尾部元素，如果链表为空，则返回空
func (q *PB_Queue) Back() *PB_Element {
	return q.Tail
}

// Remove 从链表中删除一个非空元素
func (q *PB_Queue) Remove(e *PB_Element) {
	if e == nil {
		panic("Element to remove is nil")
	}

	if e.Prev == nil {
		q.Head = e.Next
	} else {
		e.Prev.Next = e.Next
	}

	if e.Next == nil {
		q.Tail = e.Prev
	} else {
		e.Next.Prev = e.Prev
	}

	q.Size--
	e.Prev, e.Next = nil, nil
}

// Len 获取链表长度
func (q *PB_Queue) Len() uint32 {
	return q.Size
}

func (q *PB_Queue) debug() {
	var count = 0

	fmt.Printf("============== Queue Debug ==============\n")
	fmt.Printf("Len : %d\n", q.Len())
	fmt.Printf("Head: %p\n", q.Head)
	fmt.Printf("Tail: %p\n", q.Tail)
	fmt.Printf("Elements: \n")
	for itr := q.Head; itr != nil; itr = itr.Next {
		count++
		fmt.Printf("[%d] Cur:%p, Prev:%p, Next:%p\n", count, itr, itr.Prev, itr.Next)
	}
}

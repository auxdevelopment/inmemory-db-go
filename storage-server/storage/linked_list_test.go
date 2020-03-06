package storage

import (
	"fmt"
	"testing"
)

func TestAppend(t *testing.T) {
	list := NewLinkedList()

	headNode := ListNode{KeyValuePair{"Head", "HeadValue"}, nil, nil}
	firstNode := ListNode{KeyValuePair{"First", "FirstValue"}, nil, nil}

	list.AppendNode(&headNode)
	list.AppendNode(&firstNode)

	if list.Head != &headNode {
		t.Error("Head != &headNode")
	}

	if list.Head.Pair.Key != headNode.Pair.Key {
		t.Error("Head.Key != headNode.Key")
	}

	if list.Head.Next != &firstNode {
		t.Error("Head.Next != &firstNode")
	}

	if firstNode.Previous != &headNode {
		t.Error("Previous property not set correctly")
	}
}

func TestRemoveHead(t *testing.T) {
	list := NewLinkedList()

	headKey := "Head"
	headNode := ListNode{KeyValuePair{headKey, "HeadValue"}, nil, nil}
	firstNode := ListNode{KeyValuePair{"First", "FirstValue"}, nil, nil}

	list.AppendNode(&headNode)
	list.AppendNode(&firstNode)

	result := list.RemoveByKey(headKey)

	if !result {
		t.Error("Should have deleted head")
	}

	if list.Head != &firstNode {
		t.Error("First node should become head of list")
	}
}
func TestRemoveHeadOnly(t *testing.T) {
	list := NewLinkedList()

	headKey := "Head"
	headNode := ListNode{KeyValuePair{headKey, "HeadValue"}, nil, nil}

	list.AppendNode(&headNode)

	result := list.RemoveByKey(headKey)

	if !result {
		t.Error("Should have deleted head")
	}

	if list.Head != nil {
		t.Error("list.Head should be nil")
	}
}

func TestRemoveEnclosedNode(t *testing.T) {
	list := NewLinkedList()

	nodes := [4]*ListNode{
		&ListNode{KeyValuePair{"Head", "HeadValue"}, nil, nil},
		&ListNode{KeyValuePair{"First", "FirstValue"}, nil, nil},
		&ListNode{KeyValuePair{"Second", "SecondValue"}, nil, nil},
		&ListNode{KeyValuePair{"Third", "ThirdValue"}, nil, nil},
	}

	for _, value := range nodes {
		list.AppendNode(value)
	}

	list.RemoveByKey("Second")

	if list.Head.Pair.Key != "Head" {
		t.Error("Head should not have changed")
	}

	if list.Head.Next.Pair.Key != "First" {
		t.Error("First should not have changed")
	}

	if list.Head.Next.Next.Pair.Key != "Third" {
		t.Error("Second should've been replaced by third")
	}
}

func TestGetHead(t *testing.T) {
	list := NewLinkedList()

	searchKey := "Head"

	nodes := [4]*ListNode{
		&ListNode{KeyValuePair{searchKey, "HeadValue"}, nil, nil},
		&ListNode{KeyValuePair{"First", "FirstValue"}, nil, nil},
		&ListNode{KeyValuePair{"Second", "SecondValue"}, nil, nil},
		&ListNode{KeyValuePair{"Third", "ThirdValue"}, nil, nil},
	}

	for _, value := range nodes {
		list.AppendNode(value)
	}

	searchResult, err := list.Get(searchKey)

	if err != nil {
		t.Error(fmt.Sprintf("Get should find node with key '%s'", searchKey))
	}

	if searchResult.Key != "Head" {
		t.Error("Get should find the correct node")
	}
}

func TestGetArbitraryValue(t *testing.T) {
	list := NewLinkedList()

	searchKey := "Second"

	nodes := [4]*ListNode{
		&ListNode{KeyValuePair{"Head", "HeadValue"}, nil, nil},
		&ListNode{KeyValuePair{"First", "FirstValue"}, nil, nil},
		&ListNode{KeyValuePair{searchKey, "SecondValue"}, nil, nil},
		&ListNode{KeyValuePair{"Third", "ThirdValue"}, nil, nil},
	}

	for _, value := range nodes {
		list.AppendNode(value)
	}

	searchResult, err := list.Get(searchKey)

	if err != nil {
		t.Error(fmt.Sprintf("Get should find node with key '%s'", searchKey))
	}

	if searchResult.Key != searchKey {
		t.Error("Get should find the correct node")
	}
}

package storage

type ListNode struct {
	Pair     KeyValuePair
	Next     *ListNode
	Previous *ListNode
}

type LinkedList struct {
	Head *ListNode
	size int
}

func NewLinkedList() *LinkedList {
	return &LinkedList{
		nil, 0,
	}
}

func (list *LinkedList) String() string {
	current := list.Head
	stringified := ""

	for current != nil {
		stringified += current.Pair.Key + "{" + current.Pair.Value + "}"

		if current.Next != nil {
			stringified += "->"
		}

		current = current.Next
	}

	return stringified
}

func (list *LinkedList) AppendNode(node *ListNode) error {
	current := list.Head

	if current == nil {
		list.Head = node
		list.size++

		return nil
	}

	for current.Next != nil {
		current = current.Next

		if current.Pair.Key == node.Pair.Key {
			return &DuplicateKeyError{"Key already exists: " +
				node.Pair.Key, node.Pair.Key,
			}
		}
	}

	node.Previous = current
	current.Next = node

	list.size++

	return nil
}

func (list *LinkedList) Add(pair KeyValuePair) error {
	return list.AppendNode(
		&ListNode{pair, nil, nil},
	)
}

func (list *LinkedList) RemoveByKey(key string) bool {
	current := list.Head

	if current == nil {
		return false
	}

	for current != nil {
		if current.Pair.Key == key {
			hasNext := current.Next != nil
			hasPrevious := current.Previous != nil
			isHead := !hasPrevious

			if isHead {
				if !hasNext {
					list.Head = nil
					list.size--

					return true
				}

				next := list.Head.Next
				next.Previous = nil

				list.Head = next
				list.size--

				return true
			}

			if !hasNext && hasPrevious {
				previous := current.Previous
				previous.Next = nil
				list.size--

				return true
			}

			if hasNext && hasPrevious {
				previous := current.Previous
				previous.Next = nil

				next := current.Next
				next.Previous = nil

				current = nil

				next.Previous = previous
				previous.Next = next

				list.size--

				return true
			}
		}

		current = current.Next
	}

	return false
}

func (list *LinkedList) Get(key string) (*KeyValuePair, error) {
	current := list.Head

	if current == nil {
		return nil, NewKeyNotFoundError(key)
	}

	for current != nil {
		if current.Pair.Key == key {
			return &current.Pair, nil
		}

		current = current.Next
	}

	return nil, NewKeyNotFoundError(key)
}

func (list *LinkedList) pop() (KeyValuePair, error) {
	removeCandidate := list.Head.Pair
	removed := list.RemoveByKey(removeCandidate.Key)

	if !removed {
		return KeyValuePair{"error", "error"}, NewKeyNotFoundError(list.Head.Pair.Key)
	}

	return removeCandidate, nil
}

func (list *LinkedList) Size() int {
	return list.size
}

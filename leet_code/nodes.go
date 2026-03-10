package main

import "fmt"

type ListNode struct {
	Val  int
	Next *ListNode
}

func generateListNode(s []int) *ListNode {
	if len(s) == 0 {
		return nil
	}
	head := &ListNode{Val: s[0], Next: nil}
	current := head
	for i := 1; i < len(s); i++ {
		current.Next = &ListNode{Val: s[i], Next: nil}
		current = current.Next
	}
	return head
}

func (n *ListNode) getBaseNode() *ListNode {
	current := n
	for current.Next != nil {
		current = current.Next
	}
	return current
}

func AddListNode(l1, l2 *ListNode) *ListNode {
	base := l1.getBaseNode()
	base.Next = l2
	return l1
}

func (n *ListNode) printListNode() {
	current := n
	for current != nil {
		fmt.Println(current.Val, current.Next)
		current = current.Next
	}
}

func mergeListNodes(l1, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	current := dummy

	for l1 != nil && l2 != nil {
		if l1.Val < l2.Val {
			newNode := &ListNode{Val: l1.Val}
			current.Next = newNode
			l1 = l1.Next
		} else {
			newNode := &ListNode{Val: l2.Val}
			current.Next = newNode
			l2 = l2.Next
		}
		current = current.Next
	}

	remaining := l1
	if l2 != nil {
		remaining = l2
	}
	for remaining != nil {
		newNode := &ListNode{Val: remaining.Val}
		current.Next = newNode
		current = current.Next
		remaining = remaining.Next
	}

	return dummy.Next
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    current := dummy
    carry := 0

    for l1 != nil || l2 != nil || carry > 0 {
        sum := carry
        if l1 != nil {
            sum += l1.Val
            l1 = l1.Next
        }
        if l2 != nil {
            sum += l2.Val
            l2 = l2.Next
        }

        carry = sum / 10
        current.Next = &ListNode{Val: sum % 10}
        current = current.Next
    }

    return dummy.Next
}


func main() {

	// l1 := generateListNode([]int{1, 3, 5, 7})
	// l2 := generateListNode([]int{2, 4, 6, 8})

	// // ad := AddListNode(l1, l2)
	// // ad.printListNode()
	// fmt.Println("_____")
	// m := mergeListNodes(l1, l2)
	// m.printListNode()

	l1 := generateListNode([]int{2, 4, 3})
	l2 := generateListNode([]int{5, 6, 4})
	// expected [7,0,8]

	r := addTwoNumbers(l1, l2)
	r.printListNode()

}

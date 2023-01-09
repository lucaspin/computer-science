package binary_heap

/*
 * A binary heap must satisfy two properties:
 *
 *	1. The heap property:
 *	  - In a max-heap, the key of the parent must be greater than the keys of its children.
 *      The root of the heap holds the greatest value.
 *	  - In a min-heap, the key of the parent must be smaller than the keys of its children.
 *      The root of the heap holds the smallest value.
 *
 *	2. The shape property:
 *	  - A binary heap must be a complete binary tree, e.g., at every level,
 *	    except possibly the last, the tree is completely filled, and all nodes in the last level are as far left as possible.
 *
 */

type BinaryHeap[T any] struct {
	nodes      []T
	isBiggerFn func(a, b T) bool
}

/*
 * This method is suboptimal. It is here for my own educational purposes.
 * It even has a fancy name: Williams' method, after the guy who came up with binary heaps in the 60s.
 *
 * It's based on this idea:
 *   1. Start with an empty heap.
 *   2. Push() each element from the input list.
 *
 * The reason it is suboptimal is because each Push() operation runs on *O(log n)*,
 * and we do it for each element on the input list, so this method runs on *O(n log n)*.
 */
func NewBinaryHeap[T any](list []T, isBiggerFn func(a, b T) bool) BinaryHeap[T] {
	heap := BinaryHeap[T]{
		nodes:      []T{},
		isBiggerFn: isBiggerFn,
	}

	for _, e := range list {
		heap.Push(e)
	}

	return heap
}

/*
 * The method is based on this idea:
 *   1. Start by satisfying the shape property, e.g. add the element in the leftmost open bottom position
 *   2. swap(child, parent) until they satisfy the heap property.
 *      Check if the heap property is satisfied for the element and its newly parent.
 *      If yes, stop. If the heap property is satisfied there, it will be satisfied on upper levels as well.
 *      If not, swap(child, parent), and check again (one level up now).
 *
 * This operation of "moving the child up until it satisfies the heap property" is known by many names: up-heap, sift-up, heapify-up.
 * But the idea is simple, move the element up until the heap property is satisfied.
 *
 * Here, the efficiency is measured by how many swap() operations we have to do.
 *   - The worst-case scenario is the new element becoming the root of the tree.
 * 	   The height of a complete binary tree is always log(n), so O(log n).
 *   - But, if there are a lot of insertions, and they are randomized, the average-case complexity can become O(1).
 *     See: https://ieeexplore.ieee.org/document/6312854
 *
 * I'll just use up() here, because I had to look up the meaning of the word sift.
 */
func (h *BinaryHeap[T]) Push(e T) {
	// The leftmost open bottom position is always len(h) - so we append here.
	h.nodes = append(h.nodes, e)

	// Move the new element up, until heap property is satisfied.
	h.up(len(h.nodes) - 1)
}

/*
 * The Pop() operation has a similar logic to the Push() operation, but it's kind of inversed:
 *   1. Swap the current root with the last element in the tree.
 *   2. Move the new root (previously last element) down until the heap property is satisfied.
 *      We do that by comparing the parent with its children: if the parent already satisfies
 *      the heap property (for min-heap, it's smaller and its children; for max-heap, it's bigger),
 *      there's nothing to do, just stop. If it still doesn't, we swap the parent with a child
 *      (the biggest child if max-heap, and the smallest child if min-heap).
 *      We do this until the condition is no longer true.
 */
func (h *BinaryHeap[T]) Pop() *T {
	if len(h.nodes) == 0 {
		return nil
	}

	// The element to pop from the heap
	e := h.nodes[0]

	// Replace the root with the last element on the last level
	h.nodes[0] = h.nodes[len(h.nodes)-1]

	// Remove the last element
	h.nodes = h.nodes[0 : len(h.nodes)-1]

	// Re-establish heap property, if needed, and return popped element.
	h.down()
	return &e
}

func (h *BinaryHeap[T]) up(childIndex int) {
	// We iterate until we reach the root of the tree
	for childIndex > 0 {

		// The parent of node at index *i*, is always at *floor(i-1/2)*
		parentIndex := (childIndex - 1) / 2

		// If the parent already satisfies the heap property, stop.
		if h.isBiggerFn(h.nodes[parentIndex], h.nodes[childIndex]) {
			break
		}

		// If not, swap parent and child, and move one level up
		x := h.nodes[parentIndex]
		h.nodes[parentIndex] = h.nodes[childIndex]
		h.nodes[childIndex] = x
		childIndex = parentIndex
	}
}

func (h *BinaryHeap[T]) down() {
	parentIndex := 0

	for {
		// The children of a parent at index i are always at [2i+1, 2i+2]
		leftChildIndex := parentIndex*2 + 1
		rightChildIndex := parentIndex*2 + 2

		// If the current node has no children, we stop
		if leftChildIndex > len(h.nodes)-1 {
			break
		}

		var childIndex int

		// If the current node has no right child, left child is the biggest one.
		if rightChildIndex > len(h.nodes)-1 {
			childIndex = leftChildIndex
		} else if h.isBiggerFn(h.nodes[rightChildIndex], h.nodes[leftChildIndex]) {
			// If the current node has a right child, and it is bigger than the left one, we use it.
			// Note: in a max-heap, we use the bigger child and in a min-heap, we use the smaller child.
			childIndex = rightChildIndex
		} else {
			// If the current node has a right child, but it is not bigger/smaller than the left one, we use the left one.
			childIndex = leftChildIndex
		}

		// If parent already satisfies the heap property, stop.
		if h.isBiggerFn(h.nodes[parentIndex], h.nodes[childIndex]) {
			break
		}

		// Otherwise, we swap the parent with its child.
		x := h.nodes[parentIndex]
		h.nodes[parentIndex] = h.nodes[childIndex]
		h.nodes[childIndex] = x

		// And move down the tree.
		parentIndex = childIndex
	}
}

func (h *BinaryHeap[T]) Len() int {
	return len(h.nodes)
}

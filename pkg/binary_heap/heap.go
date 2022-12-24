package binary_heap

/*
 * A binary heap must satisfy two properties:
 *
 *	1. The heap property:
 *	- In a max-heap, the key of the parent must be greater than the keys of its children. The root of the heap holds the greatest value.
 *	- In a min-heap, the key of the parent must be smaller than the keys of its children. The root of the heap holds the smallest value.
 *
 *	2. The shape property:
 *	- A binary heap must be a complete binary tree, e.g., at every level,
 *	except possibly the last, the tree is completely filled, and all nodes in the last level are as far left as possible.
 *
 */

// TODO: use generics to make this great
// TODO: this has been implemented as a max-heap, but you should implement it a min-heap too.

/*
 *
 * This method is suboptimal. It is here for my own educational purposes.
 * It even has a fancy name: Williams' method, after the guy who came up with binary heaps in the 60s.
 *
 * It's based on this idea:
 *   1. Start with a empty heap.
 *   2. Push() each element from the input list.
 *
 * The reason it is suboptimal is because each Push() operation runs on *O(log n)*,
 * and we do it for each element on the input list, so this method runs on *O(n log n)*.
 *
 * */
func SuboptimalToHeap(list []int) []int {
	heap := []int{}
	for _, e := range list {
		heap = Push(heap, e)
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
 * *Rant*: why do they use words that are not familiar to non-English speakers to explain things?
 *
 * */
func Push(heap []int, e int) []int {
	// The leftmost open bottom position is always len(h) - so we append here.
	heap = append(heap, e)

	// Move the new element up, until heap property is satisfied.
	new := up(heap, len(heap)-1)
	return new
}

// TODO: this one needs a good comment.
func Pop(heap []int) ([]int, int) {
	if len(heap) == 0 {
		return []int{}, -1
	}

	// The element to pop from the heap
	e := heap[0]

	// Replace the root with the last element on the last level
	heap[0] = heap[len(heap)-1]

	// Remove the last element
	heap = heap[0 : len(heap)-1]

	// Re-establish heap property, if needed
	return down(heap), e
}

func up(heap []int, childIndex int) []int {
	// We iterate until we reach the root of the tree
	for childIndex > 0 {

		// The parent of node at index *i*, is always at *floor(i-1/2)*
		parentIndex := (childIndex - 1) / 2

		// If the parent is already greater than the child, stop -> heap property satisfied.
		if heap[parentIndex] > heap[childIndex] {
			break
		}

		// If not, swap parent and child, and move one level up
		x := heap[parentIndex]
		heap[parentIndex] = heap[childIndex]
		heap[childIndex] = x
		childIndex = parentIndex
	}

	return heap
}

func down(heap []int) []int {
	parentIndex := 0

	// iterate until the current parent has no children
	for parentIndex*2+1 < len(heap) {
		// The children of a parent at index i are always at [2i+1, 2i+2]
		left := heap[parentIndex*2+1]

		// TODO: there must be a better to write this
		right := 0
		if parentIndex*2+2 < len(heap) {
			right = heap[parentIndex*2+2]
		}

		if left > right {

			// if parent is already greater than max child, stop
			if heap[parentIndex] > left {
				break
			}

			// swap parent with left child
			x := heap[parentIndex]
			heap[parentIndex] = left
			heap[parentIndex*2+1] = x

			// and move down the tree
			parentIndex = parentIndex*2 + 1
		} else {

			// if parent is already greater than max child, stop
			if heap[parentIndex] > right {
				break
			}

			// swap parent with right child
			x := heap[parentIndex]
			heap[parentIndex] = right
			heap[parentIndex*2+2] = x

			// and move down the tree
			parentIndex = parentIndex*2 + 2
		}
	}

	return heap
}

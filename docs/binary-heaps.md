## Binary heaps

A binary tree that satisfies the [heap property](https://xlinux.nist.gov/dads/HTML/heapproperty.html). In summary:
  - In a max-heap, the key of the parent must be greater than the keys of its children. The root of the heap holds the greatest value.
  - In a min-heap, the key of the parent must be smaller than the keys of its children. The root of the heap holds the smallest value.
- A heap that uses a [binary tree](../binary-trees/README.md), for its implementation.
- The binary heap was introduced in 1964, as the data structure backing the heapsort sorting algorithm.
- Very useful when needing to access greatest/smallest values. This is the reason why they are the most efficient way to implement [priority queues](../priority-queues/README.md).
- The binary heap is not sorted. It uses a binary tree, not a binary search tree for its implementation. The heap property makes it partially ordered, but not fully sorted.
- A heap tree is a complete binary tree, so it has the smallest possible height. A heap with *N* nodes and *a* branches for each node always has *log(a)N* height.

### Implementation

- Usually implemented with an array
- Each element in the array represents a node
- The parent/child relationships are defined implicitly using the array's indexes.
- The first index contains the root element. The next 2 indices of the array contain the root's children. The next 4 indices contain the children of the root's two child nodes, and so on.
- For a node at index *i*, its children are at indexes *[2i+1,2i+2]*
- For a node at index *i*, its parent is at *floor(i-1/2)*

### Possible next steps on this subject

- Types of binary trees
- "Also very useful when insertions need to be interspersed with removals of the root node" - what does that mean?
- Heap sort
- Dijkstra algorithm for graphs shortest paths
- 2-3 heap
- b-heap
- beap
- binomial heap
- brodal queue
- d-ary heap
- Fibonnacci heap
- K-D heap
- Leftist heap
- Pairing heap
- Radix heap
- skew heaps
- soft heap
- ternary heap
- treap
- weak heap

- BST
- can be used as a dictionary and as a priority queue
- basic operations are proportional to the height of the tree
- for a complete binary tree, basic operations run in `O(lg(n))`
- if the tree is a linear chain, basic operations run in `O(n)`
- the expected height of a randomly built BST is `O(lg(n))`, but we can't always guarantee that
- each node is an object with
  - key and sattelite data
  - left, right and parent pointers
- nodes respect one property: each node has <= nodes to the left and > nodes to the right
- if a BST isn't balanced, its height is dependent on the insertion order

## Iterations

## In-order

1. print all nodes in the left subtree
2. print root node
3. print all nodes in the right subtree

## Pre-order

1. print root node
2. print all nodes in the left subtree
3. print all nodes in the right subtree

## Post-order

1. print all nodes in the left subtree
2. print all nodes in the right subtree
3. print root node
package merkletools

import (
	"crypto/sha256"
	"math"
)

type node struct {
	parent, left, right *node
	checksum            [sha256.Size]byte
}

// Tree is an append only binary merkle tree.
type Tree struct {
	root  *node
	leafs []*node
}

// NumRecords returns the number of records (= number of leaves in merkle tree)
func (t *Tree) NumRecords() int {
	return len(t.leafs)
}

// Update the checksums of nodes along path to the root
func (n *node) updateChecksums() {
	n.checksum = sha256.Sum256(append(n.left.checksum[:], n.right.checksum[:]...))
	if n.parent != nil {
		n.parent.updateChecksums()
	}
}

// MerkleRootHash returns the checksum of the root of the merkle tree.
func (t *Tree) MerkleRootHash() [sha256.Size]byte {
	if t.root == nil {
		return sha256.Sum256([]byte(""))
	}
	return t.root.checksum
}

func isPowerOfTwo(n int) bool {
	return (n != 0) && (n&(n-1)) == 0
}

// AddRecord appends newLeaf to binary tree and updates the checksums
// of nodes along root -> newLeaf branch to maintain merkle property.
func (t *Tree) AddRecord(r []byte) {
	newLeaf := new(node)
	newLeaf.checksum = sha256.Sum256(r)

	// special case for first insertion
	if t.leafs == nil {
		t.root = newLeaf
		t.leafs = append(t.leafs, newLeaf)
		return
	}

	// newLeaf is right child of newLeafParent to ensure branch factor of 2
	newLeafParent := new(node)
	newLeaf.parent = newLeafParent
	newLeafParent.right = newLeaf

	// splice in newLeafParent above the the largest filled right subtree
	cnt := len(t.leafs)
	filledSubtree := t.root
	for !isPowerOfTwo(cnt) {
		// subtract by the lower 
		cnt -= 1 << uint(math.Log2(float64(cnt)))
		filledSubtree = filledSubtree.right
	}
	newLeafParent.parent = filledSubtree.parent
	newLeafParent.left = filledSubtree
	filledSubtree.parent = newLeafParent

	if newLeafParent.parent == nil {
		t.root = newLeafParent
	}

	t.leafs = append(t.leafs, newLeaf)
	newLeaf.parent.updateChecksums()
}

func (t *Tree) Proof() {
    // TODO, along with tests :)
}

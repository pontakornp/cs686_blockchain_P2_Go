package p1

import (
	"../stack"
	"errors"
)

func (mpt *MerklePatriciaTrie) Delete(key string) (string, error) {
	node_stack, error := mpt.GetStack(key)
	if error != nil {
		return "", error
	}
	prefix := []uint8{}
	is_merging := false
	to_create_leaf := true
	new_hash_node := ""
	temp_value := ""

	for !node_stack.IsEmpty() {
		ref := node_stack.Pop().(ParentNodeRef)
		hash_node := ref.hash_node
		node := mpt.db[hash_node]
		delete(mpt.db, hash_node)
		switch node.node_type {
		case 0:
			return "", errors.New("path_not_found")
		case 1:
			branch_value_index := ref.index
			if branch_value_index == 16 {
				is_merging = true
			}
			if !is_merging {
				if branch_value_index != 17 {
					node.branch_value[branch_value_index] = new_hash_node
				}
				//node.branch_value[branch_value_index] = new_hash_node
				//delete(p1.db, hash_node)
				hash_node = node.hash_node()
				mpt.db[hash_node] = node
				new_hash_node = hash_node
				break
			}
			// set deleted node to empty string
			if branch_value_index != 17 {
				node.branch_value[branch_value_index] = ""
			}
			count := 0
			//value := ""
			tmp_prefix := []uint8{}
			for i, v := range node.branch_value[:16] {
				if v != "" {
					tmp_prefix = append(tmp_prefix, uint8(i))
					//value = v
					count++
				}
			}
			branch_last_value := node.branch_value[16]
			// delete branch node
			//delete(p1.db, hash_node)
			// if there are more than one value under branch node, no need to merge the child
			if (branch_last_value == "" && count >= 2) || (branch_last_value != "" && count >= 1) {
				// hash branch node
				hash_branch_node := node.hash_node()
				// update db
				mpt.db[hash_branch_node] = node
				// update current node
				new_hash_node = hash_branch_node
				is_merging = false
			} else { // if there are one value under branch node, merge the branch and child node prefix
				// case when there is one leaf left as branch child
				// set prefix
				prefix = tmp_prefix
				if branch_last_value == "" {
					hash_child_node := node.branch_value[tmp_prefix[0]]
					child_node := mpt.db[hash_child_node]
					switch child_node.node_type {
					case 0:
						return "", errors.New("path_not_found")
					case 1:
						temp_value = hash_child_node
						to_create_leaf = false
					case 2:
						// delete child node
						delete(mpt.db, hash_child_node)
						encoded_prefix := child_node.flag_value.encoded_prefix
						if isLeafNode(encoded_prefix) {
							to_create_leaf = true
						} else {
							to_create_leaf = false
						}
						// get the prefix and value from the child node
						child_node_prefix := compact_decode(child_node.flag_value.encoded_prefix)
						temp_value = child_node.flag_value.value
						// combine prefix
						prefix = append(prefix, child_node_prefix...)
						//// delete child node
						//delete(p1.db, hash_child_node)
					}
					// if parent is branch, then merge
					temp_str := mpt.mergeHelper(node_stack, to_create_leaf, prefix, temp_value)
					if temp_str != "" {
						new_hash_node = temp_str
						is_merging = false
					}
				} else { // case when there is one value at index 16 left in branch node
					to_create_leaf = true
					temp_value = node.branch_value[16]
					//new_hash_node = p1.mergeNodes(to_create_leaf, prefix, temp_value)
					//is_merging = false
					temp_str := mpt.mergeHelper(node_stack, to_create_leaf, prefix, temp_value)
					if temp_str != "" {
						new_hash_node = temp_str
						is_merging = false
					}
				}
			}
		case 2:
			encoded_prefix := node.flag_value.encoded_prefix
			if isLeafNode(encoded_prefix) {
				//// delete the node
				//delete(p1.db, hash_node)
				// update is merging
				is_merging = true
			} else {
				//delete(p1.db, hash_node)
				if !is_merging {
					node.flag_value.value = new_hash_node
					hash_node = node.hash_node()
					mpt.db[hash_node] = node
					new_hash_node = hash_node
					break
				}
				ext_prefix := compact_decode(node.flag_value.encoded_prefix)
				// combine prefix
				prefix = append(ext_prefix, prefix...)
				new_hash_node = mpt.mergeNodes(to_create_leaf, prefix, temp_value)
				is_merging = false
			}
		}
	}
	// if is_merging is true
	if is_merging {
		new_hash_node = mpt.mergeNodes(to_create_leaf, prefix, temp_value)
	}
	// update root
	mpt.root = new_hash_node
	return "", errors.New("path_not_found")
}

func (mpt *MerklePatriciaTrie) mergeNodes(to_create_leaf bool, prefix []uint8, value string) string {
	node := newEmptyNode()
	if to_create_leaf { // to create leaf
		// create leaf
		node = newLeafNode(prefix, value)
	} else { // to create extension
		// create extension
		node = newExtNode(prefix, value)
	}
	// hash node
	hash_node := node.hash_node()
	// update db
	mpt.db[hash_node] = node
	return hash_node
}

func (mpt *MerklePatriciaTrie) mergeHelper(node_stack *stack.Stack, to_create_leaf bool, prefix []uint8, value string) string {
	if !node_stack.IsEmpty() {
		ref := node_stack.Peek().(ParentNodeRef)
		hash_node := ref.hash_node
		node := mpt.db[hash_node]
		if node.node_type == 1 {
			new_hash_node := mpt.mergeNodes(to_create_leaf, prefix, value)
			return new_hash_node
		}
	}
	return ""
}
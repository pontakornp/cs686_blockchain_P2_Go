package project

import "reflect"

type ParentNodeRef struct {
	hash_node string
	index uint8
}

func isEmpty(node Node) bool {
	return reflect.DeepEqual(node, nil)
}

func newEmptyNode() Node {
	node := Node{
		0,
		[17]string{},
		Flag_value{[]uint8{}, ""},
	}
	return node
}

func newBranchNode(branch_value [17]string, value string) Node {
	if(value != "") {
		branch_value[len(branch_value) - 1] = value
	}
	flag_value := Flag_value{[]uint8{},""}
	node := Node {
		1,
		branch_value,
		flag_value,
	}
	return node
}

func newExtNode(prefix []uint8, value string) Node {
	encoded_prefix := compact_encode(prefix)
	flag := Flag_value {
		encoded_prefix,
		value,
	}
	node := Node {
		2,
		[17]string{},
		flag,
	}
	return node
}

func newLeafNode(prefix []uint8, value string) Node {
	prefix = append(prefix, 16)
	encoded_prefix := compact_encode(prefix)
	flag := Flag_value {
		encoded_prefix,
		value,
	}
	node := Node {
		2,
		[17]string{},
		flag,
	}
	return node
}

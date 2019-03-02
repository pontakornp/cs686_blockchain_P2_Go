package p1

import (
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
	"strings"
)

type Flag_value struct {
	encoded_prefix []uint8 // shared nibble(s) for ext node or key for leaf node
	value string // hash node or value of leaf
}

type Node struct {
	node_type int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value Flag_value
}

type MerklePatriciaTrie struct {
	db map[string]Node
	root string
}

func (mpt *MerklePatriciaTrie) GetDb() map[string]Node {
	return mpt.db
}

func (mpt *MerklePatriciaTrie) GetRoot() string {
	return mpt.root
}

// encode hex_array to ascii
func compact_encode(hex_array []uint8) []uint8 {
	term := 0
	if hex_array[len(hex_array) - 1] == 16 {
		term = 1
	}
	if term == 1 {
		hex_array = hex_array[:len(hex_array) - 1]
	}
	var odd_len int = len(hex_array) % 2
	var flags uint8 = uint8(2 * term + odd_len)
	if odd_len == 1 {
		hex_array  = append([]uint8{flags}, hex_array...)
	} else {
		hex_array = append([]uint8{flags, 0}, hex_array...)
	}
	o := []uint8{}
	for i := 0; i < len(hex_array); i+=2 {
		o = append(o, 16 * hex_array[i] + hex_array[i + 1])
	}
	return o
}

// If Leaf, ignore 16 at the end
func compact_decode(encoded_arr []uint8) []uint8 {
	hex_array := []uint8{}
	for i := 0; i < len(encoded_arr); i++ {
		hex_array = append(hex_array, encoded_arr[i] / 16)
		hex_array = append(hex_array, encoded_arr[i] % 16)
	}
	if hex_array[0] == 0 || hex_array[0] == 2 {
		hex_array = hex_array[2:]
	} else {
		hex_array = hex_array[1:]
	}
	return hex_array
}

func Test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
	fmt.Println("", compact_encode([]uint8{2, 6, 3, 16}))
	fmt.Println("", compact_decode(compact_encode([]uint8{5, 16})))
	//fmt.Println("", compact_decode(compact_encode([]uint8{2, 6, 3, 16})))
}

func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value
	}

	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

func (node *Node) String() string {
	str := "empty string"
	switch node.node_type {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.branch_value[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.branch_value[16])
	case 2:
		encoded_prefix := node.flag_value.encoded_prefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.flag_value.value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) Initial() {
	mpt.db = make(map[string]Node)
	mpt.root = ""
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0] / 16 < 2
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.root)
	for hash := range mpt.db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.db[hash]))
	}
	return content
}

func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart" + cur_hash + "HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart" + cur_hash + "HashEnd", fmt.Sprintf("Hash%v", i),  -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}
package project

// check if encoded array is a leaf
func isLeafNode(encoded_arr []uint8) bool {
	prefix := encoded_arr[0] / 16
	if prefix == 0 || prefix == 1 {
		return false
	}
	return true
}

func ConvertStringToHexArray(str string) []uint8 {
	hex_array := []uint8{}
	for i := 0; i < len(str); i++ {
		hex_array = append(hex_array, str[i]/16)
		hex_array = append(hex_array, str[i]%16)
	}
	return hex_array
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
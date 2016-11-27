package predicate

func BytesAll(array []byte, b byte) bool {
	for _, val := range array {
		if val != b {
			return false
		}
	}
	return true
}

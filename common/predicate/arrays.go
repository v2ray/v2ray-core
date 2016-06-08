package predicate

func BytesAll(array []byte, b byte) bool {
	for _, v := range array {
		if v != b {
			return false
		}
	}
	return true
}

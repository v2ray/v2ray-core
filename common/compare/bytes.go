package compare

import "v2ray.com/core/common/errors"

func BytesEqualWithDetail(a []byte, b []byte) error {
	if len(a) != len(b) {
		return errors.New("mismatch array length ", len(a), " vs ", len(b))
	}
	for idx, v := range a {
		if b[idx] != v {
			return errors.New("mismatch array value at index [", idx, "]: ", v, " vs ", b[idx])
		}
	}
	return nil
}

func BytesEqual(a []byte, b []byte) bool {
	return BytesEqualWithDetail(a, b) == nil
}

func BytesAll(arr []byte, value byte) bool {
	for _, v := range arr {
		if v != value {
			return false
		}
	}

	return true
}

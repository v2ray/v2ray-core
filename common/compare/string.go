package compare

import "v2ray.com/core/common/errors"

func StringEqualWithDetail(a string, b string) error {
	if a != b {
		return errors.New("Got ", b, " but want ", a)
	}
	return nil
}

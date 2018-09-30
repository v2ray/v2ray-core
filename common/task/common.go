package task

import "v2ray.com/core/common"

func Close(v interface{}) Task {
	return func() error {
		return common.Close(v)
	}
}

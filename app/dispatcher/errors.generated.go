package dispatcher

import "v2ray.com/core/common/errors"

func newError(values ...interface{}) *errors.Error { return errors.New(values...).Path("App", "Dispatcher") }

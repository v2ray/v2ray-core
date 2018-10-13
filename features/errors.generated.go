package features

import "v2ray.com/core/common/errors"

type errPathObjHolder struct {}
func newError(values ...interface{}) *errors.Error { return errors.New(values...).WithPathObj(errPathObjHolder{}) }

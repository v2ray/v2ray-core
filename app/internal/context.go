package internal

type contextImpl struct {
	callerTag string
}

func (this *contextImpl) CallerTag() string {
	return this.callerTag
}

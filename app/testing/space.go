package testing

type Context struct {
	CallerTagValue string
}

func (this *Context) CallerTag() string {
	return this.CallerTagValue
}

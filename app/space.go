package app

type ID int

// Context of a function call from proxy to app.
type Context interface {
	CallerTag() string
}

type Caller interface {
	Tag() string
}

// A Space contains all apps that may be available in a V2Ray runtime.
// Caller must check the availability of an app by calling HasXXX before getting its instance.
type Space interface {
	HasApp(ID) bool
	GetApp(ID) interface{}
	BindApp(ID, interface{})
}

type spaceImpl struct {
	cache map[ID]interface{}
}

func NewSpace() Space {
	return &spaceImpl{
		cache: make(map[ID]interface{}),
	}
}

func (this *spaceImpl) HasApp(id ID) bool {
	_, found := this.cache[id]
	return found
}

func (this *spaceImpl) GetApp(id ID) interface{} {
	obj, found := this.cache[id]
	if !found {
		return nil
	}
	return obj
}

func (this *spaceImpl) BindApp(id ID, object interface{}) {
	this.cache[id] = object
}

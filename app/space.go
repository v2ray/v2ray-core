package app

type ID int

// Context of a function call from proxy to app.
type Context interface {
	CallerTag() string
}

// A Space contains all apps that may be available in a V2Ray runtime.
// Caller must check the availability of an app by calling HasXXX before getting its instance.
type Space interface {
	HasApp(ID) bool
	GetApp(ID) interface{}
}

type ForContextCreator func(Context, interface{}) interface{}

var (
	metadataCache = make(map[ID]ForContextCreator)
)

func RegisterApp(id ID, creator ForContextCreator) {
	// TODO: check id
	metadataCache[id] = creator
}

type contextImpl struct {
	callerTag string
}

func (this *contextImpl) CallerTag() string {
	return this.callerTag
}

type spaceImpl struct {
	cache map[ID]interface{}
	tag   string
}

func newSpaceImpl(tag string, cache map[ID]interface{}) *spaceImpl {
	space := &spaceImpl{
		tag:   tag,
		cache: make(map[ID]interface{}),
	}
	context := &contextImpl{
		callerTag: tag,
	}
	for id, object := range cache {
		creator, found := metadataCache[id]
		if found {
			space.cache[id] = creator(context, object)
		}
	}
	return space
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

// A SpaceController is supposed to be used by a shell to create Spaces. It should not be used
// directly by proxies.
type SpaceController struct {
	objectCache map[ID]interface{}
}

func NewController() *SpaceController {
	return &SpaceController{
		objectCache: make(map[ID]interface{}),
	}
}

func (this *SpaceController) Bind(id ID, object interface{}) {
	this.objectCache[id] = object
}

func (this *SpaceController) ForContext(tag string) Space {
	return newSpaceImpl(tag, this.objectCache)
}

package cp

import (
	"net/http"
	"reflect"
	"slices"
	"strings"

	"github.com/creamsensation/cp/internal/constant/componentState"
	"github.com/creamsensation/cp/internal/constant/queryKey"
	"github.com/creamsensation/cp/internal/querystring"
	"github.com/creamsensation/cp/internal/result"
	"github.com/creamsensation/cp/internal/util"
	"github.com/creamsensation/gox"
)

type Component interface {
	Control
	Main() Control
}

type component interface {
	Component
	Model()
	Name() string
	Node() gox.Node
}

type componentLifecycle struct {
	control *control
	cv      reflect.Value
	ct      reflect.Type
	name    string
}

const (
	componentMethodModel = "Model"
	componentMethodName  = "Name"
	componentMethodNode  = "Node"
)

const (
	componentFieldStateTag = "state"
	componentFieldStateUse = "true"
)

var (
	componentMethods              = []string{componentMethodModel, componentMethodName, componentMethodNode}
	componentControlInterfaceName = util.GetInterfaceName[Component]()
)

func createComponentControl(c *control, ct component) *control {
	cc := new(control)
	*cc = *c
	cc.component = ct
	cc.main = c
	cc.flash = &flashMessenger{control: cc, name: createPrefixedComponentName(cc)}
	cc.state = createState(cc)
	return cc
}

func createComponentLifecycle(c *control) *componentLifecycle {
	ref := reflect.ValueOf(c.component)
	if ref.Type().Kind() != reflect.Ptr {
		return nil
	}
	return &componentLifecycle{
		control: c,
		cv:      ref,
		ct:      ref.Elem().Type(),
		name:    createPrefixedComponentName(c),
	}
}

func (l *componentLifecycle) run() *componentLifecycle {
	cfg := l.control.Config().Component
	useCacheState := cfg.State == componentState.Cache
	useQueryState := cfg.State == componentState.Query
	isAction := l.control.Request().Is().Action()
	if useCacheState {
		l.loadFromCache()
	}
	l.injectDeps()
	if isAction && useQueryState {
		l.prepareFromRequest(l.control.request)
	}
	l.control.component.Model()
	if isAction {
		l.callAction()
		if useCacheState {
			l.storeToCache()
		}
	}
	return l
}

func (l *componentLifecycle) node() gox.Node {
	return l.control.component.Node()
}

func (l *componentLifecycle) prepareFromRequest(r *http.Request) *componentLifecycle {
	querystring.New(l.control.component).
		Request(r).
		IgnoreInterface(componentControlInterfaceName).
		Decode()
	return l
}

func (l *componentLifecycle) injectDeps() *componentLifecycle {
	controlRef := reflect.ValueOf(l.control)
	for i := 0; i < l.cv.Elem().NumField(); i++ {
		field := l.cv.Elem().Field(i)
		if field.Type().String() == componentControlInterfaceName {
			field.Set(controlRef)
			continue
		}
		if field.Type().Kind() != reflect.Interface && field.Type().Kind() != reflect.Struct {
			continue
		}
		dep := provide[any](l.control, field.Type(), field.Type().String())
		if dep == nil {
			continue
		}
		field.Set(reflect.ValueOf(dep))
	}
	return l
}

func (l *componentLifecycle) callAction() {
	action := l.control.Request().Query(queryKey.Action)
	cn := action[:strings.LastIndex(action, linkLevelDivider)]
	if cn != l.name {
		return
	}
	an := action[strings.LastIndex(action, linkLevelDivider)+1:]
	if slices.Contains(componentMethods, an) {
		return
	}
	method := l.cv.MethodByName(an)
	if !method.IsValid() {
		return
	}
	methodResult := method.Call([]reflect.Value{})
	if len(methodResult) == 0 {
		return
	}
	if methodResult[0].Interface() == nil {
		return
	}
	l.control.main.result = methodResult[0].Interface().(result.Result)
}

func (l *componentLifecycle) loadFromCache() {
	r := reflect.New(l.cv.Elem().Type())
	l.control.State().Get(r.Interface())
	for i := 0; i < l.cv.Elem().NumField(); i++ {
		field := r.Elem().Field(i)
		fieldName := r.Elem().Type().Field(i).Name
		if field.IsZero() || field.Type().String() == componentControlInterfaceName {
			continue
		}
		fieldStruct, ok := r.Elem().Type().FieldByName(fieldName)
		if ok && fieldStruct.Tag.Get(componentFieldStateTag) != componentFieldStateTag {
			continue
		}
		l.cv.Elem().FieldByName(fieldName).Set(field)
	}
}

func (l *componentLifecycle) storeToCache() {
	c := reflect.New(l.cv.Elem().Type())
	for i := 0; i < l.cv.Elem().NumField(); i++ {
		field := l.cv.Elem().Field(i)
		fieldName := l.cv.Elem().Type().Field(i).Name
		if field.Type().String() == componentControlInterfaceName {
			continue
		}
		fieldStruct, ok := c.Elem().Type().FieldByName(fieldName)
		if ok && fieldStruct.Tag.Get(componentFieldStateTag) != componentFieldStateTag {
			continue
		}
		c.Elem().FieldByName(fieldName).Set(field)
	}
	l.control.State().Set(c.Elem().Interface())
}

func createPrefixedComponentName(c *control) string {
	name := c.component.Name()
	if len(c.route.Controller) > 0 {
		name = c.route.Controller + linkLevelDivider + name
	}
	if len(c.route.Module) > 0 {
		name = c.route.Module + linkLevelDivider + name
	}
	return name
}

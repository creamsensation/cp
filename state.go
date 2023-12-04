package cp

import (
	"encoding/json"
	"time"
	
	"github.com/dchest/uniuri"
	
	"github.com/creamsensation/cp/internal/constant/cookieName"
)

type StateManager interface {
	Get(result any)
	Set(data any)
	Reset()
}

var (
	stateExpiration = time.Hour * 24 * 30
)

type stateManager struct {
	control   Control
	name      string
	token     string
	cacheKey  string
	cookieKey string
	reset     bool
}

type state struct {
	Exist     bool   `json:"exists"`
	Lang      string `json:"lang"`
	Ip        string `json:"ip"`
	UserAgent string `json:"userAgent"`
	Data      []byte `json:"data"`
}

func createState(control *control) StateManager {
	s := &stateManager{control: control}
	if control.component != nil {
		s.name = createPrefixedComponentName(control)
	}
	s.cookieKey = s.createCookieKey()
	s.token = control.Cookie().Get(s.cookieKey)
	if len(s.token) == 0 {
		s.token = uniuri.New()
	}
	s.cacheKey = s.createCacheKey()
	return s
}

func (s *stateManager) Get(result any) {
	exists := s.control.Cache().Exists(s.cacheKey)
	if !exists {
		return
	}
	var r state
	s.control.Cache().Get(s.cacheKey, &r)
	if s.control.Request().Lang() != r.Lang ||
		s.control.Request().Ip() != r.Ip ||
		s.control.Request().UserAgent() != r.UserAgent {
		return
	}
	s.control.Error().Check(json.Unmarshal(r.Data, result))
}

func (s *stateManager) Set(data any) {
	if s.reset {
		return
	}
	dataBytes, err := json.Marshal(data)
	if err != nil {
		s.control.Error().Check(err)
	}
	s.control.Cache().Set(s.cacheKey, s.createState(dataBytes), stateExpiration)
	s.control.Cookie().Set(s.cookieKey, s.token, stateExpiration)
}

func (s *stateManager) Reset() {
	s.reset = true
	s.control.Cache().Set(s.cacheKey, "", time.Millisecond)
	s.control.Cookie().Set(s.cookieKey, "", time.Millisecond)
}

func (s *stateManager) createCacheKey() string {
	if len(s.name) == 0 {
		return "state:" + s.token
	}
	return "state-" + s.name + ":" + s.token
}

func (s *stateManager) createCookieKey() string {
	if len(s.name) == 0 {
		return cookieName.State
	}
	return cookieName.State + "-" + s.name
}

func (s *stateManager) createState(data []byte) state {
	return state{
		Exist:     true,
		Lang:      s.control.Request().Lang(),
		Ip:        s.control.Request().Ip(),
		UserAgent: s.control.Request().UserAgent(),
		Data:      data,
	}
}

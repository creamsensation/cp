package session

type Session struct {
	Id        int      `json:"id"`
	Email     string   `json:"email"`
	Roles     []string `json:"role"`
	Super     bool     `json:"super"`
	Ip        string   `json:"ip"`
	UserAgent string   `json:"userAgent"`
}

type Field struct {
	Name  string
	Value any
}

func (s Session) Map() map[string]any {
	if s.Id == 0 {
		return map[string]any{}
	}
	if len(s.Ip) == 0 {
		s.Ip = "localhost"
	}
	return map[string]any{
		"Id":        s.Id,
		"Email":     s.Email,
		"Roles":     s.Roles,
		"Super":     s.Super,
		"Ip":        s.Ip,
		"UserAgent": s.UserAgent,
	}
}

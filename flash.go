package cp

type Flash interface {
	Get() ([]Message, error)
	Success(title, value string)
	Warning(title, value string)
	Error(title, value string)
	
	MustGet() []Message
}

type flash struct {
	state *state
}

type Message struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Value string `json:"value"`
}

const (
	FlashSuccess = "success"
	FlashWarning = "warning"
	FlashError   = "error"
)

func (f flash) Get() ([]Message, error) {
	messages := f.state.Messages
	f.state.Messages = make([]Message, 0)
	return messages, f.state.save()
}

func (f flash) MustGet() []Message {
	messages, err := f.Get()
	if err != nil {
		panic(err)
	}
	return messages
}

func (f flash) Success(title, value string) {
	f.state.Messages = append(f.state.Messages, Message{Type: FlashSuccess, Title: title, Value: value})
}

func (f flash) Warning(title, value string) {
	f.state.Messages = append(f.state.Messages, Message{Type: FlashWarning, Title: title, Value: value})
}

func (f flash) Error(title, value string) {
	f.state.Messages = append(f.state.Messages, Message{Type: FlashError, Title: title, Value: value})
}

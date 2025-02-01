package push

type Push struct {
}

func NewPush() *Push {
	return &Push{}
}

func (p *Push) PushIOS() {
	// do something
}

func (p *Push) PushAndroid() {
	// do something
}

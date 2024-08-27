package helper

type Amer interface {
	doSomething() Amer
}
type Mohame interface {
	Amer
	anotherThing()
}

type Nothing struct {
}

func (n *Nothing) doSomething() Amer {
	return n
}

package interfaces

type uciInterface struct {
}

func NewUciProtocolHandler() uciInterface {
	return uciInterface{}
}

func (i uciInterface) Loop() {

}

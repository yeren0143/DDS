package edp

type IEDPListener interface {
}


// Class EDPSimplePUBReaderListener, used to define the behavior when a new WriterProxyData is received.
type EDPSimplePubListener struct {

}

func newEDPSimplePubListener(edp *EDPSimple) *EDPSimplePubListener {
	return &EDPSimplePubListener{}
}

//Class EDPSimpleSUBReaderListener, used to define the behavior when a new ReaderProxyData is received.
type EDPSimpleSubListener struct {

}

func newEDPSimpleSubListener(edp *EDPSimple) *EDPSimpleSubListener {
	return &EDPSimpleSubListener{}
}
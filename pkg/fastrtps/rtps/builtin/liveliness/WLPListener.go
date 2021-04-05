package liveliness

import "github.com/yeren0143/DDS/fastrtps/rtps/reader"

var _ reader.IReaderListener = (*WLPListener)(nil)

type WLPListener struct {
	reader.ReaderListenerBase
}

// Class WLPListener that receives the liveliness messages asserting the liveliness of remote endpoints.type WLPListener struct {}
func NewWlpListener(wlp *WLP) *WLPListener {
	return &WLPListener{}
}

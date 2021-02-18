package reader

import "github.com/yeren0143/DDS/fastrtps/rtps/history"

var _ history.IWriterProxyWithHistory = (*WriterProxy)(nil)

type WriterProxy struct {
}

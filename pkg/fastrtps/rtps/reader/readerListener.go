package reader

// ReaderListener to be used by the user to override some of is virtual method to program actions to
// certain events.
type ReaderListener struct {
}

// This method is invoked when a new reader matches
func (listener *ReaderListener) OnReaderMatched(reader *IRTPSReader) {
}

// This method is invoked when a new reader matches
func (listener *ReaderListener) OnReaderMatchedWithStatus() {

}

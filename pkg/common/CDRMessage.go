package common

const (
	//RTPSMessageDefaultSize define max size of rtps message in bytes
	RTPSMessageDefaultSize = 10500

	//RTPSMessageCommonRTPSPayloadSize define common payload a rtps message has TODO(Ricardo) It is necessary?
	RTPSMessageCommonRTPSPayloadSize = 536

	//RTPSMessageCommonDataPayloadSize define common data size
	RTPSMessageCommonDataPayloadSize = 1000

	//RTPSMessageHeaderSize define header size in bytes
	RTPSMessageHeaderSize = 20

	//RTPSMessageSubMessageHeaderSize define submessage header size
	RTPSMessageSubMessageHeaderSize = 4

	//RTPSMessageDataExtraInlineQosSize ?
	RTPSMessageDataExtraInlineQosSize = 4

	//RTPSMessageInfoTSSize ?
	RTPSMessageInfoTSSize = 12

	//RTPSMessageOcetsToInlineQosDataSubMsg ...
	RTPSMessageOcetsToInlineQosDataSubMsg = 16

	//RTPSMessageOcetsToInlineQosDataFragSubMsg ...
	RTPSMessageOcetsToInlineQosDataFragSubMsg = 28

	//RTPSMessageDataMinLength ...
	RTPSMessageDataMinLength = 24
)

//CDRMessageT contains a serialized message.
type CDRMessageT struct {
	//Buffer Pointer to the buffer where the data is stored.
	Buffer []Octet

	//Read or write position.
	Pos uint32

	//Max size of the message.
	MaxSize uint32

	//Size allocated on buffer. May be higher than max_size.
	ReservedSize uint32

	//Current length of the message.
	Length uint32

	//Whether this message is wrapping a buffer managed elsewhere.
	Wrap bool

	//Endianness of the message
	MsgEndian Endianness
}

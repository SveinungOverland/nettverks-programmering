package frame

type Frame struct {
	IsFragment bool
	Opcode     byte
	Reserved   byte
	IsMasked   bool
	Length     uint64
	Payload    []byte
}

// Text payload
func (f Frame) Text() string {
	return string(f.Payload)
}

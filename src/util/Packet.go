package util

type Packet struct {
	Form string

	StringData string

	// Variables for file transfers
	BytePos   int64
	FileData  []byte
	Completed bool
}

func (packet *Packet) GetForm() string {
	return packet.Form
}

func (packet *Packet) SetForm(newForm string) {
	packet.Form = newForm
}

func (packet *Packet) GetStringData() string {
	return packet.StringData
}

func (packet *Packet) SetStringData(newStringData string) {
	packet.StringData = newStringData
}

func (packet *Packet) GetBytePos() int64 {
	return packet.BytePos
}

func (packet *Packet) SetBytePos(newBytePos int64) {
	packet.BytePos = newBytePos
}

func (packet *Packet) GetFileData() []byte {
	return packet.FileData
}

func (packet *Packet) SetFileData(newFileData []byte) {
	packet.FileData = newFileData
}

func (packet *Packet) GetComplete() bool {
	return packet.Completed
}

func (packet *Packet) SetCompleted(newCompleted bool) {
	packet.Completed = newCompleted
}

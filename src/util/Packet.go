package util

type Packet struct {
	form string

	stringData string

	// Variables for file transfers
	bytePos   int64
	fileData  []byte
	completed bool
}

func (packet *Packet) GetForm() string {
	return packet.form
}

func (packet *Packet) SetForm(newForm string) {
	packet.form = newForm
}

func (packet *Packet) GetStringData() string {
	return packet.stringData
}

func (packet *Packet) SetStringData(newStringData string) {
	packet.stringData = newStringData
}

func (packet *Packet) GetBytePos() int64 {
	return packet.bytePos
}

func (packet *Packet) SetBytePos(newBytePos int64) {
	packet.bytePos = newBytePos
}

func (packet *Packet) GetFileData() []byte {
	return packet.fileData
}

func (packet *Packet) SetFileData(newFileData []byte) {
	packet.fileData = newFileData
}

func (packet *Packet) GetComplete() bool {
	return packet.completed
}

func (packet *Packet) SetCompleted(newCompleted bool) {
	packet.completed = newCompleted
}

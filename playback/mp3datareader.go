package playback

const _bufferSize = 1024 * 16

type MP3DataReader struct {
	data       []byte
	readOffset int
}

func NewMP3DataReader(data []byte) *MP3DataReader {
	n := new(MP3DataReader)
	n.data = data
	return n
}

func (m *MP3DataReader) IsLastChunk() bool {
	if (m.readOffset*_bufferSize)+_bufferSize > len(m.data) {
		return true
	}
	return false
}

func (m *MP3DataReader) BufferNext() []byte {
	size := len(m.data)

	var readStart, readEnd int
	// end bytes
	readStart = _bufferSize * m.readOffset

	if readStart+_bufferSize > size {
		readEnd = size
	} else {
		readEnd = readStart + _bufferSize
		m.readOffset++
	}

	buf := make([]byte, readEnd-readStart)

	copy(buf, m.data[readStart:readEnd])

	return buf
}

func (m *MP3DataReader) Reset() {
	m.readOffset = 0
}

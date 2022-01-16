package validate

import (
	"encoding/binary"
	"time"
)

type Parse struct {
}

func (p *Parse) GetTime(d []byte) time.Time {
	return time.Unix(int64(binary.LittleEndian.Uint64(d[:8])), 0)
}

func (p *Parse) GetByteTime(d []byte) []byte {
	return d[:8]
}

func (p *Parse) GetRndPart(d []byte) []byte {
	return d[8:20]
}

func (p *Parse) GetNonce(d []byte) uint32 {
	return binary.LittleEndian.Uint32(d[20:32])
}

func (p *Parse) GetByteNonce(d []byte) []byte {
	return d[20:32]
}

func (p *Parse) GetHash(d []byte) []byte {
	return d[32:]
}

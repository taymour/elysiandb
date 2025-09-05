package encoding

import "encoding/binary"

func EncodeInt64ToBytes(ts int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(ts))

	return buf
}

func DecodeInt64FromBytes(b []byte) (int64, bool) {
	if len(b) != 8 {
		return 0, false
	}

	return int64(binary.BigEndian.Uint64(b)), true
}

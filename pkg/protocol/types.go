package protocol

import "errors"

// minecraft data-types, see: https://wiki.vg/Protocol#Data_types
type (
	Boolean       bool
	Byte          int8
	UnsignedByte  uint8
	Short         int16
	UnsignedShort uint16

	Int    int32
	Long   int64
	Float  float32
	Double float64

	String     string
	Chat       string
	Identifier string

	// VarInt can be stored in a flexible maximum 6 bytes byte sequence
	VarInt int32
	// VarLong can be stored in a flexible maximum 10 bytes byte sequence
	VarLong int64

	// UUID Encoded as an unsigned 128-bit integer
	UUID uint32

	Position struct {
		x, y, z int64
	}
)

const (
	SegmentBits     int32 = 0x7f
	ContinueBit     byte  = 0x80
	VarIntMaxBytes        = 5
	VarLongMaxBytes       = 10

	MaxStringLength     = 32767
	MaxChatLength       = 262144
	MaxIdentifierLength = 32767
)

func EncodeVarInt(value VarInt) []byte {
	bytes := make([]byte, VarIntMaxBytes)
	index := 0
	num := uint32(value)
	for {
		b := byte(num & uint32(SegmentBits))
		num = num >> 7
		if num != 0 {
			b |= ContinueBit
		}
		bytes[index] = b
		index++
		if num == 0 {
			break
		}
		if index == VarIntMaxBytes {
			panic(errors.New("var int overflow"))
		}
	}
	return bytes[:index]
}

// DecodeVarInt decode a 5-bytes maxim var-int from given byte sequence
func DecodeVarInt(bytes []byte) VarInt {
	var value int32 = 0
	offset, i := 0, 0
	var more byte = 0
	for ; i < VarIntMaxBytes; i++ {
		b := bytes[i]
		// get the 7 bits of value
		v := int32(b) & SegmentBits
		// last bit tells if there's more bytes
		more = b >> 7
		value |= v << offset
		offset += 7
		if more == 0 {
			return VarInt(value)
		}
	}
	// last byte's highest bit is 1, meaning this value is negative
	if more == 1 {
		return VarInt(-value)
	}
	return VarInt(value)
}

func EncodeVarLong(value VarLong) []byte {
	bytes := make([]byte, VarLongMaxBytes)
	index := 0
	num := uint32(value)
	for {
		b := byte(num & uint32(SegmentBits))
		num = num >> 7
		if num != 0 {
			b |= ContinueBit
		}
		bytes[index] = b
		index++
		if num == 0 {
			break
		}
		if index == VarLongMaxBytes {
			panic(errors.New("var long overflow"))
		}
	}
	return bytes[:index]
}

func DecodeVarLong(bytes []byte) VarLong {
	var value int64 = 0
	offset, i := 0, 0
	var more byte = 0
	for ; i < VarLongMaxBytes; i++ {
		b := bytes[i]
		// get the 7 bits of value
		v := int64(b) & int64(SegmentBits)
		// last bit tells if there's more bytes
		more = b >> 7
		value |= v << offset
		offset += 7
		if more == 0 {
			return VarLong(value)
		}
	}
	// last byte's highest bit is 1, meaning this value is negative
	if more == 1 {
		return VarLong(-value)
	}
	return VarLong(value)
}

// DecodePosition decode a long to Position x,y,z
func DecodePosition(value int64) Position {
	x := value >> 38
	y := value & 0xfff
	z := (value >> 12) & 0x3ffffff
	return Position{x, y, z}
}

// EncodePosition encode a Position 's x,y,z to int64
func EncodePosition(pos Position) int64 {
	return (pos.x&0x3ffffff)<<38 | (pos.z&0x3ffffff)<<12 | (pos.y & 0xfff)
}

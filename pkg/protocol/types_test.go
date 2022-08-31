package protocol

import (
	"testing"
)

/*
this test-case is from Minecraft Protocol wiki:https://wiki.vg/Protocol#VarInt_and_VarLong
*/
var (
	varIntTestCases = [][]byte{
		{127},
		{128, 1},
		{255, 1},
		{221, 199, 1},
		{255, 255, 127},
		{255, 255, 255, 255, 7},
		{255, 255, 255, 255, 15},
		{128, 128, 128, 128, 8},
	}
	varIntTestAnswers = []VarInt{127, 128, 255, 25565, 2097151, 2147483647, -1, -2147483648}
	varLongTestCases  = [][]byte{
		{128, 1},
		{255, 1},
		{255, 255, 255, 255, 7},
		{255, 255, 255, 255, 255, 255, 255, 255, 127},
		{255, 255, 255, 255, 255, 255, 255, 255, 255, 1},
		{128, 128, 128, 128, 248, 255, 255, 255, 255, 1},
		{128, 128, 128, 128, 128, 128, 128, 128, 128, 1},
	}
	varLongTestAnswers = []VarLong{128, 255, 2147483647, 9223372036854775807, -1, -2147483648, -9223372036854775808}
)

// TestDecodeVarInt
func TestDecodeVarInt(t *testing.T) {
	for i := 0; i < len(varIntTestCases); i++ {
		if v := DecodeVarInt(varIntTestCases[i]); v != varIntTestAnswers[i] {
			t.Logf("Failed test case: %v, expect: %d, got: %d", varIntTestCases[i], varIntTestAnswers[i], v)
			t.FailNow()
		}
	}
}

// TestDecodeVarLong
func TestDecodeVarLong(t *testing.T) {
	for i := 0; i < len(varLongTestCases); i++ {
		if v := DecodeVarLong(varLongTestCases[i]); v != varLongTestAnswers[i] {
			t.Logf("Failed test case: %v, expect: %d, got: %d", varLongTestCases[i], varLongTestAnswers[i], v)
			t.FailNow()
		}
	}
}

// TestPositionCodec
func TestPositionCodec(t *testing.T) {
	var p1 int64 = 0b0100011000000111011000110001001111101010010010111000001100111111
	pos1 := Position{18357644, 831, 20882616}
	position := DecodePosition(p1)
	if position.x != pos1.x || position.y != pos1.y || position.z != pos1.z {
		t.FailNow()
	}
	enc := EncodePosition(position)
	if enc != p1 {
		t.FailNow()
	}
}

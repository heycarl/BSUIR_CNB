package hammingCode

import (
	"CNB/internal/bitHelper"
	"bytes"
)

const blockSize = 8
const parityBits = 4 // p = ⌈ log2(blockSize + 1) ⌉

func Encode(data []byte) []byte {
	res := bytes.Buffer{}
	for _, b := range data {
		res.Write(blockCalculate(bitHelper.EncodeByte(b)))
	}
	return res.Bytes()
}

func blockCalculate(d []byte) []byte {
	res := make([]byte, blockSize+parityBits)
	//                 0      1  2  3      4  5  6  7
	//         0   1   2  3   4  5  6  7   8  9  10 11
	// result: p1  p2  1  p3  0  1  1  p4  0  0  0  0
	// p1:     ^       ^      ^     ^      ^     ^
	// p2:         ^   ^         ^  ^         ^  ^
	// p3:                ^   ^  ^  ^               ^
	// p4:                             ^   ^  ^  ^  ^
	res[0] = (d[0] + d[1] + d[3] + d[4] + d[6]) % 2 // parity-1
	res[1] = (d[0] + d[2] + d[3] + d[5] + d[6]) % 2 // parity-2
	res[2] = d[0]
	res[3] = (d[1] + d[2] + d[3] + d[7]) % 2 // parity-3
	res[4] = d[1]
	res[5] = d[2]
	res[6] = d[3]
	res[7] = (d[4] + d[5] + d[6] + d[7]) % 2 // parity-4
	res[8] = d[4]
	res[9] = d[5]
	res[10] = d[6]
	res[11] = d[7]
	return res
}

func blockValidate(d []byte) int {
	errorPosition := 0
	if d[0] != (d[0]+d[1]+d[3]+d[4]+d[6])%2 {
		errorPosition += 1
	}
	if d[1] != (d[0]+d[2]+d[3]+d[5]+d[6])%2 {
		errorPosition += 2
	}
	if d[0] != (d[1]+d[2]+d[3]+d[7])%2 {
		errorPosition += 4
	}
	if d[0] != (d[4]+d[5]+d[6]+d[7])%2 {
		errorPosition += 8
	}
	if errorPosition == 0 {
		return -1
	} else {
		return errorPosition
	}
}

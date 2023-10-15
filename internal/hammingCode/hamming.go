package hammingCode

import (
	"CNB/internal/bitHelper"
	"bytes"
	"errors"
	"log"
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

func Decode(data []byte) ([]byte, error) {
	if len(data)%(blockSize+parityBits) != 0 {
		return nil, errors.New("incorrect bits number")
	}
	var res []byte
	for i := 0; i < len(data); i += blockSize + parityBits {
		res = append(res, bitHelper.DecodeByte(blockDecode(data[i:i+blockSize+parityBits])))
	}
	return res, nil
}

func blockCalculate(d []byte) []byte {
	res := make([]byte, blockSize+parityBits)
	log.Printf("Raw block data: %b", d)
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
	log.Printf("Parity bits: %b %b %b %b", res[0], res[1], res[3], res[7])
	return res
}

func extractPayload(d []byte) []byte {
	var r []byte
	r = append(r, d[2])
	r = append(r, d[4:7]...)
	r = append(r, d[8:12]...)
	return r
}

func blockDecode(d []byte) []byte {
	errorPosition := 0
	//                 0      1  2  3      4  5  6  7
	//         0   1   2  3   4  5  6  7   8  9  10 11
	// result: p1  p2  1  p3  0  1  1  p4  0  0  0  0
	// p1:     ^       ^      ^     ^      ^     ^
	// p2:         ^   ^         ^  ^         ^  ^
	// p3:                ^   ^  ^  ^               ^
	// p4:                             ^   ^  ^  ^  ^
	if d[0] != (d[2]+d[4]+d[6]+d[8]+d[10])%2 {
		errorPosition += 1
	}
	if d[1] != (d[2]+d[5]+d[6]+d[9]+d[10])%2 {
		errorPosition += 2
	}
	if d[3] != (d[4]+d[5]+d[6]+d[11])%2 {
		errorPosition += 4
	}
	if d[7] != (d[8]+d[9]+d[10]+d[11])%2 {
		errorPosition += 8
	}
	if errorPosition != 0 {
		log.Printf("Transmission error acquired!\nPosition: %d bit", errorPosition-1)
		bitHelper.FlipBit(&d[errorPosition-1])
	}
	return extractPayload(d)
}

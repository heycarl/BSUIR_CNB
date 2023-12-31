package bitHelper

func Encode(data []byte) []byte {
	var bits []byte

	for _, b := range data {
		for i := 7; i >= 0; i-- {
			bit := b >> uint(i) & 1
			bits = append(bits, bit)
		}
	}

	return bits
}

func EncodeByte(b byte) []byte {
	var bits []byte
	for i := 7; i >= 0; i-- {
		bit := b >> uint(i) & 1
		bits = append(bits, bit)
	}
	return bits
}

func Decode(bits []byte) []byte {
	var data []byte
	var byteAccumulator byte

	for i, bit := range bits {
		byteAccumulator = (byteAccumulator << 1) | bit

		if (i+1)%8 == 0 {
			data = append(data, byteAccumulator)
			byteAccumulator = 0
		}
	}

	return data
}

func DecodeByte(bits []byte) byte {
	var byteAccumulator byte

	for _, bit := range bits {
		byteAccumulator = (byteAccumulator << 1) | bit
	}
	return byteAccumulator
}

func FlipBit(b *byte) {
	if *b == 0 {
		*b = 1
	} else {
		*b = 0
	}
}

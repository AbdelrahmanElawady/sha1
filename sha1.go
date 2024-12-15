package sha1

const uintSize = 32 << (^uint(0) >> 63)

func Hash(message []byte) [20]byte {
	l := len(message)
	message = append(message, 0x80)
	zeros := 0
	if (l+1)%64 > 56 {
		zeros = (64 - ((l+1)%64 - 56))
	} else {
		zeros = 56 - (l+1)%64
	}
	if uintSize == 32 {
		zeros += 4
	}

	for range zeros {
		message = append(message, 0)
	}
	message = append(message, intToBytes(l*8)...)

	h := [5]uint32{
		0x67452301,
		0xEFCDAB89,
		0x98BADCFE,
		0x10325476,
		0xC3D2E1F0,
	}
	chunks := len(message) / 64
	for i := range chunks {
		chunk := message[64*i : 64*(i+1)]
		chunkHash := hashChunk(chunk, h)
		h[0] += chunkHash[0]
		h[1] += chunkHash[1]
		h[2] += chunkHash[2]
		h[3] += chunkHash[3]
		h[4] += chunkHash[4]
	}
	var out [20]byte
	putUint32(out[0:], h[0])
	putUint32(out[4:], h[1])
	putUint32(out[8:], h[2])
	putUint32(out[12:], h[3])
	putUint32(out[16:], h[4])

	return out

}

func hashChunk(chunk []byte, h [5]uint32) [5]uint32 {
	ws := make([]uint32, 80)
	for i := range 80 {
		if i < 16 {
			ws[i] = uint32(chunk[4*i])<<24 | uint32(chunk[4*i+1])<<16 | uint32(chunk[4*i+2])<<8 | uint32(chunk[4*i+3])
		} else {
			a := ws[i-3] ^ ws[i-8] ^ ws[i-14] ^ ws[i-16]
			ws[i] = (a << 1) | (a >> 31)
		}
		var k uint32
		var f uint32
		if i < 20 {
			f = (h[1] & h[2]) | (^h[1] & h[3])
			k = 0x5A827999
		} else if i < 40 {
			f = h[1] ^ h[2] ^ h[3]
			k = 0x6ED9EBA1
		} else if i < 60 {
			f = (h[1] & h[2]) | (h[1] & h[3]) | (h[2] & h[3])
			k = 0x8F1BBCDC
		} else {
			f = h[1] ^ h[2] ^ h[3]
			k = 0xCA62C1D6
		}
		h = [5]uint32{
			((h[0] << 5) | (h[0] >> 27)) + f + h[4] + k + ws[i],
			h[0],
			(h[1] << 30) | (h[1] >> 2),
			h[2],
			h[3],
		}
	}
	return h
}

func putUint32(s []byte, x uint32) {
	bytes := int32ToBytes(int32(x))
	s[0] = bytes[0]
	s[1] = bytes[1]
	s[2] = bytes[2]
	s[3] = bytes[3]
}

func intToBytes(x int) []byte {
	if uintSize == 32 {
		return int32ToBytes(int32(x))
	}
	return int64ToBytes(int64(x))
}

func int32ToBytes(x int32) []byte {
	return []byte{
		byte(x >> 24),
		byte(x >> 16),
		byte(x >> 8),
		byte(x),
	}
}

func int64ToBytes(x int64) []byte {
	return append(
		int32ToBytes(int32(x>>32)),
		int32ToBytes(int32(x))...,
	)
}

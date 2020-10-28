package utils

const (
	blocksize = 64
	s11       = 7
	s12       = 12
	s13       = 17
	s14       = 22
	s21       = 5
	s22       = 9
	s23       = 14
	s24       = 20
	s31       = 4
	s32       = 11
	s33       = 16
	s34       = 23
	s41       = 6
	s42       = 10
	s43       = 15
	s44       = 21
)

type Md5 struct {
	buffer    [blocksize]byte // bytes that didn't fit in last 64 byte chunk
	count     [2]uint32       // 64bit counter for number of bits (lo, hi)
	state     [4]uint32       // digest so far
	finalized bool
	digest    [16]byte //the result
}

func NewMd5() Md5 {
	md5 := Md5{}
	md5.state[0] = 0x67452301
	md5.state[1] = 0xefcdab89
	md5.state[2] = 0x98badcfe
	md5.state[3] = 0x10325476
	md5.finalized = false
	return md5
}

// decodes input (unsigned char) into output (uint4).
func (md5 *Md5) decode(input []byte) []uint32 {
	var output []uint32
	i := 0
	j := 0
	for j < len(input)*4 {
		a := uint32(input[j])
		b := uint32(input[j+1]) << 8
		c := uint32(input[j+2]) << 16
		d := uint32(input[j+3]) << 24
		output = append(output, a|b|c|d)
		i += 1
		j += 4
	}
	return output
}

// encodes input (uint4) into output (unsigned char).
func (md5 *Md5) encode(input []uint32) []byte {
	var output []byte
	i := 0
	j := 0
	for j < len(input)*4 {
		output = append(output, byte(input[i]&0xff))
		output = append(output, byte((input[i]>>8)&0xff))
		output = append(output, byte((input[i]>>16)&0xff))
		output = append(output, byte((input[i]>>24)&0xff))

		i += 1
		j += 4
	}
	return output
}

func (md5 *Md5) f(x, y, z uint32) uint32 {
	return (x & y) | (^x & z)
}

func (md5 *Md5) g(x, y, z uint32) uint32 {
	return (x & z) | (y &^ z)
}

func (md5 *Md5) h(x, y, z uint32) uint32 {
	return x ^ y ^ z
}

func (md5 *Md5) i(x, y, z uint32) uint32 {
	return y ^ (x | ^z)
}

func (md5 *Md5) rotate_left(x uint32, n uint32) uint32 {
	return (x << n) | (x >> (32 - n))
}

func (md5 *Md5) ff(a *uint32, b, c, d, x, s, ac uint32) {
	*a = md5.rotate_left(*a+md5.f(b, c, d)+x+ac, s) + b
}

func (md5 *Md5) gg(a *uint32, b, c, d, x, s, ac uint32) {
	*a = md5.rotate_left(*a+md5.g(b, c, d)+x+ac, s) + b
}

func (md5 *Md5) hh(a *uint32, b, c, d, x, s, ac uint32) {
	*a = md5.rotate_left(*a+md5.h(b, c, d)+x+ac, s) + b
}

func (md5 *Md5) ii(a *uint32, b, c, d, x, s, ac uint32) {
	*a = md5.rotate_left(*a+md5.i(b, c, d)+x+ac, s) + b
}

func (md5 *Md5) transform(block *[blocksize]byte) {
	a := md5.state[0]
	b := md5.state[1]
	c := md5.state[2]
	d := md5.state[3]

	x := md5.decode(block[0:])

	/* Round 1 */
	md5.ff(&a, b, c, d, x[0], s11, 0xd76aa478)  /* 1 */
	md5.ff(&d, a, b, c, x[1], s12, 0xe8c7b756)  /* 2 */
	md5.ff(&c, d, a, b, x[2], s13, 0x242070db)  /* 3 */
	md5.ff(&b, c, d, a, x[3], s14, 0xc1bdceee)  /* 4 */
	md5.ff(&a, b, c, d, x[4], s11, 0xf57c0faf)  /* 5 */
	md5.ff(&d, a, b, c, x[5], s12, 0x4787c62a)  /* 6 */
	md5.ff(&c, d, a, b, x[6], s13, 0xa8304613)  /* 7 */
	md5.ff(&b, c, d, a, x[7], s14, 0xfd469501)  /* 8 */
	md5.ff(&a, b, c, d, x[8], s11, 0x698098d8)  /* 9 */
	md5.ff(&d, a, b, c, x[9], s12, 0x8b44f7af)  /* 10 */
	md5.ff(&c, d, a, b, x[10], s13, 0xffff5bb1) /* 11 */
	md5.ff(&b, c, d, a, x[11], s14, 0x895cd7be) /* 12 */
	md5.ff(&a, b, c, d, x[12], s11, 0x6b901122) /* 13 */
	md5.ff(&d, a, b, c, x[13], s12, 0xfd987193) /* 14 */
	md5.ff(&c, d, a, b, x[14], s13, 0xa679438e) /* 15 */
	md5.ff(&b, c, d, a, x[15], s14, 0x49b40821) /* 16 */

	/* Round 2 */
	md5.gg(&a, b, c, d, x[1], s21, 0xf61e2562)  /* 17 */
	md5.gg(&d, a, b, c, x[6], s22, 0xc040b340)  /* 18 */
	md5.gg(&c, d, a, b, x[11], s23, 0x265e5a51) /* 19 */
	md5.gg(&b, c, d, a, x[0], s24, 0xe9b6c7aa)  /* 20 */
	md5.gg(&a, b, c, d, x[5], s21, 0xd62f105d)  /* 21 */
	md5.gg(&d, a, b, c, x[10], s22, 0x2441453)  /* 22 */
	md5.gg(&c, d, a, b, x[15], s23, 0xd8a1e681) /* 23 */
	md5.gg(&b, c, d, a, x[4], s24, 0xe7d3fbc8)  /* 24 */
	md5.gg(&a, b, c, d, x[9], s21, 0x21e1cde6)  /* 25 */
	md5.gg(&d, a, b, c, x[14], s22, 0xc33707d6) /* 26 */
	md5.gg(&c, d, a, b, x[3], s23, 0xf4d50d87)  /* 27 */
	md5.gg(&b, c, d, a, x[8], s24, 0x455a14ed)  /* 28 */
	md5.gg(&a, b, c, d, x[13], s21, 0xa9e3e905) /* 29 */
	md5.gg(&d, a, b, c, x[2], s22, 0xfcefa3f8)  /* 30 */
	md5.gg(&c, d, a, b, x[7], s23, 0x676f02d9)  /* 31 */
	md5.gg(&b, c, d, a, x[12], s24, 0x8d2a4c8a) /* 32 */

	/* Round 3 */
	md5.hh(&a, b, c, d, x[5], s31, 0xfffa3942)  /* 33 */
	md5.hh(&d, a, b, c, x[8], s32, 0x8771f681)  /* 34 */
	md5.hh(&c, d, a, b, x[11], s33, 0x6d9d6122) /* 35 */
	md5.hh(&b, c, d, a, x[14], s34, 0xfde5380c) /* 36 */
	md5.hh(&a, b, c, d, x[1], s31, 0xa4beea44)  /* 37 */
	md5.hh(&d, a, b, c, x[4], s32, 0x4bdecfa9)  /* 38 */
	md5.hh(&c, d, a, b, x[7], s33, 0xf6bb4b60)  /* 39 */
	md5.hh(&b, c, d, a, x[10], s34, 0xbebfbc70) /* 40 */
	md5.hh(&a, b, c, d, x[13], s31, 0x289b7ec6) /* 41 */
	md5.hh(&d, a, b, c, x[0], s32, 0xeaa127fa)  /* 42 */
	md5.hh(&c, d, a, b, x[3], s33, 0xd4ef3085)  /* 43 */
	md5.hh(&b, c, d, a, x[6], s34, 0x4881d05)   /* 44 */
	md5.hh(&a, b, c, d, x[9], s31, 0xd9d4d039)  /* 45 */
	md5.hh(&d, a, b, c, x[12], s32, 0xe6db99e5) /* 46 */
	md5.hh(&c, d, a, b, x[15], s33, 0x1fa27cf8) /* 47 */
	md5.hh(&b, c, d, a, x[2], s34, 0xc4ac5665)  /* 48 */

	/* Round 4 */
	md5.ii(&a, b, c, d, x[0], s41, 0xf4292244)  /* 49 */
	md5.ii(&d, a, b, c, x[7], s42, 0x432aff97)  /* 50 */
	md5.ii(&c, d, a, b, x[14], s43, 0xab9423a7) /* 51 */
	md5.ii(&b, c, d, a, x[5], s44, 0xfc93a039)  /* 52 */
	md5.ii(&a, b, c, d, x[12], s41, 0x655b59c3) /* 53 */
	md5.ii(&d, a, b, c, x[3], s42, 0x8f0ccc92)  /* 54 */
	md5.ii(&c, d, a, b, x[10], s43, 0xffeff47d) /* 55 */
	md5.ii(&b, c, d, a, x[1], s44, 0x85845dd1)  /* 56 */
	md5.ii(&a, b, c, d, x[8], s41, 0x6fa87e4f)  /* 57 */
	md5.ii(&d, a, b, c, x[15], s42, 0xfe2ce6e0) /* 58 */
	md5.ii(&c, d, a, b, x[6], s43, 0xa3014314)  /* 59 */
	md5.ii(&b, c, d, a, x[13], s44, 0x4e0811a1) /* 60 */
	md5.ii(&a, b, c, d, x[4], s41, 0xf7537e82)  /* 61 */
	md5.ii(&d, a, b, c, x[11], s42, 0xbd3af235) /* 62 */
	md5.ii(&c, d, a, b, x[2], s43, 0x2ad7d2bb)  /* 63 */
	md5.ii(&b, c, d, a, x[9], s44, 0xeb86d391)  /* 64 */

	md5.state[0] += a
	md5.state[1] += b
	md5.state[2] += c
	md5.state[3] += d
}

func (md5 *Md5) Finalize() {
	padding := []byte{0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if md5.finalized == false {
		var input []uint32
		for _, v := range md5.count {
			input = append(input, v)
		}
		bits := md5.encode(input)

		// pad out to 56 mod 64.
		index := md5.count[0] / 8 % 64
		var padLen uint32
		if index < 56 {
			padLen = 56 - index
		} else {
			padLen = 120 - index
		}
		md5.Update(padding, padLen)

		// Append length (before padding)
		md5.Update(bits, 8)

		// Store state in digest
		var state []uint32
		for _, v := range md5.state {
			state = append(state, v)
		}
		result := md5.encode(state)
		for i, v := range result {
			md5.digest[i] = v
		}
		md5.finalized = true
	}
}

func (md5 *Md5) Update(input []byte, len uint32) {
	index := md5.count[0] / 8 % blocksize

	md5.count[0] += (len << 3)
	if md5.count[0] < (len << 3) {
		md5.count[1] += 1
	}
	md5.count[1] += uint32(len >> 29)

	firstpart := 64 - index

	var i uint32
	if len > firstpart {
		first_data := input[:firstpart]
		copy(md5.buffer[index:], first_data)
		md5.transform(&md5.buffer)

		for i := firstpart; i+blocksize <= len; i += blocksize {
			var block [blocksize]byte
			for i := 0; i < blocksize; i += 1 {
				block[i] = input[i]
			}
			//block := input[i : i+blocksize]
			md5.transform(&block)
		}

		index = 0

	} else {
		i = 0
	}

	copy(md5.buffer[index:], input[i:len-i])
}

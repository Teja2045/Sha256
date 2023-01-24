package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
)

var (
	K = [64]uint32{
		0x428a2f98, 0x71374491, 0xb5c0fbcf, 0xe9b5dba5,
		0x3956c25b, 0x59f111f1, 0x923f82a4, 0xab1c5ed5,
		0xd807aa98, 0x12835b01, 0x243185be, 0x550c7dc3,
		0x72be5d74, 0x80deb1fe, 0x9bdc06a7, 0xc19bf174,
		0xe49b69c1, 0xefbe4786, 0x0fc19dc6, 0x240ca1cc,
		0x2de92c6f, 0x4a7484aa, 0x5cb0a9dc, 0x76f988da,
		0x983e5152, 0xa831c66d, 0xb00327c8, 0xbf597fc7,
		0xc6e00bf3, 0xd5a79147, 0x06ca6351, 0x14292967,
		0x27b70a85, 0x2e1b2138, 0x4d2c6dfc, 0x53380d13,
		0x650a7354, 0x766a0abb, 0x81c2c92e, 0x92722c85,
		0xa2bfe8a1, 0xa81a664b, 0xc24b8b70, 0xc76c51a3,
		0xd192e819, 0xd6990624, 0xf40e3585, 0x106aa070,
		0x19a4c116, 0x1e376c08, 0x2748774c, 0x34b0bcb5,
		0x391c0cb3, 0x4ed8aa4a, 0x5b9cca4f, 0x682e6ff3,
		0x748f82ee, 0x78a5636f, 0x84c87814, 0x8cc70208,
		0x90befffa, 0xa4506ceb, 0xbef9a3f7, 0xc67178f2,
	}
	h = [8]uint32{
		0x6a09e667, 0xbb67ae85, 0x3c6ef372, 0xa54ff53a,
		0x510e527f, 0x9b05688c, 0x1f83d9ab, 0x5be0cd19,
	}
)

func PaddedData(input string) []byte {
	data := []byte(input)
	var length [8]byte

	binary.BigEndian.PutUint64(length[:], uint64(len(data)*8))

	data = append(data, 0x01)
	for (len(data)+8)%64 != 0 {
		data = append(data, 0x00)
	}
	data = append(data, length[:]...)
	return data
}

func RotateRight(a uint32, b int) uint32 {
	return (a >> b) | (a << (32 - b))
}

func RotateLeft(a uint32, b int) uint32 {
	return (a << b) | (a >> (32 - b))
}

func Add(a, b uint32) (sum uint32) {
	var cin uint32
	sum = a ^ b ^ cin
	//cout = (a & b) | (cin & (a ^ b))
	return
}

func Usersha256(paddedData []byte) []byte {
	var w [64]uint32

	var s0 uint32
	var s1 uint32
	for i := 0; i < len(paddedData); i += 64 {

		for j := 0; j < 16; j++ {
			w[j] = binary.BigEndian.Uint32(paddedData[i+j*4 : i+j*4+4])
		}
		//fmt.Println(w)
		/*
						s0 = XORXOR(rotr(w[i-15], 7), rotr(w[i-15], 18), shr(w[i-15], 3) )
			      		s1 = XORXOR(rotr(w[i-2], 17), rotr(w[i-2], 19), shr(w[i-2], 10))
			      		w[i] = add(add(add(w[i-16], s0), w[i-7]), s1)
		*/

		for j := 16; j < 64; j++ {
			s0 = RotateRight(w[j-15], 7) ^ RotateRight(w[j-15], 18) ^ (w[j-15] >> 3)
			s1 = RotateRight(w[j-2], 17) ^ RotateRight(w[j-2], 19) ^ (w[j-2] >> 10)
			//w[j] = Add(Add(Add(w[j-16], s0), w[j-7]), s1)
			w[j] = w[j-16] + s0 + w[j-7] + s1
		}
		a := h[0]
		b := h[1]
		c := h[2]
		d := h[3]
		e := h[4]
		f := h[5]
		g := h[6]
		hh := h[7]

		/*
					 for j in range(64):
			      S1 = XORXOR(rotr(e, 6), rotr(e, 11), rotr(e, 25) )
			      ch = XOR(AND(e, f), AND(NOT(e), g))
			      temp1 = add(add(add(add(h, S1), ch), k[j]), w[j])
			      S0 = XORXOR(rotr(a, 2), rotr(a, 13), rotr(a, 22))
			      m = XORXOR(AND(a, b), AND(a, c), AND(b, c))
			      temp2 = add(S0, m)
			      h = g
			      g = f
			      f = e
			      e = add(d, temp1)
			      d = c
			      c = b
			      b = a
			      a = add(temp1, temp2)
		*/
		for j := 0; j < 64; j++ {
			s1 = RotateRight(e, 6) ^ RotateRight(e, 11) ^ RotateRight(e, 25)
			ch := (e & f) ^ ((^e) & g)
			//temp1 := Add(Add(Add(Add(h, s1), ch), K[j]), w[j])
			temp1 := hh + s1 + ch + K[j] + w[j]
			s0 = RotateRight(a, 2) ^ RotateRight(a, 13) ^ RotateRight(a, 22)
			m := (a & b) ^ (a & c) ^ (b & c)
			//temp2 := Add(s0, m)
			temp2 := s0 + m
			hh = g
			g = f
			f = e
			//e = Add(d, temp1)
			e = d + temp1
			d = c
			c = b
			b = a
			//a = Add(temp1, temp2)
			a = temp1 + temp2
		}

		h[0] += a
		h[1] += b
		h[2] += c
		h[3] += d
		h[4] += e
		h[5] += f
		h[6] += g
		h[7] += hh

	}

	hash := make([]byte, 32)
	for i := 0; i < 8; i++ {
		binary.BigEndian.PutUint32(hash[i*4:], h[i])
	}
	return hash
}

func main() {

	inputmessage := "a"
	message := []byte(inputmessage)

	paddedData := PaddedData(inputmessage)
	fmt.Printf("User SHA-256 hash of '%s': %x\n", message, Usersha256(paddedData))

	h := sha256.New()
	h.Write(message)
	hash := h.Sum(nil)
	fmt.Printf("inbuilt SHA-256 hash of '%s': %x\n", message, hash)
}

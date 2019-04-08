package main

type Mirroring int

const (
	horizontal Mirroring = 0
	vertical   Mirroring = 1
)

type Header struct {
	// header的前4个字节应为 4E 45 53 1A
	prgRomSize uint8 // offset 4
	chrRomSize uint8 // offset 5
	// flags6

	flags7  byte
	flags8  byte
	flags9  byte
	flags10 byte
	// 11 ~ 15 未使用
}

type Rom struct {
	header  Header // 16 bytes
	trainer []byte // 0 or 512 bytes
	prgRom  []byte // 16384 * x bytes
	chrRom  []byte // 8192 * y bytes
	instRom []byte // 0 or 8192 bytes
	pRom    []byte
}

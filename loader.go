// NES文件解析
// 参见：https://wiki.nesdev.com/w/index.php/NES_2.0
package main

import (
	"flag"
	"log"
)

type MirroringType uint8
type ConsoleType uint8
type CPUTimingMode uint8

const (
	HORIZONTAL MirroringType = 0
	VERTICAL   MirroringType = 1

	NesOrFc             ConsoleType = 0
	NintendoVsSystem    ConsoleType = 1
	NintendoPlayVoice10 ConsoleType = 2
	ExtendedConsoleType ConsoleType = 3

	RP2C02          CPUTimingMode = 0
	RP2C07          CPUTimingMode = 1
	MULTIPLE_REGION CPUTimingMode = 2
	UMC_6527P       CPUTimingMode = 3
)

type Header struct {
	prgRomSize                uint32
	chrRomSize                uint32
	mirroringType             MirroringType
	isPresentBattery          bool
	isPresentTrainer          bool
	isHardwiredFourScreenMode bool
	consoleType               ConsoleType
	nesRomVersion             uint8
	// INES
	prgRamSize uint32
	mapperNum  uint8
	// NES 2.0
	prgRamShiftCount       uint8
	prgEEPROMShiftCount    uint8
	chrRamSizeShiftCount   uint8
	chrNVRAMSizeShiftCount uint8
	cpuTimingMode          CPUTimingMode
}

type Rom struct {
	header  Header // 16 bytes
	trainer []byte // 0 or 512 bytes
	prgRom  []byte // 16384 * x bytes
	chrRom  []byte // 8192 * y bytes
	instRom []byte // 0 or 8192 bytes
	pRom    []byte
}

func loadGameRom() Rom {
	data := readGameFile()
	log.Printf("rom file is 0x%x bytes\n", len(data.data))

	rom := Rom{header: parseHeader(data.read(16))}

	if rom.header.isPresentTrainer {
		log.Fatal("尚未支持此场景")
	}

	proRomSize := rom.header.prgRomSize
	rom.prgRom = data.read(proRomSize)

	chrRomSize := rom.header.chrRomSize
	rom.chrRom = data.read(chrRomSize)
	return rom
}

// 解析rom文件的头部
func parseHeader(data []byte) Header {
	if !(data[0] == 'N' && data[1] == 'E' && data[2] == 'S' && data[3] == 0x1a) {
		log.Fatal("illegal NES game rom format")
	}
	header := Header{}

	flags6 := data[6]
	if flags6&0x01 == 0x01 {
		header.mirroringType = VERTICAL
	} else {
		header.mirroringType = HORIZONTAL
	}
	if flags6&0x02 == 0x02 {
		header.isPresentBattery = true
	}
	if flags6&0x03 == 0x03 {
		header.isPresentTrainer = true
	}
	if flags6&0x04 == 0x04 {
		header.isHardwiredFourScreenMode = true
	}

	// 获取NES文件格式的版本
	header.nesRomVersion = data[7] & 0xc >> 2

	if header.nesRomVersion == 2 { // NES 2.0
		flags7 := data[7]
		consoleType := flags7 & 0x03
		switch consoleType {
		case 0:
			header.consoleType = NesOrFc
		case 1:
			header.consoleType = NintendoVsSystem
		case 2:
			header.consoleType = NintendoPlayVoice10
		case 3:
			header.consoleType = ExtendedConsoleType
		default:
			log.Fatalf("unkonwn console type %d", consoleType)
		}

		//flags8 := data[8]
		flags9 := data[9]
		prgRomSizeMsb := flags9 & 0xf
		chrRomSizeMsb := flags9 >> 4
		prgRomSizeLsb := data[4]
		chrRomSizeLsb := data[5]

		if prgRomSizeMsb == 0xf {
			log.Fatal("暂时不支持此种情况")
		} else {
			header.prgRomSize = uint32(prgRomSizeMsb)<<8 + uint32(prgRomSizeLsb)
		}

		if chrRomSizeMsb == 0xf {
			log.Fatal("暂时不支持此种情况")
		} else {
			header.chrRomSize = uint32(chrRomSizeMsb)<<8 + uint32(chrRomSizeLsb)
		}
	} else { // INES
		header.prgRomSize = uint32(data[4]) * 16 * 1024 // 以16kb为单位
		header.chrRomSize = uint32(data[5]) * 8 * 1024  // 以8kb为单位

		lowMapperNum := flags6 >> 4
		flags7 := data[7]
		upMapperNum := flags7 >> 4

		header.mapperNum = upMapperNum<<4 + lowMapperNum

		if flags7&0x01 == 0x01 {
			header.consoleType = NintendoVsSystem
		} else if flags7&0x02 == 0x02 {
			header.consoleType = NintendoPlayVoice10
		}

		header.prgRamSize = uint32(data[8]) * 8 * 1024
	}
	return header
}

// 已经读取的rom数据，并保存了当前读取的offset
type RomData struct {
	data   []byte
	offset uint32
}

func (d *RomData) read(size uint32) []byte {
	data := d.data[d.offset : d.offset+size]
	d.offset += size
	log.Printf("total: %d, offset: %d\n", len(d.data), d.offset)
	return data
}

// 根据属性参数读取游戏文件的内容
func readGameFile() RomData {
	romName := flag.Arg(0)
	if romName == "" {
		log.Fatal("no game rom file found")
	}
	return RomData{data: readFile(romName), offset: 0}
}

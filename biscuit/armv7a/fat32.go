package main

import "unsafe"
import "fmt"

const MBR_MAGIC = 0xAA55
const FAT_TYPE = 0xb
const BLKSIZE = 512

type MBR struct {
	partitions [4]partition
	magic      uint32
}

type partition struct {
	bootflag  uint8
	chsbegin  uint32
	typecode  uint8
	chsend    uint32
	lba_begin uint32
	nsectors  uint32
}

type volume_id struct {
	bytes_per_sec   uint16
	sec_per_cluster uint8
	n_rsrv_sectors  uint16
	n_fats          uint8
	sec_per_fat     uint32
	root_cluster    uint32
	signature       uint32
}

type fd struct {
}

var theMBR MBR
var vol_id volume_id

func getmbr() (bool, uint32) {
	lba := uint32(0x0)
	good, data := read_som_sdcard(512, 0x0)
	if !good {
		return false, 0x0
	}
	theMBR.magic = uint32(uint32(data[510])) | (uint32(data[511]) << 8)
	//	copy(mbr_slice[:], data[446:511])
	if theMBR.magic != MBR_MAGIC {
		fmt.Printf("mbr magic is wrong: 0x%x\n", theMBR.magic)
		fmt.Println(unsafe.Sizeof(theMBR))
		fmt.Println(theMBR)
		fmt.Println(data)
		return false, 0x0
	}
	fmt.Printf("MBR 4 partitions:\n")
	for i := 0; i < 4; i++ {
		theMBR.partitions[i].bootflag = (uint8(data[446+(16*i)+0]))
		theMBR.partitions[i].chsbegin = (uint32(data[446+(16*i)+1]) << 0) | (uint32(data[446+(16*i)+2]) << 8) | (uint32(data[446+(16*i)+3]) << 16)
		theMBR.partitions[i].typecode = (uint8(data[446+(16*i)+4]))
		theMBR.partitions[i].chsend = (uint32(data[446+(16*i)+5]) << 0) | (uint32(data[446+(16*i)+6]) << 8) | (uint32(data[446+(16*i)+7]) << 16)
		theMBR.partitions[i].lba_begin = (uint32(data[446+(16*i)+8]) << 0) | (uint32(data[446+(16*i)+9]) << 8) | (uint32(data[446+(16*i)+10]) << 16) | (uint32(data[446+(16*i)+11]) << 24)
		theMBR.partitions[i].nsectors = (uint32(data[446+(16*i)+12]) << 0) | (uint32(data[446+(16*i)+13]) << 8) | (uint32(data[446+(16*i)+14]) << 16) | (uint32(data[446+(16*i)+15]) << 24)
		fmt.Printf("\ttype %d\n", theMBR.partitions[i].typecode)
		if theMBR.partitions[i].typecode == FAT_TYPE {
			fmt.Printf("\t\tfound FAT32 partition with LBA 0x%x\n", theMBR.partitions[i].lba_begin)
			lba = theMBR.partitions[i].lba_begin
		}
	}
	fmt.Println(data[446:])
	return true, lba
}

func getvolumeid(lba uint32) (bool, uint32) {
	good, data := read_som_sdcard(512, lba*BLKSIZE)
	if !good {
		return false, 0x0
	}
	vol_id.bytes_per_sec = (uint16(data[0xb])) | (uint16(data[0xb+1]) << 8)
	vol_id.sec_per_cluster = data[0xd]
	vol_id.n_rsrv_sectors = (uint16(data[0xe]) << 0) | (uint16(data[0xe+1]) << 8)
	vol_id.n_fats = data[0x10]
	vol_id.sec_per_fat = (uint32(data[0x24]) << 0) | (uint32(data[0x24+1]) << 8) | (uint32(data[0x24+2]) << 16) | (uint32(data[0x24+3]) << 24)
	vol_id.root_cluster = (uint32(data[0x2C]) << 0) | (uint32(data[0x2C+1]) << 8) | (uint32(data[0x2C+2]) << 16) | (uint32(data[0x2C+3]) << 24)
	vol_id.signature = (uint32(data[0x1fe]) << 0) | (uint32(data[0x1fe+1]) << 8)
	if vol_id.signature != MBR_MAGIC {
		fmt.Printf("volume id signature is wrong\n")
		fmt.Println(vol_id)
		fmt.Println(data)
		return false, 0x0
	} else {
		fmt.Printf("found FAT32 volume signature\n")
		fmt.Printf("bytes per sector: 0x%d\n", vol_id.bytes_per_sec)
	}
	return true, 0x0
}

//this returns in interface for navigating and modifying the directory structure
func fat32_som_start() (bool, int) {
	if !init_som_sdcard() {
		return false, 0
	}
	good, lba := getmbr()
	if !good {
		return false, 0
	}
	fmt.Printf("found FAT32 partition with LBA 0x%x\n", lba)
	good, _ = getvolumeid(lba)
	if !good {
		return false, 0
	}
	return true, 1
}

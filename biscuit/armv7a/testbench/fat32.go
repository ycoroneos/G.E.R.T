package main

import "unsafe"
import "fmt"

const MBR_MAGIC = 0xAA55
const FAT_TYPE = 0xb
const BLKSIZE = 512
const EOC = 0x0FFFFFF8

type readfunc func(uint32, uint32) (bool, []byte)

var readbytes readfunc

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

type file struct {
	shortname string
	extension string
	attribute uint8
	cluster   uint32
	size      uint32
}

type directory struct {
	files []file
}

type fd struct {
}

var theMBR MBR
var vol_id volume_id
var fat_begin_lba uint32
var cluster_begin_lba uint32
var sectors_per_cluster uint8
var root_dir_first_cluster uint32
var fatcache []uint32

func lba2addr(lba uint32) uint32 {
	return lba * BLKSIZE
}

func cluster2lba(clusternum uint32) uint32 {
	if clusternum < 2 {
		panic("bad cluster number in cluster2lba")
		return 0
	}
	return cluster_begin_lba + (clusternum-2)*uint32(sectors_per_cluster)
}

//this takes a cluster number for a directory as input and it follows the FAT
//chain in order to compile a list of all the files and subdirectories
func readdir_cluster(cluster uint32) (bool, directory) {
	//first make a list of all the clusters that contain data
	var clusters []uint32
	clusters = append(clusters, cluster)
	for fatcache[cluster] < EOC {
		cluster = fatcache[cluster]
		clusters = append(clusters, cluster)
	}
	fmt.Println("clusters for this dir are ", clusters)
	for i := 0; i < len(clusters); i++ {
		discluster := clusters[i]
		addr := lba2addr(cluster2lba(discluster))
		for {
			//skip over long file entries
			good, data := readbytes(32, addr)
			if data[0] == 0xE5 {
				addr += 0x20
				continue
			} else if data[0] == 0x0 {
				fmt.Println("reached end of directory")
				break
			}
			for data[11] == 0xf {
				if !good {
					return false, directory{}
				}
				addr += 0x20
				good, data = readbytes(32, addr)
			}
			fmt.Printf("cluster addr: 0x%x\n", addr)
			fmt.Println(data)
			name := string(data[0:8])
			extension := string(data[8:11])
			attrib := data[11]
			size := (uint32(data[0x1c] << 0)) | (uint32(data[0x1c+1] << 8)) | (uint32(data[0x1c+2] << 16)) | (uint32(data[0x1c+3] << 24))
			if attrib == 0xf {
				fmt.Printf("long filename condensed: ")
			}
			fmt.Println("first file name: ", name, ".", extension, " size: ", size, " ")
			cluster_begin := (uint32(data[0x14]) << 16) | (uint32(data[0x14+1]) << 24) | (uint32(data[0x1a]) << 0) | (uint32(data[0x1a+1]) << 8)
			fmt.Printf("cluster begin number 0x%x\n", cluster_begin)
			fmt.Printf("cluster data begin 0x%x\n", lba2addr(cluster2lba(cluster_begin)))
			//readdir_cluster(cluster_begin)
			addr += 0x20
		}
	}
	return true, directory{}
}

func getmbr() (bool, uint32) {
	lba := uint32(0x0)
	good, data := readbytes(512, 0x0)
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
	good, data := readbytes(512, lba*BLKSIZE)
	if !good {
		return false, 0x0
	}
	vol_id.bytes_per_sec = (uint16(data[0xb])) | (uint16(data[0xb+1]) << 8)
	vol_id.sec_per_cluster = data[13]
	vol_id.n_rsrv_sectors = (uint16(data[0xe]) << 0) | (uint16(data[0xe+1]) << 8)
	vol_id.n_fats = data[0x10]
	vol_id.sec_per_fat = (uint32(data[0x24]) << 0) | (uint32(data[0x24+1]) << 8) | (uint32(data[0x24+2]) << 16) | (uint32(data[0x24+3]) << 24)
	vol_id.root_cluster = (uint32(data[0x2C]) << 0) // | (uint32(data[0x2C+1]) << 8) | (uint32(data[0x2C+2]) << 16) | (uint32(data[0x2C+3]) << 24)
	vol_id.signature = (uint32(data[0x1fe]) << 0) | (uint32(data[0x1fe+1]) << 8)
	if vol_id.signature != MBR_MAGIC {
		fmt.Printf("volume id signature is wrong\n")
		fmt.Println(vol_id)
		fmt.Println(data)
		return false, 0x0
	} else {
		fmt.Printf("found FAT32 volume signature\n")
		fmt.Printf("bytes per sector: 0x%d\n", vol_id.bytes_per_sec)
		fmt.Printf("number of FAT's %d\n", vol_id.n_fats)
		fmt.Printf("sectors per fat %d\n", vol_id.sec_per_fat)
		fmt.Printf("root cluster number is %d\n", vol_id.root_cluster)
	}
	fat_begin_lba = lba + uint32(vol_id.n_rsrv_sectors)
	cluster_begin_lba = lba + uint32(vol_id.n_rsrv_sectors) + (2 * uint32(vol_id.sec_per_fat))
	sectors_per_cluster = vol_id.sec_per_cluster
	root_dir_first_cluster = vol_id.root_cluster
	fmt.Printf("fat addr is 0x%x\n", lba2addr(fat_begin_lba))

	//now read the fat into memory to speed up access times
	fmt.Printf("reading FAT into memory ... ")
	good, fat_bytes := readbytes(vol_id.sec_per_fat*uint32(vol_id.bytes_per_sec), lba2addr(fat_begin_lba))
	if !good {
		return false, 0x0
	}
	nfat_entries := (vol_id.sec_per_fat * uint32(vol_id.bytes_per_sec)) / 4
	fatcache = make([]uint32, nfat_entries, nfat_entries)
	for i := uint32(0); i < nfat_entries; i++ {
		fatcache[i] = (uint32(fat_bytes[(4*i)+0]) << 0) | (uint32(fat_bytes[(4*i)+1]) << 8) | (uint32(fat_bytes[(4*i)+2]) << 16) | (uint32(fat_bytes[(4*i)+3]) << 24)
	}
	fmt.Printf("done! fatcache is %d bytes\n", len(fat_bytes))

	return true, root_dir_first_cluster
}

//this returns an interface for navigating and modifying the directory structure
func fat32_som_start(sdcard_init func() bool, reader_func readfunc) (bool, directory) {
	readbytes = reader_func
	if !sdcard_init() {
		return false, directory{}
	}
	good, lba := getmbr()
	if !good {
		return false, directory{}
	}
	good, root_cluster := getvolumeid(lba)
	if !good {
		return false, directory{}
	}
	good, dir := readdir_cluster(root_cluster)
	if !good {
		return false, directory{}
	}
	return true, dir
}

package main

import "unsafe"
import "fmt"

//sd card type
const (
	CARD_SD   = 0
	CARD_MMC  = 1
	CARD_ESD  = 2
	CARD_EMMC = 3
)

//commands
const (
	CMD0   = 0
	CMD1   = 1
	CMD2   = 2
	CMD3   = 3
	CMD5   = 5
	CMD6   = 6
	ACMD6  = 6
	CMD7   = 7
	CMD8   = 8
	CMD9   = 9
	CMD12  = 12
	CMD13  = 13
	CMD16  = 16
	CMD17  = 17
	CMD18  = 18
	CMD24  = 24
	CMD25  = 25
	CMD26  = 26
	CMD32  = 32
	CMD33  = 33
	CMD35  = 35
	CMD36  = 36
	CMD37  = 37
	CMD38  = 38
	CMD39  = 39
	ACMD41 = 41
	CMD43  = 43
	ACMD51 = 51
	CMD55  = 55
	CMD60  = 60
	CMD61  = 61
	CMD62  = 62
)

//states
const (
	IDLE  = 0
	READY = 1
	IDENT = 2
	STBY  = 3
	TRAN  = 4
	DATA  = 5
	RCV   = 6
	PRG   = 7
	DIS   = 8
)

const (
	USDHC_PORT1        = 0
	USDHC_PORT2        = 1
	USDHC_PORT3        = 2
	USDHC_PORT4        = 3
	USDHC_NUMBER_PORTS = 4
)

type xfer_type_t uint32

const (
	WRITE      = 0
	READ       = 1
	SD_COMMAND = 2
)

type response_format_t uint32

const (
	RESPONSE_NONE          = 0
	RESPONSE_136           = 1
	RESPONSE_48            = 2
	RESPONSE_48_CHECK_BUSY = 3
)

type data_present_select uint32

const (
	DATA_PRESENT_NONE = 0
	DATA_PRESENT      = 1
)

type crc_check_enable uint32
type cmdindex_check_enable uint32
type block_count_enable uint32
type ddren_enable uint32

const (
	DISABLE = 0
	ENABLE  = 1
)

type multi_single_block_select uint32

const (
	SINGLE   = 0
	MULTIPLE = 1
)

//interrupt status
const (
	INTR_BUSY  = 0
	INTR_TC    = 1
	INTR_ERROR = 2
)

type command_t struct {
	command                  uint32
	arg                      uint32
	data_transfer            xfer_type_t
	response_format          response_format_t
	data_present             data_present_select
	crc_check                crc_check_enable
	cmdindex_check           cmdindex_check_enable
	block_count_enable_check block_count_enable
	multi_single_block       multi_single_block_select
	dma_enable               uint32
	acmd12_enable            uint32
	ddren                    ddren_enable
}

type command_response_t struct {
	format   response_format_t
	cmd_rsp0 uint32
	cmd_rsp1 uint32
	cmd_rsp2 uint32
	cmd_rsp3 uint32
}

type usdhc_regs struct {
	DS_ADDR              uint32
	BLT_ATT              uint32
	CMD_ARG              uint32
	CMD_XFR_TYP          uint32
	CMD_RSP0             uint32
	CMD_RSP1             uint32
	CMD_RSP2             uint32
	CMD_RSP3             uint32
	DATA_BUFF_ACC_PORT   uint32
	PRES_STATE           uint32
	PROT_CTRL            uint32
	SYS_CTRL             uint32
	INT_STATUS           uint32
	INT_STATUS_EN        uint32
	INT_SIGNAL_EN        uint32
	AUTOCMD12_ERR_STAT   uint32
	HOST_CTRL_CAP        uint32
	WTMK_LVL             uint32
	MIX_CTRL             uint32
	FORCE_EVENT          uint32
	ADMA_ERR_STATUS      uint32
	ADMA_SYS_ADDR        uint32
	DLL_CTRL             uint32
	DLL_STATUS           uint32
	CLK_TUNE_CTRL_STATUS uint32
	VEND_SPEC            uint32
	MMC_BOOT             uint32
	VEND_SPEC2           uint32
}

type usdhc_inst_t struct {
	//reg_base uint32 //register base address
	regbase  *usdhc_regs //register base address
	adma_ptr uint32      //ADMA buffer address
	//void (*isr) (void);         //interrupt service routine
	isr       uintptr //unused for now
	rca       uint16  //relative card address
	addr_mode uint8   //addressing mode
	intr_id   uint8   //interrupt ID
	status    uint8   //interrupt status
}

type sdhc_freq_t uint32

const (
	OPERATING_FREQ      = 0
	IDENTIFICATION_FREQ = 1
	HS_FREQ             = 2
)

//sdcard regs base
const (
	REGS_USDHC1_BASE = 0x02190000
	REGS_USDHC2_BASE = 0x02194000
	REGS_USDHC3_BASE = 0x02198000
	REGS_USDHC4_BASE = 0x0219C000
)

const (
	USDHC_ADMA_BUFFER1 = 0x00907000
	USDHC_ADMA_BUFFER2 = 0x00908000
	USDHC_ADMA_BUFFER3 = 0x00909000
	USDHC_ADMA_BUFFER4 = 0x0090A000
)

//sdcard interrupts
const (
	IMX_INT_USDHC1 = 54
	IMX_INT_USDHC2 = 55
	IMX_INT_USDHC3 = 56
	IMX_INT_USDHC4 = 57
)

//sdcard isrs
const ()

//sdcard randoms
const (
	ESDHC_ONE_BIT_SUPPORT    = 0x0
	ESDHC_LITTLE_ENDIAN_MODE = 0x2
)

const card_detect_test_en = 0

//statically initialize the sdcards
var usdhc_device = [...]usdhc_inst_t{
	usdhc_inst_t{((*usdhc_regs)(unsafe.Pointer(uintptr(REGS_USDHC1_BASE)))), USDHC_ADMA_BUFFER1, 0, 0, 0, IMX_INT_USDHC1, 1},
	usdhc_inst_t{((*usdhc_regs)(unsafe.Pointer(uintptr(REGS_USDHC2_BASE)))), USDHC_ADMA_BUFFER2, 0, 0, 0, IMX_INT_USDHC2, 1},
	usdhc_inst_t{((*usdhc_regs)(unsafe.Pointer(uintptr(REGS_USDHC3_BASE)))), USDHC_ADMA_BUFFER3, 0, 0, 0, IMX_INT_USDHC3, 1},
	usdhc_inst_t{((*usdhc_regs)(unsafe.Pointer(uintptr(REGS_USDHC4_BASE)))), USDHC_ADMA_BUFFER4, 0, 0, 0, IMX_INT_USDHC4, 1},
}

//go:nosplit
func usdhc_card_detected(instance uint32) bool {
	switch instance {
	case 0:
		fmt.Println("there is no such thing as sdhc0. this is confusing")
		return false
	case 1:
		fmt.Println("attempting to detect card on sdhc1")
		return true
	case 2:
		fmt.Println("there is no sd card to be detected on sdhc3")
		return false
	case 3:
		fmt.Println("sdhc3 is not yet supported")
		return false
	default:
		fmt.Printf("unrecognized sdhc%d\n", instance)
		return false
	}
}

//go:nosplit
func host_init(instance uint32) {
	/* Enable Clock Gating */
	// i think uboot does this

	/* IOMUX Configuration */
	usdhc_iomux_config(instance)

	//unclear what this does
	//usdhc_gpio_config(instance)
}

//go:nosplit
func usdhc_set_data_transfer_width(instance uint32, data_width int) {
	usdhc_device[instance-1].regbase.PROT_CTRL |= (uint32(data_width) & 0x3) << 1
}

//go:nosplit
func usdhc_set_endianness(instance uint32, endian_mode int) {
	usdhc_device[instance-1].regbase.PROT_CTRL |= (uint32(endian_mode) & 0x3) << 4
}

//go:nosplit
func host_set_bus_width(instance uint32, bus_width int) {
	usdhc_set_data_transfer_width(instance, bus_width)
}

//go:nosplit
func host_reset(instance uint32, bus_width int, endian_mode int) {
	/* Reset the eSDHC by writing 1 to RSTA bit of SYSCTRL Register */
	usdhc_device[instance-1].regbase.SYS_CTRL |= 0x1 << 24

	/* Wait until RSTA field cleared */
	for (usdhc_device[instance-1].regbase.SYS_CTRL & 0x1 << 24) > 0 {
	}

	/* Set default bus width to eSDHC */
	host_set_bus_width(instance, bus_width)

	/* Set Endianness of eSDHC */
	usdhc_set_endianness(instance, endian_mode)
}

//init sd card according to state machine in the firmware guide
//go:nosplit
func card_init(instance, bus_width uint32) int {
	init_status := -1
	/* Initialize uSDHC Controller */
	host_init(instance)

	/* Software Reset to Interface Controller */
	host_reset(instance, ESDHC_ONE_BIT_SUPPORT, ESDHC_LITTLE_ENDIAN_MODE)

	//card detection
	if card_detect_test_en > 0 {
		if usdhc_card_detected(instance) == false {
			fmt.Printf("SD card detected, but not instance %d\n", instance)
			return -1
		} else {
			fmt.Printf("SD card %d detected\n", instance)
		}
	}

	//write protect detection
	if write_protect_test_en > 0 {
		if usdhc_write_protected(instance) == true {
			fmt.Printf("Card on SD%d is write protected.\n", instance)
			return -1
		} else {
			fmt.Printf("Card on SD%d is not write protected.\n", instance)
		}
	}

	/* Initialize interrupt */
	if card_init_interrupt(instance) < 0 {
		fmt.Printf("Interrupt initialize failed.\n")
		return -1
	}

	/* Enable Identification Frequency */
	host_cfg_clock(instance, IDENTIFICATION_FREQ)

	/* Send Init 80 Clock */
	host_init_active(instance)

	fmt.Printf("Reset card.\n")

	/* Issue Software Reset to card */
	if card_software_reset(instance) < 0 {
		return -1
	}

	/* SD Voltage Validation */
	if sd_voltage_validation(instance) > 0 {
		fmt.Printf("SD voltage validation passed.\n")

		/* SD Initialization */
		init_status = sd_init(instance, bus_width)
	} else if mmc_voltage_validation(instance) > 0 { /* MMC Voltage Validation */
		fmt.Printf("MMC voltage validation passed.\n")

		/* MMC Initialization */
		init_status = mmc_init(instance, bus_width)
	}

	return init_status
}

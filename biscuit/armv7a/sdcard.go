package main

/*
* This entire damn driver was directly ported from the iMX6 bare metal sdk
* for use in the Biscuit embedded toolkit
 */

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
	BLK_ATT              uint32
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
	rca       uint32  //relative card address
	addr_mode uint8   //addressing mode
	intr_id   uint8   //interrupt ID
	status    uint8   //interrupt status
}

type sdhc_freq_t uint32

const (
	OPERATING_FREQ         = 0
	IDENTIFICATION_FREQ    = 1
	HS_FREQ                = 2
	ESDHC_IDENT_DVS        = 8
	ESDHC_IDENT_SDCLKFS    = 0x20
	ESDHC_OPERT_DVS        = 0x3
	ESDHC_OPERT_SDCLKFS    = 0x1
	ESDHC_HS_DVS           = 0x1
	ESDHC_HS_SDCLKFS       = 0x1
	ESDHC_SYSCTL_DTOCV_VAL = 0xE
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
	ESDHC_CLEAR_INTERRUPT              = 0x117F01FF
	ESDHC_INTERRUPT_ENABLE             = 0x007F013F
	ESDHC_STATUS_END_CMD_RESP_TIME_MSK = 0x100F0001
	IMX_INT_USDHC1                     = 54
	IMX_INT_USDHC2                     = 55
	IMX_INT_USDHC3                     = 56
	IMX_INT_USDHC4                     = 57
)

//sdcard isrs
const ()

//sdcard randoms
const (
	ESDHC_ONE_BIT_SUPPORT             = 0x0
	ESDHC_LITTLE_ENDIAN_MODE          = 0x2
	ESDHC_CIHB_CHK_COUNT              = 10
	ESDHC_CDIHB_CHK_COUNT             = 100
	BM_USDHC_PROT_CTRL_DMASEL         = 0x00000300
	ESDHC_MIXER_CTRL_CMD_MASK         = 0xFFFFFFC0
	BP_USDHC_MIX_CTRL_DMAEN           = 0
	BP_USDHC_MIX_CTRL_BCEN            = 1
	BP_USDHC_MIX_CTRL_AC12EN          = 2
	BP_USDHC_MIX_CTRL_DDR_EN          = 3
	BP_USDHC_MIX_CTRL_DTDSEL          = 4
	BP_USDHC_MIX_CTRL_MSBSEL          = 5
	BP_USDHC_CMD_XFR_TYP_RSPTYP       = 16
	BP_USDHC_CMD_XFR_TYP_CCCEN        = 19
	BP_USDHC_CMD_XFR_TYP_CICEN        = 20
	BP_USDHC_CMD_XFR_TYP_DPSEL        = 21
	BP_USDHC_CMD_XFR_TYP_CMDINX       = 24
	ESDHC_OPER_TIMEOUT_COUNT          = 10000
	SD_IF_CMD_ARG_COUNT               = 2
	SD_IF_HV_COND_ARG                 = 0x000001AA
	SD_IF_LV_COND_ARG                 = 0x000002AA
	SD_OCR_VALUE_HV_HC                = 0x40ff8000
	SD_OCR_VALUE_LV_HC                = 0x40000080
	SD_OCR_VALUE_HV_LC                = 0x00ff8000
	SD_OCR_VALUE_COUNT                = 3
	SD_VOLT_VALID_COUNT               = 3000
	CARD_BUSY_BIT                     = 0x80000000
	SD_OCR_HC_RES                     = 0x40000000
	SECT_MODE                         = 1
	BYTE_MODE                         = 0
	RCA_SHIFT                         = 16
	SD_R1_STATUS_APP_CMD_MSK          = 0x20
	MMC_HV_HC_OCR_VALUE               = 0x40FF8000
	MMC_VOLT_VALID_COUNT              = 3000
	MMC_OCR_HC_BIT_MASK               = 0x60000000
	MMC_OCR_HC_RESP_VAL               = 0x40000000
	BLK_LEN                           = 512
	ESDHC_FIFO_LENGTH                 = 0x80
	ESDHC_BLKATTR_WML_BLOCK           = 0x80
	ESDHC_STATUS_END_DATA_RSP_TC_MASK = 0x00700002
)

var sd_if_cmd_arg = [...]uint32{
	SD_IF_HV_COND_ARG,
	SD_IF_LV_COND_ARG,
}

var sd_ocr_value = [...]uint32{
	SD_OCR_VALUE_HV_HC,
	SD_OCR_VALUE_LV_HC,
	SD_OCR_VALUE_HV_LC,
}

const card_detect_test_en = 0
const write_protect_test_en = 0
const SDHC_ADMA_mode = 0

//statically initialize the sdcards
var usdhc_device = [...]usdhc_inst_t{
	usdhc_inst_t{((*usdhc_regs)(unsafe.Pointer(uintptr(REGS_USDHC1_BASE)))), USDHC_ADMA_BUFFER1, 0, 0, 0, IMX_INT_USDHC1, 1},
	usdhc_inst_t{((*usdhc_regs)(unsafe.Pointer(uintptr(REGS_USDHC2_BASE)))), USDHC_ADMA_BUFFER2, 0, 0, 0, IMX_INT_USDHC2, 1},
	usdhc_inst_t{((*usdhc_regs)(unsafe.Pointer(uintptr(REGS_USDHC3_BASE)))), USDHC_ADMA_BUFFER3, 0, 0, 0, IMX_INT_USDHC3, 1},
	usdhc_inst_t{((*usdhc_regs)(unsafe.Pointer(uintptr(REGS_USDHC4_BASE)))), USDHC_ADMA_BUFFER4, 0, 0, 0, IMX_INT_USDHC4, 1},
}

//go:nosplit
func usdhc_check_transfer(instance uint32) int {
	status := -1
	if instance == 0 {
		fmt.Printf("host_reset instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]

	if ((dev.regbase.INT_STATUS & (0x1 << 1)) > 0) &&
		(dev.regbase.INT_STATUS&(0x1<<20) <= 0) &&
		(dev.regbase.INT_STATUS&(0x1<<21) <= 0) {
		status = 1
	} else {
		fmt.Printf("Error transfer status: 0x%x\n", dev.regbase.INT_STATUS)
	}

	return status
}

//go:nosplit
func host_data_read(instance uint32, dst_ptr *[]uint32, length int, wml int) int {
	var idx int
	var itr int
	var loop int
	dst_spot := 0
	if instance == 0 {
		fmt.Printf("host_reset instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]

	/* Clear Interrupt */
	dev.regbase.INT_STATUS = 0xFFFFFFFF

	/* Enable Interrupt */
	dev.regbase.INT_STATUS_EN = 0xFFFFFFFF
	//HW_USDHC_INT_STATUS_EN(instance).U |= ESDHC_INTERRUPT_ENABLE;

	/* Read data to dst_ptr */
	loop = length / (4 * wml)
	for idx = 0; idx < loop; idx++ {
		/* Wait until buffer ready */
		fmt.Printf("\t wait for buffer ready\n")
		for (dev.regbase.PRES_STATE & (0x1 << 11)) == 0 {
		}

		/* Read from FIFO watermark words */
		for itr = 0; itr < wml; itr++ {
			data := dev.regbase.DATA_BUFF_ACC_PORT
			fmt.Printf("\tread %x\n", data)
			(*dst_ptr)[dst_spot] = data
			dst_spot += 1
		}
	}

	/* Read left data that not WML aligned */
	loop = (length % (4 * wml)) / 4
	if loop != 0 {
		/* Wait until buffer ready */
		fmt.Printf("\twait for buffer ready\n")
		for (dev.regbase.PRES_STATE & (0x1 << 11)) == 0 {
		}

		/* Read the left to destination buffer */
		for itr = 0; itr < loop; itr++ {
			data := dev.regbase.DATA_BUFF_ACC_PORT
			fmt.Printf("\tread %x\n", data)
			(*dst_ptr)[dst_spot] = data
			dst_spot += 1
		}

		/* Clear FIFO */
		fmt.Printf("\tclear fifo\n")
		for ; itr < wml; itr++ {
			idx = int(dev.regbase.DATA_BUFF_ACC_PORT)
		}
	}

	/* Wait until transfer complete */
	excount := 0
	fmt.Printf("\twait for transfer complete\n")
	for (dev.regbase.INT_STATUS & ESDHC_STATUS_END_DATA_RSP_TC_MASK) <= 0 {
		if (dev.regbase.INT_STATUS & 0x20) > 0 {
			extra := dev.regbase.DATA_BUFF_ACC_PORT
			excount += 1
			fmt.Printf("found extra data %x\n", extra)
		}
	}
	fmt.Printf("total bytes read: %d\n", 4*(excount+dst_spot))

	/* Check if error happened */
	return usdhc_check_transfer(instance)
}

//go:nosplit
func host_cfg_block(instance uint32, blk_len int, nob int, wml int) {
	if instance == 0 {
		fmt.Printf("host_reset instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	/* Number of blocks and block length */
	dev.regbase.BLK_ATT = (uint32(nob) << 16) | uint32(blk_len)

	/* Watermark level - for DMA transfer */
	dev.regbase.WTMK_LVL = uint32(wml)
}

//go:nosplit
func host_clear_fifo(instance uint32) {
	if instance == 0 {
		fmt.Printf("host_reset instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	var trash uint32
	/* If data present in Rx FIFO */
	if (dev.regbase.INT_STATUS & (0x1 << 5)) > 0 {
		/* Read from FIFO until empty */
		for idx := 0; idx < ESDHC_FIFO_LENGTH; idx++ {
			trash = dev.regbase.DATA_BUFF_ACC_PORT
		}
	}

	trash = trash + 2
	/* Maybe not necessary */
	dev.regbase.INT_STATUS |= 0x1 << 5
}

//go:nosplit
func card_set_blklen(instance uint32, len int) int {
	var cmd command_t
	status := -1

	/* Configure CMD16 */
	card_cmd_config(&cmd, CMD16, len, READ, RESPONSE_48, DATA_PRESENT_NONE, 1, 1)

	//fmt.Printf("Send CMD16.\n")

	/* Send CMD16 */
	if host_send_cmd(instance, &cmd) > 0 {
		status = 1
	}

	return status
}

//go:nosplit
func card_data_read(instance uint32, length int, offset uint32) (int, []byte) {

	//var port int
	var sector int
	var cmd command_t

	/* Get uSDHC port according to instance */
	if instance == 0 {
		fmt.Printf("host_reset instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]

	//fmt.Printf("card_data_read: Read 0x%x bytes from SD%d offset 0x%x\n", length, port+1, offset)

	/* Get sector number */
	if SDHC_ADMA_mode == 1 {
		/* For DMA mode, length should be sector aligned */
		if (length % BLK_LEN) != 0 {
			length = length + BLK_LEN - (length % BLK_LEN)
		}

		sector = length / BLK_LEN
	} else {
		/* For PIO mode, not neccesary */
		sector = length / BLK_LEN

		if (length % BLK_LEN) != 0 {
			sector += 1
		}
	}

	/* Offset should be sectors */
	if dev.addr_mode == SECT_MODE {
		offset = offset / BLK_LEN
	}

	/* Set block length to card */
	if card_set_blklen(instance, BLK_LEN) < 0 {
		fmt.Printf("Fail to set block length to card in reading sector %d.\n", offset/BLK_LEN)
		return -1, nil
	}

	/* Clear Rx FIFO */
	host_clear_fifo(instance)

	/* Configure block length/number and watermark */
	host_cfg_block(instance, BLK_LEN, sector, ESDHC_BLKATTR_WML_BLOCK)

	/* If DMA mode enabled, configure BD chain */
	if SDHC_ADMA_mode > 0 {
		fmt.Println("ADMA transfer is not supported. why did this happen?")
		//host_setup_adma(instance, dst_ptr, length)
		//card_buffer_flush(dst_ptr, length)
	}

	/* Use CMD18 for multi-block read */
	card_cmd_config(&cmd, CMD18, int(offset), READ, RESPONSE_48, DATA_PRESENT, 1, 1)

	//fmt.Printf("card_data_read: Send CMD18.\n")

	/* make slice */
	out_data := make([]byte, length)
	/* Send CMD18 */
	if host_send_cmd(instance, &cmd) < 0 {
		fmt.Printf("Fail to send CMD18.\n")
		return -1, nil
	} else {
		/* In polling IO mode, manually read data from Rx FIFO */
		if SDHC_ADMA_mode <= 0 {
			//fmt.Printf("Non-DMA mode, read data from FIFO.\n")
			//fmt.Printf("Block length: %d\n", BLK_LEN)
			//fmt.Printf("Nbytes to read: %d\n", length)
			//fmt.Printf("Nsectors: %d\n", sector)

			/* Clear Interrupt */
			dev.regbase.INT_STATUS = 0xFFFFFFFF

			/* Enable Interrupt */
			dev.regbase.INT_STATUS_EN = 0xFFFFFFFF

			for count := 0; count < length; count += 4 {
				for (dev.regbase.PRES_STATE & (0x1 << 11)) == 0 {
				}
				data := dev.regbase.DATA_BUFF_ACC_PORT
				for i := 0; i < 4; i++ {
					if (count + i) > length {
						break
					}
					out_data[count+i] = byte((data >> (8 * uint32(i))) & 0xFF)
				}
			}

			/* clear out the fifo */
			//var trash uint32
			for (dev.regbase.INT_STATUS & ESDHC_STATUS_END_DATA_RSP_TC_MASK) <= 0 {
				if (dev.regbase.INT_STATUS & 0x20) > 0 {
					dev.regbase.DATA_BUFF_ACC_PORT |= 0
					//trash = dev.regbase.DATA_BUFF_ACC_PORT
					//fmt.Printf("read trash %x\n", trash)
				}
			}
			//trash += 2

			if usdhc_check_transfer(instance) < 0 {
				return -1, nil
			}
			//if host_data_read(instance, dst_ptr, length, ESDHC_BLKATTR_WML_BLOCK) < 0 {
			//	fmt.Printf("Fail to read data from card.\n")
			//	return -1
			//}
		}
	}

	//fmt.Printf("card_data_read: Data read successful.\n")

	return 1, out_data

}

//go:nosplit
func mmc_switch(instance uint32, arg uint32) int {
	var cmd command_t
	status := -1

	/* Configure MMC Switch Command */
	card_cmd_config(&cmd, CMD6, int(arg), READ, RESPONSE_48, DATA_PRESENT_NONE, 1, 1)

	fmt.Printf("Send CMD6.\n")

	/* Send CMD6 */
	if host_send_cmd(instance, &cmd) > 0 {
		status = card_trans_status(instance)
	}
	return status
}

//go:nosplit
func MMC_SWITCH_SETBW_ARG(bus_width uint32) uint32 {
	return uint32(0x03b70001 | ((bus_width >> 2) << 8))
}

//go:nosplit
func mmc_set_bus_width(instance uint32, bus_width int) int {
	return mmc_switch(instance, MMC_SWITCH_SETBW_ARG(uint32(bus_width)))
}

//go:nosplit
func mmc_set_rca(instance uint32) int {
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	var cmd command_t
	var response command_response_t
	status := -1

	/* Set RCA to ONE */
	dev.rca = 1

	/* Configure CMD3 */
	card_cmd_config(&cmd, CMD3, int(dev.rca<<RCA_SHIFT), READ, RESPONSE_48, DATA_PRESENT_NONE, 1, 1)

	fmt.Printf("Send CMD3.\n")

	/* Send CMD3 */
	if host_send_cmd(instance, &cmd) > 0 {
		response.format = RESPONSE_48
		host_read_response(instance, &response)

		/* Check the IDENT card state */
		card_state := int((response.cmd_rsp0 & 0x1e00) >> 9)

		if card_state == IDENT {
			status = 1
		}
	}

	return status
}

//go:nosplit
func mmc_init(instance uint32, bus_width int) int {
	status := -1

	//fmt.Printf("Get CID.\n")

	/* Get CID */
	if card_get_cid(instance) > 0 {
		//fmt.Printf("Set RCA.\n")

		/* Set RCA */
		if mmc_set_rca(instance) > 0 {
			/* Check Card Type here */
			//fmt.Printf("Set operating frequency.\n")

			/* Switch to Operating Frequency */
			host_cfg_clock(instance, OPERATING_FREQ)

			//fmt.Printf("Enter transfer state.\n")

			/* Enter Transfer State */
			if card_enter_trans(instance) > 0 {
				//fmt.Printf("Set bus width.\n")

				/* Set Card Bus Width */
				if mmc_set_bus_width(instance, bus_width) > 0 {
					/* Set Host Bus Width */
					host_set_bus_width(instance, bus_width)

					/* Set High Speed Here */
					{
						status = 1
					}
				}
			}
		}
	}

	return status
}

//go:nosplit
func mmc_voltage_validation(instance uint32) int {
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	var cmd command_t
	var response command_response_t
	ocr_val := MMC_HV_HC_OCR_VALUE
	status := -1

	for count := 0; count < MMC_VOLT_VALID_COUNT && status < 0; {
		/* Configure CMD1 */
		card_cmd_config(&cmd, CMD1, ocr_val, READ, RESPONSE_48, DATA_PRESENT_NONE, 0, 0)

		/* Send CMD1 */
		if host_send_cmd(instance, &cmd) < 0 {
			fmt.Printf("Send CMD1 failed.\n")
			break
		} else {
			/* Check Response */
			response.format = RESPONSE_48
			host_read_response(instance, &response)

			/* Check Busy Bit Cleared or NOT */
			if (response.cmd_rsp0 & CARD_BUSY_BIT) > 0 {
				/* Check Address Mode */
				if (response.cmd_rsp0 & MMC_OCR_HC_BIT_MASK) == MMC_OCR_HC_RESP_VAL {
					dev.addr_mode = SECT_MODE
				} else {
					dev.addr_mode = BYTE_MODE
				}

				status = 1
			} else {
				count += 1
				//hal_delay_us(MMC_VOLT_VALID_DELAY);
			}
		}
	}

	return status
}

//go:nosplit
func sd_set_bus_width(instance uint32, bus_width int) int {
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	var cmd command_t
	var response command_response_t
	status := -1

	address := dev.rca << RCA_SHIFT

	/* Check Bus Width */
	if (bus_width != 4) && (bus_width != 1) {
		fmt.Printf("Invalid bus_width: %d\n", bus_width)
		return status
	}

	/* Configure CMD55 */
	card_cmd_config(&cmd, CMD55, int(address), READ, RESPONSE_48, DATA_PRESENT_NONE, 1, 1)

	//fmt.Printf("Send CMD55.\n")

	/* Send ACMD6 */
	if host_send_cmd(instance, &cmd) > 0 {
		/* Check Response of Application Command */
		response.format = RESPONSE_48
		host_read_response(instance, &response)

		if (response.cmd_rsp0 & SD_R1_STATUS_APP_CMD_MSK) > 0 {
			bus_width = bus_width >> 1

			/* Configure ACMD6 */
			card_cmd_config(&cmd, ACMD6, bus_width, READ, RESPONSE_48, DATA_PRESENT_NONE, 1, 1)

			//fmt.Printf("Send CMD6.\n")

			/* Send ACMD6 */
			if host_send_cmd(instance, &cmd) > 0 {
				status = 1
			}
		}
	}

	return status
}

//go:nosplit
func card_trans_status(instance uint32) int {
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	var cmd command_t
	var response command_response_t
	status := -1

	/* Get RCA */
	card_address := dev.rca << RCA_SHIFT

	/* Configure CMD13 */
	card_cmd_config(&cmd, CMD13, int(card_address), READ, RESPONSE_48, DATA_PRESENT_NONE, 1, 1)

	//fmt.Printf("Send CMD13.\n")

	/* Send CMD13 */
	if host_send_cmd(instance, &cmd) > 0 {
		/* Get Response */
		response.format = RESPONSE_48
		host_read_response(instance, &response)

		/* Read card state from response */
		card_state := int((response.cmd_rsp0 & 0x1e00) >> 9)
		if card_state == TRAN {
			status = 1
		}
	}

	return status
}

//go:nosplit
func card_enter_trans(instance uint32) int {
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	var cmd command_t
	status := -1

	/* Get RCA */
	card_address := dev.rca << RCA_SHIFT

	/* Configure CMD7 */
	card_cmd_config(&cmd, CMD7, int(card_address), READ, RESPONSE_48_CHECK_BUSY, DATA_PRESENT_NONE, 1, 1)

	//fmt.Printf("Send CMD7.\n")

	/* Send CMD7 */
	if host_send_cmd(instance, &cmd) > 0 {
		/* Check if the card in TRAN state */
		if card_trans_status(instance) > 0 {
			status = 1
		}
	}

	return status
}

//go:nosplit
func sd_get_rca(instance uint32) int {
	var cmd command_t
	card_state := 0
	status := -1
	var response command_response_t
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]

	/* Configure CMD3 */
	card_cmd_config(&cmd, CMD3, 0, READ, RESPONSE_48, DATA_PRESENT_NONE, 1, 1)

	//fmt.Printf("Send CMD3.\n")

	/* Send CMD3 */
	if host_send_cmd(instance, &cmd) > 0 {
		response.format = RESPONSE_48
		host_read_response(instance, &response)

		/* Set RCA to Value Read */
		//fmt.Printf("RCA is 0x%x\n", response.cmd_rsp0)
		dev.rca = response.cmd_rsp0 >> RCA_SHIFT
		//fmt.Printf("saved RCA is 0x%x\n", dev.rca)

		/* Check the IDENT card state */
		card_state = int((response.cmd_rsp0 & 0x1e00) >> 9)

		if card_state == IDENT {
			status = 1
		}
	}

	return status
}

//go:nosplit
func card_get_cid(instance uint32) int {
	var cmd command_t
	status := -1
	var response command_response_t

	/* Configure CMD2 */
	card_cmd_config(&cmd, CMD2, 0, READ, RESPONSE_136, DATA_PRESENT_NONE, 1, 0)

	//fmt.Printf("Send CMD2.\n")

	/* Send CMD2 */
	if host_send_cmd(instance, &cmd) > 0 {
		response.format = RESPONSE_136
		host_read_response(instance, &response)

		/* No Need to Save CID */

		status = 1
	}

	return status
}

//go:nosplit
func sd_init(instance uint32, bus_width int) int {
	status := -1

	//fmt.Printf("Get CID.\n")

	/* Read CID */
	if card_get_cid(instance) > 0 {
		//fmt.Printf("Get RCA.\n")

		/* Obtain RCA */
		if sd_get_rca(instance) > 0 {
			//fmt.Printf("Set operating frequency.\n")

			/* Enable Operating Freq */
			host_cfg_clock(instance, OPERATING_FREQ)

			if bus_width == 8 {
				bus_width = 4
			}

			//fmt.Printf("Enter transfer state.\n")

			/* Enter Transfer State */
			if card_enter_trans(instance) > 0 {
				//fmt.Printf("Set bus width.\n")

				/* Set Bus Width for SD card */
				if sd_set_bus_width(instance, bus_width) > 0 {
					/* Set Bus Width for Controller */
					host_set_bus_width(instance, bus_width)

					/* Set High Speed Here */
					{
						status = 1
					}
				}
			}
		}
	}

	return status
}

//go:nosplit
func host_read_response(instance uint32, response *command_response_t) {
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	/* Read response from registers */
	response.cmd_rsp0 = dev.regbase.CMD_RSP0
	response.cmd_rsp1 = dev.regbase.CMD_RSP1
	response.cmd_rsp2 = dev.regbase.CMD_RSP2
	response.cmd_rsp3 = dev.regbase.CMD_RSP3
}

//go:nosplit
func sd_voltage_validation(instance uint32) int {
	var cmd command_t
	var response command_response_t
	status := -1
	ocr_value := 0
	card_usable := -1
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]

	//fmt.Printf("Send CMD8.\n")

	for loop := uint32(0); loop < SD_IF_CMD_ARG_COUNT; {
		card_cmd_config(&cmd, CMD8, int(sd_if_cmd_arg[loop]), READ, RESPONSE_48, DATA_PRESENT_NONE, 1, 1)
		if host_send_cmd(instance, &cmd) < 0 {
			loop += 1

			if (loop >= SD_IF_CMD_ARG_COUNT) && (loop < SD_OCR_VALUE_COUNT) {
				/* Card is of SD-1.x spec with LC */
				ocr_value = int(sd_ocr_value[loop])
				card_usable = 1
			}
		} else {
			/* Card is supporting SD spec version >= 2.0 */
			response.format = RESPONSE_48
			host_read_response(instance, &response)

			/* Check if response lies in the expected volatge range */
			if (response.cmd_rsp0 & sd_if_cmd_arg[loop]) == sd_if_cmd_arg[loop] {
				ocr_value = int(sd_ocr_value[loop])
				card_usable = 1

				break
			} else {
				ocr_value = 0
				card_usable = -1

				break
			}
		}
	}

	if card_usable < 0 {
		return status
	}

	//fmt.Printf("Send ACMD41.\n")

	for loop := 0; loop < SD_VOLT_VALID_COUNT && status < 0; {
		card_cmd_config(&cmd, CMD55, 0, READ, RESPONSE_48, DATA_PRESENT_NONE, 1, 1)

		if host_send_cmd(instance, &cmd) < 0 {
			fmt.Printf("Send CMD55 failed.\n")
			break
		} else {
			card_cmd_config(&cmd, ACMD41, ocr_value, READ, RESPONSE_48, DATA_PRESENT_NONE, 0, 0)

			if host_send_cmd(instance, &cmd) < 0 {
				fmt.Printf("Send ACMD41 failed.\n")
				break
			} else {
				/* Check Response */
				response.format = RESPONSE_48
				host_read_response(instance, &response)

				/* Check Busy Bit Cleared or NOT */
				if (response.cmd_rsp0 & CARD_BUSY_BIT) > 0 {
					/* Check card is HC or LC from card response */
					if (response.cmd_rsp0 & SD_OCR_HC_RES) == SD_OCR_HC_RES {
						dev.addr_mode = SECT_MODE
					} else {
						dev.addr_mode = BYTE_MODE
					}

					status = 1
				} else {
					loop += 1
					//hal_delay_us(SD_VOLT_VALID_DELAY);
				}
			}
		}
	}

	return status
}

//go:nosplit
func usdhc_check_response(instance uint32) int {
	if instance == 0 {
		fmt.Printf("\thost_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	status := -1
	if (((dev.regbase.INT_STATUS & 0x1) > 0) ||
		((dev.regbase.MMC_BOOT & (0x1 << 6)) > 0)) &&
		((dev.regbase.INT_STATUS & (0x1 << 16)) == 0) &&
		((dev.regbase.INT_STATUS & (0x1 << 17)) == 0) &&
		((dev.regbase.INT_STATUS & (0x1 << 19)) == 0) &&
		((dev.regbase.INT_STATUS & (0x1 << 18)) == 0) {
		status = 1
		//fmt.Printf("\tresponse all good\n")
		//fmt.Printf("\tError status: 0x%x\n", dev.regbase.INT_STATUS)
		//fmt.Printf("\tMMC_BOOT = 0x%x\n", dev.regbase.MMC_BOOT)
	} else {
		fmt.Printf("\tresponse bad\n")
		fmt.Printf("\tError status: 0x%x\n", dev.regbase.INT_STATUS)
		fmt.Printf("\tMMC_BOOT = 0x%x\n", dev.regbase.MMC_BOOT)

		/* Clear CIHB and CDIHB status */
		if ((dev.regbase.PRES_STATE & 0x1) > 0) ||
			((dev.regbase.PRES_STATE & 0x2) > 0) {
			dev.regbase.SYS_CTRL |= 1 << 24
		}
	}

	return status
}

//go:nosplit
func usdhc_wait_end_cmd_resp_intr(instance uint32) {
	if instance == 0 {
		fmt.Printf("instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	count := 0

	for (dev.regbase.INT_STATUS & ESDHC_STATUS_END_CMD_RESP_TIME_MSK) <= 0 {
		if count == ESDHC_OPER_TIMEOUT_COUNT {
			fmt.Printf("Command timeout. Nothing happened at all\n")
			break
		}

		count += 1
		for junk := 0; junk < 10000; junk++ {
		}
		//hal_delay_us(ESDHC_STATUS_CHK_TIMEOUT);
	}
}

//go:nosplit
func usdhc_cmd_cfg(instance uint32, cmd *command_t) {
	var consist_status uint32
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]

	/* Write Command Argument in Command Argument Register */
	dev.regbase.CMD_ARG = cmd.arg

	/* Clear the DMAS field */
	//fmt.Printf("PROT_CTRL addr 0x%x\n", &dev.regbase.PROT_CTRL)
	//fmt.Printf("old PROT_CTRL 0x%x\n", dev.regbase.PROT_CTRL)
	dev.regbase.PROT_CTRL &= ^(uint32(BM_USDHC_PROT_CTRL_DMASEL))
	//fmt.Printf("new PROT_CTRL 0x%x\n", dev.regbase.PROT_CTRL)

	/* If ADMA mode enabled and command with DMA, enable ADMA2 */
	//    if ((cmd->dma_enable == TRUE) && (read_usdhc_adma_mode() == TRUE)) {
	//    	BW_USDHC_PROT_CTRL_DMASEL(instance, ESDHC_PRTCTL_ADMA2_VAL);
	//    }

	/* Keep bit fields other than command setting intact */
	consist_status = dev.regbase.MIX_CTRL & ESDHC_MIXER_CTRL_CMD_MASK

	dev.regbase.MIX_CTRL = (consist_status |
		(uint32(cmd.dma_enable) << BP_USDHC_MIX_CTRL_DMAEN) |
		(uint32(cmd.block_count_enable_check) << BP_USDHC_MIX_CTRL_BCEN) |
		(uint32(cmd.acmd12_enable) << BP_USDHC_MIX_CTRL_AC12EN) |
		(uint32(cmd.ddren) << BP_USDHC_MIX_CTRL_DDR_EN) |
		(uint32(cmd.data_transfer) << BP_USDHC_MIX_CTRL_DTDSEL) |
		(uint32(cmd.multi_single_block) << BP_USDHC_MIX_CTRL_MSBSEL))

	dev.regbase.CMD_XFR_TYP = ((uint32(cmd.response_format) << BP_USDHC_CMD_XFR_TYP_RSPTYP) |
		(uint32(cmd.crc_check) << BP_USDHC_CMD_XFR_TYP_CCCEN) |
		(uint32(cmd.cmdindex_check) << BP_USDHC_CMD_XFR_TYP_CICEN) |
		(uint32(cmd.data_present) << BP_USDHC_CMD_XFR_TYP_DPSEL) |
		(uint32(cmd.command) << BP_USDHC_CMD_XFR_TYP_CMDINX))
}

//go:nosplit
func usdhc_wait_cmd_data_lines(instance uint32, data_present int) int {
	count := 0
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	/* Wait for release of CMD line */
	for (dev.regbase.PRES_STATE & 0x1) > 0 {
		if count == ESDHC_CIHB_CHK_COUNT {
			fmt.Printf("wait_cmd_data_lines cmd timeout\n")
			return -1
		}
		count += 1
		//hal_delay_us(ESDHC_STATUS_CHK_TIMEOUT);
	}

	/* If data present with command, wait for release of DATA lines */
	if data_present == DATA_PRESENT {
		count = 0
		for (dev.regbase.PRES_STATE & 0x2) > 0 {
			if count == ESDHC_CDIHB_CHK_COUNT {
				fmt.Printf("wait_cmd_data_lines data timeout\n")
				return -1
			}

			count += 1
			//hal_delay_us(ESDHC_STATUS_CHK_TIMEOUT);
		}
	}

	return 1
}

//go:nosplit
func host_send_cmd(instance uint32, cmd *command_t) int {
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	/* Clear Interrupt status register */
	//dev.regbase.INT_STATUS = ESDHC_CLEAR_INTERRUPT
	dev.regbase.INT_STATUS = 0xFFFFFFFF

	/* Enable Interrupt */
	//dev.regbase.INT_STATUS_EN |= ESDHC_INTERRUPT_ENABLE
	dev.regbase.INT_STATUS_EN = 0xFFFFFFFF

	/* Wait for CMD/DATA lines to be free */
	if usdhc_wait_cmd_data_lines(instance, int(cmd.data_present)) < 0 {
		fmt.Printf("\tData/Command lines busy.\n")
		return -1
	}

	/* Clear interrupt status */
	//dev.regbase.INT_STATUS |= ESDHC_STATUS_END_CMD_RESP_TIME_MSK
	dev.regbase.INT_STATUS = 0xFFFFFFFF
	//fmt.Printf("\tINT_STATUS reads %x after clear\n", dev.regbase.INT_STATUS)
	//fmt.Printf("INT_STATUS_EN reads %x\n", dev.regbase.INT_STATUS_EN)

	/* Enable interrupt when sending DMA commands */
	//    if ((read_usdhc_intr_mode() > 0) && (cmd->dma_enable>0)) {
	//        int idx = card_get_port(instance);
	//
	//        /* Set interrupt flag to busy */
	//        usdhc_device[idx].status = INTR_BUSY;
	//
	//        /* Enable uSDHC interrupt */
	//        HW_USDHC_INT_SIGNAL_EN_WR(instance, ESDHC_STATUS_END_DATA_RSP_TC_MASK);
	//    }

	/* Configure Command */
	//fmt.Printf("\tSending command 0x%x\n", cmd.command)
	usdhc_cmd_cfg(instance, cmd)

	/* If DMA Enabled */
	if cmd.dma_enable > 0 {
		fmt.Printf("\twhy is DMA enabled? It is not supported\n")
		//        /* Return in interrupt mode */
		//        if (read_usdhc_intr_mode() == TRUE) {
		//            return SUCCESS;
		//        }
		//
		//        usdhc_wait_end_cmd_resp_dma_intr(instance);
	} else {
		//fmt.Printf("\twait for response... ")
		usdhc_wait_end_cmd_resp_intr(instance)
		//fmt.Println("got it!")
	}

	/* Mask all interrupts */
	dev.regbase.INT_SIGNAL_EN = 0

	/* Check if an error occured */
	return usdhc_check_response(instance)
}

//go:nosplit
func card_cmd_config(cmd *command_t, index int, argument int, transfer xfer_type_t,
	format response_format_t, data data_present_select,
	crc crc_check_enable, cmdindex cmdindex_check_enable) {
	cmd.command = uint32(index)
	cmd.arg = uint32(argument)
	cmd.data_transfer = transfer
	cmd.response_format = format
	cmd.data_present = data
	cmd.crc_check = crc
	cmd.cmdindex_check = cmdindex
	cmd.dma_enable = 0
	cmd.block_count_enable_check = 0
	cmd.multi_single_block = SINGLE
	cmd.acmd12_enable = 0
	cmd.ddren = 0

	/* Multi Block R/W Setting */
	if (CMD18 == index) || (CMD25 == index) {
		if SDHC_ADMA_mode > 0 {
			cmd.dma_enable = 1
		}

		cmd.block_count_enable_check = 1
		cmd.multi_single_block = MULTIPLE
		cmd.acmd12_enable = 1
	}
}

//go:nosplit
func card_software_reset(instance uint32) int {
	var cmd command_t
	response := -1

	/* Configure CMD0 */
	card_cmd_config(&cmd, CMD0, 0, READ, RESPONSE_NONE, DATA_PRESENT_NONE, 0, 0)

	//fmt.Printf("Send CMD0.\n")

	/* Issue CMD0 to Card */
	if host_send_cmd(instance, &cmd) > 0 {
		response = 1
	}

	return response
}

//go:nosplit
func host_init_active(instance uint32) {
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]

	//fmt.Println("send 80 clock ticks")
	/* Send 80 clock ticks for card to power up */
	dev.regbase.SYS_CTRL |= 1 << 27

	//fmt.Println("wait for INITA to clear")
	/* Wait until INITA field cleared */
	for (dev.regbase.SYS_CTRL & (1 << 27)) > 0 {
	}
}

//go:nosplit
func host_cfg_clock(instance uint32, frequency int) {
	/* Clear SDCLKEN bit, this bit is reserved in Rev D*/
	//esdhc_base->system_control &= ~ESDHC_SYSCTL_SDCLKEN_MASK;
	if instance == 0 {
		fmt.Printf("host_cfg_clock instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]

	//fmt.Printf("wait for clock stable\n")
	/* Wait until clock stable */
	for (dev.regbase.PRES_STATE & (0x1 << 3)) == 0 {
	}

	//fmt.Printf("clear things\n")
	/* Clear DTOCV, SDCLKFS, DVFS bits */
	dev.regbase.SYS_CTRL &= ^(uint32(0xF<<16) | uint32(0xFF<<8) | uint32(0xF<<4))

	//fmt.Printf("wait for clock stable\n")
	/* Wait until clock stable */
	for (dev.regbase.PRES_STATE & (0x1 << 3)) == 0 {
	}

	//fmt.Printf("set frequency dividers\n")
	/* Set frequency dividers */
	if frequency == IDENTIFICATION_FREQ {
		dev.regbase.SYS_CTRL |= ESDHC_IDENT_DVS << 4
		dev.regbase.SYS_CTRL |= ESDHC_IDENT_SDCLKFS << 8
	} else if frequency == OPERATING_FREQ {
		dev.regbase.SYS_CTRL |= ESDHC_OPERT_DVS << 4
		dev.regbase.SYS_CTRL |= ESDHC_OPERT_SDCLKFS << 8
	} else {
		dev.regbase.SYS_CTRL |= ESDHC_HS_DVS << 4
		dev.regbase.SYS_CTRL |= ESDHC_HS_SDCLKFS << 8
	}

	//fmt.Printf("wait for clock stable\n")
	/* Wait until clock stable */
	for (dev.regbase.PRES_STATE & (0x1 << 3)) == 0 {
	}

	//fmt.Printf("set some bit\n")
	/* Set DTOCV bit */
	dev.regbase.SYS_CTRL |= ESDHC_SYSCTL_DTOCV_VAL << 16
}

//i dont use interrupts yet for this
//go:nosplit
func card_init_interrupt(instance uint32) int {
	return 1
}

//go:nosplit
func usdhc_write_protected(instance uint32) bool {
	return true
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
	//usdhc_iomux_config(instance)

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
	if instance == 0 {
		fmt.Printf("host_reset instance 0 is not valid\n")
	}
	dev := &usdhc_device[instance-1]
	/* Reset the eSDHC by writing 1 to RSTA bit of SYSCTRL Register */
	//fmt.Printf("sent reset\n")
	dev.regbase.SYS_CTRL |= 0x1 << 24

	//fmt.Printf("wait for rsta clear\n")
	//fmt.Printf("sys_ctrl is %x\n", dev.regbase.SYS_CTRL)
	/* Wait until RSTA field cleared */
	for (dev.regbase.SYS_CTRL & (0x1 << 24)) > 0 {
	}

	//fmt.Printf("set default bus width\n")
	/* Set default bus width to eSDHC */
	host_set_bus_width(instance, bus_width)

	//fmt.Printf("set endian mode\n")
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
			//fmt.Printf("SD card %d detected\n", instance)
		}
	}

	//write protect detection
	if write_protect_test_en > 0 {
		if usdhc_write_protected(instance) == true {
			fmt.Printf("Card on SD%d is write protected.\n", instance)
			return -1
		} else {
			//fmt.Printf("Card on SD%d is not write protected.\n", instance)
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

	//fmt.Printf("\n\nReset card:\n")

	/* Issue Software Reset to card */
	if card_software_reset(instance) < 0 {
		return -1
	}

	//fmt.Printf("\n\nvalidate sdcard voltage:\n")
	/* SD Voltage Validation */
	if sd_voltage_validation(instance) > 0 {
		//fmt.Printf("This is an sd card\n")
		//fmt.Printf("SD voltage validation passed.\n")

		/* SD Initialization */
		init_status = sd_init(instance, int(bus_width))
	} else if mmc_voltage_validation(instance) > 0 { /* MMC Voltage Validation */
		//fmt.Printf("This is actually an mmc\n")
		//fmt.Printf("MMC voltage validation passed.\n")

		/* MMC Initialization */
		init_status = mmc_init(instance, int(bus_width))
	}
	return init_status
}

////////////////////////////////////////////////
//useful functions

//init the som sdcard at port 3 with 4bit bus width
//go:nosplit
func init_som_sdcard() bool {
	return card_init(3, 4) > 0
}

///read length bytes at byte offset from the start of the sdcard
//go:nosplit
func read_som_sdcard(length int, offset uint32) (bool, []byte) {
	val, data := card_data_read(uint32(3), length, offset)
	status := val > 0
	return status, data
}

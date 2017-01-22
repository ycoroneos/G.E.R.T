package embedded

/*
* Translated again from the imx6 platform sdk
*
 */

type SPI_dev uint32

const (
	DEV_SPI1 = 1
	DEV_SPI2 = 2
	DEV_SPI3 = 3
	DEV_SPI4 = 4
	DEV_SPI5 = 5
)

type spi_desc struct {
	channel    uint8
	mode       uint8
	ss_pol     uint8
	sclk_pol   uint8
	sclk_phase uint8
	pre_div    uint8
	post_div   uint8
}

//ecspi_open(dev SPI_dev, param_ecspi_t spi_desc) int {
//    // Configure IO signals
//    ecspi_iomux_config(dev)
//
//    // Ungate the module clock.
//		//Put 0x3F into CCM_CCGR1
//    clock_gating_config(REGS_ECSPI_BASE(dev), CLOCK_ON);
//
//    // Configure eCSPI registers
//    ecspi_configure(dev, param);
//
//    return SUCCESS;
//}

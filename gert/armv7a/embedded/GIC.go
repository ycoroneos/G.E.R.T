package embedded

import "unsafe"
import "runtime"
import "fmt"

/*
*This is an excellent guide that won't help at all
* you have to actually read ~3 500pg documents. you're screwed
*ftp://ftp.altera.com/up/pub/Altera_Material/14.0/Tutorials/Using_GIC.pdf
 */

const MPCOREBASE uint32 = 0xa00000

type GIC_distributor_map struct {
	distributor_control_register                     uint32
	interrupt_controller_type_register               uint32
	distributor_implmementer_identification_register uint32
	reserved1                                        [29]uint32
	interrupt_security_registers                     [8]uint32
	reserved5                                        [24]uint32
	interrupt_set_enable_registers                   [32]uint32
	interrupt_clear_enable_registers                 [32]uint32
	interrupt_set_pending_registers                  [32]uint32
	interrupt_clear_pending_registers                [32]uint32
	active_bit_registers                             [32]uint32
	reserved2                                        [32]uint32
	interrupt_priority_registers                     [255 * 4]uint8
	reserved3                                        uint32
	interrupt_processor_targets_registers            [255 * 4]uint8
	reserved4                                        uint32
	interrupt_configuration_registers                [64]uint32
	implementation_defined_registers                 [64]uint32
	reserved6                                        [64]uint32
	software_generated_interrupt_register            uint32
	reserved7                                        [51]uint32
	peripheral_id4                                   uint32
	peripheral_id5                                   uint32
	peripheral_id6                                   uint32
	peripheral_id7                                   uint32
	peripheral_id0                                   uint32
	peripheral_id1                                   uint32
	peripheral_id2                                   uint32
	peripheral_id3                                   uint32
	component_id0                                    uint32
	component_id1                                    uint32
	component_id2                                    uint32
	component_id3                                    uint32
}

type GIC_cpu_map struct {
	cpu_interface_control_register       uint32
	interrupt_priority_mask_register     uint32
	binary_point_register                uint32
	interrupt_acknowledge_register       uint32
	end_of_interrupt_register            uint32
	running_priority_register            uint32
	highest_pending_interrupt_register   uint32
	aliased_binary_point_register        uint32
	reserved1                            [8]uint32
	implementation_defined_registers     [36]uint32
	reserved2                            [11]uint32
	cpu_interface_dentification_register uint32
}

var gic_distributor *GIC_distributor_map = (*GIC_distributor_map)(unsafe.Pointer(uintptr(MPCOREBASE + 0x1000)))
var gic_cpu *GIC_cpu_map = (*GIC_cpu_map)(unsafe.Pointer(uintptr(MPCOREBASE + 0x100)))

//priorty register is
//0x080a00104

func GIC_init(checks bool) {
	runtime.DisableIRQ()
	GIC_mask_all() //mask all interrupts
	fmt.Println("interrupts off")
	/*
	   On the Freescale i.MX6Q this should produce PA:00A00000
	*/
	fmt.Printf("\tARM registers Address: %x\r\n", runtime.Getmpcorebase())
	/*
	   On the Freescale i.MX6Q this should produce:3901243B
	*/
	fmt.Printf("\tcpu_interface_dentification_register: %x\r\n", gic_cpu.cpu_interface_dentification_register)
	if checks && gic_cpu.cpu_interface_dentification_register != 0x3901243B {
		panic("cpu identification register mismatch")
	}

	/*
	   On the Freescale i.MX6Q this should produce: 0000000D 000000F0 00000005 000000B1
	*/
	fmt.Printf("\tComponent ID: %x", gic_distributor.component_id0)
	if gic_distributor.component_id0 != 0xD {
		panic("cpu distributor component ID mismatch")
	}
	fmt.Printf(" %x", gic_distributor.component_id1)
	if gic_distributor.component_id1 != 0xF0 {
		panic("cpu distributor component ID mismatch")
	}
	fmt.Printf(" %x", gic_distributor.component_id2)
	if gic_distributor.component_id2 != 0x5 {
		panic("cpu distributor component ID mismatch")
	}
	fmt.Printf(" %x \r\n", gic_distributor.component_id3)
	if gic_distributor.component_id3 != 0xB1 {
		panic("cpu distributor component ID mismatch")
	}
	/*
	   Enable all interrupt levels (0 = high, 0xFF = low)
	*/
	gic_cpu.interrupt_priority_mask_register = 0xFF

	/*
	   Enable the GIC (both secure and insecure interrupts)
	*/
	gic_cpu.cpu_interface_control_register = 0x03       // enable everything
	gic_distributor.distributor_control_register = 0x03 // enable everything

	/*
	   Enable the GPT interrupt
	*/
	//enable_interrupt(87, 2)
	//enable_interrupt(0, 1)
	runtime.EnableIRQ()
	GIC_unmask_all() // unmask all interrupts
	runtime.DMB()
	fmt.Println("interrupts on")

}

func Enable_interrupt(num uint32, cpunum uint32) {
	gic_distributor.interrupt_priority_registers[num] = 0                        // highest priority
	gic_distributor.interrupt_security_registers[num/32] &= ^(1 << (num & 0x1F)) // disable security
	//gic_distributor.interrupt_processor_targets_registers[num] |= uint8(cpunum & 0xFF) // send to CPU 0
	gic_distributor.interrupt_processor_targets_registers[num] = uint8((0x1 << cpunum) & 0xFF) // send to CPU 0
	gic_distributor.interrupt_set_enable_registers[num/32] = 1 << (num & 0x1F)                 // enable the interrupt
}

func Sgi(num uint32, cpus uint32) {
	gic_distributor.software_generated_interrupt_register = ((cpus & 0xFF) << 16) | (num & 0xF)
}

func GIC_mask_all() {
	gic_cpu.interrupt_priority_mask_register = 0x0
}

func GIC_unmask_all() {
	gic_cpu.interrupt_priority_mask_register = 0xFF
}

/*
//for triggering sgi's
func sgi(id uint32) {
	*(*uint32)(unsafe.Pointer(uintptr(ICDSGIR))) = id
}

func clrint() {
	*(*uint32)(unsafe.Pointer(uintptr(ICDICPR))) = 0
}
*/

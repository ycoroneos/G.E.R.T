`timescale 1ns / 1ps
//////////////////////////////////////////////////////////////////////////////////
// Company: PDOS
// Engineer: yance
// 
// Create Date: 01/26/2017 09:21:44 PM
// Design Name: 
// Module Name: toplevel
// Project Name: 
// Target Devices: 
// Tool Versions: 
// Description: 
// 
// Dependencies: 
// 
// Revision:
// Revision 0.01 - File Created
// Additional Comments:
// 
//////////////////////////////////////////////////////////////////////////////////


module toplevel(
input wire sysclk,
input wire [7:0] sw,
output wire [7:0] led,
output wire [7:0] ja
    );
    wire reset=sw[0];
    wire counterclk;
    wire clk100mhz=sysclk;
    wire signal=ja[0];
    assign ja[1] = counterclk;
    //100MHz/1000000 = 100Hz
    clkdiv #(.LOGLENGTH(31), .COUNTVAL(1000000)) counterclkgen(.inclk(clk100mhz), .reset(reset), .newclk(counterclk));
    counter #(.COUNTVAL(100)) counter0(.clk(counterclk), .reset(reset), .out(signal), .done(led[0]));
endmodule

`timescale 1ns / 1ps
//////////////////////////////////////////////////////////////////////////////////
// Company: PDOS
// Engineer: yance
// 
// Create Date: 01/26/2017 09:24:14 PM
// Design Name: 
// Module Name: clkdiv
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
module clkdiv(
input wire inclk, reset,
output wire newclk
    );
    parameter LOGLENGTH=31;
    parameter COUNTVAL=100000;
    reg [LOGLENGTH:0] count=0;
    reg output_reg=0;
    assign newclk=output_reg;
    always @(posedge inclk)
    begin
    if (reset)
        begin
        count <=0;
        output_reg<=0;
        end
    else
        begin
        if (count==COUNTVAL)
            begin
            count<=0;
            output_reg<=~output_reg;
            end
        else
            count<=count+1;
        end
    end
endmodule

`timescale 1ns / 1ps
//////////////////////////////////////////////////////////////////////////////////
// Company: PDOS
// Engineer: yance
// 
// Create Date: 01/26/2017 09:32:44 PM
// Design Name: 
// Module Name: counter
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


module counter(
input wire clk,
input wire reset,
output wire out,
output reg done
    );
    parameter COUNTVAL=100000;
    reg [31:0] curcount = COUNTVAL*2;
    reg level=0;
    reg done=0;
    always @(posedge clk)
    begin
        //reset logic
        if (reset)
            begin
            curcount <= COUNTVAL*2;
            level<=0;
            done<=0;
            end
        else
            //flip on every positive edge
            if (curcount == 0)
            begin
                done<=1;
            end
            else
            begin
            curcount<=curcount-1;
            level <= ~level;
            end
    end
    assign out=level;
endmodule

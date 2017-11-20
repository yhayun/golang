package main

import (
	"net"
	"fmt"
	"time"
)
const UDP_SIZE = 1500
const MPEG2TS_PACKET_LENGTH = 188

//server code:   //extractstrean() code
func main() {
	addr, _ := net.ResolveUDPAddr("udp", ":8888")
	sock, _ := net.ListenUDP("udp", addr)
	fmt.Println("working on UDP");
	i := 0
	for {
		i++
		buf := make([]byte, UDP_SIZE)
		rlen, _, err := sock.ReadFromUDP(buf)
		if err != nil {
			fmt.Println(err)
		}
		defer sock.Close();
		fmt.Println(string(buf[0:rlen]));
		fmt.Println(i);
		//go handlePacket(buf, rlen)
	}
}



type  Mpeg2TSource struct {
	tsPacket Mpeg2TSPacket
	programPID int // = 0
	endFlag bool //= false;
	detectFlag bool //= false;
	previousTime int64 // = -1
}



func (src Mpeg2TSource) extractStream(packet []byte) {
	for offset:=0; offset < len(packet); offset += MPEG2TS_PACKET_LENGTH {
		src.tsPacket.fromBytes(packet,offset, src.programPID)
		//if(src.detectFlag && ((System.currentTimeMillis() - src.previousTime) > 1000)) {



	}
}





//for (int offset = 0; offset < packet.getLength(); offset += MPEG2TS_PACKET_LENGTH) {
//
//tsPacket.fromBytes(packet.getData(), offset, programPID);
//
//if(previousTime != -1) {
//if(detectFlag && ((System.currentTimeMillis() - previousTime) > 1000)) {
//log.info("Timeout detected, now looking for new PID");
//programPID = 0;
//pmtFrame = new PMTFrame();
//detectFlag = false;
//}
//}
package main

import (
	"net"
	"fmt"
)
const UDP_SIZE = 1500;

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



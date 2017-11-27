package main

import (
	"net"
	"fmt"
	"time"
	//"log"
)
const UDP_SIZE = 1500
const MPEG2TS_PACKET_LENGTH = 188
const MAX_UDP_PACKET_SIZE = 1500;


func timeMilisecs() long {
	return long(time.Now().UnixNano()%1e6/1e3)
}

type DatagramPacket struct {
	buf []byte
	length int
}

func (dgPacket* DatagramPacket) GetData() []byte{
	return dgPacket.buf
}

func (dgPacket* DatagramPacket) GetLength() int{
	return dgPacket.length
}

func (dgPacket* DatagramPacket) SetData(data []byte) {
	dgPacket.buf = data
}

type Mpeg2TSSource struct {
	socket DatagramSocket
	reTransmitFlag bool
	videoFrames frameQueue
    programPID int
 	udpSource UdpSource
	outputPort int
	socketReceiveBufferSize int
}








//
//type  UdpSource struct {
//	tsPacket Mpeg2TSPacket
//	programPID int // = 0
//	endFlag bool //= false;
//	detectFlag bool //= false;
//	previousTime int64 // = -1
//	pmtFrame PMTFrame //
//}

type UdpSource struct {
	Mpeg2TSSource
	socket int //DatagramSocket
	outputSocket int //DatagramSocket

	endFlag bool
	detectFlag bool
	videoPID int

	videoTSParser Mpeg2TSParser
	tsPacket Mpeg2TSPacket

	pmtFrame PMTFrame
	previousTime long
	streamsMap map[int]StreamInfo //HashMap<Integer, StreamInfo>
	packet DatagramPacket

}
func (u* UdpSource) extractStreams(packet DatagramPacket) {

	for offset := 0; offset < packet.GetLength(); offset += MPEG2TS_PACKET_LENGTH {

		u.tsPacket.FromBytes(packet.GetData(), offset, u.programPID);

		if u.previousTime != -1 {
			if(u.detectFlag && ((timeMilisecs() - u.previousTime) > 1000)) {
				//log.info("Timeout detected, now looking for new PID");
				u.programPID = 0;
				u.pmtFrame = *NewPMTFrame();
				u.detectFlag = false;
			}
		}

		if (u.detectFlag) {
			if (u.tsPacket.GetPID() == u.videoPID) {
				u.previousTime = timeMilisecs();
				u.videoTSParser.Write(u.tsPacket);
			}
		} else {
			if (u.tsPacket.GetType() == Mpeg2TSPacketType_PAT) {
				u.programPID = u.tsPacket.GetProgramPID();
			}

			if (u.tsPacket.GetType() == Mpeg2TSPacketType_PMT) { // PMT
				// received

				if (u.pmtFrame.AddPacket(u.tsPacket)) {
					// pmtFrame is complete
					u.streamsMap = u.pmtFrame.GetStreamInfos();

					//for (map.Entry<Integer, StreamInfo> entry : streamsMap
					//.entrySet()) {
					for k, v := range u.streamsMap {
						var info StreamInfo = v
						if (info.streamType == TS_STREAM_TYPE_H264) {
							fmt.Println("New PID detected = " + string(info.streamPID))
							u.previousTime = timeMilisecs()
							u.detectFlag = true;
							u.videoPID = info.streamPID;
							break;
						}

					}
				}
			}

		}
	}
}
func (u* UdpSource) producer() {
	var firstPacket bool = true
	var buffer= make([]byte, 1500)
	u.packet.SetData(buffer)

	addr, _ := net.ResolveUDPAddr("udp", ":8888")
	sock, _ := net.ListenUDP("udp", addr)
	fmt.Println("working on UDP");
	for u.endFlag != true {
		
	}

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



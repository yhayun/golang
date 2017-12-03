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

func NewDatagramPacket (b []byte, l int) *DatagramPacket {return &DatagramPacket{b, l}}
func NewDatagramPacketE () *DatagramPacket {return &DatagramPacket{}}

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
	//socket DatagramSocket
	videoFrames frameQueue
    programPID int
	socket net.UDPConn
	//outputPort int
	//socketReceiveBufferSize int
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

	//outputSocket int //DatagramSocket

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
				fmt.Println("PAT found")
				u.programPID = u.tsPacket.GetProgramPID();
			}

			if (u.tsPacket.GetType() == Mpeg2TSPacketType_PMT) { // PMT received
				fmt.Println("PMT received")
				if (u.pmtFrame.AddPacket(u.tsPacket)) {// pmtFrame is complete

					u.streamsMap = u.pmtFrame.GetStreamInfos();

					//for (map.Entry<Integer, StreamInfo> entry : streamsMap
					//.entrySet()) {
					for _, v := range u.streamsMap {
						var info StreamInfo = v
						if (info.streamType == TS_STREAM_TYPE_H264) {
							fmt.Println("New PID detected = ", info.streamPID)
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
func (u* UdpSource) producer() { // equivalent
	var firstPacket bool = true
	var buffer= make([]byte, 1500)
	u.packet.SetData(buffer)
	fmt.Println("Entered Producer")

	for u.endFlag != true {
		fmt.Println("-----:: ", u.videoTSParser.tester)
		rlen, _, err := u.socket.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("No packet received")
			continue
		}
		if firstPacket {
			fmt.Println("At least one packet was received")
			firstPacket = false
		}
		var pgPacket = NewDatagramPacket(buffer, rlen)
		//fmt.Println(buffer[:30])
		u.extractStreams(*pgPacket)
	}
	fmt.Println("exited loop")
	//todo clean socket
	defer u.socket.Close();
}


func NewMpeg2TSSource(port int,queue frameQueue) *Mpeg2TSSource {
	//addr, _ := net.ResolveUDPAddr("udp", ":"+string(port))
	addr, _ := net.ResolveUDPAddr("udp", ":8888")
	sock, _ := net.ListenUDP("udp", addr)
	return &Mpeg2TSSource{queue, 0, *sock}
}
func NewUdpSource(port int,queue frameQueue) * UdpSource{
	return &UdpSource{
		*NewMpeg2TSSource(port,queue),
		false,
		false,
		0,
		*NewMpeg2TSParser(queue),
		*NewMpeg2TSPacket(),
		*NewPMTFrame(),
		-1,
		make(map[int]StreamInfo),
		*NewDatagramPacketE(),
	}
}

var Done = make(chan bool)
func main() {

	//addr, _ := net.ResolveUDPAddr("udp", ":8888")
	//sock, _ := net.ListenUDP("udp", addr)
	var videoFrames frameQueue = *NewFrameQueue(100,UDP_SIZE)
	//var tsSource Mpeg2TSSource = *NewMpeg2TSSource(8888, videoFrames)
	var uSource UdpSource = *NewUdpSource(100, videoFrames)
	fmt.Println("working on UDP");
	go uSource.producer()
	go consumer(videoFrames)
	<-Done
}



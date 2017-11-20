package main

import (
	"net"
	"fmt"
	"time"
	"log"
)
const UDP_SIZE = 1500
const MPEG2TS_PACKET_LENGTH = 188
const MAX_UDP_PACKET_SIZE = 1500;


func timeMilisecs() int64 {
	return time.Now().UnixNano()%1e6/1e3
}

type Mpeg2TSSource struct {
	socket DatagramSocket
	reTransmitFlag bool
	videoFrames PooledQueue<Frame>
    programPID int
 	udpSource UdpSource
	outputPort int
	socketReceiveBufferSize int
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
	socket int //DatagramSocket
	outputSocket int //DatagramSocket

	endFlag bool
	detectFlag bool
	videoPID int

	videoTSParser Mpeg2TSParser
	tsPacket Mpeg2TSPacket

	pmtFrame int //PMTFramee
	streamsMap Map[int]StreamInfo //HashMap<Integer, StreamInfo>

	previousTime long
}
func (u* UdpSource) extractStreams(packet int /*DatagramPacket todo*/ ) {

	for offset := 0; offset < packet/*.getLength()todo*/; offset += MPEG2TS_PACKET_LENGTH {

		u.tsPacket.fromBytes(packet.getData(), offset, programPID);

		if(previousTime != -1) {
			if(detectFlag && ((System.currentTimeMillis() - previousTime) > 1000)) {
				log.info("Timeout detected, now looking for new PID");
				programPID = 0;
				pmtFrame = new PMTFrame();
				detectFlag = false;
			}
		}

		if (detectFlag) {
			if (tsPacket.getPID() == videoPID) {
				previousTime = System.currentTimeMillis();
				videoTSParser.write(tsPacket);
			}
		} else {
			if (tsPacket.getType() == Mpeg2TSPacketType.PAT) {
				programPID = tsPacket.getProgramPID();
			}

			if (tsPacket.getType() == Mpeg2TSPacketType.PMT) { // PMT
				// received

				if (pmtFrame.addPacket(tsPacket)) {
					// pmtFrame is complete
					streamsMap = pmtFrame.getStreamInfos();

					for (Map.Entry<Integer, StreamInfo> entry : streamsMap
					.entrySet()) {
						StreamInfo info = (StreamInfo) entry.getValue();
						if (info.streamType == TS_STREAM_TYPE.TS_STREAM_TYPE_H264) {
							log.info("New PID detected = " + info.streamPID);
							previousTime = System.currentTimeMillis();
							detectFlag = true;
							videoPID = info.streamPID;
							break;
						}

					}
				}
			}

		}
	}
}


func (src Mpeg2TSource) extractStream(packet []byte) {
	for offset := 0; offset < len(packet); offset += MPEG2TS_PACKET_LENGTH {
		src.tsPacket.fromBytes(packet, offset, src.programPID)
		if src.detectFlag && (timeMilisecs()-src.previousTime) > 1000 {
			fmt.Println("Timeout detected, now looking for new PID")
			src.programPID = 0

		}
	}
}


//
//
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














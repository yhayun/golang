package main

import (
	"net"
	"fmt"
	"time"
	"io/ioutil"

)

const UDP_SIZE = 1500
const FRAME_SIZE = 300000; // 300KB
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

func (u* UdpSource) FrameQueueFiller() {
	var firstPacket bool = true
	var buffer= make([]byte, 1500)
	u.packet.SetData(buffer)
	fmt.Println("Entered FrameQueueFiller")

	for u.endFlag != true {
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
	defer u.socket.Close();
}


func NewMpeg2TSSource(port string ,queue frameQueue) *Mpeg2TSSource {
	addr, _ := net.ResolveUDPAddr("udp", ":"+ port)
	sock, _ := net.ListenUDP("udp", addr)
	return &Mpeg2TSSource{queue, 0, *sock}
}
func NewUdpSource(port string,queue frameQueue) * UdpSource{
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


//-----------------------------------------------------------//
// tester.go file code moved here
//-----------------------------------------------------------//

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteFile(buffer []byte, file string){
	err := ioutil.WriteFile(file, buffer, 0644)
	check(err)
}

func CheckIfIFrame(buffer []byte, offset int, length int) bool {
	var startOffset int = offset
	for i := startOffset; i < (length -7); i++ {
		var j int = offset + i
		if buffer[j] == 0 {
			if buffer[j+1] == 0 {
				if buffer[j+2] == 1 {
					var val short= (short) (0x001f & buffer[j + 3])
					if val == 7 {
						return true // Start GOP was detected
					}
				}
			}
		}
	}
	return false
}

func FrameQueueDispatcherFullFile(videoFrames frameQueue) {
	fmt.Println("entered  FULL FILE FrameQueueDispatcher")
	h264Buffer := make([]byte, 60*1024*1024)
	var length int = 0
	var counter = 0
	var iframeDetected bool = false
	for {
		//fmt.Println("entered consmer loop")
		if videoFrames.IsEmpty() {
			fmt.Println("consume sleep 100ms")
			time.Sleep(100 * time.Millisecond)
			if videoFrames.IsEmpty() {
				fmt.Println("consume sleep 8000ms")
				time.Sleep(8000 * time.Millisecond)
				if videoFrames.IsEmpty() {
					break
				}
			}
		}
		var frame *Frame  = videoFrames.Poll();
		if !iframeDetected {
			if CheckIfIFrame(frame.GetData(),0, frame.Size()) {
				iframeDetected = true
			} else {
				videoFrames.Recylce(frame)
				continue
			}
		}

		//we have alot 0 sized p_frames...dunno why, but lets skip this for now:
		if (frame.Size() == 0) {
			videoFrames.Recylce(frame)
			continue
		}
		counter++
		fmt.Println("counter: ",counter)
		ArrayCopy(frame.GetData(),0, h264Buffer,length,frame.Size())
		length += frame.Size()
		videoFrames.Recylce(frame)
		if (counter >= 500){
			break;
		}
	}
	fmt.Println("left FrameQueueDispatcher loop.")
	WriteFile(h264Buffer[:length],"../../full.264");
}




func FrameQueueDispatcher(videoFrames frameQueue) {
	fmt.Println("entered FrameQueueDispatcher")
	//Used to identify frames when written to files.
	var counter = 0
	var i_counter = 0
	var p_counter = 0
	var p_type = 'P'
	var i_type = 'I'
	var s string //= fmt.Sprintf("../tmp/frame[%d]-I[%d]-P[%d]_TYPE<%c>", counter,i_counter,p_counter,p_type)
	var iframeDetected bool = false

	for {
		//fmt.Println("entered consmer loop")
		if videoFrames.IsEmpty() {
			fmt.Println("consume sleep 100ms")
			time.Sleep(100 * time.Millisecond)
			if videoFrames.IsEmpty() {
				fmt.Println("consume sleep 8000ms")
				time.Sleep(8000 * time.Millisecond)
				if videoFrames.IsEmpty() {
					break
				}
			}
		}
		var frame *Frame  = videoFrames.Poll();
		if !iframeDetected {
			if CheckIfIFrame(frame.GetData(),0, frame.Size()) {
				iframeDetected = true
			} else {
				videoFrames.Recylce(frame)
				continue
			}
		}

		//we have alot 0 sized p_frames...dunno why, but lets skip this for now:
		if (frame.Size() == 0) {
			videoFrames.Recylce(frame)
			continue
		}

		counter++
		if CheckIfIFrame(frame.GetData(),0, frame.Size()) {
			i_counter++
			p_counter = 0
			s = fmt.Sprintf("../../tmp/frame[%d]-I[%d]-P[%d]_TYPE(%c)____size_%d-----", counter,i_counter,p_counter,i_type,frame.Size())
		} else {
			p_counter++;
			s = fmt.Sprintf("../../tmp/frame[%d]-I[%d]-P[%d]_TYPE(%c)____size_%d", counter,i_counter,p_counter,p_type,frame.Size())
		}

		fmt.Println("file: ",counter)
		actualData := frame.GetData()[:frame.Size()]
		WriteFile(actualData,s)
		WriteFile(actualData,fmt.Sprint("../../media/",counter))
		videoFrames.Recylce(frame)
		if (counter >= 500){
			break;
		}

	}
	fmt.Println("left FrameQueueDispatcher loop.")
}




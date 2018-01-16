package main
// https://www.slideshare.net/gamzabaw/implementing-hls-server-with-go
import (
	"net/http"
	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
	"time"
	"encoding/binary"
)

var Done = make(chan bool)
var Queue = make(chan []byte, 1)
var counter int = 1

func getMediaBase(mId int) string{
	mediaRoot := "assessts/media"
	return fmt.Sprintf("%s/%d", mediaRoot,mId)
}


func serveHlsFile( w http.ResponseWriter, r *http.Request, mediabase, segName string) {
	//mediaFile := fmt.Sprint("%s/hls/%s",mediabase,segName)
	//http.ServeFile(w,r,mediaFile)
	name := fmt.Sprint("../../media/",counter)
	counter++
	http.ServeFile(w,r,name)
	fmt.Println("%s",name)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "video/MP2T")
}


func serveHlsQueue( w http.ResponseWriter, r *http.Request, mediabase, segName string) {
	body  := <-Queue
	//fmt.Println(body)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "text/plain; charset=x-user-defined")
	pts := (body[len(body)-16:len(body)-8])//bytesToUint64
	dts := (body[:len(body)-8])
	fmt.Print("pts", pts, " dts:",0, "\n")
	w.Write(body)
}

func streamHandlerFile(response http.ResponseWriter, request *http.Request) {
	fmt.Println("streamHandler")
	vars := mux.Vars(request)
	mId, err := strconv.Atoi(vars["mId"])
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	segName, _ := vars["segName"]
	mediaBase := getMediaBase(mId)
	serveHlsFile(response, request, mediaBase, segName)
}

func streamHandlerQueue(response http.ResponseWriter, request *http.Request) {
	fmt.Println("streamHandler")
	vars := mux.Vars(request)
	mId, err := strconv.Atoi(vars["mId"])
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	segName, _ := vars["segName"]
	mediaBase := getMediaBase(mId)
	serveHlsQueue(response, request, mediaBase, segName)
}


func Handlers() *mux.Router {
	router :=  mux.NewRouter()
	router.HandleFunc("/media/stream/{segName:test}/{mId:[0-9]+}", streamHandlerQueue).Methods("GET")
	return router
}


func RunServerQueue () {
	fmt.Println("Entered Server")
	http.Handle("/",Handlers())
	http.ListenAndServe(":8000",nil)
}


// // Used to force main thread to go to sleep so we can handle when program stops running.
//func main() {
//	var videoFrames frameQueue = *NewFrameQueue(100,FRAME_SIZE)
//	//var tsSource Mpeg2TSSource = *NewMpeg2TSSource(8888, videoFrames)
//	var uSource UdpSource = *NewUdpSource(100, videoFrames)
//	fmt.Println("working on UDP");
//	go uSource.FrameQueueFiller()
//	go FrameQueueDispatcherFullFile(videoFrames)
//	<-Done
//}


func main() {
	var videoFrames frameQueue = *NewFrameQueue(100,FRAME_SIZE)
	var uSource UdpSource = *NewUdpSource(100, videoFrames)
	fmt.Println("working on UDP");
	go uSource.FrameQueueFiller()
	go FinalQueueFilller(videoFrames)
	go RunServerFiles()
	<-Done
}

// https://stackoverflow.com/questions/22452804/angularjs-http-get-request-failed-with-custom-header-alllowed-in-cors
func RunServerFiles() {
	r := mux.NewRouter()
	r.HandleFunc("/media/stream/{segName:test}/{mId:[0-9]+}", streamHandlerQueue)
	http.Handle("/", &MyServer{r})
	http.ListenAndServe(":8000", nil);
}

type MyServer struct {
	r *mux.Router
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	// Stop here if its Preflighted OPTIONS request
	if req.Method == "OPTIONS" {
		return
	}
	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}




func FinalQueueFilller(videoFrames frameQueue ) {
	fmt.Println("entered  FinalQueueFilller")
	var length int = 0
	var counter = 0
	var iframeDetected bool = false
	for {
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
		if (frame.Size() == 0) {
			videoFrames.Recylce(frame)
			continue
		}

		counter++
		fmt.Println("counter: ",counter,"  dts:",frame.GetDTS() ,"  pts:",frame.GetPTS())
		actualData := frame.GetData()[:frame.Size()]
		actualData = append(actualData, uint64TObytes((frame.GetPTS()))...)
		actualData = append(actualData, uint64TObytes((frame.GetDTS()))...)
		//fmt.Println(actualData)
		length += frame.Size()
		Queue <- actualData  // data, PTS(8bytes), DTS(8bytes)
		videoFrames.Recylce(frame)
	}
}

func uint64TObytes(num uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(num))
	return b;
}
func bytesToUint64(b []byte) uint64 {
	return uint64(binary.LittleEndian.Uint64(b))
}
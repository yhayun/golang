package main
// https://www.slideshare.net/gamzabaw/implementing-hls-server-with-go
import (
	"net/http"
	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
)


var counter int = 0

func getMediaBase(mId int) string{
	mediaRoot := "assessts/media"
	return fmt.Sprintf("%s/%d", mediaRoot,mId)
}


func serveHlsTls( w http.ResponseWriter, r *http.Request, mediabase, segName string) {
	//mediaFile := fmt.Sprint("%s/hls/%s",mediabase,segName)
	//http.ServeFile(w,r,mediaFile)

	//testing"
	name := fmt.Sprint("../../media/",counter)
	counter++

	http.ServeFile(w,r,name)
	fmt.Println("%s",name)
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "video/MP2T")
}

func streamHnandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("streamHandler")
	vars := mux.Vars(request)
	mId, err := strconv.Atoi(vars["mId"])
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	segName, _ := vars["segName"]
	mediaBase := getMediaBase(mId)
	serveHlsTls(response, request, mediaBase, segName)
}


func Handlers() *mux.Router {
	router :=  mux.NewRouter()
	router.HandleFunc("/media/stream/{segName:test}/{mId:[0-9]+}", streamHnandler).Methods("GET")
	return router
}


func RunServer () {
	fmt.Println("started")
	http.Handle("/",Handlers())
	http.ListenAndServe(":8000",nil)
}


//Used to force main thread to go to sleep so we can handle when program stops running.
//var Done = make(chan bool)
//func main() {
//	var videoFrames frameQueue = *NewFrameQueue(100,FRAME_SIZE)
//	//var tsSource Mpeg2TSSource = *NewMpeg2TSSource(8888, videoFrames)
//	var uSource UdpSource = *NewUdpSource(100, videoFrames)
//	fmt.Println("working on UDP");
//	go uSource.FrameQueueFiller()
//	go FrameQueueDispatcher(videoFrames)
//	<-Done
//}

func main() {
	RunServer()
}

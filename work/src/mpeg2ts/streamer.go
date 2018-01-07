package main

import (
	"net/http"
	//"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"strconv"
	"fmt"
)



// todo reference for hls stramer:  https://www.slideshare.net/gamzabaw/implementing-hls-server-with-go

func getMediaBase(mId int) string{
	mediaRoot := "assessts/media"
	return fmt.Sprintf("%s/%d", mediaRoot,mId)
}

func serveHlsTls( w http.ResponseWriter, r *http.Request, mediabase, segName string) {
	//mediaFile := fmt.Sprint("%s/hls/%s",mediabase,segName)
	//http.ServeFile(w,r,mediaFile)
	name := "../../tmp/frame[1]-I[1]-P[0]_TYPE<I>__size_80956"
	http.ServeFile(w,r,name)
	fmt.Println("%s",name)
	w.Header().Set("Content-Type", "video/MP2T")
}

func streamHnandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("streamHandler")
	vars := mux.Vars(request)
	mId, err := strconv.Atoi(vars["mId"])
	fmt.Println("%d",mId)
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


func main () {
	fmt.Println("started")
	http.Handle("/",Handlers())
	http.ListenAndServe(":8000",nil)
}


////Used to force main thread to go to sleep so we can handle when program stops running.
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

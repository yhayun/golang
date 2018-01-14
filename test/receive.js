//temp to emulate xmlhttprequest on node:
var XMLHttpRequest = require("xmlhttprequest").XMLHttpRequest;
var Flow = require("./flow.js")

var baseUrl = "http://localhost:8000/media/stream/test/1";
var binaryClient = new BinaryHttpClient();
var queue  = [];
var initFlag = false
var counter = 1
//Run Queue Filler:
FillQueue();
/**
 *  Calling order: FillQueue() -> (retriveData() <-> handleLoop()). in handle loop once we have enough frames
 *  we can start handling the data.
 */

    function  FillQueue() {
        var xmlhttp = new XMLHttpRequest();
        retrieveData()
        
	};

	function handleLoop() {
		if (queue.length == 500) {
			console.log("Queue ready:")
			return ProcessQueue();
		}
		retrieveData()
	}

	function ProcessQueue() {
		// for (u= 0; u < queue.length; u++) {
		// 	var units = Flow.demuxerTS._parseAVCNALu(queue[u])
		// 	//console.log("units: ",units)
		// 	for (var i =0; i <units.length; i++) {
		// 		//console.log("type:", units[i].type)
		// 		if( units[i].type === 7) {
		// 			expGolombDecoder = new Flow.ExpGolomb(units[i].data);
		// 			var config = expGolombDecoder.readSPS();
		// 			console.log("config:", config)
		// 		}
		// 	}			
		// }

		writeQueueToFIle();

		// var units = Flow.demuxerTS._parseAVCNALu(queue[0])
		// console.log("units: ",units)
		// expGolombDecoder = new Flow.ExpGolomb(units[1].data);
		// var config = expGolombDecoder.readSPS();
		// console.log("config:", config)

		return
	}

	function writeQueueToFIle(_data,_path) {
		var fs = require("fs");
		var path = "full_file_test.264";
		// var data;
		// for (u= 0; u < queue.length; u++) {
		// 	data += queue[u];	
		// 	console.log(`${u}\n\n\n`,queue[u])		
		// }

		fs.writeFile(_path, _data, function(error) {
		     if (error) {
		       console.error("write error:  " + error.message);
		     } else {
		       console.log("Successful Write to ");
		     }
		});
	}


    function BinaryHttpClient() {
        this.get = function(aUrl, aCallback) {
            var anHttpRequest = new XMLHttpRequest();
            anHttpRequest.open( "GET", aUrl, true );            
            anHttpRequest.responseType = "arraybuffer";
            
            anHttpRequest.onload = function (oEvent) {
            	 if (anHttpRequest.readyState == 4 && anHttpRequest.status == 200) {
	                  //var arrayBuffer = anHttpRequest.response; 
	                  var arrayBuffer = str2ab(anHttpRequest.responseText); 
	                   writeQueueToFIle(arrayBuffer, "output/_"+counter);
	                   counter++
	                  if (arrayBuffer) {
	                    var byteArray = new Uint8Array(arrayBuffer);
	                    aCallback(byteArray);
	                  }
              	} else {
              		console.log("Request failed with status:",anHttpRequest.status)
              		if (queue.length != 0) {
              			ProcessQueue();
              		}
              	}
            };       
            anHttpRequest.send( null );
        }
    }



    function retrieveData() {
        var byteArray = null;
        binaryClient.get(baseUrl, function(byteArray) {
            if (byteArray.byteLength == 0) {
                console.log('Skipping since byteArray.byteLength = ' + byteArray.byteLength);
                return;
            }
            if(initFlag == false) {
                console.log('Got data first time: ' + byteArray.byteLength);   
                initFlag = true;                     
                queue.push(byteArray); 
                handleLoop()
                return;
            }
            console.log('Got data: ' + byteArray.byteLength);
            queue.push(byteArray);
            console.log('queue.length: ' + queue.length);
            handleLoop()
        });     
    }


// our server sends the file as texts, we need to convert it to bytearray for this to work
// https://developers.google.com/web/updates/2012/06/How-to-convert-ArrayBuffer-to-and-from-String
function ab2str(buf) {
  return String.fromCharCode.apply(null, new Uint8Array(buf));
}
function str2ab(str) {
  var buf = new ArrayBuffer(str.length*2); // 2 bytes for each char
  var bufView = new Uint8Array(buf);
  for (var i=0, strLen=str.length; i < strLen; i++) {
    bufView[i] = str.charCodeAt(i);
  }
  return buf;
}
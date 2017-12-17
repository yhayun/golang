package main

// todo - https://github.com/video-dev/hls.js/blob/ce19344d1bfb025e4d55c5b1521f0c8f02c06b5b/src/demux/tsdemuxer.js
type VideoTrack struct {
	id uint
	pid uint
	inputTimeScale uint
	sequenceNumber uint
	samples []byte
	len uint
	dropped uint
	isAAC bool
	duration bool

	//For the first frame also:
	width int
	height int
	sps int//(GOP header)
	pps int //(I frame header)

	//added during ParseAVCNALu:
	naluState int
}

func NewVideoTrack() *VideoTrack {
	return &VideoTrack{
		id: 0,
		pid : -1,
		inputTimeScale : 90000,
		sequenceNumber: 0,
		//samples : [],
		len : 0,
		dropped: 0, // type === 'video' ? 0 : undefined,
		isAAC: false, //type === 'audio' ? true : undefined,
		duration: false, //type === 'audio' ? duration : undefined


		//First Frame only:
		width: -1,
		height: -1,
		sps: -1,
		pps: -1,


		//added during ParseAVCNALu:
		naluState: 0,
	}
}


type unit struct {

}

func ParseAVCNALu(array []byte, track VideoTrack) unit {
	i  := 0
	len := len(array)
	var value byte
	overflow := 0
	state := track.naluState
	lastState := state
	//units
	//unit
	//unitType
	lastUnitStart := -1
	var lastUnitType byte

    if (state == -1) {
    // special use case where we found 3 or 4-byte start codes exactly at the end of previous PES packet
      lastUnitStart = 0;
      // NALu type is value read from offset 0
      lastUnitType = array[0] & 0x1f;
      state = 0;
      i = 1;
    }

    for (i < len) {
      value = array[i];
      i++
      // optimization. state 0 and 1 are the predominant case. let's handle them outside of the switch/case
      if (state == 0) {
      	if (value != 0) {
      		state = 0
      		continue
		}
        state = 1
        continue;
      }
      if (state == 1) {
		  if (value != 0) {
			  state = 0
			  continue
		  }
        state =  2
        continue
      }

      // here we have state either equal to 2 or 3
      if(!value) {
        state = 3;
      } else if (value === 1) {
        if (lastUnitStart >=0) {
          unit = {data: array.subarray(lastUnitStart, i - state - 1), type: lastUnitType};
          //logger.log('pushing NALU, type/size:' + unit.type + '/' + unit.data.byteLength);
          units.push(unit);
        } else {
          // lastUnitStart is undefined => this is the first start code found in this PES packet
          // first check if start code delimiter is overlapping between 2 PES packets,
          // ie it started in last packet (lastState not zero)
          // and ended at the beginning of this PES packet (i <= 4 - lastState)
          let lastUnit = this._getLastNalUnit();
          if (lastUnit) {
            if(lastState &&  (i <= 4 - lastState)) {
              // start delimiter overlapping between PES packets
              // strip start delimiter bytes from the end of last NAL unit
                // check if lastUnit had a state different from zero
              if (lastUnit.state) {
                // strip last bytes
                lastUnit.data = lastUnit.data.subarray(0,lastUnit.data.byteLength - lastState);
              }
            }
            // If NAL units are not starting right at the beginning of the PES packet, push preceding data into previous NAL unit.
            overflow  = i - state - 1;
            if (overflow > 0) {
              //logger.log('first NALU found with overflow:' + overflow);
              let tmp = new Uint8Array(lastUnit.data.byteLength + overflow);
              tmp.set(lastUnit.data, 0);
              tmp.set(array.subarray(0, overflow), lastUnit.data.byteLength);
              lastUnit.data = tmp;
            }
          }
        }
        // check if we can read unit type
        if (i < len) {
          unitType = array[i] & 0x1f;
          //logger.log('find NALU @ offset:' + i + ',type:' + unitType);
          lastUnitStart = i;
          lastUnitType = unitType;
          state = 0;
        } else {
          // not enough byte to read unit type. let's read it on next PES parsing
          state = -1;
        }
      } else {
        state = 0;
      }
    }
    if (lastUnitStart >=0 && state >=0) {
      unit = {data: array.subarray(lastUnitStart, len), type: lastUnitType, state : state};
      units.push(unit);
      //logger.log('pushing NALU, type/size/state:' + unit.type + '/' + unit.data.byteLength + '/' + state);
    }
    // no NALu found
    if (units.length === 0) {
      // append pes.data to previous NAL unit
      let  lastUnit = this._getLastNalUnit();
      if (lastUnit) {
        let tmp = new Uint8Array(lastUnit.data.byteLength + array.byteLength);
        tmp.set(lastUnit.data, 0);
        tmp.set(array, lastUnit.data.byteLength);
        lastUnit.data = tmp;
      }
    }
    track.naluState = state;
    return units;

}





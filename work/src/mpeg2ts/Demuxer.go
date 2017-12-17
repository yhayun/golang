package main


type Unit struct {
	data []byte
	Type byte
	State int
}


type Sample struct {
	units []Unit
	pts float64
	dts float64
	length int
}

type Config struct {
	width int
	height int
	pixelRatio [2]uint
}


// todo - https://github.com/video-dev/hls.js/blob/ce19344d1bfb025e4d55c5b1521f0c8f02c06b5b/src/demux/tsdemuxer.js
type VideoTrack struct {
	id uint
	pid uint
	inputTimeScale uint
	timescale int
	sequenceNumber uint
	samples []Sample
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


type TSDemuxer struct {
	track VideoTrack
}

func NewTSDemuxer() *TSDemuxer {
	return &TSDemuxer{
		*NewVideoTrack(),
	}
}

func (this *TSDemuxer) GetLastNalUnit() *Unit {
	// todo - do we need this in our special 1 frame case????
    // avcSample = this.avcSample, lastUnit;
    //// try to fallback to previous sample if current one is empty
    //if (!avcSample || avcSample.units.length === 0) {
    //  let track = this._avcTrack, samples = track.samples;
    //  avcSample = samples[samples.length-1];
    //}
    //if (avcSample) {
    //  let units = avcSample.units;
    //  lastUnit = units[units.length - 1];
    //}
    //return lastUnit;
    return nil
  }


func (this *TSDemuxer) ParseAVCNALu(array []byte, ) []Unit {
	i  := 0
	len := len(array)
	var value byte = -1
	//overflow := 0
	state := this.track.naluState
	//lastState := state
	var units []Unit
	//unit
	var unitType byte
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
      	if (value != -1) {
      		state = 0
      		continue
		}
        state = 1
        continue;
      }
      if (state == 1) {
		  if (value != -1) {
			  state = 0
			  continue
		  }
        state =  2
        continue
      }

      // here we have state either equal to 2 or 3
      if(value == -1) {
        state = 3;
      } else if (value == 1) {
        if (lastUnitStart >=0 ) {
          unit := Unit{ array[lastUnitStart:i - state - 1], lastUnitType, 0}
          units = append(units,unit)
        } else {
          // lastUnitStart is undefined => this is the first start code found in this PES packet
          // first check if start code delimiter is overlapping between 2 PES packets,
          // ie it started in last packet (lastState not zero)
          // and ended at the beginning of this PES packet (i <= 4 - lastState)
          lastUnit := this.GetLastNalUnit();
          if (lastUnit != nil) {
			  // todo - do we need this in our special 1 frame case????
            //if(lastState != -1 &&  (i <= 4 - lastState)) {
            //  // start delimiter overlapping between PES packets
            //  // strip start delimiter bytes from the end of last NAL unit
            //    // check if lastUnit had a state different from zero
            //  if (lastUnit.state != -1) {
            //    // strip last bytes
            //    lastUnit.data = lastUnit.data.subarray(0,lastUnit.data.byteLength - lastState);
            //  }
            //}
            //// If NAL units are not starting right at the beginning of the PES packet, push preceding data into previous NAL unit.
            //overflow  = i - state - 1;
            //if (overflow > 0) {
            //  //logger.log('first NALU found with overflow:' + overflow);
            //  let tmp = new Uint8Array(lastUnit.data.byteLength + overflow);
            //  tmp.set(lastUnit.data, 0);
            //  tmp.set(array.subarray(0, overflow), lastUnit.data.byteLength);
            //  lastUnit.data = tmp;
            //}
          }
        }
        // check if we can read unit type
        if (i < len) {
          unitType = array[i] & 0x1f;
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
      unit := Unit{ array[lastUnitStart:len], lastUnitType, state}
      units = append(units,unit)
    }
    // no NALu found
    lenUnits := 0
	for range units {
		lenUnits++
	}
    if (lenUnits == 0) {
		// todo - do we need this in our special 1 frame case????
      //// append pes.data to previous NAL unit
      //let  lastUnit = this._getLastNalUnit();
      //if (lastUnit) {
      //  let tmp = new Uint8Array(lastUnit.data.byteLength + array.byteLength);
      //  tmp.set(lastUnit.data, 0);
      //  tmp.set(array, lastUnit.data.byteLength);
      //  lastUnit.data = tmp;
      //}
    }
    this.track.naluState = state;
    return units;
}





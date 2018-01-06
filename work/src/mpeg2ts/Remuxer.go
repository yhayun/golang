package main

import (
	"sort"
	"math"
	"encoding/binary"
)

type VideoData struct {
	data1    []byte //moof
	data2    []byte //mdat
	startPTS int
	endPTS   int
	_type    string
	nb       int // maybe nb is number of bytes??
	dropped  int
}
type TrackVideo_metaData struct {
	width int
	height int
}

type TracksVideo struct {
	container string//: 'video/mp4',
	codec string // :  videoTrack.codec, (might be 'avc1.42e01e').
	initSegment : MP4.initSegment([videoTrack]),
	metadata TrackVideo_metaData
}

const MaxInt = int(^uint(0) >> 1)

type Flags struct {
	isLeading int//: 0,
	isDependedOn int//: 0,
	hasRedundancy int//: 0,
	degradPrio int//: 0,
	dependsOn int//: 1,
	isNonSync int// :1,
}

type Mp4Sample struct {
	size int// unitLen,
	cts int //0,
	duration int// : 0,
	flags Flags

}

// our project note - we need to suport mp4 only.
type TypeSupported struct {
	mp4  bool //: MediaSource.isTypeSupported('video/mp4'),
	mpeg bool //: MediaSource.isTypeSupported('audio/mpeg'),
	mp3  bool //: MediaSource.isTypeSupported('audio/mp4; codecs="mp3"')
};

func Round(num float64) int64 {
	return int64(num + 0.5)
}

func PTSNormalize(value float64, reference *float64) float64 {
	var offset float64
	if (reference == nil) {
		return value;
	}
	if (*reference < value) {
		// - 2^33
		offset = -8589934592;
	} else {
		// + 2^33
		offset = 8589934592;
	}
	/* PTS is 33bit (from 0 to 2^33 -1)
	  if diff between value and reference is bigger than half of the amplitude (2^32) then it means that
	  PTS looping occured. fill the gap */
	for ; (math.Abs(value - *reference) > 4294967296); {
		value += offset;
	}
	return value;
}

type MP4Remuxer struct {
	//constructor(observer, config, typeSupported, vendor) {
	//this.observer = observer;
	//this.config = config;
	//this.typeSupported = typeSupported;
	//const userAgent = navigator.userAgent;
	//this.isSafari = vendor && vendor.indexOf('Apple') > -1 && userAgent && !userAgent.match('CriOS');
	_initPTS      int
	_initDTS      int
	_PTSNormalize int
	nextAvcDts    *float64
	ISGenerated   bool //= false;
	typeSupported TypeSupported
}

func NewMP4Remuxer() *MP4Remuxer {
	return &MP4Remuxer{
		ISGenerated: false,
		typeSupported: TypeSupported{true , false , false },
	}
}

func (_this *MP4Remuxer) Remux(videoTrack VideoTrack, timeOffset uint, contiguous uint, accurateTimeOffset uint) {
	// generate Init Segment if needed
	if (!_this.ISGenerated) {
		_this.GenerateIS(videoTrack, timeOffset);
	}

	if (_this.ISGenerated) {
		//const nbAudioSamples = audioTrack.samples.length;
		var nbVideoSamples = len(videoTrack.samples)
		var videoTimeOffset = timeOffset;
		//if (nbAudioSamples && nbVideoSamples) {
		//  // timeOffset is expected to be the offset of the first timestamp of this fragment (first DTS)
		//  // if first audio DTS is not aligned with first video DTS then we need to take that into account
		//  // when providing timeOffset to remuxAudio / remuxVideo. if we don't do that, there might be a permanent / small
		//  // drift between audio and video streams
		//  let audiovideoDeltaDts = (audioTrack.samples[0].dts - videoTrack.samples[0].dts)/videoTrack.inputTimeScale;
		//  audioTimeOffset += math.max(0,audiovideoDeltaDts);
		//  videoTimeOffset += math.max(0,-audiovideoDeltaDts);
		//}  todo we don't need this crap, we only stream video

		// Purposefully remuxing audio before video, so that remuxVideo can use nextAudioPts, which is
		// calculated in remuxAudio.
		//logger.log('nb AAC samples:' + audioTrack.samples.length);
		//if (nbAudioSamples) { todo always zero
		//  // if initSegment was generated without video samples, regenerate it again
		//  if (!audioTrack.timescale) {
		//    logger.warn('regenerate InitSegment as audio detected');
		//    this.generateIS(audioTrack,videoTrack,timeOffset);
		//  }
		//  let audioData = this.remuxAudio(audioTrack,audioTimeOffset,contiguous,accurateTimeOffset);
		//  //logger.log('nb AVC samples:' + videoTrack.samples.length);
		//  if (nbVideoSamples) {
		//    let audioTrackLength;
		//    if (audioData) {
		//      audioTrackLength = audioData.endPTS - audioData.startPTS;
		//    }
		//    // if initSegment was generated without video samples, regenerate it again
		//    if (!videoTrack.timescale) {
		//      logger.warn('regenerate InitSegment as video detected');
		//      this.generateIS(audioTrack,videoTrack,timeOffset);
		//    }
		//    this.remuxVideo(videoTrack,videoTimeOffset,contiguous,audioTrackLength, accurateTimeOffset);
		//  }
		//} else {
		var videoData VideoData
		//logger.log('nb AVC samples:' + videoTrack.samples.length);
		if nbVideoSamples != 0 {
			videoData = _this.RemuxVideo(videoTrack, videoTimeOffset, contiguous, accurateTimeOffset);
		}
		//if (videoData && audioTrack.codec) {
		//  this.remuxEmptyAudio(audioTrack, audioTimeOffset, contiguous, videoData);
		//} todo we don't need this
		//}
	}
	////logger.log('nb ID3 samples:' + audioTrack.samples.length);
	//if (id3Track.samples.length) { todo zero
	//  this.remuxID3(id3Track,timeOffset);
	//}
	////logger.log('nb ID3 samples:' + audioTrack.samples.length);
	//if (textTrack.samples.length) {
	//  this.remuxText(textTrack,timeOffset);
	//}
	//notify end of parsing
	//_this.observer.trigger(Event.FRAG_PARSED); todo observer
}

func (_this *MP4Remuxer) GenerateIS(videoTrack VideoTrack, timeOffset uint) {
	//var observer = _this.observer  //todo observer
	//audioSamples = audioTrack.samples,
	var videoSamples = videoTrack.samples
	var typeSupported = _this.typeSupported
	var container = "audio/mp4"
	var computePTSDTS = (_this._initPTS == undefined)
	var initPTS float64
	var initDTS float64

	if (computePTSDTS) {
		initPTS = math.MaxFloat64
		initDTS = math.MaxFloat64
	}
	//if (audioTrack.config && audioSamples.length) { todo no audio samples => always false
	//  // let's use audio sampling rate as MP4 time scale.
	//  // rationale is that there is a integer nb of audio frames per audio sample (1024 for AAC)
	//  // using audio sampling rate here helps having an integer MP4 frame duration
	//  // this avoids potential rounding issue and AV sync issue
	//  audioTrack.timescale = audioTrack.samplerate;
	//  logger.log (`audio sampling rate : ${audioTrack.samplerate}`);
	//  if (!audioTrack.isAAC) {
	//	if (typeSupported.mpeg) { // Chrome and Safari
	//	  container = 'audio/mpeg';
	//	  audioTrack.codec = '';
	//	} else if (typeSupported.mp3) { // Firefox
	//	  audioTrack.codec = 'mp3';
	//	}
	//  }
	//  tracks.audio = {
	//	container : container,
	//	codec :  audioTrack.codec,
	//	initSegment : !audioTrack.isAAC && typeSupported.mpeg ? new Uint8Array() : MP4.initSegment([audioTrack]),
	//	metadata : {
	//	  channelCount : audioTrack.channelCount
	//	}
	//  };
	//  if (computePTSDTS) {
	//	// remember first PTS of this demuxing context. for audio, PTS = DTS
	//	initDTS = audioSamples[0].pts - audioTrack.inputTimeScale * timeOffset;
	//	initPTS = initDTS
	//  }todo maybe this is relevant
	//}

	if (videoTrack.sps !=0 && videoTrack.pps!=0 && len(videoSamples)!=0) {
		// let's use input time scale as MP4 video timescale
		// we use input time scale straight away to avoid rounding issues on frame duration / cts computation
		var inputTimeScale = videoTrack.inputTimeScale
		videoTrack.timescale = int(inputTimeScale)
		tracks.video = Video{
		container:
			"video/mp4",
				codec :  videoTrack.codec,
			initSegment : MP4.initSegment([videoTrack]),
			metadata: {
		width: videoTrack.width,
		height: videoTrack.height
		}
		};
		if (computePTSDTS) {
			initPTS = math.Min(initPTS, videoSamples[0].pts-inputTimeScale*timeOffset);
			initDTS = math.Min(initDTS, videoSamples[0].dts-inputTimeScale*timeOffset);
			this.observer.trigger(Event.INIT_PTS_FOUND,
			{
			initPTS:
				initPTS
			});
		}
	}
	if (Object.keys(tracks).length) {
		observer.trigger(Event.FRAG_PARSING_INIT_SEGMENT, data);
		_this.ISGenerated = true
		if (computePTSDTS) {
			_this._initPTS = initPTS;
			_this._initDTS = initDTS;
		}
	} else {
		observer.trigger(Event.ERROR,
		{
			type: ErrorTypes.MEDIA_ERROR, details: ErrorDetails.FRAG_PARSING_ERROR, fatal: false, reason: 'n
			o
			audio / video
			samples
			found
			'});
		}
	}
}

type byHarta []Sample

func (s byHarta) Len() int {
	return len(s)
}
func (s byHarta) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byHarta) Less(i, j int) bool {
	var deltadts = s[i].dts - s[j].dts;
	var deltapts = s[i].pts - s[j].pts;
	if deltadts != 0 {
		return true
	} else {
		//if (deltapts != 0){
		return true
		//}else{
		//return s[i].id - s[j].id
		//}

	}
}

func (_this *MP4Remuxer) RemuxVideo(track VideoTrack, timeOffset, contiguous /*=0*/ , accurateTimeOffset int /*=0*/) {
	var offset = 8
	var timeScale = track.timescale
	var mp4SampleDuration
	var mdat
	var moof
	var firstPTS
	var firstDTS
	var nextDTS
	var lastPTS
	var lastDTS
	var inputSamples = track.samples
	var outputSamples []Mp4Sample
	var nbSamples = len(inputSamples)
	//var ptsNormalize = _this._PTSNormalize
	var initDTS = _this._initDTS

	// for (let i = 0; i < track.samples.length; i++) {
	//   let avcSample = track.samples[i];
	//   let units = avcSample.units;
	//   let unitsString = '';
	//   for (let j = 0; j < units.length ; j++) {
	//     unitsString += units[j].type + ',';
	//     if (units[j].data.length < 500) {
	//       unitsString += Hex.hexDump(units[j].data);
	//     }
	//   }
	//   logger.log(avcSample.pts + '/' + avcSample.dts + ',' + unitsString + avcSample.units.length);
	// }

	// if parsed fragment is contiguous with last one, let's use last DTS value as reference
	var nextAvcDts = _this.nextAvcDts;

	//const isSafari = _this.isSafari; todo who cares if this is safari
	//
	//// Safari does not like overlapping DTS on consecutive fragments. let's use nextAvcDts to overcome this if fragments are consecutive
	//if (isSafari) {
	//// also consider consecutive fragments as being contiguous (even if a level switch occurs),
	//// for sake of clarity:
	//// consecutive fragments are frags with
	////  - less than 100ms gaps between new time offset (if accurate) and next expected PTS OR
	////  - less than 200 ms PTS gaps (timeScale/5)
	//contiguous |= (inputSamples.length && nextAvcDts &&
	//((accurateTimeOffset && math.abs(timeOffset-nextAvcDts/timeScale) < 0.1) ||
	//math.abs((inputSamples[0].pts-nextAvcDts-initDTS)) < timeScale/5)
	//);
	//}

	if contiguous != 0 {
		// if not contiguous, let's use target timeOffset
		*nextAvcDts = timeOffset * timeScale;
	}

	// PTS is coded on 33bits, and can loop from -2^32 to 2^32
	// ptsNormalize will make PTS/DTS value monotonic, we use last known DTS value as reference value
	for i := 0; i < len(inputSamples); i++ {
		//var sample Sample = inputSamples[i]
		inputSamples[i].pts = PTSNormalize(inputSamples[i].pts-initDTS, nextAvcDts);
		inputSamples[i].dts = PTSNormalize(inputSamples[i].dts-initDTS, nextAvcDts);
	}
	//inputSamplesforEach(function(sample) {
	//});

	// sort video samples by DTS then PTS then demux id order

	//inputSamples.sort(function(a, b) {
	//const deltadts = a.dts - b.dts;
	//const deltapts = a.pts - b.pts;
	//return deltadts ? deltadts : deltapts ? deltapts : (a.id - b.id);
	//});
	sort.Sort(byHarta(inputSamples))

	// handle broken streams with PTS < DTS, tolerance up 200ms (18000 in 90kHz timescale)
	var PTSDTSshift = 0.
	for i:= 1; i< len(inputSamples); i++ {
		PTSDTSshift = PTSDTSshift + math.Max(math.Min(inputSamples[i-1], inputSamples[i].pts-inputSamples[i].dts), -18000)
	}
	if (PTSDTSshift < 0) {
		//logger.warn(`PTS < DTS detected in video samples, shifting DTS by ${math.round(PTSDTSshift/90)} ms to overcome this issue`);
		for i := 0; i < len(inputSamples);i++ {
			inputSamples[i].dts += PTSDTSshift;
		}
	}

	// compute first DTS and last DTS, normalize them against reference value
	var sample = inputSamples[0];
	firstDTS = math.Max(float64(sample.dts), 0);
	firstPTS = math.Max(float64(sample.dts), 0);

	// check timestamp continuity accross consecutive fragments (this is to remove inter-fragment gap/hole
	var delta = math.Round((firstDTS - float64(*nextAvcDts)) / 90);
	// if fragment are contiguous, detect hole/overlapping between fragments
	if (contiguous != 0) {
		if (delta) {
			if (delta > 1) {
				//logger.log(`AVC:${delta} ms hole between fragments detected,filling it`);
			} else if (delta < -1) {
				//logger.log(`AVC:${(-delta)} ms overlapping between fragments detected`);
			}
			// remove hole/gap : set DTS to next expected DTS
			firstDTS = *nextAvcDts;
			inputSamples[0].dts = firstDTS;
			// offset PTS as well, ensure that PTS is smaller or equal than new DTS
			firstPTS = math.Max(firstPTS-delta, *nextAvcDts);
			inputSamples[0].pts = firstPTS;
			//logger.log(`Video/PTS/DTS adjusted: ${math.round(firstPTS/90)}/${math.round(firstDTS/90)},delta:${delta} ms`);
		}
	}
	nextDTS = firstDTS

	// compute lastPTS/lastDTS
	sample = inputSamples[len(inputSamples)-1];
	lastDTS = math.Max(sample.dts, 0);
	lastPTS = math.Max(sample.pts, 0, lastDTS);

	// on Safari let's signal the same sample duration for all samples
	// sample duration (as expected by trun MP4 boxes), should be the delta between sample DTS
	// set this constant duration as being the avg delta between consecutive DTS.
	//if (isSafari) {
	//	mp4SampleDuration = math.round((lastDTS - firstDTS) / (inputSamples.length - 1));
	//}

	var nbNalu = 0
	var naluLen = 0
	for i := 0;	i < nbSamples; i++ {
	// compute total/avc sample length and nb of NAL units
	var sample = inputSamples[i]
	var units = sample.units
	var nbUnits = len(units)
	var sampleLen = 0
	for j := 0; j < nbUnits; j++ {
	sampleLen += len(units[j].data)
	}
	naluLen += sampleLen;
	nbNalu += nbUnits;
	sample.length = sampleLen;

	// normalize PTS/DTS
	//if (isSafari) {
	//// sample DTS is computed using a constant decoding offset (mp4SampleDuration) between samples
	//sample.dts = firstDTS + i*mp4SampleDuration;
	//} else {
	// ensure sample monotonic DTS
	sample.dts = math.Max(sample.dts, firstDTS);
	//}
	// ensure that computed value is greater or equal than sample DTS
	sample.pts = math.Max(sample.pts, sample.dts);
	}

	/* concatenate the video data and construct the mdat in place
	  (need 8 more bytes to fill length and mpdat type) */
	var mdatSize = naluLen + (4 * nbNalu) + 8;
	mdat= new([mdatSize]byte)
	binary.BigEndian.PutUint32(mdat, mdatSize);
	ArrayCopy(MP4.types.mdat,0, mdat,	4 , len(MP4.types.mdat))//todo implement MP$-generator
	for	i := 0;	i < nbSamples; i++ {
		var avcSample= inputSamples[i]
		var avcSampleUnits= avcSample.units
		var mp4SampleLength= 0
		var compositionTimeOffset

	// convert NALU bitstream to MP4 format (prepend NALU with size field)
	nbUnits := len(avcSampleUnits)
	for j := 0; j < nbUnits; j++ {
	var unit = avcSampleUnits[j]
	var unitData = unit.data
	var unitDataLen = len(unit.data)
	binary.BigEndian.PutUint32(mdat[offset:], uint32(unitDataLen));
	offset += 4;
	ArrayCopy(unitData,0, mdat, offset, unitDataLen) //mdat.set(unitData, offset);
	offset += unitDataLen;
	mp4SampleLength += 4 + unitDataLen;
	}

	//if (!isSafari) {//todo never safari
	// expected sample duration is the Decoding Timestamp diff of consecutive samples
	if (i < nbSamples - 1) {
	mp4SampleDuration = inputSamples[i+1].dts - avcSample.dts;
	} else {
	var config = _this.config
	var tempIndex
	if i>0 {
		tempIndex = i-1
	}else{
		i
	}
	var lastFrameDuration = avcSample.dts - inputSamples[tempIndex].dts;
	//if (config.stretchShortVideoTrack) {
	//// In some cases, a segment's audio track duration may exceed the video track duration.
	//// Since we've already remuxed audio, and we know how long the audio track is, we look to
	//// see if the delta to the next segment is longer than the minimum of maxBufferHole and
	//// maxSeekHole. If so, playback would potentially get stuck, so we artificially inflate
	//// the duration of the last frame to minimize any potential gap between segments.
	//var maxBufferHole = config.maxBufferHole
	//var maxSeekHole = config.maxSeekHole
	//var gapTolerance = math.Floor(math.Min(maxBufferHole, maxSeekHole) * timeScale)
	//var deltaToFrameEnd
	//deltaToFrameEnd = _this.nextAudioPts) - avcSample.pts
	//
	//if (deltaToFrameEnd > gapTolerance) {
	//// We subtract lastFrameDuration from deltaToFrameEnd to try to prevent any video
	//// frame overlap. maxBufferHole/maxSeekHole should be >> lastFrameDuration anyway.
	//mp4SampleDuration = deltaToFrameEnd - lastFrameDuration;
	//if (mp4SampleDuration < 0) {
	//mp4SampleDuration = lastFrameDuration;
	//}
	//logger.log(`It is approximately ${deltaToFrameEnd/90} ms to the next segment; using duration ${mp4SampleDuration/90} ms for the last video frame.`);
	//} else {
	//mp4SampleDuration = lastFrameDuration;
	//}
	//} else {
	//mp4SampleDuration = lastFrameDuration;
	//}
	}
	compositionTimeOffset = Round(avcSample.pts - avcSample.dts)
	//} else {//todo never safari
	//compositionTimeOffset = math.max(0, mp4SampleDuration*math.round((avcSample.pts - avcSample.dts)/mp4SampleDuration));
	//}

	var sample Mp4Sample

	if (avcSample.key) {
		sample = Mp4Sample{
			mp4SampleLength,
			compositionTimeOffset,
			-1,
			Flags{
				0,
				0,
				0,
				0,
				2,
				0,
			},
		}
	} else {
		sample = Mp4Sample{
			mp4SampleLength,
			compositionTimeOffset,
			-1,
			Flags{
				0,
				0,
				0,
				0,
				1,
				1,
			},
		}
	}

	outputSamples = append(outputSamples,sample)

	}
	// next AVC sample DTS should be equal to last sample DTS + last sample duration (in PES timescale)
	_this.nextAvcDts = lastDTS + mp4SampleDuration;
	var dropped = track.dropped;
	track.len = 0;
	track.nbNalu = 0;
	track.dropped = 0;
	if (outputSamples.length && navigator.userAgent.toLowerCase().indexOf('chrome
	') > -1) {
	let flags = outputSamples[0].flags;
	// chrome workaround, mark first sample as being a Random Access Point to avoid sourcebuffer append issue
	// https://code.google.com/p/chromium/issues/detail?id=229412
	flags.dependsOn = 2;
	flags.isNonSync = 0;
	}
	track.samples = outputSamples;
	moof = MP4.moof(track.sequenceNumber++, firstDTS, track);
	track.samples = [];
	let
	data =
	{
	data1:
		moof,
			data2: mdat,
		startPTS: firstPTS / timeScale,
		endPTS: (lastPTS + mp4SampleDuration) / timeScale,
		startDTS: firstDTS / timeScale,
		endDTS: this.nextAvcDts / timeScale,
		type: 'v
		ideo
		',
	nb:
		outputSamples.length,
			dropped : dropped
	};
	this.observer.trigger(Event.FRAG_PARSING_DATA, data);
	return data;
}

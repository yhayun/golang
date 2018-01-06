package main

import "math"

/**
 * Generate MP4 Box
*/
//import Hex from '../utils/hex';
const UINT32_MAX = math.MaxUint32;

type _Types struct {
  avc1 []byte
  avcC []byte
  btrt []byte
  dinf []byte
  dref []byte
  esds []byte
  ftyp []byte
  hdlr []byte
  mdat []byte
  mdhd []byte
  mdia []byte
  mfhd []byte
  minf []byte
  moof []byte
  moov []byte
  mp4a []byte
  mvex []byte
  mvhd []byte
  pasp []byte
  sdtp []byte
  stbl []byte
  stco []byte
  stsc []byte
  stsd []byte
  stsz []byte
  stts []byte
  tfdt []byte
  tfhd []byte
  traf []byte
  trak []byte
  trun []byte
  trex []byte
  tkhd []byte
  vmhd []byte
  smhd []byte
}


type MP4 struct {
  types _Types
  videoHdlr []uint8
  dref []uint8
  stco []uint8
  STTS []uint8
  STSC []uint8
  STCO []uint8
  STSZ []uint8
  VMHD []uint8
  SMHD []uint8
  STSD []uint8
  majorBrand []uint8
  avc1Brand []uint8
  minorVersion []uint8
  FTYP []uint8
  DINF []uint8
}

type Track struct {
	duration int // 0xffffffff
	timescale int
	id int
	width int
	height int
	_type string
	sps [][]byte
	pps [][]byte
	pixelRatio [2]byte
}



}
func (mp4* MP4) InitSegment(tracks Tracks) []byte{
  if (!mp4.types) {
    mp4.init();
  }
  var movie []byte = mp4.moov(tracks)
  var result []byte
  result = make([]uint8,len(mp4.FTYP) + len(movie));
  ArrayCopy(mp4.FTYP,0,result, 0 ,len(mp4.FTYP))
  ArrayCopy(movie,0,result, len(mp4.FTYP) ,len(movie))
  return result;
}
export default MP4;

init(mp4 *MP4) {
//var i;
//for (i in MP4.types) {
//if (MP4.types.hasOwnProperty(i)) {
//MP4.types[i] = []rune(types[i])
//  [
//  i.charCodeAt(0),
//  i.charCodeAt(1),
//  i.charCodeAt(2),
//  i.charCodeAt(3)
//];
//}
//}
  mp4.types.avc1=  [97, 118, 99, 49]
  mp4.types.avcC=  [97, 118, 99, 67]
  mp4.types.btrt=  [98, 116, 114, 116]
  mp4.types.dinf=  [100, 105, 110, 102]
  mp4.types.dref=  [100, 114, 101, 102]
  mp4.types.esds=  [101, 115, 100, 115]
  mp4.types.ftyp=  [102, 116, 121, 112]
  mp4.types.hdlr=  [104, 100, 108, 114]
  mp4.types.mdat=  [109, 100, 97, 116]
  mp4.types.mdhd=  [109, 100, 104, 100]
  mp4.types.mdia=  [109, 100, 105, 97]
  mp4.types.mfhd=  [109, 102, 104, 100]
  mp4.types.minf=  [109, 105, 110, 102]
  mp4.types.moof=  [109, 111, 111, 102]
  mp4.types.moov=  [109, 111, 111, 118]
  mp4.types.mp4a=  [109, 112, 52, 97]
  mp4.types.mvex=  [109, 118, 101, 120]
  mp4.types.mvhd=  [109, 118, 104, 100]
  mp4.types.pasp=  [112, 97, 115, 112]
  mp4.types.sdtp=  [115, 100, 116, 112]
  mp4.types.stbl=  [115, 116, 98, 108]
  mp4.types.stco=  [115, 116, 99, 111]
  mp4.types.stsc=  [115, 116, 115, 99]
  mp4.types.stsd=  [115, 116, 115, 100]
  mp4.types.stsz=  [115, 116, 115, 122]
  mp4.types.stts=  [115, 116, 116, 115]
  mp4.types.tfdt=  [116, 102, 100, 116]
  mp4.types.tfhd=  [116, 102, 104, 100]
  mp4.types.traf=  [116, 114, 97, 102]
  mp4.types.trak=  [116, 114, 97, 107]
  mp4.types.trun=  [116, 114, 117, 110]
  mp4.types.trex=  [116, 114, 101, 120]
  mp4.types.tkhd=  [116, 107, 104, 100]
  mp4.types.vmhd=  [118, 109, 104, 100]
  mp4.types.smhd=  [115, 109, 104, 100]
  mp4.videoHdlr = [
  0x00, // version 0
  0x00, 0x00, 0x00, // flags
  0x00, 0x00, 0x00, 0x00, // pre_defined
  0x76, 0x69, 0x64, 0x65, // handler_type: 'vide'
  0x00, 0x00, 0x00, 0x00, // reserved
  0x00, 0x00, 0x00, 0x00, // reserved
  0x00, 0x00, 0x00, 0x00, // reserved
  0x56, 0x69, 0x64, 0x65,
  0x6f, 0x48, 0x61, 0x6e,
  0x64, 0x6c, 0x65, 0x72, 0x00 // name: 'VideoHandler'
  ]

    //var audioHdlr = new Uint8Array([
    //  0x00, // version 0
    //  0x00, 0x00, 0x00, // flags
    //  0x00, 0x00, 0x00, 0x00, // pre_defined
    //  0x73, 0x6f, 0x75, 0x6e, // handler_type: 'soun'
    //  0x00, 0x00, 0x00, 0x00, // reserved
    //  0x00, 0x00, 0x00, 0x00, // reserved
    //  0x00, 0x00, 0x00, 0x00, // reserved
    //  0x53, 0x6f, 0x75, 0x6e,
    //  0x64, 0x48, 0x61, 0x6e,
    //  0x64, 0x6c, 0x65, 0x72, 0x00 // name: 'SoundHandler'
    //]);

    //MP4.HDLR_TYPES = { // todo always  videoHdlr
    //  'video': videoHdlr,
    //  'audio': audioHdlr
    //};

    mp4.dref = [
      0x00, // version 0
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, 0x00, 0x01, // entry_count
      0x00, 0x00, 0x00, 0x0c, // entry_size
      0x75, 0x72, 0x6c, 0x20, // 'url' type
      0x00, // version 0
      0x00, 0x00, 0x01 // entry_flags
    ]

     mp4.stco = [
      0x00, // version
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, 0x00, 0x00 // entry_count
    ]

    mp4.STTS = mp4.STSC = mp4.STCO = stco;

    mp4.STSZ = [
      0x00, // version
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, 0x00, 0x00, // sample_size
      0x00, 0x00, 0x00, 0x00, // sample_count
    ]

mp4.VMHD = [
      0x00, // version
      0x00, 0x00, 0x01, // flags
      0x00, 0x00, // graphicsmode
      0x00, 0x00,
      0x00, 0x00,
      0x00, 0x00 // opcolor
    ]
    mp4.SMHD = [
      0x00, // version
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, // balance
      0x00, 0x00 // reserved
    ]

    mp4.STSD = [
      0x00, // version 0
      0x00, 0x00, 0x00, // flags
      0x00, 0x00, 0x00, 0x01]// entry_count

    mp4.majorBrand = new Uint8Array([105,115,111,109]); // isom
    mp4.avc1Brand = new Uint8Array([97,118,99,49]); // avc1
    mp4.minorVersion = new Uint8Array([0, 0, 0, 1]);

mp4.FTYP = MP4.box(MP4.types.ftyp, mp4.majorBrand, mp4.minorVersion, mp4.majorBrand, mp4.avc1Brand);
mp4.DINF = MP4.box(MP4.types.dinf, MP4.box(MP4.types.dref, mp4.dref));
  }

func (mp4* MP4) Sdtp(track Track) []byte{
var samples = track.samples
var bytes = make([]byte,4 + len(samples))
var flags
var i
// leave the full box header (4 bytes) all zero
// write the sample table
for i = 0; i < len(samples); i++ {
flags = samples[i].flags;
bytes[i + 4] = (flags.dependsOn << 4) |
(flags.isDependedOn << 2) |
(flags.hasRedundancy);
}
return mp4.Box(mp4.types.sdtp, bytes);
}

func (mp4* MP4) Traf(track Track, baseMediaDecodeTime int) []byte{
    var sampleDependencyTable = mp4.Sdtp(track)
    var id = track.id
    var upperWordBaseMediaDecodeTime = int(math.Floor(float64(baseMediaDecodeTime / (UINT32_MAX + 1)))),
    var lowerWordBaseMediaDecodeTime = int(math.Floor(float64(baseMediaDecodeTime % (UINT32_MAX + 1))));
    return mp4.Box(mp4.types.traf,mp4.Box(mp4.types.tfhd, []byte{
  0x00,             // version 0
  0x00, 0x00, 0x00, // flags
  byte(id >> 24),
  byte(id >> 16) & 0XFF,
  byte(id >> 8) & 0XFF,
  byte(id & 0xFF)}) // track_ID
  ,mp4.Box(mp4.types.tfdt, []byte{
  0x01,             // version 1
  0x00, 0x00, 0x00, // flags
  byte(upperWordBaseMediaDecodeTime >>24),
  byte(upperWordBaseMediaDecodeTime >> 16) & 0XFF,
  byte(upperWordBaseMediaDecodeTime >> 8) & 0XFF,
  byte(upperWordBaseMediaDecodeTime & 0xFF),
  byte(lowerWordBaseMediaDecodeTime >>24),
  byte(lowerWordBaseMediaDecodeTime >> 16) & 0XFF,
  byte(lowerWordBaseMediaDecodeTime >> 8) & 0XFF,
  byte((lowerWordBaseMediaDecodeTime) & 0xFF)
  }),mp4.Trun(track,
                    len(sampleDependencyTable) +
                    16 + // tfhd
                    20 + // tfdt
                    8 +  // traf header
                    16 + // mfhd
                    8 +  // moof header
                    8),  // mdat header
               sampleDependencyTable);
}

func (mp4* MP4) Box(_type []byte, payload ...[]byte) []byte {
	size := 8
	i := len(payload)
	length := len(payload)

    // calculate the total size we need to allocate
      for (i >= 0) {
      size += len(payload[i])
      i--
    }

    var result = make([]byte,size)
    result[0] = (byte)(size >> 24) & 0xff;
    result[1] = (byte)(size >> 16) & 0xff;
    result[2] = (byte)(size >> 8) & 0xff;
    result[3] = (byte)(size)  & 0xff;
	ArrayCopy(_type, 0,  result, 4, len(_type)) //result.set(type, 4);
    // copy the payload into the result
    size = 8
    for i= 0; i < length; i++ {
      // copy payload[i] array @ offset size
      ArrayCopy(payload[i], 0, result, size, len(payload[i]))// result.set(payload[i], size);
      size += len(payload[i])
    }

    return result;
  }

func (mp4* MP4) hdlr(_type []byte) []byte {
	return mp4.Box(mp4.types.hdlr, mp4.videoHdlr);
}

func (mp4* MP4) stbl(track Track) []byte {
	return mp4.Box(mp4.types.stbl, mp4.stsd(track), mp4.Box(mp4.types.stts, mp4.STTS), mp4.Box(mp4.types.stsc, mp4.STSC), mp4.Box(mp4.types.stsz, mp4.STSZ), mp4.Box(mp4.types.stco, mp4.STCO));
}

func (mp4* MP4) moov(tracks []Track) []byte{
    var i = len(tracks)
    boxes := [][]byte{[]byte{}}
    for i>= 0 {
      boxes[i] = mp4.trak(tracks[i]);
      i--
    }
    return  append(append(mp4.types.moov, mp4.mvhd(tracks[0].timescale, tracks[0].duration)...),boxes...)
  }

  /**
   * Generate a track box.
   * @param track {object} a track definition
   * @return {Uint8Array} the track box
   */
func (mp4* MP4) trak(track Track) []byte {
    track.duration = track.duration || 0xffffffff;
    return mp4.Box(mp4.types.trak, mp4.tkhd(track), mp4.mdia(track));
  }

func (mp4* MP4) stsd(track Track) []byte {
    //if (track.type === 'audio') {
    //  if (!track.isAAC && track.codec === 'mp3') {
    //    return MP4.box(MP4.types.stsd, MP4.STSD, MP4.mp3(track));
    //  }
    //  return MP4.box(MP4.types.stsd, MP4.STSD, MP4.mp4a(track));
    //}
      return mp4.Box(mp4.types.stsd, mp4.STSD, mp4.avc1(track))
  }

func (mp4* MP4)  trex(track Track) []byte {
    var id = track.id;
    return mp4.Box(mp4.types.trex, []byte{
      0x00, // version 0
      0x00, 0x00, 0x00, // flags
		(byte)(id >> 24),
		(byte)(id >> 16) & 0XFF,
		(byte)(id >> 8) & 0XFF,
		(byte)(id & 0xFF), // track_ID
      0x00, 0x00, 0x00, 0x01, // default_sample_description_index
      0x00, 0x00, 0x00, 0x00, // default_sample_duration
      0x00, 0x00, 0x00, 0x00, // default_sample_size
      0x00, 0x01, 0x00, 0x01, // default_sample_flags
	});
  }

func (mp4* MP4) minf(track Track) []byte{
    //if (track._type == 'audio') {  <---- AUDIO ! WE DONT DO AUDIO!!
    //  return MP4.box(MP4.types.minf, MP4.box(MP4.types.smhd, MP4.SMHD), MP4.DINF, MP4.stbl(track));
    //}
      return mp4.Box(mp4.types.minf, mp4.Box(mp4.types.vmhd, mp4.VMHD), mp4.DINF, mp4.stbl(track));
  }

func (mp4* MP4) mvex(tracks []Track) []byte {
    var i = len(tracks)
	boxes := [][]byte{[]byte{}}

	for i>= 0 {
		boxes[i] = mp4.trex(tracks[i]);
		i--
	}
    return mp4.Box.apply(null, [mp4.types.mvex].concat(boxes));
  }


func (mp4* MP4) mdhd(timescale int , duration int ) []byte {
    duration *= timescale;
	upperWordDuration := int(math.Floor((float64(duration / (UINT32_MAX + 1)))))
	lowerWordDuration := int(math.Floor(float64(duration % (UINT32_MAX + 1))))
    return mp4.Box(mp4.types.mdhd,  []byte{
		0x01,                                           // version 1
		0x00, 0x00, 0x00,                               // flags
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // creation_time
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, // modification_time
		(byte)(timescale >> 24) & 0xFF,
		(byte)(timescale >> 16) & 0xFF,
		(byte)(timescale >> 8) & 0xFF,
		(byte)(timescale & 0xFF), // timescale
		(byte)(upperWordDuration >> 24),
		(byte)(upperWordDuration >> 16) & 0xFF,
		(byte)(upperWordDuration >> 8) & 0xFF,
		(byte)(upperWordDuration & 0xFF),
		(byte)(lowerWordDuration >> 24),
		(byte)(lowerWordDuration >> 16) & 0xFF,
		(byte)(lowerWordDuration >> 8) & 0xFF,
		(byte)(lowerWordDuration & 0xFF),
		0x55, 0xc4, // 'und' language (undetermined)
		0x00, 0x00,
	});
  }


func (mp4* MP4) mvhd (timescale int ,duration int) []byte {
    duration*=timescale
     upperWordDuration := int(math.Floor((float64(duration / (UINT32_MAX + 1)))))
     lowerWordDuration := int(math.Floor(float64(duration % (UINT32_MAX + 1))))
    var
      bytes =  []byte{
		0x01,                                           // version 1
		0x00, 0x00, 0x00,                               // flags
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // creation_time
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, // modification_time
		  (byte)(timescale >> 24) & 0xFF,
		  (byte)(timescale >> 16) & 0xFF,
		  (byte)(timescale >> 8) & 0xFF,
		  (byte)(timescale & 0xFF), // timescale
		  (byte)(upperWordDuration >> 24),
		  (byte)(upperWordDuration >> 16) & 0xFF,
		  (byte)(upperWordDuration >> 8) & 0xFF,
		  (byte)(upperWordDuration & 0xFF),
		  (byte)(lowerWordDuration >> 24),
		  (byte)(lowerWordDuration >> 16) & 0xFF,
	 	 (byte)(lowerWordDuration >> 8) & 0xFF,
	 	 (byte)(lowerWordDuration & 0xFF),
		0x00, 0x01, 0x00, 0x00, // 1.0 rate
		0x01, 0x00,             // 1.0 volume
		0x00, 0x00,             // reserved
		0x00, 0x00, 0x00, 0x00, // reserved
		0x00, 0x00, 0x00, 0x00, // reserved
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x40, 0x00, 0x00, 0x00, // transformation: unity matrix
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, // pre_defined
		0xff, 0xff, 0xff, 0xff, // next_track_ID
	};

    return mp4.Box(mp4.types.mvhd, bytes);
  }

func (mp4* MP4) mdia(track Track) []byte {
	return mp4.Box(mp4.types.mdia, MP4.mdhd(track.timescale, track.duration), MP4.hdlr(track.type), MP4.minf(track));
}

func (mp4* MP4) avc1(track Track) []byte {
    var sps []byte
    var pps []byte
    var length int
    var i int

    // assemble the SPSs
    for i = 0; i < len(track.sps); i++ {
      var data = track.sps[i]
	  length = len(data)
	  sps = append(sps,(byte)(length >> 8) & 0xFF)
	  sps = append(sps,(byte)(length & 0xFF))
	  sps = append(sps,data...)//SPS
    }
    // assemble the PPSs
    for i = 0; i < len(track.pps); i++ {
      var data = track.pps[i]
	  length = len(data)
	  pps = append(pps,(byte)(length >> 8) & 0xFF)
	  pps = append(pps,(byte)(length & 0xFF))
	  pps = append(pps,data...)//PPS
    }

	var avcc = mp4.Box(mp4.types.avcC, []byte{
		0x01,                   // version
		sps[3],                 // profile
		sps[4],                 // profile compat
		sps[5],                 // level
		0xfc | 3,               // lengthSizeMinusOne, hard-coded to 4 bytes
		0xE0 | (byte)(len(track.sps)), // 3bit reserved (111) + numOfSequenceParameterSets
	})
	avcc = append(avcc,sps...)
	avcc = append(avcc,(byte)(len(track.pps)))
	avcc = append(avcc, pps...)
        var width = track.width
        var height = track.height
        var hSpacing = track.pixelRatio[0]
        var vSpacing = track.pixelRatio[1]

    return mp4.Box(mp4.types.avc1,[]byte {
        0x00, 0x00, 0x00, // reserved
        0x00, 0x00, 0x00, // reserved
        0x00, 0x01, // data_reference_index
        0x00, 0x00, // pre_defined
        0x00, 0x00, // reserved
        0x00, 0x00, 0x00, 0x00,
        0x00, 0x00, 0x00, 0x00,
        0x00, 0x00, 0x00, 0x00, // pre_defined
		(byte) (width >> 8) & 0xFF,
		(byte)(width & 0xff), // width
		(byte) (height >> 8) & 0xFF,
		(byte)(height & 0xff), // height
        0x00, 0x48, 0x00, 0x00, // horizresolution
        0x00, 0x48, 0x00, 0x00, // vertresolution
        0x00, 0x00, 0x00, 0x00, // reserved
        0x00, 0x01, // frame_count
        0x12,
        0x64, 0x61, 0x69, 0x6C, //dailymotion/hls.js
        0x79, 0x6D, 0x6F, 0x74,
        0x69, 0x6F, 0x6E, 0x2F,
        0x68, 0x6C, 0x73, 0x2E,
        0x6A, 0x73, 0x00, 0x00,
        0x00, 0x00, 0x00, 0x00,
        0x00, 0x00, 0x00, 0x00,
        0x00, 0x00, 0x00, // compressorname
        0x00, 0x18,   // depth = 24
        0x11, 0x11}, // pre_defined = -1
          avcc,
          mp4.Box(mp4.types.btrt, []byte {
            0x00, 0x1c, 0x9c, 0x80, // bufferSizeDB
            0x00, 0x2d, 0xc6, 0xc0, // maxBitrate
            0x00, 0x2d, 0xc6, 0xc0}), // avgBitrate
          mp4.Box(mp4.types.pasp, []byte {
            (hSpacing >> 24),         // hSpacing
            (hSpacing >> 16) & 0xFF,
            (hSpacing >>  8) & 0xFF,
            hSpacing & 0xFF,
            (vSpacing >> 24),         // vSpacing
            (vSpacing >> 16) & 0xFF,
            (vSpacing >>  8) & 0xFF,
            vSpacing & 0xFF}),
          );
  }

func (mp4* MP4) tkhd(track Track) []byte {
    var id = track.id
    var duration = track.duration*track.timescale
    var width = track.width
    var height = track.height
	upperWordDuration := int(math.Floor((float64(duration / (UINT32_MAX + 1)))))
	lowerWordDuration := int(math.Floor(float64(duration % (UINT32_MAX + 1))))

    return mp4.Box(mp4.types.tkhd, []byte{
		0x01,                                           // version 1
		0x00, 0x00, 0x07,                               // flags
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, // creation_time
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, // modification_time
		(byte)(id >> 24) & 0xFF,
		(byte)(id >> 16) & 0xFF,
		(byte)(id >> 8) & 0xFF,
		(byte)(id & 0xFF),              // track_ID
		0x00, 0x00, 0x00, 0x00, // reserved
		(byte)(upperWordDuration >> 24),
		(byte)(upperWordDuration >> 16) & 0xFF,
		(byte)(upperWordDuration >> 8) & 0xFF,
		(byte)(upperWordDuration & 0xFF),
		(byte)(lowerWordDuration >> 24),
		(byte)(lowerWordDuration >> 16) & 0xFF,
		(byte)(lowerWordDuration >> 8) & 0xFF,
		(byte)(lowerWordDuration & 0xFF),
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, // reserved
		0x00, 0x00,             // layer
		0x00, 0x00,             // alternate_group
		0x00, 0x00,             // non-audio track volume
		0x00, 0x00,             // reserved
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x40, 0x00, 0x00, 0x00, // transformation: unity matrix
		(byte)(width >> 8) & 0xFF,
		(byte)(width & 0xFF),
		0x00, 0x00, // width
		(byte)(height >> 8) & 0xFF,
		(byte)(height & 0xFF),
		0x00, 0x00, // height
	});
  }


func (mp4 *MP4) Moof(sn int, baseMediaDecodeTime uint32, track Track) []byte{
return mp4.Box(mp4.types.moof, mp4.mfhd(sn), mp4.traf(track,baseMediaDecodeTime));
}

func (mp4 *MP4) Mdat(data []byte) []byte {
  return mp4.Box(mp4.types.mdat, data);
}

func (mp4 *MP4) Mfhd(sequenceNumber int) []byte{
return mp4.Box(mp4.types.mfhd, []byte{
0x00,
0x00, 0x00, 0x00, // flags
byte(sequenceNumber >> 24),
byte(sequenceNumber >> 16) & 0xFF,
byte(sequenceNumber >>  8) & 0xFF,
byte(sequenceNumber) & 0xFF});// sequence_number ));
}

func (mp4 *MP4) Trun(track Track, offset int) []byte{
    var samples= track.samples
    var lenn int = len(samples)
    var arraylen = 12 + (16 * lenn)
    var array = make([]uint8,arraylen)
    var i int
    var sample
    var duration
    var size
    var flags
    var cts
    offset += 8 + arraylen;
  var tempArr = []byte{0x00, // version 0
      0x00, 0x0f, 0x01, // flags
    (uint8)(lenn >> 24) & 0xFF,
    (uint8)(lenn >> 16) & 0xFF,
    (uint8)(lenn >> 8) & 0xFF,
    (uint8)(lenn & 0xFF), // sample_count
    (uint8)(offset >> 24) & 0xFF,
    (uint8)(offset >> 16) & 0xFF,
    (uint8)(offset >> 8) & 0xFF,
    (uint8)(offset & 0xFF)} // data_offset ]
  ArrayCopy(array,0,tempArr,0,12);
    for i = 0; i < lenn; i++ {
      sample = samples[i];
      duration = sample.duration;
      size = sample.size;
      flags = sample.flags;
      cts = sample.cts;
      var tempArr = []byte{
        (duration >> 24) & 0xFF,
        (duration >> 16) & 0xFF,
        (duration >> 8) & 0xFF,
        duration & 0xFF, // sample_duration
        (size >> 24) & 0xFF,
        (size >> 16) & 0xFF,
        (size >> 8) & 0xFF,
        size & 0xFF, // sample_size
        (flags.isLeading << 2) | flags.dependsOn,
        (flags.isDependedOn << 6) |
            (flags.hasRedundancy << 4) |
            (flags.paddingValue << 1) |
            flags.isNonSync,
        flags.degradPrio & 0xF0 << 8,
        flags.degradPrio & 0x0F, // sample_flags
        (cts >> 24) & 0xFF,
        (cts >> 16) & 0xFF,
        byte(cts >> 8) & 0xFF,
        cts & 0xFF };// sample_composition_time_offset

      ArrayCopy(tempArr,0,tempArr,12+16*i,len(tempArr))

    }
    return mp4.Box(mp4.types.trun, array);
  }


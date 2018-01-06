//// todo https://stackoverflow.com/questions/11184336/how-to-convert-from-byte-to-int-in-go-programming
//
///**
// * Parser for exponential Golomb codes, a variable-bitwidth number encoding scheme used by h264.
//*/
package main
//
//import (
//	"encoding/binary"
//)
//
//
//func min(x uint, y uint) uint {
//	if ( x > y) {
//		return y
//	}
//	return x
//}
//
//
//type ExpGolomb struct {
//	data []byte
//	bytesAvailable uint
//	word uint
//	bitsAvailable uint
//}
//
//func NewExpGolomb(data []byte) *ExpGolomb{
//	return &ExpGolomb{
//		data,
//		(uint)(len(data)), // the number of bytes left to examine in this.data
//		0,						// the current word being examined
//		0,					// the number of bits left to examine in the current word
//	}
//}
//
//func (this *ExpGolomb)  loadWord() {
//    position := (uint)(len(this.data)) - this.bytesAvailable
//    workingBytes := make([]byte, 32) //workingBytes = new Uint8Array(4),
//    availableBytes := min(4, this.bytesAvailable);
//	ArrayCopy(this.data[position:position+ availableBytes],0,workingBytes,0,(int)(availableBytes))
//	this.word = (uint)(binary.BigEndian.Uint32(workingBytes)) // this.word = new DataView(workingBytes.buffer).getUint32(0);
//    // track the amount of this.data that has been processed
//    this.bitsAvailable = availableBytes * 8;
//    this.bytesAvailable -= availableBytes;
//  }
//
//func (this *ExpGolomb) skipBits(count uint) {
//    var skipBytes uint
//    if (this.bitsAvailable > count) {
//      this.word <<= count;
//      this.bitsAvailable -= count;
//    } else {
//      count -= this.bitsAvailable;
//      skipBytes = count >> 3;
//      count -= (skipBytes >> 3);
//      this.bytesAvailable -= skipBytes;
//      this.loadWord();
//      this.word <<= count;
//      this.bitsAvailable -= count;
//    }
//  }
//
//func (this *ExpGolomb) readBits(size uint) uint {
//	bits := min(this.bitsAvailable, size) // :uint
//	value := this.word >> (32 - bits); // :uint
//    this.bitsAvailable -= bits;
//    if (this.bitsAvailable > 0) {
//      this.word <<= bits;
//    } else if (this.bytesAvailable > 0) {
//      this.loadWord();
//    }
//    bits = size - bits;
//    if (bits > 0 && this.bitsAvailable != 0) {
//      return value << bits | this.readBits(bits);
//    } else {
//      return value;
//    }
//  }
//
//
//
//  // ():uint
//func (this *ExpGolomb)  skipLZ() uint {
//    var leadingZeroCount uint  // :uint
//    for leadingZeroCount = 0; leadingZeroCount < this.bitsAvailable; leadingZeroCount++ {
//      if (0 != (this.word & (0x80000000 >> leadingZeroCount))) {
//        // the first bit of working word is 1
//        this.word <<= leadingZeroCount;
//        this.bitsAvailable -= leadingZeroCount;
//        return leadingZeroCount;
//      }
//    }
//    // we exhausted word and still have not found a 1
//    this.loadWord();
//    return leadingZeroCount + this.skipLZ();
//  }
//
//  // ():void
//func (this *ExpGolomb) skipUEG() {
//    this.skipBits(1 + this.skipLZ());
//  }
//
//  // ():void
//func (this *ExpGolomb) skipEG() {
//    this.skipBits(1 + this.skipLZ());
//  }
//
//  // ():uint
//func (this *ExpGolomb) readUEG() uint {
//    var clz = this.skipLZ(); // :uint
//    return this.readBits(clz + 1) - 1;
//  }
//
//  // ():int
//func (this *ExpGolomb) readEG() uint{
//    var value = this.readUEG(); // :int
//    if (0x01 & value != 0) {
//      // the number is odd if the low order bit is set
//      return (1 + value) >> 1; // add 1 to make it even, and divide by 2
//    } else {
//      return -1 * (value >> 1); // divide by two then make it negative
//    }
//  }
//
//  // Some convenience functions
//  // :Boolean
//func (this *ExpGolomb)  readBoolean() bool {
//    return 1 == this.readBits(1);
//  }
//
//  // ():int
//func (this *ExpGolomb)  readUByte() uint {
//    return this.readBits(8);
//  }
//
//  // ():int
//func (this *ExpGolomb) readUShort()uint{
//    return this.readBits(16);
//  }
//    // ():int
//func (this *ExpGolomb)  readUInt() uint {
//    return this.readBits(32);
//  }
//
//
//  /**
//   * Advance the ExpGolomb decoder past a scaling list. The scaling
//   * list is optionally transmitted as part of a sequence parameter
//   * set and is not relevant to transmuxing.
//   * @param count {number} the number of entries in this scaling list
//   * @see Recommendation ITU-T H.264, Section 7.3.2.1.1.1
//   */
//func (this *ExpGolomb)  skipScalingList(count int) {
//	var lastScale uint = 8
//	var nextScale uint = 8
//
//	for j := 0; j < count; j++ {
//		if (nextScale != 0) {
//			deltaScale := this.readEG();
//			nextScale = (lastScale + deltaScale + 256) % 256;
//		}
//		if (nextScale == 0) {
//			lastScale = lastScale
//		} else {
//			lastScale = nextScale;
//		}
//	}
//}
//
//
//
//  /**
//   * Read a sequence parameter set and return some interesting video
//   * properties. A sequence parameter set is the H264 metadata that
//   * describes the properties of upcoming video frames.
//   * @param data {Uint8Array} the bytes of a sequence parameter set
//   * @return {object} an object with configuration parsed from the
//   * sequence parameter set, including the dimensions of the
//   * associated video frames.
//   */
//func (this *ExpGolomb) readSPS() Config {
//	frameCropLeftOffset := 0
//	frameCropRightOffset := 0
//	frameCropTopOffset := 0
//	frameCropBottomOffset := 0
//	var scalingListCount, numRefFramesInPicOrderCntCycle, picWidthInMbsMinus1,
//		picHeightInMapUnitsMinus1, frameMbsOnlyFlag int
//    //	var
//      //profileIdc,
//      //profileCompat,
//      //levelIdc,
//      //numRefFramesInPicOrderCntCycle,
//      //picWidthInMbsMinus1,
//      //picHeightInMapUnitsMinus1,
//      //frameMbsOnlyFlag,
//      //scalingListCount,
//      //i,
//    this.readUByte();
//    profileIdc := this.readUByte(); // profile_idc
//    this.readBits(5); // constraint_set[0-4]_flag, u(5)
//    this.skipBits(3); // reserved_zero_3bits u(3),
//    this.readUByte(); //level_idc u(8)
//    this.skipUEG(); // seq_parameter_set_id
//    // some profiles have more optional data we don't need
//    if (profileIdc == 100 ||
//        profileIdc == 110 ||
//        profileIdc == 122 ||
//        profileIdc == 244 ||
//        profileIdc == 44  ||
//        profileIdc == 83  ||
//        profileIdc == 86  ||
//        profileIdc == 118 ||
//        profileIdc == 128) {
//      var chromaFormatIdc = this.readUEG();
//      if (chromaFormatIdc == 3) {
//        this.skipBits(1); // separate_colour_plane_flag
//      }
//      this.skipUEG(); // bit_depth_luma_minus8
//      this.skipUEG(); // bit_depth_chroma_minus8
//      this.skipBits(1); // qpprime_y_zero_transform_bypass_flag
//      if (this.readBoolean()) { // seq_scaling_matrix_present_flag
//      	if (chromaFormatIdc != 3) {
//			scalingListCount = 8
//		  }else {
//			scalingListCount = 12
//		}
//        for i := 0; i < scalingListCount; i++ {
//          if (this.readBoolean()) { // seq_scaling_list_present_flag[ i ]
//            if (i < 6) {
//              this.skipScalingList(16);
//            } else {
//              this.skipScalingList(64);
//            }
//          }
//        }
//      }
//    }
//    this.skipUEG(); // log2_max_frame_num_minus4
//    var picOrderCntType = this.readUEG();
//    if (picOrderCntType == 0) {
//      this.readUEG(); //log2_max_pic_order_cnt_lsb_minus4
//    } else if (picOrderCntType == 1) {
//      this.skipBits(1); // delta_pic_order_always_zero_flag
//      this.skipEG(); // offset_for_non_ref_pic
//      this.skipEG(); // offset_for_top_to_bottom_field
//      numRefFramesInPicOrderCntCycle = (int)(this.readUEG())
//      for i := 0; i < numRefFramesInPicOrderCntCycle; i++ {
//        this.skipEG(); // offset_for_ref_frame[ i ]
//      }
//    }
//    this.skipUEG(); // max_num_ref_frames
//    this.skipBits(1); // gaps_in_frame_num_value_allowed_flag
//    picWidthInMbsMinus1 = (int)(this.readUEG())
//    picHeightInMapUnitsMinus1 = (int)(this.readUEG())
//    frameMbsOnlyFlag = (int)(this.readBits(1))
//    if (frameMbsOnlyFlag == 0) {
//		this.skipBits(1); // mb_adaptive_frame_field_flag
//    }
//	this.skipBits(1); // direct_8x8_inference_flag
//    if (this.readBoolean()) { // frame_cropping_flag
//      frameCropLeftOffset = (int)(this.readUEG())
//      frameCropRightOffset = (int)(this.readUEG())
//      frameCropTopOffset = (int)(this.readUEG())
//      frameCropBottomOffset = (int)(this.readUEG())
//    }
//     pixelRatio := [2]uint{1,1}
//    if (this.readBoolean()) {
//      // vui_parameters_present_flag
//      if (this.readBoolean()) {
//        // aspect_ratio_info_present_flag
//         aspectRatioIdc := this.readUByte();
//        switch (aspectRatioIdc) {
//          case 1: pixelRatio = [2]uint{1,1}; break;
//          case 2: pixelRatio = [2]uint{12,11}; break;
//          case 3: pixelRatio = [2]uint{10,11}; break;
//          case 4: pixelRatio = [2]uint{16,11}; break;
//          case 5: pixelRatio = [2]uint{40,33}; break;
//          case 6: pixelRatio = [2]uint{24,11}; break;
//          case 7: pixelRatio = [2]uint{20,11}; break;
//          case 8: pixelRatio = [2]uint{32,11}; break;
//          case 9: pixelRatio = [2]uint{80,33}; break;
//          case 10: pixelRatio = [2]uint{18,11}; break;
//          case 11: pixelRatio = [2]uint{15,11}; break;
//          case 12: pixelRatio = [2]uint{64,33}; break;
//          case 13: pixelRatio = [2]uint{160,99}; break;
//          case 14: pixelRatio = [2]uint{4,3}; break;
//          case 15: pixelRatio = [2]uint{3,2}; break;
//          case 16: pixelRatio = [2]uint{2,1}; break;
//          case 255: {
//            pixelRatio = [2]uint{this.readUByte() << 8 | this.readUByte(), this.readUByte() << 8 | this.readUByte()};
//            break;
//          }
//        }
//      }
//    }
//
//	var height int
//	if (frameMbsOnlyFlag != 0) {
//		height =((2 - frameMbsOnlyFlag) * (picHeightInMapUnitsMinus1 + 1) * 16) - (2) * (frameCropTopOffset + frameCropBottomOffset)
//	}else {
//		height =((2 - frameMbsOnlyFlag) * (picHeightInMapUnitsMinus1 + 1) * 16) - (4) * (frameCropTopOffset + frameCropBottomOffset)
//	}
//
//    return Config{
//		((picWidthInMbsMinus1 + 1) * 16) - frameCropLeftOffset * 2 - frameCropRightOffset * 2,
//    	height,
//    	pixelRatio,
//	}
//  }
//
//
//func (this *ExpGolomb)  readSliceType() uint {
//    // skip NALu type
//    this.readUByte();
//    // discard first_mb_in_slice
//    this.readUEG();
//    // return slice_type
//    return this.readUEG();
//}
//
//
//

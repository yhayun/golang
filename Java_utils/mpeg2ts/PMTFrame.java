/**
 * 
 */
package com.ibm.mpeg2ts;

import java.util.HashMap;
import java.util.Map;



import com.ibm.mpeg2ts.Mpeg2TSPacket.Mpeg2TSPacketType;
import com.ibm.mpeg2ts.StreamInfo.TS_STREAM_TYPE;


public class PMTFrame{

	private	byte []data = new byte[1280];
	private int size = 0;
	private int expectedSize = -1;

	private int programPID = -1;
	private int continuityCounter = -1;

	public PMTFrame(){

	}

	private TS_STREAM_TYPE getStreamType(byte streamTypeByte){
		switch(streamTypeByte){
		case 0x02:
			return TS_STREAM_TYPE.TS_STREAM_TYPE_MPEG2V;
		case 0x1B:
			return TS_STREAM_TYPE.TS_STREAM_TYPE_H264;
		case 0x03:
			return TS_STREAM_TYPE.TS_STREAM_TYPE_MP3;
		case 0xF:
			return TS_STREAM_TYPE.TS_STREAM_TYPE_AAC;
		default:			
			return TS_STREAM_TYPE.TS_STREAM_TYPE_UNKNOWN;

		}
	}

	public void setProgramPID(int programID){this.programPID = programID;}
	public int getProgramPID(){return programPID;}
	public void setExpectedSize(int expectedSize){this.expectedSize = expectedSize;}
	public int getExpectedSize(){return expectedSize;}

	public boolean addPacket(Mpeg2TSPacket p){
		if(p.getType() != Mpeg2TSPacketType.PMT){
			return false;
		}

		if(continuityCounter == -1 && !p.isStart()){
			return false;
		}

		if(p.isStart()){
			expectedSize = p.getPMTLength();
		}

		if(continuityCounter != -1 && p.getContinuityCounter() != (continuityCounter + 1)%Mpeg2TSPacket.MAX_CONTINUITY_COUNTER){
			size = 0;
			continuityCounter = -1;
			return false;
		}else{
			System.arraycopy(p.getData(), p.getDataOffset(), data, size, p.getDataLength());
			size += p.getDataLength();
			continuityCounter = p.getContinuityCounter();
			if(size >= expectedSize){
				return true;
			}
		}

		return false;
	}

	public Map<Integer, StreamInfo> getStreamInfos(){
		HashMap<Integer, StreamInfo> result = new HashMap<Integer, StreamInfo>();

		int offset = 2;

		int programMapSectionLength= ((data[offset]&0x0f) <<8) + data[offset+1];

		//Skip everything till start of streams array
		offset+=9;
		int programInfoLength = ((data[offset]&0xf) <<8) + data[offset+1];
		offset+=programInfoLength+2;

		while(offset < programMapSectionLength-1){
			StreamInfo streamInfo = new StreamInfo();

			streamInfo.streamType = getStreamType(data[offset]);
			offset +=1;		  
			streamInfo.streamPID = (((data[offset] & 0x1F) << 8)& 0x0000ffff) + (data[offset+1]& 0x00ff);
			offset += 2;
			int esInfoLength = ((data[offset]&0xf)<<8)+data[offset+1];
			offset += esInfoLength +2;
			result.put(streamInfo.streamPID, streamInfo);
		}

		return result;

	}
	
	public int getPcrPID(){
		int pcrPID = (((data[9] & 0x1F) << 8) & 0x0000ffff) + (data[10] & 0x00ff);
		return pcrPID;
	}

}
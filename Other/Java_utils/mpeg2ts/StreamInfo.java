/**
 * 
 */
package com.ibm.mpeg2ts;



/**
 * @author evgenyh
 *
 */
public class StreamInfo {

	public enum TS_STREAM_TYPE{
		TS_STREAM_TYPE_MPEG2V,
		TS_STREAM_TYPE_H264,
		TS_STREAM_TYPE_AAC,
		TS_STREAM_TYPE_MP3,
		TS_STREAM_TYPE_AC3,
		TS_STREAM_TYPE_G711,	
		TS_STREAM_TYPE_UNKNOWN
	}

	public TS_STREAM_TYPE streamType;
	public int streamPID;
    
	@Override
	public String toString() {

		StringBuilder builder = new StringBuilder();
		
		builder.append("Type: "+streamType.name().substring(streamType.name().lastIndexOf("_")+1)+"\n");
		
		builder.append("PID: "+streamPID+"\n");
		
		return builder.toString();
		
	}
	
}

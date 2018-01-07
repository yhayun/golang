package com.ibm.mp4;

import com.ibm.frame.*;
import com.ibm.mpeg2ts.*;
import com.ibm.mpeg2ts.Mpeg2TSPacket.Mpeg2TSPacketType;
import com.ibm.mpeg2ts.StreamInfo.TS_STREAM_TYPE;
import com.ibm.utils.Sleeper;

import java.io.File;
import java.io.FileInputStream;
import java.io.FileNotFoundException;
import java.io.FileOutputStream;
import java.io.IOException;
import java.net.DatagramPacket;
import java.nio.ByteBuffer;
import java.nio.channels.FileChannel;
import java.util.HashMap;
import java.util.Map;

import org.apache.log4j.Logger;

import com.coremedia.iso.boxes.Container;
import com.googlecode.mp4parser.FileDataSourceImpl;
import com.googlecode.mp4parser.MemoryDataSourceImpl;
import com.googlecode.mp4parser.authoring.Movie;
import com.googlecode.mp4parser.authoring.builder.DefaultMp4Builder;
import com.googlecode.mp4parser.authoring.builder.FragmentedMp4Builder;
import com.googlecode.mp4parser.authoring.tracks.H264TrackImpl;


public class TestMP4 {


	private static Logger log = Logger.getLogger(TestMP4.class);

	private static final int VIDEO_FRAME_CAPACITY = 1000000;

	private byte[] h264Buffer = new byte[60*1024*1024];
	private int length = 0;
	
	
	private byte[] mp4Buffer = new byte[60*1024*1024];

	/**
	 * @param args
	 */
	public static void main(String[] args) {
		new TestMP4().start();
	}

	private boolean checkIfIFrame(byte[] buffer, int offset, int length) {


		int startOffset = offset;

		for (int i = startOffset; i < (length - 7); i++) {
			int j = offset + i;
			if (buffer[j] == 0) {
				if (buffer[j + 1] == 0) {
					if (buffer[j + 2] == 1) {
						short val = (short) (0x001f & buffer[j + 3]);
						if (val == 7) {
							return true; // Start GOP was detected
						}
					}
				}
			}
		}
		return false;
	}
	
	public void start() {
		PooledQueue<Frame> videoFrames = new FrameQueue(100, VIDEO_FRAME_CAPACITY);
		new Thread(new H264Source("d:/videos/RR_Clock_H264/Elad_H264.mpg", videoFrames)).start();
		boolean iframeDetected = false;
		int numIframes = 0;
		while(true) {
			if(videoFrames.isEmpty()) {
				Sleeper.sleep(100);
				if(videoFrames.isEmpty()) {
					Sleeper.sleep(8000);
					if(videoFrames.isEmpty()) {						
						break;
					}
				}
			}
			Frame frame = videoFrames.poll();
			if(!iframeDetected) {
				if(checkIfIFrame(frame.getData(), 0, frame.size())) {
					iframeDetected = true;
				}
				else {
					videoFrames.recycle(frame);
					continue;
				}
			}
			if(checkIfIFrame(frame.getData(), 0, frame.size())) {
				numIframes++;
			}
			if(numIframes >= 2) {
				videoFrames.recycle(frame);
				continue;
			}
			
			//System.out.println(frame.size());
			System.arraycopy(frame.getData(), 0, h264Buffer, length, frame.size());
			length += frame.size();
			videoFrames.recycle(frame);
		}
		
		
		
		/*try {
			File file = new File("d:/videos/RR_Clock_H264/out.264");
			FileInputStream is = new FileInputStream(file);
			is.read(h264Buffer, 0, (int)file.length());
			length = (int) file.length();
			System.out.println("Length = " + length);
			is.close();
		} catch (Exception e1) {
			// TODO Auto-generated catch block
			e1.printStackTrace();
		}*/
		
		
		//ByteBuffer byteBuffer = ByteBuffer.wrap(h264Buffer, 0, length);
		ByteBuffer byteBuffer = ByteBuffer.allocate(length);
		byteBuffer.put(h264Buffer, 0, length);
		MemoryDataSourceImpl dataSource = new MemoryDataSourceImpl(byteBuffer);
		
	
		try {
			
			//H264TrackImpl h264Track = new H264TrackImpl(new FileDataSourceImpl("d:/videos/RR_Clock_H264/out.264"));
			H264TrackImpl h264Track = new H264TrackImpl(dataSource);
			Movie movie = new Movie();
			movie.addTrack(h264Track);
			Container mp4file = new FragmentedMp4Builder().build(movie);
			
			MemoryChannel channel = new MemoryChannel(mp4Buffer);
			mp4file.writeContainer(channel);
			channel.close();
			
			FileOutputStream os = new FileOutputStream(new File("d:/videos/RR_Clock_H264/Elad_H264.mp4"));
			os.write(mp4Buffer, 0, channel.getLength());
			os.close();
			
			//FileChannel fc = new FileOutputStream(new File("d:/videos/RR_Clock_H264/Elad_H264.mp4")).getChannel();
			
			
		} catch (Exception e) {
			//log.fatal(e);
			e.printStackTrace();
			System.err.println(e.getMessage() + " " + e.toString());
		}
		
	}


	public class H264Source implements Runnable {

		private PMTFrame pmtFrame = new PMTFrame();
		private Map<Integer, StreamInfo> streamsMap = new HashMap<Integer, StreamInfo>();
		private int programPID = 0;
		private static final int MPEG2TS_PACKET_LENGTH = 188;
		private boolean endFlag = false;
		private boolean detectFlag = false;
		private int videoPID;
		private String fileName;
		private PooledQueue<Frame> output;

		public H264Source(String fileName, PooledQueue<Frame> output) {
			this.fileName = fileName;
			this.output = output;
		}


		@Override
		public void run() {
			try {
				Mpeg2TSParser videoTSParser = new Mpeg2TSParser(output);
				Mpeg2TSPacket tsPacket = new Mpeg2TSPacket();

				File file = new File(fileName);
				FileInputStream is = new FileInputStream(file);
				long numPackets = file.length() / MPEG2TS_PACKET_LENGTH;

				byte[] packet = new byte[MPEG2TS_PACKET_LENGTH];

				for (int index = 0; index < numPackets; index++) {

					is.read(packet, 0, MPEG2TS_PACKET_LENGTH);
					tsPacket.fromBytes(packet, 0, programPID );

					if(detectFlag) {
						if (tsPacket.getPID() == videoPID) {
							videoTSParser.write(tsPacket);
						}
					}
					else {
						if(tsPacket.getType() == Mpeg2TSPacketType.PAT) {
							programPID = tsPacket.getProgramPID();
						}

						if(tsPacket.getType() == Mpeg2TSPacketType.PMT){ //PMT received

							if(pmtFrame.addPacket(tsPacket)){
								//pmtFrame is complete
								streamsMap = pmtFrame.getStreamInfos();

								for (Map.Entry<Integer, StreamInfo> entry : streamsMap.entrySet()) {
									StreamInfo info = (StreamInfo)entry.getValue();
									if(info.streamType == TS_STREAM_TYPE.TS_STREAM_TYPE_H264) { 
										detectFlag = true;
										videoPID = info.streamPID;
										break;
									}

								}
							}
						}

					}
				}
				is.close();
			}
			catch (Exception e) {
				//log.fatal(e);
				System.err.println(e);
			}
		}
	}

}

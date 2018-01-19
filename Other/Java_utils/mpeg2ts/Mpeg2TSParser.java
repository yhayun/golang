package com.ibm.mpeg2ts;


import org.apache.log4j.Logger;



import com.ibm.frame.*;
import com.ibm.utils.Sleeper;


/**
 * Collects Mpeg2TSPackets and builds a single frame (access unit) out of them.
 * Frames are added to a PooledQueue given upon creation.
 */
public class Mpeg2TSParser {
	
	private static Logger log = Logger.getLogger(Mpeg2TSParser.class);

	private static final int MAX_COUNTER = 16;

	private PooledQueue<Frame> output;
	private Frame currentFrame;
	private int counter;
	private boolean iFrameFound = false;
	private long lastFrameTime = 0; 
	
	private boolean endFlag = false;

	public Mpeg2TSParser(PooledQueue<Frame> output) {
		this.output = output;
		counter = -1;
	}

	public void close() {
		endFlag = true;
	}
	
	private int getPreviousCounter(int counter) {
		int prevCounter = counter - 1;
		if(prevCounter == -1) {
			return 15;
		}
		return prevCounter;
	}
	
	/**
	 * Takes a single Mpeg2TSPacket and adds it to the current frame. If the
	 * packet is the start of a new frame, adds the current frame to the output
	 * and starts a new frame using this packet.
	 * 
	 * @param packet
	 *            The Mpeg2TSPacket to write
	 */
	public void write(Mpeg2TSPacket packet) {
		
		int packetCounter = packet.getContinuityCounter();	
		
		if(!packet.isPayloadExist()) {
			//log.info("Payload doens't exist");
			return;
		}
		
		if(packetCounter == getPreviousCounter(counter)) {
			log.info("Duplicate packet found");
			return;
		}
		
		if (currentFrame == null) {
			currentFrame = output.newElement();			
			currentFrame.clear();
		}
		
		if (currentFrame == null) {
			log.info("currentFrame is null");
			Sleeper.sleep(5);
			return;
		}
		
		
		
		if (packet.isStartOfPES()) {
			if (!currentFrame.isEmpty()) {
				
				if(packetCounter == counter) {
					long startTime = System.currentTimeMillis();
					//log.info("before offer took");
					while (!output.offer(currentFrame)) {
						try {
							if(endFlag) {
								break;
							}
							Thread.sleep(5);
						} catch (InterruptedException e) {							;
							break;
						}
					}
					///long diff = System.currentTimeMillis() - lastFrameTime;
					//lastFrameTime = System.currentTimeMillis();
					//log.info("After.  " + diff);
					//log.info("Offer took " + (System.currentTimeMillis() - startTime));
					currentFrame = output.newElement();
					currentFrame.clear(); 
					
					if (currentFrame == null) {
						log.info("(currentFrame == null)");
						return;
					}
					
				}
				else {
					log.info("currentFrame.clear()");
					currentFrame.clear();
				}
				
			}		
			currentFrame.setPTS(packet.getPTS()/90);
		} else if (currentFrame.isEmpty() || packetCounter != counter) {
			log.info("Continuity = " + packetCounter + " counter = " + counter);
			counter = (packetCounter + 1) % MAX_COUNTER;
			return;
		}
		
		
		if (packet.isPayloadExist()) {			
			if(currentFrame.isEmpty()) {
				currentFrame.setCurrentSize(12);
			}
			currentFrame.append(packet.getData(), packet.getPayloadOffset(),
					packet.getPayloadLength());
			counter = (packetCounter + 1) % MAX_COUNTER;
		}
	}

	/**
	 * Deletes the content of the current frame and recycles it.
	 */
	public void flush() {
		output.recycle(currentFrame);
		currentFrame = null;
	}

	

}

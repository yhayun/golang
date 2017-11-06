package com.ibm.mpeg2ts;



public class Mpeg2TSPacket {

	public enum Mpeg2TSPacketType {

		PAT, PMT, PES

	}

	public enum VideoFrameType {

		I_START, P_START, OTHER

	}

	public static final int TRANSPORT_PACKET_SIZE = 188;
	public static final int MAX_CONTINUITY_COUNTER = 16;

	private byte[] data;

	private int offset;

	private boolean valid;

	private int pid;

	private Mpeg2TSPacketType type;

	private int continuityCounter;

	private int dataOffset;

	private int dataLength;

	private boolean adaptationExists;

	private boolean pcrExists;

	private long pcr;

	// in case of PES
	private boolean startOfPES;

	private long pts;

	private boolean payloadExist;

	private int payloadOffset;

	private int payloadLength;

	// in case of PAT
	private int programPID;

	// in case of PMT
	private int pmtLength;
	private boolean start;

	private long receiveTime;

	public Mpeg2TSPacket() {
	}

	public void fromBytes(byte[] data, int offset, int programPID) {

		reset();

		if (!validMpegTsPacket(data, offset)) {
			valid = false;
			return;
		}

		this.data = data;
		this.offset = offset;

		valid = true;

		pid = getPID(data, offset);

		dataOffset = tsDataOffset(data, offset);
		dataLength = TRANSPORT_PACKET_SIZE + offset - dataOffset;
		continuityCounter = getContinuityCounter(data, offset);
		adaptationExists = isAdaptationFieldExist(data, offset);
		if (adaptationExists && isPCRExist(data, offset)) {
			pcrExists = true;
			pcr = readPCR(data, offset);
		}
		
		/*if (pid != 512) {
			Log.i(getClass().getName(), "pid = " + pid);
		}*/

		if (pid == 0) {
			type = Mpeg2TSPacketType.PAT;
			this.programPID = getProgramPID(data, offset);
		} else if (pid == programPID) {
			type = Mpeg2TSPacketType.PMT;
			start = isStart(data, offset);
			if (start) {
				pmtLength = getPMTLength(data, offset);
			}
		} else {
			type = Mpeg2TSPacketType.PES;
			startOfPES = isStartOfPES(data, offset);
			payloadExist = payloadExists(data, offset);
			if (payloadExist) {
				payloadOffset = getPayloadOffset(data, offset);
				payloadLength = TRANSPORT_PACKET_SIZE + offset - payloadOffset;
			}
			if (startOfPES) {
				pts = getPTS(data, offset);
			}
		}

	}

	private void reset() {
		valid = false;
		pid = 0;
		type = null;
		continuityCounter = 0;
		dataOffset = 0;
		dataLength = 0;
		adaptationExists = false;
		pcrExists = false;
		pcr = 0;
		startOfPES = false;
		pts = 0;
		payloadExist = false;
		payloadOffset = 0;
		payloadLength = 0;
		programPID = 0;
		pmtLength = 0;
		start = false;
		receiveTime = 0;
	}

	// getters
	public byte[] getData() {
		return data;
	}

	public boolean isValid() {
		return valid;
	}

	public int getPID() {
		return pid;
	}

	public Mpeg2TSPacketType getType() {
		return type;
	}

	public int getContinuityCounter() {
		return continuityCounter;
	}

	public int getDataOffset() {
		return dataOffset;
	}

	public int getDataLength() {
		return dataLength;
	}

	public boolean isAdaptationExist() {
		return adaptationExists;
	}

	public long getPCR() {
		return pcr;
	}

	/**
	 * @return the pcrExists
	 */
	public boolean isPcrExists() {
		return pcrExists;
	}

	public boolean isStartOfPES() {
		return startOfPES;
	}

	public long getPTS() {
		return pts;
	}

	public boolean isPayloadExist() {
		return payloadExist;
	}

	public int getPayloadOffset() {
		return payloadOffset;
	}

	public int getPayloadLength() {
		return payloadLength;
	}

	public int getProgramPID() {
		return programPID;
	}

	public int getPMTLength() {
		return pmtLength;
	};

	public boolean isStart() {
		return start;
	}

	public long getReceiveTime() {
		return receiveTime;
	}

	public void setPTS(long pts) {

		setPTS(data, pts, offset);

	}

	private boolean isAdaptationFieldExist(byte[] buffer, int offset) {
		return ((buffer[3 + offset] & 0x20) != 0);
	}

	private boolean isPCRExist(byte[] buffer, int offset) {

		if (!isAdaptationFieldExist(buffer, offset)
				|| adaptationFieldLength(buffer, offset) == 0) {
			return false;
		}

		if ((buffer[5 + offset] & 0x10) != 0) {
			return true;
		}

		return false;
	}

	private int adaptationFieldLength(byte[] buffer, int offset) {
		int length = (buffer[4 + offset] & 0x00ff);
		return length;
	}

	private int tsDataOffset(byte[] buffer, int offset) {
		int internalOffset = 4 + offset;

		if (isAdaptationFieldExist(buffer, offset))
			internalOffset += 1 + adaptationFieldLength(buffer, offset);

		return internalOffset;
	}

	private boolean isStartOfPES(byte[] buffer, int offset) {
		if ((buffer[1 + offset] & 0x40) == 0)
			return false;

		int tsDataOffset = tsDataOffset(buffer, offset);

		boolean b = (buffer[tsDataOffset] == 0 && buffer[tsDataOffset + 1] == 0 && buffer[tsDataOffset + 2] == 0x01);
		return b;
	}

	private long getPTS(byte[] buffer, int offset) {

		int dataOffset = tsDataOffset(buffer, offset);

		if ((buffer[7 + dataOffset] & 0x80) == 0) {
			return -1;
		}

		int ptsOffset = 9 + dataOffset;

		long pts;

		pts = ((long) ((buffer[ptsOffset] & 0x0e) >> 1)) << 30;
		pts += ((long) (buffer[1 + ptsOffset] & 0xff) << 22);
		pts += ((long) ((buffer[2 + ptsOffset] & 0xfe) >> 1)) << 15;
		pts += ((long) (buffer[3 + ptsOffset] & 0xff) << 7);
		pts += ((buffer[4 + ptsOffset] & 0xfe) >> 1);
		return pts;
	}

	// private long getDTS(byte []buffer, int offset) {
	//
	// int tsDataOffset = tsDataOffset(buffer, offset)+ 14;
	//
	// long dts = (buffer[tsDataOffset] & 0x0e) >> 1;
	// dts = (dts << 15) + ((buffer[1 + tsDataOffset] & 0xff) << 7) + ((buffer[2
	// + tsDataOffset] & 0xfe) >> 1);
	// dts = (dts << 15) + ((buffer[3 + tsDataOffset] & 0xff) << 7) + ((buffer[4
	// + tsDataOffset] & 0xfe) >> 1);
	//
	// return dts;
	// }

	private void setPTS(byte[] buffer, long pts, int offset) {

		int ptsOffset = 9 + tsDataOffset(buffer, offset);
		buffer[ptsOffset] = (byte) ((0x2 << 4) + ((pts >> 30) << 1) + 1);
		buffer[1 + ptsOffset] = (byte) ((pts >> 22) & 0xFF);
		buffer[2 + ptsOffset] = (byte) ((((pts >> 15) & 0x7F) << 1) + 1);
		buffer[3 + ptsOffset] = (byte) ((pts >> 7) & 0xFF);
		buffer[4 + ptsOffset] = (byte) (((pts & 0x7F) << 1) + 1);

		return;
	}

	private int getContinuityCounter(byte[] buffer, int offset) {
		int counter = buffer[3 + offset] & 0x0f;
		return counter;
	}

	public void increaseContinuityCounter() {
		int counter = data[3 + offset] & 0x0f;
		counter++;
		counter = counter % 16;
		data[3 + offset] = (byte) (data[3 + offset] & 0xf0);
		data[3 + offset] = (byte) (data[3 + offset] | (byte) counter);
	}

	private int getPID(byte[] buffer, int offset) {
		int pid = (((buffer[1 + offset] & 0x1f) << 8) & 0x0000ffff)
				+ (buffer[2 + offset] & 0x00ff);
		return pid;
	}

	private int getPayloadOffset(byte[] buffer, int offset) {
		int internalOffset = tsDataOffset(buffer, offset);

		if (isStartOfPES(buffer, offset)) {
			byte pes_header_data_length = buffer[internalOffset + 8];
			internalOffset += 9 + pes_header_data_length;
		}

		return internalOffset;
	}

	private boolean isStart(byte[] tsPacket, int offset) {
		return (tsPacket[1 + offset] & 0x40) != 0;
	}

	private int getPMTLength(byte[] buffer, int offset) {

		int tsDataOffset = tsDataOffset(buffer, offset) + 2;
		int pmtSectionLength = ((buffer[tsDataOffset] & 0xf) << 8)
				+ buffer[1 + tsDataOffset];
		return pmtSectionLength + 3;
	}

	private int getProgramPID(byte[] buffer, int offset) {
		int dataOffset = tsDataOffset(buffer, offset);
		return ((buffer[dataOffset + 11] & 0x1F) << 8)
				+ (int)((int)0x0ff & (int)buffer[dataOffset + 12]);
	}

	private boolean payloadExists(byte[] buffer, int offset) {
		return (buffer[3 + offset] & 0x10) != 0;
	}

	private boolean validMpegTsPacket(byte[] buffer, int offset) {

		return buffer[offset] == 0x47;

	}

	private long readPCR(byte[] buffer, int offset) {

		long pcrBase = ((long) (buffer[offset + 6] & 0x0ff) << 25)
				| ((long) (buffer[offset + 7] & 0x0ff) << 17)
				| ((long) (buffer[offset + 8] & 0x0ff) << 9)
				| ((long) (buffer[offset + 9] & 0x0ff) << 1)
				| ((long) (buffer[offset + 10] & 0x0ff) >> 7);

		long pcrExt = ((long) (buffer[offset + 10] & 0x001) << 8)
				| ((long) (buffer[offset + 11] & 0x0ff));

		return pcrBase * 300 + pcrExt;
	}

	public VideoFrameType getMpeg2VideoFrameType() {

		byte[] buffer = data;

		int startOffset = tsDataOffset(buffer, offset) - offset;
		for (int i = startOffset; i < (Mpeg2TSPacket.TRANSPORT_PACKET_SIZE - 7); i++) {
			int j = offset + i;
			if (buffer[j] == 0) {
				if (buffer[j + 1] == 0) {
					if (buffer[j + 2] == 1) {
						if (buffer[j + 3] == 0) {
							// System.out.println(Integer.toBinaryString(buffer[j
							// + 5]));
							if (((buffer[j + 5] >>> 3) & 0x07) == 1) {
								return VideoFrameType.I_START;
							}
							if (((buffer[j + 5] >>> 3) & 0x07) == 2) {
								return VideoFrameType.P_START;
							}
						}
					}
				}
			}
		}
		return VideoFrameType.OTHER;
	}

	public VideoFrameType getH264VideoFrameType() {
		short type = getH264type();

		if (type == 0) {
			return VideoFrameType.I_START;
		} else if ((type == 1) || (type == 6)) {
			return VideoFrameType.P_START;
		}

		return VideoFrameType.OTHER;

	}

	public short getH264type() {

		byte[] buffer = data;

		int startOffset = tsDataOffset(buffer, offset) - offset;

		// If start of PES, look for Sequence Parameter Set (SPS)
		// If it was found then it is I frame / GOP start otherwise P frame
		if (isStartOfPES()) {
			for (int i = startOffset; i < (Mpeg2TSPacket.TRANSPORT_PACKET_SIZE - 5); i++) {
				int j = offset + i;
				if (buffer[j] == 0) {
					if (buffer[j + 1] == 0) {
						if (buffer[j + 2] == 1) {
							short val = (short) (0x001f & buffer[j + 3]);
							if (val == 7) {
								return 0; // Start GOP was detected
							}
						}
					}
				}
			}
			return 1; // P frame was detected
		}
		return -1;

		// Old test - Look for an access unit delimiter NAL Unit
		/*for (int i = startOffset; i < (Mpeg2TSPacket.TRANSPORT_PACKET_SIZE - 5); i++) {
			int j = offset + i;
			if (buffer[j] == 0) {
				if (buffer[j + 1] == 0) {
					if (buffer[j + 2] == 1) {
						short val = (short) (0x001f & buffer[j + 3]);
						System.out.println("NAL type = " + val);
						if (val == 9) {
							val = (short) (0x0007 & (buffer[j + 4] >>> 5));
							System.out.println("Picture type = " + val);
							System.out.println("buffer[j + 4] = "
									+ buffer[j + 4]);
							return val;
						}
					}
				}
			}
		}*/
	}

}

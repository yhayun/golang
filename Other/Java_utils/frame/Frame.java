package com.ibm.frame;

/**
 * Wraps a byte array (similar to ByteBuffer) with additional data such as PTS.
 */
public class Frame {

	private byte[] data;
	private int offset;
	private long pts;

	public Frame(int capacity) {
		data = new byte[capacity];
		offset = 0;
	}

	/**
	 * @return A reference to the underlying byte array.
	 */
	public byte[] getData() {
		return data;
	}

	/**
	 * Copies a given byte array into this frame.
	 * 
	 * @param data
	 *            The byte array to copy
	 */
	public void setData(byte[] data) {
		clear();
		append(data, 0, data.length);
	}

	/**
	 * @return The PTS of the frame.
	 */
	public long getPTS() {
		return pts;
	}

	/**
	 * @param pts
	 *            The PTS of the frame.
	 */
	public void setPTS(long pts) {
		this.pts = pts;
	}

	/**
	 * @return true if this frame is empty, false otherwise.
	 */
	public boolean isEmpty() {
		return offset == 0;
	}

	/**
	 * @return The size of the data wrapped by this frame (in bytes).
	 */
	public int size() {
		return offset;
	}

	public void setCurrentSize(int size) {
		offset = size;
	}
	
	/**
	 * @return The capacity of this frame.
	 */
	public int capacity() {
		return data.length;
	}

	/**
	 * Clears the content of this frame.
	 */
	public void clear() {
		offset = 0;
	}

	/**
	 * Appends a given byte array to this frame.
	 * 
	 * @param buf
	 *            The byte array to append
	 * @param offset
	 *            The offset in buf
	 * @param limit
	 *            The number of bytes to copy from buf
	 */
	public void append(byte[] buf, int offset, int limit) {
		for (int i = 0; i < limit && this.offset < data.length; ++i)
			data[this.offset++] = buf[offset + i];
	}

}

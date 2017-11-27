package main

type frame struct {
	pts long
	offset int
	data []byte
}

//CONSTRUCTOR:
func NewFrame(capacity int) *frame {
	return &frame{
		offset: 0,
		data: make([]byte,capacity),
	}
}


/**
 * @return A reference to the underlying byte array.
 */
 func (f frame) GetData() []byte {
 	return f.data
 }

/**
* Clears the content of this frame.
*/
func ( f frame) Clear() {
	f.offset = 0
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
 func (f frame) Append (buf []byte, offset int ,limit int ) {
	for i:= 0; i < limit && f.offset < len(f.data); i++ {
		f.data[f.offset] = buf[offset + i]
		f.offset++
	 }
 }


/**
 * Copies a given byte array into this frame.
 *
 * @param data
 *            The byte array to copy
 */
 func ( f frame) SetData(data []byte) {
 	f.Clear()
 	f.Append(data, 0, len(data))
 }


/**
 * @return The PTS of the frame.
 */
 func (f frame) GetPTS () long {
 	return f.pts
 }


/**
 * @param pts
 *            The PTS of the frame.
 */
 func (f frame) SetPTS(pts long) {
 	f.pts = pts
 }


/**
 * @return true if this frame is empty, false otherwise.
 */
 func (f frame) IsEmpty() bool {
 	return f.offset == 0
 }

/**
 * @return The size of the data wrapped by this frame (in bytes).
 */
func (f frame) Size() int {
	return f.offset
}

func ( f frame) SetCurrentSize( size int) {
	f.offset = size
}

/**
 * @return The capacity of this frame.
 */
 func (f frame) Capacity() int {
 	return len(f.data)
 }



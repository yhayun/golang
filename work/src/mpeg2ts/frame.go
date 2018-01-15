package main

type Frame struct {
	pts uint32
	dts uint32
	offset int
	data []byte
}

//CONSTRUCTOR:
func NewFrame(capacity int) *Frame {
	return &Frame{
		offset: 0,
		data: make([]byte,capacity),
	}
}


/**
 * @return A reference to the underlying byte array.
 */
 func (f *Frame) GetData() []byte {
 	return f.data
 }

/**
* Clears the content of this Frame.
*/
func ( f *Frame) Clear() {
	f.offset = 0
}

/**
 * Appends a given byte array to this Frame.
 *
 * @param buf
 *            The byte array to append
 * @param offset
 *            The offset in buf
 * @param limit
 *            The number of bytes to copy from buf
 */
 func (f *Frame) Append (buf []byte, offset int ,limit int ) {
	for i:= 0; i < limit && f.offset < len(f.data); i++ {
		f.data[f.offset] = buf[offset  + i]
		f.offset++
	 }
 }


/**
 * Copies a given byte array into this Frame.
 *
 * @param data
 *            The byte array to copy
 */
 func ( f *Frame) SetData(data []byte) {
 	f.Clear()
 	f.Append(data, 0, len(data))
 }


/**
 * @return The PTS of the Frame.
 */
 func (f *Frame) GetPTS () uint32 {
 	return f.pts
 }


/**
 * @param pts
 *            The PTS of the Frame.
 */
 func (f *Frame) SetPTS(pts uint32) {
 	f.pts = pts
 }

/**
* @return The PTS of the Frame.
*/
func (f *Frame) GetDTS () uint32 {
	return f.dts
}


/**
 * @param pts
 *            The PTS of the Frame.
 */
func (f *Frame) SetDTS(dts uint32) {
	f.dts = dts
}


/**
 * @return true if this Frame is empty, false otherwise.
 */
 func (f *Frame) IsEmpty() bool {
 	return f.offset == 0
 }

/**
 * @return The size of the data wrapped by this Frame (in bytes).
 */
func (f *Frame) Size() int {
	return f.offset
}

func ( f *Frame) SetCurrentSize( size int) {
	f.offset = size
}

/**
 * @return The capacity of this Frame.
 */
 func (f *Frame) Capacity() int {
 	return len(f.data)
 }



package main

import (
	//"fmt",
)

type frame struct {
	pts uint64
	offset int
	data []byte
}

//func (c Creature) Dump() {
//	fmt.Printf("Name: '%s', Real: %t\n", c.Name, c.Real)
//}



/**
 * @return A reference to the underlying byte array.
 */
 func (f frame) getData() []byte {
 	return f.data;
 }

/**
* Clears the content of this frame.
*/
func ( f frame) clear() {
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
 func (f frame) append (buf []byte, offset int ,limit int ) {
	for i:= 0; i < limit && f.offset < len(f.data); i++ {
		f.data[f.offset] = buf[offset + i];
		f.offset++
	 }
 }


/**
 * Copies a given byte array into this frame.
 *
 * @param data
 *            The byte array to copy
 */
 func ( f frame) setData(data []byte) {
 	f.clear()
 	f.append(data, 0, len(data))
 }


/**
 * @return The PTS of the frame.
 */
 func (f frame) getPTS () uint64 {
 	return f.pts
 }


/**
 * @param pts
 *            The PTS of the frame.
 */
 func (f frame) setPTS(pts uint64) {
 	f.pts = pts
 }


/**
 * @return true if this frame is empty, false otherwise.
 */
 func (f frame) isEmpty() bool {
 	return f.offset == 0
 }

/**
 * @return The size of the data wrapped by this frame (in bytes).
 */
func (f frame) size() int {
	return f.offset
}

func ( f frame) setCurrentSize( size int) {
	f.offset = size
}

/**
 * @return The capacity of this frame.
 */
 func (f frame) capacity() int {
 	return len(f.data)
 }



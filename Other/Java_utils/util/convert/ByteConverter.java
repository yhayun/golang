package com.ibm.util.convert;

import com.ibm.cms.Copyright;

public abstract class ByteConverter {

	@SuppressWarnings("unused")
	private static final transient String copyright = Copyright.COPYRIGHT;
	
	/**
	 * Converts a long to bytes and inserts into a given byte array
	 * @param num the number to convert to byte array
	 * @param array the byte array where the number must be inserted
	 * @param offset the index at which the number must be inserted in <code>array</code>
	 */
	public static void longToBytesLE(long num, byte[] array, int offset) {
		for (int i = 0; i <= 7; i++) {        
			array[offset + i] = (byte)(int)(num >>> (i * 8));
		}
	}

	/**
	 * Converts an integer to bytes and inserts into a given byte array
	 * @param num the number to convert to byte array
	 * @param array the byte array where the number must be inserted
	 * @param offset the index at which the number must be inserted in <code>array</code>
	 */
	public static void intToBytesLE(int num, byte[] array, int offset) {
		for (int i = 0; i <= 3; i++) {        
			array[offset + i] = (byte)(num >>> (i * 8));
		}
	}
	
	/**
	 * Converts a short to bytes and inserts into a given byte array
	 * @param num the number to convert to byte array
	 * @param array the byte array where the number must be inserted
	 * @param offset the index at which the number must be inserted in <code>array</code>
	 */
	public static void shortToBytesLE(short num, byte[] array, int offset) {
		for (int i = 0; i <= 1; i++) {        
			array[offset + i] = (byte)(num >>> (i * 8));
		}
	}
	
	
	/**
	 * Converts a long to byte array
	 * @param num the number to convert to byte array
	 * @return the byte representation of <code>num</code>
	 */
	public static byte[] longToBytesLE(long num) {
		
		byte[] array = new byte[8];
		
		for (int i = 0; i <= 7; i++) {        
			array[i] = (byte)(int)(num >>> (i * 8));
		}
		
		return array;
	}

	/**
	 * Converts an integer to byte array
	 * @param num the number to convert to byte array
	 * @return the byte representation of <code>num</code>
	 */
	public static byte[] intToBytesLE(int num) {

		byte[] array = new byte[4];
		
		for (int i = 0; i <= 3; i++) {        
			array[i] = (byte)(num >>> (i * 8));
		}
		
		return array;
	}
	
	/**
	 * Converts a short to byte array
	 * @param num the number to convert to byte array
	 * @return the byte representation of <code>num</code>
	 */
	public static byte[] shortToBytesLE(short num) {

		byte[] array = new byte[2];
		
		for (int i = 0; i <= 1; i++) {        
			array[i] = (byte)(num >>> (i * 8));
		}
		
		return array;
	}
	
	
	/**
	 * Converts a byte array to a long
	 * @param array the array that holds the byte representation of the long
	 * @param offset the index at which the byte representation begins 
	 * in <code>array</code>
	 * @return the long coded in <code>array</code>
	 */
	public static long bytesToLongLE(byte[] array, int offset) {
		long j = 0;
		for (int i = 0; i < 8 ; i++) {
			j |= ((long)array[offset + (7 - i)] & 0xff) << 8 * (7 - i);
		}
		return j;
	}
	
	/**
	 * Converts a byte array to an integer
	 * @param array the array that holds the byte representation of the integer
	 * @param offset the index at which the byte representation begins 
	 * in <code>array</code>
	 * @return the integer coded in <code>array</code>
	 */
	public static int bytesToIntLE(byte[] array, int offset) {
		int j = 0;
		for (int i = 0; i < 4 ; i++) {
			j |= (array[offset + (3 - i)] & 0xff) << 8 * (3 - i);
		}
		return j;
	}

	/**
	 * Converts a byte array to a short
	 * @param array the array that holds the byte representation of the short
	 * @param offset the index at which the byte representation begins 
	 * in <code>array</code>
	 * @return the short coded in <code>array</code>
	 */
	public static short bytesToShortLE(byte[] array, int offset) {
		short j = 0;
		for (int i = 0; i < 2 ; i++) {
			j |= (array[offset + (1 - i)] & 0xff) << 8 * (1 - i);
		}
		return j;
	}

	/**
	 * Converts a long to bytes and inserts into a given byte array
	 * @param num the number to convert to byte array
	 * @param array the byte array where the number must be inserted
	 * @param offset the index at which the number must be inserted in <code>array</code>
	 */
	public static void longToBytesBE(long num, byte[] array, int offset) {
		for (int i = 0; i <= 7; i++) {        
			array[offset + 7 - i] = (byte)(int)(num >>> (i * 8));
		}
	}

	/**
	 * Converts an integer to bytes and inserts into a given byte array
	 * @param num the number to convert to byte array
	 * @param array the byte array where the number must be inserted
	 * @param offset the index at which the number must be inserted in <code>array</code>
	 */
	public static void intToBytesBE(int num, byte[] array, int offset) {
		for (int i = 0; i <= 3; i++) {        
			array[offset + 3 - i] = (byte)(num >>> (i * 8));
		}
	}
	
	/**
	 * Converts a short to bytes and inserts into a given byte array
	 * @param num the number to convert to byte array
	 * @param array the byte array where the number must be inserted
	 * @param offset the index at which the number must be inserted in <code>array</code>
	 */
	public static void shortToBytesBE(short num, byte[] array, int offset) {
		for (int i = 0; i <= 1; i++) {        
			array[offset + 1 - i] = (byte)(num >>> (i * 8));
		}
	}
	
	
	/**
	 * Converts a long to byte array
	 * @param num the number to convert to byte array
	 * @return the byte representation of <code>num</code>
	 */
	public static byte[] longToBytesBE(long num) {
		
		byte[] array = new byte[8];
		
		for (int i = 0; i <= 7; i++) {        
			array[7 - i] = (byte)(int)(num >>> (i * 8));
		}
		
		return array;
	}

	/**
	 * Converts an integer to byte array
	 * @param num the number to convert to byte array
	 * @return the byte representation of <code>num</code>
	 */
	public static byte[] intToBytesBE(int num) {

		byte[] array = new byte[4];
		
		for (int i = 0; i <= 3; i++) {        
			array[3 - i] = (byte)(num >>> (i * 8));
		}
		
		return array;
	}
	
	/**
	 * Converts a short to byte array
	 * @param num the number to convert to byte array
	 * @return the byte representation of <code>num</code>
	 */
	public static byte[] shortToBytesBE(short num) {

		byte[] array = new byte[2];
		
		for (int i = 0; i <= 1; i++) {        
			array[1 - i] = (byte)(num >>> (i * 8));
		}
		
		return array;
	}
	
	
	/**
	 * Converts a byte array to a long
	 * @param array the array that holds the byte representation of the long
	 * @param offset the index at which the byte representation begins 
	 * in <code>array</code>
	 * @return the long coded in <code>array</code>
	 */
	public static long bytesToLongBE(byte[] array, int offset) {
		long j = 0;
		for (int i = 0; i < 8 ; i++) {
			j |= ((long)array[offset + i] & 0xff) << 8 * (7 - i);
		}
		return j;
	}
	
	/**
	 * Converts a byte array to an integer
	 * @param array the array that holds the byte representation of the integer
	 * @param offset the index at which the byte representation begins 
	 * in <code>array</code>
	 * @return the integer coded in <code>array</code>
	 */
	public static int bytesToIntBE(byte[] array, int offset) {
		int j = 0;
		for (int i = 0; i < 4 ; i++) {
			j |= (array[offset + i] & 0xff) << 8 * (3 - i);
		}
		return j;
	}

	/**
	 * Converts a byte array to a short
	 * @param array the array that holds the byte representation of the short
	 * @param offset the index at which the byte representation begins 
	 * in <code>array</code>
	 * @return the short coded in <code>array</code>
	 */
	public static short bytesToShortBE(byte[] array, int offset) {
		short j = 0;
		for (int i = 0; i < 2 ; i++) {
			j |= (array[offset + i] & 0xff) << 8 * (1 - i);
		}
		return j;
	}

	public static short unsignedByteToShort(byte value) {
		return (short) (0xff & value);
	}

	public static int unsignedShortToInt(short value) {
		return 0xffff & value;
	}

	public static long unsignedIntToLong(int value) {
		return (0xffffffffL & value);
	}
	
	public static String bytesToHexString(byte[] array, int offset, int length) {
		   StringBuilder sb = new StringBuilder();
		   for(int i=offset;i<offset + length;i++)
		      sb.append(String.format("%02x", array[i]&0xff));
		   return sb.toString();
		}

	public static byte[] hexStringToBytes(String s) {
		
		
		
	    int len = s.length();
	    byte[] data = new byte[len / 2];
	    for (int i = 0; i < len; i += 2) {
	        data[i / 2] = (byte) ((Character.digit(s.charAt(i), 16) << 4)
	                             + Character.digit(s.charAt(i+1), 16));
	    }
	    return data;
	}
	
	public static byte hexStringToByte(String s){
		s = preprocess(s, 1);
		return hexStringToBytes(s)[0];
	}
	
	public static short hexStringToShortBE(String s){
		s = preprocess(s, 2);
		return bytesToShortBE(hexStringToBytes(s),0);
	}
	
	public static int hexStringToIntBE(String s){
		s = preprocess(s, 4);
		return bytesToIntBE(hexStringToBytes(s),0);
	}
	
	public static long hexStringToLongBE(String s){
		s = preprocess(s, 8);
		return bytesToLongBE(hexStringToBytes(s),0);
	}
	
	private static String preprocess(String s, int size){
		if(s.startsWith("0x") || s.startsWith("0X")){
			s = s.substring(2);
		}
		
		while(s.length()/2 < size){
			s = "0"+s;
		}
		
		return s;
	}
	
}

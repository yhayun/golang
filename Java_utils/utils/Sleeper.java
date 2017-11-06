package com.ibm.utils;

import com.ibm.cms.Copyright;

import java.util.Random;

public class Sleeper {

	@SuppressWarnings("unused")
	private static final transient String copyright = Copyright.COPYRIGHT;

	
	/**
	 * Performs sleep according to provided number of 'ms'
	 * @param long timeToSleep - time to sleep
	 */
	public static void sleep(long timeToSleep){
	
		try {// performing random sleep to minimize chance for such collision in future
            Thread.sleep(timeToSleep);
        } catch (InterruptedException e) {}
	}
	
	public static void randomSleep(int rangeMS){
		Random random = new Random();
		long timeToSleep = random.nextInt(rangeMS);
		
		try {// performing random sleep to minimize chance for such collision in future
            Thread.sleep(timeToSleep);
        } catch (InterruptedException e) {}
	}
	
	public static void randomSleep(int min, int max){
	
		long timeToSleep = min + (int)(Math.random() * ((max - min) + 1));
		sleep(timeToSleep);

	}
	
	
	/**
	 * Performs wait on the provided object
	 * @param Object obj - the object that should wait
	 */
	public static void wait(Object obj){
		try {
			obj.wait();
        } catch (InterruptedException e) {}
	}
}

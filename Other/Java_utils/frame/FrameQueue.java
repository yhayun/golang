package com.ibm.frame;

import java.util.concurrent.BlockingQueue;
import java.util.concurrent.LinkedBlockingQueue;

public class FrameQueue extends LinkedBlockingQueue<Frame> implements
		PooledQueue<Frame> {

	private static final long serialVersionUID = 7519377596896832697L;

	private BlockingQueue<Frame> pool;

	public FrameQueue(int capacity, int frameCapacity) {
		super(capacity);
		pool = new LinkedBlockingQueue<Frame>(capacity + 2);
		while (pool.offer(new Frame(frameCapacity)))
			;
	}

	@Override
	public Frame newElement() {
		return pool.poll();
	}

	@Override
	public void recycle(Frame f) {
		if (f != null) {
			f.clear();
			pool.offer(f);
		}
	}

	@Override
	public void recycleAll() {
		while (!isEmpty())
			pool.offer(poll());
	}

}

package com.ibm.frame;

import java.util.Queue;

/**
 * Implementors should maintain a pool of pre-allocated instances of type E, and
 * allow users to ask for a new pre-allocated instance and recycle a used
 * instance.
 */
public interface PooledQueue<E> extends Queue<E> {

	/**
	 * @return A new pre-allocated instance of type E.
	 */
	public E newElement();

	/**
	 * @param e
	 *            Recycles a given instance of type E back to the pool.
	 */
	public void recycle(E e);

	/**
	 * Recycles all instances that are currently in the container (relevant only
	 * to containers that implement this interface).
	 */
	public void recycleAll();

}

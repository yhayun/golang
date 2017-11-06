package com.ibm.mpeg2ts;

import java.io.IOException;
import java.net.DatagramPacket;
import java.net.DatagramSocket;
import java.net.InetAddress;
import java.net.InetSocketAddress;
import java.net.MulticastSocket;
import java.net.SocketAddress;
import java.net.SocketTimeoutException;
import java.util.HashMap;
import java.util.Map;

import org.apache.log4j.Logger;

import com.ibm.frame.Frame;
import com.ibm.frame.PooledQueue;
import com.ibm.mpeg2ts.Mpeg2TSPacket.Mpeg2TSPacketType;
import com.ibm.mpeg2ts.StreamInfo.TS_STREAM_TYPE;

public class Mpeg2TSSource {

	private static Logger log = Logger.getLogger(Mpeg2TSSource.class);

	private static final int MAX_UDP_PACKET_SIZE = 1500;

	private DatagramSocket socket;

	private boolean reTransmitFlag = true;

	private PooledQueue<Frame> videoFrames;

	private int programPID = 0;

	private UdpSource udpSource = null;

	private int outputPort = -1;

	private int socketReceiveBufferSize;

	public Mpeg2TSSource(String ip, int port, PooledQueue<Frame> videoFrames,
			int socketReceiveBufferSize) throws IOException {
		this.socketReceiveBufferSize = socketReceiveBufferSize;
		this.videoFrames = videoFrames;
		socket = createMulticastSocket(new InetSocketAddress(ip, port));
	}

	public Mpeg2TSSource(int port, PooledQueue<Frame> videoFrames,
			int socketReceiveBufferSize, boolean reTransmitFlag)
			throws IOException {
		this.socketReceiveBufferSize = socketReceiveBufferSize;
		outputPort = port + 1;
		this.videoFrames = videoFrames;
		this.reTransmitFlag = reTransmitFlag;
		socket = createUnicastSocket(port);
	}

	public void start() {
		udpSource = new UdpSource(socket, videoFrames, outputPort);
		Thread thread = new Thread(udpSource);
		thread.setName(udpSource.getClass().getSimpleName());
		thread.start();
	}

	public void close() {
		if (udpSource != null) {
			udpSource.close();
		}
	}

	public DatagramSocket createUnicastSocket(int port) throws IOException {

		DatagramSocket socket = null;
		socket = new DatagramSocket(port);
		socket.setReceiveBufferSize(socketReceiveBufferSize);
		socket.setSoTimeout(1000);
		return socket;
	}

	public DatagramSocket createMulticastSocket(SocketAddress address)
			throws IOException {

		DatagramSocket socket = null;

		InetSocketAddress socketAddress = (InetSocketAddress) address;

		if (socketAddress.getAddress().isMulticastAddress()) {
			String osName = System.getProperty("os.name").toLowerCase();
			if (osName.indexOf("linux") != -1) {
				socket = new MulticastSocket(address);
			} else if (osName.indexOf("windows") != -1) {
				socket = new MulticastSocket(socketAddress.getPort());
				socket.setReuseAddress(true);//redundant, already set with empty constructor
			} else {
				log.warn("Subserver is running on unknown O.S. (supporting windows and linux)");
				System.exit(1);
			}

			socket.setReceiveBufferSize(socketReceiveBufferSize);

			((MulticastSocket) socket).joinGroup(socketAddress.getAddress());

		} else {
			socket = new DatagramSocket(address);
		}
		socket.setSoTimeout(1000);
		return socket;
	}

	public class UdpSource implements Runnable {

		private DatagramSocket socket = null;
		private DatagramSocket outputSocket = null;

		private boolean endFlag = false;
		private boolean detectFlag = false;
		private int videoPID;

		private Mpeg2TSParser videoTSParser;
		private Mpeg2TSPacket tsPacket = new Mpeg2TSPacket();

		private PMTFrame pmtFrame = new PMTFrame();
		private Map<Integer, StreamInfo> streamsMap = new HashMap<Integer, StreamInfo>();
		
		private long previousTime = -1;

		DatagramPacket packet = new DatagramPacket(
				new byte[MAX_UDP_PACKET_SIZE], MAX_UDP_PACKET_SIZE);
		DatagramPacket outputPacket = new DatagramPacket(
				new byte[MAX_UDP_PACKET_SIZE], MAX_UDP_PACKET_SIZE);
		private int outputPort;

		private static final int MPEG2TS_PACKET_LENGTH = 188;

		public UdpSource(DatagramSocket socket, PooledQueue<Frame> videoFrames,
				int outputPort) {
			this.socket = socket;
			this.outputPort = outputPort;
			videoTSParser = new Mpeg2TSParser(videoFrames);
		}

		/**
		 * Extracts Mpeg2TSPackets from a given DatagramPacket and sends them to
		 * the corresponding Mpeg2TSParser according to their PID.
		 * 
		 * @param packet
		 *            The packet to extract Mpeg2TSPackets from.
		 */
		private void extractStreams(DatagramPacket packet) {

			for (int offset = 0; offset < packet.getLength(); offset += MPEG2TS_PACKET_LENGTH) {

				tsPacket.fromBytes(packet.getData(), offset, programPID);
				
				if(previousTime != -1) {
					if(detectFlag && ((System.currentTimeMillis() - previousTime) > 1000)) {
						log.info("Timeout detected, now looking for new PID");
						programPID = 0;
						pmtFrame = new PMTFrame();
						detectFlag = false;
					}
				}
				
				if (detectFlag) {				
					if (tsPacket.getPID() == videoPID) {
						previousTime = System.currentTimeMillis();
						videoTSParser.write(tsPacket);
					}
				} else {
					if (tsPacket.getType() == Mpeg2TSPacketType.PAT) {
						programPID = tsPacket.getProgramPID();
					}

					if (tsPacket.getType() == Mpeg2TSPacketType.PMT) { // PMT
																		// received

						if (pmtFrame.addPacket(tsPacket)) {
							// pmtFrame is complete
							streamsMap = pmtFrame.getStreamInfos();

							for (Map.Entry<Integer, StreamInfo> entry : streamsMap
									.entrySet()) {
								StreamInfo info = (StreamInfo) entry.getValue();
								if (info.streamType == TS_STREAM_TYPE.TS_STREAM_TYPE_H264) {
									log.info("New PID detected = " + info.streamPID);
									previousTime = System.currentTimeMillis();
									detectFlag = true;
									videoPID = info.streamPID;
									break;
								}

							}
						}
					}

				}
			}
		}

		@Override
		public void run() {

			boolean firstPacket = true;

			byte[] buffer = new byte[1500];
			packet.setData(buffer);

			if (reTransmitFlag) {
				try {
					outputSocket = new DatagramSocket();
					outputSocket.setSendBufferSize(socketReceiveBufferSize);
					outputPacket.setAddress(InetAddress.getLocalHost());
					outputPacket.setPort(outputPort);
				} catch (Exception e1) {
					e1.printStackTrace();
					log.info("Error creating output socket");
					return;
				}
			}

			while (!endFlag) {

				try {
					//log.info("Before socket.receive ");
					socket.receive(packet);
					//log.info("After socket.receive ");
					if (reTransmitFlag) {
						outputPacket.setData(packet.getData());
						outputPacket.setLength(packet.getLength());
						outputSocket.send(outputPacket);
					}

					if (firstPacket) {
						log.info("At least one packet was received");
						firstPacket = false;
					}

					extractStreams(packet);

				} catch (SocketTimeoutException e) {
					log.info("No packet received");
				} catch (Exception e) {
					log.error(e.toString());
					e.printStackTrace();
				}

			}

			log.info("Before close socket ");
			try {
				cleanSocket();
			} catch (IOException e) {
				log.info("Exception in closing socket - " + e.getMessage());
			}
			log.info("After close socket");

		}

		private void cleanSocket() throws IOException {

			if (outputSocket != null && !outputSocket.isClosed()) {
				outputSocket.close();
				outputSocket = null;
			}
			if (socket != null && !socket.isClosed()) {
				socket.close();
				socket = null;
			}
		}

		public void close() {
			videoTSParser.close();
			endFlag = true;
		}

	}

}

rem asset_dir targetVideoIP targetVideoPort rate seekTime
java -cp subserver.jar;log4j-1.2.14.jar com.ibm.cms.rtmf.test.Streamer .\Haivision_H264p  127.0.0.1 8888 1 0
pause;
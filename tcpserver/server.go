package tcpserver

import (
	"fmt"
	"net"
	"bufio"
	"syscall"
	"time"
	. "tcpserver/globa"
)

const maxRead = 1024 * 100

func ServerMain() {
	listener := initServer(":50000")
	for {
		conn, err := listener.Accept()
		CheckError(err, "Accept: ")
		go connectionHandler(conn)
	}
}

func initServer(hostAndPort string) net.Listener {
	serverAddr, err := net.ResolveTCPAddr("tcp", hostAndPort)
	CheckError(err, "Resolving address:port failed: '"+hostAndPort+"'")
	listener, err := net.ListenTCP("tcp", serverAddr)
	CheckError(err, "ListenTCP: ")
	println("Listening to: ", listener.Addr().String())
	return listener
}

func connectionHandler(conn net.Conn) {
	var ibuf []byte = make([]byte, maxRead+1)
	connFrom := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	println("Connection from: ", connFrom)

	for {
		length, err := reader.Read(ibuf[0:maxRead])
		fmt.Printf("length1111111=%d\n", length)
		switch err {
		case nil:
			if length < 12{//没收到足够的数据
				continue
			}
			dataLength := ByteToUint16(ibuf[6:8])
			desiredLength := dataLength + 10
			for {
				if length < desiredLength{
					partLength, err := reader.Read(ibuf[length:maxRead])
					length += partLength
					// fmt.Printf("partLength=%d,length=%d\n", partLength,length)
					if err != nil{ //若收不到数据就退出
						break;
					}
				}else{ //收满了就退出
					if ibuf[0] == 0x7E {
						finalBuf := ibuf[:desiredLength]
						go handleMsg(desiredLength, finalBuf,writer)
					}
					break;
				}
			}
			// sayHello(conn)
		case syscall.EAGAIN: // try again
			continue
		default:
			goto DISCONNECT
		}
	}
DISCONNECT:
	err := conn.Close()
	println("Closed connection: ", connFrom)
	CheckError(err, "Close: ")
}

func sayHello(to net.Conn) {
	obuf := []byte{'L', 'e', 't', '\'', 's', ' ', 'G', 'O', '!', '\n'}
	wrote, err := to.Write(obuf)
	CheckError(err, "Write: wrote "+string(wrote)+" bytes.")
}

func handleMsg(length int, msg []byte, writer *bufio.Writer) {
	if length > 0 {
		print("<", length, ":")
		for i := 0; i < length; i++ {
			fmt.Printf("%02x ", msg[i])
		}
		print(">\n")
	}

	var sqlToDo string

	if length > 10 {
		// zoneid := int(msg[1])
		devtype := int(msg[2])
		devindex := ByteToUint16(msg[3:5])
		// println(devindex)
		// fmt.Printf("%f\n", data)
		switch (devtype){
		case TYPE_TEMP_SENSOR:
			// 7E0101000101000441BC7AE11234
			data := ByteToFloat32(msg[8:12])
			rows, _ := Db.Query("SELECT sensor_id FROM d_temp_data WHERE sensor_id = ?", devindex)

			defer rows.Close()
			if rows != nil{
				sqlToDo = fmt.Sprintf("UPDATE d_temp_data set temp=%f where sensor_id=%d",data, devindex)
			}else{
				sqlToDo = fmt.Sprintf("INSERT d_temp_data(sensor_id,temp) VALUES(%d,%f)",devindex, data)
			}
			break
		case TYPE_VIBRATION_SENSOR:
			// 7E01020001010004000000011234
			data := ByteToUint32(msg[8:12])
			rows, _ := Db.Query("SELECT sensor_id FROM d_vibration_data WHERE sensor_id = ?", devindex)
			defer rows.Close()
			if rows != nil{
				sqlToDo = fmt.Sprintf("UPDATE d_vibration_data set status=%d where sensor_id=%d",data, devindex)
			}else{
				sqlToDo = fmt.Sprintf("INSERT d_vibration_data(sensor_id,status) VALUES(%d,%d)",devindex, data)
			}
			break		
		}
		println(sqlToDo)
		for{
			_,err := Db.Exec(sqlToDo)
			if err != nil{
		    	CheckError(err, "sql failed:")
		    	time.Sleep(100000)
			}else{
				// break
			}
		}
		writer.WriteString("ok\n")
		writer.Flush()
		println("write finish")
	}
}
// var count int = 0
// func handleMsg(length int, msg []byte, writer *bufio.Writer) {
// 	fmt.Printf("lengthfffffff=%d\n", length)
// 	if length > 12 {
// 		dataLength := ByteToUint16(msg[6:8])
// 		devindex := ByteToUint16(msg[3:5])

// 	    sqlToDo := fmt.Sprintf("INSERT d_realtime_sensor_%d (data) VALUES(%d)",devindex, ByteToUint32(msg[8:8+4]));
// 	    for i:=12;(i <= dataLength + 4) && (i <= length - 4);i=i+4 {
// 	        sqlToDo += fmt.Sprintf(",(%d)",ByteToUint32(msg[i:i+4]))
// 	    }

// 		// fmt.Printf("%d:%s\n", dataLength,sqlToDo)
// 	    for{
// 			_,err := Db.Exec(sqlToDo)
// 			if err != nil{
// 		    	CheckError(err, "sql failed:")
// 		    	time.Sleep(500000)
// 			}else{
// 				break;
// 			}
// 		}

// 		writer.WriteString("ok\n")
// 		writer.Flush()
// 		count++
// 		// println(count)
// 	}
// }
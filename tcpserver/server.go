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
			desiredLength := dataLength + 8
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

// func handleMsg(length int, msg []byte) {
// 	if length > 0 {
// 		print("<", length, ":")
// 		for i := 0; i < length; i++ {
// 			fmt.Printf("%02x ", msg[i])
// 		}
// 		print(">\n")
// 	}
	
// 	if length > 12 {
// 		// zoneid := int(msg[1])
// 		// devtype := int(msg[2])
// 		devindex := ByteToUint16(msg[3:5])
// 		// println(devindex)
// 		data := ByteToFloat32(msg[7:11])
// 		// fmt.Printf("%f\n", data)
// 		// 7E01010001010441BC7AE11234
// 		_,err := Db.Exec("UPDATE d_temp_data set temp=? where sensor_id=?",data,devindex)
// 	    CheckError(err, "sql failed:")
// 	}
// }
var count int = 0
func handleMsg(length int, msg []byte, writer *bufio.Writer) {
	fmt.Printf("lengthfffffff=%d\n", length)
	if length > 12 {
		dataLength := ByteToUint16(msg[6:8])
		devindex := ByteToUint16(msg[3:5])

	    sqlToDo := fmt.Sprintf("INSERT d_realtime_sensor_%d (data) VALUES(%d)",devindex, ByteToUint32(msg[8:8+4]));
	    for i:=12;(i <= dataLength + 4) && (i <= length - 4);i=i+4 {
	        sqlToDo += fmt.Sprintf(",(%d)",ByteToUint32(msg[i:i+4]))
	    }

		// fmt.Printf("%d:%s\n", dataLength,sqlToDo)
	    for{
			_,err := Db.Exec(sqlToDo)
			if err != nil{
		    	CheckError(err, "sql failed:")
		    	time.Sleep(500000)
			}else{
				break;
			}
		}

		writer.WriteString("ok\n")
		writer.Flush()
		count++
		// println(count)
	}
}
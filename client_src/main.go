package main

import (	
	"fmt"
	"net"
    "time"
    "bufio"
	"math/rand"
	"encoding/binary"
)

const maxRead = 1024 * 100

func main() {
	var sendData []byte
	var u16Tmp []byte = make([]byte, 2)
	var u32Tmp []byte = make([]byte, 4)

	sendData = append(sendData, 0x7E, 1, 1)

	binary.BigEndian.PutUint16(u16Tmp, 1)
	sendData = append(sendData, u16Tmp...)

	sendData = append(sendData, 1)

	binary.BigEndian.PutUint16(u16Tmp, 40000)
	sendData = append(sendData, u16Tmp...)

    for i:=0;i<10000;i++ {
    	randNum := rand.New(rand.NewSource(time.Now().UnixNano()))
		binary.BigEndian.PutUint32(u32Tmp, uint32(randNum.Intn(1000)))
		sendData = append(sendData, u32Tmp...)
		// print(binary.BigEndian.Uint32(sendData[i:i+4]))
		// print(",")
    }

    for {
		conn, _ := net.Dial("tcp", "39.108.4.211:50000")
		if conn != nil {			
			writer := bufio.NewWriter(conn)
			reader := bufio.NewReader(conn)

			defer conn.Close()
			for {
				writer.Write(sendData)
				writer.Flush()
				println("write finish")
				recv, _, err := reader.ReadLine()
				println("read finish")
				if err != nil{
					fmt.Printf("err=%v\n", err)
					conn.Close()
					// break
				}else{
					fmt.Printf("recv=%q\n", recv)
				}
			}
		}
		time.Sleep(100000)
    }
	// conn.Close()
}
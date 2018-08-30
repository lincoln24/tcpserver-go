package global

import (
	_ "fmt"
	_ "log"
	"encoding/binary"
    "math"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

const (
	SUCCESS = 0
	FAILURE = -1
)

const (
    TYPE_TEMP_SENSOR = 1
    TYPE_VIBRATION_SENSOR = 2
)

var Db *sql.DB

func CheckError(error error, info string) {
	if error != nil {
		println("ERROR: " + info + " " + error.Error()) // terminate program
	}
}

func Float32ToByte(src float32) []byte {
    bits := math.Float32bits(src)
    bytes := make([]byte, 4)
    binary.BigEndian.PutUint32(bytes, bits)
 
    return bytes
}
 
func ByteToFloat32(bytes []byte) float32 {
    bits := binary.BigEndian.Uint32(bytes)
 	// fmt.Printf("%02x %02x %02x %02x\n", bytes[0],bytes[1],bytes[2],bytes[3])
    return math.Float32frombits(bits)
}
 
func ByteToUint32(bytes []byte) int {
 	// fmt.Printf("%02x %02x %02x %02x\n", bytes[0],bytes[1],bytes[2],bytes[3])
    return int(binary.BigEndian.Uint32(bytes))
}
 
func ByteToUint16(bytes []byte) int { 
	return int(bytes[1]) | int(bytes[0])<<8
}
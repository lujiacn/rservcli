package rservcli

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/lujiacn/rservcli/assign"
	"github.com/lujiacn/rservcli/constants"
	"io"
	"net"
	"strconv"
	"strings"
)

var _ = fmt.Print

//authType
type authType int

const (
	atPlain authType = 1
	atCrypt authType = 2
)

//dataType
type dataType int

const (
	dtString dataType = 4
	dtSexp   dataType = 10
)

//command Type
type command int

const (
	cmdLogin    command = 1
	cmdVoidEval command = 2
	cmdEval     command = 3
)

type Rcli struct {
	// hostName            string
	// portNumber          int64
	conn            *net.TCPConn
	ReadWriteCloser io.ReadWriteCloser //connect interface
	ReadWriter      *bufio.ReadWriter  //buffer interface, wrap connect
}

//Input host and port
func NewRcli(hostname string, portnumber int64) (*Rcli, error) {
	var addr *net.TCPAddr
	var conn *net.TCPConn

	addr, err := net.ResolveTCPAddr("tcp", hostname+":"+strconv.FormatInt(portnumber, 10))
	if err != nil {
		return nil, err
	}

	conn, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}

	r := new(Rcli)
	r.conn = conn
	r.ReadWriteCloser = r.conn
	buffReader := bufio.NewReader(r.ReadWriteCloser)
	buffWriter := bufio.NewWriter(r.ReadWriteCloser)
	r.ReadWriter = bufio.NewReadWriter(buffReader, buffWriter)

	//handshake
	if err := r.parseInitMsg(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Rcli) Close() error {
	err := r.ReadWriteCloser.Close()
	return err
}

func (r *Rcli) readNBytes(bytes int) []byte {
	ret := make([]byte, bytes)
	for v := 0; v < bytes; v++ {
		ret[v], _ = r.ReadWriter.ReadByte()
	}
	return ret
}
func (r *Rcli) parseInitMsg() error {
	//32 bytes from initial connection
	rServeIDSig := string(r.readNBytes(4))
	rServeProtocol := string(r.readNBytes(4))
	rServeCommProtocol := string(r.readNBytes(4))
	r.readNBytes(20)
	fmt.Println(rServeCommProtocol, rServeIDSig, rServeProtocol)

	if rServeCommProtocol != "QAP1" ||
		rServeIDSig != "Rsrv" ||
		rServeProtocol != "0103" {
		return errors.New("Handshake failed")
	}
	return nil
}

func setHdr(valueType dataType, valueLength int, buf []byte) {
	buf[0] = byte(valueType)
	buf[1] = byte(valueLength & 255)
	buf[2] = byte((valueLength & 0xff00) >> 8)
	buf[3] = byte((valueLength & 0xff0000) >> 16)
}

func prepareStringCommand(cmd string) []byte {
	cmd = strings.Replace(cmd, "\r", "\n", -1) //avoid potential issue when loading external r script block
	rawCmdBytes := []byte(cmd)
	requiredLength := len(rawCmdBytes) + 1
	//make sure length is divisible by 4
	if requiredLength&3 > 0 {
		requiredLength = (requiredLength & 0xfffffc) + 4
	}
	cmdBytes := make([]byte, requiredLength+5)
	for i := 0; i < len(rawCmdBytes); i++ {
		cmdBytes[4+i] = rawCmdBytes[i]
	}
	setHdr(dtString, requiredLength, cmdBytes)
	return cmdBytes
}

func (r *Rcli) sendCommand(cmdType constants.Command, cmd string) {
	cmdBytes := prepareStringCommand(cmd)
	buf := new(bytes.Buffer)
	//command
	binary.Write(buf, binary.LittleEndian, int32(cmdType))
	//length of message (bits 0-31)
	binary.Write(buf, binary.LittleEndian, int32(len(cmdBytes)))
	//offset of message part
	binary.Write(buf, binary.LittleEndian, int32(0))
	// length of message (bits 32-63)
	binary.Write(buf, binary.LittleEndian, int32(0))
	binary.Write(buf, binary.LittleEndian, cmdBytes)

	r.ReadWriter.Write(buf.Bytes())
	r.ReadWriter.Flush()
}

func (r *Rcli) readResponse() Packet {
	rep := binary.LittleEndian.Uint32(r.readNBytes(4))
	r1 := binary.LittleEndian.Uint32(r.readNBytes(4))
	r.readNBytes(8)

	// fmt.Println(int(rep))

	if r1 <= 0 {
		return newPacket(int(rep), nil)
	}

	results := r.readNBytes(int(r1))
	return newPacket(int(rep), results)
}

func (r *Rcli) Eval(command string) (interface{}, error) {
	if r.conn == nil {
		return nil, errors.New("Connection was previously closed")
	}
	r.sendCommand(constants.CmdEval, command+"\n")
	p := r.readResponse()
	return p.GetResultObject()
}

func (r *Rcli) VoidEval(command string) error {
	if r.conn == nil {
		return errors.New("Connection was previously closed")
	}
	r.sendCommand(constants.CmdVoidEval, command+"\n")
	p := r.readResponse()
	_, err := p.GetResultObject()
	return err

}

func (r *Rcli) request(cmdType constants.Command, cont []byte, offset int, length int) Packet {
	if cont != nil {
		if offset >= len(cont) {
			cont = nil
			length = 0
		} else if length > (len(cont) - offset) {
			length = len(cont) - offset
		}
	}
	if offset < 0 {
		offset = 0
	}
	if length < 0 {
		length = 0
	}

	contlen := 0
	if cont != nil {
		contlen = length
	}

	hdr := make([]byte, 16)
	assign.SetInt(int(cmdType), hdr, 0)
	assign.SetInt(contlen, hdr, 4)
	for i := 8; i < 16; i++ {
		hdr[i] = 0
	}

	buf := new(bytes.Buffer)

	binary.Write(buf, binary.LittleEndian, hdr)
	binary.Write(buf, binary.LittleEndian, cont)

	r.ReadWriter.Write(buf.Bytes())
	r.ReadWriter.Flush()

	return r.readResponse()
}

func (r *Rcli) Assign(symbol string, value interface{}) error {
	fmt.Println("in assign")
	assignCommand, err := assign.Assign(symbol, value)
	if err != nil {
		fmt.Println("Error when assign")
		return err
	}
	rp := r.request(constants.CmdSetSexp, assignCommand, 0, len(assignCommand))
	fmt.Println("rp is", rp)
	if rp != nil && !rp.IsError() {
		return nil
	}
	return errors.New("Assign failed")
}

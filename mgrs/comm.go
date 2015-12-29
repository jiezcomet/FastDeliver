package mgrs

import (
	dm "FastDeliver/datamodel"
	. "FastDeliver/log"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"strconv"
	"github.com/golang/protobuf/proto"
	"os"
	"path/filepath"
	"io"
	"code.google.com/p/go-uuid/uuid"
)
const (
	buffer_length int =1024*4
	max_threads int=4
)
const (
	c_SendFile uint = iota + 1
	c_SendCmd
	c_ConnectReqEp
	c_ConnectPullEp
	c_Subscriber
	c_OpenReqEp
	c_OpenPubEp
)
type TransFileCmd struct{
	RootPath string
	Header *dm.FilePackageHeader
}
type cmdEvp struct {
	cmdType uint
	ep      *EndPoint
	data    interface{}
}

type EndPoint struct {
	Ip   string
	Port int
}

func (e EndPoint) ToTcp() string {
	return "tcp://" + e.Ip + ":" + strconv.Itoa(e.Port)
}

type Communicator struct {
	inputChan  chan *cmdEvp
	outputChan chan []byte
	epChannels map[string]chan *cmdEvp
}

func (c *Communicator) ConnectEp(ep *EndPoint) {
	Log.Debug("Start to connect req/rep endpoint %s", ep.ToTcp())
	var evp = &cmdEvp{c_ConnectReqEp, ep, nil}
	c.inputChan <- evp
}

func (c *Communicator) Subscriber(ep *EndPoint, myId string) {
	Log.Debug("Start to Subscriber endpoint %s with identifier [%s]", ep.ToTcp(), myId)
	var evp = &cmdEvp{c_Subscriber, ep, myId}
	c.inputChan <- evp
}
func (c *Communicator) ConnectPullEp(ep *EndPoint) {
	Log.Debug("Start to connect pull endpoint %s", ep.ToTcp())
	var evp = &cmdEvp{c_ConnectPullEp, ep, nil}
	c.inputChan <- evp
}
func (c *Communicator) OpenReqPort(ep *EndPoint){
	Log.Debug("Start to open Req/Reply endpoint %s", ep.ToTcp())
	var evp = &cmdEvp{c_OpenReqEp, ep, nil}
	c.inputChan <- evp
}
func (c *Communicator) OpenPubPort(ep *EndPoint){
	Log.Debug("Start to open pub endpoint %s", ep.ToTcp())
	var evp = &cmdEvp{c_OpenPubEp, ep, nil}
	c.inputChan <- evp
}


func (c *Communicator) SendFile(ep *EndPoint, file *TransFileCmd ) {
	Log.Debug("Send file fileName:%s,rel path:%s", file.Header.FileName, file.Header.RelPath)
	var evp = &cmdEvp{c_SendFile, ep, file}
	c.inputChan <- evp
}
func (c *Communicator) SendCmd(ep *EndPoint, cmdData []byte) {
	Log.Debug("Send cmd[%v] to endpoint [%s]", cmdData, ep.ToTcp())
	var evp = &cmdEvp{c_SendCmd, ep, cmdData}
	c.inputChan <- evp

}
func (c *Communicator) eventLoop() {
	for evp, isOk := <-c.inputChan; isOk; {
		switch evp.cmdType {
		case c_SendCmd:
			fallthrough
		case c_SendFile:
			if inputChan, isExists := c.epChannels[evp.ep.ToTcp()]; isExists {
				inputChan <- evp
			} else {
				Log.Error("Target endpoint [%s] has not been conntected or configured", evp.ep.ToTcp())
			}
		case c_ConnectReqEp:
			if err := c.connectReqEp(evp.ep); err != nil {
				Log.Error(err)
			}
		case c_ConnectPullEp:
			if ep, isOk := evp.data.(*EndPoint); isOk {
				Log.Info(ep)
			}
		case c_Subscriber:
			if id, isOk := evp.data.(string); isOk {
				Log.Info(id)
			}
		}
	}
}

func (c *Communicator) connectReqEp(ep *EndPoint) error {
	if _, isExists := c.epChannels[ep.ToTcp()]; isExists {
		return fmt.Errorf("The remote endpoint[%s] has been connected", ep.ToTcp())
	}
	var inputChan = make(chan *cmdEvp, 10000)
	c.epChannels[ep.ToTcp()] = inputChan
	var socket, err = zmq.NewSocket(zmq.REQ)
	defer socket.Close()
	socket.Connect(ep.ToTcp())
	go func(inputChan <-chan *cmdEvp) {
		for e, isOk := <-inputChan; isOk; {
			switch e.cmdType {
			case c_SendCmd:
				c.sendCmdHandler(socket,e)
			case c_SendFile:
					
			}
		}
	}(inputChan)
	return err
}
func (c *Communicator) sendCmdHandler(socket *zmq.Socket, cmd *cmdEvp) {
	if data, isOk := cmd.data.([]byte); isOk {
		Log.Debug("Cmd %v has been send to endpoint %s", cmd, cmd.ep.ToTcp())
		socket.SendBytes(data, 0)
		rData, err := socket.RecvBytes(0)
		if err == nil {
			c.outputChan <- rData
		} else {
			Log.Error("Socket err:%v", err)
		}
	} else {
		Log.Error("Cmd %v data cannot convert to []byte", cmd)
	}
}
func (c *Communicator) sendFileHandler(socket *zmq.Socket, cmd *cmdEvp) {
	if tfc, isOk := cmd.data.(*TransFileCmd); isOk {
		Log.Debug("file %v has been send to endpoint %s", tfc.Header, cmd.ep.ToTcp())
		var data,err=proto.Marshal(tfc.Header)
		if err!=nil {
			Log.Error("File heaher %v encode failed err:%v",tfc.Header,err)
			return
		}else{
			socket.SendBytes(data,zmq.SNDMORE)
		}
		rootPath,header:=tfc.RootPath,tfc.Header
		if _,err:=os.Stat(filepath.Join(rootPath,header.GetRelPath()));err==nil{
			file,err:=os.Open(filepath.Join(rootPath,header.GetRelPath()))
			if err!=nil{
				Log.Error("File %v cannot be open! err:%v",tfc,err)
				return
			}
			defer file.Close()
			var buffer=make([]byte,buffer_length)
			for {
				readCount,err:= file.Read(buffer)
				if err==nil {
					socket.SendBytes(buffer,zmq.SNDMORE)
				}else if err==io.EOF && readCount>0{
					socket.SendBytes(buffer[:readCount],0)
					break
				}else{
					Log.Error("Read file %v error err:%v ",file,err)
					socket.SendBytes([]byte{},0)
					break
				}
			}
			rData, err := socket.RecvBytes(0)
			if err == nil {
				c.outputChan <- rData
			} else {
				Log.Error("Socket err:%v", err)
			}
		}else{
			Log.Error("File %v cannot be found! err:%v",tfc,err)
		}
	} else {
		Log.Error("Cmd %v data cannot convert to []byte", tfc)
	}
}


func (c *Communicator)openReqEp(ep *EndPoint) error{
	var ipcId=fmt.Sprint("ipc://work-%s.ipc",uuid.New())
	for i:=0;i<max_threads;i++{
		go func(ipcId string){
			receiver, _ := zmq.NewSocket(zmq.REP)
		    defer receiver.Close()
		    receiver.Connect(ipcId)
		
//		    for true {
//		       	received, _ := receiver.RecvMessage()
		     	        
//		    }
			
		}(ipcId)
	}
	return nil
}

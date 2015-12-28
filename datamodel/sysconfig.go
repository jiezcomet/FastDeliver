package datamodel

import (
	. "FastDeliver/log"
	"encoding/json"
	"os"
)
const(
	Node_Server uint=iota+1
	Node_Client
	Node_ControlClient	
)
type ServerConfig struct{
	CmdReqPort uint
	CmdPubPort uint
	DataReqPort uint
	HttpPort uint 
	LanIp string
	WanIp string
	StorePath string
}

func LoadServerConfiguration(confPath string)(*ServerConfig,error){
	Log.Info("Start to load server configuration file file:%s",confPath)
	if _,err:=os.Stat(confPath);err==nil {
		var sc=&ServerConfig{}
		file, err := os.Open(confPath)
		if err != nil {
			Log.Error("File cannot be open  error:%v",err);
			return nil, err
		}
		dc := json.NewDecoder(file)
		if err = dc.Decode(&sc); err != nil {
			Log.Error("File cannot be decoded by json  error:%v",err);
			return nil, err
		} else {
			return sc,nil
		}
	}else{
		Log.Error("Configuration file cannot be found! path:%s",confPath)
		return nil,err
	}
}
func (s ServerConfig)SaveServerConfiguration(confPath string) error{
	b, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		Log.Error("Server configurataion %v cannot be encoded by json %v ",s,err)
		return err
	}
	file, err := os.Create(confPath)
	if err != nil {
		Log.Error("File create failed path%s err:%v ",confPath,err)
		return err
	}
	defer file.Close()
	_, err = file.Write(b)
	return err
}
package main

import (
	"FastDeliver/flib"
//	"path/filepath"
//	"bytes"
//	"os"
//	"fmt"
	"time"
//	"strings"
	"FastDeliver/log"
	dm "FastDeliver/datamodel"
)
var fileCount,repeatedCount int64=0,0

func main() {
//	var fileDic=make(map[string] []byte);
//	var repeDic=make(map[string] []byte);
//	var rootPath="E:/workspace/FD/RM2.0CI_CDC_20150921.2"
//	filepath.Walk(rootPath,func(p string,info os.FileInfo,e error) error{
//		if !info.IsDir(){
//			code,err:=flib.GetFileHashCode(p);
//			if err==nil{
//				if oCode,isExists:=fileDic[info.Name()];!isExists{
//					fileDic[info.Name()]=code
//				}else if bytes.Compare(oCode,code) !=0 {
//					repeatedCount++
//					repeDic[p]=code
//				}else{
//					repeatedCount++
//				}
//			}else{
//				fmt.Println(err)
//			}
			
//			fileCount++
			
			
//		}
//		return nil
//	})
//	fmt.Printf("file count:%d, re count:%d , recodue: %d\n",fileCount,repeatedCount,len(repeDic))
//	for k,c:=range repeDic{
//		fmt.Printf("File Name:%s ,code:%v\n",k,c)
//	}
	
//	var imagePath="E:/workspace/FD/image";
//	var image,err=flib.CreateImage("T1",imagePath)
//	if err!=nil{
//		fmt.Println(err)
//	}else{
//		fmt.Println(image)
//	}
//	nImage,err:=flib.LoadFromFile("E:/workspace/FD/image/configuration.json")
//	if err==nil{
//		fmt.Println(nImage)
//	}else{
//		fmt.Println(err)
//	}
//	var oldName="E:/workspace/FD/image/configuration.json"
//	var newName="E:/workspace/FD/clone/configuration.json"
//	var clonePath="E:/workspace/FD/clone";
//	os.Symlink(oldName,newName);
//	filepath.Walk(clonePath,func(p string,info os.FileInfo,e error) error{
//		fmt.Printf("File name:%s fileType %v \n",p,info.Mode().IsRegular())
//		return nil
//	})
	log.InitLog();
	log.Log.Debug( time.Now())
	
	var sc=&dm.ServerConfig{}
	sc.CmdPubPort,sc.CmdReqPort=9000,9100
	sc.DataReqPort,sc.HttpPort=10000,8080
	sc.LanIp,sc.WanIp="127.0.0.1","127.0.0.1"
	sc.SaveServerConfiguration("./config/serverconfiguration.json")
	
	var imagePath="E:\\workspace\\FD\\RM2.0CI_CDC_20151226.1";

	log.Log.Info( time.Now())
	var image,_=flib.CreateImageN("T12",imagePath)
	log.Log.Debug("FileCount :%d\n",len(image.RealFileHashMap))
	log.Log.Info( time.Now())
	log.Close()
	
	time.Sleep(time.Second)
}

package flib

import (
	"bytes"
	"encoding/json"
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	//"compress/flate"
	. "FastDeliver/log"
)

const (
	configFileName string = "/_image_configuration.json"
)
const (
	F_IsNotExists int = iota + 1
	F_LinkHashMatched
	F_LinkHashNotMatched
	F_HashNotMatched
	F_HashMatched
)

type Image struct {
	sync.RWMutex
	Name            string
	RootPath        string
	RealFileHashMap map[string][]byte
	FileLinkMap     map[string]string //key link target, value link source
	LastModifyDate  time.Time
}

func CreateImageN(name, path string) (*Image, error) {
	Log.Info("Initialize file image")
	if _, err := os.Stat(path); err == nil {
		var image = &Image{}
		image.Name, image.RootPath = name, path
		image.RealFileHashMap = make(map[string][]byte)
		image.FileLinkMap = make(map[string]string)
		image.LastModifyDate = time.Now()
		if err := image.scan(); err != nil {
			Log.Error("Image path scan process failed! err:%v",err)
			return nil, err
		} else {
			return image, image.flush()
		}
	} else {
		Log.Error("Image path cannot be found! err:%v",err)
		return nil, err
	}
}
func (i *Image) scan() error {
	Log.Info("start to scan the path %s",i.RootPath)
	var hashCodeMaps = make(map[string]map[string][]byte)

	err := filepath.Walk(i.RootPath, func(p string, info os.FileInfo, e error) error {
		if !info.IsDir() {
			var relPath = strings.TrimPrefix(p, i.RootPath)
			code, err := GetFileHashCode(p)
			if err == nil {
				Log.Debug("file related path:%s , hashcode is %s",relPath,base64.StdEncoding.EncodeToString(code))
				if inHashCodeMap, isExists := hashCodeMaps[info.Name()]; !isExists {
					inHashCodeMap = make(map[string][]byte)
					inHashCodeMap[relPath] = code
					hashCodeMaps[info.Name()] = inHashCodeMap
					i.RealFileHashMap[relPath] = code
				} else {
					var isMatched = false
					for linkSource, oCode := range inHashCodeMap {
						if bytes.Compare(oCode, code) == 0 {
							i.FileLinkMap[relPath] = linkSource
							isMatched = true
							break
						}
					}
					if !isMatched {
						i.RealFileHashMap[relPath] = code
						inHashCodeMap[relPath] = code
					}
				}
			}else{
				Log.Error("Hashcode cannot be sum for File (full path %s),err:%v ",p,err)
			}
		}
		return nil
	})
	return err
}
func (i Image) flush() error {
	var cfgFile = i.RootPath + configFileName
	i.RLock()
	defer i.RUnlock()
	b, err := json.MarshalIndent(i, "", "\t")
	if err == nil {
		if file, err := os.OpenFile(cfgFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666); err == nil {
			defer file.Close()
//			var cmpWriter,_=flate.NewWriter(file,flate.DefaultCompression)
//			defer cmpWriter.Close()
//			cmpWriter.Write(b)
//			cmpWriter.Flush()
			file.Write(b)
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}


package log

import (
	l4g "code.google.com/p/log4go"
)

var Log l4g.Logger

func InitLog(){
	Log=make(l4g.Logger)
	Log.LoadConfiguration("./log/log4go.xml");
}
func Close(){
	Log.Close()
}
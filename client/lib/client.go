package lib

import(
	//. "FastDeliver/log"
	"FastDeliver/mgrs"
)

type Client struct{
	C chan int // initialize wating gate
	IsReady bool
	ct *mgrs.Communicator
}
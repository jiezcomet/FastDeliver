package mgrs

import (
	dm "FastDeliver/datamodel"
)

type ConnMgr struct {
	clients map[string]*dm.ClientInfo
}

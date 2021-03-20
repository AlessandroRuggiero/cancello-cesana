package espconn

import (
	"fmt"
)

func (esp *EspConnection) Apricancello(id string) error {
	return esp.WriteMessage(1, []byte(fmt.Sprintf("0:%s", id)))
}

func (esp *EspConnection) Apricancelletto(id string) error {
	return esp.WriteMessage(1, []byte(fmt.Sprintf("1:%s", id)))
}

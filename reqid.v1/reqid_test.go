package reqid

import (
	"fmt"
	"testing"
)

func TestGenReqId(t *testing.T) {
	reqId := GenReqId()
	fmt.Println(reqId)

	tt := ParseTimeFromReqId(reqId)
	fmt.Println("时间：", tt)
}

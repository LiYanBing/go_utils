package reqid

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"net/http"
	"os"
	"time"

	"github.com/spaolacci/murmur3"
)

var (
	pid         = uint16(time.Now().UnixNano() & 65535)
	machineFlag uint16
)

const (
	reqIdHeaderKey = "ReqId"
)

func init() {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	machineFlag = uint16(murmur3.Sum32([]byte(hostname)) & 65535)
}

func GenReqId() string {
	var b [12]byte
	binary.LittleEndian.PutUint16(b[:], pid)
	binary.LittleEndian.PutUint16(b[2:], machineFlag)
	binary.LittleEndian.PutUint64(b[4:], uint64(time.Now().UnixNano()))
	return base64.URLEncoding.EncodeToString(b[:])
}

func ParseTimeFromReqId(reqId string) time.Time {
	b, err := base64.URLEncoding.DecodeString(reqId)
	if err != nil {
		panic(err)
	}
	b = b[4:]
	nano := int64(binary.LittleEndian.Uint64(b))
	return time.Unix(0, nano)
}

func NewContext(w http.ResponseWriter, req *http.Request) context.Context {
	reqId := GenReqId()
	req.Header.Set(reqIdHeaderKey, reqId)
	w.Header().Set(reqIdHeaderKey, reqId)
	return context.WithValue(req.Context(), reqIdHeaderKey, reqId)
}

func FromContext(ctx context.Context) (string, bool) {
	reqId, ok := ctx.Value(reqIdHeaderKey).(string)
	return reqId, ok
}

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := NewContext(w, req)
		next.ServeHTTP(w, req.WithContext(ctx))
	})
}

package util

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
)

func CheckValueErrInt(value int, err error) *wrappers.StringValue {
	return CheckValueErrStr(fmt.Sprint(value), err)
}

func CheckValueErrStr(value string, err error) *wrappers.StringValue {
	if err == nil && value != "" {
		return &wrappers.StringValue{Value: value}
	}

	return nil
}

func CheckErrUint64(value string, err error) *wrappers.UInt64Value {
	if err == nil && value != "" {
		var i uint64
		i, err = strconv.ParseUint(value, 10, 64)
		if err == nil {
			return &wrappers.UInt64Value{Value: i}
		}
	}

	return nil
}

func CheckTime(t *time.Time, err error) (*timestamp.Timestamp, error) {
	if err == nil && t != nil {
		var ts *timestamp.Timestamp
		ts, err = ptypes.TimestampProto(*t)
		if err == nil {
			return ts, nil
		}
	}

	return nil, err
}

func IsPublicIPv4(ip net.IP) bool {
	if ip == nil {
		return false
	}

	if ip.IsLoopback() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast() {
		return false
	}

	if ip4 := ip.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}

	return false
}

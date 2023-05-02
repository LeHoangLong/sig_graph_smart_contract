package encrypt

import (
	"bytes"
	"fmt"
	"reflect"
	"sig_graph/utility"
	"strconv"

	"github.com/shopspring/decimal"
)

func bytesToValue[T any](Bytes []byte) (T, utility.Error) {
	var t T
	var tAny any = t
	var err utility.Error
	var ret any
	var retUint64 uint64
	var retInt64 int64

	switch tAny.(type) {
	case string:
		ret = string(Bytes)
	case bool:
		ret = bytes.Equal(Bytes, []byte("true"))
	case uint:
		retUint64, err = bytesToUint64(Bytes)
		if err != nil {
			return t, err.AddMessage("failed to convert byte to uint")
		}
		ret = uint(retUint64)
	case uint8:
		retUint64, err = bytesToUint64(Bytes)
		if err != nil {
			return t, err.AddMessage("failed to convert byte to uint8")
		}
		ret = uint8(retUint64)
	case uint16:
		retUint64, err = bytesToUint64(Bytes)
		if err != nil {
			return t, err.AddMessage("failed to convert byte to uint16")
		}
		ret = uint16(retUint64)
	case uint32:
		retUint64, err = bytesToUint64(Bytes)
		if err != nil {
			return t, err.AddMessage("failed to convert byte to uint32")
		}
		ret = uint32(retUint64)
	case uint64:
		retUint64, err = bytesToUint64(Bytes)
		if err != nil {
			return t, err.AddMessage("failed to convert byte to uint64")
		}
		ret = uint64(retUint64)
	case int:
		retInt64, err = bytesToInt64(Bytes)
		if err != nil {
			return t, err.AddMessage("failed to convert byte to int")
		}
		ret = int(retInt64)
	case int8:
		retInt64, err = bytesToInt64(Bytes)
		if err != nil {
			return t, err.AddMessage("failed to convert byte to int8")
		}
		ret = int8(retInt64)
	case int16:
		retInt64, err = bytesToInt64(Bytes)
		if err != nil {
			return t, err.AddMessage("failed to convert byte to int16")
		}
		ret = int16(retInt64)
	case int32:
		retInt64, err = bytesToInt64(Bytes)
		if err != nil {
			return t, err.AddMessage("failed to convert byte to int32")
		}
		ret = int32(retInt64)
	case int64:
		retInt64, err = bytesToInt64(Bytes)
		if err != nil {
			return t, err.AddMessage("failed to convert byte to int64")
		}
		ret = int64(retInt64)
	case decimal.Decimal:
		valStr := string(Bytes)
		var err error
		ret, err = decimal.NewFromString(valStr)
		if err != nil {
			return t, utility.NewError(err).AddMessage("failed to convert byte to Decimal")
		}
	default:
		return t, utility.NewError(utility.ErrTypeError).AddMessage(fmt.Sprintf("unrecognized type %s", reflect.TypeOf(t)))
	}

	if val, ok := ret.(T); !ok {
		return t, utility.NewError(utility.ErrTypeError)
	} else {
		return val, nil
	}
}

func bytesToUint64(bytes []byte) (uint64, utility.Error) {
	ret, err := strconv.ParseUint(string(bytes), 10, 64)
	if err != nil {
		return ret, utility.NewError(utility.ErrTypeError).AddMessage(fmt.Sprintf("cannot convert %s to uint64", string(bytes)))
	}

	return ret, nil
}

func bytesToInt64(bytes []byte) (int64, utility.Error) {
	ret, err := strconv.ParseInt(string(bytes), 10, 64)
	if err != nil {
		return ret, utility.NewError(utility.ErrTypeError).AddMessage(fmt.Sprintf("cannot convert %s to int64", string(bytes)))
	}

	return ret, nil
}

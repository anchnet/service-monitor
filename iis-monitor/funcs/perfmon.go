package funcs

import (
	"errors"
	"fmt"
	"unsafe"

	"github.com/lxn/win"
)

func ReadPerformanceCounter(counter string) (float64, error) {

	var queryHandle win.PDH_HQUERY
	var counterHandle win.PDH_HCOUNTER

	ret := win.PdhOpenQuery(0, 0, &queryHandle)
	if ret != win.ERROR_SUCCESS {
		return 0, errors.New("Unable to open query through DLL call")
	}

	// test path
	ret = win.PdhValidatePath(counter)
	if ret == win.PDH_CSTATUS_BAD_COUNTERNAME {
		return 0, errors.New("Unable to fetch counter (this is unexpected)")
	}

	ret = win.PdhAddCounter(queryHandle, counter, 0, &counterHandle)
	if ret != win.ERROR_SUCCESS {
		return 0, errors.New(fmt.Sprintf("Unable to add process counter. Error code is %x\n", ret))
	}

	ret = win.PdhCollectQueryData(queryHandle)
	if ret != win.ERROR_SUCCESS {
		return 0, errors.New(fmt.Sprintf("Got an error: 0x%x\n", ret))
	}

	var bufSize uint32
	var bufCount uint32
	var size uint32 = uint32(unsafe.Sizeof(win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE{}))
	var emptyBuf [1]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE // need at least 1 addressable null ptr.
	var v float64
	ret = win.PdhGetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &emptyBuf[0])
	if ret == win.PDH_MORE_DATA {
		filledBuf := make([]win.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE, bufCount*size)
		ret = win.PdhGetFormattedCounterArrayDouble(counterHandle, &bufSize, &bufCount, &filledBuf[0])
		if ret == win.ERROR_SUCCESS {
			for i := 0; i < int(bufCount); i++ {
				c := filledBuf[i]
				//				s := win.UTF16PtrToString(c.SzName)
				v = c.FmtValue.DoubleValue
				//fmt.Printf("Index %d -> %s, value %v\n", i, s, c.FmtValue.DoubleValue)
			}
			filledBuf = nil
			bufCount = 0
			bufSize = 0
		}
	}

	return v, nil
}

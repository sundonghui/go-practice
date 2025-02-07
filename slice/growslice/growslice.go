package growslice

import "fmt"

func GrowSlice() {
	var s []int
	lastCap := cap(s)

	for i := 0; i < 10000; i++ {
		s = append(s, i)
		newCap := cap(s)

		if newCap != lastCap {
			fmt.Printf("len: %d, cap: %d -> %d (增长倍数: %.2f)\n",
				len(s), lastCap, newCap, float64(newCap)/float64(lastCap))
			lastCap = newCap
		}
	}

}

// go version go1.23.4 darwin/amd64
// newCap = oldCap + (oldCap + 3*256) / 4
/*
len: 2, cap: 1 -> 2 (增长倍数: 2.00)
len: 3, cap: 2 -> 4 (增长倍数: 2.00)
len: 5, cap: 4 -> 8 (增长倍数: 2.00)
len: 9, cap: 8 -> 16 (增长倍数: 2.00)
len: 17, cap: 16 -> 32 (增长倍数: 2.00)
len: 33, cap: 32 -> 64 (增长倍数: 2.00)
len: 65, cap: 64 -> 128 (增长倍数: 2.00)
len: 129, cap: 128 -> 256 (增长倍数: 2.00)
len: 257, cap: 256 -> 512 (增长倍数: 2.00)
len: 513, cap: 512 -> 848 (增长倍数: 1.66)
len: 849, cap: 848 -> 1280 (增长倍数: 1.51)
len: 1281, cap: 1280 -> 1792 (增长倍数: 1.40)
len: 1793, cap: 1792 -> 2560 (增长倍数: 1.43)
len: 2561, cap: 2560 -> 3408 (增长倍数: 1.33)
len: 3409, cap: 3408 -> 5120 (增长倍数: 1.50)
len: 5121, cap: 5120 -> 7168 (增长倍数: 1.40)
len: 7169, cap: 7168 -> 9216 (增长倍数: 1.29)
len: 9217, cap: 9216 -> 12288 (增长倍数: 1.33)
*/

/*
func growslice(oldPtr unsafe.Pointer, newLen, oldCap, num int, et *_type) slice {
	...
	// Specialize for common values of et.size.
	switch {
	case et.Size_ == 1:
		lenmem = uintptr(oldLen)
		newlenmem = uintptr(newLen)
		capmem = roundupsize(uintptr(newcap), noscan)
		overflow = uintptr(newcap) > maxAlloc
		newcap = int(capmem)
	case et.Size_ == goarch.PtrSize:
		lenmem = uintptr(oldLen) * goarch.PtrSize
		newlenmem = uintptr(newLen) * goarch.PtrSize
		capmem = roundupsize(uintptr(newcap)*goarch.PtrSize, noscan)
		overflow = uintptr(newcap) > maxAlloc/goarch.PtrSize
		newcap = int(capmem / goarch.PtrSize)
	case isPowerOfTwo(et.Size_):
		var shift uintptr
		if goarch.PtrSize == 8 {
			// Mask shift for better code generation.
			shift = uintptr(sys.TrailingZeros64(uint64(et.Size_))) & 63
		} else {
			shift = uintptr(sys.TrailingZeros32(uint32(et.Size_))) & 31
		}
		lenmem = uintptr(oldLen) << shift
		newlenmem = uintptr(newLen) << shift
		capmem = roundupsize(uintptr(newcap)<<shift, noscan)
		overflow = uintptr(newcap) > (maxAlloc >> shift)
		newcap = int(capmem >> shift)
		capmem = uintptr(newcap) << shift
	default:
		lenmem = uintptr(oldLen) * et.Size_
		newlenmem = uintptr(newLen) * et.Size_
		capmem, overflow = math.MulUintptr(et.Size_, uintptr(newcap))
		capmem = roundupsize(capmem, noscan)
		newcap = int(capmem / et.Size_)
		capmem = uintptr(newcap) * et.Size_
	}

	...
}

func roundupsize(size uintptr, noscan bool) (reqSize uintptr) {
	reqSize = size
	if reqSize <= maxSmallSize-mallocHeaderSize {
		// Small object.
		if !noscan && reqSize > minSizeForMallocHeader { // !noscan && !heapBitsInSpan(reqSize)
			reqSize += mallocHeaderSize
		}
		// (reqSize - size) is either mallocHeaderSize or 0. We need to subtract mallocHeaderSize
		// from the result if we have one, since mallocgc will add it back in.
		if reqSize <= smallSizeMax-8 {
			return uintptr(class_to_size[size_to_class8[divRoundUp(reqSize, smallSizeDiv)]]) - (reqSize - size)
		}
		return uintptr(class_to_size[size_to_class128[divRoundUp(reqSize-smallSizeMax, largeSizeDiv)]]) - (reqSize - size)
	}
	// Large object. Align reqSize up to the next page. Check for overflow.
	reqSize += pageSize - 1
	if reqSize < size {
		return size
	}
	return reqSize &^ (pageSize - 1)
}
*/

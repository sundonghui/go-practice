package leetcode

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"
)

type primeInfo struct {
	num    int64
	worker int
	time   time.Duration
}

type chunk struct {
	start int64
	size  int
	index int
}

type result struct {
	primes []primeInfo
	index  int
}

func numberNPrime(start, number int) {
	startNum := int64(start)
	required := number
	chunkSize := 500 // 每个块处理的候选数数量
	workers := 8     // 根据CPU核心数调整

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chunkCh := make(chan chunk)
	resultCh := make(chan result)

	var wg sync.WaitGroup
	// 启动worker池
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for c := range chunkCh {
				processChunk(c.start, c.size, workerID)
				primes := processChunk(c.start, c.size, workerID)
				select {
				case resultCh <- result{primes, c.index}:
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}

	// 生成块
	go func() {
		defer close(chunkCh)
		index := 0
		currentStart := startNum
		for {
			select {
			case <-ctx.Done():
				return
			default:
				c := chunk{start: currentStart, size: chunkSize, index: index}
				select {
				case chunkCh <- c:
					currentStart += 2 * int64(chunkSize)
					index++
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	// 收集结果
	var primes []primeInfo
	currentIndex := 0
	cache := make(map[int][]primeInfo)
	for res := range resultCh {
		if res.index != currentIndex {
			cache[res.index] = res.primes
			continue
		}
		primes = append(primes, res.primes...)
		currentIndex++
		// 检查缓存
		for {
			cachedPrimes, ok := cache[currentIndex]
			if !ok {
				break
			}
			primes = append(primes, cachedPrimes...)
			delete(cache, currentIndex)
			currentIndex++
		}
		if len(primes) >= required {
			cancel()
			break
		}
	}

	defer close(resultCh)
	wg.Wait()

	if len(primes) > required {
		primes = primes[:required]
	}

	// 输出结果示例
	fmt.Printf("Found %d primes\n", len(primes))
	for i, p := range primes {
		if i < 5 || i >= len(primes)-5 {
			fmt.Printf("Prime %d: %d (Worker: %d, Time: %v)\n", i+1, p.num, p.worker, p.time)
		}
	}
}

// 分块处理，按顺序收集结果，保证结果顺序正确
func processChunk(start int64, size int, workerID int) []primeInfo {
	var primes []primeInfo
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < size; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			num := start + 2*int64(i)
			startTime := time.Now()
			if isPrime(num) {
				elapsed := time.Since(startTime)
				mu.Lock()
				primes = append(primes, primeInfo{num, workerID, elapsed})
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()

	// 排序块内的质数
	sort.Slice(primes, func(i, j int) bool {
		return primes[i].num < primes[j].num
	})

	return primes
}

// 使用米勒-拉宾测试来快速判断质数
func isPrime(n int64) bool {
	if n <= 1 {
		return false
	}
	if n <= 3 {
		return true
	}
	if n%2 == 0 {
		return false
	}

	bigN := big.NewInt(n)
	d := new(big.Int).Sub(bigN, big.NewInt(1))
	s := 0
	for d.Bit(0) == 0 {
		d.Rsh(d, 1)
		s++
	}

	bases := []int64{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37}
	for _, base := range bases {
		if base >= n {
			continue
		}
		a := big.NewInt(base)
		x := new(big.Int).Exp(a, d, bigN)
		if x.Cmp(big.NewInt(1)) == 0 || x.Cmp(new(big.Int).Sub(bigN, big.NewInt(1))) == 0 {
			continue
		}
		composite := true
		for i := 0; i < s-1; i++ {
			x.Exp(x, big.NewInt(2), bigN)
			if x.Cmp(new(big.Int).Sub(bigN, big.NewInt(1))) == 0 {
				composite = false
				break
			}
		}
		if composite {
			return false
		}
	}
	return true
}

package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Batch struct {
	size  int
	flags Bitmask
}

type Bitmask uint32

func (f Bitmask) HasFlag(flag Bitmask) bool { return f&flag != 0 }
func (f *Bitmask) AddFlag(flag Bitmask)     { *f |= flag }
func (f *Bitmask) ClearFlag(flag Bitmask)   { *f &= ^flag }
func (f *Bitmask) ToggleFlag(flag Bitmask)  { *f ^= flag }

const (
	SELECT Bitmask = 1 << iota
	BUBBLE
	INSERT
	QUICK
	QUICK_P
	COUNT
	MERGE
	MERGE_P
	HEAP
	SORT_INTS
)

const RANDOM_INT = 99999

func main() {
	//bs := []int{1000, 2000, 4000, 16000, 256000, 1000000, 10000000, 100000000, 1000000000}
	bs := []int{1000, 2000, 4000, 16000, 256000, 1000000, 10000000, 100000000}
	batches := []Batch{}

	for _, val := range bs {
		f := getFlags(val)
		if f.flags != 0 {
			batches = append(batches, f)
		}
	}

	results := map[string][]time.Duration{}

	for _, batch := range batches {
		arr := gen(batch.size)
		r := exec(batch.flags, arr)
		for key, val := range r {
			results[key] = append(results[key], val)
		}
	}
	print(batches, results)
}

func getFlags(size int) Batch {
	var flags Bitmask

	switch size {
	case 1000:
		fallthrough
	case 2000:
		fallthrough
	case 4000:
		fallthrough
	case 16000:
		flags.AddFlag(SELECT)
		flags.AddFlag(BUBBLE)
		flags.AddFlag(INSERT)
		flags.AddFlag(QUICK)
		flags.AddFlag(COUNT)
		flags.AddFlag(MERGE)
		flags.AddFlag(MERGE_P)
		flags.AddFlag(HEAP)
		flags.AddFlag(SORT_INTS)
	case 256000:
		fallthrough
	case 1000000:
		flags.AddFlag(QUICK)
		flags.AddFlag(COUNT)
		flags.AddFlag(MERGE)
		flags.AddFlag(MERGE_P)
		flags.AddFlag(HEAP)
		flags.AddFlag(SORT_INTS)
	case 10000000:
		flags.AddFlag(COUNT)
		flags.AddFlag(MERGE)
		flags.AddFlag(MERGE_P)
		flags.AddFlag(HEAP)
		flags.AddFlag(SORT_INTS)
	case 100000000:
		flags.AddFlag(COUNT)
		flags.AddFlag(MERGE_P)
	case 1000000000:
		flags.AddFlag(MERGE_P)
	default:

	}
	var batch Batch
	batch.size = size
	batch.flags = flags
	return batch

}

func print(batches []Batch, results map[string][]time.Duration) {
	p := message.NewPrinter(language.English)
	table := tablewriter.NewWriter(os.Stdout)
	header := []string{"algo"}
	for _, b := range batches {
		header = append(header, p.Sprintf("%d", b.size))
	}
	table.SetHeader(header)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)

	for key, val := range results {
		line := []string{key}
		for idx, d := range val {
			if isFastest(d, idx, results) && d.String() != "0s" {
				line = append(line, "\033[0;32m"+d.String()+"\033[0m")
			} else {
				if d.String() == "0s" {
					line = append(line, "-")
					continue
				}
				line = append(line, d.String())
			}
		}
		table.Append(line)
	}
	table.Render()
}

func isFastest(d time.Duration, idx int, results map[string][]time.Duration) bool {
	res := true
	for _, vals := range results {
		if idx > len(vals)-1 {
			continue
		}
		if d > vals[idx] && (vals[idx]).String() != "0s" {
			res = false
		}
	}
	return res
}

func exec(flags Bitmask, arr []int) map[string]time.Duration {
	var m string
	var t time.Time

	results := map[string]time.Duration{}

	if flags.HasFlag(SELECT) {
		m, t = track("select")
		selectionSortResult := selectionSort(cp(arr))
		s, took := dur(m, t)
		results[s] = took
		if !sort.IntsAreSorted(selectionSortResult) {
			fmt.Println("selectin sort not sorted")
		}
	} else {
		results["select"], _ = time.ParseDuration("0")
	}

	if flags.HasFlag(BUBBLE) {
		m, t = track("bubble")
		bubbleSortResult := bubbleSort(cp(arr))
		s, took := dur(m, t)
		results[s] = took
		if !sort.IntsAreSorted(bubbleSortResult) {
			fmt.Println("bubble sort not sorted")
		}
	} else {
		results["bubble"], _ = time.ParseDuration("0")
	}

	if flags.HasFlag(INSERT) {
		m, t = track("insert")
		insertionSortResult := insertionSort(cp(arr))
		s, took := dur(m, t)
		results[s] = took

		if sort.IntsAreSorted(insertionSortResult) == false {
			fmt.Println("insertion sort not sorted")
		}
	} else {
		results["insert"], _ = time.ParseDuration("0")
	}

	if flags.HasFlag(QUICK) {
		m, t = track("quick")
		quickSortResult := quickSort(cp(arr), 0, len(arr)-1)
		s, took := dur(m, t)
		results[s] = took
		if !sort.IntsAreSorted(quickSortResult) {
			fmt.Println("quick sort not sorted")
		}
	} else {
		results["quick"], _ = time.ParseDuration("0")
	}

	if flags.HasFlag(COUNT) {
		m, t = track("count")
		countSortResult := countSort(cp(arr))
		s, took := dur(m, t)
		results[s] = took
		if !sort.IntsAreSorted(countSortResult) {
			fmt.Println("count sort not sorted")
		}
	} else {
		results["count"], _ = time.ParseDuration("0")
	}

	//if flags.HasFlag(QUICK_P) {
	//m, t = track("quick parallel")
	//quickSortParallelResult := quickSortParallel(cp(arr), 0, len(arr)-1)
	//s, took := dur(m, t)
	//results[s] = took

	//if !sort.IntsAreSorted(quickSortParallelResult) {
	//fmt.Println("quick sort parallel not sorted")
	//}
	//} else {
	//results["quick p"], _ = time.ParseDuration("0")
	//}

	if flags.HasFlag(MERGE) {
		m, t = track("merge")
		mergeSortResult := mergeSort(cp(arr))
		s, took := dur(m, t)
		results[s] = took

		if !sort.IntsAreSorted(mergeSortResult) {
			fmt.Println("merge sort not sorted")
		}
	} else {
		results["merge"], _ = time.ParseDuration("0")
	}

	if flags.HasFlag(MERGE_P) {
		m, t = track("merge p")
		mergeSortParallelResult := mergeSortParallel(cp(arr))
		s, took := dur(m, t)
		results[s] = took

		if !sort.IntsAreSorted(mergeSortParallelResult) {
			fmt.Println("merge parallel sort not sorted")
		}
	} else {
		results["merge p"], _ = time.ParseDuration("0")
	}

	if flags.HasFlag(HEAP) {
		m, t = track("heap")
		heap := newHeap(cp(arr))
		heap.sort()
		s, took := dur(m, t)
		results[s] = took

		if !sort.IntsAreSorted(heap.data) {
			fmt.Println("heap sort not sorted")
		}
	} else {
		results["heap"], _ = time.ParseDuration("0")
	}

	if flags.HasFlag(SORT_INTS) {
		m, t = track("sort.Ints")
		sort.Ints(arr)
		s, took := dur(m, t)
		results[s] = took
	} else {
		results["select"], _ = time.ParseDuration("0")
	}

	return results

}

func countSort(arr []int) []int {
	// find max
	max := math.MinInt32
	min := math.MaxInt32
	for _, val := range arr {
		if max < val {
			max = val
		}
		if min > val {
			min = val
		}
	}

	// generate array from min to max
	counter := make([]int, max-min+1)

	// count
	for _, val := range arr {
		counter[val-min] += 1
	}

	// add previous to curr
	for i := 1; i < len(counter); i++ {
		counter[i] += counter[i-1]
	}

	res := make([]int, len(arr))
	// copy to correct pos
	for i := 0; i < len(arr); i++ {
		elem := arr[i]
		t := counter[elem-min] - 1
		res[t] = elem
		counter[elem-min] = counter[elem-min] - 1
	}

	return res
}

func mergeSortParallel(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	if len(arr) == 2 {
		if arr[0] > arr[1] {
			arr[0], arr[1] = arr[1], arr[0]
		}
		return arr
	}
	if len(arr) < 2048 {

		return merge(mergeSortParallel(arr[:len(arr)/2]), mergeSortParallel(arr[len(arr)/2:]))
	} else {
		var wg sync.WaitGroup
		wg.Add(2)

		var a, b []int
		go func() {
			defer wg.Done()
			a = mergeSortParallel(arr[:len(arr)/2])
		}()

		go func() {
			defer wg.Done()
			b = mergeSortParallel(arr[len(arr)/2:])
		}()

		return merge(a, b)
	}
}

func mergeSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	if len(arr) == 2 {
		if arr[0] > arr[1] {
			arr[0], arr[1] = arr[1], arr[0]
		}
		return arr
	}
	return merge(mergeSort(arr[:len(arr)/2]), mergeSort(arr[len(arr)/2:]))
}

func merge(a []int, b []int) []int {

	if len(a) < 1 {
		return b
	}
	if len(b) < 1 {
		return a
	}

	c := make([]int, len(a)+len(b))

	j, k := 0, 0
	for i := 0; i < len(c); i++ {
		if j < len(a) && k < len(b) {
			if a[j] < b[k] {
				c[i] = a[j]
				j++

			} else {
				c[i] = b[k]
				k++
			}
			continue
		}

		if k < len(b) {
			c[i] = b[k]
			k++
		}

		if j < len(a) {
			c[i] = a[j]
			j++
		}

	}

	return c
}

func cp(arr []int) []int {
	ret := make([]int, len(arr))
	copy(ret, arr)
	return ret
}

func quickSortParallel(arr []int, low int, high int) []int {
	//fmt.Println("quicks", arr)
	//	defer duration(track("quick"))
	if low < high {

		pivot := arr[high]
		i := low - 1
		for j := low; j <= high-1; j++ {
			//fmt.Println(arr[j], pivot, i)

			if arr[j] < pivot {
				i++
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
		arr[i+1], arr[high] = arr[high], arr[i+1]
		idx := i + 1

		if idx-1-low < 2048 {

			quickSort(arr, low, idx-1)
			quickSort(arr, idx+1, high)
		} else {
			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer wg.Done()
				quickSort(arr, low, idx-1)
			}()

			go func() {
				defer wg.Done()
				quickSort(arr, idx+1, high)
			}()
			wg.Wait()
		}
	}

	return arr
}

func quickSort(arr []int, low int, high int) []int {
	//fmt.Println("quicks", arr)
	//	defer duration(track("quick"))
	if low < high {

		pivot := arr[high]
		i := low - 1
		for j := low; j <= high-1; j++ {
			//fmt.Println(arr[j], pivot, i)

			if arr[j] < pivot {
				i++
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
		arr[i+1], arr[high] = arr[high], arr[i+1]
		idx := i + 1
		quickSort(arr, low, idx-1)
		quickSort(arr, idx+1, high)
	}

	return arr
}

func insertionSort(arr []int) []int {
	//fmt.Println("insert", arr)
	//	defer duration(track("insert"))
	for i := 1; i < len(arr); i++ {
		if arr[i] < arr[i-1] {
			temp := arr[i]
			j := i
			for ; j > 0 && arr[j-1] > temp; j-- {
				arr[j] = arr[j-1]
			}
			arr[j] = temp
		}
	}
	return arr
}

func bubbleSort(arr []int) []int {
	//fmt.Println("bubble", arr)
	//	defer duration(track("bubble"))
	for i := len(arr) - 1; i > 0; i-- {
		for j := 0; j < i; j++ {
			if arr[j] > arr[j+1] {
				arr[j+1], arr[j] = arr[j], arr[j+1]
			}
		}
	}
	return arr
}

func selectionSort(arr []int) []int {
	//fmt.Println("select", arr)
	//	defer duration(track("select"))
	for i := 0; i < len(arr); i++ {
		j := i + 1
		minIdx := i
		for ; j < len(arr); j++ {
			if arr[j] < arr[minIdx] {
				minIdx = j
			}
		}
		arr[minIdx], arr[i] = arr[i], arr[minIdx]
	}
	return arr
}

func gen(size int) []int {
	//	defer duration(track("gen"))
	arr := make([]int, size)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < size; i++ {
		arr[i] = rand.Intn(RANDOM_INT) - rand.Intn(RANDOM_INT)
	}

	return arr
}

func track(msg string) (string, time.Time) {
	return msg, time.Now()
}

func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

func dur(msg string, start time.Time) (string, time.Duration) {
	return msg, time.Since(start)
}

type Heap struct {
	data []int
}

func (h *Heap) build(arr []int) {
	h.data = append([]int{}, arr...)

	for i := len(h.data) / 2; i >= 0; i-- {
		h.heapify(i, len(h.data))
	}
}

func (h *Heap) sort() {
	for l := len(h.data); l > 1; l-- {
		h.removeTop(l)
	}
}

func (h *Heap) removeTop(i int) {
	lastIndex := i - 1
	h.data[0], h.data[lastIndex] = h.data[lastIndex], h.data[0]
	h.heapify(0, lastIndex)
}

func (h *Heap) heapify(root int, length int) {
	max := root
	l := h.left(root)
	r := h.right(root)

	if l < length && h.data[l] > h.data[max] {
		max = l
	}

	if r < length && h.data[r] > h.data[max] {
		max = r
	}

	if max != root {
		h.data[root], h.data[max] = h.data[max], h.data[root]
		h.heapify(max, length)
	}
}

func (h *Heap) left(root int) int {
	return (root * 2) + 1
}
func (h *Heap) right(root int) int {
	return (root * 2) + 2
}
func newHeap(arr []int) *Heap {
	h := &Heap{}
	h.build(arr)
	return h
}

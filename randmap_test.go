package randmap

import (
	"bytes"
	"compress/gzip"
	"math/rand"
	"runtime"
	"testing"
)

// builtinInitKey selects a key by ranging over m and returning the key at the
// initial position of the iterator. It is biased when the map is sparse or
// has overflow buckets. It runs in O(1) space and time.
func builtinInitKey(m map[int]int) (n int) {
	for n = range m {
		return
	}
	panic("empty map")
}

// builtinSeekKey selects a key by advancing the map iterator a random number
// of times. It is unbiased. It runs in O(1) space and O(n) time.
func builtinSeekKey(m map[int]int) (n int) {
	r := rand.Intn(len(m)) + 1
	for n = range m {
		if r--; r <= 0 {
			return
		}
	}
	panic("empty map")
}

// builtinFlattenKey selects a key by flattening the map into a slice of its keys
// and selecting a random index. It is unbiased. It runs in O(n) space and
// O(n) time.
func builtinFlattenKey(m map[int]int) int {
	flat := make([]int, 0, len(m))
	for n := range m {
		flat = append(flat, n)
	}
	return flat[rand.Intn(len(flat))]
}

func TestBuiltinMapInit(t *testing.T) {
	const iters = 100000
	m := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	}
	counts := make([]int, len(m))
	for i := 0; i < iters; i++ {
		counts[builtinInitKey(m)]++
	}
	// 0 should be selected 45-55% of the time
	if (iters/2-iters/20) > counts[0] || counts[0] > (iters/2+iters/20) {
		t.Errorf("expected builtin map to be less random: expected ~%v for elem 0, got %v", iters/2, counts[0])
	}
}

func TestBuiltinMapSeek(t *testing.T) {
	const iters = 100000
	m := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	}
	counts := make([]int, len(m))
	for i := 0; i < iters; i++ {
		counts[builtinSeekKey(m)]++
	}

	// should be unbiased
	for n, c := range counts {
		if (iters/len(m))/2 > c || c > (iters/len(m))*2 {
			t.Errorf("suspicious count: expected %v-%v, got %v (%v)", (iters/len(m))/2, (iters/len(m))*2, c, n)
		}
	}
}

func TestBuiltinMapFlatten(t *testing.T) {
	const iters = 100000
	m := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	}
	counts := make([]int, len(m))
	for i := 0; i < iters; i++ {
		counts[builtinFlattenKey(m)]++
	}

	// should be unbiased
	for n, c := range counts {
		if (iters/len(m))/2 > c || c > (iters/len(m))*2 {
			t.Errorf("suspicious count: expected %v-%v, got %v (%v)", (iters/len(m))/2, (iters/len(m))*2, c, n)
		}
	}
}

func TestKey(t *testing.T) {
	const iters = 100000
	m := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
		6: 6,
		7: 7,
		8: 8,
		9: 9,
	}
	counts := make([]int, len(m))
	for i := 0; i < iters; i++ {
		counts[Key(m).(int)]++
	}

	for n, c := range counts {
		if (iters/len(m))/2 > c || c > (iters/len(m))*2 {
			t.Errorf("suspicious count: expected %v-%v, got %v (%v)", (iters/len(m))/2, (iters/len(m))*2, c, n)
		}
	}
}

func TestVal(t *testing.T) {
	const iters = 100000
	m := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
		6: 6,
		7: 7,
		8: 8,
		9: 9,
	}
	counts := make([]int, len(m))
	for i := 0; i < iters; i++ {
		counts[Val(m).(int)]++
	}

	for n, c := range counts {
		if (iters/len(m))/2 > c || c > (iters/len(m))*2 {
			t.Errorf("suspicious count: expected %v-%v, got %v (%v)", (iters/len(m))/2, (iters/len(m))*2, c, n)
		}
	}
}

func TestFastKey(t *testing.T) {
	const iters = 100000
	m := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
		6: 6,
		7: 7,
		8: 8,
		9: 9,
	}
	counts := make([]int, len(m))
	for i := 0; i < iters; i++ {
		counts[FastKey(m).(int)]++
	}

	for n, c := range counts {
		if (iters/len(m))/2 > c || c > (iters/len(m))*2 {
			t.Errorf("suspicious count: expected %v-%v, got %v (%v)", (iters/len(m))/2, (iters/len(m))*2, c, n)
		}
	}
}

func TestGhostIndex(t *testing.T) {
	// sometimes, an element is never selected despite thousands of
	// iterations. This affects the builtin map range as well.
	const outer = 1000
	const inner = 1000
	for i := 0; i < outer; i++ {
		m := map[int]int{
			0: 0,
			1: 1,
			2: 2,
			3: 3,
			4: 4,
			5: 5,
			6: 6,
			7: 7,
			8: 8,
		}
		counts := make([]int, len(m))
		for j := 0; j < inner; j++ {
			counts[FastKey(m).(int)]++
		}

		for n, c := range counts {
			if c == 0 {
				t.Fatalf("%v: key %v was never selected!", i, n)
			}
		}
	}
}

func TestInsert(t *testing.T) {
	// Go maps incrementally copy values after each resize. The randmap
	// functions should continue to work in the middle of an incremental copy.
	const outer = 100
	const inner = 100
	m := make(map[int]int)
	for i := 0; i < outer; i++ {
		m[i] = i

		counts := make([]int, len(m))
		for j := 0; j < inner*len(m); j++ {
			counts[FastKey(m).(int)]++
		}

		for n, c := range counts {
			if inner/2 > c || c > inner*2 {
				t.Errorf("suspicious count: expected %v-%v, got %v (%v)", inner/2, inner*2, c, n)
			}
		}
	}
}

func TestFastVal(t *testing.T) {
	const iters = 100000
	m := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	}
	counts := make([]int, len(m))
	for i := 0; i < iters; i++ {
		counts[FastVal(m).(int)]++
	}

	for n, c := range counts {
		if (iters/len(m))/2 > c || c > (iters/len(m))*2 {
			t.Errorf("suspicious count: expected %v-%v, got %v (%v)", (iters/len(m))/2, (iters/len(m))*2, c, n)
		}
	}
}

func TestEmpty(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when accessing empty map")
		}
	}()
	_ = Key(make(map[int]int))
}

func TestEntropy(t *testing.T) {
	m := make(map[int]int)
	for j := 0; j < 417; j++ { // 417 chosen to ensure unevacuated map buckets
		m[j] = j
	}
	b := make([]byte, 10000)
	for j := range b {
		b[j] = byte(FastKey(m).(int))
	}
	var buf bytes.Buffer
	w, _ := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	w.Write(b)
	w.Close()
	if buf.Len() < len(b) {
		t.Fatalf("gzip was able to compress random keys by %.2f%%! (%v total bytes)", float64(100*buf.Len())/float64(len(b)), buf.Len())
	}
}

func TestGC(t *testing.T) {
	// if pointers are not preserved, overflow buckets can get garbage
	// collected. The randmap functions should continue to work in the
	// presence of GC.
	m := make(map[int]int)
	for i := 0; i < 1000; i++ {
		m[i] = i
		FastKey(m)
		runtime.GC()
	}
}

func TestConcurrent(t *testing.T) {
	m := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
	}
	const iters = 10000
	go func() {
		for i := 0; i < iters/len(m); i++ {
			for range m {
			}
		}
	}()
	for i := 0; i < 10; i++ {
		go func() {
			for i := 0; i < iters; i++ {
				Key(m)
			}
		}()
	}
	for i := 0; i < iters; i++ {
		Key(m)
	}
}

func BenchmarkKey(b *testing.B) {
	m := make(map[int]int, 10000)
	for i := 0; i < 10000; i++ {
		m[i] = i
	}

	b.Run("key", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = Key(m).(int)
		}
	})

	b.Run("fastkey", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = FastKey(m).(int)
		}
	})

	b.Run("seek", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = builtinSeekKey(m)
		}
	})

	b.Run("flatten", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = builtinFlattenKey(m)
		}
	})
}

func TestIter(t *testing.T) {
	const iters = 1000
	m := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
		6: 6,
		7: 7,
		8: 8,
		9: 9,
	}
	counts := make([][]int, len(m))
	for i := range counts {
		counts[i] = make([]int, len(m))
	}
	var k, v int
	for i := 0; i < iters; i++ {
		it := Iter(m, &k, &v)
		for j := 0; it.Next(); j++ {
			// key k appeared at index j
			counts[k][j]++
		}
	}

	// each key should have appeared at each index about iters/len(m) times
	for k, cs := range counts {
		for i, c := range cs {
			if (iters/len(m))/2 > c || c > (iters/len(m))*2 {
				t.Errorf("suspicious count for key %v index %v: expected %v-%v, got %v", k, i, (iters/len(m))/2, (iters/len(m))*2, c)
			}
		}
	}
}

func TestFastIter(t *testing.T) {
	const iters = 1000
	m := map[int]int{
		0: 0,
		1: 1,
		2: 2,
		3: 3,
		4: 4,
		5: 5,
		6: 6,
		7: 7,
		8: 8,
		9: 9,
	}
	counts := make([][]int, len(m))
	for i := range counts {
		counts[i] = make([]int, len(m))
	}
	var k, v int
	for i := 0; i < iters; i++ {
		it := FastIter(m, &k, &v)
		for j := 0; it.Next(); j++ {
			// key k appeared at index j
			counts[k][j]++
		}
	}

	// each key should have appeared at each index about iters/len(m) times
	for k, cs := range counts {
		for i, c := range cs {
			if (iters/len(m))/2 > c || c > (iters/len(m))*2 {
				t.Errorf("suspicious count for key %v index %v: expected %v-%v, got %v", k, i, (iters/len(m))/2, (iters/len(m))*2, c)
			}
		}
	}
}

func TestIterBadType(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when passing wrong type to Iter")
		}
	}()
	_ = Iter(make(map[int]int), new(uint8), new(uint8))
}

func BenchmarkIter(b *testing.B) {
	m := make(map[int]int, 1000)
	for i := 0; i < 1000; i++ {
		m[i] = i
	}

	b.Run("iter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var k, v int
			it := FastIter(m, &k, &v)
			for it.Next() {
			}
		}
	})

	b.Run("fastiter", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			var k, v int
			it := FastIter(m, &k, &v)
			for it.Next() {
			}
		}
	})

	b.Run("flatten", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			flat := make([]int, 0, len(m))
			for n := range m {
				flat = append(flat, n)
			}
			for _, k := range rand.Perm(len(flat)) {
				_ = k
			}
		}
	})
}

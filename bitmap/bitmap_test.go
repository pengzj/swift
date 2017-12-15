package bitmap

import "testing"

var (
	bt *Bitmap = NewBitmap()
)

func TestBitmap_Get(t *testing.T) {
	bt.Set(100)
	bt.Set(88)
	bt.Set(0)
	t.Log(bt.All())
	t.Log("bt.Get(88) = ", bt.Get(88))
	t.Log("bt.Get(100) = ", bt.Get(0))
	if bt.Get(88) == 0 {
		t.Fatal("100 not in bitmap")
	}
	bt.Clear(1000)
	bt.Clear(88)
	bt.Clear(99)
	t.Log(bt.All())
}


func BenchmarkBitmap_All(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bt.Set(100)
		bt.Set(88)
		bt.Get(3000)
	}
}

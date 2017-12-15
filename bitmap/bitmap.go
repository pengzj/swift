package bitmap

const (
	SHIFT = 5
	MASK = 0x1F
)
type Bitmap struct {
	db []int
}

func NewBitmap() *Bitmap {
	return &Bitmap{db: make([]int, 0, 50)}
}

func (b *Bitmap) Set(i int) {
	offset := int(i >>SHIFT + 1) - len(b.db)
	if offset > 0 {
		b.db = append(b.db, make([]int, offset)...)
	}

	tmp := uint(i)
	b.db[tmp>>SHIFT] |= (1 << (tmp & MASK))
}

func (b *Bitmap) Clear(i int)  {
	offset := int(i >>SHIFT + 1) - len(b.db)
	if offset > 0 {
		return
	}
	tmp := uint(i)
	b.db[tmp>>SHIFT] &= ^(1<<(tmp & MASK));
}

func (b *Bitmap) Get(i int) int {
	offset := int(i >>SHIFT + 1) - len(b.db)
	if offset > 0 {
		return 0
	}
	tmp := uint(i)
	return b.db[tmp>>SHIFT] & (1<<(tmp&MASK))
}

func (b *Bitmap) All() []int {
	return b.db
}


package bitmap

const (
	SHIFT = 5
	MASK = 0x1F
)
type Bitmap struct {
	db []uint
}

func NewBitmap() *Bitmap {
	return &Bitmap{db: make([]uint, 0, 50)}
}

func (b *Bitmap) Set(i uint) {
	offset := int(i >>SHIFT + 1) - len(b.db)
	if offset > 0 {
		b.db = append(b.db, make([]uint, offset)...)
	}

	b.db[i>>SHIFT] |= (1 << (i & MASK))
}

func (b *Bitmap) Clear(i uint)  {
	offset := int(i >>SHIFT + 1) - len(b.db)
	if offset > 0 {
		return
	}

	b.db[i>>SHIFT] &= ^(1<<(i & MASK));
}

func (b *Bitmap) Get(i uint) uint {
	offset := int(i >>SHIFT + 1) - len(b.db)
	if offset > 0 {
		return 0
	}
	return b.db[i>>SHIFT] & (1<<(i&MASK))
}

func (b *Bitmap) All() []uint {
	return b.db
}


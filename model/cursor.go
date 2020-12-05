package model

// Cursor ..
type Cursor struct {
	size int64
	page int64
	sort string
}

// NewCursor ..
func NewCursor(size, page int64, sort string) Cursor {
	return Cursor{size: size, page: page, sort: sort}
}

// GetPage ..
func (c Cursor) GetPage() int64 {
	return c.page
}

// GetSort ..
func (c Cursor) GetSort() string {
	return c.sort
}

// GetSize ..
func (c Cursor) GetSize() int64 {
	if c.size < 1 {
		return 10
	}

	if c.size > 25 {
		return 25
	}

	return c.size
}

// GetOffset ..
func (c Cursor) GetOffset() int64 {
	return (c.page - 1) * c.GetSize()
}

// GetTotalPage ..
func (c Cursor) GetTotalPage(count int64) int64 {
	totalPage := count / c.GetSize()
	remainder := count % c.GetSize()
	if remainder != 0 {
		totalPage++
	}

	return totalPage
}

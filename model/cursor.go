package model

// Sort ..
type Sort string

// String string value of Sort
func (s Sort) String() string {
	return string(s)
}

// NewSorter parse sort and match it to the sorts constant
// if none is match default to SortCreatedAtDesc
func NewSorter(sort string) Sort {
	switch sort {
	case string(SortCreatedAtAsc):
		return SortCreatedAtAsc
	default:
		return SortCreatedAtDesc
	}
}

// sorts constant
const (
	SortCreatedAtDesc = Sort("CREATED_AT_DESC")
	SortCreatedAtAsc  = Sort("CREATED_AT_ASC")
)

// Cursor ..
type Cursor struct {
	size int64
	page int64
	sort Sort
}

// NewCursor ..
func NewCursor(size, page int64, sort Sort) Cursor {
	return Cursor{size: size, page: page, sort: sort}
}

// GetPage ..
func (c Cursor) GetPage() int64 {
	if c.page < 1 {
		return 1
	}

	return c.page
}

// GetSort ..
func (c Cursor) GetSort() Sort {
	return c.sort
}

func (c *Cursor) SetPage(i int64) {
	c.page = i
}

// GetSize return the current size of cursor range from 1 to 25
// if size is out of range it will default to 10 or 25
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

// GetTotalPage total page is count divied by size
// when there is reminder, it will be increased by one
func (c Cursor) GetTotalPage(count int64) int64 {
	totalPage := count / c.GetSize()
	remainder := count % c.GetSize()
	if remainder != 0 {
		totalPage++
	}

	return totalPage
}

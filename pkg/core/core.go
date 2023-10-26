package core

import "gorm.io/gorm"

type Ctx struct {
	GormDB *gorm.DB
}

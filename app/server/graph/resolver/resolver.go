//go:generate go run github.com/99designs/gqlgen

package resolver

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Resolver struct {
	DB *gorm.DB
}

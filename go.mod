module github.com/reza-neyrami/go-passport

go 1.21.6

require gorm.io/gorm v1.25.6

require github.com/mattn/go-sqlite3 v1.14.17 // indirect

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/google/uuid v1.6.0
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmoiron/sqlx v1.3.5
	golang.org/x/crypto v0.18.0
	gorm.io/driver/sqlite v1.5.4
)

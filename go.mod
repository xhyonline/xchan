module github.com/xhyonline/xchan

go 1.13

require (
	github.com/astaxie/beego v1.12.3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.4.1 // indirect
	github.com/jinzhu/gorm v1.9.16
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/qiniu/api.v7/v7 v7.8.2
	github.com/ugorji/go v1.2.4 // indirect
	github.com/xhyonline/xutil v0.0.0-20210621081203-3c08ec0b5c93
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4

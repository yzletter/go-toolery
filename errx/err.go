package errx

import "errors"

var (
	ErrInvalidParam = errors.New("PKCS7 : UnPadding 传入非法参数")
)

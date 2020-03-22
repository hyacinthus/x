package xerr

// IsErrDuplicate 是否为本模块的 ErrDuplicate 错误
func IsErrDuplicate(err error) bool {
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	if e.code == 400 && e.Key == "Duplicate" {
		return true
	}
	return false
}

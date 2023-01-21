package cfg_utils

type Config interface {
	Validate() error
}

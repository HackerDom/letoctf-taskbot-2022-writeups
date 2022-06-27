package memeCache

import "fmt"

type RedisFetchKeyError struct {
	err error
}

func (e *RedisFetchKeyError) Error() string {
	return fmt.Sprintf("unable to fetch key from redis: %s", e.Error())
}

func (e *RedisFetchKeyError) UserError() string {
	return "Невозможно получить ключ из Redis cache"
}

//////////////////////////////////////////////////////////////////////////////////////

type RedisInvalidPictureByKey struct {
	err error
}

func (e *RedisInvalidPictureByKey) Error() string {
	return fmt.Sprintf("unable to parse image from redis cache: %s", e.err.Error())
}

func (e *RedisInvalidPictureByKey) UserError() string {
	return "В redis cache по заданному ключу находилась не картинка"
}

//////////////////////////////////////////////////////////////////////////////////////

type RedisPutKeyError struct {
	err error
}

func (e *RedisPutKeyError) Error() string {
	return fmt.Sprintf("cannot put value to redis: %s", e.err.Error())
}

func (e *RedisPutKeyError) UserError() string {
	return "Невозможно обновить значение в redis cache по ключу"
}

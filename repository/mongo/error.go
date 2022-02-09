package mongo

import (
	"fmt"

	"channels-instagram-dm/domain"
)

func newErrorNotFound(collection, value string) error {
	return domain.NewErrorNotFound(fmt.Sprintf("Mongo repository %s: value [%v] not found", collection, value))
}

func newErrorInvalidValue(collection, value string, err error) error {
	return domain.NewErrorInvalidArgument(fmt.Sprintf("Mongo repository %s: value [%v] invalid. %v", collection, value, err))
}

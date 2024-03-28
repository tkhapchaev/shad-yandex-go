//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	for {
	loopGet:
		previousValue, err := c.Get(&kvapi.GetRequest{Key: key})

		var nextValue string
		var apiError *kvapi.APIError
		var authError *kvapi.AuthError
	loopKeyNotFound:
		if errors.As(err, &authError) {
			return &kvapi.APIError{Method: "get", Err: authError}
		}

		if errors.Is(err, kvapi.ErrKeyNotFound) || (errors.Is(err, errors.Unwrap(kvapi.ErrKeyNotFound)) && errors.Unwrap(kvapi.ErrKeyNotFound) != nil) {
			nextValue, err = updateFn(nil)
		} else if errors.As(err, &apiError) {
			goto loopGet
		} else {
			nextValue, err = updateFn(&previousValue.Value)
		}

		if err != nil {
			return err
		}

		if errors.As(err, &apiError) {
			goto loopGet
		}

		version := uuid.Must(uuid.NewV4())
	loopSet:
		if previousValue != nil {
			_, err = c.Set(&kvapi.SetRequest{
				Key:        key,
				Value:      nextValue,
				OldVersion: previousValue.Version,
				NewVersion: version,
			})
		} else {
			_, err = c.Set(&kvapi.SetRequest{
				Key:        key,
				Value:      nextValue,
				OldVersion: uuid.UUID{},
				NewVersion: version,
			})
		}
		if errors.As(err, &authError) {
			return &kvapi.APIError{
				Method: "set",
				Err:    authError,
			}
		}

		var conflictError *kvapi.ConflictError

		if errors.As(err, &conflictError) {
			if version == conflictError.ExpectedVersion {
				return nil
			}

			previousValue.Version = conflictError.ExpectedVersion

			goto loopGet
		}

		if errors.Is(err, kvapi.ErrKeyNotFound) || (errors.Is(err, errors.Unwrap(kvapi.ErrKeyNotFound)) && errors.Unwrap(kvapi.ErrKeyNotFound) != nil) {
			if previousValue == nil {
				previousValue = &kvapi.GetResponse{
					Value:   "",
					Version: uuid.UUID{},
				}
			}

			previousValue.Value = nextValue
			previousValue.Version = uuid.UUID{}

			goto loopKeyNotFound
		}

		if errors.As(err, &apiError) {
			goto loopSet
		}

		return nil
	}
}

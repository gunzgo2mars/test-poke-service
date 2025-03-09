package utils

import "github.com/mitchellh/mapstructure"

type UserContext struct {
	UUID string
}

func GetUserContext(src any) (*UserContext, error) {
	var output UserContext

	if err := mapstructure.Decode(src, &output); err != nil {
		return nil, err
	}

	return &output, nil
}

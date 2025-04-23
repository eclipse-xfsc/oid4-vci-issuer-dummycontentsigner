package issuance

import "errors"

type IssuanceStorage interface {
	GetCredential(code string) (map[string]interface{}, error)
	AddCredential(code string, credential map[string]interface{}) error
}

type DummyStorage struct {
}

var store = make(map[string]map[string]interface{})

func (dummy *DummyStorage) GetCredential(code string) (map[string]interface{}, error) {

	_, ok := store[code]

	if !ok {
		return nil, errors.New("no item found for " + code)
	}

	return store[code], nil
}

func (dummy *DummyStorage) AddCredential(code string, credential map[string]interface{}) error {

	store[code] = credential

	return nil
}

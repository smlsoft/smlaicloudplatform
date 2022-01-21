package utils

import "io/ioutil"

func LoadKey(signKeyPath string, verifyKeyPath string) ([]byte, []byte, error) {

	signBytes, err := ioutil.ReadFile(signKeyPath)

	if err != nil {
		return nil, nil, err
	}

	verifyBytes, err := ioutil.ReadFile(verifyKeyPath)

	if err != nil {
		return nil, nil, err
	}

	return signBytes, verifyBytes, err

}

package storage

import "errors"

type DbConfig struct {
	Username   string
	Password   string
	AppName    string
	ClusterURL string
}

func (c DbConfig) verifyMongoValidity() error {
	if c.Username == "" || c.Password == "" {
		return errors.New("database credentials are not set")
	}

	if c.ClusterURL == "" {
		return errors.New("database url is not set")
	}

	if c.AppName == "" {
		return errors.New("database appname is not set")
	}

	return nil
}

package data

import "github.com/stretchr/testify/mock"

type MockConnection struct {
	mock.Mock
}

func (c *MockConnection) IsConnected() (bool, error) {
	return true, nil
}

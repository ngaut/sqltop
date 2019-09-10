package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StatTestSuite struct {
	suite.Suite
}

func (suite *StatTestSuite) SetupTest() {
	// FIXME: hard-coded test env
	cfg = &Conf{
		Host:             "127.0.01",
		Port:             4000,
		DBUser:           "root",
		DBPwd:            "",
		NumProcessToShow: 10,
	}

	InitDB()
}

func (suite *StatTestSuite) TestRefreshProcessList() {
	err := refreshProcessList()
	suite.Nil(err)

	processList, ok := Stat().Load(PROCESS_LIST)
	suite.Equal(true, ok)
	suite.Greater(len(processList.([]ProcessRecord)), 0)

	usingDBs, ok := Stat().Load(USING_DBS)
	suite.Equal(true, ok)
	suite.Greater(usingDBs.(int), 0)
}

func TestSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(StatTestSuite))
}

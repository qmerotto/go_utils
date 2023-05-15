package saga

import (
	"fmt"
	"go_utils/saga/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"testing"
)

type SagaSuite struct {
	suite.Suite
	mockHandlers []mocks.StepHandler
	saga         *saga
}

func (s *SagaSuite) SetupSuite() {
	s.initNewTest()
}

func (s *SagaSuite) initNewTest() {
	s.mockHandlers = []mocks.StepHandler{{}, {}, {}, {}}

	s.saga = New("test_saga").WithLogger("service_name", "test")
	s.saga.AddStep(Step{
		"step1",
		nil,
		&s.mockHandlers[0],
	})
	s.saga.AddStep(Step{
		"step2",
		[]string{"step1"},
		&s.mockHandlers[1],
	})
	s.saga.AddStep(Step{
		"step3",
		[]string{"step1", "step2"},
		&s.mockHandlers[2],
	})
	s.saga.AddStep(Step{
		"step4",
		[]string{"step3"},
		&s.mockHandlers[3],
	})
}

func TestSagaSuite(t *testing.T) {
	suite.Run(t, new(SagaSuite))
}

func (s *SagaSuite) TestHappyPath() {
	//Each step (except the first) gets one or more integer as input parameters and returns their sum + 1

	//Step1 should be called without any input
	s.mockHandlers[0].On("Execute").Return(
		func(in ...interface{}) interface{} {
			return 5
		},
		func(in ...interface{}) error {
			return nil
		},
	).Once()

	//Step 1 should be called with result from step 0
	s.mockHandlers[1].On("Execute", 5).Return(
		func(in ...interface{}) interface{} {
			return in[0].(int) + 1
		},
		func(in ...interface{}) error {
			return nil
		},
	).Once()

	//Step 2 should be called with result from step 0 and step 1
	s.mockHandlers[2].On("Execute", 5, 6).Return(
		func(in ...interface{}) interface{} {
			return in[0].(int) + in[1].(int) + 1
		},
		func(in ...interface{}) error {
			return nil
		},
	).Once()

	//Step 3 should be called with result from step 2
	s.mockHandlers[3].On("Execute", 12).Return(
		func(in ...interface{}) interface{} {
			return in[0].(int) + 1
		},
		func(in ...interface{}) error {
			return nil
		},
	).Once()

	res, err := s.saga.Exec()

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 5, res["step1"])

	//Result from step 0 + 1
	assert.Equal(s.T(), 6, res["step2"])

	//Result from step 0 + step 1 + 1
	assert.Equal(s.T(), 12, res["step3"])

	//Result from step 2 + 1
	assert.Equal(s.T(), 13, res["step4"])

}

func (s *SagaSuite) TestWithLastStepErrorShouldSuccessfullyRollback() {
	s.mockHandlers[0].On("Execute").Return(
		func(in ...interface{}) interface{} {
			return 1
		},
		func(in ...interface{}) error {
			return nil
		},
	).Once()

	s.mockHandlers[1].On("Execute", 1).Return(
		func(in ...interface{}) interface{} {
			return in[0].(int) + 1
		},
		func(in ...interface{}) error {
			return nil
		},
	).Once()

	s.mockHandlers[2].On("Execute", 2).Return(
		func(in ...interface{}) interface{} {
			return nil
		},
		func(in ...interface{}) error {
			return fmt.Errorf("BOOM")
		},
	).Once()

	s.mockHandlers[1].On("Rollback").Return(
		func() error {
			return nil
		},
	).Once()

	s.mockHandlers[0].On("Rollback").Return(
		func() error {
			return nil
		},
	).Once()

	res, err := s.saga.Exec()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), res)
}

func (s *SagaSuite) TestWithRollbackError() {
	s.mockHandlers[0].On("Execute").Return(
		func(in ...interface{}) interface{} {
			return 1
		},
		func(in ...interface{}) error {
			return nil
		},
	).Once()

	s.mockHandlers[1].On("Execute", 1).Return(
		func(in ...interface{}) interface{} {
			return in[0].(int) + 1
		},
		func(in ...interface{}) error {
			return nil
		},
	).Once()

	s.mockHandlers[2].On("Execute", 2).Return(
		func(in ...interface{}) interface{} {
			return nil
		},
		func(in ...interface{}) error {
			return fmt.Errorf("BOOM")
		},
	).Once()

	s.mockHandlers[1].On("Rollback").Return(
		func() error {
			return fmt.Errorf("DOUBLE BOOM")
		},
	).Once()

	s.mockHandlers[1].On("HandleRollbackError").Return(
		func() {},
	).Once()

	s.mockHandlers[0].On("Rollback").Return(
		func() error {
			return nil
		},
	).Once()

	res, err := s.saga.Exec()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), res)
}

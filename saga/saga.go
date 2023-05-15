package saga

import (
	"fmt"
	"go_utils/logger"
)

type StepHandler interface {
	Execute(...interface{}) (interface{}, error)
	Rollback() error
	HandleRollbackError()
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Step struct {
	Name         string
	Requirements []string
	UseCase      StepHandler
}

type saga struct {
	lastSuccessIndex int
	logger           Logger
	steps            []Step
	name             string
	results          map[string]interface{}
}

func New(name string) *saga {
	return &saga{name: name}
}

func (s *saga) AddStep(step Step) {
	s.steps = append(s.steps, step)
}

func (s *saga) WithLogger(serviceName string, environment string) *saga {
	s.logger = logger.New(fmt.Sprintf("%s - %s", serviceName, s.name))

	return s
}

func (s *saga) init() {
	s.lastSuccessIndex = -1
	s.results = make(map[string]interface{}, len(s.steps))
}

func (s *saga) Exec() (results map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			s.rollback(s.lastSuccessIndex)

			results = nil
			err = fmt.Errorf(fmt.Sprintf("error during saga %s execution, error: %v", s.name, r))
			s.logger.Error(err.Error())

			return
		}

		results = s.results
		s.logger.Info(fmt.Sprintf("successfully processed saga %s", s.name))
	}()

	s.init()
	s.logger.Info(fmt.Sprintf("begin to process saga %s", s.name))

	for _, step := range s.steps {
		s.process(step, s.buildStepInput(step))
	}

	return s.results, err
}

//buildStepInput retrieves previous results that are required for the current step
func (s *saga) buildStepInput(step Step) []interface{} {
	if len(step.Requirements) == 0 {
		return nil
	}

	stepInput := make([]interface{}, len(step.Requirements))
	for i, neededStep := range step.Requirements {
		stepInput[i] = s.results[neededStep]
	}

	return stepInput
}

func (s *saga) process(step Step, input []interface{}) {
	s.logger.Info(fmt.Sprintf("begin to process step %s", step.Name))

	result, err := step.UseCase.Execute(input...)
	if err != nil {
		s.logger.Error(fmt.Sprintf("error during step %s, error: %s", step.Name, err.Error()))
		panic(err)
	}

	s.results[step.Name] = result
	s.lastSuccessIndex++

	s.logger.Info(fmt.Sprintf("successfully processed step %s", step.Name))
}

func (s *saga) rollback(from int) {
	s.logger.Info(fmt.Sprintf("begin to rollback saga %s", s.name))

	for i := from; i >= 0; i-- {
		step := s.steps[i]
		s.logger.Info(fmt.Sprintf("begin to rollback step %s", step.Name))

		if err := step.UseCase.Rollback(); err != nil {
			s.logger.Error(
				fmt.Sprintf("error during rolling back step %s, error: %s", step.Name, err.Error()),
			)
			step.UseCase.HandleRollbackError()
			continue
		}

		s.logger.Info(fmt.Sprintf("successfully rolled back step %s", step.Name))
	}

	s.logger.Info(fmt.Sprintf("successfully rolled back saga %s", s.name))
}

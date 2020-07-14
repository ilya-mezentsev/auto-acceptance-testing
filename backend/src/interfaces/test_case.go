package interfaces

import "models"

type (
	TestCaseFactory interface {
		GetAll(testCases string) ([]TestCaseRunner, error)
	}

	TestCaseRunner interface {
		Run(processing models.TestsRun)
	}

	Transaction interface {
		Execute(context TestCaseContext)
	}

	TestCaseContext interface {
		SetVariable(name string, value interface{})
		GetVariable(name string) interface{}
		GetProcessingChannels() models.TestsRun
	}
)
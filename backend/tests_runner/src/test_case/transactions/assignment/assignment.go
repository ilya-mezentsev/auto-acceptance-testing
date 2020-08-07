package assignment

import (
	"interfaces"
	"models"
	"test_case/errors"
)

const (
	unableToBuildCommand = "Unable to build assignment command"
	unableToRunCommand   = "Unable to run assignment command"
)

type Transaction struct {
	commandBuilder interfaces.CommandBuilder
	data           interfaces.AssignmentTransactionData
}

func New(
	commandBuilder interfaces.CommandBuilder,
	data interfaces.AssignmentTransactionData,
) interfaces.Transaction {
	return Transaction{commandBuilder, data}
}

func (t Transaction) Execute(context interfaces.TestCaseContext) models.TransactionError {
	command, err := t.commandBuilder.Build(t.data.GetObject(), t.data.GetCommand())
	if err != nil {
		return models.TransactionError{
			Code:            err.Error(),
			Description:     unableToBuildCommand,
			TransactionText: t.data.GetTransactionText(),
		}
	}

	result, err := command.Run(t.data.GetArguments())
	if err != nil {
		return models.TransactionError{
			Code:            err.Error(),
			Description:     unableToRunCommand,
			TransactionText: t.data.GetTransactionText(),
		}
	}

	context.SetVariable(t.data.GetVariableName(), result)
	return errors.EmptyTransactionError
}

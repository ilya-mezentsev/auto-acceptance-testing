package timeout

import (
	"api_meta/mock/services"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"services/errors"
	"services/plugins/response_factory"
	"test_utils"
	"testing"
)

var (
	testCommandsRepositoryMock = services.TestCommandsRepositoryMock{}
	s                          = New(&testCommandsRepositoryMock)
	expectedSuccessStatus      = response_factory.DefaultResponse().GetStatus()
	expectedErrorStatus        = response_factory.ErrorResponse(nil).GetStatus()
)

func init() {
	testCommandsRepositoryMock.Reset()
}

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestService_UpdateSuccess(t *testing.T) {
	defer testCommandsRepositoryMock.Reset()

	response := s.Update(services.PredefinedAccountHash, test_utils.GetReadCloser(
		fmt.Sprintf(`{
			"timeout": 15,
			"command_hashes": [{"hash": "%s"}, {"hash": "%s"}]
		}`,
			services.PredefinedCommandHash1,
			services.PredefinedCommandHash2),
	))

	test_utils.AssertEqual(expectedSuccessStatus, response.GetStatus(), t)
	test_utils.AssertFalse(response.HasData(), t)
	test_utils.AssertNil(response.GetData(), t)

	for _, command := range testCommandsRepositoryMock.Commands[services.PredefinedAccountHash] {
		test_utils.AssertEqual(15, command.Timeout, t)
	}
}

func TestService_UpdateDecodeBodyError(t *testing.T) {
	defer testCommandsRepositoryMock.Reset()

	response := s.Update(``, test_utils.GetReadCloser(`1`))

	test_utils.AssertEqual(expectedErrorStatus, response.GetStatus(), t)
	test_utils.AssertTrue(response.HasData(), t)
	test_utils.AssertEqual(
		unableToMassUpdateTimeout,
		response.GetData().(errors.ServiceError).Code,
		t,
	)
	test_utils.AssertEqual(
		errors.DecodingRequestError,
		response.GetData().(errors.ServiceError).Description,
		t,
	)
}

func TestService_UpdateInvalidRequestError(t *testing.T) {
	defer testCommandsRepositoryMock.Reset()

	response := s.Update("blah-blah", test_utils.GetReadCloser(
		fmt.Sprintf(`{
			"timeout": 100,
			"command_hashes": [{"hash": "%s"}, {"hash": "%s"}]
		}`,
			services.PredefinedCommandHash1,
			services.PredefinedCommandHash2),
	))

	test_utils.AssertEqual(expectedErrorStatus, response.GetStatus(), t)
	test_utils.AssertTrue(response.HasData(), t)
	test_utils.AssertEqual(
		unableToMassUpdateTimeout,
		response.GetData().(errors.ServiceError).Code,
		t,
	)
	test_utils.AssertEqual(
		errors.InvalidRequestError,
		response.GetData().(errors.ServiceError).Description,
		t,
	)
}

func TestService_UpdateRepositoryError(t *testing.T) {
	defer testCommandsRepositoryMock.Reset()

	response := s.Update(services.BadAccountHash, test_utils.GetReadCloser(
		fmt.Sprintf(`{
			"timeout": 5,
			"command_hashes": [{"hash": "%s"}, {"hash": "%s"}]
		}`,
			services.PredefinedCommandHash1,
			services.PredefinedCommandHash2),
	))

	test_utils.AssertEqual(expectedErrorStatus, response.GetStatus(), t)
	test_utils.AssertTrue(response.HasData(), t)
	test_utils.AssertEqual(
		unableToMassUpdateTimeout,
		response.GetData().(errors.ServiceError).Code,
		t,
	)
	test_utils.AssertEqual(
		errors.RepositoryError,
		response.GetData().(errors.ServiceError).Description,
		t,
	)
}

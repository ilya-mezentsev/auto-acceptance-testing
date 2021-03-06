package command

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	mockCommand "test_runner_meta/mock/command"
	"test_utils"
	"testing"
)

var r = mux.NewRouter()

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestCommand_RunGetAllUsers(t *testing.T) {
	mockCommand.Init(r)
	server := test_utils.GetTestServer(r)
	defer server.Close()

	res, err := New(mockCommand.MockSettings{
		Method:             http.MethodGet,
		BaseURL:            server.URL,
		Endpoint:           "users",
		PassArgumentsInURL: false,
	}).Run(``)

	test_utils.AssertNil(err, t)
	test_utils.AssertEqual(mockCommand.StatusOk, res["status"], t)
	for hash, user := range res["data"].(map[string]interface{}) {
		expectedUser := mockCommand.Users[hash]

		test_utils.AssertEqual(expectedUser.Hash, user.(map[string]interface{})["hash"], t)
		test_utils.AssertEqual(expectedUser.Name, user.(map[string]interface{})["name"], t)
	}
}

func TestCommand_RunGetParticularUser(t *testing.T) {
	mockCommand.Init(r)
	server := test_utils.GetTestServer(r)
	defer server.Close()
	expectedUser := mockCommand.Users["hash-1"]

	res, err := New(mockCommand.MockSettings{
		Method:             http.MethodGet,
		BaseURL:            server.URL,
		Endpoint:           "user",
		PassArgumentsInURL: true,
	}).Run(`hash-1`)
	user := res["data"].(map[string]interface{})

	test_utils.AssertNil(err, t)
	test_utils.AssertEqual(mockCommand.StatusOk, res["status"], t)
	test_utils.AssertEqual(expectedUser.Hash, user["hash"], t)
	test_utils.AssertEqual(expectedUser.Name, user["name"], t)
}

func TestCommand_RunCheckHeaders(t *testing.T) {
	mockCommand.Init(r)
	mockCommand.Storage.SetHash(t.Name())
	server := test_utils.GetTestServer(r)
	defer server.Close()
	expectedHeaders := map[string]string{
		"X-Test1":      "test1",
		"X-Test2":      "test2",
		"Content-Type": "application/json; charset=utf-8",
	}

	_, err := New(mockCommand.MockSettings{
		Method:             http.MethodGet,
		BaseURL:            server.URL,
		Endpoint:           "user",
		PassArgumentsInURL: true,
		Headers:            expectedHeaders,
	}).Run(`hash-1`)

	test_utils.AssertNil(err, t)
	for key, value := range expectedHeaders {
		test_utils.AssertEqual(value, mockCommand.Storage.Get(t.Name()).Header.Get(key), t)
	}
}

func TestCommand_RunCheckCookies(t *testing.T) {
	mockCommand.Init(r)
	mockCommand.Storage.SetHash(t.Name())
	server := test_utils.GetTestServer(r)
	defer server.Close()
	expectedCookies := []*http.Cookie{
		{
			Name:  "Some-Cookie",
			Value: "Some-Cookie-Value",
		},
	}

	_, err := New(mockCommand.MockSettings{
		Method:             http.MethodGet,
		BaseURL:            server.URL,
		Endpoint:           "user",
		PassArgumentsInURL: true,
		Cookies:            expectedCookies,
	}).Run(`hash-1`)
	cookie := mockCommand.Storage.Get(t.Name()).Cookies[0]

	test_utils.AssertNil(err, t)
	test_utils.AssertEqual(expectedCookies[0].Name, cookie.Name, t)
	test_utils.AssertEqual(expectedCookies[0].Value, cookie.Value, t)
}

func TestCommand_RunBadMethod(t *testing.T) {
	_, err := New(mockCommand.MockSettings{
		Method: "//",
	}).Run(`hash-1`)

	test_utils.AssertNotNil(err, t)
	test_utils.AssertEqual(`net/http: invalid method "//"`, err.Error(), t)
}

func TestCommand_RunBadURL(t *testing.T) {
	_, err := New(mockCommand.MockSettings{
		Method: http.MethodGet,
	}).Run(`hash-1`)

	test_utils.AssertNotNil(err, t)
	test_utils.AssertEqual(`Get "/": unsupported protocol scheme ""`, err.Error(), t)
}

func TestCommand_RunStatus404(t *testing.T) {
	mockCommand.Init(r)
	server := test_utils.GetTestServer(r)
	defer server.Close()

	_, err := New(mockCommand.MockSettings{
		Method:             http.MethodGet,
		BaseURL:            server.URL,
		Endpoint:           "blah-blah",
		PassArgumentsInURL: true,
	}).Run(`hash-1`)

	test_utils.AssertNotNil(err, t)
	test_utils.AssertEqual("unsuccessful-status: 404", err.Error(), t)
}

func TestCommand_RunInvalidArguments(t *testing.T) {
	mockCommand.Init(r)
	server := test_utils.GetTestServer(r)
	defer server.Close()

	_, err := New(mockCommand.MockSettings{
		Method:             http.MethodGet,
		BaseURL:            server.URL,
		Endpoint:           "user",
		PassArgumentsInURL: true,
	}).Run(`{"x": [1, 2}`)

	test_utils.AssertNotNil(err, t)
}

func TestCommand_RunInvalidResponse(t *testing.T) {
	mockCommand.Init(r)
	server := test_utils.GetTestServer(r)
	defer server.Close()

	_, err := New(mockCommand.MockSettings{
		Method:             http.MethodGet,
		BaseURL:            server.URL,
		Endpoint:           "invalid-response",
		PassArgumentsInURL: true,
	}).Run(``)

	test_utils.AssertNotNil(err, t)
}

func TestCommand_DecodeResponseBodyBadReadCloser(t *testing.T) {
	command := New(mockCommand.MockSettings{})

	_, err := command.(*Command).decodeResponseBody(mockCommand.BadReadCLoser{})

	test_utils.AssertNotNil(err, t)
}

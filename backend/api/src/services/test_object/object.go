package test_object

import (
	"api_meta/interfaces"
	"api_meta/models"
	"api_meta/types"
	"errors"
	"io"
	servicesErrors "services/errors"
	"services/plugins/hash"
	"services/plugins/logger"
	"services/plugins/request_decoder"
	"services/plugins/response_factory"
	"services/plugins/validation"
)

type service struct {
	logger     logger.CRUDEntityErrorsLogger
	repository interfaces.CRUDRepository
}

func New(repository interfaces.CRUDRepository) interfaces.CRUDService {
	return service{
		repository: repository,
		logger:     logger.CRUDEntityErrorsLogger{EntityName: entityName},
	}
}

func (s service) Create(accountHash string, request io.ReadCloser) interfaces.Response {
	var createTestObjectRequest models.CreateTestObjectRequest
	err := request_decoder.Decode(request, &createTestObjectRequest)
	if err != nil {
		s.logger.LogCreateEntityDecodeError(err)

		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToCreateTestObjectCode,
			Description: servicesErrors.DecodingRequestError,
		})
	}

	createTestObjectRequest.TestObject.Hash = hash.Md5WithTimeAsKey(
		createTestObjectRequest.TestObject.Name,
	)
	if !validation.IsMd5Hash(accountHash) || !validation.IsValid(&createTestObjectRequest) {
		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToCreateTestObjectCode,
			Description: servicesErrors.InvalidRequestError,
		})
	}

	err = s.repository.Create(accountHash, map[string]interface{}{
		"name": createTestObjectRequest.TestObject.Name,
		"hash": createTestObjectRequest.TestObject.Hash,
	})
	if errors.As(err, &types.UniqueEntityAlreadyExists{}) {
		s.logger.LogCreateEntityUniqueConstraintError(err, map[string]interface{}{
			"account_hash": accountHash,
		})

		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToCreateTestObjectCode,
			Description: servicesErrors.UniqueEntityExistsError,
		})
	} else if err != nil {
		s.logger.LogCreateEntityRepositoryError(err, map[string]interface{}{
			"account_hash": accountHash,
		})

		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToCreateTestObjectCode,
			Description: servicesErrors.RepositoryError,
		})
	}

	return response_factory.DefaultResponse()
}

func (s service) GetAll(accountHash string) interfaces.Response {
	if !validation.IsMd5Hash(accountHash) {
		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToFetchTestObjectsCode,
			Description: servicesErrors.InvalidRequestError,
		})
	}

	var testObjects []models.TestObject
	err := s.repository.GetAll(accountHash, &testObjects)
	if err != nil {
		s.logger.LogGetAllEntitiesRepositoryError(err, map[string]interface{}{
			"account_hash": accountHash,
		})

		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToFetchTestObjectsCode,
			Description: servicesErrors.RepositoryError,
		})
	}

	return response_factory.SuccessResponse(testObjects)
}

func (s service) Get(accountHash, testObjectHash string) interfaces.Response {
	if !validation.IsMd5Hash(accountHash) || !validation.IsMd5Hash(testObjectHash) {
		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToFetchTestObjectCode,
			Description: servicesErrors.InvalidRequestError,
		})
	}

	var testObject models.TestObject
	err := s.repository.Get(accountHash, testObjectHash, &testObject)
	if err != nil {
		s.logger.LogGetEntityRepositoryError(err, map[string]interface{}{
			"account_hash":     accountHash,
			"test_object_hash": testObjectHash,
		})

		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToFetchTestObjectCode,
			Description: servicesErrors.RepositoryError,
		})
	}

	return response_factory.SuccessResponse(testObject)
}

func (s service) Update(accountHash string, request io.ReadCloser) interfaces.Response {
	var updateTestObjectRequest models.UpdateRequest
	err := request_decoder.Decode(request, &updateTestObjectRequest)
	if err != nil {
		s.logger.LogUpdateEntityDecodeError(err)

		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToUpdateTestObjectCode,
			Description: servicesErrors.DecodingRequestError,
		})
	}

	newName, ok := updateTestObjectRequest.UpdatePayload[0].NewValue.(string)
	if !validation.IsMd5Hash(accountHash) || !validation.IsValid(&updateTestObjectRequest) ||
		!(ok && validation.IsRegularName(newName)) {
		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToUpdateTestObjectCode,
			Description: servicesErrors.InvalidRequestError,
		})
	}

	err = s.repository.Update(accountHash, updateTestObjectRequest.UpdatePayload)
	if err != nil {
		s.logger.LogUpdateEntityRepositoryError(err, map[string]interface{}{
			"account_hash":   accountHash,
			"update_payload": updateTestObjectRequest.UpdatePayload,
		})

		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToUpdateTestObjectCode,
			Description: servicesErrors.RepositoryError,
		})
	}

	return response_factory.DefaultResponse()
}

func (s service) Delete(accountHash, testObjectHash string) interfaces.Response {
	if !validation.IsMd5Hash(accountHash) || !validation.IsMd5Hash(testObjectHash) {
		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToDeleteTestObjectCode,
			Description: servicesErrors.InvalidRequestError,
		})
	}

	err := s.repository.Delete(accountHash, testObjectHash)
	if err != nil {
		s.logger.LogDeleteEntityRepositoryError(err, map[string]interface{}{
			"account_hash":     accountHash,
			"test_object_hash": testObjectHash,
		})

		return response_factory.ErrorResponse(servicesErrors.ServiceError{
			Code:        unableToDeleteTestObjectCode,
			Description: servicesErrors.RepositoryError,
		})
	}

	return response_factory.DefaultResponse()
}

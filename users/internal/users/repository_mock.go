// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package users

import (
	"context"
	"github.com/admarc/users/internal/models"
	"sync"
)

// Ensure, that RepositoryMock does implement Repository.
// If this is not the case, regenerate this file with moq.
var _ Repository = &RepositoryMock{}

// RepositoryMock is a mock implementation of Repository.
//
//	func TestSomethingThatUsesRepository(t *testing.T) {
//
//		// make and configure a mocked Repository
//		mockedRepository := &RepositoryMock{
//			CreateFunc: func(ctx context.Context, name string) (models.User, error) {
//				panic("mock out the Create method")
//			},
//			GetFunc: func(ctx context.Context, id string) (models.User, error) {
//				panic("mock out the Get method")
//			},
//		}
//
//		// use mockedRepository in code that requires Repository
//		// and then make assertions.
//
//	}
type RepositoryMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(ctx context.Context, name string) (models.User, error)

	// GetFunc mocks the Get method.
	GetFunc func(ctx context.Context, id string) (models.User, error)

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID string
		}
	}
	lockCreate sync.RWMutex
	lockGet    sync.RWMutex
}

// Create calls CreateFunc.
func (mock *RepositoryMock) Create(ctx context.Context, name string) (models.User, error) {
	if mock.CreateFunc == nil {
		panic("RepositoryMock.CreateFunc: method is nil but Repository.Create was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Name string
	}{
		Ctx:  ctx,
		Name: name,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(ctx, name)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockedRepository.CreateCalls())
func (mock *RepositoryMock) CreateCalls() []struct {
	Ctx  context.Context
	Name string
} {
	var calls []struct {
		Ctx  context.Context
		Name string
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *RepositoryMock) Get(ctx context.Context, id string) (models.User, error) {
	if mock.GetFunc == nil {
		panic("RepositoryMock.GetFunc: method is nil but Repository.Get was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  string
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(ctx, id)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//
//	len(mockedRepository.GetCalls())
func (mock *RepositoryMock) GetCalls() []struct {
	Ctx context.Context
	ID  string
} {
	var calls []struct {
		Ctx context.Context
		ID  string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

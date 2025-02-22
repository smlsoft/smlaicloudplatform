// Code generated by mockery v2.40.1. DO NOT EDIT.

package mocks

import (
	models "smlaicloudplatform/internal/transaction/paymentdetail/models"

	mock "github.com/stretchr/testify/mock"
)

// MockPaymentDetailRepository is an autogenerated mock type for the IPaymentDetailRepository type
type MockPaymentDetailRepository struct {
	mock.Mock
}

// Create provides a mock function with given fields: doc
func (_m *MockPaymentDetailRepository) Create(doc models.TransactionPaymentDetail) error {
	ret := _m.Called(doc)

	if len(ret) == 0 {
		panic("no return value specified for Create")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(models.TransactionPaymentDetail) error); ok {
		r0 = rf(doc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: shopID, docNo, doc
func (_m *MockPaymentDetailRepository) Delete(shopID string, docNo string, doc models.TransactionPaymentDetail) error {
	ret := _m.Called(shopID, docNo, doc)

	if len(ret) == 0 {
		panic("no return value specified for Delete")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, models.TransactionPaymentDetail) error); ok {
		r0 = rf(shopID, docNo, doc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields: shopID, docNo
func (_m *MockPaymentDetailRepository) Get(shopID string, docNo string) (*models.TransactionPaymentDetail, error) {
	ret := _m.Called(shopID, docNo)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *models.TransactionPaymentDetail
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (*models.TransactionPaymentDetail, error)); ok {
		return rf(shopID, docNo)
	}
	if rf, ok := ret.Get(0).(func(string, string) *models.TransactionPaymentDetail); ok {
		r0 = rf(shopID, docNo)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.TransactionPaymentDetail)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(shopID, docNo)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: shopID, docNo, doc
func (_m *MockPaymentDetailRepository) Update(shopID string, docNo string, doc models.TransactionPaymentDetail) error {
	ret := _m.Called(shopID, docNo, doc)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, models.TransactionPaymentDetail) error); ok {
		r0 = rf(shopID, docNo, doc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMockPaymentDetailRepository creates a new instance of MockPaymentDetailRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockPaymentDetailRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockPaymentDetailRepository {
	mock := &MockPaymentDetailRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

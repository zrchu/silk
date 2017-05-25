// This file was generated by counterfeiter
package fakes

import (
	"database/sql"
	"sync"
	"time"

	"code.cloudfoundry.org/silk/controller/database"
)

type SubnetDeleter struct {
	DeleteStub        func(database.Transaction, string, time.Duration) (sql.Result, error)
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 database.Transaction
		arg2 string
		arg3 time.Duration
	}
	deleteReturns struct {
		result1 sql.Result
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *SubnetDeleter) Delete(arg1 database.Transaction, arg2 string, arg3 time.Duration) (sql.Result, error) {
	fake.deleteMutex.Lock()
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 database.Transaction
		arg2 string
		arg3 time.Duration
	}{arg1, arg2, arg3})
	fake.recordInvocation("Delete", []interface{}{arg1, arg2, arg3})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(arg1, arg2, arg3)
	}
	return fake.deleteReturns.result1, fake.deleteReturns.result2
}

func (fake *SubnetDeleter) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *SubnetDeleter) DeleteArgsForCall(i int) (database.Transaction, string, time.Duration) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.deleteArgsForCall[i].arg1, fake.deleteArgsForCall[i].arg2, fake.deleteArgsForCall[i].arg3
}

func (fake *SubnetDeleter) DeleteReturns(result1 sql.Result, result2 error) {
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 sql.Result
		result2 error
	}{result1, result2}
}

func (fake *SubnetDeleter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return fake.invocations
}

func (fake *SubnetDeleter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ database.SubnetDeleter = new(SubnetDeleter)
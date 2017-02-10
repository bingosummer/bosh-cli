// This file was generated by counterfeiter
package jobfakes

import (
	"sync"

	"github.com/cloudfoundry/bosh-cli/release/job"
	boshman "github.com/cloudfoundry/bosh-cli/release/manifest"
)

type FakeArchiveReader struct {
	ReadStub        func(boshman.JobRef, string) (*job.Job, error)
	readMutex       sync.RWMutex
	readArgsForCall []struct {
		arg1 boshman.JobRef
		arg2 string
	}
	readReturns struct {
		result1 *job.Job
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeArchiveReader) Read(arg1 boshman.JobRef, arg2 string) (*job.Job, error) {
	fake.readMutex.Lock()
	fake.readArgsForCall = append(fake.readArgsForCall, struct {
		arg1 boshman.JobRef
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("Read", []interface{}{arg1, arg2})
	fake.readMutex.Unlock()
	if fake.ReadStub != nil {
		return fake.ReadStub(arg1, arg2)
	}
	return fake.readReturns.result1, fake.readReturns.result2
}

func (fake *FakeArchiveReader) ReadCallCount() int {
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	return len(fake.readArgsForCall)
}

func (fake *FakeArchiveReader) ReadArgsForCall(i int) (boshman.JobRef, string) {
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	return fake.readArgsForCall[i].arg1, fake.readArgsForCall[i].arg2
}

func (fake *FakeArchiveReader) ReadReturns(result1 *job.Job, result2 error) {
	fake.ReadStub = nil
	fake.readReturns = struct {
		result1 *job.Job
		result2 error
	}{result1, result2}
}

func (fake *FakeArchiveReader) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeArchiveReader) recordInvocation(key string, args []interface{}) {
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

var _ job.ArchiveReader = new(FakeArchiveReader)

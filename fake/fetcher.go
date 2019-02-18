// Code generated by counterfeiter. DO NOT EDIT.
package fake

import (
	"sync"

	"github.com/phogolabs/vault"
)

type Fetcher struct {
	SecretStub        func(path string) (interface{}, error)
	secretMutex       sync.RWMutex
	secretArgsForCall []struct {
		path string
	}
	secretReturns struct {
		result1 interface{}
		result2 error
	}
	secretReturnsOnCall map[int]struct {
		result1 interface{}
		result2 error
	}
	StopStub         func()
	stopMutex        sync.RWMutex
	stopArgsForCall  []struct{}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Fetcher) Secret(path string) (interface{}, error) {
	fake.secretMutex.Lock()
	ret, specificReturn := fake.secretReturnsOnCall[len(fake.secretArgsForCall)]
	fake.secretArgsForCall = append(fake.secretArgsForCall, struct {
		path string
	}{path})
	fake.recordInvocation("Secret", []interface{}{path})
	fake.secretMutex.Unlock()
	if fake.SecretStub != nil {
		return fake.SecretStub(path)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.secretReturns.result1, fake.secretReturns.result2
}

func (fake *Fetcher) SecretCallCount() int {
	fake.secretMutex.RLock()
	defer fake.secretMutex.RUnlock()
	return len(fake.secretArgsForCall)
}

func (fake *Fetcher) SecretArgsForCall(i int) string {
	fake.secretMutex.RLock()
	defer fake.secretMutex.RUnlock()
	return fake.secretArgsForCall[i].path
}

func (fake *Fetcher) SecretReturns(result1 interface{}, result2 error) {
	fake.SecretStub = nil
	fake.secretReturns = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *Fetcher) SecretReturnsOnCall(i int, result1 interface{}, result2 error) {
	fake.SecretStub = nil
	if fake.secretReturnsOnCall == nil {
		fake.secretReturnsOnCall = make(map[int]struct {
			result1 interface{}
			result2 error
		})
	}
	fake.secretReturnsOnCall[i] = struct {
		result1 interface{}
		result2 error
	}{result1, result2}
}

func (fake *Fetcher) Stop() {
	fake.stopMutex.Lock()
	fake.stopArgsForCall = append(fake.stopArgsForCall, struct{}{})
	fake.recordInvocation("Stop", []interface{}{})
	fake.stopMutex.Unlock()
	if fake.StopStub != nil {
		fake.StopStub()
	}
}

func (fake *Fetcher) StopCallCount() int {
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	return len(fake.stopArgsForCall)
}

func (fake *Fetcher) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.secretMutex.RLock()
	defer fake.secretMutex.RUnlock()
	fake.stopMutex.RLock()
	defer fake.stopMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Fetcher) recordInvocation(key string, args []interface{}) {
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

var _ vault.Fetcher = new(Fetcher)

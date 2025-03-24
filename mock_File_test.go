// Code generated by mockery v2.53.3. DO NOT EDIT.

package pgrx

import (
	fs "io/fs"

	mock "github.com/stretchr/testify/mock"
)

// MockFile is an autogenerated mock type for the File type
type MockFile struct {
	mock.Mock
}

type MockFile_Expecter struct {
	mock *mock.Mock
}

func (_m *MockFile) EXPECT() *MockFile_Expecter {
	return &MockFile_Expecter{mock: &_m.Mock}
}

// Close provides a mock function with no fields
func (_m *MockFile) Close() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Close")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFile_Close_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Close'
type MockFile_Close_Call struct {
	*mock.Call
}

// Close is a helper method to define mock.On call
func (_e *MockFile_Expecter) Close() *MockFile_Close_Call {
	return &MockFile_Close_Call{Call: _e.mock.On("Close")}
}

func (_c *MockFile_Close_Call) Run(run func()) *MockFile_Close_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFile_Close_Call) Return(_a0 error) *MockFile_Close_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFile_Close_Call) RunAndReturn(run func() error) *MockFile_Close_Call {
	_c.Call.Return(run)
	return _c
}

// Name provides a mock function with no fields
func (_m *MockFile) Name() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Name")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// MockFile_Name_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Name'
type MockFile_Name_Call struct {
	*mock.Call
}

// Name is a helper method to define mock.On call
func (_e *MockFile_Expecter) Name() *MockFile_Name_Call {
	return &MockFile_Name_Call{Call: _e.mock.On("Name")}
}

func (_c *MockFile_Name_Call) Run(run func()) *MockFile_Name_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFile_Name_Call) Return(_a0 string) *MockFile_Name_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFile_Name_Call) RunAndReturn(run func() string) *MockFile_Name_Call {
	_c.Call.Return(run)
	return _c
}

// Read provides a mock function with given fields: p
func (_m *MockFile) Read(p []byte) (int, error) {
	ret := _m.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for Read")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(p)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFile_Read_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Read'
type MockFile_Read_Call struct {
	*mock.Call
}

// Read is a helper method to define mock.On call
//   - p []byte
func (_e *MockFile_Expecter) Read(p interface{}) *MockFile_Read_Call {
	return &MockFile_Read_Call{Call: _e.mock.On("Read", p)}
}

func (_c *MockFile_Read_Call) Run(run func(p []byte)) *MockFile_Read_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *MockFile_Read_Call) Return(n int, err error) *MockFile_Read_Call {
	_c.Call.Return(n, err)
	return _c
}

func (_c *MockFile_Read_Call) RunAndReturn(run func([]byte) (int, error)) *MockFile_Read_Call {
	_c.Call.Return(run)
	return _c
}

// ReadAt provides a mock function with given fields: p, off
func (_m *MockFile) ReadAt(p []byte, off int64) (int, error) {
	ret := _m.Called(p, off)

	if len(ret) == 0 {
		panic("no return value specified for ReadAt")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, int64) (int, error)); ok {
		return rf(p, off)
	}
	if rf, ok := ret.Get(0).(func([]byte, int64) int); ok {
		r0 = rf(p, off)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte, int64) error); ok {
		r1 = rf(p, off)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFile_ReadAt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ReadAt'
type MockFile_ReadAt_Call struct {
	*mock.Call
}

// ReadAt is a helper method to define mock.On call
//   - p []byte
//   - off int64
func (_e *MockFile_Expecter) ReadAt(p interface{}, off interface{}) *MockFile_ReadAt_Call {
	return &MockFile_ReadAt_Call{Call: _e.mock.On("ReadAt", p, off)}
}

func (_c *MockFile_ReadAt_Call) Run(run func(p []byte, off int64)) *MockFile_ReadAt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte), args[1].(int64))
	})
	return _c
}

func (_c *MockFile_ReadAt_Call) Return(n int, err error) *MockFile_ReadAt_Call {
	_c.Call.Return(n, err)
	return _c
}

func (_c *MockFile_ReadAt_Call) RunAndReturn(run func([]byte, int64) (int, error)) *MockFile_ReadAt_Call {
	_c.Call.Return(run)
	return _c
}

// Readdir provides a mock function with given fields: count
func (_m *MockFile) Readdir(count int) ([]fs.FileInfo, error) {
	ret := _m.Called(count)

	if len(ret) == 0 {
		panic("no return value specified for Readdir")
	}

	var r0 []fs.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]fs.FileInfo, error)); ok {
		return rf(count)
	}
	if rf, ok := ret.Get(0).(func(int) []fs.FileInfo); ok {
		r0 = rf(count)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]fs.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(count)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFile_Readdir_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Readdir'
type MockFile_Readdir_Call struct {
	*mock.Call
}

// Readdir is a helper method to define mock.On call
//   - count int
func (_e *MockFile_Expecter) Readdir(count interface{}) *MockFile_Readdir_Call {
	return &MockFile_Readdir_Call{Call: _e.mock.On("Readdir", count)}
}

func (_c *MockFile_Readdir_Call) Run(run func(count int)) *MockFile_Readdir_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockFile_Readdir_Call) Return(_a0 []fs.FileInfo, _a1 error) *MockFile_Readdir_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFile_Readdir_Call) RunAndReturn(run func(int) ([]fs.FileInfo, error)) *MockFile_Readdir_Call {
	_c.Call.Return(run)
	return _c
}

// Readdirnames provides a mock function with given fields: n
func (_m *MockFile) Readdirnames(n int) ([]string, error) {
	ret := _m.Called(n)

	if len(ret) == 0 {
		panic("no return value specified for Readdirnames")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(int) ([]string, error)); ok {
		return rf(n)
	}
	if rf, ok := ret.Get(0).(func(int) []string); ok {
		r0 = rf(n)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(n)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFile_Readdirnames_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Readdirnames'
type MockFile_Readdirnames_Call struct {
	*mock.Call
}

// Readdirnames is a helper method to define mock.On call
//   - n int
func (_e *MockFile_Expecter) Readdirnames(n interface{}) *MockFile_Readdirnames_Call {
	return &MockFile_Readdirnames_Call{Call: _e.mock.On("Readdirnames", n)}
}

func (_c *MockFile_Readdirnames_Call) Run(run func(n int)) *MockFile_Readdirnames_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int))
	})
	return _c
}

func (_c *MockFile_Readdirnames_Call) Return(_a0 []string, _a1 error) *MockFile_Readdirnames_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFile_Readdirnames_Call) RunAndReturn(run func(int) ([]string, error)) *MockFile_Readdirnames_Call {
	_c.Call.Return(run)
	return _c
}

// Seek provides a mock function with given fields: offset, whence
func (_m *MockFile) Seek(offset int64, whence int) (int64, error) {
	ret := _m.Called(offset, whence)

	if len(ret) == 0 {
		panic("no return value specified for Seek")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int) (int64, error)); ok {
		return rf(offset, whence)
	}
	if rf, ok := ret.Get(0).(func(int64, int) int64); ok {
		r0 = rf(offset, whence)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(int64, int) error); ok {
		r1 = rf(offset, whence)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFile_Seek_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Seek'
type MockFile_Seek_Call struct {
	*mock.Call
}

// Seek is a helper method to define mock.On call
//   - offset int64
//   - whence int
func (_e *MockFile_Expecter) Seek(offset interface{}, whence interface{}) *MockFile_Seek_Call {
	return &MockFile_Seek_Call{Call: _e.mock.On("Seek", offset, whence)}
}

func (_c *MockFile_Seek_Call) Run(run func(offset int64, whence int)) *MockFile_Seek_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64), args[1].(int))
	})
	return _c
}

func (_c *MockFile_Seek_Call) Return(_a0 int64, _a1 error) *MockFile_Seek_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFile_Seek_Call) RunAndReturn(run func(int64, int) (int64, error)) *MockFile_Seek_Call {
	_c.Call.Return(run)
	return _c
}

// Stat provides a mock function with no fields
func (_m *MockFile) Stat() (fs.FileInfo, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Stat")
	}

	var r0 fs.FileInfo
	var r1 error
	if rf, ok := ret.Get(0).(func() (fs.FileInfo, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() fs.FileInfo); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(fs.FileInfo)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFile_Stat_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Stat'
type MockFile_Stat_Call struct {
	*mock.Call
}

// Stat is a helper method to define mock.On call
func (_e *MockFile_Expecter) Stat() *MockFile_Stat_Call {
	return &MockFile_Stat_Call{Call: _e.mock.On("Stat")}
}

func (_c *MockFile_Stat_Call) Run(run func()) *MockFile_Stat_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFile_Stat_Call) Return(_a0 fs.FileInfo, _a1 error) *MockFile_Stat_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockFile_Stat_Call) RunAndReturn(run func() (fs.FileInfo, error)) *MockFile_Stat_Call {
	_c.Call.Return(run)
	return _c
}

// Sync provides a mock function with no fields
func (_m *MockFile) Sync() error {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Sync")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFile_Sync_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Sync'
type MockFile_Sync_Call struct {
	*mock.Call
}

// Sync is a helper method to define mock.On call
func (_e *MockFile_Expecter) Sync() *MockFile_Sync_Call {
	return &MockFile_Sync_Call{Call: _e.mock.On("Sync")}
}

func (_c *MockFile_Sync_Call) Run(run func()) *MockFile_Sync_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *MockFile_Sync_Call) Return(_a0 error) *MockFile_Sync_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFile_Sync_Call) RunAndReturn(run func() error) *MockFile_Sync_Call {
	_c.Call.Return(run)
	return _c
}

// Truncate provides a mock function with given fields: size
func (_m *MockFile) Truncate(size int64) error {
	ret := _m.Called(size)

	if len(ret) == 0 {
		panic("no return value specified for Truncate")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(int64) error); ok {
		r0 = rf(size)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockFile_Truncate_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Truncate'
type MockFile_Truncate_Call struct {
	*mock.Call
}

// Truncate is a helper method to define mock.On call
//   - size int64
func (_e *MockFile_Expecter) Truncate(size interface{}) *MockFile_Truncate_Call {
	return &MockFile_Truncate_Call{Call: _e.mock.On("Truncate", size)}
}

func (_c *MockFile_Truncate_Call) Run(run func(size int64)) *MockFile_Truncate_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(int64))
	})
	return _c
}

func (_c *MockFile_Truncate_Call) Return(_a0 error) *MockFile_Truncate_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockFile_Truncate_Call) RunAndReturn(run func(int64) error) *MockFile_Truncate_Call {
	_c.Call.Return(run)
	return _c
}

// Write provides a mock function with given fields: p
func (_m *MockFile) Write(p []byte) (int, error) {
	ret := _m.Called(p)

	if len(ret) == 0 {
		panic("no return value specified for Write")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (int, error)); ok {
		return rf(p)
	}
	if rf, ok := ret.Get(0).(func([]byte) int); ok {
		r0 = rf(p)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFile_Write_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Write'
type MockFile_Write_Call struct {
	*mock.Call
}

// Write is a helper method to define mock.On call
//   - p []byte
func (_e *MockFile_Expecter) Write(p interface{}) *MockFile_Write_Call {
	return &MockFile_Write_Call{Call: _e.mock.On("Write", p)}
}

func (_c *MockFile_Write_Call) Run(run func(p []byte)) *MockFile_Write_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte))
	})
	return _c
}

func (_c *MockFile_Write_Call) Return(n int, err error) *MockFile_Write_Call {
	_c.Call.Return(n, err)
	return _c
}

func (_c *MockFile_Write_Call) RunAndReturn(run func([]byte) (int, error)) *MockFile_Write_Call {
	_c.Call.Return(run)
	return _c
}

// WriteAt provides a mock function with given fields: p, off
func (_m *MockFile) WriteAt(p []byte, off int64) (int, error) {
	ret := _m.Called(p, off)

	if len(ret) == 0 {
		panic("no return value specified for WriteAt")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, int64) (int, error)); ok {
		return rf(p, off)
	}
	if rf, ok := ret.Get(0).(func([]byte, int64) int); ok {
		r0 = rf(p, off)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func([]byte, int64) error); ok {
		r1 = rf(p, off)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFile_WriteAt_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WriteAt'
type MockFile_WriteAt_Call struct {
	*mock.Call
}

// WriteAt is a helper method to define mock.On call
//   - p []byte
//   - off int64
func (_e *MockFile_Expecter) WriteAt(p interface{}, off interface{}) *MockFile_WriteAt_Call {
	return &MockFile_WriteAt_Call{Call: _e.mock.On("WriteAt", p, off)}
}

func (_c *MockFile_WriteAt_Call) Run(run func(p []byte, off int64)) *MockFile_WriteAt_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].([]byte), args[1].(int64))
	})
	return _c
}

func (_c *MockFile_WriteAt_Call) Return(n int, err error) *MockFile_WriteAt_Call {
	_c.Call.Return(n, err)
	return _c
}

func (_c *MockFile_WriteAt_Call) RunAndReturn(run func([]byte, int64) (int, error)) *MockFile_WriteAt_Call {
	_c.Call.Return(run)
	return _c
}

// WriteString provides a mock function with given fields: s
func (_m *MockFile) WriteString(s string) (int, error) {
	ret := _m.Called(s)

	if len(ret) == 0 {
		panic("no return value specified for WriteString")
	}

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int, error)); ok {
		return rf(s)
	}
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(s)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(s)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockFile_WriteString_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'WriteString'
type MockFile_WriteString_Call struct {
	*mock.Call
}

// WriteString is a helper method to define mock.On call
//   - s string
func (_e *MockFile_Expecter) WriteString(s interface{}) *MockFile_WriteString_Call {
	return &MockFile_WriteString_Call{Call: _e.mock.On("WriteString", s)}
}

func (_c *MockFile_WriteString_Call) Run(run func(s string)) *MockFile_WriteString_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockFile_WriteString_Call) Return(ret int, err error) *MockFile_WriteString_Call {
	_c.Call.Return(ret, err)
	return _c
}

func (_c *MockFile_WriteString_Call) RunAndReturn(run func(string) (int, error)) *MockFile_WriteString_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockFile creates a new instance of MockFile. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockFile(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockFile {
	mock := &MockFile{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}

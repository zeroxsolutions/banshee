package mock_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/zeroxsolutions/banshee/mock"
)

// TestMockCache_IsConnected_IsTrue tests the MockCache's IsConnected method when the mock returns true.
func TestMockCache_IsConnected_IsTrue(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	mockCache.On("IsConnected", ctx).Return(true)

	isConnected := mockCache.IsConnected(ctx)
	if !isConnected {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_IsConnected_IsTrueWithFunction tests IsConnected using a function-based return value.
func TestMockCache_IsConnected_IsTrueWithFunction(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	mockCache.On("IsConnected", ctx).Return(func(context.Context) bool {
		return true
	})

	isConnected := mockCache.IsConnected(ctx)
	if !isConnected {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_IsConnected_IsFalse tests the MockCache's IsConnected method when the mock returns false.
func TestMockCache_IsConnected_IsFalse(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	mockCache.On("IsConnected", ctx).Return(false)

	isConnected := mockCache.IsConnected(ctx)
	if isConnected {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_Keys_Err tests the Keys method when an error is returned.
func TestMockCache_Keys_Err(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	pattern := "key*"

	r1 := errors.New("error test")

	mockCache.On("Keys", ctx, pattern).Return(nil, r1)

	r0, err := mockCache.Keys(ctx, pattern)

	if !errors.Is(err, r1) {
		t.FailNow()
	}

	if r0 != nil {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_Keys_NilErr tests the Keys method when no error is returned and keys are successfully retrieved.
func TestMockCache_Keys_NilErr(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	pattern := "key*"

	keys := []string{"key-1", "key-2"}

	mockCache.On("Keys", ctx, pattern).Return(keys, nil)

	r0, err := mockCache.Keys(ctx, pattern)

	if err != nil {
		t.FailNow()
	}

	if !reflect.DeepEqual(keys, r0) {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_Get_Err tests the Get method when an error is returned.
func TestMockCache_Get_Err(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	key := "key"

	r1 := errors.New("error test")

	mockCache.On("Get", ctx, key).Return("", r1)

	v, err := mockCache.Get(ctx, key)

	if !errors.Is(err, r1) {
		t.FailNow()
	}

	if v != "" {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_Get_NilErr tests the Get method when no error is returned and a value is successfully retrieved.
func TestMockCache_Get_NilErr(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	key := "key"
	value := "value"

	mockCache.On("Get", ctx, key).Return(value, nil)

	v, err := mockCache.Get(ctx, key)

	if err != nil {
		t.FailNow()
	}

	if v != value {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_Set_Err tests the Set method when an error is returned.
func TestMockCache_Set_Err(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	key := "key"
	value := "value"

	r0 := errors.New("error test")

	mockCache.On("Set", ctx, key, value).Return(r0)

	err := mockCache.Set(ctx, key, value)

	if !errors.Is(err, r0) {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_Set_NilErr tests the Set method when no error is returned.
func TestMockCache_Set_NilErr(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	key := "key"
	value := "value"

	mockCache.On("Set", ctx, key, value).Return(nil)

	err := mockCache.Set(ctx, key, value)

	if err != nil {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_SetWithExpiration_Err tests the SetWithExpiration method when an error is returned.
func TestMockCache_SetWithExpiration_Err(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	key := "key"
	value := "value"
	expiration := 10 * time.Second

	r0 := errors.New("error test")

	mockCache.On("SetWithExpiration", ctx, key, value, expiration).Return(r0)

	err := mockCache.SetWithExpiration(ctx, key, value, expiration)

	if !errors.Is(err, r0) {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_SetWithExpiration_NilErr tests the SetWithExpiration method when no error is returned.
func TestMockCache_SetWithExpiration_NilErr(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	key := "key"
	value := "value"
	expiration := 10 * time.Second

	mockCache.On("SetWithExpiration", ctx, key, value, expiration).Return(nil)

	err := mockCache.SetWithExpiration(ctx, key, value, expiration)

	if err != nil {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_Del_Err tests the Del method when an error is returned.
func TestMockCache_Del_Err(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	var args []interface{}
	keys := []string{"key-1", "key-2"}

	args = append(args, ctx)
	for _, key := range keys {
		args = append(args, key)
	}

	r0 := errors.New("error test")

	mockCache.On("Del", args...).Return(r0)

	err := mockCache.Del(ctx, keys...)

	if !errors.Is(err, r0) {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_Del_NilErr tests the Del method when no error is returned.
func TestMockCache_Del_NilErr(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()

	var args []interface{}
	keys := []string{"key-1", "key-2"}

	args = append(args, ctx)
	for _, key := range keys {
		args = append(args, key)
	}

	mockCache.On("Del", args...).Return(nil)

	err := mockCache.Del(ctx, keys...)

	if err != nil {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_DelWithPattern_Err tests the DelWithPattern method when an error is returned.
func TestMockCache_DelWithPattern_Err(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()
	pattern := "key*"

	r0 := errors.New("error test")

	mockCache.On("DelWithPattern", ctx, pattern).Return(r0)

	if err := mockCache.DelWithPattern(ctx, pattern); !errors.Is(err, r0) {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_DelWithPattern_NilErr tests the DelWithPattern method when no error is returned.
func TestMockCache_DelWithPattern_NilErr(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	ctx := context.Background()
	pattern := "key*"

	mockCache.On("DelWithPattern", ctx, pattern).Return(nil)

	if err := mockCache.DelWithPattern(ctx, pattern); err != nil {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_Close_Err tests the Close method when an error is returned.
func TestMockCache_Close_Err(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	r0 := errors.New("error test")

	mockCache.On("Close").Return(r0)

	if err := mockCache.Close(); !errors.Is(err, r0) {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

// TestMockCache_Close_NilErr tests the Close method when no error is returned.
func TestMockCache_Close_NilErr(t *testing.T) {
	mockCache := mock.NewMockCache(t).(*mock.MockCache)

	mockCache.On("Close").Return(nil)

	if err := mockCache.Close(); err != nil {
		t.FailNow()
	}

	mockCache.AssertExpectations(t)
}

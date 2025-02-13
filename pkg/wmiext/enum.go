//go:build windows
// +build windows

package wmiext

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

type IEnumWbemClassObjectVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr
	Reset          uintptr
	Next           uintptr
	NextAsync      uintptr
	Clone          uintptr
	Skip           uintptr
}

type Enum struct {
	enum    *ole.IUnknown
	vTable  *IEnumWbemClassObjectVtbl
	service *Service
}

func (e *Enum) Close() {
	if e != nil && e.enum != nil {
		e.enum.Release()
	}
}

func newEnum(enumerator *ole.IUnknown, service *Service) *Enum {
	return &Enum{
		enum:    enumerator,
		vTable:  (*IEnumWbemClassObjectVtbl)(unsafe.Pointer(enumerator.RawVTable)),
		service: service,
	}
}

// NextObject obtains the next instance in an enumeration and sets all fields
// of the struct pointer passed through the target parameter. Otherwise, if
// the target parameter is not a struct pointer type, an error will be
// returned.
func NextObject(enum *Enum, target interface{}) (bool, error) {
	var err error

	var instance *Instance
	if instance, err = enum.Next(); err != nil {
		return false, err
	}

	if instance == nil {
		return true, nil
	}

	//defer instance.Close()

	return false, instance.GetAll(target)
}

// NextObjects retrieves all instances in an enumeration and appends them to the provided slice pointer.
func NextObjects(enum *Enum, targetSlice interface{}) error {
	sliceValue := reflect.ValueOf(targetSlice)
	if sliceValue.Kind() != reflect.Ptr || sliceValue.Elem().Kind() != reflect.Slice {
		return errors.New("targetSlice must be a pointer to a slice")
	}

	sliceValue = sliceValue.Elem()            // 解引用到 slice 本身
	sliceElemType := sliceValue.Type().Elem() // 获取 slice 的元素类型

	// 检测 slice 元素是否是指针类型
	isPointer := sliceElemType.Kind() == reflect.Ptr
	structType := sliceElemType
	if isPointer {
		structType = sliceElemType.Elem() // 获取指针指向的 struct 类型
	}

	for {
		var instance *Instance
		var err error

		if instance, err = enum.Next(); err != nil {
			return err
		}

		if instance == nil {
			break
		}
		//defer instance.Close()

		// 创建新的 struct 实例
		newElem := reflect.New(structType).Elem()

		// 填充数据
		if err := instance.GetAll(newElem.Addr().Interface()); err != nil {
			return err
		}

		// 如果 slice 需要存储指针，则取地址
		if isPointer {
			sliceValue.Set(reflect.Append(sliceValue, newElem.Addr()))
		} else {
			sliceValue.Set(reflect.Append(sliceValue, newElem))
		}
	}

	return nil
}

// Next returns the next object instance in this iteration
func (e *Enum) Next() (instance *Instance, err error) {
	var res uintptr
	var apObjects *ole.IUnknown
	var uReturned uint32

	res, _, _ = syscall.SyscallN(
		e.vTable.Next,                       // IEnumWbemClassObject::Next()
		uintptr(unsafe.Pointer(e.enum)),     // IEnumWbemClassObject   ptr
		uintptr(WBEM_INFINITE),              // [in]  long             lTimeout,
		uintptr(1),                          // [in]  ULONG            uCount,
		uintptr(unsafe.Pointer(&apObjects)), // [out] IWbemClassObject **apObjects,
		uintptr(unsafe.Pointer(&uReturned))) // [out] ULONG            *puReturned)
	if int(res) < 0 {
		return nil, NewWmiError(res)
	}

	if uReturned < 1 {
		switch res {
		case WBEM_S_NO_ERROR, WBEM_S_FALSE:
			// No more elements
			return nil, nil
		default:
			return nil, fmt.Errorf("failure advancing enumeration (%d)", res)
		}
	}

	return newInstance(apObjects, e.service), nil
}

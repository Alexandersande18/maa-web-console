package maa

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

// 之后抽象成和Linux调用方式保持一致 先这么拆
const (
	loadDLL = "MaaCore.dll"
	sep     = string(os.PathSeparator)
)

type procCall func(a ...uintptr) (r1, r2 uintptr, lastErr error)

var createAsstFunc = func(name string, dll *syscall.DLL) (procCall, error) {
	proc, err := dll.FindProc(name)
	if err != nil {
		return nil, err
	}

	return func(a ...uintptr) (r1 uintptr, r2 uintptr, lastErr error) {
		return proc.Call(a...)
	}, nil
}

var (
	loadedDLLPath           string
	asstCreate              procCall
	asstCreateEx            procCall
	asstGetVersion          procCall
	asstAppendTask          procCall
	asstConnect             procCall
	asstSetUserDir          procCall
	asstSetTaskParams       procCall
	asstLoadResource        procCall
	asstStart               procCall
	asstStop                procCall
	asstRunning             procCall
	asstDestroy             procCall
	asstGetImage            procCall
	asstGetTasksList        procCall
	asstAsyncScreencap      procCall
	asstSetConnectionExtras procCall
)

func newSyscall(path string) (err error) {
	dllPath, err := filepath.Abs(filepath.Join(path, loadDLL))
	if err != nil {
		return err
	}
	if loadedDLLPath != "" {
		if !strings.EqualFold(loadedDLLPath, dllPath) {
			return fmt.Errorf("all instances must use the same MaaCore DLL: %s != %s", loadedDLLPath, dllPath)
		}
		return nil
	}

	// 使用 LoadLibraryEx 加载 DLL，并指定 LOAD_WITH_ALTERED_SEARCH_PATH 标志
	handle, err := windows.LoadLibraryEx(
		dllPath,
		0,                                     // 保留参数，设为 0
		windows.LOAD_WITH_ALTERED_SEARCH_PATH, // 关键标志
	)
	if err != nil {
		return err
	}

	// 手动构造 syscall.DLL 对象以兼容原有逻辑
	dll := &syscall.DLL{
		Name:   dllPath,
		Handle: syscall.Handle(handle), // 转换句柄类型
	}

	// 绑定函数
	if err = initFunc(dll); err != nil {
		return err
	}
	loadedDLLPath = dllPath
	return nil
}

func initFunc(dll *syscall.DLL) (err error) {
	// Optional for older cores when emulator enhancements are not configured.
	asstSetConnectionExtras, _ = createAsstFunc("AsstSetConnectionExtras", dll)
	asstLoadResource, err = createAsstFunc("AsstLoadResource", dll)
	if err != nil {
		return
	}
	asstCreate, err = createAsstFunc("AsstCreate", dll)
	if err != nil {
		return
	}
	asstCreateEx, err = createAsstFunc("AsstCreateEx", dll)
	if err != nil {
		return
	}
	asstSetUserDir, err = createAsstFunc("AsstSetUserDir", dll)
	if err != nil {
		return
	}

	asstConnect, err = createAsstFunc("AsstConnect", dll)
	if err != nil {
		return
	}
	asstGetVersion, err = createAsstFunc("AsstGetVersion", dll)
	if err != nil {
		return
	}
	asstAppendTask, err = createAsstFunc("AsstAppendTask", dll)
	if err != nil {
		return
	}
	asstSetTaskParams, err = createAsstFunc("AsstSetTaskParams", dll)
	if err != nil {
		return
	}
	asstStart, err = createAsstFunc("AsstStart", dll)
	if err != nil {
		return
	}
	asstStop, err = createAsstFunc("AsstStop", dll)
	if err != nil {
		return
	}
	asstRunning, err = createAsstFunc("AsstRunning", dll)
	if err != nil {
		return
	}
	asstDestroy, err = createAsstFunc("AsstDestroy", dll)
	if err != nil {
		return
	}

	asstGetImage, err = createAsstFunc("AsstGetImage", dll)
	if err != nil {
		return
	}

	asstGetTasksList, err = createAsstFunc("AsstGetTasksList", dll)
	if err != nil {
		return
	}

	asstAsyncScreencap, err = createAsstFunc("AsstAsyncScreencap", dll)
	if err != nil {
		return
	}

	return nil
}

package maa

import (
	"bytes"
	"fmt"
	"log"
	// "os"
	"syscall"
	"testing"
	"time"

	"golang.org/x/sys/windows"
    dto "maa-backend/internal/model/dto"
)

func TestMaaInstance(t *testing.T) {
	cfg, err := LoadConfig()
	if err != nil {
		t.Error(err)
	}
	maa, err := New(cfg.Asst[0])
	if err != nil {
		t.Error(err)
	}
	if maa == nil {
		t.Error("maa is nil")
	}
    maa.AppendHook(&TaskCallback{})
    log.Println(maa.AsstAppendTask(dto.StartUp, &dto.TaskStartUp{
        Enable: true,
        ClientType: "Official",
        StartGameEnabled: true,
    }))
    maa.AsstStart()
    time.Sleep(25 * time.Second)
}

func TestGetImg(t *testing.T)  {
    cfg, err := LoadConfig()
    if err != nil {
        t.Error(err)
    }
    maa, err := New(cfg.Asst[1])
    if err != nil {
        t.Error(err)
    }
    maa.AppendHook(&TaskCallback{})
    if maa == nil {
        t.Error("maa is nil")
    }
    imageData, err := maa.GetImage()
    if err != nil {
        t.Fatalf("GetImage failed: %v", err)
    }
    
    // 验证数据有效性
    if len(imageData) < 8 {
        t.Fatal("Invalid image data")
    }
    
    // 检查PNG文件头（可选）
    if !bytes.Equal(imageData[:8], []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}) {
        t.Log("Warning: Image data does not contain PNG header")
    }
    
    // 保存文件
    // filename := fmt.Sprintf("screenshot_%d.png", time.Now().Unix())
    // err = os.WriteFile(filename, imageData, 0644)
    // if err != nil {
    //     t.Fatalf("Failed to write image file: %v", err)
    // }
    
}

func UTF16PtrFromString(s string) (*uint16, error) {
    return syscall.UTF16PtrFromString(s)
}

type DLL struct {
    Name   string
    Handle syscall.Handle
}

type DLLError struct {
    Err     syscall.Errno
    ObjName string
    Msg     string
}

func (e *DLLError) Error() string {
    return e.Msg
}
func TestLoad(t *testing.T) {
    dllPath := `C:/MAA`

    // 使用 LoadLibraryEx 加载 DLL
    handle, err := windows.LoadLibraryEx(
        dllPath,
        0,
        windows.LOAD_WITH_ALTERED_SEARCH_PATH,
    )
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer windows.FreeLibrary(handle)

    fmt.Println("Successfully loaded DLL")
}
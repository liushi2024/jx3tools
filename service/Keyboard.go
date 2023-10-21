package service

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

type Keyboard struct {
	ctx      context.Context
	lib      *syscall.LazyDLL
	isRun    bool
	isParse  bool
	frontMs  int
	model    int
	key      []Key
	voice    bool
	dir      string
	dll      []byte
	startMp3 *os.File
	stopMp3  *os.File
}

type Key struct {
	Key   interface{} `json:"key"`
	KeyMs interface{} `json:"key_ms"`
}

func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

func (s *Keyboard) Startup(ctx context.Context, dll []byte) {
	s.ctx = ctx
	s.dll = dll
	//logfile, _ := os.OpenFile("log.txt", os.O_TRUNC|os.O_CREATE, 0744)
	//log.SetOutput(logfile)
	log.Println("[开启启动按键模块]")
}

// SyncFrontMs 同步前端间隔数据
func (s *Keyboard) SyncFrontMs(ms int) string {
	s.frontMs = ms
	log.Println("[修改间隔成功][", ms, "]")
	return "success"
}

// SyncFrontModel 同步前端模式数据
func (s *Keyboard) SyncFrontModel(model int) string {
	s.model = model
	log.Println("[修改模式成功][", model, "]")
	return "success"
}

// SyncFrontVoice 同步前端声音数据
func (s *Keyboard) SyncFrontVoice(voice bool) string {
	s.voice = voice
	log.Println("[修改声音成功][", voice, "]")
	return "success"
}

// SyncFrontKey 同步前端数据
func (s *Keyboard) SyncFrontKey(keys string) string {

	var keysStruct []Key

	// 解析JSON字符串并将其填充到结构体数组
	err := json.Unmarshal([]byte(keys), &keysStruct)
	if err != nil {
		log.Println("JSON解析失败:", err)
		return "按键信息同步失败"
	}

	s.key = keysStruct
	log.Println("[同步按键数据成功]", keysStruct)
	return "success"
}

// StartKeyThread 开始线程
func (s *Keyboard) StartKeyThread() {
	s.isParse = false
	if !s.isRun {
		s.isRun = true
		runtime.EventsEmit(s.ctx, "start-thread", nil)
		log.Println("[开始发送线程启动信号]")
		for i := range s.key {
			keyValue := uintptr(s.key[i].Key.(float64))
			msValue, _ := strconv.Atoi(s.key[i].KeyMs.(string))
			go s.ThreadExec(i, keyValue, s.frontMs, s.model, msValue)
		}
	}
}

// StopKeyThread 停止线程
func (s *Keyboard) StopKeyThread() {
	if s.isRun {
		s.isRun = false
		log.Println("[开始发送线程退出信号]")
		runtime.EventsEmit(s.ctx, "stop-thread", nil)
	}
}

// ParseKeyThread 暂停线程
func (s *Keyboard) ParseKeyThread() {
	s.isParse = !s.isParse
}

// DllImport 导入驱动
func (s *Keyboard) DllImport() string {
	s.voice = false
	// 获取临时目录路径
	tempDir := os.TempDir()

	tmpDLLPath := filepath.Join(tempDir, "dllfile.dll")
	_ = ioutil.WriteFile(tmpDLLPath, s.dll, 0644)

	defer os.Remove(tmpDLLPath)
	s.lib = syscall.NewLazyDLL(tmpDLLPath)
	proc := s.lib.NewProc("DD_btn")
	res, _, _ := proc.Call(0)
	if int(res) != 1 {
		log.Println("[驱动注入失败，驱动返回值:", res, "]")
		return "驱动注入失败"
	}
	log.Println("[驱动注入成功！]")
	return "success"
}

// ThreadExec 执行按键输出线程
func (s *Keyboard) ThreadExec(id int, key uintptr, ms int, model int, keyMs int) {
	proc := s.lib.NewProc("DD_key")
	num := 0
	if model == 0 && id != 0 {
		return
	}
	log.Println("[线程已启动][id:", uintptr(id), "][key:", key, "][ms:", uintptr(ms), "][model:", uintptr(model), "]")
	for {
		if !s.isRun {
			return
		}
		if !s.isParse {
			//顺序按键
			if model == 0 {
				nowKey := s.key[num].Key.(float64)
				fmt.Println("在输出", nowKey, "间隔", ms)
				proc.Call(uintptr(nowKey), 1)
				proc.Call(uintptr(nowKey), 2)
				num++
				if num >= len(s.key) {
					num = 0
				}
			}
			//连发按键
			if model == 1 {
				fmt.Println("在输出", key, "间隔", ms)
				proc.Call(key, 1)
				proc.Call(key, 2)
			}

			//按压按键
			if model == 2 {
				fmt.Println("在输出", key, "间隔", ms)
				proc.Call(key, 1)
				proc.Call(key, 2)
			}

		} else {
			fmt.Println("在暂停")
		}
		if keyMs < 50 || model == 0 || model == 2 {
			time.Sleep(time.Duration(ms) * time.Millisecond)
		} else {
			time.Sleep(time.Duration(keyMs) * time.Millisecond)
		}
	}
}

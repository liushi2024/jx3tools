package service

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Keyboard struct {
	ctx       context.Context
	lib       *syscall.LazyDLL
	isRun     bool
	isParse   bool
	frontMs   int
	model     int
	parseType string
	disabled  string
	key       []Key
	dir       string
	dll       []byte
	startMp3  *os.File
	stopMp3   *os.File
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
	if false {
		logfile, _ := os.OpenFile("log.txt", os.O_TRUNC|os.O_CREATE, 0744)
		log.SetOutput(logfile)
	}
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

// SyncDisabled 同步前端独立延迟生效数据
func (s *Keyboard) SyncDisabled(disabled string) string {
	s.disabled = disabled
	log.Println("[修改独立延迟生效成功][", disabled, "]")
	return "success"
}

// SyncParseType 同步前端暂停方式数据
func (s *Keyboard) SyncParseType(parseType string) string {
	s.parseType = parseType
	log.Println("[修改暂停方式成功][", parseType, "]")
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
			go s.ThreadExec(i, keyValue, s.frontMs, s.model)
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
func (s *Keyboard) ThreadExec(id int, key uintptr, ms int, model int) {
	proc := s.lib.NewProc("DD_key")
	//顺序按键时只保留0线程，其余线程退出
	if model == 0 && id != 0 {
		return
	}
	num := id
	log.Println("[线程已启动][id:", uintptr(id), "][key:", key, "][ms:", uintptr(ms), "][model:", uintptr(model), "]")
	for {
		if !s.isRun {
			return
		}
		if !s.isParse {
			//顺序按键
			if model == 0 {
				nowKey := s.key[num].Key.(float64)
				proc.Call(uintptr(nowKey), 1)
				proc.Call(uintptr(nowKey), 2)
				num++
				if num >= len(s.key) {
					num = 0
				}
			}
			//连发按键
			if model == 1 {
				proc.Call(key, 1)
				proc.Call(key, 2)
			}
			//按压按键
			if model == 2 {
				proc.Call(key, 1)
				proc.Call(key, 2)
			}
			log.Println("[model:", model, "][key:", key, "]", ms)
		} else {
			log.Println("[暂停执行]")
		}
		msValue, _ := strconv.Atoi(s.key[num].KeyMs.(string))
		if msValue < 50 || s.disabled == "0" {
			log.Println("统一延迟生效", ms, s.disabled)
			time.Sleep(time.Duration(ms) * time.Millisecond)
		} else {
			log.Println("独立延迟生效", msValue, s.disabled)
			time.Sleep(time.Duration(msValue) * time.Millisecond)
		}
	}
}

func (s *Keyboard) ExportPlans(plans string) string {
	dialog, err := runtime.SaveFileDialog(s.ctx, runtime.SaveDialogOptions{
		Title:   "导出文件",
		Filters: []runtime.FileFilter{{DisplayName: "按键配置", Pattern: "*.json"}},
	})
	if err != nil {
		return err.Error()
	}
	if dialog == "" {
		return "取消保存"
	}
	if !strings.HasSuffix(dialog, ".json") {
		dialog += ".json"
	}
	err = saveToFile(dialog, plans)
	if err != nil {
		return err.Error()
	}
	return "导出成功"
}

func (s *Keyboard) ImportPlans() string {
	dialog, err := runtime.OpenFileDialog(s.ctx, runtime.OpenDialogOptions{
		Title:   "导入文件",
		Filters: []runtime.FileFilter{{DisplayName: "按键配置", Pattern: "*.json"}},
	})
	if err != nil {
		return err.Error()
	}
	if dialog == "" {
		return "取消导入"
	}
	file, err := readFile(dialog)
	if err != nil {
		return err.Error()
	}
	return file
}

func saveToFile(filePath string, content string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func readFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

package service

import (
	"context"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
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
	key      []int
	voice    bool
	dir      string
	startMp3 *os.File
	stopMp3  *os.File
}

func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

func (s *Keyboard) Startup(ctx context.Context) {
	s.ctx = ctx
	logfile, _ := os.OpenFile("log.txt", os.O_TRUNC|os.O_CREATE, 0744)
	log.SetOutput(logfile)
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
func (s *Keyboard) SyncFrontKey(key []int) string {
	s.key = key
	log.Println("[同步按键数据成功]", key)
	return "success"
}

// StartKeyThread 开始线程
func (s *Keyboard) StartKeyThread() {
	s.isParse = false
	if !s.isRun {
		s.isRun = true
		go s.playMp3(true)
		runtime.EventsEmit(s.ctx, "start-thread", nil)
		log.Println("[开始发送线程启动信号]")
		for i := range s.key {
			go s.ThreadExec(i, uintptr(s.key[i]), s.frontMs, s.model)
		}
	}
}

// StopKeyThread 停止线程
func (s *Keyboard) StopKeyThread() {
	if s.isRun {
		s.isRun = false
		go s.playMp3(false)
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
	dir, _ := os.Getwd()
	s.dir = dir
	dllPath := dir + "\\dd202x.8.x64.dll"
	log.Println("[开始挂载驱动文件,路径:", dllPath, "]")
	//判断文件存不存在
	if _, err := os.Stat(dllPath); err != nil {
		log.Println("[驱动挂载失败][", err.Error(), "]")
	}
	//load dll,使用懒加载模式，也就是真正调用 API 的时候才会加载
	s.lib = syscall.NewLazyDLL(dllPath)
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
			if model == 1 {
				fmt.Println("在输出", key, "间隔", ms)
				proc.Call(key, 1)
				proc.Call(key, 2)
			} else {
				nowKey := s.key[num]
				fmt.Println("在输出", nowKey, "间隔", ms)
				proc.Call(uintptr(nowKey), 1)
				proc.Call(uintptr(nowKey), 2)
				num++
				if num >= len(s.key) {
					num = 0
				}
			}
		} else {
			fmt.Println("在暂停")
		}
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}

// 播放开始和停止音乐
func (s *Keyboard) playMp3(isStart bool) {
	if !s.voice {
		return
	}
	fileStr := "\\start.mp3"
	if !isStart {
		fileStr = "\\stop.mp3"
	}
	mp3File, err := os.Open(s.dir + fileStr)
	if err != nil {
		log.Fatal(err)
	}
	println(mp3File.Name())
	stream, format, err := mp3.Decode(mp3File)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()
	err2 := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)

	println("开始播放音乐", err2)

	speaker.Play(beep.Seq(stream, beep.Callback(func() {
		done <- true
	})))
	<-done
	println("音乐播放结束")
}

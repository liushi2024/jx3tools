package service

import (
	hook "github.com/robotn/gohook"
	"log"
	"math"
)

type HotKey struct {
	startKey int
	stopKey  int
	parseKey int
	key      *Keyboard
}

func NewHotKey() *HotKey {
	return &HotKey{}
}

func (s *HotKey) Startup(key *Keyboard) {
	s.key = key
	go s.add()
}

func (s *HotKey) Close() {
	hook.End()
	log.Println("钩子已经关闭")
}

// SyncFrontStart 同步前端开始数据
func (s *HotKey) SyncFrontStart(start int) string {
	s.startKey = start
	return "success"
}

// SyncFrontStop 同步前端停止数据
func (s *HotKey) SyncFrontStop(stop int) string {
	s.stopKey = stop
	return "success"
}

// SyncFrontParse 同步前端暂停数据
func (s *HotKey) SyncFrontParse(parse int) string {
	s.parseKey = parse
	return "success"
}

func (s *HotKey) add() {
	hooks := hook.Start()
	s.keyHook(hooks)
}

func (s *HotKey) keyHook(hooks chan hook.Event) {
	defer s.keyHook(hooks)

	for {
		select {
		case ev := <-hooks:
			s.processEvent(ev)
		}
	}
}

func (s *HotKey) processEvent(ev hook.Event) {
	s.listenKeyBoard(ev)
	s.listenMouse(ev)
}
func (s *HotKey) listenKeyBoard(ev hook.Event) {
	if s.key.model == 2 {
		s.keyPress(ev)
	} else {
		s.keyHold(ev)
	}
}

func (s *HotKey) listenMouse(ev hook.Event) {
	//监听鼠标点击
	if ev.Kind == hook.MouseHold {
		if s.startKey != s.stopKey {
			if ev.Button == uint16(s.startKey-900) {
				s.key.StartKeyThread()
			}
			if ev.Button == uint16(s.stopKey-900) {
				s.key.StopKeyThread()
			}
		} else {
			if ev.Button == uint16(s.startKey-900) {
				if !s.key.isRun {
					s.key.StartKeyThread()
				} else {
					s.key.StopKeyThread()
				}
			}
		}

		if ev.Button == uint16(s.parseKey-900) {
			s.key.isParse = !s.key.isParse
		}
	}

	//监听鼠标滚轮
	if ev.Kind == hook.MouseWheel {

		if math.Abs(float64(ev.Rotation)) != 1 {
			return
		}

		if s.startKey != s.stopKey {
			if ev.Rotation == int32(s.startKey-907) {
				s.key.StartKeyThread()
			}
			if ev.Rotation == int32(s.stopKey-907) {
				s.key.StopKeyThread()
			}
		} else {
			if ev.Rotation == int32(s.startKey-907) {
				if !s.key.isRun {
					s.key.StartKeyThread()
				} else {
					s.key.StopKeyThread()
				}
			}
		}
	}
}

// 顺序和连放模式
func (s *HotKey) keyHold(ev hook.Event) {
	//	监听键盘弹起
	if ev.Kind == hook.KeyUp {
		if s.startKey != s.stopKey {
			if ev.Rawcode == uint16(s.startKey) {
				s.key.StartKeyThread()
			}
			if ev.Rawcode == uint16(s.stopKey) {
				s.key.StopKeyThread()
			}
		} else {
			if ev.Rawcode == uint16(s.startKey) {
				if !s.key.isRun {
					s.key.StartKeyThread()
				} else {
					s.key.StopKeyThread()
				}
			}
		}
		if s.key.parseType == "0" {
			if ev.Rawcode == uint16(s.parseKey) {
				s.key.isParse = !s.key.isParse
			}
		}
		if s.key.parseType == "1" {
			if ev.Rawcode == uint16(s.parseKey) {
				s.key.isParse = false
			}
		}
	}
	if ev.Kind == hook.KeyHold {
		if s.key.parseType == "1" {
			if ev.Rawcode == uint16(s.parseKey) {
				s.key.isParse = true
			}
		}
	}

}

// 键盘按压模式
func (s *HotKey) keyPress(ev hook.Event) {
	//	监听键盘弹起
	if ev.Kind == hook.KeyHold {
		if ev.Rawcode == uint16(s.startKey) && !s.key.isRun {
			if !s.key.isRun {
				s.key.StartKeyThread()
			}
		}
	}
	if ev.Kind == hook.KeyUp {
		if ev.Rawcode == uint16(s.startKey) {
			if s.key.isRun {
				s.key.StopKeyThread()
			}
		}
	}
}

package service

import (
	hook "github.com/robotn/gohook"
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
	defer hook.End()
	for ev := range hooks {

		//	监听键盘弹起
		if ev.Kind == hook.KeyHold {
			if ev.Rawcode == uint16(s.startKey) {
				s.key.StartKeyThread()
			}
			if ev.Rawcode == uint16(s.stopKey) {
				s.key.StopKeyThread()
			}
			if ev.Rawcode == uint16(s.parseKey) {
				s.key.isParse = !s.key.isParse
			}
		}

		if ev.Kind == hook.MouseHold {
			if ev.Button == uint16(s.startKey-900) {
				s.key.StartKeyThread()
			}
			if ev.Button == uint16(s.stopKey-900) {
				s.key.StopKeyThread()
			}
			if ev.Button == uint16(s.parseKey-900) {
				s.key.isParse = !s.key.isParse
			}
		}

		if ev.Kind == hook.MouseWheel {
			if ev.Rotation == int32(s.startKey-907) {
				s.key.StartKeyThread()
			}
			if ev.Rotation == int32(s.stopKey-907) {
				s.key.StopKeyThread()
			}
		}

	}
}

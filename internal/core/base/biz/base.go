package biz

import (
	"time"

	"github.com/muidea/magicCommon/event"
	"github.com/muidea/magicCommon/task"
)

type Base struct {
	id                string
	eventHub          event.Hub
	backgroundRoutine task.BackgroundRoutine
}

type invokeTask struct {
	funcPtr func()
}

func (s *invokeTask) Run() {
	s.funcPtr()
}

func New(
	id string,
	eventHub event.Hub,
	backgroundRoutine task.BackgroundRoutine) Base {
	return Base{
		id:                id,
		eventHub:          eventHub,
		backgroundRoutine: backgroundRoutine,
	}
}

func (s *Base) ID() string {
	return s.id
}

func (s *Base) Subscribe(eventID string, observer event.Observer) {
	s.eventHub.Subscribe(eventID, observer)
}

func (s *Base) Unsubscribe(eventID string, observer event.Observer) {
	s.eventHub.Unsubscribe(eventID, observer)
}

func (s *Base) PostEvent(event event.Event) {
	s.eventHub.Post(event)
}

func (s *Base) SendEvent(event event.Event) event.Result {
	return s.eventHub.Send(event)
}

func (s *Base) CallEvent(event event.Event) event.Result {
	return s.eventHub.Call(event)
}

func (s *Base) Invoke(funcPtr func()) {
	taskPtr := &invokeTask{funcPtr: funcPtr}

	s.backgroundRoutine.Post(taskPtr)
}

func (s *Base) Timer(intervalValue time.Duration, offsetValue time.Duration, funcPtr func()) {
	taskPtr := &invokeTask{funcPtr: funcPtr}
	s.backgroundRoutine.Timer(taskPtr, intervalValue, offsetValue)
}

func (s *Base) BroadCast(eid string, header event.Values, val interface{}) {
	event := event.NewEvent(eid, s.ID(), s.RootDestination(), header, val)
	s.eventHub.Post(event)
}

func (s *Base) RootDestination() string {
	return "/#"
}

func (s *Base) InnerDestination() string {
	return s.ID()
}

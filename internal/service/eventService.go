package service

import (
	"context"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/charmbracelet/log"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Other services use this service can avoid writing a large number
// of if guard by injecting an empty EventService by default since
// 'PushEvent'
type EventService struct {
	ctx context.Context
}

func (s *EventService) InjectContext(ctx context.Context) {
	s.ctx = ctx
}

func (s *EventService) PushEvent(eventName string, v ...any) {
	// There's a chance that this method will be called before ctx injected
	// This is not a big deal so we currently skip this
	// NOTE: This guard is intentional design, see comments on EventService
	// and the way RivalInfoService uses EventService
	if s.ctx != nil {
		runtime.EventsEmit(s.ctx, eventName, v...)
	} else {
		log.Debugf("skipping event: %s", eventName)
	}
}

func (s *EventService) RefreshPage() {
	s.PushEvent("global:refresh")
}

func (s *EventService) NotifyInfo(content string) {
	s.PushEvent("global:notify", dto.NotificationDto{
		Type:    "info",
		Content: content,
	})
}

func (s *EventService) NotifySuccess(content string) {
	s.PushEvent("global:notify", dto.NotificationDto{
		Type:    "success",
		Content: content,
	})
}

func (s *EventService) NotifyError(content string) {
	s.PushEvent("global:notify", dto.NotificationDto{
		Type:    "error",
		Content: content,
	})
}

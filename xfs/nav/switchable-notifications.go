package nav

import (
	"fmt"
	"math"

	"github.com/snivilised/extendio/internal/lo"
)

type notificationBiEnum uint32

const (
	notificationUndefinedEn notificationBiEnum = 0
	notificationBeginEn                        = 1 << (iota - 1)
	notificationEndEn
	notificationDescendEn
	notificationAscendEn
	notificationStartEn
	notificationStopEn
	notificationAllEn = math.MaxUint32
)

type switchableBase struct {
	muted bool
}

func (s *switchableBase) mute(value bool) {
	s.muted = value
}

type switchableBegin struct {
	switchableBase
	handler BeginHandler
}

func (s *switchableBegin) invoke(state *NavigationState) {
	if !s.muted {
		s.handler(state)
	}
}

type switchableEnd struct {
	switchableBase
	handler EndHandler
}

func (s *switchableEnd) invoke(result *TraverseResult) {
	if !s.muted {
		s.handler(result)
	}
}

type switchableAscendancy struct {
	switchableBase
	handler AscendancyHandler
}

func (s *switchableAscendancy) invoke(item *TraverseItem) {
	if !s.muted {
		s.handler(item)
	}
}

type switchableListen struct {
	switchableBase
	handler ListenHandler
}

func (s *switchableListen) invoke(description string) {
	if !s.muted {
		s.handler(description)
	}
}

type switchable map[notificationBiEnum]*switchableBase

type notificationsSink struct {
	begin   switchableBegin
	end     switchableEnd
	descend switchableAscendancy
	ascend  switchableAscendancy
	start   switchableListen
	stop    switchableListen
	all     switchable
}

func (n *notificationsSink) init(notifications *Notifications) {
	n.begin = switchableBegin{
		handler: notifications.OnBegin,
	}
	n.end = switchableEnd{
		handler: notifications.OnEnd,
	}
	n.descend = switchableAscendancy{
		handler: notifications.OnDescend,
	}
	n.ascend = switchableAscendancy{
		handler: notifications.OnAscend,
	}
	n.start = switchableListen{
		handler: notifications.OnStart,
	}
	n.stop = switchableListen{
		handler: notifications.OnStop,
	}
	n.all = switchable{
		notificationBeginEn:   &n.begin.switchableBase,
		notificationEndEn:     &n.end.switchableBase,
		notificationDescendEn: &n.descend.switchableBase,
		notificationAscendEn:  &n.ascend.switchableBase,
		notificationStartEn:   &n.start.switchableBase,
		notificationStopEn:    &n.stop.switchableBase,
	}
}

func (n *notificationsSink) mute(notifyEn notificationBiEnum, values ...bool) {
	if notifyEn == notificationUndefinedEn {
		panic(NewInvalidNotificationMuteRequestedNativeError(fmt.Sprintf("%v", notifyEn)))
	}

	value := lo.TernaryF(len(values) > 0,
		func() bool {
			return values[0]
		},
		func() bool {
			return true
		})

	for referenceEn := range n.all {
		n.muzzle(&muzzleParams{
			notifyEn:    notifyEn,
			referenceEn: referenceEn,
			value:       value,
		})
	}
}

type muzzleParams struct {
	notifyEn    notificationBiEnum
	referenceEn notificationBiEnum
	value       bool
}

func (n *notificationsSink) muzzle(params *muzzleParams) {
	if (params.notifyEn & params.referenceEn) > 0 {
		base := n.all[params.referenceEn]
		base.mute(params.value)
	}
}

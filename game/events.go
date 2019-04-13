package game

type EventType string
type EventAction string

const PlayerEventsType = EventType("PLAYER")

const EventActionTakeDamage = EventAction("TAKE_DAMAGE")
const EventActionRegenerateAll = EventAction("REGENERATE_ALL")

type Event struct {
	Action  EventAction
	Message string
	Type    EventType
}

type EventSubscriber interface {
	On(e *Event)
}

type EventManager struct {
	Subscribers map[EventType][]EventSubscriber
}

func NewEventManager() *EventManager {
	em := &EventManager{}
	em.Subscribers = make(map[EventType][]EventSubscriber)
	return em
}

func (m *EventManager) Subscribe(subscriber EventSubscriber, eventType EventType) {
	m.Subscribers[eventType] = append(m.Subscribers[eventType], subscriber)
}

func (m *EventManager) Dispatch(e *Event) {
	subscribers, exists := m.Subscribers[e.Type]
	if exists {
		for _, subscriber := range subscribers {
			go subscriber.On(e)
		}
	}
}

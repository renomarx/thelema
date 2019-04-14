package game

type EventType string

const PlayerEventsType = EventType("PLAYER")

const ActionEat = "EAT"
const ActionHurt = "HURT"
const ActionDie = "DIE"
const ActionChangeLevel = "CHANGE_LEVEL"

type Event struct {
	Action  string
	Message string
	Type    EventType
	Payload map[string]string
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

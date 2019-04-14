package game

const ActionWalk = "WALK"
const ActionEat = "EAT"
const ActionHurt = "HURT"
const ActionDie = "DIE"
const ActionChangeLevel = "CHANGE_LEVEL"
const ActionOpenDoor = "OPEN_DOOR"
const ActionCloseDoor = "CLOSE_DOOR"
const ActionAttack = "ATTACK"
const ActionPower = "POWER"
const ActionTalk = "TALK"
const ActionTake = "TAKE"
const ActionExplode = "EXPLODE"

const ActionMenuSelect = "MENU_SELECT"
const ActionMenuConfirm = "MENU_CONFIRM"
const ActionMenuOpen = "MENU_OPEN"
const ActionMenuClose = "MENU_CLOSE"

type Event struct {
	Action  string
	Message string
	Payload map[string]string
}

type EventSubscriber interface {
	On(e *Event)
}

type EventManager struct {
	Subscribers []EventSubscriber
}

func NewEventManager() *EventManager {
	em := &EventManager{}
	return em
}

func (m *EventManager) Subscribe(subscriber EventSubscriber) {
	m.Subscribers = append(m.Subscribers, subscriber)
}

func (m *EventManager) Dispatch(e *Event) {
	for _, subscriber := range m.Subscribers {
		subscriber.On(e)
	}
}

package dialog

import (
	"context"
	"errors"
	cache "vpn_go_bot/internal/infrastructure/cache"
)

var (
	DuplicateStateError      = errors.New("duplicate state")
	NoStateWithSuchNameError = errors.New("no state with such name")
)

type State struct {
	Name string                 `json:"name"`
	data map[string]interface{} `json:"data"`
}

func (s *State) getData() map[string]interface{} {
	return s.data
}

func (s *State) updateData(key string, val interface{}) {
	if s.data == nil {
		s.data = make(map[string]interface{})
	}
	s.data[key] = val
}

type StateGroup struct {
	states map[string]*State
}

func (sg *StateGroup) GetState(name string) *State {
	if state, exists := sg.states[name]; exists {
		return state
	}
	return nil
}

func (sg *StateGroup) AddState(state *State) error {
	if sg.states == nil {
		sg.states = make(map[string]*State)
	}
	_, ok := sg.states[state.Name]
	if ok {
		return DuplicateStateError
	}
	sg.states[state.Name] = state
	return nil
}

// Приходит апдейт -> создается FSM(загружается текущий State,state Data) -> FSMcontext -> State -> StateStorage

type FSMContext struct {
	userID       int
	States       *StateGroup
	Storage      *cache.RedisClient
	currentState *State
}

// TODO #3: Добавить обработку что нет в кэше стейта
// TODO
// NewFSMContext create new FSMContext to handle chaage in state and dialog cache
func NewFSMContext(ctx context.Context, userID int, States *StateGroup, storage *cache.RedisClient) *FSMContext {
	val, _ := storage.HGet(ctx, string(userID), "current_state")
	return &FSMContext{
		userID:       userID,
		States:       States,
		Storage:      storage,
		currentState: States.GetState(val),
	}

}

// GetState return current State
func (fsm *FSMContext) GetState() *State {
	return fsm.currentState
}

// SetState set current state and save to cache
func (fsm *FSMContext) SetState(ctx context.Context, stateName string) error {
	fsm.currentState = fsm.States.GetState(stateName)
	err := fsm.Storage.HSet(ctx, string(fsm.userID), "current_state", stateName)
	if err != nil {
		return err
	}
	return nil
}

// GetCurrentStateData return current state data from cache only in string
func (fsm *FSMContext) GetCurrentStateData(ctx context.Context) map[string]string {
	res, err := fsm.Storage.HGetAll(ctx, string(fsm.userID)+":"+fsm.currentState.Name)
	if err != nil {
		return nil
	}
	return res
}

func (fsm *FSMContext) UpdateStateData(ctx context.Context, key string, val interface{}) {
	fsm.Storage.HSet(ctx, string(fsm.userID)+":"+fsm.currentState.Name, key, val)
	fsm.currentState.updateData(key, val)
	// userid + ":" + stateName + ":" + key = state_data
}

func (fsm *FSMContext) DeleteData(ctx context.Context, key string) {
	fsm.Storage.HDel(ctx, string(fsm.userID)+":"+fsm.currentState.Name, key)
	delete(fsm.currentState.getData(), key)
}

// TODO #9: Добавить метод очистки всех данных для пользователя, судя по всему   

func (fsm *FSMContext) Clear(ctx context.Context) error {
	err := fsm.Storage.Delete(ctx, string(fsm.userID))
	if err != nil {
		return err
	}
	return nil
}

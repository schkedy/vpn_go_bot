package handler

type CallbackHandlerGroup struct {
	handlers map[string]HandlerFunc
}

func NewCallbackHandlers() *CallbackHandlerGroup {
	handlers := make(map[string]HandlerFunc)
	// add callback handlers to the map
	return &CallbackHandlerGroup{
		handlers: handlers,
	}
}

func (ch *CallbackHandlerGroup) RegisterHandlers(handlers map[string]HandlerFunc) {
	if ch.handlers == nil {
		ch.handlers = make(map[string]HandlerFunc)
	}
	for k, v := range handlers {
		ch.handlers[k] = v
	}
}

func (ch *CallbackHandlerGroup) Validate(s string) (command string) {
	// validate callback data and return command
	return command
}

func (ch *CallbackHandlerGroup) GetHandler(command string) (HandlerFunc, bool) {
	if handler, exists := ch.handlers[command]; exists {
		return handler, true
	}

	var matched HandlerFunc
	longestPrefix := 0
	for prefix, handler := range ch.handlers {
		if len(prefix) == 0 {
			continue
		}
		if len(command) >= len(prefix) && command[:len(prefix)] == prefix {
			if len(prefix) > longestPrefix {
				longestPrefix = len(prefix)
				matched = handler
			}
		}
	}
	if matched != nil {
		return matched, true
	}

	return nil, false
}

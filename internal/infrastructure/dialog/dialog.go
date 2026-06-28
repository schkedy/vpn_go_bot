package dialog

type Dialog struct {
	Windows  []*Window
	handlers map[string]HandlerFunc
}

// TODO валидировать окна, чтобы не было одинаковых имен, хэндлеров и т.д.
func NewDialog(windows ...*Window) *Dialog {
	d := &Dialog{
		Windows: windows,
	}
	for _, window := range d.Windows {
		d.mergeHandlers(window.GetHandlers())
	}
	return d
}

func (d *Dialog) GetWindow(s State) *Window {
	for _, window := range d.Windows {
		if window.State == s {
			return window
		}
	}
	return nil
}

func (d *Dialog) GetHandlers() map[string]HandlerFunc {
	return d.handlers
}

func (d *Dialog) mergeHandlers(windowHandlers map[string]HandlerFunc) {
	if d.handlers == nil {
		d.handlers = make(map[string]HandlerFunc)
	}
	for k, v := range windowHandlers {
		d.handlers[k] = v
	}
}

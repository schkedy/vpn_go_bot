package dialog

import "context"

// Getter - provide dialog data and deps for render window, it is used for dynamic text and widgets, it is called every time when we need to render window, it always receives dialogManager and deps, and returns data for render
type Getter func(ctx context.Context, fsm *FSMContext, dialogSession *DialogSession, deps map[string]interface{}) map[string]interface{}

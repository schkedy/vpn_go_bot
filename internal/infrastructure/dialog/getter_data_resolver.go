package dialog

func resolveDialogManagerGetterData(dialogManager *DialogManager) map[string]interface{} {
	if dialogManager == nil || dialogManager.dialog == nil || dialogManager.Session == nil {
		return map[string]interface{}{}
	}

	window := dialogManager.dialog.GetWindow(dialogManager.Session.State)
	if window == nil {
		return map[string]interface{}{}
	}

	return window.getGetterData(dialogManager)
}


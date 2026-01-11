let editorFocused = $state(false);

export class EditorFocus {
	changeFocus(focused: boolean) {
		editorFocused = focused;
	}

	get focused() {
		return editorFocused;
	}
}

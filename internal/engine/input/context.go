package input

import "github.com/go-gl/glfw/v3.2/glfw"

type Context struct {
	Keyboard *KeyboardState
	Mouse    *MouseState
}

func InitContext(window *glfw.Window) Context {
	return Context{
		Keyboard: initKeyboard(window),
		Mouse:    initMouse(window),
	}
}

func initKeyboard(window *glfw.Window) *KeyboardState {
	keyboard := NewKeyboardState()
	window.SetKeyCallback(toGlfwKeyCallback(keyboard.OnKeyChange))

	return keyboard
}

func initMouse(window *glfw.Window) *MouseState {
	mouse := NewMouseState()
	window.SetCursorPosCallback(toGlfwCursorPositionCallback(mouse.OnPositionChange))
	window.SetScrollCallback(toGlfwScrollCallback(mouse.OnScrollChange))
	window.SetMouseButtonCallback(toGlfwMouseButtonCallback(mouse.OnButtonChange))

	return mouse
}

func toGlfwKeyCallback(f func(key glfw.Key, action glfw.Action)) glfw.KeyCallback {
	return func(_ *glfw.Window, key glfw.Key, _ int, action glfw.Action, _ glfw.ModifierKey) {
		f(key, action)
	}
}

func toGlfwCursorPositionCallback(f func(xpos float64, ypos float64)) glfw.CursorPosCallback {
	return func(_ *glfw.Window, xpos, ypos float64) { f(xpos, ypos) }
}

func toGlfwScrollCallback(f func(yoff float64)) glfw.ScrollCallback {
	return func(_ *glfw.Window, xoff, yoff float64) { f(yoff) }
}

func toGlfwMouseButtonCallback(f func(button glfw.MouseButton, action glfw.Action)) glfw.MouseButtonCallback {
	return func(_ *glfw.Window, button glfw.MouseButton, action glfw.Action, _ glfw.ModifierKey) {
		f(button, action)
	}
}

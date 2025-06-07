package input

import (
	ebiten "github.com/hajimehoshi/ebiten/v2"
)

var defaultKeyMap = map[ebiten.Key]int{
	ebiten.KeyZ:     0, // A
	ebiten.KeyX:     1, // B
	ebiten.KeyA:     2, // Select
	ebiten.KeyS:     3, // Start
	ebiten.KeyUp:    4,
	ebiten.KeyDown:  5,
	ebiten.KeyLeft:  6,
	ebiten.KeyRight: 7,
}

// Xbox 控制器按钮映射到 NES 控制器顺序
var gamepadMap = map[ebiten.GamepadButton]int{
	ebiten.GamepadButton0:  0, // A
	ebiten.GamepadButton1:  1, // B
	ebiten.GamepadButton6:  2, // Select
	ebiten.GamepadButton7:  3, // Start
	ebiten.GamepadButton12: 4, // Dpad Up
	ebiten.GamepadButton13: 5, // Dpad Down
	ebiten.GamepadButton14: 6, // Dpad Left
	ebiten.GamepadButton15: 7, // Dpad Right
}

// 每帧调用，获取输入状态
func PollInput() [8]bool {
	// 优先手柄
	gamepads := ebiten.GamepadIDs()
	if len(gamepads) > 0 {
		return pollGamepad(gamepads[0])
	}
	return pollKeyboard()
}

func pollKeyboard() [8]bool {
	var state [8]bool
	for key, index := range defaultKeyMap {
		if ebiten.IsKeyPressed(key) {
			state[index] = true
		}
	}
	return state
}

func pollGamepad(id ebiten.GamepadID) [8]bool {
	var state [8]bool
	for button, index := range gamepadMap {
		if ebiten.IsGamepadButtonPressed(id, button) {
			state[index] = true
		}
	}
	return state
}

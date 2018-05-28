package keyboard

import (
	"github.com/hajimehoshi/ebiten"
)

var ebitenASCIIMap map[ebiten.Key]uint8 // ebiten keys mapped to ASCII
var shiftMap map[uint8]uint8            // ebiten keys mapped to ASCII when shift is pressed
var controlMap map[uint8]uint8          // ebiten keys mapped to ASCII when control is pressed

var keyBoardData uint8                 // Contents of the $c000 address
var strobe uint8                       // Contents of the $c010 address
var previousKeysPressed map[uint8]bool // Keep track of what keys have been pressed in the previous round
var capsLock bool                      // Is capslock down

// Init the keyboard state and ebiten translation tables
func Init() {
	keyBoardData = 0
	strobe = 0
	capsLock = true

	ebitenASCIIMap = make(map[ebiten.Key]uint8)
	shiftMap = make(map[uint8]uint8)
	controlMap = make(map[uint8]uint8)
	previousKeysPressed = make(map[uint8]bool)

	ebitenASCIIMap[ebiten.KeyLeft] = 8
	ebitenASCIIMap[ebiten.KeyTab] = 9
	ebitenASCIIMap[ebiten.KeyDown] = 10
	ebitenASCIIMap[ebiten.KeyUp] = 11
	ebitenASCIIMap[ebiten.KeyEnter] = 13
	ebitenASCIIMap[ebiten.KeyRight] = 21
	ebitenASCIIMap[ebiten.KeyEscape] = 27
	ebitenASCIIMap[ebiten.KeyDelete] = 127

	ebitenASCIIMap[ebiten.Key0] = '0'
	ebitenASCIIMap[ebiten.Key1] = '1'
	ebitenASCIIMap[ebiten.Key2] = '2'
	ebitenASCIIMap[ebiten.Key3] = '3'
	ebitenASCIIMap[ebiten.Key4] = '4'
	ebitenASCIIMap[ebiten.Key5] = '5'
	ebitenASCIIMap[ebiten.Key6] = '6'
	ebitenASCIIMap[ebiten.Key7] = '7'
	ebitenASCIIMap[ebiten.Key8] = '8'
	ebitenASCIIMap[ebiten.Key9] = '9'
	ebitenASCIIMap[ebiten.KeyA] = 'a'
	ebitenASCIIMap[ebiten.KeyB] = 'b'
	ebitenASCIIMap[ebiten.KeyC] = 'c'
	ebitenASCIIMap[ebiten.KeyD] = 'd'
	ebitenASCIIMap[ebiten.KeyE] = 'e'
	ebitenASCIIMap[ebiten.KeyF] = 'f'
	ebitenASCIIMap[ebiten.KeyG] = 'g'
	ebitenASCIIMap[ebiten.KeyH] = 'h'
	ebitenASCIIMap[ebiten.KeyI] = 'i'
	ebitenASCIIMap[ebiten.KeyJ] = 'j'
	ebitenASCIIMap[ebiten.KeyK] = 'k'
	ebitenASCIIMap[ebiten.KeyL] = 'l'
	ebitenASCIIMap[ebiten.KeyM] = 'm'
	ebitenASCIIMap[ebiten.KeyN] = 'n'
	ebitenASCIIMap[ebiten.KeyO] = 'o'
	ebitenASCIIMap[ebiten.KeyP] = 'p'
	ebitenASCIIMap[ebiten.KeyQ] = 'q'
	ebitenASCIIMap[ebiten.KeyR] = 'r'
	ebitenASCIIMap[ebiten.KeyS] = 's'
	ebitenASCIIMap[ebiten.KeyT] = 't'
	ebitenASCIIMap[ebiten.KeyU] = 'u'
	ebitenASCIIMap[ebiten.KeyV] = 'v'
	ebitenASCIIMap[ebiten.KeyW] = 'w'
	ebitenASCIIMap[ebiten.KeyX] = 'x'
	ebitenASCIIMap[ebiten.KeyY] = 'y'
	ebitenASCIIMap[ebiten.KeyZ] = 'z'
	ebitenASCIIMap[ebiten.KeyApostrophe] = '\''
	ebitenASCIIMap[ebiten.KeyBackslash] = '\\'
	ebitenASCIIMap[ebiten.KeyComma] = ','
	ebitenASCIIMap[ebiten.KeyEqual] = '='
	ebitenASCIIMap[ebiten.KeyGraveAccent] = '`'
	ebitenASCIIMap[ebiten.KeyLeftBracket] = '['
	ebitenASCIIMap[ebiten.KeyMinus] = '-'
	ebitenASCIIMap[ebiten.KeyPeriod] = '.'
	ebitenASCIIMap[ebiten.KeyRightBracket] = ']'
	ebitenASCIIMap[ebiten.KeySemicolon] = ';'
	ebitenASCIIMap[ebiten.KeySlash] = '/'
	ebitenASCIIMap[ebiten.KeySpace] = ' '

	shiftMap['1'] = '!'
	shiftMap['2'] = '@'
	shiftMap['3'] = '#'
	shiftMap['4'] = '$'
	shiftMap['5'] = '%'
	shiftMap['6'] = '^'
	shiftMap['7'] = '&'
	shiftMap['8'] = '*'
	shiftMap['9'] = '('
	shiftMap['0'] = ')'
	shiftMap['-'] = '_'
	shiftMap['='] = '+'
	shiftMap['a'] = 'A'
	shiftMap['b'] = 'B'
	shiftMap['c'] = 'C'
	shiftMap['d'] = 'D'
	shiftMap['e'] = 'E'
	shiftMap['f'] = 'F'
	shiftMap['g'] = 'G'
	shiftMap['h'] = 'H'
	shiftMap['i'] = 'I'
	shiftMap['j'] = 'J'
	shiftMap['k'] = 'K'
	shiftMap['l'] = 'L'
	shiftMap['m'] = 'M'
	shiftMap['n'] = 'N'
	shiftMap['o'] = 'O'
	shiftMap['p'] = 'P'
	shiftMap['q'] = 'Q'
	shiftMap['r'] = 'R'
	shiftMap['s'] = 'S'
	shiftMap['t'] = 'T'
	shiftMap['u'] = 'U'
	shiftMap['v'] = 'V'
	shiftMap['w'] = 'W'
	shiftMap['x'] = 'X'
	shiftMap['y'] = 'Y'
	shiftMap['z'] = 'Z'
	shiftMap[','] = '<'
	shiftMap['.'] = '>'
	shiftMap['/'] = '?'
	shiftMap['`'] = '~'
	shiftMap['['] = '{'
	shiftMap[']'] = '}'
	shiftMap[';'] = ':'
	shiftMap['\''] = '"'
	shiftMap['\\'] = '|'
	shiftMap[' '] = ' '

	controlMap['A'] = 'A' - 0x40
	controlMap['B'] = 'B' - 0x40
	controlMap['C'] = 'C' - 0x40
	controlMap['D'] = 'D' - 0x40
	controlMap['E'] = 'E' - 0x40
	controlMap['F'] = 'F' - 0x40
	controlMap['G'] = 'G' - 0x40
	controlMap['H'] = 'H' - 0x40
	controlMap['I'] = 'I' - 0x40
	controlMap['J'] = 'J' - 0x40
	controlMap['K'] = 'K' - 0x40
	controlMap['L'] = 'L' - 0x40
	controlMap['M'] = 'M' - 0x40
	controlMap['N'] = 'N' - 0x40
	controlMap['O'] = 'O' - 0x40
	controlMap['P'] = 'P' - 0x40
	controlMap['Q'] = 'Q' - 0x40
	controlMap['R'] = 'R' - 0x40
	controlMap['S'] = 'S' - 0x40
	controlMap['T'] = 'T' - 0x40
	controlMap['U'] = 'U' - 0x40
	controlMap['V'] = 'V' - 0x40
	controlMap['W'] = 'W' - 0x40
	controlMap['X'] = 'X' - 0x40
	controlMap['Y'] = 'Y' - 0x40
	controlMap['Z'] = 'Z' - 0x40
	controlMap[']'] = 0x5d
	controlMap['`'] = 0x60

	controlMap['a'] = 'a' - 0x60
	controlMap['b'] = 'b' - 0x60
	controlMap['c'] = 'c' - 0x60
	controlMap['d'] = 'd' - 0x60
	controlMap['e'] = 'e' - 0x60
	controlMap['f'] = 'f' - 0x60
	controlMap['g'] = 'g' - 0x60
	controlMap['h'] = 'h' - 0x60
	controlMap['i'] = 'i' - 0x60
	controlMap['j'] = 'j' - 0x60
	controlMap['k'] = 'k' - 0x60
	controlMap['l'] = 'l' - 0x60
	controlMap['m'] = 'm' - 0x60
	controlMap['n'] = 'n' - 0x60
	controlMap['o'] = 'o' - 0x60
	controlMap['p'] = 'p' - 0x60
	controlMap['q'] = 'q' - 0x60
	controlMap['r'] = 'r' - 0x60
	controlMap['s'] = 's' - 0x60
	controlMap['t'] = 't' - 0x60
	controlMap['u'] = 'u' - 0x60
	controlMap['v'] = 'v' - 0x60
	controlMap['w'] = 'w' - 0x60
	controlMap['x'] = 'x' - 0x60
	controlMap['y'] = 'y' - 0x60
	controlMap['z'] = 'z' - 0x60
	controlMap['}'] = 0x5d
	controlMap['~'] = 0x60
}

// Poll queries ebiten's keyboard state and transforms that into ASCII
// values in $c000 and $c010. Keypresses from the previous round have to be
// taken into account in order to detect if a single new key has been pressed.
func Poll() {
	allKeysPressed := make(map[uint8]bool)
	newKeysPressed := make(map[uint8]bool)

	// Query ebiten for all possible keys
	for k, v := range ebitenASCIIMap {
		if ebiten.IsKeyPressed(k) {
			allKeysPressed[v] = true

			_, present := previousKeysPressed[v]
			if !present {
				newKeysPressed[v] = true
			}
		}
	}

	previousKeysPressed = allKeysPressed

	if len(allKeysPressed) == 0 {
		// No keys are pressed, clear the strobe and return
		strobe = keyBoardData & 0x7f
		return
	} else if len(newKeysPressed) == 0 {
		// No new keys pressed, do nothing
		return
	} else if len(newKeysPressed) > 1 {
		// More than one new keys pressed, do nothing
		return
	}

	// Implicit else, one new key has been pressed

	// Get the key
	keys := []uint8{}
	for k := range newKeysPressed {
		keys = append(keys, k)
	}
	key := keys[0]

	if ebiten.IsKeyPressed(ebiten.KeyControl) && ebiten.IsKeyPressed(ebiten.KeyAlt) && key == 'c' {
		// Toggle capslock
		capsLock = !capsLock
	} else {
		// Normal case. Transform the ebiten key into ASCII

		shift := ebiten.IsKeyPressed(ebiten.KeyShift)
		shift = shift || (capsLock && key >= 'a' && key <= 'z')
		if shift {
			shiftedKey, present := shiftMap[key]
			if present {
				key = shiftedKey
			}
		}

		if ebiten.IsKeyPressed(ebiten.KeyControl) {
			controlKey, present := controlMap[key]
			if present {
				key = controlKey
			}
		}

		keyBoardData = key | 0x80
		strobe = keyBoardData
	}

	return
}

// Read returns the data and strobe values from set from the Poll() call
func Read() (uint8, uint8) {
	return keyBoardData, strobe
}

// ResetStrobe clears the high bit in keyboardData
func ResetStrobe() {
	keyBoardData &= 0x7f
}

package keyboard

import (
	"github.com/hajimehoshi/ebiten"
)

var ebitenAsciiMap map[ebiten.Key]uint8
var shiftMap map[uint8]uint8
var controlMap map[uint8]uint8

var keyBoardData uint8
var strobe uint8
var previousKeysPressed map[uint8]bool

func Init() {
	keyBoardData = 0
	strobe = 0

	ebitenAsciiMap = make(map[ebiten.Key]uint8)
	shiftMap = make(map[uint8]uint8)
	controlMap = make(map[uint8]uint8)
	previousKeysPressed = make(map[uint8]bool)

	ebitenAsciiMap[ebiten.KeyLeft] = 8
	ebitenAsciiMap[ebiten.KeyTab] = 9
	ebitenAsciiMap[ebiten.KeyDown] = 10
	ebitenAsciiMap[ebiten.KeyUp] = 11
	ebitenAsciiMap[ebiten.KeyEnter] = 13
	ebitenAsciiMap[ebiten.KeyRight] = 21
	ebitenAsciiMap[ebiten.KeyEscape] = 27
	ebitenAsciiMap[ebiten.KeyDelete] = 127

	ebitenAsciiMap[ebiten.Key0] = '0'
	ebitenAsciiMap[ebiten.Key1] = '1'
	ebitenAsciiMap[ebiten.Key2] = '2'
	ebitenAsciiMap[ebiten.Key3] = '3'
	ebitenAsciiMap[ebiten.Key4] = '4'
	ebitenAsciiMap[ebiten.Key5] = '5'
	ebitenAsciiMap[ebiten.Key6] = '6'
	ebitenAsciiMap[ebiten.Key7] = '7'
	ebitenAsciiMap[ebiten.Key8] = '8'
	ebitenAsciiMap[ebiten.Key9] = '9'
	ebitenAsciiMap[ebiten.KeyA] = 'a'
	ebitenAsciiMap[ebiten.KeyB] = 'b'
	ebitenAsciiMap[ebiten.KeyC] = 'c'
	ebitenAsciiMap[ebiten.KeyD] = 'd'
	ebitenAsciiMap[ebiten.KeyE] = 'e'
	ebitenAsciiMap[ebiten.KeyF] = 'f'
	ebitenAsciiMap[ebiten.KeyG] = 'g'
	ebitenAsciiMap[ebiten.KeyH] = 'h'
	ebitenAsciiMap[ebiten.KeyI] = 'i'
	ebitenAsciiMap[ebiten.KeyJ] = 'j'
	ebitenAsciiMap[ebiten.KeyK] = 'k'
	ebitenAsciiMap[ebiten.KeyL] = 'l'
	ebitenAsciiMap[ebiten.KeyM] = 'm'
	ebitenAsciiMap[ebiten.KeyN] = 'n'
	ebitenAsciiMap[ebiten.KeyO] = 'o'
	ebitenAsciiMap[ebiten.KeyP] = 'p'
	ebitenAsciiMap[ebiten.KeyQ] = 'q'
	ebitenAsciiMap[ebiten.KeyR] = 'r'
	ebitenAsciiMap[ebiten.KeyS] = 's'
	ebitenAsciiMap[ebiten.KeyT] = 't'
	ebitenAsciiMap[ebiten.KeyU] = 'u'
	ebitenAsciiMap[ebiten.KeyV] = 'v'
	ebitenAsciiMap[ebiten.KeyW] = 'w'
	ebitenAsciiMap[ebiten.KeyX] = 'x'
	ebitenAsciiMap[ebiten.KeyY] = 'y'
	ebitenAsciiMap[ebiten.KeyZ] = 'z'
	ebitenAsciiMap[ebiten.KeyApostrophe] = '\''
	ebitenAsciiMap[ebiten.KeyBackslash] = '\\'
	ebitenAsciiMap[ebiten.KeyComma] = ','
	ebitenAsciiMap[ebiten.KeyEqual] = '='
	ebitenAsciiMap[ebiten.KeyGraveAccent] = '`'
	ebitenAsciiMap[ebiten.KeyLeftBracket] = '['
	ebitenAsciiMap[ebiten.KeyMinus] = '-'
	ebitenAsciiMap[ebiten.KeyPeriod] = '.'
	ebitenAsciiMap[ebiten.KeyRightBracket] = ']'
	ebitenAsciiMap[ebiten.KeySemicolon] = ';'
	ebitenAsciiMap[ebiten.KeySlash] = '/'
	ebitenAsciiMap[ebiten.KeySpace] = ' '

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

// Poll queries ebiten's keyboard state and transforms that into apple //e
// values in $c000 and $c010
func Poll() {
	allKeysPressed := make(map[uint8]bool)
	newKeysPressed := make(map[uint8]bool)

	for k, v := range ebitenAsciiMap {
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
	keys := []uint8{}
	for k := range newKeysPressed {
		keys = append(keys, k)
	}
	key := keys[0]

	if ebiten.IsKeyPressed(ebiten.KeyShift) {
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

	return
}

func Read() (uint8, uint8) {
	return keyBoardData, strobe
}

func ResetStrobe() {
	keyBoardData &= 0x7f
}

package video

import (
	"github.com/hajimehoshi/ebiten"
)

const charMapASCIIArt = `
0x00
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x01
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x02
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x03
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x04
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x05
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x06
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x07
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x08
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x09
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x0a
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x0b
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x0c
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x0d
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x0e
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x0f
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x10
---------
|    X  |
|   X   |
| XX XX |
|X     X|
|X    X |
|X    X |
| X X  X|
| XX XX |
---------
0x11
---------
|    XXX|
|     XX|
| XXXXXX|
|X   XX |
|X  XXXX|
|    XX |
|XXXXXX |
| X     |
---------
0x12
---------
|       |
|   XX  |
|XXX    |
|       |
|XXX    |
|  XX   |
|   X   |
|    XXX|
---------
0x13
---------
|   X   |
|       |
|  XXX  |
| X   X |
| XXXXX |
| X   X |
| X   X |
|       |
---------
0x14
---------
|   X   |
|       |
|  XXX  |
|     X |
|  XXXX |
| X   X |
|  XXXX |
|       |
---------
0x15
---------
|       |
|  XXXX |
| X     |
| X     |
|  XXXX |
|    X  |
|   X   |
|       |
---------
0x16
---------
|       |
|  X    |
|   X   |
| X   X |
| X   X |
| X   X |
| X  XX |
|  XX X |
---------
0x17
---------
|XXXXXXX|
|XXXXXXX|
|XXXXX X|
|  XX XX|
|X X XXX|
|XX XXXX|
|XX XXXX|
|XXXXXXX|
---------
0x18
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x19
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x1a
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x1b
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x1c
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x1d
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x1e
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x1f
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x20
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x21
---------
|   X   |
|   X   |
|   X   |
|   X   |
|   X   |
|       |
|   X   |
|       |
---------
0x22
---------
|  X X  |
|  X X  |
|  X X  |
|       |
|       |
|       |
|       |
|       |
---------
0x23
---------
|  X X  |
|  X X  |
| XXXXX |
|  X X  |
| XXXXX |
|  X X  |
|  X X  |
|       |
---------
0x24
---------
|   X   |
|  XXXX |
| X X   |
|  XXX  |
|   X X |
| XXXX  |
|   X   |
|       |
---------
0x25
---------
| XX    |
| XX  X |
|    X  |
|   X   |
|  X    |
| X  XX |
|    XX |
|       |
---------
0x26
---------
|  X    |
| X X   |
| X X   |
|  X    |
| X X X |
| X  X  |
|  XX X |
|       |
---------
0x27
---------
|   X   |
|   X   |
|   X   |
|       |
|       |
|       |
|       |
|       |
---------
0x28
---------
|   X   |
|  X    |
| X     |
| X     |
| X     |
|  X    |
|   X   |
|       |
---------
0x29
---------
|   X   |
|    X  |
|     X |
|     X |
|     X |
|    X  |
|   X   |
|       |
---------
0x2a
---------
|   X   |
| X X X |
|  XXX  |
|   X   |
|  XXX  |
| X X X |
|   X   |
|       |
---------
0x2b
---------
|       |
|   X   |
|   X   |
| XXXXX |
|   X   |
|   X   |
|       |
|       |
---------
0x2c
---------
|       |
|       |
|       |
|       |
|   X   |
|   X   |
|  X    |
|       |
---------
0x2d
---------
|       |
|       |
|       |
| XXXXX |
|       |
|       |
|       |
|       |
---------
0x2e
---------
|       |
|       |
|       |
|       |
|       |
|       |
|   X   |
|       |
---------
0x2f
---------
|       |
|     X |
|    X  |
|   X   |
|  X    |
| X     |
|       |
|       |
---------
0x30
---------
|  XXX  |
| X   X |
| X  XX |
| X X X |
| XX  X |
| X   X |
|  XXX  |
|       |
---------
0x31
---------
|   X   |
|  XX   |
|   X   |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0x32
---------
|  XXX  |
| X   X |
|     X |
|   XX  |
|  X    |
| X     |
| XXXXX |
|       |
---------
0x33
---------
| XXXXX |
|     X |
|    X  |
|   XX  |
|     X |
| X   X |
|  XXX  |
|       |
---------
0x34
---------
|    X  |
|   XX  |
|  X X  |
| X  X  |
| XXXXX |
|    X  |
|    X  |
|       |
---------
0x35
---------
| XXXXX |
| X     |
| XXXX  |
|     X |
|     X |
| X   X |
|  XXX  |
|       |
---------
0x36
---------
|   XXX |
|  X    |
| X     |
| XXXX  |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0x37
---------
| XXXXX |
|     X |
|    X  |
|   X   |
|  X    |
|  X    |
|  X    |
|       |
---------
0x38
---------
|  XXX  |
| X   X |
| X   X |
|  XXX  |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0x39
---------
|  XXX  |
| X   X |
| X   X |
|  XXXX |
|     X |
|    X  |
| XXX   |
|       |
---------
0x3a
---------
|       |
|       |
|   X   |
|       |
|   X   |
|       |
|       |
|       |
---------
0x3b
---------
|       |
|       |
|   X   |
|       |
|   X   |
|   X   |
|  X    |
|       |
---------
0x3c
---------
|    X  |
|   X   |
|  X    |
| X     |
|  X    |
|   X   |
|    X  |
|       |
---------
0x3d
---------
|       |
|       |
| XXXXX |
|       |
| XXXXX |
|       |
|       |
|       |
---------
0x3e
---------
|  X    |
|   X   |
|    X  |
|     X |
|    X  |
|   X   |
|  X    |
|       |
---------
0x3f
---------
|  XXX  |
| X   X |
|    X  |
|   X   |
|   X   |
|       |
|   X   |
|       |
---------
0x40
---------
|  XXX  |
| X   X |
| X X X |
| X XXX |
| X XX  |
| X     |
|  XXXX |
|       |
---------
0x41
---------
|   X   |
|  X X  |
| X   X |
| X   X |
| XXXXX |
| X   X |
| X   X |
|       |
---------
0x42
---------
| XXXX  |
| X   X |
| X   X |
| XXXX  |
| X   X |
| X   X |
| XXXX  |
|       |
---------
0x43
---------
|  XXX  |
| X   X |
| X     |
| X     |
| X     |
| X   X |
|  XXX  |
|       |
---------
0x44
---------
| XXXX  |
| X   X |
| X   X |
| X   X |
| X   X |
| X   X |
| XXXX  |
|       |
---------
0x45
---------
| XXXXX |
| X     |
| X     |
| XXXX  |
| X     |
| X     |
| XXXXX |
|       |
---------
0x46
---------
| XXXXX |
| X     |
| X     |
| XXXX  |
| X     |
| X     |
| X     |
|       |
---------
0x47
---------
|  XXXX |
| X     |
| X     |
| X     |
| X  XX |
| X   X |
|  XXXX |
|       |
---------
0x48
---------
| X   X |
| X   X |
| X   X |
| XXXXX |
| X   X |
| X   X |
| X   X |
|       |
---------
0x49
---------
|  XXX  |
|   X   |
|   X   |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0x4a
---------
|     X |
|     X |
|     X |
|     X |
|     X |
| X   X |
|  XXX  |
|       |
---------
0x4b
---------
| X   X |
| X  X  |
| X X   |
| XX    |
| X X   |
| X  X  |
| X   X |
|       |
---------
0x4c
---------
| X     |
| X     |
| X     |
| X     |
| X     |
| X     |
| XXXXX |
|       |
---------
0x4d
---------
| X   X |
| XX XX |
| X X X |
| X X X |
| X   X |
| X   X |
| X   X |
|       |
---------
0x4e
---------
| X   X |
| X   X |
| XX  X |
| X X X |
| X  XX |
| X   X |
| X   X |
|       |
---------
0x4f
---------
|  XXX  |
| X   X |
| X   X |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0x50
---------
| XXXX  |
| X   X |
| X   X |
| XXXX  |
| X     |
| X     |
| X     |
|       |
---------
0x51
---------
|  XXX  |
| X   X |
| X   X |
| X   X |
| X X X |
| X  X  |
|  XX X |
|       |
---------
0x52
---------
| XXXX  |
| X   X |
| X   X |
| XXXX  |
| X X   |
| X  X  |
| X   X |
|       |
---------
0x53
---------
|  XXX  |
| X   X |
| X     |
|  XXX  |
|     X |
| X   X |
|  XXX  |
|       |
---------
0x54
---------
| XXXXX |
|   X   |
|   X   |
|   X   |
|   X   |
|   X   |
|   X   |
|       |
---------
0x55
---------
| X   X |
| X   X |
| X   X |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0x56
---------
| X   X |
| X   X |
| X   X |
| X   X |
| X   X |
|  X X  |
|   X   |
|       |
---------
0x57
---------
| X   X |
| X   X |
| X   X |
| X X X |
| X X X |
| XX XX |
| X   X |
|       |
---------
0x58
---------
| X   X |
| X   X |
|  X X  |
|   X   |
|  X X  |
| X   X |
| X   X |
|       |
---------
0x59
---------
| X   X |
| X   X |
|  X X  |
|   X   |
|   X   |
|   X   |
|   X   |
|       |
---------
0x5a
---------
| XXXXX |
|     X |
|    X  |
|   X   |
|  X    |
| X     |
| XXXXX |
|       |
---------
0x5b
---------
| XXXXX |
| XX    |
| XX    |
| XX    |
| XX    |
| XX    |
| XXXXX |
|       |
---------
0x5c
---------
|       |
| X     |
|  X    |
|   X   |
|    X  |
|     X |
|       |
|       |
---------
0x5d
---------
| XXXXX |
|    XX |
|    XX |
|    XX |
|    XX |
|    XX |
| XXXXX |
|       |
---------
0x5e
---------
|       |
|       |
|   X   |
|  X X  |
| X   X |
|       |
|       |
|       |
---------
0x5f
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|XXXXXXX|
---------
0x60
---------
|  X    |
|   X   |
|    X  |
|       |
|       |
|       |
|       |
|       |
---------
0x61
---------
|       |
|       |
|  XXX  |
|     X |
|  XXXX |
| X   X |
|  XXXX |
|       |
---------
0x62
---------
| X     |
| X     |
| XXXX  |
| X   X |
| X   X |
| X   X |
| XXXX  |
|       |
---------
0x63
---------
|       |
|       |
|  XXXX |
| X     |
| X     |
| X     |
|  XXXX |
|       |
---------
0x64
---------
|     X |
|     X |
|  XXXX |
| X   X |
| X   X |
| X   X |
|  XXXX |
|       |
---------
0x65
---------
|       |
|       |
|  XXX  |
| X   X |
| XXXXX |
| X     |
|  XXXX |
|       |
---------
0x66
---------
|   XX  |
|  X  X |
|  X    |
| XXXX  |
|  X    |
|  X    |
|  X    |
|       |
---------
0x67
---------
|       |
|       |
|  XXX  |
| X   X |
| X   X |
|  XXXX |
|     X |
|  XXX  |
---------
0x68
---------
| X     |
| X     |
| XXXX  |
| X   X |
| X   X |
| X   X |
| X   X |
|       |
---------
0x69
---------
|   X   |
|       |
|  XX   |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0x6a
---------
|    X  |
|       |
|   XX  |
|    X  |
|    X  |
|    X  |
| X  X  |
|  XX   |
---------
0x6b
---------
| X     |
| X     |
| X   X |
| X  X  |
| XXX   |
| X  X  |
| X   X |
|       |
---------
0x6c
---------
|  XX   |
|   X   |
|   X   |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0x6d
---------
|       |
|       |
| XX XX |
| X X X |
| X X X |
| X X X |
| X   X |
|       |
---------
0x6e
---------
|       |
|       |
| XXXX  |
| X   X |
| X   X |
| X   X |
| X   X |
|       |
---------
0x6f
---------
|       |
|       |
|  XXX  |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0x70
---------
|       |
|       |
| XXXX  |
| X   X |
| X   X |
| XXXX  |
| X     |
| X     |
---------
0x71
---------
|       |
|       |
|  XXXX |
| X   X |
| X   X |
|  XXXX |
|     X |
|     X |
---------
0x72
---------
|       |
|       |
| X XXX |
| XX    |
| X     |
| X     |
| X     |
|       |
---------
0x73
---------
|       |
|       |
|  XXXX |
| X     |
|  XXX  |
|     X |
| XXXX  |
|       |
---------
0x74
---------
|  X    |
|  X    |
| XXXX  |
|  X    |
|  X    |
|  X  X |
|   XX  |
|       |
---------
0x75
---------
|       |
|       |
| X   X |
| X   X |
| X   X |
| X  XX |
|  XX X |
|       |
---------
0x76
---------
|       |
|       |
| X   X |
| X   X |
| X   X |
|  X X  |
|   X   |
|       |
---------
0x77
---------
|       |
|       |
| X   X |
| X   X |
| X X X |
| X X X |
| XX XX |
|       |
---------
0x78
---------
|       |
|       |
| X   X |
|  X X  |
|   X   |
|  X X  |
| X   X |
|       |
---------
0x79
---------
|       |
|       |
| X   X |
| X   X |
| X   X |
|  XXXX |
|     X |
|  XXX  |
---------
0x7a
---------
|       |
|       |
| XXXXX |
|    X  |
|   X   |
|  X    |
| XXXXX |
|       |
---------
0x7b
---------
|   XXX |
|  XX   |
|  XX   |
| XX    |
|  XX   |
|  XX   |
|   XXX |
|       |
---------
0x7c
---------
|   X   |
|   X   |
|   X   |
|   X   |
|   X   |
|   X   |
|   X   |
|   X   |
---------
0x7d
---------
| XXX   |
|   XX  |
|   XX  |
|    XX |
|   XX  |
|   XX  |
| XXX   |
|       |
---------
0x7e
---------
|  XX X |
| X XX  |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x7f
---------
|       |
| X X X |
|  X X  |
| X X X |
|  X X  |
| X X X |
|       |
|       |
---------
0x80
---------
|    X  |
|   X   |
| XX XX |
|XXXXXXX|
|XXXXXX |
|XXXXXX |
| XXXXXX|
| XX XX |
---------
0x81
---------
|    X  |
|   X   |
| XX XX |
|X     X|
|X    X |
|X    X |
| X X  X|
| XX XX |
---------
0x82
---------
|       |
|       |
| X     |
| XX    |
| XXX   |
| XXXX  |
| XX XX |
| X    X|
---------
0x83
---------
|XXXXXXX|
| X   X |
|  X X  |
|   X   |
|   X   |
|  X X  |
| X X X |
|XXXXXXX|
---------
0x84
---------
|       |
|      X|
|     X |
|X   X  |
| X X   |
|  X    |
|  X    |
|       |
---------
0x85
---------
|XXXXXXX|
|XXXXXX |
|XXXXX X|
|  XX XX|
|X X XXX|
|XX XXXX|
|XX XXXX|
|XXXXXXX|
---------
0x86
---------
|XXXXXX |
|XXXXXX |
|XXXXXX |
|XX XXX |
|X  XXX |
|       |
|X  XXXX|
|XX XXXX|
---------
0x87
---------
|XXXXXXX|
|       |
|XXXXXXX|
|       |
|XXXXXXX|
|       |
|       |
|XXXXXXX|
---------
0x88
---------
|   X   |
|  X    |
| X     |
|XXXXXXX|
| X     |
|  X    |
|   X   |
|       |
---------
0x89
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
| X X X |
---------
0x8a
---------
|   X   |
|   X   |
|   X   |
|   X   |
|X  X  X|
| X X X |
|  XXX  |
|   X   |
---------
0x8b
---------
|   X   |
|  XXX  |
| X X X |
|X  X  X|
|   X   |
|   X   |
|   X   |
|   X   |
---------
0x8c
---------
|XXXXXXX|
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0x8d
---------
|      X|
|      X|
|      X|
|  X   X|
| XX   X|
|XXXXXXX|
| XX    |
|  X    |
---------
0x8e
---------
|XXXXXX |
|XXXXXX |
|XXXXXX |
|XXXXXX |
|XXXXXX |
|XXXXXX |
|XXXXXX |
|XXXXXX |
---------
0x8f
---------
|XX  X  |
|   XX  |
|  XXX  |
| XXXXXX|
|  XXX  |
|   XX  |
|    X  |
|XXXX XX|
---------
0x90
---------
|  X  XX|
|  XX   |
|  XXX  |
|XXXXXX |
|  XXX  |
|  XX   |
|  X    |
|XX XXXX|
---------
0x91
---------
|      X|
|   X  X|
|   X   |
|XXXXXXX|
| XXXXX |
|  XXX  |
|   X  X|
|      X|
---------
0x92
---------
|      X|
|   X  X|
|  XXX  |
| XXXXX |
|XXXXXXX|
|   X   |
|   X  X|
|      X|
---------
0x93
---------
|       |
|       |
|       |
|XXXXXXX|
|       |
|       |
|       |
|       |
---------
0x94
---------
|X      |
|X      |
|X      |
|X      |
|X      |
|X      |
|X      |
|XXXXXXX|
---------
0x95
---------
|   X   |
|    X  |
|     X |
|XXXXXXX|
|     X |
|    X  |
|   X   |
|       |
---------
0x96
---------
| X X X |
|X X X X|
| X X X |
|X X X X|
| X X X |
|X X X X|
| X X X |
|X X X X|
---------
0x97
---------
|X X X X|
| X X X |
|X X X X|
| X X X |
|X X X X|
| X X X |
|X X X X|
| X X X |
---------
0x98
---------
|       |
| XXXXX |
|X     X|
|X      |
|X      |
|X      |
|XXXXXXX|
|       |
---------
0x99
---------
|       |
|       |
|XXXXXX |
|      X|
|      X|
|      X|
|XXXXXXX|
|       |
---------
0x9a
---------
|      X|
|      X|
|      X|
|      X|
|      X|
|      X|
|      X|
|      X|
---------
0x9b
---------
|   X   |
|  XXX  |
| XXXXX |
|XXXXXXX|
| XXXXX |
|  XXX  |
|   X   |
|       |
---------
0x9c
---------
|XXXXXXX|
|       |
|       |
|       |
|       |
|       |
|       |
|XXXXXXX|
---------
0x9d
---------
|  X X  |
|  X X  |
|XXX XXX|
|       |
|XXX XXX|
|  X X  |
|  X X  |
|       |
---------
0x9e
---------
|XXXXXXX|
|      X|
|      X|
|  XX  X|
|  XX  X|
|      X|
|      X|
|XXXXXXX|
---------
0x9f
---------
|X      |
|X      |
|X      |
|X      |
|X      |
|X      |
|X      |
|X      |
---------
0xa0
---------
|       |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0xa1
---------
|   X   |
|       |
|   X   |
|   X   |
|   X   |
|   X   |
|   X   |
|       |
---------
0xa2
---------
|   X   |
|  XXXX |
| X X   |
| X X   |
| X X   |
|  XXXX |
|   X   |
|       |
---------
0xa3
---------
|   XXX |
|  X   X|
|  X    |
| XXX   |
|  X    |
|  X    |
| X XXXX|
|       |
---------
0xa4
---------
|       |
| X   X |
|  XXX  |
|  X X  |
|  XXX  |
| X   X |
|       |
|       |
---------
0xa5
---------
| X   X |
| X   X |
|  X X  |
|   X   |
| XXXXX |
|   X   |
|   X   |
|       |
---------
0xa6
---------
|   X   |
|   X   |
|   X   |
|       |
|       |
|   X   |
|   X   |
|   X   |
---------
0xa7
---------
|  XXXX |
| X     |
|  XXX  |
| X   X |
|  XXX  |
|     X |
| XXXX  |
|       |
---------
0xa8
---------
| X   X |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0xa9
---------
|  XXX  |
| X   X |
|X  XX X|
|X X   X|
|X  XX X|
| X   X |
|  XXX  |
|       |
---------
0xaa
---------
|  XX   |
|   XX  |
|  X X  |
|   XX  |
|       |
|  XXX  |
|       |
|       |
---------
0xab
---------
|   X  X|
|  X  X |
| X  X  |
|X  X   |
| X  X  |
|  X  X |
|   X  X|
|       |
---------
0xac
---------
|       |
|       |
| XXXXX |
|     X |
|     X |
|       |
|       |
|       |
---------
0xad
---------
|       |
|       |
|       |
| XXXXX |
|       |
|       |
|       |
|       |
---------
0xae
---------
|  XXX  |
| X   X |
|X XXX X|
|X XX  X|
|X X X X|
| X   X |
|  XXX  |
|       |
---------
0xaf
---------
| XXXXX |
|       |
|       |
|       |
|       |
|       |
|       |
|       |
---------
0xb0
---------
|   X   |
|  X X  |
|   X   |
|       |
|       |
|       |
|       |
|       |
---------
0xb1
---------
|   X   |
|   X   |
| XXXXX |
|   X   |
|   X   |
|       |
| XXXXX |
|       |
---------
0xb2
---------
|  XX   |
|    X  |
|   X   |
|  XXX  |
|       |
|       |
|       |
|       |
---------
0xb3
---------
|  XXX  |
|   X   |
|    X  |
|  XX   |
|       |
|       |
|       |
|       |
---------
0xb4
---------
|    X  |
|   X   |
|  X    |
|       |
|       |
|       |
|       |
|       |
---------
0xb5
---------
|       |
|       |
| X   X |
| X   X |
| X   X |
| XX XX |
| X X X |
| X     |
---------
0xb6
---------
|  XXXX |
| X X X |
| X X X |
|  XXXX |
|   X X |
|   X X |
|   X X |
|       |
---------
0xb7
---------
|       |
|       |
|       |
|   X   |
|       |
|       |
|       |
|       |
---------
0xb8
---------
|       |
|       |
|       |
|       |
|       |
|       |
|    X  |
|   X   |
---------
0xb9
---------
|   X   |
|  XX   |
|   X   |
|  XXX  |
|       |
|       |
|       |
|       |
---------
0xba
---------
|   X   |
|  X X  |
|  X X  |
|   X   |
|       |
|  XXX  |
|       |
|       |
---------
0xbb
---------
|X  X   |
| X  X  |
|  X  X |
|   X  X|
|  X  X |
| X  X  |
|X  X   |
|       |
---------
0xbc
---------
| X     |
| X   X |
| X  X  |
|   X   |
|  X X  |
| X  XX |
|     X |
|       |
---------
0xbd
---------
| X     |
| X   X |
| X  X  |
|   X   |
|  X XX |
| X  X  |
|    XX |
|       |
---------
0xbe
---------
| XX    |
|  X  X |
| XX X  |
|   X   |
|  X X  |
| X  XX |
|     X |
|       |
---------
0xbf
---------
|   X   |
|       |
|   X   |
|   X   |
|  X    |
| X   X |
|  XXX  |
|       |
---------
0xc0
---------
|   X   |
|    X  |
|   X   |
|  X X  |
| X   X |
| XXXXX |
| X   X |
|       |
---------
0xc1
---------
|   X   |
|  X    |
|   X   |
|  X X  |
| X   X |
| XXXXX |
| X   X |
|       |
---------
0xc2
---------
|   X   |
|  X X  |
|   X   |
|  X X  |
| X   X |
| XXXXX |
| X   X |
|       |
---------
0xc3
---------
|  XX X |
| X XX  |
|   X   |
|  X X  |
| X   X |
| XXXXX |
| X   X |
|       |
---------
0xc4
---------
| X   X |
|       |
|   X   |
|  X X  |
| X   X |
| XXXXX |
| X   X |
|       |
---------
0xc5
---------
|  XXX  |
| X   X |
|  XXX  |
| X   X |
| XXXXX |
| X   X |
| X   X |
|       |
---------
0xc6
---------
|  XXXX |
| X X   |
| X X   |
| XXXXX |
| X X   |
| X X   |
| X XXX |
|       |
---------
0xc7
---------
|  XXX  |
| X   X |
| X     |
| X     |
| X   X |
|  XXX  |
|    X  |
|   X   |
---------
0xc8
---------
|    X  |
|     X |
| XXXXX |
| X     |
| XXXX  |
| X     |
| XXXXX |
|       |
---------
0xc9
---------
|     X |
|    X  |
| XXXXX |
| X     |
| XXXX  |
| X     |
| XXXXX |
|       |
---------
0xca
---------
|    X  |
|   X X |
| XXXXX |
| X     |
| XXXX  |
| X     |
| XXXXX |
|       |
---------
0xcb
---------
| X   X |
|       |
| XXXXX |
| X     |
| XXXX  |
| X     |
| XXXXX |
|       |
---------
0xcc
---------
|  X    |
|   X   |
|  XXX  |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0xcd
---------
|    X  |
|   X   |
|  XXX  |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0xce
---------
|   X   |
|  X X  |
|  XXX  |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0xcf
---------
| X   X |
|       |
|  XXX  |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0xd0
---------
|  XXX  |
|  X  X |
|  X  X |
| XXX X |
|  X  X |
|  X  X |
|  XXX  |
|       |
---------
0xd1
---------
|  XX X |
| X XX  |
|       |
| XX  X |
| X X X |
| X  XX |
| X   X |
|       |
---------
0xd2
---------
|   X   |
|    X  |
|  XXX  |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xd3
---------
|     X |
|    X  |
|  XXX  |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xd4
---------
|    X  |
|   X X |
|  XXX  |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xd5
---------
|  XX X |
| X XX  |
|  XXX  |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xd6
---------
| X   X |
|  XXX  |
| X   X |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xd7
---------
|       |
| X   X |
|  X X  |
|   X   |
|  X X  |
| X   X |
|       |
|       |
---------
0xd8
---------
|     X |
|  XXX  |
| X  XX |
| X X X |
| X X X |
| XX  X |
|  XXX  |
| X     |
---------
0xd9
---------
|  X    |
|   X   |
| X   X |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xda
---------
|    X  |
|   X   |
| X   X |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xdb
---------
|   X   |
|  X X  |
| X   X |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xdc
---------
| X   X |
|       |
| X   X |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xdd
---------
|    X  |
|   X   |
| X   X |
| X   X |
|  X X  |
|   X   |
|   X   |
|       |
---------
0xde
---------
| X     |
| X     |
| XXXX  |
| X   X |
| XXXX  |
| X     |
| X     |
|       |
---------
0xdf
---------
| XXXX  |
| X   X |
| X   X |
| XXXX  |
| X   X |
| X   X |
| X XX  |
| X     |
---------
0xe0
---------
|    X  |
|     X |
|  XXX  |
|     X |
|  XXXX |
| X   X |
|  XXXX |
|       |
---------
0xe1
---------
|     X |
|    X  |
|  XXX  |
|     X |
|  XXXX |
| X   X |
|  XXXX |
|       |
---------
0xe2
---------
|    X  |
|   X X |
|  XXX  |
|     X |
|  XXXX |
| X   X |
|  XXXX |
|       |
---------
0xe3
---------
|  XX X |
| X XX  |
|  XXX  |
|     X |
|  XXXX |
| X   X |
|  XXXX |
|       |
---------
0xe4
---------
|       |
| X   X |
|  XXX  |
|     X |
|  XXXX |
| X   X |
|  XXXX |
|       |
---------
0xe5
---------
|   X   |
|  X X  |
|  XXX  |
|     X |
|  XXXX |
| X   X |
|  XXXX |
|       |
---------
0xe6
---------
|       |
|       |
| XX X  |
|   X X |
|  XXXX |
| X X   |
|  XXXX |
|       |
---------
0xe7
---------
|       |
|       |
|  XXXX |
| X     |
| X     |
|  XXXX |
|    X  |
|   X   |
---------
0xe8
---------
|    X  |
|     X |
|  XXX  |
| X   X |
| XXXXX |
| X     |
|  XXXX |
|       |
---------
0xe9
---------
|     X |
|    X  |
|  XXX  |
| X   X |
| XXXXX |
| X     |
|  XXXX |
|       |
---------
0xea
---------
|    X  |
|   X X |
|  XXX  |
| X   X |
| XXXXX |
| X     |
|  XXXX |
|       |
---------
0xeb
---------
|       |
| X   X |
|  XXX  |
| X   X |
| XXXXX |
| X     |
|  XXXX |
|       |
---------
0xec
---------
|  X    |
|   X   |
|       |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0xed
---------
|    X  |
|   X   |
|       |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0xee
---------
|   X   |
|  X X  |
|       |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0xef
---------
|       |
| X   X |
|       |
|   X   |
|   X   |
|   X   |
|  XXX  |
|       |
---------
0xf0
---------
|  X X  |
|   X   |
|  X X  |
|  XXX  |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xf1
---------
|  XX X |
| X XX  |
|       |
| XXXX  |
| X   X |
| X   X |
| X   X |
|       |
---------
0xf2
---------
|       |
|   X   |
|    X  |
|  XXX  |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xf3
---------
|       |
|     X |
|    X  |
|  XXX  |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xf4
---------
|       |
|    X  |
|   X X |
|  XXX  |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xf5
---------
|  XX X |
| X XX  |
|       |
|  XXX  |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xf6
---------
|       |
| X   X |
|  XXX  |
| X   X |
| X   X |
| X   X |
|  XXX  |
|       |
---------
0xf7
---------
|       |
|   X   |
|       |
| XXXXX |
|       |
|   X   |
|       |
|       |
---------
0xf8
---------
|       |
|     X |
|  XXX  |
| X  XX |
| X X X |
| XX  X |
|  XXX  |
| X     |
---------
0xf9
---------
|  X    |
|   X   |
| X   X |
| X   X |
| X   X |
| X  XX |
|  XX X |
|       |
---------
0xfa
---------
|    X  |
|   X   |
| X   X |
| X   X |
| X   X |
| X  XX |
|  XX X |
|       |
---------
0xfb
---------
|   X   |
|  X X  |
| X   X |
| X   X |
| X   X |
| X  XX |
|  XX X |
|       |
---------
0xfc
---------
|       |
| X   X |
|       |
| X   X |
| X   X |
| X  XX |
|  XX X |
|       |
---------
0xfd
---------
|    X  |
|   X   |
| X   X |
| X   X |
| X   X |
|  XXXX |
|     X |
|  XXX  |
---------
0xfe
---------
| X     |
| X     |
| XXXX  |
| X   X |
| X   X |
| XXXX  |
| X     |
| X     |
---------
0xff
---------
|       |
| X   X |
|       |
| X   X |
| X   X |
|  XXXX |
|     X |
|  XXX  |
---------
`

var charMap [0x100]*ebiten.Image

// initTextCharMap initializes the text character map
func initTextCharMap() {
	for c := 0; c < 0x100; c++ {
		image, err := ebiten.NewImage(7, 8, ebiten.FilterNearest)
		if err != nil {
			panic(err)
		}

		pixels := make([]byte, 8*7*4)

		start := c*105 + 17

		for y := 0; y < 8; y++ {
			for x := 0; x < 7; x++ {
				var b float64
				set := charMapASCIIArt[start+10*y+x] == 'X'
				if set {
					b = 1
				} else {
					b = 0
				}

				p := 4 * (7*y + x)
				pixels[p+0] = byte(0xff * b)
				pixels[p+1] = byte(0xff * b)
				pixels[p+2] = byte(0xff * b)
				pixels[p+3] = 0xff

			}
		}
		image.ReplacePixels(pixels)
		charMap[c] = image
	}
}

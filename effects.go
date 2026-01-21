package elklib

// Common effect codes for ELK devices
const (
	EffectJumpRGB           uint8 = 0x87
	EffectJumpRGBYCMW       uint8 = 0x88
	EffectCrossfadeRGB      uint8 = 0x8b
	EffectCrossfadeRGBYCMW  uint8 = 0x8c
	EffectBlinkRGB          uint8 = 0x8f
	EffectBlinkRGBYCMW      uint8 = 0x90
	EffectCrossfadeRed      uint8 = 0x96
	EffectCrossfadeGreen    uint8 = 0x97
	EffectCrossfadeBlue     uint8 = 0x98
	EffectCrossfadeYellow   uint8 = 0x99
	EffectCrossfadeCyan     uint8 = 0x9a
	EffectCrossfadeMagenta  uint8 = 0x9b
	EffectCrossfadeWhite    uint8 = 0x9c
	EffectCrossfadeRedGreen uint8 = 0x9d
	EffectJumpRedGreenBlue  uint8 = 0xa0
	EffectBlinkRed          uint8 = 0xa3
	EffectBlinkGreen        uint8 = 0xa4
	EffectBlinkBlue         uint8 = 0xa5
	EffectBlinkYellow       uint8 = 0xa6
	EffectBlinkCyan         uint8 = 0xa7
	EffectBlinkMagenta      uint8 = 0xa8
	EffectBlinkWhite        uint8 = 0xa9
)

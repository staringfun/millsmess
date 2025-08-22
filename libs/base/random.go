// Mills Mess
// Licensed under the Mills Mess License Agreement
// See LICENSE.md in the root of this repository.

package base

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/staringfun/millsmess/libs/types"
	"math/big"
	mathRand "math/rand/v2"
)

func GenerateRandomStringRunes(length int, runes []rune) string {
	result := make([]rune, length)
	m := big.NewInt(int64(len(runes)))

	for i := range result {
		n, _ := rand.Int(rand.Reader, m)
		result[i] = runes[n.Int64()]
	}
	return string(result)
}

var URLSafeRunes = []rune{
	'a',
	'b',
	'c',
	'd',
	'e',
	'f',
	'g',
	'h',
	'i',
	'j',
	'k',
	'l',
	'm',
	'n',
	'o',
	'p',
	'q',
	'r',
	's',
	't',
	'u',
	'v',
	'w',
	'x',
	'y',
	'z',
	'A',
	'B',
	'C',
	'D',
	'E',
	'F',
	'G',
	'H',
	'I',
	'J',
	'K',
	'L',
	'M',
	'N',
	'O',
	'P',
	'Q',
	'R',
	'S',
	'T',
	'U',
	'V',
	'W',
	'X',
	'Y',
	'Z',
	'0',
	'1',
	'2',
	'3',
	'4',
	'5',
	'6',
	'7',
	'8',
	'9',
	'_',
	'-',
	'.',
}

func GenerateRandomString(length int) string {
	return GenerateRandomStringRunes(length, URLSafeRunes)
}

const PlayerIDLength = 8

func GeneratePlayerID() types.PlayerID {
	return types.PlayerID(GenerateRandomString(PlayerIDLength))
}

func GenerateRandomStringFast(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(mathRand.IntN(256))
	}
	return base64.RawURLEncoding.EncodeToString(b)[:n]
}

func GenerateTraceID() string {
	return GenerateRandomStringFast(8)
}

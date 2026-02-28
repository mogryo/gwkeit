package langdetector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsRuby(t *testing.T) {
	t.Run("should return true for ruby function declaration, no parentheses", func(t *testing.T) {
		src := "def main\nreturn 1\nend"
		assert.True(t, IsRuby(src))
	})
	t.Run("should return true for ruby function declaration, with parameters", func(t *testing.T) {
		src := "def isLarger?(a, b)\nreturn a>b\nend"
		assert.True(t, IsRuby(src))
	})
	t.Run("should return false for empty string", func(t *testing.T) {
		assert.False(t, IsRuby(""))
	})
	t.Run("should return false for javascript function declaration", func(t *testing.T) {
		src := "function main() {}\n"
		assert.False(t, IsRuby(src))
	})
	t.Run("should return false for golang loop", func(t *testing.T) {
		src := "{\nfor _, i := range entries {\n}\n}"
		assert.False(t, IsRuby(src))
	})
	t.Run("should return false for typescript function", func(t *testing.T) {
		src := "() => {\nlet pickedCard = Math.floor(Math.random() * 52);\nlet pickedSuit = Math.floor(pickedCard / 13);\n \nreturn { suit: this.suits[pickedSuit], card: pickedCard % 13 };\n};"
		assert.False(t, IsRuby(src))
	})
	t.Run("should return false for kotlin enum and function ", func(t *testing.T) {
		src := "enum class Bit {\nZERO, ONE\n}\n\nfun getRandomBit(): Bit {\nreturn if (Random.nextBoolean()) Bit.ONE else Bit.ZERO\n}"
		assert.False(t, IsRuby(src))
	})
}

package langdetector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsTypeScript(t *testing.T) {
	t.Run("should return true for typescript function", func(t *testing.T) {
		src := "function main(x: number, y: number): boolean {\nconst a=''\n}\n"
		assert.True(t, IsTypeScript(src))
	})
	t.Run("should return true for typescript function", func(t *testing.T) {
		src := "() => {\nlet pickedCard = Math.floor(Math.random() * 52);\nlet pickedSuit = Math.floor(pickedCard / 13);\n \nreturn { suit: this.suits[pickedSuit], card: pickedCard % 13 };\n};"
		assert.True(t, IsTypeScript(src))
	})
	t.Run("should return true for interface declaration", func(t *testing.T) {
		src := "interface Deck {\n  suits: string[];\n  cards: number[];\n  createCardPicker(this: Deck): () => Card;\n}"
		assert.True(t, IsTypeScript(src))
	})
	t.Run("should return false for empty string", func(t *testing.T) {
		assert.False(t, IsTypeScript(""))
	})
	t.Run("should return false for python function declaration", func(t *testing.T) {
		src := "def main():\nreturn 1\n"
		assert.False(t, IsTypeScript(src))
	})
	t.Run("should return false for SELECT specific columns  statement", func(t *testing.T) {
		src := "SELECT name, surname FROM users"
		assert.False(t, IsTypeScript(src))
	})
	t.Run("should return false for SELECT everything statement", func(t *testing.T) {
		src := "SELECT * FROM users"
		assert.False(t, IsTypeScript(src))
	})
}

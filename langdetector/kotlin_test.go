package langdetector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsKotlin(t *testing.T) {
	t.Run("should return true for kotlin function declaration", func(t *testing.T) {
		src := "fun main() { println(\"Hello, Kotlin!\") }"
		assert.True(t, IsKotlin(src))
	})
	t.Run("should return true for kotlin loop", func(t *testing.T) {
		src := "when (deliveryStatus) {\n\"Pending\" -> print(\"Your order is being prepared\")\n\"Shipped\" -> print(\"Your order is on the way\")\n}"
		assert.True(t, IsKotlin(src))
	})
	t.Run("should return true for kotlin enum and function ", func(t *testing.T) {
		src := "enum class Bit {\nZERO, ONE\n}\n\nfun getRandomBit(): Bit {\nreturn if (Random.nextBoolean()) Bit.ONE else Bit.ZERO\n}"
		assert.True(t, IsKotlin(src))
	})
	t.Run("should return false for empty string", func(t *testing.T) {
		assert.False(t, IsKotlin(""))
	})
	t.Run("should return false for javascript function declaration", func(t *testing.T) {
		src := "function main() {}\n"
		assert.False(t, IsKotlin(src))
	})
	t.Run("should return false for python condition", func(t *testing.T) {
		src := "if query is None or query.strip() == \"\":\nif save_query_name is not None and save_query_name.strip() != \"\":\nprint(\n\"You have provided the query name for saving, but query argument is empty!\"\n)\nreturn"
		assert.False(t, IsKotlin(src))
	})
}

package langdetector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsPython(t *testing.T) {
	t.Run("should return true for python function declaration", func(t *testing.T) {
		src := "def main():\nreturn 1\n"
		assert.True(t, IsPython(src))
	})
	t.Run("should return true for python file opening", func(t *testing.T) {
		src := "with open(\nos.path.join(file_path[0], file_path[1]),\"r\",encoding=\"utf-8\"\n) as file:\nreturn json.load(file)"
		assert.True(t, IsPython(src))
	})
	t.Run("should return true for python condition", func(t *testing.T) {
		src := "if query is None or query.strip() == \"\":\nif save_query_name is not None and save_query_name.strip() != \"\":\nprint(\n\"You have provided the query name for saving, but query argument is empty!\"\n)\nreturn"
		assert.True(t, IsPython(src))
	})
	t.Run("should return false for empty string", func(t *testing.T) {
		assert.False(t, IsPython(""))
	})
	t.Run("should return false for javascript function declaration", func(t *testing.T) {
		src := "function main() {}\n"
		assert.False(t, IsPython(src))
	})
	t.Run("should return false for golang function", func(t *testing.T) {
		src := "func main() {\nfor _, entry := range abc {\n}\n}"
		assert.False(t, IsPython(src))
	})
	t.Run("should return false for ruby ifelse", func(t *testing.T) {
		src := "if is_range_valid? from, to\n@topics = Topic.where(id: from..to).reverse\nelse\n@topics = Topic.all.last(5).reverse\nend"
		assert.False(t, IsPython(src))
	})
	t.Run("should return false for SELECT specific columns  statement", func(t *testing.T) {
		src := "SELECT name, surname FROM users"
		assert.False(t, IsPython(src))
	})
	t.Run("should return false for SELECT everything statement", func(t *testing.T) {
		src := "SELECT * FROM users"
		assert.False(t, IsPython(src))
	})
}

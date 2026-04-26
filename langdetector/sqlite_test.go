package langdetector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsSQLite(t *testing.T) {
	t.Run("should return true for SELECT specific columns  statement", func(t *testing.T) {
		src := "SELECT name, surname FROM users;"
		assert.True(t, IsSQLite(src))
	})
	t.Run("should return true for SELECT everything statement", func(t *testing.T) {
		src := "SELECT * FROM users;"
		assert.True(t, IsSQLite(src))
	})
	t.Run("should return true for PRAGMA statement", func(t *testing.T) {
		src := "PRAGMA journal_mode = WAL;"
		assert.True(t, IsSQLite(src))
	})
	t.Run("should return true for CREATE TABLE with AUTOINCREMENT", func(t *testing.T) {
		src := "CREATE TABLE users (\n  id INTEGER PRIMARY KEY AUTOINCREMENT,\n  name TEXT NOT NULL\n);"
		assert.True(t, IsSQLite(src))
	})
	t.Run("should return true for CREATE TABLE WITHOUT ROWID", func(t *testing.T) {
		src := "CREATE TABLE kv (key TEXT PRIMARY KEY, value TEXT) WITHOUT ROWID;"
		assert.True(t, IsSQLite(src))
	})
	t.Run("should return true for ATTACH DATABASE", func(t *testing.T) {
		src := "ATTACH DATABASE 'other.db' AS other;"
		assert.True(t, IsSQLite(src))
	})
	t.Run("should return false for empty string", func(t *testing.T) {
		assert.False(t, IsSQLite(""))
	})
	t.Run("should return false for PostgreSQL SERIAL type", func(t *testing.T) {
		src := "CREATE TABLE users (id SERIAL PRIMARY KEY, name TEXT);"
		assert.False(t, IsSQLite(src))
	})
	t.Run("should return false for non-SQL text", func(t *testing.T) {
		src := "func main() {\nfmt.Println(\"hello\")\n}"
		assert.False(t, IsSQLite(src))
	})
}

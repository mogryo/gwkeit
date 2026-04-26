package langdetector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsPostgreSQL(t *testing.T) {
	t.Run("should return true for SELECT specific columns  statement", func(t *testing.T) {
		src := "SELECT name, surname FROM users;"
		assert.True(t, IsPostgreSQL(src))
	})
	t.Run("should return true for SELECT everything statement", func(t *testing.T) {
		src := "SELECT * FROM users;"
		assert.True(t, IsPostgreSQL(src))
	})
	t.Run("should return true for CREATE TABLE with SERIAL", func(t *testing.T) {
		src := "CREATE TABLE users (\n  id SERIAL PRIMARY KEY,\n  name TEXT NOT NULL\n);"
		assert.True(t, IsPostgreSQL(src))
	})
	t.Run("should return true for INSERT with RETURNING", func(t *testing.T) {
		src := "INSERT INTO users (name) VALUES ('Alice') RETURNING id;"
		assert.True(t, IsPostgreSQL(src))
	})
	t.Run("should return true for CREATE EXTENSION", func(t *testing.T) {
		src := "CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"
		assert.True(t, IsPostgreSQL(src))
	})
	t.Run("should return true for dollar-quoted function", func(t *testing.T) {
		src := "CREATE OR REPLACE FUNCTION greet(name TEXT) RETURNS TEXT AS $$\nBEGIN\n  RETURN 'Hello, ' || name;\nEND;\n$$ LANGUAGE plpgsql;"
		assert.True(t, IsPostgreSQL(src))
	})
	t.Run("should return true for CREATE SEQUENCE", func(t *testing.T) {
		src := "CREATE SEQUENCE user_id_seq START 1 INCREMENT 1;"
		assert.True(t, IsPostgreSQL(src))
	})
	t.Run("should return false for empty string", func(t *testing.T) {
		assert.False(t, IsPostgreSQL(""))
	})
	t.Run("should return false for SQLite PRAGMA", func(t *testing.T) {
		src := "PRAGMA journal_mode = WAL;"
		assert.False(t, IsPostgreSQL(src))
	})
	t.Run("should return false for SQLite AUTOINCREMENT", func(t *testing.T) {
		src := "CREATE TABLE users (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT);"
		assert.False(t, IsPostgreSQL(src))
	})
	t.Run("should return false for non-SQL text", func(t *testing.T) {
		src := "func main() {\nfmt.Println(\"hello\")\n}"
		assert.False(t, IsPostgreSQL(src))
	})
}

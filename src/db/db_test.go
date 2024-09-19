package db

import (
	"database/sql"
	"testing"
)

func TestDatabase_New(t *testing.T) {
	_, err := NewDatabase("invalid url")
	if err == nil {
		t.Errorf("Expected NewDatabase to return a non-nil error for invalid url")
	}
}

func TestDatabase_Close(t *testing.T) {
	db, err := NewDatabase("postgres://user:password@localhost/database?sslmode=disable")
	if err != nil {
		t.Errorf("Expected NewDatabase to return a non-nil error")
	}
	if err := db.Close(); err != nil {
		t.Errorf("Expected Close to return a non-nil error")
	}
}

func TestDatabase_Query(t *testing.T) {
	db, err := NewDatabase("postgres://user:password@localhost/database?sslmode=disable")
	if err != nil {
		t.Errorf("Expected NewDatabase to return a non-nil error")
	}
	rows, err := db.Query("SELECT 1")
	if err != nil {
		t.Errorf("Expected Query to return a non-nil error")
	}
	defer rows.Close()
	if !rows.Next() {
		t.Errorf("Expected Next to return true")
	}
	var result int
	err = rows.Scan(&result)
	if err != nil {
		t.Errorf("Expected Scan to return a non-nil error")
	}
	if result != 1 {
		t.Errorf("Expected result to be 1")
	}
}

func TestDatabase_Exec(t *testing.T) {
	db, err := NewDatabase("postgres://user:password@localhost/database?sslmode=disable")
	if err != nil {
		t.Errorf("Expected NewDatabase to return a non-nil error")
	}
	result, err := db.Exec("INSERT INTO test (id, name) VALUES (1, 'test')")
	if err != nil {
		t.Errorf("Expected Exec to return a non-nil error")
	}
	if result.RowsAffected() != 1 {
		t.Errorf("Expected RowsAffected to return 1")
	}
}

func TestDatabase_Get(t *testing.T) {
	db, err := NewDatabase("postgres://user:password@localhost/database?sslmode=disable")
	if err != nil {
		t.Errorf("Expected NewDatabase to return a non-nil error")
	}
	var result int
	err = db.Get(&result, "SELECT 1")
	if err != nil {
		t.Errorf("Expected Get to return a non-nil error")
	}
	if result != 1 {
		t.Errorf("Expected result to be 1")
	}
}

func TestDatabase_GetJSON(t *testing.T) {
	db, err := NewDatabase("postgres://user:password@localhost/database?sslmode=disable")
	if err != nil {
		t.Errorf("Expected NewDatabase to return a non-nil error")
	}
	var result map[string]interface{}
	err = db.GetJSON(&result, "SELECT '{"key": "value"}'::json")
	if err != nil {
		t.Errorf("Expected GetJSON to return a non-nil error")
	}
	if result["key"] != "value" {
		t.Errorf("Expected result to be map[string]interface{}{\"key\": \"value\"}")
	}
}

func TestDatabase_Insert(t *testing.T) {
	db, err := NewDatabase("postgres://user:password@localhost/database?sslmode=disable")
	if err != nil {
		t.Errorf("Expected NewDatabase to return a non-nil error")
	}
	id, err := db.Insert("test", map[string]interface{}{"name": "test"})
	if err != nil {
		t.Errorf("Expected Insert to return a non-nil error")
	}
	if id != 1 {
		t.Errorf("Expected id to be 1")
	}
}

func TestDatabase_Update(t *testing.T) {
	db, err := NewDatabase("postgres://user:password@localhost/database?sslmode=disable")
	if err != nil {
		t.Errorf("Expected NewDatabase to return a non-nil error")
	}
	_, err = db.Insert("test", map[string]interface{}{"name": "test"})
	if err != nil {
		t.Errorf("Expected Insert to return a non-nil error")
	}
	affected, err := db.Update("test", map[string]interface{}{"name": "updated"}, "id = 1")
	if err != nil {
		t.Errorf("Expected Update to return a non-nil error")
	}
	if affected != 1 {
		t.Errorf("Expected affected rows to be 1")
	}
}

func TestDatabase_Delete(t *testing.T) {
	db, err := NewDatabase("postgres://user:password@localhost/database?sslmode=disable")
	if err != nil {
		t.Errorf("Expected NewDatabase to return a non-nil error")
	}
	_, err = db.Insert("test", map[string]interface{}{"name": "test"})
	if err != nil {
		t.Errorf("Expected Insert to return a non-nil error")
	}
	affected, err := db.Delete("test", "id = 1")
	if err != nil {
		t.Errorf("Expected Delete to return a non-nil error")
	}
	if affected != 1 {
		t.Errorf("Expected affected rows to be 1")
	}
}

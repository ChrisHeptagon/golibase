package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type DefaultUserSchema struct {
	Email    string `form_type:"email" required:"true" pattern:"^[^\\s]+@[^\\s]+\\.\\w+$" order:"1"`
	Password string `form_type:"password" required:"true" pattern:"^[^\\s]+$" order:"3"`
}

var SchemaFields = []string{
	"form_type",
	"required",
	"pattern",
	"order",
}

func InitializeDatabase(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schemas (
		schema_name TEXT PRIMARY KEY,
		schema BLOB
	);`)
	if err != nil {
		log.Fatal(err)
	}
	tempRows, err := db.Query(`SELECT * FROM schemas WHERE schema_name = 'user_schema';`)
	if err != nil {
		log.Fatal(err)
	}
	defer tempRows.Close()
	if !tempRows.Next() {
		fmt.Println("No user schema found, creating using default schema...")
		if err != nil {
			log.Fatal(err)
		}
		schema := make(map[string]map[string]string)
		defaultSchema := reflect.TypeOf(&DefaultUserSchema{}).Elem()
		for i := 0; i < defaultSchema.NumField(); i++ {
			field := defaultSchema.Field(i)
			schema[field.Name] = make(map[string]string)
			for _, schemaField := range SchemaFields {
				schema[field.Name][schemaField] = field.Tag.Get(schemaField)
			}
		}
		fmt.Println("Default schema: ", schema)
		stmt, err := db.Prepare(`INSERT INTO schemas (schema_name, schema) VALUES (?, ?);`)
		if err != nil {
			log.Fatal(err)
		}
		marshalledSchema, err := json.Marshal(schema)
		if err != nil {
			log.Fatal(err)
		}
		_, err = stmt.Exec("user_schema", marshalledSchema)
		if err != nil {
			log.Fatal(err)
		}
	}

}

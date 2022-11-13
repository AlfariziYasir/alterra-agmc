package migration

import (
	"api-mvc/database/postgres"
	"api-mvc/internal/model"
	"fmt"
	"reflect"
	"strings"
)

var tables = []interface{}{
	&model.User{},
	&model.Book{},
}

func Up() error {
	pg, err := postgres.NewClient()
	if err != nil {
		return err
	}

	err = pg.Conn().AutoMigrate(tables...)
	if err != nil {
		return err
	}

	return nil
}

func Drop() error {
	pg, err := postgres.NewClient()
	if err != nil {
		return err
	}

	err = pg.Conn().Migrator().DropTable(tables...)
	if err != nil {
		return err
	}

	return nil
}

func Status() error {
	var (
		colorReset  = "\033[0m"
		colorGreen  = "\033[32m"
		colorYellow = "\033[33m"
	)

	pg, err := postgres.NewClient()
	if err != nil {
		return err
	}

	fmt.Printf("In database %s:\n", pg.Conn().Migrator().CurrentDatabase())
	for _, table := range tables {
		var name string

		t := reflect.TypeOf(table)
		if t.Kind() == reflect.Ptr {
			name = strings.ToLower(t.Elem().Name()) + "s"
		} else {
			name = strings.ToLower(t.Name()) + "s"
		}

		if pg.Conn().Migrator().HasTable(table) {
			fmt.Println("\t", name, "===>", string(colorGreen), "migrated", string(colorReset))
		} else {
			fmt.Println("\t", name, "===>", string(colorYellow), "not migrated", string(colorReset))
		}
	}

	return nil
}

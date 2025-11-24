package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/internal/config"
	"github.com/mytheresa/go-hiring-challenge/internal/util/pgutil"
)

func main() {
	// Load DB configuration
	c := config.NewDB()

	// Initialize database connection
	db, close := pgutil.New(c)
	defer close()

	files, err := os.ReadDir(c.SQLDirectory)
	if err != nil {
		log.Fatalf("reading directory failed: %v", err)
	}

	// Filter and sort .sql files
	var sqlFiles []os.DirEntry
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file)
		}
	}
	sort.Slice(sqlFiles, func(i, j int) bool {
		return sqlFiles[i].Name() < sqlFiles[j].Name()
	})

	for _, file := range sqlFiles {
		path := filepath.Join(c.SQLDirectory, file.Name())

		content, err := os.ReadFile(path)
		if err != nil {
			log.Printf("reading file %s failed: %v", file.Name(), err)
		}

		sql := string(content)
		if err := db.Exec(sql).Error; err != nil {
			log.Printf("executing %s failed: %v", file.Name(), err)
			return
		}

		log.Printf("Executed %s successfully\n", file.Name())
	}
}

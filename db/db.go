package db

import (
	"fmt"
	"github.com/projectz-ro/journalz-ro/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

type Entry struct {
	ID        uint      `gorm:"primary_key"`
	Name      string    `gorm:"not null"`
	FilePath  string    `gorm:"not null"`
	Tags      []Tag     `gorm:"many2many:entry_tags;"`
	Originals []Entry   `gorm:"many2many:entry_originals;joinTableForeignKey:entry_id;joinForeignKey:original_entry_id"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Tag struct {
	ID        uint      `gorm:"primary_key"`
	TagName   string    `gorm:"not null;unique"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type DatabaseClient struct {
	DB *gorm.DB
}

var USERDB DatabaseClient

func InitializeDB() error {
	db, err := gorm.Open(sqlite.Open(config.CONFIG.ENTRY_DIR+"./metadata.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// Run AutoMigrate
	if err := db.AutoMigrate(&Entry{}, &Tag{}); err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	USERDB.DB = db
	return nil
}

func (c *DatabaseClient) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (c *DatabaseClient) InsertEntry(name string, tags []string, originals []Entry, filepath string) (*Entry, error) {
	newTags := c.batchInsertTags(tags)
	entry := Entry{
		Name:      name,
		Tags:      newTags,
		Originals: originals,
		FilePath:  filepath,
	}

	if err := c.DB.Create(&entry).Error; err != nil {
		return nil, fmt.Errorf("failed to create entry: %v", err)
	}
	return &entry, nil
}

func (c *DatabaseClient) DeleteEntry(id uint) error {
	var entry Entry
	if err := c.DB.First(&entry, id).Error; err != nil {
		return fmt.Errorf("failed to find entry: %v", err)
	}

	if entry.FilePath != "" {
		if err := os.Remove(entry.FilePath); err != nil {
			return fmt.Errorf("failed to delete file at %s: %v", entry.FilePath, err)
		}
	}

	if err := c.DB.Delete(&entry).Error; err != nil {
		return fmt.Errorf("failed to delete entry: %v", err)
	}
	return nil
}

func (c *DatabaseClient) batchInsertTags(tags []string) []Tag {
	var insertedTags []Tag
	for _, tagName := range tags {
		var tag Tag
		if err := c.DB.Where("tag_name = ?", tagName).First(&tag).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				tag = Tag{TagName: tagName}
				if err := c.DB.Create(&tag).Error; err != nil {
					log.Printf("Error creating tag %s: %v", tagName, err)
					continue
				}
			} else {
				log.Printf("Error checking tag existence for %s: %v", tagName, err)
				continue
			}
		}
		insertedTags = append(insertedTags, tag)
	}
	return insertedTags
}

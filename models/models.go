package models

import (
	"strings"
	"time"
)

// Config struct for loading parsed settings.json
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		DBName   string `json:"dbname"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
	Hugo struct {
		BaseDir  string `json:"baseDir"`
		PostsDir string `json:"postsDir"`
	}
}

// PostMetadata is the model for a blog posts attributes
type PostMetadata struct {
	PostID      string    `yaml:"postId"`
	Date        time.Time `yaml:"date"`
	IsBlogPost  string    `yaml:"blogPost"`
	Title       string    `yaml:"title"`
	Description string    `yaml:"description"`
	TagsRaw     string    `yaml:"tags"`
}

// Tags methods exports the YAML tags string as an array
func (pm *PostMetadata) Tags() []string {
	tags := make([]string, 0)
	splitTags := strings.Split(pm.TagsRaw, ",")
	for i := range splitTags {
		tags = append(tags, strings.TrimSpace(splitTags[i]))
	}
	return tags
}

// NotificationMessage object for Postgres event notifications
type NotificationMessage struct {
	Table  string                 `json:"table"`
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data"`
}

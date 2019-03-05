package hugoio

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/kafitz/codimd-hugo-sync/models"
	"gopkg.in/yaml.v2"
)

func getYAMLEnd(header string, content string) (YAMLEnd int) {
	YAMLstart := strings.Index(content, header)
	var searchIdx int
	if YAMLstart > -1 {
		searchIdx = YAMLstart + len(header)
		searchStr := content[searchIdx:]
		YAMLEnd = strings.Index(searchStr, header) + 2*len(header)
	} else {
		fmt.Println("YAML metadata not found.")
	}
	return YAMLEnd
}

func findYAMLData(content string) (YAMLData models.PostMetadata, err error) {
	YAMLHeader := "---"
	YAMLEnd := getYAMLEnd(YAMLHeader, content)
	RawYAMLContent := strings.TrimSpace(content[:YAMLEnd])
	if len(RawYAMLContent) < len(YAMLHeader)*2 {
		err = errors.New("no YAML header found")
		return YAMLData, err
	}
	YAMLContent := strings.TrimSpace(RawYAMLContent[len(YAMLHeader) : len(RawYAMLContent)-len(YAMLHeader)])
	err = yaml.Unmarshal([]byte(YAMLContent), &YAMLData)
	return YAMLData, err
}

// UpdateMetadataFromMarkdown updates an existing metadata object with values from
// a post's front matter: https://gohugo.io/content-management/front-matter/
func UpdateMetadataFromMarkdown(postContent *string, metadata *models.PostMetadata) {
	YAMLData, err := findYAMLData(*postContent)
	if err != nil {
		return
	}
	metadata.Description = YAMLData.Description
	metadata.IsBlogPost = YAMLData.IsBlogPost
	metadata.TagsRaw = YAMLData.TagsRaw
	metadata.Title = YAMLData.Title
}

// WriteMetadataToMarkdown upserts metadata map as YAML at head of markdown post content
func WriteMetadataToMarkdown(postContent string, metadata *models.PostMetadata) (updatedPost string) {
	YAMLHeader := "---"
	YAMLEnd := getYAMLEnd(YAMLHeader, postContent)
	YAMLStrRaw, err := yaml.Marshal(metadata)
	if err != nil {
		log.Fatalln(err)
	}
	YAMLStr := YAMLHeader + "\n" + strings.TrimSpace(string(YAMLStrRaw)) + "\n" + YAMLHeader + "\n"
	updatedPost = YAMLStr + postContent[YAMLEnd+1:]
	return updatedPost
}

// WritePostFile writes markdown content to a file within the hugo
// blog posts directory
func WritePostFile(outputDir string, content string, metadata *models.PostMetadata) error {
	postFilename := fmt.Sprintf("%s.md", metadata.PostID)
	postFilepath := filepath.Join(outputDir, postFilename)
	err := ioutil.WriteFile(postFilepath, []byte(content), 0644)
	return err
}

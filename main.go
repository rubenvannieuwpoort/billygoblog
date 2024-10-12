package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	qjskatex "github.com/graemephi/goldmark-qjs-katex"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

type FileEntry struct {
	Properties map[string]interface{}
	Content    *string
	Tree       *map[string][]FileEntry
}

const (
	CONTENT_FOLDER             = "public"
	TEMPLATE_FOLDER            = "templates"
	TEMPLATE_SUFFIX            = ".template"
	RENDERED_CONTENT_EXTENSION = ".html"
)

var markdown = goldmark.New(
	goldmark.WithExtensions(
		&qjskatex.Extension{},
		extension.Table,
		extension.Typographer,
	),
	goldmark.WithRendererOptions(html.WithHardWraps()),
)

var templates = make(map[string]*template.Template)

var templatesFolder, inputContentFolder string

func main() {
	if len(os.Args) < 4 || (os.Args[2] != "-o" && os.Args[2] != "--output") || os.Args[1] == "" || os.Args[3] == "" {
		programName := os.Args[0]
		fmt.Println("Usage:")
		fmt.Printf("  %s <input_folder> --output <output_folder>\n", programName)
		fmt.Printf("  %s <input_folder> -o <output_folder>\n", programName)
		os.Exit(1)
	}

	render(os.Args[1], os.Args[3])
}

func render(inputPath string, outputPath string) {
	fmt.Printf("Rendering from %s to %s...\n", inputPath, outputPath)

	templatesFolder = filepath.Join(inputPath, TEMPLATE_FOLDER)
	if err := loadTemplates(templatesFolder); err != nil {
		panic(err)
	}

	inputContentFolder = filepath.Join(inputPath, CONTENT_FOLDER)
	if err := generate(inputContentFolder, outputPath); err != nil {
		panic(err)
	}
}

func generate(inputRootPath string, outputRootPath string) error {
	_, err := handleFolder(inputRootPath, outputRootPath)
	return err
}

func loadTemplates(templatesFolder string) error {
	if _, err := os.Stat(templatesFolder); os.IsNotExist(err) {
		return nil
	}

	return filepath.WalkDir(templatesFolder, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() && strings.HasSuffix(path, TEMPLATE_SUFFIX) {
			relativePath, err := filepath.Rel(templatesFolder, path)
			if err != nil {
				return err
			}

			templateID := strings.TrimSuffix(relativePath, TEMPLATE_SUFFIX)

			template := template.New(relativePath).Funcs(map[string]any{
				"renderMarkdown": renderMarkdown,
			})

			template, err = template.ParseFiles(path)
			if err != nil {
				return fmt.Errorf("failed to parse template \"%s\": %w", path, err)
			}

			templates[templateID] = template
		}

		return nil
	})
}

func handleFolder(inputFolder string, outputFolder string) ([]FileEntry, error) {
	// Create output folder.
	if err := os.Mkdir(outputFolder, 0764); err != nil {
		return nil, fmt.Errorf("failed to make directory \"%s\": %w", outputFolder, err)
	}

	// Read subnodes of the folder.
	nodes, err := os.ReadDir(inputFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory \"%s\": %w", inputFolder, err)
	}

	// Build tree for this folder by handling subfolders.
	tree := make(map[string][]FileEntry)
	for _, node := range nodes {
		if !node.Type().IsDir() {
			continue
		}

		name := node.Name()
		if name[0] == '.' {
			continue
		}

		inputPath := filepath.Join(inputFolder, name)
		outputPath := filepath.Join(outputFolder, name)

		if subTree, err := handleFolder(inputPath, outputPath); err != nil {
			return nil, fmt.Errorf("failed to process directory \"%s\": %w", inputPath, err)
		} else {
			tree[node.Name()] = subTree
		}
	}

	// Handle files.
	var result []FileEntry
	for _, node := range nodes {
		if !node.Type().IsRegular() {
			continue
		}

		inputPath := filepath.Join(inputFolder, node.Name())
		outputPath := filepath.Join(outputFolder, node.Name())

		// Open the file.
		fileContents, err := os.ReadFile(inputPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read \"%s\": %w", inputPath, err)
		}

		if hasFrontMatter, frontMatter, body := tryParseFrontMatter(fileContents); hasFrontMatter {
			// There is front matter, so parse it to a FileEntry and add to result.
			var properties map[string]interface{}
			if err := yaml.Unmarshal([]byte(frontMatter), &properties); err != nil {
				return nil, fmt.Errorf("failed to unmarshal \"%s\": %w", inputPath, err)
			}

			uri, err := getURI(inputPath)
			if err != nil {
				return nil, fmt.Errorf("failed to get URI for \"%s\": %w", inputPath, err)
			}

			properties["uri"] = uri

			entry := FileEntry{
				Properties: properties,
				Content:    &body,
				Tree:       &tree,
			}

			var output string
			if template, hasTemplate := properties["template"]; hasTemplate {
				// If a template is specified, use it to render the file.
				templateID, templateIsString := template.(string)
				if !templateIsString {
					return nil, fmt.Errorf("error in file \"%s\" property \"template\" should be a string: %w", inputPath, err)
				}

				output, err = renderTemplate(templateID, entry)
				if err != nil {
					return nil, fmt.Errorf("failed to render template \"%s\" for file \"%s\": %w", templateID, inputPath, err)
				}

				// we replace the extension by .HTML
				outputPath = replaceExtension(outputPath, RENDERED_CONTENT_EXTENSION)
			} else {
				output = body
			}

			// Write the (possibly rendered) content to the output file.
			err = os.WriteFile(outputPath, []byte(output), 0644)
			if err != nil {
				return nil, fmt.Errorf("failed to write to \"%s\": %w", output, err)
			}

			result = append(result, entry)
		} else {
			// Just copy the file if there is no front matter.
			if err := copyFile(inputPath, outputPath); err != nil {
				return nil, fmt.Errorf("failed to copy from \"%s\" to \"%s\": %w", inputPath, outputPath, err)
			}
		}
	}

	return result, nil
}

// extractFrontMatter splits front matter from the body of the markdown.
func extractFrontMatter(content string) (string, string) {
	sections := strings.SplitN(content, "---", 3)
	if len(sections) > 2 {
		return strings.TrimSpace(sections[1]), strings.TrimSpace(sections[2])
	}
	return "", content // No front matter found
}

func tryParseFrontMatter(data []byte) (bool, string, string) {
	if !utf8.Valid(data) {
		return false, "", ""
	}

	frontMatter, body := extractFrontMatter(string(data))
	if frontMatter == "" {
		return false, "", ""
	}

	return true, frontMatter, body
}

func copyFile(sourcePath string, destinationPath string) error {
	// Open the source file.
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open \"%s\": %w", sourcePath, err)
	}
	defer sourceFile.Close()

	// Create the destination file.
	destinationFile, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create file \"%s\": %w", destinationPath, err)
	}
	defer destinationFile.Close()

	// Copy contents.
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy from \"%s\" to \"%s\": %w", sourcePath, destinationPath, err)
	}

	return nil
}

func renderTemplate(templateName string, entry FileEntry) (string, error) {
	template, templateExists := templates[templateName]
	if !templateExists {
		return "", fmt.Errorf("template \"%s\" does not exist", templateName)
	}

	buf := bytes.NewBuffer([]byte{})
	writer := bufio.NewWriter(buf)

	if err := template.Execute(writer, entry); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush writer: %w", err)
	}

	return buf.String(), nil
}

func replaceExtension(path, newExt string) string {
	ext := filepath.Ext(path)
	return strings.TrimSuffix(path, ext) + newExt
}

func renderMarkdown(source string) (template.HTML, error) {
	buf := bytes.NewBuffer([]byte{})
	writer := bufio.NewWriter(buf)

	err := markdown.Convert([]byte(source), writer)
	if err != nil {
		return "", fmt.Errorf("failed to render markdown: %w", err)
	}

	if err := writer.Flush(); err != nil {
		return "", fmt.Errorf("failed to flush writer: %w", err)
	}

	return template.HTML(buf.String()), nil
}

func getURI(path string) (string, error) {
	relativePath, err := filepath.Rel(inputContentFolder, path)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path of \"%s\" to \"%s\": %w", path, inputContentFolder, err)
	}

	return strings.TrimSuffix(relativePath, filepath.Ext(relativePath)), nil
}

package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

var templateText = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta name="go-import"
          content="go.debugged.it/{{ .Name }}
                   git {{ .URL }}" />
    <meta name="go-source"
          content="k8s.io/client-go
                   {{ .URL }}
                   {{ .URL }}/tree/main{/dir}
                   {{ .URL }}/blob/main{/dir}/{file}#L{line}" />
    <meta http-equiv="refresh" content="0; url={{ .URL }}">
</head></html>
`

var indexFileContents = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta http-equiv="refresh" content="0; url=https://github.com/haveyoudebuggedit">
</head></html>
`

type config map[string]string
type entry struct {
	Name string
	URL  string
}

func main() {
	ghPagesDir := "gh-pages"
	if err := os.MkdirAll(ghPagesDir, 0755); err != nil {
		panic(fmt.Errorf("failed to create dir %s (%w)", ghPagesDir, err))
	}

	cnameFile := filepath.Join(ghPagesDir, "CNAME")
	if err := ioutil.WriteFile(cnameFile, []byte("go.debugged.it"), 0644); err != nil {
		panic(fmt.Errorf("failed to write CNAME file %s (%w)", cnameFile, err))
	}

	indexFile := filepath.Join(ghPagesDir, "index.html")
	if err := ioutil.WriteFile(indexFile, []byte(indexFileContents), 0644); err != nil {
		panic(fmt.Errorf("failed to write index file %s (%w)", indexFile, err))
	}

	data, err := ioutil.ReadFile("packages.json")
	if err != nil {
		panic(fmt.Errorf("failed open %s (%w)", "packages.json", err))
	}
	cfg := &config{}
	if err := json.Unmarshal(data, cfg); err != nil {
		panic(fmt.Errorf("failed load %s (%w)", "packages.json", err))
	}
	tpl := template.Must(template.New("html").Parse(templateText))
	for name, url := range *cfg {
		e := entry{
			name, url,
		}
		dir := filepath.Join(ghPagesDir, name)
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(fmt.Errorf("failed to create dir %s (%w)", dir, err))
		}
		file := filepath.Join(ghPagesDir, name, "index.html")
		fh, err := os.Create(file)
		if err != nil {
			panic(fmt.Errorf("failed to open %s (%w)", file, err))
		}
		if err := tpl.Execute(fh, e); err != nil {
			panic(fmt.Errorf("failed to render template (%w)", err))
		}
	}
}

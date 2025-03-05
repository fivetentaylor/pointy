package src

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/teamreviso/code/pkg/env"
)

//go:embed dist/*
var dist embed.FS

func BuildHandler(root string) http.HandlerFunc {
	environment := os.Getenv("ENV")

	// Get the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		panic(err)
	}

	buildOptions := api.BuildOptions{
		Outdir: filepath.Join(workingDir, "pkg/admin/src/dist"),
		EntryPoints: []string{
			filepath.Join(workingDir, "pkg/admin/src/main.ts"),
			filepath.Join(workingDir, "pkg/admin/src/react-app/index.tsx"),
		},
		Platform:     api.PlatformBrowser,
		Bundle:       true,
		MinifySyntax: true,
		Sourcemap:    api.SourceMapLinked,
		Loader: map[string]api.Loader{
			".tsx": api.LoaderTSX,
			".ts":  api.LoaderTS,
		},
		Define: map[string]string{
			"process.env.NODE_ENV": `"` + os.Getenv("NODE_ENV") + `"`,
			"process.env.APP_HOST": `"` + os.Getenv("APP_HOST") + `"`,
			"process.env.WS_HOST":  `"` + os.Getenv("WS_HOST") + `"`,
		},
		Plugins: []api.Plugin{},
		Write:   true,
	}

	responseWithEmbedded := func(w http.ResponseWriter, r *http.Request) {
		log.Info("Serving embedded file", "path", r.URL.Path)

		urlPath := r.URL.Path
		requestPath := strings.TrimPrefix(urlPath, root)

		filePath := filepath.Join("dist", requestPath)
		file, err := dist.Open(filePath)
		if err != nil {
			log.Error("Failed to open embedded file", "path", filePath, "error", err)
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		_, copyErr := io.Copy(w, file)
		if copyErr != nil {
			log.Error("Failed to serve embedded file", "path", filePath, "error", copyErr)
		} else {
			log.Info("Served embedded file", "path", filePath)
		}
		return
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if environment != "development" {
			responseWithEmbedded(w, r)
			return
		}

		log := env.SLog(r.Context())

		urlPath := r.URL.Path
		requestPath := strings.TrimPrefix(urlPath, root)

		if requestPath == "" || requestPath == "/" {
			index(dist, w, r)
			return
		}

		now := time.Now()
		err := buildAndServerFromESBuild(buildOptions, requestPath, w, r)
		log.Info("Built package", "filename", requestPath, "duration", time.Since(now).Seconds())
		if err != nil {
			log.Error("Error building package", "filename", requestPath, "error", err)
			err = fmt.Errorf("Failed to build %s: %v", requestPath, err)

			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte(buildErrorScript(err)))
		}

		return
	}
}

func buildAndServerFromESBuild(
	buildOptions api.BuildOptions,
	requestPath string,
	w http.ResponseWriter,
	_ *http.Request,
) error {
	result := api.Build(buildOptions)
	if len(result.Errors) != 0 {
		return fmt.Errorf("failed to build package: %v", result.Errors)
	}
	w.Header().Set("Content-Type", "application/javascript")

	existingFiles := []string{}
	for _, outputFile := range result.OutputFiles {
		relativePath := strings.TrimPrefix(outputFile.Path, buildOptions.Outdir)
		if strings.HasSuffix(relativePath, requestPath) {
			w.Write(outputFile.Contents)
			return nil
		}
		existingFiles = append(existingFiles, outputFile.Path)
	}

	return fmt.Errorf("file not found: %s. Existing files: %v", requestPath, existingFiles)
}

func buildErrorScript(err error) string {
	return fmt.Sprintf("alert(%q)", err.Error())
}

func index(efs fs.FS, w http.ResponseWriter, _ *http.Request) {
	files, err := listEmbeddedFiles(efs)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("<html><body><ul>"))
	for _, file := range files {
		w.Write([]byte("<li><a href=\"" + file + "\">" + file + "</a></li>"))
	}
	w.Write([]byte("</ul></body></html>"))
	return
}

func listEmbeddedFiles(efs fs.FS) ([]string, error) {
	var files []string
	err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}

func getFilename(root, rawurl string) *string {
	// remove root from rawurl
	rawurl = strings.TrimPrefix(rawurl, root)

	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return nil // or handle the error as you prefer
	}

	filename := path.Base(parsedURL.Path)
	if filename == "/" || filename == "." {
		return nil
	}

	return &filename
}

package assets

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/fivetentaylor/pointy/pkg/env"
)

//go:embed static/*
var staticContent embed.FS

//go:embed src/dist/*
var dist embed.FS

// Static returns the file in the static directory.
// If it cannot be found locally on disk, it will use the embedded static directory.
//
// root is the root directory to serve static files from. ex: pkg/assets/static
// remove is the path to remove from the root. ex: /static
func Static(root, remove string, onlyEmbedded bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		trimmedPath := strings.TrimPrefix(r.URL.Path, remove)

		var (
			data []byte
			hash string
			path string
			err  error
		)

		if !onlyEmbedded {
			// Get the current working directory
			workingDir, err := os.Getwd()
			if err != nil {
				fmt.Println("Error getting working directory:", err)
				panic(err)
			}

			path = filepath.Join(workingDir, root, trimmedPath)

			// Read the file content from the os file system.
			data, err = os.ReadFile(path)
			if err == nil {
				log.Info("[assets] os file", "path", path)
			}
		}

		if data == nil {
			path := filepath.Join("static", filepath.Clean("/"+trimmedPath))
			if strings.HasSuffix(r.URL.Path, "/") {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			log.Info("[assets] embedded file", "path", path)

			var file fs.File
			file, hash, err = getFileAndHash(staticContent, path)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			data, err = io.ReadAll(file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			w.Header().Set("ETag", hash)
		}

		// Determine the content type of the file.
		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		// Set the Content-Type header.
		w.Header().Set("Content-Type", contentType)
		w.Write(data)
	}
}

func BuildAssets() error {
	buildOptions, err := buildOptions()
	if err != nil {
		log.Error("[build] Failed to create build options", "error", err)
		return err
	}

	fmt.Println("Building assets...")
	fmt.Println("  outdir:", buildOptions.Outdir)
	fmt.Println("  entrypoints:", buildOptions.EntryPoints)
	fmt.Println("  defined:", buildOptions.Define)

	err = build(buildOptions)
	if err != nil {
		log.Error("[build] Failed to build package", "error", err)
		return err
	}

	return nil
}

func SrcHandler(root string) http.HandlerFunc {
	environment := os.Getenv("ENV")

	buildOptions, err := buildOptions()
	if err != nil {
		log.Error("[build] Failed to build package", "error", err)
		panic(err)
	}

	responseWithEmbedded := func(w http.ResponseWriter, r *http.Request) {
		urlPath := r.URL.Path
		requestPath := strings.TrimPrefix(urlPath, root)
		filePath := filepath.Join("src/dist", requestPath)

		log.Info("Serving embedded file", "path", r.URL.Path, "filename", filePath)

		file, hash, err := getFileAndHash(dist, filePath)
		if err != nil {
			log.Error("Failed to open embedded file", "path", filePath, "error", err)
			http.NotFound(w, r)
			return
		}
		defer file.Close()

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("ETag", hash)

		contentType := mime.TypeByExtension(filepath.Ext(requestPath))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		w.Header().Set("Content-Type", contentType)

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

		log := env.Log(r.Context())

		urlPath := r.URL.Path
		requestPath := strings.TrimPrefix(urlPath, root)

		if requestPath == "" || requestPath == "/" {
			index(dist, w, r)
			return
		}

		now := time.Now()
		err := buildAndServerFromESBuild(buildOptions, requestPath, w, r)
		log.Info("Built package", "filename", requestPath, "duration", time.Since(now))
		if err != nil {
			log.Error("Error building package", "filename", requestPath, "error", err)
			err = fmt.Errorf("Failed to build %s: %v", requestPath, err)

			w.Header().Set("Content-Type", "application/javascript")
			w.Write([]byte(buildErrorScript(err)))
		}

		return
	}
}

func buildOptions() (api.BuildOptions, error) {
	environment := os.Getenv("ENV")
	isDevelopment := environment == "development"

	// Get the current working directory
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting working directory:", err)
		return api.BuildOptions{}, err
	}

	postcssPath := filepath.Join(workingDir, "pkg/assets/src/node_modules/.bin/postcss")

	buildOptions := api.BuildOptions{
		Outdir: filepath.Join(workingDir, "pkg/assets/src/dist"),
		EntryPoints: []string{
			filepath.Join(workingDir, "pkg/assets/src/rogue.ts"),
			filepath.Join(workingDir, "pkg/assets/src/app/index.tsx"),
			filepath.Join(workingDir, "pkg/assets/src/style/main.css"),
		},
		Platform:          api.PlatformBrowser,
		Bundle:            true,
		MinifySyntax:      !isDevelopment,
		MinifyWhitespace:  !isDevelopment,
		MinifyIdentifiers: !isDevelopment,
		Loader: map[string]api.Loader{
			".tsx":   api.LoaderTSX,
			".ts":    api.LoaderTS,
			".css":   api.LoaderCSS,
			".ttf":   api.LoaderText,
			".woff2": api.LoaderText,
			".svg":   api.LoaderText,
		},
		Define: map[string]string{
			"process.env.NODE_ENV":            `"` + os.Getenv("NODE_ENV") + `"`,
			"process.env.WEB_HOST":            `"` + os.Getenv("WEB_HOST") + `"`,
			"process.env.APP_HOST":            `"` + os.Getenv("APP_HOST") + `"`,
			"process.env.WS_HOST":             `"` + os.Getenv("WS_HOST") + `"`,
			"process.env.SEGMENT_KEY":         `"` + os.Getenv("SEGMENT_KEY") + `"`,
			"process.env.PUBLIC_POSTHOG_KEY":  `"` + os.Getenv("PUBLIC_POSTHOG_KEY") + `"`,
			"process.env.PUBLIC_POSTHOG_HOST": `"` + os.Getenv("PUBLIC_POSTHOG_HOST") + `"`,
			"process.env.IMAGE_TAG":           `"` + os.Getenv("IMAGE_TAG") + `"`,
		},
		Plugins: []api.Plugin{
			{
				Name: "postcss",
				Setup: func(build api.PluginBuild) {
					build.OnLoad(api.OnLoadOptions{Filter: `\.css$`}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
						content, err := os.ReadFile(args.Path)
						if err != nil {
							return api.OnLoadResult{}, err
						}

						cmd := exec.Command(postcssPath, args.Path)
						// TODO: set this as a var
						cmd.Dir = "pkg/assets/src"
						cmd.Stdin = bytes.NewReader(content)
						cmd.Stderr = os.Stderr
						out, err := cmd.Output()
						if err != nil {
							return api.OnLoadResult{}, err
						}

						outString := string(out)

						return api.OnLoadResult{
							Contents: &outString,
							Loader:   api.LoaderCSS,
						}, nil
					})
				},
			},
		},
		Write: true,
	}

	if isDevelopment {
		buildOptions.Sourcemap = api.SourceMapInline
	}

	return buildOptions, nil
}

func build(buildOptions api.BuildOptions) error {
	log.Info("[build] Building package", "outdir", buildOptions.Outdir, "entrypoints", buildOptions.EntryPoints)
	result := api.Build(buildOptions)
	if len(result.Errors) != 0 {
		return fmt.Errorf("failed to build package: %v", result.Errors)
	}
	return nil
}

func buildAndServerFromESBuild(
	buildOptions api.BuildOptions,
	requestPath string,
	w http.ResponseWriter,
	_ *http.Request,
) error {
	result := api.Build(buildOptions)
	if len(result.Errors) != 0 {
		// TODO: make the error better.
		// - there much more information in the error
		// - could make a custom error type that includes that information?
		return fmt.Errorf("failed to build package: %v", result.Errors)
	}

	// Determine the content type of the file.
	contentType := mime.TypeByExtension(filepath.Ext(requestPath))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Set the Content-Type header.
	w.Header().Set("Content-Type", contentType)

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

func getFileAndHash(fs fs.FS, filename string) (fs.File, string, error) {
	file, err := fs.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, "", err
	}

	outfile, err := fs.Open(filename)
	if err != nil {
		return nil, "", err
	}
	return outfile, fmt.Sprintf("%x", hash.Sum(nil)), nil
}

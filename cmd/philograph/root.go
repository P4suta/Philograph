package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"Philograph/internal/api"
	"Philograph/internal/application"
	"Philograph/internal/domain/model"
	"Philograph/internal/domain/service"
	"Philograph/internal/infrastructure/autotokenizer"
	"Philograph/internal/infrastructure/export"
	"Philograph/internal/infrastructure/graphanalyzer"
	kagometok "Philograph/internal/infrastructure/kagome"
	"Philograph/internal/infrastructure/whitespace"
	"Philograph/internal/port"

	"github.com/spf13/cobra"
)

var (
	flagPort      int
	flagLanguage  string
	flagSkipBrowser bool
	flagJSON      bool
	flagVerbose   bool
	flagStopwords string
)

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "philograph [textfile]",
		Short: "Build and visualize co-occurrence networks from texts",
		Long:  "Philograph analyzes text files (.txt, .md, or any UTF-8 text), extracts co-occurrence relationships, and visualizes them as interactive network graphs in the browser.",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runRoot,
	}

	cmd.Flags().IntVar(&flagPort, "port", 0, "HTTP server port (0=auto)")
	cmd.Flags().StringVar(&flagLanguage, "language", "auto", "Language: auto, ja, en")
	cmd.Flags().BoolVar(&flagSkipBrowser, "skip-browser", false, "Don't open browser automatically")
	cmd.Flags().BoolVar(&flagJSON, "json", false, "Output result as JSON to stdout and exit")
	cmd.Flags().BoolVar(&flagVerbose, "verbose", false, "Verbose logging")
	cmd.Flags().StringVar(&flagStopwords, "stopwords", "", "Comma-separated additional stopwords")

	return cmd
}

func runRoot(cmd *cobra.Command, args []string) error {
	// Setup logging
	level := slog.LevelInfo
	if flagVerbose {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))

	// Create AutoTokenizer (supports both languages)
	tokenizer, err := createAutoTokenizer()
	if err != nil {
		return fmt.Errorf("failed to create tokenizer: %w", err)
	}

	// Create analyzer
	analyzer := graphanalyzer.NewAnalyzer()

	// Read input file if provided
	var text string
	config := model.DefaultConfig()

	if len(args) == 1 {
		filePath := args[0]
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}
		text = string(data)

		// Detect language and set config
		lang := detectLanguageFromFlag(text)
		config.Language = lang

		if lang == model.LangJapanese {
			config.StopWords = service.DefaultStopWordsJapanese
		} else {
			config.StopWords = service.DefaultStopWordsEnglish
		}
	} else if flagJSON {
		return fmt.Errorf("--json mode requires a text file argument")
	}

	// Add custom stopwords
	if flagStopwords != "" {
		for _, w := range strings.Split(flagStopwords, ",") {
			w = strings.TrimSpace(w)
			if w != "" {
				config.StopWords = append(config.StopWords, w)
			}
		}
	}

	// JSON mode: run pipeline and output
	if flagJSON {
		pipeline := application.NewPipeline(tokenizer, analyzer, nil)
		result, err := pipeline.Run(cmd.Context(), text, config)
		if err != nil {
			return err
		}
		exporter := export.NewJSONExporter()
		return exporter.Export(os.Stdout, result.Graph)
	}

	// Server mode
	wsHub := api.NewWSHub()
	pipeline := application.NewPipeline(tokenizer, analyzer, wsHub.ProgressListener())
	session := application.NewSession(pipeline, config)

	exporters := map[string]port.Exporter{
		"json": export.NewJSONExporter(),
		"gexf": export.NewGEXFExporter(),
	}

	handler := api.NewHandler(session, exporters)
	server := api.NewServer(handler, wsHub, flagPort)

	port, err := server.Start()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	url := fmt.Sprintf("http://localhost:%d", port)
	fmt.Fprintf(os.Stderr, "Philograph server running at %s\n", url)

	// Auto-analyze if file was provided
	if text != "" {
		go func() {
			if _, err := session.Analyze(context.Background(), text); err != nil {
				slog.Error("initial analysis failed", "error", err)
			} else {
				slog.Info("initial analysis complete")
			}
		}()
	} else {
		fmt.Fprintf(os.Stderr, "No input file specified. Upload a file via the browser to start analysis.\n")
	}

	// Open browser
	if !flagSkipBrowser {
		openBrowser(url)
	}

	return nil
}

func detectLanguageFromFlag(text string) model.Language {
	switch strings.ToLower(flagLanguage) {
	case "ja", "japanese":
		return model.LangJapanese
	case "en", "english":
		return model.LangEnglish
	default:
		return model.DetectLanguage(text)
	}
}

func createAutoTokenizer() (port.Tokenizer, error) {
	ja, err := kagometok.NewTokenizer()
	if err != nil {
		return nil, err
	}
	en := whitespace.NewTokenizer()
	return autotokenizer.NewAutoTokenizer(ja, en), nil
}

func openBrowser(url string) {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "linux":
		// Check for WSL
		if isWSL() {
			cmd = "cmd.exe"
			args = []string{"/c", "start", url}
		} else {
			cmd = "xdg-open"
			args = []string{url}
		}
	case "darwin":
		cmd = "open"
		args = []string{url}
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return
	}

	exec.Command(cmd, args...).Start()
}

func isWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

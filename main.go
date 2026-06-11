package main

import (
	_ "embed"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

//go:embed dj-panel.html
var htmlContent []byte

func main() {
	// Find a free port automatically
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error al iniciar servidor:", err)
		os.Exit(1)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	ln.Close()

	addr := fmt.Sprintf("127.0.0.1:%d", port)
	url := fmt.Sprintf("http://%s", addr)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		w.Write(htmlContent)
	})

	fmt.Println()
	fmt.Println("  ╔══════════════════════════════╗")
	fmt.Println("  ║        DJFX — DJ Panel       ║")
	fmt.Println("  ╠══════════════════════════════╣")
	fmt.Printf("  ║  URL: %-23s ║\n", url)
	fmt.Println("  ║  Abriendo navegador...       ║")
	fmt.Println("  ║  Ctrl+C para cerrar          ║")
	fmt.Println("  ╚══════════════════════════════╝")
	fmt.Println()

	// Open browser after server is ready
	go func() {
		time.Sleep(400 * time.Millisecond)
		openBrowser(url)
	}()

	// Graceful shutdown on Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		fmt.Println("\n  Cerrando DJFX. ¡Hasta la próxima!")
		os.Exit(0)
	}()

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default: // linux, freebsd, etc.
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("  No se pudo abrir el navegador automáticamente.\n  Abre manualmente: %s\n", url)
	}
}

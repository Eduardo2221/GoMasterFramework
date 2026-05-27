package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type DirResult struct {
	Path   string
	Status int
}

type ScanReport struct {
	Target        string
	TotalRequests int
	Findings      []DirResult
}

func EnumerateDIR(alvo string, wordlist string, threads int, extensions []string) (ScanReport, error) {
	var report ScanReport

	if !strings.HasPrefix(alvo, "http://") && !strings.HasPrefix(alvo, "https://") {
		alvo = "http://" + alvo
	}

	report.Target = alvo
	alvo = strings.TrimSuffix(alvo, "/")

	file, err := os.Open(wordlist)
	if err != nil {
		return report, fmt.Errorf("erro ao abrir a wordlist: %v", err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var count int
	var encontrados []DirResult

	tarefas := make(chan string, threads*2)

	transporteInseguro := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transporteInseguro,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Obtém silenciosamente o comportamento do Soft 404 / Catch-All
	statusFalso, tamanhoFalso := ObterLinhaBaseHTTP(client, alvo)

	fmt.Printf("\n[+] Iniciando varredura em: %s\n", alvo)
	fmt.Printf("[+] Threads: %d | Extensões: %v\n", threads, extensions)
	fmt.Println(strings.Repeat("-", 50))

	// INICIA O WORKER POOL
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for path := range tarefas {
				fullURL := fmt.Sprintf("%s/%s", alvo, path)

				req, err := http.NewRequest("GET", fullURL, nil)
				if err != nil {
					continue
				}
				req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) GMF")

				resp, err := client.Do(req)

				mu.Lock()
				count++
				mu.Unlock()

				if err != nil && !errors.Is(err, http.ErrUseLastResponse) {
					if resp != nil {
						resp.Body.Close()
					}
					continue
				}

				if resp != nil {
					// Filtro 1: Ignora os erros explícitos de cara
					if resp.StatusCode == 404 || resp.StatusCode == 400 {
						resp.Body.Close()
						continue
					}

					// Lê o tamanho do corpo da resposta atual
					body, err := io.ReadAll(resp.Body)
					resp.Body.Close()
					if err != nil {
						continue
					}
					tamanhoAtual := int64(len(body))

					// Filtro 2: Validação contra o comportamento do Soft 404 obtido na calibração
					if resp.StatusCode == statusFalso && tamanhoAtual == tamanhoFalso {
						continue
					}

					// Se passou por tudo, guarda e exibe na tela
					mu.Lock()
					encontrados = append(encontrados, DirResult{
						Path:   path,
						Status: resp.StatusCode,
					})
					fmt.Printf("==> ENCONTRADO: /%s (Status: %d)\n", path, resp.StatusCode)
					mu.Unlock()
				}
			}
		}()
	}

	// ALIMENTA OS WORKERS
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		palavra := strings.TrimSpace(scanner.Text())

		if palavra == "" || strings.HasPrefix(palavra, "#") {
			continue
		}

		tarefas <- palavra

		for _, ext := range extensions {
			if !strings.HasPrefix(ext, ".") {
				ext = "." + ext
			}
			tarefas <- palavra + ext
		}
	}

	close(tarefas)
	wg.Wait()

	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("[+] Varredura concluída. Total de requisições: %d\n", count)

	report.TotalRequests = count
	report.Findings = encontrados

	return report, nil
}

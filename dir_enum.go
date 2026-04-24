package main

import (
	"bufio"
	"fmt"
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

func EnumerateDIR(alvo string, wordlist string, threads int) (ScanReport, error) {
	var report ScanReport
	report.Target = alvo
	alvo = strings.TrimSuffix(alvo, "/")

	file, err := os.Open(wordlist)
	if err != nil {
		// Ao invés de printar, retorna o erro para quem chamou a função lidar com ele
		return report, fmt.Errorf("erro ao abrir a wordlist: %v", err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var count int
	var encontrados []DirResult // Lista onde vai guardar os achados

	tarefas := make(chan string, threads*2)

	// INICIA O WORKER POOL
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			client := &http.Client{
				Timeout: 2 * time.Second,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}

			for path := range tarefas {
				fullURL := fmt.Sprintf("%s/%s", alvo, path)

				req, err := http.NewRequest("GET", fullURL, nil)
				if err == nil {
					req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
				}

				resp, err := client.Do(req)

				// Bloqueia para atualizar contadores e salvar dados de forma segura
				mu.Lock()
				count++

				// Se encontrou algo, não printa, apenas salva na estrutura!
				if err == nil && resp.StatusCode != 404 {
					encontrados = append(encontrados, DirResult{
						Path:   path,
						Status: resp.StatusCode,
					})
				}
				mu.Unlock()

				if err == nil {
					resp.Body.Close()
				}
			}
		}(i)
	}

	// ALIMENTA OS WORKERS COM AS TAREFAS
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		palavra := scanner.Text()
		if palavra != "" {
			tarefas <- palavra
		}
	}

	close(tarefas)
	wg.Wait()

	// Preenche o relatório final
	report.TotalRequests = count
	report.Findings = encontrados

	// Retorna o relatório completo e erro nulo (pois deu tudo certo)
	return report, nil
}

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

func EnumerateDIR(alvo string, wordlist string, threads int) {
	alvo = strings.TrimSuffix(alvo, "/")

	file, err := os.Open(wordlist)
	if err != nil {
		fmt.Printf("Erro ao abrir a wordlist: %v\n", err)
		return
	}
	defer file.Close()

	fmt.Printf("[*] Iniciando enumeração de diretórios em: %s\n", alvo)
	fmt.Printf("[*] Criando %d Workers (Threads) focados...\n", threads)
	fmt.Println("--------------------------------------------------")

	var wg sync.WaitGroup
	var mu sync.Mutex
	var count int

	// Canal por onde enviaremos os caminhos para os Workers testarem
	// O buffer (threads * 2) ajuda a manter o fluxo constante sem engarrafar
	tarefas := make(chan string, threads*2)

	// INICIA O WORKER POOL
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// Cliente HTTP com timeout curto e sem seguir redirecionamentos automáticos
			client := &http.Client{
				Timeout: 2 * time.Second,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse // Interrompe o redirect
				},
			}

			// O Worker fica em loop recebendo tarefas do canal até que ele seja fechado
			for path := range tarefas {
				fullURL := fmt.Sprintf("%s/%s", alvo, path)

				// Preparamos a requisição com nosso User-Agent (Disfarce)
				req, err := http.NewRequest("GET", fullURL, nil)
				if err == nil {
					req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
				}

				// Faz a requisição de fato
				resp, err := client.Do(req)

				// Incrementa progresso a cada 1000 tentativas para não poluir a tela
				// Usamos o Mutex (mu.Lock) para evitar que vários workers tentem somar ao mesmo tempo
				mu.Lock()
				count++
				if count%1000 == 0 {
					fmt.Printf("[*] Progresso: %d tentativas realizadas...\n", count)
				}
				mu.Unlock()

				if err != nil {
					continue // Silencioso: Alvo não respondeu ou deu timeout
				}

				// Se o status não for 404 (Not Found), encontramos algo!
				if resp.StatusCode != 404 {
					fmt.Printf("[Worker %d] FOUND: /%s (Status: %d)\n", id, path, resp.StatusCode)
				}
				resp.Body.Close()
			}
		}(i)
	}

	// ALIMENTA OS WORKERS COM AS TAREFAS
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		palavra := scanner.Text()
		if palavra != "" {
			tarefas <- palavra // Envia a palavra para o canal
		}
	}

	// ENCERRAMENTO
	// Fecha o canal para avisar aos workers que não há mais palavras
	close(tarefas)

	// Espera todos os workers terminarem o que estão fazendo
	wg.Wait()

	fmt.Println("--------------------------------------------------")
	fmt.Printf("[*] Enumeração concluída! Total de tentativas: %d\n", count)
}

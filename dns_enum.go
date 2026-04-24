package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

// Representa um subdomínio encontrado e seus IPs
type DNSResult struct {
	Subdomain string
	IPs       []string
}

// Representa o resumo final da varredura DNS
type DNSScanReport struct {
	Target       string
	TotalQueries int
	Findings     []DNSResult
}

func EnumerateDNS(alvo string, wordlist string, threads int) (DNSScanReport, error) {
	var report DNSScanReport
	report.Target = alvo

	// Tenta abrir o arquivo de texto
	file, err := os.Open(wordlist)
	if err != nil {
		return report, fmt.Errorf("erro ao abrir a wordlist: %v", err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var count int
	var encontrados []DNSResult // Lista para guardar os subdomínios válidos

	// Canal por onde envia os prefixos para os Workers testarem
	tarefas := make(chan string, threads*2)

	// INICIA O WORKER POOL
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// O Worker fica em loop recebendo tarefas do canal
			for prefixo := range tarefas {
				// Monta o subdomínio (ex: admin + . + site.com)
				subdominio := prefixo + "." + alvo

				// Faz a pergunta para os servidores DNS da internet
				ips, err := net.LookupHost(subdominio)

				// Bloqueia com Mutex para atualizar contadores e salvar o resultado
				mu.Lock()
				count++

				// Se não deu erro e retornou IPs, salva na estrutura
				if err == nil && len(ips) > 0 {
					encontrados = append(encontrados, DNSResult{
						Subdomain: subdominio,
						IPs:       ips,
					})
				}
				mu.Unlock()
			}
		}()
	}

	// ALIMENTA OS WORKERS COM AS TAREFAS
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		palavra := strings.TrimSpace(scanner.Text())
		if palavra != "" {
			tarefas <- palavra // Envia a palavra para o canal
		}
	}

	// ENCERRAMENTO
	close(tarefas)
	wg.Wait()

	// Preenche o relatório final
	report.TotalQueries = count
	report.Findings = encontrados

	return report, nil
}

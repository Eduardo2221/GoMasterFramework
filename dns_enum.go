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
	report := DNSScanReport{Target: alvo}

	// Tenta abrir o arquivo de texto
	file, err := os.Open(wordlist)
	if err != nil {
		return report, fmt.Errorf("erro ao abrir a wordlist: %v", err)
	}
	defer file.Close()

	// --- INTEGRAÇÃO DA CALIBRAÇÃO (Silenciosa) ---
	// Chama a função compartilhada do seu outro arquivo
	temWildcard, ipsFalsos := ObterLinhaBaseDNS(alvo)

	// Criamos um mapa para busca rápida (O(1)) dos IPs falsos
	wildcardIPs := make(map[string]bool)
	if temWildcard {
		for _, ip := range ipsFalsos {
			wildcardIPs[ip] = true
		}
	}
	// ---------------------------------------------

	var wg sync.WaitGroup
	var mu sync.Mutex
	var count int
	var encontrados []DNSResult

	tarefas := make(chan string, threads*2)

	// INICIA O WORKER POOL
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for prefixo := range tarefas {
				subdominio := prefixo + "." + alvo

				// Faz a consulta DNS
				ips, err := net.LookupHost(subdominio)

				mu.Lock()
				count++
				mu.Unlock()

				if err == nil && len(ips) > 0 {
					// --- FILTRAGEM DO WILDCARD ---
					if temWildcard {
						ehFalsoPositivo := false
						for _, ip := range ips {
							// Se qualquer um dos IPs retornados bater com os IPs da linha base falsa, ignora
							if wildcardIPs[ip] {
								ehFalsoPositivo = true
								break
							}
						}
						if ehFalsoPositivo {
							continue // Ignora em silêncio e vai para a próxima tarefa
						}
					}
					// ------------------------------

					// Se passou pelo filtro, salva o resultado legítimo
					mu.Lock()
					encontrados = append(encontrados, DNSResult{
						Subdomain: subdominio,
						IPs:       ips,
					})
					mu.Unlock()
				}
			}
		}()
	}

	// ALIMENTA OS WORKERS COM AS TAREFAS
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		palavra := strings.TrimSpace(scanner.Text())
		// Ignora vazios ou comentários comuns em wordlists de DNS
		if palavra != "" && !strings.HasPrefix(palavra, "#") {
			tarefas <- palavra
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

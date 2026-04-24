package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

func EnumerateDNS(alvo string, wordlist string, threads int) {
	// Tenta abrir o arquivo de texto
	file, err := os.Open(wordlist)
	if err != nil {
		fmt.Printf("Erro ao abrir a wordlist: %v\n", err)
		return
	}
	defer file.Close() // Garante que o arquivo será fechado no final

	fmt.Printf("[*] Iniciando enumeração de DNS em: %s\n", alvo)
	fmt.Printf("[*] Criando %d Workers (Threads) focados...\n", threads)
	fmt.Println("--------------------------------------------------")

	var wg sync.WaitGroup
	var mu sync.Mutex
	var count int

	// Canal por onde enviaremos os prefixos para os Workers testarem
	tarefas := make(chan string, threads*2)

	// INICIA O WORKER POOL
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// O Worker fica em loop recebendo tarefas do canal
			for prefixo := range tarefas {
				// Monta o subdomínio (ex: admin + . + site.com)
				subdominio := prefixo + "." + alvo

				// Faz a pergunta para os servidores DNS da internet
				ips, err := net.LookupHost(subdominio)

				// Incrementa progresso a cada 1000 tentativas de forma segura
				mu.Lock()
				count++
				if count%1000 == 0 {
					fmt.Printf("[*] Progresso: %d consultas realizadas...\n", count)
				}
				mu.Unlock()

				// Se não deu erro, significa que o subdomínio existe e retornou IPs!
				if err == nil {
					fmt.Printf("[Worker %d] [+] Válido: %s -> %v\n", id, subdominio, ips)
				}
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
	// Fecha o canal para avisar aos workers que o arquivo acabou
	close(tarefas)

	// O script para aqui e espera todos os trabalhadores terminarem
	wg.Wait()

	fmt.Println("--------------------------------------------------")
	fmt.Printf("[*] Enumeração concluída! Total de consultas: %d\n", count)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Menu_main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("--- Go Master Framework (GMF) ---")
	fmt.Println("Digite 'help' para comandos.")

	for {
		fmt.Print("\nGMF > ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		parts := strings.Split(input, " ")
		command := parts[0]

		switch command {

		case "set":
			// Chama a função e recebe o possível erro
			err := HandleSet(parts)
			if err != nil {
				fmt.Printf("[-] Erro: %v\n", err)
			} else {
				// Se não teve erro, avisa que deu certo!
				fmt.Printf("[+] Variável '%s' atualizada para: %s\n", parts[1], parts[2])
			}

		case "show":
			// Chama a função que monta o texto e imprime na hora
			textoConfigs := HandleShow()
			fmt.Println(textoConfigs)

		case "dns":
			if target == "" {
				fmt.Println("[-] Erro: Defina 'target' antes.")
			} else {
				fmt.Printf("[*] Iniciando enumeração DNS em %s...\n", target)

				// Recebemos os dados retornados
				relatorio, err := EnumerateDNS(target, wordlist, threads)

				// Tratamos erros (ex: wordlist não encontrada)
				if err != nil {
					fmt.Printf("[-] Erro na varredura DNS: %v\n", err)
					continue
				}

				// Imprimimos o resultado formatado
				fmt.Printf("\n=== RESULTADOS DNS (%s) ===\n", relatorio.Target)
				fmt.Printf("Consultas feitas: %d | Encontrados: %d\n", relatorio.TotalQueries, len(relatorio.Findings))
				fmt.Println("---------------------------------------")

				if len(relatorio.Findings) > 0 {
					for _, achado := range relatorio.Findings {
						fmt.Printf("[+] %s -> %v\n", achado.Subdomain, achado.IPs)
					}
				} else {
					fmt.Println("[-] Nenhum subdomínio encontrado.")
				}
				fmt.Println("=======================================")
			}

		case "dir":
			if target == "" {
				fmt.Println("[-] Erro: Defina 'target' antes.")
			} else {
				fmt.Printf("[*] Iniciando varredura em %s...\n", target)

				// Recebemos os dados retornados pela nova função
				relatorio, err := EnumerateDIR(target, wordlist, threads)

				// Tratamos o erro (ex: se a wordlist não existir)
				if err != nil {
					fmt.Printf("[-] Erro na varredura: %v\n", err)
					continue
				}

				fmt.Printf("\n=== RESULTADOS DIR (%s) ===\n", relatorio.Target)
				fmt.Printf("Tentativas: %d | Encontrados: %d\n", relatorio.TotalRequests, len(relatorio.Findings))
				fmt.Println("---------------------------------------")

				if len(relatorio.Findings) > 0 {
					for _, achado := range relatorio.Findings {
						fmt.Printf("[+] /%s (HTTP %d)\n", achado.Path, achado.Status)
					}
				} else {
					fmt.Println("[-] Nenhum diretório encontrado.")
				}
				fmt.Println("=======================================")
			}

		case "port":
			if target == "" {
				fmt.Println("[-] Erro: Defina 'target' antes.")
			} else {
				fmt.Printf("[*] Iniciando scan de portas em %s...\n", target)

				relatorio, err := PortScan(target, ports, threads)

				if err != nil {
					fmt.Printf("[-] Erro no scan de portas: %v\n", err)
					continue
				}

				fmt.Printf("\n=== RESULTADOS PORT SCAN (%s) ===\n", relatorio.Target)
				fmt.Printf("Portas testadas: %d | Abertas: %d\n", relatorio.TotalScanned, len(relatorio.OpenPorts))
				fmt.Println("---------------------------------------")

				if len(relatorio.OpenPorts) > 0 {
					for _, achado := range relatorio.OpenPorts {
						fmt.Printf("[+] Porta %s -> %s\n", achado.Port, achado.State)
					}
				} else {
					fmt.Println("[-] Nenhuma porta aberta encontrada.")
				}
				fmt.Println("=======================================")
			}

		case "help":
			fmt.Println("\nComandos Disponíveis:")
			fmt.Println("  set <campo> <valor> : Configura variáveis (target, wordlist, ports, threads)")
			fmt.Println("  show                : Exibe as configurações atuais")
			fmt.Println("  dns                 : Inicia a enumeração de subdomínios (DNS)")
			fmt.Println("  dir                 : Inicia a enumeração de diretórios (Web)")
			fmt.Println("  port                : Inicia o scan de portas TCP")
			fmt.Println("  help                : Exibe este menu de ajuda")
			fmt.Println("  exit                : Fecha o framework")

		case "exit":
			fmt.Println("[*] Saindo do GMF. Até a próxima!")
			return

		default:
			fmt.Printf("[-] Comando desconhecido: %s\n", command)
		}
	}
}

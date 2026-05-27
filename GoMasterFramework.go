package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	target    string = "127.0.0.1"
	wordlist  string = "wordlist.txt"
	ports     string = "1-1024,1433,1521,3306,3389,8000,8080"
	extensoes string = ""
	threads   int    = 10
)

func ParsePorts(portStr string) ([]string, error) {
	// Remover os espaços para evitar problemas
	portStr = strings.ReplaceAll(portStr, " ", "")

	// Separar tudo que estiver dividido por vírgula
	partes := strings.Split(portStr, ",")
	var listaFinal []string

	for _, parte := range partes {
		if parte == "" {
			continue
		}

		// Se a parte contiver um hífen, é um range (ex: "1-100")
		if strings.Contains(parte, "-") {
			rangeLimits := strings.Split(parte, "-")
			if len(rangeLimits) != 2 {
				return nil, fmt.Errorf("formato de range inválido: %s", parte)
			}

			// Convertendo os limites para números inteiros
			inicio, err1 := strconv.Atoi(rangeLimits[0])
			fim, err2 := strconv.Atoi(rangeLimits[1])

			// Validações de segurança
			if err1 != nil || err2 != nil || inicio > fim || inicio < 1 || fim > 65535 {
				return nil, fmt.Errorf("valores de porta inválidos no range: %s", parte)
			}

			// Faz o loop para adicionar todas as portas do range na lista
			for i := inicio; i <= fim; i++ {
				// Convertemos de volta para string porque o net.JoinHostPort pede string
				listaFinal = append(listaFinal, strconv.Itoa(i))
			}
		} else {
			// Se não tem hífen, é apenas uma porta única (ex: "80")
			portaNum, err := strconv.Atoi(parte)
			if err != nil || portaNum < 1 || portaNum > 65535 {
				return nil, fmt.Errorf("porta inválida: %s", parte)
			}
			listaFinal = append(listaFinal, parte)
		}
	}

	return listaFinal, nil
}

func ParseExtensions(extStr string) ([]string, error) {
	// Remover os espaços para evitar problemas
	extStr = strings.ReplaceAll(extStr, " ", "")

	// Separar tudo que estiver dividido por vírgula
	partes := strings.Split(extStr, ",")
	var listaFinal []string

	for _, parte := range partes {
		parte = strings.TrimSpace(parte)
		if parte == "" {
			continue
		}

		// Garante que não há caracteres inválidos para extensões comuns
		// Evita que o usuário passe caminhos ou sujeira na flag
		if strings.ContainsAny(parte, `/\:*?"<>|`) {
			return nil, fmt.Errorf("extensão contém caracteres inválidos: %s", parte)
		}

		// Padronização: remove o ponto inicial se o usuário digitou (ex: ".php" vira "php")
		// Já que a sua função EnumerateDIR adiciona o ponto automaticamente.
		parte = strings.TrimPrefix(parte, ".")

		listaFinal = append(listaFinal, parte)
	}

	return listaFinal, nil
}

func main() {

	// Preparamos o Go para ler argumentos de terminal (o famoso argv)
	// Sintaxe: flag.String("nome_da_flag", "valor_padrao", "Descrição")
	argTarget := flag.String("t", "", "Define o alvo (Ex: example.com)")
	argWordlist := flag.String("w", "wordlist.txt", "Caminho da wordlist")
	argModulo := flag.String("m", "", "Módulo para rodar direto (dns, dir, port)")
	argPorts := flag.String("p", "1-1024,1433,1521,3306,3389,8000,8080", "Portas para escanear (Ex: 80,443)")
	argExtensoes := flag.String("e", ".php,.php.bak,.php.old,.phtml,.jsp,.do,.action,.aspx,.asp,.config,.ascx,.py,.wsgi,.rb,.bak,.old,.swp,.tmp,.zip,.tar.gz,.rar,.7z,.tar,.json,.xml,.log,.yaml,.yml", "Extensões para varredura de diretórios (Ex: .php,.aspx)")
	argThreads := flag.Int("c", 10, "Número de threads")

	flag.Usage = func() {
		fmt.Println("=== Go Master Framework (GMF) ===")
		fmt.Println("\n[*] MODO AUTOMAÇÃO (CLI):")
		fmt.Println("  Uso: ./gmf -m <modulo> -t <alvo> [opções]")
		fmt.Println("\n  Flags:")
		flag.PrintDefaults()

		fmt.Println("\n[*] MODO INTERATIVO:")
		fmt.Println("  Rode o programa sem a flag '-m' para abrir o console.")
		fmt.Println("\n  Comandos Disponíveis (Console):")
		fmt.Println("   set <campo> <valor> : Configura variáveis (target, wordlist, ports, threads)")
		fmt.Println("   show                : Exibe as configurações atuais")
		fmt.Println("   dns                 : Inicia a enumeração de subdomínios (DNS)")
		fmt.Println("   dir                 : Inicia a enumeração de diretórios (Web)")
		fmt.Println("   port                : Inicia o scan de portas TCP")
		fmt.Println("   help                : Exibe este menu de ajuda")
		fmt.Println("   exit                : Fecha o framework")
		fmt.Println("=================================")
	}

	// Lê o que foi digitado no terminal
	flag.Parse()

	// VERIFICAÇÃO PARA AUTOMAÇÃO (BASH)
	// Se o usuário passou a flag "-m" (módulo), roda sem abrir o menu!
	if *argModulo != "" {
		if *argTarget == "" {
			fmt.Println("[-] Para automação, você precisa definir um alvo com -t")
			os.Exit(1)
		}

		// Roda o módulo escolhido silenciosamente e imprime o resultado
		switch *argModulo {
		case "dir":
			exts, err := ParseExtensions(*argExtensoes)
			if err != nil {
				fmt.Printf("[-] Erro nas extensões: %v\n", err)
				os.Exit(1)
			}
			relatorio, _ := EnumerateDIR(*argTarget, *argWordlist, *argThreads, exts)
			if len(relatorio.Findings) > 0 {
				for _, f := range relatorio.Findings {
					fmt.Printf("%s/%s\n", *argTarget, f.Path)
				}
			}

		case "dns":
			relatorio, _ := EnumerateDNS(*argTarget, *argWordlist, *argThreads)
			if len(relatorio.Findings) > 0 {
				for _, f := range relatorio.Findings {
					fmt.Printf("%s\n", f.Subdomain)
				}
			}

		case "port":
			relatorio, _ := PortScan(*argTarget, *argPorts, *argThreads)
			if len(relatorio.OpenPorts) > 0 {
				for _, p := range relatorio.OpenPorts {
					fmt.Printf("%s: Aberta\n", p.Port)
				}
			}

		case "h", "help":
			flag.Usage()
			os.Exit(0)

		default:
			fmt.Printf("[-] Módulo desconhecido: %s\n", *argModulo)
			os.Exit(0)
		}

		// Encerra o programa após a execução (não abre o menu interativo)
		return
	}

	Menu_main()
}

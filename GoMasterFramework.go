package main

import (
	"bufio"
	"flag" // Trocamos a checagem manual de os.Args pelo pacote flag nativo
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// Suas variáveis globais originais (sem valores fixos para permitir alteração dinâmica)
var (
	target    string
	wordlist  string
	ports     string
	extensoes string
	threads   int
)

func main() {
	// --- Configurando as Flags de Linha de Comando ---
	targetFlag := flag.String("t", "127.0.0.1", "Alvo para o scan (IP ou Domínio)")
	threadsFlag := flag.Int("c", 10, "Número de threads (concorrência)")
	
	wordlistFlag := flag.String("w", "", "Wordlist customizada para AMBOS (DNS e DIR)")
	dnsFlag := flag.String("wDNS", "", "Wordlist específica para DNS")
	dirFlag := flag.String("wDIR", "", "Wordlist específica para Diretorios")
	
	portProfile := flag.String("p", "mid", "Perfil de portas: low, mid, high")
	extProfile := flag.String("e", "high", "Perfil de extensões: low, mid, high")

	flag.Parse()

	// Alimentando suas variáveis globais com as flags passadas
	target = *targetFlag
	threads = *threadsFlag

	// Seleção dos perfis de Portas (Regra 4)
	switch strings.ToLower(*portProfile) {
	case "low":  ports = portsLow
	case "high": ports = portsHigh
	default:     ports = portsMid
	}

	// Seleção dos perfis de Extensões (Regra 4)
	switch strings.ToLower(*extProfile) {
	case "low":  extensoes = extsLow
	case "mid":  extensoes = extsMid
	default:     extensoes = extsHigh
	}

	// Slices que vão armazenar as palavras que vão para o fuzzer
	var dnsWords []string
	var dirWords []string

	// Lógica de Prioridade de Wordlists (Regras 1, 2 e 3)
	if *wordlistFlag != "" {
		// Regra 3: Se usar -w, ela sobrescreve tudo
		linhas, err := lerArquivoExterno(*wordlistFlag)
		if err != nil {
			fmt.Printf("[-] Erro ao ler wordlist customizada: %v\n", err)
			os.Exit(1)
		}
		dnsWords = linhas
		dirWords = linhas
		wordlist = *wordlistFlag
	} else {
		// Regra 1: DNS específica ou Embutida Padrão
		if *dnsFlag != "" {
			dnsWords, _ = lerArquivoExterno(*dnsFlag)
		} else {
			dnsWords = stringsToSlice(defaultDNSWordlist)
		}

		// Regra 2: DIR específica ou Embutida Padrão
		if *dirFlag != "" {
			dirWords, _ = lerArquivoExterno(*dirFlag)
			wordlist = *dirFlag
		} else {
			dirWords = stringsToSlice(defaultDIRWordlist)
			wordlist = "Embutida (Padrão)"
		}
	}

	// --- Seu Banner e Prints Originais ---
	fmt.Printf("\n")
	fmt.Printf("   _____    __   __    ______ \n")
	fmt.Printf("  / ___/   /  |/  /   / ____/ \n")
	fmt.Printf(" / / _ _  / /|_/ /   / /_     \n")
	fmt.Printf("/ /_/ /  / /  / /   / __/     \n")
	fmt.Printf("\\____/  /_/  /_/   /_/        \n")
	fmt.Printf("\n")
	fmt.Printf("Go Multi-Threaded Fuzzer\n")
	fmt.Printf("==================================================\n")
	fmt.Printf("[+] Target:    %s\n", target)
	fmt.Printf("[+] Wordlist:  %s (Total de linhas para DIR: %d)\n", wordlist, len(dirWords))
	fmt.Printf("[+] Threads:   %d\n", threads)
	fmt.Printf("[+] Portas:    %s\n", ports)
	fmt.Printf("[+] Extensões: %s\n", extensoes)
	fmt.Printf("==================================================\n\n")

	// --- Seu Motor de Concorrência Original (Intacto) ---
	jobs := make(chan string, len(dirWords))
	var wg sync.WaitGroup

	// Inicializa as goroutines
	for i := 0; i < threads; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	// Envia os dados para as goroutines direto da memória do slice embutido/flag
	for _, word := range dirWords {
		jobs <- word
	}
	close(jobs)

	wg.Wait()
}

// Função auxiliar para abrir arquivos externos passados por flag
func lerArquivoExterno(caminho string) ([]string, error) {
	file, err := os.Open(caminho)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var linhas []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		linha := strings.TrimSpace(scanner.Text())
		if linha != "" {
			linhas = append(linhas, linha)
		}
	}
	return linhas, scanner.Err()
}

// Sua função Worker original (Mantida sem nenhuma alteração de fluxo)
func worker(jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	for word := range jobs {
		url := fmt.Sprintf("%s/%s", target, word)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}

		req.Header.Set("User-Agent", "GMF-Fuzzer/1.0")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}

		if resp.StatusCode == 200 || resp.StatusCode == 204 || resp.StatusCode == 301 || resp.StatusCode == 302 || resp.StatusCode == 403 {
			fmt.Printf("[+] %d -> %s\n", resp.StatusCode, url)
		}
		resp.Body.Close()
	}
}

package main

import (
	_ "embed"
	"strings"
)

// Embutindo as wordlists padrão no binário
//go:embed defaults/dns.txt
var defaultDNSWordlist string

//go:embed defaults/dir.txt
var defaultDIRWordlist string

// Perfis de portas e extensões
var (
	portsLow  = "21,22,23,25,53,80,443"
	portsMid  = "1-1024,1433,1521,3306,3389,8000,8080" // Suas portas padrão atuais
	portsHigh = "1-65535"

	extsLow  = ".php,.html,.txt"
	extsMid  = ".php,.php.bak,.php.old,.phtml,.jsp,.bak,.old,.zip,.json"
	extsHigh = ".php,.php.bak,.php.old,.phtml,.jsp,.do,.action,.aspx,.asp,.config,.ascx,.py,.wsgi,.rb,.bak,.old,.swp,.tmp,.zip,.tar.gz,.rar,.7z,.tar,.json,.xml,.log,.yaml,.yml" // Suas extensões atuais
)

// Função auxiliar para limpar quebras de linha e espaços
func stringsToSlice(conteudo string) []string {
	linhas := strings.Split(conteudo, "\n")
	var limpas []string
	for _, l := range linhas {
		l = strings.TrimSpace(l)
		if l != "" {
			limpas = append(limpas, l)
		}
	}
	return limpas
}

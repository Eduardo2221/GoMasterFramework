package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

func ObterLinhaBaseDNS(dominio string) (bool, []string) {
	subdominioFalso := fmt.Sprintf("detect-wildcard-%d.%s", time.Now().UnixNano(), dominio)

	// Tenta resolver o subdomínio que sabidamente não existe
	ips, err := net.LookupHost(subdominioFalso)
	if err != nil {
		// Se deu erro de "no such host", o comportamento está normal (não tem wildcard)
		return false, nil
	}

	// Se retornou IPs, o domínio tem Wildcard DNS ativo
	return true, ips
}

func ObterLinhaBaseHTTP(client *http.Client, alvo string) (int, int64) {
	urlFalsa := fmt.Sprintf("%s/detect_wildcard_random_path_%d", alvo, time.Now().UnixNano())

	req, err := http.NewRequest("GET", urlFalsa, nil)
	if err != nil {
		return 0, 0
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) GMF")

	resp, err := client.Do(req)
	if err != nil && !errors.Is(err, http.ErrUseLastResponse) {
		return 0, 0
	}
	if resp == nil {
		return 0, 0
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, 0
	}

	return resp.StatusCode, int64(len(body))
}

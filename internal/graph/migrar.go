package graph

import (
	"context"
	"fmt"
	"os"
	"strings"
)

func (c *Cliente) ExecutarScript(ctx context.Context, caminho string) error {
	conteudo, err := os.ReadFile(caminho)
	if err != nil {
		return fmt.Errorf("lendo %s: %w", caminho, err)
	}

	for _, comando := range dividirComandos(string(conteudo)) {
		if _, err := c.Executar(ctx, comando, nil); err != nil {
			return fmt.Errorf("executando %q: %w", resumo(comando), err)
		}
	}
	return nil
}

func dividirComandos(texto string) []string {
	var comandos []string
	for _, parte := range strings.Split(texto, ";") {
		var linhas []string
		for _, linha := range strings.Split(parte, "\n") {
			if strings.HasPrefix(strings.TrimSpace(linha), "//") {
				continue
			}
			linhas = append(linhas, linha)
		}
		if cmd := strings.TrimSpace(strings.Join(linhas, "\n")); cmd != "" {
			comandos = append(comandos, cmd)
		}
	}
	return comandos
}

func resumo(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	if len(s) > 80 {
		return s[:80] + "..."
	}
	return s
}

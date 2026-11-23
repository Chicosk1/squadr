// Comando avaliar mede a qualidade do ranking contra o gabarito das personas (D-037, D-046):
// NDCG@K graduado e Precision@K, na ablação aleatório entre compatíveis → ordem do grafo →
// score unilateral → recíproco sem sinal social → motor completo, com o resultado por persona.
//
// Por padrão lê o Neo4j (rode antes o cmd/seed com os mesmos -usuarios e -semente).
// Com -offline, usa uma fonte em memória com as mesmas regras do grafo, sem banco.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Chicosk1/squadr/internal/avaliacao"
	"github.com/Chicosk1/squadr/internal/config"
	"github.com/Chicosk1/squadr/internal/domain"
	"github.com/Chicosk1/squadr/internal/graph"
	"github.com/Chicosk1/squadr/internal/matching"
	"github.com/Chicosk1/squadr/internal/seed"
)

// O *graph.Cliente precisa satisfazer a fonte do matching sem adaptador (D-034).
var _ matching.Fonte = (*graph.Cliente)(nil)

func main() {
	if err := executar(); err != nil {
		log.Fatal(err)
	}
}

func executar() error {
	usuarios := flag.Int("usuarios", 10000, "deve ser igual ao usado no seed")
	semente := flag.Uint64("semente", 42, "deve ser igual à usada no seed")
	amostra := flag.Int("amostra", 200, "quantidade de usuários avaliados")
	k := flag.Int("k", 10, "tamanho do ranking avaliado")
	offline := flag.Bool("offline", false, "avaliar sem o Neo4j, com a fonte em memória (D-046)")
	flag.Parse()

	ctx := context.Background()
	// O gabarito vem dos perfis regenerados pela semente; o motor lê só a fonte.
	g := seed.NovoGerador(*semente)
	perfis := g.GerarUsuarios(*usuarios)

	var fonte matching.Fonte
	if *offline {
		// As interações vêm do mesmo gerador, na mesma ordem do cmd/seed.
		interacoes := g.GerarInteracoes(perfis)
		pares, err := avaliacao.LerComplementaridade("scripts/referencia.cypher")
		if err != nil {
			return err
		}
		f, err := avaliacao.NovaFonteMemoria(perfis, interacoes, pares)
		if err != nil {
			return err
		}
		fonte = f
		fmt.Println("Fonte: em memória (-offline), mesmas regras do grafo")
	} else {
		cfg := config.Carregar()
		cliente, err := graph.Conectar(ctx, cfg.Neo4jURI, cfg.Neo4jUsuario, cfg.Neo4jSenha, cfg.Neo4jBanco)
		if err != nil {
			return err
		}
		defer cliente.Fechar(ctx)
		fonte = cliente
		fmt.Println("Fonte: Neo4j")
	}
	fmt.Printf("Base: %d usuários, semente %d\n", *usuarios, *semente)

	rel, err := avaliacao.Executar(ctx, fonte, perfis, avaliacao.Config{K: *k, Amostra: *amostra, Opcoes: domain.OpcoesPadrao()})
	if err != nil {
		return err
	}
	rel.Imprimir(os.Stdout)
	return nil
}

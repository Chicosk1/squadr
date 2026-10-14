package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Neo4jURI     string
	Neo4jUsuario string
	Neo4jSenha   string
	Neo4jBanco   string
	PortaAPI     string
}

func Carregar() Config {
	_ = godotenv.Load()

	return Config{
		Neo4jURI:     valor("NEO4J_URI", "bolt://localhost:7687"),
		Neo4jUsuario: valor("NEO4J_USER", "neo4j"),
		Neo4jSenha:   valor("NEO4J_PASSWORD", ""),
		Neo4jBanco:   valor("NEO4J_DATABASE", "neo4j"),
		PortaAPI:     valor("API_PORT", "8080"),
	}
}

func valor(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}

package domain

// FiltroUsuarios controla a listagem de usuários da tela de exploração (D-043).
// Os campos vazios não filtram.
type FiltroUsuarios struct {
	Jogo      Jogo
	Regiao    string // valor de Regioes
	Busca     string // trecho do nome ou começo do id, já em minúsculas
	PrefixoID string // "w" para quem veio do formulário, "u" para o seed (D-039)
	Depois    string // id do último usuário da página anterior
	Limite    int
}

// PaginaUsuarios é uma página da listagem, em ordem de id.
type PaginaUsuarios struct {
	Perfis  []Perfil
	Total   int    // quantos usuários passam no filtro, somando todas as páginas
	Proximo string // valor de Depois para a próxima página; vazio na última
}

# ADR-001: Processo de migrations e RLS no Supabase

**Status:** Aceito
**Data:** 30/06/2026

## Contexto

Durante a criação inicial do schema do banco (migrations 002 a 008), as tabelas foram criadas e aplicadas no projeto remoto via `supabase db push` antes de as políticas de Row Level Security (RLS) terem sido adicionadas ao conteúdo desses arquivos.

Posteriormente, o RLS e as `CREATE POLICY` foram adicionados editando os mesmos arquivos de migration já aplicados. Isso não teve efeito no banco remoto: o Supabase CLI rastreia migrations já executadas pelo nome do arquivo (timestamp), não pelo conteúdo atual. Como essas migrations já constavam como `Applied` no histórico (`supabase_migrations.schema_migrations`), o `db push` seguinte não reaplicou o SQL novo, e todas as 8 tabelas permaneceram com RLS desabilitado (`UNRESTRICTED` no Table Editor) mesmo com o código correto presente nos arquivos.

## Decisão

1. **Migrations já aplicadas nunca são editadas.** Uma vez que uma migration roda com sucesso no remoto, seu conteúdo é tratado como histórico imutável. Qualquer mudança de schema — incluindo políticas de segurança esquecidas — vira uma **migration nova**, nunca uma edição retroativa.

2. **RLS é habilitado na mesma migration que cria a tabela**, daqui em diante. Não há mais separação entre "criar tabela" e "proteger tabela" em momentos diferentes — isso fecha a janela de tempo em que uma tabela existe sem proteção.

3. **Policies usam `DROP POLICY IF EXISTS` antes de `CREATE POLICY`.** Isso torna a migration idempotente — pode ser reaplicada sem erro de "already exists" caso seja necessário rodar novamente em algum ambiente.

4. **Toda migration nova é commitada imediatamente** após confirmação de que rodou com sucesso no painel do Supabase (Table Editor mostrando RLS enabled, ou query em `pg_tables` confirmando `rowsecurity = true`).

## Correção aplicada

Criada a migration `010_enable_rls.sql`, consolidando `ALTER TABLE ... ENABLE ROW LEVEL SECURITY` e todas as `CREATE POLICY` das 8 tabelas (profiles, games, player_games, availability, swipes, matches, messages, reports) em um único arquivo idempotente, aplicado via `supabase db push`.

## Consequências

- Migrations futuras que alterem políticas de segurança em tabelas já existentes devem seguir o mesmo padrão: arquivo novo, nunca edição do antigo.
- Antes de qualquer `db push`, validar que as migrations esperadas aparecem como pendentes via `supabase migration list` — se já aparecerem como `Applied`, o conteúdo não será reexecutado.
- Validação pós-deploy de schema passa a incluir checagem direta em `pg_tables` (`rowsecurity`) como segunda fonte de verdade, não só a interface do Table Editor.
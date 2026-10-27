# ADR 0001 — Fundação local modular

Status: proposta implementada para revisão na issue #1.
Classificação: `METHOD_DOCUMENTED` (engenharia; nenhum método estatístico).

## Contexto e decisão

Adotar três módulos independentes: Go 1.26 (API HTTP), Python 3.12
(namespace reservado, sem processo permanente) e PHP 8.4/Laravel 13 com
Filament 5 instalado para evolução do painel. PostgreSQL 17 guarda futuramente
metadados. Docker Compose é a interface reproduzível; não exige PHP, Go ou
Composer instalados no host. Python 3.12+ é usado pelos scripts locais.

O esqueleto Laravel vem de `laravel/laravel` v13.11.0, commit
`06d016a364a37430eec9cbc52209adce3d7667ce`. Dependências PHP são fixadas por
`painel/composer.lock`. URLs de distribuição usam o endpoint oficial
`codeload.github.com` e os commits originais, pois a API GitHub está bloqueada
neste ambiente. O normalizador preserva checksums existentes e tem testes de
integridade e idempotência; TLS permanece obrigatório. Imagens têm versões de família explícitas e recebem
correções de segurança. Digests não fixados significam que reconstruções futuras
podem receber atualizações: reexecutar CI ao atualizar imagens.

A API não depende do banco porque ainda não ingere ou consulta dados. `/healthz`
verifica o processo; `/up` verifica Laravel; `/readyz` faz `SELECT 1` no banco,
retornando 503 sem detalhes sensíveis em caso de falha. O banco não expõe porta
no host; HTTP é publicado somente em loopback. O servidor PHP embutido serve
exclusivamente desenvolvimento. Não há credenciais de produção.

Filament fica reservado como dependência: não há painel administrativo público,
login, usuários de exemplo ou recursos eleitorais. Autenticação, autorização e
migrações de domínio ficam para outra issue. Não implementar um painel aberto
apenas para demonstrar que o pacote foi instalado.

## Limites e consequências

Não há ingestão TSE, Parquet gerado, dados sintéticos numéricos, equações,
regressão ou Monte Carlo. Cutoff, proveniência, invariantes eleitorais e testes
temporais serão obrigatórios quando existirem dados/experimentos; neste marco
não há fluxo de dados ao qual aplicá-los. Especificações matemáticas seguem
`A_CONFIRMAR` e exigem ADR e revisão próprios.

Volumes retêm o PostgreSQL entre reinícios. `docker compose down` preserva os
dados; não usar `down -v` automaticamente. APP_KEY local é gerada com 32 bytes
aleatórios, sem sobrescrever configuração existente. Nenhuma publicação de
simulações está ativa. Revisões científica e jurídica seguem obrigatórias para
publicação externa de simulações.

## Validação

Comandos e resultados reais são registrados em `docs/VALIDACAO_MVP00.md`.

# Evidências do incremento MVP-01

Execução em 2026-10-27. Base: commit `7166a0f` de `main`, que incorporou a
fundação pelo PR #2. Branch do incremento: `feat/mvp01-historical-data`.

Classificação: `METHOD_DOCUMENTED`. Todos os registros importados são
`HYPOTHESIS`/`SYNTHETIC`; nenhum registro é `DATA_FACT` eleitoral.

## Comandos efetivamente executados

| Comando/verificação | Resultado |
| --- | --- |
| `go vet ./...` na imagem `golang:1.26-bookworm` | passou |
| `go test -race -v ./...` na mesma imagem | 8 testes principais, 34 subtestes; todos passaram; `cmd/api` sem teste unitário |
| `DOCKER_CONFIG=/tmp/observatorio-docker ./scripts/check.sh` | exit 0; execução completa |
| Testes Python de infraestrutura | 5 passaram |
| Ruff | passou |
| Testes Python do namespace estatístico | 2 passaram |
| Composer validate/check-platform-reqs | passaram; aviso já conhecido de licença de software indefinida |
| Pint | 20 arquivos passaram |
| PHPUnit com `--fail-on-warning` | 5 testes, 10 asserções; sem avisos |
| Comando Docker de importação da seção MVP-01 do README | exit 0: `SYNTHETIC: rows=3 accepted=3 rejected=0` |
| Verificação independente em Python dos artefatos produzidos | original idêntico à fixture; SHA-256 bruto e normalizado conferem; três linhas rotuladas; marcador final e qualidade corretos |
| Leitura dos três documentos JSON de configuração/contrato | sintaxe JSON válida; não se afirma execução de um validador externo JSON Schema |
| `git diff --check` | passou |

O teste `TestCLI` compila e executa o binário real. Verifica importação bem
sucedida, destino já existente, destino ausente, cutoff divergente, hash
incorreto e exit 2 para rejeição de qualidade, sem artefato normalizado.

Os testes históricos incluem limites exatos/zeros, inteiros inválidos e overflow,
contagens negativas, votos maiores que comparecimento, comparecimento maior
que eleitorado, cargo/turno incompatíveis, identificadores inválidos, campo
faltante, duplicidade e variantes numéricas de ano/turno, metadados incompletos,
versões incompatíveis, dados com rótulo oficial, URI de rede, disponibilidade
posterior ao cutoff ou à coleta, coleta posterior com disponibilidade anterior,
CSV malformado/vazio, limite de 8 MiB, preservação dos originais, arquivos somente
leitura, retenção da saída inválida e recusa de sobrescrita.

As contagens de linhas não medem cobertura eleitoral. A regra temporal é testada
com tempos de fixtures; não valida disponibilidade histórica de dados do TSE.
Não houve alteração de Compose ou de serviços HTTP existentes; não se repetiu
seu build como condição para validar um CLI local independente.

## Bloqueios verificados

Consultas às referências do catálogo TSE, da documentação técnica e a uma
consulta exploratória CKAN falharam com HTTP 403. Um diagnóstico com headers
confirmou `CONNECT tunnel failed, response 403` no proxy para
`dadosabertos.tse.jus.br`: não houve resposta do catálogo oficial. Não foram
inferidos metadados a partir da falha de rede nem baixados grandes arquivos.

`gh issue create` e a consulta à lista de issues também falharam no endpoint
`https://api.github.com/graphql` com `Forbidden`. O rascunho de issue está em
`docs/ISSUE_MVP01.md`; não há número de issue atribuído por esta implementação.

De acordo com AGENTS.md, a implementação dependente dessas fontes foi
interrompida. Continuam pendentes inventário real, licenças, adaptador TSE,
persistência PostgreSQL/Parquet, pareamento aprovado e cobertura por eleitorado.
O aceite integral de MVP-01 **não foi atingido**. A CI remota não foi verificada.
Não houve publicação de dados, simulações ou merge automático.

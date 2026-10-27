# Contrato interno histórico — v0.1 (somente sintético)

`METHOD_DOCUMENTED`. Este arquivo e os JSON Schemas descrevem uma proposta
interna, **não o layout de um arquivo TSE**. Cada linha normalizada leva
`classification=HYPOTHESIS` e `data_kind=SYNTHETIC`. Nenhum dado oficial foi
importado. A fixture utiliza UF fictícia `ZZ` e identificadores fictícios.

## Dicionário de totais por seção

Granularidade: eleição/turno/cargo/UF/município/zona/seção. Essa composição é a
chave de duplicidade dentro de uma importação. Não corresponde a uma chave
estável entre eleições.

| Campo | Tipo e regra |
| --- | --- |
| `election_year` | inteiro 1900–9999; não comprova a existência de um pleito |
| `round` | inteiro 1 ou 2; sem inferir disponibilidade do segundo turno |
| `office` | literal `PREFEITO` no contrato v0.1 |
| `uf` | duas letras maiúsculas; cadastro oficial ainda não validado |
| `municipality_code` | string de dígitos; preservar zeros e grafia |
| `zone` | string de dígitos; preservar zeros e grafia |
| `section` | string de dígitos; preservar zeros e grafia |
| `electorate` | inteiro int64 não negativo, eleitorado da unidade sintética |
| `turnout` | inteiro int64, entre zero e eleitorado |
| `valid_votes` | inteiro int64, entre zero e comparecimento |
| `classification` | rótulo acrescentado à saída: `HYPOTHESIS` |
| `data_kind` | rótulo acrescentado à saída: `SYNTHETIC` |

CSV de entrada: UTF-8, vírgula, cabeçalho exato na ordem dos primeiros dez
campos da tabela; números inteiros sem separador de milhares. BOM, cabeçalho
TSE, campos adicionais e arquivos ZIP não são aceitos. Linhas vazias são
ignoradas pelo leitor CSV padrão. Há suporte ao escape CSV padrão; erros de
aspas impedem a importação. Não há conversão silenciosa de campos inválidos.

JSON Schema da saída: `section-totals.schema.json`. As relações entre contagens,
duplicidade e os limites temporais são verificadas pelo Go; validar apenas JSON
Schema não basta. Os identificadores não são convertidos a inteiros e não
presumem equivalência entre grafias diferentes como `01` e `1`.

## Manifesto obrigatório

Veja `manifest.schema.json` e a fixture em
`services/api/testdata/historical/manifest.json`.

| Campo | Significado |
| --- | --- |
| `source_url` | URI `fixture://` da fixture local; nenhum download implícito |
| `retrieved_at` | criação/coleta em RFC3339, com timezone |
| `available_at` | disponibilidade declarada, nunca inferida do ano eleitoral |
| `availability_evidence` | referência/explicação da disponibilidade; sintética neste marco |
| `cutoff` | limite explícito da execução, deve coincidir com `--cutoff` |
| `sha256` | 64 dígitos hexadecimais minúsculos dos bytes CSV originais |
| `dataset_version` | versão atribuída à fixture, sem se passar por versão TSE |
| `schema_version` | `section-totals/0.1.0` |
| `parser_version` | `canonical-csv/0.1.0` |
| `license` | declaração da fixture; licença de fontes oficiais permanece `A_CONFIRMAR` |
| `classification` / `data_kind` | exclusivamente `HYPOTHESIS` / `SYNTHETIC` |

`available_at <= cutoff` e `available_at <= retrieved_at`. A igualdade no cutoff
é aceita. Dados finais de pleitos antigos não são snapshots de apuração antiga.
Os testes temporais deste marco comprovam a regra do software com fixtures;
não comprovam a cronologia real de nenhum dataset.

## Artefatos e qualidade

O diretório de execução é criado exclusivamente e nunca substituído:

- `raw/<sha256>.csv`: bytes originais, somente leitura.
- `manifest.json`: manifesto validado, serializado sem alterar sua semântica.
- `quality.json`: contagens de linhas aceitas/rejeitadas e motivos por linha.
- `normalized.jsonl`: somente se todas as linhas forem válidas; inclui rótulos sintéticos.
- `artifacts.json`: SHA-256 do original e, quando gerado, do JSONL.
- `complete.json`: marcador final de execução completa e resultado de qualidade.

`quality_passed` significa apenas que esse contrato sintético passou nas
checagens. Não significa autorização de publicação, correção eleitoral do
arquivo ou validação estatística. `historical_matching=NOT_PERFORMED` é explícito;
contagens de linhas não são cobertura territorial nem cobertura por eleitorado.

CLI: exit 0 para entrada aceita; 2 para rejeições de qualidade com relatório e
sem JSONL; 1 para parâmetros inválidos, hash, manifesto, CSV estrutural, entrada
vazia, tamanho ou I/O. Erros fatais antes da criação do diretório são registrados
em stderr e não produzem artefatos. Consumidores devem exigir marcador final,
qualidade aprovada e hashes correspondentes. Não existem exceções eleitorais
silenciosas ou imputação automática.

# Inventário oficial — bloqueado, não verificado

Classificação: `METHOD_DOCUMENTED`; alvos propostos são `HYPOTHESIS`.
Em 2026-10-27, as referências existentes no README responderam HTTP 403:

| Referência | Finalidade pretendida | Resultado |
| --- | --- | --- |
| https://dadosabertos.tse.jus.br/ | catálogo, datasets e recursos | acesso bloqueado; metadados não obtidos |
| https://www.tse.jus.br/eleicoes/arquivos | layouts e documentação | acesso bloqueado; conteúdo não verificado |

Uma consulta exploratória ao endpoint CKAN `api/3/action/package_search`
também retornou 403. Não se confirmou a resposta nem o contrato dessa API.
Não foram coletados dados, hashes oficiais ou licenças. Ausência de acesso
não significa ausência dos datasets no TSE. O inventário legível por máquina
está em `config/historical_sources.json`, com `datasets=[]` intencionalmente.

Recorte inicial proposto: eleições municipais de 2020/2024, prefeito (turnos 1/2)
e vereador (turno 1). Isso é um plano, não uma afirmação sobre cobertura
real, disponibilidade de recursos ou existência de segundo turno em cada
município. A fixture interna cobre somente totais de prefeito.

## Retomada após liberar acesso

Permitir os hosts `dadosabertos.tse.jus.br` e `www.tse.jus.br` nas configurações
do ambiente e testar novamente as mesmas referências. O diagnóstico adicional
confirmou `CONNECT tunnel failed, response 403` no proxy para o catálogo: não
se obteve resposta do servidor oficial. A configuração salva
não prova acesso. Descobrir outros destinos pelos links oficiais; não inventar
URLs de arquivos nem ampliar a lista por especulação. A API GitHub também está
bloqueada e impede criar a issue/PR automaticamente (`api.github.com`).

Com acesso efetivo, registrar para cada recurso: ID de dataset/recurso, URL
canônica, pleito/cargo/turno, granularidade, cobertura territorial/temporal,
licença e evidência, tamanho, dicionário/encoding/delimitador e atualização.
Aprovar manifesto e limite de tamanho antes de qualquer download substancial.
O comando sintético atual não faz downloads nem aceita os arquivos oficiais.

Validar o adaptador em uma amostra autorizada e pequena, com hash e proveniência.
Não usar arquivos finais como evidência de disponibilidade anterior à apuração.
Pareamento entre eleições deve preservar correspondências aprovadas, motivos,
ambiguidades e mudanças territoriais. Correspondência nominal de identificador
não autoriza pareamento. Reportar cobertura por eleitorado quando essa medida
puder ser calculada a partir de fontes verificadas; até lá, `A_CONFIRMAR`.

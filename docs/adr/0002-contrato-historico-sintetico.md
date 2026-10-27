# ADR 0002 — Contrato histórico e validação sintética v0.1

Status: implementado para revisão no incremento MVP-01; integração oficial bloqueada.
Classificação: `METHOD_DOCUMENTED` para infraestrutura; linhas de teste são
`HYPOTHESIS` e `SYNTHETIC`. Nenhum método eleitoral foi implementado.

## Contexto e decisão

As duas referências oficiais já presentes no README retornaram HTTP 403 no
ambiente. Seguindo AGENTS.md, não inferir layouts TSE, licenças ou campos de
arquivos inacessíveis. Separar o contrato interno de qualquer futuro adaptador
oficial. O incremento prepara e testa a parte independente da fonte.

O contrato v0.1 representa totais por seção, cargo e turno, limitado a prefeito.
Não representa votos individuais nem linhas por candidato. Evita estender sem
revisão a relação votos válidos <= comparecimento a disputas com mais de um
voto por eleitor. UF e identificadores são apenas strings do contrato: ainda não
há validação contra cadastro oficial ou garantia de estabilidade entre eleições.
Identificadores preservam zeros à esquerda; a grafia deve ser padronizada pelo
futuro adaptador documentado. Ano e turno são normalizados numericamente antes
da verificação de duplicidade.

A importação é local e sem rede, limitada a 8 MiB. O CLI exige um cutoff separado
e idêntico ao manifesto. Aceita somente fixtures `SYNTHETIC`/`HYPOTHESIS` com
URI `fixture://`. Uma entrada não se torna oficial por trocar seu rótulo:
`DATA_FACT` é rejeitado enquanto o adaptador verificado não existir.

## Proveniência e integridade

Exigir SHA-256 dos bytes de entrada, versões de dataset/schema/parser, licença
identificada pelo autor da fixture, instante de coleta/criação, disponibilidade
e sua evidência. A evidência sintética não comprova disponibilidade histórica
real. Disponibilidade deve anteceder ou coincidir com cutoff e coleta; coleta
posterior ao cutoff é permitida se a disponibilidade documentada o antecede.
Nunca inferir disponibilidade da data da eleição nem de um arquivo final.

Cada execução usa diretório novo e grava originais pelo hash, sem sobrescrever
originais ou execuções anteriores. Os arquivos de saída são somente leitura
para reduzir alterações acidentais; isto não é armazenamento WORM contra um
administrador. Manifesto, qualidade e hashes dos artefatos acompanham a saída.
`complete.json` só é gravado ao terminar; uma falha de I/O pode deixar uma
execução parcial, que não deve ser consumida. Não apagar nem retomar uma
execução parcial automaticamente: usar outro diretório e investigar a falha.

A condição `0 <= valid_votes <= turnout <= electorate` é uma validação de
integridade documentada pelo projeto, não um estimador. Duplicatas ou linhas
inválidas geram rejeições por linha, sem copiar os valores para os diagnósticos.
Qualquer rejeição retém toda a saída normalizada. Não há imputação, pareamento
histórico automático, interpretação de abstenção, regressão ou simulação.

## Consequências e pendências

Normalização usa JSONL auditável sem bibliotecas externas neste incremento.
Persistência em PostgreSQL e exportação Parquet ficam pendentes até validar o
contrato com a fonte real. Metadados de modelo, parâmetros e semente não se
aplicam a esta importação sem modelo; não gerar valores fictícios para eles.

Não afirmar que o aceite integral do MVP-01 foi atingido: faltam inventário
verificado, licença, adaptador oficial e pareamento/qualidade com bases reais.
Não há publicação externa. A implementação oficial deverá ter nova versão do
contrato e testes vinculados ao dicionário primário.

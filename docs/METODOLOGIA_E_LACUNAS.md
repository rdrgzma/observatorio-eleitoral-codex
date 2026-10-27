# Metodologia e lacunas — documento de trabalho

## Fontes e evidência
- Fonte primária do projeto para resultados: https://dadosabertos.tse.jus.br/
- Método de Rodrigo Mahlmeister: as informações consolidadas em discussões e documentos anteriores são um **resumo de trabalho**, não uma auditoria do código original.
- Não substituir referência documental por reconstrução especulativa.

## Componentes reportados nas fontes discutidas
- Pareamento histórico de locais/seções; regressão sobre mudança eleitoral (`swing`).
- Características do histórico e porte municipal; hierarquia Brasil → região → UF.
- `Empirical Bayes shrinkage`; uso inicial de pesquisas ou referência neutra.
- Ajustes por UF de comparecimento e votos válidos.
- Incerteza aferida por testes retrospectivos.

## Lacunas que bloqueiam réplica idêntica
- Formulação matemática e transformação exata, inclusive zeros.
- Especificação da regressão, pesos e nível de agregação.
- Estimador e hiperparâmetros hierárquicos.
- Dados/pesos de pesquisas e transição entre regimes.
- Calibração exata das faixas, bootstrap e replays.
- Pareamento geográfico, empates e reconciliações.

Cada escolha própria requer ADR e classificação `METHOD_EXPERIMENTAL`.

## Distinção fundamental
**Pré-eleitoral:** simula hipóteses condicionais sem votos da eleição futura.
**Durante apuração:** estima seções pendentes a partir de resultados correntes oficiais e histórico.
Esses motores não devem compartilhar conclusões ou garantias de precisão.

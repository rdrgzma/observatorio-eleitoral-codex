# PR-03/04/05 — Simulação, incerteza e validação

## Experimentos
Pré-registrar eleição, recorte, informação disponível, corte temporal, variáveis e indicadores.
Separar treino, calibração e teste por tempo e território quando adequado.

## Modelos e cenários
Usar primeiro baselines auditáveis; qualquer regressão, modelo hierárquico e simulação Monte Carlo depende de especificação aprovada de estimando, função, pesos, distribuições e dependências.
Simulações são **condicionais às hipóteses**, não medições de preferências individuais.

## Replays
- **Real:** dados capturados com horário de disponibilização durante apuração.
- **Sintético:** ocultação controlada de dados finais, sem afirmar que replica a ordem real.
- Vários marcos de percentual apurado são pontos de **avaliação**, não promessas de acurácia.

## Métricas
Erro assinado, MAE/RMSE, conservação de votos, cobertura e largura de intervalos, calibração e robustez sob diferentes sequências e sementes.
A calibração deve ocorrer fora do conjunto de teste.

## Publicação
Exibir origem oficial e rótulos `OBSERVADO`, `HIPOTÉTICO`, `SIMULADO`; divulgar limitações e realizar revisão científica/jurídica.

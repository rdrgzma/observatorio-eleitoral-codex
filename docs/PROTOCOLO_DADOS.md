# PR-02 — Dados, integridade e proveniência

1. Catalogar bases oficiais por pleito, turno, cargo, localidade, cobertura temporal e licença.
2. Armazenar manifesto (`source_url`, `retrieved_at`, `sha256`, `dataset_version`, `schema_version`, `cutoff`).
3. Separar originais imutáveis de tabelas normalizadas e derivações.
4. Validar chaves; não presumir que códigos de seção/local são estáveis entre eleições.
5. Registrar candidatos a pareamento, ambiguidades, mudanças de local e taxa de cobertura por eleitorado.
6. Definir política explícita de imputação territorial; documentar distância, critério e peso.
7. Para apuração, capturar snapshots e verificar consistência entre configuração de seções e totais agregados.
8. Proibir o uso de arquivos finais para simular conhecimento anterior ao fim do pleito.
9. Produzir relatório de qualidade, com rejeições e advertências, a cada ingestão.

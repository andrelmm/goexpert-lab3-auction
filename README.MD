# FullCycle Auction Go

Este é um projeto de leilão desenvolvido em Go, utilizando MongoDB como banco de dados. Ele permite criar e gerenciar leilões.

## Configuração do Ambiente

1. Clone o repositório:
    ```bash
    git clone https://github.com/andrelmm/goexpert-auction-go.git
    cd goexpert-auction-go
    ```

2. Inicie o ambiente usando Docker Compose:
    ```bash
     docker-compose -f docker-compose.yml up --build
    ```

A aplicação estará disponível em `http://localhost:8080`.

## Exemplos de Chamadas de API

### Criar um Leilão
```bash
curl -X POST http://localhost:8080/auction \
-H "Content-Type: application/json" \
-d '{
    "product_name": "Laptop",
    "category": "Electronics",
    "description": "A brand new laptop",
    "condition": 0
}'
```

### Consultar Leilões

```bash
curl -X GET "http://localhost:8080/auction?category=category&condition=0&status=0"
```

### Consultar um Leilão por ID
    
```bash
curl -X GET http://localhost:8080/auction/5f7b3b7b7b7b7b7b7b7b7b7b
```

## Validação da Função closeExpiredAuctions

- A função closeExpiredAuctions é executada a cada 5 segundos.
- Após iniciar a aplicação, a função closeExpiredAuctions será executada periodicamente.
- Caso algum leilão precise ser fechado (tempo de duração expirado), o status do leilão passará de 0 (ativo) para 1 (concluído).

Para verificar essa funcionalidade:

1. Crie um leilão (A duração padrão é de 30s e pode ser alterada no arquivo .env AUCTION_DURATION).
2. Verifique o status do leilão após a duração ter expirado.

## Testes

Para executar os testes, utilize o comando:

```bash
docker exec -it <container_id> go test ./...
```
# Serviço de Processamento de Exames

Este repositório contém a implementação de um desafio técnico para a posição de **Engenheiro de Software Backend**. O projeto consiste em um serviço para o gerenciamento de exames laboratoriais, onde os exames são registrados via API e processados de forma assíncrona em background.

## Funcionalidades

* **Registro de Exames**: Endpoint para submissão de novos exames.
* **Consulta de Status**: Endpoint para verificar o status atual de um exame (`pending`, `processing`, `done`, `failed`).
* **Processamento Assíncrono**: Utiliza um pool de workers concorrentes para processar os exames sem bloquear a API.
* **Persistência de Dados**: Armazenamento em banco de dados PostgreSQL.
* **Ambiente Containerizado**: Configuração completa com Docker e Docker Compose para facilitar o setup e manter consistência entre ambientes.

## Arquitetura e Decisões de Projeto

### Clean Architecture

A aplicação foi estruturada em camadas (`domain`, `usecase`, `infra`) para garantir separação de responsabilidades.

* **Benefícios**: Alta testabilidade, baixo acoplamento e maior facilidade de manutenção.
* **Trade-off**: Pode ser considerada verbosa em um projeto simples, mas fornece uma base sólida para evolução.

### Processamento Assíncrono

Implementado com **goroutines** e **channels**, seguindo o padrão *Producer-Consumer*.

* **Abordagem**: O `CreateExamUseCase` atua como produtor, adicionando exames à fila em memória (channel). Um pool de workers consome esses exames em paralelo.
* **Benefícios**: Solução leve, performática e sem dependências externas.
* **Trade-off**: Exames em fila podem ser perdidos em caso de reinício da aplicação. Para produção, recomenda-se uma fila persistente (ex.: RabbitMQ, Kafka).

### Banco de Dados e Migrations

Banco de dados PostgreSQL gerenciado por **golang-migrate**.

* **Benefícios**: Versionamento do schema e aplicação automática das migrations ao iniciar a aplicação.
* **Trade-off**: Adiciona dependência externa, mas simplifica a manutenção do schema.

### Graceful Shutdown

A aplicação captura sinais do sistema (SIGINT, SIGTERM) para finalizar de forma ordenada.

* Interrompe novas requisições HTTP.
* Fecha a fila de jobs, evitando inclusão de novos exames.
* Aguarda os workers concluírem o processamento em andamento.
* Libera conexões e recursos de infraestrutura de forma segura.
* Timeout configurável (padrão: 10 segundos).

## Execução do Projeto

### Pré-requisitos

* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/)

### 1. Clone o repositório

```bash
git clone <URL_DO_REPOSITORIO>
cd exam-processing-service
```

### 2. Modo de Execução

#### Modo Desenvolvimento

Utiliza hot-reload com [Air](https://github.com/cosmtrek/air). Recomendado para desenvolvimento e depuração.

```bash
docker-compose up --build
```

**Características**:

* Hot-reload automático.
* Volume montado para edição em tempo real.
* Setup simples para desenvolvimento.
* Limitação no graceful shutdown devido ao proxy do Air.

#### Modo Produção

Executa o binário compilado, simulando o ambiente real. Recomendado para testes de comportamento e demonstrações.

```bash
docker-compose -f docker-compose.prod.yml up --build
```

**Características**:

* Graceful shutdown funcional e completo.
* Melhor desempenho (execução do binário compilado).
* Fiel ao ambiente de produção.
* Requer rebuild para refletir alterações no código.

### 3. Banco de Dados

Na primeira execução, as migrations serão aplicadas automaticamente para criação das tabelas necessárias.

### 4. Acesso à API

A aplicação estará disponível em:

```
http://localhost:8080
```

## Comandos Úteis

```bash
# Encerrar os serviços (desenvolvimento)
docker-compose down

# Encerrar os serviços (produção)
docker-compose -f docker-compose.prod.yml down

# Visualizar logs (desenvolvimento)
docker-compose logs -f

# Visualizar logs (produção)
docker-compose -f docker-compose.prod.yml logs -f

# Testar graceful shutdown (produção)
docker-compose -f docker-compose.prod.yml stop exam-service
```

## Exemplos de Uso da API

### Criar um Exame

```bash
curl -X POST http://localhost:8080/api/v1/exams \
-H "Content-Type: application/json" \
-d '{"patient_id":"12345", "exam_type":"sequenciamento_genetico"}'
```

**Resposta de Sucesso (201 Created):**

```json
{
  "exam_id": "E-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
```

### Consultar Status do Exame

```bash
# Substitua SEU_EXAM_ID pelo ID retornado
curl http://localhost:8080/api/v1/exams/SEU_EXAM_ID
```

**Possíveis Respostas:**

* **Pendente**

```json
{"exam_id":"...","patient_id":"12345","exam_type":"sequenciamento_genetico","status":"pending","created_at":"..."}
```

* **Processando**

```json
{"exam_id":"...","patient_id":"12345","exam_type":"sequenciamento_genetico","status":"processing","created_at":"..."}
```

* **Concluído**

```json
{"exam_id":"...","patient_id":"12345","exam_type":"sequenciamento_genetico","status":"done","created_at":"..."}
```

## Teste do Graceful Shutdown

1. Execute em modo produção:

```bash
docker-compose -f docker-compose.prod.yml up --build -d
```

2. Crie múltiplos exames:

```bash
curl -X POST http://localhost:8080/api/v1/exams \
-H "Content-Type: application/json" \
-d '{"patient_id":"test1", "exam_type":"test_shutdown"}'

curl -X POST http://localhost:8080/api/v1/exams \
-H "Content-Type: application/json" \
-d '{"patient_id":"test2", "exam_type":"test_shutdown"}'
```

3. Encerre o serviço:

```bash
docker-compose -f docker-compose.prod.yml stop exam-service
```

4. Verifique os logs:

```bash
docker-compose -f docker-compose.prod.yml logs exam-service
```

**Exemplo de saída esperada:**

```
Shutdown signal received. Shutting down gracefully...
HTTP server stopped.
Job queue closed. Waiting for workers to finish...
Worker X: finalizou o processamento do exame E-... com status done
All workers have finished.
Server gracefully stopped.
```

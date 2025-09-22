# Serviço de Processamento de Exames

Este repositório contém a implementação de um desafio técnico para a posição de Engenheiro de Software Backend. O projeto consiste em um serviço simples para o gerenciamento de exames laboratoriais, onde os exames são registrados via API e processados de forma assíncrona em background.

## Features

-   **Registro de Exames**: Endpoint para submeter novos exames para processamento.
-   **Consulta de Status**: Endpoint para verificar o status atual de um exame (`pending`, `processing`, `done`, `failed`).
-   **Processamento Assíncrono**: Utiliza um pool de workers concorrentes para processar os exames sem bloquear a API.
-   **Persistência de Dados**: Armazena os dados dos exames em um banco de dados PostgreSQL.
-   **Ambiente Containerizado**: Totalmente configurado para rodar com Docker, garantindo um setup de desenvolvimento simples e consistente.

## Decisões de Arquitetura e Trade-offs

A estrutura do projeto foi guiada por alguns princípios e tecnologias chave:

* **Clean Architecture**: A aplicação foi dividida em camadas (`domain`, `usecase`, `infra`) para garantir uma clara separação de responsabilidades.
    * **Vantagens**: Alta testabilidade, baixo acoplamento entre as camadas e facilidade de manutenção. A lógica de negócio (`domain` e `usecase`) não conhece detalhes de infraestrutura como o banco de dados ou o framework web.
    * **Trade-off**: Para um serviço simples, essa abordagem pode parecer verbosa inicialmente, mas estabelece uma base sólida para o crescimento futuro do sistema.

* **Processamento em Background com Goroutines e Channels**: O requisito de processamento assíncrono com uma fila em memória foi implementado usando as ferramentas nativas de concorrência do Go.
    * **Abordagem**: Foi utilizado o padrão *Producer-Consumer*. O `CreateExamUseCase` (producer) adiciona novos exames a um channel (a fila em memória). Um pool de workers (consumers), rodando em goroutines separadas, consome os exames dessa fila para processamento.
    * **Vantagens**: Solução leve, de alta performance e que não requer dependências externas (como RabbitMQ ou Kafka), cumprindo os requisitos do desafio.
    * **Trade-off**: Por ser uma fila em memória, se a aplicação reiniciar, os exames que estavam na fila e ainda não foram processados serão perdidos. Para um ambiente de produção, uma solução mais robusta como RabbitMQ seria mais indicada para garantir a persistência da fila.

* **Banco de Dados e Migrations**: O PostgreSQL foi escolhido como banco de dados, e o schema é gerenciado via migrations automáticas.
    * **Vantagens**: A biblioteca `golang-migrate` garante que o schema do banco de dados seja versionado e aplicado automaticamente na inicialização da aplicação. Isso elimina a necessidade de configuração manual do banco, tornando o setup do projeto mais simples e confiável.

## Como Rodar o Projeto

O projeto é totalmente containerizado, então tudo que você precisa é ter o Docker e o Docker Compose instalados.

1.  **Clone o repositório:**
    ```bash
    git clone <URL_DO_SEU_REPOSITORIO>
    cd exam-processing-service
    ```

2.  **Inicie os serviços:**
    Execute o seguinte comando na raiz do projeto. Ele irá construir a imagem da aplicação, baixar a imagem do PostgreSQL e iniciar ambos os containers.
    ```bash
    docker-compose up --build
    ```
    Na primeira vez que o comando for executado, as migrations do banco de dados serão aplicadas automaticamente para criar a tabela `exams`.

A API estará disponível em `http://localhost:8080`.

## Exemplos de Requisições

Aqui estão alguns exemplos de como interagir com a API usando `curl`.

### 1. Criar um novo exame

Envie uma requisição `POST` para `/api/v1/exams` com os dados do paciente e o tipo de exame.

```bash
curl -X POST http://localhost:8080/api/v1/exams \
-H "Content-Type: application/json" \
-d '{"patient_id":"12345", "exam_type":"sequenciamento_genetico"}'

Resposta de Sucesso (201 Created):
O exam_id retornado será único para cada requisição.

JSON

{
  "exam_id": "E-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
}
2. Consultar o status de um exame
Envie uma requisição GET para /api/v1/exams/:id, substituindo :id pelo exam_id recebido no passo anterior.

Bash

# Substitua SEU_EXAM_ID pelo ID real
curl http://localhost:8080/api/v1/exams/SEU_EXAM_ID
Possíveis Respostas (200 OK):

Logo após a criação (Pendente):

JSON

{"exam_id":"...","patient_id":"12345","exam_type":"sequenciamento_genetico","status":"pending","created_at":"..."}
Durante o processamento (Processando):
(Você pode ver este status se consultar o exame nos primeiros 5 segundos após a criação)

JSON

{"exam_id":"...","patient_id":"12345","exam_type":"sequenciamento_genetico","status":"processing","created_at":"..."}
Após a finalização (Concluído):

JSON

{"exam_id":"...","patient_id":"12345","exam_type":"sequenciamento_genetico","status":"done","created_at":"..."}
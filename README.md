# Bank API

## :package: Desenvolvimento

### :toolbox: Ferramentas necessárias:

- [Docker](https://docs.docker.com/desktop)
- [Docker Compose](https://docs.docker.com/compose)
- [Golang](https://golang.org/doc/install)

### :rocket: Executando o projeto

- Clone o repositório:
    - `git clone https://github.com/carlosrodriguesf/bank-api`
    - `cd bank-api`
- Configure o projeto
    - `cp .env.dist .env`
    - `make configure`
- Execute o serviço
    - `make run`

Após rodar o comando a lista de rotas disponíveis poderá ser acessada em: `http://localhost:8080/docs`

**Observação:**

Todos os valores na api estão em int64 para evitar problemas com pontos flutuantes, nesse caso os valores devem ser
multiplicados ou divididos por 100 ao serem informados ou lidos, respectivamente.

Exemplos:

- R$ 100,00: `100.00 * 100` = `10000`
- R$ 100,57: `100.57 * 100` = `10057`
- R$ 98.50: `98.5 * 100` = `9850`

### :hammer_and_wrench: Commando disponíveis:

- Execução local
    - `make go-generate`: Executa o comando `go generate ./...` responsável por gerar os mocks.
    - `make go-run`: Executa o comando `go run pkg/main.go`
    - `make go-test`: Executa os testes unitários.
    - `make go-test-cover`: Executa os testes unitários e abre o coverage no navegador.

- Execução dentro do docker
    - `make generate`: Roda o comando `make go-generate`.
    - `make test`: Roda o comando `make go-test`.
    - `make run-services`: Sobe o Redis e o PostgreSQL.
    - `make run`: Sobe o Redis, PostgreSQL e a api.
    - `make run-watch`: Faz a mesma coisa que o comando `make run` e também inicia o nodemon.

- Migrations
    - `make migration-create name="{name}"`: Cria uma migration
        - `name`: Nome da migration a ser criada
    - `make migration-down count={count}`: Faz rollback das migrations
        - `count`: Quantas migrations restaurar

- Documentação
    - `make swagger`: Gera a configuração do swagger

### :open_file_folder: Arquitetura

- `.env`: Arquivo de configurações para execução do projeto localmente
- `.env.dist`: Arquivo de exemplo para a criação do arquivo de configurações
- `migrations/`: Esse diretório possui todas as migration que serão necessarias para rodar a aplicação.
- `docs/`: Arquivos gerados pelo swagger, referente a documentação.
- `pkg/`: contém os pacotes e todo o código da aplicação.
    - `api/`: O codígo relacionado com a camada de api contendo rotas, middlewares e erros http.
        - `error/`: Contém o código para tratamento de erros de api e mapeamento dos códigos http
        - `middleware/`: Contém os middlewares para serem usados na camada de api.
        - `model/`: Contém modelos utilizados entre as rotas da api.
        - `swagger`: Contém o código para executar o swagger.
        - `v1/`: Contém a configuração de rotas da versão 1 da api.
    - `app/`: Aqui fica o código responsável por lidar com as regras de negócio.
    - `model/`: Aqui ficam os modelos globais utilizados entre as camadas do serviço.
    - `error/`: Aqui ficam os possíveis erros mapeados do serviço.
    - `repository/`: Aqui ficam os códigos responsáveis pela comunicação com o banco de dados.
    - `tool/`: Aqui ficam ferramentas para serem usadas na aplicação, facilitando o reaproveitamento de algumas
      funcionalidades.

## :building_construction:️ Construído com

- [echo](https://echo.labstack.com) - Microframework para gerenciamento de rotas
- [godotenv](https://github.com/joho/godotenv) - Carregamento de configurações do .env
- [go-migrate](https://github.com/golang-migrate/migrate) - Gerenciamento de migrations
- [postgresql](https://www.postgresql.org/docs) - Banco de dados da aplicação
- [redis](https://redis.io) - Banco de dados para cache
- [swaggo](https://github.com/swaggo) - Criação de documentação dos endpoints

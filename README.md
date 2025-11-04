# Desafio Técnico – Go (Golang) Backend (API REST)

## 📌 Descrição
Implemente uma **API REST** em **Go (Golang)** para gerenciar **Colaboradores** e **Departamentos**, aplicando regras de negócio de **hierarquia** de departamentos e **gestão de colaboradores**.  
Ao finalizar, publique o **repositório no GitHub** e compartilhe o link para avaliação.

---

## 🚀 Stack Obrigatória
- **Linguagem:** Go 1.22+
- **Framework HTTP:** Gin
- **ORM:** GORM
- **Banco de Dados:** PostgreSQL
- **Migrations:** Flyway
- **Documentação:** Swagger + README com instruções de instalação e uso
- **Containerização:** Docker + docker-compose (app + db)

---

## 🗂️ Modelagem de Domínio

### Colaborador
- **id** (UUIDv7)
- **nome** *(obrigatório)*
- **cpf** *(obrigatório, único, válido)*
- **rg** *(opcional, se informado deve ser único)*
- **departamento_id** *(FK para Departamento, obrigatório)*

### Departamento
- **id** (UUIDv7)
- **nome** *(obrigatório)*
- **gerente_id** *(FK para Colaborador, obrigatório)*
- **departamento_superior_id** *(FK opcional para Departamento – hierarquia)*

### Regras de negócio
1. CPF deve ser único e válido.  
2. RG, se informado, também deve ser único.  
3. O gerente deve ser um Colaborador existente e vinculado ao mesmo departamento.  
4. O Departamento Superior é opcional, mas não pode gerar ciclos na hierarquia.  

---

## 📚 Endpoints

### Colaboradores
- `POST /api/v1/colaboradores` → cria colaborador (validações: CPF, RG, depto existente).  
- `GET /api/v1/colaboradores/:id` → retorna colaborador e o **nome do gerente** do seu departamento.  
- `PUT /api/v1/colaboradores/:id` → atualiza dados.  
- `DELETE /api/v1/colaboradores/:id` → remove colaborador.  
- `POST /api/v1/colaboradores/listar` → lista colaboradores com filtros enviados no **body** (nome, cpf, rg, departamento_id) e paginação.  

### Departamentos
- `POST /api/v1/departamentos` → cria departamento (valida gerente_id).  
- `GET /api/v1/departamentos/:id` → retorna departamento, gerente e **árvore hierárquica completa** dos subdepartamentos.  
- `PUT /api/v1/departamentos/:id` → atualiza departamento (impede ciclos).  
- `DELETE /api/v1/departamentos/:id` → remove departamento.  
- `POST /api/v1/departamentos/listar` → lista departamentos com filtros enviados no **body** (nome, gerente_nome, departamento_superior_id) e paginação.  

### Gerentes
- `GET /api/v1/gerentes/:id/colaboradores` → retorna todos os colaboradores dos departamentos subordinados ao gerente, recursivamente.

---

## ⚖️ Regras Adicionais
- **Prevenção de ciclos** na hierarquia de departamentos.  
- **Constraints de unicidade** no banco (`cpf`, `rg`).  
- **Respostas de erro consistentes:**  
  - `422` → erro de validação de domínio  
  - `404` → recurso não encontrado  
  - `400` → filtros inválidos  
  - `409` → conflito de unicidade  

---

## 📦 Entregáveis
- Código em **repositório GitHub**.  
- **README.md** contendo:
  - Como rodar o projeto com Docker + PostgreSQL.  
  - Como executar migrations com Flyway.  
  - Como acessar a documentação Swagger.  
  - Exemplos de requests (via curl/Postman/Insomnia).  
- **Swagger** acessível em `/docs`.  
- `docker-compose.yml` para facilitar o setup local.

---

## ✅ Como Entregar

- Suba o código no **GitHub** (público ou privado com acesso).  
- Inclua no **README.md**:  
  - Como subir o ambiente com **Docker**.  
  - Como rodar migrations com **Flyway**.  
  - Como acessar a documentação **Swagger**.  
  - Seeds (se houver).  
- Envie o **link do repositório** para avaliação.  


---

## 🏆 Critérios de Avaliação
1. **Qualidade do código:** clareza, organização em camadas (`handlers`, `services`, `repositories`, `models`).  
2. **Corretude das regras de negócio:** unicidade de CPF/RG, gerente válido, prevenção de ciclos.  
3. **Funcionalidade:** CRUDs completos, filtros e desafios adicionais funcionando.  
4. **Banco & ORM:** modelagem correta no PostgreSQL, constraints e integridade.  
5. **Documentação:** Swagger + README detalhado.  
6. **Entrega:** execução simples via Docker, sem fricção para rodar localmente.  

---

## ⭐ Diferenciais (Bônus)
- Uso de **CTE recursivo no PostgreSQL** para montar a árvore de departamentos.  
- **Cache** para otimizar consultas hierárquicas.  
- Logs estruturados (ex.: zap/logrus).  
- **Métricas básicas** (Prometheus ou middleware Gin).  

---

## 📦 Como Rodar o Projeto

### 🐳 Com Docker + PostgreSQL

1. **Instalar dependências (opcional, para gerar Swagger localmente)**

    ```bash
    make deps
    ```

2. **Subir os serviços (app + postgres)**

    ```bash
    make up
    ```

3. **Ver logs**

    ```bash
    make logs
    ```

4. **Reconstruir imagens após alterações**

    ```bash
    make rebuild
    ```

5. **Parar e limpar containers**

    ```bash
    make down
    make clean  # remove volumes também
    ```

### 🔗 Endpoints expostos

-   API: [http://localhost:8080](http://localhost:8080)
-   Swagger: [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)

### ⚙️ Observações

-   Certifique-se de possuir um arquivo `.env` com:

    ```
    POSTGRES_HOST=localhost
    POSTGRES_DB=takehome
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    ```

-   A aplicação usa PostgreSQL; garanta que a porta `5432` esteja livre.

---

## 🧭 Migrations com Flyway

### ▶️ Executar migrations

```bash
make migrate
```

### ⚙️ Pré-requisitos

-   Serviço **flyway** definido no `docker-compose` com volume apontando para `./db/migrations`.
-   Variáveis configuradas (no `.env` ou `docker-compose`).

---

## 🧾 Documentação Swagger

### 📄 Geração/atualização

```bash
make deps
make swagger
```

### 🌐 Acesso no navegador

[http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)

## 🧪 Exemplos de Requests

### 🔹 Criar colaborador

```bash
curl -X POST http://localhost:8080/api/v1/colaboradores \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Ana Silva",
    "cpf": "12345678909",
    "rg": "MG1234567",
    "departamento_id": "018f3c3e-5c79-7b21-b7e1-d45f80cfa5ab"
  }'
```

### 🔹 Obter colaborador por ID

```bash
curl http://localhost:8080/api/v1/colaboradores/018f3c3e-5c79-7b21-b7e1-d45f80cfa5ac
```

### 🔹 Listar colaboradores com filtros

```bash
curl -X POST http://localhost:8080/api/v1/colaboradores/listar \
  -H "Content-Type: application/json" \
  -d '{
    "filtros": {
      "nome": "Ana",
      "cpf": "",
      "rg": "",
      "departamento_id": ""
    },
    "page": 1,
    "page_size": 20
  }'
```

### 🔹 Criar departamento

```bash
curl -X POST http://localhost:8080/api/v1/departamentos \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "TI",
    "gerente_id": "018f3c3e-5c79-7b21-b7e1-d45f80cfa5ad",
    "departamento_superior_id": null
  }'
```

### 🔹 Obter departamento (com hierarquia)

```bash
curl http://localhost:8080/api/v1/departamentos/018f3c3e-5c79-7b21-b7e1-d45f80cfa5ae
```

### 🔹 Listar departamentos com filtros

```bash
curl -X POST http://localhost:8080/api/v1/departamentos/listar \
  -H "Content-Type: application/json" \
  -d '{
    "filtros": {
      "nome": "TI",
      "gerente_nome": "Ana",
      "departamento_superior_id": ""
    },
    "page": 1,
    "page_size": 10
  }'
```

### 🔹 Colaboradores subordinados a um gerente

```bash
curl http://localhost:8080/api/v1/gerentes/018f3c3e-5c79-7b21-b7e1-d45f80cfa5ad/colaboradores
```

---

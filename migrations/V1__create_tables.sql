CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS departamentos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(255) NOT NULL,
    gerente_id UUID,
    departamento_superior_id UUID,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (departamento_superior_id) REFERENCES departamentos(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS colaboradores (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(255) NOT NULL,
    cpf VARCHAR(11) NOT NULL UNIQUE,
    rg VARCHAR(20) UNIQUE,
    departamento_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (departamento_id) REFERENCES departamentos(id) ON DELETE RESTRICT
);

ALTER TABLE departamentos 
ADD CONSTRAINT fk_departamentos_gerente 
FOREIGN KEY (gerente_id) REFERENCES colaboradores(id) ON DELETE RESTRICT;

CREATE INDEX idx_colaboradores_cpf ON colaboradores(cpf);
CREATE INDEX idx_colaboradores_rg ON colaboradores(rg) WHERE rg IS NOT NULL;
CREATE INDEX idx_colaboradores_departamento ON colaboradores(departamento_id);
CREATE INDEX idx_departamentos_gerente ON departamentos(gerente_id);
CREATE INDEX idx_departamentos_superior ON departamentos(departamento_superior_id) WHERE departamento_superior_id IS NOT NULL;

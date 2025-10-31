CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS departamentos (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nome VARCHAR(255) NOT NULL,
    gerente_id UUID NOT NULL,
    departamento_superior_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (departamento_superior_id) REFERENCES departamentos(id) ON DELETE SET NULL
);

CREATE INDEX idx_departamentos_gerente ON departamentos(gerente_id);
CREATE INDEX idx_departamentos_superior ON departamentos(departamento_superior_id);
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE departamentos (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  nome TEXT NOT NULL,
  gerente_id UUID,
  departamento_superior_id UUID NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT departamentos_nome_not_blank CHECK (length(trim(nome)) > 0),
  CONSTRAINT departamentos_superior_fk FOREIGN KEY (departamento_superior_id) REFERENCES departamentos (id) ON DELETE SET NULL
);

CREATE TRIGGER departamentos_set_updated_at
BEFORE UPDATE ON departamentos
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE colaboradores (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  nome TEXT NOT NULL,
  cpf VARCHAR(11) NOT NULL,
  rg TEXT NULL,
  departamento_id UUID NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT colaboradores_nome_not_blank CHECK (length(trim(nome)) > 0),
  CONSTRAINT colaboradores_cpf_len CHECK (length(cpf) = 11),
  CONSTRAINT colaboradores_depto_fk FOREIGN KEY (departamento_id) REFERENCES departamentos (id) ON DELETE RESTRICT
);

CREATE TRIGGER colaboradores_set_updated_at
BEFORE UPDATE ON colaboradores
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
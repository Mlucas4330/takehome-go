ALTER TABLE colaboradores ADD CONSTRAINT colaboradores_cpf_unique UNIQUE (cpf);

ALTER TABLE colaboradores ADD CONSTRAINT colaboradores_rg_unique UNIQUE (rg);

ALTER TABLE departamentos
ALTER COLUMN gerente_id
SET
  NOT NULL;

ALTER TABLE departamentos ADD CONSTRAINT departamentos_gerente_fk FOREIGN KEY (gerente_id) REFERENCES colaboradores (id) ON DELETE RESTRICT;
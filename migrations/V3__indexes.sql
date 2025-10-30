CREATE INDEX idx_colaboradores_nome ON colaboradores (nome);

CREATE INDEX idx_colaboradores_cpf ON colaboradores (cpf);

CREATE INDEX idx_colaboradores_rg ON colaboradores (rg);

CREATE INDEX idx_colaboradores_departamento ON colaboradores (departamento_id);

CREATE INDEX idx_departamentos_nome ON departamentos (nome);

CREATE INDEX idx_departamentos_superior ON departamentos (departamento_superior_id);

CREATE INDEX idx_departamentos_gerente ON departamentos (gerente_id);
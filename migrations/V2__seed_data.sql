INSERT INTO departamentos (id, nome, gerente_id)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Diretoria',
    NULL
);

INSERT INTO colaboradores (id, nome, cpf, departamento_id)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'João Silva',
    '12345678909',
    '00000000-0000-0000-0000-000000000001'
);

UPDATE departamentos
SET gerente_id = '00000000-0000-0000-0000-000000000001'
WHERE id = '00000000-0000-0000-0000-000000000001';

INSERT INTO colaboradores (nome, cpf, rg, departamento_id)
VALUES
    (
        'Maria Santos',
        '98765432100',
        'MG1234567',
        '00000000-0000-0000-0000-000000000001'
    ),
    (
        'Pedro Oliveira',
        '11122233344',
        'SP9876543',
        '00000000-0000-0000-0000-000000000001'
    );

INSERT INTO departamentos (id, nome, gerente_id, departamento_superior_id)
SELECT
    gen_random_uuid(),
    'TI',
    c.id,
    '00000000-0000-0000-0000-000000000001'
FROM colaboradores c
WHERE c.cpf = '98765432100';

INSERT INTO departamentos (id, nome, gerente_id, departamento_superior_id)
SELECT
    gen_random_uuid(),
    'RH',
    c.id,
    '00000000-0000-0000-0000-000000000001'
FROM colaboradores c
WHERE c.cpf = '11122233344';

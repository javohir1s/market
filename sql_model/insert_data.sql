



INSERT INTO "user"(
    id,
    first_name,
    last_name,
    login,
    password,
    client_type,
    updated_at
) VALUES (
    'f4877292-4468-44e2-a27f-9743c7a35802',
    '',
    '',
    'superadmin',
    'admin123',
    'SUPER_ADMIN',
    NOW()
)


CREATE UNIQUE INDEX user_login_idx ON "user"("login");


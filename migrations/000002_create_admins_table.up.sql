CREATE TABLE admins (
    id UUID PRIMARY KEY NOT NULL,
    full_name VARCHAR(200) NOT NULL,
    age INT,
    email TEXT NOT NULL,
    username VARCHAR(200) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role VARCHAR(100) NOT NULL,
    refresh_token TEXT
    );


INSERT INTO admins (id, full_name, age, username, email, password, role, refresh_token)
VALUES (
            'e74a31c2-ade8-444c-8aa2-4cd644d9db8f',
            'Mubina Yigitaliyeva',
            17,
            'admin',
            'mubinayigitaiyeva00@gmail.com',
            '$2a$14$STJ5xMH5Q9eaA4lBqOJ2Qu0o3BRIedl5DKhzInTKogQoPJkU6Z9bK',
            'superadmin',
            'refresh_token'
        );
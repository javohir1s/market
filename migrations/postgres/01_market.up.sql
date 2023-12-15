
CREATE TABLE "user" (
    "id" UUID NOT NULL PRIMARY KEY,
    "first_name" VARCHAR(46) NOT NULL,
    "last_name" VARCHAR(46) NOT NULL,
    "login" VARCHAR(46) NOT NULL,
    "password" VARCHAR NOT NULL,
    "active" BOOLEAN NOT NULL DEFAULT true,
    "client_type" VARCHAR(46) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP
);





-- branch
CREATE TABLE branch (
    id UUID PRIMARY KEY,
    branch_code VARCHAR(10),
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255),
    phone VARCHAR(20), 
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- salepoint
CREATE TABLE sale_point (
    id UUID PRIMARY KEY,
    branch_id UUID NOT NULL REFERENCES branch(id),
    name VARCHAR(255) NOT NULL, 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP

);

-- employee
CREATE TABLE employee (
    id UUID PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    phone VARCHAR(20),
    login VARCHAR(50),
    password VARCHAR(50),
    branch_id UUID NOT NULL REFERENCES branch(id),
    salepoint_id UUID NOT NULL REFERENCES sale_point(id),
    user_type VARCHAR(20), 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- supplier
CREATE TABLE supplier (
    id UUID PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    is_active BOOLEAN, 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- category
CREATE TABLE category (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    parent_id UUID REFERENCES category(id), 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- product
CREATE TABLE product (
    id UUID PRIMARY KEY,
    photo VARCHAR,
    title VARCHAR(255) NOT NULL,
    category_id UUID NOT NULL REFERENCES category(id),
    barcode VARCHAR(50),
    price DECIMAL(10, 2), 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- income
CREATE TABLE income (
    id UUID PRIMARY KEY,
    branch_id UUID NOT NULL REFERENCES branch(id),
    supplier_id UUID NOT NULL REFERENCES supplier(id),
    date_time TIMESTAMP,
    status VARCHAR(20), 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- income_product
CREATE TABLE income_product (
    id UUID PRIMARY KEY,
    income_id UUID NOT NULL REFERENCES income(id),
    category_id UUID NOT NULL REFERENCES category(id),
    product_name VARCHAR(255),
    barcode VARCHAR(50),
    quantity BIGINT,
    income_price DECIMAL(10, 2), 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- remainder
CREATE TABLE remainder (
    id UUID PRIMARY KEY,
    branch_id UUID NOT NULL REFERENCES branch(id),
    category_id UUID NOT NULL REFERENCES category(id),
    product_name VARCHAR(255),
    barcode VARCHAR(50),
    price_income DECIMAL(10, 2),
    quantity INT, 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);


-- shift(смена)
CREATE TABLE shift (
    id uuid PRIMARY KEY,
    branch_id uuid  NOT NULL REFERENCES branch(id),
    employee_id uuid  not null REFERENCES employee(id),
    sale_point uuid not null REFERENCES sale_point(id),
    status VARCHAR(20) CHECK (Status IN ('New', 'Open', 'Closed')),
    open_shift  TIMESTAMP,
    close_shift TIMESTAMP, 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);


-- sale
CREATE TABLE sale (
    id UUID PRIMARY KEY,
    sale_id VARCHAR(255), --  SD-000001
    branch_id UUID NOT NULL REFERENCES branch(id),
    salepoint_id UUID NOT NULL REFERENCES sale_point(id),
    shift_id UUID NOT NULL REFERENCES shift(id),
    employee_id UUID NOT NULL REFERENCES employee(id),
    barcode VARCHAR(50),
    status VARCHAR(20), 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);



-- sale_products
CREATE TABLE sale_products (
    id UUID PRIMARY KEY,
    sale_id UUID NOT NULL REFERENCES sale(id),
    category_id UUID NOT NULL REFERENCES category(id),
    product_name VARCHAR(255),
    barcode VARCHAR(50),
    remaining_quantity INT,
    quantity INT,
    allow_discount BOOLEAN,
    discount_type VARCHAR(20),
    discount DECIMAL(5, 2),
    price DECIMAL(10, 2),
    total_amount DECIMAL(10, 2), 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);


CREATE TABLE payment (
    id uuid PRIMARY KEY,
    sale_id  uuid NOT NULL REFERENCES sale(id),
    cash DECIMAL(10, 2),
    uzcard DECIMAL(10, 2),
    payme DECIMAL(10, 2),
    click DECIMAL(10, 2),
    humo DECIMAL(10, 2),
    apelsin DECIMAL(10, 2),
    total_amount DECIMAL(10, 2), 

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- transaction
CREATE TABLE transaction (
    id UUID PRIMARY KEY,
    shift_id UUID NOT NULL REFERENCES shift(id),
    cash  DECIMAL(10, 2) DEFAULT 0,
    uzcard  DECIMAL(10, 2) DEFAULT 0,
    payme  DECIMAL(10, 2) DEFAULT 0,
    click  DECIMAL(10, 2) DEFAULT 0,
    humo  DECIMAL(10, 2) DEFAULT 0,
    apelsin  DECIMAL(10, 2) DEFAULT 0,
    total_amount DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

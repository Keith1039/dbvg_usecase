CREATE TABLE IF NOT EXISTS USERS(
    USERNAME VARCHAR(255),
    PASSWORD VARCHAR(255),
    PRIMARY KEY (USERNAME),
    CREATED_AT TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS PRODUCTS(
    ID uuid PRIMARY KEY,
    NAME VARCHAR(255),
    DESCRIPTION VARCHAR(255),
    PRICE Decimal(10, 2),
    CREATED_AT TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS PURCHASES(
    USERNAME VARCHAR(255),
    PRODUCT_ID UUID,
    FOREIGN KEY (USERNAME) REFERENCES USERS,
    FOREIGN KEY (PRODUCT_ID) REFERENCES PRODUCTS,
    PRIMARY KEY (USERNAME, PRODUCT_ID)
)
IF NOT EXISTS (SELECT * FROM sys.databases WHERE name = 'shop')
BEGIN
    CREATE DATABASE shop;
END
GO

USE shop;
GO

IF OBJECT_ID('ProductCatalog', 'U') IS NULL
BEGIN
    CREATE TABLE ProductCatalog
    (
        product_id INT IDENTITY(1,1) PRIMARY KEY,
        name NVARCHAR(255) NOT NULL,
        description NVARCHAR(255) NOT NULL,
        price DECIMAL(10, 2) NOT NULL,
        quantity_stock INT NOT NULL
    );
END
ELSE
BEGIN
    IF COL_LENGTH('ProductCatalog', 'quantity_stock') IS NULL
    BEGIN
        ALTER TABLE ProductCatalog
        ADD quantity_stock INT NOT NULL DEFAULT 0;
    END
END
GO


IF OBJECT_ID('ShoppingCart', 'U') IS NULL
BEGIN
    CREATE TABLE ShoppingCart
    (
        cart_id INT IDENTITY(1,1) PRIMARY KEY,
        created_at DATETIME DEFAULT GETDATE()
    );
END
GO

IF OBJECT_ID('CartItem', 'U') IS NULL
BEGIN
    CREATE TABLE CartItem
    (
        cart_item_id INT IDENTITY(1,1) PRIMARY KEY,
        cart_id INT FOREIGN KEY REFERENCES ShoppingCart(cart_id),
        product_id INT FOREIGN KEY REFERENCES ProductCatalog(product_id),
        quantity INT NOT NULL
    );
END
GO

IF OBJECT_ID('Orders', 'U') IS NULL
BEGIN
    CREATE TABLE Orders
    (
        order_id INT IDENTITY(1,1) PRIMARY KEY,
        total_cost DECIMAL(10, 2) NOT NULL
    );
END
GO

IF OBJECT_ID('OrderItems', 'U') IS NULL
BEGIN
    CREATE TABLE OrderItems
    (
        order_item_id INT IDENTITY(1,1) PRIMARY KEY,
        order_id INT FOREIGN KEY REFERENCES Orders(order_id),
        product_id INT FOREIGN KEY REFERENCES ProductCatalog(product_id),
        quantity INT NOT NULL,
        price_at_purchase DECIMAL(10, 2) NOT NULL
    );
END
GO

CREATE DATABASE IF NOT EXISTS shop;
USE shop;

CREATE TABLE IF NOT EXISTS ProductCatalog (
    product_id INT AUTO_INCREMENT PRIMARY KEY,
    name NVARCHAR(255) NOT NULL,
    description NVARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    quantity_stock INT NOT NULL DEFAULT 0
);


INSERT INTO ProductCatalog (name, description, price, quantity_stock) VALUES
('Oculos VR', 'O futuro já é hoje', 4000.00, 10),
('Fones de Ouvido', 'Som limpo e satisfatório', 200.00, 100),
('Celular', 'A nova tecnologia em suas mãos', 2000.00, 50)
ON DUPLICATE KEY UPDATE
    name = VALUES(name),
    description = VALUES(description),
    price = VALUES(price),
    quantity_stock = VALUES(quantity_stock);


CREATE TABLE IF NOT EXISTS ShoppingCart (
    cart_id INT AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE IF NOT EXISTS CartItem (
    cart_item_id INT AUTO_INCREMENT PRIMARY KEY,
    cart_id INT,
    product_id INT,
    quantity INT NOT NULL,
    FOREIGN KEY (cart_id) REFERENCES ShoppingCart(cart_id),
    FOREIGN KEY (product_id) REFERENCES ProductCatalog(product_id)
);

CREATE TABLE IF NOT EXISTS Orders (
    order_id INT AUTO_INCREMENT PRIMARY KEY,
    total_cost DECIMAL(10, 2) NOT NULL
);


CREATE TABLE IF NOT EXISTS OrderItems (
    order_item_id INT AUTO_INCREMENT PRIMARY KEY,
    order_id INT,
    product_id INT,
    quantity INT NOT NULL,
    price_at_purchase DECIMAL(10, 2) NOT NULL,
    FOREIGN KEY (order_id) REFERENCES Orders(order_id),
    FOREIGN KEY (product_id) REFERENCES ProductCatalog(product_id)
);


INSERT INTO ShoppingCart (created_at) VALUES (NOW());


INSERT INTO CartItem (cart_id, product_id, quantity) VALUES
(1, 1, 2),
(1, 3, 1);


INSERT INTO Orders (total_cost) VALUES (3200.00); 


INSERT INTO OrderItems (order_id, product_id, quantity, price_at_purchase) VALUES
(1, 1, 2, 2000.00), 
(1, 3, 1, 200.00); 

SELECT * FROM ProductCatalog;
SELECT * FROM ShoppingCart;
SELECT * FROM CartItem;
SELECT * FROM Orders;
SELECT * FROM OrderItems;

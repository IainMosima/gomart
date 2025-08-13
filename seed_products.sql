-- Product seed data for GoMart
INSERT INTO products (product_id, product_name, description, price, sku, stock_quantity, category_id, is_active, created_at, updated_at, is_deleted) VALUES

-- FOOD & BEVERAGES (4 products)
('a1000000-1111-2222-3333-444444444001', 'Premium Kenyan Mangoes', 'Sweet and juicy mangoes from Machakos, 1kg pack', 250.00, 'FRUIT-MANGO-001', 50, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789012', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444002', 'Organic Sukuma Wiki', 'Fresh collard greens grown organically, 500g bundle', 50.00, 'VEG-SUKUMA-001', 100, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789013', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444003', 'Brookside Fresh Milk', 'Long life milk 1 liter carton', 120.00, 'DAIRY-MILK-001', 75, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789014', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444004', 'Kenyan AA Coffee Beans', 'Premium coffee beans from Kiambu, 500g pack', 850.00, 'BEV-COFFEE-001', 25, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789015', TRUE, NOW(), NULL, FALSE),

-- ELECTRONICS (4 products)
('a1000000-1111-2222-3333-444444444005', 'Samsung Galaxy A54', 'Mid-range smartphone with 128GB storage and 6GB RAM', 35000.00, 'PHONE-SAMSUNG-A54', 15, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789021', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444006', 'MacBook Air M2', 'Apple MacBook Air with M2 chip, 13-inch, 256GB SSD', 170000.00, 'LAPTOP-MAC-AIR-M2', 5, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789022', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444007', 'Sony WH-1000XM4', 'Wireless noise-canceling over-ear headphones', 28000.00, 'AUDIO-SONY-WH1000XM4', 12, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789024', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444008', 'Samsung 55" QLED TV', '55-inch 4K QLED Smart TV with HDR support', 95000.00, 'TV-SAMSUNG-55QLED', 8, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789025', TRUE, NOW(), NULL, FALSE),

-- HOME & GARDEN (3 products)
('a1000000-1111-2222-3333-444444444009', 'Ramtons Microwave Oven', '25L digital microwave with grill function', 12500.00, 'KITCHEN-RAMTONS-MW25', 20, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789031', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444010', 'Mahogany Dining Table', 'Solid mahogany 6-seater dining table, handcrafted', 45000.00, 'FURN-MAHOGANY-TABLE6', 3, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789032', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444011', 'Garden Tool Set', 'Complete gardening kit with spade, rake, pruners', 3500.00, 'GARDEN-TOOLSET-001', 15, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789033', TRUE, NOW(), NULL, FALSE),

-- FASHION & BEAUTY (2 products)
('a1000000-1111-2222-3333-444444444012', 'Kenyan Cotton Shirt', 'Locally made cotton shirt, various sizes and colors', 2800.00, 'MENS-COTTON-SHIRT-001', 30, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789042', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444013', 'Shea Butter Moisturizer', 'Natural shea butter body moisturizer, 250ml', 1200.00, 'BEAUTY-SHEA-MOIST-001', 40, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789045', TRUE, NOW(), NULL, FALSE),

-- BOOKS & MEDIA (2 products)
('a1000000-1111-2222-3333-444444444014', 'Grade 5 Mathematics Textbook', 'Kenya Institute of Education approved mathematics textbook', 650.00, 'BOOK-MATH-GRADE5-001', 60, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789052', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444015', 'The River and the Source', 'Classic Kenyan novel by Margaret Ogola', 450.00, 'BOOK-FICTION-RIVER-001', 35, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789053', TRUE, NOW(), NULL, FALSE);
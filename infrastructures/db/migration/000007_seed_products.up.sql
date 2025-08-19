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

-- FASHION & BEAUTY (4 products)
('a1000000-1111-2222-3333-444444444012', 'Nike Air Max Sneakers', 'Comfortable running shoes for men, size 42', 8500.00, 'SHOES-NIKE-AIRMAX-42', 25, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789044', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444013', 'L''Oreal Paris Foundation', 'Long-lasting liquid foundation, medium shade', 1800.00, 'BEAUTY-LOREAL-FOUND-MED', 40, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789043', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444014', 'Men''s Cotton T-Shirt', 'Premium quality cotton t-shirt, various colors', 1200.00, 'MENS-TSHIRT-COTTON-L', 60, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789041', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444015', 'Women''s Summer Dress', 'Elegant floral print dress, perfect for summer', 3500.00, 'WOMENS-DRESS-FLORAL-M', 20, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789042', TRUE, NOW(), NULL, FALSE),

-- BOOKS & MEDIA (4 products)
('a1000000-1111-2222-3333-444444444016', 'The Alchemist by Paulo Coelho', 'Bestselling novel about following your dreams', 1200.00, 'BOOK-ALCHEMIST-COELHO', 30, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789051', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444017', 'Rich Dad Poor Dad', 'Personal finance book by Robert Kiyosaki', 1500.00, 'BOOK-RICHDAD-KIYOSAKI', 25, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789052', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444018', 'The Lion King DVD', 'Classic Disney animated movie on DVD', 800.00, 'DVD-LIONKING-DISNEY', 15, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789053', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444019', 'Learn Python Programming', 'Comprehensive guide to Python programming', 2500.00, 'BOOK-PYTHON-GUIDE', 20, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789054', TRUE, NOW(), NULL, FALSE),

-- GAMING (3 products for deeper categories)
('a1000000-1111-2222-3333-444444444020', 'PlayStation 5 Console', 'Latest generation gaming console from Sony', 75000.00, 'GAMING-PS5-CONSOLE', 5, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789063', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444021', 'Xbox Series X', 'Microsoft''s flagship gaming console', 70000.00, 'GAMING-XBOX-SERIESX', 7, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789064', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444022', 'Nintendo Switch OLED', 'Handheld gaming console with OLED screen', 45000.00, 'GAMING-SWITCH-OLED', 10, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789065', TRUE, NOW(), NULL, FALSE),

-- TROPICAL FRUITS (deeper category)
('a1000000-1111-2222-3333-444444444023', 'Fresh Pineapples', 'Sweet golden pineapples from Coast region', 200.00, 'FRUIT-PINEAPPLE-001', 35, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789061', TRUE, NOW(), NULL, FALSE),
('a1000000-1111-2222-3333-444444444024', 'Passion Fruits', 'Tangy passion fruits, 1kg pack', 300.00, 'FRUIT-PASSION-001', 40, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789061', TRUE, NOW(), NULL, FALSE),

-- CITRUS FRUITS (deeper category)
('a1000000-1111-2222-3333-444444444025', 'Fresh Oranges', 'Juicy oranges from Central Kenya, 2kg pack', 180.00, 'FRUIT-ORANGE-001', 60, 'f1e2d3c4-b5a6-9c8d-7e6f-123456789062', TRUE, NOW(), NULL, FALSE);
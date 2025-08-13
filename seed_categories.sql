-- Category hierarchy for GoMart marketplace

-- Clean start - remove any existing data
DELETE FROM products WHERE is_deleted = FALSE;
DELETE FROM categories WHERE is_deleted = FALSE;

-- Main categories (root level)
INSERT INTO categories (category_id, category_name, parent_id, created_at, updated_at, is_deleted) VALUES
('f1e2d3c4-b5a6-9c8d-7e6f-123456789001', 'Food & Beverages', NULL, NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789002', 'Electronics', NULL, NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789003', 'Home & Garden', NULL, NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789004', 'Fashion & Beauty', NULL, NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789005', 'Books & Media', NULL, NOW(), NULL, FALSE);

-- Food & Beverages subcategories
INSERT INTO categories (category_id, category_name, parent_id, created_at, updated_at, is_deleted) VALUES
('f1e2d3c4-b5a6-9c8d-7e6f-123456789011', 'Fresh Produce', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789001', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789012', 'Fruits', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789011', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789013', 'Vegetables', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789011', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789014', 'Dairy & Eggs', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789001', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789015', 'Beverages', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789001', NOW(), NULL, FALSE);

-- Electronics subcategories
INSERT INTO categories (category_id, category_name, parent_id, created_at, updated_at, is_deleted) VALUES
('f1e2d3c4-b5a6-9c8d-7e6f-123456789021', 'Smartphones', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789002', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789022', 'Laptops', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789002', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789023', 'Audio & Video', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789002', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789024', 'Headphones', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789023', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789025', 'Televisions', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789023', NOW(), NULL, FALSE);

-- Home & Garden subcategories
INSERT INTO categories (category_id, category_name, parent_id, created_at, updated_at, is_deleted) VALUES
('f1e2d3c4-b5a6-9c8d-7e6f-123456789031', 'Kitchen Appliances', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789003', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789032', 'Furniture', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789003', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789033', 'Garden Tools', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789003', NOW(), NULL, FALSE);

-- Fashion & Beauty subcategories
INSERT INTO categories (category_id, category_name, parent_id, created_at, updated_at, is_deleted) VALUES
('f1e2d3c4-b5a6-9c8d-7e6f-123456789041', 'Clothing', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789004', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789042', 'Mens Wear', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789041', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789043', 'Womens Wear', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789041', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789044', 'Beauty Products', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789004', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789045', 'Skincare', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789044', NOW(), NULL, FALSE);

-- Books & Media subcategories
INSERT INTO categories (category_id, category_name, parent_id, created_at, updated_at, is_deleted) VALUES
('f1e2d3c4-b5a6-9c8d-7e6f-123456789051', 'Educational Books', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789005', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789052', 'Primary School', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789051', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789053', 'Fiction Books', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789005', NOW(), NULL, FALSE);
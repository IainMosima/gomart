-- Category hierarchy for GoMart marketplace

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
('f1e2d3c4-b5a6-9c8d-7e6f-123456789015', 'Beverages', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789001', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789016', 'Snacks & Confectionery', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789001', NOW(), NULL, FALSE);

-- Electronics subcategories
INSERT INTO categories (category_id, category_name, parent_id, created_at, updated_at, is_deleted) VALUES
('f1e2d3c4-b5a6-9c8d-7e6f-123456789021', 'Mobile Phones', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789002', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789022', 'Computers & Laptops', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789002', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789023', 'Gaming', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789002', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789024', 'Audio & Video', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789002', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789025', 'TV & Home Entertainment', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789002', NOW(), NULL, FALSE);

-- Home & Garden subcategories
INSERT INTO categories (category_id, category_name, parent_id, created_at, updated_at, is_deleted) VALUES
('f1e2d3c4-b5a6-9c8d-7e6f-123456789031', 'Kitchen Appliances', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789003', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789032', 'Furniture', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789003', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789033', 'Garden Tools', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789003', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789034', 'Home Decor', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789003', NOW(), NULL, FALSE);

-- Fashion & Beauty subcategories
INSERT INTO categories (category_id, category_name, parent_id, created_at, updated_at, is_deleted) VALUES
('f1e2d3c4-b5a6-9c8d-7e6f-123456789041', 'Men''s Fashion', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789004', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789042', 'Women''s Fashion', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789004', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789043', 'Beauty & Personal Care', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789004', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789044', 'Shoes & Accessories', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789004', NOW(), NULL, FALSE);

-- Books & Media subcategories
INSERT INTO categories (category_id, category_name, parent_id, created_at, updated_at, is_deleted) VALUES
('f1e2d3c4-b5a6-9c8d-7e6f-123456789051', 'Fiction Books', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789005', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789052', 'Non-Fiction Books', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789005', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789053', 'Movies & Music', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789005', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789054', 'Educational Materials', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789005', NOW(), NULL, FALSE);

-- Third-level categories for more depth
INSERT INTO categories (category_id, category_name, parent_id, created_at, updated_at, is_deleted) VALUES
-- Fruits subcategories
('f1e2d3c4-b5a6-9c8d-7e6f-123456789061', 'Tropical Fruits', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789012', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789062', 'Citrus Fruits', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789012', NOW(), NULL, FALSE),
-- Gaming subcategories
('f1e2d3c4-b5a6-9c8d-7e6f-123456789063', 'PlayStation', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789023', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789064', 'Xbox', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789023', NOW(), NULL, FALSE),
('f1e2d3c4-b5a6-9c8d-7e6f-123456789065', 'Nintendo', 'f1e2d3c4-b5a6-9c8d-7e6f-123456789023', NOW(), NULL, FALSE);
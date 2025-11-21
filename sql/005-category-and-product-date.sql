INSERT INTO categories (id, code, name, created_at, updated_at)
VALUES
    (DEFAULT, 'CAT001', 'Clothing', DEFAULT, DEFAULT),
    (DEFAULT, 'CAT002', 'Shoes', DEFAULT, DEFAULT),
    (DEFAULT, 'CAT003', 'Accessories', DEFAULT, DEFAULT);

Update products
set category_id = (SELECT id FROM categories WHERE name = 'Clothing')
where code in ('PROD001', 'PROD004', 'PROD007');

Update products
set category_id = (SELECT id FROM categories WHERE name = 'Shoes')
where code in ('PROD002', 'PROD006');

Update products
set category_id = (SELECT id FROM categories WHERE name = 'Accessories')
where code in ('PROD003', 'PROD005', 'PROD008');
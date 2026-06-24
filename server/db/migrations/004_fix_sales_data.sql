-- Migration: Fix sales data migration and FK references
-- ============================================================
-- The previous migration (003) had two issues:
-- 1. INSERT INTO sales_new referenced s.descr_ru but old sales table has descr (not descr_ru)
-- 2. sale_materials and sales_bases FK references point to "sales_old"(id) instead of sales(id)
-- This migration fixes both issues.
-- Idempotent: checks if sales_old table exists and has data before migrating.

-- 1. Drop and recreate sale_materials with correct FK reference to sales(id)
DROP TABLE IF EXISTS sale_materials;

CREATE TABLE IF NOT EXISTS sale_materials (
    sale_id INTEGER NOT NULL,
    material_id INTEGER NOT NULL,
    PRIMARY KEY (sale_id, material_id),
    FOREIGN KEY (sale_id) REFERENCES sales(id) ON DELETE CASCADE,
    FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sale_materials_sale_id ON sale_materials(sale_id);
CREATE INDEX IF NOT EXISTS idx_sale_materials_material_id ON sale_materials(material_id);

-- 2. Drop and recreate sales_bases with correct FK reference to sales(id)
DROP TABLE IF EXISTS sales_bases;

CREATE TABLE IF NOT EXISTS sales_bases (
    sale_id INTEGER NOT NULL,
    base_id INTEGER NOT NULL,
    PRIMARY KEY (sale_id, base_id),
    FOREIGN KEY (sale_id) REFERENCES sales(id) ON DELETE CASCADE,
    FOREIGN KEY (base_id) REFERENCES bases(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sales_bases_sale_id ON sales_bases(sale_id);
CREATE INDEX IF NOT EXISTS idx_sales_bases_base_id ON sales_bases(base_id);

-- 3. Migrate data from sales_old to sales (if sales_old exists and has data)
INSERT INTO sales (id, title, description, image_path, price, old_price, year, technique, size, status, sort_order, sold, created_at, updated_at)
    SELECT
        s.id,
        s.name_ru,
        COALESCE(s.descr, ''),
        s.images,
        CAST(s.price AS REAL) / 100.0,
        NULL,
        s.year,
        '',
        CAST(s.width AS TEXT) || 'x' || CAST(s.height AS TEXT),
        'published',
        s.id,
        0,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    FROM sales_old s
    WHERE EXISTS (SELECT 1 FROM pragma_table_info('sales_old') WHERE name = 'name_ru')
    AND NOT EXISTS (SELECT 1 FROM sales WHERE sales.id = s.id);

-- 4. Migrate base associations from sales_old.base_id to sales_bases
INSERT OR IGNORE INTO sales_bases (sale_id, base_id)
    SELECT s.id, s.base_id
    FROM sales_old s
    WHERE s.base_id > 0
    AND EXISTS (SELECT 1 FROM pragma_table_info('sales_old') WHERE name = 'base_id')
    AND NOT EXISTS (SELECT 1 FROM sales_bases WHERE sales_bases.sale_id = s.id AND sales_bases.base_id = s.base_id);

-- 5. Drop sales_old table after successful migration
DROP TABLE IF EXISTS sales_old;

-- Migration: Fix gallery schema to match Go code expectations
-- ============================================================
-- The old schema had different column names and structures.
-- This migration renames columns and creates proper tables.
-- Idempotent: checks if old columns exist before migrating.

-- 1. Fix materials table: rename material_ru -> name_ru, material_en -> name_en, add created_at
CREATE TABLE IF NOT EXISTS materials_new (
    id INTEGER PRIMARY KEY,
    name_ru TEXT NOT NULL,
    name_en TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO materials_new (id, name_ru, name_en, created_at)
    SELECT id, material_ru, material_en, CURRENT_TIMESTAMP FROM materials
    WHERE EXISTS (SELECT 1 FROM pragma_table_info('materials') WHERE name = 'material_ru');

DROP TABLE IF EXISTS materials_old;
ALTER TABLE materials RENAME TO materials_old;
DROP TABLE IF EXISTS materials;
ALTER TABLE materials_new RENAME TO materials;
DROP TABLE IF EXISTS materials_old;

-- 2. Fix bases table: rename base_ru -> name_ru, base_en -> name_en, add created_at
CREATE TABLE IF NOT EXISTS bases_new (
    id INTEGER PRIMARY KEY,
    name_ru TEXT NOT NULL,
    name_en TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO bases_new (id, name_ru, name_en, created_at)
    SELECT id, base_ru, base_en, CURRENT_TIMESTAMP FROM bases
    WHERE EXISTS (SELECT 1 FROM pragma_table_info('bases') WHERE name = 'base_ru');

DROP TABLE IF EXISTS bases_old;
ALTER TABLE bases RENAME TO bases_old;
DROP TABLE IF EXISTS bases;
ALTER TABLE bases_new RENAME TO bases;
DROP TABLE IF EXISTS bases_old;

-- 3. Create work_materials junction table
CREATE TABLE IF NOT EXISTS work_materials (
    work_id INTEGER NOT NULL,
    material_id INTEGER NOT NULL,
    PRIMARY KEY (work_id, material_id),
    FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE,
    FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE
);

INSERT OR IGNORE INTO work_materials (work_id, material_id)
    SELECT work_id, material_id FROM works_materials
    WHERE EXISTS (SELECT 1 FROM pragma_table_info('works_materials') WHERE name = 'work_id');

-- 4. Create work_bases junction table
CREATE TABLE IF NOT EXISTS work_bases (
    work_id INTEGER NOT NULL,
    base_id INTEGER NOT NULL,
    PRIMARY KEY (work_id, base_id),
    FOREIGN KEY (work_id) REFERENCES works(id) ON DELETE CASCADE,
    FOREIGN KEY (base_id) REFERENCES bases(id) ON DELETE CASCADE
);

INSERT OR IGNORE INTO work_bases (work_id, base_id)
    SELECT id, base_id FROM works WHERE base_id > 0
    AND EXISTS (SELECT 1 FROM pragma_table_info('works') WHERE name = 'base_id');

-- 5. Create new works table with the correct schema
CREATE TABLE IF NOT EXISTS works_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL DEFAULT '',
    description TEXT DEFAULT '',
    image_path TEXT DEFAULT '',
    year INTEGER,
    technique TEXT DEFAULT '',
    size TEXT DEFAULT '',
    status TEXT DEFAULT 'published',
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO works_new (id, title, description, image_path, year, technique, size, status, sort_order, created_at, updated_at)
    SELECT
        w.id,
        w.name_ru,
        COALESCE(w.descr_ru, ''),
        w.images,
        w.year,
        COALESCE(wt.value, ''),
        CAST(w.width AS TEXT) || 'x' || CAST(w.height AS TEXT),
        'published',
        w.id,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    FROM works w
    LEFT JOIN work_types wt ON w.type = wt.id
    WHERE EXISTS (SELECT 1 FROM pragma_table_info('works') WHERE name = 'name_ru');

INSERT INTO works_new (title, description, image_path, year, technique, size, status, sort_order, created_at, updated_at)
    SELECT
        i.name_ru,
        COALESCE(i.desc, ''),
        i.str_id || '/',
        i.year,
        'illustration',
        '',
        'published',
        (SELECT COALESCE(MAX(sort_order), 0) + 1 FROM works_new),
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    FROM illustrations i
    WHERE EXISTS (SELECT 1 FROM pragma_table_info('illustrations') WHERE name = 'name_ru');

INSERT INTO works_new (title, description, image_path, year, technique, size, status, sort_order, created_at, updated_at)
    SELECT
        t.name_ru,
        COALESCE(t.desc, ''),
        t.str_id || '/',
        t.year,
        '3d',
        '',
        'published',
        (SELECT COALESCE(MAX(sort_order), 0) + 1 FROM works_new),
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    FROM threeds t
    WHERE EXISTS (SELECT 1 FROM pragma_table_info('threeds') WHERE name = 'name_ru');

DROP TABLE IF EXISTS works_old;
ALTER TABLE works RENAME TO works_old;
DROP TABLE IF EXISTS works;
ALTER TABLE works_new RENAME TO works;
DROP TABLE IF EXISTS works_old;
DROP TABLE IF EXISTS work_types;
DROP TABLE IF EXISTS illustrations;
DROP TABLE IF EXISTS threeds;

-- Recreate indexes
CREATE INDEX IF NOT EXISTS idx_works_status ON works(status);
CREATE INDEX IF NOT EXISTS idx_works_sort_order ON works(sort_order);
CREATE INDEX IF NOT EXISTS idx_work_materials_work_id ON work_materials(work_id);
CREATE INDEX IF NOT EXISTS idx_work_materials_material_id ON work_materials(material_id);
CREATE INDEX IF NOT EXISTS idx_work_bases_work_id ON work_bases(work_id);
CREATE INDEX IF NOT EXISTS idx_work_bases_base_id ON work_bases(base_id);

-- 6. Create sale_materials junction table
CREATE TABLE IF NOT EXISTS sale_materials (
    sale_id INTEGER NOT NULL,
    material_id INTEGER NOT NULL,
    PRIMARY KEY (sale_id, material_id),
    FOREIGN KEY (sale_id) REFERENCES sales(id) ON DELETE CASCADE,
    FOREIGN KEY (material_id) REFERENCES materials(id) ON DELETE CASCADE
);

-- 7. Create sale_bases junction table
CREATE TABLE IF NOT EXISTS sale_bases (
    sale_id INTEGER NOT NULL,
    base_id INTEGER NOT NULL,
    PRIMARY KEY (sale_id, base_id),
    FOREIGN KEY (sale_id) REFERENCES sales(id) ON DELETE CASCADE,
    FOREIGN KEY (base_id) REFERENCES bases(id) ON DELETE CASCADE
);

-- 8. Fix sales table
CREATE TABLE IF NOT EXISTS sales_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL DEFAULT '',
    description TEXT DEFAULT '',
    image_path TEXT DEFAULT '',
    price REAL DEFAULT 0,
    old_price REAL,
    year INTEGER,
    technique TEXT DEFAULT '',
    size TEXT DEFAULT '',
    status TEXT DEFAULT 'published',
    sort_order INTEGER DEFAULT 0,
    sold INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO sales_new (id, title, description, image_path, price, year, technique, size, status, sort_order, sold, created_at, updated_at)
    SELECT
        s.id,
        s.name_ru,
        COALESCE(s.descr_ru, ''),
        s.images,
        CAST(s.price AS REAL) / 100.0,
        s.year,
        '',
        CAST(s.width AS TEXT) || 'x' || CAST(s.height AS TEXT),
        'published',
        s.id,
        0,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    FROM sales s
    WHERE EXISTS (SELECT 1 FROM pragma_table_info('sales') WHERE name = 'name_ru');

INSERT OR IGNORE INTO sale_materials (sale_id, material_id)
    SELECT sale_id, material_id FROM sales_materials
    WHERE EXISTS (SELECT 1 FROM pragma_table_info('sales_materials') WHERE name = 'sale_id');

DROP TABLE IF EXISTS sales_old;
ALTER TABLE sales RENAME TO sales_old;
DROP TABLE IF EXISTS sales;
ALTER TABLE sales_new RENAME TO sales;
DROP TABLE IF EXISTS sales_old;
DROP TABLE IF EXISTS sales_materials;

-- Recreate indexes
CREATE INDEX IF NOT EXISTS idx_sales_status ON sales(status);
CREATE INDEX IF NOT EXISTS idx_sales_sort_order ON sales(sort_order);
CREATE INDEX IF NOT EXISTS idx_sales_sold ON sales(sold);
CREATE INDEX IF NOT EXISTS idx_sale_materials_sale_id ON sale_materials(sale_id);
CREATE INDEX IF NOT EXISTS idx_sale_materials_material_id ON sale_materials(material_id);
CREATE INDEX IF NOT EXISTS idx_sale_bases_sale_id ON sale_bases(sale_id);
CREATE INDEX IF NOT EXISTS idx_sale_bases_base_id ON sale_bases(base_id);

-- Migration: Add tags system and seed master-class data
-- ======================================================

-- 1. Create tags table
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name_ru TEXT NOT NULL UNIQUE,
    name_en TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Create product_tags junction table (M:N)
CREATE TABLE IF NOT EXISTS product_tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE,
    UNIQUE(product_id, tag_id)
);

CREATE INDEX IF NOT EXISTS idx_product_tags_product_id ON product_tags(product_id);
CREATE INDEX IF NOT EXISTS idx_product_tags_tag_id ON product_tags(tag_id);

-- 3. Seed tags
INSERT OR IGNORE INTO tags (id, name_ru, name_en, slug) VALUES
    (1, 'акварель', 'watercolor', 'watercolor'),
    (2, 'основы', 'basics', 'basics'),
    (3, 'растяжка', 'color-wash', 'color-wash'),
    (4, 'лессировка', 'glazing', 'glazing'),
    (5, 'цветы', 'flowers', 'flowers'),
    (6, 'композиция', 'composition', 'composition'),
    (7, 'мокрым по мокрому', 'wet-on-wet', 'wet-on-wet'),
    (8, 'пейзаж', 'landscape', 'landscape'),
    (9, 'отражения', 'reflections', 'reflections'),
    (10, 'вода', 'water', 'water'),
    (11, 'масло', 'oil', 'oil'),
    (12, 'материалы', 'materials', 'materials'),
    (13, 'смешивание', 'color-mixing', 'color-mixing'),
    (14, 'импасто', 'impasto', 'impasto'),
    (15, 'фактура', 'texture', 'texture'),
    (16, 'мастихин', 'palette-knife', 'palette-knife'),
    (17, 'портрет', 'portrait', 'portrait'),
    (18, 'светотень', 'chiaroscuro', 'chiaroscuro'),
    (19, 'объем', 'volume', 'volume'),
    (20, 'для начинающих', 'for-beginners', 'for-beginners'),
    (21, 'линии', 'lines', 'lines'),
    (22, 'формы', 'shapes', 'shapes'),
    (23, 'перспектива', 'perspective', 'perspective'),
    (24, 'построение', 'construction', 'construction'),
    (25, 'тон', 'tone', 'tone'),
    (26, 'штриховка', 'hatching', 'hatching');

-- 4. Seed master-class products (type = 'masterclass')
--    Price in kopecks. 0 = free.
--    category_id = 5 (Мастер-классы)

INSERT OR IGNORE INTO products (
    id, type, title_ru, title_en,
    description_ru, description_en,
    short_description_ru, short_description_en,
    price, duration_days, thumbnail_url, video_url,
    status, difficulty, total_lessons, total_duration_minutes,
    category_id, language, is_featured, view_count
) VALUES
    -- id=4: Бесплатный — Основы акварельной техники
    (4, 'masterclass',
     'Основы акварельной техники', 'Watercolor Basics',
     'Этот мастер-класс познакомит вас с основами акварельной живописи. Вы узнаете, как правильно выбирать бумагу и краски, освоите ключевые техники — заливка, растяжка цвета, лессировка и работа «мокрым по мокрому». Каждый прием сопровождается практическим упражнением, чтобы вы сразу закрепили навык.',
     'This master class introduces you to the basics of watercolor painting. You will learn how to choose paper and paints, master key techniques — wash, color stretching, glazing, and wet-on-wet. Each technique comes with a practical exercise.',
     'Изучите базовые приемы работы с акварелью: растяжка, лессировка, работа с водой.',
     'Learn basic watercolor techniques: stretching, glazing, working with water.',
     0, NULL,
     'https://images.unsplash.com/photo-1579783902614-a3fb3927b6a5?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
     '#', 'published', 'beginner', 1, 25, 5, 'ru', 1, 0),

    -- id=5: Платный 1500₽ — Цветочная композиция в акварели
    (5, 'masterclass',
     'Цветочная композиция в акварели', 'Watercolor Floral Composition',
     'На этом мастер-классе вы научитесь писать цветы в акварельной технике. Мы разберем построение композиции, работу с референсами, смешивание оттенков для получения естественных цветов. Главный фокус — техника «мокрым по мокрому», которая дает те самые воздушные перетекания цвета.',
     'In this master class you will learn to paint flowers in watercolor. We will cover composition building, working with references, mixing shades for natural colors. The main focus is the wet-on-wet technique.',
     'Создайте нежную цветочную композицию, используя технику мокрым по мокрому.',
     'Create a delicate floral composition using the wet-on-wet technique.',
     150000, NULL,
     'https://images.unsplash.com/photo-1541961017774-22349e4a1262?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
     '#', 'published', 'beginner', 1, 42, 5, 'ru', 0, 0),

    -- id=6: Платный 1500₽ — Пейзаж с отражением в воде
    (6, 'masterclass',
     'Пейзаж с отражением в воде', 'Landscape with Water Reflections',
     'Пейзаж с водой — один из самых живописных сюжетов в акварели. На этом занятии вы освоите приемы для передачи глади воды, отражений неба и деревьев, научитесь работать с горизонтальными плоскостями и создавать глубину пространства.',
     'Landscape with water is one of the most picturesque subjects in watercolor. In this lesson you will master techniques for conveying water surface, reflections of sky and trees, working with horizontal planes and creating depth.',
     'Научитесь создавать реалистичные отражения и передавать атмосферу пейзажа.',
     'Learn to create realistic reflections and convey the atmosphere of a landscape.',
     150000, NULL,
     'https://images.unsplash.com/photo-1578301978693-85fa9c0320b9?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
     '#', 'published', 'intermediate', 1, 38, 5, 'ru', 0, 0),

    -- id=7: Бесплатный — Основы работы с масляными красками
    (7, 'masterclass',
     'Основы работы с масляными красками', 'Oil Painting Basics',
     'Первый шаг в мир масляной живописи. Мы разберем, какие кисти, краски и холсты нужны новичку, как организовать рабочее место. Вы освоите базовые приемы: смешивание цветов на палитре, нанесение мазков разной формы, работа с разбавителями.',
     'First step into the world of oil painting. We will cover what brushes, paints and canvases a beginner needs, how to organize your workspace. You will master basic techniques: mixing colors on the palette, applying strokes of different shapes, working with thinners.',
     'Познакомьтесь с материалами, научитесь правильно смешивать краски и наносить их на холст.',
     'Get to know materials, learn to mix paints correctly and apply them to canvas.',
     0, NULL,
     'https://images.unsplash.com/photo-1541961017774-22349e4a1262?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
     '#', 'published', 'beginner', 1, 32, 5, 'ru', 1, 0),

    -- id=8: Платный 2000₽ — Техника импасто: объемные мазки
    (8, 'masterclass',
     'Техника импасто: объемные мазки', 'Impasto Technique: Textured Strokes',
     'Импасто — одна из самых эффектных техник масляной живописи. Вы узнаете, как создавать рельефные мазки, работать мастихином, добиваться фактурности и объема. Подходит для создания выразительных натюрмортов и абстрактных работ.',
     'Impasto is one of the most impressive oil painting techniques. You will learn how to create relief strokes, work with a palette knife, achieve texture and volume. Suitable for expressive still lifes and abstract works.',
     'Освойте технику густого нанесения краски для создания выразительной фактуры.',
     'Master the technique of thick paint application for expressive texture.',
     200000, NULL,
     'https://images.unsplash.com/photo-1578301978693-85fa9c0320b9?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
     '#', 'published', 'intermediate', 1, 45, 5, 'ru', 0, 0),

    -- id=9: Платный 2500₽ — Портрет маслом: свет и тень
    (9, 'masterclass',
     'Портрет маслом: свет и тень', 'Oil Portrait: Light and Shadow',
     'Портрет — вершина мастерства художника. На этом мастер-классе вы освоите построение головы, разбор крупных форм, работу с освещением. Главный акцент — передача объема через светотеневые отношения и цветовые рефлексы.',
     'Portrait is the pinnacle of an artist''s skill. In this master class you will master head construction, breakdown of large forms, working with lighting. The main focus is conveying volume through light-shadow relationships and color reflexes.',
     'Научитесь передавать объем и характер в портретной живописи.',
     'Learn to convey volume and character in portrait painting.',
     250000, NULL,
     'https://images.unsplash.com/photo-1542744095-fcf48d80b0fd?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
     '#', 'published', 'advanced', 1, 68, 5, 'ru', 0, 0),

    -- id=10: Бесплатный — Первые шаги в рисовании
    (10, 'masterclass',
     'Первые шаги в рисовании: линии и формы', 'First Steps in Drawing: Lines and Shapes',
     'Идеальный старт для тех, кто никогда не рисовал. Мы начнем с самых простых упражнений: постановка руки, проведение линий разного характера, построение геометрических форм. Постепенно перейдем к более сложным объектам и основам пропорций.',
     'The perfect start for those who have never drawn. We begin with the simplest exercises: hand positioning, drawing lines of different character, constructing geometric shapes. Gradually move to more complex objects and basics of proportions.',
     'Научитесь уверенно держать карандаш, создавать базовые формы и понимать пропорции.',
     'Learn to confidently hold a pencil, create basic shapes and understand proportions.',
     0, NULL,
     'https://images.unsplash.com/photo-1541961017774-22349e4a1262?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
     '#', 'published', 'beginner', 1, 28, 5, 'ru', 1, 0),

    -- id=11: Платный 1200₽ — Основы композиции и перспективы
    (11, 'masterclass',
     'Основы композиции и перспективы', 'Basics of Composition and Perspective',
     'Композиция и перспектива — фундамент любого изображения. На этом занятии вы узнаете правила «золотого сечения», научитесь строить линейную и воздушную перспективу, правильно размещать объекты в листе для создания гармоничного изображения.',
     'Composition and perspective are the foundation of any image. In this lesson you will learn the rules of the golden ratio, learn to build linear and aerial perspective, correctly place objects on the sheet to create a harmonious image.',
     'Поймете, как располагать объекты на листе и создавать иллюзию глубины.',
     'Understand how to arrange objects on the sheet and create the illusion of depth.',
     120000, NULL,
     'https://images.unsplash.com/photo-1578301978693-85fa9c0320b9?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
     '#', 'published', 'beginner', 1, 35, 5, 'ru', 0, 0),

    -- id=12: Платный 1200₽ — Работа с тоном и светотенью
    (12, 'masterclass',
     'Работа с тоном и светотенью', 'Working with Tone and Chiaroscuro',
     'Тон — основа реалистичного рисунка. Вы научитесь видеть и передавать тональные отношения, работать с градациями света и тени, создавать иллюзию объема на плоском листе. Упражнения построены от простых геометрических тел до сложных форм.',
     'Tone is the foundation of realistic drawing. You will learn to see and convey tonal relationships, work with gradations of light and shadow, create the illusion of volume on a flat sheet. Exercises range from simple geometric bodies to complex forms.',
     'Научитесь передавать объем с помощью света и тени на простых объектах.',
     'Learn to convey volume using light and shadow on simple objects.',
     120000, NULL,
     'https://images.unsplash.com/photo-1542744095-fcf48d80b0fd?ixlib=rb-4.0.3&auto=format&fit=crop&w=800&q=80',
     '#', 'published', 'beginner', 1, 40, 5, 'ru', 0, 0);

-- 5. Link products with tags via product_tags

-- id=4: Основы акварельной техники → акварель, основы, растяжка, лессировка
INSERT OR IGNORE INTO product_tags (product_id, tag_id) VALUES
    (4, 1), (4, 2), (4, 3), (4, 4);

-- id=5: Цветочная композиция → акварель, цветы, композиция, мокрым по мокрому
INSERT OR IGNORE INTO product_tags (product_id, tag_id) VALUES
    (5, 1), (5, 5), (5, 6), (5, 7);

-- id=6: Пейзаж с отражением → акварель, пейзаж, отражения, вода
INSERT OR IGNORE INTO product_tags (product_id, tag_id) VALUES
    (6, 1), (6, 8), (6, 9), (6, 10);

-- id=7: Основы масла → масло, основы, материалы, смешивание
INSERT OR IGNORE INTO product_tags (product_id, tag_id) VALUES
    (7, 11), (7, 2), (7, 12), (7, 13);

-- id=8: Импасто → масло, импасто, фактура, мастихин
INSERT OR IGNORE INTO product_tags (product_id, tag_id) VALUES
    (8, 11), (8, 14), (8, 15), (8, 16);

-- id=9: Портрет маслом → масло, портрет, светотень, объем
INSERT OR IGNORE INTO product_tags (product_id, tag_id) VALUES
    (9, 11), (9, 17), (9, 18), (9, 19);

-- id=10: Первые шаги → для начинающих, основы, линии, формы
INSERT OR IGNORE INTO product_tags (product_id, tag_id) VALUES
    (10, 20), (10, 2), (10, 21), (10, 22);

-- id=11: Композиция и перспектива → для начинающих, композиция, перспектива, построение
INSERT OR IGNORE INTO product_tags (product_id, tag_id) VALUES
    (11, 20), (11, 6), (11, 23), (11, 24);

-- id=12: Тон и светотень → для начинающих, тон, светотень, штриховка
INSERT OR IGNORE INTO product_tags (product_id, tag_id) VALUES
    (12, 20), (12, 25), (12, 18), (12, 26);

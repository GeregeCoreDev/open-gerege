-- +goose Up
-- +goose StatementBegin
SET search_path TO template_backend;

-- ============================================
-- SEED: SYSTEMS
-- ============================================
INSERT INTO systems (id, code, key, name, description, is_active, icon, sequence) VALUES
    (7, 'ADMIN', 'admin', 'Админ систем', 'Админ систем', true, 'i-lucide-monitor-cog', 1),
    (8, 'GRGB', 'grgb', 'Гэрэгэ бизнес', 'Гэрэгэ бизнес', true, 'i-lucide-briefcase-business', 4),
    (9, 'GRGA', 'grga', 'Гэрэгэ апп', 'Гэрэгэ Апп', true, 'i-lucide-smartphone', 3),
    (10, 'TPAY', 'tpay', 'Гэрэгэ хэтэвч', 'Гэрэгэ төлбөрийн систем', true, 'i-lucide-wallet', 2)
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('systems', 'id'), GREATEST((SELECT MAX(id) FROM systems), 1));

-- ============================================
-- SEED: ACTIONS (Default CRUD actions)
-- ============================================
INSERT INTO actions (id, code, name, description, is_active) VALUES
    (1, 'CREATE', 'Үүсгэх', 'Create action', true),
    (2, 'READ', 'Харах', 'Read action', true),
    (3, 'UPDATE', 'Засах', 'Update action', true),
    (4, 'DELETE', 'Устгах', 'Delete action', true)
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('actions', 'id'), GREATEST((SELECT MAX(id) FROM actions), 1));

-- ============================================
-- SEED: MODULES
-- ============================================
INSERT INTO modules (id, code, name, description, is_active, system_id) VALUES
    (7, 'ADMIN-SYSTEM', 'Систем', 'Систем', true, 7),
    (8, 'ADMIN-MODULE', 'Модул', 'Модул', true, 7),
    (10, 'ADMIN-ROLE', 'Дүр', 'Дүр', true, 7),
    (11, 'ADMIN-MODULE-PERMISSION', 'Эрхийн удирдлага', 'Эрхийн удирдлага', true, 7),
    (12, 'ADMIN-ORG', 'Байгууллага', 'Байгууллага', true, 7),
    (13, 'ADMIN-USER', 'Хэрэглэгч', 'Хэрэглэгч', true, 7),
    (14, 'GRGA-ICON', 'Апп айкон', 'Апп айкон', true, 9),
    (15, 'GRGB-ICON', 'Апп айкон', 'Апп айкон', true, 8),
    (16, 'ADMIN-DASHBOARD', 'Хянах самбар', 'Админ хянах самбар', true, 7),
    (17, 'ADMIN-ORGTYPE', 'Байгууллагын төрөл', 'Байгууллагын төрөл', true, 7),
    (18, 'TPAY-DASHBOARD', 'Хянах самбар', 'T-Pay хянах самбар', true, 10),
    (19, 'GRGA-DASHBOARD', 'Хянах самбар', 'Гэрэгэ апп хянах самбар', true, 9),
    (20, 'GRGB-DASHBOARD', 'Хянах самбар', 'Гэрэгэ бизнес хянах самбар', true, 8),
    (21, 'TPAY-WALLET', 'Хэтэвч', 'TPay хэтэвч', true, 10),
    (22, 'TPAY-TRX', 'Гүйлгээ', 'TPay гүйлгээ', true, 10),
    (23, 'TPAY-BALANCE', 'Баланс', 'TPay баланс', true, 10),
    (24, 'GRGB-AGENTS', 'Агент', 'Бизнес агентууд', true, 8),
    (25, 'TPAY-ROLE', 'Дүр', 'ТPay системийн дүрүүд', true, 10),
    (26, 'GRGA-ROLE', 'Дүр', 'Гэрэгэ аппын дүрүүд', true, 9),
    (27, 'GRGB-ROLE', 'Дүр', 'Гэрэгэ бизнесийн дүрүүд', true, 8),
    (28, 'ADMIN-AI-CHAT', 'Хиймэл оюун удирдлага', 'Хиймэл оюун удирдлага', true, 7),
    (29, 'TPAY-USERS', 'Хэрэглэгчид', '', true, 10),
    (30, 'TPAY-PRODUCTS', 'Бүтээгдэхүүний каталог', '', true, 10)
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('modules', 'id'), GREATEST((SELECT MAX(id) FROM modules), 1));

-- ============================================
-- SEED: ROLES
-- ============================================
INSERT INTO roles (id, name, description, system_id, code) VALUES
    (1, 'Сүпэр админ', 'Бүрэн эрхтэй админ', 7, 'ADMIN-SUPERADMIN'),
    (2, 'Систем админ', 'Гэрэгэ хэтэвчийн систем админ', 10, 'TPAY-SYSADMIN'),
    (9, 'Менежер', 'TPay системийн менежер', 10, 'TPAY-MANAGER'),
    (10, 'Систем админ', 'Гэрэгэ апп системийн админ', 9, 'GRGA-SYSADMIN'),
    (12, 'Менежер', 'Гэрэгэ апп системийн менежер', 9, 'GRGA-MANAGER'),
    (13, 'Менежер', 'Админ системийн менежер', 7, 'ADMIN-MANAGER'),
    (14, 'Админ', 'Гэрэгэ бизнес системийн админ', 8, 'GRGB-SYSADMIN'),
    (15, 'Хэрэглэгч', 'Гэрэгэ бизнес системийн хэрэглэгч', 8, 'GRGB-USER'),
    (17, 'Хэрэглэгч', 'Гэрэгэ апп системийн хэрэглэгч', 9, 'GRGA-USER')
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('roles', 'id'), GREATEST((SELECT MAX(id) FROM roles), 1));

-- ============================================
-- SEED: ORGANIZATION_TYPES
-- ============================================
INSERT INTO organization_types (id, name, code, description) VALUES
    (1, 'Компани', 'company', 'Компани'),
    (2, 'Төрийн байгууллага', 'government', 'Төрийн байгууллага'),
    (3, 'ТББ', 'ngo', 'ТББ'),
    (4, 'Суурь систем', 'core_system', 'Зөвхөн Гэрэгэ коор ХХК-д хамаатай'),
    (6, 'Гэрэгэ салбар компани', 'gerege', 'Гэрэгэ системс ХХК-ын хамааралтай компаниуд')
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('organization_types', 'id'), GREATEST((SELECT MAX(id) FROM organization_types), 1));

-- ============================================
-- SEED: CITIZENS
-- ============================================
INSERT INTO citizens (id, civil_id, reg_no, family_name, last_name, first_name, gender, birth_date, phone_no, email, is_foreign, country_code, aimag_name, sum_name, bag_name, address_detail) VALUES
    (10000141, 623458282539, 'йк98013033', 'тайж', 'отгонбаяр', 'баярсайхан', 1, '1998-01-30', '94422470', 'craftzbay@gmail.com', 0, 'mng', 'Улаанбаатар', 'Хан-Уул', '4-р хороо', 'Улаанбаатар, Хан-Уул, 4-р хороо'),
    (10000263, 111949212017, 'ма74101813', 'харчин', 'цэнддорж', 'эрдэнэбат', 1, '1974-10-18', '99118306', 'erdenebatt@gmail.com', 0, 'mng', 'Улаанбаатар', 'Хан-Уул', '11-р хороо', 'Улаанбаатар, Хан-Уул, 11-р хороо'),
    (10000221, 111036772547, 'уп96030116', 'алагадуун', 'төвшин', 'мөнхбаяр', 1, '1996-03-01', '94942244', 'munkhbayar@gerege.com', 0, 'mng', 'Архангай', 'Батцэнгэл', '1-р баг', 'Архангай, Батцэнгэл, 1-р баг'),
    (10000081, 110770722023, 'уз98052158', 'өнгөө', 'соронзонболд', 'сэнгүм', 1, '1998-05-21', '88995566', 'sengum@gerege.mn', 0, 'mng', 'Улаанбаатар', 'Хан-Уул', '3-р хороо', 'Эрчим хотхон 10/2, 44'),
    (10162877, 823494982818, 'дк02301734', 'жалайр', 'алимаа', 'хүдэрчулуун', 1, '2002-10-17', '89908953', 'khuderchuluun@gerege.com', 0, 'mng', 'Говь-Алтай', 'Халиун', '4-р баг', 'Говь-Алтай, Халиун, 4-р баг'),
    (11949042, 111676962728, 'та00261713', 'эвэн', 'хонгор', 'чингүүн', 1, '2000-06-17', '99500509', 'sharshuwuu@gmail.com', 0, 'mng', 'Улаанбаатар', 'Сүхбаатар', '1-р хороо', 'Улаанбаатар хот, Сүхбаатар дүүрэг'),
    (10000101, 220102982359, 'лю98082031', 'алагууд', 'батжаргал', 'батзориг', 1, '1998-08-20', '88002608', 'batzorig@gerege.com', 0, 'mng', 'Сүхбаатар', '', '', 'Сүхбаатар, Баруун-Урт'),
    (10000022, 111014502067, 'уп98101639', 'тангад', 'ганболд', 'ганхөлөг', 1, '1998-10-16', '89533397', 'gankhulug.gh@gmail.com', 0, 'mng', 'Улаанбаатар', '', '', 'Улаанбаатар, Баянзүрх'),
    (10124995, 110799641398, 'чм84062431', 'монгол', 'баяраа', 'өсөхбаяр', 1, '1984-06-24', '96615552', 'default@gerege.mn', 0, 'mng', 'Улаанбаатар', 'Баянгол', '17-р хороо', 'Улаанбаатар, Баянгол, 17-р хороо'),
    (10000021, 450107940837, 'та89081251', 'хэрээ', 'чулуунпүрэв', 'мөнхдалай', 1, '1989-08-12', '99080249', 'dalai0812@gmail.com', 0, 'mng', 'Улаанбаатар', 'Баянзүрх', '8-р хороо', 'Улаанбаатар, Баянзүрх, 8-р хороо'),
    (10000424, 110706612732, 'мз00312437', 'ангиртынбэлчир', 'цогтсайхан', 'чингүүнжав', 1, '2000-11-24', '80264610', 'chingunjav25@gmail.com', 0, 'mng', 'Улаанбаатар', 'Баянгол', '1-р хороо', 'Улаанбаатар, Баянгол, 1-р хороо')
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('citizens', 'id'), GREATEST((SELECT MAX(id) FROM citizens), 1));

-- ============================================
-- SEED: USERS
-- ============================================
INSERT INTO users (id, civil_id, reg_no, family_name, last_name, first_name, gender, birth_date, phone_no, email) VALUES
    (10000141, 623458282539, 'йк98013033', 'тайж', 'отгонбаяр', 'баярсайхан', 1, '1998-01-30', '94422470', 'craftzbay@gmail.com'),
    (10000263, 111949212017, 'ма74101813', 'харчин', 'цэнддорж', 'эрдэнэбат', 1, '1974-10-18', '99118306', 'erdenebatt@gmail.com'),
    (10000221, 111036772547, 'уп96030116', 'алагадуун', 'төвшин', 'мөнхбаяр', 1, '1996-03-01', '94942244', 'munkhbayar@gerege.com'),
    (10000081, 110770722023, 'уз98052158', 'өнгөө', 'соронзонболд', 'сэнгүм', 1, '1998-05-21', '88995566', 'sengum@gerege.mn'),
    (10162877, 823494982818, 'дк02301734', 'жалайр', 'алимаа', 'хүдэрчулуун', 1, '2002-10-17', '89908953', 'khuderchuluun@gerege.com'),
    (11949042, 111676962728, 'та00261713', 'эвэн', 'хонгор', 'чингүүн', 1, '2000-06-17', '99500509', 'sharshuwuu@gmail.com'),
    (10000101, 220102982359, 'лю98082031', 'алагууд', 'батжаргал', 'батзориг', 1, '1998-08-20', '88002608', 'batzorig@gerege.com'),
    (10124995, 110799641398, 'чм84062431', 'монгол', 'баяраа', 'өсөхбаяр', 1, '1984-06-24', '96615552', 'default@gerege.mn'),
    (10000021, 450107940837, 'та89081251', 'хэрээ', 'чулуунпүрэв', 'мөнхдалай', 1, '1989-08-12', '99080249', 'dalai0812@gmail.com'),
    (10000424, 110706612732, 'мз00312437', 'ангиртынбэлчир', 'цогтсайхан', 'чингүүнжав', 1, '2000-11-24', '80264610', 'chingunjav25@gmail.com'),
    (10000461, 0, 'бю90120312', 'жадик', 'манахмет', 'жанибек', 1, '1990-12-03', '99417882', 'm.janibek@gmail.com')
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('users', 'id'), GREATEST((SELECT MAX(id) FROM users), 1));

-- ============================================
-- SEED: USER_ROLES
-- ============================================
INSERT INTO user_roles (user_id, role_id) VALUES
    (10000141, 1),
    (10000263, 1),
    (10162877, 1),
    (10000424, 1),
    (10000141, 2),
    (10000263, 2),
    (10000081, 2),
    (10162877, 2),
    (11949042, 1),
    (11949042, 2),
    (10000081, 1),
    (10000081, 9),
    (10000263, 9),
    (10124995, 1),
    (10124995, 9),
    (10124995, 14),
    (10000461, 1)
ON CONFLICT DO NOTHING;

-- ============================================
-- SEED: ORGANIZATIONS
-- ============================================
INSERT INTO organizations (id, reg_no, name, type_id, phone_no, email, is_active, address_detail) VALUES
    (20028051, '6884857', 'Гэрэгэ коор', 4, '77773773', 'gereges.mn@gmail.com', true, 'Улаанбаатар хот, Сүхбаатар дүүрэг, 1-р хороо, Аюуд цамхаг, 12 давхар'),
    (20000002, '6235972', 'Гэрэгэ системс', 6, '77773773', 'gereges.mn@gmail.com', true, 'Улаанбаатар хот, Сүхбаатар дүүрэг, 1-р хороо, Аюуд цамхаг, 12 давхар'),
    (20001044, '6537359', 'Гэрэгэ пос', 6, '72773773', 'geregepos@gerege.mn', true, 'Аюуд товер 1205'),
    (20028052, '6537332', 'Гэрэгэ пэй', 6, '12345678', 'geregepay@gmail.com', true, ''),
    (20028053, '6980945', 'Төрийн цахим үйлчилгээний зохицуулалтын газар', 2, '12345678', 'info@company.mn', true, ''),
    (20028054, '5296722', 'Улсын бүртгэлийн ерөнхий газар', 2, '', '', true, ''),
    (20028055, '5731089', 'Гэрэгэ киоск', 6, '', '', true, ''),
    (20028056, '5213339', 'Номинтрейдинг', 1, '', '', true, ''),
    (20028057, '6155804', 'Премиум нэксус', 1, '', '', true, ''),
    (20028058, '8182159', 'Сэлэнгэ сайн үйлсийн сан', 3, '', '', true, ''),
    (20028059, '1083503', 'Тэнгэрийн суут монгол нийгэмлэг ТББ', 3, '', '', true, '')
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('organizations', 'id'), GREATEST((SELECT MAX(id) FROM organizations), 1));

-- ============================================
-- SEED: ORGANIZATION_USERS
-- ============================================
INSERT INTO organization_users (org_id, user_id) VALUES
    (20000002, 10000141),
    (20000002, 10000081),
    (20000002, 10000263),
    (20000002, 10162877),
    (20001044, 10000141),
    (20001044, 10000081),
    (20001044, 10000263),
    (20001044, 10162877),
    (20028051, 10000141),
    (20028051, 10000081),
    (20028051, 10000263),
    (20028051, 10162877),
    (20028051, 10000424),
    (20028055, 10000081),
    (20028056, 10000081),
    (20028051, 10124995),
    (20001044, 10124995),
    (20000002, 10000461)
ON CONFLICT DO NOTHING;

-- ============================================
-- SEED: ORG_TYPE_SYSTEMS
-- ============================================
INSERT INTO org_type_systems (type_id, system_id) VALUES
    (2, 8),
    (2, 10),
    (6, 10),
    (6, 9),
    (6, 8),
    (4, 7),
    (4, 10),
    (4, 9),
    (4, 8),
    (3, 10),
    (1, 9),
    (1, 10)
ON CONFLICT DO NOTHING;

-- ============================================
-- SEED: ORG_TYPE_ROLES
-- ============================================
INSERT INTO org_type_roles (type_id, role_id) VALUES
    (4, 1),
    (6, 2),
    (1, 2),
    (2, 2)
ON CONFLICT DO NOTHING;

-- ============================================
-- SEED: PERMISSIONS
-- ============================================
INSERT INTO permissions (id, code, name, description, module_id) VALUES
    (5, 'ADMIN-DASHBOARD-VIEW', 'Хянах самбар харах', 'Админ хянах самбар харах эрх', 16),
    (10, 'ADMIN-ORG-TYPE-CREATE', 'Байгууллагын төрөл үүсгэх', 'Байгууллагын төрөл үүсгэх эрх', 17),
    (11, 'ADMIN-ORG-TYPE-READ', 'Байгууллагын төрөл харах', 'Байгууллагын төрөл харах эрх', 17),
    (12, 'ADMIN-ORG-TYPE-UPDATE', 'Байгууллагын төрөл засах', 'Байгууллагын төрөл засах эрх', 17),
    (13, 'ADMIN-ORG-TYPE-DELETE', 'Байгууллагын төрөл устгах', 'Байгууллагын төрөл устгах эрх', 17),
    (14, 'ADMIN-ORG-CREATE', 'Байгууллага үүсгэх', 'Байгууллага үүсгэх эрх', 12),
    (15, 'ADMIN-ORG-READ', 'Байгууллага харах', 'Байгууллага харах эрх', 12),
    (16, 'ADMIN-ORG-UPDATE', 'Байгууллага засах', 'Байгууллага засах эрх', 12),
    (17, 'ADMIN-ORG-DELETE', 'Байгууллага устгах', 'Байгууллага устгах эрх', 12),
    (18, 'ADMIN-SYSTEM-CREATE', 'Систем үүсгэх', 'Систем үүсгэх эрх', 7),
    (19, 'ADMIN-SYSTEM-READ', 'Систем харах', 'Систем харах эрх', 7),
    (20, 'ADMIN-SYSTEM-UPDATE', 'Систем засах', 'Систем засах эрх', 7),
    (21, 'ADMIN-SYSTEM-DELETE', 'Систем устгах', 'Систем устгах эрх', 7),
    (22, 'ADMIN-MODULE-CREATE', 'Модул үүсгэх', 'Модул үүсгэх эрх', 8),
    (23, 'ADMIN-MODULE-READ', 'Модул харах', 'Модул харах эрх', 8),
    (24, 'ADMIN-MODULE-UPDATE', 'Модул засах', 'Модул засах эрх', 8),
    (25, 'ADMIN-MODULE-DELETE', 'Модул устгах', 'Модул устгах эрх', 8),
    (26, 'ADMIN-MODULE-PERMISSION-CREATE', 'Эрхийн удирдлага үүсгэх', 'Эрхийн удирдлага үүсгэх эрх', 11),
    (27, 'ADMIN-MODULE-PERMISSION-READ', 'Эрхийн удирдлага харах', 'Эрхийн удирдлага харах эрх', 11),
    (28, 'ADMIN-MODULE-PERMISSION-UPDATE', 'Эрхийн удирдлага засах', 'Эрхийн удирдлага засах эрх', 11),
    (29, 'ADMIN-MODULE-PERMISSION-DELETE', 'Эрхийн удирдлага устгах', 'Эрхийн удирдлага устгах эрх', 11),
    (30, 'ADMIN-ROLE-CREATE', 'Дүр үүсгэх', 'Дүр үүсгэх эрх', 10),
    (31, 'ADMIN-ROLE-READ', 'Дүр харах', 'Дүр харах эрх', 10),
    (32, 'ADMIN-ROLE-UPDATE', 'Дүр засах', 'Дүр засах эрх', 10),
    (33, 'ADMIN-ROLE-DELETE', 'Дүр устгах', 'Дүр устгах эрх', 10),
    (34, 'ADMIN-USER-CREATE', 'Хэрэглэгч үүсгэх', 'Хэрэглэгч үүсгэх эрх', 13),
    (35, 'ADMIN-USER-READ', 'Хэрэглэгч харах', 'Хэрэглэгч харах эрх', 13),
    (36, 'ADMIN-USER-UPDATE', 'Хэрэглэгч засах', 'Хэрэглэгч засах эрх', 13),
    (37, 'ADMIN-USER-DELETE', 'Хэрэглэгч устгах', 'Хэрэглэгч устгах эрх', 13),
    (38, 'GRGA-DASHBOARD-READ', 'Хянах самбар', 'Гэрэгэ апп хянах самбар', 19),
    (39, 'GRGA-ICON-CREATE', 'Апп айкон үүсгэх', 'Апп айкон үүсгэх эрх', 14),
    (40, 'GRGA-ICON-READ', 'Апп айкон харах', 'Апп айкон харах эрх', 14),
    (41, 'GRGA-ICON-UPDATE', 'Апп айкон засах', 'Апп айкон засах эрх', 14),
    (42, 'GRGA-ICON-DELETE', 'Апп айкон устгах', 'Апп айкон устгах эрх', 14),
    (43, 'GRGA-ROLE-CREATE', 'Гэрэгэ аппын дүр үүсгэх', 'Гэрэгэ аппын дүр үүсгэх эрх', 26),
    (44, 'GRGA-ROLE-READ', 'Гэрэгэ аппын дүр харах', 'Гэрэгэ аппын дүр харах эрх', 26),
    (45, 'GRGA-ROLE-UPDATE', 'Гэрэгэ аппын дүр засах', 'Гэрэгэ аппын дүр засах эрх', 26),
    (46, 'GRGA-ROLE-DELETE', 'Гэрэгэ аппын дүр устгах', 'Гэрэгэ аппын дүр устгах эрх', 26),
    (47, 'GRGB-DASHBOARD-READ', 'Гэрэгэ бизнес хянах самбар', 'Гэрэгэ бизнес хянах самбар', 20),
    (48, 'GRGB-ICON-CREATE', 'Гэрэгэ бизнес апп айкон үүсгэх', 'Гэрэгэ бизнес апп айкон үүсгэх эрх', 15),
    (49, 'GRGB-ICON-READ', 'Гэрэгэ бизнес апп айкон харах', 'Гэрэгэ бизнес апп айкон харах эрх', 15),
    (50, 'GRGB-ICON-UPDATE', 'Гэрэгэ бизнес апп айкон засах', 'Гэрэгэ бизнес апп айкон засах эрх', 15),
    (51, 'GRGB-ICON-DELETE', 'Гэрэгэ бизнес апп айкон устгах', 'Гэрэгэ бизнес апп айкон устгах эрх', 15),
    (52, 'GRGB-ROLE-CREATE', 'Гэрэгэ бизнесийн дүрүүд үүсгэх', 'Гэрэгэ бизнесийн дүрүүд үүсгэх эрх', 27),
    (53, 'GRGB-ROLE-READ', 'Гэрэгэ бизнесийн дүрүүд харах', 'Гэрэгэ бизнесийн дүрүүд харах эрх', 27),
    (54, 'GRGB-ROLE-UPDATE', 'Гэрэгэ бизнесийн дүрүүд засах', 'Гэрэгэ бизнесийн дүрүүд засах эрх', 27),
    (55, 'GRGB-ROLE-DELETE', 'Гэрэгэ бизнесийн дүрүүд устгах', 'Гэрэгэ бизнесийн дүрүүд устгах эрх', 27),
    (56, 'GRGB-AGENTS-CREATE', 'Бизнес агентууд үүсгэх', 'Бизнес агентууд үүсгэх эрх', 24),
    (57, 'GRGB-AGENTS-READ', 'Бизнес агентууд харах', 'Бизнес агентууд харах эрх', 24),
    (58, 'GRGB-AGENTS-UPDATE', 'Бизнес агентууд засах', 'Бизнес агентууд засах эрх', 24),
    (59, 'GRGB-AGENTS-DELETE', 'Бизнес агентууд устгах', 'Бизнес агентууд устгах эрх', 24),
    (60, 'TPAY-DASHBOARD-READ', 'T-Pay хянах самбар', 'T-Pay хянах самбар', 18),
    (61, 'TPAY-BALANCE-READ', 'TPay баланс', 'TPay баланс', 23),
    (62, 'TPAY-TRX-READ', 'TPay гүйлгээ', 'TPay гүйлгээ', 22),
    (63, 'TPAY-WALLET-READ', 'TPay хэтэвч', 'TPay хэтэвч', 21),
    (64, 'TPAY-ROLE-CREATE', 'ТPay системийн дүрүүд үүсгэх', 'ТPay системийн дүрүүд үүсгэх эрх', 25),
    (65, 'TPAY-ROLE-READ', 'ТPay системийн дүрүүд харах', 'ТPay системийн дүрүүд харах эрх', 25),
    (66, 'TPAY-ROLE-UPDATE', 'ТPay системийн дүрүүд засах', 'ТPay системийн дүрүүд засах эрх', 25),
    (67, 'TPAY-ROLE-DELETE', 'ТPay системийн дүрүүд устгах', 'ТPay системийн дүрүүд устгах эрх', 25),
    (68, 'ADMIN-AI-CHAT-CREATE', 'Хиймэл оюун удирдлага үүсгэх', 'Хиймэл оюун удирдлага үүсгэх эрх', 28),
    (69, 'ADMIN-AI-CHAT-READ', 'Хиймэл оюун удирдлага харах', 'Хиймэл оюун удирдлага харах эрх', 28),
    (70, 'ADMIN-AI-CHAT-UPDATE', 'Хиймэл оюун удирдлага засах', 'Хиймэл оюун удирдлага засах эрх', 28),
    (71, 'ADMIN-AI-CHAT-DELETE', 'Хиймэл оюун удирдлага устгах', 'Хиймэл оюун удирдлага устгах эрх', 28),
    (72, 'TPAY-USER-READ-CREATE', 'тпэй хэрэглэгч үүсгэх', 'үүсгэх эрх', 29),
    (73, 'TPAY-USER-READ-READ', 'тпэй хэрэглэгч харах', 'харах эрх', 29),
    (74, 'TPAY-USER-READ-UPDATE', 'тпэй хэрэглэгч засах', 'засах эрх', 29),
    (75, 'TPAY-USER-READ-DELETE', 'тпэй хэрэглэгч устгах', 'устгах эрх', 29),
    (76, 'TPAY-PRODUCTS-CREATE', 'tpay products үүсгэх', 'үүсгэх эрх', 30),
    (77, 'TPAY-PRODUCTS-READ', 'tpay products харах', 'харах эрх', 30),
    (78, 'TPAY-PRODUCTS-UPDATE', 'tpay products засах', 'засах эрх', 30),
    (79, 'TPAY-PRODUCTS-DELETE', 'tpay products устгах', 'устгах эрх', 30)
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('permissions', 'id'), GREATEST((SELECT MAX(id) FROM permissions), 1));

-- ============================================
-- SEED: ROLE_PERMISSIONS (Super Admin role 1 gets all permissions)
-- ============================================
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions WHERE deleted_date IS NULL
ON CONFLICT DO NOTHING;

-- Additional role permissions for other roles
INSERT INTO role_permissions (role_id, permission_id) VALUES
    (2, 60), (2, 61), (2, 62), (2, 63), (2, 64), (2, 65), (2, 66), (2, 67), (2, 72), (2, 73), (2, 74), (2, 75),
    (9, 60), (9, 64), (9, 65), (9, 66), (9, 67)
ON CONFLICT DO NOTHING;

-- ============================================
-- SEED: CHAT_ITEMS
-- ============================================
INSERT INTO chat_items (id, key, answer) VALUES
    (1, 'help', 'Та манай тусламжийн төвтэй холбогдож болно.'),
    (2, 'contact', 'Утас: 77001234, Email: support@example.com')
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('chat_items', 'id'), GREATEST((SELECT MAX(id) FROM chat_items), 1));

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO template_backend;

-- Clean up in reverse order of dependencies
DELETE FROM role_permissions;
DELETE FROM permissions;
DELETE FROM org_type_roles;
DELETE FROM org_type_systems;
DELETE FROM organization_users;
DELETE FROM organizations;
DELETE FROM user_roles;
DELETE FROM users;
DELETE FROM citizens;
DELETE FROM organization_types;
DELETE FROM roles;
DELETE FROM modules;
DELETE FROM actions;
DELETE FROM systems;
DELETE FROM chat_items;
-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

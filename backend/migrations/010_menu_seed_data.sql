-- +goose Up
-- +goose StatementBegin
SET search_path TO template_backend;

-- ============================================
-- SEED: MENUS
-- Navigation menus linked to modules and permissions
-- ============================================

-- ============================================
-- ADMIN SYSTEM ROOT MENUS (System ID: 7)
-- ============================================
INSERT INTO menus (id, code, key, name, description, icon, path, sequence, parent_id, permission_id, is_active) VALUES
    -- Root menu for Admin System
    (1, 'ADMIN', 'admin', 'Админ систем', 'Админ системийн үндсэн цэс', 'i-lucide-monitor-cog', '/admin', 1, NULL, NULL, true),

    -- Dashboard (Module 16)
    (2, 'ADMIN-DASHBOARD', 'admin-dashboard', 'Хянах самбар', 'Админ хянах самбар', 'i-lucide-layout-dashboard', '/admin/dashboard', 1, 1, 5, true),

    -- System Management (Module 7)
    (3, 'ADMIN-SYSTEM', 'admin-system', 'Систем', 'Систем удирдлага', 'i-lucide-server', '/admin/system', 2, 1, 19, true),

    -- Module Management (Module 8)
    (4, 'ADMIN-MODULE', 'admin-module', 'Модул', 'Модул удирдлага', 'i-lucide-puzzle', '/admin/module', 3, 1, 23, true),

    -- Role Management (Module 10)
    (5, 'ADMIN-ROLE', 'admin-role', 'Дүр', 'Дүр удирдлага', 'i-lucide-shield', '/admin/role', 4, 1, 31, true),

    -- Permission Management (Module 11)
    (6, 'ADMIN-MODULE-PERMISSION', 'admin-permission', 'Эрхийн удирдлага', 'Эрхийн удирдлага', 'i-lucide-key', '/admin/permission', 5, 1, 27, true),

    -- Organization Management (Module 12)
    (7, 'ADMIN-ORG', 'admin-org', 'Байгууллага', 'Байгууллага удирдлага', 'i-lucide-building-2', '/admin/organization', 6, 1, 15, true),

    -- Organization Type Management (Module 17)
    (8, 'ADMIN-ORGTYPE', 'admin-orgtype', 'Байгууллагын төрөл', 'Байгууллагын төрөл удирдлага', 'i-lucide-tags', '/admin/organization-type', 7, 1, 11, true),

    -- User Management (Module 13)
    (9, 'ADMIN-USER', 'admin-user', 'Хэрэглэгч', 'Хэрэглэгч удирдлага', 'i-lucide-users', '/admin/user', 8, 1, 35, true),

    -- AI Chat Management (Module 28)
    (10, 'ADMIN-AI-CHAT', 'admin-ai-chat', 'Хиймэл оюун', 'Хиймэл оюун удирдлага', 'i-lucide-bot', '/admin/ai-chat', 9, 1, 69, true)
ON CONFLICT (key) DO NOTHING;

-- ============================================
-- TPAY SYSTEM ROOT MENUS (System ID: 10)
-- ============================================
INSERT INTO menus (id, code, key, name, description, icon, path, sequence, parent_id, permission_id, is_active) VALUES
    -- Root menu for TPay System
    (20, 'TPAY', 'tpay', 'Гэрэгэ хэтэвч', 'Гэрэгэ хэтэвч системийн цэс', 'i-lucide-wallet', '/tpay', 2, NULL, NULL, true),

    -- Dashboard (Module 18)
    (21, 'TPAY-DASHBOARD', 'tpay-dashboard', 'Хянах самбар', 'T-Pay хянах самбар', 'i-lucide-layout-dashboard', '/tpay/dashboard', 1, 20, 60, true),

    -- Wallet (Module 21)
    (22, 'TPAY-WALLET', 'tpay-wallet', 'Хэтэвч', 'TPay хэтэвч', 'i-lucide-wallet', '/tpay/wallet', 2, 20, 63, true),

    -- Transactions (Module 22)
    (23, 'TPAY-TRX', 'tpay-trx', 'Гүйлгээ', 'TPay гүйлгээ', 'i-lucide-arrow-left-right', '/tpay/transactions', 3, 20, 62, true),

    -- Balance (Module 23)
    (24, 'TPAY-BALANCE', 'tpay-balance', 'Баланс', 'TPay баланс', 'i-lucide-banknote', '/tpay/balance', 4, 20, 61, true),

    -- Users (Module 29)
    (25, 'TPAY-USERS', 'tpay-users', 'Хэрэглэгчид', 'TPay хэрэглэгчид', 'i-lucide-users', '/tpay/users', 5, 20, 73, true),

    -- Products (Module 30)
    (26, 'TPAY-PRODUCTS', 'tpay-products', 'Бүтээгдэхүүн', 'Бүтээгдэхүүний каталог', 'i-lucide-package', '/tpay/products', 6, 20, 77, true),

    -- Role Management (Module 25)
    (27, 'TPAY-ROLE', 'tpay-role', 'Дүр', 'ТPay системийн дүрүүд', 'i-lucide-shield', '/tpay/role', 7, 20, 65, true)
ON CONFLICT (key) DO NOTHING;

-- ============================================
-- GRGA SYSTEM ROOT MENUS (Gerege App - System ID: 9)
-- ============================================
INSERT INTO menus (id, code, key, name, description, icon, path, sequence, parent_id, permission_id, is_active) VALUES
    -- Root menu for Gerege App System
    (40, 'GRGA', 'grga', 'Гэрэгэ апп', 'Гэрэгэ апп системийн цэс', 'i-lucide-smartphone', '/grga', 3, NULL, NULL, true),

    -- Dashboard (Module 19)
    (41, 'GRGA-DASHBOARD', 'grga-dashboard', 'Хянах самбар', 'Гэрэгэ апп хянах самбар', 'i-lucide-layout-dashboard', '/grga/dashboard', 1, 40, 38, true),

    -- App Icons (Module 14)
    (42, 'GRGA-ICON', 'grga-icon', 'Апп айкон', 'Апп айкон удирдлага', 'i-lucide-image', '/grga/icon', 2, 40, 40, true),

    -- Role Management (Module 26)
    (43, 'GRGA-ROLE', 'grga-role', 'Дүр', 'Гэрэгэ аппын дүрүүд', 'i-lucide-shield', '/grga/role', 3, 40, 44, true)
ON CONFLICT (key) DO NOTHING;

-- ============================================
-- GRGB SYSTEM ROOT MENUS (Gerege Business - System ID: 8)
-- ============================================
INSERT INTO menus (id, code, key, name, description, icon, path, sequence, parent_id, permission_id, is_active) VALUES
    -- Root menu for Gerege Business System
    (60, 'GRGB', 'grgb', 'Гэрэгэ бизнес', 'Гэрэгэ бизнес системийн цэс', 'i-lucide-briefcase-business', '/grgb', 4, NULL, NULL, true),

    -- Dashboard (Module 20)
    (61, 'GRGB-DASHBOARD', 'grgb-dashboard', 'Хянах самбар', 'Гэрэгэ бизнес хянах самбар', 'i-lucide-layout-dashboard', '/grgb/dashboard', 1, 60, 47, true),

    -- App Icons (Module 15)
    (62, 'GRGB-ICON', 'grgb-icon', 'Апп айкон', 'Гэрэгэ бизнес апп айкон', 'i-lucide-image', '/grgb/icon', 2, 60, 49, true),

    -- Agents (Module 24)
    (63, 'GRGB-AGENTS', 'grgb-agents', 'Агент', 'Бизнес агентууд', 'i-lucide-users', '/grgb/agents', 3, 60, 57, true),

    -- Role Management (Module 27)
    (64, 'GRGB-ROLE', 'grgb-role', 'Дүр', 'Гэрэгэ бизнесийн дүрүүд', 'i-lucide-shield', '/grgb/role', 4, 60, 53, true)
ON CONFLICT (key) DO NOTHING;

-- Update sequences
SELECT setval(pg_get_serial_sequence('menus', 'id'), GREATEST((SELECT MAX(id) FROM menus), 1));

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO template_backend;

-- Clean up menus (in order to avoid FK constraint issues)
DELETE FROM menus WHERE parent_id IS NOT NULL;
DELETE FROM menus WHERE parent_id IS NULL;

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Up
-- +goose StatementBegin
SET search_path TO template_backend;

-- ============================================
-- SEED: APP_SERVICE_ICON_GROUPS
-- ============================================
INSERT INTO app_service_icon_groups (id, name, name_en, icon, type_name, seq) VALUES
    (1, 'Төрийн үйлчилгээ', 'Government Services', 'i-lucide-landmark', 'group', 1),
    (2, 'Санхүүгийн үйлчилгээ', 'Financial Services', 'i-lucide-wallet', 'group', 2),
    (3, 'Тээвэр', 'Transportation', 'i-lucide-car', 'group', 3),
    (4, 'Боловсрол', 'Education', 'i-lucide-graduation-cap', 'group', 4),
    (5, 'Эрүүл мэнд', 'Healthcare', 'i-lucide-heart-pulse', 'group', 5),
    (6, 'Худалдаа', 'Shopping', 'i-lucide-shopping-cart', 'group', 6),
    (7, 'Аялал жуулчлал', 'Travel & Tourism', 'i-lucide-plane', 'group', 7),
    (8, 'Хөдөө аж ахуй', 'Agriculture', 'i-lucide-wheat', 'group', 8)
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('app_service_icon_groups', 'id'), GREATEST((SELECT MAX(id) FROM app_service_icon_groups), 1));

-- ============================================
-- SEED: APP_SERVICE_ICONS
-- ============================================
INSERT INTO app_service_icons (id, name, name_en, icon, link, group_id, seq, is_native, is_public, is_featured, description, system_code) VALUES
    -- Төрийн үйлчилгээ
    (1, 'Иргэний бүртгэл', 'Civil Registration', 'i-lucide-user-check', '/civil-registration', 1, 1, false, true, true, 'Иргэний бүртгэлийн үйлчилгээ', 'GS'),
    (2, 'Газрын бүртгэл', 'Land Registration', 'i-lucide-map-pin', '/land-registration', 1, 2, false, true, false, 'Газрын бүртгэлийн үйлчилгээ', 'GS'),
    (3, 'Татвар', 'Tax Services', 'i-lucide-receipt', '/tax', 1, 3, false, true, false, 'Татварын үйлчилгээ', 'GS'),
    (4, 'Нийгмийн даатгал', 'Social Insurance', 'i-lucide-shield-check', '/social-insurance', 1, 4, false, true, true, 'Нийгмийн даатгалын үйлчилгээ', 'GS'),

    -- Санхүүгийн үйлчилгээ
    (5, 'Гэрэгэ хэтэвч', 'Gerege Wallet', 'i-lucide-wallet', '/wallet', 2, 1, true, true, true, 'Цахим хэтэвч', 'TP'),
    (6, 'Төлбөр төлөх', 'Payments', 'i-lucide-credit-card', '/payments', 2, 2, false, true, true, 'Төлбөр тооцоо', 'TP'),
    (7, 'Гүйлгээний түүх', 'Transaction History', 'i-lucide-history', '/transactions', 2, 3, false, true, false, 'Гүйлгээний түүх харах', 'TP'),
    (8, 'QR төлбөр', 'QR Payment', 'i-lucide-qr-code', '/qr-payment', 2, 4, true, true, true, 'QR кодоор төлбөр төлөх', 'TP'),

    -- Тээвэр
    (9, 'Авто машины бүртгэл', 'Vehicle Registration', 'i-lucide-car', '/vehicle-registration', 3, 1, false, true, false, 'Тээврийн хэрэгслийн бүртгэл', 'GS'),
    (10, 'Жолооны үнэмлэх', 'Driver License', 'i-lucide-id-card', '/driver-license', 3, 2, false, true, true, 'Жолооны үнэмлэхний үйлчилгээ', 'GS'),
    (11, 'Техникийн хяналт', 'Technical Inspection', 'i-lucide-clipboard-check', '/technical-inspection', 3, 3, false, true, false, 'Техникийн хяналтын үйлчилгээ', 'GS'),

    -- Боловсрол
    (12, 'Цахим сургалт', 'E-Learning', 'i-lucide-book-open', '/e-learning', 4, 1, false, true, true, 'Онлайн сургалт', 'GA'),
    (13, 'Гэрчилгээ баталгаажуулах', 'Certificate Verification', 'i-lucide-award', '/certificate-verify', 4, 2, false, true, false, 'Гэрчилгээ, дипломын баталгаажуулалт', 'GS'),

    -- Эрүүл мэнд
    (14, 'Эрүүл мэндийн даатгал', 'Health Insurance', 'i-lucide-heart', '/health-insurance', 5, 1, false, true, true, 'Эрүүл мэндийн даатгалын үйлчилгээ', 'GS'),
    (15, 'Цаг захиалга', 'Appointment', 'i-lucide-calendar', '/appointment', 5, 2, false, true, false, 'Эмнэлэгт цаг захиалах', 'GA'),

    -- Худалдаа
    (16, 'Онлайн дэлгүүр', 'Online Shop', 'i-lucide-store', '/shop', 6, 1, false, true, true, 'Онлайн худалдаа', 'GB'),
    (17, 'Хямдрал урамшуулал', 'Promotions', 'i-lucide-tag', '/promotions', 6, 2, false, true, true, 'Хямдрал урамшууллын мэдээлэл', 'GB'),

    -- Аялал жуулчлал
    (18, 'Зочид буудал', 'Hotels', 'i-lucide-bed', '/hotels', 7, 1, false, true, false, 'Зочид буудлын захиалга', 'GA'),
    (19, 'Нислэгийн тасалбар', 'Flight Tickets', 'i-lucide-plane-takeoff', '/flights', 7, 2, false, true, true, 'Онгоцны тасалбар захиалах', 'GA'),

    -- Хөдөө аж ахуй
    (20, 'Мал бүртгэл', 'Livestock Registry', 'i-lucide-pawprint', '/livestock', 8, 1, false, true, false, 'Малын бүртгэлийн систем', 'GS')
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('app_service_icons', 'id'), GREATEST((SELECT MAX(id) FROM app_service_icons), 1));

-- ============================================
-- SEED: NEWS
-- ============================================
INSERT INTO news (id, title, text, image_url, created_user_id, created_org_id) VALUES
    (1, 'Гэрэгэ систем шинэчлэгдлээ', 'Гэрэгэ системийн шинэ хувилбар амжилттай нэвтэрлээ. Энэхүү шинэчлэлтэнд хэрэглэгчийн интерфэйс сайжруулалт, гүйцэтгэлийн оновчлол, аюулгүй байдлын шинэ функцууд багтсан болно.', 'https://example.com/images/news/system-update.jpg', 10000141, 20028051),
    (2, 'Шинэ үйлчилгээ нэмэгдлээ', 'Иргэдэд зориулсан шинэ цахим үйлчилгээ нэмэгдлээ. Та одоо гар утаснаасаа бүх төрлийн бүртгэл хийх боломжтой боллоо.', 'https://example.com/images/news/new-service.jpg', 10000141, 20028051),
    (3, 'Хөнгөлөлттэй үнийн санал', 'Шинэ хэрэглэгчдэд зориулсан онцгой хөнгөлөлт! Бүртгүүлсэн эхний сард бүх гүйлгээний шимтгэл 50% хөнгөлөлттэй.', 'https://example.com/images/news/discount.jpg', 10000263, 20028051),
    (4, 'Системийн засвар үйлчилгээ', 'Системийн хэвийн ажиллагааг хангах үүднээс 2025 оны 3 дугаар сарын 15-ны өдрийн 00:00-06:00 цагийн хооронд техникийн засвар үйлчилгээ хийгдэнэ.', 'https://example.com/images/news/maintenance.jpg', 10000263, 20028051),
    (5, 'Хамтын ажиллагааны гэрээ байгуулав', 'Гэрэгэ системс нь Улсын бүртгэлийн ерөнхий газартай хамтран ажиллах гэрээ байгууллаа. Энэхүү хамтын ажиллагааны хүрээнд иргэдэд үзүүлэх цахим үйлчилгээ улам өргөжинө.', 'https://example.com/images/news/partnership.jpg', 10000141, 20000002),
    (6, 'Аюулгүй байдлын зөвлөмж', 'Хэрэглэгчдийн аюулгүй байдлыг хангах үүднээс нууц үгээ тогтмол солих, хоёр шатлалт нэвтрэлт идэвхжүүлэхийг зөвлөж байна.', 'https://example.com/images/news/security.jpg', 10162877, 20028051),
    (7, 'Мобайл апп шинэчлэгдлээ', 'iOS болон Android платформ дээрх Гэрэгэ аппликейшн шинэ хувилбар гарлаа. App Store болон Google Play-ээс татаж авна уу.', 'https://example.com/images/news/mobile-app.jpg', 10000081, 20000002),
    (8, 'Хэрэглэгчийн сэтгэл ханамжийн судалгаа', 'Таны санал бодол бидэнд чухал! Хэрэглэгчийн сэтгэл ханамжийн судалгаанд оролцож, үйлчилгээгээ сайжруулахад хувь нэмрээ оруулна уу.', 'https://example.com/images/news/survey.jpg', 10000141, 20028051),
    (9, 'QR төлбөрийн систем нэвтэрлээ', 'Гэрэгэ хэтэвчинд QR кодоор төлбөр төлөх боломж нэмэгдлээ. Хурдан, хялбар, аюулгүй төлбөр тооцоо хийгээрэй.', 'https://example.com/images/news/qr-payment.jpg', 10000263, 20028052),
    (10, 'Байгууллагын үйлчилгээ нээгдлээ', 'Бизнес эрхлэгчдэд зориулсан Гэрэгэ Бизнес платформ нээгдлээ. Байгууллагынхаа санхүүг удирдах, ажилчдын цалин олгох зэрэг олон үйлчилгээг нэг дороос авах боломжтой.', 'https://example.com/images/news/business.jpg', 10000141, 20000002)
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('news', 'id'), GREATEST((SELECT MAX(id) FROM news), 1));

-- ============================================
-- SEED: NOTIFICATION_GROUPS
-- ============================================
INSERT INTO notification_groups (id, user_id, title, content, type, tenant, created_username, created_user_id, created_org_id) VALUES
    (1, 0, 'Системийн мэдэгдэл', 'Системийн чухал мэдэгдлүүд', 'system', 'gerege', 'System', 0, 20028051),
    (2, 0, 'Төлбөрийн мэдэгдэл', 'Төлбөр тооцооны мэдэгдлүүд', 'payment', 'tpay', 'System', 0, 20028052),
    (3, 0, 'Урамшууллын мэдэгдэл', 'Хямдрал урамшууллын мэдэгдлүүд', 'promo', 'gerege', 'System', 0, 20028051),
    (4, 0, 'Аюулгүй байдлын мэдэгдэл', 'Аюулгүй байдалтай холбоотой мэдэгдлүүд', 'security', 'gerege', 'System', 0, 20028051)
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('notification_groups', 'id'), GREATEST((SELECT MAX(id) FROM notification_groups), 1));

-- ============================================
-- SEED: NOTIFICATIONS
-- ============================================
INSERT INTO notifications (id, user_id, title, content, is_read, type, tenant, group_id, created_username, created_user_id, created_org_id) VALUES
    -- System notifications
    (1, 10000141, 'Тавтай морил!', 'Гэрэгэ системд бүртгүүлсэнд баярлалаа. Таны бүртгэл амжилттай үүслээ.', true, 'system', 'gerege', 1, 'System', 0, 20028051),
    (2, 10000263, 'Тавтай морил!', 'Гэрэгэ системд бүртгүүлсэнд баярлалаа. Таны бүртгэл амжилттай үүслээ.', true, 'system', 'gerege', 1, 'System', 0, 20028051),
    (3, 10000081, 'Тавтай морил!', 'Гэрэгэ системд бүртгүүлсэнд баярлалаа. Таны бүртгэл амжилттай үүслээ.', false, 'system', 'gerege', 1, 'System', 0, 20028051),

    -- Payment notifications
    (4, 10000141, 'Төлбөр амжилттай', 'Таны 50,000₮ төлбөр амжилттай хийгдлээ. Гүйлгээний дугаар: TXN001234', true, 'payment', 'tpay', 2, 'TPay System', 0, 20028052),
    (5, 10000141, 'Мөнгө хүлээн авлаа', 'Танд 100,000₮ шилжүүлэг ирлээ. Илгээгч: Эрдэнэбат', true, 'payment', 'tpay', 2, 'TPay System', 0, 20028052),
    (6, 10000263, 'Төлбөр амжилттай', 'Таны 25,000₮ төлбөр амжилттай хийгдлээ. Гүйлгээний дугаар: TXN001235', false, 'payment', 'tpay', 2, 'TPay System', 0, 20028052),

    -- Promo notifications
    (7, 10000141, 'Шинэ хямдрал!', '50% хөнгөлөлттэй үйлчилгээ эхэллээ. 2025 оны 3-р сарын 31 хүртэл хүчинтэй.', false, 'promo', 'gerege', 3, 'Marketing', 0, 20028051),
    (8, 10000263, 'Шинэ хямдрал!', '50% хөнгөлөлттэй үйлчилгээ эхэллээ. 2025 оны 3-р сарын 31 хүртэл хүчинтэй.', false, 'promo', 'gerege', 3, 'Marketing', 0, 20028051),
    (9, 10000081, 'Урамшуулал', 'Найзаа урих тутам 5,000₮ бонус авах боломжтой!', false, 'promo', 'gerege', 3, 'Marketing', 0, 20028051),

    -- Security notifications
    (10, 10000141, 'Шинэ төхөөрөмжөөс нэвтэрлээ', 'Таны бүртгэлд iPhone 15 Pro төхөөрөмжөөс нэвтэрсэн байна. Хэрэв энэ та биш бол нууц үгээ шинэчилнэ үү.', true, 'security', 'gerege', 4, 'Security', 0, 20028051),
    (11, 10000263, 'Нууц үг солигдлоо', 'Таны нууц үг амжилттай солигдлоо.', true, 'security', 'gerege', 4, 'Security', 0, 20028051),
    (12, 10162877, 'Хоёр шатлалт нэвтрэлт идэвхжлээ', 'Таны бүртгэлд хоёр шатлалт нэвтрэлт амжилттай идэвхжлээ.', false, 'security', 'gerege', 4, 'Security', 0, 20028051)
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('notifications', 'id'), GREATEST((SELECT MAX(id) FROM notifications), 1));

-- ============================================
-- SEED: VEHICLES
-- ============================================
INSERT INTO vehicles (id, plate_no) VALUES
    (1, '0001УБА'),
    (2, '0002УБА'),
    (3, '1234УНА'),
    (4, '5678УНБ'),
    (5, '9999УБЕ'),
    (6, '1111ДОР'),
    (7, '2222ХОВ'),
    (8, '3333ӨМН'),
    (9, '4444БАЯ'),
    (10, '5555СЭЛ')
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('vehicles', 'id'), GREATEST((SELECT MAX(id) FROM vehicles), 1));

-- ============================================
-- SEED: PUBLIC_FILES
-- ============================================
INSERT INTO public_files (id, name, extension, description, file_url, created_user_id, created_org_id) VALUES
    (1, 'logo', 'png', 'Гэрэгэ системийн лого', 'https://cdn.gerege.mn/images/logo.png', 10000141, 20028051),
    (2, 'banner_home', 'jpg', 'Нүүр хуудасны баннер', 'https://cdn.gerege.mn/images/banner-home.jpg', 10000141, 20028051),
    (3, 'banner_promo', 'jpg', 'Урамшууллын баннер', 'https://cdn.gerege.mn/images/banner-promo.jpg', 10000263, 20028051),
    (4, 'icon_wallet', 'svg', 'Хэтэвчийн айкон', 'https://cdn.gerege.mn/icons/wallet.svg', 10000081, 20028052),
    (5, 'icon_payment', 'svg', 'Төлбөрийн айкон', 'https://cdn.gerege.mn/icons/payment.svg', 10000081, 20028052),
    (6, 'guide_pdf', 'pdf', 'Хэрэглэгчийн гарын авлага', 'https://cdn.gerege.mn/docs/user-guide.pdf', 10000141, 20028051),
    (7, 'terms_pdf', 'pdf', 'Үйлчилгээний нөхцөл', 'https://cdn.gerege.mn/docs/terms.pdf', 10000141, 20028051),
    (8, 'privacy_pdf', 'pdf', 'Нууцлалын бодлого', 'https://cdn.gerege.mn/docs/privacy.pdf', 10000141, 20028051),
    (9, 'app_screenshot_1', 'png', 'Апп дэлгэцийн зураг 1', 'https://cdn.gerege.mn/images/screenshot-1.png', 10000263, 20000002),
    (10, 'app_screenshot_2', 'png', 'Апп дэлгэцийн зураг 2', 'https://cdn.gerege.mn/images/screenshot-2.png', 10000263, 20000002)
ON CONFLICT DO NOTHING;

SELECT setval(pg_get_serial_sequence('public_files', 'id'), GREATEST((SELECT MAX(id) FROM public_files), 1));

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SET search_path TO template_backend;

-- Clean up in reverse order
DELETE FROM public_files WHERE id <= 10;
DELETE FROM vehicles WHERE id <= 10;
DELETE FROM notifications WHERE id <= 12;
DELETE FROM notification_groups WHERE id <= 4;
DELETE FROM news WHERE id <= 10;
DELETE FROM app_service_icons WHERE id <= 20;
DELETE FROM app_service_icon_groups WHERE id <= 8;

-- Reset search_path for goose
RESET search_path;

-- +goose StatementEnd

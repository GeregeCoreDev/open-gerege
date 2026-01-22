//go:build integration

// Package integration contains integration tests
//
// File: setup_test.go
// Description: Integration test setup and helpers with Testcontainers
package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"templatev25/internal/domain"
	"templatev25/internal/repository"
	"templatev25/tests/testutils"

	"git.gerege.mn/backend-packages/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var testDB *gorm.DB

// TestMain sets up and tears down the test database
func TestMain(m *testing.M) {
	// Setup
	ctx := context.Background()
	container, dsn, err := setupPostgresContainer(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to setup postgres container: %v\n", err)
		os.Exit(1)
	}

	// Connect to database
	testDB, err = gorm.Open(postgresDriver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to test database: %v\n", err)
		container.Terminate(ctx)
		os.Exit(1)
	}

	// Run migrations
	if err := runMigrations(testDB); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run migrations: %v\n", err)
		container.Terminate(ctx)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Teardown
	teardown(ctx, container)

	os.Exit(code)
}

func setupPostgresContainer(ctx context.Context) (testcontainers.Container, string, error) {
	dbName := "test_db"
	dbUser := "test"
	dbPassword := "test"

	pgContainer, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, "", err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		pgContainer.Terminate(ctx)
		return nil, "", err
	}

	return pgContainer, connStr, nil
}

func teardown(ctx context.Context, container testcontainers.Container) {
	if testDB != nil {
		if sqlDB, err := testDB.DB(); err == nil {
			sqlDB.Close()
		}
	}
	if container != nil {
		container.Terminate(ctx)
	}
}

// runMigrations creates test tables
func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Organization{},
		&domain.System{},
		&domain.Module{},
		&domain.Role{},
		&domain.Permission{},
		&domain.UserRole{},
		&domain.Menu{},
		&domain.News{},
		&domain.Notification{},
		&domain.NotificationGroup{},
		&domain.ChatItem{},
	)
}

// GetTestDB returns the test database connection
func GetTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	if testDB == nil {
		t.Fatal("test database not initialized")
	}
	return testDB
}

// GetTestDBWithTx returns a database wrapped in a transaction for isolation
func GetTestDBWithTx(t *testing.T) *gorm.DB {
	t.Helper()
	db := GetTestDB(t)

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}

	t.Cleanup(func() {
		tx.Rollback()
	})

	return tx
}

// Repositories holds all repository instances for testing
type Repositories struct {
	User         repository.UserRepository
	Role         repository.RoleRepository
	Permission   repository.PermissionRepository
	Organization repository.OrganizationRepository
	System       repository.SystemRepository
	Module       repository.ModuleRepository
	Menu         repository.MenuRepository
	News         repository.NewsRepository
	Notification repository.NotificationRepository
	ChatItem     repository.ChatItemRepository
}

// NewTestRepositories creates all repository instances with the test database
func NewTestRepositories(t *testing.T, db *gorm.DB) *Repositories {
	t.Helper()
	return &Repositories{
		User:         repository.NewUserRepository(db),
		Role:         repository.NewRoleRepository(db),
		Permission:   repository.NewPermissionRepository(db),
		Organization: repository.NewOrganizationRepository(db),
		System:       repository.NewSystemRepository(db),
		Module:       repository.NewModuleRepository(db, &config.Config{}),
		Menu:         repository.NewMenuRepository(db, &config.Config{}),
		News:         repository.NewNewsRepository(db),
		Notification: repository.NewNotificationRepository(db),
		ChatItem:     repository.NewChatItemRepository(db),
	}
}

// NewTestContext returns a context with test user/org info
func NewTestContext() context.Context {
	return context.Background()
}

// SeedTestUser creates a test user and returns it
func SeedTestUser(t *testing.T, db *gorm.DB) domain.User {
	t.Helper()

	user := domain.User{
		RegNo:     "AA12345678",
		FirstName: "Test",
		LastName:  "User",
		Email:     "test@example.com",
		PhoneNo:   "99112233",
		Gender:    1,
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to seed test user: %v", err)
	}

	return user
}

// SeedTestUsers creates multiple test users
func SeedTestUsers(t *testing.T, db *gorm.DB, count int) []domain.User {
	t.Helper()

	users := make([]domain.User, count)
	for i := 0; i < count; i++ {
		users[i] = domain.User{
			RegNo:     "AA" + string(rune('0'+i%10)) + "1234567",
			FirstName: "Test" + string(rune('A'+i%26)),
			LastName:  "User",
			Email:     "test" + string(rune('0'+i%10)) + "@example.com",
			PhoneNo:   "9911223" + string(rune('0'+i%10)),
			Gender:    i%2 + 1,
		}
	}

	if err := db.Create(&users).Error; err != nil {
		t.Fatalf("failed to seed test users: %v", err)
	}

	return users
}

// SeedTestSystem creates a test system and returns it
func SeedTestSystem(t *testing.T, db *gorm.DB) domain.System {
	t.Helper()

	isActive := true
	system := domain.System{
		Name:        "Test System",
		Description: "A test system",
		IsActive:    &isActive,
		Sequence:    1,
	}

	if err := db.Create(&system).Error; err != nil {
		t.Fatalf("failed to seed test system: %v", err)
	}

	return system
}

// SeedTestRole creates a test role and returns it
func SeedTestRole(t *testing.T, db *gorm.DB, systemID int) domain.Role {
	t.Helper()

	isActive := true
	role := domain.Role{
		SystemID:    systemID,
		Name:        "Test Role",
		Description: "A test role",
		IsActive:    &isActive,
	}

	if err := db.Create(&role).Error; err != nil {
		t.Fatalf("failed to seed test role: %v", err)
	}

	return role
}

// SeedTestOrganization creates a test organization and returns it
func SeedTestOrganization(t *testing.T, db *gorm.DB) domain.Organization {
	t.Helper()

	org := domain.Organization{
		Name:     "Test Organization",
		IsActive: boolPtr(true),
	}

	if err := db.Create(&org).Error; err != nil {
		t.Fatalf("failed to seed test organization: %v", err)
	}

	return org
}

// CleanupTable truncates a table
func CleanupTable(t *testing.T, db *gorm.DB, tableName string) {
	t.Helper()
	testutils.TruncateTable(t, db, tableName)
}

// CleanupTables truncates multiple tables
func CleanupTables(t *testing.T, db *gorm.DB, tableNames ...string) {
	t.Helper()
	testutils.TruncateTables(t, db, tableNames...)
}

// CreateTestContext returns a test context
func CreateTestContext() context.Context {
	return context.Background()
}

// SeedTestMenu creates a test menu and returns it
func SeedTestMenu(t *testing.T, db *gorm.DB) domain.Menu {
	t.Helper()
	isActive := true
	menu := domain.Menu{
		Name:     "Test Menu",
		Path:     "/test",
		Sequence: 1,
		IsActive: &isActive,
	}
	if err := db.Create(&menu).Error; err != nil {
		t.Fatalf("failed to seed test menu: %v", err)
	}
	return menu
}

// SeedTestMenus creates multiple test menus
func SeedTestMenus(t *testing.T, db *gorm.DB, count int) []domain.Menu {
	t.Helper()
	isActive := true
	menus := make([]domain.Menu, count)
	for i := 0; i < count; i++ {
		menus[i] = domain.Menu{
			Name:     fmt.Sprintf("Menu %d", i),
			Path:     fmt.Sprintf("/menu/%d", i),
			Sequence: int64(i + 1),
			IsActive: &isActive,
		}
	}
	if err := db.Create(&menus).Error; err != nil {
		t.Fatalf("failed to seed test menus: %v", err)
	}
	return menus
}

// SeedTestNews creates a test news item and returns it
func SeedTestNews(t *testing.T, db *gorm.DB) domain.News {
	t.Helper()
	news := domain.News{
		Title:    "Test News",
		Text:     "This is test news content",
		ImageUrl: "https://example.com/image.jpg",
	}
	if err := db.Create(&news).Error; err != nil {
		t.Fatalf("failed to seed test news: %v", err)
	}
	return news
}

// SeedTestNewsItems creates multiple test news items
func SeedTestNewsItems(t *testing.T, db *gorm.DB, count int) []domain.News {
	t.Helper()
	newsItems := make([]domain.News, count)
	for i := 0; i < count; i++ {
		newsItems[i] = domain.News{
			Title:    fmt.Sprintf("News %d", i),
			Text:     fmt.Sprintf("News content %d", i),
			ImageUrl: fmt.Sprintf("https://example.com/image%d.jpg", i),
		}
	}
	if err := db.Create(&newsItems).Error; err != nil {
		t.Fatalf("failed to seed test news items: %v", err)
	}
	return newsItems
}

// SeedTestNotificationGroup creates a test notification group
func SeedTestNotificationGroup(t *testing.T, db *gorm.DB, userID int) domain.NotificationGroup {
	t.Helper()
	group := domain.NotificationGroup{
		UserId:  userID,
		Title:   "Test Group",
		Content: "Test group content",
		Type:    "info",
		Tenant:  "test",
	}
	if err := db.Create(&group).Error; err != nil {
		t.Fatalf("failed to seed test notification group: %v", err)
	}
	return group
}

// SeedTestNotification creates a test notification
func SeedTestNotification(t *testing.T, db *gorm.DB, userID int, groupID int) domain.Notification {
	t.Helper()
	notification := domain.Notification{
		UserId:  userID,
		Title:   "Test Notification",
		Content: "Test notification content",
		IsRead:  false,
		Type:    "info",
		Tenant:  "test",
		GroupId: groupID,
	}
	if err := db.Create(&notification).Error; err != nil {
		t.Fatalf("failed to seed test notification: %v", err)
	}
	return notification
}

// SeedTestNotifications creates multiple test notifications
func SeedTestNotifications(t *testing.T, db *gorm.DB, userID int, groupID int, count int) []domain.Notification {
	t.Helper()
	notifications := make([]domain.Notification, count)
	for i := 0; i < count; i++ {
		notifications[i] = domain.Notification{
			UserId:  userID,
			Title:   fmt.Sprintf("Notification %d", i),
			Content: fmt.Sprintf("Notification content %d", i),
			IsRead:  false,
			Type:    "info",
			Tenant:  "test",
			GroupId: groupID,
		}
	}
	if err := db.Create(&notifications).Error; err != nil {
		t.Fatalf("failed to seed test notifications: %v", err)
	}
	return notifications
}

// SeedTestChatItem creates a test chat item
func SeedTestChatItem(t *testing.T, db *gorm.DB) domain.ChatItem {
	t.Helper()
	item := domain.ChatItem{
		Answer: "test answer",
	}
	if err := db.Create(&item).Error; err != nil {
		t.Fatalf("failed to seed test chat item: %v", err)
	}
	return item
}

// SeedTestChatItems creates multiple test chat items
func SeedTestChatItems(t *testing.T, db *gorm.DB, count int) []domain.ChatItem {
	t.Helper()
	items := make([]domain.ChatItem, count)
	for i := 0; i < count; i++ {
		items[i] = domain.ChatItem{
			Answer: fmt.Sprintf("answer %d", i),
		}
	}
	if err := db.Create(&items).Error; err != nil {
		t.Fatalf("failed to seed test chat items: %v", err)
	}
	return items
}

package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/OfficialEvsty/aa-data/commands"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	integration "github.com/OfficialEvsty/aa-data/integration_test"
	"github.com/OfficialEvsty/aa-data/queries"
	"github.com/OfficialEvsty/aa-data/repos"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/magiconair/properties/assert"
	"github.com/testcontainers/testcontainers-go"
	"log"
	"os"
	"testing"
	"time"
)

var testDB *sql.DB
var pgContainer testcontainers.Container

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../migrations", // путь к миграциям
		"postgres", driver)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up failed: %w", err)
	}

	return nil
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	container, dsn, err := integration.PostgresSetup(ctx)
	if err != nil {
		log.Fatal(err)
	}
	pgContainer = container
	defer pgContainer.Terminate(ctx)

	testDB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err := testDB.PingContext(ctx); err != nil {
		log.Fatalf("cannot ping test database: %v", err)
	}
	err = runMigrations(testDB)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("migrations success")
	code := m.Run()
	os.Exit(code)
}
func TestInsertPublish(t *testing.T) {
	ctx := context.Background()
	log.Println("starts insert publishing transaction...")
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()

	repo := repos.NewPublishRepository(testDB)
	id := uuid.New()

	key := "somekey.png"
	bucket := "somecoolbucket"
	s3name := "selectel"
	publish := domain.PublishedScreenshot{
		ID: id,
		S3Data: serializable.S3Screenshot{
			Key:    key,
			Bucket: bucket,
			S3Name: s3name,
		},
	}
	err = repo.WithTx(tx).Add(ctx, publish)
	if err != nil {
		t.Fatal(err)
	}
	pub, err := repo.WithTx(tx).Get(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	isPubDataCorrect := pub.S3Data.Bucket == bucket &&
		pub.S3Data.S3Name == s3name &&
		pub.S3Data.Key == key
	if !isPubDataCorrect {
		t.Fatal("expected published screen shot to be present")
	}
}

func TestRefreshTokensCRUD(t *testing.T) {
	ctx := context.Background()
	log.Println("starts refresh tokens CRUD...")
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	repo := repos.NewRefreshTokenRepository(testDB)
	token := domain.RefreshToken{
		Token:     uuid.New().String(),
		UserID:    uuid.New(),
		ExpiresAt: time.Now().Add(time.Hour * 48),
	}
	log.Println("token repository created ...")
	newToken, err := repo.WithTx(tx).AddOrUpdate(ctx, token)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, token.Token, newToken.Token)
	assert.Equal(t, token.ExpiresAt.Unix(), newToken.ExpiresAt.Unix())
	assert.Equal(t, token.UserID, newToken.UserID)
	log.Println("refresh token successfully added ...")
	receivedTokenByID, err := repo.WithTx(tx).GetByToken(ctx, token.Token)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, token.Token, receivedTokenByID.Token)
	assert.Equal(t, token.ExpiresAt.Unix(), receivedTokenByID.ExpiresAt.Unix())
	assert.Equal(t, token.UserID, receivedTokenByID.UserID)
	log.Println("refresh token successfully received by token raw ...")
	receivedTokenByUserID, err := repo.WithTx(tx).GetByUserID(ctx, token.UserID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, token.Token, receivedTokenByUserID.Token)
	assert.Equal(t, token.ExpiresAt.Unix(), receivedTokenByUserID.ExpiresAt.Unix())
	assert.Equal(t, token.UserID, receivedTokenByUserID.UserID)
	log.Println("refresh token successfully received by user id ...")
	err = repo.WithTx(tx).Remove(ctx, token.Token)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.WithTx(tx).GetByToken(ctx, token.Token)
	if err != sql.ErrNoRows {
		t.Fatal(err)
	}
	log.Println("refresh token successfully removed ...")
}

func TestUsersCRUD(t *testing.T) {
	ctx := context.Background()
	log.Println("starts users CRUD...")
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	repo := repos.NewUserRepository(testDB)
	log.Println("users repository created ...")
	insertUser := domain.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "testuser@gmail.com",
	}
	returningUser, err := repo.WithTx(tx).AddOrUpdate(ctx, insertUser)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, insertUser.ID, returningUser.ID)
	assert.Equal(t, insertUser.Username, returningUser.Username)
	assert.Equal(t, insertUser.Email, returningUser.Email)
	log.Println("user user successfully added ...")
	userLastActivityTime := returningUser.LastSeen.UnixNano()
	userCreatedAtTime := returningUser.CreatedAt.UnixNano()
	updatedUser, err := repo.WithTx(tx).AddOrUpdate(ctx, insertUser)
	if err != nil {
		t.Fatal(err)
	}
	if updatedUser.LastSeen.UnixNano() <= userLastActivityTime {
		t.Fatal("expected last seen time to be greater than user's last seen")
	}
	assert.Equal(t, updatedUser.CreatedAt.UnixNano(), userCreatedAtTime)
	log.Println("user successfully updated ...")
	userByID, err := repo.WithTx(tx).GetByID(ctx, insertUser.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, insertUser.ID, userByID.ID)
	assert.Equal(t, insertUser.Username, userByID.Username)
	assert.Equal(t, insertUser.Email, userByID.Email)
	log.Println("user successfully retrieved by id...")
	err = repo.WithTx(tx).Remove(ctx, insertUser.ID)
	if err != nil {
		t.Fatal(err)
	}
	_, err = repo.WithTx(tx).GetByID(ctx, insertUser.ID)
	if err != sql.ErrNoRows {
		t.Fatal(err)
	}
	log.Println("user successfully removed ...")
}

func TestAddPublishCmd(t *testing.T) {
	ctx := context.Background()
	log.Println("starts add publish command...")
	userRepo := repos.NewUserRepository(testDB)
	allianceRepo := repos.NewTenantRepository(testDB)
	testUserdata := domain.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "testuser@gmail.com",
	}
	testUser, err := userRepo.AddOrUpdate(ctx, testUserdata)
	if err != nil {
		t.Fatal(err)
	}
	testTenantData := domain.Tenant{
		ID:      uuid.New(),
		Name:    "testalliance",
		OwnerID: testUser.ID,
	}
	testTenant, err := allianceRepo.Add(ctx, testTenantData)
	if err != nil {
		t.Fatal(err)
	}
	cmd := commands.AddTenantPublishByUser{
		PublishID: uuid.New(),
		TenantID:  testTenant.ID,
		UserID:    testUser.ID,
		S3Data: serializable.S3Screenshot{
			Key:    "s3key",
			Bucket: "s3Bucket",
			S3Name: "selectel",
		},
	}
	txManager := db.NewTxManager(testDB)
	pubRepo := repos.NewPublishRepository(testDB)
	allyPubRepo := junction_repos2.NewTenantPublishRepository(testDB)
	publisher := commands.NewPublisher(
		txManager,
		pubRepo,
		allyPubRepo,
	)
	err = publisher.Handle(ctx, cmd)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("publish command successfully proceed ...")
	checkPub, err := pubRepo.Get(ctx, cmd.PublishID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, cmd.PublishID, checkPub.ID)
	log.Println("publish id retrieved ...")
	checkAllyPub, err := allyPubRepo.GetByID(ctx, cmd.PublishID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, cmd.PublishID, checkAllyPub.PublishID)
	assert.Equal(t, cmd.TenantID, checkAllyPub.TenantID)
	assert.Equal(t, cmd.UserID, checkAllyPub.UserID)
	log.Println("alliance publish successfully retrieved ...")
}

func TestTenantUserRepository(t *testing.T) {
	ctx := context.Background()
	log.Println("starts add tenant users repository...")
	tx, err := testDB.BeginTx(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	userRepo := repos.NewUserRepository(testDB)
	allianceRepo := repos.NewTenantRepository(testDB)
	testUserdata := domain.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "testuser@gmail.com",
	}
	testUser1, err := userRepo.WithTx(tx).AddOrUpdate(ctx, testUserdata)
	if err != nil {
		t.Fatal(err)
	}
	testUserdata2 := domain.User{
		ID:       uuid.New(),
		Username: "testuser2",
		Email:    "testuse2r@gmail.com",
	}
	testUser2, err := userRepo.WithTx(tx).AddOrUpdate(ctx, testUserdata2)
	if err != nil {
		t.Fatal(err)
	}
	testTenantData := domain.Tenant{
		ID:      uuid.New(),
		Name:    "testalliance",
		OwnerID: testUser1.ID,
	}
	testTenant, err := allianceRepo.WithTx(tx).Add(ctx, testTenantData)
	if err != nil {
		t.Fatal(err)
	}
	tenantUserRepo := junction_repos2.NewTenantUserRepository(testDB)
	err = tenantUserRepo.WithTx(tx).Add(ctx, testTenant.ID, testUser1.ID)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("tenant user1 successfully proceed ...")
	err = tenantUserRepo.WithTx(tx).Add(ctx, testTenant.ID, testUser2.ID)
	if err != nil {
		t.Fatal(err)
	}
	log.Println("tenant user2 successfully proceed ...")
	userIDs, err := tenantUserRepo.WithTx(tx).GetUserIDs(ctx, testTenant.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, reqUserID := range []uuid.UUID{testUser1.ID, testUser2.ID} {
		found := false
		for _, userID := range userIDs {
			if userID == reqUserID {
				found = true
			}
		}
		if !found {
			t.Fatalf("user %s not found in tenant %s", reqUserID, testTenant.ID)
		}
	}
	log.Println("tenant user1 & user2 successfully retrieved ...")
	ok, err := tenantUserRepo.WithTx(tx).CheckUser(ctx, testTenant.ID, testUser1.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatalf("user %s not found in tenant %s", testUser1.ID, testTenant.ID)
	}
	t.Log("tenant user1 & user2 successfully retrieved ...")
	found, err := tenantUserRepo.WithTx(tx).CheckUser(ctx, testTenant.ID, testUser2.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatalf("user %s not found in tenant %s", testUser2.ID, testTenant.ID)
	}
}

func TestTenantPublishQuery(t *testing.T) {
	ctx := context.Background()
	log.Println("starts add tenant publish query...")
	tenantPublishQuery := queries.NewGetTenantPublishByIDQuery(testDB)
	tenantPublishRepo := junction_repos2.NewTenantPublishRepository(testDB)
	publishRepo := repos.NewPublishRepository(testDB)
	tenantRepo := repos.NewTenantRepository(testDB)
	userRepo := repos.NewUserRepository(testDB)
	testUserdata := domain.User{
		ID:       uuid.New(),
		Username: "testuser",
		Email:    "testuser@gmail.com",
	}
	user, err := userRepo.AddOrUpdate(ctx, testUserdata)
	if err != nil {
		t.Fatal(err)
	}
	testTenantData := domain.Tenant{
		ID:      uuid.New(),
		Name:    "testalliance",
		OwnerID: user.ID,
	}
	testTenant, err := tenantRepo.Add(ctx, testTenantData)
	if err != nil {
		t.Fatal(err)
	}
	testPublish := domain.PublishedScreenshot{
		ID: uuid.New(),
		S3Data: serializable.S3Screenshot{
			Key:    "s3key",
			Bucket: "s3Bucket",
			S3Name: "selectel",
		},
	}
	err = publishRepo.Add(ctx, testPublish)
	if err != nil {
		t.Fatal(err)
	}
	testTenantPusblish := domain.TenantPublish{
		UserID:    user.ID,
		TenantID:  testTenant.ID,
		PublishID: testPublish.ID,
	}
	_, err = tenantPublishRepo.Add(ctx, testTenantPusblish)
	if err != nil {
		t.Fatal(err)
	}
	dto, err := tenantPublishQuery.Handle(ctx, testPublish.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, dto.TenantID, testTenant.ID)
	assert.Equal(t, dto.PublishID, testPublish.ID)
	assert.Equal(t, dto.UserID, user.ID)
	assert.Equal(t, dto.S3.Key, testPublish.S3Data.Key)
}

func TestFinishedPublishQuery(t *testing.T) {
	ctx := context.Background()
	log.Println("starts add finish publish query...")
	publishRepo := repos.NewPublishRepository(testDB)
	testPublish := domain.PublishedScreenshot{
		ID: uuid.New(),
		S3Data: serializable.S3Screenshot{
			Key:    "s3key",
			Bucket: "s3Bucket",
			S3Name: "selectel",
		},
	}
	err := publishRepo.Add(ctx, testPublish)
	if err != nil {
		t.Fatal(err)
	}
	finishPubRepo := junction_repos2.NewFinishedPublishRepository(testDB)
	e, err := finishPubRepo.Add(ctx, domain.FinishedPublish{
		PublishID: testPublish.ID,
		Result: serializable.NicknameResultWithConflicts{
			Conflicts: []serializable.Conflict{
				serializable.Conflict{
					Similar: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
					Box: [4]serializable.Point{
						serializable.Point{
							X: 23,
							Y: 32,
						},
						serializable.Point{
							X: 23,
							Y: 32,
						},
						serializable.Point{
							X: 23,
							Y: 32,
						},
						serializable.Point{
							X: 23,
							Y: 32,
						},
					},
				},
			},
			NicknameIDs: []uuid.UUID{uuid.New(), uuid.New()},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, e.PublishID, testPublish.ID)
	assert.Equal(t, len(e.Result.Conflicts), 1)
	assert.Equal(t, len(e.Result.Conflicts[0].Similar), 3)
	assert.Equal(t, len(e.Result.NicknameIDs), 2)
}

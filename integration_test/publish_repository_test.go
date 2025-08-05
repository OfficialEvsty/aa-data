package integration_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/OfficialEvsty/aa-data/commands"
	"github.com/OfficialEvsty/aa-data/db"
	"github.com/OfficialEvsty/aa-data/domain"
	"github.com/OfficialEvsty/aa-data/domain/serializable"
	integration "github.com/OfficialEvsty/aa-data/integration_test"
	"github.com/OfficialEvsty/aa-data/queries"
	"github.com/OfficialEvsty/aa-data/repos"
	junction_repos2 "github.com/OfficialEvsty/aa-data/repos/junction"
	"github.com/OfficialEvsty/aa-shared/golinq"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
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

func TestEventBossDropPulling(t *testing.T) {
	ctx := context.Background()
	log.Println("starts add event boss drop pulling...")
	qGetAllEventTemplates := queries.NewGetAllAvailableEventTemplates(testDB)
	eventListDTO, err := qGetAllEventTemplates.Handle(ctx)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(eventListDTO), 14)
	eventIDs := golinq.Map(eventListDTO, func(dto *queries.EventTemplateDTO) int {
		return dto.ID
	})
	qGetAllDropByEvents := queries.NewGetAllAvailableDropFromEvents(testDB)
	dropDTO, err := qGetAllDropByEvents.Handle(ctx, eventIDs)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(dropDTO), 13)
	for eID, bMap := range dropDTO {
		for bID, dropList := range bMap {
			for _, drop := range dropList {
				t.Log(fmt.Sprintf("event: %d boss: %d drop: %s", eID, bID, drop.Name))
			}
		}
	}
}

func TestLunarkRepository(t *testing.T) {
	ctx := context.Background()
	log.Println("starts add lunark repository...")
	lunarkRepo := repos.NewLunarkRepository(testDB)
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	testLunark := domain.Lunark{
		ID:        uuid.New(),
		Name:      "testalliance",
		StartDate: time.Now(),
	}
	err = lunarkRepo.WithTx(tx).Add(ctx, testLunark)
	if err != nil {
		t.Fatal(err)
	}
	receivedLunark, err := lunarkRepo.WithTx(tx).GetByID(ctx, testLunark.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, receivedLunark.Name, testLunark.Name)
	assert.Equal(t, receivedLunark.StartDate.Unix(), testLunark.StartDate.Unix())
	assert.Equal(t, receivedLunark.Opened, true)
	err = lunarkRepo.WithTx(tx).Close(ctx, testLunark.ID)
	if err != nil {
		t.Fatal(err)
	}
	err = lunarkRepo.WithTx(tx).UpdateEndDate(ctx, testLunark.ID, testLunark.StartDate)
	if err != nil {
		t.Fatal(err)
	}
	receivedLunark, err = lunarkRepo.WithTx(tx).GetByID(ctx, testLunark.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, receivedLunark.Opened, false)
	assert.Equal(t, receivedLunark.EndDate.Unix(), testLunark.StartDate.Unix())
}

func TestRaidRepository(t *testing.T) {
	ctx := context.Background()
	log.Println("starts add raid repository...")
	raidRepo := repos.NewRaidRepository(testDB)
	publishRepo := repos.NewPublishRepository(testDB)
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	defer tx.Rollback()
	testPublish := domain.PublishedScreenshot{
		ID: uuid.New(),
		S3Data: serializable.S3Screenshot{
			Key:    "s3key",
			Bucket: "s3Bucket",
			S3Name: "selectel",
		},
	}
	err = publishRepo.WithTx(tx).Add(ctx, testPublish)
	if err != nil {
		t.Fatal(err)
	}
	tm := time.Now()
	testRaid := domain.Raid{
		ID:        uuid.New(),
		PublishID: testPublish.ID,
		Status:    serializable.StatusUnrecognized,
	}
	err = raidRepo.WithTx(tx).Add(ctx, testRaid)
	if err != nil {
		t.Fatal(err)
	}
	reveivedRaid, err := raidRepo.WithTx(tx).GetById(ctx, testRaid.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, reveivedRaid.PublishID, testPublish.ID)
	require.Nil(t, reveivedRaid.RaidAt)
	assert.Equal(t, reveivedRaid.Status, serializable.StatusUnrecognized)
	err = raidRepo.WithTx(tx).UpdateTiming(ctx, testRaid.ID, tm)
	if err != nil {
		t.Fatal(err)
	}
	err = raidRepo.WithTx(tx).UpdateStatus(ctx, testRaid.ID, serializable.StatusResolved)
	if err != nil {
		t.Fatal(err)
	}
	err = raidRepo.WithTx(tx).UpdateAttendance(ctx, testRaid.ID, 45)
	if err != nil {
		t.Fatal(err)
	}
	reveivedRaid, err = raidRepo.WithTx(tx).GetById(ctx, testRaid.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, reveivedRaid.Status, serializable.StatusResolved)
	assert.Equal(t, reveivedRaid.RaidAt.Unix(), tm.Unix())
	assert.Equal(t, reveivedRaid.Attendance, 45)
}

func TestAllRaidsReceiving(t *testing.T) {
	ctx := context.Background()
	log.Println("starts add event raids receiving...")
	tx, err := testDB.Begin()
	if err != nil {
		t.Fatal(err)
	}
	ctx = db.WithTxInContext(ctx, tx)
	defer tx.Rollback()
	getAllRaidsByTenant := queries.NewGetAllCompletedRaidsByTenantID(tx)
	userRepo := repos.NewUserRepository(testDB)
	tenantRepo := repos.NewTenantRepository(testDB)
	testUser := domain.User{
		ID:       uuid.New(),
		Email:    "user@example.com",
		Username: "user",
	}
	user, err := userRepo.WithTx(tx).AddOrUpdate(ctx, testUser)
	if err != nil {
		t.Fatal(err)
	}
	testTenant := domain.Tenant{
		ID:      uuid.New(),
		Name:    "testTenant",
		OwnerID: user.ID,
	}
	tenant, err := tenantRepo.WithTx(tx).Add(ctx, testTenant)
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
	publishRepo := repos.NewPublishRepository(testDB)
	err = publishRepo.WithTx(tx).Add(ctx, testPublish)
	if err != nil {
		t.Fatal(err)
	}
	tenantPublishRepo := junction_repos2.NewTenantPublishRepository(testDB)
	testTenantPublish := domain.TenantPublish{
		UserID:    user.ID,
		TenantID:  tenant.ID,
		PublishID: testPublish.ID,
	}
	_, err = tenantPublishRepo.WithTx(tx).Add(ctx, testTenantPublish)
	if err != nil {
		t.Fatal(err)
	}

	raidRepo := repos.NewRaidRepository(testDB)
	testRaid1 := domain.Raid{
		ID:        uuid.New(),
		PublishID: testPublish.ID,
		Status:    serializable.StatusUnrecognized,
	}
	err = raidRepo.WithTx(tx).Add(ctx, testRaid1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(fmt.Sprintf("raid 1: created"))
	getAllIncompleteRaids := queries.NewGetAllIncompleteRaidQuery(testDB)
	incompleteRaids, err := getAllIncompleteRaids.WithTx(tx).Handle(ctx, testUser.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(incompleteRaids), 1)
	assert.Equal(t, incompleteRaids[0].PublishID, testPublish.ID)
	assert.Equal(t, incompleteRaids[0].Status, serializable.StatusUnrecognized)
	require.Nil(t, incompleteRaids[0].RaidAt)
	assert.Equal(t, incompleteRaids[0].UserID, testUser.ID)
	raidAt := time.Now()

	getEventTemplates := queries.NewGetAllAvailableEventTemplates(testDB)
	eventTemplates, err := getEventTemplates.WithTx(tx).Handle(ctx)
	if err != nil {
		t.Fatal(err)
	}
	chosenEvent := golinq.Where(eventTemplates, func(e *queries.EventTemplateDTO) bool {
		return e.ID == 14
	})[0]
	t.Log(fmt.Sprintf("event template %d: chosen", chosenEvent.ID))
	raidEventRepo := junction_repos2.NewRaidEventRepository(testDB)
	testRaidEvent := domain.RaidEvent{
		RaidID:  testRaid1.ID,
		EventID: chosenEvent.ID,
	}
	err = raidEventRepo.WithTx(tx).Add(ctx, testRaidEvent)
	if err != nil {
		t.Fatal(err)
	}
	_, err = raidEventRepo.All(ctx, testRaid1.ID)
	if err != nil {
		t.Fatal(err)
	}
	getAllAvailableDropFromEvents := queries.NewGetAllAvailableDropFromEvents(testDB)

	txManager := db.NewTxManager(testDB)
	bossRepo := repos.NewBossesRepository(testDB)
	itemRepo := repos.NewItemRepository(testDB)
	bossImporter := commands.NewBossesImporter(txManager, bossRepo, itemRepo)
	cmd := commands.AddBossesDropAndItemsCommand{
		Bosses: []*domain.AABoss{
			{
				ID:            1111,
				Name:          "testBoss",
				ImageGradeURL: "f",
				ImageURL:      "f",
				Loot: serializable.DropItemList{
					{
						ItemID: 1110,
						Rate:   "4",
					},
				},
			},
		},
		Items: []*domain.AAItemTemplate{
			{
				ID:       1110,
				Name:     "testItem",
				Tier:     3,
				ImageURL: "f",
				TierURL:  "f",
			},
		},
	}
	err = bossImporter.Handle(ctx, cmd)
	if err != nil {
		t.Fatal(err)
	}
	drops, err := getAllAvailableDropFromEvents.WithTx(tx).Handle(ctx, []int{14})
	if err != nil {
		t.Fatal(err)
	}
	_ = drops
	var dList []*serializable.DropItem
	for _, bossDrop := range drops {
		for _, dropList := range bossDrop {
			for _, drop := range dropList {
				t.Log(fmt.Sprintf("id: %d, %s", drop.ID, drop.Name))
				dList = append(dList, &serializable.DropItem{
					ItemID: drop.ID,
					Rate:   drop.Quantity,
				})
			}
		}
	}
	raidItemRepo := junction_repos2.NewRaidItemRepository(testDB)
	err = raidItemRepo.WithTx(tx).AddOrUpdateItems(ctx, testRaid1.ID, dList)
	if err != nil {
		t.Fatal(err)
	}
	raidDrop, err := raidItemRepo.WithTx(tx).GetItems(ctx, testRaid1.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, len(raidDrop), len(dList))
	assert.Equal(t, raidDrop[0].ItemID, int64(1110))
	getOpenedLunark := queries.NewGetOpenedLunarkByTenant(testDB)
	openedLunark, err := getOpenedLunark.WithTx(tx).Handle(ctx, testTenant.ID)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatal(err)
		}
	}
	if openedLunark == nil {
		addLunarkToTenant := commands.NewLunarkImporter(testDB)
		addLunarkCmd := commands.AddLunarkAttendedToTenantCommand{
			TenantID:  testTenant.ID,
			LunarkID:  uuid.New(),
			Name:      "testLunark",
			StartDate: raidAt,
		}
		err = addLunarkToTenant.Handle(ctx, addLunarkCmd)
		if err != nil {
			t.Fatal(err)
		}
	}
	openedLunark, err = getOpenedLunark.WithTx(tx).Handle(ctx, testTenant.ID)
	if err != nil {
		t.Fatal(err)
	}
	require.NotNil(t, openedLunark)
	addRaidInLunarkCommand := commands.NewRaidInLunarkCommand(testDB)
	raidAt = time.Now()
	raidInLunarkCommand := commands.AddRaidInLunarkCommand{
		RaidID:     testRaid1.ID,
		LunarkID:   openedLunark.LunarkID,
		Status:     serializable.StatusResolved,
		RaidAt:     &raidAt,
		Attendance: 45,
	}
	err = addRaidInLunarkCommand.Handle(ctx, raidInLunarkCommand)
	if err != nil {
		t.Fatal(err)
	}

	raid, err := raidRepo.WithTx(tx).GetById(ctx, testRaid1.ID)
	if err != nil {
		t.Fatal(err)
	}
	lunarkRepo := repos.NewLunarkRepository(testDB)
	require.NotNil(t, raid.RaidAt)
	assert.Equal(t, raid.RaidAt.Unix(), raidAt.Unix())
	assert.Equal(t, raid.Status, serializable.StatusResolved)
	lunark, err := lunarkRepo.WithTx(tx).GetByID(ctx, openedLunark.LunarkID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, lunark.ID, openedLunark.LunarkID)
	assert.Equal(t, lunark.EndDate.Unix(), raidAt.Unix())
	guildNicknameRepo := junction_repos2.NewGuildNicknameRepository(testDB)
	testGuildID, _ := uuid.Parse("5f25e254-d493-4779-8903-0e67aae7443c")
	nicknameIDs, err := guildNicknameRepo.WithTx(tx).GetMembers(ctx, testGuildID)
	if err != nil {
		t.Fatal(err)
	}
	addNicknamesToRaid := commands.NewAttendanceController(testDB)
	addNicknamesCmd := commands.AddNicknamesToRaidCommand{
		RaidID:      testRaid1.ID,
		NicknameIDs: nicknameIDs,
		Attendance:  100,
	}
	err = addNicknamesToRaid.Handle(ctx, addNicknamesCmd)
	if err != nil {
		t.Fatal(err)
	}
	allRaids, err := getAllRaidsByTenant.WithTx(tx).Handle(ctx, tenant.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range allRaids.Raids {
		t.Log(fmt.Sprintf("raid: %d user id provided: %s, time: %s", r.ID, r.ProviderID, r.OccurredAt))
	}
	assert.Equal(t, allRaids.Raids[testRaid1.ID].Attendance, 100)
	assert.Equal(t, allRaids.Raids[testRaid1.ID].ProviderID, user.ID)
	assert.Equal(t, allRaids.Raids[testRaid1.ID].ParticipantCount, uint(2))
	assert.Equal(t, len(allRaids.Raids[testRaid1.ID].EventIDs), 1)
	assert.Equal(t, allRaids.Lunark.StartDate.Unix(), raidAt.Unix())
	assert.Equal(t, allRaids.Lunark.EndDate.Unix(), raidAt.Unix())
	t.Log(allRaids.Lunark.Name)
	for _, e := range allRaids.Events {
		t.Log(fmt.Sprintf("event: %d name: %s", e.TemplateID, e.Name))
	}
}

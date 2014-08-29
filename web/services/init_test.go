package services_test

import (
    "database/sql"
    "errors"
    "net/http"
    "testing"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/cloudfoundry-incubator/notifications/web/services"

    "github.com/dgrijalva/jwt-go"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestWebHandlersServicesSuite(t *testing.T) {
    RegisterFastTokenSigningMethod()

    RegisterFailHandler(Fail)
    RunSpecs(t, "Web Handlers Services Suite")
}

const (
    UAAPrivateKey = "PRIVATE-KEY"
    UAAPublicKey  = "PUBLIC-KEY"
)

type SigningMethodFast struct{}

func (m SigningMethodFast) Alg() string {
    return "FAST"
}

func (m SigningMethodFast) Sign(signingString string, key []byte) (string, error) {
    signature := jwt.EncodeSegment([]byte(signingString + "SUPERFAST"))
    return signature, nil
}

func (m SigningMethodFast) Verify(signingString, signature string, key []byte) (err error) {
    if signature != jwt.EncodeSegment([]byte(signingString+"SUPERFAST")) {
        return errors.New("Signature is invalid")
    }

    return nil
}

func RegisterFastTokenSigningMethod() {
    jwt.RegisterSigningMethod("FAST", func() jwt.SigningMethod {
        return SigningMethodFast{}
    })
}

func BuildToken(header map[string]interface{}, claims map[string]interface{}) string {
    config.UAAPublicKey = UAAPublicKey

    alg := header["alg"].(string)
    signingMethod := jwt.GetSigningMethod(alg)
    token := jwt.New(signingMethod)
    token.Header = header
    token.Claims = claims

    signed, err := token.SignedString([]byte(UAAPrivateKey))
    if err != nil {
        panic(err)
    }

    return signed
}

type FakeCourier struct {
    Error             error
    Responses         []postal.Response
    DispatchArguments []interface{}
}

func NewFakeCourier() *FakeCourier {
    return &FakeCourier{
        Responses:         make([]postal.Response, 0),
        DispatchArguments: make([]interface{}, 0),
    }
}

func (fake *FakeCourier) Dispatch(token string, guid postal.TypedGUID, options postal.Options, conn models.ConnectionInterface) ([]postal.Response, error) {
    fake.DispatchArguments = []interface{}{token, guid, options}
    return fake.Responses, fake.Error
}

type FakeErrorWriter struct {
    Error error
}

func NewFakeErrorWriter() *FakeErrorWriter {
    return &FakeErrorWriter{}
}

func (writer *FakeErrorWriter) Write(w http.ResponseWriter, err error) {
    writer.Error = err
}

type FakeDBResult struct{}

func (fake FakeDBResult) LastInsertId() (int64, error) {
    return 0, nil
}

func (fake FakeDBResult) RowsAffected() (int64, error) {
    return 0, nil
}

type FakeDBConn struct {
    BeginWasCalled    bool
    CommitWasCalled   bool
    RollbackWasCalled bool
}

func (conn *FakeDBConn) Begin() error {
    conn.BeginWasCalled = true
    return nil
}

func (conn *FakeDBConn) Commit() error {
    conn.CommitWasCalled = true
    return nil
}

func (conn *FakeDBConn) Rollback() error {
    conn.RollbackWasCalled = true
    return nil
}

func (conn *FakeDBConn) Exec(query string, args ...interface{}) (sql.Result, error) {
    return FakeDBResult{}, nil
}

func (conn FakeDBConn) Delete(list ...interface{}) (int64, error) {
    return 0, nil
}

func (conn FakeDBConn) Insert(list ...interface{}) error {
    return nil
}

func (conn FakeDBConn) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
    return []interface{}{}, nil
}

func (conn FakeDBConn) SelectOne(i interface{}, query string, args ...interface{}) error {
    return nil
}

func (conn FakeDBConn) Update(list ...interface{}) (int64, error) {
    return 0, nil
}

type FakeClientsRepo struct {
    Clients     map[string]models.Client
    UpsertError error
    FindError   error
}

func NewFakeClientsRepo() *FakeClientsRepo {
    return &FakeClientsRepo{
        Clients: make(map[string]models.Client),
    }
}

func (fake *FakeClientsRepo) Create(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
    if _, ok := fake.Clients[client.ID]; ok {
        return client, models.ErrDuplicateRecord{}
    }
    fake.Clients[client.ID] = client
    return client, nil
}

func (fake *FakeClientsRepo) Update(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
    fake.Clients[client.ID] = client
    return client, nil
}

func (fake *FakeClientsRepo) Upsert(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
    fake.Clients[client.ID] = client
    return client, fake.UpsertError
}

func (fake *FakeClientsRepo) Find(conn models.ConnectionInterface, id string) (models.Client, error) {
    if client, ok := fake.Clients[id]; ok {
        return client, fake.FindError
    }
    return models.Client{}, models.ErrRecordNotFound{}
}

type FakeKindsRepo struct {
    Kinds         map[string]models.Kind
    UpsertError   error
    TrimError     error
    FindError     error
    TrimArguments []interface{}
}

func NewFakeKindsRepo() *FakeKindsRepo {
    return &FakeKindsRepo{
        Kinds:         make(map[string]models.Kind),
        TrimArguments: make([]interface{}, 0),
    }
}

func (fake *FakeKindsRepo) Create(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
    key := kind.ID + kind.ClientID
    if _, ok := fake.Kinds[key]; ok {
        return kind, models.ErrDuplicateRecord{}
    }
    fake.Kinds[key] = kind
    return kind, nil
}

func (fake *FakeKindsRepo) Update(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
    key := kind.ID + kind.ClientID
    fake.Kinds[key] = kind
    return kind, nil
}

func (fake *FakeKindsRepo) Upsert(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
    key := kind.ID + kind.ClientID
    fake.Kinds[key] = kind
    return kind, fake.UpsertError
}

func (fake *FakeKindsRepo) Find(conn models.ConnectionInterface, id, clientID string) (models.Kind, error) {
    key := id + clientID
    if kind, ok := fake.Kinds[key]; ok {
        return kind, fake.FindError
    }
    return models.Kind{}, models.ErrRecordNotFound{}
}

func (fake *FakeKindsRepo) Trim(conn models.ConnectionInterface, clientID string, kindIDs []string) (int, error) {
    fake.TrimArguments = []interface{}{clientID, kindIDs}
    return 0, fake.TrimError
}

type FakeFinder struct {
    Clients            map[string]models.Client
    Kinds              map[string]models.Kind
    ClientAndKindError error
}

func NewFakeFinder() *FakeFinder {
    return &FakeFinder{
        Clients: make(map[string]models.Client),
        Kinds:   make(map[string]models.Kind),
    }
}

func (finder *FakeFinder) ClientAndKind(clientID, kindID string) (models.Client, models.Kind, error) {
    return finder.Clients[clientID], finder.Kinds[kindID+"|"+clientID], finder.ClientAndKindError
}

type FakeRegistrar struct {
    RegisterArguments []interface{}
    RegisterError     error
    PruneArguments    []interface{}
    PruneError        error
}

func NewFakeRegistrar() *FakeRegistrar {
    return &FakeRegistrar{}
}

func (fake *FakeRegistrar) Register(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
    fake.RegisterArguments = []interface{}{conn, client, kinds}
    return fake.RegisterError
}

func (fake *FakeRegistrar) Prune(conn models.ConnectionInterface, client models.Client, kinds []models.Kind) error {
    fake.PruneArguments = []interface{}{conn, client, kinds}
    return fake.PruneError
}

type FakePreference struct {
    ReturnValue   services.PreferencesBuilder
    ExecuteErrors bool
    UserGUID      string
}

func NewFakePreference(returnValue services.PreferencesBuilder) *FakePreference {
    return &FakePreference{
        ReturnValue: returnValue,
    }
}

func (fake *FakePreference) Execute(userGUID string) (services.PreferencesBuilder, error) {
    fake.UserGUID = userGUID
    if fake.ExecuteErrors {
        return fake.ReturnValue, errors.New("Meltdown")
    }
    return fake.ReturnValue, nil
}

type FakePreferencesRepo struct {
    NonCriticalPreferences []models.Preference
    FindError              error
}

func NewFakePreferencesRepo(nonCriticalPreferences []models.Preference) *FakePreferencesRepo {
    return &FakePreferencesRepo{
        NonCriticalPreferences: nonCriticalPreferences,
    }
}

func (fake FakePreferencesRepo) FindNonCriticalPreferences(conn models.ConnectionInterface, userGUID string) ([]models.Preference, error) {
    return fake.NonCriticalPreferences, fake.FindError
}

type FakePreferenceUpdater struct {
    ExecuteArguments []interface{}
}

func NewFakePreferenceUpdater() *FakePreferenceUpdater {
    return &FakePreferenceUpdater{}
}

func (fake *FakePreferenceUpdater) Execute(conn models.ConnectionInterface, preferences []models.Preference, userID string) error {
    fake.ExecuteArguments = append(fake.ExecuteArguments, preferences, userID)
    return nil
}

type FakeUnsubscribesRepo struct {
    Unsubscribes map[string]models.Unsubscribe
}

func NewFakeUnsubscribesRepo() *FakeUnsubscribesRepo {
    return &FakeUnsubscribesRepo{
        Unsubscribes: map[string]models.Unsubscribe{},
    }
}

func (fake *FakeUnsubscribesRepo) Create(conn models.ConnectionInterface, unsubscribe models.Unsubscribe) (models.Unsubscribe, error) {
    key := unsubscribe.ClientID + unsubscribe.KindID + unsubscribe.UserID
    if _, ok := fake.Unsubscribes[key]; ok {
        return unsubscribe, models.ErrDuplicateRecord{}
    }
    fake.Unsubscribes[key] = unsubscribe
    return unsubscribe, nil
}

func (fake *FakeUnsubscribesRepo) Upsert(conn models.ConnectionInterface, unsubscribe models.Unsubscribe) (models.Unsubscribe, error) {
    key := unsubscribe.ClientID + unsubscribe.KindID + unsubscribe.UserID
    fake.Unsubscribes[key] = unsubscribe
    return unsubscribe, nil
}

func (fake *FakeUnsubscribesRepo) Find(conn models.ConnectionInterface, clientID string, kindID string, userID string) (models.Unsubscribe, error) {
    key := clientID + kindID + userID
    if unsubscribe, ok := fake.Unsubscribes[key]; ok {
        return unsubscribe, models.ErrDuplicateRecord{}
    }
    return models.Unsubscribe{}, models.ErrRecordNotFound{}
}

func (fake *FakeUnsubscribesRepo) Destroy(conn models.ConnectionInterface, unsubscribe models.Unsubscribe) (int, error) {
    key := unsubscribe.ClientID + unsubscribe.KindID + unsubscribe.UserID
    delete(fake.Unsubscribes, key)
    return 0, nil
}

package postal_test

import (
    "database/sql"
    "errors"
    "fmt"
    "net/http"
    "testing"

    "github.com/cloudfoundry-incubator/notifications/cf"
    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/gobble"
    "github.com/cloudfoundry-incubator/notifications/mail"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/dgrijalva/jwt-go"
    "github.com/nu7hatch/gouuid"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
    "github.com/pivotal-cf/uaa-sso-golang/uaa"
)

const (
    UAAPrivateKey = "PRIVATE-KEY"
    UAAPublicKey  = "PUBLIC-KEY"
)

func TestPostalSuite(t *testing.T) {
    RegisterFastTokenSigningMethod()

    RegisterFailHandler(Fail)
    RunSpecs(t, "Postal Suite")
}

func RegisterFastTokenSigningMethod() {
    jwt.RegisterSigningMethod("FAST", func() jwt.SigningMethod {
        return SigningMethodFast{}
    })
}

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

type FakeCloudController struct {
    CurrentToken             string
    GetUsersBySpaceGuidError error
    LoadSpaceError           error
    LoadOrganizationError    error
    UsersBySpaceGuid         map[string][]cf.CloudControllerUser
    Spaces                   map[string]cf.CloudControllerSpace
    Orgs                     map[string]cf.CloudControllerOrganization
}

func NewFakeCloudController() *FakeCloudController {
    return &FakeCloudController{
        UsersBySpaceGuid: make(map[string][]cf.CloudControllerUser),
    }
}

func (fake *FakeCloudController) GetUsersBySpaceGuid(guid, token string) ([]cf.CloudControllerUser, error) {
    fake.CurrentToken = token

    if users, ok := fake.UsersBySpaceGuid[guid]; ok {
        return users, fake.GetUsersBySpaceGuidError
    } else {
        return make([]cf.CloudControllerUser, 0), fake.GetUsersBySpaceGuidError
    }
}

func (fake *FakeCloudController) LoadSpace(guid, token string) (cf.CloudControllerSpace, error) {
    if fake.LoadSpaceError != nil {
        return cf.CloudControllerSpace{}, fake.LoadSpaceError
    }

    if space, ok := fake.Spaces[guid]; ok {
        return space, nil
    } else {
        return cf.CloudControllerSpace{}, cf.NewFailure(http.StatusNotFound, fmt.Sprintf(`{"code":40004,"description":"The app space could not be found: %s","error_code":"CF-SpaceNotFound"}`, guid))
    }
}

func (fake *FakeCloudController) LoadOrganization(guid, token string) (cf.CloudControllerOrganization, error) {
    if fake.LoadOrganizationError != nil {
        return cf.CloudControllerOrganization{}, fake.LoadOrganizationError
    }

    if org, ok := fake.Orgs[guid]; ok {
        return org, nil
    } else {
        return cf.CloudControllerOrganization{}, cf.NewFailure(http.StatusNotFound, fmt.Sprintf(`{"code":30003,"description":"The organization could not be found: %s","error_code":"CF-OrganizationNotFound"}`, guid))
    }
}

type FakeUAAClient struct {
    ClientToken      uaa.Token
    ClientTokenError error
    UsersByID        map[string]uaa.User
    ErrorForUserByID error
    AccessToken      string
}

func (fake *FakeUAAClient) SetToken(token string) {
    fake.AccessToken = token
}

func (fake FakeUAAClient) GetClientToken() (uaa.Token, error) {
    return fake.ClientToken, fake.ClientTokenError
}

func (fake FakeUAAClient) UsersByIDs(ids ...string) ([]uaa.User, error) {
    users := []uaa.User{}
    for _, id := range ids {
        if user, ok := fake.UsersByID[id]; ok {
            users = append(users, user)
        }
    }

    return users, fake.ErrorForUserByID
}

type FakeReceiptsRepo struct {
    CreateUserGUIDs     []string
    ClientID            string
    KindID              string
    CreateReceiptsError bool
}

func NewFakeReceiptsRepo() FakeReceiptsRepo {
    return FakeReceiptsRepo{}
}

func (fake *FakeReceiptsRepo) CreateReceipts(conn models.ConnectionInterface, userGUIDs []string, clientID, kindID string) error {
    if fake.CreateReceiptsError {
        return errors.New("a database error")
    }

    fake.CreateUserGUIDs = userGUIDs
    fake.ClientID = clientID
    fake.KindID = kindID

    return nil
}

type FakeMailClient struct {
    messages       []mail.Message
    errorOnSend    bool
    errorOnConnect bool
}

func (fake *FakeMailClient) Connect() error {
    if fake.errorOnConnect {
        return errors.New("BOOM!")
    }
    return nil
}

func (fake *FakeMailClient) Send(msg mail.Message) error {
    err := fake.Connect()
    if err != nil {
        return err
    }

    if fake.errorOnSend {
        return errors.New("BOOM!")
    }

    fake.messages = append(fake.messages, msg)
    return nil
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

var FakeGuidGenerator = postal.GUIDGenerationFunc(func() (*uuid.UUID, error) {
    guid := uuid.UUID([16]byte{0xDE, 0xAD, 0xBE, 0xEF, 0xAA, 0xBB, 0xCC, 0xDD, 0xEE, 0xFF, 0x00, 0x11, 0x22, 0x33, 0x44, 0x55})
    return &guid, nil
})

type FakeQueue struct {
    jobs         chan gobble.Job
    pk           int
    EnqueueError error
}

func NewFakeQueue() *FakeQueue {
    return &FakeQueue{
        jobs: make(chan gobble.Job),
    }
}

func (fake *FakeQueue) Enqueue(job gobble.Job) (gobble.Job, error) {
    if fake.EnqueueError != nil {
        return job, fake.EnqueueError
    }
    fake.pk++
    job.ID = fake.pk
    go func(job gobble.Job) {
        fake.jobs <- job
    }(job)

    return job, nil
}

func (fake *FakeQueue) Reserve(string) <-chan gobble.Job {
    return fake.jobs
}

func (fake *FakeQueue) Dequeue(job gobble.Job) {
}

func (fake *FakeQueue) Requeue(job gobble.Job) {
    go func(job gobble.Job) {
        fake.jobs <- job
    }(job)
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

func (conn *FakeDBConn) Transaction() models.TransactionInterface {
    return conn
}

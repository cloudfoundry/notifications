package handlers_test

import (
    "errors"
    "net/http"
    "testing"

    "github.com/cloudfoundry-incubator/notifications/config"
    "github.com/cloudfoundry-incubator/notifications/models"
    "github.com/cloudfoundry-incubator/notifications/postal"
    "github.com/dgrijalva/jwt-go"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestWebHandlersSuite(t *testing.T) {
    RegisterFastTokenSigningMethod()

    RegisterFailHandler(Fail)
    RunSpecs(t, "Web Handlers Suite")
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
    Error     error
    Responses []postal.Response
}

func NewFakeCourier() *FakeCourier {
    return &FakeCourier{
        Responses: make([]postal.Response, 0),
    }
}

func (fake FakeCourier) Dispatch(token string, guid postal.TypedGUID, options postal.Options) ([]postal.Response, error) {
    return fake.Responses, fake.Error
}

type FakeErrorWriter struct {
    Error error
}

func (writer *FakeErrorWriter) Write(w http.ResponseWriter, err error) {
    writer.Error = err
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
}

func NewFakeClientsRepo() *FakeClientsRepo {
    return &FakeClientsRepo{
        Clients: make(map[string]models.Client),
    }
}

func (fake *FakeClientsRepo) Create(conn models.ConnectionInterface, client models.Client) (models.Client, error) {
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
        return client, nil
    }
    return models.Client{}, models.ErrRecordNotFound{}
}

type FakeKindsRepo struct {
    Kinds         map[string]models.Kind
    UpsertError   error
    TrimError     error
    TrimArguments []interface{}
}

func NewFakeKindsRepo() *FakeKindsRepo {
    return &FakeKindsRepo{
        Kinds:         make(map[string]models.Kind),
        TrimArguments: make([]interface{}, 0),
    }
}

func (fake *FakeKindsRepo) Create(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
    fake.Kinds[kind.ID] = kind
    return kind, nil
}

func (fake *FakeKindsRepo) Update(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
    fake.Kinds[kind.ID] = kind
    return kind, nil
}

func (fake *FakeKindsRepo) Upsert(conn models.ConnectionInterface, kind models.Kind) (models.Kind, error) {
    fake.Kinds[kind.ID] = kind
    return kind, fake.UpsertError
}

func (fake *FakeKindsRepo) Find(conn models.ConnectionInterface, id string) (models.Kind, error) {
    if kind, ok := fake.Kinds[id]; ok {
        return kind, nil
    }
    return models.Kind{}, models.ErrRecordNotFound{}
}

func (fake *FakeKindsRepo) Trim(conn models.ConnectionInterface, clientID string, kindIDs []string) (int, error) {
    fake.TrimArguments = []interface{}{clientID, kindIDs}
    return 0, fake.TrimError
}

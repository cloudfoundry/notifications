package models_test

import (
	"database/sql"
	"testing"

	"github.com/cloudfoundry-incubator/notifications/application"
	"github.com/cloudfoundry-incubator/notifications/db"
	"github.com/cloudfoundry-incubator/notifications/v1/models"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestModelsSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "v1/models")
}

var sqlDB *sql.DB

var _ = BeforeEach(func() {
	env, err := application.NewEnvironment()
	Expect(err).NotTo(HaveOccurred())

	sqlDB, err = sql.Open("mysql", env.DatabaseURL)
	Expect(err).NotTo(HaveOccurred())
})

func findReceipt(conn db.ConnectionInterface, userGUID, clientID, kindID string) (models.Receipt, error) {
	receipt := models.Receipt{}
	err := conn.SelectOne(&receipt, "SELECT * FROM  `receipts` WHERE `user_guid` = ? AND `client_id` = ? AND `kind_id` = ?", userGUID, clientID, kindID)
	if err != nil {
		return models.Receipt{}, err
	}

	return receipt, nil
}

func createReceipt(conn db.ConnectionInterface, receipt models.Receipt) (models.Receipt, error) {
	err := conn.Insert(&receipt)
	if err != nil {
		return models.Receipt{}, err
	}

	return receipt, nil
}

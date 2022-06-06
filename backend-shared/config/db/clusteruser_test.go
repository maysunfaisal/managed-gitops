package db_test

import (
	"context"
	"strings"

	db "github.com/maysunfaisal/managed-gitops/backend-shared/config/db"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClusterUser Tests", func() {
	Context("It should execute all DB functions for ClusterUser", func() {
		It("Should execute all ClusterUser Functions", func() {
			err := db.SetupForTestingDBGinkgo()
			Expect(err).To(BeNil())

			ctx := context.Background()

			dbq, err := db.NewUnsafePostgresDBQueries(true, true)
			Expect(err).To(BeNil())
			defer dbq.CloseDatabase()

			user := &db.ClusterUser{
				Clusteruser_id: "test-user-id",
				User_name:      "tirthuser",
			}
			err = dbq.CreateClusterUser(ctx, user)
			Expect(err).To(BeNil())

			retrieveUser := &db.ClusterUser{
				User_name: "tirthuser",
			}
			err = dbq.GetClusterUserByUsername(ctx, retrieveUser)
			Expect(err).To(BeNil())
			Expect(user).Should(Equal(retrieveUser))
			retrieveUser = &db.ClusterUser{
				Clusteruser_id: user.Clusteruser_id,
			}
			err = dbq.GetClusterUserById(ctx, retrieveUser)
			Expect(err).To(BeNil())
			Expect(user).Should(Equal(retrieveUser))
			rowsAffected, err := dbq.DeleteClusterUserById(ctx, retrieveUser.Clusteruser_id)
			Expect(err).To(BeNil())
			Expect(rowsAffected).Should(Equal(1))
			err = dbq.GetClusterUserById(ctx, retrieveUser)
			Expect(true).To(Equal(db.IsResultNotFoundError(err)))

			// Set the invalid value
			user.User_name = strings.Repeat("abc", 100)
			err = dbq.CreateClusterUser(ctx, user)
			Expect(db.IsMaxLengthError(err)).To(Equal(true))

		})
	})
})

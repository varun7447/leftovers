package iam_test

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	awsiam "github.com/aws/aws-sdk-go/service/iam"
	"github.com/genevievelesperance/leftovers/aws/iam"
	"github.com/genevievelesperance/leftovers/aws/iam/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Roles", func() {
	var (
		iamClient *fakes.IAMClient
		logger    *fakes.Logger

		roles iam.Roles
	)

	BeforeEach(func() {
		iamClient = &fakes.IAMClient{}
		logger = &fakes.Logger{}

		roles = iam.NewRoles(iamClient, logger)
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			logger.PromptCall.Returns.Proceed = true
			iamClient.ListRolesCall.Returns.Output = &awsiam.ListRolesOutput{
				Roles: []*awsiam.Role{{
					RoleName: aws.String("banana"),
				}},
			}
		})

		It("deletes iam roles", func() {
			err := roles.Delete()
			Expect(err).NotTo(HaveOccurred())

			Expect(iamClient.DeleteRoleCall.CallCount).To(Equal(1))
			Expect(iamClient.DeleteRoleCall.Receives.Input.RoleName).To(Equal(aws.String("banana")))
			Expect(logger.PrintfCall.Messages).To(Equal([]string{"SUCCESS deleting role banana\n"}))
		})

		Context("when the client fails to list roles", func() {
			BeforeEach(func() {
				iamClient.ListRolesCall.Returns.Error = errors.New("some error")
			})

			It("does not try deleting them", func() {
				err := roles.Delete()
				Expect(err.Error()).To(Equal("Listing roles: some error"))

				Expect(iamClient.DeleteRoleCall.CallCount).To(Equal(0))
			})
		})

		Context("when the client fails to delete the role", func() {
			BeforeEach(func() {
				iamClient.DeleteRoleCall.Returns.Error = errors.New("some error")
			})

			It("returns the error", func() {
				err := roles.Delete()
				Expect(err).NotTo(HaveOccurred())

				Expect(logger.PrintfCall.Messages).To(Equal([]string{"ERROR deleting role banana: some error\n"}))
			})
		})

		Context("when the user responds no to the prompt", func() {
			BeforeEach(func() {
				logger.PromptCall.Returns.Proceed = false
			})

			It("returns the error", func() {
				err := roles.Delete()
				Expect(err).NotTo(HaveOccurred())

				Expect(logger.PromptCall.Receives.Message).To(Equal("Are you sure you want to delete role banana?"))
				Expect(iamClient.DeleteRoleCall.CallCount).To(Equal(0))
			})
		})
	})
})
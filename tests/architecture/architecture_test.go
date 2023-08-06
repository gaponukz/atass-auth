package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func TestArchitecture(t *testing.T) {
	t.Run("domain layer must have no dependencies", func(t *testing.T) {
		mockT := new(testingT)

		for _, domainPackage := range domainLayer() {
			for _, applicationPackage := range applicationLayer() {
				archtest.Package(t, domainPackage).ShouldNotDependOn(applicationPackage)
			}
		}

		for _, domainPackage := range domainLayer() {
			for _, infrastructurePackage := range infrastructureLayer() {
				archtest.Package(t, domainPackage).ShouldNotDependOn(infrastructurePackage)
			}
		}

		for _, domainPackage := range domainLayer() {
			for _, interfacePackage := range interfaceLayer() {
				archtest.Package(t, domainPackage).ShouldNotDependOn(interfacePackage)
			}
		}
		assertNoError(t, mockT)
	})

	t.Run("application layer must not depend on interface", func(t *testing.T) {
		mockT := new(testingT)

		for _, applicationPackage := range applicationLayer() {
			for _, interfacePackage := range interfaceLayer() {
				archtest.Package(t, applicationPackage).ShouldNotDependOn(interfacePackage)
			}
		}

		assertNoError(t, mockT)
	})

	t.Run("application layer must not depend on infrastructure", func(t *testing.T) {
		mockT := new(testingT)

		for _, applicationPackage := range applicationLayer() {
			for _, infrastructurePackage := range infrastructureLayer() {
				archtest.Package(t, applicationPackage).ShouldNotDependOn(infrastructurePackage)
			}
		}

		assertNoError(t, mockT)
	})
}

package architecture

import (
	"testing"

	"github.com/matthewmcnew/archtest"
)

func noDependenciesPackages() []string {
	return []string{
		"auth/src/entities",
		"auth/src/utils",
		"auth/src/errors",
	}
}

func getHelpers() []string {
	return []string{
		"auth/src/storage",
		"auth/src/notifier",
		"auth/src/security",
		"auth/src/settings",
	}
}

func getServices() []string {
	return []string{
		"auth/src/services/passreset",
		"auth/src/services/settings",
		"auth/src/services/signup",
		"auth/src/services/signin",
	}
}

func TestArchitecture(t *testing.T) {
	t.Run("some packages must have no dependencies", func(t *testing.T) {
		mockT := new(testingT)

		for _, p1 := range noDependenciesPackages() {
			for _, p := range getHelpers() {
				archtest.Package(t, p1).ShouldNotDependOn(p)
				assertNoError(t, mockT)
			}

			for _, p := range getServices() {
				archtest.Package(t, p1).ShouldNotDependOn(p)
				assertNoError(t, mockT)
			}
		}
	})

	t.Run("Services must be independent among ourselves", func(t *testing.T) {
		for _, p1 := range getServices() {
			for _, p2 := range getServices() {
				if p1 == p2 {
					continue
				}

				archtest.Package(t, p1).ShouldNotDependOn(p2)
			}
		}
	})

	t.Run("Services must not know about helpers", func(t *testing.T) {
		for _, h := range getHelpers() {
			for _, s := range getServices() {
				archtest.Package(t, s).ShouldNotDependOn(h)
			}
		}
	})
}

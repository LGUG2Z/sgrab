package sonarr_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSonarr(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sonarr Suite")
}

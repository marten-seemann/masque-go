package masque_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMasqueGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "MASQUE Suite")
}

package cli_test

import (
	. "code.cloudfoundry.org/uaa-cli/cli"

	"io/ioutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
)

var _ = Describe("JsonPrinter", func() {
	It("prints things to the Robots log", func() {
		logBuf := NewBuffer()
		printer := NewJsonPrinter(NewLogger(ioutil.Discard, logBuf, ioutil.Discard, ioutil.Discard))

		printer.Print(struct {
			Foo string
			Bar string
		}{"foo", "bar"})

		Expect(logBuf.Contents()).To(MatchJSON(`{"Foo":"foo","Bar":"bar"}`))
	})

	It("returns error when cannot marhsal into json", func() {
		printer := NewJsonPrinter(NewLogger(ioutil.Discard, ioutil.Discard, ioutil.Discard, ioutil.Discard))

		unJsonifiableObj := make(chan bool)
		err := printer.Print(unJsonifiableObj)

		Expect(err).To(HaveOccurred())
	})
})

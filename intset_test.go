package intset_test

import (
	"testing"

	gc "gopkg.in/check.v1"

	"github.com/wolf/intset"
)

func Test(t *testing.T) { gc.TestingT(t) }

type Suite struct{}

var _ = gc.Suite(&Suite{})

func (*Suite) SetUpTest(c *gc.C) {
}

func (*Suite) TearDownTest(c *gc.C) {
}

func (*Suite) TestAddInRange(c *gc.C) {
	s := intset.New(5, 1, 2, 4)
	c.Assert(s.Contains(3), gc.Equals, false)
	s.Add(3)
	c.Assert(s.Contains(3), gc.Equals, true)
}

func (*Suite) TestAddOutOfRange(c *gc.C) {
	c.Fatal("TODO")
}

func (*Suite) TestRemove(c *gc.C) {
	c.Fatal("TODO")
}

func (*Suite) TestUnion(c *gc.C) {
	c.Fatal("TODO")
}

func (*Suite) TestDifference(c *gc.C) {
	c.Fatal("TODO")
}

func (*Suite) TestChoose(c *gc.C) {
	c.Fatal("TODO: BONUS: test that Choose() provides a good source of randomness")
}

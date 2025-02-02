// +build smoke

package smoke

import (
	"os"
	"testing"

	"github.com/carolynvs/magex/shx"
	"github.com/stretchr/testify/require"
)

func TestHelloBundle(t *testing.T) {
	test, err := NewTest(t)
	defer test.Teardown()
	require.NoError(t, err, "test setup failed")

	// Build an interesting test bundle
	ref := "localhost:5000/mybuns:v0.1.1"
	shx.Copy("../testdata/mybuns", ".", shx.CopyRecursive)
	os.Chdir("mybuns")
	test.RequirePorter("build")
	test.RequirePorter("publish", "--reference", ref)

	// Do not run these commands in a bundle directory
	os.Chdir(test.TestDir)

	test.RequirePorter("install", "--reference", ref)
	test.RequirePorter("installation", "show", "mybuns")

	test.RequirePorter("upgrade", "mybuns")
	test.RequirePorter("installation", "show", "mybuns")

	test.RequirePorter("uninstall", "mybuns")
	test.RequirePorter("installation", "show", "mybuns")
}

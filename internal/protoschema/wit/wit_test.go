package wit

import (
	"strings"
	"testing"

	"github.com/straw-hat-team/bufbuild-plugins/internal/protoschema/golden"
	"github.com/stretchr/testify/require"
)

func TestWitGolden(t *testing.T) {
	t.Parallel()
	//dirPath := filepath.FromSlash("../../testdata/wit")
	testDescs, err := golden.GetTestDescriptors("../../testdata")
	require.NoError(t, err)
	for _, testDesc := range testDescs {
		output := strings.Builder{}

		for _, entry := range Generate(testDesc) {
			require.NoError(t, err)

			output.WriteString(entry.Content.String())

			//identifier := entry.ID
			//require.NotEmpty(t, identifier)
			//
			//filePath := filepath.Join(dirPath, identifier)
			//content := entry.Content.String() + "\n"
			//
			//err = golden.CheckGolden(filePath, content)
			//require.NoError(t, err)
		}

		require.NotEmpty(t, output.String())
	}
}

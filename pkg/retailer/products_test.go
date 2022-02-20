package retailer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter_HasFilterCriteria(t *testing.T) {
	t.Run("return false if no search or retailer criteria is specified", func(t *testing.T) {
		f := Filter{}
		assert.Equal(t, false, f.HasFilterCriteria())
	})

	t.Run("return true if either search or retailer criteria is specified", func(t *testing.T) {
		f := Filter{Search: "foo"}
		assert.Equal(t, true, f.HasFilterCriteria())

		f = Filter{Retailer: "bar"}
		assert.Equal(t, true, f.HasFilterCriteria())
	})
}

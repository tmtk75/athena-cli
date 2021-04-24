package main

import (
	"testing"

	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
)

func TestMergeConfigMap(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		v := viper.New()
		v.MergeConfigMap(map[string]interface{}{"a": 1})
		assert.Equal(t, 1, v.GetInt("a"))
	})
	t.Run("MergeConfigMap", func(t *testing.T) {
		v := viper.New()
		v.MergeConfigMap(map[string]interface{}{"a": 1})
		// override
		v.MergeConfigMap(map[string]interface{}{"a": 2})
		assert.Equal(t, 2, v.GetInt("a"))
	})
}

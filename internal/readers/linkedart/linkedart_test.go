package linkedart

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestModels(t *testing.T) {
	// t.Run("test GetObject", func(t *testing.T) {
	// 	lar := New()
	//
	// 	mObj, err := linkedlinkedart.GetObject("")
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, mObj)
	// 	assert.True(t, len(mObj.ID) > 0)
	// })
	//
	// t.Run("test GetNextObject", func(t *testing.T) {
	// 	lar := New()
	//
	// 	mObj, err := GetNextObject(nil)
	// 	assert.Nil(t, err)
	// 	assert.NotNil(t, mObj)
	// 	assert.True(t, len(mObj.ID) > 0)
	// })
}

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBuildStatus(t *testing.T) {

	t.Run("Succeeded", func(t *testing.T) {

		// act
		bitbucketAPIClient := newBitbucketAPIClient()
		err := bitbucketAPIClient.SetBuildStatus("KgSCeb8aZpHrUhDAVk5p8D8PA_08MbBpE3SNzFLSaTSd7uflrgOtelObPD6BRpd_AcdrvIvbu23xcWPfuws=", "xivart/inrule-catalogs-live", "3993d790d0a19fc5da450d61dc4de6473be1c440", "succeeded", "10.0.0", "", "")

		assert.Nil(t, err)
	})

	// t.Run("Failed", func(t *testing.T) {

	// 	// act
	// 	bitbucketAPIClient := newBitbucketAPIClient()
	// 	err := bitbucketAPIClient.SetBuildStatus("KgSCeb8aZpHrUhDAVk5p8D8PA_08MbBpE3SNzFLSaTSd7uflrgOtelObPD6BRpd_AcdrvIvbu23xcWPfuws=", "xivart/inrule-catalogs-live", "3993d790d0a19fc5da450d61dc4de6473be1c440", "failed")

	// 	assert.Nil(t, err)
	// })

	// t.Run("Pending", func(t *testing.T) {

	// 	// act
	// 	bitbucketAPIClient := newBitbucketAPIClient()
	// 	err := bitbucketAPIClient.SetBuildStatus("KgSCeb8aZpHrUhDAVk5p8D8PA_08MbBpE3SNzFLSaTSd7uflrgOtelObPD6BRpd_AcdrvIvbu23xcWPfuws=", "xivart/inrule-catalogs-live", "3993d790d0a19fc5da450d61dc4de6473be1c440", "pending")

	// 	assert.Nil(t, err)
	// })
}

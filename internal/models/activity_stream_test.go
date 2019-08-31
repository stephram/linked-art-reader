package models

import (
	"encoding/json"
	"linkedart-reader-golang/internal/utils"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const activityStreamInfoPayload = `{
  "@context": "https://www.w3.org/ns/activitystreams",
  "summary": "The Getty MART Repository's Recent Activity",
  "type": "OrderedCollection",
  "id": "https://mart.getty.edu/activity-stream",
  "startIndex": 1,
  "totalItems": 50291,
  "totalPages": 503,
  "maxPerPage": 100,
  "first": {
    "id": "https://mart.getty.edu/activity-stream/page/1",
    "type": "OrderedCollectionPage"
  },
  "last": {
    "id": "https://mart.getty.edu/activity-stream/page/503",
    "type": "OrderedCollectionPage"
  }
}
`
const activityStreamPagePayload = `{
  "@context": "https://www.w3.org/ns/activitystreams",
  "summary": "The Getty MART Repository's Recent Activity",
  "type": "OrderedCollection",
  "id": "https://mart.getty.edu/activity-stream/page/5",
  "startIndex": 1,
  "totalItems": 50291,
  "totalPages": 503,
  "maxPerPage": 100,
  "partOf": {
    "id": "https://mart.getty.edu/activity-stream",
    "type": "OrderedCollection"
  },
  "first": {
    "id": "https://mart.getty.edu/activity-stream/page/1",
    "type": "OrderedCollectionPage"
  },
  "last": {
    "id": "https://mart.getty.edu/activity-stream/page/503",
    "type": "OrderedCollectionPage"
  },
  "previous": {
    "id": "https://mart.getty.edu/activity-stream/page/4",
    "type": "OrderedCollectionPage"
  },
  "next": {
    "id": "https://mart.getty.edu/activity-stream/page/6",
    "type": "OrderedCollectionPage"
  },
  "orderedItems": [
    {
      "id": "https://mart.getty.edu/activity-stream/972c5636-61ed-4cfc-9c08-7a033b762db7",
      "type": "Create",
      "actor": null,
      "object": {
        "id": "https://mart.getty.edu/museum/collection/person/82a5847b-4747-456f-a5d9-88bd21e4c08a",
        "type": "Person"
      },
      "created": "2019-08-04 22:01:00+0000",
      "updated": "2019-08-04 22:01:00+0000",
      "published": "2019-08-04 22:01:00+0000"
    },
    {
      "id": "https://mart.getty.edu/activity-stream/e45cd01e-3367-4477-87be-ba2eed853ba5",
      "type": "Create",
      "actor": null,
      "object": {
        "id": "https://mart.getty.edu/museum/collection/person/40c1b584-0cde-4557-90df-a7489669ed04",
        "type": "Person"
      },
      "created": "2019-08-04 22:01:00+0000",
      "updated": "2019-08-04 22:01:00+0000",
      "published": "2019-08-04 22:02:00+0000"
    },
    {
      "id": "https://mart.getty.edu/activity-stream/83985935-ea0c-48a9-bbb5-2de709357bab",
      "type": "Create",
      "actor": null,
      "object": {
        "id": "https://mart.getty.edu/museum/collection/group/eaeaa9fc-45f4-4203-bc93-cd605266ca45",
        "type": "Group"
      },
      "created": "2019-08-04 22:01:31+0000",
      "updated": "2019-08-04 22:01:31+0000",
      "published": "2019-08-04 22:03:31+0000"
    }  
  ]
}`

const (
	activityStreamSummary = "The Getty MART Repository's Recent Activity"
	activityStreamType    = "OrderedCollection"
	activityStreamID      = "https://mart.getty.edu/activity-stream"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestActivityStreamParsing(t *testing.T) {
	t.Run("test activityStreamInfo unmarshal success", func(t *testing.T) {
		var activityStreamInfo OrderedCollection

		if err := json.Unmarshal([]byte(activityStreamInfoPayload), &activityStreamInfo); err != nil {
			assert.Fail(t, "failed to unmarshal activityStreamInfoPayload", err)
		}

		assert.Equal(t, activityStreamSummary, activityStreamInfo.Summary)
		assert.Equal(t, activityStreamID, activityStreamInfo.ID)
		assert.Equal(t, activityStreamType, activityStreamInfo.DataType)
		assert.Equal(t, "https://mart.getty.edu/activity-stream", activityStreamInfo.ID)
		assert.Equal(t, 1, activityStreamInfo.StartIndex)
		assert.Equal(t, 50291, activityStreamInfo.TotalItems)
		assert.Equal(t, 503, activityStreamInfo.TotalPages)
		assert.Equal(t, 100, activityStreamInfo.MaxPerPage)
	})

	t.Run("test activityStreamPage unmarshal success", func(t *testing.T) {
		var activityStreamPage ActivityStreamPage

		if err := json.Unmarshal([]byte(activityStreamPagePayload), &activityStreamPage); err != nil {
			assert.Fail(t, "failed to unmarshal activityStreamPagePayload", err.Error())
		}

		assert.Equal(t, 3, len(activityStreamPage.OrderedItems))

		for i := 0; i < len(activityStreamPage.OrderedItems); i++ {
			item := activityStreamPage.OrderedItems[i]

			assert.NotNil(t, item.Object)

			assert.NotNil(t, item.Object.ID)
			assert.True(t, len(item.Object.ID) > 0)

			assert.NotNil(t, item.Object.ObjType)
			assert.True(t, len(item.Object.ObjType) > 0)

			ctm := utils.ConvertPythonDateToTime(item.Created)
			log.Infof("created: %s", ctm.String())
			assert.False(t, ctm.IsZero())

			utm := utils.ConvertPythonDateToTime(item.Updated)
			log.Infof("updated: %s", utm.String())
			assert.False(t, utm.IsZero())

			ptm := utils.ConvertPythonDateToTime(item.Published)
			log.Infof("published: %s", ptm.String())
			assert.False(t, ptm.IsZero())
		}
	})
}

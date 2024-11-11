package domain_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"backend_academy_2024_project_3-go-Dabzelos/internal/domain"
)

func TestDataHolder_Parser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		testScenario string
		logs         []string
		parsedData   domain.DataHolder
		to           string
		from         string
	}{
		{
			testScenario: "normal logs",
			logs: []string{
				"80.91.33.133 - - [17/May/2015:08:05:24 +0000] \"GET /downloads/product_1 HTTP/1.1\" " +
					"304 0 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.17)\"",
				"93.180.71.3 - - [17/May/2015:08:05:23 +0000] \"GET /downloads/product_1 HTTP/1.1\" " +
					"304 0 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.21)\"",
				"80.91.33.133 - - [17/May/2015:08:05:24 +0000] \"GET /downloads/product_1 HTTP/1.1\" " +
					"304 0 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.17)\"",
			},
			parsedData: domain.DataHolder{
				TotalCounter: 3,
				UnparsedLogs: 0,
				BytesSend:    []int{0, 0, 0},
				RequestedResources: map[string]int{
					"/downloads/product_1": 3,
				},
				CommonAnswers: map[string]int{
					"304": 3,
				},
			},
		},
		{
			testScenario: "corrupted logs",
			logs: []string{
				"80.91.33.133 - - [17/May/2015:08:05:24 +0000] \"/downloads/product_1 HTTP/1.1\" " +
					"304 0 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.17)\"",
				"93.180.71.3 - - [May/2015:08:05:23 +0000] \"GET /downloads/product_1 HTTP/1.1\" " +
					"304 0 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.21)\"",
				"80.91.33.133 - - [17/May/2015:08:05:24 +0000] \"GET /downloads/product_1 HTTP/1.1\" " +
					" \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.17)\"",
			},
			parsedData: domain.DataHolder{
				TotalCounter:       0,
				UnparsedLogs:       3,
				BytesSend:          nil,
				RequestedResources: map[string]int{},
				CommonAnswers:      map[string]int{},
			},
		},
		{
			testScenario: "Time bounds",
			logs: []string{
				"80.91.33.133 - - [17/May/2015:08:15:24 +0000] \"GET /downloads/product_1 HTTP/1.1\"" +
					" 304 0 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.17)\"",
				"93.180.71.3 - - [17/May/2015:08:09:23 +0000] \"GET /downloads/product_1 HTTP/1.1\"" +
					" 304 0 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.21)\"",
				"80.91.33.133 - - [17/May/2015:08:08:24 +0000] \"GET /downloads/product_1 HTTP/1.1\"" +
					" 304 0 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.17)\"",
				"217.168.17.5 - - [17/May/2015:08:07:34 +0000] \"GET /downloads/product_1 HTTP/1.1\"" +
					" 200 490 \"-\" \"Debian APT-HTTP/1.3 (0.8.10.3)\"",
				"217.168.17.5 - - [17/May/2015:08:06:09 +0000] \"GET /downloads/product_2 HTTP/1.1\"" +
					" 200 490 \"-\" \"Debian APT-HTTP/1.3 (0.8.10.3)\"",
				"93.180.71.3 - - [17/May/2015:08:05:57 +0000] \"GET /downloads/product_1 HTTP/1.1\"" +
					" 304 0 \"-\" \"Debian APT-HTTP/1.3 (0.8.16~exp12ubuntu10.21)\"",
			},
			parsedData: domain.DataHolder{
				TotalCounter:       4,
				UnparsedLogs:       0,
				BytesSend:          []int{0, 490, 490, 0},
				RequestedResources: map[string]int{"/downloads/product_1": 3, "/downloads/product_2": 1},
				CommonAnswers:      map[string]int{"304": 2, "200": 2},
			},
			to:   "17/May/2015:08:08:24 +0000",
			from: "17/May/2015:08:05:57 +0000",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testScenario, func(t *testing.T) {
			data := domain.NewDataHolder("", "")

			for _, log := range tc.logs {
				to, _ := time.Parse("02/Jan/2006:15:04:05 -0700", tc.to)
				from, _ := time.Parse("02/Jan/2006:15:04:05 -0700", tc.from)
				data.Parser(log, from, to)
			}

			assert.Equal(t, data.TotalCounter, tc.parsedData.TotalCounter)
			assert.Equal(t, data.UnparsedLogs, tc.parsedData.UnparsedLogs)
			assert.Equal(t, data.BytesSend, tc.parsedData.BytesSend)
			assert.Equal(t, data.RequestedResources, tc.parsedData.RequestedResources)
			assert.Equal(t, data.CommonAnswers, tc.parsedData.CommonAnswers)
		})
	}
}

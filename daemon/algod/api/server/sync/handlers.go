// Copyright (C) 2019-2022 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package sync

import (
	"encoding/json"
	"fmt"
	"github.com/algorand/go-algorand/daemon/algod/api/server/lib"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
)

// Allowed Database names
var databaseNames = map[string]bool{
	"accounts":     true,
	"applications": true,
	"assets":       true,
	"blocks":       true,
	"transactions": true,
}

// handleDatabaseParam check for existence of a database
// used in high level handlers/catchalls
func handleDatabaseParam(c echo.Context) (string, error) {
	dbName := c.Param("db")
	// Check for valid database names
	if !databaseNames[dbName] {
		err := notfound(c, "database not found!")
		return dbName, err
	}

	// Return valid database
	return dbName, nil
}

type OkResponse struct {
	Ok interface{} `json:"ok"`
}

// GetBlock handles the GET /blocks/:roundNumber request
func GetBlock(ctx lib.ReqContext, c echo.Context, openRevs bool) error {
	roundNumber, err := strconv.Atoi(c.Param("docId"))
	if err != nil {
		return err
	}
	var rp struct {
		Revs  bool         `query:"revs"`
		Round basics.Round `param:"roundNumber"`
	}
	err = c.Bind(&rp)
	if err != nil {
		return err
	}
	// Check for response type
	if rp.Revs {
		// Return revisions for changes feed
		doc, err := BlockDocWithRev(ctx, basics.Round(roundNumber))
		if err != nil {
			return err
		}

		c.Response().Header().Set("ETag", "\""+doc.Rev+"\"")
		if openRevs {
			return json.NewEncoder(c.Response()).Encode([]OkResponse{{
				Ok: doc,
			}})
		} else {
			return json.NewEncoder(c.Response()).Encode(doc)
		}

	} else {
		// Return regular document
		doc, err := BlockDoc(ctx, basics.Round(roundNumber))
		if err != nil {
			return err
		}
		c.Response().Header().Set("ETag", "\""+doc.Rev+"\"")
		return json.NewEncoder(c.Response()).Encode(doc)
	}
}

// GetDocument is the unhandled GET /:db/:docId requests
func GetDocument(ctx lib.ReqContext, c echo.Context) error {
	url := c.Request().URL
	println("GetDocument:" + url.Path + "?" + url.RawQuery)
	dbName, err := handleDatabaseParam(c)
	if err != nil {
		return err
	}
	openRevs := c.QueryParam("open_revs")
	if openRevs != "" {
		switch dbName {
		case "blocks":
			return GetBlock(ctx, c, true)
		default:
			return c.JSON(http.StatusNotImplemented, DatabaseError{
				Error:  "not_implemented",
				Reason: "Function not implemented",
			})
		}
	} else {
		switch dbName {
		case "blocks":
			return GetBlock(ctx, c, false)
		default:
			return c.JSON(http.StatusNotImplemented, DatabaseError{
				Error:  "not_implemented",
				Reason: "Function not implemented",
			})
		}
	}

}

// generateChangesets creates a list of changes since a specified block to a limit length
func generateChangesets(ctx lib.ReqContext, since int, limit int) ([]Changeset, error) {
	var result []Changeset
	for i := since; i < since+limit; i++ {
		doc, err := BlockDocWithRev(ctx, basics.Round(i))
		if err != nil {
			return nil, err
		}
		changeSet := Changeset{
			// The id of the document is the round
			Id: strconv.Itoa(i),
			// Changes are immutable, always the first revision
			Changes: []ChangesetRevision{
				{
					Rev: doc.Rev,
				},
			},
			// Sequence is equal to the round
			Seq: i + 1,
		}
		result = append(result, changeSet)
	}
	return result, nil
}

// TODO: Get continuous changes optimized
func continuousChanges(ctx lib.ReqContext, c echo.Context) error {
	fmt.Println("continuousChanges")
	//ledger := ctx.Node.Ledger()
	var responseTime time.Duration
	var lastRound basics.Round
	var limit int
	var since int

	sinceQuery := c.QueryParam("since")
	if sinceQuery != "" {
		sinceInt, err := strconv.Atoi(sinceQuery)
		if err != nil {
			return c.JSON(http.StatusBadRequest, DatabaseError{
				Error:  "unknown_error",
				Reason: "bad argument",
			})
		}
		since = sinceInt
		// TODO: add in "since"
		lastRound = basics.Round(since)
	} else {
		lastRound = basics.Round(1)
	}

	timeoutQuery := c.QueryParam("timeout")
	if timeoutQuery != "" {
		timeout, err := strconv.Atoi(timeoutQuery)
		if err != nil {
			return c.JSON(http.StatusBadRequest, DatabaseError{
				Error:  "unknown_error",
				Reason: "bad argument",
			})
		}
		responseTime = time.Duration(timeout) * time.Millisecond
	} else {
		responseTime = 10000 * time.Millisecond
	}

	limitQuery := c.QueryParam("limit")
	if limitQuery != "" {
		limitQueryInt, err := strconv.Atoi(limitQuery)
		if err != nil {
			return c.JSON(http.StatusBadRequest, DatabaseError{
				Error:  "unknown_error",
				Reason: "bad argument",
			})
		}
		limit = limitQueryInt
	}
	//
	//nextRound := ledger.NextRound()
	//println(responseTime)
	//responseTimeout := time.After(responseTime)
	////nextRoundReady := ledger.Wait(nextRound)
	//enc := json.NewEncoder(c.Response())

	changeSets, err := generateChangesets(ctx, since, limit)
	if err != nil {
		return err
	}

	for i, change := range changeSets {
		fmt.Println(i, change, responseTime, lastRound)
	}
	//for lastRound < nextRound {
	//for i := 1; i < limit+1; i++ {
	//	lastRound = basics.Round(i)
	//	doc, err := blockDocWithRev(ctx, lastRound, true)
	//	if err != nil {
	//		return err
	//	}
	//
	//	err = enc.Encode(Changeset{
	//		// Sequence is equal to the round
	//		Seq: strconv.Itoa(int(lastRound)),
	//		// The id of the document is the round
	//		ID: strconv.Itoa(int(lastRound)),
	//		// Changes are immutable, always the first revision
	//		Changes: []ChangesetRevision{
	//			{
	//				Rev: doc.Rev,
	//			},
	//		},
	//	})
	//
	//	if err != nil {
	//		return err
	//	}
	//	//c.Response().Flush()
	//	// Wait for next block or abort
	//	select {
	//	//case <-ctx.Shutdown:
	//	//	println("Shutting down")
	//	//	return c.JSON(http.StatusInternalServerError, DatabaseError{
	//	//		Error:  "server_error",
	//	//		Reason: "Server is shutting down.",
	//	//	})
	//	case <-ledger.Wait(nextRound):
	//	////	println("Made it to next round: " + strconv.Itoa(int(lastRound)) + " " + strconv.Itoa(int(nextRound)))
	//	////	nextRound += 1
	//	case <-responseTimeout:
	//		println("Timeout of " + responseTime.String())
	//		err := enc.Encode(struct {
	//			LastSeq uint64 `json:"last_seq"`
	//			Pending uint64 `json:"pending"`
	//		}{
	//			LastSeq: uint64(lastRound),
	//			//Pending: uint64(nextRound - lastRound),
	//			Pending: 0,
	//		})
	//		if err != nil {
	//			return err
	//		}
	//		c.Response().Flush()
	//		return nil
	//		//return c.JSON(http.StatusOK, struct {
	//		//	LastSeq uint64 `json:"last_seq"`
	//		//	Pending uint64 `json:"pending"`
	//		//}{
	//		//	LastSeq: uint64(lastRound),
	//		//	//Pending: uint64(nextRound - lastRound),
	//		//	Pending: 0,
	//		//})
	//	}
	//}
	//err := enc.Encode(struct {
	//	LastSeq uint64 `json:"last_seq"`
	//	Pending uint64 `json:"pending"`
	//}{
	//	LastSeq: uint64(lastRound),
	//	//Pending: uint64(nextRound - lastRound),
	//	Pending: 0,
	//})
	//if err != nil {
	//	return err
	//}
	//c.Response().Flush()
	return nil
}

// normalChanges
func normalChanges(ctx lib.ReqContext, c echo.Context, p ChangesParams) error {
	limit := p.Limit
	since := p.Since
	lastRound := ctx.Node.Ledger().LastRound()
	lastSeq := limit + since
	max := int(lastRound)
	if max < lastSeq {
		lastSeq = max
		limit = max - since
	}
	changeSets, err := generateChangesets(ctx, since, limit)
	if err != nil {
		return err
	}

	changes := Changes{
		Results: changeSets,
		LastSeq: lastSeq,
		//Pending: 0,
	}

	return c.JSON(http.StatusOK, changes)
}

type ChangesParams struct {
	Limit int
	Since int
	Feed  string
}

// GetDatabaseChanges is the handler for /:db/_changes
// it follows the CouchDb replication protocol for easy data syncing
// See https://docs.couchdb.org/en/3.2.2-docs/api/database/changes.html
func GetDatabaseChanges(ctx lib.ReqContext, c echo.Context) error {
	url := c.Request().URL
	println("GetDatabaseChanges:" + url.Path + "?" + url.RawQuery)
	_, err := handleDatabaseParam(c)
	if err != nil {
		return err
	}

	//var responseTime time.Duration
	var limit int
	var since int

	// Types of change feeds
	var feedTypes = map[string]bool{
		//normal Specifies Normal Polling Mode. All past changes are returned immediately. Default.
		"normal": true,
		//longpoll Specifies Long Polling Mode. Waits until at least one change has occurred, sends the change,
		//then closes the connection. Most commonly used in conjunction with since=now, to wait for the next change.
		"longpoll": true,
		//Sets Continuous Mode. Sends a line of JSON per event. Keeps the socket open until timeout.
		"continuous": true,
		//Sets Event Source Mode. Works the same as Continuous Mode, but sends the events in EventSource format.
		"eventsource ": true,
	}
	feed := c.QueryParam("feed")
	// Check for valid database names
	if feed != "" && !feedTypes[feed] {
		return c.JSON(http.StatusBadRequest, DatabaseError{
			Error:  "bad_request",
			Reason: "invalid feed type",
		})
	} else {
		feed = "normal"
	}

	sinceQuery := c.QueryParam("since")
	if sinceQuery != "" {
		sinceInt, err := strconv.Atoi(sinceQuery)
		if err != nil {
			return c.JSON(http.StatusBadRequest, DatabaseError{
				Error:  "unknown_error",
				Reason: "bad argument",
			})
		}
		since = sinceInt
	} else {
		since = 0
	}

	//timeoutQuery := c.QueryParam("timeout")
	//if timeoutQuery != "" {
	//	timeout, err := strconv.Atoi(timeoutQuery)
	//	if err != nil {
	//		return c.JSON(http.StatusBadRequest, DatabaseError{
	//			Error:  "unknown_error",
	//			Reason: "bad argument",
	//		})
	//	}
	//	responseTime = time.Duration(timeout) * time.Millisecond
	//} else {
	//	responseTime = 10000 * time.Millisecond
	//}

	limitQuery := c.QueryParam("limit")
	if limitQuery != "" {
		limitQueryInt, err := strconv.Atoi(limitQuery)
		if err != nil {
			return c.JSON(http.StatusBadRequest, DatabaseError{
				Error:  "unknown_error",
				Reason: "bad argument",
			})
		}
		limit = limitQueryInt
	} else {
		limit = int(ctx.Node.Ledger().LastRound())
	}

	changesParams := ChangesParams{
		Limit: limit,
		Since: since,
		Feed:  feed,
	}

	switch feed {
	case "continuous":
		return continuousChanges(ctx, c)
	default:
		return normalChanges(ctx, c, changesParams)
	}
}

// GetDatabaseInfo is an httpHandler for route /{db}
// See https://docs.couchdb.org/en/3.2.2-docs/api/database/common.html
func GetDatabaseInfo(ctx lib.ReqContext, c echo.Context) error {
	url := c.Request().URL
	println("GetDatabaseInfo:" + url.Path + "?" + url.RawQuery)

	dbName, err := handleDatabaseParam(c)
	if err != nil {
		return err
	}

	// Create default database response
	database := Database{
		DatabaseName:   dbName,
		CompactRunning: false,
		DocCount:       0,
		// Blockchain is immutable 😊
		DocDelCount:       0,
		DiskFormatVersion: 0,
		// Must be "0" for legacy reasons
		InstanceStartTime: "0",
		PurgeSeq:          0,
		UpdateSeq:         0,
	}

	switch dbName {
	case "blocks":
		// Fetch node status
		stat, err := ctx.Node.Status()
		// Handle server errors
		if err != nil {
			return c.JSON(http.StatusInternalServerError, DatabaseError{
				Error:  "server_error",
				Reason: "Failed to fetch Node Status",
			})
		}
		database.DocCount = int(stat.LastRound)
		database.UpdateSeq = int(stat.LastRound)
	default:
		return c.JSON(http.StatusNotImplemented, DatabaseError{
			Error:  "not_implemented",
			Reason: "Function not implemented",
		})
	}

	// Respond with database info
	return c.JSON(http.StatusOK, database)
}

func getAllDocs(ctx lib.ReqContext, dbName string, keys []string, includeDoc bool) ([]AllDocsDocInfo, error) {
	var res []AllDocsDocInfo
	for _, keyStr := range keys {
		var doc Document
		switch dbName {
		case "blocks":
			key, err := strconv.Atoi(keyStr)
			if err != nil {
				return res, err
			}
			doc, err = BlockDoc(ctx, basics.Round(key))
			if err != nil {
				return res, err
			}
		}

		info := AllDocsDocInfo{
			Id:  keyStr,
			Key: keyStr,
			Value: struct {
				Rev string `json:"rev"`
			}{
				Rev: doc.Rev,
			},
		}
		if includeDoc {
			info.Doc = doc
		}

		res = append(res, info)
	}
	return res, nil
}

func GetAllDocs(ctx lib.ReqContext, c echo.Context) error {
	return c.JSON(http.StatusOK, AllDocsResponse{
		TotalRows: int(ctx.Node.Ledger().LastRound()),
		Offset:    0,
		Rows:      nil,
	})
}

func PostAllDocs(ctx lib.ReqContext, c echo.Context) error {
	url := c.Request().URL
	println("PostAllDocs:" + url.Path + "?" + url.RawQuery)
	dbName, err := handleDatabaseParam(c)
	var postKeys AllDocsPostRequest
	err = c.Bind(&postKeys)
	includeDocsQuery := c.QueryParam("include_docs")
	includeDocs, err := strconv.ParseBool(includeDocsQuery)
	if err != nil {
		return c.JSON(http.StatusBadRequest, DatabaseError{
			Error:  "bad_request",
			Reason: "Bad arguments",
		})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, DatabaseError{
			Error:  "bad_request",
			Reason: "Invalid post body",
		})
	}

	allDocs, err := getAllDocs(ctx, dbName, postKeys.Keys, includeDocs)
	if err != nil {
		return c.JSON(http.StatusBadRequest, DatabaseError{
			Error:  "bad_request",
			Reason: "Failed getting docs",
		})
	}
	return c.JSON(http.StatusOK, AllDocsResponse{
		TotalRows: int(ctx.Node.Ledger().LastRound()),
		Offset:    0,
		Rows:      allDocs,
	})
}

// GetNodeInfo is an httpHandler for route /
// see https://docs.couchdb.org/en/3.2.2-docs/api/server/common.html
func GetNodeInfo(ctx lib.ReqContext, context echo.Context) error {
	conf := ctx.Node.Config()
	// TODO: get fingerprint
	//id := uuid.New()
	var features []string

	if conf.Archival {
		features = append(features, "Archival")
	}
	return context.JSON(http.StatusOK, Status{
		Algod:       "Welcome",
		Version:     strconv.Itoa(int(conf.Version)),
		Network:     ctx.Node.GenesisID(),
		GenesisHash: ctx.Node.GenesisHash().String(),
		Uuid:        "d49d0c78-4e75-4e48-80b7-c6143bf805b2",
		Features:    features,
		Vendor: Vendor{
			Name:    "Algorand, Inc.",
			Version: strconv.Itoa(int(conf.Version)),
		},
	})
}

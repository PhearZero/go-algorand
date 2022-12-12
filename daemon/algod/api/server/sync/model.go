package sync

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"github.com/algorand/go-algorand/daemon/algod/api/server/lib"
	"github.com/algorand/go-algorand/daemon/algod/api/server/v1/handlers"
	"github.com/algorand/go-algorand/data/basics"
	"strconv"
)

// generateRevHash creates a MD5 hash of a struct
// this is used as the ID of a Document.Rev (ex {_rev: 1-3046b813ec5e5dcf86a93f4921b49734})
func generateRevHash(doc Document) (string, error) {
	jsonStr, err := json.Marshal(doc)

	if err != nil {
		return "", err
	}
	res := md5.Sum(jsonStr)

	return hex.EncodeToString(res[:]), err
}

// BlockDoc builds compiles a document from a basics.Round
func BlockDoc(ctx lib.ReqContext, blockRound basics.Round) (Document, error) {
	var doc Document
	block, c, err := ctx.Node.Ledger().BlockCert(blockRound)
	if err != nil {
		return doc, err
	}
	// TODO: Change encoding to v2
	v1, err := handlers.BlockEncodeV1(block, c)

	// Create the Block Document
	doc = Document{
		Id:   strconv.Itoa(int(blockRound)),
		Type: "Block",
		Data: v1,
	}

	// Create MD5 hash of document
	// TODO: Cache these somewhere
	revId, err := generateRevHash(doc)
	if err != nil {
		return doc, err
	}
	doc.MD5 = revId
	doc.Rev = "1-" + doc.MD5

	return doc, nil
}

// BlockDocWithRev fuck you GoLang and your conditional hell
func BlockDocWithRev(ctx lib.ReqContext, blockRound basics.Round) (DocumentWithRevisionTree, error) {
	var result DocumentWithRevisionTree

	// Create the Block Document
	doc, err := BlockDoc(ctx, blockRound)
	if err != nil {
		return result, err
	}

	// Create a Revision Tree
	// blockchain is immutable and there will only be one revision
	rt := RevisionTree{
		Start: 1,
		Ids:   []string{doc.MD5},
	}

	// Create the Document with revisions
	result = DocumentWithRevisionTree{
		Id:        doc.Id,
		Rev:       doc.Rev,
		Revisions: rt,
		Data:      doc.Data,
		Type:      "Block",
		MD5:       doc.MD5,
	}

	return result, nil
}

// Cluster contains Database cluster configuration
// swagger:model Cluster
type Cluster struct {
	// Replicas. The number of copies of every document.
	Replicas uint64 `json:"n"`
	// Shards. The number of range partitions.
	Shards uint64 `json:"q"`
	// Read quorum. The number of consistent copies of a document that need to be read before a successful reply.
	Read uint64 `json:"r"`
	// Write quorum. The number of copies of a document that need to be written before a successful reply.
	Write uint64 `json:"w"`
}

// Database contains the information about a collection of data
// swagger:model Database
type Database struct {
	// DatabaseName is the name of the database.
	//
	// required: true
	DatabaseName string `json:"db_name"`
	// CompactRunning set to true if the database compaction routine is operating on this database.
	//
	// required: true
	CompactRunning    bool   `json:"compact_running"`
	DiskFormatVersion int    `json:"disk_format_version"`
	DocCount          int    `json:"doc_count"`
	DocDelCount       int    `json:"doc_del_count"`
	InstanceStartTime string `json:"instance_start_time"`
	PurgeSeq          int    `json:"purge_seq"`
	//Sizes             struct {
	//	Active   int64 `json:"active"`
	//	Disk     int64 `json:"disk"`
	//	External int64 `json:"external"`
	//} `json:"sizes"`
	UpdateSeq int `json:"update_seq"`
}

// DatabaseError contains the default Couchdb error struct
type DatabaseError struct {
	// Error is the shorthand code for the error
	//
	// required: true
	Error string `json:"error"`

	// Reason explains the current error in natural language
	//
	// required: true
	Reason string `json:"reason"`
}

type Vendor struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Status of the node
type Status struct {
	Algod       string   `json:"algod"`
	Version     string   `json:"version"`
	Network     string   `json:"network"`
	GenesisHash string   `json:"genesis_hash"`
	Uuid        string   `json:"uuid"`
	Features    []string `json:"features,omitempty"`
	Vendor      Vendor   `json:"vendor"`
}

type ChangesetRevision struct {
	Rev string `json:"rev"`
}

type Changeset struct {
	Seq     int                 `json:"seq"`
	Id      string              `json:"id"`
	Changes []ChangesetRevision `json:"changes"`
}

type Changes struct {
	Results []Changeset `json:"results"`
	LastSeq int         `json:"last_seq"`
	//Pending int         `json:"pending"`
}

type RevisionTree struct {
	Start int      `json:"start"`
	Ids   []string `json:"ids"`
}

type DocumentWithRevisionTree struct {
	Id        string       `json:"_id"`
	Rev       string       `json:"_rev,omitempty"`
	Revisions RevisionTree `json:"_revisions,omitempty"`
	Type      string       `json:"type,omitempty"`
	Data      interface{}  `json:"data"`
	MD5       string       `json:"md5"`
}
type Document struct {
	Id   string      `json:"_id"`
	Rev  string      `json:"_rev,omitempty"`
	Type string      `json:"type,omitempty"`
	Data interface{} `json:"data,omitempty"`
	V1   interface{} `json:"v1,omitempty"`
	MD5  string      `json:"md5"`
}

type AllDocsPostRequest struct {
	Keys []string `json:"keys"`
}

type AllDocsDocInfo struct {
	Id    string `json:"id"`
	Key   string `json:"key"`
	Value struct {
		Rev string `json:"rev"`
	} `json:"value"`
	Doc interface{} `json:"doc,omitempty"`
}

type AllDocsResponse struct {
	TotalRows int              `json:"total_rows"`
	Offset    int              `json:"offset"`
	Rows      []AllDocsDocInfo `json:"rows"`
}

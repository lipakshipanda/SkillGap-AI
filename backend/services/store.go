package services

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"

	badger "github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"

	"skillgap-ai/models"
)

type Store struct {
	db *badger.DB
}

// OpenStore opens (or creates) the embedded BadgerDB database on disk.
// This is the direct replacement for connect_db() in mongo.py — no
// external database server required at all.
func OpenStore(path string) (*Store, error) {
	opts := badger.DefaultOptions(path).WithLoggingLevel(badger.WARNING)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger db: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

const keyPrefix = "analysis:"

func recordKey(createdAt time.Time, id string) []byte {
	// Sortable key: prefix + RFC3339Nano timestamp + id, so a prefix scan
	// naturally yields chronological order (we reverse it for "newest first").
	return []byte(fmt.Sprintf("%s%s:%s", keyPrefix, createdAt.UTC().Format(time.RFC3339Nano), id))
}

// SaveAnalysis persists a new analysis and returns its generated ID.
func (s *Store) SaveAnalysis(result models.AnalysisResult, resumeSnippet, targetRole string) (models.AnalysisRecord, error) {
	record := models.AnalysisRecord{
		ID:            uuid.NewString(),
		ResumeSnippet: resumeSnippet,
		TargetRole:    targetRole,
		Result:        result,
		CreatedAt:     time.Now().UTC(),
	}

	data, err := json.Marshal(record)
	if err != nil {
		return record, err
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(recordKey(record.CreatedAt, record.ID), data)
	})
	return record, err
}

// GetHistory returns the most recent `limit` analyses, newest first.
func (s *Store) GetHistory(limit int) ([]models.HistoryItem, error) {
	var records []models.AnalysisRecord

	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(keyPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var rec models.AnalysisRecord
				if err := json.Unmarshal(val, &rec); err != nil {
					return err
				}
				records = append(records, rec)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].CreatedAt.After(records[j].CreatedAt)
	})

	if limit > 0 && len(records) > limit {
		records = records[:limit]
	}

	items := make([]models.HistoryItem, 0, len(records))
	for _, r := range records {
		items = append(items, r.ToHistoryItem())
	}
	return items, nil
}

// findKeyByID does a prefix scan to locate the full key for a given analysis ID,
// since the primary key is timestamp-prefixed rather than the bare ID.
func (s *Store) findKeyByID(id string) ([]byte, *models.AnalysisRecord, error) {
	var foundKey []byte
	var foundRec *models.AnalysisRecord

	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(keyPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			var rec models.AnalysisRecord
			err := item.Value(func(val []byte) error {
				return json.Unmarshal(val, &rec)
			})
			if err != nil {
				return err
			}
			if rec.ID == id {
				k := item.KeyCopy(nil)
				foundKey = k
				foundRec = &rec
				return nil
			}
		}
		return nil
	})
	return foundKey, foundRec, err
}

// GetAnalysis fetches a single full analysis by ID.
func (s *Store) GetAnalysis(id string) (*models.AnalysisResponse, error) {
	_, rec, err := s.findKeyByID(id)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return nil, nil
	}
	resp := rec.ToResponse()
	return &resp, nil
}

// DeleteAnalysis removes an analysis by ID. Returns false if not found.
func (s *Store) DeleteAnalysis(id string) (bool, error) {
	key, rec, err := s.findKeyByID(id)
	if err != nil {
		return false, err
	}
	if rec == nil {
		return false, nil
	}
	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
	return true, err
}

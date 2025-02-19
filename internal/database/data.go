package database

import (
	"errors"
	"fmt"

	"github.com/AccumulateNetwork/accumulate/protocol"
	"github.com/AccumulateNetwork/accumulate/smt/storage"
)

// Data manages a data chain.
type Data struct {
	batch  *Batch
	record recordBucket
	chain  *Chain
}

// Height returns the number of entries.
func (d *Data) Height() int64 {
	return d.chain.Height()
}

// Put adds an entry to the chain.
func (d *Data) Put(hash []byte, entry *protocol.DataEntry) error {
	// Check that the entry does already not exist.
	_, err := d.batch.store.Get(d.record.Data(hash))
	if err == nil {
		return fmt.Errorf("data entry with hash %X exsits", hash)
	} else if !errors.Is(err, storage.ErrNotFound) {
		return err
	}

	// Write data entry
	data, err := entry.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %v", err)
	}
	d.batch.store.Put(d.record.Data(hash), data)

	// Add entry to the chain
	err = d.chain.AddEntry(hash)
	if err != nil {
		return err
	}

	// Write the anchor to the BPT
	var anchor [32]byte
	copy(anchor[:], d.chain.Anchor())
	d.batch.putBpt(d.record.Object().Append("Data"), anchor)

	return nil
}

// Get looks up an entry by it's hash.
func (d *Data) Get(hash []byte) (*protocol.DataEntry, error) {
	data, err := d.batch.store.Get(d.record.Data(hash))
	if err != nil {
		return nil, err
	}

	entry := new(protocol.DataEntry)
	err = entry.UnmarshalBinary(data)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// GetLatest looks up the latest entry.
func (d *Data) GetLatest() ([]byte, *protocol.DataEntry, error) {
	height := d.chain.Height()
	hash, err := d.chain.Entry(height - 1)
	if err != nil {
		return nil, nil, err
	}

	entry, err := d.Get(hash)
	if err != nil {
		return nil, nil, err
	}

	return hash, entry, nil
}

// GetHashes returns entry hashes in the given range
func (d *Data) GetHashes(start, end int64) ([][]byte, error) {
	return d.chain.Entries(start, end)
}

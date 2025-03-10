package rogue_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/charmbracelet/log"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/fivetentaylor/pointy/pkg/constants"
	"github.com/fivetentaylor/pointy/pkg/env"
	"github.com/fivetentaylor/pointy/pkg/rogue"
	"github.com/fivetentaylor/pointy/pkg/storage/s3"
	"github.com/fivetentaylor/pointy/pkg/testutils"
	v3 "github.com/fivetentaylor/pointy/rogue/v3"
)

func TestMain(m *testing.M) {
	if err := godotenv.Load("../../.env.test"); err != nil {
		log.Fatalf("Error loading .env.test file")
	}

	s3, err := s3.NewS3()
	if err != nil {
		log.Fatalf("failed to create db: %s", err)
	}
	err = s3.CreateTestBucket()
	if err != nil {
		log.Infof("failed to create bucket: %s", err)
	}

	// Run all tests
	code := m.Run()
	log.Info("stopping storage tests", "code", code)

	os.Exit(code)
}

func TestListFirst(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()

	qry := env.Query(ctx)
	s3 := env.S3(ctx)
	rds := env.Redis(ctx)

	ds := rogue.NewDocStore(s3, qry, rds)

	objBytes := []byte("hello world, I am an s3 object")
	err := s3.PutObject(s3.Bucket, "events/0000", "text/plain", objBytes)
	assert.NoError(t, err)

	err = s3.PutObject(s3.Bucket, "events/0001", "text/plain", objBytes)
	assert.NoError(t, err)

	keys, err := s3.List(s3.Bucket, "events", -1, -1)
	assert.NoError(t, err)

	fmt.Printf("keys: %+v\n", keys)

	key, err := ds.ListFirst("events")
	assert.NoError(t, err)

	fmt.Printf("first key is: %s\n", key)
	assert.Equal(t, key, "events/0000")

	// Nonexistant prefix
	key, err = ds.ListFirst("nonexistant")
	fmt.Printf("%t\n", err != nil)
	assert.Error(t, err)

	fmt.Printf("first key is: %s\n", key)
}

func TestAddDeltaLog_empty_doc(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()
	s3 := env.S3(ctx)
	rds := env.Redis(ctx)
	docId := uuid.NewString()

	ds := rogue.NewDocStore(s3, env.Query(ctx), rds)

	doc := v3.NewRogueForQuill("1")
	op, err := doc.Insert(0, "hello world")
	assert.NoError(t, err)

	seq, err := ds.AddUpdate(ctx, docId, op)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), seq)
}

func TestAddDeltaLog_exsiting_counter(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()
	s3 := env.S3(ctx)
	rds := env.Redis(ctx)
	docId := uuid.NewString()

	ds := rogue.NewDocStore(s3, env.Query(ctx), rds)

	doc := v3.NewRogueForQuill("1")
	op, err := doc.Insert(0, "hello world")
	assert.NoError(t, err)

	seq, err := ds.AddUpdate(ctx, docId, op)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), seq)

	op2, err := doc.Insert(10, "!")
	assert.NoError(t, err)

	seq2, err := ds.AddUpdate(ctx, docId, op2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), seq2)

	ops, err := ds.GetDeltaLog(ctx, docId, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ops))
}

func TestAddDeltaLog_deleted_counter_with_no_snapshot(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()
	s3 := env.S3(ctx)
	rds := env.Redis(ctx)
	docId := uuid.NewString()

	ds := rogue.NewDocStore(s3, env.Query(ctx), rds)

	doc := v3.NewRogueForQuill("1")
	op, err := doc.Insert(0, "hello world")
	assert.NoError(t, err)

	seq, err := ds.AddUpdate(ctx, docId, op)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), seq)

	counterKey := fmt.Sprintf(constants.DocCounterKeyFormat, docId)
	err = rds.Del(ctx, counterKey).Err()
	require.NoError(t, err)

	op2, err := doc.Insert(10, "!")
	assert.NoError(t, err)

	seq2, err := ds.AddUpdate(ctx, docId, op2)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), seq2)

	ops, err := ds.GetDeltaLog(ctx, docId, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(ops))
}

func TestAddDeltaLog_deleted_counter_with_snapshot(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()
	s3 := env.S3(ctx)
	rds := env.Redis(ctx)
	docId := uuid.NewString()

	ds := rogue.NewDocStore(s3, env.Query(ctx), rds)

	doc := v3.NewRogueForQuill("1")
	op, err := doc.Insert(0, "hello world")
	assert.NoError(t, err)

	seq, err := ds.AddUpdate(ctx, docId, op)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), seq)

	err = ds.SnapshotDoc(ctx, docId)
	require.NoError(t, err)

	counterKey := fmt.Sprintf(constants.DocCounterKeyFormat, docId)
	err = rds.Del(ctx, counterKey).Err()
	require.NoError(t, err)

	op2, err := doc.Insert(10, "!")
	assert.NoError(t, err)

	seq2, err := ds.AddUpdate(ctx, docId, op2)
	assert.NoError(t, err)
	assert.Equal(t, int64(2), seq2)

	ops, err := ds.GetDeltaLog(ctx, docId, 0, -1)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ops))
}

func TestGetCurrentDoc(t *testing.T) {
	testutils.EnsureStorage()
	ctx := testutils.TestContext()
	s3 := env.S3(ctx)
	rds := env.Redis(ctx)
	docID := uuid.NewString()

	ds := rogue.NewDocStore(s3, env.Query(ctx), rds)

	var op v3.Op
	var err error
	ops := []v3.Op{}

	doc := v3.NewRogueForQuill("1")
	op, err = doc.Insert(0, "hello world")
	require.NoError(t, err)
	ops = append(ops, op)

	op, err = doc.Delete(0, 5)
	require.NoError(t, err)
	ops = append(ops, op)

	op, err = doc.Format(0, 1, v3.FormatV3BulletList(0))
	require.NoError(t, err)
	ops = append(ops, op)

	addr, err := doc.GetFullAddress()
	require.NoError(t, err)

	op, err = doc.Insert(0, "goodbye")
	require.NoError(t, err)
	ops = append(ops, op)

	op, err = doc.Insert(13, "!")
	require.NoError(t, err)
	ops = append(ops, op)

	op, err = doc.Format(0, 1, v3.FormatV3Header(1))
	require.NoError(t, err)
	ops = append(ops, op)

	fullHtml, err := doc.GetFullHtml(true, false)
	require.NoError(t, err)
	fmt.Printf("fullHtml: %s\n", fullHtml)

	firstID, err := doc.GetFirstTotID()
	require.NoError(t, err)

	lastID, err := doc.GetLastTotID()
	require.NoError(t, err)

	op, err = doc.Rewind(firstID, lastID, *addr)
	require.NoError(t, err)
	ops = append(ops, op)

	for _, op := range ops {
		_, err = ds.AddUpdate(ctx, docID, op)
		require.NoError(t, err)
	}

	fullHtml, err = doc.GetFullHtml(true, false)
	require.NoError(t, err)
	fmt.Printf("fullHtml: %s\n", fullHtml)

	_, r, err := ds.GetCurrentDoc(ctx, docID)
	require.NoError(t, err)

	snap, err := r.NewSnapshotOp()
	require.NoError(t, err)

	newRogue := v3.NewRogue("1")
	_, err = newRogue.MergeOp(snap)
	require.NoError(t, err)

	rHtml, err := r.GetFullHtml(true, false)
	require.NoError(t, err)

	newHtml, err := newRogue.GetFullHtml(true, false)
	require.NoError(t, err)

	fmt.Printf("rHtml: %s\n", rHtml)
	fmt.Printf("newHtml: %s\n", newHtml)
	require.Equal(t, rHtml, newHtml)
}

func TestGetCurrentDocFuzz(t *testing.T) {
	type ContentAddressHtml struct {
		ContentAddress *v3.ContentAddress
		Html           string
		Plaintext      string
	}

	testutils.EnsureStorage()
	ctx := testutils.TestContext()
	s3 := env.S3(ctx)
	rds := env.Redis(ctx)
	docID := uuid.NewString()

	ds := rogue.NewDocStore(s3, env.Query(ctx), rds)
	doc := v3.NewRogueForQuill("1")

	var op v3.Op
	var err error
	ops := []v3.Op{}
	cas := []ContentAddressHtml{}

	rng := rand.New(rand.NewSource(42))

	for k := 0; k < 1000; k++ {
		prob := rng.Float64()

		if prob < 0.1 && doc.VisSize > 0 {
			_, op, err = doc.RandDelete(rng, 10)
			require.NoError(t, err)
		} else if prob < 0.3 && doc.VisSize > 2 {
			_, op, err = doc.RandFormat(rng, 10)
			require.NoError(t, err)
		} else if prob < 0.4 && doc.VisSize > 1 {
			op, err = doc.UndoDoc()
			require.NoError(t, err)
		} else if prob <= 1.0 {
			_, op, err = doc.RandInsert(rng, 10)
			require.NoError(t, err)
		}

		ops = append(ops, op)
		prob = rng.Float64()

		// take a content address 10% of the time
		if rng.Float64() < 0.01 {
			ca, err := doc.GetFullAddress()
			require.NoError(t, err)

			html, err := doc.GetFullHtml(true, false)
			require.NoError(t, err)

			cas = append(cas, ContentAddressHtml{
				ContentAddress: ca,
				Html:           html,
				Plaintext:      doc.GetText(),
			})
		}
	}

	for _, op := range ops {
		_, err = ds.AddUpdate(ctx, docID, op)
		require.NoError(t, err)
	}

	fullHtml, err := doc.GetFullHtml(true, false)
	require.NoError(t, err)
	fmt.Printf("fullHtml: %s\n", fullHtml)

	_, r, err := ds.GetCurrentDoc(ctx, docID)
	require.NoError(t, err)

	snap, err := r.NewSnapshotOp()
	require.NoError(t, err)

	newRogue := v3.NewRogue("1")
	_, err = newRogue.MergeOp(snap)
	require.NoError(t, err)

	rHtml, err := r.GetFullHtml(true, false)
	require.NoError(t, err)

	newHtml, err := newRogue.GetFullHtml(true, false)
	require.NoError(t, err)

	fmt.Printf("rHtml: %s\n", rHtml)
	fmt.Printf("newHtml: %s\n", newHtml)
	require.Equal(t, rHtml, newHtml)

	fmt.Printf("len(cas): %d\n", len(cas))
	for _, ca := range cas {
		_, doc, err := ds.GetCurrentDoc(ctx, docID)
		require.NoError(t, err)

		/*
			startID, err := doc.GetFirstTotID()
			require.NoError(t, err)

			endID, err := doc.GetLastTotID()
			require.NoError(t, err)
		*/

		startIx := rng.Intn(doc.VisSize)
		endIx := startIx + rng.Intn(doc.VisSize-startIx)

		if startIx >= endIx {
			continue
		}

		startID, err := doc.Rope.GetVisID(startIx)
		require.NoError(t, err)

		endID, err := doc.Rope.GetVisID(endIx)
		require.NoError(t, err)

		op, err = doc.Rewind(startID, endID, *ca.ContentAddress)
		require.NoError(t, err)

		_, err = ds.AddUpdate(ctx, docID, op)
		require.NoError(t, err)

		_, r, err := ds.GetCurrentDoc(ctx, docID)
		require.NoError(t, err)

		snap, err := r.NewSnapshotOp()
		require.NoError(t, err)

		newRogue := v3.NewRogue("1")
		_, err = newRogue.MergeOp(snap)
		require.NoError(t, err)

		rHtml, err := r.GetFullHtml(true, false)
		require.NoError(t, err)

		newHtml, err := newRogue.GetFullHtml(true, false)
		require.NoError(t, err)

		require.Equal(t, rHtml, newHtml)
	}
}

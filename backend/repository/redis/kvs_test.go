package redis

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"testing"
	"time"

	"github.com/arumakan1727/todo-app-go-react/config"
	"github.com/arumakan1727/todo-app-go-react/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newKVSForTest(t *testing.T) domain.KVS {
	t.Helper()
	k, err := NewKVS(context.Background(), config.ForTesting())
	if err != nil {
		t.Fatalf("newKVSForTest: cannot open Redis: %#v", err)
	}
	return k
}

func randToken(t *testing.T) domain.AuthToken {
	randBytes := make([]byte, 12)
	if _, err := rand.Read(randBytes); err != nil {
		t.Fatalf("cannot read crypt rand bytes: %#v", err)
	}

	out := bytes.NewBuffer(make([]byte, 0, len(randBytes)*2))
	out.Write([]byte("testing/")) // prefix 'testing/'

	enc := base64.NewEncoder(base64.StdEncoding, out)
	if _, err := enc.Write(randBytes); err != nil {
		t.Fatalf("cannot encode randBytes as base64: %#v", err)
	}
	return domain.AuthToken(out.Bytes())
}

func TestKVS(t *testing.T) {
	t.Parallel()
	k := newKVSForTest(t)

	var (
		tokenA         = randToken(t)
		tokenB         = randToken(t)
		tokenEphemeral = randToken(t)
	)
	const (
		expirationNormal = time.Second * 3
		expirationShort  = time.Millisecond * 20
	)

	in := []struct {
		key domain.AuthToken
		uid domain.UserID
		exp time.Duration
	}{
		{
			key: tokenA,
			uid: 1,
			exp: expirationNormal,
		},
		{
			key: tokenB,
			uid: 2,
			exp: expirationNormal,
		},
		{
			key: tokenA,
			uid: 3,
			exp: expirationNormal,
		},
		{
			key: tokenEphemeral,
			uid: 4,
			exp: expirationShort,
		},
	}

	ctx := context.Background()
	dict := make(map[domain.AuthToken]domain.UserID)

	t.Run("SaveAuth-OK", func(t *testing.T) {
		for i := range in {
			in := &in[i]
			err := k.SaveAuth(ctx, in.key, in.uid, in.exp)
			require.NoError(t, err, "err on testcase %d", i)
			dict[in.key] = in.uid
		}
	})
	t.Run("FetchAuth-OK", func(t *testing.T) {
		for key, wantUID := range dict {
			got, err := k.FetchAuth(ctx, key)
			if assert.NoError(t, err) {
				assert.Equal(t, wantUID, got)
			}
		}
	})
	t.Run("DeleteAuth-OK", func(t *testing.T) {
		key := tokenA
		err := k.DeleteAuth(ctx, key)
		assert.NoError(t, err)
		delete(dict, key)
	})
	t.Run("FetchAuth-deleted-auth-should-be-NotFound", func(t *testing.T) {
		key := tokenA
		got, err := k.FetchAuth(ctx, key)
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Equal(t, domain.UserID(0), got)
	})
	t.Run("FetchAuth-undeleted-auth-should-be-found", func(t *testing.T) {
		for key, wantUID := range dict {
			got, err := k.FetchAuth(ctx, key)
			if assert.NoError(t, err) {
				assert.Equal(t, wantUID, got)
			}
		}
	})
	time.Sleep(expirationShort)
	t.Run("FetchAuth-expired-auth-should-be-NotFound", func(t *testing.T) {
		key := tokenEphemeral
		got, err := k.FetchAuth(ctx, key)
		assert.ErrorIs(t, err, domain.ErrNotFound)
		assert.Equal(t, domain.UserID(0), got)
	})
}

package state

import (
	"context"
	"fmt"
	"testing"

	mh "gx/ipfs/QmZyZDi491cCNTLfAhwcaDii2Kg4pwKRkhqQzURGDvY6ua/go-multihash"
	"gx/ipfs/QmcZfnkapfECQGcLZaf9B79NRg7cRa9EnZh4LSbkCzwNvY/go-cid"
	"gx/ipfs/QmdtiofXbibTe6Day9ii5zjBZpSRm8vhfoerrNuY3sAQ7e/go-hamt-ipld"

	"github.com/filecoin-project/go-filecoin/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatePutGet(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cst := hamt.NewCborStore()
	tree := NewEmptyStateTree(cst)

	act1 := types.NewActor(types.AccountActorCodeCid, nil)
	act1.WriteStorage([]byte("hello"))
	act1.IncNonce()
	act2 := types.NewActor(types.AccountActorCodeCid, nil)
	act2.WriteStorage([]byte("world"))
	act2.IncNonce()
	act2.IncNonce()

	addrGetter := types.NewAddressForTestGetter()
	addr1 := addrGetter()
	addr2 := addrGetter()

	assert.NoError(tree.SetActor(ctx, addr1, act1))
	assert.NoError(tree.SetActor(ctx, addr2, act2))

	act1out, err := tree.GetActor(ctx, addr1)
	assert.NoError(err)
	assert.Equal(act1, act1out)
	act2out, err := tree.GetActor(ctx, addr2)
	assert.NoError(err)
	assert.Equal(act2, act2out)

	// now test it persists across recreation of tree
	tcid, err := tree.Flush(ctx)
	assert.NoError(err)

	tree2, err := LoadStateTree(ctx, cst, tcid, nil)
	assert.NoError(err)

	act1out2, err := tree2.GetActor(ctx, addr1)
	assert.NoError(err)
	assert.Equal(act1, act1out2)
	act2out2, err := tree2.GetActor(ctx, addr2)
	assert.NoError(err)
	assert.Equal(act2, act2out2)
}

func TestStateErrors(t *testing.T) {
	assert := assert.New(t)
	ctx := context.Background()
	cst := hamt.NewCborStore()
	tree := NewEmptyStateTree(cst)

	a, err := tree.GetActor(ctx, types.NewAddressForTestGetter()())
	assert.Nil(a)
	assert.Error(err)
	assert.True(IsActorNotFoundError(err))

	c, err := cid.NewPrefixV0(mh.BLAKE2B_MIN + 31).Sum([]byte("cats"))
	assert.NoError(err)

	tr2, err := LoadStateTree(ctx, cst, c, nil)
	assert.EqualError(err, "failed to load node: not found")
	assert.Nil(tr2)
}

func TestStateGetOrCreate(t *testing.T) {
	ctx := context.Background()
	cst := hamt.NewCborStore()
	tree := NewEmptyStateTree(cst)

	addr := types.NewAddressForTestGetter()()

	// no actor - error
	t.Run("no actor - error", func(t *testing.T) {
		assert := assert.New(t)

		actor, err := tree.GetOrCreateActor(ctx, addr, func() (*types.Actor, error) {
			return nil, fmt.Errorf("fail")
		})
		assert.EqualError(err, "fail")
		assert.Nil(actor)
	})

	t.Run("no actor - success", func(t *testing.T) {
		assert := assert.New(t)

		actor, err := tree.GetOrCreateActor(ctx, addr, func() (*types.Actor, error) {
			return &types.Actor{}, nil
		})
		assert.NoError(err)
		assert.Equal(actor, &types.Actor{})
	})

	t.Run("actor exists", func(t *testing.T) {
		assert := assert.New(t)

		actor := types.NewActor(nil, types.NewTokenAmount(10))
		assert.NoError(tree.SetActor(ctx, addr, actor))

		actorBack, err := tree.GetOrCreateActor(ctx, addr, func() (*types.Actor, error) {
			return &types.Actor{}, nil
		})
		assert.NoError(err)
		assert.Equal(actorBack, actor)
	})
}

func TestGetAllActors(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	ctx := context.Background()
	cst := hamt.NewCborStore()
	tree := NewEmptyStateTree(cst)

	addr := types.NewAddressForTestGetter()()

	actor := types.Actor{Code: types.AccountActorCodeCid, Nonce: 1234, Balance: types.NewTokenAmount(123)}
	err := tree.SetActor(ctx, addr, &actor)
	tree.Flush(ctx)

	assert.NoError(err)

	addrs, actors := GetAllActors(tree)

	require.Equal(1, len(addrs))
	found := *actors[0]

	assert.Equal(addr.String(), addrs[0])
	assert.Equal(actor.Code, found.Code)
	assert.Equal(actor.Nonce, found.Nonce)
	assert.Equal(actor.Balance, found.Balance)
}
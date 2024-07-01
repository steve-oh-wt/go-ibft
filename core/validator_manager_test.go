package core

import (
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func Test_CalculateQuorum(t *testing.T) {
	t.Parallel()

	vm := &ValidatorManager{
		vpLock: &sync.RWMutex{},
	}

	cases := []struct {
		validatorsVotingPower map[common.Address]*big.Int
		signers               map[common.Address]struct{}
		hasQuorum             bool
	}{
		{
			// case total voting power 4
			validatorsVotingPower: map[common.Address]*big.Int{
				common.BytesToAddress([]byte("A")): big.NewInt(1),
				common.BytesToAddress([]byte("B")): big.NewInt(1),
				common.BytesToAddress([]byte("C")): big.NewInt(1),
				common.BytesToAddress([]byte("D")): big.NewInt(1),
			},
			// all 4 signed, has quorum (quorum is 3)
			signers: map[common.Address]struct{}{
				common.BytesToAddress([]byte("A")): {},
				common.BytesToAddress([]byte("B")): {},
				common.BytesToAddress([]byte("C")): {},
				common.BytesToAddress([]byte("D")): {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power 4
			validatorsVotingPower: map[common.Address]*big.Int{
				common.BytesToAddress([]byte("A")): big.NewInt(1),
				common.BytesToAddress([]byte("B")): big.NewInt(1),
				common.BytesToAddress([]byte("C")): big.NewInt(1),
				common.BytesToAddress([]byte("D")): big.NewInt(1),
			},
			// only two signed (quorum is 3)
			signers: map[common.Address]struct{}{
				common.BytesToAddress([]byte("A")): {},
				common.BytesToAddress([]byte("B")): {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 6
			validatorsVotingPower: map[common.Address]*big.Int{
				common.BytesToAddress([]byte("A")): big.NewInt(1),
				common.BytesToAddress([]byte("B")): big.NewInt(1),
				common.BytesToAddress([]byte("C")): big.NewInt(1),
				common.BytesToAddress([]byte("D")): big.NewInt(1),
				common.BytesToAddress([]byte("E")): big.NewInt(1),
				common.BytesToAddress([]byte("F")): big.NewInt(1),
			},
			// 5 signed (quorum should be 5)
			signers: map[common.Address]struct{}{
				common.BytesToAddress([]byte("A")): {},
				common.BytesToAddress([]byte("B")): {},
				common.BytesToAddress([]byte("C")): {},
				common.BytesToAddress([]byte("D")): {},
				common.BytesToAddress([]byte("E")): {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 6
			validatorsVotingPower: map[common.Address]*big.Int{
				common.BytesToAddress([]byte("A")): big.NewInt(1),
				common.BytesToAddress([]byte("B")): big.NewInt(1),
				common.BytesToAddress([]byte("C")): big.NewInt(1),
				common.BytesToAddress([]byte("D")): big.NewInt(1),
				common.BytesToAddress([]byte("E")): big.NewInt(1),
				common.BytesToAddress([]byte("F")): big.NewInt(1),
			},
			// only 4 signed (quorum should be 5)
			signers: map[common.Address]struct{}{
				common.BytesToAddress([]byte("A")): {},
				common.BytesToAddress([]byte("B")): {},
				common.BytesToAddress([]byte("C")): {},
				common.BytesToAddress([]byte("D")): {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 9
			validatorsVotingPower: map[common.Address]*big.Int{
				common.BytesToAddress([]byte("A")): big.NewInt(2),
				common.BytesToAddress([]byte("B")): big.NewInt(2),
				common.BytesToAddress([]byte("C")): big.NewInt(2),
				common.BytesToAddress([]byte("D")): big.NewInt(3),
			},
			// 3 signed with voting power of 6 (quorum should be 7)
			signers: map[common.Address]struct{}{
				common.BytesToAddress([]byte("A")): {},
				common.BytesToAddress([]byte("C")): {},
				common.BytesToAddress([]byte("D")): {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 9
			validatorsVotingPower: map[common.Address]*big.Int{
				common.BytesToAddress([]byte("A")): big.NewInt(2),
				common.BytesToAddress([]byte("B")): big.NewInt(2),
				common.BytesToAddress([]byte("C")): big.NewInt(2),
				common.BytesToAddress([]byte("D")): big.NewInt(3),
			},
			// only 2 signed with voting power of 5 (quorum should be 7)
			signers: map[common.Address]struct{}{
				common.BytesToAddress([]byte("A")): {},
				common.BytesToAddress([]byte("D")): {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 10
			validatorsVotingPower: map[common.Address]*big.Int{
				common.BytesToAddress([]byte("A")): big.NewInt(2),
				common.BytesToAddress([]byte("B")): big.NewInt(2),
				common.BytesToAddress([]byte("C")): big.NewInt(3),
				common.BytesToAddress([]byte("D")): big.NewInt(3),
			},
			// 3 signed with voting power of 7 (quorum should be 7)
			signers: map[common.Address]struct{}{
				common.BytesToAddress([]byte("A")): {},
				common.BytesToAddress([]byte("B")): {},
				common.BytesToAddress([]byte("D")): {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 10
			validatorsVotingPower: map[common.Address]*big.Int{
				common.BytesToAddress([]byte("A")): big.NewInt(2),
				common.BytesToAddress([]byte("B")): big.NewInt(2),
				common.BytesToAddress([]byte("C")): big.NewInt(3),
				common.BytesToAddress([]byte("D")): big.NewInt(3),
			},
			// only 2 signed with voting power of 5 (quorum should be 7)
			signers: map[common.Address]struct{}{
				common.BytesToAddress([]byte("A")): {},
				common.BytesToAddress([]byte("D")): {},
			},
			hasQuorum: false,
		},
		{
			// case total voting power of 21
			validatorsVotingPower: map[common.Address]*big.Int{
				common.BytesToAddress([]byte("A")): big.NewInt(2),
				common.BytesToAddress([]byte("B")): big.NewInt(7),
				common.BytesToAddress([]byte("C")): big.NewInt(7),
				common.BytesToAddress([]byte("D")): big.NewInt(5),
			},
			// 3 signed with voting power of 16 (quorum should be 15)
			signers: map[common.Address]struct{}{
				common.BytesToAddress([]byte("A")): {},
				common.BytesToAddress([]byte("B")): {},
				common.BytesToAddress([]byte("C")): {},
			},
			hasQuorum: true,
		},
		{
			// case total voting power of 21
			validatorsVotingPower: map[common.Address]*big.Int{
				common.BytesToAddress([]byte("A")): big.NewInt(2),
				common.BytesToAddress([]byte("B")): big.NewInt(7),
				common.BytesToAddress([]byte("C")): big.NewInt(7),
				common.BytesToAddress([]byte("D")): big.NewInt(5),
			},
			// only 2 signed with voting power of 12 (quorum should be 15)
			signers: map[common.Address]struct{}{
				common.BytesToAddress([]byte("C")): {},
				common.BytesToAddress([]byte("D")): {},
			},
			hasQuorum: false,
		},
	}

	for _, c := range cases {
		require.NoError(t, vm.setCurrentVotingPower(c.validatorsVotingPower))
		require.Equal(t, c.hasQuorum, vm.HasQuorum(c.signers))
	}
}

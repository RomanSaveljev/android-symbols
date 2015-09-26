package processor

import (
	"github.com/RomanSaveljev/android-symbols/transmitter/mock"
	"github.com/RomanSaveljev/android-symbols/transmitter/signatures"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessFileSync(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sigs := signatures.NewSignatures()
	rcv := mock.NewMockReceiver(mockCtrl)
	rcv.EXPECT().Signatures().AnyTimes().Return(sigs, nil)
	write := rcv.EXPECT().Write([]byte("0etOA2)[BQ3FQB,A7]@c")).Return(20, nil)
	rcv.EXPECT().Write([]byte("\n")).After(write).Return(1, nil)

	output := make(chan int, 1)
	err := ProcessFileSync([]byte("123456789abcdefg"), rcv, output)
	assert.NoError(err)
	assert.Equal(16, <- output)
}

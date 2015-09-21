package processor

import (
	"bytes"
	"github.com/RomanSaveljev/android-symbols/transmitter/src/lib/mock"
	"github.com/RomanSaveljev/android-symbols/transmitter/src/lib/signatures"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessFileSync(t *testing.T) {
	assert := assert.New(t)
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sigs := new(signatures.Signatures)
	rcv := mock_transmitter.NewMockReceiver(mockCtrl)
	rcv.EXPECT().Signatures().AnyTimes().Return(sigs, nil)
	write := rcv.EXPECT().Write([]byte("0etOA2)[BQ3FQB,A7]@c")).Return(20, nil)
	rcv.EXPECT().Write([]byte("\n")).After(write).Return(1, nil)
	
	reader := bytes.NewReader([]byte("123456789abcdefg"))
	err := ProcessFileSync(reader, rcv)
	assert.NoError(err)
}

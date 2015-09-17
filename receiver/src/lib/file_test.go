package receiver

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestFileNewFile(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := NewMockfileSystemWorker(mockCtrl)
	mock.EXPECT().MkdirAll("/a/b/c/d.txt").Return(nil)

	file, err := newFileInjected("/a/b/c/d.txt", mock)
	assert.NoError(t, err)
	assert.NotNil(t, file)
}

func TestFileNewFileError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := NewMockfileSystemWorker(mockCtrl)
	mock.EXPECT().MkdirAll("/a/b/c/d.txt").Return(io.EOF)

	_, err := newFileInjected("/a/b/c/d.txt", mock)
	assert.Equal(t, io.EOF, err)
}

func TestFileNextSignature(t *testing.T) {
	assert := assert.New(t)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mock := NewMockfileSystemWorker(mockCtrl)
	mkdirAll := mock.EXPECT().MkdirAll("/a/b/c/d.txt").Return(io.EOF)
	rolling := []string{"abc", "123"}
	readRolling := mock.EXPECT().Readdirnames("/a/b/c/d.txt").After(mkdirAll).Return(rolling, nil)
	strongOne := []string{"def", "ghi"}
	mock.EXPECT().Readdirnames("/a/b/c/d.txt/abc").After(readRolling).Return(strongOne, nil)
	strongTwo := []string{"456", "789"}
	mock.EXPECT().Readdirnames("/a/b/c/d.txt/123").After(readRolling).Return(strongTwo, nil)

	file, err := newFileInjected("/a/b/c/d.txt", mock)
	assert.NotNil(t, file)
	signature, err := file.nextSignature()
	assert.NoError(err)
	assert.Equal(Signature{Rolling: "abc", Strong: "def"}, signature)
	signature, err = file.nextSignature()
	assert.NoError(err)
	assert.Equal(Signature{Rolling: "abc", Strong: "ghi"}, signature)
	signature, err = file.nextSignature()
	assert.NoError(err)
	assert.Equal(Signature{Rolling: "123", Strong: "456"}, signature)
	signature, err = file.nextSignature()
	assert.NoError(err)
	assert.Equal(Signature{Rolling: "123", Strong: "789"}, signature)
	_, err = file.nextSignature()
	assert.Equal(io.EOF, err)
}

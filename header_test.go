package sequencefile

import "testing"
import "bytes"
import "github.com/stretchr/testify/assert"

func TestReadBool(t *testing.T) {
	buf := new(bytes.Buffer)
	buf.Write([]byte{0x00})
	r := NewReader(buf)
	ret, err := r.readBoolean()
	assert.Equal(t, false, ret, "readBoolean should return a false value")
	assert.NoError(t, err, "readBoolean should be successful")

	buf.Write([]byte{0x01})

	ret, err = r.readBoolean()
	assert.Equal(t, true, ret, "readBoolean should return a true value")
	assert.NoError(t, err, "readBoolean should be successful")
}

func TestWriteBool(t *testing.T) {
	buf := new(bytes.Buffer)
	w := NewWriter(buf)
	num, err := w.writeBoolean(false)
	assert.Equal(t, []byte{0x00}, buf.Bytes(), "writeBoolean should write a false value")
	assert.Equal(t, 1, num, "writeBoolean should write one byte")
	assert.NoError(t, err, "writeBoolean should be successful")

	buf = new(bytes.Buffer)
	w = NewWriter(buf)
	num, err = w.writeBoolean(true)
	assert.Equal(t, []byte{0x01}, buf.Bytes(), "writeBoolean should write a true value")
	assert.Equal(t, 1, num, "writeBoolean should write one byte")
	assert.NoError(t, err, "writeBoolean should be successful")
}

type metadataSpec struct {
	Name     string
	Metadata map[string]string
	Bytes    []byte
}

var testmetadatas = []metadataSpec{
	{
		"multiple key/value pairs",
		map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
		[]byte{
			0x00, 0x00, 0x00, 0x03, // 3 pairs
			0x04, 0x6b, 0x65, 0x79, 0x31, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x31, // key1 value1
			0x04, 0x6b, 0x65, 0x79, 0x32, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32, // key1 value1
			0x04, 0x6b, 0x65, 0x79, 0x33, 0x06, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x33, // key1 value1
		},
	},
	{
		"single key/value pair",
		map[string]string{
			"foo": "bar",
		},
		[]byte{
			0x00, 0x00, 0x00, 0x01, // 1 pair
			0x03, 0x66, 0x6f, 0x6f, 0x03, 0x62, 0x61, 0x72, // foo bar
		},
	},
	{
		"no key/value pairs",
		map[string]string{},
		[]byte{
			0x00, 0x00, 0x00, 0x00, // 0 pairs
		},
	},
}

func TestWriteMetadata(t *testing.T) {
	for _, spec := range testmetadatas {
		t.Run(spec.Name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			w := NewWriter(buf)
			w.Header.Metadata = spec.Metadata
			written, err := w.writeMetadata()
			assert.NoError(t, err, "it should write successfully")
			assert.Equal(t, len(spec.Bytes), written, "it should report the correct number of bytes written")
			assert.Equal(t, spec.Bytes, buf.Bytes(), "it should write out the correct header contents")
		})
	}
}
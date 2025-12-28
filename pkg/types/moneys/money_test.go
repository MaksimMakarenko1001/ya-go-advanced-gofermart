package moneys_test

import (
	"encoding/json"
	"testing"

	"github.com/MaksimMakarenko1001/ya-go-advanced-gofermart.git/pkg/types/moneys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoney(t *testing.T) {
	tests := []struct {
		name         string // description of this test case
		input        []byte
		unmarshalErr bool
	}{
		{
			name:         "0",
			input:        []byte(`0`),
			unmarshalErr: false,
		},
		{
			name:         "0.5",
			input:        []byte(`0.5`),
			unmarshalErr: false,
		},
		{
			name:         "0.55",
			input:        []byte(`0.55`),
			unmarshalErr: false,
		},
		{
			name:         "99",
			input:        []byte(`99`),
			unmarshalErr: false,
		},
		{
			name:         "99.55",
			input:        []byte(`99.55`),
			unmarshalErr: false,
		},
		{
			name:         "-99",
			input:        []byte(`-99`),
			unmarshalErr: true,
		},
		{
			name:         "-99.45",
			input:        []byte(`-99.45`),
			unmarshalErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := assert.New(t)
			is := require.New(t)

			var m *moneys.Money
			err := json.Unmarshal(tt.input, &m)
			if err != nil && !tt.unmarshalErr {
				is.NoErrorf(err, "got error %w", err)
			}
			if tt.unmarshalErr {
				is.Errorf(err, "expects unmarshal error")
				return
			}

			output, err := json.Marshal(m)
			as.NoError(err)
			as.Equalf(tt.input, output, "output=%s", string(output))
		})
	}
}

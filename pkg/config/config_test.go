package config

import (
	"log"
	"os"
	"testing"
)

func Test_GetParam(t *testing.T) {
	cases := []struct {
		name string
		key string
		value string
		isError bool
	}{
		{
			name: "key exist, value exist",
			key: "key1",
			value: "value1",
			isError: false,
		},
		{
			name: "key exist, value empty",
			key: "key2",
			value: "value2",
			isError: true,
		},
		{
			name: "key not exist",
			key: "key3",
			value: "value3",
			isError: true,
		},
	}

	for _, cs := range cases {
		err := os.Setenv("key1", "value1")
		if err != nil {
			log.Fatal(err)
		}
		err = os.Setenv("key2", "")
		if err != nil {
			log.Fatal(err)
		}
		t.Run(cs.name, func(t *testing.T) {
			val, err := getParam(cs.key)
			if err != nil && !cs.isError {
				t.Errorf("unexpected error %v, at test %s", err, cs.name)
			}
			if err == nil && cs.isError {
				t.Errorf("expecting error, got nil, at test %s", cs.name)
			}
			if err == nil && val != cs.value {
				t.Errorf("expecting value: %s, but got: %s", val, cs.value)
			}
		})

		err = os.Unsetenv("key1")
		if err != nil {
			log.Fatal(err)
		}
		err = os.Unsetenv("key2")
		if err != nil {
			log.Fatal(err)
		}
	}
}
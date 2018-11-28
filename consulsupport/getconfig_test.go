package consulsupport

import (
    "testing"

    consulapi "github.com/hashicorp/consul/api"
)

func TestGetConfig(t *testing.T) {
    consulClient, err := NewClient(consulapi.DefaultNonPooledConfig(), "localhost:8500")
    if err != nil {
        panic(err)
    }

    const testKey = "consul/test"
    const expectedValue = "test"

    kv := &consulapi.KVPair{
        Key:   testKey,
        Flags: 0,
        Value: []byte(expectedValue),
    }

    _, err = consulClient.KV().Put(kv, nil)
    if err != nil {
        panic(err)
    }

    actualValue, err := GetConfig(consulClient, testKey)
    if err != nil {
        panic(err)
    }

    if expectedValue != actualValue {
        t.Errorf("expected value %v; but was %v", expectedValue, actualValue)
    }
}

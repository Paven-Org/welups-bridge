package manager

import (
	"bridge/libs"
	"context"
	"fmt"

	"gitlab.com/rwxrob/uniq"
	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

type Crypto interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

type SecretPropagatorConfig struct {
	Keys   []string // "Keys" as the key-value map in context.Context, not cryptographic key
	Crypto Crypto
}

type secretPropagator struct {
	keySet map[string]struct{}
	crypto Crypto
}

func MkSecretPropagator(config SecretPropagatorConfig) workflow.ContextPropagator {
	if config.Crypto == nil {
		randomKey := uniq.Hex(32)
		fmt.Printf("[MkSecretPropagator] random key: %s\n", randomKey)
		config.Crypto = libs.MkCryptor(randomKey)
	}
	keyMap := make(map[string]struct{}, len(config.Keys))
	for _, key := range config.Keys {
		keyMap[key] = struct{}{}
	}
	return &secretPropagator{
		keySet: keyMap,
		crypto: config.Crypto,
	}
}

// Inject injects values from context into headers for propagation
func (s *secretPropagator) Inject(ctx context.Context, writer workflow.HeaderWriter) error {
	for key := range s.keySet {
		if value, ok := ctx.Value(key).(string); ok {
			encryptedValue, err := s.crypto.Encrypt([]byte(value))
			//fmt.Printf("Original value: %s\n", value)
			//fmt.Printf("[Inject] Encrypted value: %x\n", encryptedValue)
			if err != nil {
				return err
			}
			encodedValue, err := converter.GetDefaultDataConverter().ToPayload(encryptedValue)
			if err != nil {
				return err
			}
			writer.Set(key, encodedValue)
		}
	}
	return nil
}

// InjectFromWorkflow injects values from context into headers for propagation
func (s *secretPropagator) InjectFromWorkflow(ctx workflow.Context, writer workflow.HeaderWriter) error {
	for key := range s.keySet {
		value := ctx.Value(key)
		if value == nil {
			return nil
		}
		encryptedValue, err := s.crypto.Encrypt([]byte(value.(string)))
		//fmt.Printf("Original value: %s\n", value)
		//fmt.Printf("[InjectFromWorkflow] Encrypted value: %x\n", encryptedValue)
		if err != nil {
			return err
		}
		encodedValue, err := converter.GetDefaultDataConverter().ToPayload(encryptedValue)
		if err != nil {
			return err
		}
		writer.Set(key, encodedValue)
	}
	return nil
}

// Extract extracts values from headers and puts them into context
func (s *secretPropagator) Extract(ctx context.Context, reader workflow.HeaderReader) (context.Context, error) {
	if err := reader.ForEachKey(func(key string, value *commonpb.Payload) error {
		if _, ok := s.keySet[key]; ok {
			var decodedValue []byte
			err := converter.GetDefaultDataConverter().FromPayload(value, &decodedValue)
			if err != nil {
				return err
			}
			//fmt.Printf("[Extract] Decoded encrypted value: %x\n", decodedValue)
			decryptedValue, err := s.crypto.Decrypt(decodedValue)
			if err != nil {
				return err
			}
			//fmt.Printf("[Extract] Decrypted value: %s\n", string(decryptedValue))
			ctx = context.WithValue(ctx, key, string(decryptedValue))
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ctx, nil
}

// ExtractToWorkflow extracts values from headers and puts them into context
func (s *secretPropagator) ExtractToWorkflow(ctx workflow.Context, reader workflow.HeaderReader) (workflow.Context, error) {
	if err := reader.ForEachKey(func(key string, value *commonpb.Payload) error {
		if _, ok := s.keySet[key]; ok {
			var decodedValue []byte
			err := converter.GetDefaultDataConverter().FromPayload(value, &decodedValue)
			if err != nil {
				return err
			}
			//fmt.Printf("[ExtractToWorkflow] Decoded encrypted value: %x\n", decodedValue)
			decryptedValue, err := s.crypto.Decrypt(decodedValue)
			if err != nil {
				return err
			}
			ctx = workflow.WithValue(ctx, key, string(decryptedValue))
			//fmt.Printf("[ExtractToWorkflow] Decrypted value: %s\n", string(decryptedValue))
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return ctx, nil
}

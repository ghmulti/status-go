package protocol

import (
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/status-im/status-go/protocol/protobuf"
	"github.com/status-im/status-go/protocol/v1"
)

// maxWhisperDrift is how many milliseconds we allow the clock value to differ
// from whisperTimestamp
const maxWhisperDriftMs uint64 = 120000

func validateClockValue(clock uint64, whisperTimestamp uint64) error {
	if clock == 0 {
		return errors.New("clock can't be 0")
	}

	c := new(big.Int).SetUint64(clock)
	w := new(big.Int).SetUint64(whisperTimestamp)
	result := new(big.Int).Sub(c, w)
	difference := result.Abs(result).Uint64()

	if difference > maxWhisperDriftMs {
		return errors.New("clock value can't be too different from whisper timestamp")
	}

	return nil
}

func ValidateMembershipUpdateMessage(message *protocol.MembershipUpdateMessage, timeNowMs uint64) error {

	for _, e := range message.Events {
		// We only compare in one direction for membership update as they are relayed
		// without the original whisper timestamp
		if e.ClockValue > timeNowMs && e.ClockValue-timeNowMs > maxWhisperDriftMs {
			return errors.New("clock value can't be too different from whisper timestamp")
		}
	}
	return nil
}

func ValidateReceivedPairInstallation(message *protobuf.PairInstallation, whisperTimestamp uint64) error {
	if err := validateClockValue(message.Clock, whisperTimestamp); err != nil {
		return err
	}

	if len(strings.TrimSpace(message.Name)) == 0 {
		return errors.New("name can't be empty")
	}

	if len(strings.TrimSpace(message.DeviceType)) == 0 {
		return errors.New("device type can't be empty")
	}

	if len(strings.TrimSpace(message.InstallationId)) == 0 {
		return errors.New("installationId can't be empty")
	}

	return nil
}

func ValidateReceivedSendTransaction(message *protobuf.SendTransaction, whisperTimestamp uint64) error {
	if err := validateClockValue(message.Clock, whisperTimestamp); err != nil {
		return err
	}

	if len(strings.TrimSpace(message.TransactionHash)) == 0 {
		return errors.New("transaction hash can't be empty")
	}

	if message.Signature == nil {
		return errors.New("signature can't be nil")
	}

	return nil
}

func ValidateReceivedRequestAddressForTransaction(message *protobuf.RequestAddressForTransaction, whisperTimestamp uint64) error {
	if err := validateClockValue(message.Clock, whisperTimestamp); err != nil {
		return err
	}

	if len(strings.TrimSpace(message.Value)) == 0 {
		return errors.New("value can't be empty")
	}

	_, err := strconv.ParseFloat(message.Value, 64)
	if err != nil {
		return err
	}

	return nil
}

func ValidateReceivedRequestTransaction(message *protobuf.RequestTransaction, whisperTimestamp uint64) error {
	if err := validateClockValue(message.Clock, whisperTimestamp); err != nil {
		return err
	}

	if len(strings.TrimSpace(message.Value)) == 0 {
		return errors.New("value can't be empty")
	}

	if len(strings.TrimSpace(message.Address)) == 0 {
		return errors.New("address can't be empty")
	}

	_, err := strconv.ParseFloat(message.Value, 64)
	if err != nil {
		return err
	}

	return nil
}

func ValidateReceivedAcceptRequestAddressForTransaction(message *protobuf.AcceptRequestAddressForTransaction, whisperTimestamp uint64) error {
	if err := validateClockValue(message.Clock, whisperTimestamp); err != nil {
		return err
	}

	if len(message.Id) == 0 {
		return errors.New("messageID can't be empty")
	}

	if len(strings.TrimSpace(message.Address)) == 0 {
		return errors.New("address can't be empty")
	}

	return nil
}

func ValidateReceivedDeclineRequestAddressForTransaction(message *protobuf.DeclineRequestAddressForTransaction, whisperTimestamp uint64) error {
	if err := validateClockValue(message.Clock, whisperTimestamp); err != nil {
		return err
	}

	if len(message.Id) == 0 {
		return errors.New("messageID can't be empty")
	}

	return nil
}

func ValidateReceivedDeclineRequestTransaction(message *protobuf.DeclineRequestTransaction, whisperTimestamp uint64) error {
	if err := validateClockValue(message.Clock, whisperTimestamp); err != nil {
		return err
	}

	if len(message.Id) == 0 {
		return errors.New("messageID can't be empty")
	}

	return nil
}

func ValidateReceivedChatMessage(message *protobuf.ChatMessage, whisperTimestamp uint64) error {
	if err := validateClockValue(message.Clock, whisperTimestamp); err != nil {
		return err
	}

	if message.Timestamp == 0 {
		return errors.New("timestamp can't be 0")
	}

	if len(strings.TrimSpace(message.Text)) == 0 {
		return errors.New("text can't be empty")
	}

	if len(message.ChatId) == 0 {
		return errors.New("chatId can't be empty")
	}

	if message.ContentType == protobuf.ChatMessage_UNKNOWN_CONTENT_TYPE {
		return errors.New("unknown content type")
	}

	if message.ContentType == protobuf.ChatMessage_TRANSACTION_COMMAND {
		return errors.New("can't receive request address for transaction from others")
	}

	if message.MessageType == protobuf.ChatMessage_UNKNOWN_MESSAGE_TYPE || message.MessageType == protobuf.ChatMessage_SYSTEM_MESSAGE_PRIVATE_GROUP {
		return errors.New("unknown message type")
	}

	if message.ContentType == protobuf.ChatMessage_STICKER {
		if message.Payload == nil {
			return errors.New("no sticker content")
		}
		sticker := message.GetSticker()
		if sticker == nil {
			return errors.New("no sticker content")
		}
		if len(sticker.Hash) == 0 {
			return errors.New("sticker hash not set")
		}
	}
	return nil
}

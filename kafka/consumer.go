package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/ledgerwatch/erigon-lib/common"
)

type KafkaConsumer struct {
	consumer sarama.ConsumerGroup
	config   KafkaConfig
}

func NewKafkaConsumer(config KafkaConfig) (*KafkaConsumer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = DEFAULT_VERSION
	saramaConfig.ClientID = config.ClientID
	// Consume the latest data
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	// Disable auto-commit to always start from latest on restart
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = false

	// Create consumer group
	consumerGroup, err := sarama.NewConsumerGroup(config.BootstrapServers, config.ClientID, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("error creating Kafka consumer: %v", err)
	}

	return &KafkaConsumer{
		consumer: consumerGroup,
		config:   config,
	}, nil
}

// ConsumeKafka starts consuming kafka messages from the specified topics
func (client *KafkaConsumer) ConsumeKafka(
	ctx context.Context,
	heightChan chan int64,
	addrBalanceChan chan common.Address,
	tokenHolderChan chan TokenHolderData,
	errorChan chan error,
	logger *log.Logger,
) {
	handler := &consumerGroupHandler{
		ctx:             ctx,
		heightChan:      heightChan,
		addrBalanceChan: addrBalanceChan,
		tokenHolderChan: tokenHolderChan,
		errorChan:       errorChan,
		logger:          logger,
	}

	topics := []string{client.config.StateTopic, client.config.NonStateTopic}
	err := client.consumer.Consume(ctx, topics, handler)
	if err != nil {
		errorChan <- fmt.Errorf("ConsumeKafka error: %v", err)
		return
	}
}

func (client *KafkaConsumer) Close() error {
	return client.consumer.Close()
}

type consumerGroupHandler struct {
	ctx             context.Context
	heightChan      chan int64
	addrBalanceChan chan common.Address
	tokenHolderChan chan TokenHolderData
	errorChan       chan error
	logger          *log.Logger
}

func (h *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-h.ctx.Done():
			err := fmt.Errorf("context cancelled - stopping consume claim")
			h.errorChan <- err
			return err
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			var kafkaData KafkaData
			if err := json.Unmarshal(msg.Value, &kafkaData); err != nil {
				if h.logger != nil {
					h.logger.Printf("kafka consume claim error, unmarshaling state message: %v\n", err)
				}
				continue
			}
			switch kafkaData.Type {
			case BlockMessageType:
				heightData, exists := kafkaData.Data["height"]
				if !exists {
					return fmt.Errorf("missing height field in block message")
				}
				heightFloat, ok := heightData.(float64)
				if !ok {
					return fmt.Errorf("height field is not a number: %T, value: %v", heightData, heightData)
				}
				height := int64(heightFloat)
				// Send message to height channel
				select {
				case h.heightChan <- height:
				case <-h.ctx.Done():
					err := fmt.Errorf("context cancelled - stopping consume claim")
					h.errorChan <- err
					return err
				}
			case AddressMessageType:
				addressData, exists := kafkaData.Data[AddressField]
				if !exists {
					return fmt.Errorf("missing address field in address messagae")
				}
				addressStr := addressData.(string)
				address := common.HexToAddress(addressStr)
				// Send address to address channel
				select {
				case h.addrBalanceChan <- address:
				case <-h.ctx.Done():
					err := fmt.Errorf("context cancelled - stopping consume claim")
					h.errorChan <- err
					return err
				}
			case TokenHolderMessageType:
				holderData, exists := kafkaData.Data[HolderAddressField]
				if !exists {
					return fmt.Errorf("missing holder address field in token holder message")
				}
				holderAddressStr := holderData.(string)
				holderAddress := common.HexToAddress(holderAddressStr)

				tokenAddressData, exists := kafkaData.Data[TokenContractAddressField]
				if !exists {
					return fmt.Errorf("missing token contract address field in token holder message")
				}
				tokenAddressStr := tokenAddressData.(string)
				tokenAddress := common.HexToAddress(tokenAddressStr)
				// Send token holder data to token holder channel
				tokenHolderData := TokenHolderData{
					Address:      holderAddress,
					TokenAddress: tokenAddress,
				}
				select {
				case h.tokenHolderChan <- tokenHolderData:
				case <-h.ctx.Done():
					err := fmt.Errorf("context cancelled - stopping consume claim")
					h.errorChan <- err
					return err
				}
			}
		}
	}
}

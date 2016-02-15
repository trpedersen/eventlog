package constants

const (
	EVENT_DATA_LEN  = 1000
	EVENT_TOPIC_LEN = 100
	EVENT_TIME_LEN  = 15 // time MarshalBinary length
	EVENT_BYTES_LEN = EVENT_TIME_LEN + EVENT_TOPIC_LEN + EVENT_DATA_LEN
)

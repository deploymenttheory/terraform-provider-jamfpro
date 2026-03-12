package guid_list_sharder

import "errors"

var (
	ErrInvalidShardCount                        = errors.New("invalid shard count")
	ErrInvalidShardName                         = errors.New("invalid shard name")
	ErrReservedIDInMultipleShards               = errors.New("reserved ID appears in multiple shards")
	ErrUnknownStrategy                          = errors.New("unknown strategy")
	ErrUnknownSourceType                        = errors.New("unknown source_type")
	ErrUnsupportedSourceTypeForGroupConvenience = errors.New("source_type does not support group membership convenience inputs")
)

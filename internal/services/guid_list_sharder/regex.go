package guid_list_sharder

import "regexp"

var numericIDRe = regexp.MustCompile(`^\d+$`)
var shardNameRe = regexp.MustCompile(`^shard_\d+$`)

package bolt_operations

func GetExactQuotaForType(quotaType string) uint64 { // unfinished, new update
	switch quotaType {
	case "default":
		return 1024 * 1024 * 1024 // 1 GB
	}
	return 0
}

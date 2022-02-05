package grpc_go_starter

// Deduplicate 将a和b中的字符串按顺序连在一起，且去重
func Deduplicate(a, b []string) []string {
	r := make([]string, 0, len(a)+len(b))
	m := make(map[string]bool)
	for _, s := range append(a, b...) {
		if _, ok := m[s]; !ok {
			m[s] = true
			r = append(r, s)
		}
	}
	return r

}

package reactapp

func mergeString(s1 string, s2 string) string {
	if s2 != "" {
		return s2
	}
	return s1
}

func mergeStringList(list1 []string, list2 []string) []string {
	if len(list2) > 0 {
		return append(list1, list2...)
	}
	return list1
}

func mergeStringMap(m1 map[string]string, m2 map[string]string) map[string]string {
	if m2 != nil {
		for k, v := range m2 {
			m1[k] = v
		}
	}
	return m1
}

func mergeObjectMap(m1 map[string]interface{}, m2 map[string]interface{}) map[string]interface{} {
	if m2 != nil {
		for k2, v2 := range m2 {
			if v1, ok := m1[k2]; ok {
				m1[k2] = mergeObject(v1, v2)
			} else {
				m1[k2] = v2
			}
		}
	}
	return m1
}

type Merger interface {
	Merge(interface{}) interface{}
}

func mergeObject(o1 interface{}, o2 interface{}) interface{} {
	if o1 == nil && o2 != nil {
		return o2
	}
	if m1, ok := o1.(Merger); ok {
		return m1.Merge(o2)
	}
	return o1
}

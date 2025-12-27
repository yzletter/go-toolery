package serializer

import "github.com/bytedance/sonic"

type JsonByBytedanceSonic struct {
}

func (j JsonByBytedanceSonic) Marshal(object any) ([]byte, error) {
	bs, err := sonic.Marshal(object)
	return bs, err
}

func (j JsonByBytedanceSonic) Unmarshal(buffer []byte, object any) error {
	err := sonic.Unmarshal(buffer, object)
	return err
}

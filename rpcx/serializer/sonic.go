package serializer

import "github.com/bytedance/sonic"

type BytedanceSonic struct {
}

func (b BytedanceSonic) Marshal(object any) ([]byte, error) {
	bs, err := sonic.Marshal(object)
	return bs, err
}

func (b BytedanceSonic) Unmarshal(buffer []byte, object any) error {
	err := sonic.Unmarshal(buffer, object)
	return err
}

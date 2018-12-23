package ranger

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"os"
	"strings"
)

var r *Ranger

func init() {
	r = New()
}

// Ranger preserves the environment.
type Ranger struct{
	defaults map[string]interface{}
	required map[string]bool
}

// New creates a new Ranger instance.
func New() *Ranger {
	r := new(Ranger)
	r.defaults = make(map[string]interface{})
	r.required = make(map[string]bool)

	return r
}

// SetDefault sets the default value for a key.
func SetDefault(key string, value interface{}) { r.SetDefault(key, value) }
func (r *Ranger) SetDefault(key string, value interface{}) {
	r.defaults[key] = value
}

// SetRequired marks a key as required. If a required key isn't found in the environment
// an error will be returned while unmarshaling.
func SetRequired(key string) { r.SetRequired(key) }
func (r *Ranger) SetRequired(key string) {
	r.required[key] = true
}

// Unmarshal pulls the configured values from the environment and preserves them into the provided interface.
func Unmarshal(output interface{}) error {
	return r.Unmarshal(output)
}
func (r *Ranger) Unmarshal(output interface{}) error {
	decodeConfig := &mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	}
	decoder, err := mapstructure.NewDecoder(decodeConfig)
	if err != nil {
		return err
	}

	input, err := r.loadSettings()
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

func (r *Ranger) loadSettings() (map[string]interface{}, error) {
	m := map[string]interface{}{}
	for k, required := range r.keys() {
		value := r.get(k)
		if value == nil{
			if required {
				return nil, fmt.Errorf("missing required key '%s'", k)
			} else {
				value = ""
			}
		}
		m[k] = value
	}
	return m, nil
}

func (r *Ranger) keys() map[string]bool {
	a := r.required
	for k := range r.defaults {
		a[k] = false
	}
	return a
}

func (r *Ranger) get(key string) interface{} {
	envKey := strings.ToUpper(key)
	envVal, ok := os.LookupEnv(envKey)
	if ok {
		return envVal
	}
	defaultVal, ok := r.defaults[key]
	if ok {
		return defaultVal
	}
	return nil
}

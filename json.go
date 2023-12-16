package rozer

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pdk/rozer/pushback"
)

func (r Roze) MarshalJSON() ([]byte, error) {
	if len(r.data) != len(r.names) {
		// marshal as array.
		return json.Marshal(r.data)
	}
	// invert the names map.
	names := make([]string, len(r.names))
	for k, v := range r.names {
		names[v] = k
	}
	// marshal as object.
	m := bytes.Buffer{}
	m.WriteByte('{')
	for i, v := range names {
		if i > 0 {
			m.WriteByte(',')
		}
		key := v
		value := r.data[i]
		keyBytes, err := json.Marshal(key)
		if err != nil {
			return m.Bytes(), err
		}
		m.Write(keyBytes)
		m.WriteByte(':')
		valueBytes, err := json.Marshal(value)
		if err != nil {
			return m.Bytes(), err
		}
		m.Write(valueBytes)
	}
	m.WriteByte('}')

	return m.Bytes(), nil
}

func (r *Roze) UnmarshalJSON(data []byte) error {

	dec := pushback.NewDecoder(bytes.NewReader(data))
	return r.parseObjectOrArray(dec)
}

func (r *Roze) parseObjectOrArray(dec *pushback.Decoder) error {

	next, err := dec.Token()
	if err != nil {
		return err
	}

	delim, ok := next.(json.Delim)
	if !ok {
		return fmt.Errorf("expect JSON object or array open with '{' or '['")
	}

	switch delim {
	case '{':
		err = r.parseObject(dec)
		if err != nil {
			return err
		}
		return confirmClose(dec, '}')
	case '[':
		err = r.parseArray(dec)
		if err != nil {
			return err
		}
		return confirmClose(dec, ']')
	default:
		return fmt.Errorf("unexpected JSON token: %T: %v", next, next)
	}
}

func confirmClose(dec *pushback.Decoder, delim json.Delim) error {
	next, err := dec.Token()
	if err != nil {
		return err
	}

	if next != delim {
		return fmt.Errorf("expected JSON object or array close with '%c'", delim)
	}

	return nil
}

func (r *Roze) parseObject(dec *pushback.Decoder) error {
	for {
		next, err := dec.Token()
		if err != nil {
			return err
		}

		delim, ok := next.(json.Delim)
		if ok {
			if delim == '}' {
				// object is complete. yay.
				dec.Pushback(next)
				return nil
			}
			return fmt.Errorf("unexpected JSON token (should be string key): %T: %v", next, next)
		}

		// TODO check if json.Decoder guarantees this is a string.
		key, ok := next.(string)
		if !ok {
			return fmt.Errorf("unexpected JSON token (check if this ever happens): %T: %v", next, next)
		}

		// json.Decoder eats the colon for us. value is next.

		next, err = dec.Token()
		if err != nil {
			return err
		}

		delim, ok = next.(json.Delim)
		if ok {
			if delim == '{' || delim == '[' {
				// value is a nested object or array
				dec.Pushback(next)
				value := &Roze{}
				err = value.parseObjectOrArray(dec)
				if err != nil {
					return err
				}
				r.Put(key, value)
			} else {
				return fmt.Errorf("unexpected JSON token (should be a value): %T: %v", next, next)
			}
		} else {
			// value is a simple value.
			r.Put(key, next)
		}

		// json.Decoder eats the comma for us.
	}
}

func (r *Roze) parseArray(dec *pushback.Decoder) error {
	for {
		next, err := dec.Token()
		if err != nil {
			return err
		}

		delim, ok := next.(json.Delim)
		if ok {
			switch delim {
			case ']':
				// array is complete. yay.
				dec.Pushback(next)
				return nil
			case '{', '[':
				// value is a nested object or array
				dec.Pushback(next)
				value := &Roze{}
				err = value.parseObjectOrArray(dec)
				if err != nil {
					return err
				}
				r.Append(value)
			default:
				return fmt.Errorf("unexpected JSON token (should be a value): %T: %v", next, next)
			}
		} else {
			// value is a simple value
			r.Append(next)
		}

		// json.Decoder eats the comma for us.
	}
}

package hashelper

import (
	"encoding"
	"encoding/json"
	"errors"
	"hash"
	"io"
)

var (
	ErrHashStatesNotMatch = errors.New("hashes and states do not match")
)

type Task struct {
	r      io.Reader
	buf    []byte
	hashes []hash.Hash
	w      io.Writer
	summed int64
}

type State struct {
	Summed int64    `json:"summed"`
	Datas  [][]byte `json:"datas"`
}

func NewTask(r io.Reader, bufferSize int64, hashes ...hash.Hash) *Task {
	var writers []io.Writer

	if bufferSize <= 0 {
		bufferSize = 32 * 1024
	}

	buf := make([]byte, bufferSize)

	for _, h := range hashes {
		writers = append(writers, h)
	}

	w := io.MultiWriter(writers...)

	return &Task{r, buf, hashes, w, 0}
}

func (t *Task) MarshalBinary() ([]byte, error) {
	s := State{t.summed, nil}

	for _, h := range t.hashes {
		m := h.(encoding.BinaryMarshaler)
		data, _ := m.MarshalBinary()
		s.Datas = append(s.Datas, data)
	}

	return json.Marshal(s)
}

func (t *Task) UnmarshalBinary(data []byte) error {
	var s State

	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	t.summed = s.Summed

	if len(t.hashes) != len(s.Datas) {
		return ErrHashStatesNotMatch
	}

	for i, h := range t.hashes {
		u := h.(encoding.BinaryUnmarshaler)
		data := s.Datas[i]

		if err := u.UnmarshalBinary(data); err != nil {
			return err
		}
	}

	return nil
}

func (t *Task) NextStep() (bool, error) {
	n, err := io.CopyBuffer(t.w, t.r, t.buf)
	if err != nil {
		return false, err
	}

	if n == 0 {
		return false, nil
	}

	t.summed += n
	return true, nil
}

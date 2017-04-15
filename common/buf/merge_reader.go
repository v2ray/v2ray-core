package buf

type MergingReader struct {
	reader        Reader
	timeoutReader TimeoutReader
}

func NewMergingReader(reader Reader) Reader {
	return &MergingReader{
		reader:        reader,
		timeoutReader: reader.(TimeoutReader),
	}
}

func (r *MergingReader) Read() (MultiBuffer, error) {
	mb, err := r.reader.Read()
	if err != nil {
		return nil, err
	}

	if r.timeoutReader == nil {
		return mb, nil
	}

	for {
		mb2, err := r.timeoutReader.ReadTimeout(0)
		if err != nil {
			break
		}
		mb.AppendMulti(mb2)
	}

	return mb, nil
}

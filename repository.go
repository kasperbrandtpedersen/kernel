package kernel

type repository[S State] struct {
	serializer Serializer
	store      Store
	publisher  Publisher
}

func (r *repository[S]) Load(stream string) (S, error) {
	var s S
	h, err := r.store.Load(stream)

	if err != nil {
		return s, err
	}

	for _, rec := range h {
		e, err := r.serializer.Deserialize(rec)

		if err != nil {
			return s, err
		}

		s.On(e)
	}

	return s, nil
}

func (r *repository[_]) Save(stream string, events []Event) error {
	if len(events) == 0 {
		return nil
	}

	var h History

	for _, e := range events {
		r, err := r.serializer.Serialize(e)

		if err != nil {
			return err
		}

		h = append(h, r)
	}

	if err := r.store.Save(stream, h); err != nil {
		return err
	}

	for _, e := range events {
		r.publisher.Publish(e)
	}

	return nil
}

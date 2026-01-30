-- +goose Up
-- +goose StatementBegin

INSERT INTO markets (name, state, country, latitude, longitude) 
VALUES
    ('Mile 12', 'Lagos', 'Nigeria', 6.6186, 3.3882),
    ('Oyingbo', 'Lagos', 'Nigeria', 6.4979, 3.3792),
    ('Bodija', 'Oyo', 'Nigeria', 7.4406, 3.9015),
    ('Wuse Market', 'FCT', 'Nigeria', 9.0765, 7.3986),
    ('Main Market', 'Kano', 'Nigeria', 12.0000, 8.5167)
ON CONFLICT (name, state, country) DO NOTHING;

INSERT INTO crops (name, unit)
VALUES
    ('Maize', 'kg'),
    ('Rice', 'kg'),
    ('Beans', 'kg'),
    ('Yam', 'tuber'),
    ('Tomato', 'basket')
ON CONFLICT (name) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM markets
WHERE (name, state, country) IN (
  ('Mile 12','Lagos','Nigeria'),
  ('Oyingbo','Lagos','Nigeria'),
  ('Bodija','Oyo','Nigeria'),
  ('Wuse Market','FCT','Nigeria'),
  ('Main Market','Kano','Nigeria')
);

DELETE FROM crops 
WHERE name IN (
    'Maize',
    'Rice',
    'Beans',
    'Yam',
    'Tomato'
);

-- +goose StatementEnd

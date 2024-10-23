CREATE SEQUENCE serial START 1000000000000;

CREATE TABLE urls (
  url varchar,
  slug varchar,

  CONSTRAINT unique_slug UNIQUE(slug)
);

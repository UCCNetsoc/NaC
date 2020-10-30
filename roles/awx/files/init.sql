DO $$
BEGIN
  CREATE USER root;
  EXCEPTION WHEN DUPLICATE_OBJECT THEN
  RAISE NOTICE 'not creating role root -- it already exists';
END
$$;
DO $$
BEGIN
  CREATE DATABASE root;
  EXCEPTION WHEN DUPLICATE_OBJECT THEN
  RAISE NOTICE 'not creating DB root -- it already exists';
END
$$;
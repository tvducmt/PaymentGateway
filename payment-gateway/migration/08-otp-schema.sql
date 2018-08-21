ALTER TABLE account  ADD COLUMN otpenable boolean default false;
ALTER TABLE account  ADD COLUMN otpsecretkey varchar;
ALTER TABLE account  ALTER COLUMN otpsecretkey SET DEFAULT '';
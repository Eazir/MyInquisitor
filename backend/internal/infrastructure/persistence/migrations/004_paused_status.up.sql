ALTER TABLE debts DROP CONSTRAINT IF EXISTS debts_status_check;
ALTER TABLE debts ADD CONSTRAINT debts_status_check CHECK (status IN ('active', 'paid', 'paused'));
UPDATE debts SET status = 'paid' WHERE status = 'settled';

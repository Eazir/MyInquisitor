ALTER TABLE recurring_expenses DROP CONSTRAINT IF EXISTS recurring_expenses_frequency_check;

ALTER TABLE recurring_expenses ADD CONSTRAINT recurring_expenses_frequency_check
  CHECK (frequency IN ('monthly', 'yearly', 'weekly', 'biweekly', 'once'));

ALTER TABLE recurring_expenses ADD COLUMN billing_month INT CHECK (billing_month BETWEEN 1 AND 12);

ALTER TABLE recurring_expenses DROP COLUMN IF EXISTS billing_month;

ALTER TABLE recurring_expenses DROP CONSTRAINT IF EXISTS recurring_expenses_frequency_check;

ALTER TABLE recurring_expenses ADD CONSTRAINT recurring_expenses_frequency_check
  CHECK (frequency IN ('monthly', 'yearly', 'weekly', 'biweekly'));

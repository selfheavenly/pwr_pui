export interface User {
  user_id: number;
  google_id?: string;
  email?: string;
  name?: string;
  balance: number;
  ytd_change_percent: number;
}

type Session = {
  session_id: string;
  user_id: string;
  created_at: Date;
  last_seen: Date;
  user_agent: string;
  ip_address: string;
  is_short_lived: boolean;
  is_api: boolean;
  location: string;
}